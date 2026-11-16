# Day 14 — Context, configuration, timeouts, and lifecycle

**Time: 55 minutes** — 5 min recall · 10 min concepts · 15 min walkthrough · 15 min exercises · 10 min review

**Version target:** Go 1.27.3; examples use features available in Go 1.26 where possible.

## Objectives
- Propagate request contexts and deadlines to cancellable work.
- Load and validate configuration at process startup.
- Configure HTTP timeouts, readiness/health, and API version routes.
- Describe graceful shutdown as stop-accepting, drain, then exit.

## Why this matters
CareLog requests must not hang indefinitely, and deployments need a controlled way to start, probe, and stop the service. Context connects client cancellation to downstream operations; configuration avoids hard-coded environment-specific behavior.

## Concepts and explanation
Use `context.Context` as the first parameter for cancellable operations; derive deadlines and always call the cancel function. Do not store request contexts in structs or use them as general-purpose optional parameters. Parse environment variables once into a typed configuration and validate required values early. `http.Server` supports read-header, read, write, and idle timeouts. Liveness asks whether the process should be restarted; readiness asks whether it can serve. Graceful shutdown uses a signal-derived context and `Server.Shutdown`, with a finite deadline. Version the API in its route contract rather than exposing internal build details.

## Runnable example
Run `go run ./days/day-14/examples` from the repository root. It demonstrates a context deadline, a versioned handler, and separate liveness/readiness routes using `httptest`.

## Project Change

### What We Are Adding
Validated startup configuration, HTTP timeouts, liveness/readiness routes, API v1 route registration, and graceful signal-driven shutdown.

### Why We Need It
Explicit deadlines bound resource use, probes distinguish process health from service readiness, and shutdown should drain active requests before exit.

### Files Added
- `project/carelog/internal/config/config_test.go` — configuration validation cases.
- `project/carelog/internal/server/server_test.go` — bounded shutdown/lifecycle tests where practical.

### Files Changed
- `project/carelog/internal/config/config.go` — validate ports, URLs, and timeout settings at startup.
- `project/carelog/internal/server/server.go` — server timeouts and `Shutdown`.
- `project/carelog/internal/httpapi/httpapi.go` — `/healthz`, `/readyz`, and `/v1/...` routes.
- `project/carelog/cmd/patient-api/main.go` and `project/carelog/cmd/observation-api/main.go` — signal and resource lifecycle wiring.

### Go Concepts Used
`context`, deadlines/cancellation, typed configuration, `http.Server`, `signal.NotifyContext`, health handlers, and deferred cleanup.

### Tests Added
Config defaults and invalid values, probe status, request cancellation reaching the service, and shutdown treating `http.ErrServerClosed` as normal.

### Architecture Impact
Commands own process lifecycle; request contexts flow into services and stores. `/readyz` checks the service's own database readiness; it does not claim every downstream dependency is healthy. No orchestration/deployment manifests are added.

## Exercises
1. Add a `loadConfig` function that validates a port and shutdown timeout; test defaults and invalid values.
2. Add a handler that waits on `r.Context().Done()` and verify a cancelled request is abandoned without mutating an observation.
3. Outline or implement a server startup/shutdown path with `signal.NotifyContext` and bounded `Shutdown`.

## Quiz
1. Why should a downstream database call receive the request context?
2. What is the difference between readiness and liveness?
3. Why is a shutdown timeout finite?

## Understanding check
Can you trace cancellation from an HTTP client through the handler to an application service, and describe what happens if graceful shutdown exceeds its deadline?

## Common mistakes
- Dropping the request context and creating unrelated background work.
- Forgetting `cancel`, or using `context.Background()` for request-scoped work.
- Treating every dependency outage as a liveness failure, causing restart loops.
- Leaving server timeouts at zero on an internet-facing server.

## Production and architecture perspective
Timeouts are resource budgets and should be coordinated across proxy, server, and downstream clients. Health endpoints must reveal operational readiness without leaking secrets. Graceful shutdown is best effort; cancellation and idempotent work still matter.

## Summary and outcome
CareLog now has bounded request lifetimes and an explicit process lifecycle contract. **Outcome:** it can start predictably, report health, and stop without abruptly dropping every active request.

## Instructor answers
See [the Day 14 instructor guide](../../instructor-guide/day-14/solutions.md) after attempting the exercises.
