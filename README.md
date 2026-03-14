# sanitize

Structured-text linter and heuristic fixer powered by [tree-sitter](https://tree-sitter.github.io/).

`sanitize` currently supports:

- YAML linting and iterative repair
- JSON linting and conservative recovery for common malformed LLM output
- a bundled web playground for parse-tree inspection and repair review

## Install

```bash
go install github.com/go-go-golems/sanitize/cmd/sanitize@latest
```

## CLI usage

```bash
# Fix a YAML file (sanitized output on stdout, fix log on stderr)
sanitize fix broken.yaml

# Pipe YAML from stdin
cat broken.yaml | sanitize fix

# Lint YAML only (no fixing) — exits 1 if issues found
sanitize lint broken.yaml

# Print the tree-sitter parse tree and structural errors for YAML
sanitize parse broken.yaml

# JSON output (full Result struct)
sanitize fix --json broken.yaml

# Custom tab width for YAML (default 2)
sanitize fix --tab-width 4 broken.yaml

# Lint malformed JSON
sanitize lint --format json broken.json

# Recover common LLM-style JSON wrappers and literals
sanitize fix --format json llm-output.txt

# List JSON rules
sanitize rules --format json

# Print JSON parse tree plus strict-parse status
sanitize parse --format json --json broken.json
```

## Web UI

```bash
sanitize serve            # http://localhost:8080
sanitize serve --port 3000 # http://localhost:3000
```

The server bundles a single-page web UI for interactive structured-text editing,
parse-tree inspection, and one-click sanitization.

The UI now supports both YAML and JSON. In JSON mode it acts as a recovery
playground for malformed LLM output, including fenced payloads, prose wrappers,
Python literals, comments, and trailing commas.

## Library usage

```go
import yamlsanitize "github.com/go-go-golems/sanitize/pkg/yaml"
import jsonsanitize "github.com/go-go-golems/sanitize/pkg/json"

// Full sanitize (lint + fix)
result := yamlsanitize.Sanitize(input,
    yamlsanitize.WithTabWidth(4),
    yamlsanitize.WithMaxIterations(5),
    yamlsanitize.WithRules("tab_indent", "missing_space_after_colon"),
)
fmt.Println(result.Sanitized)
fmt.Println(result.ParseClean, result.LintClean)

// Lint only
issues := yamlsanitize.Lint(input)

// Parse tree only
treeText, errors, err := yamlsanitize.ParseTree(input)

// Conservative JSON recovery
jsonResult := jsonsanitize.Sanitize(input)
```

## YAML rules

| Rule | Lint | Fix | Description |
|------|------|-----|-------------|
| `tab_indent` | yes | yes | Tab characters used for indentation |
| `missing_space_after_colon` | yes | yes | `key:value` instead of `key: value` |
| `list_dash_no_space` | yes | yes | `-item` instead of `- item` |
| `trailing_comma` | yes | yes | Trailing comma in flow collections `{a: 1,}` |
| `duplicate_key` | yes | yes | Duplicate sibling keys (renames with `_N` suffix) |
| `extra_colon_in_value` | yes | yes | Unquoted value contains `: ` (auto-quotes) |
| `mixed_indent` | — | yes | Inconsistent indentation depth (normalises) |

## JSON rules

| Rule | Lint | Fix | Description |
|------|------|-----|-------------|
| `markdown_fence_wrapper` | yes | yes | Markdown code fences wrapped around the JSON payload |
| `leading_or_trailing_prose` | yes | yes | Non-JSON prose before or after the payload |
| `single_quotes` | yes | yes | Single-quoted strings used where JSON requires double quotes |
| `unquoted_keys` | yes | no | Object keys are not quoted JSON strings |
| `python_literals` | yes | yes | Python-style literals such as `True`, `False`, and `None` |
| `trailing_comma` | yes | yes | Trailing comma in an object or array |
| `duplicate_comma` | yes | yes | Duplicate comma in an array or object member list |
| `comment` | yes | yes | Comment syntax present in JSON input |
| `multiple_top_level_values` | yes | no | Multiple top-level JSON values appear in one input |
| `missing_closing_delimiter` | yes | no | Likely missing closing brace or bracket |
| `duplicate_key` | yes | no | Duplicate key within the same object |
| `strict_parse_error` | yes | no | `encoding/json` strict parse error surfaced as lint |
| `structural_parse_error` | yes | no | Tree-sitter structural parse error surfaced as lint |
| `missing_syntax_node` | yes | no | Tree-sitter missing syntax node surfaced as lint |

## Options

- `WithMaxIterations(n)` — max fix iterations (default 10)
- `WithTabWidth(w)` — spaces per tab for `tab_indent` fixer (default 2)
- `WithOnlyRules(rules...)` — restrict to specific rules (default: all)
- `WithDisabledRules(rules...)` — disable specific rules
