# batch-gateway: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `.gomod-cache/github.com/grpc-ecosystem/grpc-gateway/v2@v2.28.0/.github/Dockerfile` | golang:1.26.0 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/jackc/pgx/v5@v5.9.2/.devcontainer/Dockerfile` | mcr.microsoft.com/devcontainers/go:2-1.26-trixie | 1 | vscode |  |  |  |  |
| `.gomod-cache/github.com/prometheus/client_golang@v1.23.2/Dockerfile` | quay.io/prometheus/busybox:latest | 2 |  |  |  |  | Unpinned base image: quay.io/prometheus/busybox:latest; No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/toolchain@v0.0.1-go1.25.0.linux-amd64/src/crypto/internal/boring/Dockerfile` | $ubuntu:focal | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/toolchain@v0.0.1-go1.25.0.linux-amd64/src/crypto/internal/fips140/nistec/fiat/Dockerfile` | coqorg/coq:8.13.2 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/x/net@v0.52.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/x/sys@v0.42.0/unix/linux/Dockerfile` | ubuntu:25.04 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.28.0/.github/Dockerfile` | golang:1.26.0 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/jackc/pgx/v5@v5.9.2/.devcontainer/Dockerfile` | mcr.microsoft.com/devcontainers/go:2-1.26-trixie | 1 | vscode |  |  |  |  |
| `.gopath-loader/pkg/mod/github.com/prometheus/client_golang@v1.23.2/Dockerfile` | quay.io/prometheus/busybox:latest | 2 |  |  |  |  | Unpinned base image: quay.io/prometheus/busybox:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.0.linux-amd64/src/crypto/internal/boring/Dockerfile` | $ubuntu:focal | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.0.linux-amd64/src/crypto/internal/fips140/nistec/fiat/Dockerfile` | coqorg/coq:8.13.2 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/net@v0.52.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/sys@v0.42.0/unix/linux/Dockerfile` | ubuntu:25.04 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `docker/Dockerfile.apiserver` | gcr.io/distroless/static:nonroot | 2 | 65532:65532 |  | multi-arch |  |  |
| `docker/Dockerfile.apiserver.konflux` | registry.access.redhat.com/ubi9/ubi-micro:latest | 2 | 1001:1001 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-micro:latest |
| `docker/Dockerfile.gc` | gcr.io/distroless/static:nonroot | 2 | 65532:65532 |  | multi-arch |  |  |
| `docker/Dockerfile.gc.konflux` | registry.access.redhat.com/ubi9/ubi-micro:latest | 2 | 1001:1001 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-micro:latest |
| `docker/Dockerfile.processor` | gcr.io/distroless/static:nonroot | 2 | 65532:65532 |  | multi-arch |  |  |
| `docker/Dockerfile.processor.konflux` | registry.access.redhat.com/ubi9/ubi-micro:latest | 2 | 1001:1001 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-micro:latest |

