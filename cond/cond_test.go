package cond

import (
	"sync"
	"testing"
)

func TestCond_Signal(t *testing.T) {
	var mu sync.Mutex
	c := New(&mu)
	done := make(chan bool)
	mu.Lock()
	go func() {
		mu.Lock()
		c.Wait()
		done <- true
		mu.Unlock()
	}()
	mu.Unlock()
	// Give goroutine time to block on Wait
	mu.Lock()
	mu.Unlock()
	c.Signal()
	if <-done != true {
		t.Fatal("Wait did not return")
	}
}

func TestCond_Broadcast(t *testing.T) {
	var mu sync.Mutex
	c := New(&mu)
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
