---
Title: Investigation Diary
Ticket: sanitize-008
Status: active
Topics:
    - deployment
    - kubernetes
    - ghcr
    - argocd
    - gitops
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - /home/manuel/code/wesen/corporate-headquarters/sanitize/internal/server/server.go:Go HTTP server with embedded static files
    - /home/manuel/code/wesen/corporate-headquarters/sanitize/go.mod:Module definition showing tree-sitter CGO deps
    - /home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/pretext/:Reference deployment pattern
    - /home/manuel/code/wesen/2026-03-27--hetzner-k3s/docs/public-repo-ghcr-argocd-deployment-playbook.md:GHCR deployment playbook
ExternalSources: []
Summary: "Chronological diary of investigating and planning the sanitize web UI deployment to K3s."
LastUpdated: 2026-03-29T17:01:11.520246827-04:00
WhatFor: ""
WhenToUse: ""
---

# Diary

## Goal

Document the investigation and planning process for deploying the Sanitize Playground web UI to the Hetzner K3s cluster via GHCR and Argo CD.

## Step 1: Understand the Sanitize Application Architecture

Investigated the sanitize repository to understand what needs to be deployed and how the web server works.

### Prompt Context

**User prompt (verbatim):** "Create a new docmgr ticket to deploy the sanitize web ui statically to k3s. Our k3s management gitops repository is at /home/manuel/code/wesen/2026-03-27--hetzner-k3s and contains documentation and examples on how to deploy an app to GHCR and then to the cluster. Analyze the situation and create an implementation plan to deploy sanitize's experimentation UI. Create a detailed analysis and implementation guide detailed for a new intern and explaining all the setup and necessary context, and store in the ticket and uploaed to remarkable. Keep a diary as you work."

**Assistant interpretation:** Create a docmgr ticket with a comprehensive implementation guide for deploying the sanitize web UI to k3s, following the established GHCR + ArgoCD pattern in the gitops repo.

**Inferred user intent:** Get a complete, intern-readable deployment plan that can be handed off and executed step by step.

### What I did

- Read `cmd/sanitize/main.go` — entry point, delegates to `sanitizecli.Run()`
- Read `internal/cli/root.go` — Cobra root command with subcommands: fix, lint, rules, parse, serve
- Read `internal/cli/commands.go` — all command implementations including `serveCommand` which calls `server.Run(port)`
- Read `internal/server/server.go` — the HTTP server implementation
- Listed `internal/server/static/` — contains `index.html`, `css/style.css`, `js/app.js`
- Read `internal/server/static/index.html` — the Sanitize Playground UI
- Read `internal/server/static/js/app.js` — frontend JavaScript (fetch-based API calls to `/api/sanitize`, `/api/parse`, `/api/examples`)
- Read `go.mod` — found tree-sitter dependencies requiring CGO
- Read `Makefile` — build targets, no Docker target yet
- Read `.github/workflows/push.yml` — existing CI runs tests only
- Read `.github/workflows/release.yaml` — GoReleaser binary releases (no container image)

### Why

Need to understand exactly what the deployment artifact looks like before designing the Dockerfile and Kubernetes manifests. The CGO dependency is a critical constraint that affects the entire Dockerfile design.

### What worked

- The server architecture is very clean: single binary, embedded static files via `//go:embed`, three API endpoints, no external dependencies
- This is the simplest possible deployment scenario — a self-contained binary serving HTTP

### What didn't work

- N/A (investigation phase)

### What I learned

1. **CGO is required** — tree-sitter uses native C code. This rules out `CGO_ENABLED=0` and `scratch`/`distroless/static` base images.
2. **Port 8080** is the default, configurable via `--port` flag
3. **Static files are embedded** — no volume mounts or config files needed at runtime
4. **No health endpoint** exists — but `/` (index.html) works fine for probes
5. **Examples are built into the binary** — both hardcoded examples and file-loaded examples (from `examples/` package using `//go:embed`)
6. **CORS headers are set** (`Access-Control-Allow-Origin: *`) — good for development, acceptable for this use case

