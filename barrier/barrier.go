package barrier

import (
	"github.com/abeld/syncprims/cond"
	"github.com/abeld/syncprims/mutex"
)

type Barrier struct {
	n          int
	count      int
	generation uint64
	mu         *mutex.Mutex
	cv         *cond.Cond
}

func New(n int) *Barrier {
	if n <= 0 {
		panic("barrier: n must be positive")
	}
	mu := mutex.New()
	return &Barrier{
		n:  n,
		mu: mu,
		cv: cond.New(mu),
	}
}

func (b *Barrier) Wait() {
	b.mu.Lock()
	defer b.mu.Unlock()
	gen := b.generation
	b.count++
	if b.count == b.n {
		b.count = 0
		b.generation++
		b.cv.Broadcast()
		return
	}
	for gen == b.generation {
		b.cv.Wait()
	}
}

