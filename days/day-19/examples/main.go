package main

import (
	"context"
	"fmt"
	"sync"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	jobs := make(chan int)
	results := make(chan int)
	var workers sync.WaitGroup

	for workerID := 1; workerID <= 2; workerID++ {
		workers.Add(1)
		go func(id int) {
			defer workers.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-jobs:
					if !ok {
						return
					}
					select {
					case results <- job * 2:
					case <-ctx.Done():
						return
					}
				}
			}
		}(workerID)
	}
	go func() {
		for _, job := range []int{1, 2, 3, 4} {
			select {
			case jobs <- job:
			case <-ctx.Done():
				close(jobs)
				return
			}
		}
		close(jobs)
	}()
	go func() {
		workers.Wait()
		close(results)
	}()
	for result := range results {
		fmt.Println("processed:", result)
	}
}
