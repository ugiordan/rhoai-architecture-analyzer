# kube-auth-proxy: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `.devcontainer/Dockerfile` | mcr.microsoft.com/vscode/devcontainers/go:1-1.23 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `Dockerfile` | ${RUNTIME_IMAGE} | 2 |  |  | multi-arch |  | Unpinned base image: ${BUILD_IMAGE}; Unpinned base image: ${RUNTIME_IMAGE}; No USER directive found (defaults to root) |
| `Dockerfile.redhat` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 2 | 1001 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |
| `kube-rbac-proxy/Dockerfile` | $BASEIMAGE | 1 | 65532:65532 |  |  |  | Unpinned base image: $BASEIMAGE |
| `kube-rbac-proxy/Dockerfile.ocp` | registry.ci.openshift.org/ocp/4.20:base-rhel9 | 2 | 65534 |  |  |  |  |

