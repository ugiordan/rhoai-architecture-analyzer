# kserve: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

### Secrets Referenced

| Name | Type | Referenced By |
|------|------|---------------|
| kserve-webhook-server-cert | Opaque | deployment/kserve-controller-manager |
| llmisvc-webhook-server-cert | Opaque | deployment/llmisvc-controller-manager |
| localmodel-webhook-server-cert | Opaque | deployment/kserve-localmodel-controller-manager |

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| kserve-controller-manager | manager | true | true | false | [`kustomize:config/overlays/all`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/kustomize:config/overlays/all) |
| kserve-controller-manager | kube-rbac-proxy | true | true | false | [`kustomize:config/overlays/all`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/kustomize:config/overlays/all) |
| kserve-localmodel-controller-manager | manager | true | true | false | [`kustomize:config/overlays/all`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/kustomize:config/overlays/all) |
| llmisvc-controller-manager | manager | true | true | false | [`kustomize:config/overlays/all`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/kustomize:config/overlays/all) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `Dockerfile` | gcr.io/distroless/static:nonroot | 4 |  |  |  |  | Unpinned base image: deps; Unpinned base image: deps; No USER directive found (defaults to root) |
| `tools/tf2openapi/Dockerfile` | gcr.io/distroless/static:nonroot | 2 |  |  |  |  | No USER directive found (defaults to root) |

