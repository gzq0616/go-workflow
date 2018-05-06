package go_workflow

import (
	"testing"
)

func TestMd5Sum(t *testing.T) {
	cases := []struct {
		str      string
		expected string
	}{
		{
			str:      "1aaa",
			expected: "068896804551b982511aff19931628b9",
		}, {
			str:      "123456",
			expected: "e10adc3949ba59abbe56e057f20f883e",
		},
	}
	for _, c := range cases {
		value := md5Sum(c.str)
		if value != c.expected {
			t.Errorf("md5Sum error,actual:%s,expected:%s", value, c.expected)
		}
	}
}

func TestListRemove(t *testing.T) {
	cases := []struct {
		list    []string
		element string
	}{
		{list: []string{"1", "2", "3"}, element: "3"},
		{list: []string{"1", "2", "3"}, element: "2"},
		{list: []string{"1", "2", "3"}, element: "1"},
	}

	for _, c := range cases {
		t.Logf("%+v--%s-->%+v", c.list, c.element, listRemove(c.list, c.element))
	}
}

func TestUtilsVerify(t *testing.T) {
	cases := []struct {
		name      string
		condition string
		expected  bool
	}{
		{
			name:      "测试复杂表达式",
			condition: "{NODE_A.ACTION.key1} > 88 || ( {NODE_A.ACTION.key1} == true ) && {NODE_C.ACTION.key1} == ok && {NODE_C.ACTION.key1 == null}",
			expected:  true,
		},
		{
			name:      "测试比较符写错情况",
			condition: "{NODE_A.ACTION.key1} && 88",
			expected:  false,
		},
		{
			name:      "测试变量格式不正确情况",
			condition: "{ NODE_A.ACTION.key1 } > 88",
			expected:  false,
		},
		{
			name:      "测试多一个左括号情况1",
			condition: "({NODE_A.ACTION.key1} > 88",
			expected:  false,
		},
		{
			name:      "测试多一个右括号情况1",
			condition: "{NODE_A.ACTION.key1}) > 88",
			expected:  false,
		},
		{
			name:      "测试多一个左括号情况2",
			condition: "( {NODE_A.ACTION.key1} > 88",
			expected:  false,
		},
		{
			name:      "测试多一个右括号情况2",
			condition: "( {NODE_A.ACTION.key1} > 88 ) )",
			expected:  false,
		},
		{
			name:      "测试表达式正确情况",
			condition: "{NODE_A.ACTION.key1} > 88",
			expected:  true,
		},
		{
			name:      "测试表达式不完整情况",
			condition: "{NODE_A.ACTION.key1} > 88 &&",
			expected:  false,
		},
		{
			name:      "测试表达式为空情况1",
			condition: "",
			expected:  true,
		},
		{
			name:      "测试表达式为空情况2",
			condition: "                 ",
			expected:  true,
		},
	}

	for _, c := range cases {
		err := verify(c.condition)

		if err != nil {
			if !c.expected {
				t.Logf("%s:PASS\n%s-->%t:%s\n", c.name, c.condition, c.expected, err)
			} else {
				t.Errorf("%s:FAIL\n%s-->%t:%s\n", c.name, c.condition, c.expected, err)
			}
		} else {
			t.Logf("%s:PASS\n%s-->%t\n", c.name, c.condition, c.expected)
		}
	}
}

func TestRecurCalc(t *testing.T) {
	stack := NewStack()
	stack.Push("(")
	stack.Push(true)
	stack.Push("&&")
	recurCalc(stack, false)
	peak := stack.Peak()
	if peak != false {
		t.Errorf("expect false != actual:%t", peak)
	}
	stack.Push("&&")
	stack.Push(false)
	stack.Push("&&")
	recurCalc(stack, true)
	peak = stack.Peak()
	if peak != false {
		t.Errorf("expect true != actual:%t", peak)
	}
	if stack.Len() > 2 {
		t.Errorf("expect:2 != actual:%d", stack.Len())
	}
}
