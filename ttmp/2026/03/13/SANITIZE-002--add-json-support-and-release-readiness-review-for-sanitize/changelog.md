# Changelog

## 2026-03-13

- Initial workspace created
- Added release-readiness review, JSON-support design, intern guide, and diary deliverables

## 2026-03-13

Completed JSON support design analysis, public release review, and intern onboarding guide based on full repository inspection and validation runs

### Related Files

- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/cmd/sanitize-server/main.go — Server/API reviewed
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/cmd/sanitize/main.go — CLI contract reviewed
- /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/yaml/sanitize.go — Core runtime analyzed for extraction

## 2026-03-13

- `docmgr doctor --ticket SANITIZE-002 --stale-after 30` passed
- Uploaded `SANITIZE-002 JSON support and release review.pdf` to `/ai/2026/03/13/SANITIZE-002`

## 2026-03-13

- Expanded `SANITIZE-002` tasks into a concrete implementation backlog covering release blockers, shared-core extraction, and JSON support rollout

## 2026-03-13

- Added a Glazed CLI migration track to the backlog, covering command descriptions, Cobra wiring, help/logging integration, and compatibility with the current `sanitize` interface

## 2026-03-13

- Updated the JSON implementation plan and backlog to assume a clean breaking migration: no compatibility wrappers, no legacy CLI flag aliases, and a Glazed-first command redesign

## 2026-03-13

- Uploaded refreshed bundle `SANITIZE-002 JSON support and release review v2.pdf` to `/ai/2026/03/13/SANITIZE-002`

## 2026-03-13

- Split JSON implementation work into new ticket `SANITIZE-003`
- Moved `02-json-support-architecture-and-implementation-plan.md` to `SANITIZE-003`
- Pruned `SANITIZE-002` tasks so this ticket stays focused on release-readiness fixes and analysis

## 2026-03-13

- Replaced indentation-only duplicate-key detection with tree-scoped mapping analysis
- Added regression coverage for keys repeated under different parent mappings and different sequence items

## 2026-03-13

- Refactored `cmd/sanitize` into a testable `run(...)` path so CLI exit codes can be asserted directly
- Fixed non-zero exit behavior for dirty JSON-output runs in both lint and sanitize modes
- Added command tests covering text output, JSON output, and unresolved parse errors

## 2026-03-13

- Hardened `sanitize-server` with explicit `http.Server` timeouts and bounded JSON request decoding
- Added HTTP handler tests for sanitize and parse endpoints
- Cleared the remaining `gosec` findings, including the parse offset conversion warning

## 2026-03-13

- Replaced the hand-rolled CLI with a Glazed/Cobra command tree using `fix`, `lint`, and `serve`
- Folded the standalone `sanitize-server` binary into `sanitize serve`
- Moved the HTTP server into `internal/server`, updated README usage, and removed `sanitize-server` Goreleaser artifacts
