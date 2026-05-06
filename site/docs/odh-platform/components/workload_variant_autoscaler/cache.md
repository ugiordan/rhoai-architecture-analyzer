# workload-variant-autoscaler: Cache Architecture

Controller-runtime cache configuration controls which Kubernetes resources are cached in-memory. Misconfigured caches (cluster-wide watches on high-cardinality types without filters) are a primary cause of operator OOM kills.

## Cache Architecture

### Manager Configuration

| Property | Value |
|----------|-------|
| Manager file | `cmd/main.go` |
| Cache scope | namespace-scoped |
| DefaultTransform | no |
| Memory limit | 1Gi |

### Filtered Types

| Type | Filter Kind | Filter |
|------|-------------|--------|
| corev1.ConfigMap | label | label selector |

### Implicit Informers (OOM Risk)

| Type | Source | Risk |
|------|--------|------|
| client.ListOptions | `internal/controller/configmap_bootstrap.go:53` | medium |

### Issues

- No DefaultTransform: managedFields cached for all objects (wasted memory). Add cache.DefaultTransform to strip managedFields and reduce memory footprint
- No GOMEMLIMIT set in deployment (Go GC cannot pressure-tune). Set GOMEMLIMIT to 80-90% of container memory limit for optimal GC behavior
- Type InferencePool is watched but has no cache filter (cluster-wide informer)
- Type VariantAutoscaling is watched but has no cache filter (cluster-wide informer)

