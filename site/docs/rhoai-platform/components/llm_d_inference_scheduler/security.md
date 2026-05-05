# llm-d-inference-scheduler: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

### Secrets Referenced

| Name | Type | Referenced By |
|------|------|---------------|
| cacerts | Opaque | deployment/istiod-llm-d-gateway |
| istio-kubeconfig | Opaque | deployment/istiod-llm-d-gateway |
| istiod-tls | Opaque | deployment/istiod-llm-d-gateway |

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| ${EPP_NAME} | epp | ? | ? | ? | [`deploy/components/inference-gateway/deployments.yaml`](https://github.com/llm-d/llm-d-inference-scheduler/blob/eb2ef5d06644cdf1726fcbc3276d41d8f91f70eb/deploy/components/inference-gateway/deployments.yaml) |
| ${EPP_NAME} | uds-tokenizer | ? | ? | ? | [`deploy/components/inference-gateway/deployments.yaml`](https://github.com/llm-d/llm-d-inference-scheduler/blob/eb2ef5d06644cdf1726fcbc3276d41d8f91f70eb/deploy/components/inference-gateway/deployments.yaml) |
| ${MODEL_NAME_SAFE}-vllm-sim | vllm | ? | ? | ? | [`deploy/components/vllm-sim/deployments.yaml`](https://github.com/llm-d/llm-d-inference-scheduler/blob/eb2ef5d06644cdf1726fcbc3276d41d8f91f70eb/deploy/components/vllm-sim/deployments.yaml) |
| 0 | cmd | ? | ? | ? | [`deploy/environments/kubernetes-base/common/statefulset.yaml`](https://github.com/llm-d/llm-d-inference-scheduler/blob/eb2ef5d06644cdf1726fcbc3276d41d8f91f70eb/deploy/environments/kubernetes-base/common/statefulset.yaml) |
| istiod-llm-d-gateway | discovery | true | true | ? | [`deploy/components/istio-control-plane/deployments.yaml`](https://github.com/llm-d/llm-d-inference-scheduler/blob/eb2ef5d06644cdf1726fcbc3276d41d8f91f70eb/deploy/components/istio-control-plane/deployments.yaml) |
| vllm-sim-d | vllm | ? | ? | ? | [`deploy/components/vllm-sim-epd/deployments.yaml`](https://github.com/llm-d/llm-d-inference-scheduler/blob/eb2ef5d06644cdf1726fcbc3276d41d8f91f70eb/deploy/components/vllm-sim-epd/deployments.yaml) |
| vllm-sim-d | vllm | ? | ? | ? | [`deploy/components/vllm-sim-pd/deployments.yaml`](https://github.com/llm-d/llm-d-inference-scheduler/blob/eb2ef5d06644cdf1726fcbc3276d41d8f91f70eb/deploy/components/vllm-sim-pd/deployments.yaml) |
| vllm-sim-e | vllm | ? | ? | ? | [`deploy/components/vllm-sim-epd/deployments.yaml`](https://github.com/llm-d/llm-d-inference-scheduler/blob/eb2ef5d06644cdf1726fcbc3276d41d8f91f70eb/deploy/components/vllm-sim-epd/deployments.yaml) |
| vllm-sim-p | vllm | ? | ? | ? | [`deploy/components/vllm-sim-epd/deployments.yaml`](https://github.com/llm-d/llm-d-inference-scheduler/blob/eb2ef5d06644cdf1726fcbc3276d41d8f91f70eb/deploy/components/vllm-sim-epd/deployments.yaml) |
| vllm-sim-p | vllm | ? | ? | ? | [`deploy/components/vllm-sim-pd/deployments.yaml`](https://github.com/llm-d/llm-d-inference-scheduler/blob/eb2ef5d06644cdf1726fcbc3276d41d8f91f70eb/deploy/components/vllm-sim-pd/deployments.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `Dockerfile.builder` | quay.io/projectquay/golang:1.25 | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `Dockerfile.epp` | ${BASE_IMAGE} | 2 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${BASE_IMAGE} |
| `Dockerfile.epp.konflux` | registry.access.redhat.com/ubi9/ubi-minimal:9.7@sha256:d91be7cea9f03a757d69ad7fcdfcd7849dba820110e7980d5e2a1f46ed06ea3b | 2 | 65532:65532 |  | multi-arch |  |  |
| `Dockerfile.sidecar` | ${BASE_IMAGE} | 2 | 65532:65532 |  | multi-arch |  | Unpinned base image: ${BASE_IMAGE} |
| `Dockerfile.sidecar.konflux` | registry.access.redhat.com/ubi9/ubi-minimal:9.7@sha256:d91be7cea9f03a757d69ad7fcdfcd7849dba820110e7980d5e2a1f46ed06ea3b | 2 | 65532:65532 |  | multi-arch |  |  |

