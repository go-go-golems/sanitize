---
Title: Intern guide to sanitize
Ticket: SANITIZE-002
Status: active
Topics:
    - json
    - api-design
    - release-readiness
    - documentation
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: README.md
      Note: High-level orientation for new contributors
    - Path: cmd/sanitize-server/main.go
      Note: HTTP/API walkthrough for interns
    - Path: cmd/sanitize-server/static/index.html
      Note: Frontend walkthrough and browser flow
    - Path: cmd/sanitize/main.go
      Note: CLI walkthrough for interns
    - Path: pkg/yaml/fix.go
      Note: Fixer walkthrough
    - Path: pkg/yaml/lint.go
      Note: Rule-detection walkthrough
    - Path: pkg/yaml/sanitize.go
      Note: Core orchestration walkthrough
ExternalSources: []
Summary: ""
LastUpdated: 2026-03-13T00:47:58.040748408-04:00
WhatFor: ""
WhenToUse: ""
---


# Intern guide to sanitize

## Goal

This guide explains the current `sanitize` system in enough detail that a new intern can understand how the library, CLI, server, and web UI fit together, and then use that understanding to implement JSON support safely. It is intentionally more detailed than a normal README. It is meant to be read before editing code.

## Context

`sanitize` is a Go module that detects and heuristically fixes common YAML mistakes. The key idea is simple:

1. Parse the input with tree-sitter.
2. Run lint-style heuristics to find common mistakes that a parser might not report directly.
3. Apply safe-enough fixes.
4. Repeat until the document is clean or the tool stops making progress.

Today the repository only supports YAML. The module name, CLI names, and earlier planning document suggest that JSON support is intended next. That makes the current codebase important in two ways:

- It is the implementation you will extend.
- It is also the design you will gradually generalize.

If you only remember one thing from this guide, remember this: the current code is not a generic sanitizer with a YAML plugin. It is a YAML sanitizer with a packaging shape that hints at future generalization.

## Quick Reference

### Repository map

```text
.
├── README.md                            # user-facing overview and examples
├── Makefile                             # test/lint/build/release helper targets
├── cmd/
│   ├── sanitize/main.go                 # CLI entrypoint
│   └── sanitize-server/
│       ├── main.go                      # HTTP server entrypoint
│       └── static/index.html            # embedded single-page app
├── pkg/
│   └── yaml/
│       ├── types.go                     # exported data structures
│       ├── options.go                   # functional options and config
│       ├── parse.go                     # tree-sitter parser integration
│       ├── lint.go                      # YAML heuristics that report issues
│       ├── fix.go                       # YAML heuristics that rewrite text
│       ├── sanitize.go                  # top-level orchestration loop
│       ├── examples.go                  # sample broken YAML documents
│       └── sanitize_test.go             # unit/integration tests for pkg/yaml
└── .github/workflows/                   # CI and release workflows
```

### Most important APIs today

Library API:

```go
result := yamlsanitize.Sanitize(src, yamlsanitize.WithTabWidth(4))
issues := yamlsanitize.Lint(src)
treeText, errs, err := yamlsanitize.ParseTree(src)
```

CLI API:

```bash
sanitize broken.yaml
sanitize --lint broken.yaml
sanitize --json broken.yaml
```

Server/API:

```text
GET  /api/examples
POST /api/parse
POST /api/sanitize
```

### One-page mental model

```text
           +----------------------+
           |  raw YAML text input |
           +----------+-----------+
                      |
                      v
              +---------------+
              | ParseTree()   |
              | tree-sitter   |
              +-------+-------+
                      |
          +-----------+-----------+
          |                       |
          v                       v
   parse errors             Lint() issues
                                   |
                                   v
                            applyFixes()
                                   |
                                   v
                           new candidate text
                                   |
                                   v
                      repeat until clean or stuck
```

### File-by-file purpose

`pkg/yaml/types.go`

- Defines the external data model for parsing, linting, and sanitizing.
- Important because the server and CLI both serialize these types to JSON.

`pkg/yaml/options.go`

