# workload-variant-autoscaler: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

### Secrets Referenced

| Name | Type | Referenced By |
|------|------|---------------|
| epp-metrics-token | Opaque | deployment/controller-manager |
| hf-token | Opaque | deployment/llama-deployment |
| kedaorg-certs | Opaque | deployment/keda-metrics-apiserver, deployment/keda-operator |
| prometheus-client-cert | Opaque | deployment/controller-manager |
| webhook-server-cert | Opaque | deployment/controller-manager |

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| controller-manager | manager | ? | ? | ? | [`.gopath-loader/pkg/mod/sigs.k8s.io/lws@v0.8.0/config/default/manager_metrics_patch.yaml`](https://github.com/llm-d/workload-variant-autoscaler/blob/46a611076ea6e421be0dafe4f085f3ecc80fa09e/.gopath-loader/pkg/mod/sigs.k8s.io/lws@v0.8.0/config/default/manager_metrics_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/manager/manager.yaml`](https://github.com/llm-d/workload-variant-autoscaler/blob/46a611076ea6e421be0dafe4f085f3ecc80fa09e/config/manager/manager.yaml) |
| controller-manager | manager | ? | true | ? | [`.gopath-loader/pkg/mod/sigs.k8s.io/lws@v0.8.0/config/manager/manager.yaml`](https://github.com/llm-d/workload-variant-autoscaler/blob/46a611076ea6e421be0dafe4f085f3ecc80fa09e/.gopath-loader/pkg/mod/sigs.k8s.io/lws@v0.8.0/config/manager/manager.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gopath-loader/pkg/mod/sigs.k8s.io/lws@v0.8.0/config/default/manager_webhook_patch.yaml`](https://github.com/llm-d/workload-variant-autoscaler/blob/46a611076ea6e421be0dafe4f085f3ecc80fa09e/.gopath-loader/pkg/mod/sigs.k8s.io/lws@v0.8.0/config/default/manager_webhook_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gomod-cache/sigs.k8s.io/lws@v0.8.0/config/default/manager_config_patch.yaml`](https://github.com/llm-d/workload-variant-autoscaler/blob/46a611076ea6e421be0dafe4f085f3ecc80fa09e/.gomod-cache/sigs.k8s.io/lws@v0.8.0/config/default/manager_config_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gomod-cache/sigs.k8s.io/lws@v0.8.0/config/default/manager_metrics_patch.yaml`](https://github.com/llm-d/workload-variant-autoscaler/blob/46a611076ea6e421be0dafe4f085f3ecc80fa09e/.gomod-cache/sigs.k8s.io/lws@v0.8.0/config/default/manager_metrics_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gomod-cache/sigs.k8s.io/lws@v0.8.0/config/default/manager_webhook_patch.yaml`](https://github.com/llm-d/workload-variant-autoscaler/blob/46a611076ea6e421be0dafe4f085f3ecc80fa09e/.gomod-cache/sigs.k8s.io/lws@v0.8.0/config/default/manager_webhook_patch.yaml) |
| controller-manager | manager | ? | true | ? | [`.gomod-cache/sigs.k8s.io/lws@v0.8.0/config/manager/manager.yaml`](https://github.com/llm-d/workload-variant-autoscaler/blob/46a611076ea6e421be0dafe4f085f3ecc80fa09e/.gomod-cache/sigs.k8s.io/lws@v0.8.0/config/manager/manager.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gopath-loader/pkg/mod/sigs.k8s.io/lws@v0.8.0/config/default/manager_config_patch.yaml`](https://github.com/llm-d/workload-variant-autoscaler/blob/46a611076ea6e421be0dafe4f085f3ecc80fa09e/.gopath-loader/pkg/mod/sigs.k8s.io/lws@v0.8.0/config/default/manager_config_patch.yaml) |
| keda-metrics-apiserver | keda-metrics-apiserver | true | true | ? | [`.gopath-loader/pkg/mod/github.com/kedacore/keda/v2@v2.18.0/config/metrics-server/deployment.yaml`](https://github.com/llm-d/workload-variant-autoscaler/blob/46a611076ea6e421be0dafe4f085f3ecc80fa09e/.gopath-loader/pkg/mod/github.com/kedacore/keda/v2@v2.18.0/config/metrics-server/deployment.yaml) |
| keda-metrics-apiserver | keda-metrics-apiserver | true | true | ? | [`.gomod-cache/github.com/kedacore/keda/v2@v2.18.0/config/metrics-server/deployment.yaml`](https://github.com/llm-d/workload-variant-autoscaler/blob/46a611076ea6e421be0dafe4f085f3ecc80fa09e/.gomod-cache/github.com/kedacore/keda/v2@v2.18.0/config/metrics-server/deployment.yaml) |
| keda-operator | keda-operator | true | true | ? | [`.gopath-loader/pkg/mod/github.com/kedacore/keda/v2@v2.18.0/config/manager/manager.yaml`](https://github.com/llm-d/workload-variant-autoscaler/blob/46a611076ea6e421be0dafe4f085f3ecc80fa09e/.gopath-loader/pkg/mod/github.com/kedacore/keda/v2@v2.18.0/config/manager/manager.yaml) |
| keda-operator | keda-operator | true | true | ? | [`.gomod-cache/github.com/kedacore/keda/v2@v2.18.0/config/manager/manager.yaml`](https://github.com/llm-d/workload-variant-autoscaler/blob/46a611076ea6e421be0dafe4f085f3ecc80fa09e/.gomod-cache/github.com/kedacore/keda/v2@v2.18.0/config/manager/manager.yaml) |
| llama-deployment | inference-server | ? | ? | ? | [`.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v1.2.1/tools/dynamic-lora-sidecar/deployment.yaml`](https://github.com/llm-d/workload-variant-autoscaler/blob/46a611076ea6e421be0dafe4f085f3ecc80fa09e/.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v1.2.1/tools/dynamic-lora-sidecar/deployment.yaml) |
| llama-deployment | inference-server | ? | ? | ? | [`.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v1.2.1/tools/dynamic-lora-sidecar/deployment.yaml`](https://github.com/llm-d/workload-variant-autoscaler/blob/46a611076ea6e421be0dafe4f085f3ecc80fa09e/.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v1.2.1/tools/dynamic-lora-sidecar/deployment.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `.gomod-cache/github.com/grpc-ecosystem/grpc-gateway/v2@v2.27.2/.github/Dockerfile` | golang:1.25.0 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/kedacore/keda/v2@v2.18.0/.devcontainer/Dockerfile` | golang:1.24.7 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/kedacore/keda/v2@v2.18.0/Dockerfile` | gcr.io/distroless/static:nonroot | 2 | 65532:65532 |  | multi-arch |  |  |
| `.gomod-cache/github.com/kedacore/keda/v2@v2.18.0/Dockerfile.adapter` | gcr.io/distroless/static:nonroot | 2 | 65532:65532 |  | multi-arch |  |  |
| `.gomod-cache/github.com/kedacore/keda/v2@v2.18.0/Dockerfile.webhooks` | gcr.io/distroless/static:nonroot | 2 | 65532:65532 |  | multi-arch |  |  |
| `.gomod-cache/github.com/pelletier/go-toml/v2@v2.2.4/Dockerfile` | scratch | 1 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/prometheus/client_golang@v1.23.2/Dockerfile` | quay.io/prometheus/busybox:latest | 2 |  |  |  |  | Unpinned base image: quay.io/prometheus/busybox:latest; No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/toolchain@v0.0.1-go1.25.0.linux-amd64/src/crypto/internal/boring/Dockerfile` | $ubuntu:focal | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/toolchain@v0.0.1-go1.25.0.linux-amd64/src/crypto/internal/fips140/nistec/fiat/Dockerfile` | coqorg/coq:8.13.2 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/x/net@v0.49.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/x/sys@v0.40.0/unix/linux/Dockerfile` | ubuntu:25.04 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/k8s.io/apiextensions-apiserver@v0.34.3/artifacts/simple-image/Dockerfile` | gcr.io/distroless/base-debian10:latest | 1 |  |  |  |  | Unpinned base image: gcr.io/distroless/base-debian10:latest; No USER directive found (defaults to root) |
| `.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v1.2.1/Dockerfile` | ${BASE_IMAGE} | 2 |  |  |  |  | Unpinned base image: ${BUILDER_IMAGE}; Unpinned base image: ${BASE_IMAGE}; No USER directive found (defaults to root) |
| `.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v1.2.1/tools/dynamic-lora-sidecar/Dockerfile` | python:3.10-slim-buster | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/sigs.k8s.io/lws@v0.8.0/Dockerfile` | ${BASE_IMAGE} | 2 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${BUILDER_IMAGE}; Unpinned base image: ${BASE_IMAGE} |
| `.gopath-loader/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.27.2/.github/Dockerfile` | golang:1.25.0 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/kedacore/keda/v2@v2.18.0/.devcontainer/Dockerfile` | golang:1.24.7 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/kedacore/keda/v2@v2.18.0/Dockerfile` | gcr.io/distroless/static:nonroot | 2 | 65532:65532 |  | multi-arch |  |  |
| `.gopath-loader/pkg/mod/github.com/kedacore/keda/v2@v2.18.0/Dockerfile.adapter` | gcr.io/distroless/static:nonroot | 2 | 65532:65532 |  | multi-arch |  |  |
| `.gopath-loader/pkg/mod/github.com/kedacore/keda/v2@v2.18.0/Dockerfile.webhooks` | gcr.io/distroless/static:nonroot | 2 | 65532:65532 |  | multi-arch |  |  |
| `.gopath-loader/pkg/mod/github.com/pelletier/go-toml/v2@v2.2.4/Dockerfile` | scratch | 1 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/prometheus/client_golang@v1.23.2/Dockerfile` | quay.io/prometheus/busybox:latest | 2 |  |  |  |  | Unpinned base image: quay.io/prometheus/busybox:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.0.linux-amd64/src/crypto/internal/boring/Dockerfile` | $ubuntu:focal | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.0.linux-amd64/src/crypto/internal/fips140/nistec/fiat/Dockerfile` | coqorg/coq:8.13.2 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/net@v0.49.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/sys@v0.40.0/unix/linux/Dockerfile` | ubuntu:25.04 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/k8s.io/apiextensions-apiserver@v0.34.3/artifacts/simple-image/Dockerfile` | gcr.io/distroless/base-debian10:latest | 1 |  |  |  |  | Unpinned base image: gcr.io/distroless/base-debian10:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v1.2.1/Dockerfile` | ${BASE_IMAGE} | 2 |  |  |  |  | Unpinned base image: ${BUILDER_IMAGE}; Unpinned base image: ${BASE_IMAGE}; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v1.2.1/tools/dynamic-lora-sidecar/Dockerfile` | python:3.10-slim-buster | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/sigs.k8s.io/lws@v0.8.0/Dockerfile` | ${BASE_IMAGE} | 2 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${BUILDER_IMAGE}; Unpinned base image: ${BASE_IMAGE} |
| `Dockerfile` | registry.access.redhat.com/ubi9/ubi-minimal@sha256:759f5f42d9d6ce2a705e290b7fc549e2d2cd39312c4fa345f93c02e4abb8da95 | 2 | 65532:65532 |  | multi-arch |  |  |
| `Dockerfile.konflux` | registry.access.redhat.com/ubi9/ubi-minimal@sha256:12db9874bd753eb98b1ab3d840e75de5d6842ac0604fbd68c012adefe97140be | 2 | 65532:65532 |  | multi-arch |  |  |

