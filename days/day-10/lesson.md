# Day 10 — Generics, constraints, and restraint

**Version target:** Go 1.27.3; examples use features available in Go 1.26 where possible.

## Objectives
- Read and write a generic function with a type constraint.
- Explain `comparable` and the `~` underlying-type constraint form.
- Decide when generics clarify reusable algorithms and when concrete types are clearer.
- **Time budget: 55 minutes** — read 13, code/run 12, exercises 20, quiz/review 10.

## Why
CareLog uses concrete domain types, but small utilities may operate on several numeric or comparable types. Generics offer reuse without erasing the meaning of domain-specific operations.

## Concepts and explanation
- A type parameter is introduced in brackets, e.g. `func Sum[T Number](values []T) T`. A constraint restricts valid type arguments.
- A union such as `~int | ~float64` permits types whose underlying type is `int` or `float64`, including named types. Constraints define available operations; the compiler checks them.
- `comparable` permits equality and map-key use. It does not mean values are semantically interchangeable.
- A generic type can parameterize a reusable data structure, but a generic wrapper around one concrete domain rule often adds indirection without value.
- Keep CareLog's observation validation concrete: thresholds, units, and meaning are domain-specific. Use generics for a genuinely reusable operation such as summing supported numeric values.

## Runnable example
Run `go run ./days/day-10/examples`. It sums integers and a named numeric type, then uses a generic set constrained by `comparable`.

## Project Change

### What We Are Adding
A CareLog abstraction checkpoint: keep observation validation concrete, and use generic helpers only where an algorithm is genuinely shared across types.

### Why We Need It
Generics can remove meaningful duplication at compile time, but an unnecessary abstraction makes domain behavior harder to understand without adding capability.

### Files Added
`days/day-10/examples/main.go` demonstrates a generic numeric sum and comparable set. No shared application files are added; any later application utility belongs under `project/carelog/`.

### Files Changed
No existing shared application files are changed. This checkpoint records a design choice for later work under `project/carelog/`; it does not prescribe additional project paths.

### Go Concepts Used
Type parameters, type constraints, unions, `~` underlying-type terms, `comparable`, generic functions, and generic types.

### Tests Added
None in the standalone example. Run sums with built-in and named numeric types and set membership; later tests will cover application behavior.

### Architecture Impact
Keeps domain-specific meaning in concrete types and rules while reserving generics for reusable operations whose callers and constraints are clear.

## Exercises
1. Add `float64` to the numeric constraint and sum a slice of decimals.
2. Add a named type with underlying `int` and verify it is accepted by the generic sum.
3. Add `Contains` to the generic set, using its `comparable` key.
4. Write a concrete observation-specific formatter next to a hypothetical generic formatter. Explain why one or both should remain generic.

## Quiz
1. What does `~int` allow in a constraint?
2. What operation does `comparable` permit that is needed by map keys?
3. Should a generic validator replace CareLog's domain-specific observation rules?

## Test your understanding
Try to add a struct type to the numeric constraint. Read the compile error and explain which required operation or type-set membership is missing.

## Common mistakes
- Making a helper generic when there is only one meaningful type.
- Assuming a constraint is a runtime validation.
- Using `comparable` as a claim that two values have equivalent domain meaning.
- Forgetting that floating-point summation has rounding behavior.

## Production/architecture perspective
Generics are a compile-time abstraction, not a substitute for domain modeling or interfaces at runtime boundaries. Evaluate readability, call-site clarity, and actual duplication. Measure performance rather than assuming generic code is automatically faster.

## Summary/outcome
You can use constraints for reusable algorithms and explain why CareLog's core rules should often stay concrete. Reference material: [instructor coding solutions](../../instructor-guide/day-10/solutions.md).
