---
Title: JSON detection buckets
Ticket: SANITIZE-007
Status: active
Topics:
    - json
    - linting
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Bucket classification showing which malformed JSON cases are parser-driven, heuristic-driven, or hybrid.
WhatFor: Use when deciding where tree-sitter alone is enough and where heuristic lint/fix rules must remain first-class.
WhenToUse: Read when prioritizing JSON rule work or evaluating recovery coverage.
---

# JSON Detection Buckets

## clean

- 00-valid-basic.json — Valid basic
- 01-valid-nested-objects.json — Valid nested objects
- Valid JSON

## parser-driven

- 13-missing-comma.json — Missing comma
- 14-missing-colon.json — Missing colon

## heuristic-driven

- 19-comments.json — Comments
- 21-multiple-top-level-objects.json — Multiple top level objects
- 23-duplicate-key.json — Duplicate key
- Multiple top-level values

## hybrid

- 10-trailing-comma-object.json — Trailing comma object
- 11-single-quotes.json — Single quotes
- 12-unquoted-keys.json — Unquoted keys
- 15-leading-prose-wrapper.json — Leading prose wrapper
- 16-markdown-fence-wrapper.json — Markdown fence wrapper
- 17-python-literals.json — Python literals
- 18-duplicate-comma-array.json — Duplicate comma array
- 20-missing-closing-delimiter.json — Missing closing delimiter
- 22-unterminated-string.json — Unterminated string
- 24-llm-wrapper-multi-step.json — Llm wrapper multi step
- 25-llm-commentary-comments-and-duplicate-comma.json — Llm commentary comments and duplicate comma
- Leading prose
- Markdown fence wrapper
- Python literals
- Single quotes
- Trailing comma
- Unquoted keys

