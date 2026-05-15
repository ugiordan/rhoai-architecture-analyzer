# data-science-pipelines-operator: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

### Secrets Referenced

| Name | Type | Referenced By |
|------|------|---------------|
| ds-pipeline-db-test | Opaque | deployment/mariadb |
| mariadb-certs | Opaque | deployment/mariadb |
| minio | Opaque | deployment/minio |
| minio-certs | Opaque | deployment/minio |

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| data-science-pipelines-operator-controller-manager | manager | ? | ? | ? | [`kustomize:config/overlays/odh`](https://github.com/opendatahub-io/data-science-pipelines-operator/blob/ba2d887a412d31e2f0afcebfad7fc71de3ac6521/kustomize:config/overlays/odh) |
| mariadb | mariadb | ? | ? | ? | [`.github/resources/mariadb/deployment.yaml`](https://github.com/opendatahub-io/data-science-pipelines-operator/blob/ba2d887a412d31e2f0afcebfad7fc71de3ac6521/.github/resources/mariadb/deployment.yaml) |
| minio | minio | ? | ? | ? | [`.github/resources/minio/deployment.yaml`](https://github.com/opendatahub-io/data-science-pipelines-operator/blob/ba2d887a412d31e2f0afcebfad7fc71de3ac6521/.github/resources/minio/deployment.yaml) |
| the-deployment | the-container | ? | ? | ? | [`.gomod-cache/k8s.io/cli-runtime@v0.35.3/artifacts/kustomization/deployment.yaml`](https://github.com/opendatahub-io/data-science-pipelines-operator/blob/ba2d887a412d31e2f0afcebfad7fc71de3ac6521/.gomod-cache/k8s.io/cli-runtime@v0.35.3/artifacts/kustomization/deployment.yaml) |
| the-deployment | the-container | ? | ? | ? | [`.gopath-loader/pkg/mod/k8s.io/cli-runtime@v0.35.3/artifacts/kustomization/deployment.yaml`](https://github.com/opendatahub-io/data-science-pipelines-operator/blob/ba2d887a412d31e2f0afcebfad7fc71de3ac6521/.gopath-loader/pkg/mod/k8s.io/cli-runtime@v0.35.3/artifacts/kustomization/deployment.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `.github/build/Dockerfile` | ${CI_BASE} | 2 | root |  |  |  | Unpinned base image: ${CI_BASE}; Unpinned base image: ${CI_BASE}; Container runs as root user |
| `.github/scripts/python_package_upload/Dockerfile` | docker.io/python:3.9 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/cli@v29.3.0+incompatible/Dockerfile` | scratch | 23 |  |  | multi-arch |  | Unpinned base image: build-base-alpine; Unpinned base image: build-base-debian; Unpinned base image: build-base-${BASE_VARIANT}; Unpinned base image: build-base-${BASE_VARIANT}; Unpinned base image: build-${BASE_VARIANT}; Unpinned base image: build-${BASE_VARIANT}; Unpinned base image: scratch; Unpinned base image: build-${BASE_VARIANT}; Unpinned base image: build-base-alpine; Unpinned base image: build-base-debian; Unpinned base image: e2e-base-${BASE_VARIANT}; Unpinned base image: build-base-${BASE_VARIANT}; Unpinned base image: scratch; Unpinned base image: scratch; Unpinned base image: scratch; Unpinned base image: scratch; Unpinned base image: bin-image-${TARGETOS}; Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/cli@v29.3.0+incompatible/dockerfiles/Dockerfile.authors` | gen | 3 |  |  |  |  | Unpinned base image: scratch; Unpinned base image: gen; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/cli@v29.3.0+incompatible/dockerfiles/Dockerfile.dev` | golang | 6 |  |  |  |  | Unpinned base image: golang; Unpinned base image: golang; Unpinned base image: golang; Unpinned base image: golang; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/cli@v29.3.0+incompatible/dockerfiles/Dockerfile.lint` | golang:${GO_VERSION}-alpine${ALPINE_VERSION} | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/cli@v29.3.0+incompatible/dockerfiles/Dockerfile.shellcheck` | koalaman/shellcheck-alpine:v0.7.1 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/cli@v29.3.0+incompatible/dockerfiles/Dockerfile.vendor` | base | 6 |  |  |  |  | Unpinned base image: base; Unpinned base image: scratch; Unpinned base image: vendored; Unpinned base image: base; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/cli@v29.3.0+incompatible/man/Dockerfile.5.md` | busybox | 7 | [user |  |  |  | Unpinned base image: image; Unpinned base image: ubuntu; Unpinned base image: ubuntu; Unpinned base image: ubuntu; Unpinned base image: ubuntu; Unpinned base image: busybox; Unpinned base image: busybox |
| `.gomod-cache/github.com/docker/distribution@v2.8.3+incompatible/Dockerfile` | alpine:${ALPINE_VERSION} | 8 |  |  | multi-arch |  | Unpinned base image: base; Unpinned base image: base; Unpinned base image: scratch; Unpinned base image: base; Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/distribution@v2.8.3+incompatible/contrib/compose/nginx/Dockerfile` | nginx:1.7 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/distribution@v2.8.3+incompatible/contrib/docker-integration/Dockerfile` | distribution/golem:0.1 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/distribution@v2.8.3+incompatible/contrib/docker-integration/nginx/Dockerfile` | nginx:1.9 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/distribution@v2.8.3+incompatible/contrib/docker-integration/tokenserver-oauth/Dockerfile` | dmcgowan/token-server@sha256:5a6f76d3086cdf63249c77b521108387b49d85a30c5e1c4fe82fdf5ae3b76ba7 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/distribution@v2.8.3+incompatible/contrib/docker-integration/tokenserver/Dockerfile` | dmcgowan/token-server@sha256:0eab50ebdff5b6b95b3addf4edbd8bd2f5b940f27b41b43c94afdf05863a81af | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/distribution@v2.8.3+incompatible/project/dev-image/Dockerfile` | ubuntu:14.04 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/docker-credential-helpers@v0.9.3/Dockerfile` | binaries | 20 |  |  | multi-arch |  | Unpinned base image: gobase; Unpinned base image: scratch; Unpinned base image: vendored; Unpinned base image: gobase; Unpinned base image: gobase; Unpinned base image: base; Unpinned base image: scratch; Unpinned base image: gobase; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; Unpinned base image: build-$TARGETOS; Unpinned base image: scratch; Unpinned base image: alpine; Unpinned base image: scratch; Unpinned base image: binaries; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/docker-credential-helpers@v0.9.3/deb/Dockerfile` | ${DISTRO}:${SUITE} | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/google/go-containerregistry@v0.21.3/cmd/gcrane/Dockerfile` | gcr.io/distroless/static-debian12:nonroot | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/openshift/api@v0.0.0-20260331162130-f7b3bd900c75/Dockerfile.ocp` | registry.ci.openshift.org/ocp/4.22:base-rhel9 | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/pelletier/go-toml/v2@v2.3.0/Dockerfile` | scratch | 1 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/prometheus/client_golang@v1.23.2/Dockerfile` | quay.io/prometheus/busybox:latest | 2 |  |  |  |  | Unpinned base image: quay.io/prometheus/busybox:latest; No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/toolchain@v0.0.1-go1.25.7.linux-amd64/src/crypto/internal/boring/Dockerfile` | $ubuntu:focal | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/toolchain@v0.0.1-go1.25.7.linux-amd64/src/crypto/internal/fips140/nistec/fiat/Dockerfile` | coqorg/coq:8.13.2 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/x/net@v0.52.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/x/sys@v0.42.0/unix/linux/Dockerfile` | ubuntu:25.04 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/k8s.io/apiextensions-apiserver@v0.35.3/artifacts/simple-image/Dockerfile` | gcr.io/distroless/base-debian10:latest | 1 |  |  |  |  | Unpinned base image: gcr.io/distroless/base-debian10:latest; No USER directive found (defaults to root) |
| `.gomod-cache/sigs.k8s.io/kustomize/kyaml@v0.21.1/fn/framework/example/Dockerfile` | alpine:latest | 2 |  |  |  |  | Unpinned base image: alpine:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/cli@v29.3.0+incompatible/Dockerfile` | scratch | 23 |  |  | multi-arch |  | Unpinned base image: build-base-alpine; Unpinned base image: build-base-debian; Unpinned base image: build-base-${BASE_VARIANT}; Unpinned base image: build-base-${BASE_VARIANT}; Unpinned base image: build-${BASE_VARIANT}; Unpinned base image: build-${BASE_VARIANT}; Unpinned base image: scratch; Unpinned base image: build-${BASE_VARIANT}; Unpinned base image: build-base-alpine; Unpinned base image: build-base-debian; Unpinned base image: e2e-base-${BASE_VARIANT}; Unpinned base image: build-base-${BASE_VARIANT}; Unpinned base image: scratch; Unpinned base image: scratch; Unpinned base image: scratch; Unpinned base image: scratch; Unpinned base image: bin-image-${TARGETOS}; Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/cli@v29.3.0+incompatible/dockerfiles/Dockerfile.authors` | gen | 3 |  |  |  |  | Unpinned base image: scratch; Unpinned base image: gen; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/cli@v29.3.0+incompatible/dockerfiles/Dockerfile.dev` | golang | 6 |  |  |  |  | Unpinned base image: golang; Unpinned base image: golang; Unpinned base image: golang; Unpinned base image: golang; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/cli@v29.3.0+incompatible/dockerfiles/Dockerfile.lint` | golang:${GO_VERSION}-alpine${ALPINE_VERSION} | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/cli@v29.3.0+incompatible/dockerfiles/Dockerfile.shellcheck` | koalaman/shellcheck-alpine:v0.7.1 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/cli@v29.3.0+incompatible/dockerfiles/Dockerfile.vendor` | base | 6 |  |  |  |  | Unpinned base image: base; Unpinned base image: scratch; Unpinned base image: vendored; Unpinned base image: base; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/cli@v29.3.0+incompatible/man/Dockerfile.5.md` | busybox | 7 | [user |  |  |  | Unpinned base image: image; Unpinned base image: ubuntu; Unpinned base image: ubuntu; Unpinned base image: ubuntu; Unpinned base image: ubuntu; Unpinned base image: busybox; Unpinned base image: busybox |
| `.gopath-loader/pkg/mod/github.com/docker/distribution@v2.8.3+incompatible/Dockerfile` | alpine:${ALPINE_VERSION} | 8 |  |  | multi-arch |  | Unpinned base image: base; Unpinned base image: base; Unpinned base image: scratch; Unpinned base image: base; Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/distribution@v2.8.3+incompatible/contrib/compose/nginx/Dockerfile` | nginx:1.7 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/distribution@v2.8.3+incompatible/contrib/docker-integration/Dockerfile` | distribution/golem:0.1 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/distribution@v2.8.3+incompatible/contrib/docker-integration/nginx/Dockerfile` | nginx:1.9 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/distribution@v2.8.3+incompatible/contrib/docker-integration/tokenserver-oauth/Dockerfile` | dmcgowan/token-server@sha256:5a6f76d3086cdf63249c77b521108387b49d85a30c5e1c4fe82fdf5ae3b76ba7 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/distribution@v2.8.3+incompatible/contrib/docker-integration/tokenserver/Dockerfile` | dmcgowan/token-server@sha256:0eab50ebdff5b6b95b3addf4edbd8bd2f5b940f27b41b43c94afdf05863a81af | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/distribution@v2.8.3+incompatible/project/dev-image/Dockerfile` | ubuntu:14.04 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/docker-credential-helpers@v0.9.3/Dockerfile` | binaries | 20 |  |  | multi-arch |  | Unpinned base image: gobase; Unpinned base image: scratch; Unpinned base image: vendored; Unpinned base image: gobase; Unpinned base image: gobase; Unpinned base image: base; Unpinned base image: scratch; Unpinned base image: gobase; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; Unpinned base image: build-$TARGETOS; Unpinned base image: scratch; Unpinned base image: alpine; Unpinned base image: scratch; Unpinned base image: binaries; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/docker-credential-helpers@v0.9.3/deb/Dockerfile` | ${DISTRO}:${SUITE} | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/google/go-containerregistry@v0.21.3/cmd/gcrane/Dockerfile` | gcr.io/distroless/static-debian12:nonroot | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/openshift/api@v0.0.0-20260331162130-f7b3bd900c75/Dockerfile.ocp` | registry.ci.openshift.org/ocp/4.22:base-rhel9 | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/pelletier/go-toml/v2@v2.3.0/Dockerfile` | scratch | 1 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/prometheus/client_golang@v1.23.2/Dockerfile` | quay.io/prometheus/busybox:latest | 2 |  |  |  |  | Unpinned base image: quay.io/prometheus/busybox:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.7.linux-amd64/src/crypto/internal/boring/Dockerfile` | $ubuntu:focal | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.7.linux-amd64/src/crypto/internal/fips140/nistec/fiat/Dockerfile` | coqorg/coq:8.13.2 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/net@v0.52.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/sys@v0.42.0/unix/linux/Dockerfile` | ubuntu:25.04 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/k8s.io/apiextensions-apiserver@v0.35.3/artifacts/simple-image/Dockerfile` | gcr.io/distroless/base-debian10:latest | 1 |  |  |  |  | Unpinned base image: gcr.io/distroless/base-debian10:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/sigs.k8s.io/kustomize/kyaml@v0.21.1/fn/framework/example/Dockerfile` | alpine:latest | 2 |  |  |  |  | Unpinned base image: alpine:latest; No USER directive found (defaults to root) |
| `Dockerfile` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 2 | ${USER}:${USER} |  | multi-arch | yes | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |

