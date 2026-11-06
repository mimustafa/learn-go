# Day 05 — Structs, methods, pointers, and domain types

**Version target:** Go 1.27.3; examples use features available in Go 1.26 where possible.

## Objectives
- Model an observation as a struct with a clear zero-value story.
- Define methods and choose value versus pointer receivers.
- Use named types to make domain intent visible.
- **Time budget: 55 minutes** — read 12, code/run 13, exercises 20, quiz/review 10.

## Why
CareLog's observation is more than unrelated strings and numbers. A domain type gives its fields a shared meaning and a place for behavior and validation.

## Concepts and explanation
- A struct groups named fields. Struct values are copied on assignment and argument passing; nested slices, maps, and pointers still refer to shared data.
- A defined type such as `type ObservationID string` is distinct from its underlying type. A type alias (`type X = Y`) is the same type under another name.
- Methods have receivers. Use a value receiver for small values whose method does not mutate; use a pointer receiver when mutation is intended or copying is undesirable. Keep a type's receiver choice consistent when practical.
- A nil pointer has no value to dereference. Some methods can deliberately handle a nil receiver, but ordinary domain methods should establish and document their invariants.
- Constructors are not mandatory. A useful constructor validates when it is the reliable way to keep invalid values out; exported mutable fields can bypass constructor checks.

## Runnable example
Run `go run ./days/day-05/examples`. It creates an observation, validates required fields, and uses a pointer method to attach a unit.

## Project Change

### What We Are Adding
A domain-model checkpoint: represent an observation with an identifier, label, numeric value, unit, and recorded time, then define required-field validation and display behavior.

### Why We Need It
Grouping related fields gives CareLog a meaningful value with explicit invariants instead of passing unrelated primitive values between operations.

### Files Added
`days/day-05/examples/main.go` demonstrates a standalone observation model. No shared application files are added; the application model will be implemented later under `project/carelog/`.

### Files Changed
No existing shared application files are changed. The later application checkpoint will decide concrete fields and paths under `project/carelog/`.

### Go Concepts Used
Structs, defined types and aliases, methods, value and pointer receivers, nil pointers, and value-copy semantics.

### Tests Added
None in the standalone example. Exercise invalid required fields and copies manually; automated validation tests are added in the later testing checkpoint.

### Architecture Impact
Moves behavior and validation closer to domain data while making clear that a struct or constructor alone does not guarantee validity when fields remain directly writable.

## Exercises
1. Add an `ObservationID` defined type and use it in the struct.
2. Add a method that returns a display string and a validation method that reports missing label/unit.
3. Make an invalid observation and verify validation catches it before display.
4. Decide whether an additional `NormalizeLabel` method should mutate or return a new value; defend the receiver choice.

## Quiz
1. What is copied when a struct value is assigned?
2. What is the practical difference between a defined type and a type alias?
3. When is a pointer receiver appropriate?

## Test your understanding
Copy an observation value, change a scalar field through the copy, and show that the original remains unchanged. Then explain how a slice field would differ.

## Common mistakes
- Assuming a pointer receiver changes a separately copied value.
- Using a pointer for every method without considering value semantics.
- Believing a constructor protects invariants when callers can directly construct invalid exported fields.
- Using a type alias when a distinct domain type was intended.

## Production/architecture perspective
Domain types communicate intent and centralize behavior, but a struct is not automatically a validated entity. Decide which states are allowed and where those rules live. Avoid using types to imply clinical suitability without domain authority.

## Summary/outcome
You can shape an observation as a meaningful Go value and choose receiver semantics deliberately. Reference material: [instructor coding solutions](../../instructor-guide/day-05/solutions.md).
