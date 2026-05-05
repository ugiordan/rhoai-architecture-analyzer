# model-registry: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

### Secrets Referenced

| Name | Type | Referenced By |
|------|------|---------------|
| minio-secret | Opaque | deployment/minio |
| model-catalog-hf-api-key | Opaque | deployment/model-catalog-server |
| model-catalog-postgres | Opaque | deployment/model-catalog-server |

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| controller-manager | manager | ? | ? | ? | [`manifests/kustomize/options/controller/manager/manager.yaml`](https://github.com/kubeflow/model-registry/blob/bbd3a37dfa4adfa6239250a7c0cbf9b17fe7904a/manifests/kustomize/options/controller/manager/manager.yaml) |
| minio | minio | ? | ? | ? | [`scripts/manifests/minio/deployment.yaml`](https://github.com/kubeflow/model-registry/blob/bbd3a37dfa4adfa6239250a7c0cbf9b17fe7904a/scripts/manifests/minio/deployment.yaml) |
| model-catalog-server | catalog | ? | ? | ? | [`manifests/kustomize/options/catalog/base/deployment.yaml`](https://github.com/kubeflow/model-registry/blob/bbd3a37dfa4adfa6239250a7c0cbf9b17fe7904a/manifests/kustomize/options/catalog/base/deployment.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `Dockerfile` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 3 | 65532:65532 |  | multi-arch |  | Unpinned base image: common; Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |
| `Dockerfile.konflux` | registry.access.redhat.com/ubi9/ubi-minimal@sha256:8d0a8fb39ec907e8ca62cdd24b62a63ca49a30fe465798a360741fde58437a23 | 2 | 1001 |  |  |  |  |
| `Dockerfile.odh` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 2 | 65532:65532 |  |  |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |
| `Dockerfile.testops` | registry.access.redhat.com/ubi9/python-312 | 1 | odh |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/python-312 |
| `clients/ui/Dockerfile` | ${DISTROLESS_BASE_IMAGE} | 3 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${NODE_BASE_IMAGE}; Unpinned base image: ${GOLANG_BASE_IMAGE}; Unpinned base image: ${DISTROLESS_BASE_IMAGE} |
| `clients/ui/Dockerfile.standalone` | release | 4 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${NODE_BASE_IMAGE}; Unpinned base image: ${GOLANG_BASE_IMAGE}; Unpinned base image: ${DISTROLESS_BASE_IMAGE}; Unpinned base image: release |
| `cmd/controller/Dockerfile.controller` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 2 | 65532:65532 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |
| `cmd/csi/Dockerfile.csi` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 3 | 65532:65532 |  | multi-arch |  | Unpinned base image: common; Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |
| `jobs/async-upload/Dockerfile` | registry.access.redhat.com/ubi9/python-312-minimal | 2 | 1000 |  |  |  | Unpinned base image: registry.access.redhat.com/ubi9/python-312-minimal |
| `jobs/async-upload/Dockerfile.konflux` | registry.access.redhat.com/ubi9/python-312-minimal@sha256:13c75e682cd460da9ba76cc7ebb6722bd1cb121c1c8e19e54df5c7d820a0545e | 1 | 1000 |  | multi-arch |  |  |

