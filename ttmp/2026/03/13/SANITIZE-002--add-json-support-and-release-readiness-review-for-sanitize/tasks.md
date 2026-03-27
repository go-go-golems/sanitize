# Tasks

## Completed

- [x] Create ticket `SANITIZE-002`
- [x] Inspect the current codebase and release plumbing
- [x] Write JSON support architecture and implementation plan
- [x] Write public release review with concrete findings
- [x] Write intern-oriented architecture and implementation guide
- [x] Relate key files, validate with `docmgr doctor`, and upload bundle to reMarkable

## Completed Implementation Work

- [x] Fix duplicate-key detection so keys are scoped to actual sibling mappings rather than indentation-only heuristics
- [x] Add a regression test proving `a.timeout` and `b.timeout` do not get flagged or renamed as duplicates
- [x] Fix CLI exit-code behavior so JSON-output modes still return non-zero status on dirty input
- [x] Add CLI tests covering text output, JSON output, and invalid serve configuration
- [x] Replace `http.ListenAndServe` with an explicit `http.Server` and configure read/write/idle timeouts
- [x] Bound request sizes in `sanitize-server` handlers and add HTTP handler tests
- [x] Resolve the `pkg/yaml/parse.go` integer-conversion warning reported by `gosec`
- [x] Re-run `make gosec` and capture a clean result
- [x] Hand off JSON implementation work to `SANITIZE-003` and keep this ticket focused on release readiness findings
- [x] Replace the hand-rolled CLI with a Glazed/Cobra command tree
- [x] Fold the standalone `sanitize-server` binary into `sanitize serve`
- [x] Update README and Goreleaser packaging so the public release ships one `sanitize` binary
- [x] Add a manually explorable YAML example corpus for quick lint/fix/server checks
