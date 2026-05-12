# Platform Architecture

## CRD Ownership Map

The platform defines 23 CRDs. Each CRD is owned by the component that declares it.

| Owner | CRDs | Count |
|-------|------|-------|
| **data-science-pipelines** | CompositeController, ControllerRevision, DecoratorController | 3 |
| **data-science-pipelines-operator** | DataSciencePipelinesApplication, Pipeline, PipelineVersion, ScheduledWorkflow | 4 |
| **kserve** | ClusterServingRuntime, ClusterStorageContainer, InferenceGraph, InferenceService, LLMInferenceService, LLMInferenceServiceConfig, LocalModelCache, LocalModelNamespaceCache, LocalModelNode, LocalModelNodeGroup, ServingRuntime, TrainedModel | 12 |
| **modelmesh-serving** | ClusterServingRuntime, InferenceService, Predictor, ServingRuntime | 4 |

## Cross-Component Dependencies

Relationships detected through Go module imports and CRD watch patterns.

| From | To | Type |
|------|----|------|
| modelmesh-serving | kserve | watches-crd:ServingRuntime |
| data-science-pipelines | kserve | code-ref |

## Webhooks

**Total webhooks across platform**: 16

| Component | Webhooks |
|-----------|----------|
| data-science-pipelines-operator | 1 |
| kserve | 13 |
| modelmesh-serving | 2 |

### External Webhooks

Webhooks referencing services not in the analyzed component set:

| Webhook | Owner | Service |
|---------|-------|---------|
| clusterservingruntime.kserve-webhook-server.validator | kserve | $(kserveNamespace)/$(webhookServiceName) |
| conversion-unknown | modelmesh-serving | /webhook-service |

