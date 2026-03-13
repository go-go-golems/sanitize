---
Title: Diary
Ticket: SANITIZE-004
Status: active
Topics:
    - yaml
    - linting
    - treesitter
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: cmd/sanitize/main_test.go
      Note: CLI regression coverage for parse command
    - Path: internal/cli/commands.go
      Note: Added sanitize parse command in commit 63a3d07
    - Path: internal/cli/root.go
      Note: Registered parse command in the Glazed root
    - Path: pkg/yaml/analysis.go
      Note: Added in commit 1ceb01c
    - Path: pkg/yaml/fix.go
      Note: Span-aware parse coverage and shared-analysis usage in commit 273d0ef
    - Path: pkg/yaml/indentation.go
      Note: Added in commit ee7f357
    - Path: pkg/yaml/line_index.go
      Note: Added in commit 1ceb01c
    - Path: pkg/yaml/sanitize.go
      Note: One-analysis-per-iteration orchestration in commit 273d0ef
    - Path: pkg/yaml/types.go
      Note: Expanded lint issue model in commit 273d0ef
    - Path: ttmp/2026/03/13/SANITIZE-004--improve-yaml-linting-and-sanitizing-with-tree-sitter-aware-analysis/scripts/parse_lint_matrix.go
      Note: Experiment script written during this ticket
    - Path: ttmp/2026/03/13/SANITIZE-004--improve-yaml-linting-and-sanitizing-with-tree-sitter-aware-analysis/sources/01-example-corpus-parse-vs-lint-matrix.md
      Note: Generated experiment output used in the diary
ExternalSources: []
Summary: Chronological diary for the tree-sitter-aware linting analysis ticket, including the helper CLI tooling and the evidence-gathering workflow.
LastUpdated: 2026-03-13T08:52:51.121920062-04:00
WhatFor: Record the concrete steps, commands, decisions, and review guidance for SANITIZE-004.
WhenToUse: Use when reviewing or continuing the tree-sitter-aware linting investigation.
---




# Diary

## Goal

Capture the work performed for SANITIZE-004: ticket setup, CLI/debug tooling, corpus experiments, and the final analysis documents.

## Step 1: Add parse inspection tooling and corpus experiments

I started by confirming that the ticket needed more than prose. The fastest way to analyze tree-sitter leverage points was to make the parser behavior inspectable from the existing Glazed CLI and then run that output against the manual YAML corpus. That keeps the investigation grounded in repeatable commands instead of one-off snippets.

The result of this step was a new `sanitize parse` subcommand plus a ticket-local experiment script that summarizes parse errors and lint issues across `examples/yaml/`. Together they turned the parser into something I could probe directly while writing the design.

### Prompt Context

**User prompt (verbatim):** "Create a new ticket for improving the linting/sanitizing process by leveraging treesitter. Store an extended anlysis in the ticket. You can use the scripts/ folder of the new ticket to run experiments.  You can add helper commands to the CLI as well (for example printing the parse tree, etc...) to help you along during your analysis. If you write new code, commit appropriately.

Create a detailed analysis / design / implementation guide that is very detailed for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file
  references.
  It should be very clear and detailed. Store in the ticket and the nupload to remarkable."

**Assistant interpretation:** Create a new docmgr ticket focused on tree-sitter-aware YAML diagnostics, add any tooling needed to inspect parser behavior, document the findings thoroughly, and deliver the bundle to reMarkable.

**Inferred user intent:** Produce an implementation-ready architecture investigation, not a vague brainstorming note, and leave behind enough tooling and documentation for someone else to continue the work.

**Commit (code):** `63a3d07` - `feat(cli): add parse inspection tooling`

### What I did

