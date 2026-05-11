# llm-d-kv-cache: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| 0 | cmd | ? | ? | ? | [`deploy/common/statefulset.yaml`](https://github.com/llm-d/llm-d-kv-cache/blob/d54e631afebc240807275fd702a2277448fe4db8/deploy/common/statefulset.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `Dockerfile` | registry.access.redhat.com/ubi9/ubi:latest | 2 | 65532:65532 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi:latest |
| `kv_connectors/llmd_fs_backend/Dockerfile.dev` | ${VLLM_IMAGE} | 1 |  |  |  |  | Unpinned base image: ${VLLM_IMAGE}; No USER directive found (defaults to root) |
| `kv_connectors/llmd_fs_backend/Dockerfile.wheel` | nvcr.io/nvidia/cuda:${CUDA_VERSION}-devel-ubuntu22.04 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `kv_connectors/pvc_evictor/Dockerfile` | python:3.12-slim | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `services/uds_tokenizer/Dockerfile` | python:3.12-slim | 2 | 65532:65532 |  | multi-arch |  |  |
| `services/uds_tokenizer/Dockerfile.konflux` | runtime | 4 | 65532:65532 |  | multi-arch |  | Unpinned base image: runtime; Unpinned base image: runtime |

