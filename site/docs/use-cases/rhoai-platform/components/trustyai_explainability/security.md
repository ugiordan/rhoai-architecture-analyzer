# trustyai-explainability: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `Dockerfile` | registry.access.redhat.com/ubi8/openjdk-17-runtime:latest | 2 | root |  |  |  | Unpinned base image: registry.access.redhat.com/ubi8/openjdk-17:latest; Unpinned base image: registry.access.redhat.com/ubi8/openjdk-17-runtime:latest; Container runs as root user |
| `Dockerfile.konflux` | registry.access.redhat.com/ubi9/openjdk-17-runtime@sha256:07e1fcaaeb8d4210dccd781a47e123df5d3c917d868a4ff325739118403bbc7c | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `explainability-service/src/main/docker/Dockerfile.jvm` | registry.access.redhat.com/ubi8/openjdk-17:latest | 1 | 185 |  |  |  | Unpinned base image: registry.access.redhat.com/ubi8/openjdk-17:latest |
| `explainability-service/src/main/docker/Dockerfile.legacy-jar` | registry.access.redhat.com/ubi8/openjdk-17:latest | 1 | 185 |  |  |  | Unpinned base image: registry.access.redhat.com/ubi8/openjdk-17:latest |
| `explainability-service/src/main/docker/Dockerfile.native` | registry.access.redhat.com/ubi8/ubi-minimal:8.6 | 1 | 1001 |  |  |  |  |
| `explainability-service/src/main/docker/Dockerfile.native-micro` | quay.io/quarkus/quarkus-micro-image:2.0 | 1 | 1001 |  |  |  |  |

