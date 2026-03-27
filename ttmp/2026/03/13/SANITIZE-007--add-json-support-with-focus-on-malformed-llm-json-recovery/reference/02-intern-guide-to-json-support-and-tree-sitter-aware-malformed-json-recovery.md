---
Title: Intern guide to JSON support and tree-sitter-aware malformed JSON recovery
Ticket: SANITIZE-007
Status: active
Topics:
    - json
    - linting
    - api-design
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: pkg/yaml/analysis.go
      Note: Current shared analysis pattern that JSON should likely mirror
    - Path: pkg/yaml/lint.go
      Note: Current lint assembly flow and issue model reference
    - Path: pkg/yaml/fix.go
      Note: Current fixer orchestration reference
    - Path: pkg/yaml/options.go
      Note: Current validated rule-selection model reference
    - Path: internal/cli/commands.go
      Note: Existing Glazed command surface that JSON will need to integrate with
    - Path: internal/server/server.go
      Note: Current HTTP API and embedded static UI entrypoint that must become format-aware
    - Path: internal/server/static/index.html
      Note: Current single-page browser shell that is hard-coded to YAML terminology
    - Path: internal/server/static/js/app.js
      Note: Current browser-side analysis flow and render logic for parse tree, issues, and fixes
    - Path: internal/server/static/css/style.css
      Note: Current UI layout and styling that the JSON playground will extend
    - Path: ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/scripts/json-tree-sitter-parse/main.go
      Note: Standalone tree-sitter JSON parse inspector
    - Path: ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/scripts/json-error-matrix/main.go
      Note: Replication matrix generator for malformed JSON cases
    - Path: ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/scripts/json-heuristic-probe/main.go
      Note: Heuristic probe for malformed JSON cases
    - Path: ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/01-json-parse-error-replication-matrix.md
      Note: Generated evidence comparing strict parsing and tree-sitter
    - Path: ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/02-json-heuristic-probe.md
      Note: Generated evidence showing heuristic-detectable malformed cases
ExternalSources: []
Summary: Detailed intern-oriented guide for understanding the current sanitize architecture and implementing JSON support focused on malformed LLM-generated JSON.
LastUpdated: 2026-03-13T23:45:00-04:00
WhatFor: Explain the current system, the JSON implementation that shipped in SANITIZE-007, and the remaining conservative boundaries for a new contributor.
WhenToUse: Use when onboarding to SANITIZE-007 or when extending pkg/json, the JSON playground, or the experiment/reporting scripts.
---

# Intern guide to JSON support and tree-sitter-aware malformed JSON recovery

## Goal

This guide explains the current `sanitize` system and how to extend it with JSON support that is useful for messy LLM output, not just perfect strict JSON. It is written for a new intern who needs to understand:

- what the existing codebase already does,
- which architectural patterns should be copied from YAML,
- what research and experiments already exist in `SANITIZE-007`,
- where tree-sitter helps,
- where heuristics are still necessary,
- and how to turn the current ticket into working package and CLI code.

## Context

The repository started as a YAML sanitizer and now supports both YAML and JSON. `SANITIZE-007` is the ticket that added first-release JSON support with a deliberately conservative repair boundary and a full stack of supporting evidence documents and scripts.

The important context is that the project is not just a parser wrapper. It tries to do three related jobs:

- inspect malformed structured text,
- explain what is wrong in a way that is useful to users,
- and repair a subset of common issues conservatively.

Today the project has:

- a shared analysis path,
- a lint issue model,
- a fix model,
- iterative sanitize behavior,
- a Glazed-based CLI,
- a format-aware HTTP API,
- and a shared YAML/JSON browser playground.

For JSON, the ticket research shows that strict parsing plus tree-sitter plus heuristics is likely the right blend. Strict parsing alone is too brittle for LLM-generated output, but heuristics alone are too weak for structural localization.

## Quick Reference

### High-level system map

```text
User Input
   |
   v
CLI / API surface
   |
   v
Format-specific package
   |
   +--> parse / analyze
   +--> lint issues
   +--> suggested or automatic fixes
   |
   v
Result object returned to caller
```

### Current code areas to know first

