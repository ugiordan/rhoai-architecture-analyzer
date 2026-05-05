# notebooks: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| notebook | notebook | ? | ? | ? | [`jupyter/datascience/ubi9-python-3.12/kustomize/base/statefulset.yaml`](https://github.com/red-hat-data-services/notebooks/blob/014d9abfd61762da8aa4781033a05bc0b1ec21b7/jupyter/datascience/ubi9-python-3.12/kustomize/base/statefulset.yaml) |
| notebook | notebook | ? | ? | ? | [`jupyter/minimal/ubi9-python-3.12/kustomize/base/statefulset.yaml`](https://github.com/red-hat-data-services/notebooks/blob/014d9abfd61762da8aa4781033a05bc0b1ec21b7/jupyter/minimal/ubi9-python-3.12/kustomize/base/statefulset.yaml) |
| notebook | notebook | ? | ? | ? | [`jupyter/pytorch/ubi9-python-3.12/kustomize/base/statefulset.yaml`](https://github.com/red-hat-data-services/notebooks/blob/014d9abfd61762da8aa4781033a05bc0b1ec21b7/jupyter/pytorch/ubi9-python-3.12/kustomize/base/statefulset.yaml) |
| notebook | notebook | ? | ? | ? | [`jupyter/pytorch+llmcompressor/ubi9-python-3.12/kustomize/base/statefulset.yaml`](https://github.com/red-hat-data-services/notebooks/blob/014d9abfd61762da8aa4781033a05bc0b1ec21b7/jupyter/pytorch+llmcompressor/ubi9-python-3.12/kustomize/base/statefulset.yaml) |
| notebook | notebook | ? | ? | ? | [`jupyter/rocm/pytorch/ubi9-python-3.12/kustomize/base/statefulset.yaml`](https://github.com/red-hat-data-services/notebooks/blob/014d9abfd61762da8aa4781033a05bc0b1ec21b7/jupyter/rocm/pytorch/ubi9-python-3.12/kustomize/base/statefulset.yaml) |
| notebook | notebook | ? | ? | ? | [`jupyter/rocm/tensorflow/ubi9-python-3.12/kustomize/base/statefulset.yaml`](https://github.com/red-hat-data-services/notebooks/blob/014d9abfd61762da8aa4781033a05bc0b1ec21b7/jupyter/rocm/tensorflow/ubi9-python-3.12/kustomize/base/statefulset.yaml) |
| notebook | notebook | ? | ? | ? | [`jupyter/tensorflow/ubi9-python-3.12/kustomize/base/statefulset.yaml`](https://github.com/red-hat-data-services/notebooks/blob/014d9abfd61762da8aa4781033a05bc0b1ec21b7/jupyter/tensorflow/ubi9-python-3.12/kustomize/base/statefulset.yaml) |
| notebook | notebook | ? | ? | ? | [`jupyter/trustyai/ubi9-python-3.12/kustomize/base/statefulset.yaml`](https://github.com/red-hat-data-services/notebooks/blob/014d9abfd61762da8aa4781033a05bc0b1ec21b7/jupyter/trustyai/ubi9-python-3.12/kustomize/base/statefulset.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `.devcontainer/Dockerfile.dev` | registry.fedoraproject.org/fedora:${FEDORA_TAG} | 1 |  |  |  |  | No USER directive found (defaults to root) |
| `base-images/cpu/c9s-python-3.12/Dockerfile.cpu` | quay.io/centos/centos:stream9 | 2 | ${CNB_USER_ID}:${CNB_GROUP_ID} |  | multi-arch |  |  |
| `base-images/cpu/ubi9-python-3.12/Dockerfile.cpu` | registry.access.redhat.com/ubi9/python-312:latest | 2 | 1001 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/python-312:latest; Unpinned base image: registry.access.redhat.com/ubi9/python-312:latest |
| `base-images/cuda/12.6/c9s-python-3.11/Dockerfile.cuda` | cuda-base-${TARGETARCH} | 5 | 1001 |  | multi-arch |  | Unpinned base image: base; Unpinned base image: base; Unpinned base image: cuda-base-${TARGETARCH} |
| `base-images/cuda/12.6/c9s-python-3.12/Dockerfile.cuda` | cuda-base-${TARGETARCH} | 5 | 1001 |  | multi-arch |  | Unpinned base image: base; Unpinned base image: base; Unpinned base image: cuda-base-${TARGETARCH} |
| `base-images/cuda/12.6/ubi9-python-3.12/Dockerfile.cuda` | cuda-base-${TARGETARCH} | 5 | 1001 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest; Unpinned base image: registry.access.redhat.com/ubi9/python-312:latest; Unpinned base image: base; Unpinned base image: base; Unpinned base image: cuda-base-${TARGETARCH} |
| `base-images/cuda/12.8/c9s-python-3.12/Dockerfile.cuda` | cuda-base-${TARGETARCH} | 5 | 1001 |  | multi-arch |  | Unpinned base image: base; Unpinned base image: base; Unpinned base image: cuda-base-${TARGETARCH} |
| `base-images/cuda/12.8/ubi9-python-3.12/Dockerfile.cuda` | cuda-base-${TARGETARCH} | 5 | 1001 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest; Unpinned base image: registry.access.redhat.com/ubi9/python-312:latest; Unpinned base image: base; Unpinned base image: base; Unpinned base image: cuda-base-${TARGETARCH} |
| `base-images/cuda/12.9/c9s-python-3.12/Dockerfile.cuda` | cuda-base-${TARGETARCH} | 5 | 1001 |  | multi-arch |  | Unpinned base image: base; Unpinned base image: base; Unpinned base image: cuda-base-${TARGETARCH} |
| `base-images/cuda/12.9/ubi9-python-3.12/Dockerfile.cuda` | cuda-base-${TARGETARCH} | 5 | 1001 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest; Unpinned base image: registry.access.redhat.com/ubi9/python-312:latest; Unpinned base image: base; Unpinned base image: base; Unpinned base image: cuda-base-${TARGETARCH} |
| `base-images/cuda/13.0/c9s-python-3.12/Dockerfile.cuda` | cuda-base-${TARGETARCH} | 5 | 1001 |  | multi-arch |  | Unpinned base image: base; Unpinned base image: base; Unpinned base image: cuda-base-${TARGETARCH} |
| `base-images/rocm/6.2/c9s-python-3.12/Dockerfile.rocm` | base | 3 | 1001 |  | multi-arch |  | Unpinned base image: base |
| `base-images/rocm/6.2/ubi9-python-3.12/Dockerfile.rocm` | base | 3 | 1001 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest; Unpinned base image: registry.access.redhat.com/ubi9/python-312:latest; Unpinned base image: base |
| `base-images/rocm/6.3/c9s-python-3.12/Dockerfile.rocm` | base | 3 | 1001 |  | multi-arch |  | Unpinned base image: base |
| `base-images/rocm/6.3/ubi9-python-3.12/Dockerfile.rocm` | base | 3 | 1001 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest; Unpinned base image: registry.access.redhat.com/ubi9/python-312:latest; Unpinned base image: base |
| `base-images/rocm/6.4/c9s-python-3.12/Dockerfile.rocm` | base | 3 | 1001 |  | multi-arch |  | Unpinned base image: base |
| `base-images/rocm/6.4/ubi9-python-3.12/Dockerfile.rocm` | base | 3 | 1001 |  | multi-arch |  | Unpinned base image: registry.access.redhat.com/ubi9/ubi-minimal:latest; Unpinned base image: registry.access.redhat.com/ubi9/python-312:latest; Unpinned base image: base |
| `base-images/rocm/7.1/c9s-python-3.12/Dockerfile.rocm` | base | 3 | 1001 |  | multi-arch |  | Unpinned base image: base |
| `codeserver/ubi9-python-3.12/Dockerfile.cpu` | codeserver | 6 | 1001 |  | multi-arch |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: registry.access.redhat.com/ubi9/python-312:latest; Unpinned base image: ${BASE_IMAGE}; Unpinned base image: cpu-base; Unpinned base image: codeserver; Unpinned base image: codeserver |
| `codeserver/ubi9-python-3.12/Dockerfile.konflux.cpu` | codeserver | 6 | 1001 |  | multi-arch |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: registry.access.redhat.com/ubi9/python-312:latest; Unpinned base image: ${BASE_IMAGE}; Unpinned base image: cpu-base; Unpinned base image: codeserver; Unpinned base image: codeserver |
| `jupyter/datascience/ubi9-python-3.12/Dockerfile.cpu` | jupyter-minimal | 4 | 1001 |  |  |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: cpu-base; Unpinned base image: jupyter-minimal |
| `jupyter/datascience/ubi9-python-3.12/Dockerfile.konflux.cpu` | jupyter-minimal | 4 | 1001 |  |  |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: cpu-base; Unpinned base image: jupyter-minimal |
| `jupyter/minimal/ubi9-python-3.12/Dockerfile.cpu` | cpu-base | 2 | 1001 |  |  |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: cpu-base |
| `jupyter/minimal/ubi9-python-3.12/Dockerfile.cuda` | cuda-base | 2 | 1001 |  | multi-arch |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: cuda-base |
| `jupyter/minimal/ubi9-python-3.12/Dockerfile.konflux.cpu` | cpu-base | 2 | 1001 |  |  |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: cpu-base |
| `jupyter/minimal/ubi9-python-3.12/Dockerfile.konflux.cuda` | cuda-base | 2 | 1001 |  | multi-arch |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: cuda-base |
| `jupyter/minimal/ubi9-python-3.12/Dockerfile.konflux.rocm` | rocm-base | 2 | 1001 |  |  |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: rocm-base |
| `jupyter/minimal/ubi9-python-3.12/Dockerfile.rocm` | rocm-base | 2 | 1001 |  |  |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: rocm-base |
| `jupyter/pytorch+llmcompressor/ubi9-python-3.12/Dockerfile.cuda` | cuda-jupyter-datascience | 5 | 1001 |  | multi-arch |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: cuda-base; Unpinned base image: cuda-jupyter-minimal; Unpinned base image: cuda-jupyter-datascience |
| `jupyter/pytorch+llmcompressor/ubi9-python-3.12/Dockerfile.konflux.cuda` | cuda-jupyter-datascience | 5 | 1001 |  | multi-arch |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: cuda-base; Unpinned base image: cuda-jupyter-minimal; Unpinned base image: cuda-jupyter-datascience |
| `jupyter/pytorch/ubi9-python-3.12/Dockerfile.cuda` | cuda-jupyter-datascience | 5 | 1001 |  | multi-arch |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: cuda-base; Unpinned base image: cuda-jupyter-minimal; Unpinned base image: cuda-jupyter-datascience |
| `jupyter/pytorch/ubi9-python-3.12/Dockerfile.konflux.cuda` | cuda-jupyter-datascience | 5 | 1001 |  | multi-arch |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: cuda-base; Unpinned base image: cuda-jupyter-minimal; Unpinned base image: cuda-jupyter-datascience |
| `jupyter/rocm/pytorch/ubi9-python-3.12/Dockerfile.konflux.rocm` | rocm-jupyter-datascience | 5 | 1001 |  |  |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: rocm-base; Unpinned base image: rocm-jupyter-minimal; Unpinned base image: rocm-jupyter-datascience |
| `jupyter/rocm/pytorch/ubi9-python-3.12/Dockerfile.rocm` | rocm-jupyter-datascience | 5 | 1001 |  |  |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: rocm-base; Unpinned base image: rocm-jupyter-minimal; Unpinned base image: rocm-jupyter-datascience |
| `jupyter/rocm/tensorflow/ubi9-python-3.12/Dockerfile.konflux.rocm` | rocm-jupyter-datascience | 5 | 1001 |  | multi-arch |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: rocm-base; Unpinned base image: rocm-jupyter-minimal; Unpinned base image: rocm-jupyter-datascience |
| `jupyter/rocm/tensorflow/ubi9-python-3.12/Dockerfile.rocm` | rocm-jupyter-datascience | 5 | 1001 |  | multi-arch |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: rocm-base; Unpinned base image: rocm-jupyter-minimal; Unpinned base image: rocm-jupyter-datascience |
| `jupyter/tensorflow/ubi9-python-3.12/Dockerfile.cuda` | cuda-jupyter-datascience | 5 | 1001 |  | multi-arch |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: cuda-base; Unpinned base image: cuda-jupyter-minimal; Unpinned base image: cuda-jupyter-datascience |
| `jupyter/tensorflow/ubi9-python-3.12/Dockerfile.konflux.cuda` | cuda-jupyter-datascience | 5 | 1001 |  | multi-arch |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: cuda-base; Unpinned base image: cuda-jupyter-minimal; Unpinned base image: cuda-jupyter-datascience |
| `jupyter/trustyai/ubi9-python-3.12/Dockerfile.cpu` | jupyter-datascience | 5 | 1001 |  |  |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: cpu-base; Unpinned base image: jupyter-minimal; Unpinned base image: jupyter-datascience |
| `jupyter/trustyai/ubi9-python-3.12/Dockerfile.konflux.cpu` | jupyter-datascience | 5 | 1001 |  |  |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: cpu-base; Unpinned base image: jupyter-minimal; Unpinned base image: jupyter-datascience |
| `runtimes/datascience/ubi9-python-3.12/Dockerfile.cpu` | cpu-base | 2 | 1001 |  | multi-arch |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: cpu-base |
| `runtimes/datascience/ubi9-python-3.12/Dockerfile.konflux.cpu` | cpu-base | 2 | 1001 |  | multi-arch |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: cpu-base |
| `runtimes/minimal/ubi9-python-3.12/Dockerfile.cpu` | cpu-base | 2 | 1001 |  |  |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: cpu-base |
| `runtimes/minimal/ubi9-python-3.12/Dockerfile.konflux.cpu` | cpu-base | 2 | 1001 |  |  |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: cpu-base |
| `runtimes/pytorch+llmcompressor/ubi9-python-3.12/Dockerfile.cuda` | cuda-base | 2 | 1001 |  | multi-arch |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: cuda-base |
| `runtimes/pytorch+llmcompressor/ubi9-python-3.12/Dockerfile.konflux.cuda` | cuda-base | 2 | 1001 |  | multi-arch |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: cuda-base |
| `runtimes/pytorch/ubi9-python-3.12/Dockerfile.cuda` | cuda-base | 2 | 1001 |  | multi-arch |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: cuda-base |
| `runtimes/pytorch/ubi9-python-3.12/Dockerfile.konflux.cuda` | cuda-base | 2 | 1001 |  | multi-arch |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: cuda-base |
| `runtimes/rocm-pytorch/ubi9-python-3.12/Dockerfile.konflux.rocm` | rocm-base | 2 | 1001 |  |  |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: rocm-base |
| `runtimes/rocm-pytorch/ubi9-python-3.12/Dockerfile.rocm` | rocm-base | 2 | 1001 |  |  |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: rocm-base |
| `runtimes/rocm-tensorflow/ubi9-python-3.12/Dockerfile.konflux.rocm` | rocm-base | 2 | 1001 |  |  |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: rocm-base |
| `runtimes/rocm-tensorflow/ubi9-python-3.12/Dockerfile.rocm` | rocm-base | 2 | 1001 |  |  |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: rocm-base |
| `runtimes/tensorflow/ubi9-python-3.12/Dockerfile.cuda` | cuda-base | 2 | 1001 |  | multi-arch |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: cuda-base |
| `runtimes/tensorflow/ubi9-python-3.12/Dockerfile.konflux.cuda` | cuda-base | 2 | 1001 |  | multi-arch |  | Unpinned base image: ${BASE_IMAGE}; Unpinned base image: cuda-base |
| `scripts/lockfile-generators/Dockerfile.rpm-lockfile` | ${BASE_IMAGE} | 1 | root |  |  |  | Unpinned base image: ${BASE_IMAGE}; Container runs as root user |

