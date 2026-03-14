---
Title: Diary
Ticket: SANITIZE-007
Status: active
Topics:
    - json
    - linting
    - api-design
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/tasks.md
      Note: Detailed phased backlog for JSON support, API work, and the browser playground
    - Path: ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/design-doc/01-json-support-outline-and-malformed-llm-json-error-taxonomy.md
      Note: High-level design doc updated to include the format-aware server and UI plan
    - Path: ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/reference/02-intern-guide-to-json-support-and-tree-sitter-aware-malformed-json-recovery.md
      Note: Deep implementation guide expanded with HTTP and browser UI details
    - Path: internal/server/server.go
      Note: Existing server contract that the JSON work must generalize from YAML-only request bodies
    - Path: internal/server/static/index.html
      Note: Existing browser shell that must become a format-aware playground
    - Path: internal/server/static/js/app.js
      Note: Existing browser-side analysis loop and rendering flow referenced in the guide
ExternalSources: []
Summary: Step-by-step diary for SANITIZE-007 planning, package implementation, and corpus work.
LastUpdated: 2026-03-13T20:30:00-04:00
WhatFor: Record the implementation-planning steps, code commits, failures, and review notes for JSON support research and implementation.
WhenToUse: Use when reviewing what changed in SANITIZE-007, why the task list and guide were updated, and how to validate the implemented JSON slices.
---

# Diary

## Goal

Capture the implementation-planning and implementation steps for `SANITIZE-007`, including how the JSON support backlog was expanded, how the intern guide was revised, and how the first `pkg/json` and corpus slices landed.

## Step 1: Expand the ticket into a real implementation backlog

The first problem was that `SANITIZE-007` had useful research but not yet a task list detailed enough to drive implementation. The original tasks were high-level and did not reflect the actual product surfaces already present in the repository, especially the bundled HTTP server and browser UI.

I expanded the task list into explicit phases that cover the JSON package, rule matrix, CLI surface, format-aware HTTP API, browser UI, corpus work, validation, and ticket-local experiment tooling. I also created a proper diary document and linked it from the ticket index so the ticket now has a durable place to record implementation slices.

### Prompt Context

**User prompt (verbatim):** "Alright, add detailed task list, and also handle the web UI . (update your implementation guide to have a json parse playground in the UI as well). 

Work task by task, commit at appropriate intervals, keep a diary"

**Assistant interpretation:** Expand `SANITIZE-007` into a concrete phased implementation ticket, include explicit web UI work, and record the work incrementally with commits and diary entries.

**Inferred user intent:** Turn the JSON research ticket into a real execution plan that a contributor can work from without guessing about CLI, API, or UI scope.

**Commit (code):** `527d058` — `docs(ticket): expand json implementation backlog`

### What I did

- Read the current ticket files and the current server/UI code paths.
- Confirmed the existing browser app is real and YAML-specific, not hypothetical.
- Rewrote [tasks.md](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/tasks.md) into detailed implementation phases.
- Added explicit HTTP and browser UI tasks, including a JSON parse playground mode.
- Linked the diary from [index.md](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/index.md).

### Why

- The repo already has `sanitize serve`, so JSON support must cover the browser-facing surface, not just the package and CLI.
- A vague task list is hard to execute and hard to review.
- The ticket needed a structure that makes it obvious what "done" means for JSON support.

### What worked

- The existing repo structure made it straightforward to map the work into phases.
- The current UI/server code revealed concrete integration tasks quickly.
- A phased backlog made the JSON playground work much easier to describe precisely.

### What didn't work

- N/A in this slice. No tooling or docs-generation failures occurred.

### What I learned

- `SANITIZE-007` was already strong on malformed-case evidence but weak on delivery-surface planning.
- The current UI contract is tightly coupled to `{ yaml: ... }`, so API cleanup is a core implementation task, not a polish item.

### What was tricky to build

- The subtle part was making the task list concrete without prematurely forcing a particular code abstraction. The ticket now names the product surfaces and expected behaviors, but it intentionally does not force `pkg/core` extraction up front.

### What warrants a second pair of eyes

