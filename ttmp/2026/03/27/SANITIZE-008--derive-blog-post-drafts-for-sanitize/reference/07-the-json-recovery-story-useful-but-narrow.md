---
Title: The JSON recovery story useful but narrow
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
Summary: Draft structure for an article focused on the JSON recovery work, emphasizing what shipped, what failed, and why the boundary stayed narrow.
LastUpdated: 2026-03-27T09:34:54.071883957-04:00
WhatFor: Provide a section plan for a grounded article about malformed-LLM JSON recovery as a constrained, evidence-backed project.
WhenToUse: Use when writing a JSON-focused article that balances shipped functionality with explicit limits.
---

# The JSON recovery story useful but narrow

## Goal

Prepare an article that tells the JSON story honestly: a lot of good research, a worthwhile implementation, and a much narrower repair surface than early excitement might suggest.

## Context

This article should read like a candid engineering case study. It should not be framed as failure, but it also should not overstate the fix coverage.

## Quick Reference

**Working title:** The JSON recovery story: useful, but narrow

**Core thesis:** The JSON side of `sanitize` became a strong malformed-input investigation platform and a decent wrapper-cleanup engine, but it never became a general structural JSON repair tool.

**Suggested section map:**

1. `The original promise`
   - Intent: explain why malformed LLM JSON looked like a compelling problem.
   - Content: wrappers, comments, Python literals, commas, “almost JSON” outputs from models.
   - Original resources:
     - `README.md`
     - `examples/json/README.md`
   - Potential pseudocode and examples:
     - example wrapped JSON block.

2. `The research phase`
   - Intent: show that the team investigated before claiming repair.
   - Content: error matrix, heuristic probe, detection buckets, overlap study, repair matrix.
   - Original resources:
     - `ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/01-json-parse-error-replication-matrix.md`
     - `ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/02-json-heuristic-probe.md`
     - `ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/04-json-detection-buckets.md`
     - `ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/05-json-repair-matrix.md`
     - `ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/07-json-overlap-study.md`
   - Potential pseudocode and examples:
     - compact table of research artifacts and what each one answered.

3. `What actually shipped`
   - Intent: map the implemented JSON surfaces.
   - Content: `pkg/json`, CLI support, HTTP API, browser playground.
   - Original resources:
     - `pkg/json/`
     - `internal/server/server.go`
     - `internal/server/static/js/app.js`
   - Potential pseudocode and examples:
     - `{ format, input }` request contract.

4. `The safe fixer set`
   - Intent: make the live repair boundary explicit.
   - Content: Markdown fences, prose extraction, Python literals, comments, duplicate commas, trailing commas.
   - Original resources:
     - `pkg/json/fix.go`
     - `pkg/json/sanitize_test.go`
   - Potential pseudocode and examples:
     - `applyFixes` flow summary.

5. `The stubborn cases`
   - Intent: show where the system stops.
   - Content: single quotes, unquoted keys, missing commas, missing colons, duplicate keys, multiple top-level values, missing closing delimiters.
   - Original resources:
     - `ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/05-json-repair-matrix.md`
     - `pkg/json/heuristics.go`
   - Potential pseudocode and examples:
     - show `sanitize fix --format json examples/json/11-single-quotes.json` and `sanitize lint --format json examples/json/13-missing-comma.json`.

6. `What made the work valuable anyway`
   - Intent: avoid reducing the story to “not enough fixes.”
   - Content: better API contract, UI generalization, strict vs structural validity, reusable corpus.
   - Original resources:
     - `pkg/json/types.go`
     - `internal/server/server.go`
     - `examples/json/`
   - Potential pseudocode and examples:
     - explain `StrictParseClean` versus `ParseClean`.

## Usage Examples

- This piece can use the repair matrix as its spine.
- Good ending: the repo learned where a trustworthy JSON sanitizer has to stop, and that boundary is itself a worthwhile result.

## Related

- `reference/10-what-i-learned-trying-to-repair-llm-json.md`
- `reference/12-what-not-to-auto-fix.md`
- `/home/manuel/code/wesen/obsidian-vault/Projects/2026/03/27/PROJ - Sanitize - JSON Recovery Experiments and Limits.md`
