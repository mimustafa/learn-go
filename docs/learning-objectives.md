# Learning Objectives

By the end of the course, a learner should be able to:

## Go language and tools

- Read and write idiomatic Go using packages, modules, types, slices, maps, structs, methods, pointers, interfaces and appropriate generics.
- Explain zero-value design, explicit error handling, composition, and the trade-offs of abstraction.
- Use `go run`, `go build`, `go test`, `gofmt`, `go vet`, the race detector, benchmarks and fuzzing.

## Application and backend engineering

- Design a small domain model and package structure with clear dependency direction.
- Build and test HTTP APIs with validation, context cancellation, timeouts and graceful shutdown.
- Persist data using PostgreSQL, parameterized SQL, transactions, migrations and Go's `database/sql`.
- Apply authentication, authorization, sensitive-data handling and secure logging fundamentals.

## Concurrency and performance

- Reason about goroutine ownership, channels, cancellation, shared state, races, deadlocks and goroutine leaks.
- Choose bounded concurrency and synchronization based on measured needs.
- Benchmark and profile before optimizing, and interpret allocation and profiling evidence.

## Production and architecture

- Add useful structured logs and health signals, and explain metrics and tracing.
- Design for timeouts, idempotency, bounded retries and dependency failure.
- Explain when a modular monolith is simpler, and when a small service boundary is justified.
- Review a Go service for correctness, testability, security, performance and operational complexity.
- Teach the core concepts by explaining the reasoning and common mistakes, not only the syntax.
