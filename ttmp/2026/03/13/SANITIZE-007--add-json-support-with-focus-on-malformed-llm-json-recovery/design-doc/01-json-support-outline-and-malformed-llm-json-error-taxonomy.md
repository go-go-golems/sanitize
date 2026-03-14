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
RelatedFiles:
    - Path: ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/scripts/json-error-matrix/main.go
      Note: Replicates malformed cases against encoding/json and tree-sitter
    - Path: ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/scripts/json-heuristic-probe/main.go
      Note: Probes heuristic detectability of malformed JSON cases
    - Path: ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/scripts/json-tree-sitter-parse/main.go
      Note: Standalone tree-sitter JSON parse inspector
    - Path: ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/01-json-parse-error-replication-matrix.md
      Note: Generated replication report comparing encoding/json and tree-sitter
    - Path: ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/02-json-heuristic-probe.md
      Note: Generated heuristic hit report for malformed JSON cases
    - Path: ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/local/03-json-parse-errors-import.md
      Note: Normalized imported malformed JSON case source used for experiments
    - Path: ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/local/json-parse-errors.md
      Note: Imported malformed JSON case source used for experiments
ExternalSources: []
Summary: Final JSON-support design and shipped first-release boundary for malformed LLM JSON recovery across the package, CLI, HTTP API, and browser UI.
LastUpdated: 2026-03-13T23:40:00-04:00
WhatFor: Capture the implemented JSON-support direction, the first-release recovery boundary, and the remaining lint-only cases.
WhenToUse: Use when reviewing what shipped in SANITIZE-007 and why certain malformed JSON patterns remain lint-only.
---



# JSON support outline and malformed LLM JSON error taxonomy

## Executive Summary

This ticket focuses on adding JSON support in a way that is practical for real LLM output, not just pristine hand-written JSON. The main architectural point is that strict parsing alone is not enough. The implementation should distinguish among strict parse failures, common recoverable LLM mistakes, and ambiguous situations where the tool should lint or suggest rather than silently rewrite.

## Implemented State

As of `2026-03-13`, `sanitize` now ships first-release JSON support across all public surfaces:

- `pkg/json` provides parse, lint, rule selection, and conservative sanitize flows.
- `sanitize lint --format json`, `sanitize fix --format json`, `sanitize parse --format json`, and `sanitize rules --format json` are all implemented.
- the HTTP API uses `{ "format": "...", "input": "..." }`.
- the embedded browser app is now a shared YAML/JSON playground with a JSON recovery mode.
- ticket-local experiment scripts generate:
  - example metadata,
  - detection buckets,
  - repair matrices,
  - machine-readable rule matrices,
  - heuristic overlap studies.

The first-release JSON auto-fix boundary is intentionally conservative. These patterns are auto-fixed today:

- Markdown fence wrappers
- leading or trailing prose extraction
- Python literal normalization
- comments
- duplicate commas
- trailing commas

These patterns are linted but intentionally left unfixed:

- unquoted keys
- single quotes
- missing commas
- missing colons
- missing closing delimiters
- multiple top-level values
- duplicate keys

That split is deliberate. The shipped policy is:

- safe wrapper and token normalization may auto-fix,
- ambiguous structural repairs remain lint-only,
- strict `encoding/json` validity remains the final truth for "valid JSON".

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

### Decision 4: Keep YAML and JSON packages separate behind shared CLI and server contracts

The implementation kept `pkg/yaml` and `pkg/json` as separate engines. The commonality lives in public surfaces such as:

- `sanitize <verb> --format ...`
- `/api/sanitize` and `/api/parse`
- the browser result renderer

This was the right tradeoff for the first release because YAML and JSON still differ materially in:

- parse recovery behavior,
- safe fix sets,
- strict-validity rules,
- and rule catalogs.

### Decision 5: Duplicate keys are lint-only in JSON mode

Unlike YAML, JSON duplicate keys are not automatically renamed or rewritten. The current engine reports them clearly but does not change key names because doing so would invent semantics.

