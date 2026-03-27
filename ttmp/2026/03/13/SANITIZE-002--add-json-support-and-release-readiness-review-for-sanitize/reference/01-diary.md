---
Title: Diary
Ticket: SANITIZE-002
Status: active
Topics:
    - json
    - api-design
    - release-readiness
    - documentation
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: internal/server/server.go
      Note: Hardened HTTP server implementation and `serve` subcommand target
    - Path: cmd/sanitize/main.go
      Note: CLI entrypoint delegating into the Glazed command tree
    - Path: internal/cli/root.go
      Note: Glazed/Cobra root command wiring and exit-code handling
    - Path: internal/cli/commands.go
      Note: Glazed-backed `fix`, `lint`, and `serve` command implementations
    - Path: pkg/yaml/fix.go
      Note: Investigation evidence for duplicate-key behavior
    - Path: pkg/yaml/lint.go
      Note: Investigation evidence for duplicate-key behavior
ExternalSources: []
Summary: ""
LastUpdated: 2026-03-13T00:47:58.039033936-04:00
WhatFor: ""
WhenToUse: ""
---


# Diary

## Goal

Capture the investigation, validation, and documentation work for `SANITIZE-002`, including what evidence was gathered, what was learned about JSON support, and what release-readiness problems were confirmed.

## Context

This diary started as research and analysis work, then expanded into the implementation steps needed to clear the release-readiness findings. It now records both the original evidence-gathering work and the production code changes made to bring the project closer to a public release.

## Step 1: Create The Ticket And Gather Evidence

This step established a fresh ticket workspace for the JSON-support analysis and release review, then walked the full repository so the recommendations would be anchored to the real code rather than prior assumptions. Because the repository is small, it was practical to inspect every production source file and all release-related configuration.

The investigation showed a useful pattern: the repository is already organized well enough to support a clean second-format design, but a few user-visible correctness and entrypoint issues still need attention before a public release. The most important result of this step was confirming that the duplicate-key heuristic is not just theoretically weak; it rewrites valid YAML today.

### Prompt Context

**User prompt (verbatim):**
> Create a new docmgr ticket to add json support to the sanitize package. 
>
> Analyze the current codebase and design and brainstorm how to support json. 
>
> Also, while you're at it, do a thorough code review and see if this is fit to be released publicly. 
>
> Keep both analysis / reports in the ticket as documents and upload them as a bundle to the remarkable. 
>
> Create a detailed analysis / design / implementation guide that is very detailed for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file
>   references.
>   It should be very clear and detailed. Store in the ticket and the nupload to remarkable.

**Assistant interpretation:** Create a new `docmgr` ticket, inspect the current `sanitize` repository in depth, write a JSON-support design, write a release-readiness/code-review report, write an intern-oriented architecture guide, store all of that in the ticket, and upload the bundle to reMarkable.

**Inferred user intent:** Produce a reusable research-and-planning package that can guide the next implementation phase and help judge whether the project is ready for public distribution.

**Commit (code):** N/A — no code changes in this ticket

### What I did
- Created ticket `SANITIZE-002` with the title `Add JSON support and release-readiness review for sanitize`.
- Added four documents:
  - `design-doc/01-public-release-review.md`
  - `design-doc/02-json-support-architecture-and-implementation-plan.md`
  - `reference/01-diary.md`
  - `reference/02-intern-guide-to-sanitize.md`
- Read the relevant skill instructions for `ticket-research-docmgr-remarkable`, `docmgr`, `remarkable-upload`, and `diary`.
- Inspected all key repository files, including:
  - `README.md`
  - `pkg/yaml/*.go`
  - `cmd/sanitize/main.go`
  - `cmd/sanitize-server/main.go`
  - `cmd/sanitize-server/static/index.html`
  - `.github/workflows/*`
  - `.golangci.yml`
  - `.goreleaser.yaml`
- Ran validation commands:
  - `go test ./...`
  - `go test -race ./...`
  - `make lint`
  - `make build`
  - `make gosec`
- Ran documentation validation:
  - `docmgr doctor --ticket SANITIZE-002 --stale-after 30`
- Ran reMarkable delivery commands:
  - `remarquee status`
  - `remarquee cloud account --non-interactive`
  - `remarquee upload bundle --dry-run ...`
  - `remarquee upload bundle ...`
  - `remarquee cloud ls /ai/2026/03/13/SANITIZE-002 --long --non-interactive`
- Ran behavior reproductions:
  - `printf 'a:\n  timeout: 30\nb:\n  timeout: 60\n' | go run ./cmd/sanitize --json`
  - `printf 'name:Alice\n' | go run ./cmd/sanitize --lint --json`

