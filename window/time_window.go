package window

import (
	"context"
	"log"
	"runtime"
	"time"
)

/**
* 时间滑动窗口
* 数据流： OnReduce -> OnInput -> -> OnOutput -> Sink
*/
type TimeWindow struct {
	CountWindow

	bufSize int
	timeSpan time.Duration

	ctx context.Context
	cancel context.CancelFunc
	eventChan chan Event

	handleReduce HandleFunc
}

// todo 配置 -> option的形式选择性配置
// todo 错误处理
// timeSpan 每个小窗口的时间长度
// size 窗口的数量，timeSpan*size为整个窗口的时间长度
// bufSize chan的大小，超过长度则阻塞 todo 配置化选择阻塞或者丢弃
// ctx 上下文，可通过ctx stop goroutine
func NewTimeWindow(ctx context.Context, timeSpan time.Duration, size int, bufSize int) *TimeWindow {
	cancelCtx, cancelFunc := context.WithCancel(ctx)
	return &TimeWindow{
		CountWindow: *NewWindowWithCount(size),
		bufSize:     bufSize,
		timeSpan:    timeSpan,
		ctx:         cancelCtx,
		cancel:      cancelFunc,
		eventChan:   make(chan Event, timeSpan),
	}
}

// OnReduce 聚合每个时间碎片的数据
func (window *TimeWindow) OnReduce(handler HandleFunc) *TimeWindow {
	window.handleReduce = handler
	return window
}

func (window *TimeWindow) AddEvent(input Event) {
	window.eventChan <- input
}

func (window *TimeWindow) Run() {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 64<<10)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("panic in signal proc, err: %v, stack: %s\n", r, buf)
		}
	}()
	timer := time.NewTicker(window.timeSpan)
	reduceSink := Sink{Value: window.defaultEvent().Value}
	for {
		select {
		case event := <-window.eventChan:
			reduceSink = window.handleReduce(reduceSink, event)
		case <-timer.C:
			input := Event{Value: reduceSink.Value}
			window.CountWindow.AddEvent(input)
			reduceSink.Value = window.defaultEvent().Value
		case <-window.ctx.Done():
			return
		}
	}
}

func (window *TimeWindow) Start()  {
	go window.Run()
}

func (window *TimeWindow) Stop()  {
	window.cancel()
}