- Current YAML analysis: [analysis.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/analysis.go)
- Current YAML lint engine: [lint.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/lint.go)
- Current YAML fix engine: [fix.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/fix.go)
- Current YAML options/rules: [options.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/options.go) and [rules.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/rules.go)
- Current sanitize orchestration: [sanitize.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/sanitize.go)
- Current Glazed CLI: [commands.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/internal/cli/commands.go)
- Current HTTP server: [server.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/internal/server/server.go)
- Current browser shell: [index.html](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/internal/server/static/index.html)
- Current browser logic: [app.js](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/internal/server/static/js/app.js)
- Current browser styling: [style.css](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/internal/server/static/css/style.css)

### JSON research files already created in this ticket

- Imported malformed case list: [03-json-parse-errors-import.md](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/local/03-json-parse-errors-import.md)
- Tree-sitter replication matrix: [01-json-parse-error-replication-matrix.md](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/01-json-parse-error-replication-matrix.md)
- Heuristic probe: [02-json-heuristic-probe.md](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/02-json-heuristic-probe.md)
- Example metadata: [03-json-example-metadata.json](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/03-json-example-metadata.json)
- Detection buckets: [04-json-detection-buckets.md](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/04-json-detection-buckets.md)
- Repair matrix: [05-json-repair-matrix.md](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/05-json-repair-matrix.md)
- Rule matrix: [06-json-rule-matrix.json](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/06-json-rule-matrix.json)
- Overlap study: [07-json-overlap-study.md](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/07-json-overlap-study.md)
- Standalone parse inspector: [json-tree-sitter-parse/main.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/scripts/json-tree-sitter-parse/main.go)
- Case matrix generator: [json-error-matrix/main.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/scripts/json-error-matrix/main.go)
- Heuristic probe generator: [json-heuristic-probe/main.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/scripts/json-heuristic-probe/main.go)
- Example metadata exporter: [json-example-metadata/main.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/scripts/json-example-metadata/main.go)
- Detection bucket generator: [json-detection-buckets/main.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/scripts/json-detection-buckets/main.go)
- Repair matrix generator: [json-repair-matrix/main.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/scripts/json-repair-matrix/main.go)
- Rule matrix generator: [json-rule-matrix/main.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/scripts/json-rule-matrix/main.go)
- Overlap study generator: [json-overlap-study/main.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/scripts/json-overlap-study/main.go)

## Current shipped JSON boundary

The first-release JSON engine is intentionally conservative.

It auto-fixes:

- Markdown fence wrappers
- leading or trailing prose wrappers
- Python literals such as `True`, `False`, and `None`
- comments
- duplicate commas
- trailing commas

It lints but does not auto-fix:

- single quotes
- unquoted keys
- missing commas
- missing colons
- missing closing delimiters
- multiple top-level values
- duplicate keys

That means the right mental model is "recovery for low-ambiguity LLM noise, lint for ambiguous structural repair."

## What the system is

At a product level, `sanitize` is a structured-text repair tool. It is not trying to be a general IDE parser and it is not trying to become a full schema validator. The core contract is:

- take malformed or inconsistent input,
- surface parse and lint findings,
- apply a conservative subset of repairs,
- return both the repaired text and the evidence describing what happened.

The current YAML implementation already embodies this philosophy. JSON support should follow the same mental model while adjusting for JSON's stricter syntax and the much more repetitive failure patterns that show up in LLM-generated output.

### Core data flow

```text
raw source
   |
   v
analyze document
   |
   +--> structural parser output
   +--> parse error spans
   +--> any extra format-specific structural facts
   |
   v
lint assembly
   |
   +--> parser-derived issues
   +--> heuristic issues
   |
   v
fix pass
   |
   v
final result
```

For JSON, format-specific structural facts will likely include:

- duplicate object keys,
- top-level value kind,
- any tree-sitter `ERROR` or `MISSING` nodes,
- maybe extracted spans for comments or prose-wrapped regions if pre-parse recovery is added.

## How the current YAML architecture works

Before you build `pkg/json`, you should understand the patterns that already exist in `pkg/yaml`.

### Shared analysis

The YAML package already centralizes parser-derived state in [analysis.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/analysis.go). The important idea is that `ParseTree`, `Lint`, and `Sanitize` should not all parse the same document independently if they can share one internal analysis object.

Conceptually:

```go
type documentAnalysis struct {
    TreeText      string
    ParseErrors   []ErrorNode
    DuplicateKeys []duplicateKeyOccurrence
    LineIndex     lineIndex
}
```

