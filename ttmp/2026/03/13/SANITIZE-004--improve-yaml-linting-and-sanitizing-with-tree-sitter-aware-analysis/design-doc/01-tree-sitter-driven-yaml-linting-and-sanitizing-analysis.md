---
Title: Tree-sitter-driven YAML linting and sanitizing analysis
Ticket: SANITIZE-004
Status: active
Topics:
    - yaml
    - linting
    - treesitter
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: pkg/yaml/analysis.go
      Note: Shared analysis path added during implementation
    - Path: pkg/yaml/duplicate_keys.go
      Note: Tree-based duplicate-key detection that currently reparses
    - Path: pkg/yaml/fix.go
      Note: Fix pipeline and current start-row targeting
    - Path: pkg/yaml/line_index.go
      Note: Reusable byte-to-row and column translation helper
    - Path: pkg/yaml/lint.go
      Note: Current lint rule implementation
    - Path: pkg/yaml/parse.go
      Note: ParseTree and collectErrors implementation
    - Path: pkg/yaml/sanitize.go
      Note: Current orchestration that should move to shared analysis
    - Path: pkg/yaml/types.go
      Note: Diagnostic types and their current limitations
    - Path: ttmp/2026/03/13/SANITIZE-004--improve-yaml-linting-and-sanitizing-with-tree-sitter-aware-analysis/scripts/parse_lint_matrix.go
      Note: Corpus experiment used to compare parser and linter coverage
    - Path: ttmp/2026/03/13/SANITIZE-004--improve-yaml-linting-and-sanitizing-with-tree-sitter-aware-analysis/sources/01-example-corpus-parse-vs-lint-matrix.md
      Note: Generated experiment evidence for the design conclusions
ExternalSources: []
Summary: Evidence-backed design for making YAML linting and sanitizing tree-sitter-aware without discarding necessary heuristics.
LastUpdated: 2026-03-13T08:52:51.114170223-04:00
WhatFor: Document the current parser-lint-fix architecture, explain where tree-sitter already helps, and define the next refactor for shared structural analysis.
WhenToUse: Use when refactoring pkg/yaml diagnostics, designing new lint rules, or reasoning about parse-driven fix behavior.
---



# Tree-sitter-driven YAML linting and sanitizing analysis

## Executive Summary

The current YAML sanitizer is already built on tree-sitter, but it does not yet behave like a tree-sitter-first system. `ParseTree` collects structural `ERROR` and `MISSING` nodes in `pkg/yaml/parse.go:22-35`, while `Lint` in `pkg/yaml/lint.go:19-102` still works almost entirely from line-based regular expressions. `Sanitize` in `pkg/yaml/sanitize.go:15-24` reparses and relints repeatedly, and duplicate-key detection reparses again in `pkg/yaml/duplicate_keys.go:18-49`.

The right direction is not "replace heuristics with tree-sitter." The corpus experiment in `sources/01-example-corpus-parse-vs-lint-matrix.md` shows that the repository has three classes of failures:

1. Parse-only failures, where tree-sitter already sees the breakage but `Lint` reports nothing.
2. Heuristic-only failures, where YAML still parses and tree-sitter alone will never flag the issue.
3. Hybrid failures, where tree-sitter does report an error but its span is too coarse to fully localize the problem.

The proposed design is therefore a blended model:

1. Parse the document once into a shared internal analysis object.
2. Convert structural parse failures into first-class lint issues with full spans.
3. Run heuristic rules on top of that analysis, using parse context to improve scoping and fix targeting.
4. Feed the same analysis object into `ParseTree`, `Lint`, and `Sanitize`.

This preserves the package's current strengths while removing duplicated parser work and making diagnostics richer for both the CLI and the web UI.

## Problem Statement and Scope

The repository presents itself as a tree-sitter-powered sanitizer, but the actual diagnostic flow is split across unrelated passes:

- `ParseTree` parses YAML and collects structural errors in `pkg/yaml/parse.go:22-79`.
- `Lint` scans text lines using regexes in `pkg/yaml/lint.go:9-100`.
- `findDuplicateKeys` builds a second tree-sitter parse in `pkg/yaml/duplicate_keys.go:18-49`.
- `Sanitize` calls `ParseTree` and `Lint` separately every iteration in `pkg/yaml/sanitize.go:15-24`.
- `applyFixes` uses tree-sitter only indirectly, mainly through a boolean `hasTreeErr` flag keyed on the parse error's start row in `pkg/yaml/fix.go:26-39`.

