---
Title: Common JSON parse errors from LLM output
Ticket: SANITIZE-007
Status: active
Topics:
    - json
    - linting
    - api-design
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2026-03-13T19:01:55.035803211-04:00
WhatFor: Provide a practical list of malformed JSON patterns commonly produced by LLMs so JSON support can target real inputs.
WhenToUse: Use when choosing JSON lint rules, fixer behavior, and parser recovery priorities.
---

# Common JSON parse errors from LLM output

## Goal

List the malformed JSON patterns that show up frequently in LLM output and organize them by the kind of parse failure or recovery behavior they imply.

## Context

LLMs often produce “almost JSON” rather than strict JSON. The patterns are repetitive enough that they should shape the JSON-support roadmap directly.

## Quick Reference

Common malformed patterns:

- Trailing commas in arrays or objects.
- Single-quoted strings instead of double-quoted strings.
- Unquoted object keys.
- JavaScript-style comments: `// ...` or `/* ... */`.
- Markdown code fences around otherwise valid JSON.
- Explanatory prose before or after the JSON payload.
- Missing commas between array items or object members.
- Missing closing `]` or `}` because the LLM truncated output.
- Missing closing quote in a string.
- Bare newlines inside string values without escaping.
- Invalid escape sequences such as `\_` or half-finished Unicode escapes.
- Duplicate keys in the same object.
- Colon/comma confusion, especially after nested objects.
- Accidental ellipses like `...` inside arrays or objects.
- Placeholder tokens such as `<value>` or `[insert here]` inside otherwise JSON-looking content.
- Use of `True`, `False`, or `None` instead of `true`, `false`, or `null`.
- Numbers with leading zeros or other non-JSON numeric formatting.
- Extra closing delimiters such as `}}` or `]]`.
- Concatenated top-level objects without a separating array or newline-delimited framing.

Useful high-level buckets:

- Safe candidate auto-fixes:
  - strip Markdown fences
  - remove trailing commas
  - normalize `True` / `False` / `None`
- Often recoverable but needs care:
  - comments
  - single quotes
  - prose-wrapped JSON
  - duplicate keys
- Usually ambiguous:
  - missing commas
  - truncated output
  - incomplete strings
  - placeholder fragments
  - concatenated objects with unclear intended framing

## Usage Examples

Use this list to decide what the first JSON-support MVP should do:

- `lint only`: duplicate keys, ambiguous missing commas, truncated output
- `safe auto-fix`: Markdown fences, trailing commas, Python booleans/nulls
- `suggestion only`: incomplete strings, heavy prose wrapping, conflicting duplicate-key intent

## Related

- `design-doc/01-json-support-outline-and-malformed-llm-json-error-taxonomy.md`