- Created `SANITIZE-004` and seeded the ticket documents.
- Inspected the current parser, lint, fix, sanitize loop, duplicate-key traversal, CLI, and tests.
- Added `sanitize parse` in `internal/cli/commands.go` and registered it in `internal/cli/root.go`.
- Added CLI tests in `cmd/sanitize/main_test.go`.
- Added the ticket-local experiment script `scripts/parse_lint_matrix.go`.
- Generated `sources/01-example-corpus-parse-vs-lint-matrix.md`.
- Ran:
  - `go test ./cmd/sanitize ./internal/... ./pkg/yaml`
  - `go build ./cmd/sanitize`
  - `go run ./cmd/sanitize parse --json examples/yaml/21-multiple-errors.yaml`
  - `go run ./cmd/sanitize parse examples/yaml/22-duplicate-key-different-parents.yaml`
  - `go run ./cmd/sanitize lint --json examples/yaml/22-duplicate-key-different-parents.yaml`
  - `go run ./ttmp/2026/03/13/SANITIZE-004--improve-yaml-linting-and-sanitizing-with-tree-sitter-aware-analysis/scripts/parse_lint_matrix.go > ttmp/2026/03/13/SANITIZE-004--improve-yaml-linting-and-sanitizing-with-tree-sitter-aware-analysis/sources/01-example-corpus-parse-vs-lint-matrix.md`

### Why

- The parser already existed, but there was no direct CLI path for inspecting tree-sitter output after the Glazed migration.
- A dedicated corpus matrix was the fastest way to answer whether parse locations are reliable enough to drive linting and fixing more aggressively.

### What worked

- The new `sanitize parse` command made it easy to inspect representative fixtures without writing throwaway Go snippets.
- The matrix script quickly separated parse-only, heuristic-only, and hybrid failures.
- The existing example corpus was strong enough to expose the central architectural pattern without adding new fixtures first.

### What didn't work

- `docmgr ticket list --plain` did not give useful output in this workspace, so I used `docmgr status --summary-only` plus direct filesystem inspection instead.
- The raw parse-error start row was not useful as a universal fix target. For example, `go run ./cmd/sanitize parse --json examples/yaml/10-tab-indent.yaml` reported a structural error on `service:` rather than on the tab-indented child lines.

### What I learned

- Tree-sitter is already helping more than the public API suggests because duplicate-key detection is tree-based.
- The current design does not have a tree-sitter problem; it has a data-sharing problem.
- The project needs a shared analysis object more than it needs more regexes.

### What was tricky to build

- The main trap was avoiding an overconfident conclusion from parse-error locations. Some errors, like mixed indentation, span the offending region well enough to be helpful, but others, like tab indentation, collapse to the parent mapping line. The solution was not to abandon parse locations or to overtrust them, but to frame them as structural hints that need to be combined with heuristics.

### What warrants a second pair of eyes

- The proposed future shape of `LintIssue`. Once it grows spans and source metadata, CLI and server JSON will change too.
- Whether a shared analysis object should stay internal to `pkg/yaml` or be exposed as a public API.

### What should be done in the future

- Implement the shared `analyzeDocument` pass described in the design doc.
- Upgrade `LintIssue` to span-based diagnostics.
- Add parser-derived structural lint issues and a dedicated mixed-indentation lint rule.

### Code review instructions

- Start with `internal/cli/commands.go` and review the new `parse` command.
- Then read `cmd/sanitize/main_test.go` for CLI expectations.
- Then inspect `scripts/parse_lint_matrix.go` and `sources/01-example-corpus-parse-vs-lint-matrix.md` to see how the analysis was grounded.
- Validate with:
  - `go test ./cmd/sanitize ./internal/... ./pkg/yaml`
  - `go build ./cmd/sanitize`
  - `go run ./cmd/sanitize parse --json examples/yaml/21-multiple-errors.yaml`
  - `go run ./ttmp/2026/03/13/SANITIZE-004--improve-yaml-linting-and-sanitizing-with-tree-sitter-aware-analysis/scripts/parse_lint_matrix.go`

### Technical details

- The CLI addition lives in:
  - `internal/cli/commands.go`
  - `internal/cli/root.go`
  - `cmd/sanitize/main_test.go`
- The experiment assets live in:
  - `ttmp/2026/03/13/SANITIZE-004--improve-yaml-linting-and-sanitizing-with-tree-sitter-aware-analysis/scripts/parse_lint_matrix.go`
  - `ttmp/2026/03/13/SANITIZE-004--improve-yaml-linting-and-sanitizing-with-tree-sitter-aware-analysis/sources/01-example-corpus-parse-vs-lint-matrix.md`

