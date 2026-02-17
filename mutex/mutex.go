package mutex

// Mutex is a mutual exclusion lock built on a channel.
type Mutex struct {
	ch chan struct{}
}

func New() *Mutex {
	m := &Mutex{ch: make(chan struct{}, 1)}
	m.ch <- struct{}{}
	return m
}

func (m *Mutex) Lock() {
	<-m.ch
}

func (m *Mutex) Unlock() {
	select {
	case m.ch <- struct{}{}:
	default:
		panic("syncprims/mutex: unlock of unlocked mutex")
	}
}

func (m *Mutex) TryLock() bool {
	select {
	case <-m.ch:
		return true
	default:
		return false
	}
}

