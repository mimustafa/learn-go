package main

import (
	"errors"
	"fmt"
	"sync"
)

type observationStore struct {
	mu           sync.Mutex
	observations []string
}

func (s *observationStore) importBatch(items []string, failAt int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	start := len(s.observations)
	for i, item := range items {
		if i == failAt {
			s.observations = s.observations[:start]
			return errors.New("observation insert failed; batch rolled back")
		}
		s.observations = append(s.observations, item)
	}
	return nil
}

func main() {
	store := &observationStore{}
	_ = store.importBatch([]string{"heart-rate: 72 bpm", "sleep-hours: 8 h"}, -1)
	err := store.importBatch([]string{"steps: 1200", "temperature: 37 C"}, 1)
	fmt.Printf("observations=%v rollback_error=%v\n", store.observations, err)
}
