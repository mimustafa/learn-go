# Day 15 — PostgreSQL, database/sql, and migrations

**Time: 55 minutes** — 5 min recall · 10 min concepts · 15 min walkthrough · 15 min exercises · 10 min review

**Version target:** Go 1.27.3; examples use features available in Go 1.26 where possible.

## Objectives
- Explain how `database/sql` separates application code from a driver.
- Configure PostgreSQL through pgx's `database/sql` compatibility driver.
- Use parameterized SQL for values and understand identifier limitations.
- Plan repeatable, ordered schema migrations.

## Why this matters
JSON files cannot safely coordinate multiple CareLog service instances. PostgreSQL adds durable concurrent storage and constraints, while `database/sql` provides a standard API for connections, queries, and transactions.

## Concepts and explanation
The standard library has no PostgreSQL driver. The project adds `github.com/jackc/pgx/v5/stdlib`; its `stdlib.OpenDB` returns a standard `*sql.DB` pool. `sql.Open("pgx", dsn)` is another driver-registered setup, but this course's reference uses `pgx.ParseConfig` with `stdlib.OpenDB`. Opening a handle is lazy: call `PingContext` during startup and close the pool at shutdown. PostgreSQL uses `$1`, `$2` placeholders; values must be arguments, never formatted into the SQL string. Placeholders cannot parameterize table or column names—choose identifiers from a fixed allowlist. Production migrations are ordered, versioned changes; the small CareLog reference embeds idempotent schema SQL and applies it at startup, which is a teaching simplification rather than a migration-history system.

The CareLog reference adapter in `project/carelog/internal/store/postgres.go` uses this shape (the runnable example below avoids a driver dependency):

```go
import (
    "context"
    "database/sql"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/stdlib"
)

func openDB(ctx context.Context, dsn string) (*sql.DB, error) {
    config, err := pgx.ParseConfig(dsn)
    if err != nil { return nil, err }
    db := stdlib.OpenDB(*config)
    db.SetMaxOpenConns(12)
    db.SetMaxIdleConns(4)
    if err := db.PingContext(ctx); err != nil {
        db.Close()
        return nil, err
    }
    return db, nil
}
```

## Runnable example
Run `go run ./days/day-15/examples` from the repository root. It prints a parameterized query and values without connecting to a database. This intentionally uses only the standard library; the actual project adds pgx and executes it through `database/sql`.

## Project Change

### What We Are Adding
A PostgreSQL `database/sql` adapter using pgx stdlib, with parameterized patient/observation queries and startup schema application.

### Why We Need It
The local JSON store cannot coordinate multiple service instances or enforce relational constraints; PostgreSQL supplies shared durable storage.

### Files Added
- `project/carelog/internal/store/migrations/patients.sql` — patient schema.
- `project/carelog/internal/store/migrations/observations.sql` — observation schema.
- `project/carelog/internal/store/postgres_test.go` — repository integration tests when a disposable database is configured.

### Files Changed
- `project/carelog/internal/store/postgres.go` — pgx stdlib connection, ping, parameterized operations, and migration execution.
- `project/carelog/internal/store/store.go` — persistence interfaces/error contracts.
- `project/carelog/go.mod` — pgx dependency.

### Go Concepts Used
`database/sql`, `context`, driver registration, connection pools, `QueryRowContext`, `Scan`, parameter placeholders, `embed`, and wrapped errors.

### Tests Added
Repository mapping and `sql.ErrNoRows`; PostgreSQL integration checks for constraints, inserts, reads, and query parameter behavior against disposable data only.

### Architecture Impact
Handlers/services depend on storage contracts, not SQL. PostgreSQL is the shared source of truth; database-specific behavior remains in the adapter. Embedded idempotent setup is a teaching simplification, not a migration-history tool.

## Exercises
1. Write parameterized `INSERT`, `SELECT`, and `UPDATE` statements for a CareLog observation; pass all user values as arguments.
2. Design a migration that adds a nullable `source` column to observations, then a later migration that safely makes it required with a default.
3. Implement a small `database/sql` repository method that uses `QueryRowContext`, scans a row, and maps `sql.ErrNoRows` to a domain not-found error.

## Quiz
1. Why does `sql.Open` not prove PostgreSQL is reachable?
2. Why are values bound as arguments instead of concatenated into SQL?
3. Can `$1` replace a table name?

## Understanding check
Can you explain where the driver is registered, when the pool opens connections, and how the migration version prevents a schema from being applied twice?

## Common mistakes
- Treating `sql.DB` as one connection rather than a concurrency-safe pool.
- Ignoring `PingContext`, row scan, iteration, or close errors.
- Concatenating client input into SQL or believing value placeholders escape identifiers.
- Editing an already-applied migration instead of adding a new one.

## Production and architecture perspective
Use least-privilege database credentials, TLS where appropriate, secret management, backups, migration review, and bounded connection settings. Migrations are operational changes; assess lock duration and compatibility for rolling deployments. Standard SQL APIs do not eliminate PostgreSQL-specific behavior.

## Summary and outcome
CareLog has a safe route from an application repository to PostgreSQL, with explicit driver, query, and migration boundaries. **Outcome:** the project can add a relational store without exposing SQL to HTTP handlers.

## Instructor answers
See [the Day 15 instructor guide](../../instructor-guide/day-15/solutions.md) after attempting the exercises.