This ticket focuses on improving the YAML path only. It does not implement JSON support. It also does not try to redesign the whole public API surface in one step. The scope is the diagnostic and fixing pipeline inside `pkg/yaml`, plus the CLI and HTTP surfaces that expose those diagnostics.

## Current-State Analysis

### Runtime map

The current runtime flow looks like this:

```text
CLI / server
    |
    v
pkg/yaml.Sanitize()
    |
    +-> ParseTree() ----------------------------+
    |                                           |
    +-> Lint()                                  |
    |    +-> regex line rules                   |
    |    +-> findDuplicateKeys() -> parse again |
    |
    +-> applyFixes(errors, lintIssues, cfg)
    |
    +-> repeat until clean or stuck
```

The key files are:

- `pkg/yaml/parse.go` for tree-sitter parser setup and structural error collection.
- `pkg/yaml/lint.go` for regex-style lint rules.
- `pkg/yaml/duplicate_keys.go` for the one lint rule that already traverses the parse tree.
- `pkg/yaml/fix.go` for line and document-level fix routines.
- `pkg/yaml/sanitize.go` for the orchestration loop.
- `pkg/yaml/types.go` for externally visible diagnostic types.
- `internal/cli/commands.go` for user-facing inspection surfaces, including the new `sanitize parse`.

### What tree-sitter already does well

Tree-sitter already provides several useful capabilities:

- It gives precise byte and row and column spans for `ERROR` and `MISSING` nodes in `pkg/yaml/parse.go:46-68`.
- It supports mapping-aware duplicate-key detection in `pkg/yaml/duplicate_keys.go:52-90`.
- It provides a machine-readable S-expression that is now exposed by `sanitize parse` in `internal/cli/commands.go`.

This means the repository already has a structural-analysis foundation. The missing piece is shared reuse.

### What the data types reveal

`ErrorNode` already contains byte and row and column spans in `pkg/yaml/types.go:13-23`, but `LintIssue` carries only `Rule`, `Description`, and a single zero-based `Row` in `pkg/yaml/types.go:25-31`.

That asymmetry matters because:

- parse diagnostics are already span-aware,
- lint diagnostics are still line-summary objects,
- downstream code cannot treat both as one diagnostic language.

### What the sanitize loop reveals

One sanitize iteration currently performs:

```text
Sanitize()
  -> ParseTree(src)
  -> Lint(src)
     -> regex line checks
     -> findDuplicateKeys(src)
        -> parse again
```

The duplicated parser work is visible in:

- `pkg/yaml/sanitize.go:15-24`
- `pkg/yaml/lint.go:93-100`
- `pkg/yaml/duplicate_keys.go:18-49`

This is not a catastrophic performance issue in a small CLI tool, but it is architectural duplication, and it makes future tree-aware lint rules harder to implement cleanly.

### What the fixer reveals

`applyFixes` converts parse errors into a set of start rows only:

```go
errorRows := map[uint]bool{}
for _, e := range errors {
    errorRows[e.StartRow] = true
}
```

That logic is in `pkg/yaml/fix.go:26-30`. It throws away end-row coverage immediately, and it assumes the parser start row is a useful proxy for the true problem location.

That assumption is only partially true.

## Experimental Findings

The experiment script in `scripts/parse_lint_matrix.go` generated `sources/01-example-corpus-parse-vs-lint-matrix.md`. The important result is not raw counts by themselves. It is the shape of the disagreement between parser output and lint output.

### Category A: parse-only failures

Examples:

- `examples/yaml/20-mixed-indent.yaml`
- `examples/yaml/28-unresolved-parse-error.yaml`

Observed behavior:

- tree-sitter reports structural errors,
- `Lint` emits nothing.

Interpretation:

- structural parse failures need to become first-class lint issues,
- mixed indentation should likely gain an explicit lint rule instead of being fix-only.

### Category B: heuristic-only failures

Examples:

- `examples/yaml/11-missing-space-after-colon.yaml`
- `examples/yaml/13-list-dash-no-space.yaml`
- `examples/yaml/14-trailing-comma-flow-mapping.yaml`
- `examples/yaml/17-duplicate-key-sibling.yaml`

Observed behavior:

- the parser accepts the YAML,
- lint still correctly identifies usability or semantic problems.

Interpretation:

- heuristics remain necessary,
- a tree-sitter-first design must still allow parse-clean style diagnostics.

### Category C: hybrid failures

Examples:

- `examples/yaml/10-tab-indent.yaml`
- `examples/yaml/16-extra-colon-in-value.yaml`
- `examples/yaml/21-multiple-errors.yaml`
- `examples/yaml/24-deeply-nested-mixed-errors.yaml`
- `examples/yaml/27-extra-colon-needs-quotes.yaml`

