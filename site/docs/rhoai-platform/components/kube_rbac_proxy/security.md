# kube-rbac-proxy: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `Dockerfile` | $BASEIMAGE | 1 | 65532:65532 |  |  |  | Unpinned base image: $BASEIMAGE |
| `Dockerfile.konflux` | registry.redhat.io/ubi9/ubi-minimal@sha256:8d0a8fb39ec907e8ca62cdd24b62a63ca49a30fe465798a360741fde58437a23 | 2 | 65534 |  | multi-arch |  |  |
| `Dockerfile.ocp` | registry.ci.openshift.org/ocp/4.22:base-rhel9 | 2 | 65534 |  |  |  |  |
| `Dockerfile.redhat` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 2 | 65534 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |

