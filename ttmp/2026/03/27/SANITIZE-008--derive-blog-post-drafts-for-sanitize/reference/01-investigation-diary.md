---
Title: Investigation diary
Ticket: SANITIZE-008
Status: active
Topics:
    - writing
    - documentation
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ttmp/2026/03/27/SANITIZE-008--derive-blog-post-drafts-for-sanitize/changelog.md
      Note: Ticket history for the writing packet
    - Path: ttmp/2026/03/27/SANITIZE-008--derive-blog-post-drafts-for-sanitize/index.md
      Note: Ticket overview updated during this work
    - Path: ttmp/2026/03/27/SANITIZE-008--derive-blog-post-drafts-for-sanitize/tasks.md
      Note: Task checklist for the writing packet
ExternalSources: []
Summary: Chronological diary for building the SANITIZE-008 writing packet, including ticket setup, article derivation, validation, and delivery.
LastUpdated: 2026-03-27T09:34:25.452761266-04:00
WhatFor: Record how the blog-post drafting packet was assembled and how to review or continue the work later.
WhenToUse: Use when auditing the ticket, extending the article set, or re-running delivery and validation steps.
---


# Investigation diary

## Goal

Capture the concrete steps used to turn the sanitize repo and prior project notes into a bundle of article-ready draft documents.

## Step 1: Create the ticket and define the writing packet

The first step was to convert a casual writing request into a durable workspace. The main reason to do that was volume: twelve article ideas, ticket bookkeeping, validation, and reMarkable delivery are easier to manage as a real docmgr ticket than as a couple of ad-hoc markdown files.

The output of this step was a new `SANITIZE-008` ticket with an index, task list, changelog, a design doc, and a diary. I also confirmed the current ticket landscape so the new writing ticket would sit cleanly after `SANITIZE-007`.

### Prompt Context

**User prompt (verbatim):** "Ok, create a new ticket to derive potential blog posts. Then write a flesh out draft (not full paragraphs though, more like : intent of the section, content of the section, link to original resources, potential pseudocode and examples. I'll write things out myself in my own voice later on). Create one file per article idea. upload all as bundle to remarkable"

**Assistant interpretation:** Create a docmgr ticket, turn the article ideas into outline-heavy draft files, and upload the full packet to reMarkable.

**Inferred user intent:** Prepare real writing material without replacing the user’s authorship or tone.

**Commit (code):** N/A

### What I did
- Ran `docmgr status --summary-only`.
- Created `SANITIZE-008`.
- Added the design doc and diary docs.
- Added one reference document per article idea.

### Why
- The request maps cleanly onto the existing ticket workflow.
- A ticket keeps delivery, validation, and later continuation manageable.

### What worked
- `docmgr` created the workspace and documents cleanly.
- The sanitize ticket numbering made `SANITIZE-008` the obvious next id.

### What didn't work
- No hard failure in this step. Vocabulary still needed to be checked later during `docmgr doctor`.

### What I learned
- Treating writing prep as a ticket works well when the output is a real bundle rather than an ephemeral answer.

### What was tricky to build
- The scope decision mattered: one file per article idea creates more surface area, but it is much more reusable later than a monolithic brief.

### What warrants a second pair of eyes
- Whether the final article set should later be narrowed to a “best five” shortlist.

### What should be done in the future
- Add a ranking doc if the user later wants publishing priority guidance.

### Code review instructions
- Start with `index.md`, `tasks.md`, and the design doc in this ticket.
- Confirm that `reference/02` through `reference/13` exist and map to the article ideas.

### Technical details
- Commands used:
  - `docmgr status --summary-only`
  - `docmgr ticket create-ticket --ticket SANITIZE-008 --title "Derive blog post drafts for sanitize" --topics writing,documentation`
  - `docmgr doc add --ticket SANITIZE-008 --doc-type design-doc --title "Blog post outlines and article drafting plan for sanitize"`

## Step 2: Fill the drafts from the repo and prior notes

Once the ticket existed, the real work was choosing structures that were concrete enough to be useful later but still skeletal enough to preserve the user’s own prose. I used the repo code, the SANITIZE ticket history, and the two new Obsidian project reports as the main evidence base for the drafts.

The article files were then filled with recurring building blocks: thesis, section order, section intent, section content, original repo resources, and possible pseudocode or examples. This produced a writing packet rather than a set of generic titles.

### Prompt Context

**User prompt (verbatim):** (same as Step 1)

**Assistant interpretation:** Turn each article idea into a practical draft scaffold grounded in the repo and tickets rather than leaving it as a one-line idea.

**Inferred user intent:** Reduce drafting effort later by front-loading structure and evidence now.

**Commit (code):** N/A

### What I did
- Read the current sanitize code and ticket docs.
- Used the new Obsidian YAML and JSON reports as summary anchors.
- Filled article draft docs `reference/02` through `reference/13` with:
  - title and thesis,
  - section-by-section intent,
  - likely section content,
  - original resources,
  - candidate examples and pseudocode.

### Why
- The user explicitly wanted more than titles but less than finished prose.
- A section-by-section scaffold is the most useful middle ground.

### What worked
- The repo already had strong source material: code, tests, ticket docs, generated reports, and new vault notes.
- The article set naturally split into architecture, case-study, product-surface, and engineering-lesson pieces.

