# Day 12 — HTTP handlers, ServeMux, and JSON DTOs

**Time: 55 minutes** — 5 min recall · 10 min concepts · 15 min walkthrough · 15 min exercises · 10 min review

**Version target:** Go 1.27.3; examples use features available in Go 1.26 where possible.

## Objectives
- Register method-aware routes with `net/http`'s `ServeMux`.
- Separate an API JSON DTO from the CareLog domain type.
- Decode one request body and encode a deliberate JSON response.
- Exercise an HTTP handler without opening a network listener.

## Why this matters
An HTTP API gives a mobile app or browser a stable way to create and view fictional CareLog patients. Handlers should translate HTTP input into application operations; they should not become the domain model or storage layer.

## Concepts and explanation
On Go 1.22 and later, patterns such as `"GET /patients"` and `"POST /patients"` include the method. A `ServeMux` routes requests to `http.Handler`s. A DTO defines the wire contract independently from internal fields. Decode a bounded body, reject invalid input, set `Content-Type` before writing, and encode exactly one response. `httptest.NewRecorder` and `httptest.NewRequest` make handler tests deterministic and fast. A handler in this lesson uses in-memory data to isolate HTTP behavior; it is not the final CareLog service.

## Runnable example
Run `go run ./days/day-12/examples` from the repository root. The example registers GET and POST behavior and calls it using `httptest`. Explicit method checks keep the demo compatible with legacy ServeMux behavior; method-qualified patterns are also available in Go 1.22+.

## Project Change

### What We Are Adding
A first HTTP transport for fictional patients: list and create endpoints with independent request and response DTOs.

### Why We Need It
The command-line workflow needs a stable interface for clients while preserving domain rules and storage behind application services.

### Files Added
- `project/carelog/internal/httpapi/patient.go` — patient handlers, DTOs, and route registration.
- `project/carelog/internal/httpapi/patient_test.go` — handler contract tests.

### Files Changed
- `project/carelog/cmd/carelog/main.go` — connect the handler to the application service and start the local server.

### Go Concepts Used
`net/http`, `ServeMux`, `http.Handler`, request/response DTOs, JSON encoding, `httptest`, and method/status constants.

### Tests Added
GET collection, successful POST, malformed JSON without a service mutation, content type, status, response DTO shape, and service error mapping.

### Architecture Impact
HTTP becomes an adapter at the application boundary; domain/storage types are not the public wire contract. PostgreSQL, auth, and production lifecycle remain out of scope today.

## Exercises
1. Add `GET /patients/{id}` using `Request.PathValue`, return one patient, and test a missing ID.
2. Add a request DTO distinct from the response DTO; prove an internal-only field is never serialized.
3. Add a test that submits malformed JSON and verifies the handler does not create a patient.

## Quiz
1. What responsibility belongs in a DTO rather than the domain model?
2. Why does `httptest` provide better handler tests than starting a real listener for every case?
3. When must the response content type be set?

## Understanding check
Can you trace a `POST /patients` from route match through decoding and service call to the response without coupling the handler to a file path?

## Common mistakes
- Using `json.Marshal` as a substitute for an HTTP response and forgetting status or headers.
- Decoding a request directly into a persistence/domain type.
- Assuming Go versions before 1.22 support method-qualified mux patterns.
- Writing a second response after a partial encode failure.
- Treating this sequential in-memory demo as safe for concurrent server requests.

## Production and architecture perspective
The API is a compatibility boundary. DTOs permit internal refactoring without silently changing clients. Set server timeouts and body limits in the server configuration; those concerns are introduced in Day 14. Keep handlers thin and put business rules in an application service.

## Summary and outcome
You can expose the patient collection with method-aware routes and test behavior without networking. **Outcome:** CareLog has a transport boundary that can evolve separately from its domain and store.

## Instructor answers
See [the Day 12 instructor guide](../../instructor-guide/day-12/solutions.md) after attempting the work.
