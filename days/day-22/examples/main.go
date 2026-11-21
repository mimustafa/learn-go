package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var errTemporary = errors.New("temporary example failure")

type Job struct {
	ID    string
	Value string
}

type Processor struct {
	mu        sync.Mutex
	completed map[string]string
	attempts  map[string]int
	effects   int
}

func (p *Processor) Run(ctx context.Context, job Job) (string, error) {
	p.mu.Lock()
	if result, ok := p.completed[job.ID]; ok {
		p.mu.Unlock()
		return result, nil
	}
	p.attempts[job.ID]++
	attempt := p.attempts[job.ID]
	p.mu.Unlock()

	if attempt == 1 {
		return "", errTemporary
	}
	timer := time.NewTimer(5 * time.Millisecond)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-timer.C:
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	if result, ok := p.completed[job.ID]; ok {
		return result, nil
	}
	result := "processed: " + job.Value
	p.completed[job.ID] = result
	p.effects++
	return result, nil
}

func retry(ctx context.Context, processor *Processor, job Job) (string, error) {
	delay := 10 * time.Millisecond
	var err error
	for attempt := 0; attempt < 3; attempt++ {
		var result string
		result, err = processor.Run(ctx, job)
		if err == nil || errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return result, err
		}
		if !errors.Is(err, errTemporary) {
			return "", err
		}
		if attempt == 2 {
			break
		}
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return "", ctx.Err()
		case <-timer.C:
		}
		delay *= 2
	}
	return "", err
}

func runWorkers(ctx context.Context, processor *Processor, jobs <-chan Job, count int) {
	var workers sync.WaitGroup
	for worker := 0; worker < count; worker++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			for job := range jobs {
				result, err := retry(ctx, processor, job)
				fmt.Printf("job=%s result=%q error=%v\n", job.ID, result, err)
			}
		}()
	}
	workers.Wait()
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	processor := &Processor{
		completed: make(map[string]string),
		attempts:  make(map[string]int),
	}
	jobs := make(chan Job, 3)
	jobs <- Job{ID: "report-1", Value: "patient-summary"}
	jobs <- Job{ID: "report-2", Value: "observation-summary"}
	jobs <- Job{ID: "report-3", Value: "daily-count"}
	close(jobs)

	runWorkers(ctx, processor, jobs, 2)
	result, err := retry(ctx, processor, Job{ID: "report-1", Value: "patient-summary"})
	fmt.Printf("replay=%q error=%v effects=%d\n", result, err, processor.effects)
}
