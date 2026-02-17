package semaphore

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestSemaphore_AcquireRelease(t *testing.T) {
	s := New(2)
	s.Acquire()
	s.Acquire()
	done := make(chan struct{})
	go func() {
		s.Release()
		s.Release()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Release blocked")
	}
}

func TestSemaphore_LimitsConcurrency(t *testing.T) {
	s := New(2)
	var count int
	var mu sync.Mutex
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Acquire()
			defer s.Release()
			mu.Lock()
			count++
			if count > 2 {
				t.Error("more than 2 concurrent")
			}
			mu.Unlock()
			time.Sleep(10 * time.Millisecond)
			mu.Lock()
			count--
			mu.Unlock()
		}()
	}
	wg.Wait()
}

func TestSemaphore_TryAcquire(t *testing.T) {
	s := New(1)
	if !s.TryAcquire() {
		t.Fatal("TryAcquire should succeed")
	}
	if s.TryAcquire() {
		t.Fatal("TryAcquire should fail when no permit")
	}
	s.Release()
	if !s.TryAcquire() {
		t.Fatal("TryAcquire should succeed after Release")
	}
	s.Release()
}

func TestSemaphore_AcquireContext(t *testing.T) {
	s := New(1)
	s.Acquire()
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	err := s.AcquireContext(ctx)
	if err != context.DeadlineExceeded {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
	s.Release()
	err = s.AcquireContext(context.Background())
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	s.Release()
}
