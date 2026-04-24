# trustyai-service-operator: Dataflow

## Controller Watches

Kubernetes resources this controller monitors for changes. Each watch triggers reconciliation when the watched resource is created, updated, or deleted.

| Type | GVK | Source |
|------|-----|--------|
| For | batch/v1/Job | [`controllers/evalhub/evaluation_job_failure_reconciler.go:216`](https://github.com/trustyai-explainability/trustyai-service-operator/blob/52fa0632c89259e4b8a3246ddf94bbbb5796a718/controllers/evalhub/evaluation_job_failure_reconciler.go#L216) |
| For | evalhub/v1alpha1/EvalHub | [`controllers/evalhub/evalhub_controller.go:245`](https://github.com/trustyai-explainability/trustyai-service-operator/blob/52fa0632c89259e4b8a3246ddf94bbbb5796a718/controllers/evalhub/evalhub_controller.go#L245) |
| For | gorch/v1alpha1/GuardrailsOrchestrator | [`controllers/gorch/guardrailsorchestrator_controller.go:410`](https://github.com/trustyai-explainability/trustyai-service-operator/blob/52fa0632c89259e4b8a3246ddf94bbbb5796a718/controllers/gorch/guardrailsorchestrator_controller.go#L410) |
| For | lmes/v1alpha1/LMEvalJob | [`controllers/lmes/lmevaljob_controller.go:299`](https://github.com/trustyai-explainability/trustyai-service-operator/blob/52fa0632c89259e4b8a3246ddf94bbbb5796a718/controllers/lmes/lmevaljob_controller.go#L299) |
| For | nemo_guardrails/v1alpha1/NemoGuardrails | [`controllers/nemo_guardrails/nemoguardrail_controller.go:215`](https://github.com/trustyai-explainability/trustyai-service-operator/blob/52fa0632c89259e4b8a3246ddf94bbbb5796a718/controllers/nemo_guardrails/nemoguardrail_controller.go#L215) |
| For | tas/v1alpha1/TrustyAIService | [`controllers/tas/trustyaiservice_controller.go:279`](https://github.com/trustyai-explainability/trustyai-service-operator/blob/52fa0632c89259e4b8a3246ddf94bbbb5796a718/controllers/tas/trustyaiservice_controller.go#L279) |
| Owns | /v1/ConfigMap | [`controllers/evalhub/evalhub_controller.go:248`](https://github.com/trustyai-explainability/trustyai-service-operator/blob/52fa0632c89259e4b8a3246ddf94bbbb5796a718/controllers/evalhub/evalhub_controller.go#L248) |
| Owns | /v1/Service | [`controllers/evalhub/evalhub_controller.go:247`](https://github.com/trustyai-explainability/trustyai-service-operator/blob/52fa0632c89259e4b8a3246ddf94bbbb5796a718/controllers/evalhub/evalhub_controller.go#L247) |
| Owns | apps/v1/Deployment | [`controllers/gorch/guardrailsorchestrator_controller.go:411`](https://github.com/trustyai-explainability/trustyai-service-operator/blob/52fa0632c89259e4b8a3246ddf94bbbb5796a718/controllers/gorch/guardrailsorchestrator_controller.go#L411) |
| Owns | apps/v1/Deployment | [`controllers/tas/trustyaiservice_controller.go:280`](https://github.com/trustyai-explainability/trustyai-service-operator/blob/52fa0632c89259e4b8a3246ddf94bbbb5796a718/controllers/tas/trustyaiservice_controller.go#L280) |
| Owns | apps/v1/Deployment | [`controllers/evalhub/evalhub_controller.go:246`](https://github.com/trustyai-explainability/trustyai-service-operator/blob/52fa0632c89259e4b8a3246ddf94bbbb5796a718/controllers/evalhub/evalhub_controller.go#L246) |
| Watches | /v1/Namespace | [`controllers/evalhub/evalhub_controller.go:249`](https://github.com/trustyai-explainability/trustyai-service-operator/blob/52fa0632c89259e4b8a3246ddf94bbbb5796a718/controllers/evalhub/evalhub_controller.go#L249) |
| Watches | serving/v1beta1/InferenceService | [`controllers/tas/trustyaiservice_controller.go:281`](https://github.com/trustyai-explainability/trustyai-service-operator/blob/52fa0632c89259e4b8a3246ddf94bbbb5796a718/controllers/tas/trustyaiservice_controller.go#L281) |

## Reconciliation Flow

How the controller interacts with the Kubernetes API during reconciliation.

```mermaid
sequenceDiagram
    %% Static dataflow for trustyai-service-operator

    participant KubernetesAPI as Kubernetes API
    participant controller_manager as controller-manager

    KubernetesAPI->>+controller_manager: Watch Job (reconcile)
    KubernetesAPI->>+controller_manager: Watch EvalHub (reconcile)
    KubernetesAPI->>+controller_manager: Watch GuardrailsOrchestrator (reconcile)
    KubernetesAPI->>+controller_manager: Watch LMEvalJob (reconcile)
    KubernetesAPI->>+controller_manager: Watch NemoGuardrails (reconcile)
    KubernetesAPI->>+controller_manager: Watch TrustyAIService (reconcile)
    controller_manager->>KubernetesAPI: Create/Update ConfigMap
    controller_manager->>KubernetesAPI: Create/Update Service
    controller_manager->>KubernetesAPI: Create/Update Deployment
    controller_manager->>KubernetesAPI: Create/Update Deployment
    controller_manager->>KubernetesAPI: Create/Update Deployment
    KubernetesAPI-->>+controller_manager: Watch Namespace (informer)
    KubernetesAPI-->>+controller_manager: Watch InferenceService (informer)
```

## Configuration

ConfigMaps and Helm values that control this component's runtime behavior.

