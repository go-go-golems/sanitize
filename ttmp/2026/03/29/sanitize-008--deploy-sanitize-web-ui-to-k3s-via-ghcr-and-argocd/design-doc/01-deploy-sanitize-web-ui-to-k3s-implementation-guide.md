---
Title: Deploy Sanitize Web UI to K3s - Implementation Guide
Ticket: sanitize-008
Status: active
Topics:
    - deployment
    - kubernetes
    - ghcr
    - argocd
    - gitops
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ../../../../../../../../2026-03-27--hetzner-k3s/docs/app-packaging-and-gitops-pr-standard.md:Standard packaging reference
    - Path: ../../../../../../../../2026-03-27--hetzner-k3s/docs/public-repo-ghcr-argocd-deployment-playbook.md:Deployment playbook reference
    - Path: ../../../../../../../../2026-03-27--hetzner-k3s/docs/source-app-deployment-infrastructure-playbook.md:Full deployment infrastructure playbook
    - Path: ../../../../../../../../2026-03-27--hetzner-k3s/gitops/applications/pretext.yaml:Reference ArgoCD Application for static web app
    - Path: ../../../../../../../../2026-03-27--hetzner-k3s/gitops/kustomize/pretext/deployment.yaml:Reference deployment pattern for static web app
    - Path: .github/workflows/push.yml:Existing CI workflow
    - Path: Makefile:Build targets
    - Path: go.mod:Go module definition with CGO dependency (tree-sitter)
    - Path: internal/cli/commands.go:CLI commands including serve command
    - Path: internal/server/server.go
      Note: HTTP server with embedded static files and API endpoints
    - Path: internal/server/server.go:Go HTTP server with embedded static files and API endpoints
    - Path: internal/server/static/index.html:Web UI entry point (Sanitize Playground)
ExternalSources: []
Summary: Complete implementation guide for deploying the Sanitize web UI (experimentation playground) to the Hetzner K3s cluster via GitHub Actions, GHCR, and Argo CD.
LastUpdated: 2026-03-29T17:01:11.518986968-04:00
WhatFor: ""
WhenToUse: ""
---


# Deploy Sanitize Web UI to K3s - Implementation Guide

## Executive Summary

The sanitize repository contains a Go-based structured-text linter and fixer that includes a web-based experimentation UI called "Sanitize Playground." This UI is currently only accessible by running `sanitize serve` locally. The goal is to deploy it as a publicly accessible web application on the existing Hetzner K3s cluster at `sanitize.yolo.scapegoat.dev`.

The deployment follows the established pattern used by other apps in the cluster (pretext, coinvault, draft-review): build a container image in GitHub Actions, publish it to GitHub Container Registry (GHCR), pin it in the GitOps repository, and let Argo CD reconcile it into the cluster.

This is a **Category 1 deployment** (public stateless app) — the simplest category. The sanitize web UI has no database, no secrets, and no external service dependencies. It serves static HTML/CSS/JS plus three JSON API endpoints that perform in-memory text processing.

## Problem Statement

The Sanitize Playground is a useful tool for experimenting with YAML and JSON sanitization, but it is only accessible when someone runs the binary locally with `sanitize serve`. There is no persistent, shared instance that team members or external users can access.

The K3s cluster at Hetzner already runs similar web applications (pretext-explorer, coinvault tools) and has the infrastructure to support another stateless web app: Argo CD for deployment, Traefik for ingress, cert-manager for TLS, and a DNS wildcard for `*.yolo.scapegoat.dev`.

## What You Need to Know Before Starting

This section explains the background concepts. If you are new to this cluster, read this section carefully before touching any files.

### The Sanitize Application

Sanitize is a Go CLI tool that lives at `github.com/go-go-golems/sanitize`. It has several commands:

- `sanitize fix` — fix broken YAML/JSON from stdin or a file
- `sanitize lint` — lint without fixing
- `sanitize parse` — show the tree-sitter parse tree
- `sanitize rules` — list available lint/fix rules
- `sanitize serve` — run the web server

