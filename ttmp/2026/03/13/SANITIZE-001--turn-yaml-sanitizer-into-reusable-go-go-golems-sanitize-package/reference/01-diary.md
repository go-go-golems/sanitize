---
Title: Diary
Ticket: SANITIZE-001
Status: active
Topics:
    - refactoring
    - packaging
    - testing
    - documentation
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/types.go:Extracted types (ErrorNode, LintIssue, Fix, Result, Example)
    - /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/parse.go:Tree-sitter parser and error collector
    - /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/lint.go:Lint rules extracted from sanitize.go
    - /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/fix.go:Fix rules extracted from sanitize.go
    - /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/sanitize.go:Core Sanitize() entry point with options
    - /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/options.go:Functional options (WithMaxIterations, WithTabWidth, WithRules)
    - /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/sanitize_test.go:Comprehensive test suite (27 tests)
    - /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/cmd/sanitize/main.go:CLI tool
    - /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/cmd/sanitize-server/main.go:HTTP server with go:embed
ExternalSources: []
Summary: "Implementation diary for SANITIZE-001"
LastUpdated: 2026-03-13T00:30:03.146666824-04:00
WhatFor: ""
WhenToUse: ""
---

# Diary

## Goal

Record the implementation journey of converting the yaml-sanitizer prototype into a publishable `github.com/go-go-golems/sanitize` package.

## Step 1: Analysis and Planning

Assessed the existing codebase: ~630 lines of Go (single `sanitize/` package + `main.go` HTTP server) and ~550 lines of HTML/JS web UI. Single commit, no tests, module path `yaml-sanitizer`.

### Prompt Context

**User prompt (verbatim):** "Analyze this repository and tell me how much effort it would be to turn this into a nice reusable package with docs."

**Assistant interpretation:** Evaluate the current state of the repo and estimate the work needed to make it publishable.

**Inferred user intent:** Get a clear picture of what's needed before committing to the work.

### What I did
- Read all source files (main.go, sanitize/sanitize.go, sanitize/examples.go, static/index.html)
- Assessed code quality, structure, and missing elements

### What worked
- The existing `sanitize/` package was already well-separated from the HTTP layer
- Types were already exported with JSON tags and doc comments

### What I learned
- The code was surprisingly clean for a prototype — the main issues were structural (module path, no tests, no CLI) rather than code quality

## Step 2: Full Restructure and Implementation

Implemented all 10 tasks in a single pass: module rename, package restructure, regex precompilation, functional options, go:embed, CLI, tests, CI config, and README.

### Prompt Context

**User prompt (verbatim):** "alright do it, package is github.com/go-go-golems/sanitize (because we want to support json later too). Create a docmgr ticket and add an implementation plan and tasks to that new ticket. then implement them and commit at appropriate intervals, and keep a diary (read the skills for this)."

**Assistant interpretation:** Execute the full restructure as planned, using docmgr for project tracking.

**Inferred user intent:** Get a publishable, well-structured package with CI, tests, and docs in one go.

**Commit (code):** f369fe1 — "Restructure as github.com/go-go-golems/sanitize reusable package"

### What I did
- Created `pkg/yaml/` with 7 files split from the original monolith:
  - `types.go` — all exported types
  - `parse.go` — tree-sitter parser + error collector
  - `lint.go` — lint rules with precompiled regexes
  - `fix.go` — fix rules with precompiled regexes
  - `sanitize.go` — core `Sanitize()` entry point
  - `options.go` — functional options (`WithMaxIterations`, `WithTabWidth`, `WithRules`)
  - `examples.go` — built-in examples
- Created `cmd/sanitize/` — CLI tool (stdin/file, --lint, --json, --tab-width, --max-iterations)
- Created `cmd/sanitize-server/` — HTTP server with `go:embed` for static assets
- Wrote 27 tests covering all lint rules, all fix rules, options, edge cases, helpers
- Copied CI from `~/code/wesen/corporate-headquarters/go-template/`:
  - `.golangci.yml` + `.golangci-lint-version`
  - `.github/workflows/{push,lint,release}.yml`
  - `.goreleaser.yaml` (adapted for two binaries: sanitize + sanitize-server)
  - `Makefile`, `lefthook.yml`
- Wrote `README.md` with install, CLI usage, library API, and rules table

### Why
- User wants a publishable package under go-go-golems org
- Named `sanitize` (not yaml-sanitize) to anticipate future JSON support

### What worked
- All 27 tests passed on first run
- Build succeeded immediately after module rename
- The existing code structure made the split clean — no circular dependencies

### What didn't work
- Nothing significant failed. The old `main.go` at root still had the old import path after module rename, but that was expected since it was being replaced.

### What was tricky to build
- The `WithRules` option needed to thread through `applyFixes` → `fixLine`, requiring a `config` parameter to be added to both functions. Had to be careful that the rule check happens in `fixLine` (per-rule) while `applyFixes` handles document-level fixes.
- The `.goreleaser.yaml` needed two build IDs per OS (sanitize + sanitize-server) since CGO is required for tree-sitter.

### What warrants a second pair of eyes
- The `fixExtraColonInValue` rule fires both on lint match AND on tree-sitter error (`hasRule || hasTreeErr`). This aggressive triggering might over-quote in edge cases.
- Duplicate key detection is indent-level based, which can produce false positives across separate mapping blocks at the same indent.

### What should be done in the future
- Add JSON sanitization support (the package name `sanitize` was chosen for this)
- Consider using a proper YAML library (gopkg.in/yaml.v3) alongside tree-sitter for validation
- Add benchmarks for the regex-heavy lint/fix paths

### Code review instructions
- Start at `pkg/yaml/sanitize.go` — the main entry point
- Then `pkg/yaml/options.go` for the options API
- `pkg/yaml/lint.go` and `pkg/yaml/fix.go` are the rule implementations
- Run `go test ./... -v` to validate
- Run `go build ./cmd/sanitize && echo "name:Alice" | ./cmd/sanitize` to test CLI

### Technical details

Package layout:
```
pkg/yaml/         — package yamlsanitize (core library)
cmd/sanitize/     — CLI (stdin→stdout)
cmd/sanitize-server/ — HTTP server + embedded web UI
```

Key API:
```go
yamlsanitize.Sanitize(src string, opts ...Option) Result
yamlsanitize.Lint(src string) []LintIssue
yamlsanitize.ParseTree(src string) (string, []ErrorNode, error)
```
