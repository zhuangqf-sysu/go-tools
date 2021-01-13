package window

import (
	"fmt"
	"testing"
)

func TestNewWindowWithCount(t *testing.T) {

	window := NewWindowWithCount(10).
		DefaultEvent(func() Event {return Event{Value:0}}).
		DefaultSink(func() Sink {return Sink{Value: 0}}).
		OnInput(func(oldSink Sink, event Event) (newSink Sink) {
			value := oldSink.Value.(int)
			value += event.Value.(int)
			newSink.Value = value
			return
		}).
		OnOutput(func(oldSink Sink, event Event) (newSink Sink) {
			value := oldSink.Value.(int)
			value -= event.Value.(int)
			newSink.Value = value
			return
		}).
		OnConsume(func(sink Sink) {
			fmt.Println(sink.Value)
		})

	for i:=0;i<100;i++ {
		window.AddEvent(Event{Value:i})
	}
}