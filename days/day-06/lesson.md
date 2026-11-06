# Day 06 — Composition, interfaces, embedding, and assertions

**Version target:** Go 1.27.3; examples use features available in Go 1.26 where possible.

## Objectives
- Compose domain values from smaller structs rather than relying on inheritance.
- Explain implicit interface implementation and keep interfaces focused.
- Use a type assertion safely and justify an interface seam.
- **Time budget: 55 minutes** — read 12, code/run 13, exercises 20, quiz/review 10.

## Why
CareLog may eventually write observations to different destinations. An interface becomes useful when behavior needs substitution, especially at a boundary; it is not a requirement for every domain type.

## Concepts and explanation
- Go favors composition: a type can contain another value. Embedding promotes selected fields/methods for convenient access, but does not create inheritance or override semantics.
- A type satisfies an interface implicitly by implementing its methods. Interfaces describe behavior required by a consumer; keep them small and define them near their use.
- A type assertion `x.(T)` panics if it fails. The comma-ok form `value, ok := x.(T)` is safe; a type switch handles several dynamic types.
- Introduce an interface when a caller genuinely needs a capability seam (for example, persistence that can be replaced by a fake in tests). Avoid an interface with one implementation merely to imitate a pattern.
- Composition and interfaces solve different problems: one builds values from parts; the other names behavior.

## Runnable example
Run `go run ./days/day-06/examples`. The example embeds audit metadata, accepts a small storage interface, and safely checks an optional capability.

## Project Change

### What We Are Adding
A CareLog dependency-seam checkpoint: observation operations can depend on a narrow writer capability, while the implementation can be swapped where a real persistence or testing need exists.

### Why We Need It
Substitution at a genuine boundary enables isolated tests and future storage changes without making the domain depend on infrastructure.

### Files Added
`days/day-06/examples/main.go` demonstrates composition, embedding, and small interfaces. No shared application files are added; the eventual seam belongs under `project/carelog/`.

### Files Changed
No existing shared application files are changed. The application boundary will be implemented later under `project/carelog/`; specific paths are deferred until its structure is established.

### Go Concepts Used
Composition, embedding, implicit interface implementation, consumer-oriented interfaces, type assertions, comma-ok, and type switches.

### Tests Added
None in the standalone example. The example runs against an in-memory implementation; later unit tests will use a focused fake at the selected seam.

### Architecture Impact
Introduces an interface only at a justified boundary and keeps policy independent of storage details; it does not prescribe broad abstractions or inheritance.

## Exercises
1. Add an in-memory store implementing `ObservationWriter` and pass it to the example function.
2. Extend the interface only if your caller needs another behavior; explain why each method belongs there.
3. Use a type switch for two distinct optional capabilities and a default case.
4. Remove the interface and describe which test or replacement need is lost. Decide whether the abstraction still earns its cost.

## Quiz
1. Must a type explicitly declare that it implements an interface?
2. What does embedding do, and what does it not do?
3. Which type assertion form avoids a panic?

## Test your understanding
Write down which package owns the consumer's interface in a layered design and why a storage implementation should depend inward on that contract rather than make domain code import storage.

## Common mistakes
- Creating interfaces before there is a consumer need.
- Treating embedding as subclassing.
- Using a single-result type assertion on untrusted dynamic values.
- Making a broad interface that every fake must implement.

## Production/architecture perspective
An interface seam limits coupling only when it reflects a real boundary. Keep persistence mechanics outside domain rules. A fake should model observable behavior needed by the test, not reproduce a database.

## Summary/outcome
You can combine concrete types and use a narrow behavioral seam without importing inheritance patterns. Reference material: [instructor coding solutions](../../instructor-guide/day-06/solutions.md).
