# YAML Example Corpus

This directory contains small YAML inputs for manual testing.

The files are grouped by intent:

- `00-09`: clean/control cases
- `10-19`: single-rule broken cases
- `20-29`: mixed or regression cases

Useful commands:

```bash
sanitize lint examples/yaml/11-missing-space-after-colon.yaml
sanitize fix examples/yaml/21-multiple-errors.yaml
sanitize fix --json examples/yaml/17-duplicate-key-sibling.yaml
sanitize serve
```

Notes:

- Some files are intentionally valid and should stay clean.
- Some files are intentionally broken but fixable.
- `28-unresolved-parse-error.yaml` is intentionally not cleanly recoverable and is useful for checking error reporting.
