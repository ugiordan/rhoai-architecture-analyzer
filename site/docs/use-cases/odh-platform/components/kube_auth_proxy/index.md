# kube-auth-proxy

> **Architecture snapshot: 2026-05-16** (2026-05-16)


**Repository:** opendatahub-io/kube-auth-proxy  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-05-16T03:43:04Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 0 |
| Deployments | 0 |
| Services | 0 |
| Secrets | 0 |
| Cluster Roles | 0 |
| Controller Watches | 0 |

## Component Architecture

CRDs, controllers, and owned Kubernetes resources.

```mermaid
graph LR
    %% Component architecture for kube-auth-proxy

    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef controller fill:#3498db,stroke:#2980b9,color:#fff
    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff

    subgraph controller["kube-auth-proxy Controller"]
        ctrl_1["Controller"]
        class ctrl_1 controller
    end
```

### CRDs

No CRDs defined.

## Dependencies

### Key External Dependencies

| Module | Version |
|--------|---------|
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.1 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.2.2 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.3.0 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.1 |
| github.com/go-logr/logr | v1.2.2 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.3.0 |
| github.com/go-logr/stdr | v1.2.2 |
| github.com/go-logr/stdr | v1.2.2 |
| github.com/go-logr/stdr | v1.2.2 |
| github.com/go-logr/stdr | v1.2.2 |
| github.com/go-logr/zapr | v1.3.0 |
| github.com/go-logr/zapr | v1.3.0 |
| github.com/prometheus/client_golang | v1.11.1 |
| github.com/prometheus/client_golang | v1.20.5 |
| github.com/prometheus/client_golang | v1.20.5 |
| github.com/prometheus/client_golang | v1.14.0 |
| github.com/prometheus/client_golang | v1.23.2 |
| github.com/prometheus/client_golang | v1.23.2 |
| github.com/prometheus/client_golang | v1.20.5 |
| github.com/prometheus/client_golang | v1.11.1 |
| github.com/prometheus/client_golang | v1.14.0 |
| github.com/prometheus/client_golang | v1.20.5 |
| github.com/prometheus/client_golang | v1.23.2 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.3.0 |
| github.com/prometheus/client_model | v0.6.1 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.6.1 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.3.0 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/common | v0.66.1 |
| github.com/prometheus/common | v0.66.1 |
| github.com/prometheus/common | v0.66.1 |
| github.com/prometheus/common | v0.66.1 |
| github.com/prometheus/procfs | v0.16.1 |
| github.com/prometheus/procfs | v0.16.1 |
| github.com/prometheus/procfs | v0.16.1 |
| github.com/prometheus/procfs | v0.16.1 |
| google.golang.org/grpc | v1.71.1 |
| google.golang.org/grpc | v1.75.1 |
| google.golang.org/grpc | v1.70.0 |
| google.golang.org/grpc | v1.69.4 |
| google.golang.org/grpc | v1.75.1 |
| google.golang.org/grpc | v1.72.2 |
| google.golang.org/grpc | v1.69.2 |
| google.golang.org/grpc | v1.68.0 |
| google.golang.org/grpc | v1.72.2 |
| google.golang.org/grpc | v1.75.0 |
| google.golang.org/grpc | v1.75.0 |
| google.golang.org/grpc | v1.75.1 |
| google.golang.org/grpc | v1.56.3 |
| google.golang.org/grpc | v1.71.0 |
| google.golang.org/grpc | v1.71.1 |
| google.golang.org/grpc | v1.67.1 |
| google.golang.org/grpc | v1.71.1 |
| google.golang.org/grpc | v1.71.1 |
| google.golang.org/grpc | v1.67.1 |
| google.golang.org/grpc | v1.72.2 |
| google.golang.org/grpc | v1.69.4 |
| google.golang.org/grpc | v1.56.3 |
| google.golang.org/grpc | v1.70.0 |
| google.golang.org/grpc | v1.67.1 |
| google.golang.org/grpc | v1.68.0 |
| google.golang.org/grpc | v1.71.0 |
| google.golang.org/grpc | v1.75.1 |
| google.golang.org/grpc | v1.71.0 |
| google.golang.org/grpc | v1.72.2 |
| google.golang.org/grpc | v1.71.1 |
| google.golang.org/grpc | v1.67.1 |
| google.golang.org/grpc | v1.63.2 |
| google.golang.org/grpc | v1.63.2 |
| google.golang.org/grpc | v1.69.2 |
| google.golang.org/grpc | v1.71.0 |
| google.golang.org/grpc | v1.58.2 |
| google.golang.org/grpc | v1.58.2 |
| google.golang.org/grpc | v1.71.1 |
| google.golang.org/grpc | v1.71.1 |
| google.golang.org/grpc | v1.71.1 |
| google.golang.org/grpc/examples | v0.0.0-20250407062114-b368379ef8f6 |
| google.golang.org/grpc/examples | v0.0.0-20250407062114-b368379ef8f6 |
| k8s.io/api | v0.35.3 |
| k8s.io/api | v0.35.3 |
| k8s.io/api | v0.35.3 |
| k8s.io/api | v0.35.3 |
| k8s.io/api | v0.35.3 |
| k8s.io/api | v0.35.3 |
| k8s.io/apimachinery | v0.35.3 |
| k8s.io/apimachinery | v0.35.3 |
| k8s.io/apimachinery | v0.35.3 |
| k8s.io/apimachinery | v0.35.3 |
| k8s.io/apimachinery | v0.35.3 |
| k8s.io/apimachinery | v0.35.3 |
| k8s.io/apimachinery | v0.35.3 |
| k8s.io/apimachinery | v0.35.3 |
| k8s.io/apimachinery | v0.35.3 |
| k8s.io/apimachinery | v0.35.3 |
| k8s.io/apiserver | v0.35.3 |
| k8s.io/client-go | v0.35.3 |
| k8s.io/client-go | v0.35.3 |
| k8s.io/client-go | v0.35.3 |
| k8s.io/client-go | v0.35.3 |
| k8s.io/client-go | v0.35.3 |
| k8s.io/client-go | v0.35.3 |

