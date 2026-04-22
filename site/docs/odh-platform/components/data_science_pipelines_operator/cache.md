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
| admv1.MutatingWebhookConfiguration | field | field selector |
| admv1.ValidatingWebhookConfiguration | field | field selector |
| corev1.PersistentVolumeClaim | label | label selector |
| corev1.Pod | label | label selector |
| corev1.Service | label | label selector |
| corev1.ServiceAccount | label | label selector |
| rbacv1.Role | label | label selector |
| rbacv1.RoleBinding | label | label selector |
| routev1.Route | label | label selector |

### Transform-Stripped Types

- corev1.ConfigMap
- corev1.Secret

### Issues

- Type ConfigMap is watched but has no cache filter (cluster-wide informer)
- Type DataSciencePipelinesApplication is watched but has no cache filter (cluster-wide informer)
- Type Deployment is watched but has no cache filter (cluster-wide informer)
- Type NetworkPolicy is watched but has no cache filter (cluster-wide informer)
- Type Secret is watched but has no cache filter (cluster-wide informer)

