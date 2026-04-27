# kube-auth-proxy: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

### Secrets Referenced

| Name | Type | Referenced By |
|------|------|---------------|
| kube-auth-proxy-secret | Opaque | deployment/kube-auth-proxy |
| kube-rbac-proxy-client-certificates | Opaque | deployment/kube-rbac-proxy |

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| example-app | example-app | true | ? | ? | [`examples/oidc/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/examples/oidc/deployment.yaml) |
| example-app | example-app | true | ? | ? | [`examples/openshift/manual/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/examples/openshift/manual/deployment.yaml) |
| example-app | example-app | true | ? | ? | [`examples/openshift/service-account/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/examples/openshift/service-account/deployment.yaml) |
| kube-auth-proxy | kube-auth-proxy | true | ? | ? | [`examples/oidc/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/examples/oidc/deployment.yaml) |
| kube-auth-proxy | kube-auth-proxy | true | ? | ? | [`examples/openshift/manual/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/examples/openshift/manual/deployment.yaml) |
| kube-auth-proxy | kube-auth-proxy | true | ? | ? | [`examples/openshift/service-account/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/examples/openshift/service-account/deployment.yaml) |
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`kube-rbac-proxy/test/e2e/basics/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/test/e2e/basics/deployment.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`kube-rbac-proxy/test/e2e/basics/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/test/e2e/basics/deployment.yaml) |
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`kube-rbac-proxy/test/e2e/tokenrequest/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/test/e2e/tokenrequest/deployment.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`kube-rbac-proxy/test/e2e/tokenrequest/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/test/e2e/tokenrequest/deployment.yaml) |
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`kube-rbac-proxy/examples/oidc/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/examples/oidc/deployment.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`kube-rbac-proxy/examples/oidc/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/examples/oidc/deployment.yaml) |
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`kube-rbac-proxy/examples/resource-attributes/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/examples/resource-attributes/deployment.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`kube-rbac-proxy/examples/resource-attributes/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/examples/resource-attributes/deployment.yaml) |
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`kube-rbac-proxy/examples/rewrites/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/examples/rewrites/deployment.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`kube-rbac-proxy/examples/rewrites/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/examples/rewrites/deployment.yaml) |
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`kube-rbac-proxy/examples/static-auth/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/examples/static-auth/deployment.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`kube-rbac-proxy/examples/static-auth/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/examples/static-auth/deployment.yaml) |
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`kube-rbac-proxy/test/e2e/http2/deployment-no-http2.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/test/e2e/http2/deployment-no-http2.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`kube-rbac-proxy/test/e2e/http2/deployment-no-http2.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/test/e2e/http2/deployment-no-http2.yaml) |
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`kube-rbac-proxy/test/e2e/allowpaths/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/test/e2e/allowpaths/deployment.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`kube-rbac-proxy/test/e2e/allowpaths/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/test/e2e/allowpaths/deployment.yaml) |
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`kube-rbac-proxy/examples/non-resource-url/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/examples/non-resource-url/deployment.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`kube-rbac-proxy/examples/non-resource-url/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/examples/non-resource-url/deployment.yaml) |
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`kube-rbac-proxy/test/e2e/clientcertificates/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/test/e2e/clientcertificates/deployment.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`kube-rbac-proxy/test/e2e/clientcertificates/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/test/e2e/clientcertificates/deployment.yaml) |
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`kube-rbac-proxy/test/e2e/h2c-upstream/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/test/e2e/h2c-upstream/deployment.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`kube-rbac-proxy/test/e2e/h2c-upstream/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/test/e2e/h2c-upstream/deployment.yaml) |
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`kube-rbac-proxy/test/e2e/hardcoded_authorizer/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/test/e2e/hardcoded_authorizer/deployment.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`kube-rbac-proxy/test/e2e/hardcoded_authorizer/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/test/e2e/hardcoded_authorizer/deployment.yaml) |
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`kube-rbac-proxy/test/e2e/http2/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/test/e2e/http2/deployment.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`kube-rbac-proxy/test/e2e/http2/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/test/e2e/http2/deployment.yaml) |
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`kube-rbac-proxy/test/e2e/ignorepaths/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/test/e2e/ignorepaths/deployment.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`kube-rbac-proxy/test/e2e/ignorepaths/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/test/e2e/ignorepaths/deployment.yaml) |
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`kube-rbac-proxy/test/e2e/static-auth/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/test/e2e/static-auth/deployment.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`kube-rbac-proxy/test/e2e/static-auth/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/test/e2e/static-auth/deployment.yaml) |
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`kube-rbac-proxy/test/e2e/tokenmasking/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/test/e2e/tokenmasking/deployment.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`kube-rbac-proxy/test/e2e/tokenmasking/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/test/e2e/tokenmasking/deployment.yaml) |
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`kube-rbac-proxy/examples/non-resource-url-token-request/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/examples/non-resource-url-token-request/deployment.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`kube-rbac-proxy/examples/non-resource-url-token-request/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/examples/non-resource-url-token-request/deployment.yaml) |
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`kube-rbac-proxy/test/e2e/clientcertificates/deployment-wrongca.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/test/e2e/clientcertificates/deployment-wrongca.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`kube-rbac-proxy/test/e2e/clientcertificates/deployment-wrongca.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/test/e2e/clientcertificates/deployment-wrongca.yaml) |
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`kube-rbac-proxy/test/e2e/flags/deployment-logtostderr.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/test/e2e/flags/deployment-logtostderr.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`kube-rbac-proxy/test/e2e/flags/deployment-logtostderr.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/test/e2e/flags/deployment-logtostderr.yaml) |
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`kube-rbac-proxy/test/e2e/flags/deployment-other-flags.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/test/e2e/flags/deployment-other-flags.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`kube-rbac-proxy/test/e2e/flags/deployment-other-flags.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/test/e2e/flags/deployment-other-flags.yaml) |
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`kube-rbac-proxy/test/e2e/flags/deployment-short-timeout.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/test/e2e/flags/deployment-short-timeout.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`kube-rbac-proxy/test/e2e/flags/deployment-short-timeout.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/test/e2e/flags/deployment-short-timeout.yaml) |
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`kube-rbac-proxy/test/e2e/flags/deployment-upstream-timeout.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/test/e2e/flags/deployment-upstream-timeout.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`kube-rbac-proxy/test/e2e/flags/deployment-upstream-timeout.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/test/e2e/flags/deployment-upstream-timeout.yaml) |
| kube-rbac-proxy-verb-override | kube-rbac-proxy | ? | ? | ? | [`kube-rbac-proxy/examples/verb-override/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/examples/verb-override/deployment.yaml) |
| kube-rbac-proxy-verb-override | prometheus-example-app | ? | ? | ? | [`kube-rbac-proxy/examples/verb-override/deployment.yaml`](https://github.com/opendatahub-io/kube-auth-proxy/blob/d1acdea5345235b858bf1847fa73961abc023ff6/kube-rbac-proxy/examples/verb-override/deployment.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `.devcontainer/Dockerfile` | mcr.microsoft.com/vscode/devcontainers/go:1-1.23 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `Dockerfile` | ${RUNTIME_IMAGE} | 2 |  |  | multi-arch |  | Unpinned base image: ${BUILD_IMAGE}; Unpinned base image: ${RUNTIME_IMAGE}; No USER directive found (defaults to root) |
| `Dockerfile.konflux` | registry.access.redhat.com/ubi9/ubi-minimal@sha256:7d4e47500f28ac3a2bff06c25eff9127ff21048538ae03ce240d57cf756acd00 | 2 | 1001 |  | multi-arch |  |  |
| `Dockerfile.redhat` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 2 | 1001 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |
| `kube-rbac-proxy/Dockerfile` | $BASEIMAGE | 1 | 65532:65532 |  |  |  | Unpinned base image: $BASEIMAGE |
| `kube-rbac-proxy/Dockerfile.ocp` | registry.ci.openshift.org/ocp/4.20:base-rhel9 | 2 | 65534 |  |  |  |  |
| `kube-rbac-proxy/examples/example-client-http2/Dockerfile` |  | 0 |  |  |  |  | No USER directive found (defaults to root) |
| `kube-rbac-proxy/examples/example-client-urlquery/Dockerfile` | alpine:3.7 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `kube-rbac-proxy/examples/example-client/Dockerfile` | alpine:3.12 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `kube-rbac-proxy/examples/grpcc/Dockerfile` | node:8.9.4-alpine | 1 |  |  |  |  | No USER directive found (defaults to root) |

