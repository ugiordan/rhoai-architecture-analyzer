# kserve

> **Architecture snapshot: 2026-04-24** (2026-04-24)


**Repository:** kserve/kserve  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-04-24T08:14:51Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 14 |
| Deployments | 10 |
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
        dep_2["kserve-controller-manager"]
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
        dep_8["kserve-localmodel-controller-manager"]
        class dep_8 controller
        dep_9["llmisvc-controller-manager"]
        class dep_9 controller
        dep_10["spark-pmml-iris"]
        class dep_10 controller
    end

    crd_ClusterServingRuntime{{"ClusterServingRuntime\nserving.kserve.io/v1alpha1"}}
    class crd_ClusterServingRuntime crd
    crd_ClusterStorageContainer{{"ClusterStorageContainer\nserving.kserve.io/v1alpha1"}}
    class crd_ClusterStorageContainer crd
    crd_InferenceGraph{{"InferenceGraph\nserving.kserve.io/v1alpha1"}}
    class crd_InferenceGraph crd
    crd_InferenceGraph -->|"For (reconciles)"| controller
    crd_LLMInferenceService{{"LLMInferenceService\nserving.kserve.io/v1alpha1"}}
    class crd_LLMInferenceService crd
    crd_LLMInferenceService -->|"For (reconciles)"| controller
    crd_LLMInferenceServiceConfig{{"LLMInferenceServiceConfig\nserving.kserve.io/v1alpha1"}}
    class crd_LLMInferenceServiceConfig crd
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
    controller -->|"Owns"| owned_11["Deployment"]
    class owned_11 owned
    controller -->|"Owns"| owned_12["HTTPRoute"]
    class owned_12 owned
    controller -->|"Owns"| owned_13["HorizontalPodAutoscaler"]
    class owned_13 owned
    controller -->|"Owns"| owned_14["InferencePool"]
    class owned_14 owned
    controller -->|"Owns"| owned_15["Ingress"]
    class owned_15 owned
    controller -->|"Owns"| owned_16["Job"]
    class owned_16 owned
    controller -->|"Owns"| owned_17["LeaderWorkerSet"]
    class owned_17 owned
    controller -->|"Owns"| owned_18["OpenTelemetryCollector"]
    class owned_18 owned
    controller -->|"Owns"| owned_19["PersistentVolume"]
    class owned_19 owned
    controller -->|"Owns"| owned_20["PersistentVolumeClaim"]
    class owned_20 owned
    controller -->|"Owns"| owned_21["ScaledObject"]
    class owned_21 owned
    controller -->|"Owns"| owned_22["Secret"]
    class owned_22 owned
    controller -->|"Owns"| owned_23["Service"]
    class owned_23 owned
    controller -->|"Owns"| owned_24["VariantAutoscaling"]
    class owned_24 owned
    controller -->|"Owns"| owned_25["VirtualService"]
    class owned_25 owned
    watch_26["ClusterServingRuntime"] -->|"Watches"| controller
    class watch_26 external
    watch_27["ConfigMap"] -->|"Watches"| controller
    class watch_27 external
    watch_28["Gateway"] -->|"Watches"| controller
    class watch_28 external
    watch_29["HTTPRoute"] -->|"Watches"| controller
    class watch_29 external
    watch_30["InferenceService"] -->|"Watches"| controller
    class watch_30 external
    watch_31["LLMInferenceServiceConfig"] -->|"Watches"| controller
    class watch_31 external
    watch_32["LocalModelNode"] -->|"Watches"| controller
    class watch_32 external
    watch_33["Node"] -->|"Watches"| controller
    class watch_33 external
    watch_34["Pod"] -->|"Watches"| controller
    class watch_34 external
    watch_35["ServingRuntime"] -->|"Watches"| controller
    class watch_35 external