Observed behavior:

- tree-sitter reports a structural failure,
- heuristics identify additional actionable lines,
- parser locations are sometimes broad or anchored to a parent line.

Interpretation:

- tree-sitter should inform heuristic scoping,
- tree-sitter should not replace heuristics wholesale.

### Representative example: tab indentation

`go run ./cmd/sanitize parse --json examples/yaml/10-tab-indent.yaml` reports one structural error with text `service:` on row 1, while the actual lint issues are on rows 2 and 3. This is the best short proof that parser error locations are useful but insufficient as a direct line-targeting strategy.

### Representative example: mixed indentation

`go run ./cmd/sanitize parse --json examples/yaml/20-mixed-indent.yaml` reports one structural error spanning rows 1 through 3, and `Lint` reports nothing. This shows the opposite pattern: here the parser span is useful, and the linter is the one that is incomplete.

### Representative example: unresolved parse breakage

`examples/yaml/28-unresolved-parse-error.yaml` produces two structural errors and zero lint issues. This shows that parser-derived lint issues are needed even when no repair heuristic exists.

## Gap Analysis

The repository needs four design upgrades.

### Gap 1: no shared analysis object

Everything that follows from the parse tree is recalculated piecemeal. That makes it harder to add new tree-aware rules because there is no stable place to put intermediate structural facts.

### Gap 2: no unified diagnostic shape

`ErrorNode` and `LintIssue` speak different location languages. This forces downstream code to special-case parse errors instead of treating all diagnostics uniformly.

### Gap 3: no structural lint layer

The project has parse output and heuristic lint output, but it lacks a middle layer that says, "these parser failures are user-facing lint issues too."

### Gap 4: fixers cannot ask structural questions

`fixLine` currently gets only a line string, a row number, a list of rule names, and a `hasTreeErr` boolean in `pkg/yaml/fix.go:67-149`. That is too little context for tree-aware fixing.

## Proposed Solution

Introduce an internal shared document-analysis pass and make parse, lint, and fix logic consume it.

### Design goals

1. Parse once per sanitize iteration.
2. Preserve heuristic rules for parse-clean issues.
3. Promote parser failures to first-class diagnostics.
4. Improve fix targeting without pretending tree-sitter is perfect at localizing every YAML problem.
5. Keep the public API simple even if the internal analysis object becomes richer.

### Proposed internal data model

The internal model should return plain Go structs only. Do not expose raw tree-sitter nodes outside the parse pass, because they depend on tree lifetime and are hard to use safely from the rest of the package.

```go
type analysis struct {
    TreeText       string
    ParseErrors    []ErrorNode
    DuplicateKeys  []duplicateKeyOccurrence
    LintIssues     []LintIssue
    LineStarts     []int
    ProblemWindows []lineWindow
}

type lineWindow struct {
    StartRow int
    EndRow   int
    Reason   string
}
```

The exact helper fields can change, but the design point is stable: parse once, capture reusable facts, and return plain data.

### Proposed diagnostic model

Upgrade `LintIssue` from a single-row summary to a real span-based issue:

```go
type LintIssue struct {
    Rule        string `json:"rule"`
    Source      string `json:"source"` // parse, heuristic, tree-query
    Description string `json:"description"`
    StartByte   uint   `json:"start_byte"`
    EndByte     uint   `json:"end_byte"`
    StartRow    uint   `json:"start_row"`
    StartCol    uint   `json:"start_col"`
    EndRow      uint   `json:"end_row"`
    EndCol      uint   `json:"end_col"`
}
```

`Row` can be retained temporarily as a derived convenience if the package wants a softer migration, but there is no technical reason to keep it as the canonical location once a richer span exists.

### Proposed runtime flow

```text
Analyze(src)
    -> parse once with tree-sitter
    -> collect parse errors
    -> collect duplicate-key facts from mapping nodes
    -> build line index and problem windows
    -> emit parser-derived lint issues
    -> emit heuristic lint issues using structural context

ParseTree(src)
    -> Analyze(src)
    -> return TreeText + ParseErrors

Lint(src)
    -> Analyze(src)
    -> return LintIssues

Sanitize(src)
    -> Analyze(original)
    -> loop:
         analysis := Analyze(current)
         if clean -> return
         fixes := applyFixes(current, analysis, cfg)
         if no fixes -> return
```

### Parser-derived lint layer

Every `ErrorNode` should be translated into a user-facing lint issue. Recommended initial rules:

- `structural_parse_error` for `ERROR`
- `missing_syntax_node` for `MISSING`

