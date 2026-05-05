# kubeflow: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

### Secrets Referenced

| Name | Type | Referenced By |
|------|------|---------------|
| odh-notebook-controller-webhook-cert | kubernetes.io/tls | service/webhook-service |
| webhook-server-cert | Opaque | deployment/controller-manager |

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| controller-manager | kube-rbac-proxy | ? | ? | ? | [`components/notebook-controller/config/default/manager_auth_proxy_patch.yaml`](https://github.com/red-hat-data-services/kubeflow/blob/2afcff72f7d4be00ddb2cae2f27a36cebc6d1213/components/notebook-controller/config/default/manager_auth_proxy_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`components/notebook-controller/config/default/manager_auth_proxy_patch.yaml`](https://github.com/red-hat-data-services/kubeflow/blob/2afcff72f7d4be00ddb2cae2f27a36cebc6d1213/components/notebook-controller/config/default/manager_auth_proxy_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`components/notebook-controller/config/default/manager_image_patch.yaml`](https://github.com/red-hat-data-services/kubeflow/blob/2afcff72f7d4be00ddb2cae2f27a36cebc6d1213/components/notebook-controller/config/default/manager_image_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`components/notebook-controller/config/default/manager_prometheus_metrics_patch.yaml`](https://github.com/red-hat-data-services/kubeflow/blob/2afcff72f7d4be00ddb2cae2f27a36cebc6d1213/components/notebook-controller/config/default/manager_prometheus_metrics_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`components/notebook-controller/config/default/manager_webhook_patch.yaml`](https://github.com/red-hat-data-services/kubeflow/blob/2afcff72f7d4be00ddb2cae2f27a36cebc6d1213/components/notebook-controller/config/default/manager_webhook_patch.yaml) |
| deployment | manager | ? | ? | ? | [`components/notebook-controller/config/manager/manager.yaml`](https://github.com/red-hat-data-services/kubeflow/blob/2afcff72f7d4be00ddb2cae2f27a36cebc6d1213/components/notebook-controller/config/manager/manager.yaml) |
| deployment | manager | ? | ? | ? | [`components/notebook-controller/config/overlays/openshift/manager_openshift_patch.yaml`](https://github.com/red-hat-data-services/kubeflow/blob/2afcff72f7d4be00ddb2cae2f27a36cebc6d1213/components/notebook-controller/config/overlays/openshift/manager_openshift_patch.yaml) |
| manager | manager | ? | ? | ? | [`components/odh-notebook-controller/config/manager/manager.yaml`](https://github.com/red-hat-data-services/kubeflow/blob/2afcff72f7d4be00ddb2cae2f27a36cebc6d1213/components/odh-notebook-controller/config/manager/manager.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `components/notebook-controller/Dockerfile` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 2 | 1001:0 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |
| `components/notebook-controller/Dockerfile.ci` | gcr.io/distroless/base:debug | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `components/notebook-controller/Dockerfile.konflux` | registry.access.redhat.com/ubi9/ubi-minimal@sha256:8d0a8fb39ec907e8ca62cdd24b62a63ca49a30fe465798a360741fde58437a23 | 2 | 1001:0 |  | multi-arch |  |  |
| `components/odh-notebook-controller/Dockerfile` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 2 | 1001:0 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |
| `components/odh-notebook-controller/Dockerfile.konflux` | registry.access.redhat.com/ubi9/ubi-minimal@sha256:8d0a8fb39ec907e8ca62cdd24b62a63ca49a30fe465798a360741fde58437a23 | 2 | 1001:0 |  | multi-arch |  |  |