### Why
- The ticket needed to exist before writing structured deliverables.
- The code review needed real evidence, not assumptions.
- JSON-support design decisions depend on where current responsibilities actually live.

### What worked
- `docmgr` workspace creation worked cleanly.
- The repository was small enough to inspect comprehensively.
- `go test ./...`, `go test -race ./...`, `make lint`, and `make build` all passed.
- The behavior reproductions gave concrete proof for two review findings.
- `docmgr doctor --ticket SANITIZE-002 --stale-after 30` passed cleanly.
- The reMarkable bundle dry-run succeeded, the real upload succeeded, and the remote listing confirmed the uploaded file.

### What didn't work
- `make gosec` did not pass. It reported:
  - `pkg/yaml/parse.go:56` — `G115: integer overflow conversion uint -> int`
  - `cmd/sanitize-server/main.go:99` — `G114: Use of net/http serve function that has no support for setting timeouts`
  - `cmd/sanitize-server/main.go:98` — `G706: Log injection via taint analysis`

### What I learned
- The current design is already close to a format-engine architecture; the sanitize loop is concentrated in one place.
- The riskiest bug is semantic, not mechanical: duplicate-key detection is global-by-indent and can rename valid keys in separate maps.
- The CLI’s machine-readable mode has an automation contract bug because JSON output suppresses the non-zero exit path.

### What was tricky to build
- The trickiest part was separating "promising v0 package" from "safe public release." The repository has several strong signals, such as tests, linting, CI, and Goreleaser configuration, but those positives could easily obscure the correctness issue in duplicate-key handling.
- It also took care to distinguish current behavior from future design. The ticket asked for JSON-support brainstorming, but the design had to stay grounded in the existing file boundaries and UI/API contracts rather than drifting into a greenfield architecture.

### What warrants a second pair of eyes
- The proposed shared-core extraction should be reviewed by someone familiar with expected long-term package names and public import paths.
- The release review should be checked against the intended CLI contract: specifically whether JSON output modes are expected to keep non-zero exit codes on dirty input.
- The JSON design should be reviewed by whoever will own frontend/API compatibility, because the UI currently assumes YAML-specific request shapes.

### What should be done in the future
- Implement the hardening items from the release review before a public release.
- Extract shared sanitize orchestration before implementing `pkg/json`.
- Add CLI and server tests after the format abstraction lands.

### Code review instructions
- Start with `pkg/yaml/sanitize.go`, then read `pkg/yaml/lint.go` and `pkg/yaml/fix.go`.
- Verify the duplicate-key finding by reading:
  - `pkg/yaml/lint.go:24-25`
  - `pkg/yaml/lint.go:72-85`
  - `pkg/yaml/fix.go:182-208`
- Validate behavior with:
  - `printf 'a:\n  timeout: 30\nb:\n  timeout: 60\n' | go run ./cmd/sanitize --json`
  - `printf 'name:Alice\n' | go run ./cmd/sanitize --lint --json`
- Re-run baseline checks:
  - `go test ./...`
  - `go test -race ./...`
  - `make lint`
  - `make build`
  - `make gosec`

### Technical details
- Ticket path: `ttmp/2026/03/13/SANITIZE-002--add-json-support-and-release-readiness-review-for-sanitize`
- Primary outputs:
  - `design-doc/01-public-release-review.md`
  - `design-doc/02-json-support-architecture-and-implementation-plan.md`
  - `reference/02-intern-guide-to-sanitize.md`
- Delivery path:
  - `/ai/2026/03/13/SANITIZE-002/SANITIZE-002 JSON support and release review`
- Key evidence files:
  - `pkg/yaml/*.go`
  - `cmd/sanitize/main.go`
  - `cmd/sanitize-server/main.go`
  - `cmd/sanitize-server/static/index.html`

## Related

- `../design-doc/01-public-release-review.md`
- `../design-doc/02-json-support-architecture-and-implementation-plan.md`
- `../reference/02-intern-guide-to-sanitize.md`

## Step 2: Remove Compatibility Constraints From The Plan

The follow-up clarification changed the implementation strategy materially: the project does not need to preserve CLI or library compatibility because there are no users yet. That made the earlier compatibility-preserving parts of the JSON design actively worse, because they would add wrappers, aliases, and extra tests without protecting anyone.

The resulting plan is cleaner. It now assumes a Glazed-first CLI redesign, a strict format-aware HTTP contract, and permission to remove YAML-specific public shapes that do not fit the shared core. This is the right choice for a pre-release codebase.

### Prompt Context

**User prompt (verbatim):** "you can break compatibility no problem btw, we have no users yet. Update the implementation plan and potentially the tasks."

