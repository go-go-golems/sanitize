---
Title: JSON support architecture and implementation plan
Ticket: SANITIZE-003
Status: active
Topics:
    - json
    - api-design
    - release-readiness
    - documentation
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: README.md
      Note: Current public API and packaging direction
    - Path: cmd/sanitize-server/main.go
      Note: Server request and response model
    - Path: cmd/sanitize-server/static/index.html
      Note: UI data dependencies and YAML assumptions
    - Path: cmd/sanitize/main.go
      Note: CLI format/output contract
    - Path: pkg/yaml/options.go
      Note: Shared option model candidate
    - Path: pkg/yaml/sanitize.go
      Note: Shared iterative sanitize loop candidate
    - Path: pkg/yaml/types.go
      Note: Shared result model candidate
ExternalSources: []
Summary: ""
LastUpdated: 2026-03-13T00:47:57.914150608-04:00
WhatFor: ""
WhenToUse: ""
---



# JSON support architecture and implementation plan

## Executive Summary

The current repository is organized around a single YAML engine in `pkg/yaml`, but the module name and existing `SANITIZE-001` planning both point toward a broader "sanitize structured text" tool. JSON support should be added by extracting a small shared orchestration layer and then implementing a new JSON format engine, not by cloning the YAML package and drifting into copy-paste maintenance.

The recommended path is:

1. Introduce a shared sanitize core for common result types, options, and the iterative fix loop.
2. Replace the current hand-written CLI with a Glazed/Cobra command tree that models format and output cleanly.
3. Add `pkg/json` with a conservative ruleset focused on common broken-JSON cases.
4. Update the HTTP API and UI to use explicit format-aware contracts without carrying legacy request shapes.

The design below aims to make JSON support practical while simplifying the codebase. Because there are no users yet, the plan intentionally prefers a clean break over compatibility scaffolding.

## Problem Statement

The repository currently sanitizes YAML only:

- The package namespace is `pkg/yaml`.
- The library API is `yamlsanitize.Sanitize`, `yamlsanitize.Lint`, and `yamlsanitize.ParseTree`.
- The CLI imports only `pkg/yaml` in `cmd/sanitize/main.go:10`.
- The server request payloads are shaped around a `yaml` field in `cmd/sanitize-server/main.go:56-58` and `cmd/sanitize-server/main.go:78-80`.
- The web UI presents itself entirely as "YAML Sanitizer" in `cmd/sanitize-server/static/index.html:208-267`.

That design is coherent for YAML, but it does not scale cleanly to JSON:

1. Shared result and option types would be duplicated.
2. CLI and server flags would become ambiguous.
3. Cross-format behavior would diverge without a common contract.
4. The UI would have no explicit format model.

The goal is to support JSON without turning the codebase into two near-identical packages and two incompatible frontends.

### Scope

In scope:

- JSON parse/lint/fix support.
- Shared internal orchestration.
- CLI/server/API/UI format selection.
- Intentional cleanup of current CLI and library contracts where that improves the architecture.

Out of scope:

- JSON Schema validation.
- Full JSON5 support.
- Automatic semantic transformations such as sorting keys or reformatting style-only output.

## Proposed Solution

Introduce a shared sanitize core and format engines.

### Proposed package layout

```text
github.com/go-go-golems/sanitize/
├── cmd/
│   ├── sanitize/
│   └── sanitize-server/
├── pkg/
│   ├── core/              # shared runtime, types, and options
│   │   ├── types.go
│   │   ├── options.go
│   │   ├── engine.go
│   │   └── run.go
│   ├── yaml/              # YAML adapter and rules
│   │   ├── engine.go
│   │   ├── parse.go
│   │   ├── lint.go
│   │   ├── fix.go
│   │   └── examples.go
│   └── json/              # JSON adapter and rules
│       ├── engine.go
│       ├── parse.go
│       ├── lint.go
│       ├── fix.go
│       ├── examples.go
│       └── sanitize_test.go
└── cmd/sanitize/cmds/      # Glazed/Cobra command groups and verbs
```

### Shared runtime contract

The shared core should own the generic workflow that is currently embedded in `pkg/yaml/sanitize.go:5-80`.

