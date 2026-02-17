package barrier

import "sync"

// Barrier blocks until n goroutines have called Wait, then releases them all.
type Barrier struct {
	n           int
	count       int
	mu          sync.Mutex
	cond        *sync.Cond
	generation  uint64
}

// New returns a barrier that waits for n goroutines.
func New(n int) *Barrier {
	if n <= 0 {
		panic("barrier: n must be positive")
	}
	b := &Barrier{n: n}
	b.cond = sync.NewCond(&b.mu)
	return b
}

// Wait blocks until n goroutines have called Wait. When the nth goroutine calls Wait,
// all n are released. Wait may be called again for the next "round".
func (b *Barrier) Wait() {
	b.mu.Lock()
	defer b.mu.Unlock()
	gen := b.generation
	b.count++
	if b.count == b.n {
		b.count = 0
		b.generation++
		b.cond.Broadcast()
		return
	}
	for gen == b.generation {
		b.cond.Wait()
	}
}
