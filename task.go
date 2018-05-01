package go_workflow

type TaskVariable struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`      // 定义变量名称
	ActionID int    `json:"action_id"` // 变量所属action
	Type     string `json:"type"`      // 变量类型 int,string,time,bool
	Describe string `json:"describe"`  // 变量描述
	TaskID   int    `json:"task_id"`   //
}

type TaskAction struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	TaskID int    `json:"task_id"` //
}

type TaskTransition struct {
	ID         int    `json:"id"`
	SourceNode int    `json:"source_node"`
	TargetNode int    `json:"target_node"`
	Condition  string `json:"condition"`
	TaskID     int    `json:"task_id"` //
}

type TaskWorkflow struct {
	ID     int
	Name   string `json:"name"`    // 流程名称
	TaskID int    `json:"task_id"` //
}

type TaskNode struct {
	ID         int
	Name       string  `json:"name"`         // 节点名称
	WorkflowID int     `json:"workflow_id"`  // Node属于哪个工作流
	ActionID   int     `json:"action_id"`    // 调用的action
	StatusID   int     `json:"status_id"`    // 状态
	NodeTypeID int     `json:"node_type_id"` // Node类型
	X          float64 `json:"x"`            // 坐标x
	Y          float64 `json:"y"`            // 坐标y
	TaskID     int     `json:"task_id"`      //
}

type TaskManager interface {
	ListenToken()            // 监听token
	RunAction()              // 执行action
	ActiveNode()             // 根据条件激活下一个Node
	SetToken()               // 完成后设置token
	ChangeStatus(status int) // 改变状态
}

func NewWorkflow(taskId int, workflow Workflow) {

}

func (self *TaskWorkflow) Start() {

}
