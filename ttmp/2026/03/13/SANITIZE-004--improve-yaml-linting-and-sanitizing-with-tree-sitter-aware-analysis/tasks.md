# Tasks

## Completed during analysis

- [x] Create `SANITIZE-004` and seed the design, diary, and reference docs.
- [x] Add `sanitize parse` so the CLI can print tree-sitter output and structural errors directly.
- [x] Add `scripts/parse_lint_matrix.go` and generate `sources/01-example-corpus-parse-vs-lint-matrix.md`.
- [x] Write the extended analysis and the intern-oriented implementation guide.
- [x] Upload the final document bundle to reMarkable.

## Phase 1: Shared analysis core

- [ ] Add an internal `analyzeDocument(src string)` path in `pkg/yaml` that parses once and returns plain Go data needed by parse, lint, and fix routines.
- [ ] Move duplicate-key collection into that shared analysis pass so `findDuplicateKeys` no longer reparses the document.
- [ ] Introduce a reusable line index helper so byte offsets can be translated to row and column information without repeated full-string scans.
- [ ] Keep `ParseTree` as a thin wrapper over the shared analysis object rather than as a separate parser entrypoint.

## Phase 2: Richer diagnostic model

- [ ] Expand `LintIssue` to carry start and end byte offsets, start and end row and column, and a `Source` field such as `parse`, `heuristic`, or `tree-query`.
- [ ] Add parser-derived lint issues for structural breakage so parse-only failures show up in lint output and UI output.
- [ ] Distinguish broad structural issues from narrowly-localized heuristics in descriptions and JSON output.
- [ ] Decide whether `LintIssue.Row` should be removed entirely or retained as a compatibility convenience derived from `StartRow`.

## Phase 3: Tree-sitter-aware linting

- [ ] Keep heuristic rules for parse-clean style issues such as missing spaces, trailing commas, and duplicate keys.
- [ ] Make structural or ambiguity-prone rules tree-aware by consulting parse-error spans, neighboring rows, and mapping boundaries before emitting issues.
- [ ] Add a first-class mixed-indentation lint issue so parse-only indentation failures are visible before running fixers.
- [ ] Revisit `extra_colon_in_value` so it can use parse context instead of a line-only colon heuristic.

## Phase 4: Tree-sitter-aware fixing

- [ ] Change `applyFixes` to consume the shared analysis result rather than separate `errors` and `lintIssues` slices.
- [ ] Replace the current `errorRows[e.StartRow] = true` shortcut with span-aware or neighborhood-aware targeting.
- [ ] Use structural analysis to scope duplicate-key rewrites and indentation normalization more safely.
- [ ] Preserve the iterative sanitize loop, but ensure each iteration performs one structural analysis instead of repeated unrelated parses.

## Phase 5: Surface area updates

- [ ] Update CLI `lint --json`, `fix --json`, and `parse --json` output examples to show richer issue spans.
- [ ] Update the HTTP `/api/parse` and `/api/sanitize` responses if the richer issue model is adopted there too.
- [ ] Decide whether to expose a public `Analyze` API or keep the shared analysis object internal to `pkg/yaml`.

## Phase 6: Validation

- [ ] Add golden tests or table-driven tests for parse-only, heuristic-only, and hybrid failure fixtures from `examples/yaml/`.
- [ ] Add regression coverage for parse-location-sensitive cases such as tab indentation, mixed indentation, and unresolved parse errors.
- [ ] Keep the matrix script in sync with the corpus and use it to spot regressions in parse-vs-lint coverage.
- [ ] Re-run `go test ./cmd/sanitize ./internal/... ./pkg/yaml` after each refactor phase.
