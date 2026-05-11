# data-science-pipelines-operator: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

### Secrets Referenced

| Name | Type | Referenced By |
|------|------|---------------|
| ds-pipeline-db-test | Opaque | deployment/mariadb |
| mariadb-certs | Opaque | deployment/mariadb |
| minio | Opaque | deployment/minio |
| minio-certs | Opaque | deployment/minio |

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| data-science-pipelines-operator-controller-manager | manager | ? | ? | ? | [`kustomize:config/overlays/odh`](https://github.com/opendatahub-io/data-science-pipelines-operator/blob/6d6d1b7e37c804aa93e8d73c17ab6a350af230c5/kustomize:config/overlays/odh) |
| mariadb | mariadb | ? | ? | ? | [`.github/resources/mariadb/deployment.yaml`](https://github.com/opendatahub-io/data-science-pipelines-operator/blob/6d6d1b7e37c804aa93e8d73c17ab6a350af230c5/.github/resources/mariadb/deployment.yaml) |
| minio | minio | ? | ? | ? | [`.github/resources/minio/deployment.yaml`](https://github.com/opendatahub-io/data-science-pipelines-operator/blob/6d6d1b7e37c804aa93e8d73c17ab6a350af230c5/.github/resources/minio/deployment.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `.github/build/Dockerfile` | ${CI_BASE} | 2 | root |  |  |  | Unpinned base image: ${CI_BASE}; Unpinned base image: ${CI_BASE}; Container runs as root user |
| `.github/scripts/python_package_upload/Dockerfile` | docker.io/python:3.9 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `Dockerfile` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 2 | ${USER}:${USER} |  | multi-arch | yes | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |
| `Dockerfile.konflux` | registry.access.redhat.com/ubi9/ubi-minimal@sha256:fe9e574f04371b333ed4e21d30d984f6b7fcd1046e579f5ddab4816c0c8e231d | 2 | ${USER}:${USER} |  | multi-arch |  |  |