```go
package core

type Format string

const (
    FormatYAML Format = "yaml"
    FormatJSON Format = "json"
)

type ParseError struct {
    Type      string
    StartByte uint
    EndByte   uint
    StartRow  uint
    StartCol  uint
    EndRow    uint
    EndCol    uint
    Text      string
}

type Issue struct {
    Rule        string
    Description string
    Row         int
}

type Fix struct {
    Rule        string
    Description string
    Before      string
    After       string
}

type Result struct {
    Format             Format
    Original           string
    Sanitized          string
    TreeText           string
    OriginalTreeText   string
    Errors             []ParseError
    OriginalErrors     []ParseError
    LintIssues         []Issue
    OriginalLintIssues []Issue
    Fixes              []Fix
    ParseClean         bool
    LintClean          bool
}

type Engine interface {
    Format() Format
    ParseTree(src string) (string, []ParseError, error)
    Lint(src string, cfg Config) []Issue
    ApplyFixes(src string, errs []ParseError, issues []Issue, cfg Config) (string, []Fix)
    Examples() []Example
}

func Run(engine Engine, src string, opts ...Option) Result
```

### Breaking-change stance

There are no existing users to preserve, so the implementation should prefer a clean public API:

- move shared runtime concepts out of `pkg/yaml`
- redesign the CLI around Glazed instead of preserving the current flag surface
- rename request fields and command structures directly instead of shipping transitional aliases

This reduces code, documentation burden, and test matrix size.

### JSON engine design

The JSON engine should start conservative. JSON is stricter than YAML, so aggressive heuristic fixing is riskier.

Recommended MVP rule set:

1. `trailing_comma`
2. `single_quoted_string`
3. `unquoted_object_key`
4. `comment_in_json`
5. `missing_comma_between_members` only if it can be implemented safely

Recommended MVP parse strategy:

- Use `encoding/json` for authoritative parse validity and error messages.
- Use `tree-sitter-json` only if a parse tree is needed for UI parity.

This means parse validity and tree visualization can be decoupled if needed:

- `encoding/json` answers: "is this valid JSON?"
- tree-sitter answers: "what tree should the UI render?"

If introducing `tree-sitter-json` is too much for the first increment, JSON support can ship with empty `TreeText` in CLI mode and no tree panel in the UI until phase 2. That is acceptable if documented.

### CLI contract and Glazed command tree

The current `--json` flag means "emit output as JSON." That name becomes ambiguous as soon as JSON input is supported, so the plan should remove it rather than preserve it.

Recommended Glazed/Cobra shape:

```text
sanitize parse    --format yaml file.yaml
sanitize lint     --format yaml file.yaml
sanitize fix      --format yaml file.yaml
sanitize sanitize --format json broken.json
```

Recommended flag model:

- `--format auto|yaml|json`
- `--output text|json|yaml|table` through Glazed output sections
- command-specific fields defined with `fields.New(...)`
- command settings surfaced via `cli.NewCommandSettingsSection()`

Recommended directory layout:

```text
cmd/sanitize/
  main.go
  cmds/
    parse/
      root.go
      run.go
    lint/
      root.go
      run.go
    fix/
      root.go
      run.go
    sanitize/
      root.go
      run.go
```

Implementation notes from the Glazed conventions:

- each verb should define a Glazed command struct embedding `*cmds.CommandDescription`
- each verb should decode a settings struct from `schema.DefaultSlug`
- root wiring should use `cli.BuildCobraCommandFromCommand(...)`
- the root command should wire help and logging once, not per verb

Keep default format detection simple:

- `.yaml` / `.yml` => YAML
- `.json` => JSON
- stdin with no filename => require `--format` or default explicitly in one place

### HTTP API contract

The server should stop speaking in YAML-specific payloads.

Recommended request:

```json
{
  "format": "yaml",
  "input": "name:Alice\n",
  "options": {
    "max_iterations": 10,
    "tab_width": 2,
    "rules": ["missing_space_after_colon"]
  }
}
```

Recommended response:

```json
{
  "format": "yaml",
  "sanitized": "name: Alice\n",
  "parse_clean": true,
  "lint_clean": true,
  "errors": [],
  "lint_issues": [],
  "fixes": []
}
```

Server contract strategy:

- accept only `format` + `input`
- return `400` for missing or unknown format
- update the frontend and docs at the same time rather than carrying legacy request names

### UI model

The frontend in `cmd/sanitize-server/static/index.html:214-267` currently assumes a single format. For JSON support it needs:

1. A format selector.
2. Format-specific examples.
3. Different labels for "Input" and "Sanitized output" when JSON is selected.
4. Conditional tree behavior if JSON tree output is absent in the first phase.

## Design Decisions