### Decision 6: No separate "recovery-only" preset in the first release

The first release keeps one default JSON rule set with conservative fixes and lint-only structural rules. A future preset split is still possible, but it was not necessary to ship the current functionality.

## Acceptance Contract

The shipped JSON sanitize contract is:

- preserve the original input in the result,
- expose original and sanitized tree-sitter trees,
- expose original and sanitized parse/lint state,
- expose original strict-parse status separately from final strict-parse status,
- record every applied fix in order,
- stop iterating when the document converges or the iteration cap is reached.

The concrete result fields live in [pkg/json/types.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/json/types.go).

## Rule Severity And Confidence Policy

The rule matrix generated in [06-json-rule-matrix.json](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/06-json-rule-matrix.json) formalizes the shipped policy, but the short version is:

- `warn` + safe auto-fix:
  - wrappers
  - comments
  - duplicate commas
  - trailing commas
  - Python literals
- `error` + lint-only:
  - missing commas
  - missing colons
  - missing closing delimiters
  - unquoted keys
  - single quotes
  - multiple top-level values
- `warn` + lint-only:
  - duplicate keys

Confidence levels:

- `safe-auto-fix`
  - the transform is localized and low ambiguity
- `multi-step-safe-auto-fix`
  - several individually safe transforms may be chained
- `suggestion-or-future-fix`
  - a likely future fixer exists, but current confidence is not high enough
- `lint-only`
  - the project should not rewrite this case automatically

## Alternatives Considered

### Strict JSON only

Rejected as the only mode because it ignores the most common LLM-output failure cases.

### Fully permissive JSON5-style parsing and then normalizing back to JSON

Possibly useful as an implementation aid, but not sufficient by itself because some malformed LLM output is not valid JSON5 either, and some repairs remain ambiguous.

## Implementation Summary

The implementation landed in this order:

1. Build `pkg/json` analysis, strict-parse, duplicate-key, and rule-selection foundations.
2. Add parse-derived and heuristic JSON lint rules.
3. Implement conservative JSON fixers and iterative sanitize behavior.
4. Add JSON examples and combined malformed cases under [examples/json](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/examples/json).
5. Wire JSON into the Glazed CLI.
6. Generalize the HTTP API to `{ format, input }`.
7. Convert the embedded UI into a format-aware playground with JSON recovery mode.
8. Add ticket-local experiment scripts and generated evidence documents.

### Server and browser implications

The repository already has an embedded HTTP server and single-page web UI, and this ticket now generalizes both successfully.

The server/API work now includes:

- format-aware sanitize dispatch in [internal/server/server.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/internal/server/server.go)
- format-aware parse dispatch in [internal/server/server.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/internal/server/server.go)
- mixed YAML/JSON example payloads from `/api/examples`
- JSON-specific parse responses that include `strict_parse_clean`

The UI work now includes:

- a format selector in the UI,
- format-filtered examples from `/api/examples`,
- format-aware requests using `{ "format": "json", "input": "..." }`,
- a JSON parse playground mode that exposes:
  - tree-sitter parse tree,
  - original strict parser success/failure,
  - heuristic lint findings,
  - applied fixes and repaired output.

The UI remains shared between YAML and JSON. The goal is one playground with two formats, not two separate browser apps, and that is the shape now in [internal/server/static/index.html](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/internal/server/static/index.html) and [internal/server/static/js/app.js](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/internal/server/static/js/app.js).

## Open Questions

- Whether single quotes and unquoted keys should become suggestion-backed fixers or remain lint-only permanently.
- Whether a future "LLM recovery" preset should enable more aggressive rewrites than the default JSON rule set.
- Whether the duplicate-key lint should eventually surface related spans or additional object-path context in the UI.

## References

- `SANITIZE-003` for the broader JSON-support direction already present in the repo.
- `reference/01-common-json-parse-errors-from-llm-output.md` for the concrete malformed-input list.
