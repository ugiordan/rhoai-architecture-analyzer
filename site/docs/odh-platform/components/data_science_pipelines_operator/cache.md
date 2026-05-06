# data-science-pipelines-operator: Cache Architecture

Controller-runtime cache configuration controls which Kubernetes resources are cached in-memory. Misconfigured caches (cluster-wide watches on high-cardinality types without filters) are a primary cause of operator OOM kills.

## Cache Architecture

### Manager Configuration

| Property | Value |
|----------|-------|
| Manager file | `main.go` |
| Cache scope | cluster-wide |
| DefaultTransform | yes |
| GOMEMLIMIT | 3600MiB |
| Memory limit | 4Gi |

### Filtered Types

| Type | Filter Kind | Filter |
|------|-------------|--------|
| admv1.MutatingWebhookConfiguration | field | metadata.name |
| admv1.ValidatingWebhookConfiguration | field | metadata.name |
| corev1.PersistentVolumeClaim | label | label selector |
| corev1.Pod | label | config.DSPComponentk8sLabel=config.DSPComponentk8sLabelValue (constants, resolved at runtime) |
| corev1.Service | label | label selector |
| corev1.ServiceAccount | label | label selector |
| rbacv1.Role | label | label selector |
| rbacv1.RoleBinding | label | label selector |
| routev1.Route | label | label selector |

### Cache-Bypassed Types (DisableFor)

- corev1.ConfigMap
- corev1.Secret

### Issues

- Cache bypass (DisableFor) configured for corev1.ConfigMap. This is a common fix for OOM caused by informer cache flooding from high-cardinality types (e.g., opendatahub-io/model-registry-operator#457)
- Cache bypass (DisableFor) configured for corev1.Secret. This is a common fix for OOM caused by informer cache flooding from high-cardinality types (e.g., opendatahub-io/model-registry-operator#457)
- Type DataSciencePipelinesApplication is watched but has no cache filter (cluster-wide informer)
- Type Deployment is watched but has no cache filter (cluster-wide informer)
- Type NetworkPolicy is watched but has no cache filter (cluster-wide informer)

