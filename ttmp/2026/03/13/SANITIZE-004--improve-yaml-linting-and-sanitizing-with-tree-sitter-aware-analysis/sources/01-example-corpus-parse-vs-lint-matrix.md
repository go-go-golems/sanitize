---
Title: Parse vs Lint Matrix
Ticket: SANITIZE-004
Status: active
Topics:
    - yaml
    - linting
    - treesitter
DocType: reference
Intent: short-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: "Generated matrix comparing parse errors and lint issues across examples/yaml."
LastUpdated: 2026-03-13T09:13:52.978401783-04:00
WhatFor: "Generated evidence comparing parser and linter coverage across the YAML example corpus."
WhenToUse: "Use when validating tree-sitter-aware linting changes."
---

# Parse vs Lint Matrix

Generated from `examples/yaml` with `25` file(s).

| File | Parse errors | Lint issues | Parse rows | Lint rows | Rules |
|------|--------------|-------------|------------|-----------|-------|
| examples/yaml/00-valid-basic.yaml | 0 | 0 | — | — | — |
| examples/yaml/01-valid-nested-mapping.yaml | 0 | 0 | — | — | — |
| examples/yaml/02-valid-flow-collections.yaml | 0 | 0 | — | — | — |
| examples/yaml/03-valid-block-scalar.yaml | 0 | 0 | — | — | — |
| examples/yaml/04-valid-comments-and-blank-lines.yaml | 0 | 0 | — | — | — |
| examples/yaml/10-tab-indent.yaml | 1 | 2 | 1 | 2, 3 | tab_indent |
| examples/yaml/11-missing-space-after-colon.yaml | 0 | 3 | — | 1, 2, 3 | missing_space_after_colon |
| examples/yaml/12-missing-space-after-colon-url-safe.yaml | 0 | 2 | — | 1, 2 | missing_space_after_colon |
| examples/yaml/13-list-dash-no-space.yaml | 0 | 2 | — | 2, 3 | list_dash_no_space |
| examples/yaml/14-trailing-comma-flow-mapping.yaml | 0 | 2 | — | 1, 2 | trailing_comma |
| examples/yaml/15-trailing-comma-flow-sequence.yaml | 0 | 2 | — | 1, 2 | trailing_comma |
| examples/yaml/16-extra-colon-in-value.yaml | 1 | 2 | 1 | 1, 2 | extra_colon_in_value |
| examples/yaml/17-duplicate-key-sibling.yaml | 0 | 1 | — | 3 | duplicate_key |
| examples/yaml/18-quoted-colon-safe.yaml | 0 | 0 | — | — | — |
| examples/yaml/19-commented-colon-safe.yaml | 0 | 0 | — | — | — |
| examples/yaml/20-mixed-indent.yaml | 1 | 0 | 1 | — | — |
| examples/yaml/21-multiple-errors.yaml | 1 | 11 | 1 | 2, 3, 4, 5, 6, 7 | extra_colon_in_value, list_dash_no_space, missing_space_after_colon, tab_indent |
| examples/yaml/22-duplicate-key-different-parents.yaml | 0 | 0 | — | — | — |
| examples/yaml/23-duplicate-key-separate-sequence-items.yaml | 0 | 0 | — | — | — |
| examples/yaml/24-deeply-nested-mixed-errors.yaml | 1 | 12 | 1 | 2, 3, 4, 5, 6, 7 | extra_colon_in_value, list_dash_no_space, missing_space_after_colon, tab_indent, trailing_comma |
| examples/yaml/25-clean-duplicate-like-keys.yaml | 0 | 0 | — | — | — |
| examples/yaml/26-valid-sequence-of-maps.yaml | 0 | 0 | — | — | — |
| examples/yaml/27-extra-colon-needs-quotes.yaml | 1 | 2 | 1 | 1, 2 | extra_colon_in_value |
| examples/yaml/28-unresolved-parse-error.yaml | 2 | 0 | 1, 2 | — | — |
| examples/yaml/29-valid-url-values.yaml | 0 | 0 | — | — | — |

## Details

### examples/yaml/00-valid-basic.yaml

- Parse errors: none
- Lint issues: none

### examples/yaml/01-valid-nested-mapping.yaml

- Parse errors: none
- Lint issues: none

### examples/yaml/02-valid-flow-collections.yaml

- Parse errors: none
- Lint issues: none

### examples/yaml/03-valid-block-scalar.yaml

- Parse errors: none
- Lint issues: none

### examples/yaml/04-valid-comments-and-blank-lines.yaml

- Parse errors: none
- Lint issues: none

### examples/yaml/10-tab-indent.yaml

- Parse errors:
  - `ERROR` rows 1:1 to 1:9 text="service:"
- Lint issues:
  - `tab_indent` row 2: Line 2: tab used for indentation (YAML requires spaces)
  - `tab_indent` row 3: Line 3: tab used for indentation (YAML requires spaces)

### examples/yaml/11-missing-space-after-colon.yaml

- Parse errors: none
- Lint issues:
  - `missing_space_after_colon` row 1: Line 1: missing space after colon
  - `missing_space_after_colon` row 2: Line 2: missing space after colon
  - `missing_space_after_colon` row 3: Line 3: missing space after colon

### examples/yaml/12-missing-space-after-colon-url-safe.yaml

- Parse errors: none
- Lint issues:
  - `missing_space_after_colon` row 1: Line 1: missing space after colon
  - `missing_space_after_colon` row 2: Line 2: missing space after colon

