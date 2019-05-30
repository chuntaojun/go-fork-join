package fork_join

import (
	"sync"
	"sync/atomic"
)

type LinkedQueue struct {
	head		*Node
	last		*Node
	size		int32
	takeLock	sync.Mutex
	putLock		sync.Mutex
}

type Node struct {
	item		*Item
	next		*Node
}

type Item struct {
	task		Task
	consumer	*ForkJoinTask
}

func NewLinkedQueue() *LinkedQueue {
	linkedQueue := &LinkedQueue{}
	linkedQueue.head = &Node{item:nil}
	linkedQueue.last = linkedQueue.head
	return linkedQueue
}

func (l *LinkedQueue) put(item *Item)  {

	defer func() {
		l.putLock.Unlock()
	}()

	l.putLock.Lock()
	node := &Node{item:item}
	l.last.next = node
	l.last = node
	atomic.AddInt32(&l.size, 1)
}

func (l *LinkedQueue) take() *Item {

	defer func() {
		l.takeLock.Unlock()
	}()

	l.takeLock.Lock()

	if l.getSize() == 0 {
		return nil
	}

	h := l.head
	first := h.next
	h.next = h
	l.head = first
	e := first.item
	first.item = nil
	atomic.AddInt32(&l.size, -1)
	return e

}

func (l *LinkedQueue) isEmpty() bool {
	return l.getSize() == 0
}

func (l *LinkedQueue) getSize() int32 {
	return atomic.LoadInt32(&l.size)
}