package main

import (
	"context"
	"fmt"
	"sync"
)

type result struct {
	index  int
	output int
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	inputs := []int{3, 1, 4, 1, 5, 9}
	const workerCount = 3
	type job struct {
		index int
		value int
	}
	jobs := make(chan job, 2)
	results := make(chan result, workerCount)
	var workers sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case item, ok := <-jobs:
					if !ok {
						return
					}
					select {
					case results <- result{index: item.index, output: item.value * item.value}:
					case <-ctx.Done():
						return
					}
				}
			}
		}()
	}
	go func() {
		defer close(jobs)
		for index, n := range inputs {
			select {
			case jobs <- job{index: index, value: n}:
			case <-ctx.Done():
				return
			}
		}
	}()
	go func() {
		workers.Wait()
		close(results)
	}()

	outputs := make([]int, len(inputs))
	for item := range results {
		outputs[item.index] = item.output
	}
	fmt.Println("input values:", inputs)
	fmt.Println("results in input order:", outputs)
}
