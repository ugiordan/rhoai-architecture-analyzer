# kserve: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

### Secrets Referenced

| Name | Type | Referenced By |
|------|------|---------------|
| epp-metrics-token | Opaque | deployment/controller-manager |
| hf-token | Opaque | deployment/llama-deployment |
| kedaorg-certs | Opaque | deployment/keda-metrics-apiserver, deployment/keda-operator |
| kserve-webhook-server-cert | Opaque | deployment/kserve-controller-manager |
| llmisvc-webhook-server-cert | Opaque | deployment/llmisvc-controller-manager |
| localmodel-webhook-server-cert | Opaque | deployment/kserve-localmodel-controller-manager |
| opentelemetry-operator-controller-manager-service-cert | Opaque | deployment/controller-manager |
| opentelemetry-operator-metrics | Opaque | deployment/controller-manager |
| prometheus-client-cert | Opaque | deployment/controller-manager |
| webhook-server-cert | Opaque | deployment/controller-manager |

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| keda-metrics-apiserver | keda-metrics-apiserver | true | true | ? | [`.gomod-cache/github.com/kedacore/keda/v2@v2.17.3/config/metrics-server/deployment.yaml`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/.gomod-cache/github.com/kedacore/keda/v2@v2.17.3/config/metrics-server/deployment.yaml) |
| keda-metrics-apiserver | keda-metrics-apiserver | true | true | ? | [`.gopath-loader/pkg/mod/github.com/kedacore/keda/v2@v2.17.3/config/metrics-server/deployment.yaml`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/.gopath-loader/pkg/mod/github.com/kedacore/keda/v2@v2.17.3/config/metrics-server/deployment.yaml) |
| keda-operator | keda-operator | true | true | ? | [`.gomod-cache/github.com/kedacore/keda/v2@v2.17.3/config/manager/manager.yaml`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/.gomod-cache/github.com/kedacore/keda/v2@v2.17.3/config/manager/manager.yaml) |
| keda-operator | keda-operator | true | true | ? | [`.gopath-loader/pkg/mod/github.com/kedacore/keda/v2@v2.17.3/config/manager/manager.yaml`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/.gopath-loader/pkg/mod/github.com/kedacore/keda/v2@v2.17.3/config/manager/manager.yaml) |
| kserve-controller-manager | manager | true | true | false | [`kustomize:config/overlays/odh`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/kustomize:config/overlays/odh) |
| kserve-localmodel-controller-manager | manager | true | true | false | [`config/localmodels/manager.yaml`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/config/localmodels/manager.yaml) |
| llama-deployment | inference-server | ? | ? | ? | [`.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/tools/dynamic-lora-sidecar/deployment.yaml`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/tools/dynamic-lora-sidecar/deployment.yaml) |
| llama-deployment | inference-server | ? | ? | ? | [`.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/tools/dynamic-lora-sidecar/deployment.yaml`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/tools/dynamic-lora-sidecar/deployment.yaml) |
| llmisvc-controller-manager | manager | true | true | false | [`kustomize:config/overlays/odh`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/kustomize:config/overlays/odh) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `.gomod-cache/cloud.google.com/go@v0.120.0/.devcontainer/Dockerfile` | mcr.microsoft.com/devcontainers/go:${VARIANT} | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.55.6/awstesting/sandbox/Dockerfile.golang-tip` | buildpack-deps:buster-scm | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.55.6/awstesting/sandbox/Dockerfile.test.go1.10` | golang:1.10 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.55.6/awstesting/sandbox/Dockerfile.test.go1.11` | golang:1.11 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.55.6/awstesting/sandbox/Dockerfile.test.go1.12` | golang:1.12 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.55.6/awstesting/sandbox/Dockerfile.test.go1.13` | golang:1.13 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.55.6/awstesting/sandbox/Dockerfile.test.go1.14` | golang:1.14 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.55.6/awstesting/sandbox/Dockerfile.test.go1.15` | golang:1.15 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.55.6/awstesting/sandbox/Dockerfile.test.go1.16` | golang:1.16 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.55.6/awstesting/sandbox/Dockerfile.test.go1.17` | golang:1.17 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.55.6/awstesting/sandbox/Dockerfile.test.go1.18` | golang:1.18 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.55.6/awstesting/sandbox/Dockerfile.test.go1.19` | golang:1.19 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.55.6/awstesting/sandbox/Dockerfile.test.go1.5` | golang:1.5 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.55.6/awstesting/sandbox/Dockerfile.test.go1.6` | golang:1.6 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.55.6/awstesting/sandbox/Dockerfile.test.go1.7` | golang:1.7 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.55.6/awstesting/sandbox/Dockerfile.test.go1.8` | golang:1.8 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.55.6/awstesting/sandbox/Dockerfile.test.go1.9` | golang:1.9 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/aws/aws-sdk-go@v1.55.6/awstesting/sandbox/Dockerfile.test.gotip` | aws-golang:tip | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/docker@v28.4.0+incompatible/Dockerfile` | dev-base | 71 |  |  | multi-arch |  | Unpinned base image: busybox; Unpinned base image: scratch; Unpinned base image: ${GOLANG_IMAGE}; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; Unpinned base image: binary-dummy; Unpinned base image: delve-${DELVE_SUPPORTED}; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; Unpinned base image: containerd-build; Unpinned base image: binary-dummy; Unpinned base image: containerd-${TARGETOS}; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; Unpinned base image: runc-build; Unpinned base image: binary-dummy; Unpinned base image: runc-${TARGETOS}; Unpinned base image: base; Unpinned base image: base; Unpinned base image: tini-build; Unpinned base image: binary-dummy; Unpinned base image: tini-${TARGETOS}; Unpinned base image: base; Unpinned base image: base; Unpinned base image: rootlesskit-build; Unpinned base image: binary-dummy; Unpinned base image: rootlesskit-${TARGETOS}; Unpinned base image: base; Unpinned base image: scratch; Unpinned base image: scratch; Unpinned base image: scratch; Unpinned base image: scratch; Unpinned base image: scratch; Unpinned base image: scratch; Unpinned base image: vpnkit-linux-${TARGETARCH}; Unpinned base image: vpnkit-${TARGETOS}; Unpinned base image: base; Unpinned base image: base; Unpinned base image: binary-dummy; Unpinned base image: containerutil-build; Unpinned base image: containerutil-windows-${TARGETARCH}; Unpinned base image: containerutil-${TARGETOS}; Unpinned base image: base; Unpinned base image: dev-systemd-false; Unpinned base image: dev-systemd-${SYSTEMD}; Unpinned base image: dev-systemd-true; Unpinned base image: dev-firewalld-${FIREWALLD}; Unpinned base image: base; Unpinned base image: scratch; Unpinned base image: scratch; Unpinned base image: base; Unpinned base image: dev-base; Unpinned base image: dev-base; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/docker@v28.4.0+incompatible/Dockerfile.simple` | ${GOLANG_IMAGE} | 1 |  |  |  |  | Unpinned base image: ${GOLANG_IMAGE}; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/docker@v28.4.0+incompatible/Dockerfile.windows` | ${WINDOWS_BASE_IMAGE}:${WINDOWS_BASE_IMAGE_TAG} | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/docker@v28.4.0+incompatible/contrib/busybox/Dockerfile` | ${WINDOWS_BASE_IMAGE}:${WINDOWS_BASE_IMAGE_TAG} | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/docker@v28.4.0+incompatible/contrib/httpserver/Dockerfile` | busybox | 1 |  |  |  |  | Unpinned base image: busybox; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/docker@v28.4.0+incompatible/contrib/nnp-test/Dockerfile` | debian:bookworm-slim | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/docker@v28.4.0+incompatible/contrib/syntax/nano/Dockerfile.nanorc` |  | 0 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/docker@v28.4.0+incompatible/contrib/syntax/textmate/Docker.tmbundle/Preferences/Dockerfile.tmPreferences` |  | 0 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/docker@v28.4.0+incompatible/contrib/syntax/textmate/Docker.tmbundle/Syntaxes/Dockerfile.tmLanguage` |  | 0 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/docker@v28.4.0+incompatible/contrib/syscall-test/Dockerfile` | debian:bookworm-slim | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/docker@v28.4.0+incompatible/libnetwork/cmd/diagnostic/Dockerfile.client` | alpine | 1 |  |  |  |  | Unpinned base image: alpine; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/docker@v28.4.0+incompatible/libnetwork/cmd/diagnostic/Dockerfile.dind` | docker:17.12-dind | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/docker@v28.4.0+incompatible/libnetwork/cmd/networkdb-test/Dockerfile` | alpine | 1 |  |  |  |  | Unpinned base image: alpine; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/docker@v28.4.0+incompatible/libnetwork/cmd/ssd/Dockerfile` | alpine:3.7 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/docker/docker@v28.4.0+incompatible/libnetwork/support/Dockerfile` | docker:18-dind | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/envoyproxy/protoc-gen-validate@v1.2.1/Dockerfile` | ubuntu:focal | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/grpc-ecosystem/grpc-gateway/v2@v2.27.3/.github/Dockerfile` | golang:1.25.1 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/kedacore/keda/v2@v2.17.3/.devcontainer/Dockerfile` | golang:1.23.8 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/kedacore/keda/v2@v2.17.3/Dockerfile` | gcr.io/distroless/static:nonroot | 2 | 65532:65532 |  | multi-arch |  |  |
| `.gomod-cache/github.com/kedacore/keda/v2@v2.17.3/Dockerfile.adapter` | gcr.io/distroless/static:nonroot | 2 | 65532:65532 |  | multi-arch |  |  |
| `.gomod-cache/github.com/kedacore/keda/v2@v2.17.3/Dockerfile.webhooks` | gcr.io/distroless/static:nonroot | 2 | 65532:65532 |  | multi-arch |  |  |
| `.gomod-cache/github.com/llm-d/llm-d-workload-variant-autoscaler@v0.6.0/Dockerfile` | gcr.io/distroless/static:nonroot | 2 | 65532:65532 |  | multi-arch |  |  |
| `.gomod-cache/github.com/open-telemetry/opentelemetry-operator@v0.113.0/Dockerfile` | scratch | 2 | 65532:65532 |  | multi-arch |  | Unpinned base image: scratch |
| `.gomod-cache/github.com/open-telemetry/opentelemetry-operator@v0.113.0/autoinstrumentation/apache-httpd/Dockerfile` | alpine:latest | 2 |  |  |  |  | Unpinned base image: alpine:latest; Unpinned base image: alpine:latest; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/open-telemetry/opentelemetry-operator@v0.113.0/autoinstrumentation/dotnet/Dockerfile` | busybox | 2 |  |  |  |  | Unpinned base image: busybox; Unpinned base image: busybox; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/open-telemetry/opentelemetry-operator@v0.113.0/autoinstrumentation/java/Dockerfile` | busybox | 1 |  |  |  |  | Unpinned base image: busybox; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/open-telemetry/opentelemetry-operator@v0.113.0/autoinstrumentation/nodejs/Dockerfile` | busybox | 2 |  |  |  |  | Unpinned base image: busybox; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/open-telemetry/opentelemetry-operator@v0.113.0/autoinstrumentation/python/Dockerfile` | busybox | 3 |  |  |  |  | Unpinned base image: busybox; No USER directive found (defaults to root) |
| `.gomod-cache/github.com/open-telemetry/opentelemetry-operator@v0.113.0/cmd/gather/Dockerfile` | registry.access.redhat.com/ubi9-minimal:9.2 | 1 | 65532:65532 |  | multi-arch |  |  |
| `.gomod-cache/github.com/open-telemetry/opentelemetry-operator@v0.113.0/cmd/operator-opamp-bridge/Dockerfile` | scratch | 2 | 65532:65532 |  | multi-arch |  | Unpinned base image: scratch |
| `.gomod-cache/github.com/open-telemetry/opentelemetry-operator@v0.113.0/cmd/otel-allocator/Dockerfile` | scratch | 2 | 65532:65532 |  | multi-arch |  | Unpinned base image: scratch |
| `.gomod-cache/github.com/openshift/api@v0.0.0-20240124164020-e2ce40831f2e/Dockerfile.rhel8` | registry.ci.openshift.org/ocp/4.16:base-rhel9 | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/parquet-go/parquet-go@v0.27.0/.github/workflows/Dockerfile` | golang:1.24 | 1 |  |  | multi-arch |  | No USER directive found (defaults to root) |
| `.gomod-cache/github.com/prometheus/client_golang@v1.23.2/Dockerfile` | quay.io/prometheus/busybox:latest | 2 |  |  |  |  | Unpinned base image: quay.io/prometheus/busybox:latest; No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/toolchain@v0.0.1-go1.25.8.linux-amd64/src/crypto/internal/boring/Dockerfile` | $ubuntu:focal | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/toolchain@v0.0.1-go1.25.8.linux-amd64/src/crypto/internal/fips140/nistec/fiat/Dockerfile` | coqorg/coq:8.13.2 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/x/net@v0.49.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gomod-cache/golang.org/x/sys@v0.40.0/unix/linux/Dockerfile` | ubuntu:25.04 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/k8s.io/apiextensions-apiserver@v0.34.3/artifacts/simple-image/Dockerfile` | gcr.io/distroless/base-debian10:latest | 1 |  |  |  |  | Unpinned base image: gcr.io/distroless/base-debian10:latest; No USER directive found (defaults to root) |
| `.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/Dockerfile` | ${BASE_IMAGE} | 2 |  |  |  |  | Unpinned base image: ${BUILDER_IMAGE}; Unpinned base image: ${BASE_IMAGE}; No USER directive found (defaults to root) |
| `.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/tools/dynamic-lora-sidecar/Dockerfile` | python:3.10-slim-buster | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/sigs.k8s.io/gateway-api@v1.4.2-0.20260116062110-0d0ca872766e/docker/Dockerfile.echo-advanced` | gcr.io/istio-release/app:1.21.2 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gomod-cache/sigs.k8s.io/gateway-api@v1.4.2-0.20260116062110-0d0ca872766e/docker/Dockerfile.echo-basic` | gcr.io/distroless/static:nonroot | 2 | nonroot:nonroot |  |  |  |  |
| `.gomod-cache/sigs.k8s.io/lws@v0.8.0/Dockerfile` | ${BASE_IMAGE} | 2 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${BUILDER_IMAGE}; Unpinned base image: ${BASE_IMAGE} |
| `.gopath-loader/pkg/mod/cloud.google.com/go@v0.120.0/.devcontainer/Dockerfile` | mcr.microsoft.com/devcontainers/go:${VARIANT} | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.55.6/awstesting/sandbox/Dockerfile.golang-tip` | buildpack-deps:buster-scm | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.55.6/awstesting/sandbox/Dockerfile.test.go1.10` | golang:1.10 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.55.6/awstesting/sandbox/Dockerfile.test.go1.11` | golang:1.11 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.55.6/awstesting/sandbox/Dockerfile.test.go1.12` | golang:1.12 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.55.6/awstesting/sandbox/Dockerfile.test.go1.13` | golang:1.13 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.55.6/awstesting/sandbox/Dockerfile.test.go1.14` | golang:1.14 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.55.6/awstesting/sandbox/Dockerfile.test.go1.15` | golang:1.15 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.55.6/awstesting/sandbox/Dockerfile.test.go1.16` | golang:1.16 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.55.6/awstesting/sandbox/Dockerfile.test.go1.17` | golang:1.17 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.55.6/awstesting/sandbox/Dockerfile.test.go1.18` | golang:1.18 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.55.6/awstesting/sandbox/Dockerfile.test.go1.19` | golang:1.19 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.55.6/awstesting/sandbox/Dockerfile.test.go1.5` | golang:1.5 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.55.6/awstesting/sandbox/Dockerfile.test.go1.6` | golang:1.6 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.55.6/awstesting/sandbox/Dockerfile.test.go1.7` | golang:1.7 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.55.6/awstesting/sandbox/Dockerfile.test.go1.8` | golang:1.8 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.55.6/awstesting/sandbox/Dockerfile.test.go1.9` | golang:1.9 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/aws/aws-sdk-go@v1.55.6/awstesting/sandbox/Dockerfile.test.gotip` | aws-golang:tip | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/docker@v28.4.0+incompatible/Dockerfile` | dev-base | 71 |  |  | multi-arch |  | Unpinned base image: busybox; Unpinned base image: scratch; Unpinned base image: ${GOLANG_IMAGE}; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; Unpinned base image: binary-dummy; Unpinned base image: delve-${DELVE_SUPPORTED}; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; Unpinned base image: containerd-build; Unpinned base image: binary-dummy; Unpinned base image: containerd-${TARGETOS}; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; Unpinned base image: runc-build; Unpinned base image: binary-dummy; Unpinned base image: runc-${TARGETOS}; Unpinned base image: base; Unpinned base image: base; Unpinned base image: tini-build; Unpinned base image: binary-dummy; Unpinned base image: tini-${TARGETOS}; Unpinned base image: base; Unpinned base image: base; Unpinned base image: rootlesskit-build; Unpinned base image: binary-dummy; Unpinned base image: rootlesskit-${TARGETOS}; Unpinned base image: base; Unpinned base image: scratch; Unpinned base image: scratch; Unpinned base image: scratch; Unpinned base image: scratch; Unpinned base image: scratch; Unpinned base image: scratch; Unpinned base image: vpnkit-linux-${TARGETARCH}; Unpinned base image: vpnkit-${TARGETOS}; Unpinned base image: base; Unpinned base image: base; Unpinned base image: binary-dummy; Unpinned base image: containerutil-build; Unpinned base image: containerutil-windows-${TARGETARCH}; Unpinned base image: containerutil-${TARGETOS}; Unpinned base image: base; Unpinned base image: dev-systemd-false; Unpinned base image: dev-systemd-${SYSTEMD}; Unpinned base image: dev-systemd-true; Unpinned base image: dev-firewalld-${FIREWALLD}; Unpinned base image: base; Unpinned base image: scratch; Unpinned base image: scratch; Unpinned base image: base; Unpinned base image: dev-base; Unpinned base image: dev-base; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/docker@v28.4.0+incompatible/Dockerfile.simple` | ${GOLANG_IMAGE} | 1 |  |  |  |  | Unpinned base image: ${GOLANG_IMAGE}; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/docker@v28.4.0+incompatible/Dockerfile.windows` | ${WINDOWS_BASE_IMAGE}:${WINDOWS_BASE_IMAGE_TAG} | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/docker@v28.4.0+incompatible/contrib/busybox/Dockerfile` | ${WINDOWS_BASE_IMAGE}:${WINDOWS_BASE_IMAGE_TAG} | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/docker@v28.4.0+incompatible/contrib/httpserver/Dockerfile` | busybox | 1 |  |  |  |  | Unpinned base image: busybox; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/docker@v28.4.0+incompatible/contrib/nnp-test/Dockerfile` | debian:bookworm-slim | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/docker@v28.4.0+incompatible/contrib/syntax/nano/Dockerfile.nanorc` |  | 0 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/docker@v28.4.0+incompatible/contrib/syntax/textmate/Docker.tmbundle/Preferences/Dockerfile.tmPreferences` |  | 0 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/docker@v28.4.0+incompatible/contrib/syntax/textmate/Docker.tmbundle/Syntaxes/Dockerfile.tmLanguage` |  | 0 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/docker@v28.4.0+incompatible/contrib/syscall-test/Dockerfile` | debian:bookworm-slim | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/docker@v28.4.0+incompatible/libnetwork/cmd/diagnostic/Dockerfile.client` | alpine | 1 |  |  |  |  | Unpinned base image: alpine; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/docker@v28.4.0+incompatible/libnetwork/cmd/diagnostic/Dockerfile.dind` | docker:17.12-dind | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/docker@v28.4.0+incompatible/libnetwork/cmd/networkdb-test/Dockerfile` | alpine | 1 |  |  |  |  | Unpinned base image: alpine; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/docker@v28.4.0+incompatible/libnetwork/cmd/ssd/Dockerfile` | alpine:3.7 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/docker/docker@v28.4.0+incompatible/libnetwork/support/Dockerfile` | docker:18-dind | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/envoyproxy/protoc-gen-validate@v1.2.1/Dockerfile` | ubuntu:focal | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.27.3/.github/Dockerfile` | golang:1.25.1 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/kedacore/keda/v2@v2.17.3/.devcontainer/Dockerfile` | golang:1.23.8 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/kedacore/keda/v2@v2.17.3/Dockerfile` | gcr.io/distroless/static:nonroot | 2 | 65532:65532 |  | multi-arch |  |  |
| `.gopath-loader/pkg/mod/github.com/kedacore/keda/v2@v2.17.3/Dockerfile.adapter` | gcr.io/distroless/static:nonroot | 2 | 65532:65532 |  | multi-arch |  |  |
| `.gopath-loader/pkg/mod/github.com/kedacore/keda/v2@v2.17.3/Dockerfile.webhooks` | gcr.io/distroless/static:nonroot | 2 | 65532:65532 |  | multi-arch |  |  |
| `.gopath-loader/pkg/mod/github.com/llm-d/llm-d-workload-variant-autoscaler@v0.6.0/Dockerfile` | gcr.io/distroless/static:nonroot | 2 | 65532:65532 |  | multi-arch |  |  |
| `.gopath-loader/pkg/mod/github.com/open-telemetry/opentelemetry-operator@v0.113.0/Dockerfile` | scratch | 2 | 65532:65532 |  | multi-arch |  | Unpinned base image: scratch |
| `.gopath-loader/pkg/mod/github.com/open-telemetry/opentelemetry-operator@v0.113.0/autoinstrumentation/apache-httpd/Dockerfile` | alpine:latest | 2 |  |  |  |  | Unpinned base image: alpine:latest; Unpinned base image: alpine:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/open-telemetry/opentelemetry-operator@v0.113.0/autoinstrumentation/dotnet/Dockerfile` | busybox | 2 |  |  |  |  | Unpinned base image: busybox; Unpinned base image: busybox; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/open-telemetry/opentelemetry-operator@v0.113.0/autoinstrumentation/java/Dockerfile` | busybox | 1 |  |  |  |  | Unpinned base image: busybox; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/open-telemetry/opentelemetry-operator@v0.113.0/autoinstrumentation/nodejs/Dockerfile` | busybox | 2 |  |  |  |  | Unpinned base image: busybox; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/open-telemetry/opentelemetry-operator@v0.113.0/autoinstrumentation/python/Dockerfile` | busybox | 3 |  |  |  |  | Unpinned base image: busybox; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/open-telemetry/opentelemetry-operator@v0.113.0/cmd/gather/Dockerfile` | registry.access.redhat.com/ubi9-minimal:9.2 | 1 | 65532:65532 |  | multi-arch |  |  |
| `.gopath-loader/pkg/mod/github.com/open-telemetry/opentelemetry-operator@v0.113.0/cmd/operator-opamp-bridge/Dockerfile` | scratch | 2 | 65532:65532 |  | multi-arch |  | Unpinned base image: scratch |
| `.gopath-loader/pkg/mod/github.com/open-telemetry/opentelemetry-operator@v0.113.0/cmd/otel-allocator/Dockerfile` | scratch | 2 | 65532:65532 |  | multi-arch |  | Unpinned base image: scratch |
| `.gopath-loader/pkg/mod/github.com/openshift/api@v0.0.0-20240124164020-e2ce40831f2e/Dockerfile.rhel8` | registry.ci.openshift.org/ocp/4.16:base-rhel9 | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/parquet-go/parquet-go@v0.27.0/.github/workflows/Dockerfile` | golang:1.24 | 1 |  |  | multi-arch |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/github.com/prometheus/client_golang@v1.23.2/Dockerfile` | quay.io/prometheus/busybox:latest | 2 |  |  |  |  | Unpinned base image: quay.io/prometheus/busybox:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.8.linux-amd64/src/crypto/internal/boring/Dockerfile` | $ubuntu:focal | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.8.linux-amd64/src/crypto/internal/fips140/nistec/fiat/Dockerfile` | coqorg/coq:8.13.2 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/net@v0.49.0/internal/quic/cmd/interop/Dockerfile` | martenseemann/quic-network-simulator-endpoint:latest | 2 |  |  | multi-arch |  | Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; Unpinned base image: martenseemann/quic-network-simulator-endpoint:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/golang.org/x/sys@v0.40.0/unix/linux/Dockerfile` | ubuntu:25.04 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/k8s.io/apiextensions-apiserver@v0.34.3/artifacts/simple-image/Dockerfile` | gcr.io/distroless/base-debian10:latest | 1 |  |  |  |  | Unpinned base image: gcr.io/distroless/base-debian10:latest; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/Dockerfile` | ${BASE_IMAGE} | 2 |  |  |  |  | Unpinned base image: ${BUILDER_IMAGE}; Unpinned base image: ${BASE_IMAGE}; No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/tools/dynamic-lora-sidecar/Dockerfile` | python:3.10-slim-buster | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api@v1.4.2-0.20260116062110-0d0ca872766e/docker/Dockerfile.echo-advanced` | gcr.io/istio-release/app:1.21.2 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api@v1.4.2-0.20260116062110-0d0ca872766e/docker/Dockerfile.echo-basic` | gcr.io/distroless/static:nonroot | 2 | nonroot:nonroot |  |  |  |  |
| `.gopath-loader/pkg/mod/sigs.k8s.io/lws@v0.8.0/Dockerfile` | ${BASE_IMAGE} | 2 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${BUILDER_IMAGE}; Unpinned base image: ${BASE_IMAGE} |
| `Dockerfile` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 4 | 1000:1000 |  |  |  | Unpinned base image: deps; Unpinned base image: deps; Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |
| `tools/tf2openapi/Dockerfile` | gcr.io/distroless/static:nonroot | 2 |  |  |  |  | No USER directive found (defaults to root) |

