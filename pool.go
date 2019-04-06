package fork_join

import (
	"sync"
)

type Pool struct {

	lock sync.Mutex

	workerCache sync.Pool

	PanicHandler func(interface{})
}

func NewPool() *Pool {
	p := &Pool{}
	return p
}

func (p *Pool) retrieveWorker() *Worker {

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
			}, 1),
		}
	}
	w.run()
	return w
}

