---
Title: JSON repair matrix
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
Summary: End-to-end repair matrix for JSON examples, including parser signals, lint rules, fixes, and final strict-parse status.
WhatFor: Show which malformed JSON cases are recoverable today and which remain lint-only.
WhenToUse: Use when reviewing JSON recovery coverage or selecting the next rule/fixer to implement.
---

# JSON Repair Matrix

| Example | Pattern | Original strict | Parse nodes | Rules | Fixes | Final strict |
| --- | --- | --- | --- | --- | --- | --- |
| Leading prose | `leading_or_trailing_prose` | `false` | `2` | `leading_or_trailing_prose, structural_parse_error` | `leading_or_trailing_prose` | `true` |
| Markdown fence wrapper | `markdown_fence_wrapper` | `false` | `4` | `markdown_fence_wrapper, structural_parse_error` | `markdown_fence_wrapper` | `true` |
| Multiple top-level values | `multiple_top_level_values` | `false` | `0` | `multiple_top_level_values` | `none` | `false` |
| Python literals | `python_literals` | `false` | `4` | `python_literals, structural_parse_error` | `python_literals` | `true` |
| Single quotes | `single_quotes` | `false` | `3` | `single_quotes, structural_parse_error` | `none` | `false` |
| Trailing comma | `trailing_comma` | `false` | `1` | `structural_parse_error, trailing_comma` | `trailing_comma` | `true` |
| Unquoted keys | `unquoted_keys` | `false` | `4` | `structural_parse_error, unquoted_keys` | `none` | `false` |
| Valid JSON | `valid` | `true` | `0` | `none` | `none` | `true` |
| 00-valid-basic.json | `valid_basic` | `true` | `0` | `none` | `none` | `true` |
| 01-valid-nested-objects.json | `valid_nested_objects` | `true` | `0` | `none` | `none` | `true` |
| 10-trailing-comma-object.json | `trailing_comma` | `false` | `1` | `structural_parse_error, trailing_comma` | `trailing_comma` | `true` |
| 11-single-quotes.json | `single_quotes` | `false` | `5` | `single_quotes, structural_parse_error` | `none` | `false` |
| 12-unquoted-keys.json | `unquoted_keys` | `false` | `4` | `structural_parse_error, unquoted_keys` | `none` | `false` |
| 13-missing-comma.json | `missing_comma` | `false` | `1` | `structural_parse_error` | `none` | `false` |
| 14-missing-colon.json | `missing_colon` | `false` | `1` | `structural_parse_error` | `none` | `false` |
| 15-leading-prose-wrapper.json | `leading_or_trailing_prose` | `false` | `2` | `leading_or_trailing_prose, structural_parse_error` | `leading_or_trailing_prose` | `true` |
| 16-markdown-fence-wrapper.json | `markdown_fence_wrapper` | `false` | `4` | `markdown_fence_wrapper, structural_parse_error` | `markdown_fence_wrapper` | `true` |
| 17-python-literals.json | `python_literals` | `false` | `5` | `python_literals, structural_parse_error` | `python_literals` | `true` |
| 18-duplicate-comma-array.json | `duplicate_comma` | `false` | `1` | `duplicate_comma, structural_parse_error` | `duplicate_comma` | `true` |
| 19-comments.json | `comment` | `false` | `0` | `comment, strict_parse_error` | `comment` | `true` |
| 20-missing-closing-delimiter.json | `missing_closing_delimiter` | `false` | `1` | `missing_closing_delimiter, missing_syntax_node` | `none` | `false` |
| 21-multiple-top-level-objects.json | `multiple_top_level_values` | `false` | `0` | `multiple_top_level_values` | `none` | `false` |
| 22-unterminated-string.json | `unterminated_string` | `false` | `1` | `missing_closing_delimiter, structural_parse_error` | `none` | `false` |
| 23-duplicate-key.json | `duplicate_key` | `true` | `0` | `duplicate_key` | `none` | `true` |
| 24-llm-wrapper-multi-step.json | `wrapper_python_trailing_comma_combo` | `false` | `6` | `leading_or_trailing_prose, python_literals, structural_parse_error, trailing_comma` | `leading_or_trailing_prose, python_literals, trailing_comma` | `true` |
| 25-llm-commentary-comments-and-duplicate-comma.json | `comment_python_duplicate_comma_combo` | `false` | `3` | `comment, duplicate_comma, python_literals, structural_parse_error, trailing_comma` | `comment, duplicate_comma, python_literals, trailing_comma` | `true` |

## Leading prose

