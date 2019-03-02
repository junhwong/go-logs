package logs

// Logger 日志处理器
type Logger interface {
	New(calldepth int) Entry
	Append(v Handler)
	Log(v Entry)
	Filter(v Entry) Entry
	SetLevel(v Level)
}
