---
Title: Implementation Plan
Ticket: SANITIZE-001
Status: active
Topics:
    - refactoring
    - packaging
    - testing
    - documentation
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/sanitize/sanitize.go:Core sanitize logic
    - /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/sanitize/examples.go:Built-in examples
    - /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/main.go:HTTP server entry point
    - /home/manuel/code/wesen/2026-03-05--yaml-sanitizing/static/index.html:Web UI
ExternalSources: []
Summary: "Plan to restructure yaml-sanitizer into github.com/go-go-golems/sanitize"
LastUpdated: 2026-03-13T00:30:02.611981668-04:00
WhatFor: ""
WhenToUse: ""
---

# Implementation Plan

## Executive Summary

Restructure the prototype `yaml-sanitizer` into `github.com/go-go-golems/sanitize` — a reusable Go package with a CLI, HTTP server, tests, CI, and docs. The name `sanitize` (not `yaml-sanitize`) anticipates future JSON support.

## Problem Statement

The current repo is a single-commit prototype with no tests, a non-standard module path, regexes compiled per-call, no CLI mode, and no CI. It works but isn't publishable.

## Proposed Solution

### Package layout

```
github.com/go-go-golems/sanitize/
├── pkg/
│   └── yaml/              # package yamlsanitize — the core library
│       ├── sanitize.go
│       ├── lint.go         # extracted from sanitize.go
│       ├── fix.go          # extracted from sanitize.go
│       ├── examples.go
│       └── sanitize_test.go
├── cmd/
│   ├── sanitize/           # CLI: stdin→stdout
│   │   └── main.go
│   └── sanitize-server/    # HTTP server + embedded web UI
│       ├── main.go
│       └── static/
│           └── index.html
├── .golangci.yml
├── .golangci-lint-version
├── .goreleaser.yaml
├── lefthook.yml
├── Makefile
├── .github/workflows/
│   ├── push.yml
│   ├── lint.yml
│   └── release.yaml
├── go.mod
├── go.sum
└── README.md
```

### Key decisions

1. **Package name `yamlsanitize`** — avoids collision with stdlib `yaml` and clarifies purpose when imported as `yamlsanitize`.
2. **Functional options** — `Sanitize(src, ...Option)` to control enabled rules, indent width, max iterations.
3. **Precompile regexes** — move per-call `regexp.MustCompile` into package-level `var`.
4. **go:embed** for static assets in the server binary.
5. **CI from go-template** — golangci-lint, GitHub Actions, GoReleaser, lefthook.

## Implementation Order

1. Rename module, restructure directories (`pkg/yaml/`, `cmd/`)
2. Precompile regexes, clean up error handling
3. Add functional options
4. Embed static assets, wire up server
5. Add CLI (stdin→stdout)
6. Write tests
7. Copy CI/lint/release config from go-template, adapt for `sanitize`
8. Write README
9. Commit docs
