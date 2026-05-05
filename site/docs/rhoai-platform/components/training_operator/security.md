# training-operator: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

### Secrets Referenced

| Name | Type | Referenced By |
|------|------|---------------|
| kubeflow-training-operator-webhook-cert | Opaque | deployment/training-operator |
| training-operator-webhook-cert | Opaque | deployment/training-operator |

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| training-operator | training-operator | ? | ? | ? | [`manifests/base/deployment.yaml`](https://github.com/kubeflow/training-operator/blob/8582a4b2a238e3552c6b726764580295303a3414/manifests/base/deployment.yaml) |
| training-operator | training-operator | ? | ? | ? | [`manifests/rhoai/manager_config_patch.yaml`](https://github.com/kubeflow/training-operator/blob/8582a4b2a238e3552c6b726764580295303a3414/manifests/rhoai/manager_config_patch.yaml) |
| training-operator | training-operator | ? | ? | ? | [`manifests/rhoai/manager_metrics_patch.yaml`](https://github.com/kubeflow/training-operator/blob/8582a4b2a238e3552c6b726764580295303a3414/manifests/rhoai/manager_metrics_patch.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `build/images/kubectl-delivery/Dockerfile` | alpine:3.17 | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `build/images/training-operator/Dockerfile` | gcr.io/distroless/static:latest | 2 |  |  |  |  | Unpinned base image: gcr.io/distroless/static:latest; No USER directive found (defaults to root) |
| `build/images/training-operator/Dockerfile.konflux` | registry.access.redhat.com/ubi9/ubi-minimal@sha256:8d0a8fb39ec907e8ca62cdd24b62a63ca49a30fe465798a360741fde58437a23 | 2 | ${USER}:${USER} |  |  |  |  |
| `build/images/training-operator/Dockerfile.multiarch` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 1 | 65532:65532 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |
| `build/images/training-operator/Dockerfile.rhoai` | registry.access.redhat.com/ubi9/ubi:latest | 2 | 65532:65532 |  |  |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi:latest |
| `sdk/python/Dockerfile.conformance` | python:3.10-alpine | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `sdk/python/kubeflow/storage_initializer/Dockerfile` | python:3.11 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `sdk/python/kubeflow/trainer/Dockerfile` | nvcr.io/nvidia/pytorch:24.06-py3 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `sdk/python/kubeflow/trainer/Dockerfile.cpu` | python:3.11 | 1 |  |  |  |  | No USER directive found (defaults to root) |

