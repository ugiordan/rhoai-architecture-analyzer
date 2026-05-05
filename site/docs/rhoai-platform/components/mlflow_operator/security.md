# mlflow-operator: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

### Secrets Referenced

| Name | Type | Referenced By |
|------|------|---------------|
| controller-manager-metrics-tls | Opaque | deployment/controller-manager |
| postgres-secret | Opaque | deployment/postgres-deployment |

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| mlflow-operator-controller-manager | manager | ? | true | ? | [`kustomize:config/overlays/odh`](https://github.com/opendatahub-io/mlflow-operator/blob/682055b5deae5d1cc92c0a24270aee8400704084/kustomize:config/overlays/odh) |
| postgres-deployment | postgres | ? | ? | ? | [`config/postgres/base/deployment.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/682055b5deae5d1cc92c0a24270aee8400704084/config/postgres/base/deployment.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `Dockerfile` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 2 | 1001 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |
| `Dockerfile.konflux` | registry.access.redhat.com/ubi9/ubi-minimal@sha256:8d0a8fb39ec907e8ca62cdd24b62a63ca49a30fe465798a360741fde58437a23 | 2 | 1001 |  | multi-arch |  |  |
| `mlflow-tests/images/Dockerfile.konflux` | registry.access.redhat.com/ubi9/ubi-minimal:9.7 | 1 | 1001 |  |  |  |  |

