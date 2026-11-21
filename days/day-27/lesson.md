# Day 27 — Resilience and Security Review

**Time budget: 55 minutes** — review concepts (15), run the example (10), exercises (20), review (10).

**Toolchain target:** Go 1.27.3; examples intentionally use APIs available in Go 1.26 where feasible.

## Learning Objectives

- Apply deadlines and bounded retries to outbound HTTP calls.
- Explain rate limiting and the limitations of an in-process limiter.
- Identify TLS, secret-handling, input-validation, and dependency-review requirements.
- Conduct a focused threat-oriented review of a small CareLog API.

## Why This Matters

A CareLog service depends on networks, configuration, libraries, and user input. A slow downstream can consume all request goroutines; careless retries can amplify an outage; and logging a token or patient identifier can turn diagnostics into a breach. Reliability and security need explicit boundaries.

## Concepts and Explanation

Use an `http.Client` with a timeout and derive request contexts with deadlines. The client timeout bounds the complete exchange; the context can impose a tighter per-operation budget. Always close response bodies. Retry only classified transient failures, with bounded attempts, backoff/jitter, and an overall deadline. Respect `Retry-After` where appropriate. A timeout does not prove a remote operation failed; use idempotency for side effects.

Rate limiting constrains request volume. A process-local limiter can protect one instance but its counters are not shared across replicas and are lost at restart. Choose the key carefully: IP addresses can be shared or spoofed through untrusted forwarding headers; authenticated identities may be more meaningful. Avoid exposing patient identity in limiter keys or logs.

TLS provides confidentiality and server identity when configured and verified correctly. Use a supported TLS minimum (TLS 1.2 or newer in this example), keep certificate verification enabled, and manage certificate rotation. Never “fix” certificate errors with `InsecureSkipVerify`.

Secrets belong in a managed configuration/secret system, not source control, logs, images, or error messages. Validate inputs at boundaries, use parameterized SQL, apply authorization to every protected resource, and minimize stored/sent health data. Review dependencies using the Go vulnerability tooling available in the environment (`govulncheck ./...`) and maintain an update process. `go mod verify` checks downloaded module content; it is not a vulnerability scan or proof of trust.

## Runnable Example

Run `go run ./days/day-27/examples`. The sample demonstrates an outbound timeout, a tiny per-process fixed-window limiter, constant-time comparison for a fixed-length demo token, and a TLS minimum setting. It does not configure a production secret store, distributed limiter, authentication system, or server certificate.

## Project Change

### What We Are Adding
A resilience/security review and bounded deadlines, input limits, error mapping, and retry/idempotency rules at external boundaries.

### Why We Need It
Unbounded calls, unsafe retries, invalid input, and leaked secrets can cause outages or expose sensitive data. The controls must be selected for CareLog's deployment model.

### Files Added
`project/carelog/internal/httpapi/resilience_test.go` to cover bounded behavior and error contracts, plus a focused threat-review checklist in the milestone notes.

### Files Changed
`project/carelog/internal/patientclient/client.go` for call budget propagation and `project/carelog/internal/httpapi/observation.go` for boundary validation/error handling, only where the review finds gaps.

### Go Concepts Used
HTTP client timeouts, context deadlines, bounded retry classification, TLS configuration, constant-time comparison limits, and dependency checks.

### Tests Added
Deadline/cancellation tests, retry classification tests, oversized/invalid input cases, and checks that public responses/logging do not disclose secrets.

### Architecture Impact
Timeouts and retries protect a dependency boundary but do not implement distributed resilience or security alone. In-memory rate limiting is per process; secrets, TLS lifecycle, and dependency review need deployment controls.

## Coding Exercises

1. **Timeout and retry policy.** Specify attempts, deadline, retryable failures, and behavior after timeout for an idempotent GET and a side-effecting POST.
2. **Threat review.** Review an endpoint accepting observation ID, patient ID, and value. List at least five threats and mitigations spanning validation, authentication/authorization, logs, transport, and resource exhaustion.
3. **Dependency review.** Write a repeatable dependency-check procedure for a Go module, including what `go mod verify`, `govulncheck`, and review/update policy each do not prove.

## Quiz

1. Why can retries worsen an overloaded dependency?
2. Does TLS replace authorization?
3. What is lost when a local in-memory rate limiter is used with several replicas?
4. What does `go mod verify` establish—and what does it not establish?

## Test Your Understanding

A POST times out after the downstream commits. Describe the safe retry behavior and the role of an idempotency key.

## Common Mistakes

- Retrying forever, retrying permanent failures, or retrying without an overall deadline.
- Disabling certificate verification to get a connection working.
- Keeping credentials in source or logging authorization headers.
- Assuming one instance's rate limit protects a multi-instance deployment.
- Trusting client-supplied forwarding headers from arbitrary peers.
- Treating a vulnerability scan as proof that a dependency is safe.

## Production and Architecture Perspective

Security is layered and threat-model dependent. A short demo cannot implement whole-service security: production also needs identity, authorization policy, secure configuration, key/certificate lifecycle, audit controls, incident response, and tested recovery. Resilience controls need capacity and failure-mode testing; retries and rate limits can themselves create load or deny legitimate traffic.

## Day Summary

Timeouts bound waiting, retries need idempotency and limits, rate limits need a deployment-aware key, and security requires layered controls and ongoing dependency review.

## What You Should Be Able to Do Now

Review a service boundary for common reliability and security failures and distinguish illustrative code from production controls.

**Instructor material:** [full exercise solutions and quiz key](../../instructor-guide/day-27/solutions.md).
