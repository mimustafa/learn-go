# Day 16 — Transactions, pools, and repository integration

**Time: 55 minutes** — 5 min recall · 10 min concepts · 15 min walkthrough · 15 min exercises · 10 min review

**Version target:** Go 1.27.3; examples use features available in Go 1.26 where possible.

## Objectives
- Use a transaction when multiple database writes must succeed or fail together.
- Set pool limits and understand their effect on concurrency.
- Choose transaction isolation based on invariants, not folklore.
- Keep repository methods context-aware and testable.

## Why this matters
A multi-observation CareLog import should not leave only half its rows committed when one insert fails. Transactions provide an atomic boundary for that batch; pool limits protect the database from an unbounded number of open connections.

## Concepts and explanation
Start with `BeginTx(ctx, opts)`, defer rollback for every early return, perform all related work through `*sql.Tx`, and `Commit` only after all steps succeed. A deferred rollback after commit is harmless. Never mix `db.Exec` with `tx.Exec` inside a transaction by accident. A batch-import repository method follows this shape:

```go
tx, err := db.BeginTx(ctx, nil)
if err != nil { return err }
defer tx.Rollback()
for _, observation := range observations {
    if _, err := tx.ExecContext(ctx, insertObservation,
        observation.ID, observation.PatientID, observation.Code,
        observation.Value, observation.Unit, observation.ObservedAt,
        observation.CreatedAt); err != nil {
        return err
    }
}
return tx.Commit()
```

The existing CareLog store's single-row `CreateObservation` is one statement and does not need an explicit transaction; use a transaction when adding an all-or-nothing batch import or another multi-statement invariant. `SetMaxOpenConns`, `SetMaxIdleConns`, and connection lifetime settings are workload-specific; too few connections queue requests, too many overload the database. Isolation controls visibility and anomalies, but higher isolation can increase contention or require retrying serialization failures. Keep transaction duration short and never wait on user input while holding one.

## Runnable example
Run `go run ./days/day-16/examples` from the repository root. A mutex-protected in-memory observation batch demonstrates all-or-nothing update and rollback behavior; it is a teaching model, not an implementation of PostgreSQL isolation.

## Project Change

### What We Are Adding
An all-or-nothing observation batch import, pool capacity settings, and PostgreSQL integration coverage for commit/rollback behavior.

### Why We Need It
A failed batch must not persist a misleading partial import; pool limits must also protect PostgreSQL under concurrent API requests.

### Files Added
- `project/carelog/internal/store/postgres_test.go` — transaction and pool integration checks using disposable PostgreSQL.
- `project/carelog/internal/application/import_test.go` — batch behavior tests with a fake repository.

### Files Changed
- `project/carelog/internal/store/store.go` — batch-import repository contract.
- `project/carelog/internal/store/postgres.go` — short transaction for multi-row import and pool settings.
- `project/carelog/internal/config/config.go` — validated pool settings if made configurable.
- `project/carelog/internal/application/import.go` — import orchestration.

### Go Concepts Used
`sql.Tx`, deferred rollback, commit errors, isolation options, pool stats, interfaces, contexts, and transaction-scoped SQL.

### Tests Added
Commit all items, rollback on middle insert failure, duplicate-ID constraint, cancellation, and pool configuration. Test database isolation only against PostgreSQL.

### Architecture Impact
Transactions align with a business atomicity boundary (batch import), not every repository call. The single-row observation insert remains one statement; in-process mutexes do not coordinate service replicas.

## Exercises
1. Implement `ImportObservations(ctx, observations)` using `BeginTx`; insert every observation through the transaction and roll back if any insert fails.
2. Set and justify pool limits for a small local service; expose pool stats as metrics without logging credentials.
3. Describe or test a concurrent update invariant for an observation import (e.g. duplicate observation IDs are rejected), selecting a unique constraint and mapping the conflict to a stable domain error.

## Quiz
1. Why defer `Rollback` immediately after a successful `BeginTx`?
2. What can happen when `MaxOpenConns` is far above database capacity?
3. Does serializable isolation mean a transaction can never fail?

## Understanding check
Can you identify all statements that must use the same `*sql.Tx`, all error paths that abort, and the invariant that the chosen isolation level protects?

## Common mistakes
- Beginning a transaction and then issuing part of the work on `db`.
- Forgetting rollback on an intermediate error or ignoring commit failure.
- Holding a transaction across network calls or expensive computation.
- Treating larger pools or stronger isolation as universally better.

## Production and architecture perspective
Repository transactions should reflect business atomicity, not arbitrary controller boundaries. Monitor wait time, open/in-use connections, transaction latency, deadlocks, and retry rates. Retrying is safe only when the operation is designed to be idempotent and the failure is retryable.

## Summary and outcome
Related CareLog writes now share an explicit atomic boundary, while connection concurrency has an intentional budget. **Outcome:** consistency requirements guide repository design and database capacity settings.

## Instructor answers
See [the Day 16 instructor guide](../../instructor-guide/day-16/solutions.md) after attempting the exercises.
