---
Title: Designing a corpus for malformed structured text
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
Summary: Draft structure for an article about how the sanitize example corpus and generated reports shaped rule development and project confidence.
LastUpdated: 2026-03-27T09:34:54.222705306-04:00
WhatFor: Provide a writing plan for an article about corpus design as an engineering tool rather than just demo material.
WhenToUse: Use when writing about examples, fixtures, evidence loops, and report-generating scripts.
---

# Designing a corpus for malformed structured text

## Goal

Prepare an article about the role the YAML and JSON corpora played in making `sanitize` more disciplined and more explainable.

## Context

This article should be useful to people building linters, parsers, or recovery tools. The hook is that a corpus is not just for tests; it can drive architecture decisions.

## Quick Reference

**Working title:** Designing a corpus for malformed structured text

**Core thesis:** The `sanitize` corpus worked because it was treated as a research instrument, a regression suite, and a UI/demo source all at once.

**Suggested section map:**

1. `Why examples matter more than intuition`
   - Intent: frame the problem.
   - Content: without examples, repair tools become argument-by-anecdote.
   - Original resources:
     - `examples/yaml/README.md`
     - `examples/json/README.md`
   - Potential pseudocode and examples:
     - simple table of example classes.

2. `The YAML corpus`
   - Intent: show how YAML examples exposed parser-versus-linter gaps.
   - Content: valid cases, malformed cases, mixed-error cases.
   - Original resources:
     - `examples/yaml/`
     - `ttmp/2026/03/13/SANITIZE-004--improve-yaml-linting-and-sanitizing-with-tree-sitter-aware-analysis/sources/01-example-corpus-parse-vs-lint-matrix.md`
   - Potential pseudocode and examples:
     - cite `20-mixed-indent.yaml` and `24-deeply-nested-mixed-errors.yaml`.

3. `The JSON corpus`
   - Intent: show how the JSON side was organized around malformed-LLM patterns.
   - Content: valid controls, single-pattern failures, multi-pattern combo failures.
   - Original resources:
     - `examples/json/`
     - `pkg/json/examples.go`
   - Potential pseudocode and examples:
     - `00-09`, `10-19`, `20-29` naming convention.

4. `Generated reports as corpus lenses`
   - Intent: explain why the project generated markdown reports from the corpus.
   - Content: error matrix, heuristic probe, repair matrix, overlap study.
   - Original resources:
     - `ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/scripts/`
     - corresponding `sources/` docs
   - Potential pseudocode and examples:
     - a tiny “corpus -> script -> report -> design decision” diagram.

5. `How the corpus changed the code`
   - Intent: tie fixtures back to implementation.
   - Content: mixed indentation rule, JSON repair boundary, CLI and UI examples.
   - Original resources:
     - `pkg/yaml/lint.go`
     - `pkg/json/fix.go`
     - `internal/server/server.go`
   - Potential pseudocode and examples:
     - examples endpoint surfacing corpus files.

6. `Advice for building your own corpus`
   - Intent: make the article transferable.
   - Content: include clean controls, single-pattern cases, combo cases, generated metadata, and human-readable reports.
   - Original resources:
     - repo structure itself
   - Potential pseudocode and examples:
     - short checklist.

## Usage Examples

- Strong visual: folder tree plus generated reports.
- Strong sidebar: “five ways a corpus can be more than a test fixture directory.”

## Related

- `reference/07-the-json-recovery-story-useful-but-narrow.md`
- `reference/06-anatomy-of-the-yaml-pipeline.md`
- `reference/09-how-the-browser-playground-changed-the-project.md`
