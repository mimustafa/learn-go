# Day 29 — A Small Service Boundary: HTTP, Failure, and Idempotency

**Time budget: 55 minutes** — concepts (15), run the example (10), exercises (20), review (10).

**Toolchain target:** Go 1.27.3; examples intentionally use APIs available in Go 1.26 where feasible.

## Learning Objectives

- Explain what changes when a local function call becomes a network call.
- Set HTTP client timeouts, propagate context, and map downstream failures deliberately.
- Use an idempotency key to make a retry safe after an ambiguous timeout.
- Compare REST/HTTP JSON with gRPC without treating either as universally superior.

## Why This Matters

Day 28 kept CareLog as a modular monolith by default. This day introduces one deliberate boundary: the CareLog API forwards a recording operation to an observation service. The network boundary adds independent failure, latency, versioning, deployment, and data-ownership concerns that an in-process interface does not have.

## Concepts and Explanation

A remote call can fail before sending, during execution, or after the remote side commits but before the response reaches the caller. A timeout is a local observation—not proof that the operation did not happen. The caller must pass its deadline, use a bounded client timeout, close response bodies, and translate downstream errors into stable API errors without exposing internals.

For retryable side effects, propagate a stable idempotency key across the boundary. The receiver stores the key with the result atomically and returns the original result for a duplicate. The educational example uses a process-local map; real services need durable storage and a retention/reuse policy. Keys must be scoped to an authorized caller and operation, and request payload conflicts for the same key need a defined response.

REST typically uses HTTP resources and JSON, is broadly interoperable, and is straightforward to inspect. gRPC uses Protocol Buffers and generated contracts, with strong typing, streaming, and efficient binary encoding, but adds code generation, tooling, and client/runtime considerations. Both require deadlines, auth, versioning, observability, and failure handling. Choose based on team, consumers, and requirements.

Do not split services merely to exercise networking. A boundary is justified by independent ownership, deployment, scaling, or security requirements that exceed the operational cost.

## Runnable Example

Run `go run ./days/day-29/examples`. Two `httptest` servers model a CareLog API and a downstream observation service. The downstream commits a new idempotency key, delays its first response to force a timeout, and returns the prior result on replay. The effect count remains one.

## Project Change

### What We Are Adding
A narrow HTTP contract for recording an observation remotely, with request schema, stable idempotency key, deadlines, and explicit status mapping.

### Why We Need It
A remote boundary can fail after commit but before response, so a client timeout does not determine whether the write occurred. The exercise makes distributed failure semantics visible.

### Files Added
`project/carelog/internal/patientclient/client_test.go` additions for timeout and error behavior, and `project/carelog/architecture/ADR-006-service-boundaries.md` if the architecture decision is newly recorded.

### Files Changed
`project/carelog/internal/patientclient/client.go` and the observation API composition point only if the remote-call contract is adopted for the experiment.

### Go Concepts Used
`http.Client`, request contexts, timeout/error handling, response-body closure, stable idempotency keys, and `httptest`.

### Tests Added
Tests simulate delayed downstream responses, stable-key replay, conflicting payload behavior, and client cancellation. Verify effects rather than assuming a timeout means failure.

### Architecture Impact
This is an experiment, not a required extraction. Keep the modular monolith by default; REST/JSON versus gRPC depends on consumer and operating requirements, not protocol fashion.

## Coding Exercises

1. **Failure matrix.** List behavior for refused connection, timeout before commit, timeout after commit, 4xx, and 5xx. State which cases retry and why.
2. **Idempotency contract.** Specify key scope, payload mismatch behavior, persistence/atomicity, retention, and authorization. Add a test that replays a request and asserts only one effect.
3. **REST or gRPC decision.** Recommend one for CareLog's current internal boundary and list the evidence that would change your choice.

## Quiz

1. What does a client timeout tell the caller about remote commit state?
2. Where must idempotency state be stored to survive service restart?
3. Name one strength and one cost of gRPC compared with REST/JSON.
4. Does using HTTP make an application a microservice architecture?

## Test Your Understanding

The client retries a POST with a new key after a timeout. Why does this defeat deduplication? What should happen if the same key is reused with a different payload?

## Common Mistakes

- Retrying every error, including validation and authorization failures.
- Generating a new idempotency key for each attempt.
- Assuming a timeout cancels a remote side effect.
- Returning raw downstream response bodies or internal addresses to API callers.
- Omitting an overall deadline while adding per-attempt timeouts.
- Choosing gRPC solely because it is “faster” without measuring or considering consumers.

## Production and Architecture Perspective

The two-server example is an in-process simulation, not deployment guidance. A real extraction needs service identity/authentication, TLS/mTLS decisions, durable idempotency, compatibility/versioning, retries/circuit/backpressure policy, telemetry propagation, deployment and ownership, and operational SLOs. Distributed transactions are not implied; design for partial failure and reconciliation.

## Day Summary

Remote calls create ambiguous outcomes. Deadlines bound waiting; stable idempotency makes safe replay possible; protocol choice follows consumer needs and operating costs.

## What You Should Be Able to Do Now

Design a narrow HTTP boundary and explain its failure semantics before deploying a second service.

**Instructor material:** [full exercise solutions and quiz key](../../instructor-guide/day-29/solutions.md).
