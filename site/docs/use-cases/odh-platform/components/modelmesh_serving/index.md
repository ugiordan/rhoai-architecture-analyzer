# modelmesh-serving

> **Architecture snapshot: 2026-05-15** (2026-05-15)


**Repository:** kserve/modelmesh-serving  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-05-15T11:38:02Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 4 |
| Deployments | 13 |
| Services | 9 |
| Secrets | 2 |
| Cluster Roles | 0 |
| Controller Watches | 48 |

## Component Architecture

CRDs, controllers, and owned Kubernetes resources.

```mermaid
graph LR
    %% Component architecture for modelmesh-serving

    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef controller fill:#3498db,stroke:#2980b9,color:#fff
    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff

    subgraph controller["modelmesh-serving Controller"]
        dep_1["controller-manager"]
        class dep_1 controller
        dep_2["etcd"]
        class dep_2 controller
        dep_3["kserve-controller-manager"]
        class dep_3 controller
        dep_4["kserve-controller-manager"]
        class dep_4 controller
        dep_5["kserve-controller-manager"]
        class dep_5 controller
        dep_6["kserve-controller-manager"]
        class dep_6 controller
        dep_7["kserve-controller-manager"]
        class dep_7 controller
        dep_8["kserve-controller-manager"]
        class dep_8 controller
        dep_9["kserve-controller-manager"]
        class dep_9 controller
        dep_10["kserve-controller-manager"]
        class dep_10 controller
        dep_11["kserve-controller-manager"]
        class dep_11 controller
        dep_12["kserve-controller-manager"]
        class dep_12 controller
        dep_13["modelmesh-controller"]
        class dep_13 controller
    end

    crd_ClusterServingRuntime{{"ClusterServingRuntime\nserving.kserve.io/v1alpha1"}}
    class crd_ClusterServingRuntime crd
    crd_Predictor{{"Predictor\nserving.kserve.io/v1alpha1"}}
    class crd_Predictor crd
    crd_Predictor -->|"For (reconciles)"| controller
    crd_ServingRuntime{{"ServingRuntime\nserving.kserve.io/v1alpha1"}}
    class crd_ServingRuntime crd
    crd_ServingRuntime -->|"For (reconciles)"| controller
    crd_InferenceService{{"InferenceService\nserving.kserve.io/v1beta1"}}
    class crd_InferenceService crd
    crd_InferenceService -->|"For (reconciles)"| controller
    controller -->|"Owns"| owned_14["Deployment"]
    class owned_14 owned
    controller -->|"Owns"| owned_15["Service"]
    class owned_15 owned
    controller -->|"Owns"| owned_16["VirtualService"]
    class owned_16 owned
    watch_17["ClusterServingRuntime"] -->|"Watches"| controller
    class watch_17 external
    watch_18["ConfigMap"] -->|"Watches"| controller
    class watch_18 external
    watch_19["InferenceService"] -->|"Watches"| controller
    class watch_19 external
    watch_20["Kind"] -->|"Watches"| controller
    class watch_20 external
    watch_21["Namespace"] -->|"Watches"| controller
    class watch_21 external
    watch_22["Predictor"] -->|"Watches"| controller
    class watch_22 external
    watch_23["Secret"] -->|"Watches"| controller
    class watch_23 external
    watch_24["ServiceMonitor"] -->|"Watches"| controller
    class watch_24 external
```

### CRDs

