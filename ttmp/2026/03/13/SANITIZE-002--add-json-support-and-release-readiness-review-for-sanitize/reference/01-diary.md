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
    - Path: cmd/sanitize-server/main.go
      Note: Investigation evidence for gosec findings
    - Path: cmd/sanitize/main.go
      Note: Investigation evidence for CLI exit-code behavior
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

This diary documents research and analysis work only. No production code was modified as part of this ticket. The deliverables are ticket documents, evidence-backed findings, and reMarkable delivery artifacts.

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

**Commit (code):** pending at time of writing this entry

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
