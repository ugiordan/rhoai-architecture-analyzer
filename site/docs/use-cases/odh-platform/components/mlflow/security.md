# mlflow: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `Dockerfile.konflux` | registry.access.redhat.com/ubi9/ubi-minimal:9.7 | 3 | 1001 |  |  |  |  |
| `dev/Dockerfile.protos` | python:3.10-slim | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `dev/Dockerfile.protos.dockerignore` |  | 0 |  |  |  |  | No USER directive found (defaults to root) |
| `docker/Dockerfile` | python:3.10-slim-bullseye@sha256:f1fb49e4d5501ac93d0ca519fb7ee6250842245aba8612926a46a0832a1ed089 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `docker/Dockerfile.full` | python:3.10-slim-bullseye@sha256:f1fb49e4d5501ac93d0ca519fb7ee6250842245aba8612926a46a0832a1ed089 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `docker/Dockerfile.full.dev` | python:3.10-slim-bullseye@sha256:f1fb49e4d5501ac93d0ca519fb7ee6250842245aba8612926a46a0832a1ed089 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `mlflow/R/mlflow/Dockerfile.dev` | rocker/r-ver:4.4.2 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `mlflow/R/mlflow/Dockerfile.r-devel` | rocker/r-ver:devel | 1 |  |  |  |  | No USER directive found (defaults to root) |

