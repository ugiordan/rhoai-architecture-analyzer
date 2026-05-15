# kueue: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

### Secrets Referenced

| Name | Type | Referenced By |
|------|------|---------------|
| training-operator-v2-webhook-cert | Opaque | deployment/training-operator-v2 |
| training-operator-webhook-cert | Opaque | deployment/training-operator |
| webhook-server-cert | Opaque | deployment/controller-manager, deployment/kuberay-operator |

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| bind | bind | ? | ? | ? | [`.gopath-loader/pkg/mod/github.com/cert-manager/cert-manager@v1.17.1/make/config/bind/deployment.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/.gopath-loader/pkg/mod/github.com/cert-manager/cert-manager@v1.17.1/make/config/bind/deployment.yaml) |
| bind | bind | ? | ? | ? | [`.gomod-cache/github.com/cert-manager/cert-manager@v1.17.1/make/config/bind/deployment.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/.gomod-cache/github.com/cert-manager/cert-manager@v1.17.1/make/config/bind/deployment.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gopath-loader/pkg/mod/sigs.k8s.io/lws@v0.5.1/config/default/manager_webhook_patch.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/.gopath-loader/pkg/mod/sigs.k8s.io/lws@v0.5.1/config/default/manager_webhook_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/components/manager/manager.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/config/components/manager/manager.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/rhoai/manager_webhook_patch.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/config/rhoai/manager_webhook_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/rhoai/manager_metrics_patch.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/config/rhoai/manager_metrics_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/rhoai/manager_config_patch.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/config/rhoai/manager_config_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/dev/manager_config_patch.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/config/dev/manager_config_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/default/manager_webhook_patch.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/config/default/manager_webhook_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gomod-cache/github.com/project-codeflare/appwrapper@v1.1.0/config/default/manager_config_patch.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/.gomod-cache/github.com/project-codeflare/appwrapper@v1.1.0/config/default/manager_config_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gomod-cache/github.com/project-codeflare/appwrapper@v1.1.0/config/default/manager_webhook_patch.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/.gomod-cache/github.com/project-codeflare/appwrapper@v1.1.0/config/default/manager_webhook_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gomod-cache/github.com/project-codeflare/appwrapper@v1.1.0/config/manager/manager.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/.gomod-cache/github.com/project-codeflare/appwrapper@v1.1.0/config/manager/manager.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/default/manager_visibility_patch.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/config/default/manager_visibility_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/default/manager_metrics_patch.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/config/default/manager_metrics_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gomod-cache/sigs.k8s.io/jobset@v0.8.0/config/components/manager/manager.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/.gomod-cache/sigs.k8s.io/jobset@v0.8.0/config/components/manager/manager.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gomod-cache/sigs.k8s.io/jobset@v0.8.0/config/default/manager_config_patch.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/.gomod-cache/sigs.k8s.io/jobset@v0.8.0/config/default/manager_config_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gomod-cache/sigs.k8s.io/jobset@v0.8.0/config/default/manager_webhook_patch.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/.gomod-cache/sigs.k8s.io/jobset@v0.8.0/config/default/manager_webhook_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gomod-cache/sigs.k8s.io/lws@v0.5.1/config/default/manager_config_patch.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/.gomod-cache/sigs.k8s.io/lws@v0.5.1/config/default/manager_config_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gomod-cache/sigs.k8s.io/lws@v0.5.1/config/default/manager_webhook_patch.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/.gomod-cache/sigs.k8s.io/lws@v0.5.1/config/default/manager_webhook_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gomod-cache/sigs.k8s.io/lws@v0.5.1/config/manager/manager.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/.gomod-cache/sigs.k8s.io/lws@v0.5.1/config/manager/manager.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/default/manager_config_patch.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/config/default/manager_config_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gopath-loader/pkg/mod/github.com/project-codeflare/appwrapper@v1.1.0/config/default/manager_config_patch.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/.gopath-loader/pkg/mod/github.com/project-codeflare/appwrapper@v1.1.0/config/default/manager_config_patch.yaml) |
| controller-manager | kube-rbac-proxy | ? | ? | ? | [`config/default/manager_auth_proxy_patch.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/config/default/manager_auth_proxy_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gopath-loader/pkg/mod/github.com/project-codeflare/appwrapper@v1.1.0/config/manager/manager.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/.gopath-loader/pkg/mod/github.com/project-codeflare/appwrapper@v1.1.0/config/manager/manager.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gopath-loader/pkg/mod/github.com/project-codeflare/appwrapper@v1.1.0/config/default/manager_webhook_patch.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/.gopath-loader/pkg/mod/github.com/project-codeflare/appwrapper@v1.1.0/config/default/manager_webhook_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/alpha-enabled/manager_config_patch.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/config/alpha-enabled/manager_config_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.8.0/config/components/manager/manager.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.8.0/config/components/manager/manager.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.8.0/config/default/manager_config_patch.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.8.0/config/default/manager_config_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.8.0/config/default/manager_webhook_patch.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.8.0/config/default/manager_webhook_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gopath-loader/pkg/mod/sigs.k8s.io/lws@v0.5.1/config/default/manager_config_patch.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/.gopath-loader/pkg/mod/sigs.k8s.io/lws@v0.5.1/config/default/manager_config_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`.gopath-loader/pkg/mod/sigs.k8s.io/lws@v0.5.1/config/manager/manager.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/.gopath-loader/pkg/mod/sigs.k8s.io/lws@v0.5.1/config/manager/manager.yaml) |
| kuberay-operator | kuberay-operator | ? | ? | ? | [`.gopath-loader/pkg/mod/github.com/ray-project/kuberay/ray-operator@v1.3.1/config/manager/manager.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/.gopath-loader/pkg/mod/github.com/ray-project/kuberay/ray-operator@v1.3.1/config/manager/manager.yaml) |
| kuberay-operator | kuberay-operator | ? | ? | ? | [`.gopath-loader/pkg/mod/github.com/ray-project/kuberay/ray-operator@v1.3.1/config/default-with-webhooks/manager_webhook_patch.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/.gopath-loader/pkg/mod/github.com/ray-project/kuberay/ray-operator@v1.3.1/config/default-with-webhooks/manager_webhook_patch.yaml) |
| kuberay-operator | kuberay-operator | ? | ? | ? | [`.gomod-cache/github.com/ray-project/kuberay/ray-operator@v1.3.1/config/manager/manager.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/.gomod-cache/github.com/ray-project/kuberay/ray-operator@v1.3.1/config/manager/manager.yaml) |
| kuberay-operator | kuberay-operator | ? | ? | ? | [`.gomod-cache/github.com/ray-project/kuberay/ray-operator@v1.3.1/config/default-with-webhooks/manager_webhook_patch.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/.gomod-cache/github.com/ray-project/kuberay/ray-operator@v1.3.1/config/default-with-webhooks/manager_webhook_patch.yaml) |
| mpi-operator | mpi-operator | ? | ? | ? | [`.gomod-cache/github.com/kubeflow/mpi-operator@v0.6.0/manifests/base/deployment.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/.gomod-cache/github.com/kubeflow/mpi-operator@v0.6.0/manifests/base/deployment.yaml) |
| mpi-operator | mpi-operator | ? | ? | ? | [`.gopath-loader/pkg/mod/github.com/kubeflow/mpi-operator@v0.6.0/manifests/base/deployment.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/.gopath-loader/pkg/mod/github.com/kubeflow/mpi-operator@v0.6.0/manifests/base/deployment.yaml) |
| the-deployment | the-container | ? | ? | ? | [`.gomod-cache/k8s.io/cli-runtime@v0.32.3/artifacts/kustomization/deployment.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/.gomod-cache/k8s.io/cli-runtime@v0.32.3/artifacts/kustomization/deployment.yaml) |
| the-deployment | the-container | ? | ? | ? | [`.gopath-loader/pkg/mod/k8s.io/cli-runtime@v0.32.3/artifacts/kustomization/deployment.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/.gopath-loader/pkg/mod/k8s.io/cli-runtime@v0.32.3/artifacts/kustomization/deployment.yaml) |
| training-operator | training-operator | ? | ? | ? | [`.gopath-loader/pkg/mod/github.com/kubeflow/training-operator@v1.9.0/manifests/base/deployment.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/.gopath-loader/pkg/mod/github.com/kubeflow/training-operator@v1.9.0/manifests/base/deployment.yaml) |
| training-operator | training-operator | ? | ? | ? | [`.gomod-cache/github.com/kubeflow/training-operator@v1.9.0/manifests/base/deployment.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/.gomod-cache/github.com/kubeflow/training-operator@v1.9.0/manifests/base/deployment.yaml) |
| training-operator-v2 | manager | ? | ? | ? | [`.gopath-loader/pkg/mod/github.com/kubeflow/training-operator@v1.9.0/manifests/v2/base/manager/manager.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/.gopath-loader/pkg/mod/github.com/kubeflow/training-operator@v1.9.0/manifests/v2/base/manager/manager.yaml) |
| training-operator-v2 | manager | ? | ? | ? | [`.gomod-cache/github.com/kubeflow/training-operator@v1.9.0/manifests/v2/base/manager/manager.yaml`](https://github.com/opendatahub-io/kueue/blob/97024bd289d2cc5c9369b40d9f3483ab1483143d/.gomod-cache/github.com/kubeflow/training-operator@v1.9.0/manifests/v2/base/manager/manager.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `.gomod-cache/github.com/cert-manager/cert-manager@v1.17.1/make/config/pebble/Containerfile.pebble` | $BASE_IMAGE | 1 | 1000 |  |  |  | Unpinned base image: $BASE_IMAGE |
| `.gomod-cache/github.com/cert-manager/cert-manager@v1.17.1/make/config/samplewebhook/Containerfile.samplewebhook` | $BASE_IMAGE | 1 | 1000 |  |  |  | Unpinned base image: $BASE_IMAGE |
| `.gomod-cache/github.com/grpc-ecosystem/grpc-gateway/v2@v2.25.1/.github/Dockerfile` | golang:1.23.4 | 1 | vscode |  |  |  |  |
| `.gomod-cache/github.com/kubeflow/mpi-operator@v0.6.0/Dockerfile` | gcr.io/distroless/base-debian12:latest | 2 |  |  |  |  | Unpinned base image: gcr.io/distroless/base-debian12:latest; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/kubeflow/mpi-operator@v0.6.0/build/base/Dockerfile` | debian:bookworm | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/kubeflow/training-operator@v1.9.0/build/images/kubectl-delivery/Dockerfile` | alpine:3.17 | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/kubeflow/training-operator@v1.9.0/build/images/training-operator/Dockerfile` | gcr.io/distroless/static:latest | 2 |  |  |  |  | Unpinned base image: gcr.io/distroless/static:latest; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/kubeflow/training-operator@v1.9.0/cmd/initializer_v2/dataset/Dockerfile` | python:3.11-alpine | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/kubeflow/training-operator@v1.9.0/cmd/initializer_v2/model/Dockerfile` | python:3.11-alpine | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/kubeflow/training-operator@v1.9.0/cmd/training-operator.v2alpha1/Dockerfile` | gcr.io/distroless/static:nonroot | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/kubeflow/training-operator@v1.9.0/sdk/python/Dockerfile.conformance` | python:3.10-alpine | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/kubeflow/training-operator@v1.9.0/sdk/python/kubeflow/storage_initializer/Dockerfile` | python:3.11 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/kubeflow/training-operator@v1.9.0/sdk/python/kubeflow/trainer/Dockerfile` | nvcr.io/nvidia/pytorch:24.06-py3 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/kubeflow/training-operator@v1.9.0/sdk/python/kubeflow/trainer/Dockerfile.cpu` | python:3.11 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/project-codeflare/appwrapper@v1.1.0/Dockerfile` | gcr.io/distroless/static:nonroot | 2 | 65532:65532 |  | multi-arch |  |  |
| `.gomod-cache/github.com/prometheus/client_golang@v1.21.1/Dockerfile` | quay.io/prometheus/busybox:latest | 2 |  |  |  |  | Unpinned base image: quay.io/prometheus/busybox:latest; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/ray-project/kuberay/ray-operator@v1.3.1/Dockerfile` | gcr.io/distroless/base-debian12:nonroot | 2 | 65532:65532 |  |  |  |  |
| `.gomod-cache/github.com/ray-project/kuberay/ray-operator@v1.3.1/Dockerfile.buildx` | gcr.io/distroless/base-debian12:nonroot | 1 | 65532:65532 |  | multi-arch |  |  |
| `.gomod-cache/golang.org/x/net@v0.38.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/x/sys@v0.31.0/unix/linux/Dockerfile` | ubuntu:24.10 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/k8s.io/apiextensions-apiserver@v0.32.3/artifacts/simple-image/Dockerfile` | gcr.io/distroless/base-debian10:latest | 1 |  |  |  |  | Unpinned base image: gcr.io/distroless/base-debian10:latest; No USER directive found (defaults to root) |
| `.gomod-cache/sigs.k8s.io/gateway-api@v1.1.0/docker/Dockerfile.echo-advanced` | gcr.io/istio-release/app:1.21.2 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/sigs.k8s.io/gateway-api@v1.1.0/docker/Dockerfile.echo-basic` | gcr.io/distroless/static:nonroot | 2 | nonroot:nonroot |  |  |  |  |
| `.gomod-cache/sigs.k8s.io/jobset@v0.8.0/Dockerfile` | ${BASE_IMAGE} | 2 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${BUILDER_IMAGE}; Unpinned base image: ${BASE_IMAGE} |
| `.gomod-cache/sigs.k8s.io/jobset@v0.8.0/sdk/python/Dockerfile` | python:3.13 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/sigs.k8s.io/kustomize/kyaml@v0.18.1/fn/framework/example/Dockerfile` | alpine:latest | 2 |  |  |  |  | Unpinned base image: alpine:latest; No USER directive found (defaults to root) |
| `.gomod-cache/sigs.k8s.io/lws@v0.5.1/Dockerfile` | ${BASE_IMAGE} | 2 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${BUILDER_IMAGE}; Unpinned base image: ${BASE_IMAGE} |
| `.gopath-loader/pkg/mod/github.com/cert-manager/cert-manager@v1.17.1/make/config/pebble/Containerfile.pebble` | $BASE_IMAGE | 1 | 1000 |  |  |  | Unpinned base image: $BASE_IMAGE |
| `.gopath-loader/pkg/mod/github.com/cert-manager/cert-manager@v1.17.1/make/config/samplewebhook/Containerfile.samplewebhook` | $BASE_IMAGE | 1 | 1000 |  |  |  | Unpinned base image: $BASE_IMAGE |
| `.gopath-loader/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.25.1/.github/Dockerfile` | golang:1.23.4 | 1 | vscode |  |  |  |  |
| `.gopath-loader/pkg/mod/github.com/kubeflow/mpi-operator@v0.6.0/Dockerfile` | gcr.io/distroless/base-debian12:latest | 2 |  |  |  |  | Unpinned base image: gcr.io/distroless/base-debian12:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/kubeflow/mpi-operator@v0.6.0/build/base/Dockerfile` | debian:bookworm | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/kubeflow/training-operator@v1.9.0/build/images/kubectl-delivery/Dockerfile` | alpine:3.17 | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/kubeflow/training-operator@v1.9.0/build/images/training-operator/Dockerfile` | gcr.io/distroless/static:latest | 2 |  |  |  |  | Unpinned base image: gcr.io/distroless/static:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/kubeflow/training-operator@v1.9.0/cmd/initializer_v2/dataset/Dockerfile` | python:3.11-alpine | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/kubeflow/training-operator@v1.9.0/cmd/initializer_v2/model/Dockerfile` | python:3.11-alpine | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/kubeflow/training-operator@v1.9.0/cmd/training-operator.v2alpha1/Dockerfile` | gcr.io/distroless/static:nonroot | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/kubeflow/training-operator@v1.9.0/sdk/python/Dockerfile.conformance` | python:3.10-alpine | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/kubeflow/training-operator@v1.9.0/sdk/python/kubeflow/storage_initializer/Dockerfile` | python:3.11 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/kubeflow/training-operator@v1.9.0/sdk/python/kubeflow/trainer/Dockerfile` | nvcr.io/nvidia/pytorch:24.06-py3 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/kubeflow/training-operator@v1.9.0/sdk/python/kubeflow/trainer/Dockerfile.cpu` | python:3.11 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/project-codeflare/appwrapper@v1.1.0/Dockerfile` | gcr.io/distroless/static:nonroot | 2 | 65532:65532 |  | multi-arch |  |  |
| `.gopath-loader/pkg/mod/github.com/prometheus/client_golang@v1.21.1/Dockerfile` | quay.io/prometheus/busybox:latest | 2 |  |  |  |  | Unpinned base image: quay.io/prometheus/busybox:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/ray-project/kuberay/ray-operator@v1.3.1/Dockerfile` | gcr.io/distroless/base-debian12:nonroot | 2 | 65532:65532 |  |  |  |  |
| `.gopath-loader/pkg/mod/github.com/ray-project/kuberay/ray-operator@v1.3.1/Dockerfile.buildx` | gcr.io/distroless/base-debian12:nonroot | 1 | 65532:65532 |  | multi-arch |  |  |
| `.gopath-loader/pkg/mod/golang.org/x/net@v0.38.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/sys@v0.31.0/unix/linux/Dockerfile` | ubuntu:24.10 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/k8s.io/apiextensions-apiserver@v0.32.3/artifacts/simple-image/Dockerfile` | gcr.io/distroless/base-debian10:latest | 1 |  |  |  |  | Unpinned base image: gcr.io/distroless/base-debian10:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api@v1.1.0/docker/Dockerfile.echo-advanced` | gcr.io/istio-release/app:1.21.2 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api@v1.1.0/docker/Dockerfile.echo-basic` | gcr.io/distroless/static:nonroot | 2 | nonroot:nonroot |  |  |  |  |
| `.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.8.0/Dockerfile` | ${BASE_IMAGE} | 2 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${BUILDER_IMAGE}; Unpinned base image: ${BASE_IMAGE} |
| `.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.8.0/sdk/python/Dockerfile` | python:3.13 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/sigs.k8s.io/kustomize/kyaml@v0.18.1/fn/framework/example/Dockerfile` | alpine:latest | 2 |  |  |  |  | Unpinned base image: alpine:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/sigs.k8s.io/lws@v0.5.1/Dockerfile` | ${BASE_IMAGE} | 2 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${BUILDER_IMAGE}; Unpinned base image: ${BASE_IMAGE} |
| `Dockerfile` | ${BASE_IMAGE} | 2 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${BUILDER_IMAGE}; Unpinned base image: ${BASE_IMAGE} |
| `Dockerfile.rhoai` | registry.access.redhat.com/ubi9/ubi:latest | 3 | 65532:65532 |  |  |  | Unpinned base image: ${GOLANG_IMAGE}; Unpinned base image: registry.access.redhat.com/ubi9/ubi:latest; Unpinned base image: registry.access.redhat.com/ubi9/ubi:latest |
| `cmd/experimental/kueue-viz/backend/Dockerfile` | ${BASE_IMAGE} | 2 | 65532:65532 |  |  |  | Unpinned base image: ${BUILDER_IMAGE}; Unpinned base image: ${BASE_IMAGE} |
| `cmd/experimental/kueue-viz/frontend/Dockerfile` | node:23 | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `cmd/importer/Dockerfile` | ${BASE_IMAGE} | 2 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${BUILDER_IMAGE}; Unpinned base image: ${BASE_IMAGE} |

