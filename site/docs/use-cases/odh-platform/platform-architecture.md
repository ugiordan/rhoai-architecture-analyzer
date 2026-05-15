# Platform Architecture

## CRD Ownership Map

The platform defines 80 CRDs. Each CRD is owned by the component that declares it.

| Owner | CRDs | Count |
|-------|------|-------|
| **data-science-pipelines** | CompositeController, ControllerRevision, DecoratorController | 3 |
| **data-science-pipelines-operator** | DataSciencePipelinesApplication, Pipeline, PipelineVersion, ScheduledWorkflow | 4 |
| **gateway-api-inference-extension** | InferenceModelRewrite, InferenceObjective, InferencePool, InferencePoolImport | 4 |
| **kserve** | ClusterServingRuntime, ClusterStorageContainer, InferenceGraph, InferenceService, LLMInferenceService, LLMInferenceServiceConfig, LocalModelCache, LocalModelNamespaceCache, LocalModelNode, LocalModelNodeGroup, ServingRuntime, TrainedModel | 12 |
| **kserve-autogluon-server** | ClusterServingRuntime, ClusterStorageContainer, InferenceGraph, InferenceService, LLMInferenceService, LLMInferenceServiceConfig, LocalModelCache, LocalModelNamespaceCache, LocalModelNode, LocalModelNodeGroup, ServingRuntime, TrainedModel | 12 |
| **kueue** | ClusterQueue, LocalQueue | 2 |
| **llama-stack-k8s-operator** | LlamaStackDistribution, OGXServer | 2 |
| **mlflow-operator** | MLflow, MLflowConfig | 2 |
| **modelmesh-serving** | ClusterServingRuntime, InferenceService, Predictor, ServingRuntime | 4 |
| **spark-operator** | ScheduledSparkApplication, SparkApplication, SparkConnect | 3 |
| **trainer** | ClusterTrainingRuntime, TrainJob, TrainingRuntime | 3 |
| **workload-variant-autoscaler** | VariantAutoscaling | 1 |

## Cross-Component Dependencies

Relationships detected through Go module imports and CRD watch patterns.

| From | To | Type |
|------|----|------|
| kserve-autogluon-server | gateway-api-inference-extension | watches-crd:InferencePool |
| kserve-autogluon-server | gateway-api-inference-extension | watches-crd:InferenceModelRewrite |
| kserve-autogluon-server | gateway-api-inference-extension | watches-crd:InferenceObjective |
| kserve-autogluon-server | kserve | watches-crd:InferenceGraph |
| kserve-autogluon-server | kserve | watches-crd:LocalModelCache |
| kserve-autogluon-server | kserve | watches-crd:LocalModelNamespaceCache |
| kserve-autogluon-server | kserve | watches-crd:LocalModelNode |
| kserve-autogluon-server | kserve | watches-crd:TrainedModel |
| kserve-autogluon-server | kserve | watches-crd:LLMInferenceService |
| kserve-autogluon-server | kserve | watches-crd:InferenceService |
| kubeflow | data-science-pipelines-operator | go-module |
| llm-d-inference-scheduler | gateway-api-inference-extension | watches-crd:InferencePool |
| llm-d-inference-scheduler | gateway-api-inference-extension | watches-crd:InferenceModelRewrite |
| llm-d-inference-scheduler | gateway-api-inference-extension | watches-crd:InferenceObjective |
| mlflow-operator | mlflow-operator | go-module |
| model-registry | kserve | watches-crd:InferenceService |
| modelmesh-serving | kserve | watches-crd:InferenceGraph |
| modelmesh-serving | kserve | watches-crd:ServingRuntime |
| modelmesh-serving | kserve | watches-crd:TrainedModel |
| modelmesh-serving | kserve | watches-crd:InferenceService |
| models-as-a-service | kserve | go-module |
| workload-variant-autoscaler | gateway-api-inference-extension | watches-crd:InferencePool |
| workload-variant-autoscaler | gateway-api-inference-extension | watches-crd:InferenceObjective |
| kserve | modelmesh-serving | webhook-ref |
| kserve | kserve-autogluon-server | webhook-ref |
| kserve-autogluon-server | modelmesh-serving | webhook-ref |
| modelmesh-serving | kuberay | webhook-ref |
| spark-operator | kuberay | webhook-ref |

**Tightest coupling:** `kserve-autogluon-server -> kserve` (7 dependency edges).

## Webhooks

**Total webhooks across platform**: 153

| Component | Webhooks |
|-----------|----------|
| data-science-pipelines-operator | 1 |
| distributed-workloads | 40 |
| kserve | 31 |
| kserve-autogluon-server | 29 |
| kubeflow | 2 |
| kuberay | 3 |
| kueue | 23 |
| llama-stack-k8s-operator | 1 |
| llm-d-inference-scheduler | 5 |
| modelmesh-serving | 2 |
| spark-operator | 6 |
| trainer | 6 |
| workload-variant-autoscaler | 4 |

