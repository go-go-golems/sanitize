---
Title: JSON heuristic overlap study
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
Summary: Overlap study comparing heuristic issue spans with tree-sitter parse-error spans across the JSON corpus.
WhatFor: Show where parser spans reinforce heuristic confidence and where heuristics must stand on their own.
WhenToUse: Use when deciding which heuristic fixes should remain automatic versus suggestion-only.
---

# JSON Heuristic Overlap Study

- Total heuristic issues: `34`
- Byte-span overlap with parse errors: `30`
- Same-row overlap with parse errors: `31`

| Example | Heuristic issues | Byte overlap | Row overlap |
| --- | --- | --- | --- |
| 00-valid-basic.json | 0 | 0 | 0 |
| 01-valid-nested-objects.json | 0 | 0 | 0 |
| 10-trailing-comma-object.json | 1 | 1 | 1 |
| 11-single-quotes.json | 3 | 3 | 3 |
| 12-unquoted-keys.json | 2 | 2 | 2 |
| 13-missing-comma.json | 0 | 0 | 0 |
| 14-missing-colon.json | 0 | 0 | 0 |
| 15-leading-prose-wrapper.json | 1 | 1 | 1 |
| 16-markdown-fence-wrapper.json | 1 | 1 | 1 |
| 17-python-literals.json | 3 | 3 | 3 |
| 18-duplicate-comma-array.json | 1 | 1 | 1 |
| 19-comments.json | 1 | 0 | 0 |
| 20-missing-closing-delimiter.json | 1 | 0 | 1 |
| 21-multiple-top-level-objects.json | 0 | 0 | 0 |
| 22-unterminated-string.json | 1 | 1 | 1 |
| 23-duplicate-key.json | 1 | 0 | 0 |
| 24-llm-wrapper-multi-step.json | 5 | 5 | 5 |
| 25-llm-commentary-comments-and-duplicate-comma.json | 4 | 3 | 3 |
| Leading prose | 1 | 1 | 1 |
| Markdown fence wrapper | 1 | 1 | 1 |
| Multiple top-level values | 0 | 0 | 0 |
| Python literals | 2 | 2 | 2 |
| Single quotes | 2 | 2 | 2 |
| Trailing comma | 1 | 1 | 1 |
| Unquoted keys | 2 | 2 | 2 |
| Valid JSON | 0 | 0 | 0 |

