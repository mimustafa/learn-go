# Day 21 — Synchronization and Shared State

**Time budget: 55 minutes** — learn (15), run and inspect (10), exercises (20), quiz and review (10).

**Toolchain target:** Go 1.27.3; examples intentionally use APIs available in Go 1.26 where feasible.

## Learning Objectives

- Choose between ownership by message passing and synchronization around shared state.
- Use `sync.Mutex`, `sync.RWMutex`, `sync.WaitGroup`, and typed atomics correctly.
- Explain the happens-before guarantee provided by a mutex and why atomics do not make a whole data structure safe.
- Identify a lock-order deadlock and a goroutine leak.

## Why This Matters

CareLog's in-memory observation index and background status counters may be accessed by multiple HTTP requests at once. Without a synchronization policy, map writes can race or crash. Adding a lock everywhere is not the answer either: locks introduce waiting, deadlock risk, and complexity. We need to make ownership and invariants explicit.

## Concepts and Explanation

**Shared state** is data reachable by more than one goroutine. Either give one goroutine ownership and communicate through channels, or protect shared data with a synchronization primitive. For a small in-process cache, a mutex is often simpler than a channel-based actor.

`Mutex.Lock`/`Unlock` protect an invariant, not merely a line of code. Keep the critical section short and never call unknown or blocking code while holding a lock. `RWMutex` permits concurrent readers but still excludes writers; it is useful only when read contention is meaningful and measurements support it. A plain `Mutex` can be faster and easier to reason about.

`WaitGroup` joins a known set of goroutines. Call `Add` before launching work, and pair each unit with `Done` (often deferred). It does not collect errors or cancel work. `atomic.Int64` provides indivisible operations on a single value; it cannot safely publish a related map and counter as one transaction.

Deadlocks commonly result from inconsistent lock ordering, waiting for a goroutine that needs a lock you hold, or waiting on a channel nobody will service. A goroutine leak occurs when work outlives its usefulness because nobody cancels or drains it. The race detector (`go test -race`) finds data races exercised by a test run; it does not prove that the program is race-free.

## Runnable Example

Run `go run ./days/day-21/examples`. Then run `gofmt -w days/day-21/examples/main.go` and inspect where the lock is acquired and released.

The program protects a map with `RWMutex`, joins writers with `WaitGroup`, and uses an atomic only for an independent count.

## Project Change

### What We Are Adding
Safe concurrent access to the in-memory observation cache, a detached read result, and a documented stop/join policy for cache workers.

### Why We Need It
Requests may overlap. The cache's compound map state needs protection; an independent request counter does not justify broadening that lock.

### Files Added
`project/carelog/internal/store/memory.go` for the teaching cache and `project/carelog/internal/store/memory_test.go` for concurrent behavior.

### Files Changed
`project/carelog/internal/httpapi/observation.go` should consume the cache through its existing boundary if the learner integrates this checkpoint; the standalone lesson example remains independently runnable.

### Go Concepts Used
`sync.Mutex`, `sync.RWMutex`, `sync.WaitGroup`, typed atomics, goroutine ownership, and copied values.

### Tests Added
Concurrent writer/reader test asserting the invariant after a `WaitGroup`; run it with `go test -race`.

### Architecture Impact
The cache owns its map and exposes operations, not mutable internals. This protects one process only and does not provide cross-instance coordination or durable storage.

## Coding Exercises

1. **Protect a cache.** Add a `Get` method to the example's store and make the example launch concurrent readers and writers safely.
2. **Find the deadlock.** Two functions acquire `patientMu` and `observationMu` in opposite orders. State the deadlock interleaving and fix it with one documented global lock order (or a single combined lock).
3. **Stop the leak.** A producer sends forever to a channel after its only consumer returns. Add cancellation and demonstrate that the producer exits; explain why a `WaitGroup` alone does not cancel it.

## Quiz

1. True or false: `RWMutex` is always faster than `Mutex` for read-heavy code.
2. Why should `WaitGroup.Add(1)` happen before starting the goroutine?
3. Is `atomic.Int64` sufficient to make concurrent map writes safe?
4. Name one way a goroutine can remain blocked forever even though the main function returned.

## Test Your Understanding

Without running the code, identify the invariant in `Store`, say which operations may run concurrently, and explain what can still go wrong if a function returns a map or pointer to mutable internal state.

## Common Mistakes

- Unlocking on only the success path; prefer `defer Unlock()` immediately after locking.
- Calling `Wait` before all `Add` calls have been made.
- Copying a mutex after first use or storing one in a copied value.
- Treating an atomic field as protection for adjacent fields.
- Using an `RWMutex` without evidence, or holding it across I/O.
- Assuming the race detector catches deadlocks and goroutine leaks.

## Production and Architecture Perspective

Synchronization has an operational cost and an ownership cost. Prefer immutable values and narrow lock scope; use channels when a single owner and ordered work are natural. Add timeouts and cancellation to work that can block. Measure lock contention before redesigning a cache. Never use a process-local mutex as a substitute for cross-process coordination.

## Day Summary

Correct concurrent code starts with a clear shared-state invariant. Locks protect compound state; atomics protect individual values; wait groups join work but do not stop it.

## What You Should Be Able to Do Now

Choose and justify a synchronization strategy, reason about deadlock/leak paths, and run a focused race-enabled test.

**Instructor material:** [full exercise solutions and quiz key](../../instructor-guide/day-21/solutions.md).
