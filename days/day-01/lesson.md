# Day 01 — Start CareLog: Go, modules, and a tiny CLI

**Version target:** Go 1.27.3; examples use features available in Go 1.26 where possible.

## Objectives
- Explain what the Go toolchain and a module do.
- Format, run, build, and inspect a small Go program.
- Trace a CLI from its arguments to its output.
- **Time budget: 55 minutes** — read 12, code/run 15, exercises 18, quiz/review 10.

## Why
CareLog begins as a command-line tool for recording a simple healthcare observation. A small executable keeps attention on Go itself before the course adds domain rules and architecture.

## Concepts and explanation
- `package main` declares an executable package; `func main()` is its entry point. Go source in a package shares one directory.
- A module is the unit of dependency and version management, normally declared by `go.mod`. `go run` compiles and runs, `go build` produces a binary, `go fmt` formats, `go vet` checks suspicious constructs, and `go test` runs tests.
- Go statements usually need no semicolon. Imports must be used. `:=` declares and initializes a local variable; `var` is useful when the type or zero value matters.
- `os.Args` includes the executable name at index 0. Check input before indexing it; a CLI should explain valid usage rather than panic.
- From the repository root, run the standard-library-only example with `go run ./days/day-01/examples`. The same command works for each day's example; no external dependencies are needed.

## Runnable example
Run `go run ./days/day-01/examples`, then try `go run ./days/day-01/examples "blood pressure" 118`.
`main.go` prints usage for missing arguments and accepts an observation name and value.

## Project Change

### What We Are Adding
The first CareLog checkpoint: a command-line boundary that accepts an observation label and value, reports usage when input is incomplete, and keeps parsing separate from future domain validation. This is a project design checkpoint; it does not implement the application.

### Why We Need It
A runnable starting point makes the toolchain concrete and gives later lessons a single evolving healthcare observation CLI to refine.

### Files Added
`days/day-01/examples/main.go` is the standalone lesson demonstration. No shared application files are added; the eventual application work belongs under `project/carelog/`.

### Files Changed
No existing shared application files are changed. The CLI checkpoint is for later implementation under `project/carelog/`; concrete application paths are deferred.

### Go Concepts Used
`package main`, `main`, imports, `os.Args`, argument-count checks, `fmt`, and running/building a package.

### Tests Added
None for the standalone CLI yet. Verify behavior manually with no arguments and with a label/value pair; automated tests are introduced later.

### Architecture Impact
Establishes a clear input boundary without prematurely deciding the future package layout or embedding clinical rules in argument parsing.

## Exercises
1. Run the example with no arguments, then with a name and value. Explain the difference.
2. Add an optional unit (for example `mmHg`) and print it only when supplied.
3. Use `go fmt` on the example and `go vet` on its package. Build a binary into a disposable location of your choice and run it.
4. Change the CLI to reject an empty observation name without indexing a missing argument.

## Quiz
1. What does `package main` mean, and what function starts the program?
2. Why should a CLI check `len(os.Args)` before reading an argument?
3. Which tool formats source, and which runs package tests?

## Test your understanding
Without looking at the source, predict the output for zero, one, and two user arguments. Then run each case and reconcile your prediction.

## Common mistakes
- Treating `os.Args[0]` as the first user argument.
- Assuming `go run` skips compilation.
- Adding unused imports or forgetting to format after editing.
- Treating input text as medically valid simply because it is non-empty.

## Production/architecture perspective
A CLI boundary should report actionable usage errors and keep domain validation separate from argument parsing. This tiny program is intentionally not a patient-record system and makes no clinical decisions.

## Summary/outcome
You can run, format, inspect, and build a minimal Go CLI. CareLog now has a concrete executable starting point. Reference material: [instructor coding solutions](../../instructor-guide/day-01/solutions.md).
