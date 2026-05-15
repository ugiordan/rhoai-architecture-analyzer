# kube-auth-proxy: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `.devcontainer/Dockerfile` | mcr.microsoft.com/vscode/devcontainers/go:1-1.23 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/cpuguy83/go-md2man/v2@v2.0.6/Dockerfile` | scratch | 2 |  |  | multi-arch |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/creack/pty@v1.1.18/Dockerfile.golang` | golang:${GOVERSION} | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/creack/pty@v1.1.18/Dockerfile.riscv` | golang:1.13 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/envoyproxy/go-control-plane@v0.14.0/Dockerfile.ci` | golang:1.23@sha256:77a21b3e354c03e9f66b13bc39f4f0db8085c70f8414406af66b29c6d6c4dd85 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/envoyproxy/protoc-gen-validate@v1.3.0/Dockerfile` | ubuntu:focal | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/grpc-ecosystem/grpc-gateway/v2@v2.26.3/.github/Dockerfile` | golang:1.24.0 | 1 | vscode |  |  |  |  |
| `.gomod-cache/github.com/pelletier/go-toml/v2@v2.2.4/Dockerfile` | scratch | 1 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/prometheus/client_golang@v1.23.2/Dockerfile` | quay.io/prometheus/busybox:latest | 2 |  |  |  |  | Unpinned base image: quay.io/prometheus/busybox:latest; No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/toolchain@v0.0.1-go1.25.3.linux-amd64/src/crypto/internal/boring/Dockerfile` | $ubuntu:focal | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/toolchain@v0.0.1-go1.25.3.linux-amd64/src/crypto/internal/fips140/nistec/fiat/Dockerfile` | coqorg/coq:8.13.2 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/x/net@v0.49.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/x/sys@v0.40.0/unix/linux/Dockerfile` | ubuntu:25.04 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/cpuguy83/go-md2man/v2@v2.0.6/Dockerfile` | scratch | 2 |  |  | multi-arch |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/creack/pty@v1.1.18/Dockerfile.golang` | golang:${GOVERSION} | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/creack/pty@v1.1.18/Dockerfile.riscv` | golang:1.13 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/envoyproxy/go-control-plane@v0.14.0/Dockerfile.ci` | golang:1.23@sha256:77a21b3e354c03e9f66b13bc39f4f0db8085c70f8414406af66b29c6d6c4dd85 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/envoyproxy/protoc-gen-validate@v1.3.0/Dockerfile` | ubuntu:focal | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.26.3/.github/Dockerfile` | golang:1.24.0 | 1 | vscode |  |  |  |  |
| `.gopath-loader/pkg/mod/github.com/pelletier/go-toml/v2@v2.2.4/Dockerfile` | scratch | 1 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/prometheus/client_golang@v1.23.2/Dockerfile` | quay.io/prometheus/busybox:latest | 2 |  |  |  |  | Unpinned base image: quay.io/prometheus/busybox:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.3.linux-amd64/src/crypto/internal/boring/Dockerfile` | $ubuntu:focal | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.3.linux-amd64/src/crypto/internal/fips140/nistec/fiat/Dockerfile` | coqorg/coq:8.13.2 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/net@v0.49.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/sys@v0.40.0/unix/linux/Dockerfile` | ubuntu:25.04 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `Dockerfile` | ${RUNTIME_IMAGE} | 2 |  |  | multi-arch |  | Unpinned base image: ${BUILD_IMAGE}; Unpinned base image: ${RUNTIME_IMAGE}; No USER directive found (defaults to root) |
| `Dockerfile.konflux` | registry.access.redhat.com/ubi9/ubi-minimal@sha256:12db9874bd753eb98b1ab3d840e75de5d6842ac0604fbd68c012adefe97140be | 2 | 1001 |  | multi-arch |  |  |
| `Dockerfile.redhat` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 2 | 1001 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |
| `kube-rbac-proxy/Dockerfile` | $BASEIMAGE | 1 | 65532:65532 |  |  |  | Unpinned base image: $BASEIMAGE |
| `kube-rbac-proxy/Dockerfile.ocp` | registry.ci.openshift.org/ocp/4.20:base-rhel9 | 2 | 65534 |  |  |  |  |

