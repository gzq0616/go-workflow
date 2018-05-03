package go_workflow

import (
	"time"
)

type TplAction struct {
	Id        int       `json:"id" xorm:"pk autoincr"`
	Name      string    `json:"name" xorm:"varchar(100) notnull"`
	CreatedAt time.Time `json:"created_at" xorm:"created"`
	UpdatedAt time.Time `json:"updated_at" xorm:"updated"`
}

func (self *TplAction) Add() (*TplAction, error) {
	_, err := xe.Insert(self)
	return self, err
}

type TplVariable struct {
	Id        int       `json:"id" xorm:"pk autoincr"`
	Name      string    `json:"name" xorm:"varchar(100) unique(name,action_id) notnull"` // 定义变量名称
	ActionId  int       `json:"action_id" xorm:"unique(name,action_id) notnull"`         // 变量所属action
	Type      string    `json:"type" xorm:"varchar(100) notnull"`                        // 变量类型 int,string,time,bool
	Describe  string    `json:"describe" xorm:"text"`                                    // 变量描述
	CreatedAt time.Time `json:"created_at" xorm:"created"`
	UpdatedAt time.Time `json:"updated_at" xorm:"updated"`
}

func (self *TplVariable) Add() (*TplVariable, error) {
	_, err := xe.Insert(self)
	return self, err
}

// 工作流
type TplWorkflow struct {
	Id        int       `json:"id" xorm:"pk autoincr"`
	Name      string    `json:"name" xorm:"varchar(100) unique"` // 流程名称,最好用字母
	Alias     string    `json:"alias"`                           // 流程别名或者中文名称
	CreatedAt time.Time `json:"created_at" xorm:"created"`
	UpdatedAt time.Time `json:"updated_at" xorm:"updated"`
}

func (self *TplWorkflow) Add() (*TplWorkflow, error) {
	_, err := xe.Insert(self)
	return self, err
}

// 节点 A,B,C,D
type TplNode struct {
	Id           int       `json:"id" xorm:"pk autoincr"`
	Name         string    `json:"name" xorm:"unique(name,workflow_id)"`        // 节点名称
	Alias        string    `json:"alias"`                                       // Node中文名称
	WorkflowId   int       `json:"workflow_id" xorm:"unique(name,workflow_id)"` // Node属于哪个工作流
	ActionId     int       `json:"action_id" xorm:"notnull"`                    // 调用的action
	NodeType     int       `json:"node_type" xorm:"notnull"`                    // Node类型
	PreCondition string    `json:"pre_condition"`                               // 前置条件
	X            float64   `json:"x"`                                           // 坐标x
	Y            float64   `json:"y"`                                           // 坐标y
	CreatedAt    time.Time `json:"created_at" xorm:"created"`
	UpdatedAt    time.Time `json:"updated_at" xorm:"updated"`
}

func (self *TplNode) Add() (*TplNode, error) {
	_, err := xe.Insert(self)
	return self, err
}

// 节点与节点的关系
type TplTransition struct {
	Id           int       `json:"id" xorm:"pk autoincr"`
	WorkflowId   int       `json:"workflow_id" xorm:"notnull"`
	SourceNodeId int       `json:"source_node_id" xorm:"notnull"` // 源节点
	TargetNodeId int       `json:"target_node" xorm:"notnull"`    // 目标节点
	Condition    string    `json:"condition" xorm:"default('1')"` // (执行条件 ( { value1 } > 88 and { value2 } != true )
	CreatedAt    time.Time `json:"created_at" xorm:"created"`
	UpdatedAt    time.Time `json:"updated_at" xorm:"updated"`
}

func (self *TplTransition) Add() (*TplTransition, error) {
	_, err := xe.Insert(self)
	return self, err
}

func NewAction(name string) *TplAction {
	return &TplAction{Name: name}
}

func NewTplVariable(varName, varType, varDesc string, actionId int) *TplVariable {
	return &TplVariable{
		Name:     varName,
		ActionId: actionId,
		Type:     varType,
		Describe: varDesc,
	}
}

func NewTplWorkflow(name, alias string) *TplWorkflow {
	return &TplWorkflow{
		Name:  name,
		Alias: alias,
	}
}

func NewTplNode(name, alias, preCondition string, workflowId, actionId, nodeType int, x, y float64) *TplNode {
	return &TplNode{
		Name:         name,
		Alias:        alias,
		WorkflowId:   workflowId,
		ActionId:     actionId,
		NodeType:     nodeType,
		PreCondition: preCondition,
		X:            x,
		Y:            y,
	}
}

func NewTplTransition(sourceId, targetId, workflowId int, condition string) *TplTransition {
	return &TplTransition{SourceNodeId: sourceId, TargetNodeId: targetId, Condition: condition, WorkflowId: workflowId}
}
