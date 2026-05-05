# opendatahub-operator: Cache Architecture

Controller-runtime cache configuration controls which Kubernetes resources are cached in-memory. Misconfigured caches (cluster-wide watches on high-cardinality types without filters) are a primary cause of operator OOM kills.

## Cache Architecture

### Manager Configuration

| Property | Value |
|----------|-------|
| Manager file | `cmd/main.go` |
| Cache scope | namespace-scoped |
| DefaultTransform | yes |
| Memory limit | 4Gi |

### Filtered Types

| Type | Filter Kind | Filter |
|------|-------------|--------|
| appsv1.Deployment | namespace | namespace-scoped |
| corev1.ConfigMap | namespace | namespace-scoped |
| corev1.Secret | namespace | namespace-scoped |
| extv1.CustomResourceDefinition | label | label selector |
| networkingv1.NetworkPolicy | namespace | namespace-scoped |
| rbacv1.ClusterRole | label | label selector |
| rbacv1.ClusterRoleBinding | label | label selector |
| rbacv1.Role | namespace | namespace-scoped |
| rbacv1.RoleBinding | namespace | namespace-scoped |

### Cache-Bypassed Types (DisableFor)

- authorizationv1.SelfSubjectRulesReview
- corev1.Pod
- ofapiv1alpha1.CatalogSource
- ofapiv1alpha1.Subscription
- userv1.Group

### Issues

- No GOMEMLIMIT set in deployment (Go GC cannot pressure-tune). Set GOMEMLIMIT to 80-90% of container memory limit for optimal GC behavior
- Type Auth is watched but has no cache filter (cluster-wide informer)
- Type ConsoleLink is watched but has no cache filter (cluster-wide informer)
- Type Dashboard is watched but has no cache filter (cluster-wide informer)
- Type DataSciencePipelines is watched but has no cache filter (cluster-wide informer)
- Type FeastOperator is watched but has no cache filter (cluster-wide informer)
- Type HTTPRoute is watched but has no cache filter (cluster-wide informer)
- Type Job is watched but has no cache filter (cluster-wide informer)
- Type Kserve is watched but has no cache filter (cluster-wide informer)
- Type Kueue is watched but has no cache filter (cluster-wide informer)
- Type LlamaStackOperator is watched but has no cache filter (cluster-wide informer)
- Type MLflowOperator is watched but has no cache filter (cluster-wide informer)
- Type ModelController is watched but has no cache filter (cluster-wide informer)
- Type ModelRegistry is watched but has no cache filter (cluster-wide informer)
- Type MutatingWebhookConfiguration is watched but has no cache filter (cluster-wide informer)
- Type Namespace is watched but has no cache filter (cluster-wide informer)
- Type PodDisruptionBudget is watched but has no cache filter (cluster-wide informer)
- Type PodMonitor is watched but has no cache filter (cluster-wide informer)
- Type PrometheusRule is watched but has no cache filter (cluster-wide informer)
- Type Ray is watched but has no cache filter (cluster-wide informer)
- Type Route is watched but has no cache filter (cluster-wide informer)
- Type SecurityContextConstraints is watched but has no cache filter (cluster-wide informer)
- Type Service is watched but has no cache filter (cluster-wide informer)
- Type ServiceAccount is watched but has no cache filter (cluster-wide informer)
- Type ServiceMonitor is watched but has no cache filter (cluster-wide informer)
- Type SparkOperator is watched but has no cache filter (cluster-wide informer)
- Type Template is watched but has no cache filter (cluster-wide informer)
- Type Trainer is watched but has no cache filter (cluster-wide informer)
- Type TrainingOperator is watched but has no cache filter (cluster-wide informer)
- Type TrustyAI is watched but has no cache filter (cluster-wide informer)
- Type ValidatingAdmissionPolicy is watched but has no cache filter (cluster-wide informer)
- Type ValidatingAdmissionPolicyBinding is watched but has no cache filter (cluster-wide informer)
- Type ValidatingWebhookConfiguration is watched but has no cache filter (cluster-wide informer)
- Type Workbenches is watched but has no cache filter (cluster-wide informer)

