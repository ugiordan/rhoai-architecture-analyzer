# distributed-workloads: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `benchmarks/osu-benchmarks/Dockerfile` | quay.io/opendatahub/odh-midstream-python-base-3-12:11b0d2c14a1d8de8171d428172aef0cde54ec7a7 | 2 | 1001 |  |  |  |  |
| `benchmarks/osu-benchmarks/Dockerfile.cuda` | ${TRAINING_BASE_IMAGE} | 1 | 1001 |  |  |  | Unpinned base image: ${TRAINING_BASE_IMAGE} |
| `images/dataset/alpaca/Dockerfile` | registry.access.redhat.com/ubi9:latest | 2 |  |  |  |  | Unpinned base image: registry.access.redhat.com/ubi9/python-311:latest; Unpinned base image: registry.access.redhat.com/ubi9:latest; No USER directive found (defaults to root) |
| `images/model/bloom560m/Dockerfile` | registry.access.redhat.com/ubi9:9.4 | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `images/runtime/ray/cpu/2.52.1-py311-cpu/Dockerfile` | registry.access.redhat.com/ubi9/python-${PYTHON_VERSION}:${IMAGE_TAG} | 1 | 1001 |  |  |  |  |
| `images/runtime/ray/cpu/2.52.1-py312-cpu/Dockerfile` | registry.access.redhat.com/ubi9/python-${PYTHON_VERSION}:${IMAGE_TAG} | 1 | 1001 |  |  |  |  |
| `images/runtime/ray/cuda/2.52.1-py311-cu121/Dockerfile` | registry.access.redhat.com/ubi9/python-${PYTHON_VERSION}:${IMAGE_TAG} | 1 | 1001 |  |  |  |  |
| `images/runtime/ray/cuda/2.54.1-py312-cu128/Dockerfile` | registry.access.redhat.com/ubi9/python-${PYTHON_VERSION}:${IMAGE_TAG} | 1 | 1001 |  | multi-arch |  |  |
| `images/runtime/ray/rocm/2.52.1-py311-rocm61/Dockerfile` | registry.access.redhat.com/ubi9/python-${PYTHON_VERSION}:${IMAGE_TAG} | 1 | 1001 |  |  |  |  |
| `images/runtime/ray/rocm/2.54.1-py312-rocm64/Dockerfile` | registry.access.redhat.com/ubi9/python-${PYTHON_VERSION}:${IMAGE_TAG} | 1 | 1001 |  |  |  |  |
| `images/runtime/training/py311-cuda121-torch241/Dockerfile` | registry.access.redhat.com/ubi9/python-${PYTHON_VERSION}:${IMAGE_TAG} | 1 | 1001 |  |  |  |  |
| `images/runtime/training/py311-cuda121-torch241/Dockerfile.konflux` | registry.access.redhat.com/ubi9/python-311:latest | 1 | 1001 |  |  |  | Unpinned base image: registry.access.redhat.com/ubi9/python-311:latest |
| `images/runtime/training/py311-cuda124-torch251/Dockerfile` | registry.access.redhat.com/ubi9/python-${PYTHON_VERSION}:${IMAGE_TAG} | 1 | 1001 |  |  |  |  |
| `images/runtime/training/py311-cuda124-torch251/Dockerfile.konflux` | registry.access.redhat.com/ubi9/python-311:latest | 1 | 1001 |  |  |  | Unpinned base image: registry.access.redhat.com/ubi9/python-311:latest |
| `images/runtime/training/py311-rocm62-torch241/Dockerfile` | registry.access.redhat.com/ubi9/python-${PYTHON_VERSION}:${IMAGE_TAG} | 1 | 1001 |  |  |  |  |
| `images/runtime/training/py311-rocm62-torch241/Dockerfile.konflux` | registry.access.redhat.com/ubi9/python-311:latest | 1 | 1001 |  |  |  | Unpinned base image: registry.access.redhat.com/ubi9/python-311:latest |
| `images/runtime/training/py311-rocm62-torch251/Dockerfile` | registry.access.redhat.com/ubi9/python-${PYTHON_VERSION}:${IMAGE_TAG} | 1 | 1001 |  |  |  |  |
| `images/runtime/training/py311-rocm62-torch251/Dockerfile.konflux` | registry.access.redhat.com/ubi9/python-311:latest | 1 | 1001 |  |  |  | Unpinned base image: registry.access.redhat.com/ubi9/python-311:latest |
| `images/runtime/training/py312-cuda128-torch280/Dockerfile` | registry.access.redhat.com/ubi9/python-${PYTHON_VERSION}:${IMAGE_TAG} | 1 | 1001 |  |  |  |  |
| `images/runtime/training/py312-cuda128-torch280/Dockerfile.konflux` | registry.access.redhat.com/ubi9/python-312:latest | 1 | 1001 |  |  |  | Unpinned base image: registry.access.redhat.com/ubi9/python-312:latest |
| `images/runtime/training/py312-cuda128-torch290/Dockerfile` | registry.access.redhat.com/ubi9/python-${PYTHON_VERSION}:${IMAGE_TAG} | 1 | 1001 |  |  |  |  |
| `images/runtime/training/py312-cuda128-torch290/Dockerfile.konflux` | registry.access.redhat.com/ubi9/python-312@sha256:a0a5885769d5a8c5123d3b15d5135b254541d4da8e7bc445d95e1c90595de470 | 1 | 1001 |  |  |  |  |
| `images/runtime/training/py312-cuda130-torch210-openmpi41/Dockerfile` | quay.io/opendatahub/odh-midstream-cuda-base-13-0:odh-midstream-cuda-base-13-0-on-push-p7nn6-build-images | 1 | 1001 |  |  |  |  |
| `images/runtime/training/py312-cuda130-torch210-openmpi41/Dockerfile.konflux` | quay.io/aipcc/base-images/cuda-13.0-el9.6@sha256:3de8a37c9172aea6a15fe12aeeb9fd6be09a5a5ca419ec2e9fc2e16c3f0b138d | 1 | 1001 |  |  |  |  |
| `images/runtime/training/py312-rocm64-torch280/Dockerfile` | registry.access.redhat.com/ubi9/python-${PYTHON_VERSION}:${IMAGE_TAG} | 1 | 1001 |  |  |  |  |
| `images/runtime/training/py312-rocm64-torch280/Dockerfile.konflux` | registry.access.redhat.com/ubi9/python-312:latest | 1 | 1001 |  |  |  | Unpinned base image: registry.access.redhat.com/ubi9/python-312:latest |
| `images/runtime/training/py312-rocm64-torch29-openmpi41/Dockerfile.konflux` | quay.io/aipcc/base-images/rocm-6.4-el9.6@sha256:444b1345c0bcf09ae1cb9a4f618c349522930021c8fc3e105276768beaf322cd | 1 | 1001 |  |  |  |  |
| `images/runtime/training/py312-rocm64-torch290/Dockerfile` | registry.access.redhat.com/ubi9/python-${PYTHON_VERSION}:${IMAGE_TAG} | 1 | 1001 |  |  |  |  |
| `images/runtime/training/py312-rocm64-torch290/Dockerfile.konflux` | registry.access.redhat.com/ubi9/python-312@sha256:a0a5885769d5a8c5123d3b15d5135b254541d4da8e7bc445d95e1c90595de470 | 1 | 1001 |  |  |  |  |
| `images/universal/training/th06-cpu-torch210-py312/Dockerfile` | ${BASE_IMAGE} | 2 | 1001 |  |  |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: ${BASE_IMAGE} |
| `images/universal/training/th06-cpu-torch210-py312/Dockerfile.konflux.cpu` | ${BASE_IMAGE} | 2 | 1001 |  |  |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: ${BASE_IMAGE} |
| `images/universal/training/th06-cuda130-torch210-py312/Dockerfile` | ${BASE_IMAGE} | 2 | 1001 |  |  |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: ${BASE_IMAGE} |
| `images/universal/training/th06-cuda130-torch210-py312/Dockerfile.konflux.cuda` | ${BASE_IMAGE} | 2 | 1001 |  |  |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: ${BASE_IMAGE} |
| `images/universal/training/th06-rocm64-torch291-py312/Dockerfile` | ${BASE_IMAGE} | 2 | 1001 |  |  |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: ${BASE_IMAGE} |
| `images/universal/training/th06-rocm64-torch291-py312/Dockerfile.konflux.rocm` | ${BASE_IMAGE} | 2 | 1001 |  |  |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: ${BASE_IMAGE} |
| `images/util/mc-cli/Dockerfile` | registry.access.redhat.com/ubi9:latest | 1 | 1001 |  |  |  | Unpinned base image: registry.access.redhat.com/ubi9:latest |