- The decision to use a single format-aware browser app rather than separate YAML and JSON UIs is correct in my view, but it is still a product decision worth confirming.
- The proposed request contract `{ format, input }` is a breaking API change and should stay explicit.

### What should be done in the future

- Keep marking phases complete as implementation lands.
- Add a machine-readable rule matrix once JSON rule definitions are settled.

### Code review instructions

- Start with [tasks.md](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/tasks.md).
- Confirm that the phases cover package, CLI, server, UI, tests, and corpus work.
- Check [index.md](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/index.md) to verify the diary is discoverable.

### Technical details

- Relevant existing runtime paths:
  - [server.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/internal/server/server.go)
  - [app.js](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/internal/server/static/js/app.js)
  - [index.html](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/internal/server/static/index.html)

## Step 2: Add the JSON playground plan to the intern guide and design doc

The next step was to turn the guide from "JSON package and CLI notes" into a full system guide. The earlier version described `pkg/json`, tree-sitter, and CLI work well enough, but it did not explain the existing HTTP server and SPA in enough detail for a new contributor to build a JSON mode confidently.

I expanded the guide to cover the current server and browser architecture, proposed a format-aware request contract, and described a concrete JSON parse playground. I also updated the high-level design doc so the UI is part of the official implementation plan instead of being implied only by the tasks.

### Prompt Context

**User prompt (verbatim):** (same as Step 1)

**Assistant interpretation:** Update the implementation guide so it explicitly includes a JSON playground in the browser UI, and keep the ticket docs synchronized with that plan.

**Inferred user intent:** Make the ticket usable by a new intern who needs to understand both the codebase and the product surface they are expected to build.

**Commit (code):** `f78529a` — `docs(ticket): add json playground implementation plan`

### What I did

- Added server and SPA file references to the guide frontmatter.
- Added a new section in the guide describing the current HTTP routes and browser flow.
- Added a proposed format-aware request body and dispatch pseudocode.
- Added a dedicated "JSON parse playground" section with UI goals, layout, state, and render behavior.
- Extended the implementation phases with explicit server and UI phases.
- Updated the high-level design doc to include the format-aware API contract and embedded browser playground work.

### Why

- An intern cannot implement the UI cleanly if the guide only discusses package internals.
- The current UI is hard-coded to YAML labels and a YAML-shaped request contract.
- The design doc needed to reflect the full implementation scope so ticket planning and architecture stay aligned.

### What worked

- The current SPA is small and easy to explain, which makes it realistic to evolve rather than rewrite.
- The existing render flow already maps well to JSON mode, especially the parse tree and applied-fixes panes.
- Describing the UI as a "format-aware playground" clarified the product direction.

### What didn't work

- N/A in this slice. I did not hit formatting or validation failures while updating the docs.

### What I learned

- The browser-side analysis flow is already disciplined enough that JSON mode can be added incrementally.
- The strongest architectural seam is the HTTP request body and example payload format, not the DOM rendering.

### What was tricky to build

- The tricky part was staying specific enough for implementation without pretending the JSON response schema is fully settled. I handled that by documenting the likely normalized fields and explicitly marking optional JSON-only fields like `strict_parse_clean`.

### What warrants a second pair of eyes

- Whether the browser should expose strict parser failures as a separate badge or fold them into the general issue list.
- Whether `/api/examples` should return one combined list with a `format` field or separate per-format collections.

### What should be done in the future

- Implement the format-aware API contract before touching the browser render logic heavily.
- Add browser tests once the JSON playground work begins.
- Decide the exact result schema the browser will consume across YAML and JSON modes.

### Code review instructions

- Start with the new server/UI section in [02-intern-guide-to-json-support-and-tree-sitter-aware-malformed-json-recovery.md](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/reference/02-intern-guide-to-json-support-and-tree-sitter-aware-malformed-json-recovery.md).
- Then read the updated implementation plan section in [01-json-support-outline-and-malformed-llm-json-error-taxonomy.md](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-error-taxonomy.md).
- Compare the proposed request contract to the current `yamlRequest` in [server.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/internal/server/server.go).

### Technical details

- Current browser flow:

```text
editor/example picker
   -> triggerAnalysis()
   -> POST /api/sanitize
   -> render tree/issues/output/fixes
```

