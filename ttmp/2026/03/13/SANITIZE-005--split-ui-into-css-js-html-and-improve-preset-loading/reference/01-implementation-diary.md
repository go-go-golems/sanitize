---
Title: Implementation Diary
Ticket: SANITIZE-005
Status: active
Topics:
    - frontend
    - ui
    - refactoring
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: "Step-by-step narrative of implementing SANITIZE-005: splitting the monolithic UI and adding file-based example presets."
LastUpdated: 2026-03-13T09:55:05.059985425-04:00
WhatFor: ""
WhenToUse: ""
---

# Implementation Diary — SANITIZE-005

## 2026-03-13 — Session 1

### Phase 1: File Split (Tasks 2-5)

**What:** Split monolithic `internal/server/static/index.html` (555 lines) into three files:
- `css/style.css` — lines 8-202 (all CSS from the `<style>` block)
- `js/app.js` — lines 273-552 (all JS from the `<script>` block)
- `index.html` — reduced to ~80 lines, just HTML shell with `<link>` and `<script src>` refs

**How it went:** Straightforward extraction. The Go server uses `//go:embed static` + `http.FileServer`, so adding subdirectories required zero Go code changes — new files are automatically embedded and served at their paths (`/css/style.css`, `/js/app.js`).

**Verified:** `go build ./cmd/sanitize` succeeds. The embed.FS picks up the new directory structure automatically.

**What's tricky:** Nothing — the embed.FS pattern handles this cleanly. The only thing to watch is that `<link>` and `<script>` paths must be absolute (`/css/style.css` not `css/style.css`) since the HTML is served at `/`.
