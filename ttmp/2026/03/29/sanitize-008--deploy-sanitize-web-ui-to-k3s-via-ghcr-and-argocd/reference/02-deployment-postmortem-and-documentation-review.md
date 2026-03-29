---
Title: Deployment Postmortem and Documentation Review
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
    - /home/manuel/code/wesen/2026-03-27--hetzner-k3s/docs/source-app-deployment-infrastructure-playbook.md:Primary deployment playbook - reviewed
    - /home/manuel/code/wesen/2026-03-27--hetzner-k3s/docs/public-repo-ghcr-argocd-deployment-playbook.md:GHCR deployment playbook - reviewed
    - /home/manuel/code/wesen/2026-03-27--hetzner-k3s/docs/app-packaging-and-gitops-pr-standard.md:Packaging standard - reviewed
    - /home/manuel/code/wesen/2026-03-27--hetzner-k3s/docs/argocd-app-setup.md:ArgoCD app setup - reviewed
    - /home/manuel/code/wesen/2026-03-27--hetzner-k3s/docs/operator-quickstart.md:Operator quickstart - reviewed
    - /home/manuel/code/wesen/2026-03-27--hetzner-k3s/docs/platform-baseline-overview.md:Platform baseline - reviewed
    - /home/manuel/code/wesen/2026-03-27--hetzner-k3s/docs/hetzner-k3s-server-setup.md:Server setup - reviewed
    - /home/manuel/code/wesen/2026-03-27--hetzner-k3s/docs/tailscale-k3s-admin-access-playbook.md:Tailscale access - reviewed
    - /home/manuel/code/wesen/2026-03-27--hetzner-k3s/README.md:Repository README - reviewed
ExternalSources: []
Summary: "Postmortem of the sanitize web UI deployment to k3s, with detailed documentation review covering gaps, confusion points, and improvement recommendations."
LastUpdated: 2026-03-29T19:30:00.000000000-04:00
WhatFor: ""
WhenToUse: ""
---

# Deployment Postmortem and Documentation Review

## Goal

This document captures the full postmortem of deploying the Sanitize Playground to the K3s cluster, with a particular focus on where the existing documentation helped, where it created confusion, what was missing, and what should be improved or removed.

The perspective is that of an AI assistant working through the deployment for the first time, using only the repository documentation and code as guidance. This closely mirrors the experience of a new intern or operator who does not have the repository author available to answer questions.

## Outcome Summary

The deployment succeeded end-to-end in a single session:

| Metric | Value |
|--------|-------|
| Time from ticket creation to live URL | ~30 minutes |
| Files created in sanitize repo | 2 (Dockerfile, publish-image.yaml) |
| Files created in k3s repo | 6 (5 Kustomize + 1 Application) |
| CI failures | 0 |
| Cluster errors | 0 |
| Documentation-caused blockers | 2 (kubeconfig path, Application bootstrap) |
| Rollbacks needed | 0 |

Live URL: `https://sanitize.yolo.scapegoat.dev`

## What Went Right

### 1. The deployment category system is excellent

The `source-app-deployment-infrastructure-playbook.md` and `app-packaging-and-gitops-pr-standard.md` define four clear deployment categories:

| Category | Example | Complexity |
|----------|---------|------------|
| 1: Stateless | pretext | 5 files, no secrets |
| 2: Vault/VSO | coinvault | ~10 files, Vault secrets |
| 3: Bootstrap | keycloak | ~15 files, bootstrap jobs |
| 4: Data service | postgres | ~8 files, statefulset |

Being able to immediately classify sanitize as Category 1 saved significant time. Without this taxonomy, I would have had to study every deployment pattern to figure out which one to follow.

### 2. The pretext reference pattern was a perfect template

Having a concrete working example of a Category 1 deployment (`gitops/kustomize/pretext/`) made the manifest creation almost mechanical. Every file had a clear analog. The sync waves, labels, and probe patterns were all directly reusable.

### 3. The GHCR deployment playbook is thorough and accurate

The `public-repo-ghcr-argocd-deployment-playbook.md` document is genuinely excellent for its intended audience. It explains not just the steps but the conceptual model behind each layer of the deployment chain. The sha-tag mismatch incident from mysql-ide is particularly valuable — it preemptively answered a question I would have had.

### 4. The GitHub Actions workflow patterns are well-documented

The publish-image workflow examples in the docs were directly usable. The `docker/metadata-action` tag format, the push conditions, and the `packages: write` permission were all clearly specified.

### 5. The README "Start Here" section is well-organized

The updated README has a decision-tree-style "Start Here" section that points operators to the right doc for their situation. This is genuinely helpful and reduces the "where do I even begin" problem.

## What Went Wrong

### Blocker 1: Kubeconfig Path — Stale References Everywhere

**Severity: High — caused a complete blocker until resolved**