- Defines the functional option pattern.
- Important because the sanitize loop is configurable but keeps a small public API.

`pkg/yaml/parse.go`

- Owns parser creation and parse-error extraction.
- Important because this is the place most likely to need abstraction for JSON support.

`pkg/yaml/lint.go`

- Reports issues that are not necessarily parse failures.
- Important because heuristics live here, and heuristics are where false positives can happen.

`pkg/yaml/fix.go`

- Applies actual text edits.
- Important because this is where the tool can accidentally change user meaning.

`pkg/yaml/sanitize.go`

- Orchestrates the loop.
- Important because it is the best candidate for extraction into a shared format-independent runtime.

`cmd/sanitize/main.go`

- Converts CLI flags and file/stdin input into package calls.
- Important because public UX and exit-code behavior live here.

`cmd/sanitize-server/main.go`

- Wraps the package in HTTP handlers.
- Important because API contracts and security hardening live here.

`cmd/sanitize-server/static/index.html`

- The complete frontend.
- Important because it reveals what data the server must keep returning.

### Key current references

Use these line ranges when reading the code:

- `pkg/yaml/sanitize.go:5-80` for the orchestration loop.
- `pkg/yaml/lint.go:19-110` for issue detection rules.
- `pkg/yaml/fix.go:16-316` for fixer behavior.
- `cmd/sanitize/main.go:13-70` for CLI behavior.
- `cmd/sanitize-server/main.go:23-99` for HTTP behavior.
- `cmd/sanitize-server/static/index.html:325-552` for the browser-driven sanitize flow.

## System Walkthrough

### 1. The exported data model

Before you understand the control flow, understand the data model in `pkg/yaml/types.go:13-60`.

`ErrorNode`

- Represents tree-sitter `ERROR` or `MISSING` nodes.
- Carries byte offsets, row/column information, and the original text when available.

`LintIssue`

- Represents a heuristic issue.
- Not every lint issue is a parser error.
- Example: duplicate keys are often semantically wrong even if the parser accepts them.

`Fix`

- Represents one text transformation.
- The server UI displays these in the "Applied fixes" pane.

`Result`

- Bundles the entire sanitize run.
- Includes both original and final parse/lint state.
- This is why the frontend can show "what was wrong originally" and "what the sanitizer produced" at the same time.

Important subtlety:

The result intentionally keeps both original and final error sets. That means any future shared JSON runtime should preserve this dual-state structure instead of collapsing it into a single "current issues" field.

### 2. The option system

`pkg/yaml/options.go:3-54` uses the standard Go functional-option pattern.

Why this matters:

- The public API stays small.
- Internals can evolve without changing every callsite.
- JSON support should reuse this pattern rather than invent a second configuration style.

Current options:

- `WithMaxIterations(n)`
- `WithTabWidth(w)`
- `WithRules(rules...)`

If you add JSON support, think carefully about which options are:

- cross-format options, such as max iterations
- YAML-only options, such as tab width
- JSON-only options, such as whether to strip comments

### 3. The parse phase

`pkg/yaml/parse.go:10-35` creates a YAML tree-sitter parser and returns a sexp tree plus collected errors.

What happens in detail:

1. `newParser()` creates a tree-sitter parser and installs the YAML grammar.
2. `ParseTree()` parses the byte slice and gets the root node.
3. `collectErrors()` walks the tree recursively.
4. Any `ERROR` or `MISSING` node becomes an `ErrorNode`.

Pseudocode:

```go
func ParseTree(src string) (treeText string, errs []ErrorNode, err error) {
    parser := newParser()
    tree := parser.Parse([]byte(src), nil)
    root := tree.RootNode()
    errs = collectErrors(root, []byte(src))
    return root.ToSexp(), errs, nil
}
```

What to notice as an engineer:

- The parser is created fresh on every call.
- The code walks the full tree every time.
- The sanitize loop calls `ParseTree()` repeatedly.

That is acceptable for the current repository size, but it is exactly the kind of logic you would want to centralize or optimize once a second format is added.

### 4. The lint phase

`pkg/yaml/lint.go:19-110` runs a set of regex-based heuristics over the raw source lines.

