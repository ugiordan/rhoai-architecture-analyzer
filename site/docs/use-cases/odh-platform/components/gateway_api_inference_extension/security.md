# gateway-api-inference-extension: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `.gomod-cache/golang.org/toolchain@v0.0.1-go1.25.0.linux-amd64/src/crypto/internal/boring/Dockerfile` | $ubuntu:focal | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/toolchain@v0.0.1-go1.25.0.linux-amd64/src/crypto/internal/fips140/nistec/fiat/Dockerfile` | coqorg/coq:8.13.2 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/!shopify/toxiproxy@v2.1.4+incompatible/Dockerfile` | alpine | 1 |  |  |  |  | Unpinned base image: alpine; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/apache/thrift@v0.17.0/build/docker/msvc2017/Dockerfile` | microsoft/dotnet-framework:4.7.1 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/apache/thrift@v0.17.0/build/docker/old/centos-7.3/Dockerfile` | centos:7.3.1611 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/apache/thrift@v0.17.0/build/docker/old/debian-jessie/Dockerfile` | buildpack-deps:jessie-scm | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/apache/thrift@v0.17.0/build/docker/old/debian-stretch/Dockerfile` | buildpack-deps:stretch-scm | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/apache/thrift@v0.17.0/build/docker/old/ubuntu-artful/Dockerfile` | buildpack-deps:artful-scm | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/apache/thrift@v0.17.0/build/docker/old/ubuntu-trusty/Dockerfile` | buildpack-deps:trusty-scm | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/apache/thrift@v0.17.0/build/docker/ubuntu-bionic/Dockerfile` | buildpack-deps:bionic-scm | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/apache/thrift@v0.17.0/build/docker/ubuntu-disco/Dockerfile` | buildpack-deps:disco-scm | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/apache/thrift@v0.17.0/build/docker/ubuntu-xenial/Dockerfile` | buildpack-deps:xenial-scm | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.27.0/awstesting/sandbox/Dockerfile.golang-tip` | buildpack-deps:buster-scm | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.27.0/awstesting/sandbox/Dockerfile.test.go1.10` | golang:1.10 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.27.0/awstesting/sandbox/Dockerfile.test.go1.11` | golang:1.11 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.27.0/awstesting/sandbox/Dockerfile.test.go1.12` | golang:1.12 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.27.0/awstesting/sandbox/Dockerfile.test.go1.13` | golang:1.13 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.27.0/awstesting/sandbox/Dockerfile.test.go1.5` | golang:1.5 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.27.0/awstesting/sandbox/Dockerfile.test.go1.5-novendorexp` | golang:1.5 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.27.0/awstesting/sandbox/Dockerfile.test.go1.6` | golang:1.6 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.27.0/awstesting/sandbox/Dockerfile.test.go1.7` | golang:1.7 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.27.0/awstesting/sandbox/Dockerfile.test.go1.8` | golang:1.8 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.27.0/awstesting/sandbox/Dockerfile.test.go1.9` | golang:1.9 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.27.0/awstesting/sandbox/Dockerfile.test.gotip` | aws-golang:tip | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/buger/jsonparser@v1.1.1/Dockerfile` | golang:1.6 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/client9/misspell@v0.3.4/Dockerfile` | golang:1.10.0-alpine | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/cpuguy83/go-md2man/v2@v2.0.6/Dockerfile` | scratch | 2 |  |  | multi-arch |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/creack/pty@v1.1.18/Dockerfile.golang` | golang:${GOVERSION} | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/creack/pty@v1.1.18/Dockerfile.riscv` | golang:1.13 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/envoyproxy/go-control-plane@v0.14.0/Dockerfile.ci` | golang:1.23@sha256:77a21b3e354c03e9f66b13bc39f4f0db8085c70f8414406af66b29c6d6c4dd85 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/envoyproxy/protoc-gen-validate@v1.3.0/Dockerfile` | ubuntu:focal | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.28.0/.github/Dockerfile` | golang:1.26.0 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/.circleci/Dockerfile` | golang:1.15.3 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/prometheus/client_golang@v1.23.2/Dockerfile` | quay.io/prometheus/busybox:latest | 2 |  |  |  |  | Unpinned base image: quay.io/prometheus/busybox:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/prometheus/prometheus@v0.310.0/Dockerfile` | quay.io/prometheus/busybox-${OS}-${ARCH}:latest | 1 | nobody |  |  |  | Unpinned base image: quay.io/prometheus/busybox-${OS}-${ARCH}:latest |
| `.gopath-loader/pkg/mod/github.com/prometheus/prometheus@v0.310.0/Dockerfile.distroless` | gcr.io/distroless/static-debian13:nonroot-${DISTROLESS_ARCH} | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.0.linux-amd64/src/crypto/internal/boring/Dockerfile` | $ubuntu:focal | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.0.linux-amd64/src/crypto/internal/fips140/nistec/fiat/Dockerfile` | coqorg/coq:8.13.2 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/net@v0.52.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/sys@v0.42.0/unix/linux/Dockerfile` | ubuntu:25.04 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/k8s.io/apiextensions-apiserver@v0.35.5/artifacts/simple-image/Dockerfile` | gcr.io/distroless/base-debian10:latest | 1 |  |  |  |  | Unpinned base image: gcr.io/distroless/base-debian10:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/sigs.k8s.io/kustomize/kyaml@v0.21.1/fn/framework/example/Dockerfile` | alpine:latest | 2 |  |  |  |  | Unpinned base image: alpine:latest; No USER directive found (defaults to root) |
| `Dockerfile` | ${BASE_IMAGE} | 2 |  |  |  |  | Unpinned base image: ${BUILDER_IMAGE}; Unpinned base image: ${BASE_IMAGE}; No USER directive found (defaults to root) |

