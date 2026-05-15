# llm-d: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| interactive-pod | benchmark-runner | ? | ? | ? | [`helpers/interactive-pod/manifests/deployment.yaml`](https://github.com/llm-d/llm-d/blob/3d04e73d481491695c0ffcdfee300628afb3f404/helpers/interactive-pod/manifests/deployment.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `docker/Dockerfile.cpu` | vllm-build | 5 |  |  | multi-arch |  | Unpinned base image: base-common; Unpinned base image: base-common; Unpinned base image: base-${TARGETARCH}; Unpinned base image: vllm-build; No USER directive found (defaults to root) |
| `docker/Dockerfile.cuda` | runtime | 4 | 2000 |  | multi-arch |  | Unpinned base image: runtime; Unpinned base image: runtime |
| `docker/Dockerfile.hpu` | ${DOCKER_URL}/${VERSION}/${BASE_NAME}/${REPO_TYPE}/pytorch-${TORCH_TYPE_SUFFIX}installer-${PT_VERSION}:${REVISION} | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `docker/Dockerfile.rdma-tools` | nvcr.io/nvidia/cuda:${CUDA_MAJOR}.${CUDA_MINOR}.${CUDA_PATCH}-runtime-ubi9 | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `docker/Dockerfile.rocm` | base | 3 |  |  |  |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: base; Unpinned base image: base; No USER directive found (defaults to root) |
| `helpers/interactive-pod/build/Dockerfile` | python:3.12-slim | 1 |  |  |  |  | No USER directive found (defaults to root) |

