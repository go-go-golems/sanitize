---
Title: What I learned trying to repair LLM JSON
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
Summary: Draft structure for a reflective engineering article about malformed-LLM JSON and the limits of automatic repair.
LastUpdated: 2026-03-27T09:35:05.500981054-04:00
WhatFor: Provide a more personal but still evidence-backed structure for writing about the practical lessons from the JSON work.
WhenToUse: Use when writing a reflective article or retrospective in first person about LLM output repair.
---

# What I learned trying to repair LLM JSON

## Goal

Prepare a reflective article about what the JSON work taught the project, especially where “LLM output cleanup” sounds easier than it really is.

## Context

This is a good candidate for a more personal voice later. It can carry more opinion than the formal JSON case-study article, as long as it stays anchored to repo evidence.

## Quick Reference

**Working title:** What I learned trying to repair LLM JSON

**Core thesis:** LLM JSON is messy in repetitive ways, but the difference between repetitive mess and safely repairable structure is much bigger than it first appears.

**Suggested section map:**

1. `Why I thought this would be straightforward`
   - Intent: make the opening relatable.
   - Content: JSON is small, the error patterns look repetitive, tree-sitter exists, so why not just normalize the output?
   - Original resources:
     - `README.md`
   - Potential pseudocode and examples:
     - one “obviously fixable” wrapped JSON example.

2. `The easy wins really are easy`
   - Intent: show what worked.
   - Content: code fences, prose wrappers, Python literals, comments, commas.
   - Original resources:
     - `pkg/json/fix.go`
     - `pkg/json/sanitize_test.go`
   - Potential pseudocode and examples:
     - show the multi-step recovery test case.

3. `The middle ground is deceptive`
   - Intent: explain why some cases are detectable but still not safe to rewrite.
   - Content: single quotes, unquoted keys, missing closing delimiters.
   - Original resources:
     - `pkg/json/heuristics.go`
     - `ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/04-json-detection-buckets.md`
   - Potential pseudocode and examples:
     - detectable does not equal fixable table.

4. `The hard failures are semantic`
   - Intent: explain the real wall.
   - Content: missing commas, missing colons, duplicate keys, multiple top-level values.
   - Original resources:
     - `ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/02-json-heuristic-probe.md`
     - `ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/sources/05-json-repair-matrix.md`
   - Potential pseudocode and examples:
     - show the missing-comma CLI lint output.

5. `What ended up mattering more than more fixers`
   - Intent: show the positive outcome.
   - Content: better result types, clearer API, browser mode, evidence loop, rule confidence.
   - Original resources:
     - `pkg/json/types.go`
     - `internal/server/server.go`
     - `ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/reference/03-diary.md`
   - Potential pseudocode and examples:
     - `StrictParseClean` and `OriginalStrictParseClean` fields.

6. `My durable takeaway`
   - Intent: land the article in a strong reflective close.
   - Content: a trustworthy repair tool is partly defined by what it refuses to do.
   - Original resources:
     - `README.md`
   - Potential pseudocode and examples:
     - short “what I would try next / what I still would not automate” table.

## Usage Examples

- This can be written more personally than the other drafts.
- Good close: the project did not “solve LLM JSON,” but it did identify where disciplined tooling beats hopeful normalization.

## Related

- `reference/07-the-json-recovery-story-useful-but-narrow.md`
- `reference/12-what-not-to-auto-fix.md`
- `/home/manuel/code/wesen/obsidian-vault/Projects/2026/03/27/PROJ - Sanitize - JSON Recovery Experiments and Limits.md`