Current rules:

- `tab_indent`
- `missing_space_after_colon`
- `list_dash_no_space`
- `trailing_comma`
- `duplicate_key`
- `extra_colon_in_value`

Why lint exists even though parsing already exists:

Some user mistakes either:

- are accepted by parsers but still undesirable, or
- produce parse failures that are hard to explain directly.

Lint gives the sanitizer domain-specific hints about what to fix.

Important design lesson:

Lint heuristics can be more dangerous than parser errors because they can over-match. The duplicate-key issue in the release review is a concrete example: the lint stage is where a public tool can start lying about user input.

### 5. The fix phase

`pkg/yaml/fix.go:16-316` applies text edits based on lint issues and parse errors.

The structure is:

1. `applyFixes()` splits the source into lines.
2. It indexes lint issues by row.
3. It indexes parse errors by start row.
4. It runs `fixLine()` on each line.
5. If no line changed, it tries document-level fixes:
   - duplicate key renaming
   - mixed indentation normalization

Pseudocode:

```go
func applyFixes(src string, errors []ErrorNode, issues []LintIssue, cfg *config) (string, []Fix) {
    lines := strings.Split(src, "\n")

    for row := range lines {
        lines[row], lineFixes = fixLine(lines[row], row, rulesForRow, hasTreeErr, cfg)
        fixes = append(fixes, lineFixes...)
    }

    if nothingChanged {
        try document-level fixes
    }

    return strings.Join(lines, "\n"), fixes
}
```

The key engineering judgment here is not "can we rewrite the text?" It is "is the rewrite conservative enough that users will trust it?" For JSON support, that question is even more important because JSON is often machine-generated and machine-consumed.

### 6. The orchestration loop

The most important file in the repo is `pkg/yaml/sanitize.go:5-80`.

What it does:

1. Builds config from options.
2. Captures original parse and lint state.
3. Loops up to `maxIterations`.
4. If the current document is clean, returns success.
5. Otherwise applies fixes.
6. If no fixes were possible, returns the current state.
7. After the loop, returns the last observed state.

This is the correct place to look if you want to understand the product concept. The rest of the package exists to support this loop.

### 7. The CLI layer

`cmd/sanitize/main.go:13-70` is intentionally thin:

1. Parse flags.
2. Read input from a file or stdin.
3. Call `yamlsanitize.Lint()` or `yamlsanitize.Sanitize()`.
4. Emit text or JSON.
5. Set the process exit code.

This thinness is good. Entry points should stay thin. If JSON support lands cleanly, the CLI should still mostly just:

- resolve format
- pick an engine
- render output

### 8. The server and browser UI

The server in `cmd/sanitize-server/main.go:23-99` exists mainly to expose three things to the static frontend:

- examples
- parse results
- sanitize results

The UI in `cmd/sanitize-server/static/index.html` is monolithic but readable if you treat it as three blocks:

1. DOM and styling (`1-289`)
2. browser->server flow (`297-342`)
3. rendering helpers (`345-552`)

The browser flow is:

```text
user types in textarea
  -> debounce
  -> POST /api/sanitize { yaml: src }
  -> receive Result JSON
  -> render badges, parse tree, issues, sanitized output, fixes
```

That means JSON support is not just a backend package task. The UI contract is part of the system.

## What You Need To Understand Before Adding JSON

### The system is split by responsibility, not by deployment target

Do not start with the CLI or server. Start in `pkg/yaml`. The CLI and server are thin wrappers. If you add JSON support by branching logic only in entrypoints, you will create duplication fast.

### There are two kinds of "errors"

The codebase distinguishes:

- parse errors from tree-sitter
- lint issues from heuristics

Do not blur them together in a new API. The distinction is useful to users and to the UI.

### The UI expects rich result objects

The frontend does not just want "clean or not clean." It wants:

- original tree
- current tree
- original errors
- current errors
- original lint issues
- current lint issues
- fixes

If you change the server contract, you must account for the frontend expectations explicitly.

### Exit codes matter

