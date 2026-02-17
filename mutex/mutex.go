package mutex

import "sync"

// Mutex provides mutual exclusion with Lock, Unlock, and TryLock.
type Mutex struct {
	mu sync.Mutex
}

// Lock locks the mutex. If the lock is already in use, the calling goroutine blocks.
func (m *Mutex) Lock() {
	m.mu.Lock()
}

// Unlock unlocks the mutex.
func (m *Mutex) Unlock() {
	m.mu.Unlock()
}

// TryLock attempts to lock the mutex without blocking.
// It returns true if the lock was acquired, false otherwise.
func (m *Mutex) TryLock() bool {
	return m.mu.TryLock()
}
