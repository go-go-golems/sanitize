---
Title: Anatomy of the YAML pipeline
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
Summary: Draft structure for a detailed technical article focused specifically on the YAML package architecture and repair loop.
LastUpdated: 2026-03-27T09:34:54.061653415-04:00
WhatFor: Provide a section-by-section plan for a deep technical article about the YAML subsystem.
WhenToUse: Use when writing a code-heavy engineering post focused on `pkg/yaml`.
---

# Anatomy of the YAML pipeline

## Goal

Prepare a technical deep dive that walks a reader through the YAML package from parse to lint to fix to final result.

## Context

This is the most code-centric YAML article candidate. It should assume the reader is comfortable reading small Go snippets and wants implementation detail rather than just a project story.

## Quick Reference

**Working title:** Anatomy of the YAML pipeline

**Core thesis:** The YAML subsystem works because it layers parse-aware analysis, span-rich diagnostics, and narrow iterative fixers into one coherent loop.

**Suggested section map:**

1. `Map the package`
   - Intent: orient the reader quickly.
   - Content: `analysis.go`, `parse.go`, `lint.go`, `duplicate_keys.go`, `fix.go`, `sanitize.go`, `types.go`.
   - Original resources:
     - `pkg/yaml/`
   - Potential pseudocode and examples:
     - file map table.

2. `Shared analysis`
   - Intent: explain the central architectural object.
   - Content: `documentAnalysis`, line index, tree text, parse errors, duplicate keys.
   - Original resources:
     - `pkg/yaml/analysis.go`
     - `pkg/yaml/line_index.go`
   - Potential pseudocode and examples:
     - struct snippet and one sentence on each field.

3. `Lint assembly`
   - Intent: explain how issues are constructed.
   - Content: parse-derived issues, line heuristics, mixed indentation, duplicate keys.
   - Original resources:
     - `pkg/yaml/lint.go`
     - `pkg/yaml/types.go`
   - Potential pseudocode and examples:
     - `issues = lintFromParseErrors + lintLineIssues + lintMixedIndentation`

4. `Fix application`
   - Intent: walk through one repair round.
   - Content: row-indexed lint lookup, parse-error row expansion, line fixers, document-level duplicate-key and mixed-indent repair.
   - Original resources:
     - `pkg/yaml/fix.go`
   - Potential pseudocode and examples:
     - simplified `applyFixes` pseudocode.

5. `Iteration and stopping`
   - Intent: explain why repeated passes matter.
   - Content: original state capture, iteration cap, no-progress exit, original vs final diagnostics.
   - Original resources:
     - `pkg/yaml/sanitize.go`
   - Potential pseudocode and examples:
     - main sanitize loop.

6. `Proof by example`
   - Intent: show the system handling a messy real input.
   - Content: `examples/yaml/24-deeply-nested-mixed-errors.yaml` and the current CLI output.
   - Original resources:
     - `examples/yaml/24-deeply-nested-mixed-errors.yaml`
     - `pkg/yaml/sanitize_test.go`
   - Potential pseudocode and examples:
     - before/after plus fix log excerpt.

7. `Testing and confidence`
   - Intent: close on reviewability.
   - Content: tests for duplicate-key scope, structural issues, mixed indentation, and span-carrying heuristics.
   - Original resources:
     - `pkg/yaml/sanitize_test.go`
   - Potential pseudocode and examples:
     - list of key tests and what they prove.

## Usage Examples

- Strong diagram candidate: a mermaid flow from `Sanitize` to `analyzeDocument` to `lintIssuesFromAnalysis` to `applyFixes`.
- Strong side panel: “seven YAML rules and what each one dares to change.”

## Related

- `reference/02-from-regexes-to-a-real-sanitizer.md`
- `reference/05-tree-sitter-as-infrastructure-not-a-magic-wand.md`
- `/home/manuel/code/wesen/obsidian-vault/Projects/2026/03/27/PROJ - Sanitize - YAML Sanitizing Deep Dive.md`
