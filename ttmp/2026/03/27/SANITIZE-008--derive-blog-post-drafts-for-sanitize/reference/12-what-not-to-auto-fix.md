---
Title: What not to auto-fix
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
Summary: Draft structure for an article about ambiguity, restraint, and why some malformed inputs should remain lint-only.
LastUpdated: 2026-03-27T09:35:05.666270444-04:00
WhatFor: Provide an article scaffold focused on non-goals, ambiguity, and the cost of overconfident automatic rewriting.
WhenToUse: Use when writing a principle-driven article about the negative space of repair systems.
---

# What not to auto-fix

## Goal

Prepare an article arguing that the most important rule in a repair engine may be the rule that prevents repair.

## Context

This article should be opinionated but evidence-backed. It is less about how sanitize works and more about where it responsibly stops.

## Quick Reference

**Working title:** What not to auto-fix

**Core thesis:** Automatic repair becomes untrustworthy the moment it starts inventing semantics, so the boundary between auto-fix and lint-only must be explicit and defended.

**Suggested section map:**

1. `The temptation to fix everything`
   - Intent: start with the familiar pressure.
   - Content: users want the document to become valid; engineers want clean demos; ambiguous fixes are seductive.
   - Original resources:
     - `README.md`
   - Potential pseudocode and examples:
     - a “make it valid somehow” anti-pattern quote or paraphrase.

2. `Good auto-fixes are boring`
   - Intent: define the safe zone.
   - Content: tabs, missing spaces, trailing commas, wrappers, comments, Python literals.
   - Original resources:
     - `pkg/yaml/fix.go`
     - `pkg/json/fix.go`
   - Potential pseudocode and examples:
     - before/after examples for a few safe rules.

3. `Ambiguous fixes change meaning`
   - Intent: show where the danger starts.
   - Content: missing comma, missing colon, duplicate keys, multiple top-level values, single quotes in mixed contexts.
   - Original resources:
     - `ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/05-json-repair-matrix.md`
     - `pkg/json/heuristics.go`
   - Potential pseudocode and examples:
     - list of ambiguous inputs with multiple plausible repairs.

4. `How sanitize encodes restraint`
   - Intent: tie the philosophy to the live repo.
   - Content: rule catalogs, lint-versus-fix split, no-progress exits, explicit diagnostics.
   - Original resources:
     - `pkg/json/rules.go`
     - `pkg/yaml/rules.go`
     - `pkg/json/sanitize.go`
     - `pkg/yaml/sanitize.go`
   - Potential pseudocode and examples:
     - policy matrix or rule table.

5. `Why refusal is a feature`
   - Intent: reframe user-facing failure.
   - Content: trust, auditability, safer pipelines, easier debugging.
   - Original resources:
     - current CLI behavior from repo examples
   - Potential pseudocode and examples:
     - side-by-side of an explicit lint result versus a risky guessed rewrite.

6. `A checklist for future fixers`
   - Intent: leave the reader with a tool.
   - Content: localization, ambiguity, reversibility, explainability, test coverage, corpus evidence.
   - Original resources:
     - lessons across `SANITIZE-004` and `SANITIZE-007`
   - Potential pseudocode and examples:
     - five-question checklist block.

## Usage Examples

- Strong pull quote: “A sanitizer earns trust as much by what it refuses to rewrite as by what it fixes.”
- Strong visual: auto-fix / lint-only / never-fix triangle.

## Related

- `reference/04-building-a-conservative-repair-engine.md`
- `reference/07-the-json-recovery-story-useful-but-narrow.md`
- `reference/03-why-yaml-sanitizing-worked-better-than-json-repair.md`
