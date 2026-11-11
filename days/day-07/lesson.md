# Day 07 — Errors, wrapping, panic, and recovery

**Version target:** Go 1.27.3; examples use features available in Go 1.26 where possible.

## Objectives
- Return and propagate errors as ordinary values.
- Add context with `%w` and inspect causes using `errors.Is` and `errors.As`.
- Distinguish recoverable failures from programmer-fault panics.
- **Time budget: 55 minutes** — read 13, code/run 13, exercises 19, quiz/review 10.

## Why
CareLog must report failed parsing and persistence without losing context. Explicit errors let the caller decide whether to retry, report, or stop.

## Concepts and explanation
- The `error` interface exposes `Error() string`. Return an error when a caller can handle the failure; keep successful values and failures explicit.
- Add operation context with `fmt.Errorf("...: %w", err)`. `errors.Is` searches an error chain for a sentinel or matching error; `errors.As` finds a value of a specified error type. Avoid parsing error strings.
- Sentinel errors can identify stable conditions. Custom error types carry structured context. Choose one based on what callers must distinguish.
- A panic is for unrecoverable programmer/runtime conditions, not ordinary invalid user input or I/O errors. `recover` only works in a deferred function in the same goroutine that panicked.
- A recovery boundary may protect a process from a known isolation boundary, but it should log/report the panic and preserve enough diagnostics; silently recovering hides defects.

## Runnable example
Run `go run ./days/day-07/examples`. It wraps a parsing error, classifies it, demonstrates a typed validation error, and isolates a deliberate panic behind a recovery boundary.

## Project Change

### What We Are Adding
A CareLog failure-handling checkpoint: return ordinary input and operation failures as errors, wrap them with useful context, and classify them at the appropriate boundary.

### Why We Need It
Callers need to distinguish invalid input from infrastructure or programming failures without depending on fragile error-string matching.

### Files Added
`days/day-07/examples/main.go` demonstrates wrapped errors and a deliberately isolated recovery case. No shared application files are added; application error handling will be implemented later under `project/carelog/`.

### Files Changed
No existing shared application files are changed. Future error boundaries belong under `project/carelog/`; concrete source paths are deferred.

### Go Concepts Used
The `error` interface, sentinels, custom error types, `%w`, `errors.Is`, `errors.As`, `panic`, `defer`, and `recover`.

### Tests Added
None in the standalone example. Run the error classifications and recovery demonstration; later tests will cover stable error categories and propagation.

### Architecture Impact
Defines explicit error propagation and handling boundaries. Panic/recover remains exceptional and goroutine-local; it does not replace ordinary validation or I/O error returns.

## Exercises
1. Add a sentinel for missing observation label and verify it with `errors.Is`.
2. Wrap a `ValidationError` with more context and retrieve it using `errors.As`.
3. Make the parsing caller print an action-oriented message while retaining the wrapped cause.
4. Move the deliberate panic into a new goroutine, observe why the current recovery does not catch it, then remove the panic. Do not leave an unrecovered panic in the finished example.

## Quiz
1. What does `%w` enable?
2. When do you use `errors.Is` versus `errors.As`?
3. Can a deferred recover in one goroutine catch a panic in another?

## Test your understanding
For each failure in the example, state whether it is an error return or a panic, who handles it, and whether the original cause remains inspectable.

## Common mistakes
- Comparing error strings instead of using `errors.Is`/`errors.As`.
- Discarding an error or logging it repeatedly at every layer.
- Using panic for bad user input.
- Expecting `recover` to work outside a deferred function or across goroutines.

## Production/architecture perspective
Wrap errors at meaningful boundaries with operation context; log once where the failure is handled. Keep sensitive observation values out of logs unless there is a reviewed, justified policy. Recovery is not a substitute for error handling or supervision.

## Summary/outcome
You can propagate errors with context and reserve panic/recover for exceptional boundaries. Reference material: [instructor coding solutions](../../instructor-guide/day-07/solutions.md).
