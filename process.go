package go_workflow

type Variable struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`      // 定义变量名称
	ActionID int    `json:"action_id"` // 变量所属action
	Type     string `json:"type"`      // 变量类型 int,string,time,bool
	Describe string `json:"describe"`  // 变量描述
}

type Action struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// 节点与节点的关系
type Transition struct {
	ID         int    `json:"id"`
	SourceNode int    `json:"source_node"`
	TargetNode int    `json:"target_node"`
	Condition  string `json:"condition"`
}

// 工作流
type Workflow struct {
	ID   int
	Name string `json:"name"` // 流程名称
}

// 节点 A,B,C,D
type Node struct {
	ID         int
	Name       string  `json:"name"`         // 节点名称
	WorkflowID int     `json:"workflow_id"`  // Node属于哪个工作流
	ActionID   int     `json:"action_id"`    // 调用的action
	StatusID   int     `json:"status_id"`    // 状态
	NodeTypeID int     `json:"node_type_id"` // Node类型
	X          float64 `json:"x"`            // 坐标x
	Y          float64 `json:"y"`            // 坐标y
}