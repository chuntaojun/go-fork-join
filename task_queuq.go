package fork_join

import (
	"math/rand"
	"sync"
)

type TaskQueue struct {
	cap		int32
	nodes	[][]*TaskNode
	lock	[]sync.Mutex
}

type TaskNode struct {
	task		Task
	consumer	*ForkJoinTask
}

func NewTaskQueue(workerCap int32) *TaskQueue {
	tq := &TaskQueue{cap:workerCap}
	tq.nodes = make([][]*TaskNode, workerCap)
	tq.lock = make([]sync.Mutex, workerCap)
	return tq
}

func (tq *TaskQueue) enqueue(e Task, c *ForkJoinTask) {
	id := (c.id - 1) % tq.cap
	defer func() {
		tq.lock[id].Unlock()
	}()

	// 分段锁
	tq.lock[id].Lock()
	node := &TaskNode{task:e, consumer:c}
	tq.nodes[id] = append(tq.nodes[id], node)
}

// 实现任务偷取算法，当 wId 对应的任务队列不存在任务时，采取随机策略任选一个
// 任务队列进行任务的获取，当执行到一定的次数过后，如果仍然没有找到一个有任务
// 的队列，则本次任务出队列失败，返回 (nil, nil) 值

func (tq *TaskQueue) dequeueByTali(wId int32) (bool, Task, *ForkJoinTask) {
	tq.lock[wId].Lock()

	var element *TaskNode

	if tq.HashNext(wId) {

		n := len(tq.nodes[wId]) - 1
		element = tq.nodes[wId][n]

		// Task 出队列
		tq.nodes[wId] = tq.nodes[wId][:n]
		tq.lock[wId].Unlock()

	} else {

		tq.lock[wId].Unlock()

		cnt := tq.cap
		randomWID := rand.Int31n(tq.cap)
		for len(tq.nodes[randomWID]) == 0 && cnt != 0 {
			randomWID = rand.Int31n(tq.cap)
			cnt --
		}

		if cnt == 0 {
			return false, nil, nil
		}

		tq.lock[randomWID].Lock()

		n := len(tq.nodes[randomWID]) - 1
		element = tq.nodes[randomWID][n]

		tq.nodes[randomWID] = tq.nodes[randomWID][:n]
		tq.lock[randomWID].Unlock()

	}
	return true, element.task, element.consumer
}

// 判断任务队列中是否还存在未处理的任务
func (tq *TaskQueue) HashNext(wId int32) bool {
	return len(tq.nodes[wId]) != 0
}

// 通过Hash计算key值
func (tq *TaskQueue) hashKey(v interface{}) int32 {
	return 0
}