The `serve` command starts an HTTP server (default port 8080) that serves:

1. **Static files** — an embedded single-page application (HTML + CSS + JS) at `/`
2. **`/api/examples`** — returns built-in example YAML/JSON inputs
3. **`/api/sanitize`** — accepts `{format, input}` and returns sanitized output with parse tree, fixes, and lint issues
4. **`/api/parse`** — accepts `{format, input}` and returns just the parse tree and errors

The static files are embedded into the Go binary using `//go:embed static` in `internal/server/server.go`. This means the final container image is a single binary with everything baked in — no external file mounts needed.

### Important: CGO Dependency

The sanitize binary depends on **tree-sitter** for parsing YAML and JSON. The Go bindings for tree-sitter (`github.com/tree-sitter/go-tree-sitter`) use **CGO**. This means:

- The build cannot use `CGO_ENABLED=0`
- The Dockerfile must either use a base image with a C compiler (like `golang:1.25-alpine` with `build-base`) or use a glibc-based image
- The runtime image must have compatible C libraries (cannot use `scratch` or `distroless/static`)

This is the single most important difference between sanitize and the demo app in the k3s repository's `app/` directory, which uses `CGO_ENABLED=0` and runs on bare `alpine`. You will need `alpine` with its standard C library, or a similar minimal image with libc.

### The Deployment Architecture

The deployment follows this chain:

```text
sanitize repo (github.com/go-go-golems/sanitize)
  -> GitHub Actions builds container image
    -> pushes to ghcr.io/go-go-golems/sanitize:sha-<commit>
      -> GitOps repo (wesen/2026-03-27--hetzner-k3s) pins image tag
        -> Argo CD Application watches gitops/kustomize/sanitize/
          -> Kubernetes Deployment runs the container
            -> Service exposes port 8080
              -> Ingress routes sanitize.yolo.scapegoat.dev -> Service
                -> cert-manager provisions TLS certificate
```

Each layer has a different owner and a different failure mode. Understanding this chain is essential for debugging.

### The Three Repositories Involved

1. **Sanitize app repo** (`github.com/go-go-golems/sanitize`) — owns source code, Dockerfile, CI workflows
2. **GitOps repo** (`wesen/2026-03-27--hetzner-k3s`) — owns Kubernetes manifests, Argo CD Applications
3. **The cluster itself** — runs Argo CD, Traefik, cert-manager, and your workload

You will need to make changes in repositories 1 and 2.

### Reference Pattern: pretext-explorer

The `pretext-explorer` deployment is the closest existing pattern to what we need. It is a stateless web application with:

- A Kustomize package at `gitops/kustomize/pretext/`
- An Argo CD Application at `gitops/applications/pretext.yaml`
- Namespace, Deployment, Service, and Ingress resources
- TLS via cert-manager with `letsencrypt-prod`
- Ingress via Traefik at `pretext.yolo.scapegoat.dev`

Currently, pretext uses `imagePullPolicy: Never` with a locally-imported image (`pretext-explorer:hk3s-0012`). Our sanitize deployment will use the proper GHCR registry pattern from the start.

## Proposed Solution

### Phase 1: Add a Dockerfile to the Sanitize Repository

Create a multi-stage Dockerfile that:

1. Builds the Go binary with CGO enabled
2. Copies it into a minimal Alpine runtime image

```dockerfile
# Build stage
FROM golang:1.25-alpine AS build
RUN apk add --no-cache build-base
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -trimpath -ldflags='-s -w' -o /out/sanitize ./cmd/sanitize

# Runtime stage
FROM alpine:3.21
RUN apk add --no-cache ca-certificates
RUN adduser -D -H -u 10001 appuser
USER 10001
WORKDIR /
COPY --from=build /out/sanitize /sanitize
EXPOSE 8080
ENTRYPOINT ["/sanitize", "serve"]
```

