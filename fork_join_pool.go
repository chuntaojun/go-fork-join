package fork_join

import "sync"

type ForkJoinPool struct {
	Name		string
	cap			int32
	taskQueue	*TaskQueue
	wp			*Pool
	lock		sync.Mutex
	signal		*sync.Cond
	once		sync.Once
}

func NewForkJoinPool(name string, workerCap int32) *ForkJoinPool {
	fp := &ForkJoinPool{
		Name:      name,
		cap:       workerCap,
		taskQueue: NewTaskQueue(workerCap),
	}
	fp.wp = NewPool()
	fp.signal = sync.NewCond(new(sync.Mutex))
	fp.run()
	return fp
}

func (fp *ForkJoinPool) pushTask(t Task, f *ForkJoinTask) {
	fp.taskQueue.Enqueue(t, f)
}

func (fp *ForkJoinPool) run() {
	go func() {
		wId := int32(0)
		for {
			if fp.taskQueue.HashNext(wId) {
				job, ft := fp.taskQueue.DequeueByTali(wId)
				fp.wp.retrieveWorker().job <- &struct {
					T Task
					F *ForkJoinTask
				}{T: job, F: ft}
			}
			wId = (wId + 1) % fp.cap
		}
	}()
}