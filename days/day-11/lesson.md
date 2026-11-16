# Day 11 — Persisting CareLog observations as JSON

**Time: 55 minutes** — 5 min recall · 10 min concepts · 15 min walkthrough · 15 min exercises · 10 min review

**Version target:** Go 1.27.3; examples use features available in Go 1.26 where possible.

## Objectives
- Explain the boundary between an in-memory model and a durable representation.
- Encode and decode fictional CareLog observations with `encoding/json`.
- Read and write files safely, including handling missing, malformed, and partial data.
- Keep persistence behind a small interface so HTTP or CLI code need not know file details.

## Why this matters
An in-memory CareLog forgets everything when it exits. JSON is useful for a first durable store of fictional patients and observations: it is inspectable, portable, and supported by Go's standard library. It is not a substitute for concurrent, transactional database storage.

## Concepts and explanation
Use a persistence DTO with explicit JSON field names rather than serializing internal implementation details accidentally. Model the CareLog observation fields (`id`, `patient_id`, `code`, `value`, `unit`, `observed_at`, `created_at`). `json.MarshalIndent` produces readable output; `json.Decoder` parses input. Treat file I/O and decoding as fallible operations and return errors to the caller. A missing file can mean “no observations” only when that is an intentional first-run policy; malformed data should not silently become an empty log. Write to a temporary file in the same directory, flush/close it, then rename it to reduce the chance of leaving a truncated file after a crash. This example keeps the mechanics small and does not claim crash-proof durability.

## Runnable example
Run `go run ./days/day-11/examples` from the repository root. The program writes and reads `carelog-demo.json` in the current directory and removes it on exit.

## Project Change

### What We Are Adding
A local JSON observation store with explicit first-run behavior and sibling-temp-file replacement.

### Why We Need It
CareLog must retain fictional observations after process exit. JSON is inspectable and useful for a local single-process checkpoint, but not safe for multiple writers or service replicas.

### Files Added
- `project/carelog/internal/store/json.go` — JSON store implementation.
- `project/carelog/internal/store/json_test.go` — persistence boundary tests.

### Files Changed
- `project/carelog/internal/domain/domain.go` — only if the existing observation DTO/domain conversion needs an explicit store representation.
- `project/carelog/cmd/carelog/main.go` — wire the JSON store into the local command.

### Go Concepts Used
`encoding/json`, `os` file APIs, wrapped errors, DTO conversion, interfaces, `filepath`, temporary files, and `os.Rename`.

### Tests Added
Round trip, legacy JSON with an absent optional `source`, missing-file policy, malformed JSON, and a simulated save/rename failure; use `t.TempDir()` and injected paths.

### Architecture Impact
The command depends on a store boundary rather than a path/format. Persistence remains single-process and local; PostgreSQL replaces this adapter in the later storage checkpoint.

## Exercises
1. Extend the example observation with an optional `Source` field and prove older JSON without that field still loads.
2. Change `Load` to distinguish a missing file from malformed JSON, then add a test for each case.
3. Implement a save strategy that writes a sibling temporary file and renames it; clean up the temp file on every failure path.

## Quiz
1. What does `omitempty` change, and does it make a field private?
2. Why should a malformed existing file not automatically be treated as an empty log?
3. Why is a temporary file created beside the destination before rename?

## Understanding check
Can you explain which layer owns JSON tags, what error reaches the user when decoding fails, and why two simultaneous observation writers are not safe with this design?

## Common mistakes
- Ignoring `Marshal`, `Decode`, `Write`, `Close`, or `Rename` errors.
- Replacing corrupt data with an empty collection and thereby losing evidence.
- Assuming rename alone coordinates multiple processes or provides database transactions.
- Writing internal fields to disk without considering compatibility or privacy.

## Production and architecture perspective
A file store is appropriate for a local prototype or export, not a multi-replica service. Production persistence needs explicit schema evolution, backups, access controls, atomicity and concurrency policy. Hide storage behind an interface so a later PostgreSQL implementation can replace it without rewriting business rules.

## Summary and outcome
You now have a clear JSON persistence boundary, explicit failure behavior, and test cases for the store. **Outcome:** CareLog can retain fictional observations across a restart in a single-process development setup.

## Instructor answers
See [the Day 11 instructor guide](../../instructor-guide/day-11/solutions.md); do not consult it until you have attempted the exercises and quiz.