| Group | Version | Kind | Scope | Fields | Validation Rules | Discovery | Source |
|-------|---------|------|-------|--------|------------------|-----------|--------|
| serving.kserve.io | v1alpha1 | ClusterServingRuntime | Cluster | 559 | 0 | YAML | [`config/crd/bases/serving.kserve.io_clusterservingruntimes.yaml`](https://github.com/kserve/modelmesh-serving/blob/1fcf541d867ceb459fbc76aa1e2bef102c4816db/config/crd/bases/serving.kserve.io_clusterservingruntimes.yaml) |
| serving.kserve.io | v1alpha1 | Predictor | Namespaced | 40 | 0 | YAML + Go AST | [`config/crd/bases/serving.kserve.io_predictors.yaml`](https://github.com/kserve/modelmesh-serving/blob/1fcf541d867ceb459fbc76aa1e2bef102c4816db/config/crd/bases/serving.kserve.io_predictors.yaml) |
| serving.kserve.io | v1alpha1 | ServingRuntime | Namespaced | 1140 | 0 | YAML | [`config/crd/bases/serving.kserve.io_servingruntimes.yaml`](https://github.com/kserve/modelmesh-serving/blob/1fcf541d867ceb459fbc76aa1e2bef102c4816db/config/crd/bases/serving.kserve.io_servingruntimes.yaml) |
| serving.kserve.io | v1beta1 | InferenceService | Namespaced | 6195 | 0 | YAML | [`config/crd/bases/serving.kserve.io_inferenceservices.yaml`](https://github.com/kserve/modelmesh-serving/blob/1fcf541d867ceb459fbc76aa1e2bef102c4816db/config/crd/bases/serving.kserve.io_inferenceservices.yaml) |

## Dependencies

### Key External Dependencies

| Module | Version |
|--------|---------|
| github.com/go-logr/logr | v1.2.4 |
| github.com/go-logr/logr | v1.3.0 |
| github.com/go-logr/logr | v1.2.2 |
| github.com/go-logr/logr | v1.2.0 |
| github.com/go-logr/logr | v1.3.0 |
| github.com/go-logr/logr | v1.4.1 |
| github.com/go-logr/logr | v1.3.0 |
| github.com/go-logr/logr | v1.2.4 |
| github.com/go-logr/logr | v1.4.1 |
| github.com/go-logr/logr | v1.2.2 |
| github.com/go-logr/logr | v1.3.0 |
| github.com/go-logr/logr | v1.3.0 |
| github.com/go-logr/logr | v1.3.0 |
| github.com/go-logr/logr | v1.4.1 |
| github.com/go-logr/logr | v1.2.4 |
| github.com/go-logr/logr | v1.2.0 |
| github.com/go-logr/logr | v1.2.4 |
| github.com/go-logr/zapr | v1.2.3 |
| github.com/go-logr/zapr | v1.2.4 |
| github.com/go-logr/zapr | v1.2.3 |
| github.com/go-logr/zapr | v1.2.4 |
| github.com/operator-framework/api | v0.10.0 |
| github.com/operator-framework/api | v0.10.0 |
| github.com/operator-framework/operator-lib | v0.10.0 |
| github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring | v0.55.0 |
| github.com/prometheus/client_golang | v1.17.0 |
| github.com/prometheus/client_golang | v1.11.1 |
| github.com/prometheus/client_golang | v1.17.0 |
| github.com/prometheus/client_golang | v1.16.0 |
| github.com/prometheus/client_golang | v1.11.1 |
| github.com/prometheus/client_golang | v1.11.0 |
| github.com/prometheus/client_golang | v1.11.0 |
| github.com/prometheus/client_golang | v1.17.0 |
| github.com/prometheus/client_golang | v1.17.0 |
| github.com/prometheus/client_golang | v1.16.0 |
| github.com/prometheus/client_golang | v1.16.0 |
| github.com/prometheus/client_golang | v1.16.0 |
| github.com/prometheus/client_model | v0.4.0 |
| github.com/prometheus/client_model | v0.4.1-0.20230718164431-9a2bf3000d16 |
| github.com/prometheus/client_model | v0.4.1-0.20230718164431-9a2bf3000d16 |
| github.com/prometheus/client_model | v0.4.0 |
| github.com/prometheus/client_model | v0.4.0 |
| github.com/prometheus/client_model | v0.4.1-0.20230718164431-9a2bf3000d16 |
| github.com/prometheus/client_model | v0.4.1-0.20230718164431-9a2bf3000d16 |
| github.com/prometheus/client_model | v0.2.0 |
| github.com/prometheus/client_model | v0.4.0 |
| github.com/prometheus/client_model | v0.2.0 |
| github.com/prometheus/common | v0.44.0 |
| github.com/prometheus/common | v0.45.0 |
| github.com/prometheus/common | v0.44.0 |
| github.com/prometheus/common | v0.44.0 |
| github.com/prometheus/common | v0.44.0 |
| github.com/prometheus/common | v0.45.0 |
| github.com/prometheus/procfs | v0.11.1 |
| github.com/prometheus/procfs | v0.10.1 |
| github.com/prometheus/procfs | v0.11.1 |
| github.com/prometheus/procfs | v0.10.1 |
| google.golang.org/grpc | v1.41.0 |
| google.golang.org/grpc | v1.33.2 |
| google.golang.org/grpc | v1.58.3 |
| google.golang.org/grpc | v1.59.0 |
| google.golang.org/grpc | v1.59.0 |
| google.golang.org/grpc | v1.41.0 |
| google.golang.org/grpc | v1.59.0 |
| google.golang.org/grpc | v1.33.2 |
| google.golang.org/grpc | v1.59.0 |
| google.golang.org/grpc | v1.56.3 |
| google.golang.org/grpc | v1.59.0 |
| google.golang.org/grpc | v1.56.3 |
| google.golang.org/grpc | v1.59.0 |
| google.golang.org/grpc | v1.56.1 |
| google.golang.org/grpc | v1.57.0 |
| google.golang.org/grpc | v1.59.0 |
| google.golang.org/grpc | v1.59.0 |
| google.golang.org/grpc | v1.59.0 |
| google.golang.org/grpc | v1.59.0 |
| google.golang.org/grpc | v1.59.0 |
| google.golang.org/grpc | v1.59.0 |
| google.golang.org/grpc | v1.59.0 |
| google.golang.org/grpc | v1.59.0 |
| google.golang.org/grpc | v1.41.0 |
| google.golang.org/grpc | v1.59.0 |
| google.golang.org/grpc | v1.58.3 |
| google.golang.org/grpc | v1.56.1 |
| google.golang.org/grpc | v1.41.0 |
| google.golang.org/grpc | v1.59.0 |
| google.golang.org/grpc | v1.59.0 |
| google.golang.org/grpc | v1.57.0 |
| google.golang.org/grpc | v1.59.0 |
| google.golang.org/grpc | v1.59.0 |
| k8s.io/api | v0.28.4 |
| k8s.io/api | v0.27.6 |
| k8s.io/api | v0.22.5 |
| k8s.io/api | v0.22.5 |
| k8s.io/api | v0.28.4 |
| k8s.io/api | v0.23.0 |
| k8s.io/api | v0.28.4 |
| k8s.io/api | v0.28.3 |
| k8s.io/api | v0.28.4 |
| k8s.io/api | v0.27.6 |
| k8s.io/api | v0.28.4 |
| k8s.io/api | v0.23.0 |
| k8s.io/api | v0.27.6 |
| k8s.io/api | v0.27.6 |
| k8s.io/api | v0.27.6 |
| k8s.io/api | v0.23.0 |
| k8s.io/api | v0.28.3 |
| k8s.io/api | v0.28.4 |
| k8s.io/api | v0.28.4 |
| k8s.io/api | v0.27.6 |
| k8s.io/api | v0.23.0 |
| k8s.io/apiextensions-apiserver | v0.27.6 |
| k8s.io/apiextensions-apiserver | v0.28.3 |
| k8s.io/apiextensions-apiserver | v0.23.0 |
| k8s.io/apiextensions-apiserver | v0.27.6 |
| k8s.io/apiextensions-apiserver | v0.28.3 |
| k8s.io/apiextensions-apiserver | v0.27.6 |
| k8s.io/apiextensions-apiserver | v0.27.6 |
| k8s.io/apiextensions-apiserver | v0.23.0 |
| k8s.io/apimachinery | v0.27.6 |
| k8s.io/apimachinery | v0.28.4 |
| k8s.io/apimachinery | v0.28.4 |
| k8s.io/apimachinery | v0.23.0 |
| k8s.io/apimachinery | v0.19.7 |
| k8s.io/apimachinery | v0.28.4 |
| k8s.io/apimachinery | v0.28.3 |
| k8s.io/apimachinery | v0.22.5 |
| k8s.io/apimachinery | v0.27.6 |
| k8s.io/apimachinery | v0.28.4 |
| k8s.io/apimachinery | v0.28.4 |
| k8s.io/apimachinery | v0.27.6 |
| k8s.io/apimachinery | v0.27.6 |
| k8s.io/apimachinery | v0.19.7 |
| k8s.io/apimachinery | v0.22.5 |
| k8s.io/apimachinery | v0.27.6 |
| k8s.io/apimachinery | v0.27.6 |
| k8s.io/apimachinery | v0.28.4 |
| k8s.io/apimachinery | v0.23.0 |
| k8s.io/apimachinery | v0.23.0 |
| k8s.io/apimachinery | v0.28.4 |
| k8s.io/apimachinery | v0.28.3 |
| k8s.io/apimachinery | v0.30.13 |
| k8s.io/apimachinery | v0.23.0 |
| k8s.io/apimachinery | v0.28.4 |
| k8s.io/apimachinery | v0.28.4 |
| k8s.io/apimachinery | v0.28.4 |
| k8s.io/apiserver | v0.28.3 |
| k8s.io/apiserver | v0.28.4 |
| k8s.io/apiserver | v0.28.3 |
| k8s.io/apiserver | v0.28.4 |
| k8s.io/client-go | v0.27.6 |
| k8s.io/client-go | v0.27.6 |
| k8s.io/client-go | v0.28.4 |
| k8s.io/client-go | v0.22.5 |
| k8s.io/client-go | v0.27.6 |
| k8s.io/client-go | v0.27.6 |
| k8s.io/client-go | v0.28.4 |
| k8s.io/client-go | v0.28.3 |
| k8s.io/client-go | v0.27.6 |
| k8s.io/client-go | v0.23.0 |
| k8s.io/client-go | v0.23.0 |
| k8s.io/client-go | v0.28.4 |
| k8s.io/client-go | v0.28.4 |
| k8s.io/client-go | v0.27.6 |
| k8s.io/client-go | v0.28.4 |
| k8s.io/client-go | v0.28.4 |
| k8s.io/client-go | v0.28.3 |
| k8s.io/client-go | v0.28.4 |
| k8s.io/client-go | v0.22.5 |
| sigs.k8s.io/controller-runtime | v0.16.3 |
| sigs.k8s.io/controller-runtime | v0.7.2 |
| sigs.k8s.io/controller-runtime | v0.16.3 |
| sigs.k8s.io/controller-runtime | v0.11.0 |
| sigs.k8s.io/controller-runtime | v0.11.0 |
| sigs.k8s.io/controller-runtime | v0.16.3 |
| sigs.k8s.io/controller-runtime | v0.7.2 |

