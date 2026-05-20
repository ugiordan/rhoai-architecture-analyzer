# pipelines-components: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `Dockerfile` | registry.redhat.io/ubi9/python-311@sha256:d7620b96616955d78425518143affdc9463fb1e71d00aa2b7dc2785c54621a0b | 1 | 1001 |  |  |  |  |
| `Dockerfile.konflux.automl` | ${BASE_IMAGE} | 1 |  |  |  |  | Unpinned base image: ${BASE_IMAGE}; No USER directive found (defaults to root) |
| `Dockerfile.konflux.autorag` | ${BASE_IMAGE} | 1 | default |  |  |  | Unpinned base image: ${BASE_IMAGE} |
| `Dockerfile.konflux.pipelines-components` | registry.redhat.io/ubi9/python-311@sha256:89d4f0c9c39bee97104b7f16bcd15c332c3d4e7cdc80cad81d73ab6ca59055ff | 1 | 1001 |  |  |  |  |
| `components/training/autorag/rag_templates_optimization/Dockerfile` | redhat/ubi10-minimal:10.1@sha256:a74a7a92d3069bfac09c6882087771fc7db59fa9d8e16f14f4e012fe7288554c | 1 |  |  | multi-arch |  | No USER directive found (defaults to root) |
| `pipelines/training/automl/autogluon_tabular_training_pipeline/Containerfile` | ${BASE_IMAGE} | 1 |  |  |  |  | Unpinned base image: ${BASE_IMAGE}; No USER directive found (defaults to root) |
| `pipelines/training/autorag/documents_rag_optimization_pipeline/Containerfile` | ${BASE_IMAGE} | 1 | default |  |  |  | Unpinned base image: ${BASE_IMAGE} |

