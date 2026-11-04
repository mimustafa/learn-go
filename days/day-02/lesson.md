# Day 02 — Values, types, and control flow

**Version target:** Go 1.27.3; examples use features available in Go 1.26 where possible.

## Objectives
- Use declarations, zero values, explicit conversions, and typed constants.
- Select `if`, `switch`, and `for` to express validation and classification.
- Explain why Go does not silently convert numeric types.
- **Time budget: 55 minutes** — read 12, code/run 13, exercises 20, quiz/review 10.

## Why
An observation's number is meaningful only with a defined scale and range. Go's explicit types make assumptions visible before CareLog accepts a measurement.

## Concepts and explanation
- Every variable has a type and zero value (`0`, `false`, `""`, or `nil` for reference-like values). `:=` infers a local type; a conversion such as `float64(count)` is explicit and may lose information.
- Constants are compile-time values. Untyped constants can adapt to a compatible context; typed constants cannot be used as another type without conversion.
- `if` may include a short-lived initializer. `switch` selects one matching case and does not fall through by default. `for` is Go's loop; `range` iterates a collection.
- Comparisons and boolean operators are explicit. Prefer a simple branch over a clever expression.
- A valid Go numeric value is not automatically a clinically valid observation. This example only teaches programming mechanics.

## Runnable example
Run `go run ./days/day-02/examples`. It checks a sample value, classifies it against a deliberately fictional threshold, and counts values with a loop.

## Project Change

### What We Are Adding
A validation and classification checkpoint for CareLog: parse and validate an observation value before accepting it. The example's threshold is intentionally fictional and is not medical advice.

### Why We Need It
Explicit value types and branches make assumptions visible and prevent malformed values from flowing unnoticed into later application behavior.

### Files Added
`days/day-02/examples/main.go` demonstrates typed values and control flow. No shared application files are added; the future implementation belongs under `project/carelog/`.

### Files Changed
No existing shared application files are changed. The validation checkpoint will be implemented later under `project/carelog/`; concrete paths are deferred.

### Go Concepts Used
Zero values, variables, constants, explicit numeric types, conversions, `if`, `switch`, `for`, `range`, and `continue`.

### Tests Added
None in the standalone example. Run the provided sample values and exercise negative/threshold cases manually; automated domain tests arrive later.

### Architecture Impact
Separates input validity from any clinical interpretation. Real thresholds require appropriate domain authority and must not be inferred from this teaching example.

## Exercises
1. Change the sample values and predict the classification before running.
2. Add a check that rejects a negative measurement and reports a clear message.
3. Replace the fixed sample with a loop over three values; count how many pass the example threshold.
4. Add a second classification branch and explain why it is a software demonstration, not a medical guideline.

## Quiz
1. What is the zero value of a `string` and of a `bool`?
2. Does Go implicitly convert an `int` to a `float64`?
3. Does a Go `switch` fall through automatically?

## Test your understanding
For each value used by your revised example, write down its type, zero value, and branch taken. Confirm by running the program.

## Common mistakes
- Assuming a declaration has no value until assigned.
- Comparing values of different numeric types without conversion.
- Using a threshold from a teaching example as a clinical rule.
- Writing a loop that never changes its condition or index.

## Production/architecture perspective
Parsing, range validation, and clinical interpretation are distinct responsibilities. In a real system, permissible ranges and their provenance require domain review; a language-level type alone cannot establish clinical validity.

## Summary/outcome
You can use Go values and control flow to express explicit validation and classification. Reference material: [instructor coding solutions](../../instructor-guide/day-02/solutions.md).
