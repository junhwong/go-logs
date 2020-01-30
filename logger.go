package logs

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"
)

// // Logger 日志处理器
// type Logger interface {
// 	New(calldepth int) Entry
// 	Append(v Handler)
// 	Log(v Entry)
// 	Filter(v Entry) Entry
// 	SetLevel(v Level)
// }

type SortedMap struct {
	Keys   []string
	Values []interface{}
	Less   func(string, string) int
}

func (m *SortedMap) find(key string) int {
	for i, it := range m.Keys {
		if m.Less(it, key) == 0 {
			return i
		}
	}
	return -1
}
func (m *SortedMap) Sort(less ...func(string, string) int) {
	c := m.Less
	if len(less) == 1 {
		c = less[0]
	}
	//strings.Compare
	n := len(m.Keys)

	for i := 0; i < n-1; i++ {
		k := m.Keys[i+1]
		v := m.Values[i+1]

		j := i
		for j >= 0 && c(k, m.Keys[j]) < 0 {
			m.Keys[j+1] = m.Keys[j]
			m.Values[j+1] = m.Values[j]
			j--
		}
		m.Keys[j+1] = k
		m.Values[j+1] = v
	}
}

func (m *SortedMap) Set(key string, value interface{}) interface{} {
	index := m.find(key)

	if index == -1 {
		m.Keys = append(m.Keys, key)
		m.Values = append(m.Values, value)

		m.Sort()
	} else {
		m.Values[index] = value
	}

	return value
}

type Extra = map[string]interface{}
type Tags = map[string]string

type LogEntry struct {
	Logger  `json:"-"`
	Extra   `json:"extra"`
	Tags    `json:"tags"`
	Prefix  string    `json:"prefix"`
	Level   Level     `json:"level"`
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
}

// Log
type Log interface {
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

// ExtendsLog
type WithLog interface {
	Log
	// Set is Deprecated
	Set(k string, v interface{}) *LogEntry

	// Set is Deprecated
	Prefix(v string) *LogEntry

	//
	WithFiled(k string, v interface{}) *LogEntry
	WithTag(name, value string) *LogEntry
}

type Logger interface {
	WithLog
	Fork(name string) Logger
	Log(entry *LogEntry)
}

var Default = &DefaultLogger{}

type DefaultLogger struct {
	name        string
	parent      *DefaultLogger
	engine      *StandardLogger
	lastLogTime time.Time
}

func (dl *DefaultLogger) Name() string {
	s := ""
	if dl.parent != nil {
		s = dl.parent.Name()
	}
	if s == "" {
		return dl.name
	}
	if dl.name != "" {
		return s + "." + dl.name
	}
	return s
}

func (dl *DefaultLogger) Fork(name string) Logger {
	return &DefaultLogger{
		parent: dl,
		name:   name,
	}
}

func (dl *DefaultLogger) New() *LogEntry {
	log := &LogEntry{
		Logger: dl,
		Tags:   make(Tags),
		Extra:  make(Extra),
	}
	log.Tags["prefix"] = dl.Name()
	log.Extra["source"] = "..."

	// 获得相对唯一的日志时间，不考虑线程竞争问题
	for {
		log.Time = time.Now()
		if !log.Time.Equal(dl.lastLogTime) {
			dl.lastLogTime = log.Time
			break
		}
		runtime.Gosched()
	}

	return log
}
func (dl *DefaultLogger) Log(log *LogEntry) {
	data, _ := json.Marshal(log)
	fmt.Println(string(data))
}

func (dl *DefaultLogger) log(l Level, msg string) {
	log := dl.New()
	log.Level = l
	log.Message = msg
	dl.Log(log)
}

func (dl *DefaultLogger) Set(k string, v interface{}) *LogEntry {
	log := dl.New()
	return log
}
func (dl *DefaultLogger) WithFiled(k string, v interface{}) *LogEntry {
	log := dl.New()
	return log
}
func (dl *DefaultLogger) WithTag(k, v string) *LogEntry {
	log := dl.New()
	return log
}

func (dl *DefaultLogger) Prefix(v string) *LogEntry {
	log := dl.New()
	return log
}

func (dl *DefaultLogger) Debug(v ...interface{}) { dl.log(DEBUG, fmt.Sprint(v...)) }
func (dl *DefaultLogger) Debugf(format string, v ...interface{}) {
	dl.log(DEBUG, fmt.Sprintf(format, v...))
}

func (dl *DefaultLogger) Info(v ...interface{}) { dl.log(INFO, fmt.Sprint(v...)) }
func (dl *DefaultLogger) Infof(format string, v ...interface{}) {
	dl.log(INFO, fmt.Sprintf(format, v...))
}

func (dl *DefaultLogger) Warn(v ...interface{}) { dl.log(WARN, fmt.Sprint(v...)) }
func (dl *DefaultLogger) Warnf(format string, v ...interface{}) {
	dl.log(WARN, fmt.Sprintf(format, v...))
}

func (dl *DefaultLogger) Error(v ...interface{}) { dl.log(ERROR, fmt.Sprint(v...)) }
func (dl *DefaultLogger) Errorf(format string, v ...interface{}) {
	dl.log(ERROR, fmt.Sprintf(format, v...))
}

func (dl *DefaultLogger) Fatal(v ...interface{}) { dl.log(FATAL, fmt.Sprint(v...)) }
func (dl *DefaultLogger) Fatalf(format string, v ...interface{}) {
	dl.log(FATAL, fmt.Sprintf(format, v...))
}
