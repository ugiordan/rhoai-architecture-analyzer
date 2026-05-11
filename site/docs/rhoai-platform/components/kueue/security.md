# kueue: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

### Secrets Referenced

| Name | Type | Referenced By |
|------|------|---------------|
| webhook-server-cert | Opaque | deployment/controller-manager |

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| controller-manager | manager | ? | ? | ? | [`config/alpha-enabled/manager_config_patch.yaml`](https://github.com/red-hat-data-services/kueue/blob/b20fe017c16e8a59cd2ab44e3c37b24cea1f0a21/config/alpha-enabled/manager_config_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/components/manager/manager.yaml`](https://github.com/red-hat-data-services/kueue/blob/b20fe017c16e8a59cd2ab44e3c37b24cea1f0a21/config/components/manager/manager.yaml) |
| controller-manager | kube-rbac-proxy | ? | ? | ? | [`config/default/manager_auth_proxy_patch.yaml`](https://github.com/red-hat-data-services/kueue/blob/b20fe017c16e8a59cd2ab44e3c37b24cea1f0a21/config/default/manager_auth_proxy_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/default/manager_config_patch.yaml`](https://github.com/red-hat-data-services/kueue/blob/b20fe017c16e8a59cd2ab44e3c37b24cea1f0a21/config/default/manager_config_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/default/manager_metrics_patch.yaml`](https://github.com/red-hat-data-services/kueue/blob/b20fe017c16e8a59cd2ab44e3c37b24cea1f0a21/config/default/manager_metrics_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/default/manager_visibility_patch.yaml`](https://github.com/red-hat-data-services/kueue/blob/b20fe017c16e8a59cd2ab44e3c37b24cea1f0a21/config/default/manager_visibility_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/default/manager_webhook_patch.yaml`](https://github.com/red-hat-data-services/kueue/blob/b20fe017c16e8a59cd2ab44e3c37b24cea1f0a21/config/default/manager_webhook_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/dev/manager_config_patch.yaml`](https://github.com/red-hat-data-services/kueue/blob/b20fe017c16e8a59cd2ab44e3c37b24cea1f0a21/config/dev/manager_config_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/rhoai/manager_config_patch.yaml`](https://github.com/red-hat-data-services/kueue/blob/b20fe017c16e8a59cd2ab44e3c37b24cea1f0a21/config/rhoai/manager_config_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/rhoai/manager_metrics_patch.yaml`](https://github.com/red-hat-data-services/kueue/blob/b20fe017c16e8a59cd2ab44e3c37b24cea1f0a21/config/rhoai/manager_metrics_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/rhoai/manager_webhook_patch.yaml`](https://github.com/red-hat-data-services/kueue/blob/b20fe017c16e8a59cd2ab44e3c37b24cea1f0a21/config/rhoai/manager_webhook_patch.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `Dockerfile` | ${BASE_IMAGE} | 2 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${BUILDER_IMAGE}; Unpinned base image: ${BASE_IMAGE} |
| `Dockerfile.konflux` | registry.access.redhat.com/ubi9/ubi-minimal@sha256:fe9e574f04371b333ed4e21d30d984f6b7fcd1046e579f5ddab4816c0c8e231d | 2 | 65532:65532 |  |  |  |  |
| `Dockerfile.rhoai` | registry.access.redhat.com/ubi9/ubi:latest | 3 | 65532:65532 |  |  |  | Unpinned base image: ${GOLANG_IMAGE}; Unpinned base image: registry.access.redhat.com/ubi9/ubi:latest; Unpinned base image: registry.access.redhat.com/ubi9/ubi:latest |
| `cmd/experimental/kueue-viz/backend/Dockerfile` | ${BASE_IMAGE} | 2 | 65532:65532 |  |  |  | Unpinned base image: ${BUILDER_IMAGE}; Unpinned base image: ${BASE_IMAGE} |
| `cmd/experimental/kueue-viz/frontend/Dockerfile` | node:23 | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `cmd/importer/Dockerfile` | ${BASE_IMAGE} | 2 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${BUILDER_IMAGE}; Unpinned base image: ${BASE_IMAGE} |

