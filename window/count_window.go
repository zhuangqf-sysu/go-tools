package window

// window 滑动窗口，sink流向: handleOutput -> handleInput -> consume
type CountWindow struct {
	// 窗口大小
	size int
	// 窗口数据
	events []Event
	// sink 计算结果
	sink Sink
	// index 事件index
	index int

	defaultSink func() Sink
	defaultEvent func() Event

	// 处理Input事件
	handleInput HandleFunc
	// 处理Output事件
	handleOutput HandleFunc
	// 订阅消费Sink
	consume func(sink Sink)
}

func NewWindowWithCount(size int) *CountWindow {
	return &CountWindow{
		size:         size,
		events:       make([]Event,size,size),
		sink:         Sink{},
		index:        0,
	}
}

func (window *CountWindow) DefaultSink(defaultSink func() Sink) *CountWindow {
	window.defaultSink = defaultSink
	window.sink = defaultSink()
	return window
}

func (window *CountWindow) DefaultEvent(defaultEvent func() Event) *CountWindow {
	window.defaultEvent = defaultEvent
	for i:=0;i<window.size;i++ {
		window.events[i] = defaultEvent()
	}
	return window
}

func (window *CountWindow) OnInput(handleFunc HandleFunc) *CountWindow {
	window.handleInput = handleFunc
	return window
}

func (window *CountWindow) OnOutput(handleFunc HandleFunc) *CountWindow {
	window.handleOutput = handleFunc
	return window
}

func (window *CountWindow) OnConsume(f func(sink Sink)) *CountWindow {
	window.consume = f
	return window
}

func (window *CountWindow) AddEvent(input Event) {
	sink := window.sink
	window.index = (window.index+1) % window.size
	output := window.events[window.index]
	window.events[window.index] = input

	if window.handleOutput != nil {
		sink = window.handleOutput(sink, output)
	}
	if window.handleInput != nil {
		sink = window.handleInput(sink, input)
	}
	window.sink = sink

	if window.consume != nil {
		 window.consume(sink)
	}
}

func (window *CountWindow) GetSink() Sink {
	return window.sink
}

func (window *CountWindow) GetEvent() []Event {
	return window.events
}

