package go_workflow

import "testing"

func TestGetVariable(t *testing.T) {
	TestInitWorkflow(t)
	v := getVariableBy(1, "nginx_up", "nginxUp", "ok")
	t.Logf("%+v", v)
	if v.WorkflowId != 1 {
		t.Errorf("期望能查询到数据，实际没有查询结果")
	}
}
