# vllm-orchestrator-gateway: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `Dockerfile` | ${UBI_MINIMAL_BASE_IMAGE}:${UBI_BASE_IMAGE_TAG} | 6 |  |  |  |  | Unpinned base image: rust-builder; Unpinned base image: gateway-builder; Unpinned base image: gateway-builder; Unpinned base image: gateway-builder; No USER directive found (defaults to root) |
| `Dockerfile.konflux` | ${UBI_MINIMAL_BASE_IMAGE}:${UBI_BASE_IMAGE_TAG} | 6 |  |  |  |  | Unpinned base image: rust-builder; Unpinned base image: gateway-builder; Unpinned base image: gateway-builder; Unpinned base image: gateway-builder; No USER directive found (defaults to root) |

