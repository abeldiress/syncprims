package cond

import (
	"testing"

	"github.com/abeld/syncprims/mutex"
)

func TestCond_Signal(t *testing.T) {
	mu := mutex.New()
	c := New(mu)
	done := make(chan bool, 1)

	mu.Lock()
	go func() {
		mu.Lock()
		c.Wait()
		done <- true
		mu.Unlock()
	}()
	mu.Unlock()

	// Give goroutine time to block on Wait by briefly
	// taking and releasing the lock again.
	mu.Lock()
	mu.Unlock()

	c.Signal()
	if !<-done {
		t.Fatal("Wait did not return")
	}
}

func TestCond_Broadcast(t *testing.T) {
	mu := mutex.New()
	c := New(mu)
	n := 5
	done := make(chan int, n)

	for i := 0; i < n; i++ {
		go func() {
			mu.Lock()
			c.Wait()
			done <- 1
			mu.Unlock()
		}()
	}

	mu.Lock()
	mu.Unlock()

	c.Broadcast()
	for i := 0; i < n; i++ {
		<-done
	}
}

