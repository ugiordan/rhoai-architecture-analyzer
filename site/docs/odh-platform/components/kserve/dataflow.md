# kserve: Dataflow

## Controller Watches

Kubernetes resources this controller monitors for changes. Each watch triggers reconciliation when the watched resource is created, updated, or deleted.

| Type | GVK | Source |
|------|-----|--------|
| For | serving/v1alpha1/InferenceGraph | [`pkg/controller/v1alpha1/inferencegraph/controller.go:373`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha1/inferencegraph/controller.go#L373) |
| For | serving/v1alpha1/LocalModelCache | [`pkg/controller/v1alpha1/localmodel/reconcilers/localmodelcache_reconciler.go:292`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha1/localmodel/reconcilers/localmodelcache_reconciler.go#L292) |
| For | serving/v1alpha1/LocalModelNamespaceCache | [`pkg/controller/v1alpha1/localmodel/reconcilers/localmodelnamespacecache_reconciler.go:295`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha1/localmodel/reconcilers/localmodelnamespacecache_reconciler.go#L295) |
| For | serving/v1alpha1/LocalModelNode | [`pkg/controller/v1alpha1/localmodelnode/controller.go:613`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha1/localmodelnode/controller.go#L613) |
| For | serving/v1alpha1/TrainedModel | [`pkg/controller/v1alpha1/trainedmodel/controller.go:306`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha1/trainedmodel/controller.go#L306) |
| For | serving/v1alpha2/LLMInferenceService | [`pkg/controller/v1alpha2/llmisvc/controller.go:263`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha2/llmisvc/controller.go#L263) |
| For | serving/v1beta1/InferenceService | [`pkg/controller/v1beta1/inferenceservice/controller.go:657`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1beta1/inferenceservice/controller.go#L657) |
| Owns | /v1/PersistentVolume | [`pkg/controller/v1alpha1/localmodel/reconcilers/localmodelcache_reconciler.go:293`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha1/localmodel/reconcilers/localmodelcache_reconciler.go#L293) |
| Owns | /v1/PersistentVolumeClaim | [`pkg/controller/v1alpha1/localmodel/reconcilers/localmodelcache_reconciler.go:294`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha1/localmodel/reconcilers/localmodelcache_reconciler.go#L294) |
| Owns | /v1/PersistentVolumeClaim | [`pkg/controller/v1alpha1/localmodel/reconcilers/localmodelnamespacecache_reconciler.go:296`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha1/localmodel/reconcilers/localmodelnamespacecache_reconciler.go#L296) |
| Owns | /v1/Secret | [`pkg/controller/v1alpha2/llmisvc/controller.go:267`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha2/llmisvc/controller.go#L267) |
| Owns | /v1/Service | [`pkg/controller/v1alpha2/llmisvc/controller.go:268`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha2/llmisvc/controller.go#L268) |
| Owns | /v1/Service | [`pkg/controller/v1beta1/inferenceservice/controller.go:659`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1beta1/inferenceservice/controller.go#L659) |
| Owns | api/v1/InferencePool | [`pkg/controller/v1alpha2/llmisvc/controller.go:288`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha2/llmisvc/controller.go#L288) |
| Owns | api/v1alpha1/VariantAutoscaling | [`pkg/controller/v1alpha2/llmisvc/controller.go:296`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha2/llmisvc/controller.go#L296) |
| Owns | apis/v1/HTTPRoute | [`pkg/controller/v1alpha2/llmisvc/controller.go:280`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha2/llmisvc/controller.go#L280) |
| Owns | apis/v1/HTTPRoute | [`pkg/controller/v1beta1/inferenceservice/controller.go:703`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1beta1/inferenceservice/controller.go#L703) |
| Owns | apis/v1beta1/OpenTelemetryCollector | [`pkg/controller/v1beta1/inferenceservice/controller.go:685`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1beta1/inferenceservice/controller.go#L685) |
| Owns | apix/v1alpha2/InferencePool | [`pkg/controller/v1alpha2/llmisvc/controller.go:292`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha2/llmisvc/controller.go#L292) |
| Owns | apps/v1/Deployment | [`pkg/controller/v1alpha2/llmisvc/controller.go:266`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha2/llmisvc/controller.go#L266) |
| Owns | apps/v1/Deployment | [`pkg/controller/v1beta1/inferenceservice/controller.go:658`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1beta1/inferenceservice/controller.go#L658) |
| Owns | apps/v1/Deployment | [`pkg/controller/v1alpha1/inferencegraph/controller.go:374`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha1/inferencegraph/controller.go#L374) |
| Owns | autoscaling/v2/HorizontalPodAutoscaler | [`pkg/controller/v1alpha2/llmisvc/controller.go:269`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha2/llmisvc/controller.go#L269) |
| Owns | batch/v1/Job | [`pkg/controller/v1alpha1/localmodelnode/controller.go:614`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha1/localmodelnode/controller.go#L614) |
| Owns | keda/v1alpha1/ScaledObject | [`pkg/controller/v1beta1/inferenceservice/controller.go:668`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1beta1/inferenceservice/controller.go#L668) |
| Owns | keda/v1alpha1/ScaledObject | [`pkg/controller/v1alpha2/llmisvc/controller.go:300`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha2/llmisvc/controller.go#L300) |
| Owns | leaderworkerset/v1/LeaderWorkerSet | [`pkg/controller/v1alpha2/llmisvc/controller.go:304`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha2/llmisvc/controller.go#L304) |
| Owns | networking.k8s.io/v1/Ingress | [`pkg/controller/v1alpha2/llmisvc/controller.go:265`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha2/llmisvc/controller.go#L265) |
| Owns | networking.k8s.io/v1/Ingress | [`pkg/controller/v1beta1/inferenceservice/controller.go:709`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1beta1/inferenceservice/controller.go#L709) |
| Owns | networking/v1beta1/VirtualService | [`pkg/controller/v1beta1/inferenceservice/controller.go:691`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1beta1/inferenceservice/controller.go#L691) |
| Owns | serving/v1/Service | [`pkg/controller/v1alpha1/inferencegraph/controller.go:377`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha1/inferencegraph/controller.go#L377) |
| Owns | serving/v1/Service | [`pkg/controller/v1beta1/inferenceservice/controller.go:662`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1beta1/inferenceservice/controller.go#L662) |
| Watches | /v1/ConfigMap | [`pkg/controller/v1alpha2/llmisvc/controller.go:270`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha2/llmisvc/controller.go#L270) |
| Watches | /v1/Node | [`pkg/controller/v1alpha1/localmodel/reconcilers/localmodelcache_reconciler.go:301`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha1/localmodel/reconcilers/localmodelcache_reconciler.go#L301) |
| Watches | /v1/Node | [`pkg/controller/v1alpha1/localmodel/reconcilers/localmodelnamespacecache_reconciler.go:303`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha1/localmodel/reconcilers/localmodelnamespacecache_reconciler.go#L303) |
| Watches | /v1/Pod | [`pkg/controller/v1beta1/inferenceservice/controller.go:714`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1beta1/inferenceservice/controller.go#L714) |
| Watches | /v1/Pod | [`pkg/controller/v1alpha2/llmisvc/controller.go:271`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha2/llmisvc/controller.go#L271) |
| Watches | apis/v1/Gateway | [`pkg/controller/v1alpha2/llmisvc/controller.go:284`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha2/llmisvc/controller.go#L284) |
| Watches | apis/v1/HTTPRoute | [`pkg/controller/v1alpha2/llmisvc/controller.go:281`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha2/llmisvc/controller.go#L281) |
| Watches | serving/v1alpha1/ClusterServingRuntime | [`pkg/controller/v1beta1/inferenceservice/controller.go:713`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1beta1/inferenceservice/controller.go#L713) |
| Watches | serving/v1alpha1/LocalModelNode | [`pkg/controller/v1alpha1/localmodel/reconcilers/localmodelnamespacecache_reconciler.go:304`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha1/localmodel/reconcilers/localmodelnamespacecache_reconciler.go#L304) |
| Watches | serving/v1alpha1/LocalModelNode | [`pkg/controller/v1alpha1/localmodel/reconcilers/localmodelcache_reconciler.go:303`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha1/localmodel/reconcilers/localmodelcache_reconciler.go#L303) |
| Watches | serving/v1alpha1/ServingRuntime | [`pkg/controller/v1beta1/inferenceservice/controller.go:712`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1beta1/inferenceservice/controller.go#L712) |
| Watches | serving/v1alpha2/LLMInferenceServiceConfig | [`pkg/controller/v1alpha2/llmisvc/controller.go:264`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha2/llmisvc/controller.go#L264) |
| Watches | serving/v1beta1/InferenceService | [`pkg/controller/v1alpha1/localmodel/reconcilers/localmodelnamespacecache_reconciler.go:299`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha1/localmodel/reconcilers/localmodelnamespacecache_reconciler.go#L299) |
| Watches | serving/v1beta1/InferenceService | [`pkg/controller/v1alpha1/localmodel/reconcilers/localmodelcache_reconciler.go:297`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha1/localmodel/reconcilers/localmodelcache_reconciler.go#L297) |

## Reconciliation Flow

How the controller interacts with the Kubernetes API during reconciliation.

```mermaid
sequenceDiagram
    %% Static dataflow for kserve

    participant KubernetesAPI as Kubernetes API
    participant kserve_controller_manager as kserve-controller-manager
    participant kserve_localmodel_controller_manager as kserve-localmodel-controller-manager
    participant llmisvc_controller_manager as llmisvc-controller-manager

    KubernetesAPI->>+kserve_controller_manager: Watch InferenceGraph (reconcile)
    KubernetesAPI->>+kserve_controller_manager: Watch LocalModelCache (reconcile)
    KubernetesAPI->>+kserve_controller_manager: Watch LocalModelNamespaceCache (reconcile)
    KubernetesAPI->>+kserve_controller_manager: Watch LocalModelNode (reconcile)
    KubernetesAPI->>+kserve_controller_manager: Watch TrainedModel (reconcile)
    KubernetesAPI->>+kserve_controller_manager: Watch LLMInferenceService (reconcile)
    KubernetesAPI->>+kserve_controller_manager: Watch InferenceService (reconcile)
    kserve_controller_manager->>KubernetesAPI: Create/Update PersistentVolume
    kserve_controller_manager->>KubernetesAPI: Create/Update PersistentVolumeClaim
    kserve_controller_manager->>KubernetesAPI: Create/Update PersistentVolumeClaim
    kserve_controller_manager->>KubernetesAPI: Create/Update Secret
    kserve_controller_manager->>KubernetesAPI: Create/Update Service
    kserve_controller_manager->>KubernetesAPI: Create/Update Service
    kserve_controller_manager->>KubernetesAPI: Create/Update InferencePool
    kserve_controller_manager->>KubernetesAPI: Create/Update VariantAutoscaling
    kserve_controller_manager->>KubernetesAPI: Create/Update HTTPRoute
    kserve_controller_manager->>KubernetesAPI: Create/Update HTTPRoute
    kserve_controller_manager->>KubernetesAPI: Create/Update OpenTelemetryCollector
    kserve_controller_manager->>KubernetesAPI: Create/Update InferencePool
    kserve_controller_manager->>KubernetesAPI: Create/Update Deployment
    kserve_controller_manager->>KubernetesAPI: Create/Update Deployment
    kserve_controller_manager->>KubernetesAPI: Create/Update Deployment
    kserve_controller_manager->>KubernetesAPI: Create/Update HorizontalPodAutoscaler
    kserve_controller_manager->>KubernetesAPI: Create/Update Job
    kserve_controller_manager->>KubernetesAPI: Create/Update ScaledObject
    kserve_controller_manager->>KubernetesAPI: Create/Update ScaledObject
    kserve_controller_manager->>KubernetesAPI: Create/Update LeaderWorkerSet
    kserve_controller_manager->>KubernetesAPI: Create/Update Ingress
    kserve_controller_manager->>KubernetesAPI: Create/Update Ingress
    kserve_controller_manager->>KubernetesAPI: Create/Update VirtualService
    kserve_controller_manager->>KubernetesAPI: Create/Update Service
    kserve_controller_manager->>KubernetesAPI: Create/Update Service
    KubernetesAPI-->>+kserve_controller_manager: Watch ConfigMap (informer)
    KubernetesAPI-->>+kserve_controller_manager: Watch Node (informer)
    KubernetesAPI-->>+kserve_controller_manager: Watch Node (informer)
    KubernetesAPI-->>+kserve_controller_manager: Watch Pod (informer)
    KubernetesAPI-->>+kserve_controller_manager: Watch Pod (informer)
    KubernetesAPI-->>+kserve_controller_manager: Watch Gateway (informer)
    KubernetesAPI-->>+kserve_controller_manager: Watch HTTPRoute (informer)
    KubernetesAPI-->>+kserve_controller_manager: Watch ClusterServingRuntime (informer)
    KubernetesAPI-->>+kserve_controller_manager: Watch LocalModelNode (informer)
    KubernetesAPI-->>+kserve_controller_manager: Watch LocalModelNode (informer)
    KubernetesAPI-->>+kserve_controller_manager: Watch ServingRuntime (informer)
    KubernetesAPI-->>+kserve_controller_manager: Watch LLMInferenceServiceConfig (informer)
    KubernetesAPI-->>+kserve_controller_manager: Watch InferenceService (informer)
    KubernetesAPI-->>+kserve_controller_manager: Watch InferenceService (informer)

    Note over kserve_controller_manager: Exposed Services
    Note right of kserve_controller_manager: kserve-controller-manager-metrics-service:8443/TCP [https]
    Note right of kserve_controller_manager: kserve-controller-manager-service:8443/TCP []
    Note right of kserve_controller_manager: kserve-webhook-server-service:443/TCP []
    Note right of kserve_controller_manager: llmisvc-controller-manager-service:8443/TCP [https]
    Note right of kserve_controller_manager: llmisvc-webhook-server-service:443/TCP [https]
    Note right of kserve_controller_manager: localmodel-webhook-server-service:443/TCP []

    Note over KubernetesAPI: Defined CRDs
    Note right of KubernetesAPI: ClusterServingRuntime (serving.kserve.io/v1alpha1)
    Note right of KubernetesAPI: ClusterStorageContainer (serving.kserve.io/v1alpha1)
    Note right of KubernetesAPI: InferenceGraph (serving.kserve.io/v1alpha1)
    Note right of KubernetesAPI: LocalModelCache (serving.kserve.io/v1alpha1)
    Note right of KubernetesAPI: LocalModelNamespaceCache (serving.kserve.io/v1alpha1)
    Note right of KubernetesAPI: LocalModelNode (serving.kserve.io/v1alpha1)
    Note right of KubernetesAPI: LocalModelNodeGroup (serving.kserve.io/v1alpha1)
    Note right of KubernetesAPI: ServingRuntime (serving.kserve.io/v1alpha1)
    Note right of KubernetesAPI: TrainedModel (serving.kserve.io/v1alpha1)
    Note right of KubernetesAPI: LLMInferenceService (serving.kserve.io/v1alpha2)
    Note right of KubernetesAPI: LLMInferenceServiceConfig (serving.kserve.io/v1alpha2)
    Note right of KubernetesAPI: InferenceService (serving.kserve.io/v1beta1)
```

### Webhooks

| Name | Type | Path | Failure Policy | Service | Source |
|------|------|------|----------------|---------|--------|
| clusterservingruntime.kserve-webhook-server.validator | validating | /validate-serving-kserve-io-v1alpha1-clusterservingruntime | Fail | kserve/kserve-webhook-server-service | [`kustomize:config/overlays/all (clusterservingruntime.serving.kserve.io)`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/kustomize:config/overlays/all (clusterservingruntime.serving.kserve.io)) |
| inferencegraph.kserve-webhook-server.validator | validating | /validate-serving-kserve-io-v1alpha1-inferencegraph | Fail | kserve/kserve-webhook-server-service | [`kustomize:config/overlays/all (inferencegraph.serving.kserve.io)`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/kustomize:config/overlays/all (inferencegraph.serving.kserve.io)) |
| inferenceservice.kserve-webhook-server.defaulter | mutating | /mutate-serving-kserve-io-v1beta1-inferenceservice | Fail | kserve/kserve-webhook-server-service | [`kustomize:config/overlays/all (inferenceservice.serving.kserve.io)`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/kustomize:config/overlays/all (inferenceservice.serving.kserve.io)) |
| inferenceservice.kserve-webhook-server.pod-mutator | mutating | /mutate-pods | Fail | kserve/kserve-webhook-server-service | [`kustomize:config/overlays/all (inferenceservice.serving.kserve.io)`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/kustomize:config/overlays/all (inferenceservice.serving.kserve.io)) |
| inferenceservice.kserve-webhook-server.validator | validating | /validate-serving-kserve-io-v1beta1-inferenceservice | Fail | kserve/kserve-webhook-server-service | [`kustomize:config/overlays/all (inferenceservice.serving.kserve.io)`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/kustomize:config/overlays/all (inferenceservice.serving.kserve.io)) |
| llminferenceservice.kserve-webhook-server.v1alpha1.validator | validating | /validate-serving-kserve-io-v1alpha1-llminferenceservice | Fail | kserve/llmisvc-webhook-server-service | [`kustomize:config/overlays/all (llminferenceservice.serving.kserve.io)`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/kustomize:config/overlays/all (llminferenceservice.serving.kserve.io)) |
| llminferenceservice.kserve-webhook-server.v1alpha2.validator | validating | /validate-serving-kserve-io-v1alpha2-llminferenceservice | Fail | kserve/llmisvc-webhook-server-service | [`kustomize:config/overlays/all (llminferenceservice.serving.kserve.io)`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/kustomize:config/overlays/all (llminferenceservice.serving.kserve.io)) |
| llminferenceserviceconfig.kserve-webhook-server.v1alpha1.validator | validating | /validate-serving-kserve-io-v1alpha1-llminferenceserviceconfig | Fail | kserve/llmisvc-webhook-server-service | [`kustomize:config/overlays/all (llminferenceserviceconfig.serving.kserve.io)`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/kustomize:config/overlays/all (llminferenceserviceconfig.serving.kserve.io)) |
| llminferenceserviceconfig.kserve-webhook-server.v1alpha2.validator | validating | /validate-serving-kserve-io-v1alpha2-llminferenceserviceconfig | Fail | kserve/llmisvc-webhook-server-service | [`kustomize:config/overlays/all (llminferenceserviceconfig.serving.kserve.io)`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/kustomize:config/overlays/all (llminferenceserviceconfig.serving.kserve.io)) |
| localmodelcache.kserve-webhook-server.validator | validating | /validate-serving-kserve-io-v1alpha1-localmodelcache | Fail | kserve/localmodel-webhook-server-service | [`kustomize:config/overlays/all (localmodelcache.serving.kserve.io)`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/kustomize:config/overlays/all (localmodelcache.serving.kserve.io)) |
| servingruntime.kserve-webhook-server.validator | validating | /validate-serving-kserve-io-v1alpha1-servingruntime | Fail | kserve/kserve-webhook-server-service | [`kustomize:config/overlays/all (servingruntime.serving.kserve.io)`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/kustomize:config/overlays/all (servingruntime.serving.kserve.io)) |
| trainedmodel.kserve-webhook-server.validator | validating | /validate-serving-kserve-io-v1alpha1-trainedmodel | Fail | kserve/kserve-webhook-server-service | [`kustomize:config/overlays/all (trainedmodel.serving.kserve.io)`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/kustomize:config/overlays/all (trainedmodel.serving.kserve.io)) |

### HTTP Endpoints

| Method | Path | Source |
|--------|------|--------|
| * | / | [`cmd/router/main.go:506`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/cmd/router/main.go#L506) |
| * | gateway.networking.k8s.io | [`pkg/controller/v1alpha2/llmisvc/config_merge.go:394`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha2/llmisvc/config_merge.go#L394) |
| * | gateway.networking.k8s.io | [`pkg/controller/v1alpha2/llmisvc/fixture/gwapi_builders.go:210`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha2/llmisvc/fixture/gwapi_builders.go#L210) |
| * | gateway.networking.k8s.io | [`pkg/controller/v1alpha2/llmisvc/fixture/gwapi_builders.go:228`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha2/llmisvc/fixture/gwapi_builders.go#L228) |
| * | gateway.networking.k8s.io | [`pkg/controller/v1alpha2/llmisvc/fixture/gwapi_builders.go:398`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha2/llmisvc/fixture/gwapi_builders.go#L398) |
| * | gateway.networking.k8s.io | [`pkg/controller/v1alpha2/llmisvc/fixture/gwapi_builders.go:706`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha2/llmisvc/fixture/gwapi_builders.go#L706) |
| * | inference.networking.k8s.io | [`pkg/controller/v1alpha2/llmisvc/fixture/gwapi_builders.go:290`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha2/llmisvc/fixture/gwapi_builders.go#L290) |
| * | inference.networking.x-k8s.io | [`pkg/controller/v1alpha2/llmisvc/fixture/gwapi_builders.go:304`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/pkg/controller/v1alpha2/llmisvc/fixture/gwapi_builders.go#L304) |

## Configuration

ConfigMaps and Helm values that control this component's runtime behavior.

### ConfigMaps

| Name | Data Keys | Source |
|------|-----------|--------|
| inferenceservice-config | agent, autoscaler, batcher, credentials, deploy, explainers, inferenceService, ingress, localModel, logger, metricsAggregator, opentelemetryCollector, router, security, service, storageInitializer | [`charts/_common/common-patches/configmap-patch.yaml`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/charts/_common/common-patches/configmap-patch.yaml) |
| inferenceservice-config | agent, autoscaler, batcher, credentials, deploy, explainers, inferenceService, ingress, localModel, logger, metricsAggregator, opentelemetryCollector, router, security, service, storageInitializer | [`charts/kserve-llmisvc-resources/files/common/configmap-patch.yaml`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/charts/kserve-llmisvc-resources/files/common/configmap-patch.yaml) |
| inferenceservice-config | _example, agent, autoscaler, batcher, credentials, deploy, explainers, inferenceService, ingress, localModel, logger, metricsAggregator, opentelemetryCollector, router, security, storageInitializer | [`charts/kserve-llmisvc-resources/files/common/configmap.yaml`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/charts/kserve-llmisvc-resources/files/common/configmap.yaml) |
| inferenceservice-config | agent, autoscaler, batcher, credentials, deploy, explainers, inferenceService, ingress, localModel, logger, metricsAggregator, opentelemetryCollector, router, security, service, storageInitializer | [`charts/kserve-resources/files/common/configmap-patch.yaml`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/charts/kserve-resources/files/common/configmap-patch.yaml) |
| inferenceservice-config | _example, agent, autoscaler, batcher, credentials, deploy, explainers, inferenceService, ingress, localModel, logger, metricsAggregator, opentelemetryCollector, router, security, storageInitializer | [`charts/kserve-resources/files/common/configmap.yaml`](https://github.com/kserve/kserve/blob/d5aea2c6d8f2f2c8dcf22897e23e5d929cf654dd/charts/kserve-resources/files/common/configmap.yaml) |

### Helm

**Chart:** kserve-crd vv0.18.0

