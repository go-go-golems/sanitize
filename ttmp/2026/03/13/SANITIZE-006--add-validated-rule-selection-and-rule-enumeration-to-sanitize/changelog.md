# Changelog

## 2026-03-13

- Initial workspace created


## 2026-03-13

Added the YAML rule registry and validated package-level rule selection (commit 25f17e9)

### Related Files

- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/lint.go — Option-aware lint entrypoint and rule filtering
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/options.go — Validated allowlist and disable-list config building
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/rules.go — Canonical rule metadata for rule enumeration
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/rules_test.go — Regression coverage for validation behavior
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/sanitize.go — Validated sanitize entrypoint


## 2026-03-13

Added Glazed CLI rule selection flags and the sanitize rules command (commit c7a6f53)

### Related Files

- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/cmd/sanitize/main_test.go — Regression coverage for CLI rule selection
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/internal/cli/commands.go — Added --rule/--disable-rule handling and rule listing
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/internal/cli/root.go — Registered sanitize rules in the root command

