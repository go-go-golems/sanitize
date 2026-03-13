---
Title: JSON support outline and malformed LLM JSON error taxonomy
Ticket: SANITIZE-007
Status: active
Topics:
    - json
    - linting
    - api-design
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2026-03-13T19:01:54.917234698-04:00
WhatFor: Capture the initial JSON-support direction with emphasis on malformed JSON commonly produced by LLMs.
WhenToUse: Use when scoping JSON support, parser strategy, and error recovery behavior.
---

# JSON support outline and malformed LLM JSON error taxonomy

## Executive Summary

This ticket focuses on adding JSON support in a way that is practical for real LLM output, not just pristine hand-written JSON. The main architectural point is that strict parsing alone is not enough. The implementation should distinguish among strict parse failures, common recoverable LLM mistakes, and ambiguous situations where the tool should lint or suggest rather than silently rewrite.

## Problem Statement

LLM-generated JSON often fails in repetitive, recognizable ways: trailing commas, comments, single quotes, missing commas, Markdown fences, prose before or after the object, and incomplete closing delimiters. A JSON-support feature that only accepts strict valid JSON will be less useful for the actual input this project is likely to receive. A JSON-support feature that rewrites too aggressively will be risky.

The project therefore needs:

- a list of common malformed JSON cases,
- a decision about which cases are safe to auto-fix,
- and a parser/error model that can preserve structural context for linting and repair.

## Proposed Solution

Treat malformed LLM JSON as a first-class input category.

The initial implementation should:

- keep a strict JSON parse path as the source of truth,
- layer an LLM-oriented malformed-input taxonomy on top of that strict path,
- and classify malformed cases into:
  - safe auto-fix,
  - lint only,
  - suggestion only.

Examples of likely safe auto-fix:

- stripping Markdown code fences,
- removing trailing commas,
- normalizing `True` / `False` / `None`,
- removing obvious comments when the payload is otherwise clearly JSON.

Examples of likely suggestion-only behavior:

- missing commas versus accidental string concatenation,
- repairing truncated output where multiple completions are plausible,
- rewriting duplicate keys when intent is unclear.

## Design Decisions

### Decision 1: Use the LLM-error taxonomy to drive scope

The JSON feature should be designed around actual malformed-output patterns instead of starting from a purely strict-language view of JSON.

### Decision 2: Keep strict JSON and recovery behavior conceptually separate

That separation keeps it clear which diagnostics come from the parser and which come from recovery heuristics.

### Decision 3: Prefer conservative fixes for ambiguous malformed output

If the intended JSON meaning is unclear, the system should lint or suggest instead of rewriting silently.

## Alternatives Considered

### Strict JSON only

Rejected as the only mode because it ignores the most common LLM-output failure cases.

### Fully permissive JSON5-style parsing and then normalizing back to JSON

Possibly useful as an implementation aid, but not sufficient by itself because some malformed LLM output is not valid JSON5 either, and some repairs remain ambiguous.

## Implementation Plan

1. Build and maintain the malformed-JSON taxonomy.
2. Choose the parser and error-reporting strategy.
3. Define the JSON issue/fix model.
4. Add example fixtures for common malformed LLM output.
5. Implement JSON parse/lint/fix in `pkg/json`.
6. Wire JSON into the CLI and API surfaces.

## Open Questions

- Whether the JSON engine should accept prose-wrapped JSON as an input-recovery case.
- Whether duplicate keys should ever be auto-rewritten in JSON mode.
- Whether comments and trailing commas should be treated as JSON-specific recovery or as a broader “JSON-ish” compatibility mode.

## References

- `SANITIZE-003` for the broader JSON-support direction already present in the repo.
- `reference/01-common-json-parse-errors-from-llm-output.md` for the concrete malformed-input list.
