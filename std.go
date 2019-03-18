package logs

import (
	"context"
	"errors"
	"log"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// StandardLogger 默认日志器
type StandardLogger struct {
	*sync.Mutex
	stream   chan Entry
	cancel   context.CancelFunc
	wg       *sync.WaitGroup
	handlers []Handler
	level    Level
}

func (l *StandardLogger) pathFilter(s string) string {
	return s
}

// New 创建一条日志。
//
// calldepth 表示调用栈的深度，-1 将不分析。
func (l *StandardLogger) New(calldepth int) Entry {
	e := Entry{
		"timestamp": time.Now(),
		"__$logger": l,
	}
	if calldepth > -1 {
		pc, file, line, ok := runtime.Caller(calldepth)
		if ok {
			if line > 0 {
				e["source"] = l.pathFilter(file + ":" + strconv.Itoa(line))
			} else {
				e["source"] = l.pathFilter(file)
			}
			e["method"] = l.pathFilter(runtime.FuncForPC(pc).Name())
		}
	}
	return e
}

// Log 实现 Logger 的接口方法
func (l *StandardLogger) Log(v Entry) {
	l.stream <- v
}

// Append 添加一个处理器或过滤器
func (l *StandardLogger) Append(v Handler) {
	if v == nil {
		panic(errors.New("Handler doent a nil"))
	}
	l.Lock()
	defer l.Unlock()
	if err := v.Init(); err != nil {
		panic(err)
	}
	l.handlers = append(l.handlers, v)
}

func (l *StandardLogger) run(ctx context.Context) {
	l.wg.Add(1)
	defer l.wg.Done()
	var e Entry
	for {
		select {
		case e = <-l.stream:
			e = l.Filter(e)
			if e == nil {
				break
			}
			if len(l.handlers) == 0 {
				formatToConsole(e, "")
				break
			}
			for _, h := range l.handlers {
				err := h.Handle(e)
				if err != nil {
					log.Printf("logs.handle: %+v\n", err)
				}
				if err != nil && h.FailToPrint() {
					formatToConsole(e, "")
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func (l *StandardLogger) Filter(v Entry) Entry {
	delete(v, "__$logger")
	if v.GetLevel() < l.level {
		return nil
	}
	return v
}

func (l *StandardLogger) SetLevel(v Level) {
	l.level = v
}

func (l *StandardLogger) GracefulStop() {
	l.cancel()
	l.wg.Wait()
}

// New 创建一个新的
func New() *StandardLogger {
	l := &StandardLogger{
		Mutex:    new(sync.Mutex),
		stream:   make(chan Entry, 1024*4),
		wg:       &sync.WaitGroup{},
		handlers: []Handler{},
		level:    -1,
	}
	ctx, cancel := context.WithCancel(context.Background())
	l.cancel = cancel
	go l.run(ctx)
	return l
}
