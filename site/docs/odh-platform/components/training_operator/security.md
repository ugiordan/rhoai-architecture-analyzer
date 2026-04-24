# training-operator: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

### Secrets Referenced

| Name | Type | Referenced By |
|------|------|---------------|
| kubeflow-training-operator-webhook-cert | Opaque | deployment/training-operator |
| training-operator-webhook-cert | Opaque | deployment/training-operator |

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| training-operator | training-operator | ? | ? | ? | [`manifests/base/deployment.yaml`](https://github.com/kubeflow/training-operator/blob/d843b21596b8c611b9747fcbf3aa1a3fc23e8625/manifests/base/deployment.yaml) |
| training-operator | training-operator | ? | ? | ? | [`manifests/rhoai/manager_config_patch.yaml`](https://github.com/kubeflow/training-operator/blob/d843b21596b8c611b9747fcbf3aa1a3fc23e8625/manifests/rhoai/manager_config_patch.yaml) |
| training-operator | training-operator | ? | ? | ? | [`manifests/rhoai/manager_metrics_patch.yaml`](https://github.com/kubeflow/training-operator/blob/d843b21596b8c611b9747fcbf3aa1a3fc23e8625/manifests/rhoai/manager_metrics_patch.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `build/images/kubectl-delivery/Dockerfile` | alpine:3.17 | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `build/images/training-operator/Dockerfile` | gcr.io/distroless/static:latest | 2 |  |  |  |  | Unpinned base image: gcr.io/distroless/static:latest; No USER directive found (defaults to root) |
| `build/images/training-operator/Dockerfile.multiarch` | registry.access.redhat.com/ubi9/ubi-minimal:latest | 1 | 65532:65532 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest |
| `build/images/training-operator/Dockerfile.rhoai` | registry.access.redhat.com/ubi9/ubi:latest | 2 | 65532:65532 |  |  |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi:latest |
| `examples/jax/cpu-demo/Dockerfile` | python:3.12 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `examples/jax/jax-dist-spmd-mnist/Dockerfile` | python:3.13 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `examples/pytorch/cpu-demo/Dockerfile` | python:3.8 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `examples/pytorch/deepspeed-demo/Dockerfile` | deepspeed/deepspeed:v072_torch112_cu117 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `examples/pytorch/elastic/echo/Dockerfile` | python:3.8-buster | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `examples/pytorch/elastic/imagenet/Dockerfile` | $BASE_IMAGE | 1 | root |  |  |  | Unpinned base image: $BASE_IMAGE; Container runs as root user |
| `examples/pytorch/mnist/Dockerfile` | nvcr.io/nvidia/pytorch:24.01-py3 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `examples/pytorch/smoke-dist/Dockerfile` | nvcr.io/nvidia/pytorch:24.01-py3 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `examples/tensorflow/dist-mnist/Dockerfile` | tensorflow/tensorflow:2.17.0 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `examples/tensorflow/distribution_strategy/Dockerfile` | python:3.9 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `examples/tensorflow/mnist_with_summaries/Dockerfile` | tensorflow/tensorflow:2.17.0 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `examples/tensorflow/tf_sample/Dockerfile` | tensorflow/tensorflow:2.17.0 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `examples/xgboost/lightgbm-dist/Dockerfile` | python:3.7 | 2 |  |  |  |  | No USER directive found (defaults to root) |
| `examples/xgboost/smoke-dist/Dockerfile` | python:3.6 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `examples/xgboost/xgboost-dist/Dockerfile` | python:3.6 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `sdk/python/Dockerfile.conformance` | python:3.10-alpine | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `sdk/python/kubeflow/storage_initializer/Dockerfile` | python:3.11 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `sdk/python/kubeflow/trainer/Dockerfile` | nvcr.io/nvidia/pytorch:24.06-py3 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `sdk/python/kubeflow/trainer/Dockerfile.cpu` | python:3.11 | 1 |  |  |  |  | No USER directive found (defaults to root) |

