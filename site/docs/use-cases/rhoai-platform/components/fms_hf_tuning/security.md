# fms-hf-tuning: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `build/Dockerfile` | release-base | 6 | ${USER} |  |  |  | Unpinned base image: base; Unpinned base image: base; Unpinned base image: cuda-base; Unpinned base image: cuda-devel; Unpinned base image: release-base |