- Proposed request contract:

```json
{
  "format": "json",
  "input": "{\"ok\": true}"
}
```

- Primary files discussed:
  - [server.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/internal/server/server.go)
  - [index.html](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/internal/server/static/index.html)
  - [app.js](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/internal/server/static/js/app.js)
  - [style.css](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/internal/server/static/css/style.css)

## Step 3: Build the `pkg/json` parse and analysis foundation

The first code slice turned the JSON research into a real package. I added `pkg/json` with tree-sitter parsing, strict `encoding/json` validation, duplicate-key collection, a line index, JSON-native result and issue types, and a rule catalog plus option validation. This established the package seam that later CLI and server work can target.

I kept this slice intentionally narrow: no CLI wiring, no HTTP changes, and no browser changes. That made the review about whether the package shape and parser behavior were correct, rather than mixing transport and UI concerns into the first code landing.

### Prompt Context

**User prompt (verbatim):** "go work on it task by task, commit at appropriate intervals, keep a detailed diary."

**Assistant interpretation:** Start implementing the JSON-support ticket in focused slices, committing each one and recording the exact work in the ticket diary.

**Inferred user intent:** Move from planning into real JSON-support code while preserving reviewability and a strong implementation record.

**Commit (code):** `38669df` — `feat(json): add parse analysis foundation`

### What I did

- Added the new package files under [pkg/json](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/json).
- Added tree-sitter JSON parsing and `ParseTree(...)` in [parse.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/json/parse.go).
- Added shared document analysis in [analysis.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/json/analysis.go).
- Added duplicate-key extraction in [duplicate_keys.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/json/duplicate_keys.go).
- Added JSON rule metadata and option validation in [rules.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/json/rules.go) and [options.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/json/options.go).
- Added tests in [parse_test.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/json/parse_test.go) and [rules_test.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/json/rules_test.go).
- Added the `tree-sitter-json` dependency in `go.mod` and `go.sum`.
- Fixed the pre-existing `nonamedreturns` lint error in [examples.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/examples/examples.go) because it was blocking pre-commit and sits directly in the examples path this ticket uses.

### Why

- The rest of the JSON work needs a real package surface, not just scripts in the ticket.
- Strict parsing and tree-sitter parsing need to coexist from the start so later lint and fix work can distinguish structural localization from validity.
- Duplicate-key detection is one of the easier high-value JSON checks and fits naturally into the analysis object.

### What worked

- Tree-sitter JSON uses a simple `object` / `pair` structure, which made duplicate-key traversal straightforward.
- The existing YAML package gave a good template for analysis, rules, and option validation without forcing code sharing too early.
- The targeted tests were enough to shake out the package shape before touching any higher-level surface.

### What didn't work

- The first commit attempt failed in pre-commit because `golangci-lint` reported:

```text
examples/examples.go:58:1: named return "name" with type "string" found (nonamedreturns)
pkg/json/parse.go:49:2: QF1002: could use tagged switch on err (staticcheck)
pkg/json/line_index.go:19:21: func lineIndex.rowAtByte is unused (unused)
pkg/json/line_index.go:38:21: func lineIndex.rowColAtByte is unused (unused)
pkg/json/options.go:94:18: func (*config).ruleEnabled is unused (unused)
```

- I fixed the JSON issues by removing unused helpers until they were needed, tightening the switch in `strictParseBytes`, and fixing the `examples/examples.go` named return so pre-commit could pass cleanly.

### What I learned

- `encoding/json` accepts duplicate keys, so duplicate-key detection really does need to live outside strict parsing.
- The simplest stable shape is still "YAML and JSON packages with similar public contracts", not a shared package abstraction.

### What was tricky to build

- The tricky part was deciding what to include before linting existed. I chose to include the rule catalog and option validation early because that avoids stringly-typed drift later, but I deferred sanitize/fix behavior until there is enough rule logic to justify it.

### What warrants a second pair of eyes

- The current JSON types include `StrictParseClean` in the result shape before sanitize/fix is implemented. That is defensible, but it should stay consistent with whatever server/browser response contract is chosen later.
- The small lint cleanup in [examples.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/examples/examples.go) was necessary and low-risk, but it is outside `pkg/json`, so it deserves a quick glance.

