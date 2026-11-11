# Day 08 — Packages, exports, modules, and dependency direction

**Version target:** Go 1.27.3; examples use features available in Go 1.26 where possible.

## Objectives
- Explain package names, import paths, and exported identifiers.
- Describe what `go.mod` declares and why dependencies should point inward.
- Sketch a small CareLog package boundary without splitting code prematurely.
- **Time budget: 55 minutes** — read 13, code/run 12, exercises 20, quiz/review 10.

## Why
As CareLog grows, placing every concern in one package obscures ownership. Packages should group code by responsibility and make the direction of dependencies understandable.

## Concepts and explanation
- A package is a set of source files in one directory; files in a package share declarations. Import paths identify packages within a module. The package clause names the package used in source.
- An identifier beginning with an uppercase Unicode letter is exported; lowercase names are package-private. Export only what another package needs.
- `go.mod` declares the module path and Go language version and records dependencies. `go mod tidy` reconciles imports and module requirements. Package names and module import paths are related but not identical concepts.
- A useful dependency direction is CLI/application → domain behavior, with infrastructure depending on a small inward-facing contract. Domain code should not import a database implementation simply to save an observation.
- The examples use only the standard library and run from the repository root with `go run ./days/day-08/examples`. The CareLog application's own module and package structure will be designed under the shared project location `project/carelog/`.

## Runnable example
Run `go run ./days/day-08/examples`. Inspect how `main` uses standard-library packages and a small local function. The example remains one package so this lesson can focus on package boundaries before splitting code.

## Project Change

### What We Are Adding
A package and dependency-direction checkpoint: distinguish the CLI, observation policy, and persistence responsibilities and decide what should be exported across package boundaries.

### Why We Need It
As CareLog grows, cohesive packages and inward-facing dependencies reduce accidental coupling and keep domain rules independent of transport and storage choices.

### Files Added
`days/day-08/examples/main.go` remains one runnable package for this lesson. No shared application files are added; any future module or package files will live under `project/carelog/`.

### Files Changed
No existing shared application files are changed. The eventual application module and package layout are deferred to work under `project/carelog/`; this lesson does not promise particular paths.

### Go Concepts Used
Packages, imports, import paths, exported and unexported identifiers, modules, `go.mod`, dependency direction, and consumer-owned interfaces.

### Tests Added
None in the standalone example. The package-boundary exercise is a design sketch; later package tests will verify public behavior.

### Architecture Impact
Sets a rule for future application evolution: dependencies point toward stable policy, with infrastructure implementing behavior required by inward-facing consumers.

## Exercises
1. Draw a three-box dependency sketch for CLI, domain, and storage. Add arrows and justify their direction.
2. Propose one exported function and one unexported helper for a hypothetical `observation` package; state who needs each.
3. In a separate scratch workspace, initialize a module with `go mod init example.com/carelog` and inspect the generated `go.mod`. The CareLog application module will be designed later under `project/carelog/`.
4. Identify a dependency that would be backwards if a domain package imported it, and explain how a consumer-owned interface could avoid that dependency.

## Quiz
1. How is an identifier exported from its package?
2. What is the purpose of `go.mod`?
3. Should the domain package import a concrete SQL driver to express its storage needs?

## Test your understanding
Explain the difference between a package directory, the package clause, and a module import path. Use the example's standard-library imports to ground your explanation.

## Common mistakes
- Treating every directory name, package name, and module path as interchangeable.
- Exporting every type and helper.
- Creating tiny packages before the boundaries are stable.
- Letting domain rules depend directly on transport or storage packages.

## Production/architecture perspective
Package structure is a tool for expressing ownership and reducing accidental coupling, not a goal in itself. Keep changes cohesive, expose the minimum API, and make dependencies point toward stable policy rather than volatile details.

## Summary/outcome
You can reason about Go package visibility, modules, and dependency direction. Reference material: [instructor coding solutions](../../instructor-guide/day-08/solutions.md).
