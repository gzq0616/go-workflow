package go_workflow

import "testing"

func TestNewWorkflow(t *testing.T) {
	TestInitEngine(t)
	err := NewFlowDiagram(2, 1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStartTaskWorkflow(t *testing.T) {
	TestInitEngine(t)
	err := StartTaskWorkflow(1, "nginx_workflow")
	if err != nil {
		t.Fatal(err)
	}
}
