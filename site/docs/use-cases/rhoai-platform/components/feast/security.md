# feast: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| controller-manager | manager | ? | ? | ? | [`infra/feast-operator/config/default/manager_config_patch.yaml`](https://github.com/feast-dev/feast/blob/0d2984d427b0773c832dc7ff0b20ae3193d9ca04/infra/feast-operator/config/default/manager_config_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`infra/feast-operator/config/manager/manager.yaml`](https://github.com/feast-dev/feast/blob/0d2984d427b0773c832dc7ff0b20ae3193d9ca04/infra/feast-operator/config/manager/manager.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `.gomod-cache/github.com/apache/thrift@v0.21.0/build/docker/msvc2017/Dockerfile` | microsoft/dotnet-framework:4.7.1 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/apache/thrift@v0.21.0/build/docker/ubuntu-focal/Dockerfile` | buildpack-deps:focal-scm | 1 | ${user} |  |  |  |  |
| `.gomod-cache/github.com/apache/thrift@v0.21.0/build/docker/ubuntu-jammy/Dockerfile` | buildpack-deps:jammy-scm | 1 | ${user} |  |  |  |  |
| `.gomod-cache/github.com/aws/aws-sdk-go-v2@v1.36.4/internal/awstesting/sandbox/Dockerfile.golang-tip` | buildpack-deps:buster-scm | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go-v2@v1.36.4/internal/awstesting/sandbox/Dockerfile.test.gotip` | aws-golang:tip | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go-v2@v1.36.4/internal/awstesting/sandbox/Dockerfile.test.goversion` | golang:${GO_VERSION} | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/envoyproxy/protoc-gen-validate@v1.2.1/Dockerfile` | ubuntu:focal | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/grpc-ecosystem/grpc-gateway/v2@v2.27.2/.github/Dockerfile` | golang:1.25.0 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/mattn/go-sqlite3@v1.14.23/_example/simple/Dockerfile` | scratch | 2 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/prometheus/client_golang@v1.23.2/Dockerfile` | quay.io/prometheus/busybox:latest | 2 |  |  |  |  | Unpinned base image: quay.io/prometheus/busybox:latest; No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/toolchain@v0.0.1-go1.25.0.linux-amd64/src/crypto/internal/boring/Dockerfile` | $ubuntu:focal | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/toolchain@v0.0.1-go1.25.0.linux-amd64/src/crypto/internal/fips140/nistec/fiat/Dockerfile` | coqorg/coq:8.13.2 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/x/net@v0.47.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/x/sys@v0.42.0/unix/linux/Dockerfile` | ubuntu:25.04 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/apache/thrift@v0.21.0/build/docker/msvc2017/Dockerfile` | microsoft/dotnet-framework:4.7.1 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/apache/thrift@v0.21.0/build/docker/ubuntu-focal/Dockerfile` | buildpack-deps:focal-scm | 1 | ${user} |  |  |  |  |
| `.gopath-loader/pkg/mod/github.com/apache/thrift@v0.21.0/build/docker/ubuntu-jammy/Dockerfile` | buildpack-deps:jammy-scm | 1 | ${user} |  |  |  |  |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go-v2@v1.36.4/internal/awstesting/sandbox/Dockerfile.golang-tip` | buildpack-deps:buster-scm | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go-v2@v1.36.4/internal/awstesting/sandbox/Dockerfile.test.gotip` | aws-golang:tip | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go-v2@v1.36.4/internal/awstesting/sandbox/Dockerfile.test.goversion` | golang:${GO_VERSION} | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/envoyproxy/protoc-gen-validate@v1.2.1/Dockerfile` | ubuntu:focal | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.27.2/.github/Dockerfile` | golang:1.25.0 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/mattn/go-sqlite3@v1.14.23/_example/simple/Dockerfile` | scratch | 2 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/prometheus/client_golang@v1.23.2/Dockerfile` | quay.io/prometheus/busybox:latest | 2 |  |  |  |  | Unpinned base image: quay.io/prometheus/busybox:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.0.linux-amd64/src/crypto/internal/boring/Dockerfile` | $ubuntu:focal | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.0.linux-amd64/src/crypto/internal/fips140/nistec/fiat/Dockerfile` | coqorg/coq:8.13.2 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/net@v0.47.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/sys@v0.42.0/unix/linux/Dockerfile` | ubuntu:25.04 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `Dockerfiles/Dockerfile.feast-operator.konflux` | registry.access.redhat.com/ubi9/ubi-minimal@sha256:e1c4703364c5cb58f5462575dc90345bcd934ddc45e6c32f9c162f2b5617681c | 2 | 65532:65532 |  | multi-arch |  |  |
| `Dockerfiles/Dockerfile.feature-server.konflux` | registry.access.redhat.com/ubi9/python-312-minimal:1 | 2 | 1001 |  |  |  |  |
| `go/infra/docker/feature-server/Dockerfile` | golang:1.24.12 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `infra/feast-operator/Dockerfile` | registry.access.redhat.com/ubi9/ubi-minimal:9.5 | 2 | 65532:65532 |  | multi-arch |  |  |
| `java/infra/docker/feature-server/Dockerfile` | amazoncorretto:11 | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `java/infra/docker/feature-server/Dockerfile.dev` | openjdk:11-jre | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `sdk/python/feast/infra/compute_engines/aws_lambda/Dockerfile` | public.ecr.aws/lambda/python:3.9 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `sdk/python/feast/infra/compute_engines/kubernetes/Dockerfile` | debian:11-slim | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `sdk/python/feast/infra/feature_servers/multicloud/Dockerfile` | registry.access.redhat.com/ubi9/python-312-minimal:1 | 1 | 1001 |  |  |  |  |
| `sdk/python/feast/infra/feature_servers/multicloud/Dockerfile.dev` | registry.access.redhat.com/ubi9/python-312-minimal:1 | 1 | 1001 |  |  |  |  |
| `sdk/python/feast/infra/feature_servers/multicloud/offline/Dockerfile.binary` | yarn-builder:latest | 1 |  |  |  |  | Unpinned base image: yarn-builder:latest; No USER directive found (defaults to root) |
| `sdk/python/feast/infra/feature_servers/multicloud/offline/Dockerfile.binary.release` | registry.access.redhat.com/ubi9/python-312-minimal:1 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `sdk/python/feast/infra/feature_servers/multicloud/offline/Dockerfile.builder.yarn` | registry.access.redhat.com/ubi9/python-312-minimal:1 | 1 | 1001 |  |  |  |  |
| `sdk/python/feast/infra/feature_servers/multicloud/offline/Dockerfile.builder.yum` | registry.access.redhat.com/ubi9/python-312-minimal:1 | 1 | 1001 |  |  |  |  |
| `sdk/python/feast/infra/feature_servers/multicloud/offline/Dockerfile.sdist` | yarn-builder:latest | 1 | 1001 |  |  |  | Unpinned base image: yarn-builder:latest |
| `sdk/python/feast/infra/feature_servers/multicloud/offline/Dockerfile.sdist.release` | registry.access.redhat.com/ubi9/python-312-minimal:1 | 1 | 1001 |  |  |  |  |
| `sdk/python/feast/infra/transformation_servers/Dockerfile` | python:3.11-slim | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `ui/docker/Dockerfile` | node:17.9.0-slim | 1 |  |  |  |  | No USER directive found (defaults to root) |

