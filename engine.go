package go_workflow

import (
	"github.com/go-xorm/xorm"
	"sync"
)

var (
	xe     *xorm.Engine
	engine *Engine
)

type Engine struct {
	status    int
	workflows map[string]*TaskWorkflow
}

func InitWorkflow(e *xorm.Engine) error {
	var once sync.Once
	var err error
	once.Do(func() {
		xe = e
		engine = &Engine{
			status: UndoState,
		}
		xe.ShowSQL(true)
		err = xe.Sync2(new(TplAction), new(TplVariable), new(TplTransition), new(TplWorkflow), new(TplNode),
			new(TaskVariable), new(TaskTransition), new(TaskWorkflow), new(TaskNode), )
	})
	return err
}

func Start() error {
	// todo:启动流程引擎
	if engine.status != RunningState {
		engine.status = RunningState
	}

	return nil
}

func Stop() error {
	// todo:关闭流程引擎
	return nil
}
