# argo-workflows: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| the-deployment | the-container | ? | ? | ? | [`.gomod-cache/k8s.io/cli-runtime@v0.32.2/artifacts/kustomization/deployment.yaml`](https://github.com/argoproj/argo-workflows/blob/003ed2b35a398772211441cb7c866c51f6f87e2d/.gomod-cache/k8s.io/cli-runtime@v0.32.2/artifacts/kustomization/deployment.yaml) |
| the-deployment | the-container | ? | ? | ? | [`.gopath-loader/pkg/mod/k8s.io/cli-runtime@v0.32.2/artifacts/kustomization/deployment.yaml`](https://github.com/argoproj/argo-workflows/blob/003ed2b35a398772211441cb7c866c51f6f87e2d/.gopath-loader/pkg/mod/k8s.io/cli-runtime@v0.32.2/artifacts/kustomization/deployment.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `.gomod-cache/cloud.google.com/go@v0.119.0/.devcontainer/Dockerfile` | mcr.microsoft.com/devcontainers/go:${VARIANT} | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/argoproj/argo-events@v1.9.6/Dockerfile` | scratch | 2 |  |  | multi-arch |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/argoproj/argo-events@v1.9.6/third_party/nats-streaming-docker/amd64/Dockerfile` | scratch | 1 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/argoproj/argo-events@v1.9.6/third_party/prometheus-nats-exporter-docker/amd64/Dockerfile` | scratch | 2 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go-v2@v1.36.3/internal/awstesting/sandbox/Dockerfile.golang-tip` | buildpack-deps:buster-scm | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go-v2@v1.36.3/internal/awstesting/sandbox/Dockerfile.test.gotip` | aws-golang:tip | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go-v2@v1.36.3/internal/awstesting/sandbox/Dockerfile.test.goversion` | golang:${GO_VERSION} | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/cpuguy83/go-md2man/v2@v2.0.6/Dockerfile` | scratch | 2 |  |  | multi-arch |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/creack/pty@v1.1.21/Dockerfile.golang` | golang:${GOVERSION} | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/cli@v27.1.1+incompatible/Dockerfile` | scratch | 23 |  |  | multi-arch |  | Unpinned base image: build-base-alpine; Unpinned base image: build-base-debian; Unpinned base image: build-base-${BASE_VARIANT}; Unpinned base image: build-base-${BASE_VARIANT}; Unpinned base image: build-${BASE_VARIANT}; Unpinned base image: build-${BASE_VARIANT}; Unpinned base image: scratch; Unpinned base image: build-${BASE_VARIANT}; Unpinned base image: build-base-alpine; Unpinned base image: build-base-debian; Unpinned base image: e2e-base-${BASE_VARIANT}; Unpinned base image: build-base-${BASE_VARIANT}; Unpinned base image: scratch; Unpinned base image: scratch; Unpinned base image: scratch; Unpinned base image: scratch; Unpinned base image: bin-image-${TARGETOS}; Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/cli@v27.1.1+incompatible/dockerfiles/Dockerfile.authors` | gen | 3 |  |  |  |  | Unpinned base image: scratch; Unpinned base image: gen; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/cli@v27.1.1+incompatible/dockerfiles/Dockerfile.dev` | golang | 6 |  |  |  |  | Unpinned base image: golang; Unpinned base image: golang; Unpinned base image: golang; Unpinned base image: golang; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/cli@v27.1.1+incompatible/dockerfiles/Dockerfile.lint` | golang:${GO_VERSION}-alpine${ALPINE_VERSION} | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/cli@v27.1.1+incompatible/dockerfiles/Dockerfile.shellcheck` | koalaman/shellcheck-alpine:v0.7.1 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/cli@v27.1.1+incompatible/dockerfiles/Dockerfile.vendor` | base | 6 |  |  |  |  | Unpinned base image: base; Unpinned base image: scratch; Unpinned base image: vendored; Unpinned base image: base; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/distribution@v2.8.2+incompatible/Dockerfile` | alpine:${ALPINE_VERSION} | 8 |  |  | multi-arch |  | Unpinned base image: base; Unpinned base image: base; Unpinned base image: scratch; Unpinned base image: base; Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/distribution@v2.8.2+incompatible/contrib/compose/nginx/Dockerfile` | nginx:1.7 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/distribution@v2.8.2+incompatible/contrib/docker-integration/Dockerfile` | distribution/golem:0.1 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/distribution@v2.8.2+incompatible/contrib/docker-integration/nginx/Dockerfile` | nginx:1.9 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/distribution@v2.8.2+incompatible/contrib/docker-integration/tokenserver-oauth/Dockerfile` | dmcgowan/token-server@sha256:5a6f76d3086cdf63249c77b521108387b49d85a30c5e1c4fe82fdf5ae3b76ba7 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/distribution@v2.8.2+incompatible/contrib/docker-integration/tokenserver/Dockerfile` | dmcgowan/token-server@sha256:0eab50ebdff5b6b95b3addf4edbd8bd2f5b940f27b41b43c94afdf05863a81af | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/distribution@v2.8.2+incompatible/project/dev-image/Dockerfile` | ubuntu:14.04 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/docker-credential-helpers@v0.8.2/Dockerfile` | binaries | 20 |  |  | multi-arch |  | Unpinned base image: gobase; Unpinned base image: scratch; Unpinned base image: vendored; Unpinned base image: gobase; Unpinned base image: gobase; Unpinned base image: base; Unpinned base image: scratch; Unpinned base image: gobase; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; Unpinned base image: build-$TARGETOS; Unpinned base image: scratch; Unpinned base image: alpine; Unpinned base image: scratch; Unpinned base image: binaries; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/docker-credential-helpers@v0.8.2/deb/Dockerfile` | ${DISTRO}:${SUITE} | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/envoyproxy/protoc-gen-validate@v1.2.1/Dockerfile` | ubuntu:focal | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/grpc-ecosystem/grpc-gateway/v2@v2.26.3/.github/Dockerfile` | golang:1.24.0 | 1 | vscode |  |  |  |  |
| `.gomod-cache/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/.circleci/Dockerfile` | golang:1.15.3 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/itchyny/gojq@v0.12.17/Dockerfile` | gcr.io/distroless/static:debug | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/pelletier/go-toml/v2@v2.2.3/Dockerfile` | scratch | 1 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/pjbgf/sha1cd@v0.3.2/Dockerfile.arm` | golang:1.23@sha256:51a6466e8dbf3e00e422eb0f7a97ac450b2d57b33617bbe8d2ee0bddcd9d0d37 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/pjbgf/sha1cd@v0.3.2/Dockerfile.arm64` | golang:1.23@sha256:51a6466e8dbf3e00e422eb0f7a97ac450b2d57b33617bbe8d2ee0bddcd9d0d37 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/prometheus/client_golang@v1.21.1/Dockerfile` | quay.io/prometheus/busybox:latest | 2 |  |  |  |  | Unpinned base image: quay.io/prometheus/busybox:latest; No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/x/net@v0.38.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/x/sys@v0.31.0/unix/linux/Dockerfile` | ubuntu:24.10 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/sigs.k8s.io/kustomize/kyaml@v0.18.1/fn/framework/example/Dockerfile` | alpine:latest | 2 |  |  |  |  | Unpinned base image: alpine:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/cloud.google.com/go@v0.119.0/.devcontainer/Dockerfile` | mcr.microsoft.com/devcontainers/go:${VARIANT} | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/argoproj/argo-events@v1.9.6/Dockerfile` | scratch | 2 |  |  | multi-arch |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/argoproj/argo-events@v1.9.6/third_party/nats-streaming-docker/amd64/Dockerfile` | scratch | 1 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/argoproj/argo-events@v1.9.6/third_party/prometheus-nats-exporter-docker/amd64/Dockerfile` | scratch | 2 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go-v2@v1.36.3/internal/awstesting/sandbox/Dockerfile.golang-tip` | buildpack-deps:buster-scm | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go-v2@v1.36.3/internal/awstesting/sandbox/Dockerfile.test.gotip` | aws-golang:tip | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go-v2@v1.36.3/internal/awstesting/sandbox/Dockerfile.test.goversion` | golang:${GO_VERSION} | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/cpuguy83/go-md2man/v2@v2.0.6/Dockerfile` | scratch | 2 |  |  | multi-arch |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/creack/pty@v1.1.21/Dockerfile.golang` | golang:${GOVERSION} | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/cli@v27.1.1+incompatible/Dockerfile` | scratch | 23 |  |  | multi-arch |  | Unpinned base image: build-base-alpine; Unpinned base image: build-base-debian; Unpinned base image: build-base-${BASE_VARIANT}; Unpinned base image: build-base-${BASE_VARIANT}; Unpinned base image: build-${BASE_VARIANT}; Unpinned base image: build-${BASE_VARIANT}; Unpinned base image: scratch; Unpinned base image: build-${BASE_VARIANT}; Unpinned base image: build-base-alpine; Unpinned base image: build-base-debian; Unpinned base image: e2e-base-${BASE_VARIANT}; Unpinned base image: build-base-${BASE_VARIANT}; Unpinned base image: scratch; Unpinned base image: scratch; Unpinned base image: scratch; Unpinned base image: scratch; Unpinned base image: bin-image-${TARGETOS}; Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/cli@v27.1.1+incompatible/dockerfiles/Dockerfile.authors` | gen | 3 |  |  |  |  | Unpinned base image: scratch; Unpinned base image: gen; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/cli@v27.1.1+incompatible/dockerfiles/Dockerfile.dev` | golang | 6 |  |  |  |  | Unpinned base image: golang; Unpinned base image: golang; Unpinned base image: golang; Unpinned base image: golang; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/cli@v27.1.1+incompatible/dockerfiles/Dockerfile.lint` | golang:${GO_VERSION}-alpine${ALPINE_VERSION} | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/cli@v27.1.1+incompatible/dockerfiles/Dockerfile.shellcheck` | koalaman/shellcheck-alpine:v0.7.1 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/cli@v27.1.1+incompatible/dockerfiles/Dockerfile.vendor` | base | 6 |  |  |  |  | Unpinned base image: base; Unpinned base image: scratch; Unpinned base image: vendored; Unpinned base image: base; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/distribution@v2.8.2+incompatible/Dockerfile` | alpine:${ALPINE_VERSION} | 8 |  |  | multi-arch |  | Unpinned base image: base; Unpinned base image: base; Unpinned base image: scratch; Unpinned base image: base; Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/distribution@v2.8.2+incompatible/contrib/compose/nginx/Dockerfile` | nginx:1.7 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/distribution@v2.8.2+incompatible/contrib/docker-integration/Dockerfile` | distribution/golem:0.1 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/distribution@v2.8.2+incompatible/contrib/docker-integration/nginx/Dockerfile` | nginx:1.9 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/distribution@v2.8.2+incompatible/contrib/docker-integration/tokenserver-oauth/Dockerfile` | dmcgowan/token-server@sha256:5a6f76d3086cdf63249c77b521108387b49d85a30c5e1c4fe82fdf5ae3b76ba7 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/distribution@v2.8.2+incompatible/contrib/docker-integration/tokenserver/Dockerfile` | dmcgowan/token-server@sha256:0eab50ebdff5b6b95b3addf4edbd8bd2f5b940f27b41b43c94afdf05863a81af | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/distribution@v2.8.2+incompatible/project/dev-image/Dockerfile` | ubuntu:14.04 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/docker-credential-helpers@v0.8.2/Dockerfile` | binaries | 20 |  |  | multi-arch |  | Unpinned base image: gobase; Unpinned base image: scratch; Unpinned base image: vendored; Unpinned base image: gobase; Unpinned base image: gobase; Unpinned base image: base; Unpinned base image: scratch; Unpinned base image: gobase; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; Unpinned base image: build-$TARGETOS; Unpinned base image: scratch; Unpinned base image: alpine; Unpinned base image: scratch; Unpinned base image: binaries; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/docker-credential-helpers@v0.8.2/deb/Dockerfile` | ${DISTRO}:${SUITE} | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/envoyproxy/protoc-gen-validate@v1.2.1/Dockerfile` | ubuntu:focal | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.26.3/.github/Dockerfile` | golang:1.24.0 | 1 | vscode |  |  |  |  |
| `.gopath-loader/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/.circleci/Dockerfile` | golang:1.15.3 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/itchyny/gojq@v0.12.17/Dockerfile` | gcr.io/distroless/static:debug | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/pelletier/go-toml/v2@v2.2.3/Dockerfile` | scratch | 1 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/pjbgf/sha1cd@v0.3.2/Dockerfile.arm` | golang:1.23@sha256:51a6466e8dbf3e00e422eb0f7a97ac450b2d57b33617bbe8d2ee0bddcd9d0d37 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/pjbgf/sha1cd@v0.3.2/Dockerfile.arm64` | golang:1.23@sha256:51a6466e8dbf3e00e422eb0f7a97ac450b2d57b33617bbe8d2ee0bddcd9d0d37 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/prometheus/client_golang@v1.21.1/Dockerfile` | quay.io/prometheus/busybox:latest | 2 |  |  |  |  | Unpinned base image: quay.io/prometheus/busybox:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/net@v0.38.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/sys@v0.31.0/unix/linux/Dockerfile` | ubuntu:24.10 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/sigs.k8s.io/kustomize/kyaml@v0.18.1/fn/framework/example/Dockerfile` | alpine:latest | 2 |  |  |  |  | Unpinned base image: alpine:latest; No USER directive found (defaults to root) |
| `Dockerfile` | gcr.io/distroless/static | 10 | 8737 |  |  |  | Unpinned base image: builder; Unpinned base image: builder; Unpinned base image: builder; Unpinned base image: gcr.io/distroless/static; Unpinned base image: argoexec-base; Unpinned base image: argoexec-base; Unpinned base image: gcr.io/distroless/static; Unpinned base image: gcr.io/distroless/static |
| `Dockerfile.windows` | argoexec-base | 4 | Administrator |  |  |  | Unpinned base image: builder; Unpinned base image: argoexec-base |
| `argo-argoexec/Dockerfile.ODH` | registry.redhat.io/ubi9/ubi-minimal:9.5 | 2 | 2000 |  |  |  |  |
| `argo-argoexec/Dockerfile.konflux` | registry.redhat.io/ubi9/ubi-minimal@sha256:8d905a93f1392d4a8f7fb906bd49bf540290674b28d82de3536bb4d0898bf9d7 | 2 | 2000 |  |  |  |  |
| `argo-workflowcontroller/Dockerfile.ODH` | registry.redhat.io/ubi9/ubi-minimal:9.5 | 2 | 8737 |  |  |  |  |
| `argo-workflowcontroller/Dockerfile.konflux` | registry.redhat.io/ubi9/ubi-minimal@sha256:8d905a93f1392d4a8f7fb906bd49bf540290674b28d82de3536bb4d0898bf9d7 | 2 | 8737 |  |  |  |  |
| `rhoai/Dockerfile.argoexec` | registry.redhat.io/ubi8/ubi-minimal:latest | 2 | 2000 |  |  |  | Unpinned base image: registry.redhat.io/ubi8/ubi-minimal:latest |
| `rhoai/Dockerfile.workflowcontroller` | registry.redhat.io/ubi8/ubi-minimal:latest | 2 | 8737 |  |  |  | Unpinned base image: registry.redhat.io/ubi8/ubi-minimal:latest |

