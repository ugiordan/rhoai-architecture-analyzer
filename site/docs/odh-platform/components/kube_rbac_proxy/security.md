# kube-rbac-proxy: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`examples/non-resource-url/deployment.yaml`](https://github.com/brancz/kube-rbac-proxy/blob/d1195a883e22af75d26a0dd7e31e6172c659f81c/examples/non-resource-url/deployment.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`examples/non-resource-url/deployment.yaml`](https://github.com/brancz/kube-rbac-proxy/blob/d1195a883e22af75d26a0dd7e31e6172c659f81c/examples/non-resource-url/deployment.yaml) |
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`examples/non-resource-url-token-request/deployment.yaml`](https://github.com/brancz/kube-rbac-proxy/blob/d1195a883e22af75d26a0dd7e31e6172c659f81c/examples/non-resource-url-token-request/deployment.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`examples/non-resource-url-token-request/deployment.yaml`](https://github.com/brancz/kube-rbac-proxy/blob/d1195a883e22af75d26a0dd7e31e6172c659f81c/examples/non-resource-url-token-request/deployment.yaml) |
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`examples/oidc/deployment.yaml`](https://github.com/brancz/kube-rbac-proxy/blob/d1195a883e22af75d26a0dd7e31e6172c659f81c/examples/oidc/deployment.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`examples/oidc/deployment.yaml`](https://github.com/brancz/kube-rbac-proxy/blob/d1195a883e22af75d26a0dd7e31e6172c659f81c/examples/oidc/deployment.yaml) |
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`examples/resource-attributes/deployment.yaml`](https://github.com/brancz/kube-rbac-proxy/blob/d1195a883e22af75d26a0dd7e31e6172c659f81c/examples/resource-attributes/deployment.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`examples/resource-attributes/deployment.yaml`](https://github.com/brancz/kube-rbac-proxy/blob/d1195a883e22af75d26a0dd7e31e6172c659f81c/examples/resource-attributes/deployment.yaml) |
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`examples/rewrites/deployment.yaml`](https://github.com/brancz/kube-rbac-proxy/blob/d1195a883e22af75d26a0dd7e31e6172c659f81c/examples/rewrites/deployment.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`examples/rewrites/deployment.yaml`](https://github.com/brancz/kube-rbac-proxy/blob/d1195a883e22af75d26a0dd7e31e6172c659f81c/examples/rewrites/deployment.yaml) |
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`examples/static-auth/deployment.yaml`](https://github.com/brancz/kube-rbac-proxy/blob/d1195a883e22af75d26a0dd7e31e6172c659f81c/examples/static-auth/deployment.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`examples/static-auth/deployment.yaml`](https://github.com/brancz/kube-rbac-proxy/blob/d1195a883e22af75d26a0dd7e31e6172c659f81c/examples/static-auth/deployment.yaml) |
| kube-rbac-proxy | kube-rbac-proxy | ? | ? | ? | [`test/kubetest/testtemplates/data/deployment.yaml`](https://github.com/brancz/kube-rbac-proxy/blob/d1195a883e22af75d26a0dd7e31e6172c659f81c/test/kubetest/testtemplates/data/deployment.yaml) |
| kube-rbac-proxy | prometheus-example-app | ? | ? | ? | [`test/kubetest/testtemplates/data/deployment.yaml`](https://github.com/brancz/kube-rbac-proxy/blob/d1195a883e22af75d26a0dd7e31e6172c659f81c/test/kubetest/testtemplates/data/deployment.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `Dockerfile` | $BASEIMAGE | 1 | 65532:65532 |  |  |  | Unpinned base image: $BASEIMAGE |
| `examples/example-client-http2/Dockerfile` |  | 0 |  |  |  |  | No USER directive found (defaults to root) |
| `examples/example-client-urlquery/Dockerfile` | alpine:3.7 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `examples/example-client/Dockerfile` | alpine:3.20 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `examples/grpcc/Dockerfile` | node:8.9.4-alpine | 1 |  |  |  |  | No USER directive found (defaults to root) |

