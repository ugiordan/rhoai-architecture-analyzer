# kserve

> **Architecture snapshot: 2026-05-05** (2026-05-05)


**Repository:** kserve/kserve  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-05-05T15:11:01Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 12 |
| Deployments | 3 |
| Services | 6 |
| Secrets | 3 |
| Cluster Roles | 2 |
| Controller Watches | 46 |

## Component Architecture

CRDs, controllers, and owned Kubernetes resources.

```mermaid
graph LR
    %% Component architecture for kserve

    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef controller fill:#3498db,stroke:#2980b9,color:#fff
    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff

    subgraph controller["kserve Controller"]
        dep_1["kserve-controller-manager"]
        class dep_1 controller
        dep_2["kserve-localmodel-controller-manager"]
        class dep_2 controller
        dep_3["llmisvc-controller-manager"]
        class dep_3 controller
    end

    crd_ClusterServingRuntime{{"ClusterServingRuntime\nserving.kserve.io/v1alpha1"}}
    class crd_ClusterServingRuntime crd
    crd_ClusterStorageContainer{{"ClusterStorageContainer\nserving.kserve.io/v1alpha1"}}
    class crd_ClusterStorageContainer crd
    crd_InferenceGraph{{"InferenceGraph\nserving.kserve.io/v1alpha1"}}
    class crd_InferenceGraph crd
    crd_InferenceGraph -->|"For (reconciles)"| controller
    crd_LocalModelCache{{"LocalModelCache\nserving.kserve.io/v1alpha1"}}
    class crd_LocalModelCache crd
    crd_LocalModelCache -->|"For (reconciles)"| controller
    crd_LocalModelNamespaceCache{{"LocalModelNamespaceCache\nserving.kserve.io/v1alpha1"}}
    class crd_LocalModelNamespaceCache crd
    crd_LocalModelNamespaceCache -->|"For (reconciles)"| controller
    crd_LocalModelNode{{"LocalModelNode\nserving.kserve.io/v1alpha1"}}
    class crd_LocalModelNode crd
    crd_LocalModelNode -->|"For (reconciles)"| controller
    crd_LocalModelNodeGroup{{"LocalModelNodeGroup\nserving.kserve.io/v1alpha1"}}
    class crd_LocalModelNodeGroup crd
    crd_ServingRuntime{{"ServingRuntime\nserving.kserve.io/v1alpha1"}}
    class crd_ServingRuntime crd
    crd_TrainedModel{{"TrainedModel\nserving.kserve.io/v1alpha1"}}
    class crd_TrainedModel crd
    crd_TrainedModel -->|"For (reconciles)"| controller
    crd_LLMInferenceService{{"LLMInferenceService\nserving.kserve.io/v1alpha2"}}
    class crd_LLMInferenceService crd
    crd_LLMInferenceService -->|"For (reconciles)"| controller
    crd_LLMInferenceServiceConfig{{"LLMInferenceServiceConfig\nserving.kserve.io/v1alpha2"}}
    class crd_LLMInferenceServiceConfig crd
    crd_InferenceService{{"InferenceService\nserving.kserve.io/v1beta1"}}
    class crd_InferenceService crd
    crd_InferenceService -->|"For (reconciles)"| controller
    controller -->|"Owns"| owned_4["Deployment"]
    class owned_4 owned
    controller -->|"Owns"| owned_5["HTTPRoute"]
    class owned_5 owned
    controller -->|"Owns"| owned_6["HorizontalPodAutoscaler"]
    class owned_6 owned
    controller -->|"Owns"| owned_7["InferencePool"]
    class owned_7 owned
    controller -->|"Owns"| owned_8["Ingress"]
    class owned_8 owned
    controller -->|"Owns"| owned_9["Job"]
    class owned_9 owned
    controller -->|"Owns"| owned_10["LeaderWorkerSet"]
    class owned_10 owned
    controller -->|"Owns"| owned_11["OpenTelemetryCollector"]
    class owned_11 owned
    controller -->|"Owns"| owned_12["PersistentVolume"]
    class owned_12 owned
    controller -->|"Owns"| owned_13["PersistentVolumeClaim"]
    class owned_13 owned
    controller -->|"Owns"| owned_14["ScaledObject"]
    class owned_14 owned
    controller -->|"Owns"| owned_15["Secret"]
    class owned_15 owned
    controller -->|"Owns"| owned_16["Service"]
    class owned_16 owned
    controller -->|"Owns"| owned_17["VariantAutoscaling"]
    class owned_17 owned
    controller -->|"Owns"| owned_18["VirtualService"]
    class owned_18 owned
    watch_19["ClusterServingRuntime"] -->|"Watches"| controller
    class watch_19 external
    watch_20["ConfigMap"] -->|"Watches"| controller
    class watch_20 external
    watch_21["Gateway"] -->|"Watches"| controller
    class watch_21 external
    watch_22["HTTPRoute"] -->|"Watches"| controller
    class watch_22 external
    watch_23["InferenceService"] -->|"Watches"| controller
    class watch_23 external
    watch_24["LLMInferenceServiceConfig"] -->|"Watches"| controller
    class watch_24 external
    watch_25["LocalModelNode"] -->|"Watches"| controller
    class watch_25 external
    watch_26["Node"] -->|"Watches"| controller
    class watch_26 external
    watch_27["Pod"] -->|"Watches"| controller
    class watch_27 external
    watch_28["ServingRuntime"] -->|"Watches"| controller
    class watch_28 external
```

