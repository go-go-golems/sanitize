---
Title: Intern guide to tree-sitter-aware linting in sanitize
Ticket: SANITIZE-004
Status: active
Topics:
    - yaml
    - linting
    - treesitter
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: README.md
      Note: High-level user-facing entry point for the intern guide
    - Path: internal/cli/commands.go
      Note: CLI surfaces explained in the guide
    - Path: pkg/yaml/analysis.go
      Note: New shared analysis entrypoint the guide now references
    - Path: pkg/yaml/fix.go
      Note: Fix pipeline described in the guide
    - Path: pkg/yaml/line_index.go
      Note: Location helper now part of the implementation story
    - Path: pkg/yaml/lint.go
      Note: Current linter behavior the guide explains
    - Path: pkg/yaml/parse.go
      Note: Tree-sitter parsing entrypoint for the intern guide
    - Path: pkg/yaml/sanitize.go
      Note: Core orchestration loop described for onboarding
    - Path: pkg/yaml/types.go
      Note: Core API types explained to new contributors
ExternalSources: []
Summary: Detailed onboarding guide for an intern who needs to understand the sanitize parser, linter, fixer, and the proposed tree-sitter-aware redesign.
LastUpdated: 2026-03-13T08:58:45.541255348-04:00
WhatFor: Teach a new engineer how the sanitize system works today and how to implement the proposed tree-sitter-aware linting refactor.
WhenToUse: Use before changing pkg/yaml, the CLI parse/lint/fix commands, or tree-sitter-backed diagnostic behavior.
---



# Intern guide to tree-sitter-aware linting in sanitize

## Goal

This guide is for a new intern joining the project cold. It explains what the repository does, how YAML moves through the system today, where tree-sitter fits, why the current linting design is incomplete, and how to approach the next implementation step without guessing.

## Context

`sanitize` is a Go package and CLI for cleaning up broken YAML. It has four user-facing operations:

- `sanitize fix` to apply heuristic repairs,
- `sanitize lint` to report issues,
- `sanitize parse` to expose the tree-sitter parse tree and structural errors,
- `sanitize serve` to run the web UI and HTTP API.

At a high level, the system tries to answer three different questions:

1. Can tree-sitter parse this YAML structurally?
2. Even if it parses, does it contain suspicious or undesirable patterns?
3. If it is broken, can we repair it safely enough to make progress?

Those questions map to different files, and the current implementation does not yet answer them through one unified analysis pass. That is the main architectural theme you should keep in your head while reading the code.

### Read order

If you have never seen this codebase before, read these files in this order:

1. `README.md`
2. `pkg/yaml/types.go`
3. `pkg/yaml/parse.go`
4. `pkg/yaml/lint.go`
5. `pkg/yaml/duplicate_keys.go`
6. `pkg/yaml/fix.go`
7. `pkg/yaml/sanitize.go`
8. `pkg/yaml/sanitize_test.go`
9. `internal/cli/commands.go`
10. `internal/server/server.go`

That order moves from data model to parser to lint to fix orchestration to external surfaces.

## Quick Reference

### System map

```text
CLI / HTTP server
      |
      v
  pkg/yaml public API
      |
      +-> ParseTree(src)
      +-> Lint(src)
      +-> Sanitize(src, opts...)
               |
               +-> ParseTree(current)
               +-> Lint(current)
               +-> applyFixes(...)
               +-> repeat
```

### Current API surface

From `pkg/yaml`:

```go
func ParseTree(src string) (string, []ErrorNode, error)
func Lint(src string) []LintIssue
func Sanitize(src string, opts ...Option) Result
```

Key types from `pkg/yaml/types.go`:

```go
type ErrorNode struct {
    Type      string
    StartByte uint
    EndByte   uint
    StartRow  uint
    StartCol  uint
    EndRow    uint
    EndCol    uint
    Text      string
}

type LintIssue struct {
    Rule        string
    Description string
    Row         int
}

type Result struct {
    Original           string
    Sanitized          string
    TreeText           string
    OriginalTreeText   string
    Errors             []ErrorNode
    OriginalErrors     []ErrorNode
    LintIssues         []LintIssue
    OriginalLintIssues []LintIssue
    Fixes              []Fix
    ParseClean         bool
    LintClean          bool
}
```