## Step 2: Write the design and intern guide

With the tooling and corpus results in place, I moved to the actual deliverables. The goal of this step was to explain the current codebase precisely enough that a new engineer could implement the refactor without rediscovering the same evidence.

The main design conclusion is that the sanitizer should become tree-sitter-aware through a shared analysis pass, not through a wholesale replacement of heuristics. The documents in this ticket are written around that conclusion and around the three failure categories exposed by the corpus matrix.

### Prompt Context

**User prompt (verbatim):** (same as Step 1)

**Assistant interpretation:** Turn the investigation into a durable, implementation-ready ticket package and make it usable by an intern.

**Inferred user intent:** Leave behind clear engineering guidance, not just the answer to a single question.

**Commit (code):** N/A

### What I did

- Wrote the primary design doc with evidence-backed architecture notes, a problem taxonomy, a proposed shared-analysis design, pseudocode, and phased implementation work.
- Wrote the intern guide with a read order, system map, file-by-file explanation, diagrams, and command references.
- Updated the ticket index and task list so the workspace is continuation-friendly.
- Added vocabulary entries for `yaml`, `linting`, and `treesitter` so `docmgr doctor` can pass cleanly.

### Why

- The user explicitly asked for a very detailed guide for a new intern.
- The current codebase is small, but the interaction between parser output and heuristics is subtle enough that an underspecified plan would create rework.

### What worked

- The corpus matrix provided concrete examples for every major design claim.
- The small size of `pkg/yaml` made it possible to cover the whole runtime in one design document without hand-waving.

### What didn't work

- N/A

### What I learned

- The most useful explanation for new contributors is not "here are the files," but "here are the three classes of failure and why the design needs all of them."
- The new `sanitize parse` command is already a valuable debugging surface even before any deeper refactor lands.

### What was tricky to build

- The intern guide needed to be detailed without reading like a changelog. The way to make it coherent was to structure it around flow, invariants, and review order instead of around file dumps.

### What warrants a second pair of eyes

- The migration plan for richer `LintIssue` spans.
- Any proposal to expose a public `Analyze` API.

### What should be done in the future

- Keep the matrix script and corpus in sync with new rules.
- Add golden expectations once the richer issue model lands.

### Code review instructions

- Read `design-doc/01-tree-sitter-driven-yaml-linting-and-sanitizing-analysis.md` first.
- Then read `reference/02-intern-guide-to-tree-sitter-aware-linting-in-sanitize.md`.
- Use the matrix in `sources/01-example-corpus-parse-vs-lint-matrix.md` to verify the claims.

### Technical details

- Core evidence files:
  - `pkg/yaml/parse.go`
  - `pkg/yaml/lint.go`
  - `pkg/yaml/duplicate_keys.go`
  - `pkg/yaml/fix.go`
  - `pkg/yaml/sanitize.go`
  - `pkg/yaml/sanitize_test.go`

## Step 3: Validate and upload the bundle

Once the ticket docs were written and related, I ran `docmgr doctor` and moved to reMarkable delivery. The ticket passed doctor after I fixed two hygiene issues: the shared vocabulary needed a `treesitter` topic entry, and the generated matrix script needed to emit valid frontmatter so docmgr would treat the source artifact as a proper Markdown document.

The first upload attempt failed because the raw matrix document was not PDF-friendly enough for pandoc. Rather than hiding that failure, I kept the matrix in the ticket as source evidence and rebuilt the reMarkable bundle around the human-written deliverables only. That produced a clean PDF upload while preserving the underlying experiment artifact locally.

### Prompt Context

**User prompt (verbatim):** (same as Step 1)

**Assistant interpretation:** Validate the ticket, fix any delivery issues, and complete the reMarkable upload.

**Inferred user intent:** Finish the work end to end, including the final delivery path.

**Commit (code):** N/A

### What I did

