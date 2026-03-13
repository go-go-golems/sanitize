---
Title: Imported JSON parse errors note
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
Summary: Imported source note capturing common malformed JSON cases from LLM-generated output.
WhatFor: Provide the raw external malformed JSON examples that the experiment scripts replay.
WhenToUse: Use when regenerating the JSON tree-sitter replication matrix or heuristic probe outputs.
---

# Imported JSON parse errors note

Common JSON parse errors from LLM-generated output:

1. **Trailing commas**

   ```json
   {"a": 1,}
   ```

   Error: unexpected `}` or invalid trailing comma.

2. **Single quotes instead of double quotes**

   ```json
   {'a': 1}
   ```

   Error: invalid JSON token.

3. **Unquoted keys**

   ```json
   {a: 1}
   ```

   Error: expected string key.

4. **Missing comma between fields**

   ```json
   {"a": 1 "b": 2}
   ```

   Error: expected `,` delimiter.

5. **Missing colon after key**

   ```json
   {"a" 1}
   ```

   Error: expected `:` after property name.

6. **Unescaped double quotes inside strings**

   ```json
   {"text": "He said "hello""}
   ```

   Error: unterminated string / unexpected token.

7. **Invalid escape sequences**

   ```json
   {"path": "C:\new\test"}
   ```

   Error: bad escape character like `\n`, `\t`, or invalid `\x`.

8. **Unterminated strings**

   ```json
   {"text": "hello}
   ```

   Error: unexpected end of input.

9. **Mismatched braces or brackets**

   ```json
   {"a": [1, 2}
   ```

   Error: expected `]` or `}`.

10. **Extra text before or after JSON**

```text
Here is your JSON:
{"a":1}
Thanks!
```

Error: unexpected token at position 0 or after valid JSON.

11. **Using comments**

```json
{"a": 1 // comment
}
```

Error: unexpected `/`.

12. **Multiple top-level objects**

```json
{"a":1} {"b":2}
```

Error: unexpected token after end of JSON input.

13. **Invalid numbers**

```json
{"x": NaN}
{"y": Infinity}
{"z": 01}
```

Error: invalid numeric literal.

14. **Using Python/JS literals instead of JSON literals**

```json
{"ok": True, "value": None}
```

Error: unexpected token `T` or `N`. JSON requires `true`, `false`, `null`.

15. **Markdown code fences included**

````
```json
{"a":1}
```
````

Error: parser sees backticks, not JSON.

16. **Ellipses or placeholders**

```json
{"items": [1, 2, ...]}
```

Error: unexpected token `.`

17. **Duplicate commas or malformed arrays**

```json
{"items": [1,,2]}
```

Error: unexpected `,`.

18. **Unicode/control character issues**
    Hidden newlines, smart quotes, or raw control chars in strings can trigger parse failures.

A useful way to think about it:

* **Parse errors**: JSON is not syntactically valid.
* **Validation errors**: JSON parses, but shape is wrong.
  Example:

  ```json
  {"age": "twenty"}
  ```

  This parses, but may fail schema validation if `age` must be a number.

For LLM output, the most common real-world causes are:

* extra prose around the JSON,
* trailing commas,
* single quotes,
* unescaped quotes,
* Python/JavaScript literals,
* truncated output.

I can also give you a “messy LLM JSON → parser error → fixed JSON” cheat sheet.
