# Day 18 — Goroutines, shared state, and the race detector

**Time: 55 minutes** — 5 min recall · 10 min concepts · 15 min walkthrough · 15 min exercises · 10 min review

**Version target:** Go 1.27.3; examples use features available in Go 1.26 where possible.

## Objectives
- Explain goroutine lifecycle and synchronization of shared state.
- Identify data races and distinguish them from logical bugs.
- Protect shared CareLog counters/cache state with a mutex or redesign ownership.
- Run the race detector in tests.

## Why this matters
An HTTP server handles requests concurrently. A map or counter that was safe in a single-threaded prototype may race when handlers overlap, corrupting state or crashing the process.

## Concepts and explanation
A goroutine is a lightweight concurrent execution path, not a guarantee of parallelism or completion. A data race is unsynchronized concurrent access where at least one access writes. `sync.Mutex` protects an invariant; lock and unlock around the smallest complete critical section, and do not copy a mutex after use. `sync.RWMutex` is useful only when read contention justifies its complexity. The race detector instruments executed paths; it is valuable but cannot prove absence of races in paths it did not exercise. Prefer immutable data or a single owner when that makes invariants clearer.

## Runnable example
Run `go run ./days/day-18/examples` from the repository root. It starts concurrent increments guarded by a mutex. Validate the example with `go test -race ./days/day-18/examples`.

## Project Change

### What We Are Adding
Race-safe in-memory observation-processing state and concurrency tests for shared mutable data.

### Why We Need It
An HTTP server runs handlers concurrently; unsynchronized maps, counters, or slices may race even when a serial demo appears correct.

### Files Added
- `project/carelog/internal/store/memory_test.go` — concurrent store tests.
- `project/carelog/internal/application/counter_test.go` — concurrent counter tests if the operation has shared state.

### Files Changed
- `project/carelog/internal/store/memory.go` — synchronization or single-owner state design.
- `project/carelog/internal/application/processor.go` — explicit lifecycle for any spawned goroutines.

### Go Concepts Used
Goroutines, mutexes, critical sections, `WaitGroup`, happens-before, race reports, and immutable snapshots.

### Tests Added
Concurrent operations with deterministic final invariants, reads and writes under load, and `go test -race` on affected packages.

### Architecture Impact
State ownership is explicit and synchronization is process-local only. Database constraints/transactions remain necessary for cross-process consistency; race freedom does not establish correct business ordering.

## Exercises
1. Add a concurrent read method to the counter and ensure it is synchronized.
2. Write a test that starts many goroutines, waits for all of them, and checks the exact final value.
3. Deliberately run an unprotected version under `-race`, observe the report, then remove the race and explain the reported accesses.

## Quiz
1. What two properties make conflicting accesses a data race?
2. Does the race detector prove the whole program is race-free?
3. Why must a `WaitGroup` be incremented before launching the goroutine?

## Understanding check
Can you identify the data protected by a lock, every access that needs the same lock, and the lifecycle mechanism that ensures spawned work is observed?

## Common mistakes
- Protecting writes but not reads.
- Holding a lock across slow I/O or calling code that can re-enter the lock.
- Launching goroutines without a completion/cancellation plan.
- Assuming a race-free program has correct business ordering.

## Production and architecture perspective
Concurrency safety is an ownership/design concern before it is a mutex concern. Instrument and test realistic concurrent paths; keep critical sections short. Database transactions remain necessary for cross-process invariants—the process mutex does not coordinate replicas.

## Summary and outcome
You can protect shared in-process state and use race-enabled tests to expose unsynchronized access. **Outcome:** CareLog's concurrent request paths have an explicit synchronization or ownership strategy.

## Instructor answers
See [the Day 18 instructor guide](../../instructor-guide/day-18/solutions.md) after attempting the exercises.
