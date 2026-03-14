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
