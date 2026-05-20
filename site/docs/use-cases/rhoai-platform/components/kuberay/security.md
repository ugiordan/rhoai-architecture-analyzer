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
| kuberay-operator | kuberay-operator | ? | ? | ? | [`ray-operator/config/default-with-webhooks/manager_webhook_patch.yaml`](https://github.com/ray-project/kuberay/blob/923e422d3a7ba5713d151851551a653afb8fe7ee/ray-operator/config/default-with-webhooks/manager_webhook_patch.yaml) |
| kuberay-operator | kuberay-operator | ? | ? | ? | [`ray-operator/config/manager/manager.yaml`](https://github.com/ray-project/kuberay/blob/923e422d3a7ba5713d151851551a653afb8fe7ee/ray-operator/config/manager/manager.yaml) |
| the-deployment | the-container | ? | ? | ? | [`.gomod-cache/k8s.io/cli-runtime@v0.33.1/artifacts/kustomization/deployment.yaml`](https://github.com/ray-project/kuberay/blob/923e422d3a7ba5713d151851551a653afb8fe7ee/.gomod-cache/k8s.io/cli-runtime@v0.33.1/artifacts/kustomization/deployment.yaml) |
| the-deployment | the-container | ? | ? | ? | [`.gopath-loader/pkg/mod/k8s.io/cli-runtime@v0.33.1/artifacts/kustomization/deployment.yaml`](https://github.com/ray-project/kuberay/blob/923e422d3a7ba5713d151851551a653afb8fe7ee/.gopath-loader/pkg/mod/k8s.io/cli-runtime@v0.33.1/artifacts/kustomization/deployment.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `.gomod-cache/github.com/grpc-ecosystem/grpc-gateway/v2@v2.26.3/.github/Dockerfile` | golang:1.24.0 | 1 | vscode |  |  |  |  |
| `.gomod-cache/github.com/openshift/api@v0.0.0-20250602203052-b29811a290c7/Dockerfile.ocp` | registry.ci.openshift.org/ocp/4.19:base-rhel9 | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/prometheus/client_golang@v1.23.0/Dockerfile` | quay.io/prometheus/busybox:latest | 2 |  |  |  |  | Unpinned base image: quay.io/prometheus/busybox:latest; No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/x/net@v0.43.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/x/sys@v0.35.0/unix/linux/Dockerfile` | ubuntu:24.10 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/k8s.io/apiextensions-apiserver@v0.34.1/artifacts/simple-image/Dockerfile` | gcr.io/distroless/base-debian10:latest | 1 |  |  |  |  | Unpinned base image: gcr.io/distroless/base-debian10:latest; No USER directive found (defaults to root) |
| `.gomod-cache/sigs.k8s.io/kustomize/kyaml@v0.19.0/fn/framework/example/Dockerfile` | alpine:latest | 2 |  |  |  |  | Unpinned base image: alpine:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.26.3/.github/Dockerfile` | golang:1.24.0 | 1 | vscode |  |  |  |  |
| `.gopath-loader/pkg/mod/github.com/openshift/api@v0.0.0-20250602203052-b29811a290c7/Dockerfile.ocp` | registry.ci.openshift.org/ocp/4.19:base-rhel9 | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/prometheus/client_golang@v1.23.0/Dockerfile` | quay.io/prometheus/busybox:latest | 2 |  |  |  |  | Unpinned base image: quay.io/prometheus/busybox:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/net@v0.43.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/sys@v0.35.0/unix/linux/Dockerfile` | ubuntu:24.10 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/k8s.io/apiextensions-apiserver@v0.34.1/artifacts/simple-image/Dockerfile` | gcr.io/distroless/base-debian10:latest | 1 |  |  |  |  | Unpinned base image: gcr.io/distroless/base-debian10:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/sigs.k8s.io/kustomize/kyaml@v0.19.0/fn/framework/example/Dockerfile` | alpine:latest | 2 |  |  |  |  | Unpinned base image: alpine:latest; No USER directive found (defaults to root) |
| `apiserver/Dockerfile` | scratch | 2 | 65532:65532 |  |  |  | Unpinned base image: scratch |
| `benchmark/perf-tests/images/ray-pytorch/Dockerfile` | rayproject/ray:2.46.0 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `dashboard/Dockerfile` | base | 4 | nextjs |  |  |  | Unpinned base image: base; Unpinned base image: base; Unpinned base image: base |
| `experimental/Dockerfile` | scratch | 2 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `proto/Dockerfile` | golang:1.24.0-bullseye | 1 | 65532:65532 |  |  |  |  |
| `ray-operator/Dockerfile` | gcr.io/distroless/base-debian12:nonroot | 2 | 65532:65532 |  |  |  |  |
| `ray-operator/Dockerfile.buildx` | gcr.io/distroless/base-debian12:nonroot | 1 | 65532:65532 |  | multi-arch |  |  |
| `ray-operator/Dockerfile.konflux` | registry.access.redhat.com/ubi9/ubi-minimal@sha256:24650313873554b6ba16c1a1b6b9f9142604f6ab735113e1695faf2dd07fdede | 2 | 65532:65532 |  | multi-arch |  |  |
| `ray-operator/Dockerfile.rhoai` | registry.access.redhat.com/ubi9/ubi:latest | 2 | 65532:65532 |  |  |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi:latest |

