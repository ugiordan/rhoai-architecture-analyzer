# model-registry-operator: Cache Architecture

Controller-runtime cache configuration controls which Kubernetes resources are cached in-memory. Misconfigured caches (cluster-wide watches on high-cardinality types without filters) are a primary cause of operator OOM kills.

## Cache Architecture

### Manager Configuration

| Property | Value |
|----------|-------|
| Manager file | `cmd/main.go` |
| Cache scope | cluster-wide |
| DefaultTransform | no |
| Memory limit | 512Mi |

### Filtered Types

| Type | Filter Kind | Filter |
|------|-------------|--------|
| appsv1.Deployment | label | app.kubernetes.io/created-by=model-registry-operator |
| corev1.ConfigMap | namespace | namespace-scoped |
| corev1.PersistentVolumeClaim | label | app.kubernetes.io/created-by=model-registry-operator |
| corev1.Service | label | app.kubernetes.io/created-by=model-registry-operator |
| corev1.ServiceAccount | label | app.kubernetes.io/created-by=model-registry-operator |
| networkingv1.NetworkPolicy | label | app.kubernetes.io/created-by=model-registry-operator |
| rbacv1.ClusterRoleBinding | label | app.kubernetes.io/created-by=model-registry-operator |
| rbacv1.Role | label | app.kubernetes.io/created-by=model-registry-operator |
| rbacv1.RoleBinding | label | app.kubernetes.io/created-by=model-registry-operator |

### Cache-Bypassed Types (DisableFor)

- corev1.Secret

### Issues

- Cache bypass (DisableFor) configured for corev1.Secret. This is a common fix for OOM caused by informer cache flooding from high-cardinality types (e.g., opendatahub-io/model-registry-operator#457)
- No DefaultTransform: managedFields cached for all objects (wasted memory). Add cache.DefaultTransform to strip managedFields and reduce memory footprint
- No GOMEMLIMIT set in deployment (Go GC cannot pressure-tune). Set GOMEMLIMIT to 80-90% of container memory limit for optimal GC behavior
- Type ModelRegistry is watched but has no cache filter (cluster-wide informer)