JSON should probably mirror this shape:

```go
type documentAnalysis struct {
    TreeText      string
    ParseErrors   []ErrorNode
    StrictError   error
    DuplicateKeys []duplicateKeyOccurrence
    Signals       heuristicSignals
}
```

Not all of those fields need to exist in the first implementation. The point is that all downstream behavior should depend on a single analysis product.

### Lint assembly

In YAML, [lint.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/lint.go) merges:

- parser-derived issues,
- heuristic issues,
- rule filtering.

JSON should preserve the same separation:

- parser-derived structural issues:
  - mismatched delimiters
  - truncated objects or arrays
  - invalid token spans
- heuristic issues:
  - Markdown fences
  - prose before or after the JSON
  - Python literals
  - comments
  - duplicate commas
  - unquoted keys

### Fix engine

In YAML, [fix.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/fix.go) works by applying one conservative fix round at a time, then re-analyzing. JSON should likely use the same iterative pattern because many malformed inputs become easier to interpret after one small repair.

Example:

```text
```json
{"ok": True,}
```
```

A reasonable JSON repair sequence could be:

1. remove Markdown fence
2. normalize `True -> true`
3. remove trailing comma

Trying to solve all of that in one giant heuristic is harder to reason about and harder to test.

### Rule filtering and CLI integration

Recent work in the YAML package added a rule registry and validated rule selection in [rules.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/rules.go) and [options.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/options.go).

JSON should reuse that pattern:

- every rule has a canonical name,
- the package validates rule names,
- the CLI consumes package metadata rather than inventing rule strings on its own.

## What the JSON research already proves

`SANITIZE-007` is not starting from zero. The scripts and generated reports already answer some architectural questions.

### Finding 1: Tree-sitter is useful, but not sufficient

The replication matrix in [01-json-parse-error-replication-matrix.md](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/01-json-parse-error-replication-matrix.md) shows that tree-sitter gives useful structural localization for:

- trailing commas,
- missing commas,
- missing colons,
- mismatched delimiters,
- prose-wrapped JSON,
- unterminated strings.

But it is not enough by itself for every LLM failure mode:

- comments may parse as comment nodes depending on the grammar behavior rather than as plain errors,
- Python literals need semantic normalization,
- Markdown fences are not JSON syntax in the first place,
- prose wrapping requires extraction or boundary heuristics,
- invalid escapes may parse structurally even when strict `encoding/json` semantics differ.

That means the JSON engine should not be just `parse with tree-sitter and report errors`. It needs at least a hybrid model.

### Finding 2: Heuristics catch some high-value cases cheaply

The heuristic probe in [02-json-heuristic-probe.md](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/02-json-heuristic-probe.md) shows obvious easy wins:

- Markdown fences
- comments
- single quotes
- Python literals
- trailing commas
- duplicate commas
- prose wrapping
- multiple top-level objects
- unquoted keys

These should likely become first-wave JSON lint or fix rules.

### Finding 3: Some cases are structurally ambiguous

The imported source in [03-json-parse-errors-import.md](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/local/03-json-parse-errors-import.md) includes cases where a silent rewrite would be risky:

- missing commas
- truncated output
- malformed strings
- duplicate keys with unclear intended precedence
- placeholder fragments

These should likely be lint-only or suggestion-only in the first implementation.

## Proposed JSON package shape

The cleanest shape is a new `pkg/json` package that mirrors the successful parts of `pkg/yaml` without forcing the two engines to share too much code too early.

### Suggested package structure

```text
pkg/json/
  analysis.go
  parse.go
  lint.go
  fix.go
  sanitize.go
  types.go
  options.go
  rules.go
  duplicate_keys.go
  heuristics.go
```

### Why not over-share immediately

You may be tempted to factor a common `pkg/core` right away. Do not do that first unless duplication becomes obvious and stable.

Reasons:

- YAML and JSON have different grammar and recovery properties.
- YAML is indentation-sensitive and far more permissive.
- JSON is token-rigid but very repetitive in LLM failure modes.
- Premature abstraction could hide important differences and slow the ticket.

A better plan is:

1. get `pkg/json` working with a similar external shape,
2. compare the YAML and JSON engines,
3. extract shared code only after the similarities are proven.

