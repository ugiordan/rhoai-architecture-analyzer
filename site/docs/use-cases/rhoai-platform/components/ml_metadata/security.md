# ml-metadata: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `ml_metadata/tools/dev_debug/Dockerfile` | ubuntu:20.04 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `ml_metadata/tools/docker_build/Dockerfile.manylinux2010` | gcr.io/tfx-oss-public/manylinux2014-bazel:bazel-5.3.0 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `ml_metadata/tools/docker_server/Dockerfile` | ubuntu:20.04 | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `ml_metadata/tools/docker_server/Dockerfile.fedora` | fedora:38 | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `ml_metadata/tools/docker_server/Dockerfile.konflux` | registry.access.redhat.com/ubi9/ubi-minimal@sha256:24650313873554b6ba16c1a1b6b9f9142604f6ab735113e1695faf2dd07fdede | 2 | 65534:65534 |  |  |  |  |
| `ml_metadata/tools/docker_server/Dockerfile.redhat` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 2 | 65534:65534 |  |  |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi:latest; Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |

