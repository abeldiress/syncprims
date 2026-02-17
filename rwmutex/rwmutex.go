package rwmutex

// RWMutex is a channel-based reader/writer lock.
type RWMutex struct {
	w       chan struct{}
	m       chan struct{}
	readers int
}

func New() *RWMutex {
	r := &RWMutex{
		w: make(chan struct{}, 1),
		m: make(chan struct{}, 1),
	}
	r.w <- struct{}{}
	r.m <- struct{}{}
	return r
}

func (r *RWMutex) lockM() {
	<-r.m
}

func (r *RWMutex) unlockM() {
	r.m <- struct{}{}
}

func (r *RWMutex) Lock() {
	<-r.w
}

func (r *RWMutex) Unlock() {
	select {
	case r.w <- struct{}{}:
	default:
		panic("syncprims/rwmutex: unlock of unlocked RWMutex")
	}
}

func (r *RWMutex) RLock() {
	r.lockM()
	defer r.unlockM()
	if r.readers == 0 {
		<-r.w
	}
	r.readers++
}

func (r *RWMutex) RUnlock() {
	r.lockM()
	defer r.unlockM()
	if r.readers <= 0 {
		panic("syncprims/rwmutex: RUnlock of unlocked RWMutex")
	}
	r.readers--
	if r.readers == 0 {
		r.w <- struct{}{}
	}
}

func (r *RWMutex) TryLock() bool {
	select {
	case <-r.w:
		return true
	default:
		return false
	}
}

func (r *RWMutex) TryRLock() bool {
	select {
	case <-r.m:
	default:
		return false
	}
	if r.readers == 0 {
		select {
		case <-r.w:
		default:
			r.m <- struct{}{}
			return false
		}
	}
	r.readers++
	r.m <- struct{}{}
	return true
}

