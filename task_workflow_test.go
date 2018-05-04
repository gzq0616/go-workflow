package go_workflow

import "testing"

func TestNewWorkflow(t *testing.T) {
	TestInitWorkflow(t)
	err := NewFlowDiagram(2, 1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStartTaskWorkflow(t *testing.T) {
	TestInitWorkflow(t)
	err := StartTaskWorkflow(1, "nginx_workflow")
	if err != nil {
		t.Fatal(err)
	}
}
