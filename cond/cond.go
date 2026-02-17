package cond

import "github.com/abeld/syncprims/mutex"

type Cond struct {
	m       *mutex.Mutex
	waiters []chan struct{}
}

func New(m *mutex.Mutex) *Cond {
	return &Cond{m: m}
}

func (c *Cond) Wait() {
	ch := make(chan struct{}, 1)
	c.waiters = append(c.waiters, ch)
	c.m.Unlock()
	<-ch
	c.m.Lock()
}

func (c *Cond) Signal() {
	if len(c.waiters) == 0 {
		return
	}
	ch := c.waiters[0]
	copy(c.waiters[0:], c.waiters[1:])
	c.waiters = c.waiters[:len(c.waiters)-1]
	ch <- struct{}{}
}

func (c *Cond) Broadcast() {
	for _, ch := range c.waiters {
		ch <- struct{}{}
	}
	c.waiters = nil
}

