package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	mu    sync.RWMutex
	value int
}

func (c *Counter) Add(delta int) {
	// Both reads and writes of value use the same lock.
	c.mu.Lock()
	c.value += delta
	c.mu.Unlock()
}

func (c *Counter) Read() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.value
}

func main() {
	const workers, increments = 12, 1000
	var wg sync.WaitGroup
	counter := &Counter{}
	wg.Add(workers)
	for range workers {
		go func() {
			defer wg.Done()
			for range increments {
				counter.Add(1)
				_ = counter.Read()
			}
		}()
	}
	wg.Wait()
	want := workers * increments
	if got := counter.Read(); got != want {
		panic(fmt.Sprintf("got %d, want %d", got, want))
	}
	fmt.Printf("concurrent read/write invariant passed: %d\n", want)
}
