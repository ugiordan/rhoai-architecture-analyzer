# model-registry: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

### Secrets Referenced

| Name | Type | Referenced By |
|------|------|---------------|
| minio-secret | Opaque | deployment/minio |
| model-catalog-hf-api-key | Opaque | deployment/model-catalog-server |
| model-catalog-postgres | Opaque | deployment/model-catalog-server |

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| controller-manager | manager | ? | ? | ? | [`manifests/kustomize/options/controller/manager/manager.yaml`](https://github.com/kubeflow/model-registry/blob/2bccb683c2f077c6d39db5588d7cb908885ac975/manifests/kustomize/options/controller/manager/manager.yaml) |
| minio | minio | ? | ? | ? | [`scripts/manifests/minio/deployment.yaml`](https://github.com/kubeflow/model-registry/blob/2bccb683c2f077c6d39db5588d7cb908885ac975/scripts/manifests/minio/deployment.yaml) |
| model-catalog-server | catalog | ? | ? | ? | [`manifests/kustomize/options/catalog/base/deployment.yaml`](https://github.com/kubeflow/model-registry/blob/2bccb683c2f077c6d39db5588d7cb908885ac975/manifests/kustomize/options/catalog/base/deployment.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `.gomod-cache/github.com/golang-migrate/migrate/v4@v4.19.1/Dockerfile` | alpine:3.21 | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/golang-migrate/migrate/v4@v4.19.1/Dockerfile.circleci` | $DOCKER_IMAGE | 1 |  |  |  |  | Unpinned base image: $DOCKER_IMAGE; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/golang-migrate/migrate/v4@v4.19.1/Dockerfile.github-actions` | alpine:3.19 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/jackc/pgx/v5@v5.9.2/.devcontainer/Dockerfile` | mcr.microsoft.com/devcontainers/go:2-1.26-trixie | 1 | vscode |  |  |  |  |
| `.gomod-cache/github.com/moby/moby/api@v1.54.1/Dockerfile` | base | 3 |  |  | multi-arch |  | Unpinned base image: base; Unpinned base image: base; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/pelletier/go-toml/v2@v2.2.4/Dockerfile` | scratch | 1 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gomod-cache/go.mongodb.org/mongo-driver@v1.17.4/Dockerfile` | ubuntu:20.04 | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/toolchain@v0.0.1-go1.25.7.linux-amd64/src/crypto/internal/boring/Dockerfile` | $ubuntu:focal | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/toolchain@v0.0.1-go1.25.7.linux-amd64/src/crypto/internal/fips140/nistec/fiat/Dockerfile` | coqorg/coq:8.13.2 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/x/sys@v0.42.0/unix/linux/Dockerfile` | ubuntu:25.04 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/golang-migrate/migrate/v4@v4.19.1/Dockerfile` | alpine:3.21 | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/golang-migrate/migrate/v4@v4.19.1/Dockerfile.circleci` | $DOCKER_IMAGE | 1 |  |  |  |  | Unpinned base image: $DOCKER_IMAGE; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/golang-migrate/migrate/v4@v4.19.1/Dockerfile.github-actions` | alpine:3.19 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/jackc/pgx/v5@v5.9.2/.devcontainer/Dockerfile` | mcr.microsoft.com/devcontainers/go:2-1.26-trixie | 1 | vscode |  |  |  |  |
| `.gopath-loader/pkg/mod/github.com/moby/moby/api@v1.54.1/Dockerfile` | base | 3 |  |  | multi-arch |  | Unpinned base image: base; Unpinned base image: base; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/pelletier/go-toml/v2@v2.2.4/Dockerfile` | scratch | 1 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/go.mongodb.org/mongo-driver@v1.17.4/Dockerfile` | ubuntu:20.04 | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.7.linux-amd64/src/crypto/internal/boring/Dockerfile` | $ubuntu:focal | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.7.linux-amd64/src/crypto/internal/fips140/nistec/fiat/Dockerfile` | coqorg/coq:8.13.2 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/sys@v0.42.0/unix/linux/Dockerfile` | ubuntu:25.04 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `Dockerfile` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 3 | 65532:65532 |  | multi-arch |  | Unpinned base image: common; Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |
| `Dockerfile.konflux` | registry.access.redhat.com/ubi9/ubi-minimal@sha256:12db9874bd753eb98b1ab3d840e75de5d6842ac0604fbd68c012adefe97140be | 2 | 1001 |  |  |  |  |
| `Dockerfile.odh` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 2 | 65532:65532 |  |  |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |
| `Dockerfile.testops` | registry.access.redhat.com/ubi9/python-312 | 1 | odh |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/python-312 |
| `clients/ui/Dockerfile` | ${DISTROLESS_BASE_IMAGE} | 3 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${NODE_BASE_IMAGE}; Unpinned base image: ${GOLANG_BASE_IMAGE}; Unpinned base image: ${DISTROLESS_BASE_IMAGE} |
| `clients/ui/Dockerfile.standalone` | release | 4 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${NODE_BASE_IMAGE}; Unpinned base image: ${GOLANG_BASE_IMAGE}; Unpinned base image: ${DISTROLESS_BASE_IMAGE}; Unpinned base image: release |
| `cmd/controller/Dockerfile.controller` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 2 | 65532:65532 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |
| `cmd/csi/Dockerfile.csi` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 3 | 65532:65532 |  | multi-arch |  | Unpinned base image: common; Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |
| `jobs/async-upload/Dockerfile` | registry.access.redhat.com/ubi9/python-312-minimal | 2 | 1000 |  |  |  | Unpinned base image: registry.access.redhat.com/ubi9/python-312-minimal |
| `jobs/async-upload/Dockerfile.konflux` | registry.access.redhat.com/ubi9/python-312-minimal@sha256:445709dc989a00efecaa10244f5c502ecb1604b5b1fbb8fe21e9e9ffb3e36254 | 1 | 1000 |  | multi-arch |  |  |

