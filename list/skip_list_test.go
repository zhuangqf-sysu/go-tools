package list

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
)

type Integer int

func (i Integer) Compare(b Comparable) int {
	return int(i) - int(b.(Integer))
}

func TestSkipNode_insertBefore(t *testing.T) {
	level := newSkipLevel(nil)
	for i := 0; i < 10; i++ {
		level.head.insertBefore(Integer(i), nil)
	}
	fmt.Println(level)
}

func TestSkipList_Insert(t *testing.T) {
	n := 100
	arr := make([]int, n, n)
	list := NewSkipList(5)
	for i := 0; i < n; i++ {
		arr[i] = rand.Intn(n)
		list.Insert(Integer(arr[i]))
	}
	sort.Ints(arr)
	fmt.Println(arr)
	fmt.Println(list)

	i := 0
	list.Scan(func(value Comparable) {
		if int(value.(Integer)) != arr[i] {
			t.Errorf("err: arr[%d](%d) != list[%d](%v)", i, arr[i], i, value)
		}
		i++
	})
}

func TestSkipList_Find(t *testing.T) {
	n := 100
	arr := make([]int, n, n)
	list := NewSkipList(5)
	for i := 0; i < n; i++ {
		arr[i] = rand.Intn(n)
		list.Insert(Integer(arr[i]))
	}
	sort.Ints(arr)
	fmt.Println(arr)
	fmt.Println(list)
	for j := 0; j < 5; j++ {
		a := rand.Intn(n+100) - 50
		node := list.Find(Integer(a))
		fmt.Printf("find %d :", a)
		for i := node; i != nil; i = i.next {
			fmt.Printf("%v -> ", i.value)
		}
		fmt.Println()
		if a < 0 || a == 0 && arr[0] != 0 {
			if node != list.last.head.next {
				t.Errorf("err!!! node(%v) should be head.next", node.value)
			}
			continue
		}
		if a >= 100 {
			if node != list.last.tail {
				t.Errorf("err!!! node(%v) should be tail", node.value)
			}
			continue
		}
		num := 0
		for i := list.last.head.next; i != node; i = i.next {
			if i.value.Compare(node.value) >= 0 {
				t.Errorf("find err: list[%d](%v) >= node(%v)", num, i.value, node.value)
			}
			num++
		}
	}
}

func TestSkipList_DeleteNode(t *testing.T) {
	n := 100
	arr := make([]int, n, n)
	list := NewSkipList(5)
	for i := 0; i < n; i++ {
		arr[i] = rand.Intn(n)
		list.Insert(Integer(arr[i]))
	}

	sort.Ints(arr)

	node := list.last.head.next
	for j := 0; j < n; j++ {
		if j%10 == 0 {
			arr[j] = -1
			tmp := node
			node = node.next
			list.DeleteNode(tmp)
		} else {
			node = node.next
		}
	}

	fmt.Println(arr)
	fmt.Println(list)

	i := 0
	list.Scan(func(value Comparable) {
		if i%10 == 0 {
			i++
		}
		if int(value.(Integer)) != arr[i] {
			t.Errorf("err: arr[%d](%d) != list[%d](%v)", i, arr[i], i, value)
		}
		i++
	})
}

func TestSkipList_DeleteOnce(t *testing.T) {
	n := 100
	arr := make([]int, n, n)
	list := NewSkipList(5)
	for i := 0; i < n; i++ {
		arr[i] = rand.Intn(n)
		list.Insert(Integer(arr[i]))
	}

	for j := 0; j < n; j++ {
		if j%10 == 0 {
			list.DeleteOnce(Integer(arr[j]))
			arr[j] = -1
		}
	}

	sort.Ints(arr)
	fmt.Println(arr)
	fmt.Println(list)

	i := 0
	list.Scan(func(value Comparable) {
		for arr[i] == -1 {
			i++
		}
		if int(value.(Integer)) != arr[i] {
			t.Errorf("err: arr[%d](%d) != list[%d](%v)", i, arr[i], i, value)
		}
		i++
	})
}

func TestSkipList_DeleteAll(t *testing.T) {
	n := 100
	arr := make([]int, n, n)
	list := NewSkipList(5)
	for i := 0; i < n; i++ {
		arr[i] = rand.Intn(n)
		list.Insert(Integer(arr[i]))
	}

	sort.Ints(arr)
	fmt.Println(arr)
	fmt.Println(list)

	list.DeleteAll(Integer(10))

	i := 0
	list.Scan(func(value Comparable) {
		for arr[i] == 10 {
			i++
		}
		if int(value.(Integer)) != arr[i] {
			t.Errorf("err: arr[%d](%d) != list[%d](%v)", i, arr[i], i, value)
		}
		i++
	})
}
