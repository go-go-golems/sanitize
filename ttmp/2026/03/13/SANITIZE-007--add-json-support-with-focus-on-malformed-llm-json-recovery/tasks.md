# Tasks

## Phase 1: Scope, product contract, and release boundary

- [ ] Lock the JSON scope to "strict final output, tolerant recovery for common LLM malformed input".
- [ ] Write the acceptance contract for JSON sanitize results.
  - Result should state whether input was structurally parseable, which rules fired, and which fixes were applied.
- [ ] Decide which malformed patterns are first-release in-scope.
  - Minimum candidates: trailing commas, single quotes, unquoted keys, comments, prose wrappers, markdown fences, Python literals, duplicate commas, truncated closing delimiters.
- [ ] Decide which malformed patterns are first-release out-of-scope or suggestion-only.
  - Likely candidates: ambiguous missing commas inside long arrays, invalid escape reconstruction, semantic duplicate-key merges.
- [ ] Define the success bar for "release-worthy JSON support".
  - Package API stable enough to use directly.
  - CLI supports lint/fix/parse/rules for JSON.
  - HTTP API and browser UI support JSON explicitly rather than treating it as YAML text.

## Phase 2: Malformed JSON taxonomy and evidence

- [x] Capture the common malformed-JSON patterns produced by LLMs.
- [x] Import the external `json-parse-errors.md` note into ticket-local sources.
- [x] Generate an initial tree-sitter replication matrix for the malformed JSON cases.
- [x] Generate an initial heuristic-probe report to show which malformed cases are detectable without a full parser.
- [ ] Turn the malformed-case list into a rule matrix.
  - For each case, record: parser signal, heuristic signal, fix confidence, auto-fixability, and test fixture name.
- [ ] Add severity guidance per malformed pattern.
  - `error`: structural parse blockers.
  - `warn`: recoverable wrappers or lenient-normalization patterns.
  - `info` or suggestion-only: ambiguous cleanup patterns that should not rewrite automatically.
- [ ] Add a "confidence policy" table for fixes.
  - Safe auto-fix.
  - Auto-fix only when parser span overlaps the heuristic span.
  - Suggestion only.
- [ ] Expand the ticket-local corpus under `scripts/` or `sources/` with one canonical example per malformed pattern plus expected repaired output.

## Phase 3: Shared architecture and `pkg/json` design

- [ ] Design `pkg/json` to mirror the mature parts of `pkg/yaml` without copying YAML-specific assumptions.
- [ ] Decide the internal analysis object for JSON.
  - Strict parse status.
  - Tree-sitter parse tree text.
  - Tree-sitter error spans.
  - Line index.
  - Duplicate-key findings.
  - Heuristic signals such as code fences, prose wrappers, and comment ranges.
- [ ] Decide whether to extract a cross-format internal analysis layer or keep YAML and JSON analysis packages separate behind shared CLI/server contracts.
- [ ] Define JSON-native issue and fix models.
  - Reuse span fields and source labels where possible.
  - Avoid YAML-specific field names in any shared surface.
- [ ] Define the JSON rule registry.
  - Canonical rule name.
  - Summary.
  - Whether it lints.
  - Whether it fixes.
  - Whether it is parse-aware.
  - Whether it is default-enabled.
- [ ] Write the iteration strategy for `Sanitize`.
  - One conservative fix batch.
  - Re-analyze.
  - Stop on convergence or iteration cap.
- [ ] Define how strict parser errors and tree-sitter errors are reconciled.
  - Strict parser should remain the final truth for "valid JSON".
  - Tree-sitter should supply localization and partial structure.

## Phase 4: JSON rule implementation plan

- [ ] Create an ordered first-release rule list.
  - `markdown_fence_wrapper`
  - `leading_or_trailing_prose`
  - `single_quotes`
  - `unquoted_keys`
  - `python_literals`
  - `trailing_comma`
  - `duplicate_comma`
  - `missing_closing_delimiter`
  - `duplicate_key`
- [ ] For each rule, define:
  - detection strategy,
  - parser/tree-sitter dependencies,
  - safety level,
  - exact fix transform,
  - regression fixtures.
- [ ] Split rule work into low-risk and high-risk buckets.
  - Low-risk first: wrappers, literals, trailing commas, duplicate commas.
  - Higher-risk later: unquoted keys and missing closing delimiters.
