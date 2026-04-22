# opendatahub-operator: Cache Architecture

Controller-runtime cache configuration controls which Kubernetes resources are cached in-memory. Misconfigured caches (cluster-wide watches on high-cardinality types without filters) are a primary cause of operator OOM kills.

## Cache Architecture

### Manager Configuration

| Property | Value |
|----------|-------|
| Manager file | `main.go` |
| Cache scope | cluster-wide |
| DefaultTransform | no |
| Memory limit | 4Gi |

### Issues

- No GOMEMLIMIT set in deployment (Go GC cannot pressure-tune)
- No cache configuration: all informers are cluster-wide (OOM risk)
- Type ClusterRole is watched but has no cache filter (cluster-wide informer)
- Type ClusterRoleBinding is watched but has no cache filter (cluster-wide informer)
- Type ConfigMap is watched but has no cache filter (cluster-wide informer)
- Type DSCInitialization is watched but has no cache filter (cluster-wide informer)
- Type DataScienceCluster is watched but has no cache filter (cluster-wide informer)
- Type Deployment is watched but has no cache filter (cluster-wide informer)
- Type Namespace is watched but has no cache filter (cluster-wide informer)
- Type NetworkPolicy is watched but has no cache filter (cluster-wide informer)
- Type Pod is watched but has no cache filter (cluster-wide informer)
- Type ReplicaSet is watched but has no cache filter (cluster-wide informer)
- Type Role is watched but has no cache filter (cluster-wide informer)
- Type RoleBinding is watched but has no cache filter (cluster-wide informer)
- Type Secret is watched but has no cache filter (cluster-wide informer)
- Type Service is watched but has no cache filter (cluster-wide informer)
- Type ServiceAccount is watched but has no cache filter (cluster-wide informer)