### Cross-Component Webhooks

Webhooks whose service reference points to a different component:

| Webhook | Owner | Target Component | Service |
|---------|-------|------------------|---------|
| clusterservingruntime.kserve-webhook-server.validator | kserve | modelmesh-serving | kserve/kserve-webhook-server-service |
| conversion-unknown | kserve | kserve-autogluon-server | kserve/llmisvc-webhook-server-service |
| inferencegraph.kserve-webhook-server.validator | kserve | modelmesh-serving | kserve/kserve-webhook-server-service |
| inferenceservice.kserve-webhook-server.defaulter | kserve | modelmesh-serving | kserve/kserve-webhook-server-service |
| inferenceservice.kserve-webhook-server.pod-mutator | kserve | modelmesh-serving | kserve/kserve-webhook-server-service |
| inferenceservice.kserve-webhook-server.validator | kserve | modelmesh-serving | kserve/kserve-webhook-server-service |
| llminferenceservice.kserve-webhook-server.v1alpha1.defaulter | kserve | kserve-autogluon-server | kserve/llmisvc-webhook-server-service |
| llminferenceservice.kserve-webhook-server.v1alpha1.validator | kserve | kserve-autogluon-server | kserve/llmisvc-webhook-server-service |
| llminferenceservice.kserve-webhook-server.v1alpha2.defaulter | kserve | kserve-autogluon-server | kserve/llmisvc-webhook-server-service |
| llminferenceservice.kserve-webhook-server.v1alpha2.validator | kserve | kserve-autogluon-server | kserve/llmisvc-webhook-server-service |
| llminferenceserviceconfig.kserve-webhook-server.v1alpha1.validator | kserve | kserve-autogluon-server | kserve/llmisvc-webhook-server-service |
| llminferenceserviceconfig.kserve-webhook-server.v1alpha2.validator | kserve | kserve-autogluon-server | kserve/llmisvc-webhook-server-service |
| localmodelcache.kserve-webhook-server.validator | kserve | kserve-autogluon-server | kserve/localmodel-webhook-server-service |
| servingruntime.kserve-webhook-server.validator | kserve | modelmesh-serving | kserve/kserve-webhook-server-service |
| trainedmodel.kserve-webhook-server.validator | kserve | modelmesh-serving | kserve/kserve-webhook-server-service |
| clusterservingruntime.kserve-webhook-server.validator | kserve-autogluon-server | modelmesh-serving | kserve/kserve-webhook-server-service |
| inferencegraph.kserve-webhook-server.validator | kserve-autogluon-server | modelmesh-serving | kserve/kserve-webhook-server-service |
| inferenceservice.kserve-webhook-server.defaulter | kserve-autogluon-server | modelmesh-serving | kserve/kserve-webhook-server-service |
| inferenceservice.kserve-webhook-server.pod-mutator | kserve-autogluon-server | modelmesh-serving | kserve/kserve-webhook-server-service |
| inferenceservice.kserve-webhook-server.validator | kserve-autogluon-server | modelmesh-serving | kserve/kserve-webhook-server-service |
| servingruntime.kserve-webhook-server.validator | kserve-autogluon-server | modelmesh-serving | kserve/kserve-webhook-server-service |
| trainedmodel.kserve-webhook-server.validator | kserve-autogluon-server | modelmesh-serving | kserve/kserve-webhook-server-service |
| conversion-unknown | modelmesh-serving | kuberay | /webhook-service |
| conversion-unknown | spark-operator | kuberay | system/webhook-service |

### External Webhooks

Webhooks referencing services not in the analyzed component set:

| Webhook | Owner | Service |
|---------|-------|---------|
| mappwrapper.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |
| mclusterqueue.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |
| mdeployment.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |
| mjaxjob.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |
| mjob.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |
| mjobset.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |
| mleaderworkerset.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |
| mmpijob.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |
| mpaddlejob.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |
| mpod.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |
| mpytorchjob.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |
| mraycluster.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |
| mrayjob.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |
| mresourceflavor.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |
| mstatefulset.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |
| mtfjob.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |
| mtrainjob.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |
| mworkload.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |
| mxgboostjob.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |
| vappwrapper.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |
| vclusterqueue.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |
| vcohort.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |
| vdeployment.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |
| vjaxjob.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |
| vjob.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |
| vleaderworkerset.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |
| vmpijob.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |
| vpaddlejob.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |
| vpod.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |
| vpytorchjob.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |
| vraycluster.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |
| vresourceflavor.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |
| vstatefulset.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |
| vtfjob.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |
| vtrainjob.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |
| vworkload.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |
| vxgboostjob.kb.io | distributed-workloads | openshift-kueue-operator/kueue-webhook-service |

