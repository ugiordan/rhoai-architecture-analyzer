# model-registry: Cache Architecture

Controller-runtime cache configuration controls which Kubernetes resources are cached in-memory. Misconfigured caches (cluster-wide watches on high-cardinality types without filters) are a primary cause of operator OOM kills.

## Cache Architecture

### Manager Configuration

| Property | Value |
|----------|-------|
| Manager file | `cmd/controller/main.go` |
| Cache scope | cluster-wide |
| DefaultTransform | no |
| Memory limit | 128Mi |

### Issues

- No GOMEMLIMIT set in deployment (Go GC cannot pressure-tune)
- No cache configuration: all informers are cluster-wide (OOM risk)
- Type InferenceService is watched but has no cache filter (cluster-wide informer)

