# trainer: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

### Secrets Referenced

| Name | Type | Referenced By |
|------|------|---------------|
| kubeflow-trainer-webhook-cert | Opaque | deployment/kubeflow-trainer-controller-manager |
| webhook-server-cert | Opaque | deployment/controller-manager |

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| controller-manager | manager | ? | ? | ? | [`.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.10.1/config/default/manager_config_patch.yaml`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.10.1/config/default/manager_config_patch.yaml) |
| controller-manager | kube-rbac-proxy | ? | ? | ? | [`.gopath-loader/pkg/mod/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/config/default/manager_auth_proxy_patch.yaml`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/config/default/manager_auth_proxy_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gopath-loader/pkg/mod/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/config/default/manager_auth_proxy_patch.yaml`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/config/default/manager_auth_proxy_patch.yaml) |
| controller-manager | manager | ? | true | ? | [`.gomod-cache/sigs.k8s.io/jobset@v0.10.1/config/components/manager/manager.yaml`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/sigs.k8s.io/jobset@v0.10.1/config/components/manager/manager.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gomod-cache/sigs.k8s.io/jobset@v0.10.1/config/default/manager_config_patch.yaml`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/sigs.k8s.io/jobset@v0.10.1/config/default/manager_config_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gomod-cache/sigs.k8s.io/jobset@v0.10.1/config/default/manager_metrics_patch.yaml`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/sigs.k8s.io/jobset@v0.10.1/config/default/manager_metrics_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gomod-cache/sigs.k8s.io/jobset@v0.10.1/config/default/manager_webhook_patch.yaml`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/sigs.k8s.io/jobset@v0.10.1/config/default/manager_webhook_patch.yaml) |
| controller-manager | kube-rbac-proxy | ? | ? | ? | [`.gomod-cache/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/config/default/manager_auth_proxy_patch.yaml`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/config/default/manager_auth_proxy_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gomod-cache/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/config/default/manager_auth_proxy_patch.yaml`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/config/default/manager_auth_proxy_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gomod-cache/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/config/default/manager_config_patch.yaml`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/config/default/manager_config_patch.yaml) |
| controller-manager | manager | ? | true | ? | [`.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.10.1/config/components/manager/manager.yaml`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.10.1/config/components/manager/manager.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gopath-loader/pkg/mod/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/config/manager/manager.yaml`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/config/manager/manager.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.10.1/config/default/manager_webhook_patch.yaml`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.10.1/config/default/manager_webhook_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.10.1/config/default/manager_metrics_patch.yaml`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.10.1/config/default/manager_metrics_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gopath-loader/pkg/mod/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/config/default/manager_config_patch.yaml`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/config/default/manager_config_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gomod-cache/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/config/manager/manager.yaml`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/config/manager/manager.yaml) |
| kubeflow-trainer-controller-manager | manager | true | ? | ? | [`manifests/base/manager/manager.yaml`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/manifests/base/manager/manager.yaml) |
| kubeflow-trainer-controller-manager | manager | ? | ? | ? | [`manifests/base/manager/manager_config_patch.yaml`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/manifests/base/manager/manager_config_patch.yaml) |
| kubeflow-trainer-controller-manager | manager | ? | ? | ? | [`manifests/rhoai/manager_config_patch.yaml`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/manifests/rhoai/manager_config_patch.yaml) |
| kubeflow-trainer-controller-manager | manager | ? | ? | ? | [`manifests/rhoai/manager_metrics_patch.yaml`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/manifests/rhoai/manager_metrics_patch.yaml) |
| peaks | peaks | ? | ? | ? | [`.gomod-cache/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/pkg/trimaran/peaks/deployment/deployment.yaml`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/pkg/trimaran/peaks/deployment/deployment.yaml) |
| peaks | peaks | ? | ? | ? | [`.gopath-loader/pkg/mod/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/pkg/trimaran/peaks/deployment/deployment.yaml`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/pkg/trimaran/peaks/deployment/deployment.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `.gomod-cache/github.com/pelletier/go-toml@v1.9.5/Dockerfile` | scratch | 2 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/prometheus/client_golang@v1.23.2/Dockerfile` | quay.io/prometheus/busybox:latest | 2 |  |  |  |  | Unpinned base image: quay.io/prometheus/busybox:latest; No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/x/net@v0.43.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/x/sys@v0.35.0/unix/linux/Dockerfile` | ubuntu:24.10 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/k8s.io/apiextensions-apiserver@v0.34.1/artifacts/simple-image/Dockerfile` | gcr.io/distroless/base-debian10:latest | 1 |  |  |  |  | Unpinned base image: gcr.io/distroless/base-debian10:latest; No USER directive found (defaults to root) |
| `.gomod-cache/sigs.k8s.io/jobset@v0.10.1/Dockerfile` | ${BASE_IMAGE} | 2 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${BUILDER_IMAGE}; Unpinned base image: ${BASE_IMAGE} |
| `.gomod-cache/sigs.k8s.io/jobset@v0.10.1/sdk/python/Dockerfile` | python:3.13 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/sigs.k8s.io/kind@v0.30.0/images/base/Dockerfile` | scratch | 9 |  |  | multi-arch |  | Unpinned base image: $BASE_IMAGE; Unpinned base image: $BASE_IMAGE; Unpinned base image: go-build; Unpinned base image: go-build; Unpinned base image: go-build; Unpinned base image: go-build; Unpinned base image: go-build; Unpinned base image: base; Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gomod-cache/sigs.k8s.io/kind@v0.30.0/images/haproxy/Dockerfile` | "gcr.io/distroless/static-debian11" | 2 |  |  |  |  | Unpinned base image: ${BASE}; Unpinned base image: "gcr.io/distroless/static-debian11"; No USER directive found (defaults to root) |
| `.gomod-cache/sigs.k8s.io/kind@v0.30.0/images/local-path-helper/Dockerfile` | "gcr.io/distroless/static-debian11" | 2 |  |  |  |  | Unpinned base image: ${BASE}; Unpinned base image: "gcr.io/distroless/static-debian11"; No USER directive found (defaults to root) |
| `.gomod-cache/sigs.k8s.io/kind@v0.30.0/images/local-path-provisioner/Dockerfile` | gcr.io/distroless/base-debian11 | 2 |  |  | multi-arch |  | Unpinned base image: docker.io/library/golang:latest; Unpinned base image: gcr.io/distroless/base-debian11; No USER directive found (defaults to root) |
| `.gomod-cache/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/build/controller/Dockerfile` | $DISTROLESS_BASE_IMAGE | 2 | 65532:65532 |  | multi-arch |  | Unpinned base image: $GO_BASE_IMAGE; Unpinned base image: $DISTROLESS_BASE_IMAGE |
| `.gomod-cache/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/build/scheduler/Dockerfile` | $DISTROLESS_BASE_IMAGE | 2 | 65532:65532 |  | multi-arch |  | Unpinned base image: $GO_BASE_IMAGE; Unpinned base image: $DISTROLESS_BASE_IMAGE |
| `.gopath-loader/pkg/mod/github.com/pelletier/go-toml@v1.9.5/Dockerfile` | scratch | 2 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/prometheus/client_golang@v1.23.2/Dockerfile` | quay.io/prometheus/busybox:latest | 2 |  |  |  |  | Unpinned base image: quay.io/prometheus/busybox:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/net@v0.43.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/sys@v0.35.0/unix/linux/Dockerfile` | ubuntu:24.10 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/k8s.io/apiextensions-apiserver@v0.34.1/artifacts/simple-image/Dockerfile` | gcr.io/distroless/base-debian10:latest | 1 |  |  |  |  | Unpinned base image: gcr.io/distroless/base-debian10:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.10.1/Dockerfile` | ${BASE_IMAGE} | 2 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${BUILDER_IMAGE}; Unpinned base image: ${BASE_IMAGE} |
| `.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.10.1/sdk/python/Dockerfile` | python:3.13 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/sigs.k8s.io/kind@v0.30.0/images/base/Dockerfile` | scratch | 9 |  |  | multi-arch |  | Unpinned base image: $BASE_IMAGE; Unpinned base image: $BASE_IMAGE; Unpinned base image: go-build; Unpinned base image: go-build; Unpinned base image: go-build; Unpinned base image: go-build; Unpinned base image: go-build; Unpinned base image: base; Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/sigs.k8s.io/kind@v0.30.0/images/haproxy/Dockerfile` | "gcr.io/distroless/static-debian11" | 2 |  |  |  |  | Unpinned base image: ${BASE}; Unpinned base image: "gcr.io/distroless/static-debian11"; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/sigs.k8s.io/kind@v0.30.0/images/local-path-helper/Dockerfile` | "gcr.io/distroless/static-debian11" | 2 |  |  |  |  | Unpinned base image: ${BASE}; Unpinned base image: "gcr.io/distroless/static-debian11"; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/sigs.k8s.io/kind@v0.30.0/images/local-path-provisioner/Dockerfile` | gcr.io/distroless/base-debian11 | 2 |  |  | multi-arch |  | Unpinned base image: docker.io/library/golang:latest; Unpinned base image: gcr.io/distroless/base-debian11; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/build/controller/Dockerfile` | $DISTROLESS_BASE_IMAGE | 2 | 65532:65532 |  | multi-arch |  | Unpinned base image: $GO_BASE_IMAGE; Unpinned base image: $DISTROLESS_BASE_IMAGE |
| `.gopath-loader/pkg/mod/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/build/scheduler/Dockerfile` | $DISTROLESS_BASE_IMAGE | 2 | 65532:65532 |  | multi-arch |  | Unpinned base image: $GO_BASE_IMAGE; Unpinned base image: $DISTROLESS_BASE_IMAGE |
| `cmd/data_cache/Dockerfile` | debian:bookworm-slim | 2 | cache_user |  |  |  |  |
| `cmd/initializers/dataset/Dockerfile` | python:3.11-slim-bookworm | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `cmd/initializers/model/Dockerfile` | python:3.11-slim-bookworm | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `cmd/runtimes/deepspeed/Dockerfile` | nvidia/cuda:12.8.1-devel-ubuntu22.04 | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `cmd/runtimes/mlx/Dockerfile` | nvidia/cuda:12.8.1-devel-ubuntu22.04 | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `cmd/trainer-controller-manager/Dockerfile` | gcr.io/distroless/static:nonroot | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `cmd/trainer-controller-manager/Dockerfile.odh` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 2 | 65532:65532 |  |  |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |
| `cmd/trainers/torchtune/Dockerfile` | pytorch/pytorch:2.7.1-cuda12.8-cudnn9-runtime | 1 |  |  |  |  | No USER directive found (defaults to root) |