Key decisions in this Dockerfile:

- **`build-base` in build stage**: provides gcc, musl-dev, and make — needed for CGO/tree-sitter
- **`CGO_ENABLED=1`**: required for tree-sitter native bindings
- **`alpine:3.21` runtime**: provides musl libc needed by the CGO binary. We cannot use `scratch` or `distroless/static`
- **`ca-certificates`**: not strictly needed now (no outbound HTTPS), but good hygiene for future use
- **Non-root user**: security best practice, runs as UID 10001
- **`ENTRYPOINT ["/sanitize", "serve"]`**: starts the web server by default on port 8080

### Phase 2: Add a GitHub Actions Publish Workflow

Create `.github/workflows/publish-image.yaml` that builds and pushes the container image to GHCR.

```yaml
name: publish-image

on:
  pull_request:
  push:
    branches:
      - main

permissions:
  contents: read
  packages: write

jobs:
  docker:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6

      - uses: actions/setup-go@v6
        with:
          go-version-file: go.mod
          cache: true

      - name: Run tests
        run: go test ./...

      - uses: docker/setup-buildx-action@v4

      - id: meta
        uses: docker/metadata-action@v6
        with:
          images: ghcr.io/go-go-golems/sanitize
          tags: |
            type=sha,prefix=sha-
            type=ref,event=branch
            type=raw,value=latest,enable={{is_default_branch}}

      - uses: docker/login-action@v4
        if: github.event_name == 'push' && github.ref == 'refs/heads/main'
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - uses: docker/build-push-action@v7
        with:
          context: .
          push: ${{ github.event_name == 'push' && github.ref == 'refs/heads/main' }}
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          platforms: linux/amd64
          cache-from: type=gha
          cache-to: type=gha,mode=max
```

Important notes:

- **`packages: write`** permission is required for GHCR push
- **`GITHUB_TOKEN`** is used — no additional secrets needed for public repos
- **Pull requests** build and test but do not push
- **Only `linux/amd64`** platform — the K3s node is x86_64
- **Tags emitted**: `sha-<7char>`, `main`, `latest` — deployment will pin `sha-<7char>`
- **Build cache**: GitHub Actions cache via `type=gha` speeds up rebuilds

### Phase 3: Make the GHCR Package Public

After the first successful push to GHCR, the package needs to be made publicly pullable. By default, GHCR packages from public repos may still require authentication to pull.

Steps:
1. Go to the package settings at `github.com/orgs/go-go-golems/packages/container/sanitize/settings` (or the user-level equivalent)
2. Under "Danger Zone", change visibility to "Public"
3. Verify with an anonymous pull: `docker pull ghcr.io/go-go-golems/sanitize:latest`

If the package is not public, the K3s cluster will get `ImagePullBackOff` errors because there is no image pull secret configured for the sanitize namespace.

### Phase 4: Create GitOps Manifests in the K3s Repository

Create a new Kustomize package at `gitops/kustomize/sanitize/` with these files:

#### `namespace.yaml`

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: sanitize
  annotations:
    argocd.argoproj.io/sync-wave: "0"
```

#### `deployment.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sanitize
  annotations:
    argocd.argoproj.io/sync-wave: "1"
  labels:
    app.kubernetes.io/name: sanitize
    app.kubernetes.io/component: web
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: sanitize
      app.kubernetes.io/component: web
  template:
    metadata:
      labels:
        app.kubernetes.io/name: sanitize
        app.kubernetes.io/component: web
    spec:
      enableServiceLinks: false
      containers:
        - name: sanitize
          image: ghcr.io/go-go-golems/sanitize:sha-REPLACE_ME
          imagePullPolicy: IfNotPresent
          ports:
            - containerPort: 8080
              name: http
          readinessProbe:
            httpGet:
              path: /
              port: http
            initialDelaySeconds: 3
            periodSeconds: 10
          livenessProbe:
            httpGet:
              path: /
              port: http
            initialDelaySeconds: 10
            periodSeconds: 20
          resources:
            requests:
              cpu: 25m
              memory: 64Mi
            limits:
              memory: 128Mi
