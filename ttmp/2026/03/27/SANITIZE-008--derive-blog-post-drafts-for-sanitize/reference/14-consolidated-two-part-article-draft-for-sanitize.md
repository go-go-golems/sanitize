---
Title: Consolidated two-part article draft for sanitize
Ticket: SANITIZE-008
Status: active
Topics:
    - writing
    - documentation
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: README.md
      Note: Project framing and CLI behavior used in the article opening and scope
    - Path: internal/server/static/js/app.js
      Note: Playground workflow and inspection UI discussed in Part II
    - Path: pkg/yaml/analysis.go
      Note: Shared analysis architecture and tree-sitter-driven document analysis
    - Path: pkg/yaml/fix.go
      Note: YAML fix rules and sanitize loop repair behavior
    - Path: pkg/yaml/lint.go
      Note: Span-aware diagnostics and heuristic lint rules described in Part I
    - Path: pkg/yaml/sanitize.go
      Note: Top-level sanitize orchestration and iteration flow
    - Path: ttmp/2026/03/13/SANITIZE-004--improve-yaml-linting-and-sanitizing-with-tree-sitter-aware-analysis/design-doc/01-tree-sitter-driven-yaml-linting-and-sanitizing-analysis.md
      Note: Primary design record for the tree-sitter and shared-analysis refactor
    - Path: ttmp/2026/03/13/SANITIZE-005--split-ui-into-css-js-html-and-improve-preset-loading/reference/01-implementation-diary.md
      Note: Diary for the playground refactor and iteration workflow
ExternalSources: []
Summary: Friend-reviewable two-part article draft that combines the deep YAML implementation story with the project workflow that made sanitize productive.
LastUpdated: 2026-03-27T10:34:18.416943227-04:00
WhatFor: Provide a substantial but not fully written article draft that can be shared for feedback before the final prose pass.
WhenToUse: Use when preparing one flagship sanitize article that combines the technical internals with the surrounding development workflow.
---


# Consolidated two-part article draft for sanitize

## Goal

Prepare one consolidated article draft that is detailed enough to hand to friends for feedback, but still outline-forward enough that the final prose and voice can be written later by hand.

## Context

This draft should read like a serious long-form article plan rather than a list of ideas. The structure is intentionally two-part:

1. Part I covers the YAML sanitizer implementation in real technical detail.
2. Part II covers the workflow that made the project work in practice: prototype, playground, tickets, corpora, docmgr, and iterative research with LLM help.

The target reader is an engineer who likes implementation detail but also wants to understand how the surrounding workflow shaped the code.

## Quick Reference

**Working title options:**

- `How I Built a Tree-Sitter-Aware YAML Sanitizer and the Workflow Around It`
- `Sanitize: Deep YAML Repair Internals and the Workflow That Made It Work`
- `From YAML Repair Engine to Research Workflow`

**Primary title recommendation:** `How I Built a Tree-Sitter-Aware YAML Sanitizer and the Workflow Around It`

**One-sentence promise:** This article explains both how the YAML side of `sanitize` actually works under the hood and how the project workflow around corpora, playgrounds, tickets, and iterative analysis made the implementation sharper.

**What this article is trying to do:**

- show that the YAML implementation is real and technically coherent,
- show that the surrounding workflow is part of the implementation story, not just project-management overhead,
- avoid over-indexing on the weaker JSON story while still using it as contrast where helpful,
- give enough structure and evidence that early readers can react to the shape of the piece before the final prose pass.

**Suggested article length:** `3,500–5,500` words once written out.

**Recommended tone:** technical, candid, reflective, and slightly narrative. Not academic. Not startup-hype. Not “look at my cool toy.” The tone should be “this is what actually made the project work.”

**Proposed article shape:**

### Opening

#### `Cold open: almost-valid text is a real software problem`
- Intent:
  - hook readers with the pain point before naming the solution.
- Content:
  - “almost valid” YAML and JSON,
  - strict parsers as the final boss,
  - why this kept happening in config, generated output, and LLM text.
- Original resources:
  - `README.md`
  - `examples/yaml/README.md`
  - `examples/json/README.md`
- Potential pseudocode and examples:
  - a two-line malformed YAML example,
  - a prose-wrapped JSON example,
  - one sentence framing: “I didn’t start by wanting a parser project; I started by being annoyed.”

#### `Set the article contract`
- Intent:
  - tell the reader there are two halves to the story.
- Content:
  - first half is the YAML engine itself,
  - second half is the workflow that made the engine better,
  - mention that YAML became the clearest success case.
