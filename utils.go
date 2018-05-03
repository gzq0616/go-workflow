package go_workflow

import (
	"strings"
	"fmt"
)

// 条件表达式解析
// ( { value1 } > 88 and { value2 } != true ) && { value3 } == "ok"

func parse(condition string) (bool, error) {
	condis := strings.Split(condition, " ")
	fmt.Println(condis)
	fmt.Println(condis[2])
	return false, nil
}

func (self *TaskTransition) Parse() (bool, error) {
	cds := strings.Split(self.Condition, " ")
	fmt.Println(cds)
	return false, nil
}
