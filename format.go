package logs

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

func foramtTime(t time.Time, timeFormat string) string {
	if timeFormat == "timestamp" {
		return fmt.Sprintf("%d", t.Unix())
	}
	if timeFormat != "" {
		return t.Format(timeFormat)
	}
	return t.Format("2006-01-02 15:04:05")
}

func formatText(e Entry, timeFormat, appendDelimiter string) string {
	t := foramtTime(e.GetTimestamp(), timeFormat)
	msg := formatMessage(e)
	l := e.GetLevel()
	return fmt.Sprintf("%s\t[%s]\t%s\t%s\n%+v%s", t, l, e["source"], e["method"], msg, appendDelimiter)
}

func formatToConsole(e Entry, timeFormat string) {
	w := os.Stdout
	t := foramtTime(e.GetTimestamp(), timeFormat)
	msg := formatMessage(e)
	l := e.GetLevel()
	var cl string

	// https://en.wikipedia.org/wiki/ANSI_escape_code#Colors
	switch l {
	case DEBUG:
		cl = "\033[1;34;40m" + l.String() + "\033[0m"
	case INFO:
		cl = "\033[1;32;40m" + l.String() + "\033[0m"
	case WARN:
		cl = "\033[1;33;40m" + l.String() + "\033[0m"
		w = os.Stderr
	case ERROR:
		cl = "\033[1;31;40m" + l.String() + "\033[0m"
		w = os.Stderr
	case FATAL:
		cl = "\033[1;31;43m" + l.String() + "\033[0m"
		w = os.Stderr
	default:
		cl = l.String()
	}

	fmt.Fprintf(w, "[%s %s] %s\n%+v\n", cl, t, e["source"], msg) // 忽略打印失败
}

func formatMessage(e Entry) string {
	s := e.GetMessage()
	prefix := e.GetPrefix()
	arr := []string{}
	for k, v := range e {
		if k == "level" || k == "source" || k == "method" || k == "message" || k == "timestamp" || k == "time" {
			continue
		}
		arr = append(arr, fmt.Sprintf("%s=%+v", k, v))
	}
	if prefix != "" {
		s = prefix + ":" + s
	}
	if len(arr) > 0 {
		s += "\n"
		s += strings.Join(arr, "\n")
	}
	return s
}

func formatJson(e Entry, timeFormat string) ([]byte, error) {
	if timeFormat == "" {
		timeFormat = "timestamp"
	}
	t := e.GetTimestamp()
	defer func() {
		e.Set("timestamp", t)
	}()
	e.Set("timestamp", foramtTime(t, timeFormat))
	return json.Marshal(e)
}

func replaceTimeFormat(src string) string {
	//%Y年 %m月 %d日 %H时 %M分 %S秒 %z时区
	//2006-01-02 15:04:05
	src = strings.Replace(src, "%Y", "2006", -1)
	src = strings.Replace(src, "%m", "01", -1)
	src = strings.Replace(src, "%d", "02", -1)
	src = strings.Replace(src, "%H", "15", -1)
	src = strings.Replace(src, "%M", "04", -1)
	src = strings.Replace(src, "%S", "05", -1)
	src = strings.Replace(src, "%z", "-0700", -1)
	return src
}
