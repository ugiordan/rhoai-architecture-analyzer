# data-science-pipelines: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

### Secrets Referenced

| Name | Type | Referenced By |
|------|------|---------------|
| kfp-api-webhook-cert | Opaque | deployment/ml-pipeline |
| mlpipeline-minio-artifact | Opaque | deployment/kubeflow-pipelines-profile-controller |

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| cache-server | server | ? | ? | ? | [`manifests/kustomize/base/installs/multi-user/cache/deployment-patch.yaml`](https://github.com/kubeflow/data-science-pipelines/blob/e61fa54e17eb9a52898792f7554ea3e00dc8eb0b/manifests/kustomize/base/installs/multi-user/cache/deployment-patch.yaml) |
| kubeflow-pipelines-profile-controller | profile-controller | true | ? | ? | [`manifests/kustomize/base/installs/multi-user/pipelines-profile-controller/deployment.yaml`](https://github.com/kubeflow/data-science-pipelines/blob/e61fa54e17eb9a52898792f7554ea3e00dc8eb0b/manifests/kustomize/base/installs/multi-user/pipelines-profile-controller/deployment.yaml) |
| metadata-writer | main | ? | ? | ? | [`manifests/kustomize/base/installs/multi-user/metadata-writer/deployment-patch.yaml`](https://github.com/kubeflow/data-science-pipelines/blob/e61fa54e17eb9a52898792f7554ea3e00dc8eb0b/manifests/kustomize/base/installs/multi-user/metadata-writer/deployment-patch.yaml) |
| ml-pipeline | ml-pipeline-api-server | ? | ? | ? | [`manifests/kustomize/env/cert-manager/platform-agnostic-k8s-native/patches/deployment.yaml`](https://github.com/kubeflow/data-science-pipelines/blob/e61fa54e17eb9a52898792f7554ea3e00dc8eb0b/manifests/kustomize/env/cert-manager/platform-agnostic-k8s-native/patches/deployment.yaml) |
| ml-pipeline | ml-pipeline-api-server | ? | ? | ? | [`manifests/kustomize/env/cert-manager/platform-agnostic-multi-user-k8s-native/patches/deployment.yaml`](https://github.com/kubeflow/data-science-pipelines/blob/e61fa54e17eb9a52898792f7554ea3e00dc8eb0b/manifests/kustomize/env/cert-manager/platform-agnostic-multi-user-k8s-native/patches/deployment.yaml) |
| ml-pipeline | ml-pipeline-api-server | ? | ? | ? | [`manifests/kustomize/base/installs/multi-user/api-service/deployment-patch.yaml`](https://github.com/kubeflow/data-science-pipelines/blob/e61fa54e17eb9a52898792f7554ea3e00dc8eb0b/manifests/kustomize/base/installs/multi-user/api-service/deployment-patch.yaml) |
| ml-pipeline-persistenceagent | ml-pipeline-persistenceagent | ? | ? | ? | [`manifests/kustomize/base/installs/multi-user/persistence-agent/deployment-patch.yaml`](https://github.com/kubeflow/data-science-pipelines/blob/e61fa54e17eb9a52898792f7554ea3e00dc8eb0b/manifests/kustomize/base/installs/multi-user/persistence-agent/deployment-patch.yaml) |
| ml-pipeline-scheduledworkflow | ml-pipeline-scheduledworkflow | ? | ? | ? | [`manifests/kustomize/base/installs/multi-user/scheduled-workflow/deployment-patch.yaml`](https://github.com/kubeflow/data-science-pipelines/blob/e61fa54e17eb9a52898792f7554ea3e00dc8eb0b/manifests/kustomize/base/installs/multi-user/scheduled-workflow/deployment-patch.yaml) |
| ml-pipeline-ui | ml-pipeline-ui | ? | ? | ? | [`manifests/kustomize/base/installs/multi-user/pipelines-ui/deployment-patch.yaml`](https://github.com/kubeflow/data-science-pipelines/blob/e61fa54e17eb9a52898792f7554ea3e00dc8eb0b/manifests/kustomize/base/installs/multi-user/pipelines-ui/deployment-patch.yaml) |
| ml-pipeline-viewer-crd | ml-pipeline-viewer-crd | ? | ? | ? | [`manifests/kustomize/base/installs/multi-user/viewer-controller/deployment-patch.yaml`](https://github.com/kubeflow/data-science-pipelines/blob/e61fa54e17eb9a52898792f7554ea3e00dc8eb0b/manifests/kustomize/base/installs/multi-user/viewer-controller/deployment-patch.yaml) |
| squid | squid | ? | ? | ? | [`.github/resources/squid/manifests/deployment.yaml`](https://github.com/kubeflow/data-science-pipelines/blob/e61fa54e17eb9a52898792f7554ea3e00dc8eb0b/.github/resources/squid/manifests/deployment.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `.github/resources/squid/Containerfile` | quay.io/fedora/fedora:41 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/cloud.google.com/go@v0.121.2/.devcontainer/Dockerfile` | mcr.microsoft.com/devcontainers/go:${VARIANT} | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/argoproj/argo-workflows/v3@v3.7.11/.devcontainer/builder/Dockerfile` | mcr.microsoft.com/vscode/devcontainers/base:ubuntu-22.04 | 1 |  |  | multi-arch |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/argoproj/argo-workflows/v3@v3.7.11/Dockerfile` | gcr.io/distroless/static | 10 | 8737 |  |  |  | Unpinned base image: builder; Unpinned base image: builder; Unpinned base image: builder; Unpinned base image: gcr.io/distroless/static; Unpinned base image: argoexec-base; Unpinned base image: argoexec-base; Unpinned base image: gcr.io/distroless/static; Unpinned base image: gcr.io/distroless/static |
| `.gomod-cache/github.com/argoproj/argo-workflows/v3@v3.7.11/Dockerfile.windows` | argoexec-base | 4 | Administrator |  |  |  | Unpinned base image: builder; Unpinned base image: argoexec-base |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.55.5/awstesting/sandbox/Dockerfile.golang-tip` | buildpack-deps:buster-scm | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.55.5/awstesting/sandbox/Dockerfile.test.go1.10` | golang:1.10 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.55.5/awstesting/sandbox/Dockerfile.test.go1.11` | golang:1.11 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.55.5/awstesting/sandbox/Dockerfile.test.go1.12` | golang:1.12 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.55.5/awstesting/sandbox/Dockerfile.test.go1.13` | golang:1.13 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.55.5/awstesting/sandbox/Dockerfile.test.go1.14` | golang:1.14 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.55.5/awstesting/sandbox/Dockerfile.test.go1.15` | golang:1.15 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.55.5/awstesting/sandbox/Dockerfile.test.go1.16` | golang:1.16 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.55.5/awstesting/sandbox/Dockerfile.test.go1.17` | golang:1.17 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.55.5/awstesting/sandbox/Dockerfile.test.go1.18` | golang:1.18 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.55.5/awstesting/sandbox/Dockerfile.test.go1.19` | golang:1.19 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.55.5/awstesting/sandbox/Dockerfile.test.go1.5` | golang:1.5 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.55.5/awstesting/sandbox/Dockerfile.test.go1.6` | golang:1.6 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.55.5/awstesting/sandbox/Dockerfile.test.go1.7` | golang:1.7 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.55.5/awstesting/sandbox/Dockerfile.test.go1.8` | golang:1.8 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.55.5/awstesting/sandbox/Dockerfile.test.go1.9` | golang:1.9 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.55.5/awstesting/sandbox/Dockerfile.test.gotip` | aws-golang:tip | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/envoyproxy/protoc-gen-validate@v1.3.0/Dockerfile` | ubuntu:focal | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/google/addlicense@v0.0.0-20200906110928-a0294312aa76/Dockerfile` | golang:1-alpine | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/grpc-ecosystem/grpc-gateway/v2@v2.27.1/.github/Dockerfile` | golang:1.24.4 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/.circleci/Dockerfile` | golang:1.15.3 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/mattn/go-sqlite3@v1.14.34/_example/simple/Dockerfile` | scratch | 2 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/pelletier/go-toml/v2@v2.2.4/Dockerfile` | scratch | 1 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/prometheus/client_golang@v1.23.2/Dockerfile` | quay.io/prometheus/busybox:latest | 2 |  |  |  |  | Unpinned base image: quay.io/prometheus/busybox:latest; No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/toolchain@v0.0.1-go1.25.7.linux-amd64/src/crypto/internal/boring/Dockerfile` | $ubuntu:focal | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/toolchain@v0.0.1-go1.25.7.linux-amd64/src/crypto/internal/fips140/nistec/fiat/Dockerfile` | coqorg/coq:8.13.2 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/x/net@v0.51.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/x/sys@v0.41.0/unix/linux/Dockerfile` | ubuntu:25.04 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/k8s.io/apiextensions-apiserver@v0.35.2/artifacts/simple-image/Dockerfile` | gcr.io/distroless/base-debian10:latest | 1 |  |  |  |  | Unpinned base image: gcr.io/distroless/base-debian10:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/cloud.google.com/go@v0.121.2/.devcontainer/Dockerfile` | mcr.microsoft.com/devcontainers/go:${VARIANT} | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/argoproj/argo-workflows/v3@v3.7.11/.devcontainer/builder/Dockerfile` | mcr.microsoft.com/vscode/devcontainers/base:ubuntu-22.04 | 1 |  |  | multi-arch |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/argoproj/argo-workflows/v3@v3.7.11/Dockerfile` | gcr.io/distroless/static | 10 | 8737 |  |  |  | Unpinned base image: builder; Unpinned base image: builder; Unpinned base image: builder; Unpinned base image: gcr.io/distroless/static; Unpinned base image: argoexec-base; Unpinned base image: argoexec-base; Unpinned base image: gcr.io/distroless/static; Unpinned base image: gcr.io/distroless/static |
| `.gopath-loader/pkg/mod/github.com/argoproj/argo-workflows/v3@v3.7.11/Dockerfile.windows` | argoexec-base | 4 | Administrator |  |  |  | Unpinned base image: builder; Unpinned base image: argoexec-base |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.55.5/awstesting/sandbox/Dockerfile.golang-tip` | buildpack-deps:buster-scm | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.55.5/awstesting/sandbox/Dockerfile.test.go1.10` | golang:1.10 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.55.5/awstesting/sandbox/Dockerfile.test.go1.11` | golang:1.11 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.55.5/awstesting/sandbox/Dockerfile.test.go1.12` | golang:1.12 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.55.5/awstesting/sandbox/Dockerfile.test.go1.13` | golang:1.13 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.55.5/awstesting/sandbox/Dockerfile.test.go1.14` | golang:1.14 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.55.5/awstesting/sandbox/Dockerfile.test.go1.15` | golang:1.15 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.55.5/awstesting/sandbox/Dockerfile.test.go1.16` | golang:1.16 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.55.5/awstesting/sandbox/Dockerfile.test.go1.17` | golang:1.17 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.55.5/awstesting/sandbox/Dockerfile.test.go1.18` | golang:1.18 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.55.5/awstesting/sandbox/Dockerfile.test.go1.19` | golang:1.19 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.55.5/awstesting/sandbox/Dockerfile.test.go1.5` | golang:1.5 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.55.5/awstesting/sandbox/Dockerfile.test.go1.6` | golang:1.6 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.55.5/awstesting/sandbox/Dockerfile.test.go1.7` | golang:1.7 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.55.5/awstesting/sandbox/Dockerfile.test.go1.8` | golang:1.8 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.55.5/awstesting/sandbox/Dockerfile.test.go1.9` | golang:1.9 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.55.5/awstesting/sandbox/Dockerfile.test.gotip` | aws-golang:tip | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/envoyproxy/protoc-gen-validate@v1.3.0/Dockerfile` | ubuntu:focal | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/google/addlicense@v0.0.0-20200906110928-a0294312aa76/Dockerfile` | golang:1-alpine | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.27.1/.github/Dockerfile` | golang:1.24.4 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/.circleci/Dockerfile` | golang:1.15.3 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/mattn/go-sqlite3@v1.14.34/_example/simple/Dockerfile` | scratch | 2 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/pelletier/go-toml/v2@v2.2.4/Dockerfile` | scratch | 1 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/prometheus/client_golang@v1.23.2/Dockerfile` | quay.io/prometheus/busybox:latest | 2 |  |  |  |  | Unpinned base image: quay.io/prometheus/busybox:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.7.linux-amd64/src/crypto/internal/boring/Dockerfile` | $ubuntu:focal | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.7.linux-amd64/src/crypto/internal/fips140/nistec/fiat/Dockerfile` | coqorg/coq:8.13.2 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/net@v0.51.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/sys@v0.41.0/unix/linux/Dockerfile` | ubuntu:25.04 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/k8s.io/apiextensions-apiserver@v0.35.2/artifacts/simple-image/Dockerfile` | gcr.io/distroless/base-debian10:latest | 1 |  |  |  |  | Unpinned base image: gcr.io/distroless/base-debian10:latest; No USER directive found (defaults to root) |
| `backend/Dockerfile` | registry.access.redhat.com/ubi9/ubi-minimal:9.5 | 3 | 1001 |  | multi-arch | yes | Unpinned base image: registry.access.redhat.com/ubi9/python-311 |
| `backend/Dockerfile.cacheserver` | alpine:3.21 | 2 | appuser |  |  |  |  |
| `backend/Dockerfile.conformance` | alpine:3.21 | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `backend/Dockerfile.driver` | registry.access.redhat.com/ubi9/ubi-minimal:9.5 | 2 | 65534 |  | multi-arch | yes |  |
| `backend/Dockerfile.konflux.api` | registry.access.redhat.com/ubi9/ubi-minimal@sha256:12db9874bd753eb98b1ab3d840e75de5d6842ac0604fbd68c012adefe97140be | 2 | 1001 |  |  |  |  |
| `backend/Dockerfile.konflux.driver` | registry.access.redhat.com/ubi9/ubi-minimal@sha256:12db9874bd753eb98b1ab3d840e75de5d6842ac0604fbd68c012adefe97140be | 2 | root |  | multi-arch |  | Container runs as root user |
| `backend/Dockerfile.konflux.launcher` | registry.access.redhat.com/ubi9/ubi-minimal@sha256:12db9874bd753eb98b1ab3d840e75de5d6842ac0604fbd68c012adefe97140be | 2 | root |  | multi-arch |  | Container runs as root user |
| `backend/Dockerfile.konflux.persistenceagent` | registry.access.redhat.com/ubi9/ubi-minimal@sha256:12db9874bd753eb98b1ab3d840e75de5d6842ac0604fbd68c012adefe97140be | 2 | root |  |  |  | Container runs as root user |
| `backend/Dockerfile.konflux.scheduledworkflow` | registry.access.redhat.com/ubi9/ubi-minimal@sha256:12db9874bd753eb98b1ab3d840e75de5d6842ac0604fbd68c012adefe97140be | 2 | root |  |  |  | Container runs as root user |
| `backend/Dockerfile.launcher` | registry.access.redhat.com/ubi9/ubi-minimal:9.5 | 2 | 65534 |  | multi-arch | yes |  |
| `backend/Dockerfile.persistenceagent` | registry.access.redhat.com/ubi9/ubi-minimal:9.5 | 2 | root |  | multi-arch | yes | Container runs as root user |
| `backend/Dockerfile.scheduledworkflow` | registry.access.redhat.com/ubi9/ubi-minimal:9.5 | 2 | root |  | multi-arch | yes | Container runs as root user |
| `backend/Dockerfile.viewercontroller` | alpine:3.21 | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `backend/Dockerfile.visualization` | python:3.11-slim | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `backend/api/Dockerfile` | golang:1.25.7 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `backend/metadata_writer/Dockerfile` | python:3.11-slim | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `backend/src/cache/deployer/Dockerfile` | gcr.io/google.com/cloudsdktool/google-cloud-cli:alpine | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `components/aws/athena/Dockerfile` | ubuntu:16.04 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `components/aws/emr/Dockerfile` | ubuntu:16.04 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `components/aws/sagemaker/Dockerfile` | public.ecr.aws/amazonlinux/amazonlinux:2 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `components/google-cloud/Dockerfile` | marketplace.gcr.io/google/ubuntu2404:latest | 1 |  |  |  |  | Unpinned base image: marketplace.gcr.io/google/ubuntu2404:latest; No USER directive found (defaults to root) |
| `components/kserve/Dockerfile` | python:3.11-slim-bullseye | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `frontend/Dockerfile` | node:${NODE_VERSION}-${BASE_IMAGE} | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `manifests/gcp_marketplace/deployer/Dockerfile` | gcr.io/cloud-marketplace-tools/k8s/deployer_helm/onbuild:0.11.3 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `proxy/Dockerfile` | gcr.io/inverting-proxy/agent@sha256:694d6c1bf299585b530c923c3728cd2c45083f3b396ec83ff799cef1c9dc7474 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `test_data/sdk_compiled_pipelines/valid/critical/modelcar/Dockerfile` | alpine:3.21 | 2 | 0 |  |  |  | Container runs as root user |
| `third_party/metadata_envoy/Dockerfile` | envoyproxy/envoy:v1.37.0 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `third_party/minio/Dockerfile` | minio/minio:RELEASE.2019-08-14T20-37-41Z | 2 |  |  |  |  | Unpinned base image: alpine; No USER directive found (defaults to root) |
| `third_party/ml-metadata/Dockerfile` | gcr.io/tfx-oss-public/ml_metadata_store_server:1.14.0 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `tools/bazel_builder/Dockerfile` | gcr.io/cloud-marketplace/google/rbe-ubuntu16-04@sha256:69c9f1652941d64a46f6f7358a44c1718f25caa5cb1ced4a58ccc5281cd183b5 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `tools/commit_checker/Dockerfile` | registry.access.redhat.com/ubi9/ubi-minimal:9.5 | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `tools/kind/Dockerfile.webhook-proxy` | registry.access.redhat.com/ubi9/nginx-124 | 1 |  |  |  |  | Unpinned base image: registry.access.redhat.com/ubi9/nginx-124; No USER directive found (defaults to root) |

