# NeMo-Guardrails: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `.devcontainer/Dockerfile` | mcr.microsoft.com/vscode/devcontainers/python:0-${VARIANT} | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `Dockerfile` | python:3.12-slim | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `Dockerfile.konflux` | registry.access.redhat.com/ubi9/python-312:latest | 8 | 1001 |  | multi-arch |  | Unpinned base image: packages-build; Unpinned base image: packages-build; Unpinned base image: packages-build; Unpinned base image: packages-build; Unpinned base image: packages-build; Unpinned base image: registry.access.redhat.com/ubi9/python-312:latest |
| `Dockerfile.server` | registry.access.redhat.com/ubi9/python-312:latest | 2 | 1001 |  |  |  | Unpinned base image: registry.access.redhat.com/ubi9/python-312:latest; Unpinned base image: registry.access.redhat.com/ubi9/python-312:latest |
| `nemoguardrails/library/factchecking/align_score/Dockerfile` | python:3.10 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `nemoguardrails/library/jailbreak_detection/Dockerfile` | python:3.11-slim | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `qa/Dockerfile.qa` | python:3.10 | 1 |  |  |  |  | No USER directive found (defaults to root) |

