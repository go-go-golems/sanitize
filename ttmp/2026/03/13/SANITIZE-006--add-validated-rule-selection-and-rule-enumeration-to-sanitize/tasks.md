# Tasks

## Phase 1: Ticket setup and design

- [x] Create `SANITIZE-006` and seed the ticket workspace.
- [x] Write the implementation plan for validated rule selection, rule enumeration, and CLI support.
- [x] Relate the core YAML and CLI files to the ticket docs.

## Phase 2: Canonical rule registry

- [x] Add a `pkg/yaml/rules.go` registry that defines every known rule exactly once.
- [x] Include rule metadata needed for CLI/help output, at minimum: name, summary, whether it produces lint issues, and whether it can apply fixes.
- [x] Expose a stable API for listing rule specs and checking whether a rule name is known.
- [x] Add unit tests that fail if rule metadata and rule names drift.

## Phase 3: Rule selection and validation in the YAML package

- [x] Replace the current unvalidated `WithRules(...)` behavior with validated rule selection built on the registry.
- [x] Add support for disabling individual rules separately from selecting an allowlist.
- [x] Add error-returning entrypoints for linting and sanitizing so callers can detect invalid rule names cleanly.
- [x] Keep the existing convenience wrappers for callers that do not need custom rule selection.
- [x] Add tests covering unknown rule names, allowlist behavior, disable-list behavior, and conflicting combinations.

## Phase 4: Glazed CLI integration

- [x] Add `--rule` and `--disable-rule` flags to `sanitize lint`.
- [x] Add `--rule` and `--disable-rule` flags to `sanitize fix`.
- [x] Validate CLI rule names before running lint or fix and return a clear non-zero error for unknown rules.
- [x] Add a `sanitize rules` command that lists the available rules in text and JSON forms.
- [x] Add CLI regression tests for valid rule selection, invalid rule names, and rule listing.

## Phase 5: Docs and validation

- [x] Update the SANITIZE-006 diary after each code slice with commit hashes, commands, and failures.
- [x] Update the ticket changelog after each commit.
- [x] Run `go test ./pkg/yaml ./cmd/sanitize ./internal/...` after each implementation slice.
- [ ] Run `docmgr doctor --ticket SANITIZE-006 --stale-after 30` before closing the turn.
