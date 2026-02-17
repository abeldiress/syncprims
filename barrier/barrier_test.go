package barrier

import (
	"sync"
	"testing"
)

func TestBarrier_Wait(t *testing.T) {
	b := New(3)
	var wg sync.WaitGroup
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()
			b.Wait()
		}()
	}
	wg.Wait()
}

func TestBarrier_MultipleRounds(t *testing.T) {
	b := New(2)
	var wg sync.WaitGroup
	rounds := 5
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < rounds; i++ {
			b.Wait()
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < rounds; i++ {
			b.Wait()
		}
	}()
	wg.Wait()
}
