package logs

type Handler interface {
	Init() error          // 初始化
	Handle(v Entry) error // 异步时调用
	FailToPrint() bool
}

type HandlerBase struct {
	FormatType    string // 日志格式类型：json(默认),text,syslog
	TimeFormat    string // 日期格式
	FailToConsole bool   // 记录失败时打印到控制台
}

func (h *HandlerBase) FailToPrint() bool {
	return h.FailToConsole
}