```

Key differences from the pretext deployment:

- **`containerPort: 8080`** not 80 — sanitize serves on 8080
- **`imagePullPolicy: IfNotPresent`** not `Never` — we pull from GHCR, not local import
- **Image from GHCR** — `ghcr.io/go-go-golems/sanitize:sha-XXXXX`
- **Health checks on `/`** — the static file server returns 200 for the index page

#### `service.yaml`

```yaml
apiVersion: v1
kind: Service
metadata:
  name: sanitize
  annotations:
    argocd.argoproj.io/sync-wave: "1"
  labels:
    app.kubernetes.io/name: sanitize
    app.kubernetes.io/component: web
spec:
  selector:
    app.kubernetes.io/name: sanitize
    app.kubernetes.io/component: web
  ports:
    - name: http
      port: 80
      targetPort: http
```

Note: the Service listens on port 80 externally but routes to the container's port 8080 via the named port `http`.

#### `ingress.yaml`

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: sanitize
  annotations:
    argocd.argoproj.io/sync-wave: "2"
    cert-manager.io/cluster-issuer: letsencrypt-prod
  labels:
    app.kubernetes.io/name: sanitize
    app.kubernetes.io/component: web
spec:
  ingressClassName: traefik
  tls:
    - hosts:
        - sanitize.yolo.scapegoat.dev
      secretName: sanitize-tls
  rules:
    - host: sanitize.yolo.scapegoat.dev
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: sanitize
                port:
                  name: http
```

This gives us:
- TLS termination with a Let's Encrypt certificate
- Hostname-based routing via Traefik
- The app will be accessible at `https://sanitize.yolo.scapegoat.dev`

#### `kustomization.yaml`

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: sanitize

resources:
  - namespace.yaml
  - deployment.yaml
  - service.yaml
  - ingress.yaml
```

### Phase 5: Create the Argo CD Application

Create `gitops/applications/sanitize.yaml`:

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: sanitize
  namespace: argocd
  finalizers:
    - resources-finalizer.argocd.argoproj.io
spec:
  project: default
  destination:
    server: https://kubernetes.default.svc
    namespace: sanitize
  source:
    repoURL: https://github.com/wesen/2026-03-27--hetzner-k3s.git
    targetRevision: main
    path: gitops/kustomize/sanitize
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
      - CreateNamespace=true
      - ServerSideApply=true
```

This follows the exact pattern of the pretext and draft-review applications. The `automated` sync policy with `prune` and `selfHeal` means Argo CD will automatically deploy changes and clean up removed resources.

### Phase 6: Deploy and Validate

After pushing the GitOps changes:

```bash
cd /home/manuel/code/wesen/2026-03-27--hetzner-k3s
export KUBECONFIG=$PWD/kubeconfig-91.98.46.169.yaml

# Check Argo CD picked it up
kubectl -n argocd get application sanitize \
  -o jsonpath='{.status.sync.status}{" "}{.status.health.status}{"\n"}'

# Check the pod is running
kubectl -n sanitize get pods

# Check the image
kubectl -n sanitize get deploy sanitize \
  -o jsonpath='{.spec.template.spec.containers[0].image}{"\n"}'

# Check the ingress
kubectl -n sanitize get ingress

# Validate the app responds
curl -fsSL https://sanitize.yolo.scapegoat.dev/
curl -fsSL https://sanitize.yolo.scapegoat.dev/api/examples | head -c 200
```

## Design Decisions

### Decision 1: Single binary with embedded static files

**Choice**: Use the existing `//go:embed` approach where static files are baked into the Go binary.

**Why**: This is already how the server works. No need for nginx, no sidecar containers, no volume mounts. One container, one binary, one process. This is the simplest operational surface.

