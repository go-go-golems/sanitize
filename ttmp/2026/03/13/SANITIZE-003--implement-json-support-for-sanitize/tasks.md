# Tasks

## Seed Work

- [x] Create ticket `SANITIZE-003`
- [x] Move `02-json-support-architecture-and-implementation-plan.md` from `SANITIZE-002` into this ticket

## Phase 1: Lock The JSON Shape

- [ ] Finalize the JSON MVP scope and confirm which fixer rules ship in v1 (`trailing_comma`, `single_quoted_string`, `unquoted_object_key`, `comment_in_json`, optional `missing_comma_between_members`)
- [ ] Finalize the shared result/config interfaces that move out of `pkg/yaml` into the format-agnostic core
- [ ] Finalize the breaking Glazed CLI verb tree (`parse`, `lint`, `fix`, `sanitize`) and format flag behavior
- [ ] Finalize the strict HTTP request/response contract using `format` + `input`
- [ ] Decide whether JSON tree rendering is required in the first release or deferred behind parse-only support

## Phase 2: Extract Shared Runtime

- [ ] Create `pkg/core/types.go` for shared `Format`, parse error, lint issue, fix, example, and result types
- [ ] Create `pkg/core/options.go` for shared config and option handling
- [ ] Create `pkg/core/engine.go` for the format-engine interface
- [ ] Move the iterative sanitize loop from `pkg/yaml/sanitize.go` into `pkg/core/run.go`
- [ ] Refactor the YAML implementation to become a format engine that plugs into the shared runtime
- [ ] Remove YAML-specific exported shapes that do not fit the new shared public surface
- [ ] Add shared-core tests that exercise the generic sanitize loop across engine implementations

## Phase 3: Rebuild The CLI With Glazed

- [ ] Create a Glazed/Cobra root command for `sanitize`
- [ ] Add Glazed verb commands for `parse`, `lint`, `fix`, and `sanitize`
- [ ] Define command settings structs with `glazed` tags and decode them from `schema.DefaultSlug`
- [ ] Add Glazed output sections and command settings sections consistently across all verbs
- [ ] Wire help and logging at the root command only
- [ ] Remove the old single-command CLI path and old ambiguous flags such as `--json`
- [ ] Add CLI subprocess tests covering verb routing, `--format`, Glazed output modes, and exit-code behavior

## Phase 4: Build The JSON Engine

- [ ] Create `pkg/json`
- [ ] Implement JSON parse validation using `encoding/json`
- [ ] Decide whether to add `tree-sitter-json` immediately or return empty tree output in the first iteration
- [ ] Implement JSON lint detection for the chosen MVP rules
- [ ] Implement JSON fixers conservatively for the chosen MVP rules
- [ ] Add JSON examples mirroring the YAML example style
- [ ] Add `pkg/json` unit tests for parse, lint, and fix behavior
- [ ] Add regression fixtures for edge cases where JSON fixes could accidentally change semantics

## Phase 5: Update Server And UI

- [ ] Refactor the server request types to accept only `format` + `input`
- [ ] Route server sanitize/parse handlers to the correct engine based on format
- [ ] Add handler tests for YAML success, JSON success, invalid format, and malformed bodies
- [ ] Update the web UI to expose format selection
- [ ] Add format-specific examples in the UI
- [ ] Update UI labels, placeholders, and render paths so JSON no longer appears as a YAML special case
- [ ] Decide how the parse-tree panel behaves when JSON tree output is unavailable

## Phase 6: End-To-End Validation

- [ ] Add cross-format end-to-end tests covering CLI, package, and HTTP paths
- [ ] Verify YAML behavior still works after the shared-core extraction
- [ ] Verify JSON behavior works in CLI text output and structured output modes
- [ ] Update README examples and command documentation to the new CLI shape
- [ ] Re-run `go test ./...`, `go test -race ./...`, `make lint`, and any JSON-specific smoke checks