This immediately fixes the current parse-only blind spot for files like `examples/yaml/20-mixed-indent.yaml` and `examples/yaml/28-unresolved-parse-error.yaml`.

### Heuristic rules that stay line-based

Some rules remain fundamentally heuristic and should continue to run even on parse-clean documents:

- `missing_space_after_colon`
- `list_dash_no_space`
- `trailing_comma`
- `duplicate_key`

The point of the refactor is not to delete these rules. It is to make them consume shared structural context and emit better location data.

### Heuristic rules that should become tree-aware

Two rules are especially good candidates for tree-aware scoping.

#### `extra_colon_in_value`

Current behavior:

- line-only colon matching in `pkg/yaml/lint.go:69-90`,
- line-only quoting in `pkg/yaml/fix.go:135-179`.

Problem:

- the rule has no notion of neighboring structural breakage or mapping boundaries.

Proposed improvement:

- only elevate this rule aggressively inside parse-problem windows or on lines whose structure resembles an unquoted plain scalar value,
- keep the actual fixer conservative.

#### mixed indentation

Current behavior:

- no lint rule,
- fix-only document-level repair in `pkg/yaml/fix.go:217-295`.

Problem:

- users do not see a dedicated lint issue explaining why the file is structurally broken.

Proposed improvement:

- add a `mixed_indent` lint issue that is emitted when indentation widths conflict in a parse-problem window.

### Fixing with structural context

`applyFixes` should accept the full analysis object:

```go
func applyFixes(src string, a analysis, cfg *config) (string, []Fix)
```

That unlocks several improvements:

- use parse windows instead of a single start-row boolean,
- let line fixers ask whether a row is inside a structural problem window,
- let document-level fixers consume duplicate-key facts from the same analysis pass.

### Proposed pseudocode

```go
func analyzeDocument(src string) (analysis, error) {
    tree, root := parseYAML(src)

    a := analysis{
        TreeText:      root.ToSexp(),
        ParseErrors:   collectErrors(root, []byte(src)),
        LineStarts:    buildLineIndex(src),
        DuplicateKeys: collectDuplicateKeys(root, []byte(src)),
    }

    a.ProblemWindows = buildProblemWindows(a.ParseErrors, src)
    a.LintIssues = append(a.LintIssues, lintFromParseErrors(a.ParseErrors)...)
    a.LintIssues = append(a.LintIssues, lintHeuristics(src, a)...)
    a.LintIssues = dedupeAndSortIssues(a.LintIssues)
    return a, nil
}

func lintHeuristics(src string, a analysis) []LintIssue {
    issues := lintLineRules(src, a)
    issues = append(issues, lintDuplicateKeys(a.DuplicateKeys)...)
    issues = append(issues, lintMixedIndentIfNeeded(src, a)...)
    return issues
}

func sanitize(src string, cfg config) Result {
    original := mustAnalyze(src)
    current := src

    for iter := 0; iter < cfg.maxIterations; iter++ {
        a := mustAnalyze(current)
        if len(a.ParseErrors) == 0 && len(a.LintIssues) == 0 {
            return buildResult(original, a, current)
        }

        next, fixes := applyFixes(current, a, &cfg)
        if len(fixes) == 0 {
            return buildResult(original, a, current)
        }
        current = next
    }

    final := mustAnalyze(current)
    return buildResult(original, final, current)
}
```

## Design Decisions

### Decision 1: keep heuristics

Rationale:

- multiple invalid examples parse cleanly enough that tree-sitter will not flag them,
- duplicate keys are semantic and not syntax errors,
- trailing commas and spacing issues are still important to users.

### Decision 2: use tree-sitter as the shared structural layer

Rationale:

- the repository already uses tree-sitter everywhere structural information exists,
- `duplicate_keys.go` proves the parser can already support more than raw parse errors,
- the new `sanitize parse` helper makes this layer inspectable from the CLI.

### Decision 3: return plain analysis data, not raw nodes

Rationale:

- tree-sitter node lifetimes are coupled to the tree and parser,
- plain Go data is safer for fix routines, tests, CLI JSON, and server responses.

### Decision 4: improve fix targeting gradually

Rationale:

- the experiment matrix proves parse locations are useful but not universally precise,
- trying to make every fixer fully syntax-aware in one step would slow the project down.

## Alternatives Considered

### Alternative A: keep the current split and only add more regex rules

Rejected because it ignores the best structural signal already available in the codebase and preserves the repeated parser work.

### Alternative B: trust parse error rows directly and remove heuristics