- Added `treesitter` to `ttmp/vocabulary.yaml`.
- Updated `scripts/parse_lint_matrix.go` to emit docmgr frontmatter.
- Regenerated `sources/01-example-corpus-parse-vs-lint-matrix.md`.
- Ran:
  - `docmgr validate frontmatter --doc /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-004--improve-yaml-linting-and-sanitizing-with-tree-sitter-aware-analysis/sources/01-example-corpus-parse-vs-lint-matrix.md --suggest-fixes`
  - `docmgr doctor --ticket SANITIZE-004 --stale-after 30`
  - `remarquee status`
  - `remarquee cloud account --non-interactive`
  - `remarquee upload bundle --dry-run ...`
  - `remarquee upload bundle ...`
  - `remarquee cloud ls /ai/2026/03/13/SANITIZE-004 --long --non-interactive`

### Why

- The user explicitly asked for the ticket bundle to be uploaded to reMarkable.
- Doctor-clean ticket data is the safest point to perform that upload.

### What worked

- `docmgr validate frontmatter` confirmed the regenerated matrix file was valid.
- `docmgr doctor --ticket SANITIZE-004 --stale-after 30` passed cleanly after the vocabulary and frontmatter fixes.
- The reduced bundle containing the index, design doc, intern guide, and diary uploaded successfully.
- `remarquee cloud ls /ai/2026/03/13/SANITIZE-004 --long --non-interactive` showed `SANITIZE-004 tree-sitter linting analysis`.

### What didn't work

- The first bundle upload that included the generated matrix failed with:
  - `Error: pandoc failed: Error producing PDF.`
  - `! Undefined control sequence.`
  - `l.2776 ...tt{ERROR} rows 1:1 to 3:11 text=“root:\n`
  - `: exit status 43`

### What I learned

- The raw experimental matrix is valuable as ticket evidence, but not every generated artifact is suitable for a polished PDF bundle.
- If a ticket has both source evidence and polished docs, it is worth treating those as separate delivery layers.

### What was tricky to build

- The tricky part was deciding whether to sanitize the generated matrix for PDF output or to keep it as a local evidence artifact. The safer choice was to preserve the original evidence in the ticket and upload only the polished documents.

### What warrants a second pair of eyes

- If future tickets want to upload source artifacts too, the matrix format may need a PDF-safe rendering mode.

### What should be done in the future

- Consider adding an optional `--pdf-safe` mode to the matrix script if source artifacts are expected to ship in reMarkable bundles regularly.

### Code review instructions

- Check `ttmp/vocabulary.yaml` for the added topic.
- Check `scripts/parse_lint_matrix.go` for the frontmatter emission change.
- Verify the final upload path with:
  - `remarquee cloud ls /ai/2026/03/13/SANITIZE-004 --long --non-interactive`

### Technical details

- Delivery target: `/ai/2026/03/13/SANITIZE-004`
- Uploaded document name: `SANITIZE-004 tree-sitter linting analysis`

## Step 4: Implement the shared analysis core

After the research phase, I started with the safest internal refactor: centralize parser-derived facts without changing the external lint and fix semantics yet. That meant adding one internal analysis path, moving duplicate-key collection onto it, and introducing a reusable line index for byte-to-line translation.

This step was intentionally scoped to Phase 1 of the ticket. I did not touch the public lint issue model yet because the first priority was removing duplicated parser work and creating one place where future tree-sitter-aware behavior could live.

### Prompt Context

**User prompt (verbatim):** "ok, add tasks to SANITIZE-004 and implement pthem task by task, commit at appropriate intervals and keep a diary"

**Assistant interpretation:** Turn the design ticket into active implementation work, commit focused slices, and keep the ticket docs synchronized with the code.

**Inferred user intent:** Make steady, reviewable progress rather than dropping one large unstructured refactor.

**Commit (code):** `1ceb01c` - `refactor(yaml): add shared document analysis`

### What I did

