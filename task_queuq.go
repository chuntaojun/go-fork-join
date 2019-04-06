package fork_join

import (
	"sync"
)

type TaskQueue struct {
	cap		int32
	nodes	[][]*TaskNode
	lock	[]sync.Mutex
}

type TaskNode struct {
	task	Task
	consumer	*ForkJoinTask
}

func NewTaskQueue(workerCap int32) *TaskQueue {
	tq := &TaskQueue{cap:workerCap}
	tq.nodes = make([][]*TaskNode, workerCap)
	tq.lock = make([]sync.Mutex, workerCap)
	return tq
}

func (tq *TaskQueue) Enqueue(e Task, c *ForkJoinTask) {
	id := (c.id - 1) % tq.cap
	defer func() {
		tq.lock[id].Unlock()
	}()

	// 分段锁
	tq.lock[id].Lock()
	node := &TaskNode{task:e, consumer:c}
	tq.nodes[id] = append(tq.nodes[id], node)
}

func (tq *TaskQueue) DequeueByTali(wId int32) (Task, *ForkJoinTask) {
	defer func() {
		tq.lock[wId].Unlock()
	}()
	tq.lock[wId].Lock()
	n := len(tq.nodes[wId]) - 1
	element := tq.nodes[wId][n]
	tq.nodes[wId] = tq.nodes[wId][:n]
	return element.task, element.consumer
}

func (tq *TaskQueue) HashNext(wId int32) bool {
	return len(tq.nodes[wId]) != 0
}

func (tq *TaskQueue) IsEmpty() bool {
	return len(tq.nodes) == 0
}