### What was tricky to build

The main subtlety is the CGO dependency chain. A developer unfamiliar with tree-sitter might try `CGO_ENABLED=0` (as the demo app Dockerfile in the k3s repo does), get a build failure, and not understand why. The design doc explains this explicitly.

### What warrants a second pair of eyes

- Whether `alpine:3.21` provides the right musl version for the tree-sitter bindings. This should be validated by building the Docker image locally before pushing to CI.

### What should be done in the future

- N/A for this step

### Code review instructions

- Check `internal/server/server.go` lines 1-30 for the `//go:embed` directive and default port
- Check `go.mod` for tree-sitter dependencies that require CGO

### Technical details

Server setup from `server.go`:

```go
//go:embed static
var staticFiles embed.FS

const DefaultPort = 8080

// Routes:
// /          -> http.FileServer (embedded static files)
// /api/examples -> examplesHandler
// /api/sanitize -> sanitizeHandler
// /api/parse    -> parseHandler
```

## Step 2: Analyze the K3s GitOps Repository Patterns

Investigated the gitops repository to understand existing deployment patterns and identify the correct template to follow.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Map the existing deployment patterns in the gitops repo to determine which one sanitize should follow.

**Inferred user intent:** Ensure the sanitize deployment follows established conventions and does not introduce unnecessary divergence.

### What I did

- Listed `gitops/applications/` — found 15 Argo CD Applications
- Listed `gitops/kustomize/` — found 15 Kustomize packages
- Read `gitops/applications/pretext.yaml` — standard Application manifest
- Read `gitops/kustomize/pretext/` — all 5 files (namespace, deployment, service, ingress, kustomization)
- Read `gitops/applications/draft-review.yaml` — complex app with Vault/VSO
- Read `gitops/kustomize/draft-review/kustomization.yaml` — 15 resources (vault, bootstrap, secrets)
- Read `app/Dockerfile` — demo app Dockerfile (CGO_ENABLED=0, bare alpine)
- Read `docs/public-repo-ghcr-argocd-deployment-playbook.md` — full GHCR deployment tutorial
- Read `docs/app-packaging-and-gitops-pr-standard.md` — standard packaging model
- Read `docs/source-app-deployment-infrastructure-playbook.md` — infrastructure playbook
- Read `docs/argocd-app-setup.md` — Argo CD Application setup guide

### Why

The gitops repo has four categories of deployment complexity. Need to identify which category sanitize falls into and use the right pattern.

### What worked

- The documentation in the gitops repo is excellent and clearly categorizes deployments
- Pretext is the perfect reference pattern — same shape (stateless web app), same infrastructure needs

### What didn't work

- N/A

### What I learned

1. **Category 1 (public stateless app)** is the right fit — namespace, deployment, service, ingress, kustomization
2. **Pretext still uses `imagePullPolicy: Never`** with a locally-imported image — our deployment will be the first in this category to use GHCR properly
3. **The Argo CD Application format is standardized** — all apps use the same structure with automated sync, prune, selfHeal, CreateNamespace, ServerSideApply
4. **TLS is automatic** — cert-manager with `letsencrypt-prod` cluster-issuer handles certificates
5. **Ingress uses Traefik** — `ingressClassName: traefik`
6. **Wildcard DNS** — `*.yolo.scapegoat.dev` routes to the cluster
7. **Sync waves** — namespace at wave 0, deployment/service at wave 1, ingress at wave 2
8. **The sha-tag mismatch problem** is documented — mysql-ide had a real incident where the full SHA was used in GitOps but GHCR published a 7-char short SHA

### What was tricky to build

Understanding the four deployment categories and confirming sanitize fits Category 1. The gitops repo has complex deployments (draft-review with 15 resources, vault integration, bootstrap jobs) that could mislead someone into over-engineering the sanitize deployment.

