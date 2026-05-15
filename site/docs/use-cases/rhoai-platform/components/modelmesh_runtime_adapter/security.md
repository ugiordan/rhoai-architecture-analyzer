# modelmesh-runtime-adapter: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `.gomod-cache/github.com/!i!b!m/ibm-cos-sdk-go@v1.9.1/awstesting/sandbox/Dockerfile.golang-tip` | buildpack-deps:buster-scm | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/!i!b!m/ibm-cos-sdk-go@v1.9.1/awstesting/sandbox/Dockerfile.test.go1.10` | golang:1.10 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/!i!b!m/ibm-cos-sdk-go@v1.9.1/awstesting/sandbox/Dockerfile.test.go1.11` | golang:1.11 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/!i!b!m/ibm-cos-sdk-go@v1.9.1/awstesting/sandbox/Dockerfile.test.go1.12` | golang:1.12 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/!i!b!m/ibm-cos-sdk-go@v1.9.1/awstesting/sandbox/Dockerfile.test.go1.13` | golang:1.13 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/!i!b!m/ibm-cos-sdk-go@v1.9.1/awstesting/sandbox/Dockerfile.test.go1.14` | golang:1.14 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/!i!b!m/ibm-cos-sdk-go@v1.9.1/awstesting/sandbox/Dockerfile.test.go1.15` | golang:1.15 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/!i!b!m/ibm-cos-sdk-go@v1.9.1/awstesting/sandbox/Dockerfile.test.go1.16` | golang:1.16 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/!i!b!m/ibm-cos-sdk-go@v1.9.1/awstesting/sandbox/Dockerfile.test.go1.5` | golang:1.5 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/!i!b!m/ibm-cos-sdk-go@v1.9.1/awstesting/sandbox/Dockerfile.test.go1.5-novendorexp` | golang:1.5 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/!i!b!m/ibm-cos-sdk-go@v1.9.1/awstesting/sandbox/Dockerfile.test.go1.6` | golang:1.6 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/!i!b!m/ibm-cos-sdk-go@v1.9.1/awstesting/sandbox/Dockerfile.test.go1.7` | golang:1.7 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/!i!b!m/ibm-cos-sdk-go@v1.9.1/awstesting/sandbox/Dockerfile.test.go1.8` | golang:1.8 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/!i!b!m/ibm-cos-sdk-go@v1.9.1/awstesting/sandbox/Dockerfile.test.go1.9` | golang:1.9 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/!i!b!m/ibm-cos-sdk-go@v1.9.1/awstesting/sandbox/Dockerfile.test.gotip` | aws-golang:tip | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/x/net@v0.33.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/x/sys@v0.28.0/unix/linux/Dockerfile` | ubuntu:24.10 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/google.golang.org/grpc@v1.56.3/interop/xds/client/Dockerfile` | alpine | 2 |  |  |  |  | Unpinned base image: alpine; No USER directive found (defaults to root) |
| `.gomod-cache/google.golang.org/grpc@v1.56.3/interop/xds/server/Dockerfile` | alpine | 2 |  |  |  |  | Unpinned base image: alpine; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/!i!b!m/ibm-cos-sdk-go@v1.9.1/awstesting/sandbox/Dockerfile.golang-tip` | buildpack-deps:buster-scm | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/!i!b!m/ibm-cos-sdk-go@v1.9.1/awstesting/sandbox/Dockerfile.test.go1.10` | golang:1.10 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/!i!b!m/ibm-cos-sdk-go@v1.9.1/awstesting/sandbox/Dockerfile.test.go1.11` | golang:1.11 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/!i!b!m/ibm-cos-sdk-go@v1.9.1/awstesting/sandbox/Dockerfile.test.go1.12` | golang:1.12 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/!i!b!m/ibm-cos-sdk-go@v1.9.1/awstesting/sandbox/Dockerfile.test.go1.13` | golang:1.13 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/!i!b!m/ibm-cos-sdk-go@v1.9.1/awstesting/sandbox/Dockerfile.test.go1.14` | golang:1.14 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/!i!b!m/ibm-cos-sdk-go@v1.9.1/awstesting/sandbox/Dockerfile.test.go1.15` | golang:1.15 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/!i!b!m/ibm-cos-sdk-go@v1.9.1/awstesting/sandbox/Dockerfile.test.go1.16` | golang:1.16 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/!i!b!m/ibm-cos-sdk-go@v1.9.1/awstesting/sandbox/Dockerfile.test.go1.5` | golang:1.5 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/!i!b!m/ibm-cos-sdk-go@v1.9.1/awstesting/sandbox/Dockerfile.test.go1.5-novendorexp` | golang:1.5 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/!i!b!m/ibm-cos-sdk-go@v1.9.1/awstesting/sandbox/Dockerfile.test.go1.6` | golang:1.6 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/!i!b!m/ibm-cos-sdk-go@v1.9.1/awstesting/sandbox/Dockerfile.test.go1.7` | golang:1.7 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/!i!b!m/ibm-cos-sdk-go@v1.9.1/awstesting/sandbox/Dockerfile.test.go1.8` | golang:1.8 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/!i!b!m/ibm-cos-sdk-go@v1.9.1/awstesting/sandbox/Dockerfile.test.go1.9` | golang:1.9 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/!i!b!m/ibm-cos-sdk-go@v1.9.1/awstesting/sandbox/Dockerfile.test.gotip` | aws-golang:tip | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/net@v0.33.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/sys@v0.28.0/unix/linux/Dockerfile` | ubuntu:24.10 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/google.golang.org/grpc@v1.56.3/interop/xds/client/Dockerfile` | alpine | 2 |  |  |  |  | Unpinned base image: alpine; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/google.golang.org/grpc@v1.56.3/interop/xds/server/Dockerfile` | alpine | 2 |  |  |  |  | Unpinned base image: alpine; No USER directive found (defaults to root) |
| `Dockerfile` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 3 | ${USER} |  | multi-arch |  | Unpinned base image: $BUILD_BASE; Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |
| `Dockerfile.konflux` | registry.access.redhat.com/ubi9/ubi-minimal@sha256:12db9874bd753eb98b1ab3d840e75de5d6842ac0604fbd68c012adefe97140be | 2 | ${USER} |  | multi-arch |  |  |