Rejected because the corpus disproves it. `examples/yaml/10-tab-indent.yaml` and `examples/yaml/21-multiple-errors.yaml` show parse errors anchored to the parent key rather than the actual offending lines, while `examples/yaml/11-missing-space-after-colon.yaml` and `examples/yaml/17-duplicate-key-sibling.yaml` need heuristics even though the parse succeeds.

### Alternative C: expose raw tree-sitter nodes as a public API

Rejected because it leaks parser lifetime concerns and makes downstream code harder to test and serialize.

## Implementation Plan

### Phase 0: baseline instrumentation

Already completed in this ticket:

- add `sanitize parse`,
- add the corpus matrix experiment script,
- write down the observed disagreement between parse and lint output.

### Phase 1: shared analysis refactor

1. Add an internal `analyzeDocument` helper in `pkg/yaml`.
2. Move duplicate-key traversal into that helper.
3. Add reusable line-index helpers.
4. Make `ParseTree` and `Lint` wrappers over shared analysis.

### Phase 2: richer issues

1. Expand `LintIssue` to full spans plus source metadata.
2. Convert parse errors into lint issues.
3. Update JSON tests and CLI output tests accordingly.

### Phase 3: tree-aware rules

1. Add `mixed_indent` linting.
2. Make `extra_colon_in_value` consult structural context.
3. Replace the current `hasTreeErr` start-row shortcut with analysis-driven windows.

### Phase 4: sanitize loop cleanup

1. Change `applyFixes` to accept the shared analysis object.
2. Ensure each sanitize iteration analyzes once.
3. Verify that fix ordering still behaves correctly.

### Phase 5: surface updates

1. Update README examples.
2. Update server payload examples if the issue model changes.
3. Decide whether a public `Analyze` API is needed.

## Testing and Validation Strategy

The validation strategy should mirror the three failure classes from the corpus.

### Parse-only tests

Use:

- `examples/yaml/20-mixed-indent.yaml`
- `examples/yaml/28-unresolved-parse-error.yaml`

Assertions:

- parser-derived lint issues exist,
- spans are sensible,
- fixes stay conservative when no clear repair exists.

### Heuristic-only tests

Use:

- `examples/yaml/11-missing-space-after-colon.yaml`
- `examples/yaml/13-list-dash-no-space.yaml`
- `examples/yaml/17-duplicate-key-sibling.yaml`

Assertions:

- parse remains clean,
- lint still emits issues,
- fixes still apply as before.

### Hybrid tests

Use:

- `examples/yaml/10-tab-indent.yaml`
- `examples/yaml/21-multiple-errors.yaml`
- `examples/yaml/24-deeply-nested-mixed-errors.yaml`

Assertions:

- parse-derived issues and heuristic issues both appear,
- fix targeting improves,
- repeated sanitize iterations converge.

### CLI and tooling tests

Keep and extend:

- `cmd/sanitize/main_test.go` for `sanitize parse`, `sanitize lint`, and `sanitize fix`.
- `go run ./ttmp/.../scripts/parse_lint_matrix.go` as a quick regression probe for corpus behavior.

## Risks, Tradeoffs, and Open Questions

### Risk 1: richer issue objects ripple into callers

If `LintIssue` changes shape, CLI JSON, server JSON, and any downstream code that decodes the struct will change too. The repo currently has no external users, so this is acceptable, but the change should still be explicit.

### Risk 2: parse windows may still be too broad

Some YAML failures are reported at a parent mapping line. The design should therefore treat parse windows as strong hints, not absolute truth.

### Risk 3: fix ordering may shift behavior

Once fixers start consuming richer analysis, the current ordering of line fixes versus document-level fixes may expose new interactions. This needs table-driven regression coverage.

### Open question: public `Analyze` API or internal helper only

There is a good internal reason to create `analyzeDocument`, but it is still optional whether the package exposes a public `Analyze` function. My recommendation is:

1. implement it internally first,
2. expose it publicly only if the CLI, server, or future JSON support truly need the same raw analysis bundle.

## References

- `pkg/yaml/parse.go:10-79`
- `pkg/yaml/types.go:13-54`
- `pkg/yaml/lint.go:18-102`
- `pkg/yaml/duplicate_keys.go:18-142`
- `pkg/yaml/fix.go:15-295`
- `pkg/yaml/sanitize.go:5-80`
- `pkg/yaml/options.go:3-54`
- `pkg/yaml/sanitize_test.go:48-340`
- `internal/cli/commands.go`
- `cmd/sanitize/main_test.go`
- `sources/01-example-corpus-parse-vs-lint-matrix.md`
- `scripts/parse_lint_matrix.go`