**Alternative rejected**: Serving static files from nginx with a separate API backend. This would add complexity (two containers, or a separate static build step) for no real benefit at this scale.

### Decision 2: Alpine runtime image (not scratch/distroless)

**Choice**: Use `alpine:3.21` as the runtime base.

**Why**: The tree-sitter CGO dependency links against musl libc. The runtime image must provide compatible shared libraries. Alpine provides musl libc, is small (~7MB), and is widely understood.

**Alternative rejected**: `distroless/base` (Debian-based, provides glibc) — would work but is less common in this cluster's existing patterns. `scratch` — would not work because the binary is dynamically linked.

### Decision 3: GHCR public package (no image pull secret)

**Choice**: Make the GHCR package public and pull without authentication.

**Why**: The sanitize repository is already public on GitHub. There is no reason to keep the container image private. A public package avoids the complexity of wiring Vault-backed image pull secrets (which is the Category 2 pattern used by draft-review).

**Alternative rejected**: Private package with Vault/VSO pull secret — unnecessary complexity for a public tool.

### Decision 4: Dedicated namespace `sanitize`

**Choice**: Deploy into its own `sanitize` namespace.

**Why**: Follows the pattern of other apps (pretext, coinvault, draft-review). Namespace isolation provides resource quotas, RBAC boundaries, and clean `kubectl` targeting.

### Decision 5: Health checks on `/` (the index page)

**Choice**: Use HTTP GET on `/` for both readiness and liveness probes.

**Why**: The static file server always returns the index page at `/`. No special health endpoint is needed. If the Go process is alive, `/` will respond with 200.

**Alternative considered**: Adding a `/healthz` endpoint. Not needed — the existing behavior is sufficient and adding a health endpoint would require code changes.

## Implementation Plan — Step by Step

This section is the ordered checklist for an intern to follow. Do each step in order. Do not skip ahead.

### Step 1: Create and Test the Dockerfile Locally

1. Create `Dockerfile` in the sanitize repo root (see Phase 1 above)
2. Build it locally:
   ```bash
   cd /home/manuel/code/wesen/corporate-headquarters/sanitize
   docker build -t sanitize:local .
   ```
3. Run it locally:
   ```bash
   docker run --rm -p 8080:8080 sanitize:local
   ```
4. Open `http://localhost:8080` in a browser — you should see the Sanitize Playground
5. Test the API:
   ```bash
   curl -s http://localhost:8080/api/examples | jq '. | length'
   curl -s -X POST http://localhost:8080/api/sanitize \
     -H 'Content-Type: application/json' \
     -d '{"format":"yaml","input":"foo: bar\n  baz: qux"}' | jq .
   ```

### Step 2: Add the Publish Workflow

1. Create `.github/workflows/publish-image.yaml` (see Phase 2 above)
2. Commit and push to a branch
3. Open a PR to verify the workflow runs tests and builds (but does not push)
4. Merge the PR
5. Verify the merged push triggers image publication to GHCR
6. Check the GHCR package page to confirm the image exists with a `sha-` tag

### Step 3: Make the Package Public

1. Navigate to the package settings on GitHub
2. Change visibility to Public
3. Verify: `docker pull ghcr.io/go-go-golems/sanitize:<sha-tag>`

### Step 4: Create GitOps Manifests

In the `wesen/2026-03-27--hetzner-k3s` repository:

1. Create directory `gitops/kustomize/sanitize/`
2. Create the five manifest files (see Phase 4 above)
3. Replace `sha-REPLACE_ME` with the actual SHA tag from GHCR
4. Create `gitops/applications/sanitize.yaml` (see Phase 5)
5. Commit and push

### Step 5: Validate Deployment

1. Wait for Argo CD to detect the change (or force refresh)
2. Check sync status
3. Check pod status
4. Verify the URL `https://sanitize.yolo.scapegoat.dev` is reachable
5. Test the UI in a browser
6. Test the API endpoints

### Step 6 (Future): Add CI-Created GitOps PRs

