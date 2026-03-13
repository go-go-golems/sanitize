# Tasks

## Phase 1: Scope and parser choice

- [ ] Confirm whether JSON support should target strict RFC 8259 JSON only or also tolerate common LLM deviations during lint/fix.
- [ ] Choose the structural parser strategy for JSON analysis and recovery.
- [ ] Define the initial supported JSON rule set and output contract.

## Phase 2: LLM-oriented error taxonomy

- [x] Capture the common malformed-JSON patterns produced by LLMs.
- [ ] Map each common malformed pattern to whether it should be lint-only, auto-fixable, or suggestion-only.
- [ ] Decide which malformed patterns are too ambiguous to auto-rewrite.

## Phase 3: Shared engine and package design

- [ ] Design a `pkg/json` package aligned with the existing YAML package shape where practical.
- [ ] Decide whether to share the current analysis/fix pipeline across YAML and JSON or keep separate engines behind a common CLI surface.
- [ ] Define the JSON equivalents of parse errors, lint issues, and fix records.

## Phase 4: CLI and API integration

- [ ] Add explicit format selection to the CLI and API surface.
- [ ] Add JSON-oriented parse inspection tooling similar to `sanitize parse`.
- [ ] Add JSON examples that cover the common LLM-generated malformed cases.

## Phase 5: Validation and docs

- [ ] Expand the ticket design with a concrete implementation plan once the parser strategy is chosen.
- [ ] Add implementation diary entries as JSON support work begins.
- [ ] Run `docmgr doctor --ticket SANITIZE-007 --stale-after 30`.
