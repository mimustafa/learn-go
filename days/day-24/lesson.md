# Day 24 — Integration, Concurrency, and Fuzz Tests

**Time budget: 55 minutes** — learn (15), run (10), exercises (20), review (10).

**Toolchain target:** Go 1.27.3; examples intentionally use APIs available in Go 1.26 where feasible.

## Learning Objectives

- Choose unit, integration, and end-to-end test boundaries for CareLog.
- Exercise HTTP behavior with `httptest` and concurrent behavior with goroutines.
- Write a fuzz target for parsers/validators and understand seed-corpus behavior.
- Use `go test -race` and recognize what it can and cannot establish.

## Why This Matters

A unit test can verify a domain rule quickly, but it does not prove that JSON decoding, routing, middleware, and that rule work together. Conversely, an enormous integration test is hard to diagnose. CareLog needs layers of tests with clear ownership of failures.

## Concepts and Explanation

Unit tests isolate a behavior using real small values and fakes only at meaningful boundaries. Integration tests connect components such as handler + application service + in-memory repository. `httptest.NewRecorder` tests a handler directly; `httptest.NewServer` exercises the real HTTP client/server path without opening a public port. Avoid testing implementation details that may change during refactoring.

For concurrent code, arrange concurrent callers, wait for them all, and assert on observable invariants. Run `go test -race ./...` for race detection; the detector only reports races encountered in executed paths and does not prove absence of all concurrency defects.

Fuzzing supplies generated inputs to a function and checks a property, such as “validator never panics” or “accepted IDs round-trip.” A fuzz test lives in `_test.go` and has signature `func FuzzName(f *testing.F)`. Seed it with meaningful edge cases, run `go test -fuzz=FuzzName -fuzztime=10s`, and keep minimized failing corpus entries. Fuzzing is not a replacement for expected-value tests.

## Runnable Example

Run `go run ./days/day-24/examples`. The sample exercises a handler through `httptest.NewServer` and verifies that concurrent creates are visible. Add `*_test.go` files and run `go test -race ./days/day-24/...`.

## Project Change

### What We Are Adding
Table-driven validation tests, HTTP contract tests, concurrent store tests, and fuzz properties for observation input.

### Why We Need It
Domain rules, HTTP translation, and concurrent storage fail at different boundaries; layered tests localize regressions while fuzzing explores unusual inputs.

### Files Added
`project/carelog/internal/domain/domain_test.go` additions, `project/carelog/internal/httpapi/observation_test.go` additions, and `project/carelog/internal/store/store_test.go` for concurrent invariants.

### Files Changed
`project/carelog/internal/domain/domain.go` or `project/carelog/internal/httpapi/observation.go` only if tests reveal a behavior gap; keep each assertion aligned with the public contract.

### Go Concepts Used
Table-driven tests, `httptest`, goroutine synchronization, race detection, and native fuzz targets.

### Tests Added
Valid/invalid validation cases, status/JSON contract tests, concurrent inserts under `-race`, and a fuzz target with explicit accepted-input invariants.

### Architecture Impact
Tests exercise the domain, HTTP, and storage seams independently; an in-memory test cannot establish PostgreSQL SQL or transaction behavior.

## Coding Exercises

1. **HTTP contract test.** Test malformed JSON, a valid observation, and an unknown route; assert status and response shape.
2. **Concurrency test.** Run at least 100 concurrent inserts and verify all unique IDs are present under `-race`.
3. **Fuzz a validator.** Write a fuzz target for a short observation code. Assert it never panics and accepted values satisfy documented length/character constraints.

## Quiz

1. What does `httptest.NewServer` add over directly calling a handler?
2. Does a clean race-detector run prove no races exist?
3. Where does Go save fuzz inputs that reproduce a failure?
4. Which test level should verify a pure validation rule?

## Test Your Understanding

For “invalid observation is rejected with HTTP 400,” identify one unit assertion and one handler integration assertion. Explain why they are not duplicates.

## Common Mistakes

- Using sleeps to coordinate goroutine tests instead of synchronization.
- Running concurrent writes and reading results before all goroutines finish.
- Mocking every dependency so no real integration is tested.
- Treating coverage percentage as a measure of test quality.
- Writing a fuzz target without any property or invariant.

## Production and Architecture Perspective

Tests define contracts at boundaries and make architectural refactoring safer. Integration tests should be deterministic, isolated, and representative. Database integration tests need their own lifecycle/cleanup strategy; an in-memory substitute cannot validate SQL behavior or transaction isolation.

## Day Summary

Good test suites combine fast unit tests, realistic boundary tests, race-enabled concurrent tests, and fuzzed input properties.

## What You Should Be Able to Do Now

Choose the right test seam and verify a concurrent or HTTP behavior without timing guesses.

**Instructor material:** [full exercise solutions and quiz key](../../instructor-guide/day-24/solutions.md).
