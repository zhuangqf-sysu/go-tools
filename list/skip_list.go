package list

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// 比较接口
type Comparable interface {
	// a > b return 1; a=b return 0; a < b return -1
	Compare(b Comparable) int
}

type SkipNode struct {
	value Comparable
	up    *SkipNode
	down  *SkipNode
	next  *SkipNode
	prev  *SkipNode
}

// 由下往上删
func (node *SkipNode) delete() {
	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	}
	if node.up != nil {
		node.up.delete()
	}
	node.up = nil
	node.prev = nil
	node.next = nil
	node.down = nil
}

func (node *SkipNode) Next() *SkipNode {
	return node.next
}

func (node *SkipNode) Prev() *SkipNode {
	return node.prev
}

func (node *SkipNode) insertBefore(value Comparable, up *SkipNode) *SkipNode {
	before := &SkipNode{
		value: value,
		up:    up,
		next:  node,
		prev:  node.prev,
	}
	if node.prev != nil {
		node.prev.next = before
	}
	if up != nil {
		up.down = before
	}
	node.prev = before
	return before
}

type skipLevel struct {
	head *SkipNode
	tail *SkipNode
	next *skipLevel
}

func newSkipLevel(up *skipLevel) *skipLevel {
	head := &SkipNode{}
	tail := &SkipNode{}
	head.next = tail
	tail.prev = head
	level := &skipLevel{head: head, tail: tail}
	if up != nil {
		head.up = up.head
		tail.up = up.tail
		up.head.down = head
		up.tail.down = tail
		up.next = level
	}
	return level
}

func (level *skipLevel) Scan(action func(value Comparable)) {
	for node := level.head.next; node != level.tail; node = node.next {
		action(node.value)
	}
}

func (level *skipLevel) String() string {
	var builder strings.Builder
	builder.WriteString("*head -> ")
	level.Scan(func(value Comparable) {
		builder.WriteString(fmt.Sprintf("%v -> ", value))
	})
	builder.WriteString("*tail")
	return builder.String()
}

type SkipList struct {
	level int
	count int
	first *skipLevel
	last  *skipLevel
}

func NewSkipList(level int) *SkipList {
	if level <= 0 {
		panic("SkipList level cann`t <= 0")
	}
	skipList := &SkipList{level: level}
	first := newSkipLevel(nil)
	last := first
	for i := 1; i < level; i++ {
		l := newSkipLevel(last)
		last = l
	}
	skipList.first = first
	skipList.last = last
	rand.Seed(time.Now().Unix())
	return skipList
}

func (list *SkipList) Len() int {
	return list.count
}

// 找到第一个 >=value 的最底层节点或者tail
func (list *SkipList) Find(value Comparable) *SkipNode {
	node := list.first.head
	l := 0
	for {
		next := node.next
		// tail || 找到了
		if next.value == nil || next.value.Compare(value) >= 0 {
			if node.down == nil {
				return next
			}
			node = node.down
			l++
			continue
		}
		node = next
	}
}

func (list *SkipList) Insert(value Comparable) *SkipNode {
	n := rand.Int()
	l := list.level - 1
	for n%2 == 0 && l >= 0 {
		l--
		n = n >> 1
	}

	var (
		up    *SkipNode
		index *SkipNode
	)

	for level := list.first; level != nil; level = level.next {
		if index == nil || index == level.head {
			index = level.head.next
		}
		for index != level.tail && index.value.Compare(value) < 0 {
			index = index.next
		}
		prev := index.prev
		if l <= 0 {
			up = index.insertBefore(value, up)
		}
		index = prev.down
		l--
	}
	return up
}

func (list *SkipList) DeleteNode(node *SkipNode) {
	node.delete()
}

func (list *SkipList) DeleteOnce(value Comparable) *SkipNode {
	node := list.Find(value)
	if node == list.last.tail || node.value.Compare(value) != 0 {
		return nil
	}
	node.delete()
	return node
}

func (list *SkipList) DeleteAll(value Comparable) int {
	node := list.Find(value)
	num := 0
	for node != list.last.tail && node.value.Compare(value) != 0 {
		temp := node.next
		node.delete()
		node = temp
		num++
	}
	return num
}

func (list *SkipList) Scan(action func(value Comparable)) {
	list.last.Scan(action)
}

func (list *SkipList) String() string {
	var builder strings.Builder
	for level := list.first; level != nil; level = level.next {
		builder.WriteString(level.String())
		builder.WriteByte('\n')
	}

	return builder.String()
}
