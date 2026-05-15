# Platform Architecture

## CRD Ownership Map

The platform defines 56 CRDs. Each CRD is owned by the component that declares it.

| Owner | CRDs | Count |
|-------|------|-------|
| **ai-gateway-payload-processing** | ExternalModel, ExternalProvider | 2 |
| **data-science-pipelines** | CompositeController, ControllerRevision, DecoratorController | 3 |
| **data-science-pipelines-operator** | DataSciencePipelinesApplication, Pipeline, PipelineVersion, ScheduledWorkflow | 4 |
| **gateway-api-inference-extension** | InferenceModelRewrite, InferenceObjective, InferencePool, InferencePoolImport | 4 |
| **kserve** | ClusterServingRuntime, ClusterStorageContainer, InferenceGraph, InferenceService, LLMInferenceService, LLMInferenceServiceConfig, LocalModelCache, LocalModelNamespaceCache, LocalModelNode, LocalModelNodeGroup, ServingRuntime, TrainedModel | 12 |
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
| ai-gateway-payload-processing | gateway-api-inference-extension | watches-crd:InferencePool |
| ai-gateway-payload-processing | gateway-api-inference-extension | watches-crd:InferenceModelRewrite |
| ai-gateway-payload-processing | gateway-api-inference-extension | watches-crd:InferenceObjective |
| kserve | gateway-api-inference-extension | watches-crd:InferencePool |
| kserve | gateway-api-inference-extension | watches-crd:InferenceModelRewrite |
| kserve | gateway-api-inference-extension | watches-crd:InferenceObjective |
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
| models-as-a-service | ai-gateway-payload-processing | watches-crd:ExternalModel |
| odh-cli | opendatahub-operator | go-module |
| workload-variant-autoscaler | gateway-api-inference-extension | watches-crd:InferencePool |
| workload-variant-autoscaler | gateway-api-inference-extension | watches-crd:InferenceObjective |
| kserve | modelmesh-serving | webhook-ref |
| modelmesh-serving | workload-variant-autoscaler | webhook-ref |
| spark-operator | workload-variant-autoscaler | webhook-ref |

**Tightest coupling:** `modelmesh-serving -> kserve` (4 dependency edges).

## Webhooks

**Total webhooks across platform**: 121

| Component | Webhooks |
|-----------|----------|
| data-science-pipelines-operator | 1 |
| distributed-workloads | 40 |
| kserve | 29 |
| kubeflow | 2 |
| kuberay | 2 |
| kueue | 23 |
| llama-stack-k8s-operator | 1 |
| llm-d-inference-scheduler | 5 |
| modelmesh-serving | 2 |
| spark-operator | 6 |
| trainer | 6 |
| workload-variant-autoscaler | 4 |

### Cross-Component Webhooks

Webhooks whose service reference points to a different component:

| Webhook | Type | Owner | Target Component | Target Type | Path |
|---------|------|-------|------------------|-------------|------|
| inferencegraph.kserve-webhook-server.validator | validating | kserve | modelmesh-serving |  | /validate-serving-kserve-io-v1alpha1-inferencegraph |
| inferenceservice.kserve-webhook-server.defaulter | mutating | kserve | modelmesh-serving |  | /mutate-serving-kserve-io-v1beta1-inferenceservice |
| inferenceservice.kserve-webhook-server.pod-mutator | mutating | kserve | modelmesh-serving | Mutator | /mutate-pods |
| inferenceservice.kserve-webhook-server.validator | validating | kserve | modelmesh-serving |  | /validate-serving-kserve-io-v1beta1-inferenceservice |
| servingruntime.kserve-webhook-server.validator | validating | kserve | modelmesh-serving | ServingRuntimeValidator | /validate-serving-kserve-io-v1alpha1-servingruntime |
| trainedmodel.kserve-webhook-server.validator | validating | kserve | modelmesh-serving |  | /validate-serving-kserve-io-v1alpha1-trainedmodel |
| conversion-unknown | conversion | modelmesh-serving | workload-variant-autoscaler |  | /convert |
| conversion-unknown | conversion | spark-operator | workload-variant-autoscaler |  | /convert |

### External Webhooks

