package go_workflow

import "testing"

func TestAction_Add(t *testing.T) {
	TestInitWorkflow(t)
	act := []string{"nginxDown", "nginxUp"}
	var actions []*TplAction
	for _, a := range act {
		ac, err := NewAction(a).Add()
		if err != nil {
			t.Fatalf("add tpl action fail %s", err)
		}
		actions = append(actions, ac)
	}

	for _, a := range actions {
		_, err := NewTplVariable("ok", "bool", "是否完成", a.Id).Add()
		if err != nil {
			t.Fatalf("add tpl variable fail %s", err)
		}
	}
}

func TestTplWorkflow_Add(t *testing.T) {
	TestInitWorkflow(t)
	_, err := NewTplWorkflow("nginx_workflow", "").Add()
	if err != nil {
		t.Fatalf("add tpl workflow fail %s", err)
	}
}

func TestTplNode_Add(t *testing.T) {
	TestInitWorkflow(t)
	ns := []map[string]interface{}{{"name": "nginx_up", "action": 2}, {"name": "nginx_down", "action": 1}}
	for _, n := range ns {
		_, err := NewTplNode(n["name"].(string), "", "", 1, n["action"].(int), AutoTrigger, 0, 0).Add()
		if err != nil {
			t.Fatalf("add tpl node fail %s", err)
		}
	}
}

func TestTplTransition_Add(t *testing.T) {
	TestInitWorkflow(t)
	_, err := NewTplTransition(2, 1, 1, "{ok} == true", ).Add()
	if err != nil {
		t.Fatalf("add tpl node fail %s", err)
	}
}
