# fms-guardrails-orchestrator: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `Dockerfile.amd64` | ${UBI_MINIMAL_BASE_IMAGE}:${UBI_BASE_IMAGE_TAG} | 3 | orchestr8 |  |  |  | Unpinned base image: rust-builder |
| `Dockerfile.ppc64le` | ${UBI_MINIMAL_BASE_IMAGE}:${UBI_BASE_IMAGE_TAG} | 3 | orchestr8 |  |  |  | Unpinned base image: rust-builder |
| `Dockerfile.s390x` | ${UBI_MINIMAL_BASE_IMAGE}:${UBI_BASE_IMAGE_TAG} | 3 | orchestr8 |  |  |  | Unpinned base image: rust-builder |

