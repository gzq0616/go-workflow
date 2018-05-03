package go_workflow

import (
	"errors"
	"fmt"
	"time"
)

type TaskVariable struct {
	Id        int       `json:"id" xorm:"pk autoincr"`
	Name      string    `json:"name" xorm:"varchar(100) unique(name,node_id,action_id) notnull"` // 定义变量名称
	NodeId    int       `json:"node_id" xorm:"unique(name,node_id,action_id) notnull"`           // 所属node
	ActionId  int       `json:"action_id" xorm:"unique(name,node_id,action_id) notnull"`         // 变量所属action
	Type      string    `json:"type" xorm:"varchar(100) notnull"`                                // 变量类型 int,string,time,bool
	Describe  string    `json:"describe" xorm:"text"`                                            // 变量描述
	Value     string    `json:"value"`                                                           // 变量实际值，存入数据库转为string类型
	CreatedAt time.Time `json:"created_at" xorm:"created"`
	UpdatedAt time.Time `json:"updated_at" xorm:"updated"`
}

// 节点与节点的关系
type TaskTransition struct {
	Id           int       `json:"id" xorm:"pk autoincr"`
	WorkflowId   int       `json:"workflow_id" xorm:"notnull"`
	SourceNodeId int       `json:"source_node_id" xorm:"notnull"` // 源节点
	TargetNodeId int       `json:"target_node" xorm:"notnull"`    // 目标节点
	Condition    string    `json:"condition" xorm:"default('1')"` // (执行条件 ( { value1 } > 88 and { value2 } != true )
	CreatedAt    time.Time `json:"created_at" xorm:"created"`
	UpdatedAt    time.Time `json:"updated_at" xorm:"updated"`
}

type TaskWorkflow struct {
	Id        int       `json:"id" xorm:"pk autoincr"`
	TaskId    int       `json:"task_id" xorm:"notnull"`          // task
	Name      string    `json:"name" xorm:"varchar(100) unique"` // 流程名称,最好用字母
	Alias     string    `json:"alias"`                           // 流程别名或者中文名称
	Status    int       `json:"status"`                          // 流程状态
	CreatedAt time.Time `json:"created_at" xorm:"created"`
	UpdatedAt time.Time `json:"updated_at" xorm:"updated"`
}

type TaskNode struct {
	Id           int       `json:"id" xorm:"pk autoincr"`
	Name         string    `json:"name" xorm:"unique(name,workflow_id)"`        // 节点名称
	Alias        string    `json:"alias"`                                       // Node中文名称
	WorkflowId   int       `json:"workflow_id" xorm:"unique(name,workflow_id)"` // Node属于哪个工作流
	ActionId     int       `json:"action_id" xorm:"notnull"`                    // 调用的action
	Status       int       `json:"status" xorm:"notnull"`                       // 状态
	NodeType     int       `json:"node_type" xorm:"notnull"`                    // Node类型
	PreCondition string    `json:"pre_condition"`                               // 前置条件
	X            float64   `json:"x"`                                           // 坐标x
	Y            float64   `json:"y"`                                           // 坐标y
	CreatedAt    time.Time `json:"created_at" xorm:"created"`
	UpdatedAt    time.Time `json:"updated_at" xorm:"updated"`
}

