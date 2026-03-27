---
Title: 'Implementation Plan: Split UI and Improve Presets'
Ticket: SANITIZE-005
Status: active
Topics:
    - frontend
    - ui
    - refactoring
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - internal/server/static/index.html:Current monolithic UI (555 lines, inline CSS+JS)
    - internal/server/server.go:Go server with embed.FS static serving and API endpoints
    - pkg/yaml/examples.go:Hardcoded example YAML snippets served via /api/examples
    - examples/yaml/:25 YAML test corpus files (not yet served to frontend)
ExternalSources: []
Summary: "Split the monolithic index.html into separate CSS, JS, and HTML files. Load example YAML files from the examples/yaml/ directory via the backend and present them as selectable presets with descriptions."
LastUpdated: 2026-03-13T09:52:28.624330563-04:00
WhatFor: ""
WhenToUse: ""
---

# Implementation Plan: Split UI and Improve Presets

## Executive Summary

The current YAML Sanitizer web UI is a single 555-line `index.html` with inline CSS (203 lines) and JavaScript (280 lines). This ticket splits it into separate files for maintainability, and enhances the example/preset system to load YAML files from the `examples/yaml/` corpus (25 files) instead of only the hardcoded snippets in `pkg/yaml/examples.go`.

## Problem Statement

1. **Monolithic HTML** — All CSS, JS, and HTML live in one file, making it hard to edit, lint, or cache independently.
2. **Limited presets** — The frontend only shows the ~10 hardcoded examples from `pkg/yaml/examples.go`. The project has 25 richer YAML test files in `examples/yaml/` that aren't exposed to the UI.
3. **No categorization** — Examples aren't grouped (valid vs. broken vs. edge-case), making the dropdown less useful as the corpus grows.

## Proposed Solution

### 1. File Split

Create this structure under `internal/server/static/`:

```
internal/server/static/
├── index.html          # Slim HTML shell (~70 lines)
├── css/
│   └── style.css       # All CSS extracted from <style> block
└── js/
    └── app.js          # All JS extracted from <script> block
```

The Go `//go:embed static` directive and `http.FileServer` require **zero changes** — new files are automatically embedded and served.

### 2. Backend: Serve File-Based Examples

Add a new endpoint or extend `/api/examples` to also serve the `examples/yaml/*.yaml` files:

- **Option A (chosen): Extend `/api/examples`** — Return both hardcoded and file-based examples in a single list, with a `source` field (`"builtin"` vs `"file"`) and a `category` field derived from the filename prefix convention.
- Embed `examples/yaml/` via a second `//go:embed` directive in the server (or in `pkg/yaml/`) and parse filenames to extract names/categories.
- Filename convention: `NN-description.yaml` → name = "Description", category = "valid" if file starts with `00-04`, "error" otherwise.

### 3. Frontend: Grouped Preset Selector

- Replace the flat `<select>` dropdown with an `<optgroup>`-based selector that groups examples by category (Valid, Parse Errors, Lint Issues, Mixed).
- Show example description as tooltip or subtitle.
- Auto-analyze on selection (existing behavior, preserved).

## Design Decisions

### Separate files vs. bundler
**Decision:** Simple file split, no bundler (Vite/webpack).
**Rationale:** The JS is ~280 lines — a bundler adds complexity with no benefit. The `embed.FS` approach keeps everything in a single binary. If the frontend grows significantly, a bundler can be added later (see SANITIZE-002 design doc for the `go-web-frontend-embed` pattern).

### Extend `/api/examples` vs. new endpoint
**Decision:** Extend the existing endpoint.
**Rationale:** The frontend already calls `/api/examples`. Adding a `source` field is backwards-compatible. One fetch is simpler than two.

### Embed examples in server vs. read from disk
**Decision:** Embed via `//go:embed`.
**Rationale:** Consistent with existing `static` embedding. Single binary deployment. No runtime file-path issues.

### File-based example metadata
**Decision:** Derive name from filename, read description from a YAML comment header if present, otherwise use filename.
**Rationale:** The `examples/yaml/` files already follow a `NN-slug.yaml` naming convention. Parsing a `# description:` header comment is simple and avoids a separate metadata file.

## Implementation Plan

### Phase 1: File Split (no behavior change)

1. Extract CSS from `index.html` → `css/style.css`
2. Extract JS from `index.html` → `js/app.js`
3. Update `index.html` to `<link>` and `<script src>` the new files
4. Verify: `go build && sanitize serve` — UI works identically

### Phase 2: Backend Example Loading

5. Add `//go:embed` for `examples/yaml/` in appropriate package
6. Write a loader that reads embedded YAML files, extracts name/description/category
7. Extend `/api/examples` handler to merge built-in + file-based examples
8. Add `source`, `category`, and `filename` fields to the Example JSON response

### Phase 3: Frontend Preset UI

9. Update JS to handle the new `category` field — group examples with `<optgroup>`
10. Add category styling/badges in the dropdown or a richer selector
11. Show file-based example descriptions (tooltip or inline)
12. Test with all 25 YAML files

### Phase 4: Polish

13. Add a "category" filter or search to the preset selector if the list gets long
14. Visual polish — consistent spacing, responsive behavior
15. Final build + manual testing

## Alternatives Considered

### Full Vite/React SPA
Rejected — overkill for the current UI size. The `go-web-frontend-embed` skill exists if we need this later.

### Separate metadata JSON for examples
Rejected — the filename convention and optional comment headers are sufficient. A separate JSON file would need to stay in sync with the YAML files.

### Load examples from disk at runtime
Rejected — breaks single-binary deployment model. `embed.FS` is the established pattern in this project.

## Open Questions

1. Should file-based examples fully replace the hardcoded `Examples` slice in `pkg/yaml/examples.go`, or coexist? (Recommendation: coexist for now, the hardcoded ones are concise single-issue demos while the file-based ones are richer multi-issue scenarios.)
2. Category derivation — should we use filename prefix ranges, or add explicit `# category:` comments in each YAML file?

## References

- Current UI: `internal/server/static/index.html`
- Server: `internal/server/server.go` (embed.FS, FileServer, API handlers)
- Hardcoded examples: `pkg/yaml/examples.go`
- File corpus: `examples/yaml/` (25 files)
- SANITIZE-002 design doc: discusses `go-web-frontend-embed` pattern for future
