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
| kserve-controller-manager | kube-rbac-proxy | true | true | false | [`config/default/manager_auth_proxy_patch.yaml`](https://github.com/kserve/kserve/blob/50337bf162feeeca021123b3822076ef93d66731/config/default/manager_auth_proxy_patch.yaml) |
| kserve-controller-manager | manager | ? | ? | ? | [`config/default/manager_auth_proxy_patch.yaml`](https://github.com/kserve/kserve/blob/50337bf162feeeca021123b3822076ef93d66731/config/default/manager_auth_proxy_patch.yaml) |
| kserve-controller-manager | manager | ? | ? | ? | [`config/default/manager_image_patch.yaml`](https://github.com/kserve/kserve/blob/50337bf162feeeca021123b3822076ef93d66731/config/default/manager_image_patch.yaml) |
| kserve-controller-manager | manager | ? | ? | ? | [`config/default/manager_prometheus_metrics_patch.yaml`](https://github.com/kserve/kserve/blob/50337bf162feeeca021123b3822076ef93d66731/config/default/manager_prometheus_metrics_patch.yaml) |
| kserve-controller-manager | manager | ? | ? | ? | [`config/default/manager_resources_patch.yaml`](https://github.com/kserve/kserve/blob/50337bf162feeeca021123b3822076ef93d66731/config/default/manager_resources_patch.yaml) |
| kserve-controller-manager | manager | true | true | false | [`config/manager/manager.yaml`](https://github.com/kserve/kserve/blob/50337bf162feeeca021123b3822076ef93d66731/config/manager/manager.yaml) |
| kserve-controller-manager | manager | ? | ? | ? | [`config/overlays/test/manager_image_patch.yaml`](https://github.com/kserve/kserve/blob/50337bf162feeeca021123b3822076ef93d66731/config/overlays/test/manager_image_patch.yaml) |
| kserve-controller-manager | manager | ? | ? | ? | [`config/overlays/version-template/manager_image_patch.yaml`](https://github.com/kserve/kserve/blob/50337bf162feeeca021123b3822076ef93d66731/config/overlays/version-template/manager_image_patch.yaml) |
| kserve-localmodel-controller-manager | manager | true | true | false | [`config/localmodels/manager.yaml`](https://github.com/kserve/kserve/blob/50337bf162feeeca021123b3822076ef93d66731/config/localmodels/manager.yaml) |
| llmisvc-controller-manager | manager | true | true | false | [`config/llmisvc/manager.yaml`](https://github.com/kserve/kserve/blob/50337bf162feeeca021123b3822076ef93d66731/config/llmisvc/manager.yaml) |
| spark-pmml-iris | kfserving-container | ? | ? | ? | [`docs/samples/v1beta1/spark/deployment.yaml`](https://github.com/kserve/kserve/blob/50337bf162feeeca021123b3822076ef93d66731/docs/samples/v1beta1/spark/deployment.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `Dockerfile` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 4 | 1000:1000 |  |  |  | Unpinned base image: deps; Unpinned base image: deps; Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |
| `docs/apis/Dockerfile` | registry.access.redhat.com/ubi9/go-toolset:1.25 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `docs/kfp/Dockerfile` | python:3.9-slim-bullseye | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `docs/samples/explanation/aif/germancredit/server/Dockerfile` | python:3.10-slim | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `docs/samples/graph/bgtest/Dockerfile` | scratch | 2 |  |  |  |  | Unpinned base image: scratch; No USER directive found (defaults to root) |
| `docs/samples/v1beta1/custom/paddleserving/Dockerfile` | registry.baidubce.com/paddlepaddle/serving:0.5.0-devel | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `docs/samples/v1beta1/custom/torchserve/torchserve-image/Dockerfile` | ${BASE_IMAGE} | 2 | model-server |  |  |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: ${BASE_IMAGE} |
| `docs/samples/v1beta1/torchserve/model-archiver/model-archiver-image/Dockerfile` | ubuntu:18.04 | 1 | model-server |  |  |  |  |
| `docs/samples/v1beta1/triton/fastertransformer/transformer/Dockerfile` | python:3.9-slim-bullseye | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `hack/kserve_migration/Dockerfile` | alpine:latest | 1 |  |  |  |  | Unpinned base image: alpine:latest; No USER directive found (defaults to root) |
| `tools/tf2openapi/Dockerfile` | gcr.io/distroless/static:nonroot | 2 |  |  |  |  | No USER directive found (defaults to root) |

