# Day 03 — Strings, bytes, runes, slices, and maps

**Version target:** Go 1.27.3; examples use features available in Go 1.26 where possible.

## Objectives
- Distinguish UTF-8 bytes from Unicode code points (`rune`).
- Use slices for ordered values and maps for keyed lookup.
- Avoid assuming a string index addresses a whole character.
- **Time budget: 55 minutes** — read 12, code/run 13, exercises 20, quiz/review 10.

## Why
Observation names can contain accents and non-English text. CareLog must handle text without confusing encoded bytes with user-perceived characters, while collections organize readings and counts.

## Concepts and explanation
- A Go string is immutable bytes, conventionally UTF-8. `len(s)` counts bytes; `range s` decodes UTF-8 and yields byte offsets and `rune` values. A rune is a Unicode code point, not necessarily a displayed grapheme.
- `[]byte(s)` exposes encoded bytes; `[]rune(s)` gives code points. Convert only when the operation needs that representation.
- A slice is a descriptor over an array: length and capacity differ, and `append` may allocate a new backing array. Do not rely on mutation through aliases unless intended.
- A map associates comparable keys with values. A missing key returns the value type's zero value; use `value, ok := m[key]` when absence matters. Map iteration order is unspecified.
- Array and slice equality is not built in; compare elements or use an appropriate library when needed.

## Runnable example
Run `go run ./days/day-03/examples`. Observe the byte length, rune count, slice values, and map lookup for a name containing an accented character.

## Project Change

### What We Are Adding
A CareLog data-representation checkpoint: observation labels are text, ordered readings belong in slices, and keyed counts can use maps. Future identity and normalization rules remain an explicit design decision.

### Why We Need It
Health observation labels can contain non-ASCII text, and collection choices affect ordering, lookup, and behavior around missing values.

### Files Added
`days/day-03/examples/main.go` demonstrates UTF-8 inspection, slices, and maps. No shared application files are added; later implementation stays under `project/carelog/`.

### Files Changed
No existing shared application files are changed. The future application must decide text normalization and identity within `project/carelog/`; exact paths are deferred.

### Go Concepts Used
Strings, UTF-8 bytes, `range` and runes, `[]byte`, `[]rune`, slices, `append`, map lookup with `ok`, and unspecified map iteration order.

### Tests Added
None in the standalone example. Run with non-ASCII labels and inspect offsets and map presence; later tests will codify chosen domain text rules.

### Architecture Impact
Avoids treating byte representation as a complete text policy and keeps data-structure choices aligned to access needs rather than assuming one collection fits all.

## Exercises
1. Replace the label with a string containing a non-Latin character and compare `len` with the number of runes.
2. Add another numeric reading to the slice and calculate a total using `range`.
3. Add a second observation count to the map; use the two-result lookup form to distinguish absent from present-with-zero.
4. Write a helper that counts runes using `range`, then explain why it still may not count displayed symbols as a person perceives them.

## Quiz
1. What does `len` count for a string?
2. What values does `range` over a string produce?
3. Is map iteration order guaranteed?

## Test your understanding
Print the byte offsets from `range` for a multibyte string. Explain why the offsets can skip numbers, then inspect the map lookup's `ok` value for a missing key.

## Common mistakes
- Treating a string index as a rune index.
- Assuming `len([]rune(s))` counts grapheme clusters.
- Reading a missing map key without checking `ok` when absence matters.
- Depending on map iteration order for stable output.

## Production/architecture perspective
Unicode normalization, case folding, and display segmentation are separate text concerns. Do not invent identity rules for clinical labels at the CLI boundary. Keep ordered presentation in a slice; use a map when keyed lookup is the actual requirement.

## Summary/outcome
You can choose appropriate representations for text and collections and state their limits. Reference material: [instructor coding solutions](../../instructor-guide/day-03/solutions.md).
