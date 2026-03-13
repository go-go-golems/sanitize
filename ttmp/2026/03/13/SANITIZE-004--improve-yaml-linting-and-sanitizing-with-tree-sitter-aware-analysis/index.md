---
Title: Improve YAML linting and sanitizing with tree-sitter-aware analysis
Ticket: SANITIZE-004
Status: active
Topics:
    - yaml
    - linting
    - treesitter
DocType: index
Intent: long-term
Owners: []
RelatedFiles:
    - Path: internal/cli/commands.go
      Note: CLI inspection tooling including sanitize parse
    - Path: pkg/yaml/fix.go
      Note: Current fix targeting and tree-error usage
    - Path: pkg/yaml/lint.go
      Note: Current regex-driven lint pipeline
    - Path: pkg/yaml/parse.go
      Note: Tree-sitter parser and structural error extraction
    - Path: pkg/yaml/sanitize.go
      Note: Sanitize orchestration loop
ExternalSources: []
Summary: Overview ticket for the tree-sitter-aware linting investigation and implementation plan.
LastUpdated: 2026-03-13T08:52:51.107024116-04:00
WhatFor: Track the investigation and implementation plan for making YAML linting and sanitizing more tree-sitter-aware.
WhenToUse: Use when changing YAML parse diagnostics, lint rule design, sanitize fix flow, or diagnostic tooling.
---


# Improve YAML linting and sanitizing with tree-sitter-aware analysis

## Overview

This ticket investigates how the YAML sanitizer can use tree-sitter more directly instead of treating parsing, linting, and fixing as mostly separate passes. The current code already parses with tree-sitter, but `Lint` is still dominated by line heuristics, and `Sanitize` reparses and relints repeatedly without a shared analysis object.

The ticket now contains:

- a detailed design doc describing the current architecture, the experimental findings, and a concrete refactor path,
- a detailed intern guide that explains the system from first principles,
- a chronological diary of the work performed for this ticket,
- a ticket-local experiment script plus generated corpus results.

## Key Links

- Design doc: `design-doc/01-tree-sitter-driven-yaml-linting-and-sanitizing-analysis.md`
- Diary: `reference/01-diary.md`
- Intern guide: `reference/02-intern-guide-to-tree-sitter-aware-linting-in-sanitize.md`
- Corpus matrix: `sources/01-example-corpus-parse-vs-lint-matrix.md`
- Experiment script: `scripts/parse_lint_matrix.go`

## Status

Current status: **active**

Research and documentation are complete. The next implementation step is to introduce a shared analysis pass in `pkg/yaml` and migrate lint and fix logic to consume it.

## Topics

- yaml
- linting
- treesitter

## Tasks

See [tasks.md](./tasks.md) for the current task list.

## Changelog

See [changelog.md](./changelog.md) for recent changes and decisions.

## Structure

- design/ - Architecture and design documents
- design-doc/ - Ticket-level design deliverables
- reference/ - Prompt packs, API contracts, context summaries
- playbooks/ - Command sequences and test procedures
- scripts/ - Temporary code and tooling
- sources/ - Generated experiment output and supporting artifacts
- various/ - Working notes and research
- archive/ - Deprecated or reference-only artifacts
