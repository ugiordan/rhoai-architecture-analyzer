# kube-rbac-proxy: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`examples/non-resource-url/deployment.yaml`](https://github.com/brancz/kube-rbac-proxy/blob/aa08563309397def0fbb015f7bc6dd0a3e1ec856/examples/non-resource-url/deployment.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`examples/non-resource-url/deployment.yaml`](https://github.com/brancz/kube-rbac-proxy/blob/aa08563309397def0fbb015f7bc6dd0a3e1ec856/examples/non-resource-url/deployment.yaml) |
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`examples/non-resource-url-token-request/deployment.yaml`](https://github.com/brancz/kube-rbac-proxy/blob/aa08563309397def0fbb015f7bc6dd0a3e1ec856/examples/non-resource-url-token-request/deployment.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`examples/non-resource-url-token-request/deployment.yaml`](https://github.com/brancz/kube-rbac-proxy/blob/aa08563309397def0fbb015f7bc6dd0a3e1ec856/examples/non-resource-url-token-request/deployment.yaml) |
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`examples/oidc/deployment.yaml`](https://github.com/brancz/kube-rbac-proxy/blob/aa08563309397def0fbb015f7bc6dd0a3e1ec856/examples/oidc/deployment.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`examples/oidc/deployment.yaml`](https://github.com/brancz/kube-rbac-proxy/blob/aa08563309397def0fbb015f7bc6dd0a3e1ec856/examples/oidc/deployment.yaml) |
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`examples/resource-attributes/deployment.yaml`](https://github.com/brancz/kube-rbac-proxy/blob/aa08563309397def0fbb015f7bc6dd0a3e1ec856/examples/resource-attributes/deployment.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`examples/resource-attributes/deployment.yaml`](https://github.com/brancz/kube-rbac-proxy/blob/aa08563309397def0fbb015f7bc6dd0a3e1ec856/examples/resource-attributes/deployment.yaml) |
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`examples/rewrites/deployment.yaml`](https://github.com/brancz/kube-rbac-proxy/blob/aa08563309397def0fbb015f7bc6dd0a3e1ec856/examples/rewrites/deployment.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`examples/rewrites/deployment.yaml`](https://github.com/brancz/kube-rbac-proxy/blob/aa08563309397def0fbb015f7bc6dd0a3e1ec856/examples/rewrites/deployment.yaml) |
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`examples/static-auth/deployment.yaml`](https://github.com/brancz/kube-rbac-proxy/blob/aa08563309397def0fbb015f7bc6dd0a3e1ec856/examples/static-auth/deployment.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`examples/static-auth/deployment.yaml`](https://github.com/brancz/kube-rbac-proxy/blob/aa08563309397def0fbb015f7bc6dd0a3e1ec856/examples/static-auth/deployment.yaml) |
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`test/kubetest/testtemplates/data/deployment.yaml`](https://github.com/brancz/kube-rbac-proxy/blob/aa08563309397def0fbb015f7bc6dd0a3e1ec856/test/kubetest/testtemplates/data/deployment.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`test/kubetest/testtemplates/data/deployment.yaml`](https://github.com/brancz/kube-rbac-proxy/blob/aa08563309397def0fbb015f7bc6dd0a3e1ec856/test/kubetest/testtemplates/data/deployment.yaml) |
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`test/e2e/flags/deployment-short-timeout.yaml`](https://github.com/brancz/kube-rbac-proxy/blob/aa08563309397def0fbb015f7bc6dd0a3e1ec856/test/e2e/flags/deployment-short-timeout.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`test/e2e/flags/deployment-short-timeout.yaml`](https://github.com/brancz/kube-rbac-proxy/blob/aa08563309397def0fbb015f7bc6dd0a3e1ec856/test/e2e/flags/deployment-short-timeout.yaml) |
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`test/e2e/flags/deployment-upstream-timeout.yaml`](https://github.com/brancz/kube-rbac-proxy/blob/aa08563309397def0fbb015f7bc6dd0a3e1ec856/test/e2e/flags/deployment-upstream-timeout.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`test/e2e/flags/deployment-upstream-timeout.yaml`](https://github.com/brancz/kube-rbac-proxy/blob/aa08563309397def0fbb015f7bc6dd0a3e1ec856/test/e2e/flags/deployment-upstream-timeout.yaml) |
| kube-rbac-proxy-verb-override | kube-rbac-proxy | true | true | ? | [`examples/verb-override/deployment.yaml`](https://github.com/brancz/kube-rbac-proxy/blob/aa08563309397def0fbb015f7bc6dd0a3e1ec856/examples/verb-override/deployment.yaml) |
| kube-rbac-proxy-verb-override | prometheus-example-app | true | true | ? | [`examples/verb-override/deployment.yaml`](https://github.com/brancz/kube-rbac-proxy/blob/aa08563309397def0fbb015f7bc6dd0a3e1ec856/examples/verb-override/deployment.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `Dockerfile` | $BASEIMAGE | 1 | 65532:65532 |  |  |  | Unpinned base image: $BASEIMAGE |
| `Dockerfile.konflux` | registry.redhat.io/ubi9/ubi-minimal@sha256:7d4e47500f28ac3a2bff06c25eff9127ff21048538ae03ce240d57cf756acd00 | 2 | 65534 |  | multi-arch |  |  |
| `Dockerfile.ocp` | registry.ci.openshift.org/ocp/4.22:base-rhel9 | 2 | 65534 |  |  |  |  |
| `Dockerfile.redhat` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 2 | 65534 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |
| `examples/example-client-http2/Dockerfile` |  | 0 |  |  |  |  | No USER directive found (defaults to root) |
| `examples/example-client-urlquery/Dockerfile` | alpine:3.7 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `examples/example-client/Dockerfile` | alpine:3.12 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `examples/grpcc/Dockerfile` | node:8.9.4-alpine | 1 |  |  |  |  | No USER directive found (defaults to root) |