## Suggested result and issue model

JSON should probably reuse the current shape of `ErrorNode`, `LintIssue`, `Fix`, and `Result` from [types.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/types.go), or define equivalent JSON-specific versions first and consolidate later.

The important fields are:

- rule name
- source of the issue:
  - parse
  - heuristic
  - maybe strict-parser
- byte range
- row and column range
- description

### Candidate JSON rules

- `trailing_comma`
- `single_quoted_string`
- `unquoted_key`
- `python_literal`
- `markdown_fence`
- `comment`
- `prose_wrapped_json`
- `multiple_top_level_values`
- `duplicate_key`
- `missing_comma`
- `missing_colon`
- `unterminated_string`
- `mismatched_delimiter`
- `invalid_numeric_literal`

Not all of these should be auto-fixable.

### Candidate fix classification

Safe early auto-fix:

- `markdown_fence`
- `trailing_comma`
- `python_literal`
- maybe `comment`

Likely lint-only at first:

- `missing_comma`
- `missing_colon`
- `duplicate_key`
- `unterminated_string`
- `mismatched_delimiter`

Suggestion-only or delayed:

- prose extraction if there are multiple plausible JSON spans,
- truncated completion repair,
- semantic duplicate-key rewrites.

## Proposed analyze-lint-fix flow for JSON

### End-to-end flow

```text
raw input
   |
   v
pre-analysis heuristics
   |
   +--> fence, prose, and comment detection
   |
   v
tree-sitter JSON parse
   |
   +--> structural errors
   +--> tree text
   +--> object traversal facts
   |
   v
strict parsing check
   |
   +--> encoding/json semantic failures
   |
   v
lint assembly
   |
   +--> parser-derived
   +--> strict-parser-derived
   +--> heuristic-derived
   |
   v
safe fix pass
   |
   v
re-analyze
```

### Pseudocode

```go
func Analyze(src string) (documentAnalysis, error) {
    signals := detectHeuristicSignals(src)
    tree := parseWithTreeSitterJSON(src)
    strictErr := parseWithEncodingJSON(src)
    duplicateKeys := collectDuplicateKeys(tree)

    return documentAnalysis{
        TreeText:      tree.Sexp,
        ParseErrors:   tree.Errors,
        StrictError:   strictErr,
        DuplicateKeys: duplicateKeys,
        Signals:       signals,
    }, nil
}

func LintWithOptions(src string, opts ...Option) ([]LintIssue, error) {
    cfg := validateOptions(opts)
    doc := Analyze(src)

    issues := []LintIssue{}
    issues = append(issues, lintFromTreeErrors(doc, cfg)...)
    issues = append(issues, lintFromStrictParser(doc, cfg)...)
    issues = append(issues, lintFromHeuristics(doc, cfg)...)
    return issues, nil
}

func SanitizeWithOptions(src string, opts ...Option) (Result, error) {
    cfg := validateOptions(opts)
    current := src
    fixes := []Fix{}

    for i := 0; i < cfg.MaxIterations; i++ {
        doc := Analyze(current)
        next, roundFixes := applySafeFixes(current, doc, cfg)
        if len(roundFixes) == 0 {
            break
        }
        current = next
        fixes = append(fixes, roundFixes...)
    }

    final := Analyze(current)
    return buildResult(src, current, final, fixes), nil
}
```

## API and CLI integration plan

The CLI already has a format-neutral structure in [commands.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/internal/cli/commands.go). JSON support can be integrated in one of two ways.

### Option A: Add `--format yaml|json`

Pros:

- one binary surface
- easy mental model

Cons:

- command code gets more conditional logic

### Option B: Add explicit subcommands later

Example:

```text
sanitize lint --format json
sanitize fix --format json
sanitize parse --format json
```

This is probably the simplest first step because it fits the existing `lint`, `fix`, and `parse` verbs.

### JSON-specific parse inspection

The ticket already includes a standalone parse inspector in [json-tree-sitter-parse/main.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/scripts/json-tree-sitter-parse/main.go), and the shipped CLI now exposes JSON parse inspection directly:

```text
sanitize parse --format json
sanitize lint --format json
sanitize fix --format json
sanitize rules --format json
```

## Current server and browser UI architecture

