---
Title: Why YAML sanitizing worked better than JSON repair
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
Summary: Draft structure for a comparison article explaining why sanitize achieved more convincing results on YAML than on malformed JSON repair.
LastUpdated: 2026-03-27T09:34:45.547110042-04:00
WhatFor: Provide a comparative article outline contrasting YAML repair success with JSON repair limits.
WhenToUse: Use when writing about format differences, ambiguity, and why conservative repair stops at different boundaries.
---

# Why YAML sanitizing worked better than JSON repair

## Goal

Prepare an article that answers the obvious question a reader will have after seeing both subsystems: why does the YAML side feel mature while the JSON side still feels tentative?

## Context

This is a comparison piece. It should feel part case study, part design lesson, and part honesty about project limits.

## Quick Reference

**Working title:** Why YAML sanitizing worked better than JSON repair

**Core thesis:** YAML gives the project more locally fixable mistakes, while malformed JSON forces the tool to choose between ambiguity and restraint much earlier.

**Suggested section map:**

1. `Start with the surprising result`
   - Intent: state plainly that YAML repair ended up stronger than JSON recovery.
   - Content: describe the current repo split and why this is counterintuitive.
   - Original resources:
     - `README.md`
     - `/home/manuel/code/wesen/obsidian-vault/Projects/2026/03/27/PROJ - Sanitize - YAML Sanitizing Deep Dive.md`
     - `/home/manuel/code/wesen/obsidian-vault/Projects/2026/03/27/PROJ - Sanitize - JSON Recovery Experiments and Limits.md`
   - Potential pseudocode and examples:
     - short comparison table: YAML vs JSON vs safe auto-fix surface.

2. `What YAML mistakes look like`
   - Intent: show that many YAML problems are formatting-local.
   - Content: tabs, missing spaces, list dashes, trailing commas in flow collections, duplicate keys, colon-in-value quoting.
   - Original resources:
     - `pkg/yaml/fix.go`
     - `pkg/yaml/sanitize_test.go`
   - Potential pseudocode and examples:
     - include the output from `sanitize fix examples/yaml/24-deeply-nested-mixed-errors.yaml`.

3. `What malformed JSON mistakes look like`
   - Intent: show the narrower set of safe JSON rewrites.
   - Content: wrappers, comments, Python literals, commas versus missing commas and missing colons.
   - Original resources:
     - `pkg/json/fix.go`
     - `pkg/json/sanitize_test.go`
     - `ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/05-json-repair-matrix.md`
   - Potential pseudocode and examples:
     - include the output from `sanitize fix --format json examples/json/24-llm-wrapper-multi-step.json`.

4. `Ambiguity is the real villain`
   - Intent: make the design logic explicit.
   - Content: explain why `{name:"Alice"}` or `{"name":"Alice" "age":30}` is a riskier rewrite than `key:value` in YAML.
   - Original resources:
     - `ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/design-doc/01-json-support-outline-and-malformed-llm-json-error-taxonomy.md`
   - Potential pseudocode and examples:
     - ambiguous JSON examples listed in a small blockquote or table.

5. `What the project learned`
   - Intent: generalize beyond this repo.
   - Content: strict grammars do not automatically make repair easier; conservative tools should stop where semantics become guesswork.
   - Original resources:
     - `pkg/json/rules.go`
     - `README.md`
   - Potential pseudocode and examples:
     - policy table: auto-fix / lint-only / avoid.

## Usage Examples

- This article can open with two CLI outputs side by side: the YAML deep-error case that recovers and the JSON single-quote case that still exits 1.
- A strong closing sentence is that the repo’s most valuable JSON result was not “more fixes,” but a sharper understanding of where fixes should stop.

## Related

- `reference/07-the-json-recovery-story-useful-but-narrow.md`
- `reference/12-what-not-to-auto-fix.md`
- `reference/06-anatomy-of-the-yaml-pipeline.md`
