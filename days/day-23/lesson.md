# Day 23 — Logs, Request Identity, Metrics, and Tracing

**Time budget: 55 minutes** — concepts (15), run and inspect (10), exercises (20), review (10).

**Toolchain target:** Go 1.27.3; examples intentionally use APIs available in Go 1.26 where feasible.

## Learning Objectives

- Emit structured events with `log/slog` and distinguish operational logs from metrics.
- Carry a request ID through middleware and downstream work without trusting arbitrary input blindly.
- Describe what metrics and distributed traces answer, and how cardinality affects cost.
- Implement a minimal readiness/liveness endpoint without overstating its guarantees.

## Why This Matters

When a CareLog request fails in production, an operator needs to find related events and tell whether a problem affects one request or the whole service. Unstructured strings are difficult to query. A request identifier, bounded metrics, and trace context provide complementary views of the same execution.

## Concepts and Explanation

Structured logging stores named attributes (`slog.String`, `slog.Int`) rather than encoding everything into a sentence. Log at meaningful boundaries; avoid logging secrets, health details, or full patient payloads. An error should be logged once at the boundary that has useful context, not at every layer.

A request ID is a correlation label, not authentication and not a trace. Middleware may accept a validated client-provided ID or generate its own; always prevent control characters and unreasonable length, and propagate the chosen value explicitly. `context.Context` can carry request-scoped metadata, but it should not become an optional-parameter bag.

Metrics aggregate numeric measurements over time: request counts, error counts, and duration histograms. Keep label values bounded (route template, status class), never patient IDs or request IDs. Tracing records spans and causal relationships across components; propagating a trace requires compatible instrumentation and context propagation. A log field called `trace_id` alone does not implement tracing.

Health endpoints need precise semantics. Liveness asks whether the process should be restarted; readiness asks whether it should receive traffic. A `200` from `/healthz` does not prove dependencies, data correctness, or clinical safety.

## Runnable Example

Run `go run ./days/day-23/examples`. The example uses `httptest` to call a small handler, adds a generated request ID, writes one structured log event, counts requests, and exposes `/healthz` and `/metrics`. Its counter is deliberately tiny and is not Prometheus exposition or full telemetry instrumentation.

## Project Change

### What We Are Adding
Request-ID middleware, structured request events, a bounded request counter, and distinct liveness/readiness endpoints.

### Why We Need It
Operators need to correlate one request and distinguish a process-health signal from traffic readiness without putting sensitive identifiers in metrics.

### Files Added
`project/carelog/internal/httpapi/observability.go` and `project/carelog/internal/httpapi/observability_test.go`.

### Files Changed
`project/carelog/internal/httpapi/httpapi.go` should wrap routes with the middleware; `project/carelog/internal/server/server.go` should expose the documented health handlers.

### Go Concepts Used
`log/slog`, middleware composition, typed context keys, atomics, HTTP handlers, and bounded metric dimensions.

### Tests Added
Request-ID validation/response tests and health-route status tests; verify logs omit patient payloads and metrics use fixed route/status labels.

### Architecture Impact
Logging, metrics, and health are distinct signals. The examples are illustrative only; they do not supply a telemetry exporter, complete health evaluation, or distributed tracing.

## Coding Exercises

1. **Correlate a request.** Extend the example to accept an incoming request ID only if it matches a strict length/character policy; otherwise generate one. Return it in the response.
2. **Bound metric labels.** Record status class and route template; explain why a patient ID is an unsafe label.
3. **Design a trace.** Draw spans for API → application service → external practitioner directory, indicating which context and trace ID are propagated. State what instrumentation is still missing.

## Quiz

1. Which is better for searching one request: an unbounded request-ID metric label or a structured log field?
2. Does a request ID automatically create a distributed trace?
3. How does readiness differ from liveness?
4. Why should a health handler avoid returning detailed database errors to unauthenticated callers?

## Test Your Understanding

For one failed observation submission, name the log fields, metric dimensions, and trace spans that help an operator without disclosing patient data.

## Common Mistakes

- Logging complete request bodies or sensitive headers.
- Putting unique IDs in metric labels and causing unbounded time-series cardinality.
- Treating a correlation ID as authorization.
- Logging an error at every layer and producing duplicate noise.
- Calling a superficial process response a comprehensive health check.

## Production and Architecture Perspective

The sample is an in-process teaching sketch. Production observability requires an agreed telemetry schema, export pipeline, sampling/retention policy, redaction rules, alert thresholds, and ownership. Health checks should not create a retry storm by synchronously probing every dependency on every orchestrator poll.

## Day Summary

Logs explain individual events, metrics show aggregate behavior, traces connect distributed execution, and health endpoints communicate carefully defined service states.

## What You Should Be Able to Do Now

Add correlation and minimal observability without leaking data or introducing high-cardinality metrics.

**Instructor material:** [full exercise solutions and quiz key](../../instructor-guide/day-23/solutions.md).
