---
Title: Public release review
Ticket: SANITIZE-002
Status: active
Topics:
    - json
    - api-design
    - release-readiness
    - documentation
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: README.md
      Note: Published contract and release-facing claims
    - Path: cmd/sanitize-server/main.go
      Note: HTTP server hardening and API behavior
    - Path: cmd/sanitize/main.go
      Note: CLI behavior and exit-code handling
    - Path: pkg/yaml/fix.go
      Note: Duplicate-key renaming and fixer behavior
    - Path: pkg/yaml/lint.go
      Note: Duplicate-key detection and lint heuristics
    - Path: pkg/yaml/parse.go
      Note: Parse extraction and gosec finding
    - Path: pkg/yaml/sanitize_test.go
      Note: Existing automated coverage baseline
ExternalSources: []
Summary: ""
LastUpdated: 2026-03-13T00:47:57.912136143-04:00
WhatFor: ""
WhenToUse: ""
---


# Public release review

## Executive Summary

The repository is close to a publishable v0 in terms of packaging and baseline engineering hygiene: `go test ./...`, `go test -race ./...`, `make lint`, and `make build` all passed on March 13, 2026. The module path, README, CLI entrypoint, server entrypoint, GitHub Actions, Goreleaser config, and unit tests are all present.

The code is not yet fit for an unconstrained public release. The main blockers are correctness and operational issues that would surprise downstream users: the duplicate-key heuristic mutates valid YAML, the CLI's JSON mode hides failure via exit code `0`, and the HTTP server is missing basic hardening such as timeouts. `gosec` also reports an integer-conversion bug risk in parse error extraction. Those items should be addressed before calling the package publicly safe and automation-friendly.

## Problem Statement

The user asked for a thorough release review of the current `sanitize` codebase and a judgment about whether it is ready for public release. That requires more than checking whether the code compiles. It requires reviewing:

1. Behavioral correctness of the YAML sanitizer.
2. CLI/API contracts exposed to users and automation.
3. Server hardening for a public binary.
4. Test coverage relative to the public surface area.
5. Release plumbing such as CI and packaging.

This document focuses on release readiness of the code as it exists today, not on the future JSON-support design. JSON support is covered in the companion design document.

## Proposed Solution

Treat the release as "nearly ready, but not yet releasable" and gate publication on a short pre-release hardening pass:

1. Fix the duplicate-key algorithm so it only considers actual sibling keys, not every key at the same indentation depth across the document.
2. Make CLI exit codes consistent in machine-readable modes.
3. Harden the HTTP server with explicit timeouts and bounded request bodies.
4. Remove or justify the `uint -> int` conversion in parse error extraction.
5. Add targeted regression tests for the issues above.

Once those changes are in place, rerun the validation commands listed below and re-evaluate whether any remaining gaps are acceptable for a `v0.x` public release.

## Design Decisions

### Findings are prioritized by user impact

This review is ordered by the risk each issue poses to external users. A release review that starts with stylistic preferences would hide the real blockers.

### Passing tests are treated as necessary but insufficient

The current repository has a healthy baseline: the package tests pass, race tests pass, and CI config exists. Those are positive signals, but they do not outweigh user-visible correctness bugs.

### Security and operability issues are counted as release issues even for a small local server

`sanitize-server` may start as a local tool, but public binaries tend to get used in less controlled ways. Missing server timeouts and unbounded request handling should be treated as release-readiness issues rather than "nice to have" cleanup.

## Findings

### 1. High: duplicate-key detection rewrites valid YAML in different mappings

Evidence:

- `pkg/yaml/lint.go:24-25` keeps a single `seenKeys` map for the entire document.
- `pkg/yaml/lint.go:72-85` uses `indent + key` as the deduplication key.
- `pkg/yaml/fix.go:182-208` uses the same heuristic to rename keys during fixing.

Why this is a release blocker:

Two different mappings can legitimately reuse the same field name at the same indentation level. The current algorithm treats those as duplicates even when they are under different parents.

Observed reproduction:

```text
printf 'a:\n  timeout: 30\nb:\n  timeout: 60\n' | go run ./cmd/sanitize --json
```

Observed result:

```text
"sanitized": "a:\n  timeout: 30\nb:\n  timeout_2: 60\n"
```

That mutation is semantically incorrect. A public sanitizer cannot silently rewrite valid input this way.

Recommended fix:

- Derive duplicate detection from the parse tree instead of indentation alone.
- At minimum, track a path stack of parent mappings so the deduplication key is `parent-path + key`, not `indent + key`.
- Add a regression test for the exact reproduction above.

### 2. Medium: `sanitize --lint --json` returns exit code 0 even when lint issues are present

Evidence:

- `cmd/sanitize/main.go:35-47` returns immediately after JSON encoding and only exits non-zero in the text-output path.
- The README claims lint mode exits `1` when issues are found: `README.md:24-31`.

Observed reproduction:

```text
printf 'name:Alice\n' | go run ./cmd/sanitize --lint --json
```

The command emitted a JSON issue list and exited with code `0`.

Why this matters:

Machine-readable modes are usually the ones used in automation. A CI step or editor integration relying on the exit code would incorrectly treat dirty input as success.

Recommended fix:

- Compute the exit code independently of output format.
- Apply the same rule to `--json` in normal sanitize mode when `ParseClean` or `LintClean` is false.
- Add CLI tests that assert both payload and exit code.

### 3. Medium: `sanitize-server` is missing HTTP timeouts and request-size bounds

Evidence:

