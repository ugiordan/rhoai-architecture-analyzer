# vllm-cpu: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `Dockerfile.cpu.ubi` | python-install | 3 | 2000 |  |  |  | Unpinned base image: base; Unpinned base image: python-install |
| `Dockerfile.hpu.ubi` | python-install | 3 | 2000 |  |  |  | Unpinned base image: base; Unpinned base image: python-install |
| `Dockerfile.konflux.cpu` | registry.access.redhat.com/ubi9/ubi-minimal:${BASE_UBI_IMAGE_TAG} | 2 | 2000 |  |  |  |  |
| `Dockerfile.ppc64le.ubi` | vllm-openai | 12 | 2000 |  |  |  | Unpinned base image: centos-deps-builder; Unpinned base image: base-builder; Unpinned base image: base-builder; Unpinned base image: base-builder; Unpinned base image: base-builder; Unpinned base image: base-builder; Unpinned base image: base-builder; Unpinned base image: base-builder; Unpinned base image: vllm-openai |
| `Dockerfile.rocm.ubi` | rocm_base | 3 | 2000 |  |  |  | Unpinned base image: base; Unpinned base image: rocm_base |
| `Dockerfile.s390x.ubi` | vllm-openai | 14 | 2000 |  |  |  | Unpinned base image: base; Unpinned base image: python-install; Unpinned base image: python-install; Unpinned base image: python-install; Unpinned base image: python-install; Unpinned base image: python-install; Unpinned base image: python-install; Unpinned base image: python-install; Unpinned base image: python-install; Unpinned base image: python-install; Unpinned base image: python-install; Unpinned base image: vllm-openai; Unpinned base image: vllm-openai |
| `Dockerfile.tpu.ubi` | registry.access.redhat.com/ubi9/ubi-minimal:${BASE_UBI_IMAGE_TAG} | 1 | 2000 |  |  |  |  |
| `Dockerfile.ubi` | python-install | 4 | 2000 |  |  |  | Unpinned base image: base; Unpinned base image: cuda-base; Unpinned base image: python-install |
| `docker/Dockerfile` | vllm-openai-base | 10 |  |  | multi-arch |  | Unpinned base image: ${BUILD_BASE_IMAGE}; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; Unpinned base image: ${FINAL_BASE_IMAGE}; Unpinned base image: vllm-base; Unpinned base image: vllm-base; Unpinned base image: vllm-openai-base; Unpinned base image: vllm-openai-base; No USER directive found (defaults to root) |
| `docker/Dockerfile.cpu` | vllm-openai | 10 |  |  | multi-arch |  | Unpinned base image: base-common; Unpinned base image: base-common; Unpinned base image: base-${TARGETARCH}; Unpinned base image: base; Unpinned base image: base; Unpinned base image: vllm-build; Unpinned base image: vllm-test-deps; Unpinned base image: base; Unpinned base image: vllm-openai; No USER directive found (defaults to root) |
| `docker/Dockerfile.nightly_torch` | vllm-base | 4 |  |  | multi-arch |  | Unpinned base image: base; Unpinned base image: vllm-base; No USER directive found (defaults to root) |
| `docker/Dockerfile.ppc64le` | registry.access.redhat.com/ubi9/ubi-minimal:${BASE_UBI_IMAGE_TAG} | 11 |  |  |  |  | Unpinned base image: centos-deps-builder; Unpinned base image: base-builder; Unpinned base image: base-builder; Unpinned base image: base-builder; Unpinned base image: base-builder; Unpinned base image: base-builder; Unpinned base image: base-builder; Unpinned base image: base-builder; No USER directive found (defaults to root) |
| `docker/Dockerfile.rocm` | final | 13 |  |  |  |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: base; Unpinned base image: base; Unpinned base image: fetch_vllm_${REMOTE_VLLM}; Unpinned base image: fetch_vllm; Unpinned base image: scratch; Unpinned base image: base; Unpinned base image: base; Unpinned base image: fetch_vllm; Unpinned base image: scratch; Unpinned base image: base; Unpinned base image: base; Unpinned base image: final; No USER directive found (defaults to root) |
| `docker/Dockerfile.rocm_base` | base | 10 |  |  |  |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; No USER directive found (defaults to root) |
| `docker/Dockerfile.s390x` | python-install | 11 | 2000 |  |  |  | Unpinned base image: base; Unpinned base image: python-install; Unpinned base image: python-install; Unpinned base image: python-install; Unpinned base image: python-install; Unpinned base image: python-install; Unpinned base image: python-install; Unpinned base image: python-install; Unpinned base image: python-install; Unpinned base image: python-install |
| `docker/Dockerfile.tpu` | $BASE_IMAGE | 1 |  |  |  |  | Unpinned base image: $BASE_IMAGE; No USER directive found (defaults to root) |
| `docker/Dockerfile.xpu` | vllm-base | 2 |  |  |  |  | Unpinned base image: vllm-base; No USER directive found (defaults to root) |

