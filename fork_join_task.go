package fork_join

import (
	"sync/atomic"
)

type Task interface {
	Compute() interface{}
}

type ForkJoinTask struct {
	id	int32
	result 	chan interface{}
	taskPool *ForkJoinPool
	t	Task
}

func (f *ForkJoinTask)Build(taskPool *ForkJoinPool, parentId *int32) *ForkJoinTask {
	f.taskPool = taskPool
	f.result = make(chan interface{}, 1)
	f.id = atomic.AddInt32(parentId, 1)
	return f
}

func (f *ForkJoinTask) GetTaskID() int32 {
	return f.id
}

func (f *ForkJoinTask) Run(t Task) {
	f.taskPool.taskQueue.Enqueue(t, f)
}

func (f *ForkJoinTask) Join() interface{} {
	return <-f.result
}
