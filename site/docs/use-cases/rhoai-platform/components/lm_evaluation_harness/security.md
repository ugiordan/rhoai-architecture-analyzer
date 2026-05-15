# lm-evaluation-harness: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `Dockerfile.konflux.lmes-job` | builder | 5 | 65532:65532 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/python-311:latest; Unpinned base image: builder; Unpinned base image: builder; Unpinned base image: builder; Unpinned base image: builder |
| `Dockerfile.lmes-job` | builder | 5 | 65532:65532 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/python-311:latest; Unpinned base image: builder; Unpinned base image: builder; Unpinned base image: builder; Unpinned base image: builder |

