# CareLog Milestones

The eight checkpoints are cumulative, and each milestone points to the runnable reference implementation in [`../carelog`](../carelog/). The day lessons introduce the work incrementally; the final module is intentionally not copied eight times, which would create stale snapshots and obscure the one-project progression.

For each checkpoint, meet the requirements, run the listed verification, and review the architecture questions before moving on. Checkpoint numbers map to the corresponding topics in [COURSE.md](../../COURSE.md).

| Milestone | Target days | Working state |
|---|---:|---|
| [01 — CLI and domain seed](./milestone-01.md) | 1–5 | Runnable CLI validates fictional patient/observation input |
| [02 — Packages and tests](./milestone-02.md) | 6–10 | Domain constructors and focused automated tests |
| [03 — HTTP API](./milestone-03.md) | 11–14 | Versioned handlers, validation, lifecycle and health |
| [04 — PostgreSQL](./milestone-04.md) | 15–16 | Database/sql stores, parameterized SQL and migrations |
| [05 — Authentication](./milestone-05.md) | 17 | Password-hash login and signed, expiring bearer tokens |
| [06 — Concurrent processing](./milestone-06.md) | 18–22 | Cancellable bounded processing exercises and failure analysis |
| [07 — Production signals](./milestone-07.md) | 23–27 | Logs, tests, race/benchmark/profile workflow and resilience |
| [08 — Service boundary](./milestone-08.md) | 28–30 | Two independently runnable APIs with explicit network dependency |

## What “working checkpoint” means

The canonical final implementation remains buildable throughout. At a checkpoint, learners run the relevant subset of the final system and its tests, then compare their exercise implementation with the separately stored instructor reference. This gives every checkpoint a runnable reference while avoiding eight diverging copies of the same source tree.
