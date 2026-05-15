# openvino_model_server: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `Dockerfile.konflux` | $RELEASE_BASE_IMAGE | 5 | ovms |  |  |  | Unpinned base image: $BASE_IMAGE; Unpinned base image: base_build; Unpinned base image: $BUILD_IMAGE; Unpinned base image: $BUILD_IMAGE; Unpinned base image: $RELEASE_BASE_IMAGE |
| `Dockerfile.redhat` | $RELEASE_BASE_IMAGE | 5 | ovms |  |  |  | Unpinned base image: $BASE_IMAGE; Unpinned base image: base_build; Unpinned base image: $BUILD_IMAGE; Unpinned base image: $BUILD_IMAGE; Unpinned base image: $RELEASE_BASE_IMAGE |
| `Dockerfile.ubuntu` | $BASE_IMAGE | 5 | ovms |  |  |  | Unpinned base image: $BASE_IMAGE; Unpinned base image: base_build; Unpinned base image: $BUILD_IMAGE; Unpinned base image: $BUILD_IMAGE; Unpinned base image: $BASE_IMAGE |
| `ci/Dockerfile` | $BUILD_IMAGE | 1 |  |  |  |  | Unpinned base image: $BUILD_IMAGE; No USER directive found (defaults to root) |
| `ci/Dockerfile.coverity` | $BUILD_IMAGE | 1 |  |  |  |  | Unpinned base image: $BUILD_IMAGE; No USER directive found (defaults to root) |
| `client/go/kserve-api/Dockerfile` | golang:1.24.0 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `demos/benchmark/python/Dockerfile` | haproxy:2.8.2 | 1 | root |  |  |  | Container runs as root user |
| `demos/bert_question_answering/python/Dockerfile` | ubuntu:22.04 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `demos/c_api_minimal_app/capi_files/Dockerfile.redhat` | registry.access.redhat.com/ubi9/ubi-minimal:9.7 | 2 | ovms |  |  |  | Unpinned base image: $BASE_IMAGE |
| `demos/c_api_minimal_app/capi_files/Dockerfile.ubuntu` | $BASE_IMAGE | 2 | ovms |  |  |  | Unpinned base image: $BASE_IMAGE; Unpinned base image: $BASE_IMAGE |
| `demos/common/stream_client/Dockerfile` | ubuntu:22.04 | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `demos/image_classification/go/Dockerfile` | golang:1.24.0 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `demos/python_demos/Dockerfile.redhat` | ${IMAGE_NAME} | 1 | ovms |  |  |  | Unpinned base image: ${IMAGE_NAME} |
| `demos/python_demos/Dockerfile.ubuntu` | ${IMAGE_NAME} | 1 | ovms |  |  |  | Unpinned base image: ${IMAGE_NAME} |
| `extras/nginx-mtls-auth/Dockerfile.redhat` | $BASE_IMAGE | 1 | ovms |  |  |  | Unpinned base image: $BASE_IMAGE |
| `extras/nginx-mtls-auth/Dockerfile.ubuntu` | $BASE_IMAGE | 1 | ovms |  |  |  | Unpinned base image: $BASE_IMAGE |
| `src/custom_nodes/Dockerfile.redhat` | registry.access.redhat.com/ubi8/ubi:8.10 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `src/custom_nodes/Dockerfile.ubuntu` | $BASE_IMAGE | 1 |  |  |  |  | Unpinned base image: $BASE_IMAGE; No USER directive found (defaults to root) |
| `src/example/SampleCpuExtension/Dockerfile.redhat` | $BASE_IMAGE | 1 |  |  |  |  | Unpinned base image: $BASE_IMAGE; No USER directive found (defaults to root) |
| `src/example/SampleCpuExtension/Dockerfile.ubuntu` | $BASE_IMAGE | 1 |  |  |  |  | Unpinned base image: $BASE_IMAGE; No USER directive found (defaults to root) |
| `tools/deps/redhat/Dockerfile` | registry.access.redhat.com/ubi8/ubi:8.10 | 2 |  |  |  |  | No USER directive found (defaults to root) |

