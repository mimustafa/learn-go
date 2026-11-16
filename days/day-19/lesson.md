# Day 19 — Channels, select, ownership, and cancellation

**Time: 55 minutes** — 5 min recall · 10 min concepts · 15 min walkthrough · 15 min exercises · 10 min review

**Version target:** Go 1.27.3; examples use features available in Go 1.26 where possible.

## Objectives
- Use channels to pass values and signal completion.
- Use `select` to respond to cancellation and channel readiness.
- Assign clear channel ownership, including who closes a channel.
- Compose worker fan-out and result fan-in without leaking goroutines.

## Why this matters
CareLog can process independent reminders or export records concurrently. Channels let stages communicate without sharing mutable slices, but every goroutine and channel needs a clear owner and termination path.

## Concepts and explanation
A channel send/receive synchronizes communication. Unbuffered channels rendezvous; buffered channels permit a bounded number of queued values. The sender (or stage that knows no more values will be sent) closes a channel; receivers generally should not close it. Range until close for finite streams. `select` chooses among ready cases and can include `ctx.Done()` so blocked workers can stop. Cancellation is a request, not forced termination: workers must check and return. Fan-out distributes jobs; fan-in collects results. Avoid closing a shared channel while another goroutine may still send.

## Runnable example
Run `go run ./days/day-19/examples` from the repository root. Two workers consume a finite job stream, respect context cancellation, and a collector prints completion results.

## Project Change

### What We Are Adding
A cancellable in-process pipeline for a CareLog observation import/export, with explicit producer, worker, and collector ownership.

### Why We Need It
Independent observations can be processed concurrently, but blocked sends, abandoned results, and unclear channel closing can leak goroutines.

### Files Added
- `project/carelog/internal/application/pipeline_test.go` — result, cancellation, and shutdown tests.

### Files Changed
- `project/carelog/internal/application/pipeline.go` — jobs, fan-out workers, result fan-in, and context handling.
- `project/carelog/internal/application/import.go` — use the pipeline only for the defined batch operation.

### Go Concepts Used
Channels, `select`, close/ownership rules, `context.Done`, `WaitGroup`, fan-out/fan-in, and bounded buffering.

### Tests Added
Each input processed once, cancellation while blocked, result channel closure after workers stop, and completion of every goroutine.

### Architecture Impact
Pipeline is an in-process computation detail, not a durable queue. API semantics remain explicit; work that must survive restart needs a persistent job system and idempotency.

## Exercises
1. Fan out entry IDs to two workers; fan in results and verify each input appears exactly once.
2. Cancel while workers are blocked waiting for jobs and show the pipeline terminates.
3. Add a bounded output channel; explain which goroutine closes it and why no send can occur afterward.

## Quiz
1. Who should close a channel, and what does close communicate?
2. Does sending on a channel automatically stop when a context is canceled?
3. What is the difference between fan-out and fan-in?

## Understanding check
Can you draw each goroutine's input/output channels, identify its exit cases, and prove who closes every channel after all sends finish?

## Common mistakes
- Closing a channel from the receiver while senders still run.
- Forgetting to check cancellation in a potentially blocking send or receive.
- Starting goroutines whose results are never drained.
- Treating a closed channel's zero value as a valid item without checking `ok`.

## Production and architecture perspective
Every queue needs a capacity and overload policy. Cancellation must propagate through all stages. For durable or cross-process work, use a queue with persistence, retry, and idempotency semantics rather than in-memory channels.

## Summary and outcome
You can organize CareLog batch work as owned, cancellable stages with safe fan-in/fan-out. **Outcome:** concurrent work completes deliberately instead of leaving blocked goroutines behind.

## Instructor answers
See [the Day 19 instructor guide](../../instructor-guide/day-19/solutions.md) after attempting the exercises.
