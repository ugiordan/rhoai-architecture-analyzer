# guardrails-detectors: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `detectors/Dockerfile.builtIn` | builder | 3 |  |  |  |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal; Unpinned base image: base; Unpinned base image: builder; No USER directive found (defaults to root) |
| `detectors/Dockerfile.hf` | builder | 5 | root |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal; Unpinned base image: base; Unpinned base image: base; Unpinned base image: base; Unpinned base image: builder; Container runs as root user |
| `detectors/Dockerfile.judge` | builder | 3 |  |  |  |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal; Unpinned base image: base; Unpinned base image: builder; No USER directive found (defaults to root) |