The repository already contains a browser-facing playground. This matters because JSON support is not only a package and CLI concern. `SANITIZE-007` now brings the server and browser along with the package and CLI, so the product no longer has a split personality.

The current server lives in [server.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/internal/server/server.go). It serves static assets and exposes three routes:

- `GET /api/examples`
- `POST /api/sanitize`
- `POST /api/parse`

The shipped request model is now format-aware:

```go
type analyzeRequest struct {
    Format string `json:"format"`
    Input  string `json:"input"`
}
```

That change matters because the browser, tests, and any future API clients can now speak the same contract for both formats.

### How the browser works today

The current browser shell is [index.html](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/internal/server/static/index.html). The stateful browser logic is in [app.js](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/internal/server/static/js/app.js). The styling is in [style.css](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/internal/server/static/css/style.css).

The current UI flow is:

```text
Example dropdown or textarea input
   |
   v
debounced triggerAnalysis()
   |
   v
POST /api/sanitize with { format, input }
   |
   v
renderResult(data)
   |
   +--> status badge
   +--> parse tree panel
   +--> issue list
   +--> sanitized output
   +--> applied fixes list
```

That architecture turned out to be close to what JSON needed. The important changes were:

- add a format selector
- filter examples by format
- expose original strict-parse state for JSON
- rename the app into format-neutral language
- keep original/sanitized tree toggling identical across formats

### Existing browser invariants you should preserve

- The input editor is the source of truth.
- Analysis is debounced so the UI stays responsive while typing.
- The parse tree and issue panes reflect the original input by default.
- The user can toggle the parse tree between original and sanitized output.
- Sanitized output is copyable as plain text.

JSON mode should preserve those interaction patterns because they are already coherent.

## Format-aware HTTP contract

The shipped clean request model is:

```json
{
  "format": "yaml",
  "input": "name: Alice\n"
}
```

and for JSON:

```json
{
  "format": "json",
  "input": "{\"name\": \"Alice\"}"
}
```

The response is structurally similar across formats. That does not require YAML and JSON to share the same internal package, but it does require the server to normalize the result shape enough that the browser does not need entirely separate rendering stacks.

### Server request and dispatch pseudocode

```go
type analyzeRequest struct {
    Format string `json:"format"`
    Input  string `json:"input"`
}

func sanitizeHandler(w http.ResponseWriter, r *http.Request) {
    var req analyzeRequest
    decodeJSONBody(w, r, &req)

    switch req.Format {
    case "yaml":
        encode(yamlsanitize.Sanitize(req.Input))
    case "json":
        encode(jsonsanitize.Sanitize(req.Input))
    default:
        http.Error(w, "unsupported format", http.StatusBadRequest)
    }
}
```

### Why the API contract matters for the UI

If the server keeps a YAML-shaped request body like `{ "yaml": "..." }`, then every client ends up doing awkward remapping. The browser cannot easily support multiple formats cleanly, and tests become harder to read. A format-aware request body lets the same browser action dispatch to YAML or JSON without branching all over the rendering logic.

## JSON parse playground in the browser UI

The JSON work now includes a real parse playground mode in the embedded browser app. This is not a decorative addition. It is one of the best ways to inspect malformed LLM JSON and verify whether tree-sitter spans and lint rules line up with user expectations.

### Product goal for the playground

The user should be able to paste messy LLM output like this:

```text
Here is your JSON:
```json
{"ok": True,}
```
```

and immediately see:

- that the original payload is not strict JSON,
- what the structural parse issues are,
- which heuristic rules fired,
- which recovery fixes were applied,
- and whether the sanitized result is now strict valid JSON.

### UI elements that shipped

The current `index.html` evolved rather than being replaced. The main shipped changes are:

- rename the application to `Sanitize Playground`
- add a format picker next to the example picker
- make the example picker format-aware
- keep format-specific panel labels such as `Input YAML` and `Input JSON`
- add a strict-parse status indicator in the parse panel for JSON mode

### Toolbar layout

```text
[Format: YAML | JSON] [Example dropdown filtered by format] [status badge] [spinner]
```

### Panel behavior in JSON mode

- Input panel mentions malformed or prose-wrapped JSON in the placeholder.
- Example picker loads JSON fixtures only when JSON is selected.
- Parse tree panel shows tree-sitter errors and strict JSON status separately.
- Issues panel distinguishes parse, strict-parser, and heuristic signals.
- Sanitized output panel shows final text plus applied fixes and final strict-clean state.

