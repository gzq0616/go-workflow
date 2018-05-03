package go_workflow

import (
	"time"
	"errors"
	"fmt"
)

type TplAction struct {
	Id        int       `json:"id" xorm:"pk autoincr"`
	Name      string    `json:"name" xorm:"varchar(100) notnull"`
	CreatedAt time.Time `json:"created_at" xorm:"created"`
	UpdatedAt time.Time `json:"updated_at" xorm:"updated"`
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

// 工作流
type TplWorkflow struct {
	Id        int       `json:"id" xorm:"pk autoincr"`
	Name      string    `json:"name" xorm:"varchar(100) unique"` // 流程名称,最好用字母
	Alias     string    `json:"alias"`                           // 流程别名或者中文名称
	CreatedAt time.Time `json:"created_at" xorm:"created"`
	UpdatedAt time.Time `json:"updated_at" xorm:"updated"`
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

// 节点与节点的关系
type TplTransition struct {
	Id           int       `json:"id" xorm:"pk autoincr"`
	WorkflowId   int       `json:"workflow_id" xorm:"unique(WorkflowId,source_node_id,target_node_id) notnull"`
	SourceNodeId int       `json:"source_node_id" xorm:"unique(WorkflowId,source_node_id,target_node_id) notnull"` // 源节点
	TargetNodeId int       `json:"target_node_id" xorm:"unique(WorkflowId,source_node_id,target_node_id) notnull"` // 目标节点
	Condition    string    `json:"condition" xorm:"default('1')"`                                                  // (执行条件 ( { value1 } > 88 and { value2 } != true )
	CreatedAt    time.Time `json:"created_at" xorm:"created"`
	UpdatedAt    time.Time `json:"updated_at" xorm:"updated"`
}

func AddAction(name string) (*TplAction, error) {
	action := &TplAction{Name: name}
	_, err := xe.Insert(action)
	if err != nil {
		return nil, err
	}
	return action, nil
}

func AddTplVariable(varName, varType, varDesc string, actionId int) (*TplVariable, error) {
	has, err := xe.Exist(&TplAction{Id: actionId})
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New(fmt.Sprintf("action not exsit,actionId:%d", actionId))
	}

	tplVariable := &TplVariable{
		Name:     varName,
		ActionId: actionId,
		Type:     varType,
		Describe: varDesc,
	}

	_, err = xe.Insert(tplVariable)
	if err != nil {
		return nil, err
	}
	return tplVariable, nil
}

func AddTplWorkflow(name, alias string) (*TplWorkflow, *TplNode, *TplNode, error) {
	tplWorkflow := &TplWorkflow{
		Name:  name,
		Alias: alias,
	}
	_, err := xe.Insert(tplWorkflow)
	if err != nil {
		return nil, nil, nil, err
	}

	startNode := &TplNode{
		Name:       StartNodeName,
		Alias:      "开始",
		WorkflowId: tplWorkflow.Id,
		NodeType:   StartNode,
	}
	_, err = xe.Insert(startNode)
	if err != nil {
		return nil, nil, nil, err
	}

	endNode := &TplNode{
		Name:       EndNodeName,
		Alias:      "结束",
		WorkflowId: tplWorkflow.Id,
		NodeType:   EndNode,
	}
	_, err = xe.Insert(endNode)
	if err != nil {
		return nil, nil, nil, err
	}

	return tplWorkflow, startNode, endNode, nil
}

func AddTplNode(name, alias, preCondition string, workflowId, actionId, nodeType int, x, y float64) (*TplNode, error) {
	has, err := xe.Exist(&TplAction{Id: actionId})
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New(fmt.Sprintf("action not exsit,actionId:%d", actionId))
	}

	has, err = xe.Exist(&TplWorkflow{Id: workflowId})
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New(fmt.Sprintf("workflow not exsit,workflowId:%d", workflowId))
	}

	tplNode := &TplNode{
		Name:         name,
		Alias:        alias,
		WorkflowId:   workflowId,
		ActionId:     actionId,
		NodeType:     nodeType,
		PreCondition: preCondition,
		X:            x,
		Y:            y,
	}

	_, err = xe.Insert(tplNode)
	if err != nil {
		return nil, err
	}
	return tplNode, nil
}

func AddTplTransition(sourceId, targetId, workflowId int, condition string) (*TplTransition, error) {
	// 校验是否存在
	sourceNode := &TplNode{Id: sourceId}
	has, err := xe.Get(sourceNode)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New(fmt.Sprintf("not found this template node by source id :%d", sourceId))
	}

	has, err = xe.Exist(&TplNode{Id: targetId})
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New(fmt.Sprintf("not found this template node by target id :%d", targetId))
	}
	tplTran := &TplTransition{SourceNodeId: sourceId, TargetNodeId: targetId, WorkflowId: workflowId}

	// 开始节点的条件强制设为1,即直接跳转到下一个节点
	if sourceNode.NodeType == StartNode {
		tplTran.Condition = "1"
	} else {
		tplTran.Condition = condition
	}

	_, err = xe.Insert(tplTran)
	if err != nil {
		return nil, err
	}

	return tplTran, nil
}
