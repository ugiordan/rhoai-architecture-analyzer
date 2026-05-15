# llm-d-inference-scheduler: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

### Secrets Referenced

| Name | Type | Referenced By |
|------|------|---------------|
| cacerts | Opaque | deployment/istiod-llm-d-gateway |
| istio-kubeconfig | Opaque | deployment/istiod-llm-d-gateway |
| istiod-tls | Opaque | deployment/istiod-llm-d-gateway |

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| ${EPP_NAME} | epp | ? | ? | ? | [`deploy/components/inference-gateway/deployments.yaml`](https://github.com/llm-d/llm-d-inference-scheduler/blob/e004b4d81f11b319f97f559a7e12b0a00c08fa58/deploy/components/inference-gateway/deployments.yaml) |
| ${EPP_NAME} | uds-tokenizer | ? | ? | ? | [`deploy/components/inference-gateway/deployments.yaml`](https://github.com/llm-d/llm-d-inference-scheduler/blob/e004b4d81f11b319f97f559a7e12b0a00c08fa58/deploy/components/inference-gateway/deployments.yaml) |
| ${MODEL_NAME_SAFE}-vllm-sim | vllm | ? | ? | ? | [`deploy/components/vllm-sim/deployments.yaml`](https://github.com/llm-d/llm-d-inference-scheduler/blob/e004b4d81f11b319f97f559a7e12b0a00c08fa58/deploy/components/vllm-sim/deployments.yaml) |
| 0 | cmd | ? | ? | ? | [`.gomod-cache/github.com/llm-d/llm-d-kv-cache@v0.7.1/deploy/common/statefulset.yaml`](https://github.com/llm-d/llm-d-inference-scheduler/blob/e004b4d81f11b319f97f559a7e12b0a00c08fa58/.gomod-cache/github.com/llm-d/llm-d-kv-cache@v0.7.1/deploy/common/statefulset.yaml) |
| 0 | cmd | ? | ? | ? | [`.gopath-loader/pkg/mod/github.com/llm-d/llm-d-kv-cache@v0.7.1/deploy/common/statefulset.yaml`](https://github.com/llm-d/llm-d-inference-scheduler/blob/e004b4d81f11b319f97f559a7e12b0a00c08fa58/.gopath-loader/pkg/mod/github.com/llm-d/llm-d-kv-cache@v0.7.1/deploy/common/statefulset.yaml) |
| 0 | cmd | ? | ? | ? | [`deploy/environments/kubernetes-base/common/statefulset.yaml`](https://github.com/llm-d/llm-d-inference-scheduler/blob/e004b4d81f11b319f97f559a7e12b0a00c08fa58/deploy/environments/kubernetes-base/common/statefulset.yaml) |
| istiod-llm-d-gateway | discovery | true | true | ? | [`deploy/components/istio-control-plane/deployments.yaml`](https://github.com/llm-d/llm-d-inference-scheduler/blob/e004b4d81f11b319f97f559a7e12b0a00c08fa58/deploy/components/istio-control-plane/deployments.yaml) |
| vllm-sim-d | vllm | ? | ? | ? | [`deploy/components/vllm-sim-epd/deployments.yaml`](https://github.com/llm-d/llm-d-inference-scheduler/blob/e004b4d81f11b319f97f559a7e12b0a00c08fa58/deploy/components/vllm-sim-epd/deployments.yaml) |
| vllm-sim-d | vllm | ? | ? | ? | [`deploy/components/vllm-sim-pd/deployments.yaml`](https://github.com/llm-d/llm-d-inference-scheduler/blob/e004b4d81f11b319f97f559a7e12b0a00c08fa58/deploy/components/vllm-sim-pd/deployments.yaml) |
| vllm-sim-e | vllm | ? | ? | ? | [`deploy/components/vllm-sim-epd/deployments.yaml`](https://github.com/llm-d/llm-d-inference-scheduler/blob/e004b4d81f11b319f97f559a7e12b0a00c08fa58/deploy/components/vllm-sim-epd/deployments.yaml) |
| vllm-sim-p | vllm | ? | ? | ? | [`deploy/components/vllm-sim-epd/deployments.yaml`](https://github.com/llm-d/llm-d-inference-scheduler/blob/e004b4d81f11b319f97f559a7e12b0a00c08fa58/deploy/components/vllm-sim-epd/deployments.yaml) |
| vllm-sim-p | vllm | ? | ? | ? | [`deploy/components/vllm-sim-pd/deployments.yaml`](https://github.com/llm-d/llm-d-inference-scheduler/blob/e004b4d81f11b319f97f559a7e12b0a00c08fa58/deploy/components/vllm-sim-pd/deployments.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `.gomod-cache/github.com/envoyproxy/protoc-gen-validate@v1.3.0/Dockerfile` | ubuntu:focal | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/grpc-ecosystem/grpc-gateway/v2@v2.28.0/.github/Dockerfile` | golang:1.26.0 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/llm-d/llm-d-kv-cache@v0.7.1/Dockerfile` | registry.access.redhat.com/ubi9/ubi:latest | 3 | 65532:65532 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi:latest |
| `.gomod-cache/github.com/llm-d/llm-d-kv-cache@v0.7.1/kv_connectors/llmd_fs_backend/Dockerfile.wheel` | nvcr.io/nvidia/cuda:${CUDA_VERSION}-devel-ubuntu22.04 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/llm-d/llm-d-kv-cache@v0.7.1/kv_connectors/pvc_evictor/Dockerfile` | python:3.12-slim | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/llm-d/llm-d-kv-cache@v0.7.1/services/uds_tokenizer/Dockerfile` | python:3.12-slim | 2 | 65532:65532 |  | multi-arch |  |  |
| `.gomod-cache/github.com/prometheus/client_golang@v1.23.2/Dockerfile` | quay.io/prometheus/busybox:latest | 2 |  |  |  |  | Unpinned base image: quay.io/prometheus/busybox:latest; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/prometheus/prometheus@v0.310.0/Dockerfile` | quay.io/prometheus/busybox-${OS}-${ARCH}:latest | 1 | nobody |  |  |  | Unpinned base image: quay.io/prometheus/busybox-${OS}-${ARCH}:latest |
| `.gomod-cache/github.com/prometheus/prometheus@v0.310.0/Dockerfile.distroless` | gcr.io/distroless/static-debian13:nonroot-${DISTROLESS_ARCH} | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/toolchain@v0.0.1-go1.25.8.linux-amd64/src/crypto/internal/boring/Dockerfile` | $ubuntu:focal | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/toolchain@v0.0.1-go1.25.8.linux-amd64/src/crypto/internal/fips140/nistec/fiat/Dockerfile` | coqorg/coq:8.13.2 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/x/net@v0.52.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/x/sys@v0.42.0/unix/linux/Dockerfile` | ubuntu:25.04 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/k8s.io/apiextensions-apiserver@v0.35.4/artifacts/simple-image/Dockerfile` | gcr.io/distroless/base-debian10:latest | 1 |  |  |  |  | Unpinned base image: gcr.io/distroless/base-debian10:latest; No USER directive found (defaults to root) |
| `.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v1.5.0/Dockerfile` | ${BASE_IMAGE} | 2 |  |  |  |  | Unpinned base image: ${BUILDER_IMAGE}; Unpinned base image: ${BASE_IMAGE}; No USER directive found (defaults to root) |
| `.gomod-cache/sigs.k8s.io/gateway-api@v1.5.1/docker/Dockerfile.echo-advanced` | gcr.io/istio-release/app:1.21.2 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/sigs.k8s.io/gateway-api@v1.5.1/docker/Dockerfile.echo-basic` | gcr.io/distroless/static:nonroot | 2 | nonroot:nonroot |  |  |  |  |
| `.gomod-cache/sigs.k8s.io/kustomize/kyaml@v0.21.1/fn/framework/example/Dockerfile` | alpine:latest | 2 |  |  |  |  | Unpinned base image: alpine:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/envoyproxy/protoc-gen-validate@v1.3.0/Dockerfile` | ubuntu:focal | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.28.0/.github/Dockerfile` | golang:1.26.0 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/llm-d/llm-d-kv-cache@v0.7.1/Dockerfile` | registry.access.redhat.com/ubi9/ubi:latest | 3 | 65532:65532 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi:latest |
| `.gopath-loader/pkg/mod/github.com/llm-d/llm-d-kv-cache@v0.7.1/kv_connectors/llmd_fs_backend/Dockerfile.wheel` | nvcr.io/nvidia/cuda:${CUDA_VERSION}-devel-ubuntu22.04 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/llm-d/llm-d-kv-cache@v0.7.1/kv_connectors/pvc_evictor/Dockerfile` | python:3.12-slim | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/llm-d/llm-d-kv-cache@v0.7.1/services/uds_tokenizer/Dockerfile` | python:3.12-slim | 2 | 65532:65532 |  | multi-arch |  |  |
| `.gopath-loader/pkg/mod/github.com/prometheus/client_golang@v1.23.2/Dockerfile` | quay.io/prometheus/busybox:latest | 2 |  |  |  |  | Unpinned base image: quay.io/prometheus/busybox:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/prometheus/prometheus@v0.310.0/Dockerfile` | quay.io/prometheus/busybox-${OS}-${ARCH}:latest | 1 | nobody |  |  |  | Unpinned base image: quay.io/prometheus/busybox-${OS}-${ARCH}:latest |
| `.gopath-loader/pkg/mod/github.com/prometheus/prometheus@v0.310.0/Dockerfile.distroless` | gcr.io/distroless/static-debian13:nonroot-${DISTROLESS_ARCH} | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.8.linux-amd64/src/crypto/internal/boring/Dockerfile` | $ubuntu:focal | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.8.linux-amd64/src/crypto/internal/fips140/nistec/fiat/Dockerfile` | coqorg/coq:8.13.2 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/net@v0.52.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/sys@v0.42.0/unix/linux/Dockerfile` | ubuntu:25.04 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/k8s.io/apiextensions-apiserver@v0.35.4/artifacts/simple-image/Dockerfile` | gcr.io/distroless/base-debian10:latest | 1 |  |  |  |  | Unpinned base image: gcr.io/distroless/base-debian10:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v1.5.0/Dockerfile` | ${BASE_IMAGE} | 2 |  |  |  |  | Unpinned base image: ${BUILDER_IMAGE}; Unpinned base image: ${BASE_IMAGE}; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api@v1.5.1/docker/Dockerfile.echo-advanced` | gcr.io/istio-release/app:1.21.2 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api@v1.5.1/docker/Dockerfile.echo-basic` | gcr.io/distroless/static:nonroot | 2 | nonroot:nonroot |  |  |  |  |
| `.gopath-loader/pkg/mod/sigs.k8s.io/kustomize/kyaml@v0.21.1/fn/framework/example/Dockerfile` | alpine:latest | 2 |  |  |  |  | Unpinned base image: alpine:latest; No USER directive found (defaults to root) |
| `Dockerfile.builder` | quay.io/projectquay/golang:1.25 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `Dockerfile.epp` | ${BASE_IMAGE} | 2 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${BASE_IMAGE} |
| `Dockerfile.epp.konflux` | registry.access.redhat.com/ubi9/ubi-minimal:9.7@sha256:d91be7cea9f03a757d69ad7fcdfcd7849dba820110e7980d5e2a1f46ed06ea3b | 2 | 65532:65532 |  | multi-arch |  |  |
| `Dockerfile.sidecar` | ${BASE_IMAGE} | 2 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${BASE_IMAGE} |
| `Dockerfile.sidecar.konflux` | registry.access.redhat.com/ubi9/ubi-micro:9.7 | 2 | 65532:65532 |  | multi-arch |  |  |

