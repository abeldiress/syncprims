package cond

import "sync"

// Cond implements a condition variable for goroutines waiting for a condition.
type Cond struct {
	*sync.Cond
}

// New returns a new Cond associated with the given Locker.
func New(mu sync.Locker) *Cond {
	return &Cond{Cond: sync.NewCond(mu)}
}

// Wait atomically unlocks the locker and suspends the goroutine until Signal or Broadcast.
// The locker must be locked before calling Wait.
func (c *Cond) Wait() {
	c.Cond.Wait()
}

// Signal wakes one goroutine waiting on the condition.
func (c *Cond) Signal() {
	c.Cond.Signal()
}

// Broadcast wakes all goroutines waiting on the condition.
func (c *Cond) Broadcast() {
	c.Cond.Broadcast()
}