- Pattern: `leading_or_trailing_prose`
- Original strict parse clean: `false`
- Original tree-sitter error count: `2`
- Original rules: `leading_or_trailing_prose, structural_parse_error`
- Applied fixes: `leading_or_trailing_prose`
- Final strict parse clean: `true`
- Final lint clean: `true`

```text
{"name":"Alice"}
```

## Markdown fence wrapper

- Pattern: `markdown_fence_wrapper`
- Original strict parse clean: `false`
- Original tree-sitter error count: `4`
- Original rules: `markdown_fence_wrapper, structural_parse_error`
- Applied fixes: `markdown_fence_wrapper`
- Final strict parse clean: `true`
- Final lint clean: `true`

```text
{"name":"Alice"}
```

## Multiple top-level values

- Pattern: `multiple_top_level_values`
- Original strict parse clean: `false`
- Original tree-sitter error count: `0`
- Original rules: `multiple_top_level_values`
- Applied fixes: `none`
- Final strict parse clean: `false`
- Final lint clean: `false`

```text
{"a":1} {"b":2}
```

## Python literals

- Pattern: `python_literals`
- Original strict parse clean: `false`
- Original tree-sitter error count: `4`
- Original rules: `python_literals, structural_parse_error`
- Applied fixes: `python_literals`
- Final strict parse clean: `true`
- Final lint clean: `true`

```text
{"ok": true, "value": null}
```

## Single quotes

- Pattern: `single_quotes`
- Original strict parse clean: `false`
- Original tree-sitter error count: `3`
- Original rules: `single_quotes, structural_parse_error`
- Applied fixes: `none`
- Final strict parse clean: `false`
- Final lint clean: `false`

```text
{'name': 'Alice'}
```

## Trailing comma

- Pattern: `trailing_comma`
- Original strict parse clean: `false`
- Original tree-sitter error count: `1`
- Original rules: `structural_parse_error, trailing_comma`
- Applied fixes: `trailing_comma`
- Final strict parse clean: `true`
- Final lint clean: `true`

```text
{"name":"Alice"}
```

## Unquoted keys

- Pattern: `unquoted_keys`
- Original strict parse clean: `false`
- Original tree-sitter error count: `4`
- Original rules: `structural_parse_error, unquoted_keys`
- Applied fixes: `none`
- Final strict parse clean: `false`
- Final lint clean: `false`

```text
{name: "Alice", age: 30}
```

## Valid JSON

- Pattern: `valid`
- Original strict parse clean: `true`
- Original tree-sitter error count: `0`
- Original rules: `none`
- Applied fixes: `none`
- Final strict parse clean: `true`
- Final lint clean: `true`

```text
{"name":"Alice","age":30,"tags":["ops","infra"]}
```

## 00-valid-basic.json

- Pattern: `valid_basic`
- Original strict parse clean: `true`
- Original tree-sitter error count: `0`
- Original rules: `none`
- Applied fixes: `none`
- Final strict parse clean: `true`
- Final lint clean: `true`

```text
{"name":"Alice","age":30}
```

## 01-valid-nested-objects.json

- Pattern: `valid_nested_objects`
- Original strict parse clean: `true`
- Original tree-sitter error count: `0`
- Original rules: `none`
- Applied fixes: `none`
- Final strict parse clean: `true`
- Final lint clean: `true`

```text
{"user":{"name":"Alice","roles":["admin","operator"]},"active":true}
```

## 10-trailing-comma-object.json

- Pattern: `trailing_comma`
- Original strict parse clean: `false`
- Original tree-sitter error count: `1`
- Original rules: `structural_parse_error, trailing_comma`
- Applied fixes: `trailing_comma`
- Final strict parse clean: `true`
- Final lint clean: `true`

```text
{"name":"Alice","age":30}
```

## 11-single-quotes.json

- Pattern: `single_quotes`
- Original strict parse clean: `false`
- Original tree-sitter error count: `5`
- Original rules: `single_quotes, structural_parse_error`
- Applied fixes: `none`
- Final strict parse clean: `false`
- Final lint clean: `false`

```text
{'name':'Alice','age':30}
```

## 12-unquoted-keys.json

- Pattern: `unquoted_keys`
- Original strict parse clean: `false`
- Original tree-sitter error count: `4`
- Original rules: `structural_parse_error, unquoted_keys`
- Applied fixes: `none`
- Final strict parse clean: `false`
- Final lint clean: `false`

```text
{name:"Alice",age:30}
```

## 13-missing-comma.json

- Pattern: `missing_comma`
- Original strict parse clean: `false`
- Original tree-sitter error count: `1`
- Original rules: `structural_parse_error`
- Applied fixes: `none`
- Final strict parse clean: `false`
- Final lint clean: `false`

```text
{"name":"Alice" "age":30}
```

