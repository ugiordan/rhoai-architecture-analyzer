# vllm: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `Dockerfile` | vllm-openai-base | 8 |  |  | multi-arch |  | Unpinned base image: base; Unpinned base image: base; Unpinned base image: vllm-base; Unpinned base image: vllm-base; Unpinned base image: vllm-openai-base; Unpinned base image: vllm-openai-base; No USER directive found (defaults to root) |
| `Dockerfile.arm` | cpu-test-arm | 2 |  |  |  |  | Unpinned base image: cpu-test-arm; No USER directive found (defaults to root) |
| `Dockerfile.cpu` | cpu-test-1 | 2 |  |  |  |  | Unpinned base image: cpu-test-1; No USER directive found (defaults to root) |
| `Dockerfile.hpu` | vault.habana.ai/gaudi-docker/1.19.1/ubuntu22.04/habanalabs/pytorch-installer-2.5.1:latest | 1 |  |  |  |  | Unpinned base image: vault.habana.ai/gaudi-docker/1.19.1/ubuntu22.04/habanalabs/pytorch-installer-2.5.1:latest; No USER directive found (defaults to root) |
| `Dockerfile.neuron` | $BASE_IMAGE | 1 |  |  |  |  | Unpinned base image: $BASE_IMAGE; No USER directive found (defaults to root) |
| `Dockerfile.openvino` | ubuntu:22.04 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `Dockerfile.ppc64le` | mambaorg/micromamba | 1 | root |  |  |  | Unpinned base image: mambaorg/micromamba; Container runs as root user |
| `Dockerfile.ppc64le.ubi` | vllm-openai | 9 | 2000 |  |  |  | Unpinned base image: base-builder; Unpinned base image: base-builder; Unpinned base image: base-builder; Unpinned base image: base-builder; Unpinned base image: base-builder; Unpinned base image: base-builder; Unpinned base image: vllm-openai |
| `Dockerfile.rocm` | base | 8 |  |  |  |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: base; Unpinned base image: base; Unpinned base image: fetch_vllm_${REMOTE_VLLM}; Unpinned base image: fetch_vllm; Unpinned base image: scratch; Unpinned base image: base; Unpinned base image: base; No USER directive found (defaults to root) |
| `Dockerfile.rocm.ubi` | vllm-openai | 9 | 2000 |  |  |  | Unpinned base image: base; Unpinned base image: rocm_base; Unpinned base image: rocm_devel; Unpinned base image: rocm_devel; Unpinned base image: rocm_devel; Unpinned base image: rocm_base; Unpinned base image: rocm_base; Unpinned base image: vllm-openai |
| `Dockerfile.rocm_base` | base | 7 |  |  |  |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; No USER directive found (defaults to root) |
| `Dockerfile.s390x.ubi` | vllm-openai | 9 | 2000 |  |  |  | Unpinned base image: base; Unpinned base image: python-install; Unpinned base image: python-install; Unpinned base image: python-install; Unpinned base image: python-install; Unpinned base image: python-install; Unpinned base image: vllm-openai; Unpinned base image: vllm-openai |
| `Dockerfile.tpu` | $BASE_IMAGE | 1 |  |  |  |  | Unpinned base image: $BASE_IMAGE; No USER directive found (defaults to root) |
| `Dockerfile.ubi` | vllm-openai | 9 | 2000 |  |  |  | Unpinned base image: base; Unpinned base image: python-install; Unpinned base image: cuda-base; Unpinned base image: python-cuda-base; Unpinned base image: dev; Unpinned base image: base; Unpinned base image: python-install; Unpinned base image: vllm-openai |
| `Dockerfile.xpu` | vllm-base | 2 |  |  |  |  | Unpinned base image: vllm-base; No USER directive found (defaults to root) |

