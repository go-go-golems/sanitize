# JSON Example Corpus

These files are the JSON-side equivalent of `examples/yaml/`.

Filename convention:

- `00-09`: valid controls
- `10-19`: single malformed pattern
- `20-29`: mixed or more ambiguous malformed cases

The intent is to make CLI, server, and browser testing easier by keeping one
canonical file per common malformed LLM-generated JSON pattern.

Notable combined recovery cases:

- `24-llm-wrapper-multi-step.json`: prose wrapper + Markdown fence + Python
  literal + trailing commas. This exercises iterative recovery across several
  low-risk fixers.
- `25-llm-commentary-comments-and-duplicate-comma.json`: comments + duplicate
  comma + Python boolean + trailing comma. This is a more "LLM-ish" object that
  still stays within the current conservative fixer scope.
