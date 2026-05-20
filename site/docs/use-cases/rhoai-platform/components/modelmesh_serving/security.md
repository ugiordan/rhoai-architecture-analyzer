# modelmesh-serving: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

### Secrets Referenced

| Name | Type | Referenced By |
|------|------|---------------|
| kserve-webhook-server-cert | Opaque | deployment/kserve-controller-manager |
| modelmesh-webhook-server-cert | Opaque | deployment/modelmesh-controller |

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| controller-manager | kube-rbac-proxy | ? | ? | ? | [`config/default/manager_auth_proxy_patch.yaml`](https://github.com/kserve/modelmesh-serving/blob/4e16417034a9fc02561b6cdb0356d337805589b1/config/default/manager_auth_proxy_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/default/manager_auth_proxy_patch.yaml`](https://github.com/kserve/modelmesh-serving/blob/4e16417034a9fc02561b6cdb0356d337805589b1/config/default/manager_auth_proxy_patch.yaml) |
| etcd | etcd | ? | ? | ? | [`kustomize:config/overlays/odh`](https://github.com/kserve/modelmesh-serving/blob/4e16417034a9fc02561b6cdb0356d337805589b1/kustomize:config/overlays/odh) |
| kserve-controller-manager | manager | ? | ? | ? | [`.gopath-loader/pkg/mod/github.com/kserve/kserve@v0.12.0/config/default/manager_image_patch.yaml`](https://github.com/kserve/modelmesh-serving/blob/4e16417034a9fc02561b6cdb0356d337805589b1/.gopath-loader/pkg/mod/github.com/kserve/kserve@v0.12.0/config/default/manager_image_patch.yaml) |
| kserve-controller-manager | manager | ? | ? | ? | [`.gomod-cache/github.com/kserve/kserve@v0.12.0/config/default/manager_resources_patch.yaml`](https://github.com/kserve/modelmesh-serving/blob/4e16417034a9fc02561b6cdb0356d337805589b1/.gomod-cache/github.com/kserve/kserve@v0.12.0/config/default/manager_resources_patch.yaml) |
| kserve-controller-manager | manager | ? | ? | ? | [`.gomod-cache/github.com/kserve/kserve@v0.12.0/config/manager/manager.yaml`](https://github.com/kserve/modelmesh-serving/blob/4e16417034a9fc02561b6cdb0356d337805589b1/.gomod-cache/github.com/kserve/kserve@v0.12.0/config/manager/manager.yaml) |
| kserve-controller-manager | kube-rbac-proxy | ? | ? | ? | [`.gopath-loader/pkg/mod/github.com/kserve/kserve@v0.12.0/config/default/manager_auth_proxy_patch.yaml`](https://github.com/kserve/modelmesh-serving/blob/4e16417034a9fc02561b6cdb0356d337805589b1/.gopath-loader/pkg/mod/github.com/kserve/kserve@v0.12.0/config/default/manager_auth_proxy_patch.yaml) |
| kserve-controller-manager | manager | ? | ? | ? | [`.gopath-loader/pkg/mod/github.com/kserve/kserve@v0.12.0/config/default/manager_auth_proxy_patch.yaml`](https://github.com/kserve/modelmesh-serving/blob/4e16417034a9fc02561b6cdb0356d337805589b1/.gopath-loader/pkg/mod/github.com/kserve/kserve@v0.12.0/config/default/manager_auth_proxy_patch.yaml) |
| kserve-controller-manager | kube-rbac-proxy | ? | ? | ? | [`.gomod-cache/github.com/kserve/kserve@v0.12.0/config/default/manager_auth_proxy_patch.yaml`](https://github.com/kserve/modelmesh-serving/blob/4e16417034a9fc02561b6cdb0356d337805589b1/.gomod-cache/github.com/kserve/kserve@v0.12.0/config/default/manager_auth_proxy_patch.yaml) |
| kserve-controller-manager | manager | ? | ? | ? | [`.gomod-cache/github.com/kserve/kserve@v0.12.0/config/default/manager_auth_proxy_patch.yaml`](https://github.com/kserve/modelmesh-serving/blob/4e16417034a9fc02561b6cdb0356d337805589b1/.gomod-cache/github.com/kserve/kserve@v0.12.0/config/default/manager_auth_proxy_patch.yaml) |
| kserve-controller-manager | manager | ? | ? | ? | [`.gopath-loader/pkg/mod/github.com/kserve/kserve@v0.12.0/config/default/manager_prometheus_metrics_patch.yaml`](https://github.com/kserve/modelmesh-serving/blob/4e16417034a9fc02561b6cdb0356d337805589b1/.gopath-loader/pkg/mod/github.com/kserve/kserve@v0.12.0/config/default/manager_prometheus_metrics_patch.yaml) |
| kserve-controller-manager | manager | ? | ? | ? | [`.gopath-loader/pkg/mod/github.com/kserve/kserve@v0.12.0/config/default/manager_resources_patch.yaml`](https://github.com/kserve/modelmesh-serving/blob/4e16417034a9fc02561b6cdb0356d337805589b1/.gopath-loader/pkg/mod/github.com/kserve/kserve@v0.12.0/config/default/manager_resources_patch.yaml) |
| kserve-controller-manager | manager | ? | ? | ? | [`.gopath-loader/pkg/mod/github.com/kserve/kserve@v0.12.0/config/manager/manager.yaml`](https://github.com/kserve/modelmesh-serving/blob/4e16417034a9fc02561b6cdb0356d337805589b1/.gopath-loader/pkg/mod/github.com/kserve/kserve@v0.12.0/config/manager/manager.yaml) |
| kserve-controller-manager | manager | ? | ? | ? | [`.gomod-cache/github.com/kserve/kserve@v0.12.0/config/default/manager_prometheus_metrics_patch.yaml`](https://github.com/kserve/modelmesh-serving/blob/4e16417034a9fc02561b6cdb0356d337805589b1/.gomod-cache/github.com/kserve/kserve@v0.12.0/config/default/manager_prometheus_metrics_patch.yaml) |
| kserve-controller-manager | manager | ? | ? | ? | [`.gomod-cache/github.com/kserve/kserve@v0.12.0/config/default/manager_image_patch.yaml`](https://github.com/kserve/modelmesh-serving/blob/4e16417034a9fc02561b6cdb0356d337805589b1/.gomod-cache/github.com/kserve/kserve@v0.12.0/config/default/manager_image_patch.yaml) |
| modelmesh-controller | manager | ? | ? | ? | [`kustomize:config/overlays/odh`](https://github.com/kserve/modelmesh-serving/blob/4e16417034a9fc02561b6cdb0356d337805589b1/kustomize:config/overlays/odh) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `.gomod-cache/cloud.google.com/go@v0.110.10/.devcontainer/Dockerfile` | mcr.microsoft.com/devcontainers/go:${VARIANT} | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.48.0/awstesting/sandbox/Dockerfile.golang-tip` | buildpack-deps:buster-scm | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.48.0/awstesting/sandbox/Dockerfile.test.go1.10` | golang:1.10 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.48.0/awstesting/sandbox/Dockerfile.test.go1.11` | golang:1.11 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.48.0/awstesting/sandbox/Dockerfile.test.go1.12` | golang:1.12 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.48.0/awstesting/sandbox/Dockerfile.test.go1.13` | golang:1.13 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.48.0/awstesting/sandbox/Dockerfile.test.go1.14` | golang:1.14 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.48.0/awstesting/sandbox/Dockerfile.test.go1.15` | golang:1.15 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.48.0/awstesting/sandbox/Dockerfile.test.go1.16` | golang:1.16 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.48.0/awstesting/sandbox/Dockerfile.test.go1.17` | golang:1.17 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.48.0/awstesting/sandbox/Dockerfile.test.go1.18` | golang:1.18 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.48.0/awstesting/sandbox/Dockerfile.test.go1.19` | golang:1.19 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.48.0/awstesting/sandbox/Dockerfile.test.go1.5` | golang:1.5 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.48.0/awstesting/sandbox/Dockerfile.test.go1.6` | golang:1.6 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.48.0/awstesting/sandbox/Dockerfile.test.go1.7` | golang:1.7 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.48.0/awstesting/sandbox/Dockerfile.test.go1.8` | golang:1.8 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.48.0/awstesting/sandbox/Dockerfile.test.go1.9` | golang:1.9 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.48.0/awstesting/sandbox/Dockerfile.test.gotip` | aws-golang:tip | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/kserve/kserve@v0.12.0/Dockerfile` | gcr.io/distroless/static:nonroot | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/kserve/kserve@v0.12.0/tools/tf2openapi/Dockerfile` | gcr.io/distroless/static:nonroot | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/pelletier/go-toml/v2@v2.0.0-beta.8/Dockerfile` | scratch | 1 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/pelletier/go-toml@v1.9.4/Dockerfile` | scratch | 2 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/prometheus/client_golang@v1.17.0/Dockerfile` | quay.io/prometheus/busybox:latest | 2 |  |  |  |  | Unpinned base image: quay.io/prometheus/busybox:latest; No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/x/net@v0.33.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/x/sys@v0.28.0/unix/linux/Dockerfile` | ubuntu:24.10 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/google.golang.org/grpc@v1.59.0/interop/xds/client/Dockerfile` | alpine | 2 |  |  |  |  | Unpinned base image: alpine; No USER directive found (defaults to root) |
| `.gomod-cache/google.golang.org/grpc@v1.59.0/interop/xds/server/Dockerfile` | alpine | 2 |  |  |  |  | Unpinned base image: alpine; No USER directive found (defaults to root) |
| `.gomod-cache/k8s.io/apiextensions-apiserver@v0.28.4/artifacts/simple-image/Dockerfile` | gcr.io/distroless/base-debian10:latest | 1 |  |  |  |  | Unpinned base image: gcr.io/distroless/base-debian10:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/cloud.google.com/go@v0.110.10/.devcontainer/Dockerfile` | mcr.microsoft.com/devcontainers/go:${VARIANT} | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.48.0/awstesting/sandbox/Dockerfile.golang-tip` | buildpack-deps:buster-scm | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.48.0/awstesting/sandbox/Dockerfile.test.go1.10` | golang:1.10 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.48.0/awstesting/sandbox/Dockerfile.test.go1.11` | golang:1.11 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.48.0/awstesting/sandbox/Dockerfile.test.go1.12` | golang:1.12 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.48.0/awstesting/sandbox/Dockerfile.test.go1.13` | golang:1.13 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.48.0/awstesting/sandbox/Dockerfile.test.go1.14` | golang:1.14 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.48.0/awstesting/sandbox/Dockerfile.test.go1.15` | golang:1.15 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.48.0/awstesting/sandbox/Dockerfile.test.go1.16` | golang:1.16 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.48.0/awstesting/sandbox/Dockerfile.test.go1.17` | golang:1.17 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.48.0/awstesting/sandbox/Dockerfile.test.go1.18` | golang:1.18 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.48.0/awstesting/sandbox/Dockerfile.test.go1.19` | golang:1.19 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.48.0/awstesting/sandbox/Dockerfile.test.go1.5` | golang:1.5 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.48.0/awstesting/sandbox/Dockerfile.test.go1.6` | golang:1.6 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.48.0/awstesting/sandbox/Dockerfile.test.go1.7` | golang:1.7 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.48.0/awstesting/sandbox/Dockerfile.test.go1.8` | golang:1.8 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.48.0/awstesting/sandbox/Dockerfile.test.go1.9` | golang:1.9 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.48.0/awstesting/sandbox/Dockerfile.test.gotip` | aws-golang:tip | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/kserve/kserve@v0.12.0/Dockerfile` | gcr.io/distroless/static:nonroot | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/kserve/kserve@v0.12.0/tools/tf2openapi/Dockerfile` | gcr.io/distroless/static:nonroot | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/pelletier/go-toml/v2@v2.0.0-beta.8/Dockerfile` | scratch | 1 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/pelletier/go-toml@v1.9.4/Dockerfile` | scratch | 2 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/prometheus/client_golang@v1.17.0/Dockerfile` | quay.io/prometheus/busybox:latest | 2 |  |  |  |  | Unpinned base image: quay.io/prometheus/busybox:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/net@v0.33.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/sys@v0.28.0/unix/linux/Dockerfile` | ubuntu:24.10 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/google.golang.org/grpc@v1.59.0/interop/xds/client/Dockerfile` | alpine | 2 |  |  |  |  | Unpinned base image: alpine; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/google.golang.org/grpc@v1.59.0/interop/xds/server/Dockerfile` | alpine | 2 |  |  |  |  | Unpinned base image: alpine; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/k8s.io/apiextensions-apiserver@v0.28.4/artifacts/simple-image/Dockerfile` | gcr.io/distroless/base-debian10:latest | 1 |  |  |  |  | Unpinned base image: gcr.io/distroless/base-debian10:latest; No USER directive found (defaults to root) |
| `Dockerfile` | registry.access.redhat.com/ubi9/ubi-minimal:9.5 | 2 | ${USER} |  | multi-arch |  | Unpinned base image: ${DEV_IMAGE} |
| `Dockerfile.develop` | registry.access.redhat.com/ubi9/go-toolset:$GOLANG_VERSION | 1 | root |  | multi-arch |  | Container runs as root user |
| `Dockerfile.develop.ci` | registry.access.redhat.com/ubi9/go-toolset:$GOLANG_VERSION | 1 | root |  | multi-arch |  | Container runs as root user |
| `Dockerfile.konflux` | registry.access.redhat.com/ubi9/ubi-minimal@sha256:24650313873554b6ba16c1a1b6b9f9142604f6ab735113e1695faf2dd07fdede | 2 | ${USER} |  | multi-arch |  |  |