### Browser-state model

The current [app.js](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/internal/server/static/js/app.js) still uses lightweight globals, but JSON mode added explicit format state and format-filtered example behavior.

```js
const state = {
  format: "yaml",
  examples: [],
  examplesByFormat: {
    yaml: [],
    json: [],
  },
  lastResult: null,
  showingOriginalTree: true,
};
```

### Browser analysis pseudocode

```js
async function triggerAnalysis() {
  const input = editor.value;
  if (!input.trim()) {
    resetUI();
    return;
  }

  const res = await fetch("/api/sanitize", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      format: state.format,
      input,
    }),
  });

  renderResult(await res.json(), state.format);
}
```

### `renderResult` contract

The browser should avoid format-specific rendering branches except for labels and a few JSON-only badges. Aim for a common result interface like:

```json
{
  "format": "json",
  "tree_text": "(document ...)",
  "original_tree_text": "(document ...)",
  "errors": [],
  "original_errors": [],
  "lint_issues": [],
  "original_lint_issues": [],
  "fixes": [],
  "sanitized": "{\"ok\": true}",
  "parse_clean": true,
  "lint_clean": true,
  "strict_parse_clean": true
}
```

JSON now exposes both `strict_parse_clean` and `original_strict_parse_clean`. YAML does not need to grow those fields; the browser treats them as JSON-only optional state.

### Suggested UI diagram

```text
+---------------------------------------------------------------+
| sanitize playground                                           |
| format: [json]   example: [llm prose wrapper]   status: lint |
+----------------------+----------------------+-----------------+
| Input                | Parse Tree           | Sanitized       |
|                      |                      | Output          |
| Here is your JSON:   | (document            | {"ok": true}    |
| ```json              |   (ERROR ...)        |                 |
| {"ok": True,}        | )                    | fixes:          |
| ```                  | strict parse: fail   | - strip fence   |
|                      | errors: 2            | - True->true    |
|                      |                      | - rm comma      |
+----------------------+----------------------+-----------------+
| Issues: parse error · strict parser · heuristic lint          |
+---------------------------------------------------------------+
```

### UI implementation order

Do not start by rewriting the whole SPA. The better sequence is:

1. change the API contract to accept `format` and `input`
2. add format metadata to examples
3. add a format picker to the toolbar
4. make the existing render path format-aware
5. add JSON-only badges and labels
6. then add richer issue grouping if needed

That sequence reduces the chance of mixing product redesign with transport-contract changes.

## Implementation phases for an intern

### Phase 1: Build `pkg/json` analysis and parse reporting

Deliverables:

- `pkg/json/parse.go`
- `pkg/json/analysis.go`
- `pkg/json/types.go`

Tasks:

- wire `tree-sitter-json`
- collect `ERROR` and `MISSING` nodes
- expose `ParseTree(src)` for JSON
- add tests for malformed imported cases

Validation:

- parse each case from [03-json-parse-errors-import.md](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/local/03-json-parse-errors-import.md)
- confirm spans resemble the replication matrix

### Phase 2: Add first-wave heuristic lint rules

Deliverables:

- `pkg/json/lint.go`
- `pkg/json/heuristics.go`
- `pkg/json/rules.go`

First rules:

- `markdown_fence`
- `trailing_comma`
- `python_literal`
- `comment`
- `single_quoted_string`
- `prose_wrapped_json`

Validation:

- compare lint output against [02-json-heuristic-probe.md](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/02-json-heuristic-probe.md)

### Phase 3: Add safe fixers

Deliverables:

- `pkg/json/fix.go`
- `pkg/json/sanitize.go`

Start with:

- remove Markdown fences
- normalize Python literals
- remove trailing commas

Do not start with ambiguous structure repair.

Validation:

- table-driven tests for each fix
- multiple-iteration tests for combinations like fence plus trailing comma plus python literal

### Phase 4: Wire CLI support

Deliverables:

- format selection in [commands.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/internal/cli/commands.go)
- CLI tests in [main_test.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/cmd/sanitize/main_test.go)

Tasks:

- add JSON mode to `parse`, `lint`, and `fix`
- add JSON rules to rule listing
- add JSON-specific examples or fixtures

### Phase 5: Make the HTTP server format-aware

