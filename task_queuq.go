package fork_join

import (
	"math/rand"
)

type TaskQueue struct {
	cap		int32
	nodes	[]*LinkedQueue
}

func NewTaskQueue(workerCap int32) *TaskQueue {
	tq := &TaskQueue{
		cap:workerCap,
		nodes:make([]*LinkedQueue, workerCap, workerCap),
	}
	for i := int32(0); i < workerCap; i ++ {
		tq.nodes[i] = NewLinkedQueue()
	}
	return tq
}

func (tq *TaskQueue) enqueue(e Task, c *ForkJoinTask) {

	id := (c.id - 1) % tq.cap

	// 分段锁
	item := &Item{task:e, consumer:c}
	tq.nodes[id].put(item)
}

// 实现任务偷取算法，当 wId 对应的任务队列不存在任务时，采取随机策略任选一个
// 任务队列进行任务的获取，当执行到一定的次数过后，如果仍然没有找到一个有任务
// 的队列，则本次任务出队列失败，返回 (nil, nil) 值

func (tq *TaskQueue) dequeueByTali(wId int32) (bool, Task, *ForkJoinTask) {

	var element *Item

	if !tq.nodes[wId].isEmpty() {

		element = tq.nodes[wId].take()

	} else {

		cnt := tq.cap
		randomWID := rand.Int31n(tq.cap)
		for tq.nodes[randomWID].isEmpty() && cnt != 0 {
			randomWID = rand.Int31n(tq.cap)
			cnt --
		}

		if cnt == 0 || randomWID == wId {
			return false, nil, nil
		}

		if tq.nodes[randomWID].isEmpty() {
			return false, nil, nil
		}

		element = tq.nodes[randomWID].take()

	}
	if element == nil {
		return false, nil, nil
	}
	return true, element.task, element.consumer
}


// 通过Hash计算key值
func (tq *TaskQueue) hashKey(v interface{}) int32 {
	return 0
}