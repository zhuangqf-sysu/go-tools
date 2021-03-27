### go-tools
*go 工具包，实现一些常用的数据结构*

- window 滑动窗口
- skip list 跳跃表

##### window
- 简单滑动窗口
example: 计算滑动窗口（大小为10）内的和 
```go
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

```
- 时间滑动窗口 TimeWindow
example: 统计1s内请求量+失败数+成功数
```go
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

```

##### SkipList 跳跃表  
example :
```go
    // 构造5层的跳跃表
    list := NewSkipList(5)
	list.Insert(Integer(1))
    node := list.Find(Integer(1))
    list.DeleteNode(node)
    list.DeleteOnce(Integer(1))
    list.DeleteAll(Integer(1))
```


