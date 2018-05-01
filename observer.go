package go_workflow

type Observer interface {
	//更新事件
	Update()
}

// 被观察的对象接口
type Subject interface {
	//注册观察者
	Register(Observer)
	//注销观察者
	Deregist(Observer)
	//通知观察者事件
	Notify()
}
