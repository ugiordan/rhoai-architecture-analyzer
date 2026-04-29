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
| notebook | ClusterIP | 8888/TCP | [`jupyter/datascience/ubi9-python-3.12/kustomize/base/service.yaml`](https://github.com/red-hat-data-services/notebooks/blob/1fb86d8e3d342f8b789342720ef8fdc57c03d6a8/jupyter/datascience/ubi9-python-3.12/kustomize/base/service.yaml) |
| notebook | ClusterIP | 8888/TCP | [`jupyter/minimal/ubi9-python-3.12/kustomize/base/service.yaml`](https://github.com/red-hat-data-services/notebooks/blob/1fb86d8e3d342f8b789342720ef8fdc57c03d6a8/jupyter/minimal/ubi9-python-3.12/kustomize/base/service.yaml) |
| notebook | ClusterIP | 8888/TCP | [`jupyter/pytorch/ubi9-python-3.12/kustomize/base/service.yaml`](https://github.com/red-hat-data-services/notebooks/blob/1fb86d8e3d342f8b789342720ef8fdc57c03d6a8/jupyter/pytorch/ubi9-python-3.12/kustomize/base/service.yaml) |
| notebook | ClusterIP | 8888/TCP | [`jupyter/pytorch+llmcompressor/ubi9-python-3.12/kustomize/base/service.yaml`](https://github.com/red-hat-data-services/notebooks/blob/1fb86d8e3d342f8b789342720ef8fdc57c03d6a8/jupyter/pytorch+llmcompressor/ubi9-python-3.12/kustomize/base/service.yaml) |
| notebook | ClusterIP | 8888/TCP | [`jupyter/rocm/pytorch/ubi9-python-3.12/kustomize/base/service.yaml`](https://github.com/red-hat-data-services/notebooks/blob/1fb86d8e3d342f8b789342720ef8fdc57c03d6a8/jupyter/rocm/pytorch/ubi9-python-3.12/kustomize/base/service.yaml) |
| notebook | ClusterIP | 8888/TCP | [`jupyter/rocm/tensorflow/ubi9-python-3.12/kustomize/base/service.yaml`](https://github.com/red-hat-data-services/notebooks/blob/1fb86d8e3d342f8b789342720ef8fdc57c03d6a8/jupyter/rocm/tensorflow/ubi9-python-3.12/kustomize/base/service.yaml) |
| notebook | ClusterIP | 8888/TCP | [`jupyter/tensorflow/ubi9-python-3.12/kustomize/base/service.yaml`](https://github.com/red-hat-data-services/notebooks/blob/1fb86d8e3d342f8b789342720ef8fdc57c03d6a8/jupyter/tensorflow/ubi9-python-3.12/kustomize/base/service.yaml) |
| notebook | ClusterIP | 8888/TCP | [`jupyter/trustyai/ubi9-python-3.12/kustomize/base/service.yaml`](https://github.com/red-hat-data-services/notebooks/blob/1fb86d8e3d342f8b789342720ef8fdc57c03d6a8/jupyter/trustyai/ubi9-python-3.12/kustomize/base/service.yaml) |

!!! warning "No Network Policies"
    No NetworkPolicy resources found. All pod-to-pod traffic is allowed by default.

