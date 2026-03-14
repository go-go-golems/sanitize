# JSON Example Corpus

These files are the JSON-side equivalent of `examples/yaml/`.

Filename convention:

- `00-09`: valid controls
- `10-19`: single malformed pattern
- `20-29`: mixed or more ambiguous malformed cases

The intent is to make CLI, server, and browser testing easier by keeping one
canonical file per common malformed LLM-generated JSON pattern.
