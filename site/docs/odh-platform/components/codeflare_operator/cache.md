# codeflare-operator: Cache Architecture

Controller-runtime cache configuration controls which Kubernetes resources are cached in-memory. Misconfigured caches (cluster-wide watches on high-cardinality types without filters) are a primary cause of operator OOM kills.

## Cache Architecture

### Manager Configuration

| Property | Value |
|----------|-------|
| Manager file | `main.go` |
| Cache scope | cluster-wide |
| DefaultTransform | no |
| Memory limit | 1Gi |

### Issues

- No GOMEMLIMIT set in deployment (Go GC cannot pressure-tune)
- No cache configuration: all informers are cluster-wide (OOM risk)
- Type ClusterRoleBinding is watched but has no cache filter (cluster-wide informer)
- Type Ingress is watched but has no cache filter (cluster-wide informer)
- Type NetworkPolicy is watched but has no cache filter (cluster-wide informer)
- Type RayCluster is watched but has no cache filter (cluster-wide informer)
- Type Route is watched but has no cache filter (cluster-wide informer)
- Type Secret is watched but has no cache filter (cluster-wide informer)
- Type Service is watched but has no cache filter (cluster-wide informer)
- Type ServiceAccount is watched but has no cache filter (cluster-wide informer)

