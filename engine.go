package go_workflow

import (
	"fmt"
	"sync"
	"errors"
	"github.com/go-xorm/xorm"
)

var (
	xe        *xorm.Engine
	engine    *Engine
	actionMap map[string]Action
)

type Engine struct {
	status     int
	startQueue chan *TaskWorkflow
	stopQueue  chan *TaskWorkflow
}

type Action interface {
	Run() (map[string]string, error)
}

type RegAction struct {
	Name      string
	Action    interface{}
	Variables []RegVariable
}

type RegVariable struct {
	Name     string
	Type     string
	Describe string
}

func InitEngine(e *xorm.Engine, actions ...RegAction) error {
	var once sync.Once
	var err error
	once.Do(func() {
		xe = e
		engine = &Engine{
			status:     UndoState,
			startQueue: make(chan *TaskWorkflow, 1),
			stopQueue:  make(chan *TaskWorkflow, 1),
		}
		xe.ShowSQL(true)
		err = xe.Sync2(new(TplAction), new(TplVariable), new(TplTransition), new(TplWorkflow), new(TplNode),
			new(TaskVariable), new(TaskTransition), new(TaskWorkflow), new(TaskNode))
		if err != nil {
			return
		}
		// 自动注册action和variable
		err = registryAction(actions...)
	})
	return err
}

//自动注册action和variable
func registryAction(actions ...RegAction) error {
	session := xe.NewSession()
	session.Begin()
	actionMap = make(map[string]Action)
	for _, regAction := range actions {
		// action name 首字母必须大写
		actionName := regAction.Name
		if actionName == "" {
			session.Rollback()
			return errors.New("没有Action名称")
		}
		if a := []byte(actionName)[0]; a < 97 && a > 113 {
			session.Rollback()
			return errors.New("action name 必须以大写字母开头")
		}

		actionInterface, ok := regAction.Action.(Action)
		if !ok {
			session.Rollback()
			return errors.New(fmt.Sprintf("%s:没有实现Action接口", actionName))
		}
		actionMap[actionName] = actionInterface
		action := &TplAction{Name: actionName}
		has, err := session.Get(action)
		if err != nil {
			session.Rollback()
			return err
		}
		if !has {
			_, err = session.Insert(action)
			if err != nil {
				session.Rollback()
				return err
			}
		}

		for _, v := range regAction.Variables {
			variable := &TplVariable{Name: v.Name, ActionId: action.Id, Type: v.Type, Describe: v.Describe}
			has, err = session.Exist(variable)
			if err != nil {
				session.Rollback()
				return err
			}
			if !has {
				_, err = session.Insert(variable)
				if err != nil {
					session.Rollback()
					return err
				}
			}
		}
	}
	session.Commit()
	return nil
}

func Start() {
	if engine.status != RunningState {
		engine.status = RunningState
	}
}

func Stop() error {
	// todo:关闭流程引擎
	return nil
}
