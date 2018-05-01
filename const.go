package go_workflow

// Status
const (
	UndoState      = iota // 未执行状态
	ActivatedState        // 激活状态
	RunningState          // 运行中状态
	CompletedState        // 完成状态
)

const (
	AutoTrigger    = iota // 自动执行
	ManualTrigger         // 人工执行
	MessageTrigger        // 消息
	TimingTrigger         // 定时执行
)

//const (
//	StartNode = "开始节点"
//	TaskNode  = "人工节点"
//	AutoNode  = "自动节点"
//	DeciNode  = "决策节点"
//	ForkNode  = "发散节点"
//	JoinNode  = "聚合节点"
//	SubNode   = "子流程节点"
//	SignNode  = "会签节点"
//	WaitNode  = "等待节点"
//	EndNode   = "结束节点"
//)