**What happened:** After deploying the GitOps manifests, I attempted to validate the deployment using the kubeconfig at `kubeconfig-91.98.46.169.yaml` (the public IP). The kubectl commands timed out because public Kubernetes API access on port 6443 is disabled.

**Why it happened:** Multiple documents still reference the public IP kubeconfig path:

1. **`argocd-app-setup.md` Step 4** explicitly shows:
   ```bash
   export KUBECONFIG=$PWD/kubeconfig-91.98.46.169.yaml
   ```
   This is the first command in the "Apply the Application and Trigger a Refresh" section — the exact step I was trying to follow.

2. **`hetzner-k3s-server-setup.md`** mentions `scripts/get-kubeconfig.sh` (the public IP version) in the repo layout section.

3. **The README Quick Start** Step 6 shows:
   ```bash
   ./scripts/get-kubeconfig.sh <server-ip>
   export KUBECONFIG=$PWD/kubeconfig-<server-ip>.yaml
   ```

4. **The README Day-2 Operations** section shows:
   ```bash
   export KUBECONFIG=$PWD/kubeconfig-<server-ip>.yaml
   ```

5. **The README Vault bootstrap notes** section shows:
   ```bash
   export KUBECONFIG=$PWD/kubeconfig-91.98.46.169.yaml
   ```

**What fixed it:** The user pointed me to the `.envrc` file, which sets up the Tailscale kubeconfig environment. After running `scripts/get-kubeconfig-tailscale.sh`, everything worked.

**Impact:** This would completely block a new intern who follows the argocd-app-setup doc literally. They would see timeouts and not understand why, because the doc doesn't mention that the kubeconfig path has changed.

### Blocker 2: Argo CD Application Bootstrap Not Obvious

**Severity: Medium — caused confusion but was resolvable**

**What happened:** After pushing the Application manifest to git and waiting, `kubectl -n argocd get application sanitize` returned "NotFound". I initially expected Argo CD to auto-discover new Application manifests from git.

**Why it happened:** The distinction between "Application manifest exists in git" and "Application object exists in the cluster" is not prominently explained in the deployment flow docs. 

The `argocd-app-setup.md` does explain this in Step 4, but it's framed as a migration step ("Apply the Application and Trigger a Refresh"), not as a mandatory first-time bootstrap step. The `source-app-deployment-infrastructure-playbook.md` recently added a paragraph about this:

> There is one more boundary that matters during the first rollout of a brand-new app:
> Git repo contains gitops/applications/<app>.yaml != cluster already has Application/<app>

This is good, but it appears mid-document in a conceptual section. An operator following the step-by-step flow might not reach it in time.

**What fixed it:** The `kubectl apply -f gitops/applications/sanitize.yaml` command from the argocd-app-setup doc.

**Impact:** A new intern would push the Application YAML, wait for something to happen, and eventually ask for help. The fix is simple once you know it, but the gap between expectation and reality is confusing.

## Documentation Gaps

### Gap 1: No Single "Deploy a New App End-to-End" Checklist

There are three deployment docs that each cover part of the story:

1. `source-app-deployment-infrastructure-playbook.md` — the conceptual model and full system
2. `public-repo-ghcr-argocd-deployment-playbook.md` — GHCR-specific steps
3. `argocd-app-setup.md` — Argo CD Application creation

But there is no single checklist that says: "You have a public Go app, you want it on the cluster. Here are the 10 steps in order." An operator currently has to read all three docs, mentally merge them, and figure out the correct execution order.

**Recommendation:** Add a "New Public Stateless App Deployment Checklist" section to `source-app-deployment-infrastructure-playbook.md` or create a short new doc. Something like:

```markdown
## Quick Checklist: Deploy a New Public Stateless App

1. [ ] App has a working Dockerfile
2. [ ] App has a publish-image.yaml workflow
3. [ ] Image is published to GHCR with sha- tags
4. [ ] GHCR package is publicly pullable
5. [ ] gitops/kustomize/<app>/ created (namespace, deployment, service, ingress, kustomization)
6. [ ] gitops/applications/<app>.yaml created
7. [ ] GitOps changes pushed to main
8. [ ] Application bootstrapped: kubectl apply -f gitops/applications/<app>.yaml
9. [ ] Argo CD shows Synced Healthy
10. [ ] Public URL responds with TLS
```

### Gap 2: No CGO / Non-Trivial Dockerfile Guidance

The demo app Dockerfile in `app/Dockerfile` uses `CGO_ENABLED=0` and a bare Alpine runtime. All the deployment docs implicitly assume this simple build model. There is no guidance for apps that require CGO, native dependencies, or multi-stage builds with build tools.

For sanitize, the tree-sitter CGO dependency is the single most important Dockerfile design constraint. A new intern copying the demo Dockerfile would get a build failure and not understand why.

