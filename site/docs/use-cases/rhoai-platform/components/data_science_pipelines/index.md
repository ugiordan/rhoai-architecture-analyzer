# data-science-pipelines

> **Architecture snapshot: 2026-05-15** (2026-05-15)


**Repository:** kubeflow/data-science-pipelines  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-05-15T11:39:28Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 3 |
| Deployments | 11 |
| Services | 2 |
| Secrets | 2 |
| Cluster Roles | 14 |
| Controller Watches | 1 |

## Component Architecture

CRDs, controllers, and owned Kubernetes resources.

```mermaid
graph LR
    %% Component architecture for data-science-pipelines

    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef controller fill:#3498db,stroke:#2980b9,color:#fff
    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff

    subgraph controller["data-science-pipelines Controller"]
        dep_1["cache-server"]
        class dep_1 controller
        dep_2["kubeflow-pipelines-profile-controller"]
        class dep_2 controller
        dep_3["metadata-writer"]
        class dep_3 controller
        dep_4["ml-pipeline"]
        class dep_4 controller
        dep_5["ml-pipeline"]
        class dep_5 controller
        dep_6["ml-pipeline"]
        class dep_6 controller
        dep_7["ml-pipeline-persistenceagent"]
        class dep_7 controller
        dep_8["ml-pipeline-scheduledworkflow"]
        class dep_8 controller
        dep_9["ml-pipeline-ui"]
        class dep_9 controller
        dep_10["ml-pipeline-viewer-crd"]
        class dep_10 controller
        dep_11["squid"]
        class dep_11 controller
    end

    crd_CompositeController{{"CompositeController\nmetacontroller.k8s.io/v1alpha1"}}
    class crd_CompositeController crd
    crd_ControllerRevision{{"ControllerRevision\nmetacontroller.k8s.io/v1alpha1"}}
    class crd_ControllerRevision crd
    crd_DecoratorController{{"DecoratorController\nmetacontroller.k8s.io/v1alpha1"}}
    class crd_DecoratorController crd
```

### CRDs

