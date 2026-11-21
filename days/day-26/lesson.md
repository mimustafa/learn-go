# Day 26 — Profiling and Evidence-Based Optimization

**Time budget: 55 minutes** — concepts (15), run and inspect (10), exercises (20), review (10).

**Toolchain target:** Go 1.27.3; examples intentionally use APIs available in Go 1.26 where feasible.

## Learning Objectives

- Capture and interpret CPU and heap profiles using Go's `pprof` tooling.
- Distinguish CPU hotspots, allocation sites, and retained heap.
- Form a measurable optimization hypothesis and compare before/after behavior.
- Explain why profiles need representative workload and safe access controls.

## Why This Matters

CareLog's response path may be slow because of JSON encoding, a database query, lock contention, or an accidental quadratic algorithm. Guessing can make the wrong code more complicated. Profiling identifies where the sampled cost is occurring so we can target one bottleneck and verify the result.

## Concepts and Explanation

CPU profiles sample where the program spends CPU time. Heap profiles help inspect allocation and live-object behavior; select the appropriate profile view (`alloc_space` versus `inuse_space`) because “allocated over time” and “still retained” answer different questions. Profiles are statistical evidence, not exact accounting.

`net/http/pprof` exposes runtime profiles on HTTP endpoints when registered. Protect those endpoints: profiles can reveal code paths, resource usage, and implementation details. For a local service, bind to loopback or use a controlled admin listener; never casually expose profiling publicly. The `go tool pprof` command can inspect a saved profile, and `go tool pprof -http` can display a local web UI.

An optimization loop is: define a service-level target, capture a baseline under a representative workload, inspect the dominant cost, make one change, run correctness/race tests, and compare profiles and benchmark results under the same conditions. A faster microbenchmark is not automatically a faster service. Keep the simpler version if the gain is noise or irrelevant.

## Runnable Example

Run `go run ./days/day-26/examples`. It profiles a CPU workload into memory and writes a heap profile to a byte buffer, reporting profile sizes. It creates no profile files and is only a mechanics demonstration; use a benchmark or realistic load to investigate an actual service.

## Project Change

### What We Are Adding
A controlled CPU/heap profiling procedure and one evidence-backed optimization experiment for a representative endpoint.

### Why We Need It
The visible slow path may not be the dominant cost. Profiles help target work and before/after measurements verify relevance.

### Files Added
`project/carelog/internal/httpapi/observation_benchmark_test.go` if a focused benchmark is missing, plus a short profile capture/run procedure kept with the checkpoint notes.

### Files Changed
Change only the measured hotspot in its owning package; otherwise no production file changes are warranted.

### Go Concepts Used
`runtime/pprof`, `go tool pprof`, CPU sampling, heap profile views, benchmark comparisons, and controlled hypotheses.

### Tests Added
No new behavior is required for profiling itself; run existing tests and race tests after any optimization, and retain/add a focused benchmark.

### Architecture Impact
Profiling may inform implementation efficiency but should not prompt broader abstraction or service decomposition without separate evidence. Protect profiling endpoints and profile artifacts.

## Coding Exercises

1. **CPU profile.** Capture a profile from a representative CareLog handler, inspect the top functions, and identify one dominant application-owned cost.
2. **Heap profile.** Compare allocation-space and in-use-space views; describe a case where a short-lived high-allocation path differs from a retained-memory leak.
3. **Optimize one hotspot.** Propose, implement, and benchmark one small change. Include baseline and new measurements, environment, correctness tests, and a reason to keep or revert it.

## Quiz

1. What question does a CPU profile answer?
2. What is the difference between allocation space and in-use heap?
3. Is a 20% microbenchmark improvement proof of equivalent service-level improvement?
4. Why should pprof endpoints not be open to the public internet?

## Test Your Understanding

Heap usage grows after every request but returns to baseline after a GC. Is this evidence of retained-memory growth? What additional evidence would you gather?

## Common Mistakes

- Profiling an idle or unrepresentative system.
- Optimizing the first function in a profile without understanding the call path.
- Comparing different machines, inputs, concurrency, or build flags.
- Treating `alloc_space` as a leak report.
- Leaving pprof endpoints unauthenticated on a public listener.
- Making several changes at once and losing causal evidence.

## Production and Architecture Perspective

Profiling is an operational capability with security implications. Capture short, controlled profiles; label workloads; restrict access; and avoid collecting sensitive request payloads. Reduce algorithmic complexity or unnecessary I/O before reaching for low-level micro-optimizations. Service-level latency and error budgets determine whether a local speedup matters.

## Day Summary

Profiles guide attention; controlled benchmarks and service measurements establish whether a change helped.

## What You Should Be Able to Do Now

Capture a profile, ask the right question of it, and make or reject an optimization based on repeatable evidence.

**Instructor material:** [full exercise solutions and quiz key](../../instructor-guide/day-26/solutions.md).
