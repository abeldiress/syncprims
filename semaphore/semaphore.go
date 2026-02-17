package semaphore

import "context"

// Semaphore is a counting semaphore. It limits concurrent access to a resource.
type Semaphore struct {
	ch chan struct{}
}

// New returns a semaphore with n permits.
func New(n int) *Semaphore {
	if n <= 0 {
		panic("semaphore: n must be positive")
	}
	return &Semaphore{ch: make(chan struct{}, n)}
}

// Acquire acquires one permit, blocking until one is available.
func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

// Release releases one permit.
func (s *Semaphore) Release() {
	<-s.ch
}

// TryAcquire attempts to acquire one permit without blocking.
// It returns true if a permit was acquired, false otherwise.
func (s *Semaphore) TryAcquire() bool {
	select {
	case s.ch <- struct{}{}:
		return true
	default:
		return false
	}
}

// AcquireContext acquires one permit, blocking until the context is done or a permit is available.
// It returns nil on success, or ctx.Err() if the context is cancelled or times out.
func (s *Semaphore) AcquireContext(ctx context.Context) error {
	select {
	case s.ch <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
