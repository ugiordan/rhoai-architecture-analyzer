# odh-dashboard: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

### Secrets Referenced

| Name | Type | Referenced By |
|------|------|---------------|
| dashboard-proxy-tls | kubernetes.io/tls | deployment/odh-dashboard, service/odh-dashboard |
| webhook-server-cert | Opaque | deployment/workspaces-controller |

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| odh-dashboard | odh-dashboard | ? | ? | ? | [`manifests/core-bases/base/deployment.yaml`](https://github.com/opendatahub-io/odh-dashboard/blob/37ad44c6f0e918c8bf8994312c0b99aa2e403a0c/manifests/core-bases/base/deployment.yaml) |
| odh-dashboard | kube-rbac-proxy | ? | ? | ? | [`manifests/core-bases/base/deployment.yaml`](https://github.com/opendatahub-io/odh-dashboard/blob/37ad44c6f0e918c8bf8994312c0b99aa2e403a0c/manifests/core-bases/base/deployment.yaml) |
| workspaces-backend | workspaces-backend | ? | ? | ? | [`packages/notebooks/upstream/workspaces/backend/manifests/kustomize/base/deployment.yaml`](https://github.com/opendatahub-io/odh-dashboard/blob/37ad44c6f0e918c8bf8994312c0b99aa2e403a0c/packages/notebooks/upstream/workspaces/backend/manifests/kustomize/base/deployment.yaml) |
| workspaces-controller | manager | ? | ? | ? | [`packages/notebooks/upstream/workspaces/controller/manifests/kustomize/components/certmanager/deployment_patch.yaml`](https://github.com/opendatahub-io/odh-dashboard/blob/37ad44c6f0e918c8bf8994312c0b99aa2e403a0c/packages/notebooks/upstream/workspaces/controller/manifests/kustomize/components/certmanager/deployment_patch.yaml) |
| workspaces-controller | manager | ? | ? | ? | [`packages/notebooks/upstream/workspaces/controller/manifests/kustomize/components/istio/deployment_patch.yaml`](https://github.com/opendatahub-io/odh-dashboard/blob/37ad44c6f0e918c8bf8994312c0b99aa2e403a0c/packages/notebooks/upstream/workspaces/controller/manifests/kustomize/components/istio/deployment_patch.yaml) |
| workspaces-controller | manager | ? | ? | ? | [`packages/notebooks/upstream/workspaces/controller/manifests/kustomize/components/prometheus/deployment_patch.yaml`](https://github.com/opendatahub-io/odh-dashboard/blob/37ad44c6f0e918c8bf8994312c0b99aa2e403a0c/packages/notebooks/upstream/workspaces/controller/manifests/kustomize/components/prometheus/deployment_patch.yaml) |
| workspaces-controller | manager | true | ? | ? | [`packages/notebooks/upstream/workspaces/controller/manifests/kustomize/base/manager/manager.yaml`](https://github.com/opendatahub-io/odh-dashboard/blob/37ad44c6f0e918c8bf8994312c0b99aa2e403a0c/packages/notebooks/upstream/workspaces/controller/manifests/kustomize/base/manager/manager.yaml) |
| workspaces-frontend | workspaces-frontend | ? | ? | ? | [`packages/notebooks/upstream/workspaces/frontend/manifests/kustomize/base/deployment.yaml`](https://github.com/opendatahub-io/odh-dashboard/blob/37ad44c6f0e918c8bf8994312c0b99aa2e403a0c/packages/notebooks/upstream/workspaces/frontend/manifests/kustomize/base/deployment.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `Dockerfile` | ${BASE_IMAGE} | 2 | 1001:0 |  |  |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: ${BASE_IMAGE} |
| `packages/automl/Dockerfile` | ${DISTROLESS_BASE_IMAGE} | 3 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${NODE_BASE_IMAGE}; Unpinned base image: ${GOLANG_BASE_IMAGE}; Unpinned base image: ${DISTROLESS_BASE_IMAGE} |
| `packages/automl/Dockerfile.workspace` | ${DISTROLESS_BASE_IMAGE} | 3 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${NODE_BASE_IMAGE}; Unpinned base image: ${GOLANG_BASE_IMAGE}; Unpinned base image: ${DISTROLESS_BASE_IMAGE} |
| `packages/autorag/Dockerfile` | ${DISTROLESS_BASE_IMAGE} | 3 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${NODE_BASE_IMAGE}; Unpinned base image: ${GOLANG_BASE_IMAGE}; Unpinned base image: ${DISTROLESS_BASE_IMAGE} |
| `packages/autorag/Dockerfile.workspace` | ${DISTROLESS_BASE_IMAGE} | 3 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${NODE_BASE_IMAGE}; Unpinned base image: ${GOLANG_BASE_IMAGE}; Unpinned base image: ${DISTROLESS_BASE_IMAGE} |
| `packages/eval-hub/Dockerfile` | ${DISTROLESS_BASE_IMAGE} | 3 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${NODE_BASE_IMAGE}; Unpinned base image: ${GOLANG_BASE_IMAGE}; Unpinned base image: ${DISTROLESS_BASE_IMAGE} |
| `packages/eval-hub/Dockerfile.workspace` | ${DISTROLESS_BASE_IMAGE} | 3 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${NODE_BASE_IMAGE}; Unpinned base image: ${GOLANG_BASE_IMAGE}; Unpinned base image: ${DISTROLESS_BASE_IMAGE} |
| `packages/gen-ai/Dockerfile` | registry.access.redhat.com/ubi9/ubi-minimal:9.3 | 3 | 1001 |  |  |  | Unpinned base image: registry.access.redhat.com/ubi9/nodejs-20 |
| `packages/gen-ai/Dockerfile.workspace` | ${DISTROLESS_BASE_IMAGE} | 3 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${NODE_BASE_IMAGE}; Unpinned base image: ${GOLANG_BASE_IMAGE}; Unpinned base image: ${DISTROLESS_BASE_IMAGE} |
| `packages/gen-ai/bff/Dockerfile` | registry.access.redhat.com/ubi9/ubi-minimal:9.3 | 2 | 1001 |  |  |  |  |
| `packages/maas/Dockerfile` | ${DISTROLESS_BASE_IMAGE} | 3 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${NODE_BASE_IMAGE}; Unpinned base image: ${GOLANG_BASE_IMAGE}; Unpinned base image: ${DISTROLESS_BASE_IMAGE} |
| `packages/maas/Dockerfile.workspace` | ${DISTROLESS_BASE_IMAGE} | 3 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${NODE_BASE_IMAGE}; Unpinned base image: ${GOLANG_BASE_IMAGE}; Unpinned base image: ${DISTROLESS_BASE_IMAGE} |
| `packages/mlflow/Dockerfile` | ${DISTROLESS_BASE_IMAGE} | 3 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${NODE_BASE_IMAGE}; Unpinned base image: ${GOLANG_BASE_IMAGE}; Unpinned base image: ${DISTROLESS_BASE_IMAGE} |
| `packages/mlflow/Dockerfile.workspace` | ${DISTROLESS_BASE_IMAGE} | 3 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${NODE_BASE_IMAGE}; Unpinned base image: ${GOLANG_BASE_IMAGE}; Unpinned base image: ${DISTROLESS_BASE_IMAGE} |
| `packages/model-registry/Dockerfile.workspace` | ${DISTROLESS_BASE_IMAGE} | 3 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${NODE_BASE_IMAGE}; Unpinned base image: ${GOLANG_BASE_IMAGE}; Unpinned base image: ${DISTROLESS_BASE_IMAGE} |
| `packages/model-registry/upstream/Dockerfile` | ${DISTROLESS_BASE_IMAGE} | 3 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${NODE_BASE_IMAGE}; Unpinned base image: ${GOLANG_BASE_IMAGE}; Unpinned base image: ${DISTROLESS_BASE_IMAGE} |
| `packages/model-registry/upstream/Dockerfile.standalone` | release | 4 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${NODE_BASE_IMAGE}; Unpinned base image: ${GOLANG_BASE_IMAGE}; Unpinned base image: ${DISTROLESS_BASE_IMAGE}; Unpinned base image: release |
| `packages/notebooks/Dockerfile.workspace` | ${DISTROLESS_BASE_IMAGE} | 3 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${NODE_BASE_IMAGE}; Unpinned base image: ${GOLANG_BASE_IMAGE}; Unpinned base image: ${DISTROLESS_BASE_IMAGE} |
| `packages/notebooks/upstream/workspaces/backend/Dockerfile` | gcr.io/distroless/static:nonroot | 2 | 65532:65532 |  | multi-arch |  |  |
| `packages/notebooks/upstream/workspaces/controller/Dockerfile` | gcr.io/distroless/static:nonroot | 2 | 65532:65532 |  | multi-arch |  |  |
| `packages/notebooks/upstream/workspaces/frontend/Dockerfile` | nginx:alpine | 2 | 101:101 |  |  |  |  |
| `packages/notebooks/upstream/workspaces/frontend/Dockerfile.dev` | node:20-slim | 1 | 1000:1000 |  |  |  |  |
| `packages/plugin-template/Dockerfile.workspace` | ${DISTROLESS_BASE_IMAGE} | 3 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${NODE_BASE_IMAGE}; Unpinned base image: ${GOLANG_BASE_IMAGE}; Unpinned base image: ${DISTROLESS_BASE_IMAGE} |
| `scripts/ci/Dockerfile` | quay.io/fedora/fedora:43 | 1 | $USER |  |  |  |  |