### What warrants a second pair of eyes

- Confirm the wildcard DNS `*.yolo.scapegoat.dev` is actually configured
- Confirm cert-manager's `letsencrypt-prod` issuer is working (check existing certificates)

### What should be done in the future

- After sanitize is deployed, update pretext to also use GHCR instead of `imagePullPolicy: Never`

### Code review instructions

- Compare the proposed sanitize manifests against `gitops/kustomize/pretext/` for consistency
- Verify the Argo CD Application matches the standard pattern in `gitops/applications/pretext.yaml`

### Technical details

Deployment categories in the gitops repo:

| Category | Example | Resources | Secrets |
|----------|---------|-----------|---------|
| 1: Stateless | pretext | 5 files | None |
| 2: Vault/VSO | coinvault | ~10 files | Vault secrets |
| 3: Bootstrap | keycloak | ~15 files | Vault + bootstrap jobs |
| 4: Data service | postgres | ~8 files | Vault + statefulset |

Sanitize is Category 1.

## Step 3: Write the Implementation Guide

Synthesized findings from Steps 1 and 2 into a comprehensive design document.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Write the design doc with enough detail for a new intern to execute the deployment independently.

**Inferred user intent:** A self-contained document that serves as both explanation and step-by-step runbook.

### What I did

- Created docmgr ticket `sanitize-008`
- Created design doc: `Deploy Sanitize Web UI to K3s - Implementation Guide`
- Created this diary document
- Wrote the complete design doc covering:
  - Executive summary
  - Background concepts (the application, CGO, the architecture, the three repos, the reference pattern)
  - Full Dockerfile with CGO-aware build
  - GitHub Actions publish workflow
  - GHCR public package setup
  - All five Kustomize manifests (namespace, deployment, service, ingress, kustomization)
  - Argo CD Application manifest
  - Deployment validation commands
  - Design decisions with alternatives considered
  - Step-by-step implementation plan
  - Risks and mitigations
  - Troubleshooting guide
  - References

### Why

The user asked for an intern-level guide. That means explaining not just "what to do" but "why each piece exists" and "what to check when it breaks."

### What worked

- Having the pretext pattern as a concrete template made the manifest design straightforward
- The existing gitops documentation covered most of the conceptual ground — the design doc could reference it rather than re-explain everything

### What didn't work

- N/A

### What I learned

- The sanitize deployment is genuinely simple from an infrastructure perspective. The only interesting wrinkle is the CGO dependency affecting the Dockerfile.

### What was tricky to build

Writing at the right level of detail for an intern audience without being condescending or omitting critical steps. The CGO/musl dependency chain and the GHCR visibility step are the two most likely failure points that need extra emphasis.

### What warrants a second pair of eyes

- The Dockerfile — specifically whether `build-base` in Alpine provides everything tree-sitter needs
- The resource limits (64Mi request, 128Mi limit) — should be validated after deployment with actual memory usage

### What should be done in the future

1. Actually implement the Dockerfile and publish workflow
2. Create the GitOps manifests
3. Deploy and validate
4. Consider adding CI-created GitOps PRs (Phase 6 in the design doc)
5. Consider adding OCI metadata labels to the Dockerfile

### Code review instructions

- Read the design doc from top to bottom
- Verify each manifest against the pretext reference
- Check the Dockerfile build stages for correctness
- Confirm the workflow YAML is valid GitHub Actions syntax

### Technical details

Files created in this step:
- `ttmp/2026/03/29/sanitize-008--deploy-sanitize-web-ui-to-k3s-via-ghcr-and-argocd/design-doc/01-deploy-sanitize-web-ui-to-k3s-implementation-guide.md`
- `ttmp/2026/03/29/sanitize-008--deploy-sanitize-web-ui-to-k3s-via-ghcr-and-argocd/reference/01-investigation-diary.md`