**Assistant interpretation:** Update the JSON-support implementation plan and backlog to remove compatibility-preserving requirements and reflect a clean breaking migration.

**Inferred user intent:** Optimize the architecture for internal cleanliness and future extensibility instead of preserving transitional APIs that have no consumers.

**Commit (code):** N/A — documentation-only update

### What I did
- Updated `design-doc/02-json-support-architecture-and-implementation-plan.md` to:
  - remove compatibility-wrapper guidance
  - replace legacy CLI aliasing with a breaking Glazed/Cobra verb design
  - change the server contract plan to accept only `format` + `input`
  - reorder phases so release hardening and CLI redesign happen cleanly
- Updated `tasks.md` to:
  - remove compatibility-preservation tasks
  - add breaking-change tasks for the Glazed CLI, API contract, and public shape cleanup
- Updated `changelog.md` to record the shift to a clean breaking migration

### Why
- Compatibility scaffolding only makes sense when there are users or dependent code to protect.
- The project is still pre-release, so this is the cheapest and safest time to simplify the future architecture.

### What worked
- The ticket structure already separated research from implementation, so the clarification could be applied cleanly.
- The Glazed skill guidance fit naturally into the rewritten CLI plan.

### What didn't work
- N/A

### What I learned
- The existence or absence of users should directly affect architecture decisions. In this case it changes the correct answer from "wrap and migrate" to "reshape and simplify."

### What was tricky to build
- The main challenge was updating the plan consistently rather than only deleting a few compatibility bullets. Once the breaking-change assumption changed, the CLI contract, phase ordering, HTTP contract, and backlog wording all needed to move together.

### What warrants a second pair of eyes
- The proposed Glazed verb tree should still be reviewed before implementation starts, because it will become the public CLI shape.
- The shared-core extraction should be checked to ensure package naming still feels coherent after removing wrapper-based compatibility.

### What should be done in the future
- Turn the updated backlog into ordered milestones if implementation will be split across multiple sessions or tickets.

### Code review instructions
- Review the updated plan in `design-doc/02-json-support-architecture-and-implementation-plan.md`.
- Check that `tasks.md` now matches the breaking-change direction and no longer assumes legacy aliases or wrappers.

### Technical details
- Directional changes made:
  - no compatibility wrappers
  - no preserved `--json` alias
  - Glazed-first CLI redesign
  - strict format-aware HTTP request model

## Step 3: Fix Duplicate-Key Scoping

The first implementation step focused on the highest-risk correctness bug in the release review: duplicate-key detection was using indentation-only bookkeeping, which caused valid YAML in separate mappings and separate sequence items to be rewritten incorrectly. The fix moved duplicate detection onto actual tree-sitter mapping nodes so duplicates are now scoped to real sibling mappings instead of global indentation buckets.

This was the right first cut because it addressed the most dangerous behavior: silent corruption of valid user data. It also kept the change locally bounded to `pkg/yaml`, which made it possible to validate with targeted tests and concrete reproductions before moving on to CLI and server work.

### Prompt Context

**User prompt (verbatim):** "work on SANITIZE-002 for making this release worthy. Task by task, commit as you go, keep a detailed diary"

**Assistant interpretation:** Execute the release-readiness tasks in `SANITIZE-002` as focused implementation steps, committing each step separately and recording the work in the ticket diary.

**Inferred user intent:** Move from analysis into disciplined implementation, with traceable progress and commit-by-commit documentation.

**Commit (code):** `0d19eb7` — `fix(yaml): scope duplicate keys to actual mappings`

### What I did
- Added `pkg/yaml/duplicate_keys.go` to detect duplicate keys by traversing actual tree-sitter mapping nodes.
- Replaced the line-based duplicate-key lint logic in `pkg/yaml/lint.go` with results from the new tree-based duplicate detector.
- Replaced the line-based duplicate-key fixer in `pkg/yaml/fix.go` with byte-range replacements against exact duplicate key nodes.
- Added regression tests in `pkg/yaml/sanitize_test.go` covering:
  - duplicate keys under the same mapping
  - same key reused under different parent mappings
  - same key reused in separate sequence items
- Ran:
  - `gofmt -w pkg/yaml/duplicate_keys.go pkg/yaml/lint.go pkg/yaml/fix.go pkg/yaml/sanitize_test.go`
  - `go test ./pkg/yaml`
  - `printf 'a:\n  timeout: 30\nb:\n  timeout: 60\n' | go run ./cmd/sanitize --json`
  - `printf 'items:\n  - name: a\n    id: 1\n  - name: b\n    id: 2\n' | go run ./cmd/sanitize --json`
  - `printf 'config:\n  timeout: 30\n  timeout: 60\n' | go run ./cmd/sanitize --json`

