# modelmesh-serving

> **Architecture snapshot: 2026-05-05** (2026-05-05)


**Repository:** kserve/modelmesh-serving  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-05-05T15:09:04Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 4 |
| Deployments | 3 |
| Services | 3 |
| Secrets | 1 |
| Cluster Roles | 0 |
| Controller Watches | 17 |

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
        dep_3["modelmesh-controller"]
        class dep_3 controller
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
    controller -->|"Owns"| owned_4["Deployment"]
    class owned_4 owned
    controller -->|"Owns"| owned_5["Service"]
    class owned_5 owned
    watch_6["ClusterServingRuntime"] -->|"Watches"| controller
    class watch_6 external
    watch_7["ConfigMap"] -->|"Watches"| controller
    class watch_7 external
    watch_8["InferenceService"] -->|"Watches"| controller
    class watch_8 external
    watch_9["Namespace"] -->|"Watches"| controller
    class watch_9 external
    watch_10["Predictor"] -->|"Watches"| controller
    class watch_10 external
    watch_11["Secret"] -->|"Watches"| controller
    class watch_11 external
    watch_12["ServiceMonitor"] -->|"Watches"| controller
    class watch_12 external
```

### CRDs

| Group | Version | Kind | Scope | Fields | Validation Rules | Source |
|-------|---------|------|-------|--------|------------------|--------|
| serving.kserve.io | v1alpha1 | ClusterServingRuntime | Cluster | 559 | 0 | [`config/crd/bases/serving.kserve.io_clusterservingruntimes.yaml`](https://github.com/kserve/modelmesh-serving/blob/056bc2e855779c02536db9ef786b26cc73c63f20/config/crd/bases/serving.kserve.io_clusterservingruntimes.yaml) |
| serving.kserve.io | v1alpha1 | Predictor | Namespaced | 40 | 0 | [`config/crd/bases/serving.kserve.io_predictors.yaml`](https://github.com/kserve/modelmesh-serving/blob/056bc2e855779c02536db9ef786b26cc73c63f20/config/crd/bases/serving.kserve.io_predictors.yaml) |
| serving.kserve.io | v1alpha1 | ServingRuntime | Namespaced | 1140 | 0 | [`config/crd/bases/serving.kserve.io_servingruntimes.yaml`](https://github.com/kserve/modelmesh-serving/blob/056bc2e855779c02536db9ef786b26cc73c63f20/config/crd/bases/serving.kserve.io_servingruntimes.yaml) |
| serving.kserve.io | v1beta1 | InferenceService | Namespaced | 6195 | 0 | [`config/crd/bases/serving.kserve.io_inferenceservices.yaml`](https://github.com/kserve/modelmesh-serving/blob/056bc2e855779c02536db9ef786b26cc73c63f20/config/crd/bases/serving.kserve.io_inferenceservices.yaml) |

## Dependencies

### Key External Dependencies

| Module | Version |
|--------|---------|
| github.com/go-logr/logr | v1.4.1 |
| github.com/operator-framework/operator-lib | v0.10.0 |
| github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring | v0.55.0 |
| google.golang.org/grpc | v1.59.0 |
| k8s.io/api | v0.28.4 |
| k8s.io/apimachinery | v0.30.13 |
| k8s.io/client-go | v0.28.4 |
| sigs.k8s.io/controller-runtime | v0.16.3 |

