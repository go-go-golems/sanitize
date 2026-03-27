---
Title: A weekend project that turned into a small research lab
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
Summary: Draft structure for a narrative article about sanitize growing from a small weekend idea into a corpus-driven research-and-tooling project.
LastUpdated: 2026-03-27T09:35:05.508055804-04:00
WhatFor: Provide a narrative article plan that emphasizes project evolution, changing scope, and how documentation and tooling reinforced each other.
WhenToUse: Use when writing a higher-level retrospective about how sanitize grew over time.
---

# A weekend project that turned into a small research lab

## Goal

Prepare a narrative retrospective about how a narrow YAML cleanup idea gradually turned into a CLI, package, UI, example corpus, and research workflow.

## Context

This article should feel broader and more human than the subsystem deep dives. It is a project-history piece with enough technical anchors to stay concrete.

## Quick Reference

**Working title:** A weekend project that turned into a small research lab

**Core thesis:** `sanitize` became interesting not because it got bigger, but because each new surface made the earlier work more disciplined: tickets clarified decisions, examples clarified boundaries, and the UI clarified contracts.

**Suggested section map:**

1. `The small original project`
   - Intent: make the beginning modest.
   - Content: YAML sanitizing experiment, heuristic cleanup, reusable package idea.
   - Original resources:
     - `ttmp/2026/03/13/SANITIZE-001--turn-yaml-sanitizer-into-reusable-go-go-golems-sanitize-package/`
   - Potential pseudocode and examples:
     - early YAML rule examples.

2. `The project kept sprouting surfaces`
   - Intent: show growth by layers rather than chronology alone.
   - Content: CLI, package, examples, UI, JSON support, rules enumeration, validated selection.
   - Original resources:
     - `README.md`
     - `internal/cli/commands.go`
     - `internal/server/server.go`
   - Potential pseudocode and examples:
     - timeline or milestone table.

3. `Tickets turned experiments into decisions`
   - Intent: explain the value of `ttmp`.
   - Content: SANITIZE-004 for YAML tree-sitter-aware analysis, SANITIZE-005 for UI split, SANITIZE-007 for JSON recovery.
   - Original resources:
     - `ttmp/2026/03/13/SANITIZE-004--improve-yaml-linting-and-sanitizing-with-tree-sitter-aware-analysis/`
     - `ttmp/2026/03/13/SANITIZE-005--split-ui-into-css-js-html-and-improve-preset-loading/`
     - `ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/`
   - Potential pseudocode and examples:
     - quick “ticket -> what it changed” table.

4. `The research loop`
   - Intent: show how code and docs fed each other.
   - Content: corpus, report-generating scripts, design docs, implementation diaries, browser playground.
   - Original resources:
     - `examples/`
     - `ttmp/.../scripts/`
     - `ttmp/.../sources/`
   - Potential pseudocode and examples:
     - loop diagram: code -> examples -> reports -> design -> code.

5. `The asymmetry that made the project honest`
   - Intent: highlight that not every direction paid off equally.
   - Content: YAML feels more mature than JSON; this improved the project’s credibility instead of weakening it.
   - Original resources:
     - new Obsidian notes from `2026-03-27`
   - Potential pseudocode and examples:
     - YAML-vs-JSON comparison table.

6. `What I would carry into the next project`
   - Intent: end on transferable lessons.
   - Content: small tools get better when they gain inspectability, examples, and explicit boundaries.
   - Original resources:
     - `README.md`
   - Potential pseudocode and examples:
     - checklist of project-growth signals worth encouraging.

## Usage Examples

- This is a strong article for a more narrative, first-person voice.
- Good image candidate: milestone timeline from SANITIZE-001 through SANITIZE-007.

## Related

- `reference/02-from-regexes-to-a-real-sanitizer.md`
- `reference/09-how-the-browser-playground-changed-the-project.md`
- `reference/08-designing-a-corpus-for-malformed-structured-text.md`
