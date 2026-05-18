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
| kuberay-operator | kuberay-operator | ? | ? | ? | [`ray-operator/config/overlays/test-overrides/deployment-override.yaml`](https://github.com/ray-project/kuberay/blob/1466289a1bcff3df4b0be5f0c804e178b7aa8e05/ray-operator/config/overlays/test-overrides/deployment-override.yaml) |
| kuberay-operator | kuberay-operator | ? | ? | ? | [`ray-operator/config/default-with-webhooks/manager_webhook_patch.yaml`](https://github.com/ray-project/kuberay/blob/1466289a1bcff3df4b0be5f0c804e178b7aa8e05/ray-operator/config/default-with-webhooks/manager_webhook_patch.yaml) |
| kuberay-operator | kuberay-operator | ? | true | ? | [`ray-operator/config/manager/manager.yaml`](https://github.com/ray-project/kuberay/blob/1466289a1bcff3df4b0be5f0c804e178b7aa8e05/ray-operator/config/manager/manager.yaml) |
| the-deployment | the-container | ? | ? | ? | [`.gomod-cache/k8s.io/cli-runtime@v0.36.0/artifacts/kustomization/deployment.yaml`](https://github.com/ray-project/kuberay/blob/1466289a1bcff3df4b0be5f0c804e178b7aa8e05/.gomod-cache/k8s.io/cli-runtime@v0.36.0/artifacts/kustomization/deployment.yaml) |
| the-deployment | the-container | ? | ? | ? | [`.gopath-loader/pkg/mod/k8s.io/cli-runtime@v0.36.0/artifacts/kustomization/deployment.yaml`](https://github.com/ray-project/kuberay/blob/1466289a1bcff3df4b0be5f0c804e178b7aa8e05/.gopath-loader/pkg/mod/k8s.io/cli-runtime@v0.36.0/artifacts/kustomization/deployment.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `.gomod-cache/github.com/grpc-ecosystem/grpc-gateway/v2@v2.27.7/.github/Dockerfile` | golang:1.25.6 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/prometheus/client_golang@v1.23.2/Dockerfile` | quay.io/prometheus/busybox:latest | 2 |  |  |  |  | Unpinned base image: quay.io/prometheus/busybox:latest; No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/toolchain@v0.0.1-go1.26.0.linux-amd64/src/crypto/internal/boring/Dockerfile` | $ubuntu:focal | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/toolchain@v0.0.1-go1.26.0.linux-amd64/src/crypto/internal/fips140/nistec/fiat/Dockerfile` | coqorg/coq:8.13.2 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/x/net@v0.49.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/x/sys@v0.40.0/unix/linux/Dockerfile` | ubuntu:25.04 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/k8s.io/apiextensions-apiserver@v0.36.0/artifacts/simple-image/Dockerfile` | gcr.io/distroless/base-debian10:latest | 1 |  |  |  |  | Unpinned base image: gcr.io/distroless/base-debian10:latest; No USER directive found (defaults to root) |
| `.gomod-cache/sigs.k8s.io/gateway-api@v1.4.1/docker/Dockerfile.echo-advanced` | gcr.io/istio-release/app:1.21.2 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/sigs.k8s.io/gateway-api@v1.4.1/docker/Dockerfile.echo-basic` | gcr.io/distroless/static:nonroot | 2 | nonroot:nonroot |  |  |  |  |
| `.gomod-cache/sigs.k8s.io/kustomize/kyaml@v0.21.1/fn/framework/example/Dockerfile` | alpine:latest | 2 |  |  |  |  | Unpinned base image: alpine:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.27.7/.github/Dockerfile` | golang:1.25.6 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/prometheus/client_golang@v1.23.2/Dockerfile` | quay.io/prometheus/busybox:latest | 2 |  |  |  |  | Unpinned base image: quay.io/prometheus/busybox:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.0.linux-amd64/src/crypto/internal/boring/Dockerfile` | $ubuntu:focal | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.0.linux-amd64/src/crypto/internal/fips140/nistec/fiat/Dockerfile` | coqorg/coq:8.13.2 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/net@v0.49.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/sys@v0.40.0/unix/linux/Dockerfile` | ubuntu:25.04 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/k8s.io/apiextensions-apiserver@v0.36.0/artifacts/simple-image/Dockerfile` | gcr.io/distroless/base-debian10:latest | 1 |  |  |  |  | Unpinned base image: gcr.io/distroless/base-debian10:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api@v1.4.1/docker/Dockerfile.echo-advanced` | gcr.io/istio-release/app:1.21.2 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api@v1.4.1/docker/Dockerfile.echo-basic` | gcr.io/distroless/static:nonroot | 2 | nonroot:nonroot |  |  |  |  |
| `.gopath-loader/pkg/mod/sigs.k8s.io/kustomize/kyaml@v0.21.1/fn/framework/example/Dockerfile` | alpine:latest | 2 |  |  |  |  | Unpinned base image: alpine:latest; No USER directive found (defaults to root) |
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