### What didn't work
- No hard tooling failures in this step.

### What I learned
- The strongest article ideas are the ones that combine implementation and decision history.
- The JSON story is especially strong as a “limits of repair” article because the repo has evidence for the boundary.

### What was tricky to build
- The sharp edge was avoiding assistant-style prose completion. I kept the files in planning mode by writing intent and content bullets instead of finished paragraphs.

### What warrants a second pair of eyes
- Which article candidates should be prioritized first; not all twelve are equally strong.

### What should be done in the future
- Add publication-target notes if the user later wants variants for blog, newsletter, or talk formats.

### Code review instructions
- Open any two or three draft files and confirm they contain real section plans rather than placeholders.
- Spot-check that the cited repo resources match the article claim.

### Technical details
- Main sources used:
  - `README.md`
  - `pkg/yaml/*`
  - `pkg/json/*`
  - `internal/server/*`
  - `ttmp/.../SANITIZE-004/...`
  - `ttmp/.../SANITIZE-007/...`
  - `/home/manuel/code/wesen/obsidian-vault/Projects/2026/03/27/PROJ - Sanitize - YAML Sanitizing Deep Dive.md`
  - `/home/manuel/code/wesen/obsidian-vault/Projects/2026/03/27/PROJ - Sanitize - JSON Recovery Experiments and Limits.md`

## Step 3: Consolidate the packet into one friend-reviewable flagship article draft

After the article packet existed, the next request was not for more ideas but for synthesis. The user wanted one larger article that already felt substantial enough to circulate, while still stopping short of fully written prose. That changed the job from “derive article candidates” to “compose a flagship draft from the strongest strands.”

The document I added in this step merges two lines that had been separate in the packet: the deep YAML implementation story and the surrounding workflow story. That produced a stronger article than either side alone because the code and the workflow explain each other.

### Prompt Context

**User prompt (verbatim):** "now let's create a consolidate article that is in two big parts:

- the implementation of the yaml sanitize part itself, in deep technical detail: spans, heuristics, test corpus, early experiments with regexp, treesitter setup, etc...

the second part is then about the workflow:
- early prototype
- web playground to experiment and test
- using docmgr and making tickets after initial import
- iterating and fixing and using LLMs to assess stuff
- etc...

Gther all the work and research you did for the initial draft to create a great fleshed out version for this article, notyet fully written out, but something that can already be handed out to friends for some feedback. 

Save and then upload to remarkable."

**Assistant interpretation:** Build one substantial two-part article draft out of the existing SANITIZE-008 material, update the ticket to track it, and upload that new draft to reMarkable.

**Inferred user intent:** Move from a set of candidate articles to one flagship piece that is already strong enough to circulate for structural feedback.

**Commit (code):** N/A

### What I did
- Added `reference/14-consolidated-two-part-article-draft-for-sanitize.md`.
- Folded material from the YAML implementation draft, the workflow/UI/corpus drafts, and the new Obsidian notes into one article structure.
- Added a “feedback prompts for friends” section so the draft can be circulated immediately.
- Prepared the new draft for standalone reMarkable upload rather than rebundling the entire ticket packet.

### Why
- The strongest current article is the one that combines internals and workflow.
- A friend-reviewable draft needs more narrative shape than the individual article idea files.

### What worked
- The earlier SANITIZE-008 packet made synthesis much faster because the core sections, resources, and examples were already decomposed.
- The two-part structure gives the article a clean backbone without forcing full prose yet.

### What didn't work
- The existing diary had placeholder scaffolding left at the end from the original doc template. I removed that dead structure while updating the diary so the ticket remains continuation-friendly.

### What I learned
- The workflow material is not just “Part II filler.” It explains why the YAML implementation got as good as it did.
- The best article shape for this project is not “code only” or “retrospective only,” but a braid of both.

### What was tricky to build
- The main challenge was staying below “finished prose” while still making the document feel substantial enough to share. I handled that by writing stronger section-level narrative guidance, explicit article promises, and a feedback checklist instead of polished paragraphs.

### What warrants a second pair of eyes
- Whether the two-part article is the right final publication shape or whether it should eventually split into two companion essays.

### What should be done in the future
- If the draft gets positive feedback, the next step is not more outlining but a prose pass plus selective cutting for length.

### Code review instructions
- Read `reference/14-consolidated-two-part-article-draft-for-sanitize.md` top to bottom and check whether the Part I / Part II split feels coherent.
- Confirm that each major section still points back to repo-local source material.

### Technical details
- New document:
  - `ttmp/2026/03/27/SANITIZE-008--derive-blog-post-drafts-for-sanitize/reference/14-consolidated-two-part-article-draft-for-sanitize.md`
- Source drafts folded in:
  - `reference/02-from-regexes-to-a-real-sanitizer.md`
  - `reference/05-tree-sitter-as-infrastructure-not-a-magic-wand.md`
  - `reference/06-anatomy-of-the-yaml-pipeline.md`
  - `reference/08-designing-a-corpus-for-malformed-structured-text.md`
  - `reference/09-how-the-browser-playground-changed-the-project.md`
  - `reference/11-a-weekend-project-that-turned-into-a-small-research-lab.md`
