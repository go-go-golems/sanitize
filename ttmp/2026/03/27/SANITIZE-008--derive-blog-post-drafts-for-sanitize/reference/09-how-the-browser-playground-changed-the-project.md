---
Title: How the browser playground changed the project
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
Summary: Draft structure for an article about the browser playground as an architectural forcing function rather than just a demo UI.
LastUpdated: 2026-03-27T09:34:54.233209327-04:00
WhatFor: Provide a writing scaffold for explaining how the local web UI improved API design, observability, and product thinking.
WhenToUse: Use when writing a project-surface article about developer tools, browser demos, and local inspection workflows.
---

# How the browser playground changed the project

## Goal

Prepare an article showing that the browser playground was not ornamental; it pushed the repo toward clearer APIs, better result shapes, and more usable diagnostics.

## Context

This article should appeal to engineers who have built “just a quick UI” and discovered that it changes the architecture underneath it.

## Quick Reference

**Working title:** How the browser playground changed the project

**Core thesis:** The browser UI turned `sanitize` from a package-and-CLI experiment into an inspection tool, and that forced better contracts across YAML and JSON.

**Suggested section map:**

1. `The UI started as a convenience`
   - Intent: establish humble origins.
   - Content: local playground, parse-tree inspection, fix review, example loading.
   - Original resources:
     - `internal/server/static/index.html`
     - `internal/server/static/js/app.js`
     - `internal/server/static/css/style.css`
   - Potential pseudocode and examples:
     - screenshot placeholder or layout sketch.

2. `The API had to become honest`
   - Intent: explain the biggest backend consequence.
   - Content: old YAML-shaped request versus new `{ format, input }` contract.
   - Original resources:
     - `internal/server/server.go`
     - `ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/reference/03-diary.md`
   - Potential pseudocode and examples:
     - request/response JSON blocks.

3. `The UI exposed result-shape gaps`
   - Intent: show how visual surfaces reveal missing fields.
   - Content: original strict JSON validity versus final strict validity, parse source versus heuristic source.
   - Original resources:
     - `pkg/json/types.go`
     - `pkg/json/sanitize.go`
   - Potential pseudocode and examples:
     - mention `OriginalStrictParseClean`.

4. `Why a shared YAML/JSON playground matters`
   - Intent: show the product-level gain.
   - Content: one interface, two engines, explicit differences.
   - Original resources:
     - `internal/server/server.go`
     - `examples/examples.go`
   - Potential pseudocode and examples:
     - examples filtered by format.

5. `The UI as a research lab`
   - Intent: tie the browser back to corpus and ticket work.
   - Content: inspect parse tree, compare fixes, load tricky examples, reason about repair safety.
   - Original resources:
     - `examples/json/`
     - `examples/yaml/`
     - `ttmp/2026/03/13/SANITIZE-005--split-ui-into-css-js-html-and-improve-preset-loading/`
   - Potential pseudocode and examples:
     - user flow: select example -> sanitize -> inspect issues.

6. `The broader lesson`
   - Intent: make it useful outside this repo.
   - Content: small UIs can sharpen the architecture of backend tools by making missing abstractions visible.
   - Original resources:
     - `README.md`
   - Potential pseudocode and examples:
     - optional “what the UI forced us to formalize” list.

## Usage Examples

- Good image or diagram candidate: editor pane -> API request -> parse/issues/output panes.
- Good anecdote: manual browser validation revealed the need for `OriginalStrictParseClean`.

## Related

- `reference/07-the-json-recovery-story-useful-but-narrow.md`
- `reference/08-designing-a-corpus-for-malformed-structured-text.md`
- `reference/11-a-weekend-project-that-turned-into-a-small-research-lab.md`
