package fork_join

import (
	"context"
	"sync"
)

type Pool struct {
	lock 			sync.Mutex
	workerCache 	sync.Pool
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

	defer func() {
		p.lock.Unlock()
	}()

	var w *Worker
	p.lock.Lock()
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
	return w
}

