---
Title: From regexes to a real sanitizer
Ticket: SANITIZE-008
Status: active
Topics:
    - writing
    - documentation
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Draft structure for an article about sanitize evolving from a regex-heavy cleaner into a parse-aware structured-text repair tool.
LastUpdated: 2026-03-27T09:34:45.541568027-04:00
WhatFor: Provide a section-level writing plan for telling the sanitize origin story as an engineering evolution narrative.
WhenToUse: Use when writing an origin-story article about why regex-only approaches were not enough and what architecture replaced them.
---

# From regexes to a real sanitizer

## Goal

Prepare an article draft that tells the project’s most intuitive story: starting from line-based cleanup and ending with a parser-assisted, evidence-driven sanitizer.

## Context

This is a broad audience article with enough technical detail to satisfy engineers. It should feel like an evolution narrative rather than a subsystem deep dive.

## Quick Reference

**Working title:** From regexes to a real sanitizer

**Core thesis:** `sanitize` became interesting when it stopped treating tree-sitter as a separate debugging tool and started using shared analysis to drive linting and repair.

**Suggested section map:**

1. `The original problem`
   - Intent: explain the annoying class of "almost valid" YAML and JSON inputs.
   - Content: malformed config, generated text, LLM wrappers, why strict parsers fail hard.
   - Original resources:
     - `README.md`
     - `examples/yaml/`
     - `examples/json/`
   - Potential pseudocode and examples:
     - show a tiny malformed YAML example and a wrapped JSON example.

2. `The regex phase`
   - Intent: make the early approach sound useful but limited.
   - Content: line regexes for tabs, missing spaces, list dashes, commas; why these rules are attractive at first.
   - Original resources:
     - `pkg/yaml/lint.go`
     - `pkg/yaml/fix.go`
   - Potential pseudocode and examples:
     - `for each line { run regex rules }`

3. `Where regexes broke down`
   - Intent: show the failure modes that demand structure.
   - Content: duplicate keys, mixed indentation, parse errors anchored on the wrong line, JSON wrappers versus structural JSON breakage.
   - Original resources:
     - `ttmp/2026/03/13/SANITIZE-004--improve-yaml-linting-and-sanitizing-with-tree-sitter-aware-analysis/sources/01-example-corpus-parse-vs-lint-matrix.md`
     - `ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/05-json-repair-matrix.md`
   - Potential pseudocode and examples:
     - side-by-side “parser sees more / parser sees less” examples.

4. `The architectural turn`
   - Intent: explain shared analysis as the pivotal move.
   - Content: `documentAnalysis`, one parse per iteration, parser spans reused by lint and fix flows.
   - Original resources:
     - `pkg/yaml/analysis.go`
     - `pkg/yaml/sanitize.go`
     - `pkg/json/analysis.go`
   - Potential pseudocode and examples:
     - show the `documentAnalysis` struct.

5. `What the system looks like now`
   - Intent: land the current shape of the project.
   - Content: CLI, packages, browser playground, example corpus, ticket docs.
   - Original resources:
     - `README.md`
     - `internal/server/server.go`
     - `internal/server/static/js/app.js`
   - Potential pseudocode and examples:
     - a mermaid pipeline diagram.

6. `The engineering lesson`
   - Intent: make the article useful beyond this repo.
   - Content: regexes are not bad; they just need a structural home and a stopping rule.
   - Original resources:
     - design docs in `SANITIZE-004` and `SANITIZE-007`
   - Potential pseudocode and examples:
     - a short “bad / better / best” summary table.

## Usage Examples

- Good opening anecdote: “I didn’t set out to build a parser-assisted sanitizer; I just kept tripping over text that was almost valid.”
- Good visual: before/after diagram from line regexes to parse-aware loop.
- Good closing move: connect the story to any tool that sits between permissive humans and strict downstream parsers.

## Related

- `reference/06-anatomy-of-the-yaml-pipeline.md`
- `reference/05-tree-sitter-as-infrastructure-not-a-magic-wand.md`
- `/home/manuel/code/wesen/obsidian-vault/Projects/2026/03/27/PROJ - Sanitize - YAML Sanitizing Deep Dive.md`