```

### CRDs

| Group | Version | Kind | Scope | Fields | Validation Rules | Source |
|-------|---------|------|-------|--------|------------------|--------|
| serving.kserve.io | v1alpha1 | ClusterServingRuntime | Cluster | 1183 | 0 | [`config/crd/full/serving.kserve.io_clusterservingruntimes.yaml`](https://github.com/kserve/kserve/blob/7b52c79e4a1520709d784c755b028a41be371072/config/crd/full/serving.kserve.io_clusterservingruntimes.yaml) |
| serving.kserve.io | v1alpha1 | ClusterStorageContainer | Cluster | 216 | 0 | [`config/crd/full/clusterstoragecontainer/serving.kserve.io_clusterstoragecontainers.yaml`](https://github.com/kserve/kserve/blob/7b52c79e4a1520709d784c755b028a41be371072/config/crd/full/clusterstoragecontainer/serving.kserve.io_clusterstoragecontainers.yaml) |
| serving.kserve.io | v1alpha1 | InferenceGraph | Namespaced | 150 | 0 | [`config/crd/full/serving.kserve.io_inferencegraphs.yaml`](https://github.com/kserve/kserve/blob/7b52c79e4a1520709d784c755b028a41be371072/config/crd/full/serving.kserve.io_inferencegraphs.yaml) |
| serving.kserve.io | v1alpha1 | LLMInferenceService | Namespaced | 5731 | 108 | [`config/crd/full/llmisvc/serving.kserve.io_llminferenceservices.yaml`](https://github.com/kserve/kserve/blob/7b52c79e4a1520709d784c755b028a41be371072/config/crd/full/llmisvc/serving.kserve.io_llminferenceservices.yaml) |
| serving.kserve.io | v1alpha1 | LLMInferenceServiceConfig | Namespaced | 5710 | 108 | [`config/crd/full/llmisvc/serving.kserve.io_llminferenceserviceconfigs.yaml`](https://github.com/kserve/kserve/blob/7b52c79e4a1520709d784c755b028a41be371072/config/crd/full/llmisvc/serving.kserve.io_llminferenceserviceconfigs.yaml) |
| serving.kserve.io | v1alpha1 | LocalModelCache | Cluster | 20 | 1 | [`config/crd/full/localmodel/serving.kserve.io_localmodelcaches.yaml`](https://github.com/kserve/kserve/blob/7b52c79e4a1520709d784c755b028a41be371072/config/crd/full/localmodel/serving.kserve.io_localmodelcaches.yaml) |
| serving.kserve.io | v1alpha1 | LocalModelNamespaceCache | Namespaced | 20 | 1 | [`config/crd/full/localmodel/serving.kserve.io_localmodelnamespacecaches.yaml`](https://github.com/kserve/kserve/blob/7b52c79e4a1520709d784c755b028a41be371072/config/crd/full/localmodel/serving.kserve.io_localmodelnamespacecaches.yaml) |
| serving.kserve.io | v1alpha1 | LocalModelNode | Cluster | 15 | 0 | [`config/crd/full/localmodel/serving.kserve.io_localmodelnodes.yaml`](https://github.com/kserve/kserve/blob/7b52c79e4a1520709d784c755b028a41be371072/config/crd/full/localmodel/serving.kserve.io_localmodelnodes.yaml) |
| serving.kserve.io | v1alpha1 | LocalModelNodeGroup | Cluster | 220 | 0 | [`config/crd/full/localmodel/serving.kserve.io_localmodelnodegroups.yaml`](https://github.com/kserve/kserve/blob/7b52c79e4a1520709d784c755b028a41be371072/config/crd/full/localmodel/serving.kserve.io_localmodelnodegroups.yaml) |
| serving.kserve.io | v1alpha1 | ServingRuntime | Namespaced | 1183 | 0 | [`config/crd/full/serving.kserve.io_servingruntimes.yaml`](https://github.com/kserve/kserve/blob/7b52c79e4a1520709d784c755b028a41be371072/config/crd/full/serving.kserve.io_servingruntimes.yaml) |
| serving.kserve.io | v1alpha1 | TrainedModel | Namespaced | 25 | 0 | [`config/crd/full/serving.kserve.io_trainedmodels.yaml`](https://github.com/kserve/kserve/blob/7b52c79e4a1520709d784c755b028a41be371072/config/crd/full/serving.kserve.io_trainedmodels.yaml) |
| serving.kserve.io | v1alpha2 | LLMInferenceService | Namespaced | 5733 | 110 | [`config/crd/full/llmisvc/serving.kserve.io_llminferenceservices.yaml`](https://github.com/kserve/kserve/blob/7b52c79e4a1520709d784c755b028a41be371072/config/crd/full/llmisvc/serving.kserve.io_llminferenceservices.yaml) |
| serving.kserve.io | v1alpha2 | LLMInferenceServiceConfig | Namespaced | 5712 | 95 | [`config/crd/full/llmisvc/serving.kserve.io_llminferenceserviceconfigs.yaml`](https://github.com/kserve/kserve/blob/7b52c79e4a1520709d784c755b028a41be371072/config/crd/full/llmisvc/serving.kserve.io_llminferenceserviceconfigs.yaml) |
| serving.kserve.io | v1beta1 | InferenceService | Namespaced | 6547 | 0 | [`config/crd/full/serving.kserve.io_inferenceservices.yaml`](https://github.com/kserve/kserve/blob/7b52c79e4a1520709d784c755b028a41be371072/config/crd/full/serving.kserve.io_inferenceservices.yaml) |

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

