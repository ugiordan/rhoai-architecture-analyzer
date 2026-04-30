# models-as-a-service: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

### Secrets Referenced

| Name | Type | Referenced By |
|------|------|---------------|
| maas-api-serving-cert | kubernetes.io/tls | deployment/maas-api, service/maas-api |

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| maas-api | maas-api | true | true | ? | [`deployment/base/maas-api/core/deployment.yaml`](https://github.com/red-hat-data-services/models-as-a-service/blob/28670ac4b5807bc0d4b38ea40ae66616d16bcc3e/deployment/base/maas-api/core/deployment.yaml) |
| maas-api | maas-api | ? | ? | ? | [`deployment/base/maas-api/overlays/tls/deployment-patch.yaml`](https://github.com/red-hat-data-services/models-as-a-service/blob/28670ac4b5807bc0d4b38ea40ae66616d16bcc3e/deployment/base/maas-api/overlays/tls/deployment-patch.yaml) |
| maas-controller | manager | ? | ? | ? | [`deployment/base/maas-controller/manager/manager.yaml`](https://github.com/red-hat-data-services/models-as-a-service/blob/28670ac4b5807bc0d4b38ea40ae66616d16bcc3e/deployment/base/maas-controller/manager/manager.yaml) |
| payload-processing | payload-processing | ? | true | ? | [`deployment/base/payload-processing/manager/deployment.yaml`](https://github.com/red-hat-data-services/models-as-a-service/blob/28670ac4b5807bc0d4b38ea40ae66616d16bcc3e/deployment/base/payload-processing/manager/deployment.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `maas-api/Dockerfile` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 2 | 1001 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |
| `maas-api/Dockerfile.konflux` | registry.access.redhat.com/ubi9/ubi-minimal@sha256:7d4e47500f28ac3a2bff06c25eff9127ff21048538ae03ce240d57cf756acd00 | 2 | 1001 |  | multi-arch |  |  |
| `maas-controller/Dockerfile` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 2 | 1001 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |
| `maas-controller/Dockerfile.konflux` | registry.access.redhat.com/ubi9/ubi-minimal@sha256:7d4e47500f28ac3a2bff06c25eff9127ff21048538ae03ce240d57cf756acd00 | 2 | 1001 |  | multi-arch |  |  |

