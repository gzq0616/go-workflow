package go_workflow

import (
	"github.com/go-xorm/xorm"
)

var xe *xorm.Engine

func InitWorkflow(e *xorm.Engine) error {
	xe = e
	xe.ShowSQL(true)
	return xe.Sync2(new(TplAction), new(TplVariable), new(TplTransition), new(TplWorkflow), new(TplNode),
		new(TaskVariable), new(TaskTransition), new(TaskWorkflow), new(TaskNode), )
}

func Start() {
	// todo
}

func Stop() {
	// todo
}