### What should be done in the future

- Add parser-derived lint issues next.
- Add corpus fixtures that match the rule catalog.
- Only then start CLI/server integration.

### Code review instructions

- Start with [analysis.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/json/analysis.go) and [parse.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/json/parse.go).
- Then review [duplicate_keys.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/json/duplicate_keys.go), [rules.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/json/rules.go), and [options.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/json/options.go).
- Validate with:
  - `go test ./pkg/json`
  - `go test ./cmd/sanitize ./internal/... ./pkg/yaml`

### Technical details

- Strict parser path: `encoding/json.Decoder` with `UseNumber()`
- Structural parser path: `tree-sitter-json`
- Core analysis fields:
  - tree text
  - tree-sitter error spans
  - strict parse error
  - duplicate keys
  - line index

## Step 4: Add a file-backed JSON malformed corpus

The next slice added the file-backed JSON example corpus and a JSON loader in the shared `examples` package. This moves JSON examples from scattered notes into the same kind of stable regression asset the YAML side already has.

I also added package-level built-in JSON examples in [pkg/json/examples.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/json/examples.go) so future CLI and server work can offer concise demos before the browser starts consuming the full file corpus.

### Prompt Context

**User prompt (verbatim):** (same as Step 3)

**Assistant interpretation:** Continue implementing the next focused JSON-support task and keep the work split into small, reviewable commits.

**Inferred user intent:** Build implementation assets that later CLI, API, and UI work can rely on instead of leaving the ticket at the package-only stage.

**Commit (code):** `0e8e0d3` — `feat(examples): add json malformed corpus`

### What I did

- Added the JSON corpus under [examples/json](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/examples/json).
- Added [examples/json/README.md](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/examples/json/README.md).
- Added valid controls and malformed single-pattern cases such as:
  - [10-trailing-comma-object.json](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/examples/json/10-trailing-comma-object.json)
  - [16-markdown-fence-wrapper.json](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/examples/json/16-markdown-fence-wrapper.json)
  - [17-python-literals.json](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/examples/json/17-python-literals.json)
  - [21-multiple-top-level-objects.json](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/examples/json/21-multiple-top-level-objects.json)
- Extended [examples.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/examples/examples.go) with `LoadJSONExamples()`.
- Added tests in [examples_test.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/examples/examples_test.go).

### Why

- The browser and server work need a stable corpus that can be filtered by format.
- Rule development is much easier when each malformed pattern already has a canonical file fixture.
- Built-in package examples and file-backed corpus examples serve different needs; both are useful.

### What worked

- The YAML examples loader pattern generalized cleanly to JSON.
- The file naming convention already fit the malformed-case taxonomy from the ticket.

### What didn't work

- The first validation run failed because the YAML loader still called `parseFilename` with the old one-argument signature after I generalized it:

```text
examples/examples.go:44:35: not enough arguments in call to parseFilename
	have (string)
	want (string, string)
```

- I fixed that call site and re-ran the same narrow validation loop.

### What I learned

- The corpus work benefits from staying format-neutral in the shared `examples` package, but the example payload types should remain format-specific for now.
- The JSON example set is already enough to support parse/lint work before any fixers exist.

### What was tricky to build

- The main sharp edge was balancing built-in examples versus file-backed examples. I kept both, but made the file-backed corpus the more comprehensive source of malformed cases.

### What warrants a second pair of eyes

- Whether some of the current `20-29` cases should be renamed as mixed versus ambiguous; right now they are useful but the taxonomy may evolve once fix safety is decided.

### What should be done in the future

- Add mixed multi-failure JSON cases.
- Add UI-oriented metadata export once the server examples endpoint becomes format-aware.

### Code review instructions

- Start with [examples.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/examples/examples.go) and [examples_test.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/examples/examples_test.go).
- Then scan the new `examples/json/` files and [pkg/json/examples.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/json/examples.go).
- Validate with:
  - `go test ./examples ./pkg/json ./cmd/sanitize ./internal/... ./pkg/yaml`
  - `golangci-lint run -v ./examples ./pkg/json ./cmd/sanitize ./internal/... ./pkg/yaml`

