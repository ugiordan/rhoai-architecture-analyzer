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
| controller-manager | manager | ? | true | ? | [`config/manager/manager.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/5a8e5dbe95e606d5d964968f933b3037399805c8/config/manager/manager.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/overlays/odh/manager_patch.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/5a8e5dbe95e606d5d964968f933b3037399805c8/config/overlays/odh/manager_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/overlays/openshift/manager_patch.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/5a8e5dbe95e606d5d964968f933b3037399805c8/config/overlays/openshift/manager_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/overlays/rhoai/manager_patch.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/5a8e5dbe95e606d5d964968f933b3037399805c8/config/overlays/rhoai/manager_patch.yaml) |
| mlflow-operator-controller-manager | manager | ? | ? | ? | [`config/overlays/kind/manager-patch.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/5a8e5dbe95e606d5d964968f933b3037399805c8/config/overlays/kind/manager-patch.yaml) |
| postgres-deployment | postgres | ? | ? | ? | [`config/postgres/base/deployment.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/5a8e5dbe95e606d5d964968f933b3037399805c8/config/postgres/base/deployment.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `Dockerfile` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 2 | 1001 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |
| `mlflow-tests/images/Dockerfile.konflux` | registry.access.redhat.com/ubi9/ubi-minimal:9.7 | 1 | 1001 |  |  |  |  |