### Why
- The duplicate-key bug was the strongest release blocker because it changed valid YAML semantics.
- Tree-scoped detection is a more defensible release-quality behavior than any indentation-only heuristic.

### What worked
- The tree-sitter Go API exposed enough structure (`Kind`, `NamedChildren`, `ChildByFieldName`) to scope duplicate checks to actual mapping nodes.
- The targeted package tests passed after the refactor.
- The real-world reproductions now behave correctly:
  - `a.timeout` and `b.timeout` no longer collide
  - repeated keys in separate list items no longer collide
  - true sibling duplicates still get renamed

### What didn't work
- N/A

### What I learned
- The tree-sitter YAML grammar is good enough to support release-quality duplicate-key scoping without adding a second parser dependency.
- Sequence-item boundaries matter just as much as parent mapping boundaries for this rule.

### What was tricky to build
- The sharp edge was making the fixer precise enough to rename only the duplicate key token rather than doing another line-level string heuristic. Using exact byte ranges from the key node solved that problem cleanly.

### What warrants a second pair of eyes
- The duplicate-key identity currently normalizes simple quoted keys by stripping the outer quote characters before comparison. That behavior is reasonable, but it should still be reviewed with a few more quoted-key edge cases in mind.

### What should be done in the future
- Add broader fixture coverage for quoted keys and flow mappings if duplicate-key behavior expands further.

### Code review instructions
- Start with `pkg/yaml/duplicate_keys.go`.
- Then review the call sites in `pkg/yaml/lint.go` and `pkg/yaml/fix.go`.
- Validate with:
  - `go test ./pkg/yaml`
  - `printf 'a:\n  timeout: 30\nb:\n  timeout: 60\n' | go run ./cmd/sanitize --json`
  - `printf 'items:\n  - name: a\n    id: 1\n  - name: b\n    id: 2\n' | go run ./cmd/sanitize --json`

### Technical details
- Files changed:
  - `pkg/yaml/duplicate_keys.go`
  - `pkg/yaml/lint.go`
  - `pkg/yaml/fix.go`
  - `pkg/yaml/sanitize_test.go`

## Step 4: Fix CLI JSON Exit Codes And Add Command Tests

The second implementation step fixed the machine-readable CLI contract. Before this change, `sanitize --lint --json` and `sanitize --json` could emit structured output but still exit with status `0` even when issues remained. That behavior would make CI, editors, and shell automation silently accept dirty input, which is not acceptable for a release-worthy command.

I refactored the command into a reusable `run(...)` function so the tests could assert exit codes directly without depending on shell wrappers. That made it possible to cover both text and JSON modes in a fast unit-style test file, while still validating the real `go run` behavior from the shell.

### Prompt Context

**User prompt (verbatim):** "work on SANITIZE-002 for making this release worthy. Task by task, commit as you go, keep a detailed diary"

**Assistant interpretation:** Continue executing the next release-readiness task in `SANITIZE-002`, commit the result as a focused step, and record both the code change and its validation in the diary.

**Inferred user intent:** Make the public CLI predictable for automation, not just for interactive use.

**Commit (code):** `6991e6b` — `fix(cli): honor exit codes in json modes`

### What I did
- Refactored `cmd/sanitize/main.go` so `main()` delegates to `run(args, stdin, stdout, stderr)`.
- Switched flag parsing to a local `flag.FlagSet` so the command can be tested without mutating global process state.
- Made lint-only mode return exit status `1` even when `--json` is used and issues are present.
- Made sanitize mode return exit status `1` when parse or lint issues remain even when `--json` is used.
- Added `cmd/sanitize/main_test.go` with coverage for:
  - lint text mode on dirty input
  - lint JSON mode on dirty input
  - sanitize JSON mode after a successful fix
  - sanitize JSON mode when parse errors remain
- Ran:
  - `gofmt -w cmd/sanitize/main.go cmd/sanitize/main_test.go`
  - `go test ./cmd/sanitize`
  - `printf 'name:Alice\n' | go run ./cmd/sanitize --lint --json; echo EXIT:$?`
  - `printf 'a: [1,2\n' | go run ./cmd/sanitize --json; echo EXIT:$?`

### Why
- Exit codes are part of the CLI contract, especially for machine-readable modes.
- Refactoring to a testable `run(...)` path is a low-cost way to harden the command before the larger Glazed redesign in `SANITIZE-003`.

### What worked
- The `flag.FlagSet` refactor cleanly separated argument parsing from process exit.
- The new command tests run quickly and cover the previously broken cases directly.
- Both shell reproductions now return `EXIT:1` as expected.

