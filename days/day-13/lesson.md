# Day 13 — Validation, errors, middleware, and pagination

**Time: 55 minutes** — 5 min recall · 10 min concepts · 15 min walkthrough · 15 min exercises · 10 min review

**Version target:** Go 1.27.3; examples use features available in Go 1.26 where possible.

## Objectives
- Validate syntactic and domain constraints before changing CareLog state.
- Return stable status codes and a consistent JSON error envelope.
- Wrap handlers with middleware while preserving response behavior.
- Bound collection reads with limit/offset pagination.

## Why this matters
Clients need predictable feedback, and a server must not accept invalid or unbounded work. Validation, error mapping, and pagination make the CareLog API usable as the number of records grows.

## Concepts and explanation
Decode errors concern representation; validation errors concern meaning (for example a blank observation code or non-finite value). Keep those checks explicit and map known failures to 400, 422, 404, or 500 without leaking internal details. A small error DTO can contain a stable code and message. Middleware is a handler wrapper; composition order matters. Pagination parameters are untrusted: parse them, choose defaults, clamp maximum page size, and handle out-of-range offsets. Never build SQL by interpolating untrusted values; pagination values should eventually be bound parameters.

## Runnable example
Run `go run ./days/day-13/examples` from the repository root. It validates an observation code, caps `limit`, emits JSON, and demonstrates a request-counting middleware.

## Project Change

### What We Are Adding
Observation request validation, a stable JSON error envelope, request-ID middleware, and bounded observation-list pagination.

### Why We Need It
Clients need distinguishable errors; limits bound CPU, memory, and database work regardless of caller input.

### Files Added
- `project/carelog/internal/httpapi/middleware.go` — request ID/logging wrapper.
- `project/carelog/internal/httpapi/observation_test.go` — validation and pagination contract tests.

### Files Changed
- `project/carelog/internal/httpapi/observation.go` — validation, status mapping, and list limit.
- `project/carelog/internal/httpapi/httpapi.go` — shared JSON/error response helpers if required.

### Go Concepts Used
Wrapped errors, `errors.Is/As`, `strconv`, UTF-8/length policy, middleware composition, response headers, and safe slice bounds.

### Tests Added
Blank/overlong code, malformed JSON, valid input, invalid/maximum limit, offset boundaries in the in-memory exercise, stable error codes, and request-ID propagation.

### Architecture Impact
Transport owns HTTP status/schema mapping; domain errors remain independent. Middleware adds observability, not identity or authorization. Database-backed pagination uses a stable order and bound values.

## Exercises
1. Reject a whitespace-only observation code and a code longer than the domain's 64-byte limit with distinct stable error codes.
2. Add `offset` parsing to the in-memory example, rejecting invalid/negative input and testing offset equal to and beyond collection length. For a live database API, choose stable ordering and document whether offset pagination is supported.
3. Add middleware that assigns a request ID and includes it in both a response header and structured server log.

## Quiz
1. What is the distinction between decoding failure and domain validation failure?
2. Why cap `limit` instead of trusting the client?
3. Why should a 500 response avoid including the raw database error?

## Understanding check
Can you demonstrate valid observation input, malformed JSON, semantically invalid input, and an unavailable internal dependency as distinct outcomes?

## Common mistakes
- Treating client validation as a replacement for server validation.
- Returning 200 for every outcome or leaking stack/database details.
- Allowing negative or arbitrarily large pagination values.
- Writing middleware that consumes the body or response unexpectedly.
- Treating this sequential in-memory demo as safe for concurrent production requests.

## Production and architecture perspective
Use a stable error schema and observability fields, not unstable error strings as machine-readable codes. Centralize mapping at the transport edge while preserving wrapped causes for logs. Pagination policy is part of API design; cursor pagination may be preferable when records change frequently.

## Summary and outcome
CareLog rejects invalid work predictably and bounds list responses. **Outcome:** clients receive actionable errors without relying on implementation internals.

## Instructor answers
See [the Day 13 instructor guide](../../instructor-guide/day-13/solutions.md) after attempting the exercises.