The asymmetry between `ErrorNode` and `LintIssue` is one of the biggest clues that the system is not structurally unified yet.

### File-by-file job descriptions

`pkg/yaml/parse.go`

- creates a tree-sitter parser for YAML,
- parses one source string,
- walks the tree for `ERROR` and `MISSING` nodes,
- returns the S-expression and structural errors.

`pkg/yaml/lint.go`

- runs regex-style lint checks line by line,
- currently has no shared parse-analysis input,
- calls duplicate-key detection after line scanning.

`pkg/yaml/duplicate_keys.go`

- reparses YAML,
- walks mapping nodes,
- detects duplicate sibling keys,
- is already proof that tree traversal helps linting.

`pkg/yaml/fix.go`

- applies line-level and document-level fixes,
- currently gets only raw parse errors plus lint issues,
- uses a start-row boolean to decide whether tree errors are relevant on a line.

`pkg/yaml/sanitize.go`

- is the orchestrator,
- captures original state,
- loops through parse, lint, fix until clean or stuck.

`internal/cli/commands.go`

- wires the public `fix`, `lint`, `parse`, and `serve` subcommands,
- is the easiest place to expose debugging and inspection helpers.

### Why tree-sitter matters here

Tree-sitter gives the project three things that regexes do not:

- a full parse tree,
- byte and row and column spans,
- access to structural node kinds such as mappings and pairs.

The project already uses that power in two places:

- `ParseTree` for structural error extraction,
- duplicate-key detection for mapping traversal.

The next refactor is about using the same structural layer for all diagnostics, not just those two isolated cases.

### Current failure taxonomy

When you run the corpus matrix in `sources/01-example-corpus-parse-vs-lint-matrix.md`, you should think in these categories.

#### Parse-only failures

These are files where tree-sitter sees a structural problem but `Lint` stays silent.

Examples:

- `examples/yaml/20-mixed-indent.yaml`
- `examples/yaml/28-unresolved-parse-error.yaml`

Implication:

- lint needs parser-derived structural issues.

#### Heuristic-only failures

These are files that still parse but should still be flagged.

Examples:

- `examples/yaml/11-missing-space-after-colon.yaml`
- `examples/yaml/14-trailing-comma-flow-mapping.yaml`
- `examples/yaml/17-duplicate-key-sibling.yaml`

Implication:

- tree-sitter cannot replace heuristics.

#### Hybrid failures

These are files where tree-sitter notices breakage, but heuristics are still needed to localize or explain it.

Examples:

- `examples/yaml/10-tab-indent.yaml`
- `examples/yaml/21-multiple-errors.yaml`
- `examples/yaml/24-deeply-nested-mixed-errors.yaml`

Implication:

- the future design must combine parse context with heuristics.

### Important invariant

Do not expose raw tree-sitter nodes broadly unless you have to. The safe pattern is:

1. parse with tree-sitter,
2. extract plain Go facts,
3. close the parser and tree,
4. operate on plain Go facts everywhere else.

That keeps tests simple and avoids lifetime bugs.

### Current sanitize loop in pseudocode

```go
func Sanitize(src string, opts ...Option) Result {
    cfg := defaultConfig()
    original := src
    origTree, origErrors, _ := ParseTree(original)
    origIssues := Lint(original)

    for iter := 0; iter < cfg.maxIterations; iter++ {
        tree, errors, _ := ParseTree(src)
        issues := Lint(src)
        if len(errors) == 0 && len(issues) == 0 {
            return cleanResult(...)
        }

        fixed, fixes := applyFixes(src, errors, issues, &cfg)
        if len(fixes) == 0 {
            return stuckResult(...)
        }
        src = fixed
    }

    return finalResult(...)
}
```

Notice the duplication:

- `ParseTree` parses once,
- `Lint` reparses for duplicate keys,
- the loop repeats that work each iteration.

### Proposed future loop in pseudocode

