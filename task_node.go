package go_workflow

import (
	"errors"
	"time"
	"strings"
)

type TaskNode struct {
	Id           int       `json:"id" xorm:"pk autoincr"`
	Name         string    `json:"name" xorm:"unique(name,workflow_id)"`        // 节点名称,同一个流程中不能出现相同的名称
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

func (self *TaskNode) preconditionCheck() bool {
	// todo 判断前置条件是否满足
	return false
}

func (self *TaskNode) run(workflow *TaskWorkflow, q chan<- *TaskNode) error {
	// 设置状态为激活状态,如果当前状态为正在运行则跳过!
	if self.Status == RunningState {
		return nil
	}

	// 设置流程中节点的状态为false
	workflow.Result[self.Name] = false
	self.Status = ActivatedState
	_, err := xe.ID(self.Id).Cols("status").Update(self)
	if err != nil {
		return err
	}

	// 接收信号，根据前置条件触发执行action
	// 没有任何前置条件，直接执行action
	if strings.Trim(self.PreCondition, "") != "" {
		// 首先注册监听
		workflow.register(self.Name)
		for {
			<-workflow.nodeChannel[self.Name]
			// 判断前置条件是否满足
			if self.preconditionCheck() {
				// 注销监听
				workflow.cancel(self.Name)
				break
			}
		}
	}

	// 执行action,设置状态为执行中,没有action直接跳转到下一个节点
	if self.ActionId != 0 {
		// todo : 调用action
	}

	// 根据条件激活下一个node,即将下一个node加入队列
	err = activateNextNode(self, q)
	if err != nil {
		return err
	}

	// 打上自己的token,并发出信号
	workflow.Result[self.Name] = true
	workflow.sign <- self.Name

	// 设置状态为完成
	self.Status = CompletedState
	_, err = xe.ID(self.Id).Cols("status").Update(self)
	if err != nil {
		return err
	}

	// 如果为最后一个节点,关闭流程
	if self.NodeType == EndNode {
		closeTaskWorkflow(workflow)
	}

	return nil
}

// 条件表达式解析
// ( { NODE_A.ACTION.key1 } > 88 and { NODE_A.ACTION.key1 } != true ) && { NODE_C.ACTION.key1 } == "ok" && { NODE_C.ACTION.key1 != "" }
// ( {1.3.k} > 55 || {4.5.k} == ok ) &&  2.3.k != 0
func (self *TaskNode) decide(condition string) bool {
	// todo:https://studygolang.com/articles/3970
	// 如果没有设置条件，即无条件执行下一步
	if condition == "" {
		return true
	}
	return false
}

func (self *TaskNode) verify(condition string) bool {
	//for _,s := range condition {
	//	ch := string(s)
	//
	//}
	// todo:https://studygolang.com/articles/3970
	return false
}

// 根据条件激活下一个节点
func activateNextNode(current *TaskNode, q chan<- *TaskNode) error {
	// 节点和节点的关系对象
	relations := make([]*TaskTransition, 0)
	err := xe.Where("workflow_id = ? AND source_node_id = ?", current.WorkflowId, current.Id).Find(&relations)
	if err != nil {
		return err
	}

	for _, relation := range relations {
		condition := relation.Condition
		// 判断条件是否满足,如果满足，将下一个节点加入队列
		if current.decide(condition) {
			nextNode := &TaskNode{Id: relation.TargetNodeId}
			has, err := xe.Get(nextNode)
			if err != nil {
				return err
			}
			if !has {
				return errors.New("not find next node")
			}
			q <- nextNode
		}
	}

	return nil
}
