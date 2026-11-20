package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

type Job struct {
	Index int
	Value int
}

type Outcome struct {
	Index int
	Value int
	Err   error
}

func runPool(ctx context.Context, values []int, workers, queue int) ([]Outcome, error) {
	if workers < 1 || queue < 0 {
		return nil, errors.New("workers must be positive and queue nonnegative")
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	jobs := make(chan Job, queue)
	// Buffering limits queued work; the worker count limits active work.
	results := make(chan Outcome, workers)
	var wg sync.WaitGroup
	var active, maxActive atomic.Int64
	updateMax := func(current int64) {
		for previous := maxActive.Load(); current > previous; previous = maxActive.Load() {
			if maxActive.CompareAndSwap(previous, current) {
				return
			}
		}
	}
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-jobs:
					if !ok {
						return
					}
					current := active.Add(1)
					updateMax(current)
					outcome := Outcome{Index: job.Index, Value: job.Value * job.Value}
					if job.Value < 0 {
						outcome.Err = fmt.Errorf("negative observation value at index %d", job.Index)
					}
					active.Add(-1)
					select {
					// Cancellation also releases a worker blocked on a full result queue.
					case results <- outcome:
					case <-ctx.Done():
						return
					}
				}
			}
		}()
	}
	go func() {
		defer close(jobs)
		for index, value := range values {
			select {
			case jobs <- Job{Index: index, Value: value}:
			case <-ctx.Done():
				return
			}
		}
	}()
	go func() {
		wg.Wait()
		close(results)
	}()
	outcomes := make([]Outcome, len(values))
	for outcome := range results {
		// Indexing restores input order even if workers finish out of order.
		outcomes[outcome.Index] = outcome
	}
	if maxActive.Load() > int64(workers) {
		return outcomes, errors.New("concurrency ceiling exceeded")
	}
	if err := ctx.Err(); err != nil {
		return outcomes, err
	}
	return outcomes, nil
}

func main() {
	outcomes, err := runPool(context.Background(), []int{3, 1, -4, 5, 9}, 3, 2)
	if err != nil {
		panic(err)
	}
	for _, outcome := range outcomes {
		if outcome.Err != nil {
			fmt.Printf("item %d rejected: %v\n", outcome.Index, outcome.Err)
		} else {
			fmt.Printf("item %d result=%d\n", outcome.Index, outcome.Value)
		}
	}
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := runPool(canceled, []int{1, 2, 3}, 2, 1); err == nil {
		panic("canceled worker pool unexpectedly succeeded")
	}
	fmt.Println("bounded worker pool: fixed worker count, buffered queue, indexed outcomes")
}
