---
Title: Add JSON support and release-readiness review for sanitize
Ticket: SANITIZE-002
Status: active
Topics:
    - json
    - api-design
    - release-readiness
    - documentation
DocType: index
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2026-03-13T00:47:47.36828062-04:00
WhatFor: ""
WhenToUse: ""
---

# Add JSON support and release-readiness review for sanitize

## Overview

This ticket captures a follow-on analysis for the `sanitize` repository: how to add JSON support cleanly, whether the current code is ready for public release, and what a new engineer needs to understand before implementing the next phase. The deliverables are evidence-backed design and review documents rather than code changes.

## Key Links

- `design-doc/01-public-release-review.md`
- `reference/01-diary.md`
- `reference/02-intern-guide-to-sanitize.md`

## Status

Current status: **active**

## Topics

- json
- api-design
- release-readiness
- documentation

## Tasks

See [tasks.md](./tasks.md) for the current task list.

## Changelog

See [changelog.md](./changelog.md) for recent changes and decisions.

## Current Conclusion

The repository has solid packaging and validation basics, but it should not be released publicly yet. The most important blockers are the duplicate-key false positive and fixer, inconsistent CLI exit-code behavior in JSON mode, and missing HTTP server hardening. JSON support should be implemented by extracting a shared sanitize core and adding a format-specific JSON engine rather than duplicating `pkg/yaml`.

## Scope Note

JSON implementation planning and execution now live under `SANITIZE-003`. This ticket remains the release-readiness and analysis umbrella for the current codebase.

## Structure

- design/ - Architecture and design documents
- reference/ - Prompt packs, API contracts, context summaries
- playbooks/ - Command sequences and test procedures
- scripts/ - Temporary code and tooling
- various/ - Working notes and research
- archive/ - Deprecated or reference-only artifacts