- Original resources:
  - `/home/manuel/code/wesen/obsidian-vault/Projects/2026/03/27/PROJ - Sanitize - YAML Sanitizing Deep Dive.md`
  - `ttmp/2026/03/27/SANITIZE-008--derive-blog-post-drafts-for-sanitize/reference/06-anatomy-of-the-yaml-pipeline.md`
- Potential pseudocode and examples:
  - short “Part I / Part II” roadmap block.

---

### Part I: The YAML Sanitizer Itself

#### `Why YAML was the right place to start`
- Intent:
  - explain why YAML made a better first target than JSON.
- Content:
  - YAML has many locally fixable formatting failures,
  - parseability versus health,
  - duplicate keys and style-shaped errors still matter even when the parser accepts them.
- Original resources:
  - `/home/manuel/code/wesen/obsidian-vault/Projects/2026/03/27/PROJ - Sanitize - YAML Sanitizing Deep Dive.md`
  - `pkg/yaml/sanitize_test.go`
- Potential pseudocode and examples:
  - short YAML-vs-JSON comparison table,
  - examples: `tab_indent`, `missing_space_after_colon`, `duplicate_key`.

#### `The early regexp phase`
- Intent:
  - make it clear that the project did not begin with grand architecture.
- Content:
  - line-by-line cleanup,
  - why regexes were attractive,
  - the first rules that were genuinely useful.
- Original resources:
  - `pkg/yaml/lint.go`
  - `pkg/yaml/fix.go`
  - `ttmp/2026/03/13/SANITIZE-001--turn-yaml-sanitizer-into-reusable-go-go-golems-sanitize-package/`
- Potential pseudocode and examples:
  - `for each line -> run rules -> rewrite line`,
  - examples of tab indentation and list-dash fixes.

#### `Where regexes stopped being enough`
- Intent:
  - explain the architectural pressure for structural parsing.
- Content:
  - duplicate keys need structure,
  - parser failures can span multiple rows,
  - some problems are heuristic-only and some are parse-only,
  - direct line targeting can be misleading.
- Original resources:
  - `ttmp/2026/03/13/SANITIZE-004--improve-yaml-linting-and-sanitizing-with-tree-sitter-aware-analysis/design-doc/01-tree-sitter-driven-yaml-linting-and-sanitizing-analysis.md`
  - `ttmp/2026/03/13/SANITIZE-004--improve-yaml-linting-and-sanitizing-with-tree-sitter-aware-analysis/sources/01-example-corpus-parse-vs-lint-matrix.md`
- Potential pseudocode and examples:
  - parse-only / heuristic-only / hybrid matrix,
  - `20-mixed-indent.yaml`,
  - tab-indentation example where parser location and actionable location differ.

#### `Tree-sitter setup: what it gives and what it does not`
- Intent:
  - explain tree-sitter as infrastructure rather than magic.
- Content:
  - parser setup,
  - `ERROR` and `MISSING` nodes,
  - S-expression trees for inspection,
  - why tree-sitter alone does not tell you how to repair a file.
- Original resources:
  - `pkg/yaml/parse.go`
  - `pkg/yaml/analysis.go`
  - `ttmp/2026/03/27/SANITIZE-008--derive-blog-post-drafts-for-sanitize/reference/05-tree-sitter-as-infrastructure-not-a-magic-wand.md`
- Potential pseudocode and examples:
  - `ParseTree(src) -> treeText + errors`,
  - one tiny S-expression example,
  - a callout box: “parser witness, not repair oracle.”

#### `The real turning point: shared analysis`
- Intent:
  - identify the genuine architectural pivot.
- Content:
  - `documentAnalysis`,
  - shared tree text,
  - parse errors,
  - duplicate keys,
  - line index.
- Original resources:
  - `pkg/yaml/analysis.go`
  - `pkg/yaml/line_index.go`
  - `/home/manuel/code/wesen/obsidian-vault/Projects/2026/03/27/PROJ - Sanitize - YAML Sanitizing Deep Dive.md`
- Potential pseudocode and examples:
  - the `documentAnalysis` struct,
  - short explanation of why each field exists.

#### `Spans, line indexes, and unified diagnostics`
- Intent:
  - show that diagnostic shape matters as much as parser presence.
- Content:
  - row/column spans,
  - byte ranges,
  - `LintIssue` becoming richer,
  - parser and heuristic issues sharing one model.
- Original resources:
  - `pkg/yaml/types.go`
  - `pkg/yaml/lint.go`
  - `ttmp/2026/03/27/SANITIZE-008--derive-blog-post-drafts-for-sanitize/reference/13-how-to-write-span-aware-diagnostics-for-text-repair-tools.md`
- Potential pseudocode and examples:
  - issue struct snippet,
  - before/after comparison: row-only versus span-aware diagnostics.

