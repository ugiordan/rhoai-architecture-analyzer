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

### Implicit Informers (OOM Risk)

| Type | Source | Risk |
|------|--------|------|
| client.ListOptions | `internal/controller/configmap_bootstrap.go:52` | medium |

### Issues

- No GOMEMLIMIT set in deployment (Go GC cannot pressure-tune)
- No cache configuration: all informers are cluster-wide (OOM risk)
- Type ConfigMap is watched but has no cache filter (cluster-wide informer)
- Type InferencePool is watched but has no cache filter (cluster-wide informer)
- Type VariantAutoscaling is watched but has no cache filter (cluster-wide informer)

