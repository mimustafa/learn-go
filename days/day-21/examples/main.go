package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Store struct {
	mu   sync.RWMutex
	data map[string]int
}

func (s *Store) Put(key string, value int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *Store) Snapshot() map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	copy := make(map[string]int, len(s.data))
	for key, value := range s.data {
		copy[key] = value
	}
	return copy
}

func main() {
	store := &Store{data: make(map[string]int)}
	var writes atomic.Int64
	var workers sync.WaitGroup

	for worker := 0; worker < 4; worker++ {
		workers.Add(1)
		go func(worker int) {
			defer workers.Done()
			for item := 0; item < 100; item++ {
				store.Put(fmt.Sprintf("worker-%d-item-%d", worker, item), item)
				writes.Add(1)
			}
		}(worker)
	}
	workers.Wait()
	fmt.Printf("entries=%d writes=%d\n", len(store.Snapshot()), writes.Load())
}
