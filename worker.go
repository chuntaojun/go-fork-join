package fork_join

import (
	"context"
	"fmt"
	"time"
)

type Worker struct {
	pool        *Pool
	isRunning   bool
	recycleTime time.Time
	job         chan *struct {
		T Task
		F *ForkJoinTask
		C context.Context
	}
}

func (w *Worker) run(ctx context.Context) {
	go func() {

		var tmpTask *ForkJoinTask

		defer func() {
			if p := recover(); p != nil {
				w.pool.panicHandler(p)
				if tmpTask != nil {
					w.pool.err = p
					close(tmpTask.result)
				}
			}
		}()

		for {
			select {
			case <-ctx.Done():
				fmt.Println("An exception occurred and the task has stopped")
				return
			default:
				for job := range w.job {
					if job == nil {
						w.pool.workerCache.Put(w)
						return
					}
					tmpTask = job.F
					job.F.result <- job.T.Compute()
					panic("异常出现")
				}
			}
		}
	}()
}