- `cmd/sanitize-server/main.go:23-99` builds handlers on a `ServeMux`.
- `cmd/sanitize-server/main.go:99` calls `http.ListenAndServe`.
- `cmd/sanitize-server/main.go:59` and `cmd/sanitize-server/main.go:81` decode request bodies directly with `json.NewDecoder`.
- `make gosec` reported `G114` for the timeout-free server.

Why this matters:

Without `ReadHeaderTimeout`, `ReadTimeout`, `WriteTimeout`, and `IdleTimeout`, the server is vulnerable to slow-client and resource-exhaustion problems. Without `http.MaxBytesReader`, a client can send an arbitrarily large JSON body and force large allocations.

Recommended fix:

- Replace `http.ListenAndServe` with an explicit `http.Server`.
- Set sane default timeouts.
- Wrap request bodies with `http.MaxBytesReader`.
- Document that the server is for local use unless additional auth/origin controls are added.

### 4. Medium: parse error extraction uses a flagged `uint -> int` conversion

Evidence:

- `pkg/yaml/parse.go:54-58` does `int(eb) <= len(src)` before slicing.
- `make gosec` reported `G115` on `pkg/yaml/parse.go:56`.

Why this matters:

In typical inputs this will probably not overflow, but public libraries should avoid unchecked narrowing conversions in parser-adjacent code. The fix is small and the warning is reasonable.

Recommended fix:

- Compare `eb` against `uint(len(src))` instead of converting `eb` to `int`.
- Add a comment that the bounds check is intentionally kept in unsigned space.

### 5. Medium: the public surface lacks tests for the CLI and server behaviors that matter most

Evidence:

- `pkg/yaml/sanitize_test.go` covers the core YAML package extensively.
- `cmd/sanitize` and `cmd/sanitize-server` have no tests (`go test ./...` reported `[no test files]` for both commands).
- The web UI in `cmd/sanitize-server/static/index.html` is 554 lines and has no automated coverage.

Why this matters:

The public-facing regressions found in this review are outside the well-tested package core:

- CLI exit-code behavior.
- HTTP request/response behavior.
- UI/API coupling.

That means the current test suite is strong on local rule helpers but weaker on user-facing contracts.

Recommended fix:

- Add CLI tests using subprocess execution against fixtures.
- Add HTTP handler tests with `httptest`.
- Add at least one smoke test that exercises the `/api/sanitize` contract end-to-end.

### 6. Low: JSON output currently emits `null` for absent slices

Evidence:

- `pkg/yaml/types.go:47-53` exposes slices directly.
- Observed CLI JSON output emitted `"errors": null` and `"lint_issues": null` for clean results.

Why this matters:

This is not a correctness bug, but it makes client code noisier and can surprise users expecting empty arrays. The frontend already works around this with `|| []` in `cmd/sanitize-server/static/index.html:349-391`.

Recommended fix:

- Normalize slices to empty slices in results intended for JSON output, or document the `null` behavior explicitly.

## Alternatives Considered

### "Release now because tests and lint pass"

Rejected because the duplicate-key bug changes user data incorrectly. A green build is not enough when the core tool can rewrite valid documents.

### "Treat the server hardening issues as out of scope because this is only a local tool"

Rejected because the repository ships `cmd/sanitize-server` as a public installable binary and documents it in `README.md:33-41`. If it is publicly distributed, baseline hardening should be part of release quality.

### "Ignore machine-readable exit codes because JSON callers can inspect the payload"

Rejected because shell automation and CI routinely depend on exit status. The fix is straightforward and the current behavior is surprising.

## Implementation Plan

### Release gate checklist

1. Fix duplicate-key detection and add a regression test.
2. Fix CLI exit-code handling for JSON output modes.
3. Harden the HTTP server and add handler tests.
4. Resolve the `G115` parse warning cleanly.
5. Re-run:
   - `go test ./...`
   - `go test -race ./...`
   - `make lint`
   - `make build`
   - `make gosec`

### Validation performed for this review

The following commands were run on March 13, 2026:

```bash
go test ./...
go test -race ./...
make lint
make build
make gosec
printf 'a:\n  timeout: 30\nb:\n  timeout: 60\n' | go run ./cmd/sanitize --json
printf 'name:Alice\n' | go run ./cmd/sanitize --lint --json
```

Results summary:

- `go test ./...`: passed
- `go test -race ./...`: passed
- `make lint`: passed
- `make build`: passed
- `make gosec`: failed with `G115`, `G114`, and `G706`
- duplicate-key reproduction: demonstrated incorrect key renaming
- lint JSON reproduction: demonstrated incorrect exit code behavior

## Open Questions

1. Is the intended contract that JSON output mode should preserve non-zero exit codes on dirty input? This review assumes yes.
2. Is `sanitize-server` intended for local-only use, or should the public docs position it as safe to expose on a network?
3. Has `pkg/yaml` already been consumed externally enough that API-moving cleanup should wait until JSON support lands?

## Release Verdict

The repository is promising and operationally close, but it is not yet fit for a public release without a short hardening pass. The main reason is that valid YAML can currently be rewritten incorrectly. Fixing the issues above should be a bounded piece of work, not a rewrite.

## References

- `README.md:1-82`
- `cmd/sanitize/main.go:13-70`
- `cmd/sanitize-server/main.go:17-101`
- `cmd/sanitize-server/static/index.html:325-391`
- `pkg/yaml/lint.go:19-110`
- `pkg/yaml/fix.go:16-316`
- `pkg/yaml/parse.go:20-79`
- `pkg/yaml/sanitize.go:5-80`
- `pkg/yaml/sanitize_test.go:10-364`
- `.github/workflows/lint.yml:1-34`
- `.github/workflows/push.yml:1-26`
- `.github/workflows/release.yaml:1-117`
- `.golangci.yml:1-37`
- `.goreleaser.yaml:1-111`
