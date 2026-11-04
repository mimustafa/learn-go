# Day 04 — Functions, multiple returns, variadics, closures, and defer

**Version target:** Go 1.27.3; examples use features available in Go 1.26 where possible.

## Objectives
- Design small functions with explicit inputs and outputs.
- Use multiple returns for a value plus a success/error signal.
- Apply variadic parameters, closures, and `defer` without obscuring behavior.
- **Time budget: 55 minutes** — read 12, code/run 13, exercises 20, quiz/review 10.

## Why
CareLog will need to parse values, summarize readings, and reliably finish small operations. Functions make those responsibilities testable without introducing a framework.

## Concepts and explanation
- A function signature states parameter and result types. Multiple results are ordinary Go; `value, ok` and `value, err` are common patterns.
- A variadic parameter (`...T`) is received as a slice and must be last. To pass an existing slice as variadic arguments, use `slice...`.
- A closure is a function value that can refer to variables from its enclosing scope. Captured state is shared, so lifetime and mutation deserve attention.
- `defer` schedules a call for function return; deferred calls run last-in-first-out. Arguments are evaluated when the defer statement executes. Use it for cleanup, not routine control flow.
- Named results are legal but not automatically clearer. Prefer explicit returns unless naming materially improves the function.

## Runnable example
Run `go run ./days/day-04/examples`. It parses a decimal with an error result, summarizes variadic numbers, creates a labeler closure, and demonstrates deferred cleanup output.

## Project Change

### What We Are Adding
A function-design checkpoint: parsing and summarization are focused operations; CLI output remains at the boundary. The lesson extends those functions with blank-input handling and an explicit empty-summary contract.

### Why We Need It
Small functions make expected behavior, failures, and edge cases easier to reason about and test as CareLog grows.

### Files Added
`days/day-04/examples/main.go` demonstrates function features. No shared application files are added; the later application implementation belongs under `project/carelog/`.

### Files Changed
No existing shared application files are changed. Future CareLog functions will be introduced under `project/carelog/`; concrete paths are deferred.

### Go Concepts Used
Function signatures, multiple results, variadic parameters, closures, captured values, and `defer`.

### Tests Added
None in the standalone example. Try valid, malformed, and blank parser inputs and zero-value aggregation; automated tests are developed in a later checkpoint.

### Architecture Impact
Encourages pure computation at the core and side effects at the CLI edge, while making empty-input behavior part of the function contract.

## Exercises
1. Make `parseValue` reject a blank string and return a useful error.
2. Extend the summary function to return both count and average; define its result for zero values.
3. Create a closure that prefixes all labels with one encounter identifier.
4. Add two deferred calls and predict their output order. Explain when defer arguments are evaluated.

## Quiz
1. Where must a variadic parameter appear?
2. In what order do deferred calls execute?
3. When are a deferred call's arguments evaluated?

## Test your understanding
Call the parser with valid, invalid, and blank inputs. Check that callers can distinguish all outcomes without relying on printed output.

## Common mistakes
- Ignoring the error result from a function.
- Placing a variadic parameter before a regular parameter.
- Expecting deferred calls to run in the order written.
- Capturing and mutating state in a closure without making that state explicit to callers.

## Production/architecture perspective
Pure functions are easy to reason about and test. Functions that print, read global state, or depend on the clock are harder to reuse. `defer` is especially valuable for resource release, where every early return must still clean up.

## Summary/outcome
You can model parsing and aggregation with composable functions and predictable cleanup. Reference material: [instructor coding solutions](../../instructor-guide/day-04/solutions.md).
