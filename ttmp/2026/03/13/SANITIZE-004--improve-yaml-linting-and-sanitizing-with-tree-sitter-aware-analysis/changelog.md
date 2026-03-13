# Changelog

## 2026-03-13

- Initial workspace created


## 2026-03-13

Documented the tree-sitter-aware linting design, intern guide, and evidence matrix

### Related Files

- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-004--improve-yaml-linting-and-sanitizing-with-tree-sitter-aware-analysis/design-doc/01-tree-sitter-driven-yaml-linting-and-sanitizing-analysis.md — Primary design deliverable
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-004--improve-yaml-linting-and-sanitizing-with-tree-sitter-aware-analysis/reference/01-diary.md — Chronological record of the work
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-004--improve-yaml-linting-and-sanitizing-with-tree-sitter-aware-analysis/reference/02-intern-guide-to-tree-sitter-aware-linting-in-sanitize.md — Detailed onboarding and implementation guide
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-004--improve-yaml-linting-and-sanitizing-with-tree-sitter-aware-analysis/sources/01-example-corpus-parse-vs-lint-matrix.md — Generated evidence used by the design


## 2026-03-13

Added parse inspection tooling and corpus experiment assets (commit 63a3d07)

### Related Files

- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/cmd/sanitize/main_test.go — Added parse command coverage
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/internal/cli/commands.go — Added sanitize parse command
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/internal/cli/root.go — Registered parse command
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-004--improve-yaml-linting-and-sanitizing-with-tree-sitter-aware-analysis/scripts/parse_lint_matrix.go — Added corpus experiment script


## 2026-03-13

Validated SANITIZE-004, fixed source frontmatter/vocabulary hygiene, and uploaded the bundle to reMarkable

### Related Files

- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-004--improve-yaml-linting-and-sanitizing-with-tree-sitter-aware-analysis/reference/01-diary.md — Recorded validation and upload outcomes
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-004--improve-yaml-linting-and-sanitizing-with-tree-sitter-aware-analysis/scripts/parse_lint_matrix.go — Emit docmgr frontmatter for generated matrix output
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/vocabulary.yaml — Added treesitter topic for ticket validation


## 2026-03-13

Added structural lint diagnostics and reused analysis during sanitize (commit 273d0ef)

### Related Files

- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/fix.go — Span-aware parse-row coverage and shared duplicate-key analysis
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/lint.go — Parser-derived lint issues and span-aware heuristic diagnostics
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/sanitize.go — One shared analysis per sanitize iteration
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/sanitize_test.go — Coverage for structural lint issues and spans
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/types.go — Richer lint issue shape with spans and source


## 2026-03-13

Implemented the Phase 1 shared analysis core (commit 1ceb01c)

### Related Files

- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/analysis.go — Shared analysis object and parser reuse
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/duplicate_keys.go — Duplicate-key extraction moved onto shared analysis
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/line_index.go — Reusable line index for byte-to-row translation
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/parse.go — ParseTree now wraps shared analysis

