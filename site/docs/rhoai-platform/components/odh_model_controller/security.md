# odh-model-controller: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

### Secrets Referenced

| Name | Type | Referenced By |
|------|------|---------------|
| model-serving-api-tls | kubernetes.io/tls | service/model-serving-api |
| odh-model-controller-webhook-cert | kubernetes.io/tls | deployment/odh-model-controller, service/odh-model-controller-webhook-service |

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| odh-model-controller | manager | ? | ? | ? | [`config/default/manager_webhook_patch.yaml`](https://github.com/opendatahub-io/odh-model-controller/blob/55c98bf18a3fa4334d31305c836593a7f6dc4d6d/config/default/manager_webhook_patch.yaml) |
| odh-model-controller | manager | ? | ? | ? | [`config/manager/manager.yaml`](https://github.com/opendatahub-io/odh-model-controller/blob/55c98bf18a3fa4334d31305c836593a7f6dc4d6d/config/manager/manager.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `Containerfile` | registry.access.redhat.com/ubi9/ubi-minimal:9.5 | 2 | 65532:65532 |  | multi-arch |  |  |
| `Containerfile.server` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 2 | 1000:1000 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |
| `Containerfile.server.konflux` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 2 | ${USER} |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |
| `Dockerfile.konflux` | registry.access.redhat.com/ubi9/ubi-minimal@sha256:8d0a8fb39ec907e8ca62cdd24b62a63ca49a30fe465798a360741fde58437a23 | 2 | ${USER} |  | multi-arch |  |  |
| `dev_tools/Containerfile.devspace` | registry.access.redhat.com/ubi9/go-toolset:1.25@sha256:8c5aeac74b4b60dc2e5e44f6b639186b7ec2fec8f0eb9a36d4a32dcf8e255f52 | 1 | root |  |  |  | Container runs as root user |

