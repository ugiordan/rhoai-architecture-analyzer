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
| controller-manager | manager | ? | true | ? | [`config/manager/manager.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/46d31981c7867fd944d4ab8a9e0d8434648d379e/config/manager/manager.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/overlays/odh/manager_patch.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/46d31981c7867fd944d4ab8a9e0d8434648d379e/config/overlays/odh/manager_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/overlays/openshift/manager_patch.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/46d31981c7867fd944d4ab8a9e0d8434648d379e/config/overlays/openshift/manager_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/overlays/rhoai/manager_patch.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/46d31981c7867fd944d4ab8a9e0d8434648d379e/config/overlays/rhoai/manager_patch.yaml) |
| mlflow-operator-controller-manager | manager | ? | ? | ? | [`config/overlays/kind/manager-patch.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/46d31981c7867fd944d4ab8a9e0d8434648d379e/config/overlays/kind/manager-patch.yaml) |
| postgres-deployment | postgres | ? | ? | ? | [`config/postgres/base/deployment.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/46d31981c7867fd944d4ab8a9e0d8434648d379e/config/postgres/base/deployment.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `Dockerfile` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 2 | 1001 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |
| `Dockerfile.konflux` | registry.access.redhat.com/ubi9/ubi-minimal@sha256:7d4e47500f28ac3a2bff06c25eff9127ff21048538ae03ce240d57cf756acd00 | 2 | 1001 |  | multi-arch |  |  |
| `mlflow-tests/images/Dockerfile.konflux` | registry.access.redhat.com/ubi9/ubi-minimal:9.7 | 1 | 1001 |  |  |  |  |

