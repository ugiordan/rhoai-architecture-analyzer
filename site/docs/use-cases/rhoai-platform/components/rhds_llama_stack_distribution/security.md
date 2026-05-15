# rhds-llama-stack-distribution: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `Dockerfile.konflux` | ${BASE_IMAGE} | 1 | ${CNB_USER_ID}:${CNB_GROUP_ID} |  |  |  | Unpinned base image: ${BASE_IMAGE} |
| `distribution/Containerfile` | registry.access.redhat.com/ubi9/python-312@sha256:95ec8d3ee9f875da011639213fd254256c29bc58861ac0b11f290a291fa04435 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `distribution/Containerfile.in` | registry.access.redhat.com/ubi9/python-312@sha256:95ec8d3ee9f875da011639213fd254256c29bc58861ac0b11f290a291fa04435 | 1 |  |  |  |  | No USER directive found (defaults to root) |

