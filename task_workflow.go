package go_workflow

import (
	"time"
	"fmt"
	"github.com/go-xorm/xorm"
	"errors"
)

type TaskWorkflow struct {
	Id              int                    `json:"id" xorm:"pk autoincr"`
	TaskId          int                    `json:"task_id" xorm:"unique(task_id,name) notnull"`           // task
	Name            string                 `json:"name" xorm:"unique(task_id,name) varchar(100) notnull"` // 流程名称,最好用字母
	Alias           string                 `json:"alias"`                                                 // 流程别名或者中文名称
	Status          int                    `json:"status"`                                                // 流程状态
	sign            chan string            `json:"-" xorm:"-"`                                            // 信号
	Result          map[string]bool        `json:"-" xorm:"json"`                                         // 存储节点完成状态,key为节点的名称
	nodeChannel     map[string]chan string `json:"-" xorm:"-"`                                            // 节点注册的channel,key为节点的名称
	nextNodeChannel chan *TaskNode         `json:"-" xorm:"-"`
	CreatedAt       time.Time              `json:"created_at" xorm:"created"`
	UpdatedAt       time.Time              `json:"updated_at" xorm:"updated"`
}

func getWorkflowByTaskIdAndWorkflowName(taskId int, workflowName string) (*TaskWorkflow, error) {
	// 根据taskId和workflowName查询TaskWorkflow
	workflow := &TaskWorkflow{TaskId: taskId, Name: workflowName}
	_, err := xe.Get(workflow)
	if err != nil {
		return nil, err
	}
	return workflow, nil
}

// workflow发布订阅消息
func (self *TaskWorkflow) register(nodeName string) {
	self.nodeChannel[nodeName] = make(chan string, 1)
}

func (self *TaskWorkflow) notify() {
	for self.Status == RunningState {
		signal := <-self.sign
		for _, channel := range self.nodeChannel {
			channel <- signal
		}
	}
}

func (self *TaskWorkflow) cancel(nodeName string) {
	channel := self.nodeChannel[nodeName]
	delete(self.nodeChannel, nodeName)
	close(channel)
}

// 根据流程模板新建工作流
func NewFlowDiagram(taskId int, tplWorkflowId int) error {
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
				Name:       tplVariable.Name,
				WorkflowId: taskWorkflow.Id,
				NodeId:     taskNode.Id,
				ActionId:   tplVariable.ActionId,
				Type:       tplVariable.Type,
				Describe:   tplVariable.Describe,
			}
			_, err := session.Insert(taskVariable)
			if err != nil {
				return err
			}
		}
	}

	session.Commit()
	return nil
}

func StartTaskWorkflow(taskId int, workflowName string) error {
	// 根据taskId和workflowName查询TaskWorkflow
	workflow, err := getWorkflowByTaskIdAndWorkflowName(taskId, workflowName)
	if err != nil {
		return err
	}
	return startTaskWorkflow(workflow)
}

func CloseTaskWorkflow(taskId int, workflowName string) error {
	// 根据taskId和workflowName查询TaskWorkflow
	workflow, err := getWorkflowByTaskIdAndWorkflowName(taskId, workflowName)
	if err != nil {
		return err
	}
	return closeTaskWorkflow(workflow)
}

func startTaskWorkflow(workflow *TaskWorkflow) error {
	session := xe.NewSession()
	err := session.Begin()
	if err != nil {
		return err
	}

	// 重置整个流程状态为初始值
	err = resetWorkflow(workflow, session)
	if err != nil {
		session.Rollback()
		return err
	}

	// 设置流程状态为运行中
	workflow.Status = RunningState
	_, err = session.ID(workflow.Id).Cols("status").Update(workflow)
	if err != nil {
		return err
	}

	// 寻找流程入口node节点
	startTaskNode := &TaskNode{WorkflowId: workflow.Id, Name: StartNodeName, NodeType: StartNode}
	_, err = session.Get(startTaskNode)
	if err != nil {
		session.Rollback()
		return err
	}
	session.Commit()

	workflow.sign = make(chan string)
	//workflow.Result = make(map[string]bool) 可能不需要初始化,数据库查询出来应该已经初始化过了
	workflow.nodeChannel = make(map[string]chan string)
	workflow.nextNodeChannel = make(chan *TaskNode, 10)
	// 将入口Node加入队列处理
	workflow.nextNodeChannel <- startTaskNode

	// 加入引擎中启动
	go func() {
		workflow.notify()
	}()

	return nil
}

func closeTaskWorkflow(workflow *TaskWorkflow) error {
	// todo:  关闭工作流
	return nil
}

func resetWorkflow(workflow *TaskWorkflow, session *xorm.Session) error {
	// 重置TaskNode
	_, err := session.Table(new(TaskNode)).Cols("status").Where("workflow_id=? AND status!=?", workflow.Id, UndoState).Update(map[string]int{"status": UndoState})
	if err != nil {
		return err
	}

	// 重置TaskVariable
	_, err = session.Table(new(TaskVariable)).Cols("value").Where("workflow_id=? AND value!=?", workflow.Id, "").Update(map[string]string{"value": ""})
	if err != nil {
		return err
	}
	return nil
}
