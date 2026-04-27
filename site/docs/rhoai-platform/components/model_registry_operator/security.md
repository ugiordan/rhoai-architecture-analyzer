# model-registry-operator: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

### Secrets Referenced

| Name | Type | Referenced By |
|------|------|---------------|
| controller-manager-metrics-service | Opaque | deployment/controller-manager |
| webhook-server-cert | Opaque | deployment/controller-manager |

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| controller-manager | manager | ? | ? | ? | [`config/default/manager_auth_proxy_patch.yaml`](https://github.com/opendatahub-io/model-registry-operator/blob/c21ecc2f94a0910690e7e2f9aa72250639004022/config/default/manager_auth_proxy_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/default/manager_config_patch.yaml`](https://github.com/opendatahub-io/model-registry-operator/blob/c21ecc2f94a0910690e7e2f9aa72250639004022/config/default/manager_config_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/default/manager_webhook_patch.yaml`](https://github.com/opendatahub-io/model-registry-operator/blob/c21ecc2f94a0910690e7e2f9aa72250639004022/config/default/manager_webhook_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/manager/manager.yaml`](https://github.com/opendatahub-io/model-registry-operator/blob/c21ecc2f94a0910690e7e2f9aa72250639004022/config/manager/manager.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/overlays/odh/patches/manager_auth_proxy_patch.yaml`](https://github.com/opendatahub-io/model-registry-operator/blob/c21ecc2f94a0910690e7e2f9aa72250639004022/config/overlays/odh/patches/manager_auth_proxy_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/overlays/odh/patches/manager_istio_config_patch.yaml`](https://github.com/opendatahub-io/model-registry-operator/blob/c21ecc2f94a0910690e7e2f9aa72250639004022/config/overlays/odh/patches/manager_istio_config_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/overlays/odh/patches/manager_migration_env_patch.yaml`](https://github.com/opendatahub-io/model-registry-operator/blob/c21ecc2f94a0910690e7e2f9aa72250639004022/config/overlays/odh/patches/manager_migration_env_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/overlays/odh/patches/manager_webhook_patch.yaml`](https://github.com/opendatahub-io/model-registry-operator/blob/c21ecc2f94a0910690e7e2f9aa72250639004022/config/overlays/odh/patches/manager_webhook_patch.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `Dockerfile` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 2 | 65532:65532 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |
| `Dockerfile.konflux` | registry.access.redhat.com/ubi9/ubi-minimal@sha256:7d4e47500f28ac3a2bff06c25eff9127ff21048538ae03ce240d57cf756acd00 | 2 | 65532:65532 |  |  |  |  |

