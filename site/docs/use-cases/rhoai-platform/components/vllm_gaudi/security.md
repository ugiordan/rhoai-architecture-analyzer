# vllm-gaudi: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `.cd/Dockerfile.rhel.ubi.vllm` | gaudi-pytorch | 3 | 2000 |  |  |  | Unpinned base image: gaudi-base; Unpinned base image: gaudi-pytorch |
| `.cd/Dockerfile.ubuntu.pytorch.vllm` | ${DOCKER_URL}/${VERSION}/${BASE_NAME}/${REPO_TYPE}/pytorch-${TORCH_TYPE_SUFFIX}installer-${PT_VERSION}:${REVISION} | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `.cd/Dockerfile.ubuntu.pytorch.vllm.nixl.latest` | ${DOCKER_URL}/${VERSION}/${BASE_NAME}/${REPO_TYPE}/pytorch-${TORCH_TYPE_SUFFIX}installer-${PT_VERSION}:${REVISION} | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `Dockerfile.hpu.ubi` | gaudi-pytorch | 3 | 2000 |  |  |  | Unpinned base image: gaudi-base; Unpinned base image: gaudi-pytorch |
| `Dockerfile.konflux.gaudi` | gaudi-pytorch | 3 | 2000 |  |  |  | Unpinned base image: gaudi-base; Unpinned base image: gaudi-pytorch |

