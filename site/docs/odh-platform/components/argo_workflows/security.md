# argo-workflows: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `Dockerfile` | gcr.io/distroless/static | 10 | 8737 |  |  |  | Unpinned base image: builder; Unpinned base image: builder; Unpinned base image: builder; Unpinned base image: gcr.io/distroless/static; Unpinned base image: argoexec-base; Unpinned base image: argoexec-base; Unpinned base image: gcr.io/distroless/static; Unpinned base image: gcr.io/distroless/static |
| `Dockerfile.windows` | argoexec-base | 4 | Administrator |  |  |  | Unpinned base image: builder; Unpinned base image: argoexec-base |
| `argo-argoexec/Dockerfile.ODH` | registry.redhat.io/ubi9/ubi-minimal:9.5 | 2 | 2000 |  |  |  |  |
| `argo-workflowcontroller/Dockerfile.ODH` | registry.redhat.io/ubi9/ubi-minimal:9.5 | 2 | 8737 |  |  |  |  |

