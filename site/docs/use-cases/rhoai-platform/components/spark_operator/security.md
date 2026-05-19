# spark-operator: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

### Secrets Referenced

| Name | Type | Referenced By |
|------|------|---------------|
| webhook-server-cert | Opaque | deployment/controller-manager |

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| controller-manager | manager | ? | ? | ? | [`.gopath-loader/pkg/mod/sigs.k8s.io/scheduler-plugins@v0.32.7/config/default/manager_config_patch.yaml`](https://github.com/kubeflow/spark-operator/blob/1ff69d896dc2b7c0769f5bde06d3ab6f25089228/.gopath-loader/pkg/mod/sigs.k8s.io/scheduler-plugins@v0.32.7/config/default/manager_config_patch.yaml) |
| controller-manager | kube-rbac-proxy | ? | ? | ? | [`.gomod-cache/sigs.k8s.io/scheduler-plugins@v0.32.7/config/default/manager_auth_proxy_patch.yaml`](https://github.com/kubeflow/spark-operator/blob/1ff69d896dc2b7c0769f5bde06d3ab6f25089228/.gomod-cache/sigs.k8s.io/scheduler-plugins@v0.32.7/config/default/manager_auth_proxy_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gomod-cache/sigs.k8s.io/scheduler-plugins@v0.32.7/config/default/manager_auth_proxy_patch.yaml`](https://github.com/kubeflow/spark-operator/blob/1ff69d896dc2b7c0769f5bde06d3ab6f25089228/.gomod-cache/sigs.k8s.io/scheduler-plugins@v0.32.7/config/default/manager_auth_proxy_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gomod-cache/sigs.k8s.io/scheduler-plugins@v0.32.7/config/default/manager_config_patch.yaml`](https://github.com/kubeflow/spark-operator/blob/1ff69d896dc2b7c0769f5bde06d3ab6f25089228/.gomod-cache/sigs.k8s.io/scheduler-plugins@v0.32.7/config/default/manager_config_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gomod-cache/sigs.k8s.io/scheduler-plugins@v0.32.7/config/manager/manager.yaml`](https://github.com/kubeflow/spark-operator/blob/1ff69d896dc2b7c0769f5bde06d3ab6f25089228/.gomod-cache/sigs.k8s.io/scheduler-plugins@v0.32.7/config/manager/manager.yaml) |
| controller-manager | kube-rbac-proxy | ? | ? | ? | [`.gopath-loader/pkg/mod/sigs.k8s.io/scheduler-plugins@v0.32.7/config/default/manager_auth_proxy_patch.yaml`](https://github.com/kubeflow/spark-operator/blob/1ff69d896dc2b7c0769f5bde06d3ab6f25089228/.gopath-loader/pkg/mod/sigs.k8s.io/scheduler-plugins@v0.32.7/config/default/manager_auth_proxy_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gopath-loader/pkg/mod/sigs.k8s.io/scheduler-plugins@v0.32.7/config/default/manager_auth_proxy_patch.yaml`](https://github.com/kubeflow/spark-operator/blob/1ff69d896dc2b7c0769f5bde06d3ab6f25089228/.gopath-loader/pkg/mod/sigs.k8s.io/scheduler-plugins@v0.32.7/config/default/manager_auth_proxy_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gopath-loader/pkg/mod/sigs.k8s.io/scheduler-plugins@v0.32.7/config/manager/manager.yaml`](https://github.com/kubeflow/spark-operator/blob/1ff69d896dc2b7c0769f5bde06d3ab6f25089228/.gopath-loader/pkg/mod/sigs.k8s.io/scheduler-plugins@v0.32.7/config/manager/manager.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/default/manager_webhook_patch.yaml`](https://github.com/kubeflow/spark-operator/blob/1ff69d896dc2b7c0769f5bde06d3ab6f25089228/config/default/manager_webhook_patch.yaml) |
| peaks | peaks | ? | ? | ? | [`.gomod-cache/sigs.k8s.io/scheduler-plugins@v0.32.7/pkg/trimaran/peaks/deployment/deployment.yaml`](https://github.com/kubeflow/spark-operator/blob/1ff69d896dc2b7c0769f5bde06d3ab6f25089228/.gomod-cache/sigs.k8s.io/scheduler-plugins@v0.32.7/pkg/trimaran/peaks/deployment/deployment.yaml) |
| peaks | peaks | ? | ? | ? | [`.gopath-loader/pkg/mod/sigs.k8s.io/scheduler-plugins@v0.32.7/pkg/trimaran/peaks/deployment/deployment.yaml`](https://github.com/kubeflow/spark-operator/blob/1ff69d896dc2b7c0769f5bde06d3ab6f25089228/.gopath-loader/pkg/mod/sigs.k8s.io/scheduler-plugins@v0.32.7/pkg/trimaran/peaks/deployment/deployment.yaml) |
| spark-operator-controller | controller | true | true | false | [`kustomize:config/overlays/odh`](https://github.com/kubeflow/spark-operator/blob/1ff69d896dc2b7c0769f5bde06d3ab6f25089228/kustomize:config/overlays/odh) |
| spark-operator-webhook | webhook | true | true | false | [`kustomize:config/overlays/odh`](https://github.com/kubeflow/spark-operator/blob/1ff69d896dc2b7c0769f5bde06d3ab6f25089228/kustomize:config/overlays/odh) |
| the-deployment | the-container | ? | ? | ? | [`.gopath-loader/pkg/mod/k8s.io/cli-runtime@v0.32.5/artifacts/kustomization/deployment.yaml`](https://github.com/kubeflow/spark-operator/blob/1ff69d896dc2b7c0769f5bde06d3ab6f25089228/.gopath-loader/pkg/mod/k8s.io/cli-runtime@v0.32.5/artifacts/kustomization/deployment.yaml) |
| the-deployment | the-container | ? | ? | ? | [`.gomod-cache/k8s.io/cli-runtime@v0.32.5/artifacts/kustomization/deployment.yaml`](https://github.com/kubeflow/spark-operator/blob/1ff69d896dc2b7c0769f5bde06d3ab6f25089228/.gomod-cache/k8s.io/cli-runtime@v0.32.5/artifacts/kustomization/deployment.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `.gomod-cache/github.com/containerd/containerd@v1.7.29/.github/workflows/release/Dockerfile` | scratch | 7 |  |  | multi-arch |  | Unpinned base image: $GO_IMAGE; Unpinned base image: ${BASE_IMAGE}; Unpinned base image: base; Unpinned base image: base; Unpinned base image: ${TARGETOS}; Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/containerd/containerd@v1.7.29/contrib/Dockerfile.test` | build-env | 11 |  |  |  |  | Unpinned base image: golang; Unpinned base image: golang; Unpinned base image: golang; Unpinned base image: golang; Unpinned base image: build-env; Unpinned base image: integration; Unpinned base image: cri-integration; Unpinned base image: critest; Unpinned base image: golang; Unpinned base image: build-env; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/containerd/containerd@v1.7.29/integration/images/volume-copy-up/Dockerfile` | $BASE | 1 |  |  |  |  | Unpinned base image: $BASE; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/containerd/containerd@v1.7.29/integration/images/volume-ownership/Dockerfile` | ubuntu | 1 |  |  |  |  | Unpinned base image: ubuntu; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/pelletier/go-toml/v2@v2.2.4/Dockerfile` | scratch | 1 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/prometheus/client_golang@v1.23.2/Dockerfile` | quay.io/prometheus/busybox:latest | 2 |  |  |  |  | Unpinned base image: quay.io/prometheus/busybox:latest; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/rubenv/sql-migrate@v1.8.0/Dockerfile` | alpine:${ALPINE_VERSION} | 3 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/x/net@v0.44.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/x/sys@v0.36.0/unix/linux/Dockerfile` | ubuntu:24.10 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/k8s.io/apiextensions-apiserver@v0.32.5/artifacts/simple-image/Dockerfile` | gcr.io/distroless/base-debian10:latest | 1 |  |  |  |  | Unpinned base image: gcr.io/distroless/base-debian10:latest; No USER directive found (defaults to root) |
| `.gomod-cache/sigs.k8s.io/kustomize/kyaml@v0.19.0/fn/framework/example/Dockerfile` | alpine:latest | 2 |  |  |  |  | Unpinned base image: alpine:latest; No USER directive found (defaults to root) |
| `.gomod-cache/sigs.k8s.io/scheduler-plugins@v0.32.7/build/controller/Dockerfile` | $DISTROLESS_BASE_IMAGE | 2 | 65532:65532 |  | multi-arch |  | Unpinned base image: $GO_BASE_IMAGE; Unpinned base image: $DISTROLESS_BASE_IMAGE |
| `.gomod-cache/sigs.k8s.io/scheduler-plugins@v0.32.7/build/scheduler/Dockerfile` | $DISTROLESS_BASE_IMAGE | 2 | 65532:65532 |  | multi-arch |  | Unpinned base image: $GO_BASE_IMAGE; Unpinned base image: $DISTROLESS_BASE_IMAGE |
| `.gopath-loader/pkg/mod/github.com/containerd/containerd@v1.7.29/.github/workflows/release/Dockerfile` | scratch | 7 |  |  | multi-arch |  | Unpinned base image: $GO_IMAGE; Unpinned base image: ${BASE_IMAGE}; Unpinned base image: base; Unpinned base image: base; Unpinned base image: ${TARGETOS}; Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/containerd/containerd@v1.7.29/contrib/Dockerfile.test` | build-env | 11 |  |  |  |  | Unpinned base image: golang; Unpinned base image: golang; Unpinned base image: golang; Unpinned base image: golang; Unpinned base image: build-env; Unpinned base image: integration; Unpinned base image: cri-integration; Unpinned base image: critest; Unpinned base image: golang; Unpinned base image: build-env; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/containerd/containerd@v1.7.29/integration/images/volume-copy-up/Dockerfile` | $BASE | 1 |  |  |  |  | Unpinned base image: $BASE; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/containerd/containerd@v1.7.29/integration/images/volume-ownership/Dockerfile` | ubuntu | 1 |  |  |  |  | Unpinned base image: ubuntu; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/pelletier/go-toml/v2@v2.2.4/Dockerfile` | scratch | 1 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/prometheus/client_golang@v1.23.2/Dockerfile` | quay.io/prometheus/busybox:latest | 2 |  |  |  |  | Unpinned base image: quay.io/prometheus/busybox:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/rubenv/sql-migrate@v1.8.0/Dockerfile` | alpine:${ALPINE_VERSION} | 3 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/net@v0.44.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/sys@v0.36.0/unix/linux/Dockerfile` | ubuntu:24.10 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/k8s.io/apiextensions-apiserver@v0.32.5/artifacts/simple-image/Dockerfile` | gcr.io/distroless/base-debian10:latest | 1 |  |  |  |  | Unpinned base image: gcr.io/distroless/base-debian10:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/sigs.k8s.io/kustomize/kyaml@v0.19.0/fn/framework/example/Dockerfile` | alpine:latest | 2 |  |  |  |  | Unpinned base image: alpine:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/sigs.k8s.io/scheduler-plugins@v0.32.7/build/controller/Dockerfile` | $DISTROLESS_BASE_IMAGE | 2 | 65532:65532 |  | multi-arch |  | Unpinned base image: $GO_BASE_IMAGE; Unpinned base image: $DISTROLESS_BASE_IMAGE |
| `.gopath-loader/pkg/mod/sigs.k8s.io/scheduler-plugins@v0.32.7/build/scheduler/Dockerfile` | $DISTROLESS_BASE_IMAGE | 2 | 65532:65532 |  | multi-arch |  | Unpinned base image: $GO_BASE_IMAGE; Unpinned base image: $DISTROLESS_BASE_IMAGE |
| `Dockerfile` | ${SPARK_IMAGE} | 2 | ${SPARK_UID}:${SPARK_GID} |  | multi-arch |  | Unpinned base image: ${SPARK_IMAGE} |
| `Dockerfile.konflux` | ${BASE_IMAGE} | 2 | ${SPARK_UID}:${SPARK_GID} |  | multi-arch |  | Unpinned base image: ${GO_BUILDER_IMAGE}; Unpinned base image: ${BASE_IMAGE} |
| `docker/Dockerfile.kubectl` | ${BASE_IMAGE} | 2 |  |  | multi-arch |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: ${BASE_IMAGE}; No USER directive found (defaults to root) |
| `spark-docker/Dockerfile` | ${SPARK_IMAGE} | 1 | ${spark_uid} |  |  |  | Unpinned base image: ${SPARK_IMAGE} |