- Added `pkg/yaml/analysis.go` with `documentAnalysis` and `analyzeDocument(src string)`.
- Added `pkg/yaml/line_index.go` with a reusable line index and byte-to-row and column helpers.
- Changed `ParseTree` to delegate to `analyzeDocument`.
- Moved duplicate-key collection into the shared analysis path.
- Updated `Lint` to consume duplicate keys from the shared analysis path instead of forcing another parse.
- Added `TestLineIndexAtByte` in `pkg/yaml/sanitize_test.go`.
- Ran:
  - `gofmt -w pkg/yaml/analysis.go pkg/yaml/line_index.go pkg/yaml/parse.go pkg/yaml/duplicate_keys.go pkg/yaml/lint.go pkg/yaml/sanitize_test.go`
  - `go test ./pkg/yaml ./cmd/sanitize ./internal/...`
  - `go run ./cmd/sanitize parse --json examples/yaml/17-duplicate-key-sibling.yaml`

### Why

- The design doc identified shared analysis as the first prerequisite for every later task.
- Duplicate-key detection was already tree-based, so moving it onto the shared parse path was a low-risk way to reduce duplicated parsing immediately.

### What worked

- The refactor stayed behavior-preserving for existing lint and sanitize behavior.
- The new line index gives the package a reusable location primitive instead of a repeated linear scan helper.
- Tests stayed green without having to rewrite the public API first.

### What didn't work

- N/A

### What I learned

- The shared-analysis step was smaller than it looked because the package already had most of the raw tree-sitter walking logic; the missing piece was central orchestration.

### What was tricky to build

- The main sharp edge was keeping the refactor plain-data-only. It would have been easy to start passing raw tree nodes around, but that would have created lifetime and testability problems immediately.

### What warrants a second pair of eyes

- `pkg/yaml/analysis.go` because it becomes the new internal hub for later work.
- `pkg/yaml/line_index.go` because future span-based diagnostics will depend on it.

### What should be done in the future

- Layer the richer lint issue model on top of the new analysis path.
- Use the shared analysis inside the sanitize loop itself so each iteration only parses once.

### Code review instructions

- Start with `pkg/yaml/analysis.go`.
- Then read `pkg/yaml/parse.go`, `pkg/yaml/duplicate_keys.go`, and `pkg/yaml/lint.go`.
- Validate with:
  - `go test ./pkg/yaml ./cmd/sanitize ./internal/...`
  - `go run ./cmd/sanitize parse --json examples/yaml/17-duplicate-key-sibling.yaml`

### Technical details

- New files:
  - `pkg/yaml/analysis.go`
  - `pkg/yaml/line_index.go`
- Updated files:
  - `pkg/yaml/parse.go`
  - `pkg/yaml/duplicate_keys.go`
  - `pkg/yaml/lint.go`
  - `pkg/yaml/sanitize_test.go`

## Step 5: Add structural lint diagnostics and reuse analysis during sanitize

Once the shared analysis path existed, I used it to make the public diagnostics more tree-sitter-aware. This step expanded `LintIssue`, added parser-derived issues for parse-only failures, and changed the sanitize loop so each iteration reuses one structural analysis instead of calling parse and lint independently.

I also tightened fix targeting by expanding parse-error coverage from a single start row to every row touched by the parse-error span. That is the first actual fix behavior change driven by tree-sitter location data in this ticket.

### Prompt Context

**User prompt (verbatim):** (same as Step 4)

**Assistant interpretation:** Keep implementing SANITIZE-004 in focused slices with real commits and ticket bookkeeping.

**Inferred user intent:** Land the first meaningful user-visible improvements from the design, not just scaffolding.

**Commit (code):** `273d0ef` - `feat(yaml): add structural lint diagnostics`

### What I did

- Expanded `LintIssue` in `pkg/yaml/types.go` with source and span fields while keeping `Row` as a compatibility alias.
- Added parser-derived lint issues in `pkg/yaml/lint.go` with rules `structural_parse_error` and `missing_syntax_node`.
- Added heuristic issue span population for line-based rules and duplicate-key rules.
- Added `lintIssuesFromAnalysis(...)` so `Lint` and `Sanitize` can share the same diagnostic assembly logic.
- Changed `applyFixes` to consume `documentAnalysis` for parse-error spans and duplicate-key occurrences.
- Changed sanitize orchestration so each iteration uses one `analyzeDocument` result rather than independent parse and lint calls.
- Removed now-dead duplicate-key wrapper functions once the analysis path fully owned that work.
- Added tests:
  - `TestLint_ParseErrorProducesStructuralIssue`
  - `TestLint_HeuristicIssueCarriesSpan`
