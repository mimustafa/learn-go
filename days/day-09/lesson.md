# Day 09 — Unit tests, table tests, and fakes

**Version target:** Go 1.27.3; examples use features available in Go 1.26 where possible.

## Objectives
- Write a standard-library Go test and run it with `go test`.
- Use table-driven cases to cover behavior, including errors.
- Replace an external dependency with a small fake at a justified interface seam.
- **Time budget: 55 minutes** — read 12, code/run 15, exercises 18, quiz/review 10.

## Why
CareLog needs confidence that observation rules remain stable as its design evolves. Fast unit tests make the behavior executable documentation without requiring a database.

## Concepts and explanation
- Test files end in `_test.go`; test functions are named `TestXxx` and receive `*testing.T`. `go test` compiles and runs them.
- Table-driven tests make related inputs and expected outcomes visible. Give cases descriptive names and stop a case with `t.Fatalf` when later checks cannot be meaningful.
- Test observable behavior, not private implementation details. Use `t.Errorf` to report an issue and continue, `t.Fatalf` to stop the current test.
- A fake is a small in-memory implementation that gives controlled behavior. Use one when it tests a real seam; avoid elaborate mocks that mirror implementation internals.
- Tests should not depend on execution order, current time, map iteration order, or external services unless deliberately testing an integration boundary.

## Runnable example
Run `go run ./days/day-09/examples` for the demonstration. Then run `go test ./days/day-09/examples` to execute the included table-driven tests.

## Project Change

### What We Are Adding
A tested CareLog checkpoint: table-driven tests specify observation validation, and a fake writer verifies save behavior and error propagation without a database.

### Why We Need It
Fast deterministic tests catch regressions in domain rules and application orchestration while keeping external services out of unit-test feedback loops.

### Files Added
`days/day-09/examples/main.go` and `days/day-09/examples/main_test.go` are the standalone runnable demonstration and tests. No shared application files are added; later application tests belong under `project/carelog/`.

### Files Changed
No existing shared application files are changed. Future tests and implementation are added under `project/carelog/`; exact paths are deferred.

### Go Concepts Used
The `testing` package, `_test.go`, `TestXxx`, table-driven cases, subtests, `t.Errorf`, `t.Fatalf`, interfaces, and fake dependencies.

### Tests Added
Table-driven validation tests plus fake-writer tests for one successful save and wrapped storage-error propagation.

### Architecture Impact
Validates a narrow application boundary independently of storage and identifies where integration tests are still needed to verify actual database and wiring behavior.

## Exercises
1. Add table cases for a whitespace-only label and a negative value.
2. Add a fake writer that returns a configured error; test that the caller propagates it.
3. Add a case proving a valid observation is saved exactly once.
4. Explain which additional behavior needs an integration test instead of this unit test, and what external resource it would require.

## Quiz
1. What filename suffix makes Go treat a file as test code?
2. What is the difference between `t.Errorf` and `t.Fatalf`?
3. When is a fake preferable to a live database?

## Test your understanding
Temporarily break one validation rule. Confirm that a focused test fails and identifies the case, then restore the rule and rerun all tests.

## Common mistakes
- Writing a test that only repeats the implementation's control flow.
- Checking only that an error is non-nil when its category matters.
- Sharing mutable fake state across unrelated tests.
- Making a fake so elaborate that it becomes a second implementation.

## Production/architecture perspective
Unit tests provide fast feedback, not proof of system behavior. Keep the seam narrow; reserve integration tests for serialization, actual database semantics, and wiring that a fake cannot validate. Avoid storing unnecessary health information in test output.

## Summary/outcome
You can write table-driven tests and use a small fake to test application behavior independently. Reference material: [instructor coding solutions](../../instructor-guide/day-09/solutions.md).