This is optional and follows the pattern described in `docs/app-packaging-and-gitops-pr-standard.md`:

1. Add `deploy/gitops-targets.json` to the sanitize repo
2. Add `scripts/open_gitops_pr.py`
3. Add a workflow job that opens a PR against the GitOps repo on successful image publish

This automates the image tag bump so you do not have to manually edit the deployment manifest.

## Risks and Mitigations

### Risk 1: CGO build fails on GitHub Actions runner

**Likelihood**: Low — Alpine's `build-base` package is well-tested.

**Mitigation**: The existing CI (`push.yml`) already runs `go test ./...` which exercises CGO. If tests pass, the Docker build should too.

### Risk 2: Image tag mismatch between GHCR and GitOps

**Likelihood**: Medium — this happened before with mysql-ide (see the public-repo playbook).

**Mitigation**: After the first image publishes, note the exact `sha-XXXXXXX` tag from GHCR. Use that exact string in the deployment manifest. Do not guess or compute it separately.

### Risk 3: GHCR package not publicly pullable

**Likelihood**: High if you forget Step 3.

**Mitigation**: Always verify with an anonymous `docker pull` before assuming the cluster can pull. The cluster does not have GHCR credentials for the sanitize namespace.

### Risk 4: DNS not configured for sanitize.yolo.scapegoat.dev

**Likelihood**: Low — `*.yolo.scapegoat.dev` should already be a wildcard DNS record pointing to the cluster IP.

**Mitigation**: Verify with `dig sanitize.yolo.scapegoat.dev` or `nslookup sanitize.yolo.scapegoat.dev`.

## Troubleshooting Guide

### Pod stuck in ImagePullBackOff

```bash
kubectl -n sanitize describe pod <pod-name>
# Look at Events section for the exact error
# Common causes:
# - Package not public
# - Tag does not exist
# - Registry path wrong
```

### Pod starts but CrashLoopBackOff

```bash
kubectl -n sanitize logs deploy/sanitize --tail=50
# Common causes:
# - Missing libc (wrong base image)
# - Port conflict
# - Binary crash on startup
```

### Ingress not working / 404

```bash
kubectl -n sanitize get ingress sanitize -o yaml
# Check: host matches, service name matches, port name matches
# Also check cert-manager:
kubectl get certificate -n sanitize
kubectl describe certificate sanitize-tls -n sanitize
```

### Argo CD shows OutOfSync

```bash
kubectl -n argocd get application sanitize -o yaml
# Check: source path correct, targetRevision correct
# Force sync if needed:
kubectl -n argocd annotate application sanitize \
  argocd.argoproj.io/refresh=hard --overwrite
```

## References

### Sanitize Repository Files

- `internal/server/server.go` — HTTP server with embedded static files, API handlers
- `internal/server/static/` — HTML, CSS, JS for the Playground UI
- `internal/cli/commands.go` — CLI command definitions including `serve`
- `cmd/sanitize/main.go` — binary entrypoint
- `go.mod` — module definition, shows tree-sitter dependencies
- `.github/workflows/push.yml` — existing CI (tests only, no image publish)
- `.github/workflows/release.yaml` — GoReleaser release workflow (binary releases, not container)
- `Makefile` — build targets

### K3s GitOps Repository Files

- `docs/public-repo-ghcr-argocd-deployment-playbook.md` — how to publish to GHCR and deploy via Argo CD
- `docs/app-packaging-and-gitops-pr-standard.md` — standard packaging model and CI-created GitOps PRs
- `docs/source-app-deployment-infrastructure-playbook.md` — full deployment infrastructure playbook
- `docs/argocd-app-setup.md` — how to create Argo CD Applications
- `gitops/kustomize/pretext/` — reference deployment (stateless web app, closest pattern)
- `gitops/applications/pretext.yaml` — reference Argo CD Application
- `app/Dockerfile` — demo app Dockerfile (simpler, no CGO)