- Ran:
  - `gofmt -w pkg/yaml/types.go pkg/yaml/lint.go pkg/yaml/fix.go pkg/yaml/sanitize.go pkg/yaml/duplicate_keys.go pkg/yaml/sanitize_test.go`
  - `go test ./pkg/yaml ./cmd/sanitize ./internal/...`
  - `go run ./cmd/sanitize lint --json examples/yaml/20-mixed-indent.yaml`
  - `go run ./cmd/sanitize fix --json examples/yaml/20-mixed-indent.yaml`

### Why

- Parse-only failures were still invisible to `Lint`, which contradicted the design goals of the ticket.
- The sanitize loop was still reparsing more than necessary even after the first refactor.

### What worked

- `sanitize lint --json examples/yaml/20-mixed-indent.yaml` now returns a structural parse issue instead of `null`.
- `sanitize fix --json examples/yaml/20-mixed-indent.yaml` now reports the parser-derived original lint issue and still fixes the indentation successfully.
- The span coverage change gives `applyFixes` better structural context without requiring a full rewrite of the fix pipeline.

### What didn't work

- I tried to split the richer lint diagnostics and the sanitize-loop reuse into two separate commits, but the intermediate state left dead wrapper functions behind and failed `golangci-lint` with:
  - `pkg/yaml/duplicate_keys.go:18:6: func findDuplicateKeys is unused (unused)`
  - `pkg/yaml/fix.go:184:6: func fixDuplicateKeys is unused (unused)`
- I collapsed those two slices into one coherent commit instead of forcing a broken midpoint.

### What I learned

- The first user-visible benefit of tree-sitter awareness is not a fancy selector DSL. It is simply making parse-only failures visible in the same lint stream as heuristic issues.
- Reusing the same analysis object inside `Sanitize` makes later fix-pipeline cleanup much easier because the plumbing is already in place.

### What was tricky to build

- The hardest part was choosing a safe migration for `LintIssue`. I kept `Row` as a compatibility field so the current fix pipeline and tests would not have to be rewritten all at once, while still making the richer span data available immediately.

### What warrants a second pair of eyes

- `pkg/yaml/lint.go` because it now defines the emerging contract for parse-derived issues and heuristic issue spans.
- `pkg/yaml/sanitize.go` because the orchestration path changed even though the outward API stayed the same.
- `pkg/yaml/fix.go` because span-aware row coverage subtly changes when tree-driven fixes can trigger.

### What should be done in the future

- Finish the remaining Phase 4 task by removing the separate `lintIssues` parameter from `applyFixes` entirely if that still feels worthwhile after the next refactor.
- Add a dedicated `mixed_indent` lint rule instead of relying only on the generic structural parse issue.
- Decide whether `LintIssue.Row` can now be removed or should remain as a compatibility alias.

### Code review instructions

- Start with `pkg/yaml/types.go`.
- Then read `pkg/yaml/lint.go`.
- Then read `pkg/yaml/fix.go` and `pkg/yaml/sanitize.go`.
- Validate with:
  - `go test ./pkg/yaml ./cmd/sanitize ./internal/...`
  - `go run ./cmd/sanitize lint --json examples/yaml/20-mixed-indent.yaml`
  - `go run ./cmd/sanitize fix --json examples/yaml/20-mixed-indent.yaml`

### Technical details

- Updated files:
  - `pkg/yaml/types.go`
  - `pkg/yaml/lint.go`
  - `pkg/yaml/fix.go`
  - `pkg/yaml/sanitize.go`
  - `pkg/yaml/duplicate_keys.go`
  - `pkg/yaml/sanitize_test.go`

## Step 6: Add a dedicated mixed-indentation lint rule

The next missing piece from the design was a first-class `mixed_indent` lint issue. Up to this point, mixed indentation only surfaced indirectly through a generic structural parse error and then got fixed later by `fixMixedIndentation`. That meant users could not see the specific cause in `Lint`, even though the fixer already knew how to repair it.

