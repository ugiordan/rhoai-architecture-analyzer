# llama-stack: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `containers/Containerfile` | ${BASE_IMAGE} | 1 |  |  |  |  | Unpinned base image: ${BASE_IMAGE}; No USER directive found (defaults to root) |
| `src/ogx_ui/Containerfile` | node:22.5.1-alpine | 1 | nextjs |  |  |  |  |

