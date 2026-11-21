# Day 28 — Modular Monoliths and Service Boundaries

**Time budget: 55 minutes** — concepts (15), run and inspect (10), exercises (20), review (10).

**Toolchain target:** Go 1.27.3; examples intentionally use APIs available in Go 1.26 where feasible.

## Learning Objectives

- Define a modular monolith and distinguish module boundaries from deployment boundaries.
- Choose dependency direction around CareLog domain behavior.
- Use an interface at a real boundary without creating an interface for every type.
- Write a concise ADR that records context, options, decision, and consequences.

## Why This Matters

CareLog has accumulated concurrency, jobs, observability, and security needs, but those needs do not automatically justify multiple services. A modular monolith can preserve simple deployment while keeping patient and observation responsibilities understandable and testable. It also makes later extraction possible if a genuine independent scaling or ownership need appears.

## Concepts and Explanation

A **modular monolith** is one deployable application with deliberate internal boundaries. A module owns its domain rules and data access; other modules interact through a small application-facing API rather than reaching into its tables or mutable state. A package boundary can make this visible to the compiler, but package count alone does not create modularity.

Keep dependencies pointing toward stable domain/application rules. An HTTP handler translates transport into an application call; it should not own domain decisions. Infrastructure implements a capability needed by the application. In Go, a small interface is often defined by the consumer and satisfied implicitly. Introduce it when it isolates an unstable/external dependency or enables a useful test seam—not for every concrete type.

Clean/hexagonal architecture vocabulary describes dependency direction and ports/adapters; it is not a requirement for layers of boilerplate. DDD concepts such as bounded context help name ownership; they do not require elaborate aggregates for a small learning project. A modular monolith is usually a better default than splitting a small team and data model across network services.

An ADR captures a decision at the time it is made. Keep it short: context, problem, options, decision, consequences. Record the rejected alternatives and why; revisit when assumptions change rather than pretending the choice is permanent.

## Runnable Example

Run `go run ./days/day-28/examples`. `ObservationModule` uses a narrow `PatientReader` capability; the in-memory implementation is composed in `main`. The single-file sample demonstrates dependency direction, not a complete package layout.

## Project Change

### What We Are Adding
A module ownership map, dependency rules, and a concise ADR for keeping CareLog modular-monolithic or extracting a service.

### Why We Need It
Concurrency, jobs, and observability introduce responsibilities, but separate deployment adds network and operations costs. Clear ownership is useful even in one binary.

### Files Added
`project/carelog/architecture/ADR-006-service-boundaries.md` and `project/carelog/architecture/module-map.md`.

### Files Changed
If module access violates the chosen boundary, adjust only the consuming package's capability or composition point, such as `project/carelog/internal/patientclient/client.go`; do not move files solely to imitate a diagram.

### Go Concepts Used
Packages, exported APIs, consumer-defined interfaces, composition, and dependency direction.

### Tests Added
No behavior change is required. Add package-level contract tests only if a boundary change changes observable behavior; otherwise run the existing suite.

### Architecture Impact
Retain one deployable process with explicit module ownership. Clean/hexagonal/DDD language is used to reason about dependencies, not to mandate layers or premature abstractions.

## Coding Exercises

1. **Boundary map.** Draw CareLog modules, owned data, public operations, and allowed dependencies. Identify one dependency that currently crosses an ownership boundary.
2. **Interface judgment.** For `PatientReader`, decide whether an interface is justified. Name a concrete alternative and the testing/runtime trade-off.
3. **ADR.** Write a concise ADR deciding whether to keep CareLog modular-monolithic or extract an observation service. Include context, options, decision, and consequences.

## Quiz

1. Does a package boundary imply a separately deployed service?
2. Where should a small interface usually be defined in Go?
3. Name one signal that could justify service extraction and one that does not.
4. What information makes an ADR useful later?

## Test Your Understanding

If the observations package imports the HTTP server package to write responses, which way does the dependency point, and why is that a boundary smell?

## Common Mistakes

- Creating interfaces before there is an alternative implementation or useful seam.
- Sharing internal mutable maps or database tables across modules.
- Conflating deployment units with Go packages.
- Applying “clean architecture” as a fixed folder template rather than dependency reasoning.
- Extracting services to solve code organization or imagined scale alone.
- Writing ADRs that record the conclusion but omit the context and consequences.

## Production and Architecture Perspective

Service extraction adds network failure, retries, deployment, observability, data ownership, and operational on-call boundaries. A well-modularized monolith can scale vertically/horizontally as one service and can later extract a capability with a stable contract. The right boundary is discovered through changing requirements and team ownership, not by diagram fashion.

## Day Summary

Modularity is about ownership and permitted dependencies; microservices add runtime and organizational boundaries that must earn their cost.

## What You Should Be Able to Do Now

Describe CareLog module ownership, justify an interface, and record an architecture decision with consequences.

**Instructor material:** [full exercise solutions and quiz key](../../instructor-guide/day-28/solutions.md).
