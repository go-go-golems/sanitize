---
Title: JSON parse error replication matrix
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
Summary: Replication matrix comparing strict encoding/json failures with tree-sitter JSON error nodes for the imported malformed LLM JSON cases.
WhatFor: Show which malformed JSON cases tree-sitter can localize structurally and where heuristics will still be required.
WhenToUse: Use when designing JSON parse-aware lint rules and fix heuristics.
---

# JSON Parse Error Replication Matrix

Source: `../sources/local/03-json-parse-errors-import.md`

## Case 1: Trailing commas

- Language tag: `json`
- Strict `encoding/json` valid: `false`
- Strict parse error: `invalid character '}' looking for beginning of object key string`
- Tree-sitter error node count: `1`
- First tree-sitter error: `ERROR` rows 0:10 to 0:11 text=","
- Source note: unexpected `}` or invalid trailing comma.
- Heuristic hits: `trailing_comma`

```text
   {"a": 1,}
```

```text
(document (object (pair key: (string (string_content)) value: (number)) (ERROR)))
```

## Case 2: Single quotes instead of double quotes

- Language tag: `json`
- Strict `encoding/json` valid: `false`
- Strict parse error: `invalid character '\'' looking for beginning of object key string`
- Tree-sitter error node count: `3`
- First tree-sitter error: `ERROR` rows 0:3 to 0:8 text="{'a':"
- Source note: invalid JSON token.
- Heuristic hits: `single_quotes`

```text
   {'a': 1}
```

```text
(document (ERROR (UNEXPECTED ''')) (number) (ERROR))
```

## Case 3: Unquoted keys

- Language tag: `json`
- Strict `encoding/json` valid: `false`
- Strict parse error: `invalid character 'a' looking for beginning of object key string`
- Tree-sitter error node count: `3`
- First tree-sitter error: `ERROR` rows 0:3 to 0:6 text="{a:"
- Source note: expected string key.
- Heuristic hits: `none`

```text
   {a: 1}
```

```text
(document (ERROR (UNEXPECTED 'a')) (number) (ERROR))
```

## Case 4: Missing comma between fields

- Language tag: `json`
- Strict `encoding/json` valid: `false`
- Strict parse error: `invalid character '"' after object key:value pair`
- Tree-sitter error node count: `1`
- First tree-sitter error: `ERROR` rows 0:4 to 0:10 text="\"a\": 1"
- Source note: expected `,` delimiter.
- Heuristic hits: `none`

```text
   {"a": 1 "b": 2}
```

```text
(document (object (ERROR (pair key: (string (string_content)) value: (number))) (pair key: (string (string_content)) value: (number))))
```

## Case 5: Missing colon after key

- Language tag: `json`
- Strict `encoding/json` valid: `false`
- Strict parse error: `invalid character '1' after object key`
- Tree-sitter error node count: `1`
- First tree-sitter error: `ERROR` rows 0:4 to 0:9 text="\"a\" 1"
- Source note: expected `:` after property name.
- Heuristic hits: `none`

```text
   {"a" 1}
```

```text
(document (object (ERROR (string (string_content)) (number))))
```

## Case 6: Unescaped double quotes inside strings

- Language tag: `json`
- Strict `encoding/json` valid: `false`
- Strict parse error: `invalid character 'h' after object key:value pair`
- Tree-sitter error node count: `3`
- First tree-sitter error: `ERROR` rows 0:3 to 0:30 text="{\"text\": \"He said \"hello\"\"}"
- Source note: unterminated string / unexpected token.
- Heuristic hits: `none`

```text
   {"text": "He said "hello""}
```

```text
(document (ERROR (ERROR (pair key: (string (string_content)) value: (string (string_content))) (UNEXPECTED 'h')) (string_content)))
```

## Case 7: Invalid escape sequences

- Language tag: `json`
- Strict `encoding/json` valid: `true`
- Tree-sitter error node count: `0`
- Source note: bad escape character like `\n`, `\t`, or invalid `\x`.
- Heuristic hits: `none`

```text
   {"path": "C:\new\test"}
```

```text
(document (object (pair key: (string (string_content)) value: (string (string_content) (escape_sequence) (string_content) (escape_sequence) (string_content)))))
```

## Case 8: Unterminated strings

- Language tag: `json`
- Strict `encoding/json` valid: `false`
- Strict parse error: `unexpected EOF`
- Tree-sitter error node count: `1`
- First tree-sitter error: `ERROR` rows 0:3 to 0:19 text="{\"text\": \"hello}"
- Source note: unexpected end of input.
- Heuristic hits: `none`

```text
   {"text": "hello}
```

```text
(ERROR (string (string_content)) (string_content))
```

## Case 9: Mismatched braces or brackets

- Language tag: `json`
- Strict `encoding/json` valid: `false`
- Strict parse error: `invalid character '}' after array element`
- Tree-sitter error node count: `1`
- First tree-sitter error: `MISSING` rows 0:14 to 0:14 text=""
- Source note: expected `]` or `}`.
- Heuristic hits: `none`

```text
   {"a": [1, 2}
```

