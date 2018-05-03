package go_workflow

// Status
const (
	UndoState      = iota // 未执行状态
	ActivatedState        // 激活状态
	RunningState          // 运行中状态
	FailureState          // 运行失败
	SuccessState          // 运行成功
	CompletedState        // 完成状态
)

const (
	StartNode      = iota // 开始节点
	EndNode               // 终止节点
	AutoTrigger           // 自动执行
	ManualTrigger         // 人工执行
	MessageTrigger        // 消息
	TimingTrigger         // 定时执行
)

const (
	StartNodeName = "start"
	EndNodeName   = "end"
)
