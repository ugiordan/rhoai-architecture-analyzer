# ai-gateway-payload-processing

> **Architecture snapshot: 2026-05-18** (2026-05-18)


**Repository:** opendatahub-io/ai-gateway-payload-processing  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-05-18T04:24:55Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 2 |
| Deployments | 0 |
| Services | 1 |
| Secrets | 0 |
| Cluster Roles | 0 |
| Controller Watches | 17 |

## Component Architecture

CRDs, controllers, and owned Kubernetes resources.

```mermaid
graph LR
    %% Component architecture for ai-gateway-payload-processing

    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef controller fill:#3498db,stroke:#2980b9,color:#fff
    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff

    subgraph controller["ai-gateway-payload-processing Controller"]
        ctrl_1["Controller"]
        class ctrl_1 controller
    end

    crd_ExternalModel{{"ExternalModel\ninference.opendatahub.io/v1alpha1"}}
    class crd_ExternalModel crd
    crd_ExternalModel -->|"For (reconciles)"| controller
    crd_ExternalProvider{{"ExternalProvider\ninference.opendatahub.io/v1alpha1"}}
    class crd_ExternalProvider crd
    crd_ExternalProvider -->|"For (reconciles)"| controller
    controller -->|"Owns"| owned_2["HTTPRoute"]
    class owned_2 owned
    controller -->|"Owns"| owned_3["Service"]
    class owned_3 owned
    watch_4["ExternalProvider"] -->|"Watches"| controller
    class watch_4 external
```

### CRDs

| Group | Version | Kind | Scope | Fields | Validation Rules | Discovery | Source |
|-------|---------|------|-------|--------|------------------|-----------|--------|
| inference.opendatahub.io | v1alpha1 | ExternalModel | Namespaced | 18 | 0 | YAML | [`config/crd/bases/inference.opendatahub.io_externalmodels.yaml`](https://github.com/opendatahub-io/ai-gateway-payload-processing/blob/d873739d504086159cbe9a1bf0c410fbd908196b/config/crd/bases/inference.opendatahub.io_externalmodels.yaml) |
| inference.opendatahub.io | v1alpha1 | ExternalProvider | Namespaced | 19 | 0 | YAML | [`config/crd/bases/inference.opendatahub.io_externalproviders.yaml`](https://github.com/opendatahub-io/ai-gateway-payload-processing/blob/d873739d504086159cbe9a1bf0c410fbd908196b/config/crd/bases/inference.opendatahub.io_externalproviders.yaml) |

## Dependencies

### Key External Dependencies

| Module | Version |
|--------|---------|
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.1 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.3.0 |
| github.com/go-logr/logr | v1.3.0 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.1 |
| github.com/go-logr/logr | v1.2.2 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.2.2 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/stdr | v1.2.2 |
| github.com/go-logr/stdr | v1.2.2 |
| github.com/go-logr/stdr | v1.2.2 |
| github.com/go-logr/stdr | v1.2.2 |
| github.com/go-logr/zapr | v1.3.0 |
| github.com/go-logr/zapr | v1.3.0 |
| github.com/go-logr/zapr | v1.3.0 |
| github.com/go-logr/zapr | v1.3.0 |
| github.com/go-logr/zapr | v1.3.0 |
| github.com/go-logr/zapr | v1.3.0 |
| github.com/prometheus/client_golang | v1.23.2 |
| github.com/prometheus/client_golang | v1.11.1 |
| github.com/prometheus/client_golang | v1.23.2 |
| github.com/prometheus/client_golang | v1.11.1 |
| github.com/prometheus/client_golang | v1.23.2 |
| github.com/prometheus/client_golang | v1.23.2 |
| github.com/prometheus/client_golang | v1.23.2 |
| github.com/prometheus/client_golang | v1.23.2 |
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
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/common | v0.67.5 |
| github.com/prometheus/common | v0.67.5 |
| github.com/prometheus/common | v0.66.1 |
| github.com/prometheus/common | v0.66.1 |
| github.com/prometheus/common | v0.66.1 |
| github.com/prometheus/common | v0.66.1 |
| github.com/prometheus/procfs | v0.16.1 |
| github.com/prometheus/procfs | v0.16.1 |
| github.com/prometheus/procfs | v0.16.1 |
| github.com/prometheus/procfs | v0.16.1 |
| github.com/prometheus/prometheus | v0.310.0 |
| github.com/prometheus/prometheus | v0.310.0 |
| google.golang.org/grpc | v1.80.0 |
| google.golang.org/grpc | v1.79.2 |
| google.golang.org/grpc | v1.79.1 |
| google.golang.org/grpc | v1.56.3 |
| google.golang.org/grpc | v1.72.2 |
| google.golang.org/grpc | v1.80.0 |
| google.golang.org/grpc | v1.56.3 |
| google.golang.org/grpc | v1.78.0 |
| google.golang.org/grpc | v1.79.3 |
| google.golang.org/grpc | v1.72.2 |
| google.golang.org/grpc | v1.72.2 |
| google.golang.org/grpc | v1.79.3 |
| google.golang.org/grpc | v1.68.0 |
| google.golang.org/grpc | v1.78.0 |
| google.golang.org/grpc | v1.79.2 |
| google.golang.org/grpc | v1.79.1 |
| google.golang.org/grpc | v1.58.2 |
| google.golang.org/grpc | v1.72.2 |
| google.golang.org/grpc | v1.58.2 |
| google.golang.org/grpc | v1.68.0 |
| google.golang.org/grpc | v1.80.0 |
| google.golang.org/grpc | v1.80.0 |
| k8s.io/api | v0.35.4 |
| k8s.io/api | v0.35.4 |
| k8s.io/api | v0.35.5 |
| k8s.io/api | v0.35.4 |
| k8s.io/api | v0.35.5 |
| k8s.io/api | v0.35.5 |
| k8s.io/api | v0.35.4 |
| k8s.io/api | v0.35.5 |
| k8s.io/api | v0.35.0 |
| k8s.io/api | v0.35.5 |
| k8s.io/api | v0.35.4 |
| k8s.io/api | v0.35.4 |
| k8s.io/api | v0.35.0 |
| k8s.io/apiextensions-apiserver | v0.35.4 |
| k8s.io/apiextensions-apiserver | v0.35.4 |
| k8s.io/apiextensions-apiserver | v0.35.0 |
| k8s.io/apiextensions-apiserver | v0.35.0 |
| k8s.io/apimachinery | v0.35.4 |
| k8s.io/apimachinery | v0.35.5 |
| k8s.io/apimachinery | v0.35.4 |
| k8s.io/apimachinery | v0.35.1 |
| k8s.io/apimachinery | v0.35.5 |
| k8s.io/apimachinery | v0.35.5 |
| k8s.io/apimachinery | v0.35.5 |
| k8s.io/apimachinery | v0.35.4 |
| k8s.io/apimachinery | v0.35.4 |
| k8s.io/apimachinery | v0.35.4 |
| k8s.io/apimachinery | v0.35.5 |
| k8s.io/apimachinery | v0.35.5 |
| k8s.io/apimachinery | v0.35.0 |
| k8s.io/apimachinery | v0.35.5 |
| k8s.io/apimachinery | v0.35.5 |
| k8s.io/apimachinery | v0.35.5 |
| k8s.io/apimachinery | v0.35.4 |
| k8s.io/apimachinery | v0.35.1 |
| k8s.io/apimachinery | v0.35.0 |
| k8s.io/apiserver | v0.35.0 |
| k8s.io/apiserver | v0.35.4 |
| k8s.io/apiserver | v0.35.0 |
| k8s.io/apiserver | v0.35.4 |
| k8s.io/client-go | v0.35.1 |
| k8s.io/client-go | v0.35.5 |
| k8s.io/client-go | v0.35.1 |
| k8s.io/client-go | v0.35.4 |
| k8s.io/client-go | v0.35.0 |
| k8s.io/client-go | v0.35.5 |
| k8s.io/client-go | v0.35.4 |
| k8s.io/client-go | v0.35.5 |
| k8s.io/client-go | v0.35.4 |
| k8s.io/client-go | v0.35.4 |
| k8s.io/client-go | v0.35.4 |
| k8s.io/client-go | v0.35.5 |
| k8s.io/client-go | v0.35.5 |
| k8s.io/client-go | v0.35.0 |
| k8s.io/client-go | v0.35.4 |
| sigs.k8s.io/controller-runtime | v0.23.3 |
| sigs.k8s.io/controller-runtime | v0.23.3 |
| sigs.k8s.io/controller-runtime | v0.23.3 |

