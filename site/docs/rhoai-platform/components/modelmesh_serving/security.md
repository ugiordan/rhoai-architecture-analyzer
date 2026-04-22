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
| controller-manager | kube-rbac-proxy | ? | ? | ? | [`config/default/manager_auth_proxy_patch.yaml`](https://github.com/kserve/modelmesh-serving/blob/794cbf39ff585034ae3ed8b73953e65b2524a738/config/default/manager_auth_proxy_patch.yaml) |
| controller-manager | manager | ? | ? | ? | [`config/default/manager_auth_proxy_patch.yaml`](https://github.com/kserve/modelmesh-serving/blob/794cbf39ff585034ae3ed8b73953e65b2524a738/config/default/manager_auth_proxy_patch.yaml) |
| modelmesh-controller | manager | ? | ? | ? | [`config/default/manager_webhook_patch.yaml`](https://github.com/kserve/modelmesh-serving/blob/794cbf39ff585034ae3ed8b73953e65b2524a738/config/default/manager_webhook_patch.yaml) |
| modelmesh-controller | manager | ? | ? | ? | [`config/manager/manager.yaml`](https://github.com/kserve/modelmesh-serving/blob/794cbf39ff585034ae3ed8b73953e65b2524a738/config/manager/manager.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `Dockerfile` | registry.access.redhat.com/ubi9/ubi-minimal:9.5 | 2 | ${USER} |  | multi-arch |  | Unpinned base image: ${DEV_IMAGE} |
| `Dockerfile.develop` | registry.access.redhat.com/ubi9/go-toolset:$GOLANG_VERSION | 1 | root |  | multi-arch |  | Container runs as root user |
| `Dockerfile.develop.ci` | registry.access.redhat.com/ubi9/go-toolset:$GOLANG_VERSION | 1 | root |  | multi-arch |  | Container runs as root user |
| `Dockerfile.konflux` | registry.access.redhat.com/ubi9/ubi-minimal@sha256:7d4e47500f28ac3a2bff06c25eff9127ff21048538ae03ce240d57cf756acd00 | 2 | ${USER} |  | multi-arch |  |  |
| `docs/examples/python-custom-runtime/custom-model/Dockerfile` | python:3.9.13 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `tests/Dockerfile` | quay.io/centos/centos:stream8 | 1 |  |  |  |  | No USER directive found (defaults to root) |