## 14-missing-colon.json

- Pattern: `missing_colon`
- Original strict parse clean: `false`
- Original tree-sitter error count: `1`
- Original rules: `structural_parse_error`
- Applied fixes: `none`
- Final strict parse clean: `false`
- Final lint clean: `false`

```text
{"name" "Alice"}
```

## 15-leading-prose-wrapper.json

- Pattern: `leading_or_trailing_prose`
- Original strict parse clean: `false`
- Original tree-sitter error count: `2`
- Original rules: `leading_or_trailing_prose, structural_parse_error`
- Applied fixes: `leading_or_trailing_prose`
- Final strict parse clean: `true`
- Final lint clean: `true`

```text
{"name":"Alice","age":30}
```

## 16-markdown-fence-wrapper.json

- Pattern: `markdown_fence_wrapper`
- Original strict parse clean: `false`
- Original tree-sitter error count: `4`
- Original rules: `markdown_fence_wrapper, structural_parse_error`
- Applied fixes: `markdown_fence_wrapper`
- Final strict parse clean: `true`
- Final lint clean: `true`

```text
{"name":"Alice","age":30}
```

## 17-python-literals.json

- Pattern: `python_literals`
- Original strict parse clean: `false`
- Original tree-sitter error count: `5`
- Original rules: `python_literals, structural_parse_error`
- Applied fixes: `python_literals`
- Final strict parse clean: `true`
- Final lint clean: `true`

```text
{"ok": true, "value": null, "enabled": false}
```

## 18-duplicate-comma-array.json

- Pattern: `duplicate_comma`
- Original strict parse clean: `false`
- Original tree-sitter error count: `1`
- Original rules: `duplicate_comma, structural_parse_error`
- Applied fixes: `duplicate_comma`
- Final strict parse clean: `true`
- Final lint clean: `true`

```text
{"items":[1,2]}
```

## 19-comments.json

- Pattern: `comment`
- Original strict parse clean: `false`
- Original tree-sitter error count: `0`
- Original rules: `comment, strict_parse_error`
- Applied fixes: `comment`
- Final strict parse clean: `true`
- Final lint clean: `true`

```text
{"name":"Alice" 
}
```

## 20-missing-closing-delimiter.json

- Pattern: `missing_closing_delimiter`
- Original strict parse clean: `false`
- Original tree-sitter error count: `1`
- Original rules: `missing_closing_delimiter, missing_syntax_node`
- Applied fixes: `none`
- Final strict parse clean: `false`
- Final lint clean: `false`

```text
{"name":"Alice","tags":["ops","infra"}
```

## 21-multiple-top-level-objects.json

- Pattern: `multiple_top_level_values`
- Original strict parse clean: `false`
- Original tree-sitter error count: `0`
- Original rules: `multiple_top_level_values`
- Applied fixes: `none`
- Final strict parse clean: `false`
- Final lint clean: `false`

```text
{"a":1} {"b":2}
```

## 22-unterminated-string.json

- Pattern: `unterminated_string`
- Original strict parse clean: `false`
- Original tree-sitter error count: `1`
- Original rules: `missing_closing_delimiter, structural_parse_error`
- Applied fixes: `none`
- Final strict parse clean: `false`
- Final lint clean: `false`

```text
{"text":"hello}
```

## 23-duplicate-key.json

- Pattern: `duplicate_key`
- Original strict parse clean: `true`
- Original tree-sitter error count: `0`
- Original rules: `duplicate_key`
- Applied fixes: `none`
- Final strict parse clean: `true`
- Final lint clean: `false`

```text
{"timeout":30,"timeout":60,"retries":3}
```

## 24-llm-wrapper-multi-step.json

- Pattern: `wrapper_python_trailing_comma_combo`
- Original strict parse clean: `false`
- Original tree-sitter error count: `6`
- Original rules: `leading_or_trailing_prose, python_literals, structural_parse_error, trailing_comma`
- Applied fixes: `leading_or_trailing_prose, python_literals, trailing_comma`
- Final strict parse clean: `true`
- Final lint clean: `true`

```text
{"ok": true, "items": [1,2]}
```

## 25-llm-commentary-comments-and-duplicate-comma.json

- Pattern: `comment_python_duplicate_comma_combo`
- Original strict parse clean: `false`
- Original tree-sitter error count: `3`
- Original rules: `comment, duplicate_comma, python_literals, structural_parse_error, trailing_comma`
- Applied fixes: `comment, duplicate_comma, python_literals, trailing_comma`
- Final strict parse clean: `true`
- Final lint clean: `true`

```text
{
  
  "items": [1,2],
  "active": false}
```

