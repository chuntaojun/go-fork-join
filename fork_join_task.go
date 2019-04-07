package fork_join

import (
	"context"
	"sync/atomic"
)

type Task interface {
	Compute() interface{}
}

type ForkJoinTask struct {
	id			int32
	result 		chan interface{}
	taskPool 	*ForkJoinPool
	t			Task
	ctx         context.Context
}

func (f *ForkJoinTask) Build(taskPool *ForkJoinPool) *ForkJoinTask {
	f.taskPool = taskPool
	f.result = make(chan interface{}, 1)
	f.ctx = taskPool.ctx
	f.id = atomic.AddInt32(&taskPool.goroutineID, 1)
	return f
}

func (f *ForkJoinTask) GetTaskID() int32 {
	return f.id
}

func (f *ForkJoinTask) Run(t Task) {
	f.taskPool.taskQueue.enqueue(t, f)
}

func (f *ForkJoinTask) Join() (bool, interface{}) {
	for {
		select {
		case data, ok := <-f.result:
			if ok {
				return true, data
			}
		case <-f.ctx.Done():
			panic(f.taskPool.err)
		}
	}
}