Deliverables:

- updated request model in [server.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/internal/server/server.go)
- format-aware `/api/examples`
- format-aware `/api/sanitize`
- format-aware `/api/parse`

Tasks:

- replace the YAML-only request body with `{ format, input }`
- dispatch to YAML or JSON packages by format
- include format metadata in example payloads
- add server tests for YAML and JSON flows

Validation:

- `go test ./internal/server`
- manual request checks with `curl`

### Phase 6: Add the JSON parse playground to the browser UI

Deliverables:

- updated [index.html](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/internal/server/static/index.html)
- updated [app.js](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/internal/server/static/js/app.js)
- updated [style.css](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/internal/server/static/css/style.css)

Tasks:

- add a format picker
- filter examples by format
- rename YAML-specific labels to generic labels
- preserve the original/sanitized tree toggle in JSON mode
- show strict JSON parse status
- keep issue and fix rendering format-aware but mostly shared

Validation:

- open `sanitize serve`
- exercise YAML mode and JSON mode manually
- verify that malformed LLM JSON examples render useful parse and fix feedback

## Review and testing checklist

When you implement code for this ticket, use this review order:

1. analysis and parse layer
2. issue model
3. heuristic rules
4. fix safety
5. CLI wiring

Commands you should run often:

```bash
go test ./pkg/yaml ./cmd/sanitize ./internal/...
cd ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/scripts
go run ./json-error-matrix > ../sources/01-json-parse-error-replication-matrix.md
go run ./json-heuristic-probe > ../sources/02-json-heuristic-probe.md
docmgr doctor --ticket SANITIZE-007 --stale-after 30
```

## Common mistakes to avoid

- Do not assume tree-sitter and `encoding/json` mean the same thing for every malformed case.
- Do not auto-fix ambiguous structural errors just because you can guess a likely repair.
- Do not add CLI-specific rule names that are not registered in the package.
- Do not over-abstract YAML and JSON too early.
- Do not throw away imported malformed examples; keep the experiment corpus as a regression asset.

## Recommended first tasks for a new intern

- Read [analysis.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/analysis.go), [lint.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/lint.go), and [sanitize.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/sanitize.go).
- Run the JSON experiment scripts from `SANITIZE-007/scripts/`.
- Read the generated matrix and heuristic probe outputs.
- Implement `pkg/json/parse.go` and `pkg/json/analysis.go`.
- Add tests for the first three easy recovery cases:
  - Markdown fences
  - trailing commas
  - Python literals

## Usage examples

### Example mental model

If the input is:

```text
```json
{"ok": True,}
```
```

Then the engine should think:

- this is formatting pollution plus almost-JSON
- tree-sitter can parse some structure, but the backticks are outside JSON
- heuristics should strip fences first
- then normalize Python literals
- then remove the trailing comma
- then re-parse strictly

### Example review question

When you add a new fixer, ask:

```text
Is this repair deterministic enough that two reasonable humans would choose the same rewrite?
```

If the answer is no, it probably should not be an automatic fix yet.

## Usage Examples

Use this guide as the onboarding document before touching code for `SANITIZE-007`.

Suggested reading order:

1. [01-json-support-outline-and-malformed-llm-json-error-taxonomy.md](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/design-doc/01-json-support-outline-and-malformed-llm-json-error-taxonomy.md)
2. [01-common-json-parse-errors-from-llm-output.md](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/reference/01-common-json-parse-errors-from-llm-output.md)
3. [01-json-parse-error-replication-matrix.md](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/01-json-parse-error-replication-matrix.md)
4. [02-json-heuristic-probe.md](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/02-json-heuristic-probe.md)
5. [tasks.md](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/tasks.md)
6. [03-diary.md](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/reference/03-diary.md)

## Related

- [01-json-support-outline-and-malformed-llm-json-error-taxonomy.md](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/design-doc/01-json-support-outline-and-malformed-llm-json-error-taxonomy.md)
- [01-common-json-parse-errors-from-llm-output.md](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/reference/01-common-json-parse-errors-from-llm-output.md)
- [01-json-parse-error-replication-matrix.md](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/01-json-parse-error-replication-matrix.md)
- [02-json-heuristic-probe.md](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/02-json-heuristic-probe.md)
- [03-diary.md](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/reference/03-diary.md)
