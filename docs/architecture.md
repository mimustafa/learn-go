# CareLog Architecture Overview

CareLog is a teaching system, not a clinical product. It uses fictional patient data and keeps the model deliberately small: patients and observations.

## Evolution

The first lessons use a single command and simple values. Types and domain rules arrive before interfaces. Tests motivate a small persistence seam. HTTP, PostgreSQL and authentication appear only after the learner can read and test the core Go code. Concurrency is introduced when background observation processing creates a real need to coordinate work.

The expected final shape is a small Patient API and Observation service. Each owns its data and exposes HTTP endpoints. They share authentication conventions, but a network call is not a function call: timeout, cancellation, retries, identity propagation and partial failure are part of the design. The course compares this split with a modular monolith and does not add a message broker, Kubernetes, an ORM or a large framework.

## Boundaries

- **Domain:** identifiers and validation rules that do not depend on HTTP or SQL.
- **Application:** use cases and orchestration.
- **Adapters:** HTTP handlers, PostgreSQL, and calls to the other service.
- **Commands:** configuration and process lifecycle.

These are teaching boundaries, not a mandate to create a package for every noun. Keep the first versions simpler and refactor only when a concrete problem appears.

## Operational constraints

Use explicit timeouts, bounded concurrency, graceful shutdown, structured logs, request IDs and health endpoints. Do not log credentials, bearer tokens, patient details, or raw request paths. Run tests with `-race`; use benchmarks and `pprof` before optimization. Database integration tests require a disposable PostgreSQL database and should never point at a real healthcare system.

## Limitations

CareLog does not implement a healthcare standard, clinical workflow, regulatory controls, production-grade identity management, or real-world patient-data handling. It is not safe for production or clinical use.