### What didn't work
- N/A

### What I learned
- Small CLIs become much easier to harden once `main()` stops owning all of the logic.
- The current release ticket can improve the existing CLI contract without blocking on the future Glazed migration.

### What was tricky to build
- The main edge case was preserving normal help/parse behavior while moving off the package-level `flag` state. Handling `flag.ErrHelp` explicitly keeps the command behavior sane while still making it testable.

### What warrants a second pair of eyes
- The future Glazed migration should revisit which exit codes should be used for parse errors vs flag parse failures, but the current `0` for help, `2` for flag parse errors, and `1` for dirty input is a defensible baseline.

### What should be done in the future
- Keep the CLI tests around as a behavior reference when the Glazed rewrite lands.

### Code review instructions
- Start with `cmd/sanitize/main.go`.
- Then review `cmd/sanitize/main_test.go`.
- Validate with:
  - `go test ./cmd/sanitize`
  - `printf 'name:Alice\n' | go run ./cmd/sanitize --lint --json; echo EXIT:$?`
  - `printf 'a: [1,2\n' | go run ./cmd/sanitize --json; echo EXIT:$?`

### Technical details
- Files changed:
  - `cmd/sanitize/main.go`
  - `cmd/sanitize/main_test.go`

## Step 5: Harden The HTTP Server And Clear The Remaining Gosec Findings

The third implementation step closed the remaining release-hardening findings from the initial review. The HTTP server had been relying on `http.ListenAndServe` without timeouts, request bodies were effectively unbounded, and `gosec` was still flagging an unsafe integer conversion in YAML parse error collection. None of those are deep architectural problems, but together they are the kind of sharp edges that make a small tool look unfinished in a public release.

This step tightened the server in-place before the larger CLI migration. That sequencing kept the validation surface narrow: first make the existing HTTP behavior defensible, then move the entrypoint wiring later without mixing security posture changes and command-tree changes in one commit.

### Prompt Context

**User prompt (verbatim):** "work on SANITIZE-002 for making this release worthy. Task by task, commit as you go, keep a detailed diary"

**Assistant interpretation:** Continue with the next release-readiness findings from the ticket review, commit them as an isolated task, and document both the remediation and the proof that the security scan is now clean.

**Inferred user intent:** Eliminate the remaining concrete release blockers before moving into design-driven refactors like the Glazed CLI migration.

**Commit (code):** `0849cc8` — `refactor(cli): move sanitize to glazed subcommands`

### What I did
- Reworked `cmd/sanitize-server/main.go` so startup builds an explicit `http.Server` instead of calling `http.ListenAndServe` directly.
- Added read header, read, write, and idle timeouts to the server configuration.
- Added `serverPort(...)` validation so `PORT` values are parsed and range-checked explicitly.
- Centralized request decoding through `decodeJSONBody(...)`, using `http.MaxBytesReader` and `json.Decoder.DisallowUnknownFields()`.
- Added `cmd/sanitize-server/main_test.go` with handler coverage for:
  - successful sanitize requests
  - oversized request rejection
  - parse endpoint method handling
  - successful parse requests
- Reworked duplicate-key line reporting to compute line numbers from byte offsets without the `uint -> int` conversion that `gosec` flagged.
- Tightened `pkg/yaml/parse.go` to compare byte offsets without converting `uint` to `int`.
- Ran:
  - `gofmt -w cmd/sanitize-server/main.go cmd/sanitize-server/main_test.go pkg/yaml/duplicate_keys.go pkg/yaml/fix.go pkg/yaml/lint.go pkg/yaml/parse.go`
  - `go test ./...`
  - `make gosec`

### Why
- Public-facing HTTP servers should not ship with zero timeout protection.
- Request-size bounds and strict JSON decoding reduce the chance of accidental abuse or confusing handler behavior.
- A clean `gosec` run was one of the explicit unfinished findings from the release review.

### What worked
- The explicit `http.Server` configuration resolved the timeout-related finding cleanly.
- The request decoder abstraction made it easy to add both size bounds and unknown-field rejection in one place.
- The new handler tests passed without needing end-to-end browser fixtures.
- `make gosec` now reports `Issues : 0`.

### What didn't work
- `gosec` initially still flagged the startup log line when it included the parsed port value, so the startup message had to stay static to keep the scan clean.

### What I learned
- Even in a small repo, lightweight security tooling pushes the code toward better seams: request decoding became clearer once it had to own bounds and validation centrally.
- The duplicate-key refactor still had one release-quality loose end in how it computed line numbers; fixing the scanner findings forced that code to become a little more robust too.

