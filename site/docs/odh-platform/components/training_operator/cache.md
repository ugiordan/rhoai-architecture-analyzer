# training-operator: Cache Architecture

Controller-runtime cache configuration controls which Kubernetes resources are cached in-memory. Misconfigured caches (cluster-wide watches on high-cardinality types without filters) are a primary cause of operator OOM kills.

## Cache Architecture

### Manager Configuration

| Property | Value |
|----------|-------|
| Manager file | `cmd/training-operator.v1/main.go` |
| Cache scope | cluster-wide |
| DefaultTransform | yes |
| GOMEMLIMIT | 460MiB |
| Memory limit | 512Mi |

### Filtered Types

| Type | Filter Kind | Filter |
|------|-------------|--------|
| corev1.ConfigMap | label | label selector |
| corev1.Pod | label | label selector |
| corev1.Secret | label | label selector |
| corev1.Service | label | label selector |
| corev1.ServiceAccount | label | label selector |
| rbacv1.Role | label | label selector |
| rbacv1.RoleBinding | label | label selector |

