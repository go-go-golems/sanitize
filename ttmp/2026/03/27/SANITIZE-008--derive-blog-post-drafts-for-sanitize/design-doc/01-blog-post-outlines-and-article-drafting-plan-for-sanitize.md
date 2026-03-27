---
Title: Blog post outlines and article drafting plan for sanitize
Ticket: SANITIZE-008
Status: active
Topics:
    - writing
    - documentation
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ../../../../../../../obsidian-vault/Projects/2026/03/27/PROJ - Sanitize - JSON Recovery Experiments and Limits.md
      Note: Vault summary used to shape JSON article directions
    - Path: ../../../../../../../obsidian-vault/Projects/2026/03/27/PROJ - Sanitize - YAML Sanitizing Deep Dive.md
      Note: Vault summary used to shape YAML article directions
    - Path: README.md
      Note: Repo-level framing for the article set
    - Path: ttmp/2026/03/13/SANITIZE-004--improve-yaml-linting-and-sanitizing-with-tree-sitter-aware-analysis/design-doc/01-tree-sitter-driven-yaml-linting-and-sanitizing-analysis.md
      Note: YAML architecture evidence for multiple articles
    - Path: ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/design-doc/01-json-support-outline-and-malformed-llm-json-error-taxonomy.md
      Note: JSON recovery evidence for multiple articles
ExternalSources: []
Summary: Evidence-backed plan for turning the sanitize repo, ticket history, and new project reports into a bundle of article-ready draft outlines.
LastUpdated: 2026-03-27T09:34:25.45384931-04:00
WhatFor: Define the article set, the structure each draft should follow, the repo resources each article can draw from, and the delivery plan for the bundle.
WhenToUse: Use when drafting blog posts or essays about sanitize without wanting the assistant to write final prose in place of the author.
---


# Blog post outlines and article drafting plan for sanitize

## Executive Summary

This ticket turns the sanitize project into a writing packet rather than a code-change ticket. The output is a collection of article-ready draft files: one file per article idea, each structured around section intent, likely content, repo-local resources, and candidate examples or pseudocode.

The packet is deliberately outline-heavy and paragraph-light. The goal is to preserve research, framing, and reusable examples while leaving voice, emphasis, and final prose to the human author.

## Problem Statement

The sanitize repo now has enough substance to support multiple articles, but the material is spread across code, tests, ticket docs, generated reports, and new Obsidian project notes. Without a structured writing packet, the likely outcome is either:

1. repeatedly rediscovering the same framing and evidence while drafting, or
2. asking an assistant to over-write the final article voice instead of helping prepare the material.

The problem this ticket addresses is therefore a documentation packaging problem: turn existing project evidence into article-shaped draft assets that a human can later expand in their own voice.

## Proposed Solution

Create a single ticket with three layers of deliverables:

1. a ticket-level design doc that explains the article set and the drafting approach,
2. a diary that records how the packet was assembled,
3. one reference document per article idea, each containing:
   - working title and thesis,
   - recommended structure,
   - section-by-section intent,
   - section-by-section content guidance,
   - repo-local source material,
   - candidate examples, snippets, or pseudocode.

The article draft files should be copy-forward documents, not final essays. They are intended to reduce activation energy for later writing while keeping factual grounding close to the code and ticket history.

## Design Decisions

### Decision 1: One file per article idea

This keeps each writing direction self-contained and avoids a giant omnibus document that would be hard to scan or bundle.

### Decision 2: Drafts should stop before full prose

The user explicitly wants section intent and content planning, not finished paragraphs. This keeps authorship and voice with the human writer.

### Decision 3: Prefer repo-local sources over fresh speculation

Every article draft should point back to concrete files:

- repo code such as `pkg/yaml/*`, `pkg/json/*`, `internal/server/*`,
- ticket docs such as `SANITIZE-004` and `SANITIZE-007`,
- example corpora under `examples/yaml` and `examples/json`,
- newly written Obsidian project reports.

### Decision 4: Keep the bundle reMarkable-friendly

The set of docs should be individually useful but also coherent as one bundled PDF with a table of contents.

## Alternatives Considered

### One giant writing brief

Rejected because it would make the packet harder to navigate and much harder to split by writing mood or publication target.

### Finished assistant-written blog posts

Rejected because the user explicitly wants to write the prose later in their own voice.

### A spreadsheet or table-only plan

Rejected because it would flatten article flow too much and lose the narrative shape that makes drafting easy.

## Implementation Plan

### Phase 1: Create ticket and writing scaffolding

- create `SANITIZE-008`,
- add the primary design doc,
- add the diary,
- create one doc per article idea.

### Phase 2: Fill ticket-level framing

- update the index,
- replace placeholder tasks,
- update the changelog,
- define the writing packet’s scope and goals.

### Phase 3: Fill article drafts

For each article draft:

1. set the working thesis,
2. define a section map,
3. capture intent and content for each section,
4. attach likely repo resources,
5. suggest examples, CLI outputs, diagrams, or pseudocode.

### Phase 4: Validate and deliver

- run `docmgr doctor --ticket SANITIZE-008 --stale-after 30`,
- add any missing vocabulary if required,
- dry-run reMarkable bundle upload,
- upload real bundle,
- verify remote listing.

## Open Questions

- Which of the twelve article directions are strongest enough to promote into full essays first?
- Should later drafts add publication-target variants such as “personal retrospective” versus “deep technical post”?
- Should the Obsidian project reports be included in future bundles, or kept separate from the ticket docs?

## References

- `README.md`
- `pkg/yaml/analysis.go`
- `pkg/yaml/lint.go`
- `pkg/yaml/fix.go`
- `pkg/yaml/sanitize.go`
- `pkg/json/analysis.go`
- `pkg/json/heuristics.go`
- `pkg/json/fix.go`
- `pkg/json/sanitize.go`
- `internal/server/server.go`
- `internal/server/static/js/app.js`
- `ttmp/2026/03/13/SANITIZE-004--improve-yaml-linting-and-sanitizing-with-tree-sitter-aware-analysis/`
- `ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/`
- `/home/manuel/code/wesen/obsidian-vault/Projects/2026/03/27/PROJ - Sanitize - YAML Sanitizing Deep Dive.md`
- `/home/manuel/code/wesen/obsidian-vault/Projects/2026/03/27/PROJ - Sanitize - JSON Recovery Experiments and Limits.md`
