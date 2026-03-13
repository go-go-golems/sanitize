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
LastUpdated: 2026-03-13T20:05:00-04:00
WhatFor: Explain the current system, the JSON research performed so far, and a concrete implementation path for a new contributor.
WhenToUse: Use when onboarding to SANITIZE-007 or when implementing pkg/json, JSON lint rules, fixers, and CLI integration.
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

The repository started as a YAML sanitizer and now has a more structured internal architecture. The current production code is still YAML-only, but `SANITIZE-007` adds research artifacts for JSON support.

The important context is that the project is not just a parser wrapper. It tries to do three related jobs:

- inspect malformed structured text,
- explain what is wrong in a way that is useful to users,
- and repair a subset of common issues conservatively.

For YAML, the project already has:

- a shared analysis path,
- a lint issue model,
- a fix model,
- iterative sanitize behavior,
- and a Glazed-based CLI.

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

### JSON research files already created in this ticket

- Imported malformed case list: [03-json-parse-errors-import.md](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/local/03-json-parse-errors-import.md)
- Tree-sitter replication matrix: [01-json-parse-error-replication-matrix.md](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/01-json-parse-error-replication-matrix.md)
- Heuristic probe: [02-json-heuristic-probe.md](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/02-json-heuristic-probe.md)
- Standalone parse inspector: [json-tree-sitter-parse/main.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/scripts/json-tree-sitter-parse/main.go)
- Case matrix generator: [json-error-matrix/main.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/scripts/json-error-matrix/main.go)
- Heuristic probe generator: [json-heuristic-probe/main.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/scripts/json-heuristic-probe/main.go)

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

The ticket already includes a standalone parse inspector in [json-tree-sitter-parse/main.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/scripts/json-tree-sitter-parse/main.go). That should likely become a real CLI mode eventually.

Suggested future CLI:

```text
sanitize parse --format json
sanitize lint --format json
sanitize fix --format json
sanitize rules --format json
```

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

## Related

- [01-json-support-outline-and-malformed-llm-json-error-taxonomy.md](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/design-doc/01-json-support-outline-and-malformed-llm-json-error-taxonomy.md)
- [01-common-json-parse-errors-from-llm-output.md](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/reference/01-common-json-parse-errors-from-llm-output.md)
- [01-json-parse-error-replication-matrix.md](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/01-json-parse-error-replication-matrix.md)
- [02-json-heuristic-probe.md](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/02-json-heuristic-probe.md)