```go
func sanitize(src string, cfg config) Result {
    orig := analyzeDocument(src)
    current := src

    for iter := 0; iter < cfg.maxIterations; iter++ {
        a := analyzeDocument(current)
        if a.IsClean() {
            return buildResult(orig, a, current)
        }

        next, fixes := applyFixes(current, a, &cfg)
        if len(fixes) == 0 {
            return buildResult(orig, a, current)
        }
        current = next
    }

    final := analyzeDocument(current)
    return buildResult(orig, final, current)
}
```

This is the design target you should keep coming back to.

### Suggested internal types for the refactor

```go
type analysis struct {
    TreeText      string
    ParseErrors   []ErrorNode
    DuplicateKeys []duplicateKeyOccurrence
    LintIssues    []LintIssue
    LineStarts    []int
}

type LintIssue struct {
    Rule        string
    Source      string
    Description string
    StartByte   uint
    EndByte     uint
    StartRow    uint
    StartCol    uint
    EndRow      uint
    EndCol      uint
}
```

### Implementation checklist for an intern

1. Run the existing tests first so you know the baseline.
2. Read `sources/01-example-corpus-parse-vs-lint-matrix.md` to understand current behavior.
3. Add the shared analysis helper without changing public behavior yet.
4. Move duplicate-key detection onto the shared analysis helper.
5. Expand `LintIssue` and update tests.
6. Add parser-derived lint issues.
7. Change `applyFixes` to consume the shared analysis helper.
8. Re-run the corpus matrix after each milestone.

### Diagram: current vs proposed architecture

```text
Current
-------
ParseTree()   Lint()                    applyFixes()
    |           |                            |
    |           +-> regex rules              |
    |           +-> duplicate-key reparse    |
    +----------------------------------------+
                 loosely connected by caller

Proposed
--------
            analyzeDocument()
                   |
    +--------------+--------------+
    |              |              |
ParseTree()      Lint()       applyFixes()
```

## Usage Examples

### Inspect a single file structurally

```bash
go run ./cmd/sanitize parse examples/yaml/21-multiple-errors.yaml
go run ./cmd/sanitize parse --json examples/yaml/21-multiple-errors.yaml
```

Use this when you want to see whether tree-sitter itself sees a structural failure and where it anchors that failure.

### Generate the corpus matrix

```bash
go run ./ttmp/2026/03/13/SANITIZE-004--improve-yaml-linting-and-sanitizing-with-tree-sitter-aware-analysis/scripts/parse_lint_matrix.go \
  > ./ttmp/2026/03/13/SANITIZE-004--improve-yaml-linting-and-sanitizing-with-tree-sitter-aware-analysis/sources/01-example-corpus-parse-vs-lint-matrix.md
```

Use this when you change lint rules or parse handling and want to see whether coverage shifts.

### Run the relevant tests

```bash
go test ./cmd/sanitize ./internal/... ./pkg/yaml
```

### Review the CLI behavior

```bash
go run ./cmd/sanitize lint examples/yaml/20-mixed-indent.yaml
go run ./cmd/sanitize fix --json examples/yaml/21-multiple-errors.yaml
go run ./cmd/sanitize serve --port 8080
```

### How to review a future refactor

Start with:

1. `pkg/yaml/types.go`
2. the new shared analysis helper
3. `pkg/yaml/lint.go`
4. `pkg/yaml/fix.go`
5. `pkg/yaml/sanitize.go`
6. `cmd/sanitize/main_test.go`

Then validate with:

```bash
go test ./cmd/sanitize ./internal/... ./pkg/yaml
go run ./ttmp/2026/03/13/SANITIZE-004--improve-yaml-linting-and-sanitizing-with-tree-sitter-aware-analysis/scripts/parse_lint_matrix.go
```

## Related

- Design doc: `design-doc/01-tree-sitter-driven-yaml-linting-and-sanitizing-analysis.md`
- Diary: `reference/01-diary.md`
- Corpus matrix: `sources/01-example-corpus-parse-vs-lint-matrix.md`
- Experiment script: `scripts/parse_lint_matrix.go`
- Core implementation files:
  - `pkg/yaml/parse.go`
  - `pkg/yaml/lint.go`
  - `pkg/yaml/duplicate_keys.go`
  - `pkg/yaml/fix.go`
  - `pkg/yaml/sanitize.go`
