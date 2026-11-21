package day21

import (
	"context"
	"sync"
	"testing"
	"time"
)

type concurrentStore struct {
	mu   sync.RWMutex
	data map[int]int
}

func (s *concurrentStore) put(key, value int) {
	// A write lock protects the map from concurrent reads and writes.
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *concurrentStore) get(key int) (int, bool) {
	// Multiple readers can hold the read lock at the same time.
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.data[key]
	return value, ok
}

func TestConcurrentStore(t *testing.T) {
	store := &concurrentStore{data: make(map[int]int)}
	var workers sync.WaitGroup
	for i := 0; i < 100; i++ {
		workers.Add(1)
		go func(i int) {
			defer workers.Done()
			store.put(i, i*i)
			if _, ok := store.get(i); !ok {
				t.Errorf("missing key %d", i)
			}
		}(i)
	}
	workers.Wait()
	if len(store.data) != 100 {
		t.Fatalf("entries = %d, want 100", len(store.data))
	}
}

func TestConsistentLockOrder(t *testing.T) {
	var patientMu, observationMu sync.Mutex
	updateBoth := func() {
		// Every caller takes these locks in the same order to avoid a cycle.
		patientMu.Lock()
		defer patientMu.Unlock()
		observationMu.Lock()
		defer observationMu.Unlock()
	}
	done := make(chan struct{}, 2)
	go func() { updateBoth(); done <- struct{}{} }()
	go func() { updateBoth(); done <- struct{}{} }()
	for range 2 {
		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("consistent lock order unexpectedly blocked")
		}
	}
}

func TestCancelableProducerExits(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	out := make(chan int)
	var producer sync.WaitGroup
	producer.Add(1)
	go func() {
		defer producer.Done()
		for value := 0; ; value++ {
			select {
			case <-ctx.Done():
				return
			case out <- value:
			}
		}
	}()
	<-out
	cancel()
	done := make(chan struct{})
	go func() { producer.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("producer did not observe cancellation")
	}
}