### What was tricky to build
- The subtlest part was preserving the existing API behavior while hardening it. The handlers still need to feel simple from the frontend’s point of view, even though body bounding, stricter decoding, and timeout management now happen under the hood.

### What warrants a second pair of eyes
- The `maxRequestBodyBytes` limit is reasonable for the current UI, but it should still be reviewed if the server is expected to accept very large YAML documents later.

### What should be done in the future
- When the standalone server binary is folded into the main CLI, preserve these server tests and the request-decoding helpers rather than reintroducing looser behavior.

### Code review instructions
- Start with `cmd/sanitize-server/main.go`.
- Then review `cmd/sanitize-server/main_test.go`.
- Finally check the scanner-driven cleanup in `pkg/yaml/duplicate_keys.go`, `pkg/yaml/lint.go`, `pkg/yaml/fix.go`, and `pkg/yaml/parse.go`.
- Validate with:
  - `go test ./...`
  - `make gosec`

### Technical details
- Files changed:
  - `cmd/sanitize-server/main.go`
  - `cmd/sanitize-server/main_test.go`
  - `pkg/yaml/duplicate_keys.go`
  - `pkg/yaml/lint.go`
  - `pkg/yaml/fix.go`
  - `pkg/yaml/parse.go`

## Step 6: Replace The CLI With Glazed And Fold The Server Into `serve`

The fourth implementation step removed the last obviously pre-release part of the user-facing surface: the project had two binaries with divergent command plumbing, one hand-rolled CLI, and one server-only entrypoint. The new shape is a single `sanitize` binary with Glazed/Cobra-backed `fix`, `lint`, and `serve` subcommands, with the HTTP server moved into a shared internal package instead of living behind a second install target.

This is a breaking CLI change, but that is consistent with the earlier decision that the repository has no external users yet. For a first public release, the cleaner outcome is more important than preserving a transient pre-release interface.

### Prompt Context

**User prompt (verbatim):** "Use glazed to implement the CLI (see `glaze help build-first-command` and the skill), and make sanitize-server a `serve` subverb."

**Assistant interpretation:** Rebuild the command-line entrypoint around Glazed/Cobra, preserve the existing YAML fix/lint behavior under explicit subcommands, and make the web server reachable as `sanitize serve` instead of as a separate binary.

**Inferred user intent:** Ship a cleaner public CLI surface with one install target, a more maintainable command tree, and a server implementation that no longer sits outside the main product entrypoint.

**Commit (code):** `pending at time of writing this entry`

### What I did
- Added `internal/cli/root.go` to own the Glazed/Cobra root command, logging initialization, and controlled exit-code mapping.
- Added `internal/cli/commands.go` with Glazed command descriptions and `Run(...)` implementations for:
  - `sanitize fix`
  - `sanitize lint`
  - `sanitize serve`
- Kept file/stdin behavior intact by moving input loading into a shared helper used by the new command implementations.
- Moved the hardened HTTP server from `cmd/sanitize-server` into `internal/server/server.go`.
- Moved the web UI assets into `internal/server/static/` so they remain embedded while no longer requiring a second binary.
- Updated `cmd/sanitize/main.go` to delegate to the new Glazed command tree.
- Updated `cmd/sanitize/main_test.go` to assert the new subcommand-based interface and added coverage for invalid `serve --port`.
- Updated `README.md` so installation and usage examples now describe the single-binary CLI.
- Updated `.goreleaser.yaml` to stop producing `sanitize-server` artifacts and ship only the `sanitize` binary.
- Ran:
  - `gofmt -w cmd/sanitize/main.go cmd/sanitize/main_test.go internal/cli/root.go internal/cli/commands.go internal/server/server.go internal/server/server_test.go`
  - `go mod tidy`
  - `GOROOT=/home/manuel/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.7.linux-amd64 PATH=/home/manuel/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.7.linux-amd64/bin:$PATH GOTOOLCHAIN=local GOSUMDB=off GOPROXY=file:///home/manuel/go/pkg/mod/cache/download GOMODCACHE=/tmp/go-mod GOCACHE=/tmp/go-build GOWORK=off go test ./cmd/sanitize ./internal/... ./pkg/yaml`
  - `GOROOT=/home/manuel/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.7.linux-amd64 PATH=/home/manuel/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.7.linux-amd64/bin:$PATH GOTOOLCHAIN=local GOSUMDB=off GOPROXY=file:///home/manuel/go/pkg/mod/cache/download GOMODCACHE=/tmp/go-mod GOCACHE=/tmp/go-build GOWORK=off go build ./cmd/sanitize`
  - `GOROOT=/home/manuel/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.7.linux-amd64 PATH=/home/manuel/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.7.linux-amd64/bin:$PATH GOTOOLCHAIN=local GOCACHE=/tmp/go-build GOWORK=off /home/manuel/go/bin/gosec -exclude-generated -exclude=G101,G304,G301,G306 ./cmd/sanitize ./internal/... ./pkg/...`
  - `GOROOT=/home/manuel/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.7.linux-amd64 PATH=/home/manuel/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.7.linux-amd64/bin:$PATH GOTOOLCHAIN=local GOSUMDB=off GOPROXY=file:///home/manuel/go/pkg/mod/cache/download GOMODCACHE=/tmp/go-mod GOCACHE=/tmp/go-build GOWORK=off golangci-lint run -v ./cmd/sanitize ./internal/... ./pkg/yaml`

