package fork_join

import "time"

type Worker struct {
	pool	*Pool
	isRunning	bool
	recycleTime time.Time
	job		chan *struct{
		T	Task
		F	*ForkJoinTask
	}
}

func (w *Worker) run() {
	go func() {
		defer func() {
			if p := recover(); p != nil {
			}
		}()
		for job := range w.job {
			if job == nil {
				w.pool.workerCache.Put(w)
				return
			}
			job.F.result <- job.T.Compute()
		}
	}()
}