**Recommendation:** Add a "Non-trivial Dockerfile Patterns" section to the source-app playbook, or at minimum add a note like:

> If your app uses CGO (e.g., for native C library bindings), you cannot use `CGO_ENABLED=0`. You will need `build-base` or equivalent in your build stage, and your runtime image must provide compatible shared libraries (e.g., Alpine's musl libc, not `scratch` or `distroless/static`).

### Gap 3: Tailscale Kubeconfig Is the Only Working Path, but Not the Default in Examples

The `operator-quickstart.md` correctly identifies Tailscale as the normal operator path. But many other docs still show the public IP kubeconfig or generic `<server-ip>` patterns. This creates a split-brain problem where the quick reference and the detailed docs disagree.

**Recommendation:** Update all kubectl examples in these docs to use the Tailscale kubeconfig pattern:
- `argocd-app-setup.md` Step 4
- `hetzner-k3s-server-setup.md` repo layout and examples
- `README.md` Quick Start and Day-2 Operations
- `README.md` Vault bootstrap notes

The simplest fix is to replace:
```bash
export KUBECONFIG=$PWD/kubeconfig-91.98.46.169.yaml
```
with:
```bash
# See operator-quickstart.md for Tailscale setup
export KUBECONFIG=$PWD/.cache/kubeconfig-tailnet.yaml
```

### Gap 4: No Mention of `.envrc` in Any Doc

The `.envrc` file is the fastest way to get a working environment, but it is not mentioned in any documentation. The `operator-quickstart.md` manually walks through setting each environment variable, which is fine for understanding, but a one-liner like "if you use direnv, the `.envrc` in this repo sets up the right environment automatically" would save time.

**Recommendation:** Add a mention to `operator-quickstart.md`:

> If you use `direnv`, the `.envrc` in this repo sets up the Tailscale environment variables and provides `kcfg-refresh` and `kcfg-use` helper functions.

### Gap 5: The README Quick Start Still Describes a First-Boot Flow

The README Quick Start section (steps 1-10) is a Terraform + cloud-init first-boot guide. For someone who wants to deploy a new app to an already-running cluster, the Quick Start is the wrong starting point. The "Start Here" decision tree above the Quick Start is the right entry point, but the visual weight of the numbered Quick Start steps can draw attention away from it.

**Recommendation:** Rename "Quick start" to something like "Initial Cluster Bring-Up" and add a prominent note:

> If the cluster is already running and you want to deploy a new app, skip this section and go to [source-app-deployment-infrastructure-playbook.md](docs/source-app-deployment-infrastructure-playbook.md).

## What Is Confusing

### Confusion 1: Three Overlapping Deployment Docs

There are three docs that cover deploying an app:

| Doc | Focus | Length |
|-----|-------|--------|
| `source-app-deployment-infrastructure-playbook.md` | Full system model, all categories | ~800 lines |
| `public-repo-ghcr-argocd-deployment-playbook.md` | Public GHCR path specifically | ~500 lines |
| `app-packaging-and-gitops-pr-standard.md` | Packaging standards and CI-created PRs | ~400 lines |

These three docs have significant conceptual overlap. All three explain the deployment chain, the three control planes, and the tag strategy. The source-app playbook tries to be the "canonical" doc but also says "use the public-repo GHCR page when the repository and image are intentionally public." But the public-repo GHCR page also explains the full model from scratch.

A new reader doesn't know which one to read. If they read all three, they encounter the same concepts explained three times with slightly different framing.

**Recommendation:** Consider whether the public-repo GHCR doc can be folded into the source-app playbook as a section or shortened to just the GHCR-specific delta. Alternatively, make the source-app playbook the only entry point and have the other two docs explicitly say "this doc assumes you've read [source-app playbook] already."

### Confusion 2: `imagePullPolicy: Never` in Pretext (the Reference Pattern)

Pretext is the closest reference pattern for a Category 1 deployment, but it still uses `imagePullPolicy: Never` with a locally-imported image (`pretext-explorer:hk3s-0012`). This creates a confusing situation: the docs say "use GHCR and pin sha- tags," but the actual reference pattern in the cluster does the opposite.

An intern copying pretext's deployment.yaml verbatim would get `imagePullPolicy: Never`, which would fail for a GHCR image.

**Recommendation:** Either update pretext to use GHCR (making it a true reference pattern), or add a note to the deployment docs saying "pretext still uses the old local-import flow; adapt as shown below for GHCR."

### Confusion 3: The Demo App vs Real Apps

The `app/` directory in the k3s repo contains a simple Go demo app with a Dockerfile. This is referenced in docs as an example. But its Dockerfile uses `CGO_ENABLED=0` and a pattern that doesn't work for most real Go apps with native dependencies. The demo app creates a misleading mental model.

**Recommendation:** Add a comment to `app/Dockerfile`:
```dockerfile
# NOTE: This demo uses CGO_ENABLED=0. Real apps with native C dependencies
# (e.g., tree-sitter, SQLite) need CGO_ENABLED=1 and a build-base stage.
```

### Confusion 4: Which Kubeconfig Script to Use

There are two kubeconfig scripts:
- `scripts/get-kubeconfig.sh` — public IP path
- `scripts/get-kubeconfig-tailscale.sh` — Tailscale path

The public IP version is the one shown in the README and older docs. The Tailscale version is the one that actually works. The naming doesn't indicate which is current.

**Recommendation:** Either:
- Rename `get-kubeconfig.sh` to `get-kubeconfig-public-ip.sh` and add a deprecation notice
- Or add a comment at the top of `get-kubeconfig.sh` saying "Prefer get-kubeconfig-tailscale.sh — public K8s API may be blocked"

## What Could Be Removed or Consolidated

### 1. The Legacy Helm Chart References

Multiple docs carefully explain the `gitops/charts/demo-stack` legacy compatibility. This context is important for the person who might accidentally delete it, but it occupies significant space in docs that should focus on the current deployment model.

**Recommendation:** Create one short doc called `legacy-bootstrap-compatibility.md` that explains the demo-stack chart, why it exists, and when it can be removed. Then replace the inline explanations in other docs with a single sentence: "See [legacy-bootstrap-compatibility.md](./legacy-bootstrap-compatibility.md) for the demo-stack chart context."

### 2. Repeated Explanations of the Deployment Chain

The deployment chain diagram:
```text
source repo -> CI -> GHCR -> GitOps -> Argo CD -> Kubernetes
```

appears in at least four docs with slight variations. This repetition makes each doc feel self-contained, which is good, but it also makes the doc set feel bloated for someone reading multiple docs.

**Recommendation:** Keep the diagram in the source-app playbook (the canonical doc) and reference it from the others rather than reproducing it.

### 3. The mysql-ide SHA Tag Mismatch Story

The story of the sha-tag mismatch (full SHA vs 7-char short SHA) appears in both the public-repo GHCR playbook and indirectly in the source-app playbook. It's a great cautionary tale, but it could be told once and referenced.

**Recommendation:** Keep it in the public-repo GHCR playbook (where it's most detailed) and reference it from other docs.

## Specific Doc Fixes Needed

| Doc | Issue | Fix |
|-----|-------|-----|
| `argocd-app-setup.md` Step 4 | Uses `kubeconfig-91.98.46.169.yaml` | Change to Tailscale kubeconfig |
| `argocd-app-setup.md` Step 4 | Framed as migration, not bootstrap | Add a "First-time bootstrap" note |
| `README.md` Quick Start Step 6 | Uses `get-kubeconfig.sh` | Change to `get-kubeconfig-tailscale.sh` or note both |
| `README.md` Day-2 Operations | Uses `kubeconfig-<server-ip>.yaml` | Change to Tailscale kubeconfig |
| `README.md` Vault bootstrap | Hardcodes `kubeconfig-91.98.46.169.yaml` | Change to Tailscale kubeconfig |
| `README.md` Quick Start | Named "Quick start" | Rename to "Initial Cluster Bring-Up" |
| `hetzner-k3s-server-setup.md` | References `scripts/get-kubeconfig.sh` | Add note about Tailscale alternative |
| `source-app-deployment-infrastructure-playbook.md` | No CGO guidance | Add non-trivial Dockerfile note |
| `app/Dockerfile` | No CGO caveat | Add comment about CGO_ENABLED=0 limitation |

## What Was Missing That Would Have Saved the Most Time

If I had to pick the single change that would have saved the most time during this deployment, it would be:

**A consistent, working kubeconfig pattern across all docs.**

The Tailscale migration happened, the old public IP path was disabled, but the docs weren't fully updated. This created a situation where following the docs exactly led to failure. The fix is straightforward (update the KUBECONFIG lines in each doc), and it would prevent every new operator from hitting the same wall.

The second most impactful change would be:

**A single "New App Deployment Checklist" — 10 numbered steps, no prose, just commands.**

The existing docs are excellent for understanding. But when you're actually deploying, you need a reference card, not a textbook. The two complement each other.

## Conclusion

The documentation in the k3s repository is remarkably thorough and well-written for its age. The category system, the deployment chain model, the cautionary tales from real incidents — these are all genuinely valuable. The main issues are:

1. **Stale kubeconfig references** — the #1 operational trap
2. **No single end-to-end checklist** — forces mental merging of multiple docs
3. **Doc overlap** — three docs covering similar ground with different emphasis
4. **Reference pattern (pretext) doesn't match the documented pattern** — `imagePullPolicy: Never` vs GHCR

None of these prevented the deployment from succeeding, but they each added friction that a clean doc set would eliminate.
