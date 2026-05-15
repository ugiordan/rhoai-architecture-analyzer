# rest-proxy: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `.gomod-cache/github.com/grpc-ecosystem/grpc-gateway/v2@v2.15.0/.github/Dockerfile` | golang:1.19.4 | 1 | vscode |  |  |  |  |
| `.gomod-cache/github.com/grpc-ecosystem/grpc-gateway/v2@v2.15.0/.github/plugins/protoc-gen-grpc-gateway/Dockerfile` | scratch | 2 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/grpc-ecosystem/grpc-gateway/v2@v2.15.0/.github/plugins/protoc-gen-openapiv2/Dockerfile` | scratch | 2 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/x/net@v0.33.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/x/sys@v0.28.0/unix/linux/Dockerfile` | ubuntu:24.10 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/google.golang.org/grpc@v1.56.3/interop/xds/client/Dockerfile` | alpine | 2 |  |  |  |  | Unpinned base image: alpine; No USER directive found (defaults to root) |
| `.gomod-cache/google.golang.org/grpc@v1.56.3/interop/xds/server/Dockerfile` | alpine | 2 |  |  |  |  | Unpinned base image: alpine; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.15.0/.github/Dockerfile` | golang:1.19.4 | 1 | vscode |  |  |  |  |
| `.gopath-loader/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.15.0/.github/plugins/protoc-gen-grpc-gateway/Dockerfile` | scratch | 2 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.15.0/.github/plugins/protoc-gen-openapiv2/Dockerfile` | scratch | 2 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/net@v0.33.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/sys@v0.28.0/unix/linux/Dockerfile` | ubuntu:24.10 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/google.golang.org/grpc@v1.56.3/interop/xds/client/Dockerfile` | alpine | 2 |  |  |  |  | Unpinned base image: alpine; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/google.golang.org/grpc@v1.56.3/interop/xds/server/Dockerfile` | alpine | 2 |  |  |  |  | Unpinned base image: alpine; No USER directive found (defaults to root) |
| `Dockerfile` | registry.access.redhat.com/ubi9/ubi-micro:9.5 | 3 | ${USER} |  | multi-arch |  | Unpinned base image: develop |
| `Dockerfile.konflux` | registry.access.redhat.com/ubi9/ubi-minimal@sha256:12db9874bd753eb98b1ab3d840e75de5d6842ac0604fbd68c012adefe97140be | 2 | ${USER} |  | multi-arch |  |  |

