package main

import (
	"fmt"
	"sync"
)

type safeCounter struct {
	mu    sync.Mutex
	value int
}

func (c *safeCounter) add(delta int) {
	c.mu.Lock()
	c.value += delta
	c.mu.Unlock()
}

func (c *safeCounter) read() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func main() {
	const workers = 8
	const increments = 1000
	var wg sync.WaitGroup
	counter := &safeCounter{}
	wg.Add(workers)
	for range workers {
		go func() {
			defer wg.Done()
			for range increments {
				counter.add(1)
			}
		}()
	}
	wg.Wait()
	fmt.Printf("counter=%d (expected %d)\n", counter.read(), workers*increments)
}