#### `The heuristics that survived the refactor`
- Intent:
  - explain that tree-sitter did not kill heuristics; it disciplined them.
- Content:
  - tab indentation,
  - missing space after colon,
  - list dash,
  - trailing comma,
  - duplicate keys,
  - extra colon in value,
  - mixed indentation.
- Original resources:
  - `pkg/yaml/lint.go`
  - `pkg/yaml/fix.go`
  - `pkg/yaml/rules.go`
- Potential pseudocode and examples:
  - table with rule, why it exists, whether it lints, whether it fixes, and what makes it safe.

#### `The sanitize loop`
- Intent:
  - walk the reader through runtime behavior.
- Content:
  - capture original state,
  - analyze,
  - collect issues,
  - apply small fixes,
  - repeat,
  - stop on convergence or no progress.
- Original resources:
  - `pkg/yaml/sanitize.go`
  - `pkg/yaml/fix.go`
- Potential pseudocode and examples:
  - main sanitize loop pseudocode,
  - mermaid flowchart.

#### `Corpus and experiments: why the YAML design got better`
- Intent:
  - connect code to evidence.
- Content:
  - corpus as experiment harness,
  - parse-versus-lint matrix,
  - why categories A/B/C mattered,
  - how the corpus justified the shared-analysis refactor.
- Original resources:
  - `examples/yaml/`
  - `ttmp/2026/03/13/SANITIZE-004--.../scripts/parse_lint_matrix.go`
  - `ttmp/2026/03/13/SANITIZE-004--.../sources/01-example-corpus-parse-vs-lint-matrix.md`
- Potential pseudocode and examples:
  - tiny “run corpus -> generate matrix -> refine design” loop diagram,
  - a short list of representative fixtures.

#### `A full walkthrough of one broken file`
- Intent:
  - prove the system with a concrete example.
- Content:
  - take `examples/yaml/24-deeply-nested-mixed-errors.yaml`,
  - walk through what the parser sees,
  - what lint sees,
  - what fixes apply,
  - what the final result becomes.
- Original resources:
  - `examples/yaml/24-deeply-nested-mixed-errors.yaml`
  - current CLI output
  - `pkg/yaml/sanitize_test.go`
- Potential pseudocode and examples:
  - copy the CLI output block,
  - annotate each fix with why it is safe.

#### `What is still inelegant or unfinished`
- Intent:
  - keep the article credible.
- Content:
  - internal-only analysis types,
  - duplicate-key rename policy,
  - remaining public-API questions,
  - places where YAML still relies on judgment more than abstraction.
- Original resources:
  - `pkg/yaml/analysis.go`
  - `pkg/yaml/types.go`
  - `SANITIZE-004` design doc open questions
- Potential pseudocode and examples:
  - short “if I revisited this now” note.

---

### Part II: The Workflow Around the Code

#### `Early prototype and import mindset`
- Intent:
  - explain how the project began in a rougher, more exploratory mode.
- Content:
  - early sanitizer idea,
  - why an import / prototype mindset was useful,
  - how rough code became acceptable because the problem was still being discovered.
- Original resources:
  - `SANITIZE-001`
  - `README.md`
- Potential pseudocode and examples:
  - milestone table instead of prose-heavy history.

#### `The web playground changed the kind of project this was`
- Intent:
  - explain the workflow value of the UI.
- Content:
  - inspect parse trees,
  - load example fixtures,
  - compare original and sanitized output,
  - use the UI as an experiment loop instead of just a demo.
- Original resources:
  - `internal/server/server.go`
  - `internal/server/static/index.html`
  - `internal/server/static/js/app.js`
  - `SANITIZE-005`
- Potential pseudocode and examples:
  - UI flow diagram: example select -> analyze -> inspect issues -> compare output.

#### `Splitting the UI made experimentation easier`
- Intent:
  - show one concrete workflow improvement.
- Content:
  - splitting monolithic HTML into `index.html`, `css/style.css`, `js/app.js`,
  - why that mattered for iteration, not just neatness.
- Original resources:
  - `ttmp/2026/03/13/SANITIZE-005--split-ui-into-css-js-html-and-improve-preset-loading/reference/01-implementation-diary.md`
- Potential pseudocode and examples:
  - before/after file layout.

#### `Docmgr changed the project from hacking to research`
- Intent:
  - make the ticket workflow feel integral, not bureaucratic.
- Content:
  - new tickets after initial import,
  - diaries,
  - design docs,
  - tasks,
  - changelogs,
  - relating files back to docs.
- Original resources:
  - `ttmp/2026/03/13/SANITIZE-004--...`
  - `ttmp/2026/03/13/SANITIZE-005--...`
  - `ttmp/2026/03/13/SANITIZE-007--...`
