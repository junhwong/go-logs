package logs

import (
	"fmt"
	"time"
)

// ILog 日志记录接口
type ILog interface {
	Set(k string, v interface{}) ILog
	Prefix(v string) ILog
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
}

// Entry 日志
type Entry map[string]interface{}

func (e Entry) get(k string, defaultValue interface{}) interface{} {
	if k == "" {
		panic(ErrInvalidKey)
	}
	v, ok := e[k]
	if !ok {
		return defaultValue
	}
	return v
}
func (e Entry) log(l Level, msg string) {
	logger := e.get("__$logger", std).(Logger)
	e.Level(l)
	e.Message(msg)
	logger.Log(e)
}

func (e Entry) Set(k string, v interface{}) ILog {
	if k == "" {
		panic(ErrInvalidKey)
	}
	if v == nil {
		delete(e, k)
	} else {
		e[k] = v
	}
	return e
}

func (e Entry) GetTimestamp() time.Time { return e.get("timestamp", time.Now()).(time.Time) }
func (e Entry) GetLevel() Level         { return e.get("level", UNSET).(Level) }
func (e Entry) Level(v Level) ILog      { return e.Set("level", v) }
func (e Entry) GetMessage() string      { return e.get("message", "").(string) }
func (e Entry) Message(v string) ILog   { return e.Set("message", v) }
func (e Entry) GetPrefix() string       { return e.get("prefix", "").(string) }
func (e Entry) Prefix(v string) ILog    { return e.Set("prefix", v) }

func (e Entry) Debug(v ...interface{})                 { e.log(DEBUG, fmt.Sprint(v...)) }
func (e Entry) Debugf(format string, v ...interface{}) { e.log(DEBUG, fmt.Sprintf(format, v...)) }

func (e Entry) Info(v ...interface{})                 { e.log(INFO, fmt.Sprint(v...)) }
func (e Entry) Infof(format string, v ...interface{}) { e.log(INFO, fmt.Sprintf(format, v...)) }

func (e Entry) Warn(v ...interface{})                 { e.log(WARN, fmt.Sprint(v...)) }
func (e Entry) Warnf(format string, v ...interface{}) { e.log(WARN, fmt.Sprintf(format, v...)) }

func (e Entry) Error(v ...interface{})                 { e.log(ERROR, fmt.Sprint(v...)) }
func (e Entry) Errorf(format string, v ...interface{}) { e.log(ERROR, fmt.Sprintf(format, v...)) }

func (e Entry) Fatal(v ...interface{})                 { e.log(FATAL, fmt.Sprint(v...)) }
func (e Entry) Fatalf(format string, v ...interface{}) { e.log(FATAL, fmt.Sprintf(format, v...)) }
