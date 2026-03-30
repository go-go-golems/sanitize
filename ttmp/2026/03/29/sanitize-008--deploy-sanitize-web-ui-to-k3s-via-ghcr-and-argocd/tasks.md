# Tasks

## Phase 1: Container Image (sanitize repo)

- [x] 1. Create Dockerfile with CGO-aware multi-stage build
- [x] 2. Test Dockerfile locally (docker build)
- [x] 3. Test container locally (docker run + verify UI + API)
- [x] 4. Add OCI metadata labels to Dockerfile
- [x] 5. Add .github/workflows/publish-image.yaml
- [x] 6. Commit Dockerfile and workflow to sanitize repo
- [x] 7. Push to main, verify GHCR image published with sha- tag (sha-cebe68d)
- [x] 8. Verify GHCR package pullable (docker pull succeeded anonymously)

## Phase 2: GitOps Manifests (k3s repo)

- [x] 9. Create gitops/kustomize/sanitize/namespace.yaml
- [x] 10. Create gitops/kustomize/sanitize/deployment.yaml (pinned ghcr.io/go-go-golems/sanitize:sha-cebe68d)
- [x] 11. Create gitops/kustomize/sanitize/service.yaml
- [x] 12. Create gitops/kustomize/sanitize/ingress.yaml (sanitize.yolo.scapegoat.dev)
- [x] 13. Create gitops/kustomize/sanitize/kustomization.yaml
- [x] 14. Create gitops/applications/sanitize.yaml (ArgoCD Application)
- [x] 15. Commit and push GitOps manifests to k3s repo

## Phase 3: Deploy and Validate

- [x] 16. Bootstrap Argo CD Application (kubectl apply -f gitops/applications/sanitize.yaml)
- [x] 17. Verify Argo CD syncs: Synced Healthy
- [x] 18. Verify pod running (sanitize-5c44dbdcfd-hrhd2 Running)
- [x] 19. Verify TLS certificate provisioned (sanitize-tls Ready)
- [x] 20. Validate UI accessible at https://sanitize.yolo.scapegoat.dev
- [x] 21. Validate /api/examples returns 61 examples
- [x] 22. Validate /api/sanitize processes YAML correctly

## Phase 4: CI GitOps PR Automation (future)

- [ ] 23. Add deploy/gitops-targets.json to sanitize repo
- [ ] 24. Add scripts/open_gitops_pr.py
- [ ] 25. Add GitOps PR job to publish-image workflow
