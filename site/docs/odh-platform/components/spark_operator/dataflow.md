# spark-operator: Dataflow

## Controller Watches

Kubernetes resources this controller monitors for changes. Each watch triggers reconciliation when the watched resource is created, updated, or deleted.

| Type | GVK | Source |
|------|-----|--------|
| For | api/v1alpha1/SparkConnect | [`internal/controller/sparkconnect/reconciler.go:97`](https://github.com/kubeflow/spark-operator/blob/b8a995788a0bd700354170600d0813db8a597241/internal/controller/sparkconnect/reconciler.go#L97) |

## Reconciliation Flow

How the controller interacts with the Kubernetes API during reconciliation.

```mermaid
sequenceDiagram
    %% Static dataflow for spark-operator

    participant KubernetesAPI as Kubernetes API
    participant controller_manager as controller-manager
    participant spark_operator_controller as spark-operator-controller
    participant spark_operator_webhook as spark-operator-webhook

    KubernetesAPI->>+controller_manager: Watch SparkConnect (reconcile)

    Note over controller_manager: Exposed Services
    Note right of controller_manager: spark-operator-webhook-svc:443/TCP [webhook]

    Note over KubernetesAPI: Defined CRDs
    Note right of KubernetesAPI: SparkConnect (sparkoperator.k8s.io/v1alpha1)
    Note right of KubernetesAPI: ScheduledSparkApplication (sparkoperator.k8s.io/v1beta2)
    Note right of KubernetesAPI: SparkApplication (sparkoperator.k8s.io/v1beta2)
```

### Webhooks

| Name | Type | Path | Failure Policy | Service | Source |
|------|------|------|----------------|---------|--------|
| mutate-pod.sparkoperator.k8s.io | mutating | /mutate--v1-pod | Fail | opendatahub/spark-operator-webhook-svc | [`kustomize:config/overlays/odh (spark-operator-webhook)`](https://github.com/kubeflow/spark-operator/blob/b8a995788a0bd700354170600d0813db8a597241/kustomize:config/overlays/odh (spark-operator-webhook)) |
| mutate-scheduledsparkapplication.sparkoperator.k8s.io | mutating | /mutate-sparkoperator-k8s-io-v1beta2-scheduledsparkapplication | Fail | opendatahub/spark-operator-webhook-svc | [`kustomize:config/overlays/odh (spark-operator-webhook)`](https://github.com/kubeflow/spark-operator/blob/b8a995788a0bd700354170600d0813db8a597241/kustomize:config/overlays/odh (spark-operator-webhook)) |
| mutate-sparkapplication.sparkoperator.k8s.io | mutating | /mutate-sparkoperator-k8s-io-v1beta2-sparkapplication | Fail | opendatahub/spark-operator-webhook-svc | [`kustomize:config/overlays/odh (spark-operator-webhook)`](https://github.com/kubeflow/spark-operator/blob/b8a995788a0bd700354170600d0813db8a597241/kustomize:config/overlays/odh (spark-operator-webhook)) |
| validate-scheduledsparkapplication.sparkoperator.k8s.io | validating | /validate-sparkoperator-k8s-io-v1beta2-scheduledsparkapplication | Fail | opendatahub/spark-operator-webhook-svc | [`kustomize:config/overlays/odh (spark-operator-webhook)`](https://github.com/kubeflow/spark-operator/blob/b8a995788a0bd700354170600d0813db8a597241/kustomize:config/overlays/odh (spark-operator-webhook)) |
| validate-sparkapplication.sparkoperator.k8s.io | validating | /validate-sparkoperator-k8s-io-v1beta2-sparkapplication | Fail | opendatahub/spark-operator-webhook-svc | [`kustomize:config/overlays/odh (spark-operator-webhook)`](https://github.com/kubeflow/spark-operator/blob/b8a995788a0bd700354170600d0813db8a597241/kustomize:config/overlays/odh (spark-operator-webhook)) |

## Configuration

ConfigMaps and Helm values that control this component's runtime behavior.

### Helm

**Chart:** spark-operator v2.4.0

