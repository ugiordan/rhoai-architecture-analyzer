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
| kserve-controller-manager | manager | true | true | false | [`kustomize:config/overlays/odh`](https://github.com/kserve/kserve/blob/5d509f23f903a2657dbe2394e785b3bd84c4c40d/kustomize:config/overlays/odh) |
| kserve-localmodel-controller-manager | manager | true | true | false | [`config/localmodels/manager.yaml`](https://github.com/kserve/kserve/blob/5d509f23f903a2657dbe2394e785b3bd84c4c40d/config/localmodels/manager.yaml) |
| llmisvc-controller-manager | manager | true | true | false | [`kustomize:config/overlays/odh`](https://github.com/kserve/kserve/blob/5d509f23f903a2657dbe2394e785b3bd84c4c40d/kustomize:config/overlays/odh) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `Dockerfile` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 4 | 1000:1000 |  |  |  | Unpinned base image: deps; Unpinned base image: deps; Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |
| `tools/tf2openapi/Dockerfile` | gcr.io/distroless/static:nonroot | 2 |  |  |  |  | No USER directive found (defaults to root) |

