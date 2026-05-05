# mlflow-operator: Cache Architecture

Controller-runtime cache configuration controls which Kubernetes resources are cached in-memory. Misconfigured caches (cluster-wide watches on high-cardinality types without filters) are a primary cause of operator OOM kills.

## Cache Architecture

### Manager Configuration

| Property | Value |
|----------|-------|
| Manager file | `cmd/main.go` |
| Cache scope | namespace-scoped |
| DefaultTransform | no |
| Memory limit | 4Gi |

### Filtered Types

| Type | Filter Kind | Filter |
|------|-------------|--------|
| appsv1.Deployment | label | label selector |
| corev1.PersistentVolumeClaim | label | label selector |
| corev1.Secret | label | label selector |
| corev1.Service | label | label selector |
| corev1.ServiceAccount | label | label selector |
| rbacv1.ClusterRole | label | label selector |
| rbacv1.ClusterRoleBinding | label | label selector |

### Issues

- No DefaultTransform: managedFields cached for all objects (wasted memory). Add cache.DefaultTransform to strip managedFields and reduce memory footprint
- No GOMEMLIMIT set in deployment (Go GC cannot pressure-tune). Set GOMEMLIMIT to 80-90% of container memory limit for optimal GC behavior
- Type ConsoleLink is watched but has no cache filter (cluster-wide informer)
- Type HTTPRoute is watched but has no cache filter (cluster-wide informer)
- Type MLflow is watched but has no cache filter (cluster-wide informer)
- Type ServiceMonitor is watched but has no cache filter (cluster-wide informer)

