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
Summary: Step-by-step diary for SANITIZE-007 planning and documentation work.
LastUpdated: 2026-03-13T20:30:00-04:00
WhatFor: Record the implementation-planning steps, commit hashes, and review notes for JSON support research and design.
WhenToUse: Use when reviewing what changed in SANITIZE-007, why the task list and guide were updated, and how to validate the ticket docs.
---

# Diary

## Goal

Capture the implementation-planning steps for `SANITIZE-007`, including how the JSON support backlog was expanded and how the intern guide was revised to include the HTTP server and browser playground work.

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
