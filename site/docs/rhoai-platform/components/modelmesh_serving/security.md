# modelmesh-serving: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

### Secrets Referenced

| Name | Type | Referenced By |
|------|------|---------------|
| modelmesh-webhook-server-cert | Opaque | deployment/modelmesh-controller |

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| controller-manager | kube-rbac-proxy | ? | ? | ? | [`config/default/manager_auth_proxy_patch.yaml`](https://github.com/kserve/modelmesh-serving/blob/e9a796fca454a6cea0134c4525eb926aae9c2691/config/default/manager_auth_proxy_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/default/manager_auth_proxy_patch.yaml`](https://github.com/kserve/modelmesh-serving/blob/e9a796fca454a6cea0134c4525eb926aae9c2691/config/default/manager_auth_proxy_patch.yaml) |
| etcd | etcd | ? | ? | ? | [`kustomize:config/overlays/odh`](https://github.com/kserve/modelmesh-serving/blob/e9a796fca454a6cea0134c4525eb926aae9c2691/kustomize:config/overlays/odh) |
| modelmesh-controller | manager | ? | ? | ? | [`kustomize:config/overlays/odh`](https://github.com/kserve/modelmesh-serving/blob/e9a796fca454a6cea0134c4525eb926aae9c2691/kustomize:config/overlays/odh) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `Dockerfile` | registry.access.redhat.com/ubi9/ubi-minimal:9.5 | 2 | ${USER} |  | multi-arch |  | Unpinned base image: ${DEV_IMAGE} |
| `Dockerfile.develop` | registry.access.redhat.com/ubi9/go-toolset:$GOLANG_VERSION | 1 | root |  | multi-arch |  | Container runs as root user |
| `Dockerfile.develop.ci` | registry.access.redhat.com/ubi9/go-toolset:$GOLANG_VERSION | 1 | root |  | multi-arch |  | Container runs as root user |
| `Dockerfile.konflux` | registry.access.redhat.com/ubi9/ubi-minimal@sha256:fe9e574f04371b333ed4e21d30d984f6b7fcd1046e579f5ddab4816c0c8e231d | 2 | ${USER} |  | multi-arch |  |  |