### Why
- A public CLI should not require users to discover and install a second binary just to reach the web UI.
- Glazed gives the command tree a real schema-driven structure instead of continuing to grow an ad hoc `flag`-based entrypoint.
- Moving the server into a shared package removes duplicated entrypoint responsibility and makes future subcommands or API tests easier to add.

### What worked
- The new `sanitize fix` and `sanitize lint` commands preserve the previous YAML behavior while making the verb structure explicit.
- The `sanitize serve` command now reuses the same hardened server implementation and tests from the earlier hardening step.
- The custom Glazed-backed Cobra adapter preserves Glazed flag/argument schemas without inheriting the hard process exits from `cobra.CheckErr`.
- The targeted test, lint, build, and `gosec` runs all passed after the migration.

### What didn't work
- The stock `cli.BuildCobraCommandFromCommand` path in Glazed calls `cobra.CheckErr`, which exits the process directly. That made the new `run(...)` wrapper untestable until the command builder layer was replaced with a custom adapter that still uses Glazed parsing.
- The sandbox environment blocked writes to the shared Go and golangci-lint caches, so validation had to run against a cached Go `1.25.7` toolchain with temp caches and a local file-backed module proxy.

### What I learned
- Glazed’s command-description and parser layers are useful independently of its default Cobra run helper. Using the schemas but taking control of execution gave the right balance here.
- Converging on one binary simplified more than the install story; it also forced the server package boundaries to become more honest.

### What was tricky to build
- The hardest part was not the command descriptions themselves. The real edge was exit handling: the release ticket already depended on testable exit codes, so the Glazed migration had to preserve that contract rather than falling back to library-managed process exits.

### What warrants a second pair of eyes
- The choice of `fix`, `lint`, and `serve` as the public verb set should be reviewed once there is a broader roadmap for future JSON or parse-tree subcommands.
- The dependency expansion from adding Glazed is justified by the new CLI shape, but it increases the transitive module surface and should be kept in mind for long-term maintenance.

### What should be done in the future
- If the project gains more documentation, wire a Glazed help system on top of the current root command so examples and tutorials can live beside the CLI.
- When JSON support lands in `SANITIZE-003`, add sibling JSON verbs or format-aware subcommands to the same command tree instead of introducing a second entrypoint pattern.

### Code review instructions
- Start with `internal/cli/root.go`, then `internal/cli/commands.go`.
- Then review `internal/server/server.go` to see how the `serve` command now reuses the hardened server implementation.
- Finally check `README.md` and `.goreleaser.yaml` to confirm the public surface is now single-binary.
- Validate with:
  - `go test ./cmd/sanitize`
  - `go test ./internal/server ./pkg/yaml`
  - `go build ./cmd/sanitize`
  - `gosec -exclude-generated -exclude=G101,G304,G301,G306 ./cmd/sanitize ./internal/... ./pkg/...`

### Technical details
- Files changed:
  - `cmd/sanitize/main.go`
  - `cmd/sanitize/main_test.go`
  - `internal/cli/root.go`
  - `internal/cli/commands.go`
  - `internal/server/server.go`
  - `internal/server/server_test.go`
  - `internal/server/static/index.html`
  - `README.md`
  - `.goreleaser.yaml`
  - `go.mod`
  - `go.sum`

## Step 7: Add A Manual YAML Example Corpus

The next step added a dedicated example corpus for manual testing and demos. Up to this point the repository had the built-in UI snippets in `pkg/yaml/examples.go` and focused unit tests, but it did not have an easy-to-browse directory of real input files that someone could lint, sanitize, diff, or load into the web UI directly from disk.

The new `examples/yaml/` directory fills that gap. It gives the repo a lightweight manual test fixture set without entangling those files with the Go test suite, and it covers both happy-path controls and the rule-specific or regression-specific cases that matter most for release confidence.

