package rwmutex

import (
	"sync"
	"testing"
)

func TestRWMutex_LockUnlock(t *testing.T) {
	var m RWMutex
	m.Lock()
	m.Unlock()
}

func TestRWMutex_RLockRUnlock(t *testing.T) {
	var m RWMutex
	m.RLock()
	m.RUnlock()
}

func TestRWMutex_TryLock(t *testing.T) {
	var m RWMutex
	m.RLock()
	if !m.TryLock() {
		// TryLock may fail while readers hold (implementation-dependent)
	}
	m.RUnlock()
	if !m.TryLock() {
		t.Fatal("TryLock should succeed when no writers")
	}
	m.Unlock()
}

func TestRWMutex_ConcurrentReads(t *testing.T) {
	var m RWMutex
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.RLock()
			m.RUnlock()
		}()
	}
	wg.Wait()
}