- Potential pseudocode and examples:
  - “ticket -> experiment -> design doc -> code -> diary -> review” loop.

#### `The corpus-workflow loop`
- Intent:
  - show how experimentation stayed grounded.
- Content:
  - add example,
  - run script,
  - inspect matrix,
  - adjust rule or design,
  - verify in tests and UI.
- Original resources:
  - `examples/`
  - `ttmp/.../scripts/`
  - `ttmp/.../sources/`
- Potential pseudocode and examples:
  - a very short numbered loop with concrete commands.

#### `Using LLMs as critics and accelerators, not final authors`
- Intent:
  - describe the real assistant workflow without making it mystical.
- Content:
  - using LLMs to suggest gaps, summarize tickets, compare approaches, and assess edge cases,
  - still grounding decisions in code and corpus,
  - why the best use was evaluative and organizational rather than “just write the feature.”
- Original resources:
  - current SANITIZE-008 writing ticket itself,
  - prior ticket docs and diaries,
  - repo artifacts produced from those loops.
- Potential pseudocode and examples:
  - optional sidebar: “Where LLMs were useful / where they were not.”

#### `Why the workflow made the implementation better`
- Intent:
  - tie Part II back to Part I.
- Content:
  - UI exposed missing result fields,
  - tickets clarified what the parser was actually good for,
  - corpus reports prevented hand-wavy claims,
  - docs kept later work from rediscovering old reasoning.
- Original resources:
  - `SANITIZE-004`, `SANITIZE-005`, `SANITIZE-007`
  - `internal/server/server.go`
  - `pkg/yaml/sanitize.go`
- Potential pseudocode and examples:
  - “problem -> workflow artifact -> code improvement” table.

#### `What I would keep if I started over`
- Intent:
  - close with strong transferable lessons.
- Content:
  - start with one format where local repairs are actually safe,
  - build a corpus early,
  - make a UI earlier than expected,
  - give experiments a ticket and a diary,
  - use LLMs to sharpen judgment, not replace it.
- Original resources:
  - the whole repo history
- Potential pseudocode and examples:
  - final checklist or callout box.

---

### Sections To Potentially Cut If The Article Runs Long

- Most JSON-specific repair detail beyond using JSON as contrast.
- Deep discussion of rule enumeration and validated rule selection.
- Detailed reMarkable / delivery workflow.
- Full timeline of every SANITIZE ticket; keep only the most illustrative ones.

### Feedback Prompts For Friends

Use these questions when circulating the draft:

1. Does the split between “implementation” and “workflow” feel natural, or does it feel like two separate articles awkwardly fused together?
2. Does Part I feel concrete enough to satisfy a technical reader who wants actual implementation detail?
3. Does Part II feel like a real engineering workflow story, or does it drift into process for process’s sake?
4. Is the opening hook strong enough before the article gets deep into internal details?
5. After reading the outline alone, is the main claim of the article obvious?

### Resources To Keep Nearby During Final Writing

- `README.md`
- `pkg/yaml/analysis.go`
- `pkg/yaml/lint.go`
- `pkg/yaml/fix.go`
- `pkg/yaml/sanitize.go`
- `pkg/yaml/sanitize_test.go`
- `examples/yaml/24-deeply-nested-mixed-errors.yaml`
- `ttmp/2026/03/13/SANITIZE-004--improve-yaml-linting-and-sanitizing-with-tree-sitter-aware-analysis/`
- `ttmp/2026/03/13/SANITIZE-005--split-ui-into-css-js-html-and-improve-preset-loading/`
- `ttmp/2026/03/13/SANITIZE-007--add-json-support-with-focus-on-malformed-llm-json-recovery/`
- `/home/manuel/code/wesen/obsidian-vault/Projects/2026/03/27/PROJ - Sanitize - YAML Sanitizing Deep Dive.md`
- `/home/manuel/code/wesen/obsidian-vault/Projects/2026/03/27/PROJ - Sanitize - JSON Recovery Experiments and Limits.md`

## Usage Examples

- This draft is ready to share with friends for structure and scope feedback before the prose pass.
- A very practical next step would be to mark which sections are definitely in, definitely out, and “appendix only.”
- If you want a shorter article later, split this draft at the Part I / Part II boundary and keep Part II as a companion post.

## Related

- `reference/02-from-regexes-to-a-real-sanitizer.md`
- `reference/06-anatomy-of-the-yaml-pipeline.md`
- `reference/09-how-the-browser-playground-changed-the-project.md`
- `reference/11-a-weekend-project-that-turned-into-a-small-research-lab.md`