### Prompt Context

**User prompt (verbatim):** "make a directory with a ton of different example cases, in a yaml/ subfolder for now, to test the different cases easily."

**Assistant interpretation:** Create a browsable directory of YAML fixture files, organized for manual testing, with enough coverage to exercise the implemented sanitizer behavior and the known regression boundaries.

**Inferred user intent:** Make it easy to manually sanity-check the CLI and server against representative files without having to reconstruct examples from tests or from embedded Go strings.

**Commit (code):** `pending at time of writing this entry`

### What I did
- Added `examples/yaml/README.md` describing the corpus and how to run common manual checks.
- Added clean/control fixtures:
  - `00-valid-basic.yaml`
  - `01-valid-nested-mapping.yaml`
  - `02-valid-flow-collections.yaml`
  - `03-valid-block-scalar.yaml`
  - `04-valid-comments-and-blank-lines.yaml`
  - `25-clean-duplicate-like-keys.yaml`
  - `26-valid-sequence-of-maps.yaml`
  - `29-valid-url-values.yaml`
- Added focused broken fixtures for implemented rules:
  - `10-tab-indent.yaml`
  - `11-missing-space-after-colon.yaml`
  - `13-list-dash-no-space.yaml`
  - `14-trailing-comma-flow-mapping.yaml`
  - `15-trailing-comma-flow-sequence.yaml`
  - `16-extra-colon-in-value.yaml`
  - `17-duplicate-key-sibling.yaml`
  - `20-mixed-indent.yaml`
  - `27-extra-colon-needs-quotes.yaml`
- Added regression and edge-control fixtures:
  - `12-missing-space-after-colon-url-safe.yaml`
  - `18-quoted-colon-safe.yaml`
  - `19-commented-colon-safe.yaml`
  - `22-duplicate-key-different-parents.yaml`
  - `23-duplicate-key-separate-sequence-items.yaml`
  - `28-unresolved-parse-error.yaml`
- Added mixed-error compound fixtures:
  - `21-multiple-errors.yaml`
  - `24-deeply-nested-mixed-errors.yaml`
- Ran:
  - `go run ./cmd/sanitize lint examples/yaml/11-missing-space-after-colon.yaml`
  - `go run ./cmd/sanitize fix --json examples/yaml/22-duplicate-key-different-parents.yaml`
  - `go run ./cmd/sanitize fix examples/yaml/21-multiple-errors.yaml`

### Why
- Manual fixtures make release checks faster for humans, especially when debugging CLI behavior or demoing the tool.
- The regression cases are easier to verify from real files than from embedded strings in tests.
- The example corpus complements the unit tests without forcing those files into a specific test harness design yet.

### What worked
- The missing-space fixture reports the expected lint failures.
- The different-parent duplicate-key regression case stays clean and unchanged.
- The mixed-error fixture exercises multiple fixers in one run and sanitizes into the expected readable form.

### What didn't work
- N/A

### What I learned
- The repo benefits from having three layers of examples now:
  - small built-in snippets for the web UI
  - targeted unit tests for correctness
  - file-based manual fixtures for operator-style validation

### What was tricky to build
- The main judgment call was deciding which files should be intentionally valid controls versus intentionally broken fixtures. The useful corpus is not just "lots of broken YAML"; it also needs regression boundaries where the sanitizer must *not* change anything.

### What warrants a second pair of eyes
- The corpus naming scheme is intentionally simple and numeric for easy sorting. If the repo later adds JSON fixtures beside it, it may be worth reviewing whether the numbering and grouping should be expanded into subdirectories.

### What should be done in the future
- Add a small smoke script that runs `sanitize lint` or `sanitize fix` across the corpus and summarizes expected clean/dirty outcomes.
- When JSON support lands, mirror this structure under an adjacent `examples/json/` directory.

### Code review instructions
- Start with `examples/yaml/README.md`.
- Spot-check:
  - `examples/yaml/11-missing-space-after-colon.yaml`
  - `examples/yaml/21-multiple-errors.yaml`
  - `examples/yaml/22-duplicate-key-different-parents.yaml`
  - `examples/yaml/23-duplicate-key-separate-sequence-items.yaml`
- Validate with:
  - `go run ./cmd/sanitize lint examples/yaml/11-missing-space-after-colon.yaml`
  - `go run ./cmd/sanitize fix examples/yaml/21-multiple-errors.yaml`
  - `go run ./cmd/sanitize fix --json examples/yaml/22-duplicate-key-different-parents.yaml`

### Technical details
- Files changed:
  - `examples/yaml/README.md`
  - `examples/yaml/*.yaml`
