# Day 25 — Benchmarks, Allocations, and Memory

**Time budget: 55 minutes** — concepts (15), run and inspect (10), exercises (20), review (10).

**Toolchain target:** Go 1.27.3; examples intentionally use APIs available in Go 1.26 where feasible.

## Learning Objectives

- Write a repeatable Go benchmark and report allocations.
- Explain slice length/capacity, backing-array aliasing, and growth behavior.
- Read compiler escape-analysis output without treating it as a performance verdict.
- Describe the roles of the allocator and garbage collector and the cost of retaining memory.

## Why This Matters

CareLog's observation lists may be large and frequently decoded, copied, or serialized. An optimization that reduces allocations can help throughput, but premature changes can make the code less safe and still miss the actual bottleneck. Measure representative work before changing it.

## Concepts and Explanation

A benchmark named `BenchmarkX(b *testing.B)` repeats a workload. Use `b.ReportAllocs()` or `-benchmem`, keep setup out of the timed loop when appropriate, and consume results so work is not accidentally optimized away. Benchmark on representative data and compare with `benchstat` where available; one run is not proof.

A slice is a descriptor over a backing array: pointer, length, and capacity. Appending within capacity may reuse storage; growth may allocate and copy. A small subslice can retain a large backing array. Copying can reduce long-lived retention but costs an allocation and CPU. Maps and strings also have implementation-dependent allocation behavior.

Escape analysis decides whether values can live on a stack or must remain reachable on the heap. Use `go test -gcflags=-m` or `go build -gcflags=-m` to inspect compiler decisions, not as a requirement to eliminate every escape. Heap allocation is not inherently bad; clarity, lifetime, and total workload matter.

The garbage collector reclaims unreachable heap objects and trades CPU/memory overhead for automatic memory management. Reducing temporary allocation can lower GC work, while retaining large objects or pooling incorrectly can raise live heap. `runtime.MemStats` is a coarse snapshot; production monitoring should use service-level metrics and profiles.

## Runnable Example

Run `go run ./days/day-25/examples`. It compares repeated string concatenation with `strings.Builder`, reports coarse allocation deltas, and shows slice capacity. Treat timings and counts as observations, not stable promises; formal benchmarks belong in `*_test.go`.

## Project Change

### What We Are Adding
Representative projection and filtering benchmarks with allocation reporting and a recorded baseline.

### Why We Need It
The service should optimize the measured cost of representative work, not assumed allocation patterns or tiny synthetic inputs.

### Files Added
`project/carelog/internal/httpapi/observation_benchmark_test.go` for response projection and filtering benchmarks.

### Files Changed
No production code changes until a baseline and a specific bottleneck justify one; if justified, change the narrow projection/filter implementation only.

### Go Concepts Used
`testing.B`, `b.ReportAllocs`, slice length/capacity, backing-array ownership, escape analysis, and heap/GC concepts.

### Tests Added
Benchmarks report ns/op, B/op, and allocs/op; ordinary correctness tests remain the guard against behavior changes.

### Architecture Impact
Measurement guides optimization and does not change module boundaries. Avoid pooling or copying without evidence about lifetime and workload.

## Coding Exercises

1. **Benchmark a projection.** Add a benchmark that projects 1,000 observations to response values, calls `ReportAllocs`, and keeps a package-level result sink.
2. **Slice retention.** A function returns `records[:10]` from a 1-million-element allocation. Explain retention and implement a copy-based alternative; discuss when the copy is worse.
3. **Inspect an escape.** Run `go test -gcflags=-m` on a small function returning a pointer. Explain why the value escapes and why forcing a different result is not automatically an improvement.

## Quiz

1. What do slice capacity and length each mean?
2. Does every `append` allocate?
3. Does escape analysis prove an allocation is harmful?
4. Which benchmark flag reports allocations and bytes per operation?

## Test Your Understanding

You reduce allocations by retaining a 4 MiB slice for a tiny result instead of copying 24 bytes. Which memory measurement could reveal the regression, and what workload behavior makes it serious?

## Common Mistakes

- Timing setup and I/O unintentionally.
- Benchmarking one tiny input that does not represent production.
- Assuming a slice is an independent copy.
- Optimizing stack-vs-heap output without measurement.
- Using `sync.Pool` as durable storage or assuming values will remain in it.
- Calling `runtime.GC` in normal request paths to “fix” memory.

## Production and Architecture Perspective

Performance is a property of a workload and service objective, not a code-style preference. Track latency distributions, throughput, and memory under realistic concurrency. The standard library provides enough tooling to start; compare before/after under the same environment and retain only changes that preserve behavior and materially improve the target.

## Day Summary

Benchmarks measure a controlled operation; allocation counts and profiles help explain cost; slices and object lifetimes determine retention.

## What You Should Be Able to Do Now

Establish a benchmark baseline, reason about slice ownership, and investigate escapes without cargo-cult optimization.

**Instructor material:** [full exercise solutions and quiz key](../../instructor-guide/day-25/solutions.md).