I fixed that by extracting shared indentation analysis into a helper, reusing it in both linting and fixing, and finishing the `applyFixes` API cleanup so it now derives its own lint rule view from the shared document analysis.

### Prompt Context

**User prompt (verbatim):** (same as Step 4)

**Assistant interpretation:** Continue the ticket task by task and land the next coherent slice with a focused commit.

**Inferred user intent:** Keep closing explicit ticket tasks, not just making opportunistic changes.

**Commit (code):** `ee7f357` - `feat(yaml): add mixed indentation lint rule`

### What I did

- Added `pkg/yaml/indentation.go` with shared helpers for:
  - leading-space counting,
  - dominant indentation unit detection,
  - mixed-indentation offender row detection.
- Added `lintMixedIndentation(...)` to `pkg/yaml/lint.go`.
- Updated `lintIssuesFromAnalysis(...)` to include the new `mixed_indent` rule.
- Changed `applyFixes(...)` to consume only `src`, `documentAnalysis`, and `cfg`, deriving lint issues internally instead of taking a separate `lintIssues` argument.
- Updated `fixMixedIndentation(...)` to reuse the new indentation helper instead of maintaining its own duplicate analysis.
- Added `TestLint_MixedIndentProducesDedicatedIssue`.
- Ran:
  - `gofmt -w pkg/yaml/indentation.go pkg/yaml/lint.go pkg/yaml/fix.go pkg/yaml/sanitize.go pkg/yaml/sanitize_test.go`
  - `go test ./pkg/yaml ./cmd/sanitize ./internal/...`
  - `go run ./cmd/sanitize lint --json examples/yaml/20-mixed-indent.yaml`
  - `go run ./cmd/sanitize fix --json examples/yaml/20-mixed-indent.yaml`

### Why

- The ticket explicitly called for a first-class mixed-indentation lint issue.
- `applyFixes(...)` was still carrying an avoidable extra parameter even after the shared analysis refactor.
- The existing fixer logic already had the core indentation reasoning, so it made sense to make lint reuse it too.

### What worked

- `sanitize lint --json examples/yaml/20-mixed-indent.yaml` now emits both:
  - `structural_parse_error`
  - `mixed_indent`
- `sanitize fix --json examples/yaml/20-mixed-indent.yaml` still repairs the file to clean YAML and now reports the dedicated original lint issue alongside the structural parse issue.
- The fix pipeline is simpler because `applyFixes(...)` no longer needs a separately threaded lint slice.

### What didn't work

- N/A

### What I learned

- Mixed indentation was a good example of the intended architecture: structural parse failure plus a more specific heuristic explanation layered on top.
- Shared helper extraction made the lint and fix behavior line up more naturally than keeping two similar indentation analyses in separate places.

### What was tricky to build

- The main design choice was how broadly to emit `mixed_indent`. I kept it gated behind existing parse errors so it stays focused on structural failures instead of becoming an over-eager style rule on parse-clean input.

### What warrants a second pair of eyes

- `pkg/yaml/indentation.go` because it is now shared by linting and fixing.
- `pkg/yaml/fix.go` because the `applyFixes(...)` signature changed and the rule-derivation path moved inward.

### What should be done in the future

- Revisit whether `mixed_indent` should remain parse-gated or eventually become a broader standalone lint rule.
- Use the same pattern for `extra_colon_in_value`, which is still more line-based than it should be.

### Code review instructions

- Start with `pkg/yaml/indentation.go`.
- Then read `pkg/yaml/lint.go` and the new `lintMixedIndentation(...)`.
- Then read `pkg/yaml/fix.go` for the `applyFixes(...)` cleanup.
- Validate with:
  - `go test ./pkg/yaml ./cmd/sanitize ./internal/...`
  - `go run ./cmd/sanitize lint --json examples/yaml/20-mixed-indent.yaml`
  - `go run ./cmd/sanitize fix --json examples/yaml/20-mixed-indent.yaml`

### Technical details

- New file:
  - `pkg/yaml/indentation.go`
- Updated files:
  - `pkg/yaml/lint.go`
  - `pkg/yaml/fix.go`
  - `pkg/yaml/sanitize.go`
  - `pkg/yaml/sanitize_test.go`