### Decision: extract shared orchestration before adding JSON rules

Rationale:

The sanitize loop in `pkg/yaml/sanitize.go:18-79` is already a reusable algorithm. Duplicating it into `pkg/json` would create two subtly diverging runtimes and make future fixes harder.

### Decision: break compatibility where it reduces scaffolding

Rationale:

There are no users yet. Transitional wrappers, aliases, and duplicate contracts would add code and tests without creating user value.

### Decision: make format an explicit input concept

Rationale:

The current CLI already uses `--json` for output formatting. Preserving that flag would force a confusing distinction between "JSON output" and "JSON input". A breaking rename is cleaner.

### Decision: move to a Glazed-backed CLI while reshaping the verbs

Rationale:

The current CLI is small enough that this is the right time to adopt `cmds.CommandDescription`, `fields.New`, Glazed output sections, and structured Cobra groups. Doing this during JSON work avoids migrating the CLI twice.

### Decision: start JSON support with conservative fixers

Rationale:

YAML is designed for hand-authored, messy text. JSON usually exists in stricter machine pipelines. Conservative fixes minimize the risk of silently changing intended semantics.

## Current-State Analysis

### Current runtime flow

The current YAML runtime works like this:

```text
input
  -> ParseTree(original)
  -> Lint(original)
  -> repeat up to maxIterations
       -> ParseTree(current)
       -> Lint(current)
       -> applyFixes(current)
  -> Result
```

Evidence:

- `pkg/yaml/sanitize.go:14-18` captures original state.
- `pkg/yaml/sanitize.go:19-63` runs the iterative loop.
- `pkg/yaml/fix.go:17-64` applies line-level and document-level fixes.

### Current architectural boundaries

Library:

- `pkg/yaml/types.go` defines exported data structures.
- `pkg/yaml/options.go` defines functional options.
- `pkg/yaml/parse.go` owns tree-sitter integration.
- `pkg/yaml/lint.go` owns heuristic issue detection.
- `pkg/yaml/fix.go` owns rule implementations.

Entry points:

- `cmd/sanitize/main.go` is the CLI.
- `cmd/sanitize-server/main.go` is the HTTP server.
- `cmd/sanitize-server/static/index.html` is the embedded UI.

### Constraints that shape JSON support

1. The codebase is small, so a full abstraction layer is feasible.
2. The package already uses tree-sitter for YAML parse trees.
3. The UI consumes `tree_text`, `errors`, `original_errors`, `lint_issues`, and `fixes`.
4. There is no external-user compatibility requirement, so package and CLI cleanup can be folded into the same change set.

## Gap Analysis

### Gap 1: types are YAML-package-local

`pkg/yaml/types.go:13-60` defines general-purpose concepts such as result, issue, and fix, but they are trapped inside the YAML package.

### Gap 2: the CLI has no format model

`cmd/sanitize/main.go:14-18` exposes flags for output JSON, lint-only mode, tab width, and max iterations, but nothing about input format.

### Gap 3: the HTTP contract is YAML-specific

`cmd/sanitize-server/main.go:56-58` and `cmd/sanitize-server/main.go:78-80` declare request bodies with a `YAML` field, which makes JSON support awkward and semantically wrong.

### Gap 4: the UI text and examples are hardcoded to YAML

The static page title, button labels, placeholders, and example loader all assume YAML.

## Pseudocode And Key Flows

### Shared sanitize loop

```go
func Run(engine Engine, src string, opts ...Option) Result {
    cfg := defaultConfig()
    for _, opt := range opts {
        opt(&cfg)
    }

    original := src
    origTree, origErrs, _ := engine.ParseTree(original)
    origIssues := engine.Lint(original, cfg)

    for iter := 0; iter < cfg.maxIterations; iter++ {
        tree, errs, err := engine.ParseTree(src)
        if err != nil {
            break
        }

        issues := engine.Lint(src, cfg)
        if len(errs) == 0 && len(issues) == 0 {
            return buildResult(engine.Format(), original, src, tree, origTree, errs, origErrs, issues, origIssues, allFixes)
        }

        next, fixes := engine.ApplyFixes(src, errs, issues, cfg)
        if len(fixes) == 0 {
            return finalizeCurrentState(...)
        }

        src = next
        allFixes = append(allFixes, fixes...)
    }

    return finalizeCurrentState(...)
}
```

### CLI flow with explicit format

