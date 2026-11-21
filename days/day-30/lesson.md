# Day 30 — Final Production-Oriented Build and Review

**Time budget: 55 minutes** — independent build (30), verification (15), architecture review (10).

**Toolchain target:** Go 1.27.3; examples intentionally use APIs available in Go 1.26 where feasible.

## Learning Objectives

- Integrate Go types, HTTP, validation, concurrency safety, context, tests, and operational concerns in one small service.
- Make explicit decisions about idempotency, error boundaries, readiness, and data sensitivity.
- Review a service for correctness, reliability, security, observability, and operational gaps.
- Distinguish an educational production-oriented design from a production-ready healthcare system.

## Why This Matters

The course has progressed from Go syntax and domain behavior through HTTP, persistence concepts, concurrency, observability, performance, resilience, and boundaries. A production-oriented engineer can now build a coherent small component, explain its guarantees, test those guarantees, and state what remains unimplemented.

## Concepts and Explanation

The final build is intentionally small. A good service slice has a clear request/response contract, boundary validation, domain behavior, safe shared-state handling, bounded work, stable errors, and tests at its seams. Operational details such as config, deadlines, logging, metrics, health, deployment, and recovery are part of the design—not optional polish.

Review the system as a set of promises. Does a successful response mean the write is durable? Can the caller safely retry? Are unauthorized patient records inaccessible? Can a downstream stall indefinitely? Does shutdown drain or abandon work? Can an operator correlate a failure without seeing sensitive data? Is the reported health meaningful?

No sample in this course implements complete healthcare security, privacy, auditability, high availability, clinical safety, or regulatory compliance. A real system requires domain experts, threat modeling, legal/regulatory review, durable storage, identity and authorization, lifecycle controls, and operational validation.

## Runnable Example

Run `go run ./days/day-30/examples`. The sample composes a small standard-library HTTP handler, validates JSON input, protects an in-memory idempotency store, returns a request ID, and exercises the contract via `httptest`. It is a reviewable reference slice, not a deployable complete system.

## Project Change

### What We Are Adding
An independently designed observation API slice with validation, request correlation, idempotency, readiness, concurrency-safe storage, and a review report.

### Why We Need It
The final exercise integrates prior lessons and asks learners to state the exact guarantee boundary instead of equating “runs” with “production-ready.”

### Files Added
For the standalone challenge, add `project/carelog/internal/httpapi/final_observation_test.go` and any focused package tests needed for the chosen implementation.

### Files Changed
Depending on the design, change the corresponding files under `project/carelog/internal/domain/`, `internal/httpapi/`, and `internal/store/`; the reference demonstration remains in `days/day-30/examples/main.go`.

### Go Concepts Used
Go types and errors, HTTP handlers, context, mutexes, idempotency, JSON validation, `httptest`, and concurrent tests.

### Tests Added
At minimum: domain validation, HTTP malformed/valid contract, duplicate replay/conflict, and concurrent idempotent writes; run `go test -race`.

### Architecture Impact
Deliver a reviewed slice with explicit limitations. Do not promote the in-memory sample to real patient use; durable persistence, access control, TLS/secrets, audit, recovery, and relevant domain/compliance reviews remain separate requirements.

**Independent challenge (30 minutes):** without looking at the instructor solution, build a `POST /observations` endpoint with a validated patient ID, code, numeric value, request ID, and idempotency key. Add a readiness route, one unit test, one HTTP contract test, and a concurrent test. Keep data synthetic. State whether persistence is durable, how shutdown behaves, and what the service does not guarantee.

Then review with this checklist:

- **Correctness:** clear invariants; validation and error/status mapping; idempotency behavior.
- **Concurrency:** ownership/locks; cancellation; race test; no unbounded goroutine growth.
- **Reliability:** bounded deadlines; retry rules; shutdown and recovery behavior.
- **Security/privacy:** authentication/authorization boundary; input/resource limits; no sensitive logs; TLS and secrets plan.
- **Observability:** structured contextual logs; bounded metrics; trace plan; health semantics.
- **Performance:** representative benchmark/profile only where a requirement calls for it.
- **Architecture:** module ownership, dependency direction, service boundary rationale, ADR.

## Coding Exercises

1. **Build the service slice.** Implement the independent challenge above. Run formatting, tests, race detection, and vet; write down assumptions and limitations.
2. **Production review.** Review the runnable example against the checklist. Identify at least eight gaps, prioritize them by risk, and propose the smallest next change for the top three.
3. **Architecture decision.** Decide whether the observation service remains inside the modular monolith or is independently deployed. Write a short ADR with options, evidence, and consequences.

## Quiz

1. Does returning `201 Created` from an in-memory store prove durability?
2. Which tests are needed to exercise a handler's HTTP contract and shared-state race safety?
3. Name two security guarantees that this sample does not provide.
4. What evidence should justify extracting CareLog into a separately deployed service?

## Test Your Understanding

Explain the guarantee boundary for each: request accepted, response returned, idempotency record saved, process shutdown, and client retry. What evidence would you need before telling an operator that an observation is durable?

## Common Mistakes

- Calling the sample production-ready because it compiles and has a health route.
- Mixing validation, HTTP status selection, persistence, and logging in one opaque handler.
- Treating an in-memory map as durable or cross-instance state.
- Ignoring request-body limits, authorization, privacy, and shutdown behavior.
- Adding architecture layers or a second service without a demonstrated need.
- Treating passing tests as evidence of requirements the tests never exercise.

## Production and Architecture Perspective

“Production-oriented” means the design considers deployment, security, observability, failure, and operations. It does not mean these educational snippets are certified, complete, or suitable for real patient data. Before real use: select durable storage, implement and independently review identity/authorization, establish privacy/audit policy, validate backups/recovery, test load/failure modes, scan/update dependencies, define SLOs, and obtain the appropriate clinical, security, and compliance reviews.

## Day Summary

The course outcome is not memorized APIs; it is the ability to build a small Go service, test its observable contracts, reason about its boundaries, and state precisely what it guarantees.

## What You Should Be Able to Do Now

Independently implement and review a production-oriented Go backend slice, explain trade-offs, and identify what evidence and controls are still required before real deployment.

**Instructor material:** [full reference solution and quiz key](../../instructor-guide/day-30/solutions.md).
