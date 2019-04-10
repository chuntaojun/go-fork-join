package fork_join

import (
	"context"
	"sync"
)

type Pool struct {
	lock 			sync.Mutex
	workerCache 	sync.Pool
	workers			[]*Worker
	cancel			context.CancelFunc
	panicHandler 	func(interface{})
	err				interface{}
}

func newPool(cancel context.CancelFunc) *Pool {
	p := &Pool{
		cancel:cancel,
	}
	p.panicHandler = func(i interface{}) {
		p.cancel()
		p.err = i
	}
	return p
}

func (p *Pool) retrieveWorker(ctx context.Context) *Worker {

	var w *Worker

	p.lock.Lock()
	idleWorker := p.workers
	n := len(idleWorker) - 1
	if n >= 0 {
		w = idleWorker[n]
		p.workers = idleWorker[:n]
		p.lock.Unlock()
	} else {
		p.lock.Unlock()
		if cacheWorker := p.workerCache.Get(); cacheWorker != nil {
			w = cacheWorker.(*Worker)
		} else {
			w = &Worker{
				pool: p,
				job: make(chan *struct {
					T Task
					F *ForkJoinTask
					C context.Context
				}, 1),
			}
		}
		w.run(ctx)
	}
	return w
}

func (p *Pool) releaseWorker(worker *Worker) {
	p.lock.Lock()
	p.workers = append(p.workers, worker)
	p.lock.Unlock()
}

