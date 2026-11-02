# Glossary

| Term | Meaning in this course |
|---|---|
| Context | A value passed through a call chain to carry cancellation, deadlines and request-scoped metadata. |
| Domain | The application concepts and rules, such as a Patient or a validated Observation. |
| Error wrapping | Adding context to an error while preserving the original cause for `errors.Is` and `errors.As`. |
| Escape analysis | Compiler analysis that decides whether a value may outlive its stack frame and therefore needs another allocation strategy. |
| Fan-in / fan-out | Distributing work across concurrent producers or workers and then combining results. |
| Fuzzing | Supplying generated inputs to a function to search for crashes or violated invariants. |
| Idempotency | Repeating an operation with the same identity does not apply its effect more than once. |
| Interface | A type set of methods; a value implements an interface implicitly by having those methods. |
| Modular monolith | One deployable application with deliberate internal module boundaries. |
| Observation | A fictional measured value associated with a patient and a time. |
| Repository | A boundary for loading and storing domain data; use one when it clarifies a real dependency. |
| Race | Unsynchronized concurrent access to the same memory where at least one access writes. |
| Service boundary | An independently owned component's API; network calls introduce latency and partial failure. |
| Zero value | The default value produced when a variable is declared without an explicit initializer. |

