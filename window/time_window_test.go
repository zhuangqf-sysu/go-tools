package window

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"math/rand"
	"testing"
	"time"
)

type Counter struct {
	total int64
	failure int64
	success int64
}

func TestNewTimeWindow(t *testing.T) {
	group, _ := errgroup.WithContext(context.Background())

	duration := 100 * time.Microsecond
	window := NewTimeWindow(context.Background() ,duration, 10, 1024)
	window.
		DefaultEvent(func() Event {return Event{&Counter{}}}).
		DefaultSink(func() Sink {return Sink{&Counter{}}}).
		OnInput(func(oldSink Sink, event Event) (newSink Sink) {
			counter := oldSink.Value.(*Counter)
			value := event.Value.(*Counter)
			counter.failure += value.failure
			counter.success += value.success
			counter.total += value.total

			newSink.Value = counter
			return
		}).
		OnOutput(func(oldSink Sink, event Event) (newSink Sink) {
			counter := oldSink.Value.(*Counter)
			value := event.Value.(*Counter)
			counter.failure -= value.failure
			counter.success -= value.success
			counter.total -= value.total

			newSink.Value = counter
			return
		}).
		OnConsume(func(sink Sink) {
			fmt.Printf("%v: %v \n",time.Now(), sink.Value)
		})
	window.OnReduce(func(oldSink Sink, event Event) (newSink Sink) {
		counter := oldSink.Value.(*Counter)
		counter.total++
		value := event.Value.(bool)
		if value {
			counter.success++
		}else {
			counter.failure++
		}
		newSink.Value = counter
		return
	})
	window.Start()

	for i:=0;i<10;i++ {
		group.Go(func() error {
			rand.Seed(time.Now().Unix())
			for i:=0;i<1000+rand.Intn(100000);i++{
				window.AddEvent(Event{
					Value: rand.Intn(2) % 2 == 0,
				})
				time.Sleep(time.Microsecond)
			}
			return nil
		})
	}
	_ = group.Wait()
	window.Stop()
	fmt.Println(window.GetSink().Value)
}