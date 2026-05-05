# trainer: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

### Secrets Referenced

| Name | Type | Referenced By |
|------|------|---------------|
| kubeflow-trainer-webhook-cert | Opaque | deployment/kubeflow-trainer-controller-manager |

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| kubeflow-trainer-controller-manager | manager | true | ? | ? | [`manifests/base/manager/manager.yaml`](https://github.com/kubeflow/trainer/blob/c1b8b98e7aa8ed01098deb38104c335fdb68a35c/manifests/base/manager/manager.yaml) |
| kubeflow-trainer-controller-manager | manager | ? | ? | ? | [`manifests/base/manager/manager_config_patch.yaml`](https://github.com/kubeflow/trainer/blob/c1b8b98e7aa8ed01098deb38104c335fdb68a35c/manifests/base/manager/manager_config_patch.yaml) |
| kubeflow-trainer-controller-manager | manager | ? | ? | ? | [`manifests/rhoai/manager_config_patch.yaml`](https://github.com/kubeflow/trainer/blob/c1b8b98e7aa8ed01098deb38104c335fdb68a35c/manifests/rhoai/manager_config_patch.yaml) |
| kubeflow-trainer-controller-manager | manager | ? | ? | ? | [`manifests/rhoai/manager_metrics_patch.yaml`](https://github.com/kubeflow/trainer/blob/c1b8b98e7aa8ed01098deb38104c335fdb68a35c/manifests/rhoai/manager_metrics_patch.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `cmd/data_cache/Dockerfile` | debian:bookworm-slim | 2 | cache_user |  |  |  |  |
| `cmd/initializers/dataset/Dockerfile` | python:3.11-slim-bookworm | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `cmd/initializers/model/Dockerfile` | python:3.11-slim-bookworm | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `cmd/runtimes/deepspeed/Dockerfile` | nvidia/cuda:12.8.1-devel-ubuntu22.04 | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `cmd/runtimes/mlx/Dockerfile` | nvidia/cuda:12.8.1-devel-ubuntu22.04 | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `cmd/trainer-controller-manager/Dockerfile` | gcr.io/distroless/static:nonroot | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `cmd/trainer-controller-manager/Dockerfile.odh` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 2 | 65532:65532 |  |  |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |
| `cmd/trainer-controller-manager/Dockerfile.rhoai.konflux` | registry.access.redhat.com/ubi9/ubi-minimal@sha256:61d5ad475048c2e655cd46d0a55dfeaec182cc3faa6348cb85989a7c9e196483 | 2 | 65532:65532 |  |  |  |  |
| `cmd/trainers/torchtune/Dockerfile` | pytorch/pytorch:2.7.1-cuda12.8-cudnn9-runtime | 1 |  |  |  |  | No USER directive found (defaults to root) |

