---
Title: JSON heuristic probe
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
Summary: Generated heuristic hit report showing which malformed JSON cases are detectable by simple string/regex heuristics before structural parsing.
WhatFor: Show where simple heuristics are enough and where tree-sitter or stricter parsing is still required.
WhenToUse: Use when deciding which malformed JSON cases should become heuristic lint rules.
---

# JSON Heuristic Probe

Source: `../sources/local/03-json-parse-errors-import.md`

## Case 1: Trailing commas

- `trailing_comma`

## Case 2: Single quotes instead of double quotes

- `single_quotes`

## Case 3: Unquoted keys

- `unquoted_keys`

## Case 4: Missing comma between fields

- `none`

## Case 5: Missing colon after key

- `none`

## Case 6: Unescaped double quotes inside strings

- `none`

## Case 7: Invalid escape sequences

- `none`

## Case 8: Unterminated strings

- `none`

## Case 9: Mismatched braces or brackets

- `none`

## Case 10: Extra text before or after JSON

- `prose_wrapped`

## Case 11: Using comments

- `comment`

## Case 12: Multiple top-level objects

- `multiple_top_level`

## Case 13: Invalid numbers

- `none`

## Case 14: Using Python/JS literals instead of JSON literals

- `python_literals`

## Case 15: Markdown code fences included

- `markdown_fence`

## Case 16: Ellipses or placeholders

- `ellipsis_or_placeholder`

## Case 17: Duplicate commas or malformed arrays

- `duplicate_comma`

## Case 18: Unicode/control character issues

- `none`

