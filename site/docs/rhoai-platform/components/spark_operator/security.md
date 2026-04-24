# spark-operator: Security

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
| controller-manager | manager | ? | ? | ? | [`config/default/manager_webhook_patch.yaml`](https://github.com/kubeflow/spark-operator/blob/5366d3d2fe80d4a3e972ea010e3556631d52f017/config/default/manager_webhook_patch.yaml) |
| spark-operator-controller | controller | true | true | false | [`config/manager/manager.yaml`](https://github.com/kubeflow/spark-operator/blob/5366d3d2fe80d4a3e972ea010e3556631d52f017/config/manager/manager.yaml) |
| spark-operator-webhook | webhook | true | true | false | [`config/webhook/deployment.yaml`](https://github.com/kubeflow/spark-operator/blob/5366d3d2fe80d4a3e972ea010e3556631d52f017/config/webhook/deployment.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `Dockerfile` | ${SPARK_IMAGE} | 2 | ${SPARK_UID}:${SPARK_GID} |  | multi-arch |  | Unpinned base image: ${SPARK_IMAGE} |
| `Dockerfile.konflux` | ${BASE_IMAGE} | 2 | ${SPARK_UID}:${SPARK_GID} |  | multi-arch |  | Unpinned base image: ${GO_BUILDER_IMAGE}; Unpinned base image: ${BASE_IMAGE} |
| `docker/Dockerfile.kubectl` | ${BASE_IMAGE} | 2 |  |  | multi-arch |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: ${BASE_IMAGE}; No USER directive found (defaults to root) |
| `examples/openshift/Dockerfile` | apache/spark:3.5.7-java17-python3 | 1 | 0 |  |  |  | Container runs as root user |
| `examples/openshift/Dockerfile.odh` | ${BASE_IMAGE} | 2 | ${SPARK_UID}:${SPARK_GID} |  | multi-arch |  | Unpinned base image: ${GO_BUILDER_IMAGE}; Unpinned base image: ${BASE_IMAGE} |
| `spark-docker/Dockerfile` | ${SPARK_IMAGE} | 1 | ${spark_uid} |  |  |  | Unpinned base image: ${SPARK_IMAGE} |

