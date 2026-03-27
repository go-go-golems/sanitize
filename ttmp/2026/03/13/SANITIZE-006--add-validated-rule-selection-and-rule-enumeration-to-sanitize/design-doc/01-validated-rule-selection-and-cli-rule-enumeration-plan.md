---
Title: Validated rule selection and CLI rule enumeration plan
Ticket: SANITIZE-006
Status: active
Topics:
    - yaml
    - cli
    - linting
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2026-03-13T10:20:39.55048815-04:00
WhatFor: Define a clean architecture for validating YAML sanitizer rule names and exposing rule selection and enumeration in the Glazed CLI.
WhenToUse: Use when implementing or reviewing rule registry, rule filtering, and CLI rule selection behavior.
---

# Validated rule selection and CLI rule enumeration plan

## Executive Summary

The sanitize package already has a partial rule-selection mechanism in `WithRules(...)`, but it is not backed by a canonical rule registry and it is not surfaced in the Glazed CLI. As a result, callers can pass unknown rule names silently, the CLI cannot validate or enumerate rules, and the codebase duplicates rule knowledge across linting, fixing, tests, and command help.

This ticket introduces a single rule catalog in `pkg/yaml`, builds validated rule filtering on top of it, and wires that capability into the CLI via `--rule`, `--disable-rule`, and a dedicated `sanitize rules` command. The goal is to make rule selection explicit, discoverable, and safe without scattering more string literals across the codebase.

## Problem Statement

The current implementation has three concrete gaps:

- Rule names are string literals spread across `pkg/yaml/lint.go`, `pkg/yaml/fix.go`, tests, and the existing `WithRules(...)` helper in `pkg/yaml/options.go`.
- `WithRules(...)` does not validate its input, so typos like `WithRules("missing-space-after-colon")` quietly produce surprising behavior.
- The CLI does not let users constrain lint or fix runs to a subset of rules, nor can it list what rules exist.

That creates both correctness and usability problems. Correctness suffers because configuration mistakes are not surfaced. Usability suffers because users and future contributors have no canonical place to inspect available rules or understand which ones lint only versus which ones can also fix.

## Proposed Solution

Add a canonical rule registry in `pkg/yaml` and make every rule-related behavior derive from it.

The concrete shape is:

- Add `pkg/yaml/rules.go` with a `RuleSpec` type and a stable list of all known rules.
- Expose helpers such as:
  - `RuleCatalog() []RuleSpec`
  - `KnownRule(name string) bool`
  - `ValidateRuleNames(names ...string) error`
- Replace the current single `enabledRules` map with a more explicit selection model:
  - optional allowlist of enabled rules
  - optional disable-list of rules to suppress
- Introduce error-returning option-aware entrypoints for package callers, for example:
  - `LintWithOptions(src string, opts ...Option) ([]LintIssue, error)`
  - `SanitizeWithOptions(src string, opts ...Option) (Result, error)`
- Keep `Lint(src)` and `Sanitize(src, ...)` as convenience wrappers so existing internal callers remain simple.
- Add CLI flags on `lint` and `fix`:
  - `--rule <name>` repeated
  - `--disable-rule <name>` repeated
- Add a `sanitize rules` command that prints all rule metadata in text or JSON.

The key invariant is that the CLI and the YAML package both derive rule validation and help text from the same registry.

## Design Decisions

### Decision 1: Put the canonical registry in `pkg/yaml`

The CLI should not maintain its own copy of rule names. The package owns the rules, so the package should own the canonical metadata too. That keeps tests and future APIs aligned around one source of truth.

### Decision 2: Separate selection from validation

Rule-selection state belongs in config, but validation belongs in registry helpers. That keeps rule checks consistent across package APIs and CLI entrypoints.

### Decision 3: Support both allowlist and disable-list flows

Two common workflows matter:

- “run only these rules”
- “run all normal rules except this one”

Supporting only `WithRules(...)` forces callers into the first style even when they only want a small exclusion.

### Decision 4: Add explicit error-returning APIs instead of encoding config errors inside `Result`

Invalid rule names are configuration failures, not lint findings. They should return regular Go errors and regular CLI failures, not masquerade as YAML diagnostics.

### Decision 5: Keep rule listing as a dedicated CLI command

`sanitize rules` is easier to discover and document than overloading `lint --list-rules` or `fix --list-rules`. It also avoids awkward interactions between “list rules” and “run a lint/fix pass”.

## Alternatives Considered

### Keep `WithRules(...)` as-is and validate only in the CLI

Rejected because package callers would still get silent misconfiguration, and the CLI would still need to reverse-engineer the list of known rules.

### Infer the rule list from lint and fix code at runtime

Rejected because the rule definitions are intentionally not all in one place today, and reflection or regex-based inference would be brittle.

### Make unknown rules warnings instead of hard errors

Rejected because a mistyped rule name is almost always operator error. Hard failure is safer and easier to debug.

## Implementation Plan

### Step 1: Create the registry

- Add `pkg/yaml/rules.go`.
- Define `RuleSpec`.
- Encode every current rule name exactly once.
- Add helpers for listing and validating rule names.

### Step 2: Refactor config and package entrypoints

- Update `pkg/yaml/options.go` so config can carry:
  - allowlist
  - disable-list
- Add validation against the registry.
- Add error-returning entrypoints for lint and sanitize.
- Update internal lint and fix logic to consult a shared rule-enabled check.

### Step 3: Update CLI commands

- Extend `fixSettings` and `lintSettings` in `internal/cli/commands.go`.
- Add repeated string-list flags using Glazed `fields.TypeStringList`.
- Pass the selected rules into the YAML package.
- Add `sanitize rules` with text and JSON output.

### Step 4: Add tests and docs

- Add package tests for rule validation and rule filtering.
- Add CLI tests for:
  - valid allowlist selection
  - valid disable-list selection
  - invalid rule failure
  - `sanitize rules`
- Record each slice in the diary and changelog.

## Open Questions

- Whether `Lint(src)` should remain all-rules-always or eventually gain a package-global config object. For this ticket, the simplest approach is to keep `Lint(src)` as the all-default-rules convenience wrapper.
- Whether parse-only rules like `structural_parse_error` should be disableable. My current assumption is yes: if a rule has a stable name and appears in lint output, it should participate in the same rule-selection system.

## References

- Current rule-selection helper: `pkg/yaml/options.go`
- Current rule literals and lint behavior: `pkg/yaml/lint.go`
- Current fix behavior: `pkg/yaml/fix.go`
- Current Glazed command wiring: `internal/cli/commands.go`