### examples/yaml/13-list-dash-no-space.yaml

- Parse errors: none
- Lint issues:
  - `list_dash_no_space` row 2: Line 2: list dash not followed by a space
  - `list_dash_no_space` row 3: Line 3: list dash not followed by a space

### examples/yaml/14-trailing-comma-flow-mapping.yaml

- Parse errors: none
- Lint issues:
  - `trailing_comma` row 1: Line 1: trailing comma in flow collection
  - `trailing_comma` row 2: Line 2: trailing comma in flow collection

### examples/yaml/15-trailing-comma-flow-sequence.yaml

- Parse errors: none
- Lint issues:
  - `trailing_comma` row 1: Line 1: trailing comma in flow collection
  - `trailing_comma` row 2: Line 2: trailing comma in flow collection

### examples/yaml/16-extra-colon-in-value.yaml

- Parse errors:
  - `ERROR` rows 1:1 to 1:13 text="message: key"
- Lint issues:
  - `extra_colon_in_value` row 1: Line 1: plain scalar value contains a colon — may need quoting
  - `extra_colon_in_value` row 2: Line 2: plain scalar value contains a colon — may need quoting

### examples/yaml/17-duplicate-key-sibling.yaml

- Parse errors: none
- Lint issues:
  - `duplicate_key` row 3: Line 3: duplicate key 'timeout'

### examples/yaml/18-quoted-colon-safe.yaml

- Parse errors: none
- Lint issues: none

### examples/yaml/19-commented-colon-safe.yaml

- Parse errors: none
- Lint issues: none

### examples/yaml/20-mixed-indent.yaml

- Parse errors:
  - `ERROR` rows 1:1 to 3:11 text="root:\n  child_a: 1\n   child_b"
- Lint issues: none

### examples/yaml/21-multiple-errors.yaml

- Parse errors:
  - `ERROR` rows 1:1 to 1:9 text="service:"
- Lint issues:
  - `tab_indent` row 2: Line 2: tab used for indentation (YAML requires spaces)
  - `missing_space_after_colon` row 2: Line 2: missing space after colon
  - `tab_indent` row 3: Line 3: tab used for indentation (YAML requires spaces)
  - `missing_space_after_colon` row 3: Line 3: missing space after colon
  - `tab_indent` row 4: Line 4: tab used for indentation (YAML requires spaces)
  - `tab_indent` row 5: Line 5: tab used for indentation (YAML requires spaces)
  - `list_dash_no_space` row 5: Line 5: list dash not followed by a space
  - `tab_indent` row 6: Line 6: tab used for indentation (YAML requires spaces)
  - `list_dash_no_space` row 6: Line 6: list dash not followed by a space
  - `tab_indent` row 7: Line 7: tab used for indentation (YAML requires spaces)
  - `extra_colon_in_value` row 7: Line 7: plain scalar value contains a colon — may need quoting

### examples/yaml/22-duplicate-key-different-parents.yaml

- Parse errors: none
- Lint issues: none

### examples/yaml/23-duplicate-key-separate-sequence-items.yaml

- Parse errors: none
- Lint issues: none

### examples/yaml/24-deeply-nested-mixed-errors.yaml

- Parse errors:
  - `ERROR` rows 1:1 to 1:10 text="pipeline:"
- Lint issues:
  - `tab_indent` row 2: Line 2: tab used for indentation (YAML requires spaces)
  - `tab_indent` row 3: Line 3: tab used for indentation (YAML requires spaces)
  - `list_dash_no_space` row 3: Line 3: list dash not followed by a space
  - `tab_indent` row 4: Line 4: tab used for indentation (YAML requires spaces)
  - `list_dash_no_space` row 4: Line 4: list dash not followed by a space
  - `tab_indent` row 5: Line 5: tab used for indentation (YAML requires spaces)
  - `tab_indent` row 6: Line 6: tab used for indentation (YAML requires spaces)
  - `extra_colon_in_value` row 6: Line 6: plain scalar value contains a colon — may need quoting
  - `tab_indent` row 7: Line 7: tab used for indentation (YAML requires spaces)
  - `missing_space_after_colon` row 7: Line 7: missing space after colon
  - `trailing_comma` row 7: Line 7: trailing comma in flow collection
  - `extra_colon_in_value` row 7: Line 7: plain scalar value contains a colon — may need quoting

### examples/yaml/25-clean-duplicate-like-keys.yaml

- Parse errors: none
- Lint issues: none

### examples/yaml/26-valid-sequence-of-maps.yaml

- Parse errors: none
- Lint issues: none

### examples/yaml/27-extra-colon-needs-quotes.yaml

- Parse errors:
  - `ERROR` rows 1:1 to 1:16 text="alert: severity"
- Lint issues:
  - `extra_colon_in_value` row 1: Line 1: plain scalar value contains a colon — may need quoting
  - `extra_colon_in_value` row 2: Line 2: plain scalar value contains a colon — may need quoting

### examples/yaml/28-unresolved-parse-error.yaml

- Parse errors:
  - `ERROR` rows 1:1 to 4:18 text="broken:\n  - one\n  - two\n  trailing: [1, 2"
  - `ERROR` rows 2:3 to 3:8 text="- one\n  - two"
- Lint issues: none

### examples/yaml/29-valid-url-values.yaml

- Parse errors: none
- Lint issues: none

