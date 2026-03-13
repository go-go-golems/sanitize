---
Title: Split UI into CSS/JS/HTML and improve preset loading
Ticket: SANITIZE-005
Status: active
Topics:
    - frontend
    - ui
    - refactoring
DocType: index
Intent: long-term
Owners: []
RelatedFiles:
    - Path: examples/yaml
      Note: 25 YAML corpus files to expose as presets
    - Path: internal/server/server.go
      Note: Server with embed.FS and API handlers - no changes needed for file split
    - Path: internal/server/static/index.html
      Note: Monolithic UI file to split into CSS/JS/HTML
    - Path: pkg/yaml/examples.go
      Note: Hardcoded examples - will coexist with file-based examples
ExternalSources: []
Summary: ""
LastUpdated: 2026-03-13T09:51:10.157747971-04:00
WhatFor: ""
WhenToUse: ""
---


# Split UI into CSS/JS/HTML and improve preset loading

## Overview

<!-- Provide a brief overview of the ticket, its goals, and current status -->

## Key Links

- **Related Files**: See frontmatter RelatedFiles field
- **External Sources**: See frontmatter ExternalSources field

## Status

Current status: **active**

## Topics

- frontend
- ui
- refactoring

## Tasks

See [tasks.md](./tasks.md) for the current task list.

## Changelog

See [changelog.md](./changelog.md) for recent changes and decisions.

## Structure

- design/ - Architecture and design documents
- reference/ - Prompt packs, API contracts, context summaries
- playbooks/ - Command sequences and test procedures
- scripts/ - Temporary code and tooling
- various/ - Working notes and research
- archive/ - Deprecated or reference-only artifacts
