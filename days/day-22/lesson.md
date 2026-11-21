# Day 22 — Background Jobs, Retries, and Shutdown

**Time budget: 55 minutes** — concepts (15), run the example (10), exercises (20), review (10).

**Toolchain target:** Go 1.27.3; examples intentionally use APIs available in Go 1.26 where feasible.

## Learning Objectives

- Separate accepting a job from executing it, and bound concurrency with a worker count.
- Design retry behavior with an explicit limit, delay, and cancellation path.
- Explain idempotency and why a retry can repeat an operation after an ambiguous failure.
- Shut down a service without abandoning accepted work or waiting forever.

## Why This Matters

CareLog may need to generate a report or process an observation after responding to an API request. Doing all work synchronously can make requests slow; launching an unbounded goroutine per request can exhaust memory or downstream capacity. A bounded worker queue gives us explicit backpressure and a place to define completion and shutdown semantics.

## Concepts and Explanation

A job should carry stable identity and the input needed to execute it. A worker pool limits concurrent execution; the queue capacity limits buffered work. When full, the API must deliberately reject, wait, or persist the job—not silently drop it. An in-memory queue is lost on process crash and is not durable messaging.

Retry only failures that may recover. Use bounded attempts, exponential backoff with jitter in production, and respect `context.Context` during both work and delay. Retrying is safe only when the operation is idempotent or deduplicated. A request may have succeeded remotely while its response was lost; the caller cannot infer that it did not happen.

An idempotency key identifies the logical operation, not an individual network attempt. Persist the key and outcome atomically with the effect when correctness must survive process failure. A map in this lesson demonstrates the idea but is process-local and not transactional.

Graceful shutdown stops intake first, then drains accepted work up to a deadline. Decide what happens to remaining jobs when the deadline expires. Cancellation, joining, and durable recovery are distinct responsibilities.

## Runnable Example

Run `go run ./days/day-22/examples`. The example uses a bounded worker pool, cancellable retry/backoff, queue draining, and an in-memory idempotency registry; observe that replaying a completed job ID does not repeat its effect.

## Project Change

### What We Are Adding
A bounded report-job queue, stable job IDs, retry policy, and a shutdown sequence that stops intake before draining workers.

### Why We Need It
Synchronous report work can consume request time, while unbounded goroutines can exhaust resources. The queue makes capacity and rejection behavior explicit.

### Files Added
`project/carelog/internal/jobs/worker.go` and `project/carelog/internal/jobs/worker_test.go`.

### Files Changed
`project/carelog/internal/httpapi/observation.go` may submit a job after recording an observation; define whether queue-full rejects or persists instead of silently dropping work.

### Go Concepts Used
Channels, bounded worker pools, `context`, `WaitGroup`, retries with backoff, idempotency, and graceful shutdown.

### Tests Added
Tests for bounded parallelism, transient retry, repeated job IDs, cancellation, and queue drain on shutdown.

### Architecture Impact
This is explicitly an in-memory teaching queue: accepted jobs can be lost on process failure. Durable delivery requires persistent state or a broker and a recovery policy.

## Coding Exercises

1. **Bound a worker pool.** Change the example to accept four jobs on two workers and report how many executions may run simultaneously.
2. **Make retries cancellable.** Add a capped exponential delay that stops promptly when the context is canceled. Do not retry validation errors.
3. **Plan shutdown.** Write and implement the order of operations for stopping intake, draining jobs, and timing out after a deadline. State what happens to unprocessed in-memory jobs.

## Quiz

1. Why can a timeout after sending a request leave the caller unsure whether the side effect occurred?
2. Does an in-memory idempotency map guarantee exactly-once processing across restarts?
3. What is backpressure in a bounded queue?
4. Which should stop first during graceful shutdown: job intake or workers?

## Test Your Understanding

Trace one job through enqueue, execution, transient failure, retry, and shutdown. Mark the point at which the system can lose the job and the point where a duplicate can occur.

## Common Mistakes

- Infinite immediate retries that amplify an outage.
- Retrying permanent errors such as invalid input.
- Treating an idempotency key as a substitute for atomic persistence.
- Starting an unbounded goroutine per job.
- Closing a queue while senders may still be sending.
- Claiming graceful shutdown while silently dropping queued work.

## Production and Architecture Perspective

The sample is an educational, in-memory worker, not a production queue. Production choices include a database-backed outbox, broker, or durable job store, each with operational costs. Retries need limits, jitter, visibility, and dead-letter/reconciliation policy. Shutdown must fit the orchestrator's termination grace period.

## Day Summary

Reliable background work needs bounds, cancellation, idempotency, and an explicit durability contract—not just a goroutine.

## What You Should Be Able to Do Now

Explain the trade-off between synchronous work and a queue, design safe retries, and describe exactly what shutdown guarantees.

**Instructor material:** [full exercise solutions and quiz key](../../instructor-guide/day-22/solutions.md).