```text
(document (object (pair key: (string (string_content)) value: (array (number) (number) (MISSING "]")))))
```

## Case 10: Extra text before or after JSON

- Language tag: `text`
- Strict `encoding/json` valid: `false`
- Strict parse error: `invalid character 'H' looking for beginning of value`
- Tree-sitter error node count: `4`
- First tree-sitter error: `ERROR` rows 0:0 to 0:18 text="Here is your JSON:"
- Source note: unexpected token at position 0 or after valid JSON.
- Heuristic hits: `prose_wrapped`

```text
Here is your JSON:
{"a":1}
Thanks!
```

```text
(document (ERROR (UNEXPECTED 'H')) (object (pair key: (string (string_content)) value: (number))) (ERROR (UNEXPECTED 'T')))
```

## Case 11: Using comments

- Language tag: `json`
- Strict `encoding/json` valid: `false`
- Strict parse error: `invalid character '/' after object key:value pair`
- Tree-sitter error node count: `0`
- Source note: unexpected `/`.
- Heuristic hits: `comment`

```text
{"a": 1 // comment
}
```

```text
(document (object (pair key: (string (string_content)) value: (number)) (comment)))
```

## Case 12: Multiple top-level objects

- Language tag: `json`
- Strict `encoding/json` valid: `false`
- Strict parse error: `multiple top-level JSON values`
- Tree-sitter error node count: `0`
- Source note: unexpected token after end of JSON input.
- Heuristic hits: `multiple_top_level`

```text
{"a":1} {"b":2}
```

```text
(document (object (pair key: (string (string_content)) value: (number))) (object (pair key: (string (string_content)) value: (number))))
```

## Case 13: Invalid numbers

- Language tag: `json`
- Strict `encoding/json` valid: `false`
- Strict parse error: `invalid character 'N' looking for beginning of value`
- Tree-sitter error node count: `5`
- First tree-sitter error: `ERROR` rows 0:6 to 0:10 text="NaN}"
- Source note: invalid numeric literal.
- Heuristic hits: `none`

```text
{"x": NaN}
{"y": Infinity}
{"z": 01}
```

```text
(document (object (pair key: (string (string_content)) (ERROR (UNEXPECTED 'N')) value: (object (pair key: (string (string_content)) (ERROR (UNEXPECTED 'I') (pair key: (string (string_content)) value: (number))) value: (number)))) (MISSING "}")))
```

## Case 14: Using Python/JS literals instead of JSON literals

- Language tag: `json`
- Strict `encoding/json` valid: `false`
- Strict parse error: `invalid character 'T' looking for beginning of value`
- Tree-sitter error node count: `4`
- First tree-sitter error: `ERROR` rows 0:7 to 0:12 text="True,"
- Source note: unexpected token `T` or `N`. JSON requires `true`, `false`, `null`.
- Heuristic hits: `python_literals`

```text
{"ok": True, "value": None}
```

```text
(document (object (pair key: (string (string_content)) (ERROR (UNEXPECTED 'T')) value: (string (string_content))) (ERROR (UNEXPECTED 'N'))))
```

## Case 15: Markdown code fences included

- Language tag: ``
- Strict `encoding/json` valid: `false`
- Strict parse error: `invalid character '`' looking for beginning of value`
- Tree-sitter error node count: `4`
- First tree-sitter error: `ERROR` rows 0:0 to 0:7 text="```json"
- Source note: parser sees backticks, not JSON.
- Heuristic hits: `markdown_fence`

````text
```json
{"a":1}
```
````

```text
(document (ERROR (UNEXPECTED '`')) (object (pair key: (string (string_content)) value: (number))) (ERROR (UNEXPECTED '`')))
```

## Case 16: Ellipses or placeholders

- Language tag: `json`
- Strict `encoding/json` valid: `false`
- Strict parse error: `invalid character '.' looking for beginning of value`
- Tree-sitter error node count: `2`
- First tree-sitter error: `ERROR` rows 0:15 to 0:20 text=", ..."
- Source note: unexpected token `.`
- Heuristic hits: `ellipsis_or_placeholder`

```text
{"items": [1, 2, ...]}
```

```text
(document (object (pair key: (string (string_content)) value: (array (number) (number) (ERROR (UNEXPECTED '.'))))))
```

## Case 17: Duplicate commas or malformed arrays

- Language tag: `json`
- Strict `encoding/json` valid: `false`
- Strict parse error: `invalid character ',' looking for beginning of value`
- Tree-sitter error node count: `1`
- First tree-sitter error: `ERROR` rows 0:12 to 0:13 text=","
- Source note: unexpected `,`.
- Heuristic hits: `duplicate_comma`

```text
{"items": [1,,2]}
```

```text
(document (object (pair key: (string (string_content)) value: (array (number) (ERROR) (number)))))
```

## Case 18: Unicode/control character issues

- Language tag: `json`
- Strict `encoding/json` valid: `true`
- Tree-sitter error node count: `0`
- Heuristic hits: `none`

```text
  {"age": "twenty"}
```

```text
(document (object (pair key: (string (string_content)) value: (string (string_content)))))
```

