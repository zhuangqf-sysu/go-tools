package window

type Window interface {
	AddEvent(input Event)
	GetSink() Sink
	GetEvent() []Event
} 

// Event 事件，计数主体
type Event struct {
	Value interface{}
}

// Sink 计算结果
type Sink struct {
	Value interface{}
}

type HandleFunc func(oldSink Sink, event Event) (newSink Sink)
