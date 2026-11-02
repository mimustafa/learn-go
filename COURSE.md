# 30-Day Course Guide

The target is approximately 50–60 minutes per day: 35–40 minutes for concepts and examples, then 15–20 minutes for a coding exercise and quiz. If an exercise takes longer, finish it before moving on rather than rushing the concepts.

The entire course grows **CareLog**, a fictional patient-and-observation backend. Use only invented data. CareLog is a learning project and is not suitable for real clinical use.

## Daily route

Open the [linked lesson index](./days/README.md) to navigate directly to any day.

| Day | Learn | CareLog checkpoint |
|---|---|---|
| 01 | Toolchain, modules, `go run`, `go build`, formatting and vet | First command-line program |
| 02 | Values, zero values, types, conversion and control flow | Validate an observation |
| 03 | Strings, runes, bytes, arrays, slices and maps | Work with measurement data |
| 04 | Functions, multiple returns, variadics, closures and `defer` | Extract testable operations |
| 05 | Structs, methods, pointers, custom types and zero-value design | Define Patient and Observation |
| 06 | Composition, interfaces, embedding and type assertions | Add only a useful storage boundary |
| 07 | Errors, wrapping, `errors.Is/As`, panic and recovery | Establish validation and error boundaries |
| 08 | Packages, visibility, modules and dependency direction | Separate domain from command entrypoints |
| 09 | Unit tests, table tests, fakes and coverage | Protect domain rules |
| 10 | Generics, constraints and deliberate non-use | Evaluate a reusable helper |
| 11 | JSON and file-backed persistence | Complete a simple runnable workflow |
| 12 | `net/http`, `ServeMux`, handlers and JSON DTOs | Expose the first API |
| 13 | Validation, error responses, middleware and pagination | Make the API predictable |
| 14 | Context, deadlines, config, health, API versioning and shutdown | Bound request and process lifetime |
| 15 | PostgreSQL, SQL, `database/sql`, `pgx`, migrations and parameters | Persist domain data |
| 16 | Transactions, pools, isolation, repository choices and DB tests | Make multi-step writes reliable |
| 17 | Authentication, authorization, password hashes, tokens and sensitive data | Protect API operations |
| 18 | Goroutines, parallelism, races and race detector | Find unsafe concurrent behavior |
| 19 | Channels, `select`, ownership, cancellation and fan-in/out | Coordinate cancellable processing |
| 20 | Worker pools, pipelines, backpressure and bounded concurrency | Process observation imports safely |
| 21 | Mutexes, RWMutex, WaitGroups, atomics, leaks and deadlocks | Choose coordination from ownership needs |
| 22 | Background work, idempotency, retries and shutdown | Handle duplicate and interrupted work |
| 23 | `log/slog`, request IDs, health, metrics and tracing concepts | Make requests diagnosable |
| 24 | Integration, concurrent and fuzz testing | Test real boundaries |
| 25 | Benchmarks, allocations, escape analysis, memory and GC | Measure before changing code |
| 26 | CPU/heap profiling with `pprof` | Find and verify a real bottleneck |
| 27 | Timeouts, bounded retries, rate limits, TLS, secrets and dependency checks | Review reliability and security |
| 28 | Modular monolith, dependency direction, ADRs and service boundaries | Explain why the architecture is still small |
| 29 | Two-service communication, timeouts, failure, idempotency; REST vs gRPC | Extract one justified boundary |
| 30 | Independent build, tests, race checks, profiling and architecture review | Demonstrate a production-oriented system |

## Project route

The project and its runnable source live in [`project/carelog`](./project/carelog/). The checkpoints are cumulative, not unrelated sample applications:

1. CLI and basic domain workflow
2. Packaged, tested domain behavior
3. HTTP API and request validation
4. PostgreSQL persistence and transaction handling
5. Authentication and authorization
6. Cancellable, bounded background processing
7. Structured observability, reliability and measured performance
8. A small Patient API and Observation service, with clear ownership and explicit network failure behavior

The final split is intentionally late. The course compares it with a modular monolith and explains the operational costs of distributed systems. PostgreSQL is the only relational database; SQL is explicit and uses `database/sql` with the `pgx` driver, not an ORM.

## Course conventions

- The student lesson never embeds its answer key. Open `instructor-guide/day-NN/` only after attempting the work.
- Code examples use invented data and avoid real patient identifiers.
- Every daily example is runnable from the repository root with `go run ./days/day-NN/examples`.
- The CareLog Go module and exact setup are documented in [`project/carelog/README.md`](./project/carelog/README.md).
- Use `go test`, `go test -race`, `go test -bench`, `go test -fuzz`, and `go vet` when the day's topic calls for them.
- Infrastructure is deliberately limited. There is no CI/CD, certification, or deployment pipeline curriculum.
