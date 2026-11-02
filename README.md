# Learn Go by Building CareLog

A 30-day, self-paced Go course for experienced developers who are new to Go. The course follows one small healthcare information system from a command-line program to a tested, observable backend with two deliberately small services.

## Start here

1. Read [Course Guide](./COURSE.md) and [the learning objectives](./docs/learning-objectives.md).
2. Install Go 1.27.3 or a later compatible patch release. See [Go installation](https://go.dev/doc/install).
3. Begin with [Day 01](./days/day-01/lesson.md); work one day at a time.
4. Follow the [linked day-by-day lesson index](./days/README.md).
5. Do exercises before opening their answers in [the instructor guide](./instructor-guide/README.md).

The project is in [`project/carelog`](./project/carelog/). Milestone descriptions explain how the same codebase grows across the course. Architecture decisions are recorded in [`architecture`](./architecture/).

## Prerequisites

- Experience programming in another language; no prior Go experience is assumed.
- A terminal and editor.
- PostgreSQL for database lessons (a local installation or the provided Compose setup).
- Docker is optional; the language and standard-library examples do not require it.

## Common commands

Run an individual standard-library example:

```sh
go run ./days/day-01/examples
```

Run the CareLog tests:

```sh
cd project/carelog
go test ./...
```

Format, vet, and race-check the project:

```sh
cd project/carelog
gofmt -w .
go vet ./...
go test -race ./...
```

PostgreSQL-backed integration tests are opt-in and use `CARELOG_TEST_DATABASE_URL`; see [the project guide](./project/carelog/README.md). Never put real patient information or credentials in course exercises, logs, screenshots, or commits. CareLog is educational software, not a clinical system.

## Course map

| Days | Focus |
|---|---|
| 01–05 | Toolchain, values, collections, functions, structs |
| 06–10 | Interfaces, errors, packages, tests, generics |
| 11–14 | JSON, HTTP, API design, context and lifecycle |
| 15–17 | PostgreSQL, transactions, authentication and security |
| 18–22 | Concurrency, cancellation, workers and reliability |
| 23–27 | Observability, integration tests, performance and resilience |
| 28–30 | Architecture, service boundaries and final review |

Every day has a student lesson, a runnable example, exercises and a quiz. Solutions and answer keys are separate from student material.

## Learning workflow

For each day: **Learn → Code → Test → Review → Architect**. Keep your own progress in Git commits if useful; Git workflows are not part of the curriculum. The complete outline is in [COURSE.md](./COURSE.md).
