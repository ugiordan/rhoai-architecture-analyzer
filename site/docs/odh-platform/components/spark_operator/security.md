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
| controller-manager | manager | ? | ? | ? | [`config/default/manager_webhook_patch.yaml`](https://github.com/kubeflow/spark-operator/blob/b8a995788a0bd700354170600d0813db8a597241/config/default/manager_webhook_patch.yaml) |
| spark-operator-controller | controller | true | true | false | [`kustomize:config/overlays/odh`](https://github.com/kubeflow/spark-operator/blob/b8a995788a0bd700354170600d0813db8a597241/kustomize:config/overlays/odh) |
| spark-operator-webhook | webhook | true | true | false | [`kustomize:config/overlays/odh`](https://github.com/kubeflow/spark-operator/blob/b8a995788a0bd700354170600d0813db8a597241/kustomize:config/overlays/odh) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `Dockerfile` | ${SPARK_IMAGE} | 2 | ${SPARK_UID}:${SPARK_GID} |  | multi-arch |  | Unpinned base image: ${SPARK_IMAGE} |
| `docker/Dockerfile.kubectl` | ${BASE_IMAGE} | 2 |  |  | multi-arch |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: ${BASE_IMAGE}; No USER directive found (defaults to root) |
| `spark-docker/Dockerfile` | ${SPARK_IMAGE} | 1 | ${spark_uid} |  |  |  | Unpinned base image: ${SPARK_IMAGE} |

