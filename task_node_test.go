package go_workflow

import "testing"

func TestPreConditionCheck(t *testing.T) {
	TestInitEngine(t)
	node := &TaskNode{Id:3}
	xe.Get(node)
	ok := node.preConditionCheck()
	if !ok {
		t.Errorf("TestPreConditionCheck fail")
	}
}