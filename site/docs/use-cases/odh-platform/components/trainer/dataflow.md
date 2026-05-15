# trainer: Dataflow

## Controller Watches

Kubernetes resources this controller monitors for changes. Each watch triggers reconciliation when the watched resource is created, updated, or deleted.

| Type | GVK | Source |
|------|-----|--------|
| For | /v1/Pod | [`.gomod-cache/sigs.k8s.io/jobset@v0.10.1/pkg/controllers/pod_controller.go:65`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/sigs.k8s.io/jobset@v0.10.1/pkg/controllers/pod_controller.go#L65) |
| For | /v1/Pod | [`.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.10.1/pkg/controllers/pod_controller.go:65`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.10.1/pkg/controllers/pod_controller.go#L65) |
| For | jobset/v1alpha2/JobSet | [`.gomod-cache/sigs.k8s.io/jobset@v0.10.1/pkg/controllers/jobset_controller.go:231`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/sigs.k8s.io/jobset@v0.10.1/pkg/controllers/jobset_controller.go#L231) |
| For | jobset/v1alpha2/JobSet | [`.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.10.1/pkg/controllers/jobset_controller.go:231`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.10.1/pkg/controllers/jobset_controller.go#L231) |
| For | scheduling/v1alpha1/ElasticQuota | [`.gopath-loader/pkg/mod/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/pkg/controllers/elasticquota_controller.go:175`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/pkg/controllers/elasticquota_controller.go#L175) |
| For | scheduling/v1alpha1/ElasticQuota | [`.gomod-cache/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/pkg/controllers/elasticquota_controller.go:175`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/pkg/controllers/elasticquota_controller.go#L175) |
| For | scheduling/v1alpha1/PodGroup | [`.gopath-loader/pkg/mod/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/pkg/controllers/podgroup_controller.go:193`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/pkg/controllers/podgroup_controller.go#L193) |
| For | scheduling/v1alpha1/PodGroup | [`.gomod-cache/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/pkg/controllers/podgroup_controller.go:193`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/pkg/controllers/podgroup_controller.go#L193) |
| Owns | /v1/Service | [`.gomod-cache/sigs.k8s.io/jobset@v0.10.1/pkg/controllers/jobset_controller.go:233`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/sigs.k8s.io/jobset@v0.10.1/pkg/controllers/jobset_controller.go#L233) |
| Owns | /v1/Service | [`.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.10.1/pkg/controllers/jobset_controller.go:233`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.10.1/pkg/controllers/jobset_controller.go#L233) |
| Owns | batch/v1/Job | [`.gomod-cache/sigs.k8s.io/jobset@v0.10.1/pkg/controllers/jobset_controller.go:232`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/sigs.k8s.io/jobset@v0.10.1/pkg/controllers/jobset_controller.go#L232) |
| Owns | batch/v1/Job | [`.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.10.1/pkg/controllers/jobset_controller.go:232`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.10.1/pkg/controllers/jobset_controller.go#L232) |
| Watches | /v1/Pod | [`.gomod-cache/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/pkg/controllers/elasticquota_controller.go:174`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/pkg/controllers/elasticquota_controller.go#L174) |
| Watches | /v1/Pod | [`.gopath-loader/pkg/mod/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/pkg/controllers/elasticquota_controller.go:174`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/pkg/controllers/elasticquota_controller.go#L174) |
| Watches | /v1/Pod | [`.gopath-loader/pkg/mod/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/pkg/controllers/podgroup_controller.go:192`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/pkg/controllers/podgroup_controller.go#L192) |
| Watches | /v1/Pod | [`.gomod-cache/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/pkg/controllers/podgroup_controller.go:192`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/sigs.k8s.io/scheduler-plugins@v0.34.1-devel/pkg/controllers/podgroup_controller.go#L192) |

### Programmatic Resource Operations

| Verb | Kind | Group | Condition |
|------|------|-------|----------|
| update | ClusterTrainingRuntime | trainer |  |
| update | TrainingRuntime | trainer |  |
| update | TrainJob | trainer |  |

## Reconciliation Flow

How the controller interacts with the Kubernetes API during reconciliation.

```mermaid
sequenceDiagram
    %% Static dataflow for trainer

    participant KubernetesAPI as Kubernetes API
    participant controller_manager as controller-manager
    participant kubeflow_trainer_controller_manager as kubeflow-trainer-controller-manager
    participant peaks as peaks

    KubernetesAPI->>+controller_manager: Watch Pod (reconcile)
    KubernetesAPI->>+controller_manager: Watch Pod (reconcile)
    KubernetesAPI->>+controller_manager: Watch JobSet (reconcile)
    KubernetesAPI->>+controller_manager: Watch JobSet (reconcile)
    KubernetesAPI->>+controller_manager: Watch ElasticQuota (reconcile)
    KubernetesAPI->>+controller_manager: Watch ElasticQuota (reconcile)
    KubernetesAPI->>+controller_manager: Watch PodGroup (reconcile)
    KubernetesAPI->>+controller_manager: Watch PodGroup (reconcile)
    controller_manager->>KubernetesAPI: Create/Update Service
    controller_manager->>KubernetesAPI: Create/Update Service
    controller_manager->>KubernetesAPI: Create/Update Job
    controller_manager->>KubernetesAPI: Create/Update Job
    KubernetesAPI-->>+controller_manager: Watch Pod (informer)
    KubernetesAPI-->>+controller_manager: Watch Pod (informer)
    KubernetesAPI-->>+controller_manager: Watch Pod (informer)
    KubernetesAPI-->>+controller_manager: Watch Pod (informer)

    Note over controller_manager: Exposed Services
    Note right of controller_manager: webhook-service:443/TCP []
    Note right of controller_manager: webhook-service:443/TCP []

    Note over KubernetesAPI: Defined CRDs
    Note right of KubernetesAPI: ClusterTrainingRuntime (trainer.kubeflow.org/v1alpha1)
    Note right of KubernetesAPI: TrainJob (trainer.kubeflow.org/v1alpha1)
    Note right of KubernetesAPI: TrainingRuntime (trainer.kubeflow.org/v1alpha1)
```

### Webhooks

| Name | Type | Path | Failure Policy | Service | Overlays | Enable Condition | Sources |
|------|------|------|----------------|---------|----------|------------------|----------|
| ClusterTrainingRuntimeWebhook-webhook | validating | /validate-trainer-kubeflow-org-v1alpha1-clustertrainingruntime |  |  |  |  |  |
| TrainJobWebhook-webhook | validating | /validate-trainer-kubeflow-org-v1alpha1-trainjob |  |  |  |  |  |
| TrainingRuntimeWebhook-webhook | validating | /validate-trainer-kubeflow-org-v1alpha1-trainingruntime |  |  |  |  |  |
| mjobset.kb.io | mutating | /mutate-jobset-x-k8s-io-v1alpha2-jobset | fail |  |  |  | [`.gomod-cache/sigs.k8s.io/jobset@v0.10.1/pkg/webhooks/jobset_webhook.go`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/sigs.k8s.io/jobset@v0.10.1/pkg/webhooks/jobset_webhook.go), [`.gomod-cache/sigs.k8s.io/jobset@v0.10.1/pkg/webhooks/jobset_webhook.go`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/sigs.k8s.io/jobset@v0.10.1/pkg/webhooks/jobset_webhook.go) |
| mjobset.kb.io | mutating | /mutate-jobset-x-k8s-io-v1alpha2-jobset | fail |  |  |  | [`.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.10.1/pkg/webhooks/jobset_webhook.go`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.10.1/pkg/webhooks/jobset_webhook.go), [`.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.10.1/pkg/webhooks/jobset_webhook.go`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.10.1/pkg/webhooks/jobset_webhook.go) |
| vjobset.kb.io | validating | /validate-jobset-x-k8s-io-v1alpha2-jobset | fail |  |  |  | [`.gomod-cache/sigs.k8s.io/jobset@v0.10.1/pkg/webhooks/jobset_webhook.go`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/sigs.k8s.io/jobset@v0.10.1/pkg/webhooks/jobset_webhook.go), [`.gomod-cache/sigs.k8s.io/jobset@v0.10.1/pkg/webhooks/jobset_webhook.go`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/sigs.k8s.io/jobset@v0.10.1/pkg/webhooks/jobset_webhook.go) |
| vjobset.kb.io | validating | /validate-jobset-x-k8s-io-v1alpha2-jobset | fail |  |  |  | [`.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.10.1/pkg/webhooks/jobset_webhook.go`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.10.1/pkg/webhooks/jobset_webhook.go), [`.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.10.1/pkg/webhooks/jobset_webhook.go`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.10.1/pkg/webhooks/jobset_webhook.go) |
| vpod.kb.io | validating | /validate--v1-pod | fail |  |  |  | [`.gomod-cache/sigs.k8s.io/jobset@v0.10.1/pkg/webhooks/pod_admission_webhook.go`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/sigs.k8s.io/jobset@v0.10.1/pkg/webhooks/pod_admission_webhook.go), [`.gomod-cache/sigs.k8s.io/jobset@v0.10.1/pkg/webhooks/pod_admission_webhook.go`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/sigs.k8s.io/jobset@v0.10.1/pkg/webhooks/pod_admission_webhook.go) |
| vpod.kb.io | validating | /validate--v1-pod | fail |  |  |  | [`.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.10.1/pkg/webhooks/pod_admission_webhook.go`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.10.1/pkg/webhooks/pod_admission_webhook.go), [`.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.10.1/pkg/webhooks/pod_admission_webhook.go`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.10.1/pkg/webhooks/pod_admission_webhook.go) |

### HTTP Endpoints

| Method | Path | Source |
|--------|------|--------|
| * | / | [`.gomod-cache/golang.org/x/net@v0.43.0/webdav/litmus_test_server.go:83`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/golang.org/x/net@v0.43.0/webdav/litmus_test_server.go#L83) |
| * | / | [`.gomod-cache/golang.org/x/tools@v0.36.0/cmd/present/dir.go:23`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/golang.org/x/tools@v0.36.0/cmd/present/dir.go#L23) |
| * | / | [`.gopath-loader/pkg/mod/golang.org/x/tools@v0.36.0/cmd/godoc/handlers.go:42`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/golang.org/x/tools@v0.36.0/cmd/godoc/handlers.go#L42) |
| * | / | [`.gomod-cache/github.com/google/pprof@v0.0.0-20250403155104-27863c87afa6/internal/driver/webui.go:212`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/github.com/google/pprof@v0.0.0-20250403155104-27863c87afa6/internal/driver/webui.go#L212) |
| * | / | [`.gopath-loader/pkg/mod/golang.org/x/tools@v0.36.0/cmd/godoc/handlers.go:31`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/golang.org/x/tools@v0.36.0/cmd/godoc/handlers.go#L31) |
| * | / | [`.gopath-loader/pkg/mod/golang.org/x/net@v0.43.0/webdav/litmus_test_server.go:83`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/golang.org/x/net@v0.43.0/webdav/litmus_test_server.go#L83) |
| * | / | [`.gopath-loader/pkg/mod/golang.org/x/tools@v0.36.0/cmd/present/dir.go:23`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/golang.org/x/tools@v0.36.0/cmd/present/dir.go#L23) |
| * | / | [`.gopath-loader/pkg/mod/golang.org/x/tools@v0.36.0/go/types/internal/play/play.go:46`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/golang.org/x/tools@v0.36.0/go/types/internal/play/play.go#L46) |
| * | / | [`.gopath-loader/pkg/mod/github.com/google/pprof@v0.0.0-20250403155104-27863c87afa6/internal/driver/webui.go:212`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/github.com/google/pprof@v0.0.0-20250403155104-27863c87afa6/internal/driver/webui.go#L212) |
| * | / | [`.gopath-loader/pkg/mod/golang.org/x/tools@v0.36.0/godoc/pres.go:130`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/golang.org/x/tools@v0.36.0/godoc/pres.go#L130) |
| * | / | [`.gomod-cache/golang.org/x/tools@v0.36.0/go/types/internal/play/play.go:46`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/golang.org/x/tools@v0.36.0/go/types/internal/play/play.go#L46) |
| * | / | [`.gomod-cache/golang.org/x/tools@v0.36.0/cmd/godoc/handlers.go:31`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/golang.org/x/tools@v0.36.0/cmd/godoc/handlers.go#L31) |
| * | / | [`.gomod-cache/golang.org/x/tools@v0.36.0/cmd/godoc/handlers.go:42`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/golang.org/x/tools@v0.36.0/cmd/godoc/handlers.go#L42) |
| * | / | [`.gomod-cache/golang.org/x/tools@v0.36.0/godoc/pres.go:130`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/golang.org/x/tools@v0.36.0/godoc/pres.go#L130) |
| * | /abort | [`.gomod-cache/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go:63`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go#L63) |
| * | /abort | [`.gopath-loader/pkg/mod/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go:63`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go#L63) |
| * | /aggregated-nonprimary-procs-report | [`.gomod-cache/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go:60`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go#L60) |
| * | /aggregated-nonprimary-procs-report | [`.gopath-loader/pkg/mod/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go:60`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go#L60) |
| * | /before-suite-completed | [`.gopath-loader/pkg/mod/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go:57`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go#L57) |
| * | /before-suite-completed | [`.gomod-cache/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go:57`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go#L57) |
| * | /before-suite-state | [`.gomod-cache/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go:58`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go#L58) |
| * | /before-suite-state | [`.gopath-loader/pkg/mod/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go:58`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go#L58) |
| * | /compile | [`.gopath-loader/pkg/mod/golang.org/x/tools@v0.36.0/playground/playground.go:23`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/golang.org/x/tools@v0.36.0/playground/playground.go#L23) |
| * | /compile | [`.gomod-cache/golang.org/x/tools@v0.36.0/playground/playground.go:23`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/golang.org/x/tools@v0.36.0/playground/playground.go#L23) |
| * | /counter | [`.gomod-cache/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go:61`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go#L61) |
| * | /counter | [`.gopath-loader/pkg/mod/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go:61`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go#L61) |
| * | /debug/pprof/ | [`.gopath-loader/pkg/mod/sigs.k8s.io/controller-runtime@v0.22.3/pkg/manager/internal.go:316`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/sigs.k8s.io/controller-runtime@v0.22.3/pkg/manager/internal.go#L316) |
| * | /debug/pprof/ | [`.gomod-cache/sigs.k8s.io/controller-runtime@v0.22.3/pkg/manager/internal.go:316`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/sigs.k8s.io/controller-runtime@v0.22.3/pkg/manager/internal.go#L316) |
| * | /debug/pprof/cmdline | [`.gopath-loader/pkg/mod/sigs.k8s.io/controller-runtime@v0.22.3/pkg/manager/internal.go:317`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/sigs.k8s.io/controller-runtime@v0.22.3/pkg/manager/internal.go#L317) |
| * | /debug/pprof/cmdline | [`.gomod-cache/sigs.k8s.io/controller-runtime@v0.22.3/pkg/manager/internal.go:317`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/sigs.k8s.io/controller-runtime@v0.22.3/pkg/manager/internal.go#L317) |
| * | /debug/pprof/profile | [`.gopath-loader/pkg/mod/sigs.k8s.io/controller-runtime@v0.22.3/pkg/manager/internal.go:318`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/sigs.k8s.io/controller-runtime@v0.22.3/pkg/manager/internal.go#L318) |
| * | /debug/pprof/profile | [`.gomod-cache/sigs.k8s.io/controller-runtime@v0.22.3/pkg/manager/internal.go:318`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/sigs.k8s.io/controller-runtime@v0.22.3/pkg/manager/internal.go#L318) |
| * | /debug/pprof/symbol | [`.gopath-loader/pkg/mod/sigs.k8s.io/controller-runtime@v0.22.3/pkg/manager/internal.go:319`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/sigs.k8s.io/controller-runtime@v0.22.3/pkg/manager/internal.go#L319) |
| * | /debug/pprof/symbol | [`.gomod-cache/sigs.k8s.io/controller-runtime@v0.22.3/pkg/manager/internal.go:319`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/sigs.k8s.io/controller-runtime@v0.22.3/pkg/manager/internal.go#L319) |
| * | /debug/pprof/trace | [`.gopath-loader/pkg/mod/sigs.k8s.io/controller-runtime@v0.22.3/pkg/manager/internal.go:320`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/sigs.k8s.io/controller-runtime@v0.22.3/pkg/manager/internal.go#L320) |
| * | /debug/pprof/trace | [`.gomod-cache/sigs.k8s.io/controller-runtime@v0.22.3/pkg/manager/internal.go:320`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/sigs.k8s.io/controller-runtime@v0.22.3/pkg/manager/internal.go#L320) |
| * | /did-run | [`.gopath-loader/pkg/mod/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go:49`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go#L49) |
| * | /did-run | [`.gomod-cache/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go:49`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go#L49) |
| * | /emit-output | [`.gopath-loader/pkg/mod/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go:51`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go#L51) |
| * | /emit-output | [`.gomod-cache/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go:51`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go#L51) |
| * | /fmt | [`.gopath-loader/pkg/mod/golang.org/x/tools@v0.36.0/cmd/godoc/handlers.go:39`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/golang.org/x/tools@v0.36.0/cmd/godoc/handlers.go#L39) |
| * | /fmt | [`.gomod-cache/golang.org/x/tools@v0.36.0/cmd/godoc/handlers.go:39`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/golang.org/x/tools@v0.36.0/cmd/godoc/handlers.go#L39) |
| * | /have-nonprimary-procs-finished | [`.gomod-cache/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go:59`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go#L59) |
| * | /have-nonprimary-procs-finished | [`.gopath-loader/pkg/mod/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go:59`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go#L59) |
| * | /main.css | [`.gomod-cache/golang.org/x/tools@v0.36.0/go/types/internal/play/play.go:48`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/golang.org/x/tools@v0.36.0/go/types/internal/play/play.go#L48) |
| * | /main.css | [`.gopath-loader/pkg/mod/golang.org/x/tools@v0.36.0/go/types/internal/play/play.go:48`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/golang.org/x/tools@v0.36.0/go/types/internal/play/play.go#L48) |
| * | /main.js | [`.gomod-cache/golang.org/x/tools@v0.36.0/go/types/internal/play/play.go:47`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/golang.org/x/tools@v0.36.0/go/types/internal/play/play.go#L47) |
| * | /main.js | [`.gopath-loader/pkg/mod/golang.org/x/tools@v0.36.0/go/types/internal/play/play.go:47`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/golang.org/x/tools@v0.36.0/go/types/internal/play/play.go#L47) |
| * | /opensearch.xml | [`.gopath-loader/pkg/mod/golang.org/x/tools@v0.36.0/godoc/pres.go:133`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/golang.org/x/tools@v0.36.0/godoc/pres.go#L133) |
| * | /opensearch.xml | [`.gomod-cache/golang.org/x/tools@v0.36.0/godoc/pres.go:133`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/golang.org/x/tools@v0.36.0/godoc/pres.go#L133) |
| * | /pkg/C/ | [`.gomod-cache/golang.org/x/tools@v0.36.0/cmd/godoc/handlers.go:38`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/golang.org/x/tools@v0.36.0/cmd/godoc/handlers.go#L38) |
| * | /pkg/C/ | [`.gopath-loader/pkg/mod/golang.org/x/tools@v0.36.0/cmd/godoc/handlers.go:38`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/golang.org/x/tools@v0.36.0/cmd/godoc/handlers.go#L38) |
| * | /play.js | [`.gopath-loader/pkg/mod/golang.org/x/tools@v0.36.0/cmd/present/play.go:43`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/golang.org/x/tools@v0.36.0/cmd/present/play.go#L43) |
| * | /play.js | [`.gomod-cache/golang.org/x/tools@v0.36.0/cmd/present/play.go:43`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/golang.org/x/tools@v0.36.0/cmd/present/play.go#L43) |
| * | /progress-report | [`.gomod-cache/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go:52`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go#L52) |
| * | /progress-report | [`.gopath-loader/pkg/mod/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go:52`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go#L52) |
| * | /report-before-suite-completed | [`.gomod-cache/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go:55`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go#L55) |
| * | /report-before-suite-completed | [`.gopath-loader/pkg/mod/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go:55`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go#L55) |
| * | /report-before-suite-state | [`.gomod-cache/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go:56`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go#L56) |
| * | /report-before-suite-state | [`.gopath-loader/pkg/mod/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go:56`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go#L56) |
| * | /search | [`.gomod-cache/golang.org/x/tools@v0.36.0/godoc/pres.go:131`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/golang.org/x/tools@v0.36.0/godoc/pres.go#L131) |
| * | /search | [`.gopath-loader/pkg/mod/golang.org/x/tools@v0.36.0/godoc/pres.go:131`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/golang.org/x/tools@v0.36.0/godoc/pres.go#L131) |
| * | /select.json | [`.gomod-cache/golang.org/x/tools@v0.36.0/go/types/internal/play/play.go:49`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/golang.org/x/tools@v0.36.0/go/types/internal/play/play.go#L49) |
| * | /select.json | [`.gopath-loader/pkg/mod/golang.org/x/tools@v0.36.0/go/types/internal/play/play.go:49`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/golang.org/x/tools@v0.36.0/go/types/internal/play/play.go#L49) |
| * | /socket | [`.gomod-cache/golang.org/x/tools@v0.36.0/cmd/present/play.go:59`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/golang.org/x/tools@v0.36.0/cmd/present/play.go#L59) |
| * | /socket | [`.gopath-loader/pkg/mod/golang.org/x/tools@v0.36.0/cmd/present/play.go:59`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/golang.org/x/tools@v0.36.0/cmd/present/play.go#L59) |
| * | /src/pkg/ | [`.gopath-loader/pkg/mod/golang.org/x/tools@v0.36.0/godoc/redirect/redirect.go:21`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/golang.org/x/tools@v0.36.0/godoc/redirect/redirect.go#L21) |
| * | /src/pkg/ | [`.gomod-cache/golang.org/x/tools@v0.36.0/godoc/redirect/redirect.go:21`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/golang.org/x/tools@v0.36.0/godoc/redirect/redirect.go#L21) |
| * | /static/ | [`.gomod-cache/golang.org/x/tools@v0.36.0/cmd/present/main.go:98`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/golang.org/x/tools@v0.36.0/cmd/present/main.go#L98) |
| * | /static/ | [`.gopath-loader/pkg/mod/golang.org/x/tools@v0.36.0/cmd/present/main.go:98`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/golang.org/x/tools@v0.36.0/cmd/present/main.go#L98) |
| * | /suite-did-end | [`.gomod-cache/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go:50`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go#L50) |
| * | /suite-did-end | [`.gopath-loader/pkg/mod/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go:50`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go#L50) |
| * | /suite-will-begin | [`.gomod-cache/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go:48`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go#L48) |
| * | /suite-will-begin | [`.gopath-loader/pkg/mod/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go:48`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go#L48) |
| * | /ui/ | [`.gopath-loader/pkg/mod/github.com/google/pprof@v0.0.0-20250403155104-27863c87afa6/internal/driver/webui.go:211`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/github.com/google/pprof@v0.0.0-20250403155104-27863c87afa6/internal/driver/webui.go#L211) |
| * | /ui/ | [`.gomod-cache/github.com/google/pprof@v0.0.0-20250403155104-27863c87afa6/internal/driver/webui.go:211`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/github.com/google/pprof@v0.0.0-20250403155104-27863c87afa6/internal/driver/webui.go#L211) |
| * | /up | [`.gopath-loader/pkg/mod/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go:62`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go#L62) |
| * | /up | [`.gomod-cache/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go:62`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/github.com/onsi/ginkgo/v2@v2.25.3/internal/parallel_support/http_server.go#L62) |
| GET | /{user-id} | [`.gopath-loader/pkg/mod/github.com/emicklei/go-restful/v3@v3.12.2/doc.go:19`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/github.com/emicklei/go-restful/v3@v3.12.2/doc.go#L19) |
| GET | /{user-id} | [`.gomod-cache/github.com/emicklei/go-restful/v3@v3.12.2/doc.go:19`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/github.com/emicklei/go-restful/v3@v3.12.2/doc.go#L19) |
| GET | /{user-id} | [`.gopath-loader/pkg/mod/github.com/emicklei/go-restful/v3@v3.12.2/doc.go:83`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/github.com/emicklei/go-restful/v3@v3.12.2/doc.go#L83) |
| GET | /{user-id} | [`.gomod-cache/github.com/emicklei/go-restful/v3@v3.12.2/doc.go:83`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/github.com/emicklei/go-restful/v3@v3.12.2/doc.go#L83) |
| * | header | [`.gomod-cache/golang.org/x/net@v0.43.0/quic/qlog.go:267`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/golang.org/x/net@v0.43.0/quic/qlog.go#L267) |
| * | header | [`.gomod-cache/golang.org/x/net@v0.43.0/quic/qlog.go:211`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/golang.org/x/net@v0.43.0/quic/qlog.go#L211) |
| * | header | [`.gomod-cache/golang.org/x/net@v0.43.0/quic/qlog.go:165`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/golang.org/x/net@v0.43.0/quic/qlog.go#L165) |
| * | header | [`.gopath-loader/pkg/mod/golang.org/x/net@v0.43.0/quic/qlog.go:211`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/golang.org/x/net@v0.43.0/quic/qlog.go#L211) |
| * | header | [`.gomod-cache/golang.org/x/net@v0.43.0/quic/qlog.go:187`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/golang.org/x/net@v0.43.0/quic/qlog.go#L187) |
| * | header | [`.gopath-loader/pkg/mod/golang.org/x/net@v0.43.0/quic/qlog.go:187`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/golang.org/x/net@v0.43.0/quic/qlog.go#L187) |
| * | header | [`.gopath-loader/pkg/mod/golang.org/x/net@v0.43.0/quic/qlog.go:267`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/golang.org/x/net@v0.43.0/quic/qlog.go#L267) |
| * | header | [`.gopath-loader/pkg/mod/golang.org/x/net@v0.43.0/quic/qlog.go:165`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/golang.org/x/net@v0.43.0/quic/qlog.go#L165) |
| * | raw | [`.gomod-cache/golang.org/x/net@v0.43.0/quic/qlog.go:217`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/golang.org/x/net@v0.43.0/quic/qlog.go#L217) |
| * | raw | [`.gopath-loader/pkg/mod/golang.org/x/net@v0.43.0/quic/qlog.go:172`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/golang.org/x/net@v0.43.0/quic/qlog.go#L172) |
| * | raw | [`.gomod-cache/golang.org/x/net@v0.43.0/quic/qlog.go:193`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/golang.org/x/net@v0.43.0/quic/qlog.go#L193) |
| * | raw | [`.gopath-loader/pkg/mod/golang.org/x/net@v0.43.0/quic/qlog.go:193`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/golang.org/x/net@v0.43.0/quic/qlog.go#L193) |
| * | raw | [`.gomod-cache/golang.org/x/net@v0.43.0/quic/qlog.go:172`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/golang.org/x/net@v0.43.0/quic/qlog.go#L172) |
| * | raw | [`.gopath-loader/pkg/mod/golang.org/x/net@v0.43.0/quic/qlog.go:217`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/golang.org/x/net@v0.43.0/quic/qlog.go#L217) |
| * | vantage_point | [`.gopath-loader/pkg/mod/golang.org/x/net@v0.43.0/quic/qlog.go:96`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/golang.org/x/net@v0.43.0/quic/qlog.go#L96) |
| * | vantage_point | [`.gomod-cache/golang.org/x/net@v0.43.0/quic/qlog.go:96`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/golang.org/x/net@v0.43.0/quic/qlog.go#L96) |

## Configuration

ConfigMaps and Helm values that control this component's runtime behavior.

### Helm

**Chart:** jobset v0.10.1

