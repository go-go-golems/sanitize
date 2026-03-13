# sanitize

YAML linter and heuristic fixer powered by [tree-sitter](https://tree-sitter.github.io/).

Parses YAML with tree-sitter, detects structural errors and common mistakes,
then iteratively applies fixes until the document is clean (or no more progress
can be made).

## Install

```bash
go install github.com/go-go-golems/sanitize/cmd/sanitize@latest
```

## CLI usage

```bash
# Fix a file (sanitized output on stdout, fix log on stderr)
sanitize fix broken.yaml

# Pipe from stdin
cat broken.yaml | sanitize fix

# Lint only (no fixing) — exits 1 if issues found
sanitize lint broken.yaml

# JSON output (full Result struct)
sanitize fix --json broken.yaml

# Custom tab width (default 2)
sanitize fix --tab-width 4 broken.yaml
```

## Web UI

```bash
sanitize serve            # http://localhost:8080
sanitize serve --port 3000 # http://localhost:3000
```

The server bundles a single-page web UI for interactive YAML editing, parse-tree
inspection, and one-click sanitization.

## Library usage

```go
import yamlsanitize "github.com/go-go-golems/sanitize/pkg/yaml"

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
```

## Supported rules

| Rule | Lint | Fix | Description |
|------|------|-----|-------------|
| `tab_indent` | yes | yes | Tab characters used for indentation |
| `missing_space_after_colon` | yes | yes | `key:value` instead of `key: value` |
| `list_dash_no_space` | yes | yes | `-item` instead of `- item` |
| `trailing_comma` | yes | yes | Trailing comma in flow collections `{a: 1,}` |
| `duplicate_key` | yes | yes | Duplicate sibling keys (renames with `_N` suffix) |
| `extra_colon_in_value` | yes | yes | Unquoted value contains `: ` (auto-quotes) |
| `mixed_indent` | — | yes | Inconsistent indentation depth (normalises) |

## Options

- `WithMaxIterations(n)` — max fix iterations (default 10)
- `WithTabWidth(w)` — spaces per tab for `tab_indent` fixer (default 2)
- `WithRules(rules...)` — restrict to specific rules (default: all)