func NewWorkflow(taskId int, tplWorkflowId int) error {
	session := xe.NewSession()
	defer session.Close()
	// begin transaction
	err := session.Begin()
	if err != nil {
		return err
	}
	// 创建TaskWorkflow
	tplWorkflow := &TplWorkflow{Id: tplWorkflowId}
	has, err := session.Get(tplWorkflow)
	if err != nil {
		return err
	}
	if !has {
		return errors.New(fmt.Sprintf("not found tplWorkflow by Id:%d", tplWorkflowId))
	}
	taskWorkflow := &TaskWorkflow{TaskId: taskId, Name: tplWorkflow.Name, Alias: tplWorkflow.Alias, Status: UndoState}
	_, err = session.Insert(taskWorkflow)
	if err != nil {
		session.Rollback()
		return err
	}

	// 创建taskNode
	tplNodes := make([]*TplNode, 0)
	err = session.Where("workflow_id=?", tplWorkflowId).Find(&tplNodes)
	if err != nil {
		session.Rollback()
		return err
	}
	taskNodes := make([]*TaskNode, 0)
	for _, tplNode := range tplNodes {
		taskNode := &TaskNode{
			Name:       tplNode.Name,
			Alias:      tplNode.Alias,
			WorkflowId: taskWorkflow.Id,
			ActionId:   tplNode.ActionId,
			Status:     UndoState,
			NodeType:   tplNode.NodeType,
			X:          tplNode.X,
			Y:          tplNode.Y,
		}
		_, err = session.Insert(taskNode)
		if err != nil {
			session.Rollback()
			return err
		}
		taskNodes = append(taskNodes, taskNode)
	}

	// 创建路由
	tplTrans := make([]*TplTransition, 0)
	err = session.Where("workflow_id=?", tplWorkflowId).Find(&tplTrans)
	if err != nil {
		session.Rollback()
		return err
	}
	for _, tplTran := range tplTrans {
		tplSourceNode := &TplNode{Id: tplTran.SourceNodeId, WorkflowId: tplWorkflowId}
		has, err := session.Get(tplSourceNode)
		if err != nil {
			session.Rollback()
			return err
		}
		if !has {
			session.Rollback()
			return errors.New(fmt.Sprintf("not found tplSourceNode by Id:%d", tplTran.SourceNodeId))
		}

		taskSourceNode := &TaskNode{Name: tplSourceNode.Name, WorkflowId: taskWorkflow.Id}
		has, err = session.Get(taskSourceNode)
		if err != nil {
			session.Rollback()
			return err
		}
		if !has {
			session.Rollback()
			return errors.New(fmt.Sprintf("not found taskSourceNode by Name:%s and WorkflowId:%d", tplSourceNode.Name, taskWorkflow.Id))
		}

		tplTargetNode := &TplNode{Id: tplTran.TargetNodeId, WorkflowId: tplWorkflowId}
		has, err = session.Get(tplTargetNode)
		if err != nil {
			session.Rollback()
			return err
		}
		if !has {
			session.Rollback()
			return errors.New(fmt.Sprintf("not found tplTargetNode by Id:%d", tplTran.SourceNodeId))
		}

		taskTargetNode := &TaskNode{Name: tplTargetNode.Name, WorkflowId: taskWorkflow.Id}
		has, err = session.Get(taskTargetNode)
		if err != nil {
			session.Rollback()
			return err
		}
		if !has {
			session.Rollback()
			return errors.New(fmt.Sprintf("not found taskTargetNode by Name:%s and WorkflowId:%d", tplTargetNode.Name, taskWorkflow.Id))
		}

		taskTransition := &TaskTransition{WorkflowId: taskWorkflow.Id, SourceNodeId: taskSourceNode.Id, TargetNodeId: taskTargetNode.Id, Condition: tplTran.Condition}
		_, err = session.Insert(taskTransition)
		if err != nil {
			session.Rollback()
			return err
		}
	}

	// 创建变量
	for _, taskNode := range taskNodes {
		tplVariables := make([]*TplVariable, 0)
		err = session.Where("action_id=?", taskNode.ActionId).Find(&tplVariables)
		if err != nil {
			session.Rollback()
			return err
		}
		for _, tplVariable := range tplVariables {
			taskVariable := &TaskVariable{
				Name:     tplVariable.Name,
				NodeId:   taskNode.Id,
				ActionId: tplVariable.ActionId,
				Type:     tplVariable.Type,
				Describe: tplVariable.Describe,
			}
			_, err := session.Insert(taskVariable)
			if err != nil {
				return err
			}
		}
	}

	err = session.Commit()
	if err != nil {
		return err
	}
	return nil
}

func TaskWorkflowStart() {

}
