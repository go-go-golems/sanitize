---
Title: How to write span-aware diagnostics for text repair tools
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
Summary: Draft structure for an article about byte spans, row/column mapping, and unified diagnostics in parser-assisted repair tools.
LastUpdated: 2026-03-27T09:35:05.666877071-04:00
WhatFor: Provide a technical article plan centered on diagnostic shapes and why span-aware issue models matter.
WhenToUse: Use when writing for readers interested in tool UX, diagnostics, and parser-assisted editor or CLI tooling.
---

# How to write span-aware diagnostics for text repair tools

## Goal

Prepare an article that explains why “line 7 is wrong” is not enough for serious repair tooling, and how `sanitize` moved toward richer diagnostics.

## Context

This article is a narrower technical piece for tooling engineers. It should use `sanitize` as the concrete case study but end with general advice for linters, formatters, and repair tools.

## Quick Reference

**Working title:** How to write span-aware diagnostics for text repair tools

**Core thesis:** Span-aware diagnostics are the bridge between parser output, heuristics, CLI UX, and interactive tooling, and `sanitize` became much more coherent once its issue types carried real location data.

**Suggested section map:**

1. `Why row-only diagnostics break down`
   - Intent: explain the problem.
   - Content: parser spans can be wide, heuristics can be line-local, and UIs need a common model.
   - Original resources:
     - `pkg/yaml/types.go`
     - `pkg/json/types.go`
   - Potential pseudocode and examples:
     - “row only” versus “start/end span” comparison.

2. `The data shapes in sanitize`
   - Intent: make the case concrete.
   - Content: `ErrorNode`, `LintIssue`, `Fix`, `Result`.
   - Original resources:
     - `pkg/yaml/types.go`
     - `pkg/json/types.go`
   - Potential pseudocode and examples:
     - struct snippets.

3. `Line indexes make parser offsets usable`
   - Intent: explain the glue layer.
   - Content: byte offsets, row/column conversion, duplicate-key location mapping, strict-parse offsets in JSON.
   - Original resources:
     - `pkg/yaml/line_index.go`
     - `pkg/json/line_index.go`
     - `pkg/json/lint.go`
   - Potential pseudocode and examples:
     - `rowColAtByte` style pseudocode.

4. `Unifying parse and heuristic issues`
   - Intent: show the architectural payoff.
   - Content: parser-derived issues and heuristic-derived issues can now share one display contract.
   - Original resources:
     - `pkg/yaml/lint.go`
     - `pkg/json/lint.go`
   - Potential pseudocode and examples:
     - issue assembly pseudocode.

5. `Why the UI cares`
   - Intent: connect diagnostics to user-facing tooling.
   - Content: highlighting, issue panels, parse source versus heuristic source, strict parser badges.
   - Original resources:
     - `internal/server/static/js/app.js`
     - `internal/server/server.go`
   - Potential pseudocode and examples:
     - render pipeline sketch.

6. `Lessons for other tools`
   - Intent: make it portable.
   - Content: carry spans early, unify issue types, keep parser errors and heuristic warnings distinguishable.
   - Original resources:
     - full repo as case study
   - Potential pseudocode and examples:
     - a short “diagnostic model checklist.”

## Usage Examples

- Strong visual: one issue shown as byte span, row/column span, source, and rule.
- Strong example pair: YAML mixed indentation and JSON strict-parse offset reporting.

## Related

- `reference/05-tree-sitter-as-infrastructure-not-a-magic-wand.md`
- `reference/06-anatomy-of-the-yaml-pipeline.md`
- `reference/09-how-the-browser-playground-changed-the-project.md`
