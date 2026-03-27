---
Title: Building a conservative repair engine
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
Summary: Draft structure for an article about the philosophy and mechanics of conservative automatic repair in sanitize.
LastUpdated: 2026-03-27T09:34:45.733569784-04:00
WhatFor: Provide a writing scaffold for explaining the rule policy, stop conditions, and trust model behind sanitize.
WhenToUse: Use when writing a design-principles post about safe automatic repair.
---

# Building a conservative repair engine

## Goal

Prepare an article centered on the product and engineering philosophy of “repair only what you can explain.”

## Context

This article should read like a design-values essay with concrete implementation examples. It is a good candidate if the goal is to make the repo’s choices feel principled rather than arbitrary.

## Quick Reference

**Working title:** Building a conservative repair engine

**Core thesis:** A sanitizer earns trust by making small, justified edits and by refusing repairs once intent becomes ambiguous.

**Suggested section map:**

1. `Define conservative repair`
   - Intent: make the term precise early.
   - Content: explain localized transform, low ambiguity, reversible reasoning, and stop conditions.
   - Original resources:
     - `README.md`
     - `pkg/yaml/fix.go`
     - `pkg/json/fix.go`
   - Potential pseudocode and examples:
     - simple rule definition table.

2. `The rule catalog as policy`
   - Intent: show that rule lists are not just implementation detail.
   - Content: YAML rule surface versus JSON rule surface; lints versus fixes; parse-aware flags.
   - Original resources:
     - `pkg/yaml/rules.go`
     - `pkg/json/rules.go`
   - Potential pseudocode and examples:
     - rule matrix snippet.

3. `Iterative repair instead of one giant rewrite`
   - Intent: explain the sanitize loop.
   - Content: analyze, lint, apply small fixes, re-run, stop on convergence.
   - Original resources:
     - `pkg/yaml/sanitize.go`
     - `pkg/json/sanitize.go`
   - Potential pseudocode and examples:
     - core sanitize loop pseudocode.

4. `Good examples of conservative fixes`
   - Intent: make the philosophy concrete.
   - Content: tab indentation, missing space after colon, prose extraction, Markdown fence removal, trailing comma removal.
   - Original resources:
     - `pkg/yaml/fix.go`
     - `pkg/json/fix.go`
     - `pkg/json/sanitize_test.go`
   - Potential pseudocode and examples:
     - before/after snippets.

5. `Where the engine refuses to guess`
   - Intent: show the boundary.
   - Content: missing commas, missing colons, single quotes in JSON, duplicate keys in JSON, multiple top-level values.
   - Original resources:
     - `ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/05-json-repair-matrix.md`
   - Potential pseudocode and examples:
     - lint-only examples block.

6. `Why this matters`
   - Intent: end on trust rather than syntax.
   - Content: explain why conservative tools are more reusable in pipelines and why explicit failure is a feature.
   - Original resources:
     - `README.md`
   - Potential pseudocode and examples:
     - optional short checklist for deciding whether to auto-fix a new rule.

## Usage Examples

- Good diagram candidate: `input -> analyze -> lint -> fix -> reanalyze -> stop`.
- Good sidebar: “Questions to ask before adding a new fixer.”

## Related

- `reference/12-what-not-to-auto-fix.md`
- `reference/05-tree-sitter-as-infrastructure-not-a-magic-wand.md`
- `reference/07-the-json-recovery-story-useful-but-narrow.md`
