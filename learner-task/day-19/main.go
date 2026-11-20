package main

import (
	"context"
	"fmt"
	"sync"
)

type result struct{ job, worker int }

func process(ctx context.Context, jobs []int, workerCount int) ([]result, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	input := make(chan int)
	output := make(chan result)
	var workers sync.WaitGroup
	for id := 1; id <= workerCount; id++ {
		workers.Add(1)
		go func(worker int) {
			defer workers.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-input:
					if !ok {
						return
					}
					select {
					case output <- result{job: job, worker: worker}:
					case <-ctx.Done():
						return
					}
				}
			}
		}(id)
	}
	go func() {
		// The producer is the only sender, so it owns closing the input channel.
		defer close(input)
		for _, job := range jobs {
			select {
			case input <- job:
			case <-ctx.Done():
				return
			}
		}
	}()
	go func() {
		// Results close only after every worker has stopped sending.
		workers.Wait()
		close(output)
	}()
	results := make([]result, 0, len(jobs))
	for item := range output {
		results = append(results, item)
	}
	if err := ctx.Err(); err != nil {
		return results, err
	}
	return results, nil
}

func main() {
	results, err := process(context.Background(), []int{1, 2, 3, 4}, 2)
	if err != nil || len(results) != 4 {
		panic("fan-out/fan-in failed")
	}
	seen := make(map[int]bool, len(results))
	for _, item := range results {
		if seen[item.job] {
			panic("job processed more than once")
		}
		seen[item.job] = true
		fmt.Printf("job %d handled by worker %d\n", item.job, item.worker)
	}
	if len(seen) != 4 {
		panic("job result missing")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := process(ctx, []int{1, 2, 3}, 2); err == nil {
		panic("canceled pipeline unexpectedly succeeded")
	}
	fmt.Println("cancellation propagated and workers joined")
}
