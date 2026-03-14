---
Title: Add JSON support with focus on malformed LLM JSON recovery
Ticket: SANITIZE-007
Status: active
Topics:
    - json
    - linting
    - api-design
DocType: index
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources:
    - local:03-json-parse-errors-import.md
Summary: Ticket for designing JSON support around malformed LLM-generated JSON, including parse-error experiments, heuristic probes, and implementation guidance.
LastUpdated: 2026-03-13T19:32:38.618853469-04:00
WhatFor: Track the research, design, and implementation planning for JSON linting and sanitizing support.
WhenToUse: Use when reviewing the JSON-support roadmap, experiment outputs, and intern implementation guidance.
---


# Add JSON support with focus on malformed LLM JSON recovery

## Overview

This ticket investigates how `sanitize` should support malformed JSON commonly produced by LLMs. It combines a malformed-case taxonomy, tree-sitter experiments, heuristic probes, and an intern-oriented implementation guide that maps the current YAML architecture onto a proposed `pkg/json` engine.

## Key Links

- **Related Files**: See frontmatter RelatedFiles field
- **External Sources**: See frontmatter ExternalSources field
- Design: [design-doc/01-json-support-outline-and-malformed-llm-json-error-taxonomy.md](./design-doc/01-json-support-outline-and-malformed-llm-json-error-taxonomy.md)
- Common malformed cases: [reference/01-common-json-parse-errors-from-llm-output.md](./reference/01-common-json-parse-errors-from-llm-output.md)
- Intern guide: [reference/02-intern-guide-to-json-support-and-tree-sitter-aware-malformed-json-recovery.md](./reference/02-intern-guide-to-json-support-and-tree-sitter-aware-malformed-json-recovery.md)
- Diary: [reference/03-diary.md](./reference/03-diary.md)
- Tree-sitter matrix: [sources/01-json-parse-error-replication-matrix.md](./sources/01-json-parse-error-replication-matrix.md)
- Heuristic probe: [sources/02-json-heuristic-probe.md](./sources/02-json-heuristic-probe.md)

## Status

Current status: **active**

## Topics

- json
- linting
- api-design

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
