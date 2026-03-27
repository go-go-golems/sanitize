---
Title: Tree-sitter as infrastructure, not a magic wand
Ticket: SANITIZE-008
Status: active
Topics:
    - writing
    - documentation
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Draft structure for an article about how sanitize uses tree-sitter as one source of structural evidence rather than a complete repair solution.
LastUpdated: 2026-03-27T09:34:45.741225258-04:00
WhatFor: Provide a writing scaffold for explaining the practical role of tree-sitter inside sanitize.
WhenToUse: Use when writing a parser-tooling article that pushes back on “just use tree-sitter” simplifications.
---

# Tree-sitter as infrastructure, not a magic wand

## Goal

Prepare an article that explains the real role tree-sitter plays in the repo: essential, but insufficient by itself.

## Context

This article is for readers who understand parsers or editor tooling and may be tempted to assume tree-sitter alone solves recovery. It should be calmly corrective rather than argumentative.

## Quick Reference

**Working title:** Tree-sitter as infrastructure, not a magic wand

**Core thesis:** Tree-sitter made `sanitize` structurally aware, but the project only worked once parser facts were integrated with heuristics, strict parsing, and explicit rule policy.

**Suggested section map:**

1. `What people hope tree-sitter will do`
   - Intent: start from a recognizable misconception.
   - Content: “parse the thing and fix it” fantasy versus what parser errors actually provide.
   - Original resources:
     - `pkg/yaml/parse.go`
     - `pkg/json/parse.go`
   - Potential pseudocode and examples:
     - tiny `ParseTree` example.

2. `What tree-sitter really gives sanitize`
   - Intent: describe the actual benefits.
   - Content: error nodes, missing nodes, spans, tree text, duplicate-key traversal targets.
   - Original resources:
     - `pkg/yaml/analysis.go`
     - `pkg/json/analysis.go`
     - `pkg/yaml/duplicate_keys.go`
     - `pkg/json/duplicate_keys.go`
   - Potential pseudocode and examples:
     - `documentAnalysis` struct excerpt.

3. `Why parser output still needs help`
   - Intent: show the gap.
   - Content: heuristic-only issues in YAML, strict-parser-only or hybrid issues in JSON, coarse parser location on some errors.
   - Original resources:
     - `ttmp/2026/03/13/SANITIZE-004--improve-yaml-linting-and-sanitizing-with-tree-sitter-aware-analysis/design-doc/01-tree-sitter-driven-yaml-linting-and-sanitizing-analysis.md`
     - `ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/04-json-detection-buckets.md`
   - Potential pseudocode and examples:
     - parse-only / heuristic-only / hybrid table.

4. `How sanitize turned parser data into infrastructure`
   - Intent: explain the architectural move.
   - Content: shared analysis objects, line indexes, unified `LintIssue` spans, parser-derived lint issues.
   - Original resources:
     - `pkg/yaml/types.go`
     - `pkg/json/types.go`
     - `pkg/yaml/lint.go`
     - `pkg/json/lint.go`
   - Potential pseudocode and examples:
     - issue-shape snippet.

5. `The broader lesson`
   - Intent: close with transferable advice.
   - Content: parsers are great witnesses, not omniscient repair authors.
   - Original resources:
     - `README.md`
   - Potential pseudocode and examples:
     - optional decision checklist for parser-assisted tools.

## Usage Examples

- A good visual is a three-column matrix: parser sees, heuristics see, fixer acts.
- A useful anecdote is YAML tab indentation, where parser location and actionable location do not fully coincide.

## Related

- `reference/04-building-a-conservative-repair-engine.md`
- `reference/13-how-to-write-span-aware-diagnostics-for-text-repair-tools.md`
- `reference/06-anatomy-of-the-yaml-pipeline.md`
