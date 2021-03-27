package list

import (
	"container/heap"
	"math/rand"
	"testing"
)

func randList(len int) []int {
	ans := make([]int, len, len)
	for i := 0; i < len; i++ {
		ans[i] = rand.Int()
	}
	return ans
}

func BenchmarkSkipList_Insert(b *testing.B) {
	b.StopTimer()
	arr := randList(100)
	b.StartTimer()
	list := NewSkipList(20)
	for i := 0; i < b.N; i++ {
		list.Insert(Integer(i + arr[i%100]))
	}
}

type Ints []int

func (arr Ints) Len() int {
	return len(arr)
}

func (arr Ints) Less(i, j int) bool {
	return arr[i] < arr[j]
}

func (arr Ints) Swap(i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

func (arr Ints) Push(x interface{}) {
	arr = append(arr, x.(int))
}

func (arr Ints) Pop() interface{} {
	if len(arr) == 0 {
		return nil
	}
	a := arr[0]
	arr = arr[1:]
	return a
}

func BenchmarkSkipList_Insert2(b *testing.B) {
	list := Ints(make([]int, 0))
	heap.Init(list)
	for i := 0; i < b.N; i++ {
		heap.Push(list, rand.Intn(b.N))
	}
}
