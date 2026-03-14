# Changelog

## 2026-03-13

- Initial workspace created


## 2026-03-13

Imported malformed JSON cases and added tree-sitter/heuristic experiment scripts

### Related Files

- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/scripts/json-error-matrix/main.go — Replication matrix generator
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/scripts/json-heuristic-probe/main.go — Heuristic probe generator
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/scripts/json-tree-sitter-parse/main.go — Tree-sitter JSON inspection helper
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/01-json-parse-error-replication-matrix.md — Generated tree-sitter replication matrix
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/02-json-heuristic-probe.md — Generated heuristic hit report
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/local/json-parse-errors.md — Imported source note


## 2026-03-13

Added a detailed intern-oriented analysis, design, and implementation guide for JSON support

### Related Files

- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/reference/02-intern-guide-to-json-support-and-tree-sitter-aware-malformed-json-recovery.md — Detailed implementation guide for a new contributor
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/design-doc/01-json-support-outline-and-malformed-llm-json-error-taxonomy.md — Primary JSON-support design document
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/01-json-parse-error-replication-matrix.md — Structural evidence used in the guide
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/02-json-heuristic-probe.md — Heuristic evidence used in the guide


## 2026-03-13

Expanded the JSON-support backlog into a phased implementation plan that explicitly covers the CLI, HTTP API, browser UI, corpus work, and validation strategy (commit `527d058`)

### Related Files

- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/tasks.md — Detailed phased backlog for the JSON implementation
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/index.md — Ticket index updated to surface the diary


## 2026-03-13

Extended the implementation guide and design doc to include the existing server/UI architecture and a JSON parse playground plan for the embedded browser app (commit `f78529a`)

### Related Files

- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/reference/02-intern-guide-to-json-support-and-tree-sitter-aware-malformed-json-recovery.md — Intern guide updated with server/API/UI guidance
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/design-doc/01-json-support-outline-and-malformed-llm-json-error-taxonomy.md — Design doc updated to cover the browser playground
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/internal/server/server.go — Current YAML-only API contract referenced by the new plan
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/internal/server/static/js/app.js — Current browser analysis flow referenced by the new plan


## 2026-03-13

Added the first `pkg/json` foundation with tree-sitter parsing, strict parsing, duplicate-key analysis, rule metadata, and option validation (commit `38669df`)

### Related Files

- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/json/analysis.go — Shared JSON document analysis
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/json/parse.go — JSON tree-sitter parse and strict parser entrypoints
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/json/duplicate_keys.go — Duplicate-key extraction for JSON objects
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/json/rules.go — JSON rule catalog
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/json/options.go — JSON rule-selection validation
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/json/parse_test.go — Parse and duplicate-key tests


## 2026-03-13

Added a file-backed JSON malformed corpus plus loader support in the shared examples package (commit `0e8e0d3`)

### Related Files

- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/examples/examples.go — JSON file loader added alongside the YAML loader
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/examples/examples_test.go — Loader tests for YAML and JSON examples
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/examples/json/README.md — JSON corpus conventions
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/examples/json/10-trailing-comma-object.json — Representative malformed JSON example
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/examples/json/16-markdown-fence-wrapper.json — Representative LLM-style wrapper case
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/json/examples.go — Built-in JSON examples for future CLI/server use


## 2026-03-13

Added parse-aware JSON lint diagnostics that surface tree-sitter errors, strict-parser failures, and duplicate keys (commit `5ac5f6c`)

### Related Files

- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/json/lint.go — JSON lint assembly and diagnostics
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/json/parse.go — Strict-parser multi-value error handling used by linting
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/json/rules.go — `strict_parse_error` added to the JSON rule catalog
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/json/parse_test.go — Lint coverage for duplicate keys and strict parse failures
