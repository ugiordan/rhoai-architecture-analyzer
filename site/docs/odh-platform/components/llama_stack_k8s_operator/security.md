# llama-stack-k8s-operator: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| controller-manager | kube-rbac-proxy | ? | ? | ? | [`config/default/manager_auth_proxy_patch.yaml`](https://github.com/llamastack/llama-stack-k8s-operator/blob/ba8020a4fc5b6ac86e14aea251992ee2ccdde5ef/config/default/manager_auth_proxy_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/default/manager_auth_proxy_patch.yaml`](https://github.com/llamastack/llama-stack-k8s-operator/blob/ba8020a4fc5b6ac86e14aea251992ee2ccdde5ef/config/default/manager_auth_proxy_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/default/manager_config_patch.yaml`](https://github.com/llamastack/llama-stack-k8s-operator/blob/ba8020a4fc5b6ac86e14aea251992ee2ccdde5ef/config/default/manager_config_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/manager/manager.yaml`](https://github.com/llamastack/llama-stack-k8s-operator/blob/ba8020a4fc5b6ac86e14aea251992ee2ccdde5ef/config/manager/manager.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `Dockerfile` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 2 | 1001 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |

