package go_workflow

import "time"

// 节点与节点的关系
type TaskTransition struct {
	Id           int       `json:"id" xorm:"pk autoincr"`
	WorkflowId   int       `json:"workflow_id" xorm:"unique(WorkflowId,source_node_id,target_node_id) notnull"`
	SourceNodeId int       `json:"source_node_id" xorm:"unique(WorkflowId,source_node_id,target_node_id) notnull"` // 源节点
	TargetNodeId int       `json:"target_node_id" xorm:"unique(WorkflowId,source_node_id,target_node_id) notnull"` // 目标节点
	Condition    string    `json:"condition"`                                                                      // (执行条件 ( { value1 } > 88 and { value2 } != true )
	CreatedAt    time.Time `json:"created_at" xorm:"created"`
	UpdatedAt    time.Time `json:"updated_at" xorm:"updated"`
}

func (self *TaskTransition) Verify() error {
	return conditionVerify(self.Condition)
}
