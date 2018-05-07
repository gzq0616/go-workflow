package go_workflow

import (
	"fmt"
	"log"
	"sync"
	"errors"
	"github.com/go-xorm/xorm"
)

var (
	xe        *xorm.Engine      // 全局数据库引擎对象
	engine    *Engine           // 全局engine对象
	actionMap map[string]Action // 全局已注册的action对象
)

type Engine struct {
	status      int                      // 引擎运行状态
	wg          sync.WaitGroup           //
	mutex       *sync.RWMutex            // 读写锁
	startQueue  chan *TaskWorkflow       // 启动的流程队列
	stopQueue   chan *TaskWorkflow       // 停止的流程队列
	workflowMap map[string]*TaskWorkflow // 运行中的workflow
}

func (self *Engine) start() error {
	if self.status == RunningState {
		return errors.New("请不要重复启动引擎")
	}

	// 查询正在执行的workflow
	var runWorkflow []*TaskWorkflow
	err := xe.Where("status=?", RunningState).Find(&runWorkflow)
	if err != nil {
		return err
	}

	for _, w := range runWorkflow {
		w.run()
		self.workflowMap[""] = w
	}

	go func() {
		go func() {
			for self.status == RunningState {
				select {
				case taskWorkflow := <-self.startQueue:
					self.wg.Add(1)
					go taskWorkflow.start()
				case taskWorkflow := <-self.stopQueue:
					go taskWorkflow.stop()
					self.wg.Done()
				}
			}
		}()
		self.wg.Wait()
	}()

	return nil
}

func (self *Engine) stop() error {
	if self.status != RunningState {
		return nil
	}
	for _, v := range engine.workflowMap {
		v.stop()
	}
	self.status = UndoState
	return nil
}

type Action interface {
	Run() (map[string]string, error) // action的运行结果保存在map中
}

type RegAction struct {
	Name      string        // action名称，第一个字母必须为大写
	Action    interface{}   // 实现Action方法的对象接口
	Variables []RegVariable // action执行后相关的变量及属性
}

type RegVariable struct {
	Name     string // 变量名称
	Type     string // 变量类型
	Describe string // 变量描述
}

// 初始化engine和表结构方法
func initEngine() error {
	engine = &Engine{
		status:      UndoState,
		startQueue:  make(chan *TaskWorkflow, 1),
		stopQueue:   make(chan *TaskWorkflow, 1),
		workflowMap: make(map[string]*TaskWorkflow),
	}
	xe.ShowSQL(true)
	return xe.Sync2(new(TplAction), new(TplVariable), new(TplTransition), new(TplWorkflow), new(TplNode),
		new(TaskVariable), new(TaskTransition), new(TaskWorkflow), new(TaskNode))
}

//自动注册action和variable方法
func registryAction(actions ...RegAction) error {
	if len(actions) == 0 {
		return nil
	}
	session := xe.NewSession()
	session.Begin()
	actionMap = make(map[string]Action)
	for _, regAction := range actions {
		actionName := regAction.Name
		// 校验名称是否为空
		if actionName == "" {
			session.Rollback()
			return errors.New("没有Action名称")
		}
		// 校验是否以大小字母开头
		if !hasTitle(actionName) {
			session.Rollback()
			return errors.New("action必须以大写字母开头")
		}
		// 校验是否实现Action接口
		actionInterface, ok := regAction.Action.(Action)
		if !ok {
			session.Rollback()
			return errors.New(fmt.Sprintf("%s:没有实现Action接口", actionName))
		}

		// 自动注册action
		actionMap[actionName] = actionInterface

		// 首先查询是否已存在action
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

		// 自动注册变量名和属性
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

// 初始化引擎
func EngineInit(e *xorm.Engine, actions ...RegAction) error {
	var once sync.Once
	var err error
	once.Do(func() {
		xe = e
		// 初始化engine和表结构
		err = initEngine()
		if err != nil {
			return
		}
		// 自动注册action和variable
		err = registryAction(actions...)
	})
	return err
}

// 启动引擎
func EngineStart() error {
	log.Println("正在启动引擎。。。")
	err := engine.start()
	if err != nil {
		log.Println("引擎启动失败")
		return err
	}
	log.Println("引擎启动成功")
	return nil
}

// 关闭引擎
func EngineStop() error {
	log.Println("正在关闭引擎。。。")
	err := engine.stop()
	if err != nil {
		log.Println("无法关闭引擎")
		return err
	}

	log.Println("引擎已关闭。。。")
	return nil
}
