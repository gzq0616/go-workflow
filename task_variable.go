package go_workflow

import (
	"time"
)

type TaskVariable struct {
	Id         int       `json:"id" xorm:"pk autoincr"`
	Name       string    `json:"name" xorm:"varchar(100) unique(name,node_id,action_id) notnull"` // 定义变量名称
	WorkflowId int       `json:"workflow_id" xorm:"notnull"`                                      // 所属流程
	NodeId     int       `json:"node_id" xorm:"unique(name,node_id,action_id) notnull"`           // 所属node
	ActionId   int       `json:"action_id" xorm:"unique(name,node_id,action_id) notnull"`         // 变量所属action
	Type       string    `json:"type" xorm:"varchar(100) notnull"`                                // 变量类型 int,string,time,bool
	Describe   string    `json:"describe" xorm:"text"`                                            // 变量描述
	Value      string    `json:"value"`                                                           // 变量实际值，存入数据库转为string类型
	CreatedAt  time.Time `json:"created_at" xorm:"created"`
	UpdatedAt  time.Time `json:"updated_at" xorm:"updated"`
}

func getVariableBy(workflowId int, nodeName, actionName, variableName string) *TaskVariable {
	/*
	 select * from task_variable tv
	 left join task_node as tn on tn.id = tv.node_id
	 left join tpl_action as ta on ta.id = tv.action_id
	 where tv.workflow_id = workflowId and tn.name = nodeName  and ta.name = actionName and tv.name=variableName
	 */
	taskVariable := new(TaskVariable)
	has, err := xe.Table("task_variable").Alias("tv").
		Join("LEFT OUTER", "task_node", "tn.id=tv.node_id").Alias("tn").
		Join("LEFT OUTER", "task_action", "ta.id=tv.action_id").Alias("ta").
		Where("tv.workflow_id = ? and tn.name = ?  and ta.name = ? and tv.name=?", workflowId, nodeName, actionName, variableName).
		Get(taskVariable)
	if err != nil {
		return nil
	}
	if has {
		return taskVariable
	}
	return nil
}
