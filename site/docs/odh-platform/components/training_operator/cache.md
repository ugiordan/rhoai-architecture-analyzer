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

### Cache-Bypassed Types (DisableFor)

- corev1.ConfigMap
- corev1.Secret

### Issues

- Cache bypass (DisableFor) configured for corev1.ConfigMap. This is a common fix for OOM caused by informer cache flooding from high-cardinality types (e.g., opendatahub-io/model-registry-operator#457)
- Cache bypass (DisableFor) configured for corev1.Secret. This is a common fix for OOM caused by informer cache flooding from high-cardinality types (e.g., opendatahub-io/model-registry-operator#457)

