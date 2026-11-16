# Day 20 — Worker pools, pipelines, and backpressure

**Time: 55 minutes** — 5 min recall · 10 min concepts · 15 min walkthrough · 15 min exercises · 10 min review

**Version target:** Go 1.27.3; examples use features available in Go 1.26 where possible.

## Objectives
- Build a worker pool with a fixed concurrency bound.
- Use bounded queues to apply backpressure.
- Propagate cancellation and collect errors/results.
- Choose a policy for overload rather than spawning unbounded goroutines.

## Why this matters
Bulk CareLog exports, reminders, and imports may contain thousands of independent operations. Launching one goroutine per record can exhaust connections, memory, or downstream capacity. A worker pool makes concurrency and queueing explicit.

## Concepts and explanation
A worker pool has a producer, a bounded jobs queue, a fixed number of workers, and a result collector. Backpressure occurs when producers wait because a bounded queue is full; this limits outstanding work. Alternatives include reject, shed, persist, or rate-limit, depending on whether work is optional or must be durable. Cancellation must unblock producers and workers; errors need an explicit policy (stop, skip, aggregate, retry). Keep results tagged with input identity if completion order differs from input order. For DB work, set pool capacity consistently with DB connection limits; worker count is not free throughput.

## Runnable example
Run `go run ./days/day-20/examples` from the repository root. A fixed worker count processes a finite batch, bounds queued jobs, tags results by input position, and prints them in input order. Cancellation branches are ready for a caller to stop the pipeline.

## Project Change

### What We Are Adding
A fixed-size worker pool for CareLog observation imports with a bounded queue, cancellation, and an outcome per source observation.

### Why We Need It
One goroutine per record can exhaust connections and memory. A bounded queue makes overload visible and slows or rejects upstream work instead of accumulating without limit.

### Files Added
- `project/carelog/internal/application/workerpool_test.go` — concurrency-ceiling, cancellation, and outcome tests.

### Files Changed
- `project/carelog/internal/application/workerpool.go` — bounded workers, queue, and indexed result collection.
- `project/carelog/internal/config/config.go` — worker/queue limit validation if configurable.
- `project/carelog/internal/application/import.go` — defined batch import policy.

### Go Concepts Used
Worker pools, bounded channels, backpressure, cancellation, indexed fan-in, explicit error policy, and bounded concurrency.

### Tests Added
Observed active workers never exceed configuration; jobs/results are neither lost nor duplicated; failure policy is tested; producer/worker exit after cancellation.

### Architecture Impact
Concurrency is bounded per process and coordinated with the database pool. The pool is not durable; jobs requiring restart survival, cross-replica execution, or retries need a persistent queue and idempotent handlers.

## Exercises
1. Add an atomic active-worker counter and test observed concurrency never exceeds the configured count.
2. Add a failing input and choose/document stop-on-first-error versus collect-all behavior; test the contract.
3. Add configurable queue capacity and explain behavior when the producer outruns consumers.

## Quiz
1. How does a bounded queue create backpressure?
2. Does a worker pool guarantee input-order results?
3. Why should worker count be coordinated with database pool capacity?

## Understanding check
Can you state the maximum number of active tasks and queued tasks, name the overload policy, and prove all goroutines stop after cancellation?

## Common mistakes
- Creating one goroutine per item and calling that a worker pool.
- Using an unbounded channel or collecting unlimited results in memory.
- Forgetting to close jobs or drain results, leaving goroutines blocked.
- Returning only the first result while workers continue running.

## Production and architecture perspective
Capacity is end-to-end: API rate, worker count, queue capacity, DB pool, and dependency quotas must fit together. Monitor queue depth, wait time, failures, retries, and cancellation. If jobs must survive restarts, be retried, or run across replicas, use a durable queue and idempotent handlers.

## Summary and outcome
CareLog batch work now has a fixed concurrency budget, bounded buffering, and a clear overload/cancellation contract. **Outcome:** throughput is controlled rather than purchased with unbounded resource consumption.

## Instructor answers
See [the Day 20 instructor guide](../../instructor-guide/day-20/solutions.md) after attempting the exercises.
