package go_workflow

import (
	"fmt"
	"errors"
	"strings"
	"crypto/md5"
	"encoding/hex"
)

func md5Sum(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	cipherStr := hash.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func listRemove(list []string, item string) []string {
	var ret []string
	for index, value := range list {
		if value == item {
			end := index + 1
			ret = append(ret, list[:index]...)
			ret = append(ret, list[end:]...)
		}
	}
	return ret
}

var compareOperator = []string{"==", "!=", ">", ">=", "<", "<="}
var logicOperator = []string{"&&", "||"}

func isOperator(ch string, list []string) bool {
	for _, op := range list {
		if op == ch {
			return true
		}
	}
	return false
}

// {NODE_A.ACTION.key1} > 88 && ( {NODE_A.ACTION.key1} == true && {NODE_C.ACTION.key1} == ok ) && {NODE_C.ACTION.key1} == null
// 三个为一组,变量使用node+action+name方式表示，以{}包裹，中间不能有空格，而其它元素必须以空格分隔，字符串空值用null表示
func conditionVerify(condition string) error {
	exprStr := strings.Split(condition, " ")

	firstStack := NewStack()
	secondStack := NewStack()

	for _, s := range exprStr {
		ch := strings.Trim(string(s), " ")
		switch {
		case ch == "":
		case ch == "(":
			secondStack.Push(ch)
		case strings.HasPrefix(ch, "{") && strings.HasSuffix(ch, "}"):
			if len(strings.Split(ch, ".")) == 3 {
				firstStack.Push(ch)
			} else {
				return errors.New("表达式错误，变量名格式不符,应该包含node,action,variable")
			}
		case ch == ")":
			last1 := secondStack.Pop().(string)
			if secondStack.Empty() {
				return errors.New("括号不匹配")
			}
			last2 := secondStack.Pop().(string)
			if last2 != "(" {
				return errors.New("括号不匹配")
			}
			comb := fmt.Sprintf("%s %s %s", last2, last1, ch)
			recurStack(secondStack, comb)
		case isOperator(ch, logicOperator):
			secondStack.Push(ch)
		case strings.HasPrefix(ch, "(") && len(ch) > 1:
			return errors.New("(格式不符，'('必须以空格分开")
		case strings.HasSuffix(ch, ")") && len(ch) > 1:
			return errors.New("(格式不符，')'必须以空格分开")
		default:
			if firstStack.Len() == 1 && !isOperator(ch, compareOperator) {
				return errors.New("表达式错误，比较符不正确")
			}
			firstStack.Push(ch)
		}

		if firstStack.Len() == 3 {
			firstLast1 := firstStack.Pop().(string)
			firstLast2 := firstStack.Pop().(string)
			firstLast3 := firstStack.Pop().(string)
			comb := fmt.Sprintf("%s %s %s", firstLast3, firstLast2, firstLast1)
			if !strings.Contains(comb, "{") && !strings.Contains(comb, "}") {
				return errors.New("表达式错误，变量名格式不符")
			}
			recurStack(secondStack, comb)
		}
	}

	if secondStack.Len() > 1 {
		return errors.New("括号不匹配或者表达式不完整")
	}

	return nil
}

func recurStack(stack *Stack, element string) {
	if stack.Empty() {
		stack.Push(element)
		return
	}
	yes := isOperator(stack.Peak().(string), logicOperator)
	stack.Push(element)
	if yes {
		first := stack.Pop().(string)
		second := stack.Pop().(string)
		third := stack.Pop().(string)
		comb := fmt.Sprintf("%s %s %s", third, second, first)
		recurStack(stack, comb)
	}
}

func recurCalc(stack *Stack, element bool) {
	if stack.Empty() {
		stack.Push(element)
		return
	}
	op, ok := stack.Peak().(string)
	if ok && op != "(" {
		var ret bool
		stack.Pop()
		obj := stack.Pop().(bool)
		switch op {
		case "&&":
			ret = element && obj
		case "||":
			ret = element || obj
		}
		recurCalc(stack, ret)
	} else {
		stack.Push(element)
	}
}

func compareInt(actual, expect int64, op string, reverse bool) (bool, error) {
	switch op {
	case "==":
		return actual == expect, nil
	case "!=":
		return actual != expect, nil
	case ">":
		ret := actual > expect
		if reverse {
			return !ret, nil
		}
		return ret, nil
	case ">=":
		ret := actual >= expect
		if reverse {
			return !ret, nil
		}
		return ret, nil
	case "<":
		ret := actual < expect
		if reverse {
			return !ret, nil
		}
		return ret, nil
	case "<=":
		ret := actual <= expect
		if reverse {
			return !ret, nil
		}
		return ret, nil
	default:
		return false, errors.New("无法比较结果")
	}
}

func compareString(actual, expect string, op string) (bool, error) {
	switch op {
	case "==":
		return actual == expect, nil
	case "!=":
		return actual != expect, nil
	default:
		return false, errors.New(fmt.Sprintf("字符串不能比较大小:%s", op))
	}
}

func compareBool(actual, expect bool, op string) (bool, error) {
	switch op {
	case "==":
		return actual == expect, nil
	case "!=":
		return actual != expect, nil
	default:
		return false, errors.New(fmt.Sprintf("布尔类型不能比较大小:%s", op))
	}
}

func hasVariableExpr(str string) bool {
	return strings.HasPrefix(str, "{")
}

func hasTitle(str string) bool {
	if str == "" {
		return false
	}
	a := []byte(str)[0]
	if a >= 97 && a <= 113 {
		return true
	}
	return false
}

func PreConditionVerify(workflowId int, condition []string) (bool, error) {
	/*
	condition: [node1,node2,node3]
	依赖节点
	 */
	for _, nodeStr := range condition {
		node := &TaskNode{WorkflowId: workflowId, Name: nodeStr}
		has, err := xe.Exist(node)
		if err != nil {
			return false, err
		}
		if !has {
			return false, errors.New(fmt.Sprintf("not found node %s", nodeStr))
		}
	}

	return true, nil
}
