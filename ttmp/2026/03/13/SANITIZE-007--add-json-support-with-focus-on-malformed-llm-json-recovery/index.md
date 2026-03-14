---
Title: Add JSON support with focus on malformed LLM JSON recovery
Ticket: SANITIZE-007
Status: complete
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
Summary: Ticket for the completed first-release JSON support implementation, including malformed-case experiments, conservative recovery, CLI/API/UI delivery, and intern-oriented implementation guidance.
LastUpdated: 2026-03-13T23:55:00-04:00
WhatFor: Track the research, implementation, evidence, and validation for JSON linting and sanitizing support.
WhenToUse: Use when reviewing what shipped in SANITIZE-007, how the JSON engine works, and where the conservative recovery boundary sits.
---


# Add JSON support with focus on malformed LLM JSON recovery

## Overview

This ticket delivered first-release JSON support for `sanitize` with an emphasis on malformed LLM output. It combines a malformed-case taxonomy, tree-sitter experiments, heuristic probes, conservative fixers, CLI/API/UI integration, and an intern-oriented implementation guide that maps the YAML architecture onto the shipped `pkg/json` engine.

## Key Links

- **Related Files**: See frontmatter RelatedFiles field
- **External Sources**: See frontmatter ExternalSources field
- Design: [design-doc/01-json-support-outline-and-malformed-llm-json-error-taxonomy.md](./design-doc/01-json-support-outline-and-malformed-llm-json-error-taxonomy.md)
- Common malformed cases: [reference/01-common-json-parse-errors-from-llm-output.md](./reference/01-common-json-parse-errors-from-llm-output.md)
- Intern guide: [reference/02-intern-guide-to-json-support-and-tree-sitter-aware-malformed-json-recovery.md](./reference/02-intern-guide-to-json-support-and-tree-sitter-aware-malformed-json-recovery.md)
- Diary: [reference/03-diary.md](./reference/03-diary.md)
- Tree-sitter matrix: [sources/01-json-parse-error-replication-matrix.md](./sources/01-json-parse-error-replication-matrix.md)
- Heuristic probe: [sources/02-json-heuristic-probe.md](./sources/02-json-heuristic-probe.md)
- Detection buckets: [sources/04-json-detection-buckets.md](./sources/04-json-detection-buckets.md)
- Repair matrix: [sources/05-json-repair-matrix.md](./sources/05-json-repair-matrix.md)
- Rule matrix: [sources/06-json-rule-matrix.json](./sources/06-json-rule-matrix.json)
- Overlap study: [sources/07-json-overlap-study.md](./sources/07-json-overlap-study.md)

## Status

Current status: **complete**

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