### CRDs

| Group | Version | Kind | Scope | Fields | Validation Rules | Source |
|-------|---------|------|-------|--------|------------------|--------|
| serving.kserve.io | v1alpha1 | ClusterServingRuntime | Cluster | 1183 | 0 | [`config/crd/full/serving.kserve.io_clusterservingruntimes.yaml`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/config/crd/full/serving.kserve.io_clusterservingruntimes.yaml) |
| serving.kserve.io | v1alpha1 | ClusterStorageContainer | Cluster | 216 | 0 | [`config/crd/full/clusterstoragecontainer/serving.kserve.io_clusterstoragecontainers.yaml`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/config/crd/full/clusterstoragecontainer/serving.kserve.io_clusterstoragecontainers.yaml) |
| serving.kserve.io | v1alpha1 | InferenceGraph | Namespaced | 150 | 0 | [`config/crd/full/serving.kserve.io_inferencegraphs.yaml`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/config/crd/full/serving.kserve.io_inferencegraphs.yaml) |
| serving.kserve.io | v1alpha1 | LocalModelCache | Cluster | 20 | 1 | [`config/crd/full/localmodel/serving.kserve.io_localmodelcaches.yaml`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/config/crd/full/localmodel/serving.kserve.io_localmodelcaches.yaml) |
| serving.kserve.io | v1alpha1 | LocalModelNamespaceCache | Namespaced | 20 | 1 | [`config/crd/full/localmodel/serving.kserve.io_localmodelnamespacecaches.yaml`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/config/crd/full/localmodel/serving.kserve.io_localmodelnamespacecaches.yaml) |
| serving.kserve.io | v1alpha1 | LocalModelNode | Cluster | 15 | 0 | [`config/crd/full/localmodel/serving.kserve.io_localmodelnodes.yaml`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/config/crd/full/localmodel/serving.kserve.io_localmodelnodes.yaml) |
| serving.kserve.io | v1alpha1 | LocalModelNodeGroup | Cluster | 220 | 0 | [`config/crd/full/localmodel/serving.kserve.io_localmodelnodegroups.yaml`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/config/crd/full/localmodel/serving.kserve.io_localmodelnodegroups.yaml) |
| serving.kserve.io | v1alpha1 | ServingRuntime | Namespaced | 1183 | 0 | [`config/crd/full/serving.kserve.io_servingruntimes.yaml`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/config/crd/full/serving.kserve.io_servingruntimes.yaml) |
| serving.kserve.io | v1alpha1 | TrainedModel | Namespaced | 25 | 0 | [`config/crd/full/serving.kserve.io_trainedmodels.yaml`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/config/crd/full/serving.kserve.io_trainedmodels.yaml) |
| serving.kserve.io | v1alpha2 | LLMInferenceService | Namespaced | 5733 | 110 | [`config/crd/full/llmisvc/serving.kserve.io_llminferenceservices.yaml`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/config/crd/full/llmisvc/serving.kserve.io_llminferenceservices.yaml) |
| serving.kserve.io | v1alpha2 | LLMInferenceServiceConfig | Namespaced | 5712 | 95 | [`config/crd/full/llmisvc/serving.kserve.io_llminferenceserviceconfigs.yaml`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/config/crd/full/llmisvc/serving.kserve.io_llminferenceserviceconfigs.yaml) |
| serving.kserve.io | v1beta1 | InferenceService | Namespaced | 6547 | 0 | [`config/crd/full/serving.kserve.io_inferenceservices.yaml`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/config/crd/full/serving.kserve.io_inferenceservices.yaml) |

## Dependencies

### Key External Dependencies

| Module | Version |
|--------|---------|
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/zapr | v1.3.0 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/common | v0.67.4 |
| k8s.io/api | v0.34.5 |
| k8s.io/apiextensions-apiserver | v0.34.3 |
| k8s.io/apimachinery | v0.34.5 |
| k8s.io/client-go | v0.34.5 |
| sigs.k8s.io/controller-runtime | v0.19.7 |