- [ ] Decide whether duplicate keys are lint-only or optionally rewritten with explicit naming like YAML.
- [ ] Decide whether there should be a "recovery only" preset for LLM JSON versus a stricter default preset.

## Phase 5: CLI integration

- [ ] Extend the Glazed CLI to support explicit format selection.
  - `sanitize lint --format json`
  - `sanitize fix --format json`
  - `sanitize parse --format json`
- [ ] Decide whether format-specific subcommands are needed in addition to `--format`.
- [ ] Add JSON rule enumeration to `sanitize rules`.
- [ ] Add JSON examples and sample commands to CLI help output.
- [ ] Add JSON-oriented parse inspection output that remains readable when the tree contains many `ERROR` nodes.
- [ ] Add CLI tests for format selection, unknown rule validation, parse output, and JSON fix output.

## Phase 6: HTTP API and browser UI

- [ ] Replace the YAML-only request body with an explicit format-aware request contract.
  - Proposed request shape: `{ "format": "yaml" | "json", "input": "..." }`
- [ ] Update `/api/sanitize` and `/api/parse` to dispatch by format.
- [ ] Update `/api/examples` to return format metadata so the browser can filter examples cleanly.
- [ ] Keep backward-compatibility out unless product requirements change.
- [ ] Update server tests for format-aware request/response flows.
- [ ] Convert the current browser UI from "YAML Sanitizer" into a format-aware structured-text playground.
- [ ] Add a JSON parse playground mode in the UI.
  - Format switcher near the example picker.
  - JSON-specific placeholder text.
  - JSON-specific example subset.
  - Parse-tree pane still available.
  - Issue pane distinguishes structural parse errors, heuristic lint issues, and applied fixes.
- [ ] Add a dedicated "recovery playground" experience for malformed LLM JSON.
  - Users should be able to paste prose-wrapped or fenced JSON and immediately see extraction and repair results.
- [ ] Decide how the parse-tree panel behaves in JSON mode.
  - Show raw tree-sitter tree.
  - Show strict-parse status.
  - Show original vs sanitized tree when both are available.
- [ ] Update UI labels and copy so the app no longer reads as YAML-specific.
- [ ] Add browser-side state tests or end-to-end tests for the JSON playground mode.

## Phase 7: Examples, corpus, and experiments

- [ ] Add `examples/json/` with one file per malformed pattern plus a README describing what each case is meant to exercise.
- [ ] Add generated or curated "LLM-ish" JSON cases that combine multiple failure modes in one sample.
- [ ] Add scripts under the ticket to compare:
  - strict parser failures,
  - tree-sitter error spans,
  - rule hits,
  - repaired output,
  - final strict-parse success.
- [ ] Add a script to dump UI-ready example metadata from the corpus so server/UI fixtures do not drift.
- [ ] Add at least one experiment that measures how often tree-sitter spans overlap the heuristic spans for the first-release rules.

## Phase 8: Validation and release readiness

- [ ] Add package-level tests for each JSON rule.
- [ ] Add sanitize iteration tests for multi-step recovery.
- [ ] Add CLI tests for lint/fix/parse/rules in JSON mode.
- [ ] Add HTTP server tests for format-aware sanitize/parse/examples endpoints.
- [ ] Add browser UI tests for the JSON playground and example switching.
- [ ] Update README and release docs to describe YAML + JSON support accurately.
- [ ] Run the release validation suite after implementation.
  - `go test ./...`
  - `go test -race ./...`
  - `golangci-lint run -v`
  - `docmgr doctor --ticket SANITIZE-007 --stale-after 30`

## Phase 9: Ticket maintenance

- [ ] Keep the design doc synchronized with implementation decisions.
- [x] Add an implementation diary and keep it current as work lands.
- [x] Run `docmgr doctor --ticket SANITIZE-007 --stale-after 30`.

## Phase 10: Experiment tooling

- [x] Add a ticket-local tree-sitter JSON parse inspection script under `scripts/`.
- [x] Add a ticket-local malformed-case matrix script under `scripts/`.
- [x] Add a ticket-local heuristic-probe script under `scripts/`.
- [ ] Add a ticket-local script that classifies malformed cases into parser-driven, heuristic-driven, or hybrid-detection buckets.
- [ ] Add a ticket-local script that emits a machine-readable rule matrix for the first-release malformed cases.
