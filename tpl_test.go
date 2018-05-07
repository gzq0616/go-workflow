package go_workflow

import "testing"

func TestAddAction(t *testing.T) {
	TestInitEngine(t)
	act := []string{"nginxDown", "nginxUp"}
	for _, a := range act {
		_, err := AddAction(a)
		if err != nil {
			t.Fatalf("add tpl action fail %s", err)
		}
	}
}

func TestAddTplVariable(t *testing.T) {
	TestInitEngine(t)
	actionIds := []int{1, 2}
	for _, actionId := range actionIds {
		_, err := AddTplVariable("ok", "bool", "是否完成", actionId)
		if err != nil {
			t.Fatalf("add tpl variable fail %s", err)
		}
	}
}

func TestAddTplWorkflow(t *testing.T) {
	TestInitEngine(t)
	_, startNode, endNode, err := AddTplWorkflow("nginx_workflow", "NGINX上下线流程")
	if err != nil {
		t.Fatalf("add tpl workflow fail %s", err)
	}
	if startNode.NodeType != StartNode {
		t.Fatal("start node type error")
	}
	if endNode.NodeType != EndNode {
		t.Fatal("end node type error")
	}
}

func TestAddTplNode(t *testing.T) {
	TestInitEngine(t)
	ns := []map[string]interface{}{{"name": "nginx_up", "action": 2}, {"name": "nginx_down", "action": 1}}
	for _, n := range ns {
		_, err := AddTplNode(n["name"].(string), n["name"].(string), "", 1, n["action"].(int), AutoTrigger, 0, 0)
		if err != nil {
			t.Fatalf("add tpl node fail %s", err)
		}
	}
}

func TestAddTplTransition(t *testing.T) {
	TestInitEngine(t)
	cases := []struct {
		source    int
		target    int
		workflow  int
		condition string
		expected  string
	}{
		{
			source:    1,
			target:    4,
			workflow:  1,
			condition: "dfdfdfdf",
			expected:  "",
		},
		{
			source:    4,
			target:    3,
			workflow:  1,
			condition: "dfdfdfdf",
			expected:  "dfdfdfdf",
		},
		{
			source:    3,
			target:    2,
			workflow:  1,
			condition: "dfdfdfdf",
			expected:  "dfdfdfdf",
		},
	}

	for _, c := range cases {
		tran, err := AddTplTransition(c.source, c.target, c.workflow, c.condition)
		if err != nil {
			t.Fatalf("add tpl node fail %s", err)
		}
		if tran.Condition != c.expected {
			t.Fatalf("条件值和期望值不一样,实际值:%s，期望值:%s", tran.Condition, c.expected)
		}
	}

}
