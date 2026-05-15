# mlflow-operator: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

### Secrets Referenced

| Name | Type | Referenced By |
|------|------|---------------|
| controller-manager-metrics-tls | Opaque | deployment/controller-manager |
| postgres-secret | Opaque | deployment/postgres-deployment |

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| mlflow-operator-controller-manager | manager | ? | true | ? | [`kustomize:config/overlays/odh`](https://github.com/opendatahub-io/mlflow-operator/blob/e68add073b5755c9417d37dc4cd2f5e1cda6dbab/kustomize:config/overlays/odh) |
| postgres-deployment | postgres | ? | ? | ? | [`config/postgres/base/deployment.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/e68add073b5755c9417d37dc4cd2f5e1cda6dbab/config/postgres/base/deployment.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `.gomod-cache/github.com/grpc-ecosystem/grpc-gateway/v2@v2.27.7/.github/Dockerfile` | golang:1.25.6 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/openshift/api@v0.0.0-20260209232644-126cbbe24427/Dockerfile.ocp` | registry.ci.openshift.org/ocp/4.22:base-rhel9 | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/pelletier/go-toml/v2@v2.2.4/Dockerfile` | scratch | 1 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/prometheus/client_golang@v1.23.0/Dockerfile` | quay.io/prometheus/busybox:latest | 2 |  |  |  |  | Unpinned base image: quay.io/prometheus/busybox:latest; No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/x/net@v0.49.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/x/sys@v0.40.0/unix/linux/Dockerfile` | ubuntu:25.04 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/k8s.io/apiextensions-apiserver@v0.34.3/artifacts/simple-image/Dockerfile` | gcr.io/distroless/base-debian10:latest | 1 |  |  |  |  | Unpinned base image: gcr.io/distroless/base-debian10:latest; No USER directive found (defaults to root) |
| `.gomod-cache/sigs.k8s.io/gateway-api@v1.4.1/docker/Dockerfile.echo-advanced` | gcr.io/istio-release/app:1.21.2 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/sigs.k8s.io/gateway-api@v1.4.1/docker/Dockerfile.echo-basic` | gcr.io/distroless/static:nonroot | 2 | nonroot:nonroot |  |  |  |  |
| `.gopath-loader/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.27.7/.github/Dockerfile` | golang:1.25.6 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/openshift/api@v0.0.0-20260209232644-126cbbe24427/Dockerfile.ocp` | registry.ci.openshift.org/ocp/4.22:base-rhel9 | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/pelletier/go-toml/v2@v2.2.4/Dockerfile` | scratch | 1 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/prometheus/client_golang@v1.23.0/Dockerfile` | quay.io/prometheus/busybox:latest | 2 |  |  |  |  | Unpinned base image: quay.io/prometheus/busybox:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/net@v0.49.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/sys@v0.40.0/unix/linux/Dockerfile` | ubuntu:25.04 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/k8s.io/apiextensions-apiserver@v0.34.3/artifacts/simple-image/Dockerfile` | gcr.io/distroless/base-debian10:latest | 1 |  |  |  |  | Unpinned base image: gcr.io/distroless/base-debian10:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api@v1.4.1/docker/Dockerfile.echo-advanced` | gcr.io/istio-release/app:1.21.2 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api@v1.4.1/docker/Dockerfile.echo-basic` | gcr.io/distroless/static:nonroot | 2 | nonroot:nonroot |  |  |  |  |
| `Dockerfile` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 2 | 1001 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |
| `Dockerfile.konflux` | registry.access.redhat.com/ubi9/ubi-minimal@sha256:12db9874bd753eb98b1ab3d840e75de5d6842ac0604fbd68c012adefe97140be | 2 | 1001 |  | multi-arch |  |  |
| `mlflow-tests/images/Dockerfile.konflux` | registry.access.redhat.com/ubi9/ubi-minimal:9.7 | 1 | 1001 |  |  |  |  |

