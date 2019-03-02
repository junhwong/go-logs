package logs

import "fmt"

var std = New()

func Std() *StandardLogger { return std }

func Debug(v ...interface{})                 { std.New(2).Debug(v...) }
func Debugf(format string, v ...interface{}) { std.New(2).Debugf(format, v...) }

func Info(v ...interface{})                 { std.New(2).Info(v...) }
func Infof(format string, v ...interface{}) { std.New(2).Infof(format, v...) }

func Warn(v ...interface{})                 { std.New(2).Warn(v...) }
func Warnf(format string, v ...interface{}) { std.New(2).Warnf(format, v...) }

func Error(v ...interface{})                 { std.New(2).Error(v...) }
func Errorf(format string, v ...interface{}) { std.New(2).Errorf(format, v...) }

func Fatal(v ...interface{})                 { std.New(2).Fatal(v...) }
func Fatalf(format string, v ...interface{}) { std.New(2).Fatalf(format, v...) }

func Set(k string, v interface{}) ILog { return std.New(2).Set(k, v) }
func Prefix(v string) ILog             { return std.New(2).Prefix(v) }

// ==== 兼容旧程序

func Print(v ...interface{}) {
	e := std.New(2)
	e.Level(INFO)
	e.Message(fmt.Sprint(v...))
	e = std.Filter(e)
	if e == nil {
		return
	}
	formatToConsole(e, "")
}
func Printf(format string, v ...interface{}) {
	e := std.New(2)
	e.Level(INFO)
	e.Message(fmt.Sprintf(format, v...))
	e = std.Filter(e)
	if e == nil {
		return
	}
	formatToConsole(e, "")
}

func PanicIf(v error) {
	if v == nil {
		return
	}
	e := std.New(2)
	e.Error(v)
	panic(v)
}
