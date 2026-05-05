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
| mlflow-operator-controller-manager | manager | ? | true | ? | [`kustomize:config/overlays/odh`](https://github.com/opendatahub-io/mlflow-operator/blob/c47bea907957a9eeb35ee4ed3b0403e855d096cc/kustomize:config/overlays/odh) |
| postgres-deployment | postgres | ? | ? | ? | [`config/postgres/base/deployment.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/c47bea907957a9eeb35ee4ed3b0403e855d096cc/config/postgres/base/deployment.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `Dockerfile` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 2 | 1001 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |
| `mlflow-tests/images/Dockerfile.konflux` | registry.access.redhat.com/ubi9/ubi-minimal:9.7 | 1 | 1001 |  |  |  |  |

