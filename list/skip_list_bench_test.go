package list

import (
	"container/heap"
	"math/rand"
	"sort"
	"testing"
)

func randList(len int) []int {
	ans := make([]int, len, len)
	for i := 0; i < len; i++ {
		ans[i] = rand.Int()
	}
	return ans
}

func randIntegerList(len int) []Integer {
	ans := make([]Integer, len, len)
	for i := 0; i < len; i++ {
		ans[i] = Integer(rand.Int())
	}
	return ans
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

func BenchmarkSkipList_Insert(b *testing.B) {
	b.StopTimer()
	arr := randList(100)
	b.StartTimer()
	list := NewSkipList(17)
	for i := 0; i < b.N; i++ {
		list.Insert(Integer(i + arr[i%100]))
	}
}

func BenchmarkSkipList_Insert2(b *testing.B) {
	b.StopTimer()
	arr := randList(100)
	b.StartTimer()

	list := Ints(make([]int, 0))
	heap.Init(list)
	for i := 0; i < b.N; i++ {
		heap.Push(list, i+arr[i%100])
	}
}

func BenchmarkSkipList_Find(b *testing.B) {
	b.StopTimer()
	list := NewSkipList(6)
	for i := 0; i < 100; i++ {
		list.Insert(Integer(rand.Int()))
	}
	arr := randIntegerList(100)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		list.Find(arr[i%100])
	}
}

func BenchmarkSkipList_Find2(b *testing.B) {
	b.StopTimer()
	list := make([]int, 1000, 1000)
	for i := 0; i < 1000; i++ {
		list = append(list, rand.Int())
	}
	arr := randList(1000)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sort.SearchInts(list, arr[i%1000])
	}
}