| Group | Version | Kind | Scope | Fields | Validation Rules | Discovery | Source |
|-------|---------|------|-------|--------|------------------|-----------|--------|
| metacontroller.k8s.io | v1alpha1 | CompositeController | Cluster | 109 | 0 | YAML | [`manifests/kustomize/third-party/metacontroller/base/crd.yaml`](https://github.com/kubeflow/data-science-pipelines/blob/e61fa54e17eb9a52898792f7554ea3e00dc8eb0b/manifests/kustomize/third-party/metacontroller/base/crd.yaml) |
| metacontroller.k8s.io | v1alpha1 | ControllerRevision | Namespaced | 8 | 0 | YAML | [`manifests/kustomize/third-party/metacontroller/base/crd.yaml`](https://github.com/kubeflow/data-science-pipelines/blob/e61fa54e17eb9a52898792f7554ea3e00dc8eb0b/manifests/kustomize/third-party/metacontroller/base/crd.yaml) |
| metacontroller.k8s.io | v1alpha1 | DecoratorController | Cluster | 75 | 0 | YAML | [`manifests/kustomize/third-party/metacontroller/base/crd.yaml`](https://github.com/kubeflow/data-science-pipelines/blob/e61fa54e17eb9a52898792f7554ea3e00dc8eb0b/manifests/kustomize/third-party/metacontroller/base/crd.yaml) |

## Dependencies

### Key External Dependencies

| Module | Version |
|--------|---------|
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.1 |
| github.com/go-logr/logr | v1.2.2 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.1 |
| github.com/go-logr/logr | v1.2.2 |
| github.com/go-logr/logr | v1.3.0 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.3.0 |
| github.com/go-logr/stdr | v1.2.2 |
| github.com/go-logr/stdr | v1.2.2 |
| github.com/go-logr/stdr | v1.2.2 |
| github.com/go-logr/stdr | v1.2.2 |
| github.com/go-logr/zapr | v1.3.0 |
| github.com/go-logr/zapr | v1.3.0 |
| github.com/prometheus/client_golang | v1.23.2 |
| github.com/prometheus/client_golang | v1.23.2 |
| github.com/prometheus/client_golang | v1.22.0 |
| github.com/prometheus/client_golang | v1.22.0 |
| github.com/prometheus/client_golang | v1.23.2 |
| github.com/prometheus/client_golang | v1.22.0 |
| github.com/prometheus/client_golang | v1.22.0 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/common | v0.64.0 |
| github.com/prometheus/common | v0.64.0 |
| github.com/prometheus/common | v0.66.1 |
| github.com/prometheus/common | v0.64.0 |
| github.com/prometheus/common | v0.64.0 |
| github.com/prometheus/common | v0.66.1 |
| github.com/prometheus/procfs | v0.16.1 |
| github.com/prometheus/procfs | v0.16.1 |
| google.golang.org/grpc | v1.72.2 |
| google.golang.org/grpc | v1.65.0 |
| google.golang.org/grpc | v1.71.1 |
| google.golang.org/grpc | v1.71.0 |
| google.golang.org/grpc | v1.74.2 |
| google.golang.org/grpc | v1.71.1 |
| google.golang.org/grpc | v1.72.1 |
| google.golang.org/grpc | v1.58.2 |
| google.golang.org/grpc | v1.71.1 |
| google.golang.org/grpc | v1.43.0 |
| google.golang.org/grpc | v1.72.2 |
| google.golang.org/grpc | v1.72.2 |
| google.golang.org/grpc | v1.72.1 |
| google.golang.org/grpc | v1.65.0 |
| google.golang.org/grpc | v1.72.1 |
| google.golang.org/grpc | v1.72.1 |
| google.golang.org/grpc | v1.72.1 |
| google.golang.org/grpc | v1.72.0 |
| google.golang.org/grpc | v1.72.2 |
| google.golang.org/grpc | v1.43.0 |
| google.golang.org/grpc | v1.72.1 |
| google.golang.org/grpc | v1.33.2 |
| google.golang.org/grpc | v1.68.0 |
| google.golang.org/grpc | v1.33.1 |
| google.golang.org/grpc | v1.72.1 |
| google.golang.org/grpc | v1.75.0 |
| google.golang.org/grpc | v1.71.1 |
| google.golang.org/grpc | v1.72.0 |
| google.golang.org/grpc | v1.75.1 |
| google.golang.org/grpc | v1.75.1 |
| google.golang.org/grpc | v1.73.0 |
| google.golang.org/grpc | v1.71.1 |
| google.golang.org/grpc | v1.72.1 |
| google.golang.org/grpc | v1.63.2 |
| google.golang.org/grpc | v1.72.2 |
| google.golang.org/grpc | v1.71.0 |
| google.golang.org/grpc | v1.75.0 |
| google.golang.org/grpc | v1.33.2 |
| google.golang.org/grpc | v1.72.2 |
| google.golang.org/grpc | v1.68.0 |
| google.golang.org/grpc | v1.72.2 |
| google.golang.org/grpc | v1.65.0 |
| google.golang.org/grpc | v1.65.0 |
| google.golang.org/grpc | v1.58.2 |
| google.golang.org/grpc | v1.74.2 |
| google.golang.org/grpc | v1.56.3 |
| google.golang.org/grpc | v1.63.2 |
| google.golang.org/grpc | v1.71.1 |
| google.golang.org/grpc | v1.72.1 |
| google.golang.org/grpc | v1.79.3 |
| google.golang.org/grpc | v1.67.1 |
| google.golang.org/grpc | v1.72.0 |
| google.golang.org/grpc | v1.67.1 |
| google.golang.org/grpc | v1.72.1 |
| google.golang.org/grpc | v1.33.1 |
| google.golang.org/grpc | v1.73.0 |
| google.golang.org/grpc | v1.72.0 |
| google.golang.org/grpc | v1.72.2 |
| google.golang.org/grpc/cmd/protoc-gen-go-grpc | v1.5.1 |
| google.golang.org/grpc/examples | v0.0.0-20250407062114-b368379ef8f6 |
| google.golang.org/grpc/examples | v0.0.0-20250407062114-b368379ef8f6 |
| k8s.io/api | v0.35.2 |
| k8s.io/api | v0.35.2 |
| k8s.io/api | v0.33.1 |
| k8s.io/api | v0.35.2 |
| k8s.io/api | v0.35.2 |
| k8s.io/api | v0.35.0 |
| k8s.io/api | v0.35.0 |
| k8s.io/api | v0.33.1 |
| k8s.io/api | v0.35.2 |
| k8s.io/apiextensions-apiserver | v0.35.0 |
| k8s.io/apiextensions-apiserver | v0.35.0 |
| k8s.io/apimachinery | v0.35.2 |
| k8s.io/apimachinery | v0.35.2 |
| k8s.io/apimachinery | v0.35.2 |
| k8s.io/apimachinery | v0.33.1 |
| k8s.io/apimachinery | v0.33.1 |
| k8s.io/apimachinery | v0.35.2 |
| k8s.io/apimachinery | v0.35.2 |
| k8s.io/apimachinery | v0.35.0 |
| k8s.io/apimachinery | v0.35.0 |
| k8s.io/apimachinery | v0.35.2 |
| k8s.io/apimachinery | v0.35.2 |
| k8s.io/apimachinery | v0.35.2 |
| k8s.io/apimachinery | v0.35.2 |
| k8s.io/apiserver | v0.35.0 |
| k8s.io/apiserver | v0.35.2 |
| k8s.io/apiserver | v0.35.2 |
| k8s.io/apiserver | v0.35.0 |
| k8s.io/client-go | v0.35.0 |
| k8s.io/client-go | v0.35.2 |
| k8s.io/client-go | v0.35.0 |
| k8s.io/client-go | v0.33.1 |
| k8s.io/client-go | v0.35.2 |
| k8s.io/client-go | v0.33.1 |
| k8s.io/client-go | v0.35.2 |
| sigs.k8s.io/controller-runtime | v0.23.3 |

