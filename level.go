package logs

// Level 日志级别
type Level int

func (l Level) MarshalJSON() ([]byte, error) {
	return []byte("\"" + LevelText[l] + "\""), nil
}
func (l Level) String() string {
	return LevelText[l]
}

// 日志级别
const (
	UNSET Level = iota // 未明确设置
	DEBUG              // 调试(用于开发时的辅助)
	INFO               // 普通(正常需要输出)
	WARN               // 警告(不符合预期)
	ERROR              // 错误(可恢复)
	FATAL              // 故障(程序应该挂掉)
)

var LevelText = map[Level]string{
	UNSET: "",
	DEBUG: "DEBUG",
	INFO:  "INFO ",
	WARN:  "WARN ",
	ERROR: "ERROR",
	FATAL: "FATAL",
}