```text
parse flags
  -> resolve format
  -> choose engine
  -> read input
  -> run sanitize or lint-only
  -> render text or JSON output
  -> set exit code from cleanliness, not from output mode
```

### Server flow

```text
POST /api/sanitize
  -> decode request {format,input,options}
  -> resolve engine
  -> run sanitize
  -> encode Result
```

## Alternatives Considered

### Alternative 1: duplicate `pkg/yaml` into `pkg/json`

Rejected because the sanitize loop, result types, option types, CLI handling, and server wiring would be duplicated immediately.

### Alternative 2: add JSON support only in the CLI by converting JSON to YAML internally

Rejected because JSON and YAML have different failure modes and different safe-fix sets. Converting formats would blur semantics and make the server/UI model harder to reason about.

### Alternative 3: use only `encoding/json` and give up parse-tree support

Possible for an MVP, but rejected as the long-term design because the server UI is built around tree display. This remains a valid phase-1 compromise if tree-sitter JSON integration would delay the more important API work.

### Alternative 4: introduce a registry and plugin architecture immediately

Rejected as over-engineered for a two-format codebase. A simple `switch` on format or a small engine map is enough.

## Implementation Plan

### Phase 1: release hardening before expansion

1. Fix duplicate-key detection.
2. Fix CLI exit-code behavior on dirty input.
3. Harden `sanitize-server` with timeouts and bounded request sizes.
4. Resolve the `gosec` parse warning.
5. Add regression, CLI, and HTTP tests for those issues.

### Phase 2: extract the shared core

1. Move common types from `pkg/yaml/types.go` into `pkg/core/types.go`.
2. Move option/config logic from `pkg/yaml/options.go` into `pkg/core/options.go`.
3. Move the generic sanitize loop from `pkg/yaml/sanitize.go` into `pkg/core/run.go`.
4. Add a YAML engine that reuses the existing parse/lint/fix functions.
5. Remove or rename YAML-specific public shapes as needed instead of preserving wrappers.

### Phase 3: replace the CLI with Glazed/Cobra commands

1. Design the verb tree (`parse`, `lint`, `fix`, `sanitize`).
2. Implement each verb as a Glazed command with `CommandDescription`, `fields.New`, and decoded settings structs.
3. Build the root/group Cobra tree with `cli.BuildCobraCommandFromCommand(...)`.
4. Add help and logging setup at the root.
5. Remove the old single-file CLI path rather than preserving alias flags.

### Phase 4: add JSON parse and lint support

1. Create `pkg/json`.
2. Implement JSON parse validation.
3. Define JSON issue and fix rules.
4. Add examples and tests.
5. Ship a read-only JSON mode first if fixers are not all ready.

### Phase 5: wire the server and UI

1. Change request models to `format` + `input`.
2. Add format selector and examples in the UI.
3. Support format-aware labels and rendering.
4. Add `httptest` coverage for both formats.

### Phase 6: publishable cleanup

1. Update README examples and install docs to the new CLI shape.
2. Re-run validation and capture clean results.
3. Re-evaluate release readiness after the JSON MVP lands.

## Open Questions

1. Should JSON support include comment stripping, or should comments be treated as a hard error because they are not valid JSON?
2. Should auto-detection inspect file contents when reading from stdin, or should stdin require an explicit `--format`?
3. Is parse-tree parity required for JSON in the first release, or can the UI degrade gracefully?

## Testing Strategy

1. Preserve all existing YAML tests.
2. Add a parallel `pkg/json/sanitize_test.go`.
3. Add cross-format table tests at the shared-core level.
4. Add CLI subprocess tests for verb wiring, output mode, and exit code behavior.
5. Add HTTP tests for request decoding, invalid format handling, and JSON/YAML success cases.
6. Add regression fixtures for every fixer rule to avoid silent semantic drift.

## References

- `README.md:1-82`
- `cmd/sanitize/main.go:13-70`
- `cmd/sanitize-server/main.go:23-99`
- `cmd/sanitize-server/static/index.html:214-552`
- `pkg/yaml/options.go:3-54`
- `pkg/yaml/types.go:13-60`
- `pkg/yaml/parse.go:20-79`
- `pkg/yaml/lint.go:19-110`
- `pkg/yaml/fix.go:16-316`
- `pkg/yaml/sanitize.go:5-80`
- `ttmp/2026/03/13/SANITIZE-001--turn-yaml-sanitizer-into-reusable-go-go-golems-sanitize-package/design-doc/01-implementation-plan.md:1-82`
