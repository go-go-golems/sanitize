---
Title: Implement JSON support for sanitize
Ticket: SANITIZE-003
Status: active
Topics:
    - json
    - api-design
    - documentation
DocType: index
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2026-03-13T07:57:56.057889893-04:00
WhatFor: ""
WhenToUse: ""
---

# Implement JSON support for sanitize

## Overview

This ticket isolates the JSON implementation work from the broader `SANITIZE-002` analysis ticket. It owns the architecture, CLI redesign, shared-core extraction, server/UI contract updates, and test work required to add first-class JSON support to `sanitize`.

## Key Links

- `design-doc/02-json-support-architecture-and-implementation-plan.md`

## Status

Current status: **active**

## Topics

- json
- api-design
- documentation

## Tasks

See [tasks.md](./tasks.md) for the current task list.

## Changelog

See [changelog.md](./changelog.md) for recent changes and decisions.

## Current Direction

This ticket assumes a clean breaking migration. There are no users to preserve, so the implementation should prefer a Glazed-first CLI, a shared format-agnostic core, and a strict format-aware server contract instead of compatibility wrappers or legacy aliases.

## Structure

- design/ - Architecture and design documents
- reference/ - Prompt packs, API contracts, context summaries
- playbooks/ - Command sequences and test procedures
- scripts/ - Temporary code and tooling
- various/ - Working notes and research
- archive/ - Deprecated or reference-only artifacts
