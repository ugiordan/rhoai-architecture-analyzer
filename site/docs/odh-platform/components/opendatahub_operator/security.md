# opendatahub-operator: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

### Secrets Referenced

| Name | Type | Referenced By |
|------|------|---------------|
| opendatahub-operator-controller-webhook-cert | kubernetes.io/tls | deployment/controller-manager, service/webhook-service |
| redhat-ods-operator-controller-webhook-cert | kubernetes.io/tls | deployment/rhods-operator, service/webhook-service |

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| azure-cloud-manager-operator | manager | ? | ? | ? | [`config/cloudmanager/azure/local/manager_pull_policy_patch.yaml`](https://github.com/opendatahub-io/opendatahub-operator/blob/607e20f6a959b75625a7313721aa1ced964187d6/config/cloudmanager/azure/local/manager_pull_policy_patch.yaml) |
| azure-cloud-manager-operator | manager | ? | true | ? | [`config/cloudmanager/azure/manager/manager.yaml`](https://github.com/opendatahub-io/opendatahub-operator/blob/607e20f6a959b75625a7313721aa1ced964187d6/config/cloudmanager/azure/manager/manager.yaml) |
| azure-cloud-manager-operator | manager | ? | ? | ? | [`config/cloudmanager/azure/rhoai/manager_rhoai_patch.yaml`](https://github.com/opendatahub-io/opendatahub-operator/blob/607e20f6a959b75625a7313721aa1ced964187d6/config/cloudmanager/azure/rhoai/manager_rhoai_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/rhaii/odh-operator/manager_webhook_patch.yaml`](https://github.com/opendatahub-io/opendatahub-operator/blob/607e20f6a959b75625a7313721aa1ced964187d6/config/rhaii/odh-operator/manager_webhook_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/rhaii/odh-operator/manager_auth_proxy_patch.yaml`](https://github.com/opendatahub-io/opendatahub-operator/blob/607e20f6a959b75625a7313721aa1ced964187d6/config/rhaii/odh-operator/manager_auth_proxy_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/rhaii/odh-operator/manager_patch.yaml`](https://github.com/opendatahub-io/opendatahub-operator/blob/607e20f6a959b75625a7313721aa1ced964187d6/config/rhaii/odh-operator/manager_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/default/manager_auth_proxy_patch.yaml`](https://github.com/opendatahub-io/opendatahub-operator/blob/607e20f6a959b75625a7313721aa1ced964187d6/config/default/manager_auth_proxy_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/default/manager_webhook_patch.yaml`](https://github.com/opendatahub-io/opendatahub-operator/blob/607e20f6a959b75625a7313721aa1ced964187d6/config/default/manager_webhook_patch.yaml) |
| controller-manager | manager | ? | true | ? | [`config/manager/manager.yaml`](https://github.com/opendatahub-io/opendatahub-operator/blob/607e20f6a959b75625a7313721aa1ced964187d6/config/manager/manager.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/rhaii/odh-local/manager_pull_policy_patch.yaml`](https://github.com/opendatahub-io/opendatahub-operator/blob/607e20f6a959b75625a7313721aa1ced964187d6/config/rhaii/odh-local/manager_pull_policy_patch.yaml) |
| coreweave-cloud-manager-operator | manager | ? | true | ? | [`config/cloudmanager/coreweave/manager/manager.yaml`](https://github.com/opendatahub-io/opendatahub-operator/blob/607e20f6a959b75625a7313721aa1ced964187d6/config/cloudmanager/coreweave/manager/manager.yaml) |
| coreweave-cloud-manager-operator | manager | ? | ? | ? | [`config/cloudmanager/coreweave/rhoai/manager_rhoai_patch.yaml`](https://github.com/opendatahub-io/opendatahub-operator/blob/607e20f6a959b75625a7313721aa1ced964187d6/config/cloudmanager/coreweave/rhoai/manager_rhoai_patch.yaml) |
| coreweave-cloud-manager-operator | manager | ? | ? | ? | [`config/cloudmanager/coreweave/local/manager_pull_policy_patch.yaml`](https://github.com/opendatahub-io/opendatahub-operator/blob/607e20f6a959b75625a7313721aa1ced964187d6/config/cloudmanager/coreweave/local/manager_pull_policy_patch.yaml) |
| rhods-operator | rhods-operator | ? | ? | ? | [`config/rhaii/rhoai/operator/manager_auth_proxy_patch.yaml`](https://github.com/opendatahub-io/opendatahub-operator/blob/607e20f6a959b75625a7313721aa1ced964187d6/config/rhaii/rhoai/operator/manager_auth_proxy_patch.yaml) |
| rhods-operator | rhods-operator | ? | ? | ? | [`config/rhaii/rhoai/operator/manager_patch.yaml`](https://github.com/opendatahub-io/opendatahub-operator/blob/607e20f6a959b75625a7313721aa1ced964187d6/config/rhaii/rhoai/operator/manager_patch.yaml) |
| rhods-operator | rhods-operator | ? | ? | ? | [`config/rhaii/rhoai/operator/manager_webhook_patch.yaml`](https://github.com/opendatahub-io/opendatahub-operator/blob/607e20f6a959b75625a7313721aa1ced964187d6/config/rhaii/rhoai/operator/manager_webhook_patch.yaml) |
| rhods-operator | rhods-operator | ? | ? | ? | [`config/rhoai/default/manager_auth_proxy_patch.yaml`](https://github.com/opendatahub-io/opendatahub-operator/blob/607e20f6a959b75625a7313721aa1ced964187d6/config/rhoai/default/manager_auth_proxy_patch.yaml) |
| rhods-operator | rhods-operator | ? | ? | ? | [`config/rhoai/default/manager_webhook_patch.yaml`](https://github.com/opendatahub-io/opendatahub-operator/blob/607e20f6a959b75625a7313721aa1ced964187d6/config/rhoai/default/manager_webhook_patch.yaml) |
| rhods-operator | rhods-operator | ? | true | ? | [`config/rhoai/manager/manager.yaml`](https://github.com/opendatahub-io/opendatahub-operator/blob/607e20f6a959b75625a7313721aa1ced964187d6/config/rhoai/manager/manager.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `Dockerfiles/Dockerfile` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 3 | 1001 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/toolbox; Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |

