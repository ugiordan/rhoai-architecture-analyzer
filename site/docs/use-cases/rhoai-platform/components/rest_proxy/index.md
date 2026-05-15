# rest-proxy

> **Architecture snapshot: 2026-05-15** (2026-05-15)


**Repository:** kserve/rest-proxy  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-05-15T11:48:39Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 0 |
| Deployments | 0 |
| Services | 0 |
| Secrets | 0 |
| Cluster Roles | 0 |
| Controller Watches | 4 |

## Component Architecture

CRDs, controllers, and owned Kubernetes resources.

```mermaid
graph LR
    %% Component architecture for rest-proxy

    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef controller fill:#3498db,stroke:#2980b9,color:#fff
    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff

    subgraph controller["rest-proxy Controller"]
        ctrl_1["Controller"]
        class ctrl_1 controller
    end

    watch_2["Kind"] -->|"Watches"| controller
    class watch_2 external
```

### CRDs

No CRDs defined.

## Dependencies

### Key External Dependencies

| Module | Version |
|--------|---------|
| github.com/go-logr/logr | v1.4.1 |
| github.com/go-logr/logr | v1.2.3 |
| github.com/go-logr/logr | v1.2.2 |
| github.com/go-logr/logr | v1.2.2 |
| github.com/go-logr/logr | v1.4.1 |
| github.com/go-logr/logr | v1.2.3 |
| github.com/go-logr/zapr | v1.2.3 |
| github.com/go-logr/zapr | v1.2.3 |
| github.com/prometheus/client_golang | v1.14.0 |
| github.com/prometheus/client_golang | v1.14.0 |
| github.com/prometheus/client_model | v0.3.0 |
| github.com/prometheus/client_model | v0.3.0 |
| google.golang.org/grpc | v1.54.0 |
| google.golang.org/grpc | v1.56.3 |
| google.golang.org/grpc | v1.51.0 |
| google.golang.org/grpc | v1.54.0 |
| google.golang.org/grpc | v1.51.0 |
| k8s.io/api | v0.26.0 |
| k8s.io/api | v0.26.0 |
| k8s.io/apiextensions-apiserver | v0.26.0 |
| k8s.io/apiextensions-apiserver | v0.26.0 |
| k8s.io/apimachinery | v0.26.0 |
| k8s.io/apimachinery | v0.26.0 |
| k8s.io/client-go | v0.26.0 |
| k8s.io/client-go | v0.26.0 |
| sigs.k8s.io/controller-runtime | v0.14.1 |