The release review found that JSON output mode can mask failure in the CLI. That is a reminder that behavior outside the core library still matters. A technically correct library is not enough if the entrypoints give confusing signals.

## Suggested Implementation Strategy For An Intern

### Step 1: read the code in this order

1. `README.md`
2. `pkg/yaml/types.go`
3. `pkg/yaml/options.go`
4. `pkg/yaml/sanitize.go`
5. `pkg/yaml/lint.go`
6. `pkg/yaml/fix.go`
7. `cmd/sanitize/main.go`
8. `cmd/sanitize-server/main.go`
9. `cmd/sanitize-server/static/index.html`

### Step 2: write down the existing contracts before editing anything

Capture:

- what `Result` means
- what the CLI promises
- what the HTTP API expects
- what the frontend reads from the result payload

Doing this first prevents accidental breaking changes.

### Step 3: extract shared logic before adding new format code

Resist the temptation to create `pkg/json` by copying files from `pkg/yaml`. The repo is small now, so this is exactly the right moment to extract a core runtime.

### Step 4: ship JSON in safe increments

Recommended order:

1. shared core
2. JSON parse-only support
3. conservative JSON fixes
4. CLI/server/UI wiring
5. release hardening

## API reference sketches

### Current API

```go
package yamlsanitize

func ParseTree(src string) (string, []ErrorNode, error)
func Lint(src string) []LintIssue
func Sanitize(src string, opts ...Option) Result

func WithMaxIterations(n int) Option
func WithTabWidth(w int) Option
func WithRules(rules ...string) Option
```

### Recommended future API

```go
package core

type Engine interface {
    Format() Format
    ParseTree(src string) (string, []ParseError, error)
    Lint(src string, cfg Config) []Issue
    ApplyFixes(src string, errs []ParseError, issues []Issue, cfg Config) (string, []Fix)
}

func Run(engine Engine, src string, opts ...Option) Result
```

## Commands you should know

```bash
go test ./...
go test -race ./...
make lint
make build
make gosec
sanitize --lint broken.yaml
sanitize --json broken.yaml
sanitize-server
```

## Common pitfalls

1. Do not assume a lint heuristic is correct just because it seems reasonable. Validate it against realistic examples.
2. Do not change the shape of `Result` casually. The UI depends on it.
3. Do not add JSON support by only touching the CLI. That produces a shallow feature with poor internal structure.
4. Do not forget exit-code behavior when adding machine-readable output.
5. Do not expose the server more broadly without adding timeouts and bounded request sizes.

## Usage Examples

### Example: tracing one sanitize request end-to-end

If a user runs:

```bash
printf 'name:Alice\n' | sanitize
```

the code path is:

1. `cmd/sanitize/main.go` reads stdin.
2. `yamlsanitize.Sanitize()` is called.
3. `pkg/yaml/sanitize.go` captures original parse/lint state.
4. `pkg/yaml/lint.go` reports `missing_space_after_colon`.
5. `pkg/yaml/fix.go` changes `name:Alice` to `name: Alice`.
6. `pkg/yaml/sanitize.go` reruns parse/lint, sees the document is clean, and returns.
7. The CLI prints sanitized text to stdout and fix information to stderr.

### Example: where to add a new YAML rule today

If you wanted to add a new YAML rule right now, you would usually touch:

1. `pkg/yaml/lint.go` to detect the issue.
2. `pkg/yaml/fix.go` to implement the rewrite.
3. `pkg/yaml/examples.go` to add a sample.
4. `pkg/yaml/sanitize_test.go` to add tests.

### Example: where JSON support should land

After extracting a shared core, JSON support should touch:

1. `pkg/core/*`
2. `pkg/json/*`
3. `cmd/sanitize/main.go`
4. `cmd/sanitize-server/main.go`
5. `cmd/sanitize-server/static/index.html`

## Related

- `../../SANITIZE-001--turn-yaml-sanitizer-into-reusable-go-go-golems-sanitize-package/design-doc/01-implementation-plan.md`
- `../design-doc/01-public-release-review.md`
- `../design-doc/02-json-support-architecture-and-implementation-plan.md`
- `README.md`