| Webhook | Type | Owner | Target Type | Path | Failure Policy |
|---------|------|-------|-------------|------|----------------|
| mappwrapper.kb.io | mutating | distributed-workloads |  | /mutate-workload-codeflare-dev-v1beta2-appwrapper | Fail |
| mclusterqueue.kb.io | mutating | distributed-workloads |  | /mutate-kueue-x-k8s-io-v1beta2-clusterqueue | Fail |
| mdeployment.kb.io | mutating | distributed-workloads |  | /mutate-apps-v1-deployment | Fail |
| mjaxjob.kb.io | mutating | distributed-workloads |  | /mutate-kubeflow-org-v1-jaxjob | Fail |
| mjob.kb.io | mutating | distributed-workloads |  | /mutate-batch-v1-job | Fail |
| mjobset.kb.io | mutating | distributed-workloads |  | /mutate-jobset-x-k8s-io-v1alpha2-jobset | Fail |
| mleaderworkerset.kb.io | mutating | distributed-workloads |  | /mutate-leaderworkerset-x-k8s-io-v1-leaderworkerset | Fail |
| mmpijob.kb.io | mutating | distributed-workloads |  | /mutate-kubeflow-org-v2beta1-mpijob | Fail |
| mpaddlejob.kb.io | mutating | distributed-workloads |  | /mutate-kubeflow-org-v1-paddlejob | Fail |
| mpod.kb.io | mutating | distributed-workloads |  | /mutate--v1-pod | Fail |
| mpytorchjob.kb.io | mutating | distributed-workloads |  | /mutate-kubeflow-org-v1-pytorchjob | Fail |
| mraycluster.kb.io | mutating | distributed-workloads |  | /mutate-ray-io-v1-raycluster | Fail |
| mrayjob.kb.io | mutating | distributed-workloads |  | /mutate-ray-io-v1-rayjob | Fail |
| mresourceflavor.kb.io | mutating | distributed-workloads |  | /mutate-kueue-x-k8s-io-v1beta2-resourceflavor | Fail |
| mstatefulset.kb.io | mutating | distributed-workloads |  | /mutate-apps-v1-statefulset | Fail |
| mtfjob.kb.io | mutating | distributed-workloads |  | /mutate-kubeflow-org-v1-tfjob | Fail |
| mtrainjob.kb.io | mutating | distributed-workloads |  | /mutate-trainer-kubeflow-org-v1alpha1-trainjob | Fail |
| mworkload.kb.io | mutating | distributed-workloads |  | /mutate-kueue-x-k8s-io-v1beta2-workload | Fail |
| mxgboostjob.kb.io | mutating | distributed-workloads |  | /mutate-kubeflow-org-v1-xgboostjob | Fail |
| vappwrapper.kb.io | validating | distributed-workloads |  | /validate-workload-codeflare-dev-v1beta2-appwrapper | Fail |
| vclusterqueue.kb.io | validating | distributed-workloads |  | /validate-kueue-x-k8s-io-v1beta2-clusterqueue | Fail |
| vcohort.kb.io | validating | distributed-workloads |  | /validate-kueue-x-k8s-io-v1beta2-cohort | Fail |
| vdeployment.kb.io | validating | distributed-workloads |  | /validate-apps-v1-deployment | Fail |
| vjaxjob.kb.io | validating | distributed-workloads |  | /validate-kubeflow-org-v1-jaxjob | Fail |
| vjob.kb.io | validating | distributed-workloads |  | /validate-batch-v1-job | Fail |
| vleaderworkerset.kb.io | validating | distributed-workloads |  | /validate-leaderworkerset-x-k8s-io-v1-leaderworkerset | Fail |
| vmpijob.kb.io | validating | distributed-workloads |  | /validate-kubeflow-org-v2beta1-mpijob | Fail |
| vpaddlejob.kb.io | validating | distributed-workloads |  | /validate-kubeflow-org-v1-paddlejob | Fail |
| vpod.kb.io | validating | distributed-workloads |  | /validate--v1-pod | Fail |
| vpytorchjob.kb.io | validating | distributed-workloads |  | /validate-kubeflow-org-v1-pytorchjob | Fail |
| vraycluster.kb.io | validating | distributed-workloads |  | /validate-ray-io-v1-raycluster | Fail |
| vresourceflavor.kb.io | validating | distributed-workloads |  | /validate-kueue-x-k8s-io-v1beta2-resourceflavor | Fail |
| vstatefulset.kb.io | validating | distributed-workloads |  | /validate-apps-v1-statefulset | Fail |
| vtfjob.kb.io | validating | distributed-workloads |  | /validate-kubeflow-org-v1-tfjob | Fail |
| vtrainjob.kb.io | validating | distributed-workloads |  | /validate-trainer-kubeflow-org-v1alpha1-trainjob | Fail |
| vworkload.kb.io | validating | distributed-workloads |  | /validate-kueue-x-k8s-io-v1beta2-workload | Fail |
| vxgboostjob.kb.io | validating | distributed-workloads |  | /validate-kubeflow-org-v1-xgboostjob | Fail |
| clusterservingruntime.kserve-webhook-server.validator | validating | kserve |  | /validate-serving-kserve-io-v1alpha1-clusterservingruntime | Fail |
| mraycluster.kb.io | mutating | kuberay |  | /mutate-ray-io-v1-raycluster | Fail |

