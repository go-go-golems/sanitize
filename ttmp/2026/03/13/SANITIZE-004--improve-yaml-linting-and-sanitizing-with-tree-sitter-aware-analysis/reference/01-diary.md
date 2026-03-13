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