### Technical details

- Naming convention:
  - `00-09` valid
  - `10-19` single malformed pattern
  - `20-29` mixed or more ambiguous

## Step 5: Add parse-aware JSON lint diagnostics

With the analysis object and corpus in place, the next package slice was linting. I added parser-derived lint issues, strict-parser diagnostics, and duplicate-key lint issues. This is the first point where the JSON package starts to look like something the future CLI can expose directly.

This slice stays conservative. There are still no JSON fixers, and there are no heuristics yet for wrappers, comments, or Python literals. The goal was to expose the structural and strict-parser signals first so later heuristics can build on a stable baseline.

### Prompt Context

**User prompt (verbatim):** (same as Step 3)

**Assistant interpretation:** Continue the next implementation slice by turning the JSON analysis foundation into a linter surface.

**Inferred user intent:** Keep moving toward usable JSON support in small increments instead of waiting for a giant all-at-once feature branch.

**Commit (code):** `5ac5f6c` — `feat(json): add parse-aware lint diagnostics`

### What I did

- Added [lint.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/json/lint.go).
- Added `Lint(...)` and `LintWithOptions(...)`.
- Added parser-derived rules:
  - `structural_parse_error`
  - `missing_syntax_node`
- Added `strict_parse_error` to the rule catalog and emitted it from the strict parser path.
- Added duplicate-key lint issues from the analysis object.
- Restored and used the line-index row/column helpers so strict-parser and duplicate-key issues carry spans.
- Expanded [parse_test.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/json/parse_test.go) with lint coverage.

### Why

- The CLI and server will need lintable output before fixers exist.
- Strict-parser diagnostics catch cases where tree-sitter and `encoding/json` differ in usefulness.
- Duplicate keys are already present in the analysis object, so surfacing them as lint issues is a natural next step.

### What worked

- The analysis object was already rich enough to drive lint assembly without re-parsing.
- The row/column spans for duplicate keys and strict-parser errors were good enough for package-level diagnostics.
- The focused rule set kept the slice coherent and easy to validate.

### What didn't work

- The first lint validation failed on a `staticcheck` warning:

```text
pkg/json/parse.go:63:9: S1039: unnecessary use of fmt.Sprintf (staticcheck)
```

- After removing the unnecessary formatting call, the next run failed because the `fmt` import remained in [parse.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/json/parse.go):

```text
pkg/json/parse.go:6:2: "fmt" imported and not used
```

- I removed the stale import and re-ran the same test and lint commands until the slice was clean.

### What I learned

- Even before heuristics exist, JSON benefits from having both tree-sitter-derived and strict-parser-derived issue sources.
- The future UI should probably show strict-parse diagnostics distinctly from tree-sitter `ERROR` nodes because they are conceptually different signals.

### What was tricky to build

- The hardest design choice in this slice was how much strict-parser location data to expose. I chose a practical middle ground: use offset-based spans when available, but keep the rule coarse rather than pretending we have perfect semantic locations for every strict-parser failure.

### What warrants a second pair of eyes

- Whether `strict_parse_error` should remain a separate rule or later be folded into a broader structural bucket for the UI.
- Whether the current multi-value strict parse error offset is the right span for display.

### What should be done in the future

- Add heuristic JSON lint rules next.
- Start wiring the CLI for `parse`, `lint`, and `rules` once the rule surface is broad enough.

### Code review instructions

- Start with [lint.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/json/lint.go).
- Then compare it to [analysis.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/json/analysis.go), [parse.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/json/parse.go), and [rules.go](/home/manuel/code/wesen/2026-03-05--yaml-sanitizing/pkg/json/rules.go).
- Validate with:
  - `go test ./pkg/json ./examples ./cmd/sanitize ./internal/... ./pkg/yaml`
  - `golangci-lint run -v ./pkg/json ./examples ./cmd/sanitize ./internal/... ./pkg/yaml`

### Technical details

- Issue sources now in use:
  - `parse`
  - `strict-parser`
  - `heuristic` for duplicate keys
