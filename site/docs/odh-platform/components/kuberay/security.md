# kuberay: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

### Secrets Referenced

| Name | Type | Referenced By |
|------|------|---------------|
| webhook-server-cert | Opaque | deployment/kuberay-operator |

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| kuberay-operator | kuberay-operator | ? | ? | ? | [`ray-operator/config/overlays/test-overrides/deployment-override.yaml`](https://github.com/ray-project/kuberay/blob/acbf7e027447a2ca3057213fc4ebba83ac1547c7/ray-operator/config/overlays/test-overrides/deployment-override.yaml) |
| kuberay-operator | kuberay-operator | ? | ? | ? | [`ray-operator/config/default-with-webhooks/manager_webhook_patch.yaml`](https://github.com/ray-project/kuberay/blob/acbf7e027447a2ca3057213fc4ebba83ac1547c7/ray-operator/config/default-with-webhooks/manager_webhook_patch.yaml) |
| kuberay-operator | kuberay-operator | ? | true | ? | [`ray-operator/config/manager/manager.yaml`](https://github.com/ray-project/kuberay/blob/acbf7e027447a2ca3057213fc4ebba83ac1547c7/ray-operator/config/manager/manager.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `apiserver/Dockerfile` | scratch | 2 | 65532:65532 |  |  |  | Unpinned base image: scratch |
| `apiserver/Dockerfile.buildx` | scratch | 1 | 65532:65532 |  | multi-arch |  | Unpinned base image: scratch |
| `benchmark/perf-tests/images/ray-pytorch/Dockerfile` | rayproject/ray:2.46.0 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `dashboard/Dockerfile` | base | 4 | nextjs |  |  |  | Unpinned base image: base; Unpinned base image: base; Unpinned base image: base |
| `experimental/Dockerfile` | scratch | 2 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `experimental/Dockerfile.buildx` | scratch | 1 | 65532:65532 |  | multi-arch |  | Unpinned base image: scratch |
| `historyserver/Dockerfile.collector` | ubuntu:22.04 | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `historyserver/Dockerfile.historyserver` | ubuntu:22.04 | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `historyserver/cmd/collector/Dockerfile` |  | 0 |  |  |  |  | No USER directive found (defaults to root) |
| `historyserver/cmd/historyserver/Dockerfile` |  | 0 |  |  |  |  | No USER directive found (defaults to root) |
| `proto/Dockerfile` | golang:1.26-bookworm | 1 | 65532:65532 |  |  |  |  |
| `ray-operator/Dockerfile` | gcr.io/distroless/base-debian12:nonroot | 3 | 65532:65532 |  |  |  | Unpinned base image: scratch |
| `ray-operator/Dockerfile.buildx` | gcr.io/distroless/base-debian12:nonroot | 1 | 65532:65532 |  | multi-arch |  |  |
| `ray-operator/Dockerfile.submitter.buildx` | scratch | 1 | 65532:65532 |  | multi-arch |  | Unpinned base image: scratch |

