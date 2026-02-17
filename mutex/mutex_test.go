package mutex

import (
	"sync"
	"testing"
)

func TestMutex_LockUnlock(t *testing.T) {
	var m Mutex
	m.Lock()
	m.Unlock()
}

func TestMutex_TryLock(t *testing.T) {
	var m Mutex
	if !m.TryLock() {
		t.Fatal("TryLock should succeed when unlocked")
	}
	if m.TryLock() {
		t.Fatal("TryLock should fail when locked")
	}
	m.Unlock()
	if !m.TryLock() {
		t.Fatal("TryLock should succeed after Unlock")
	}
	m.Unlock()
}

func TestMutex_Concurrent(t *testing.T) {
	var m Mutex
	var counter int
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.Lock()
			counter++
			m.Unlock()
		}()
	}
	wg.Wait()
	if counter != 100 {
		t.Errorf("counter = %d, want 100", counter)
	}
}
