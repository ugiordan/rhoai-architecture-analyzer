# notebooks: Network

## Service Map

*1 unique services (8 total, duplicates from test fixtures collapsed).*

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    notebooks["notebooks"]:::component
    notebooks --> svc_0["notebook\nClusterIP: 8888/TCP"]:::svc
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| notebook | ClusterIP | 8888/TCP | [`jupyter/datascience/ubi9-python-3.12/kustomize/base/service.yaml`](https://github.com/opendatahub-io/notebooks/blob/b3803ad4e49d09613f01f8362d5455fd70da204b/jupyter/datascience/ubi9-python-3.12/kustomize/base/service.yaml) |
| notebook | ClusterIP | 8888/TCP | [`jupyter/minimal/ubi9-python-3.12/kustomize/base/service.yaml`](https://github.com/opendatahub-io/notebooks/blob/b3803ad4e49d09613f01f8362d5455fd70da204b/jupyter/minimal/ubi9-python-3.12/kustomize/base/service.yaml) |
| notebook | ClusterIP | 8888/TCP | [`jupyter/pytorch/ubi9-python-3.12/kustomize/base/service.yaml`](https://github.com/opendatahub-io/notebooks/blob/b3803ad4e49d09613f01f8362d5455fd70da204b/jupyter/pytorch/ubi9-python-3.12/kustomize/base/service.yaml) |
| notebook | ClusterIP | 8888/TCP | [`jupyter/pytorch+llmcompressor/ubi9-python-3.12/kustomize/base/service.yaml`](https://github.com/opendatahub-io/notebooks/blob/b3803ad4e49d09613f01f8362d5455fd70da204b/jupyter/pytorch+llmcompressor/ubi9-python-3.12/kustomize/base/service.yaml) |
| notebook | ClusterIP | 8888/TCP | [`jupyter/rocm/pytorch/ubi9-python-3.12/kustomize/base/service.yaml`](https://github.com/opendatahub-io/notebooks/blob/b3803ad4e49d09613f01f8362d5455fd70da204b/jupyter/rocm/pytorch/ubi9-python-3.12/kustomize/base/service.yaml) |
| notebook | ClusterIP | 8888/TCP | [`jupyter/rocm/tensorflow/ubi9-python-3.12/kustomize/base/service.yaml`](https://github.com/opendatahub-io/notebooks/blob/b3803ad4e49d09613f01f8362d5455fd70da204b/jupyter/rocm/tensorflow/ubi9-python-3.12/kustomize/base/service.yaml) |
| notebook | ClusterIP | 8888/TCP | [`jupyter/tensorflow/ubi9-python-3.12/kustomize/base/service.yaml`](https://github.com/opendatahub-io/notebooks/blob/b3803ad4e49d09613f01f8362d5455fd70da204b/jupyter/tensorflow/ubi9-python-3.12/kustomize/base/service.yaml) |
| notebook | ClusterIP | 8888/TCP | [`jupyter/trustyai/ubi9-python-3.12/kustomize/base/service.yaml`](https://github.com/opendatahub-io/notebooks/blob/b3803ad4e49d09613f01f8362d5455fd70da204b/jupyter/trustyai/ubi9-python-3.12/kustomize/base/service.yaml) |

!!! warning "No Network Policies"
    No NetworkPolicy resources were found in the analyzed sources. Network policies may exist in overlays, Helm values, or cluster-level configurations not captured by static analysis.

