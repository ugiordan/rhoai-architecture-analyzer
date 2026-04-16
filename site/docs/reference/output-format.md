# Output Format

## component-architecture.json

The core output format. Contains all data extracted by the 17 extractors.

### Top-level structure

```json
{
  "component": "my-operator",
  "repo": "github.com/org/my-operator",
  "extracted_at": "2026-04-14T10:30:00Z",
  "analyzer_version": "0.2.0",
  "crds": [],
  "rbac": {},
  "services": [],
  "deployments": [],
  "network_policies": [],
  "controller_watches": {},
  "dependencies": {},
  "secrets_referenced": [],
  "dockerfiles": [],
  "helm": {},
  "webhooks": [],
  "configmaps": [],
  "http_endpoints": [],
  "ingress_routing": [],
  "external_connections": [],
  "feature_gates": [],
  "cache_config": {}
}
```

### Key types

#### CRD

```json
{
  "group": "datasciencecluster.opendatahub.io",
  "version": "v1",
  "kind": "DataScienceCluster",
  "scope": "Cluster",
  "field_count": 42,
  "cel_rules": 3,
  "source_file": "config/crd/bases/datasciencecluster.yaml"
}
```

#### RBAC

```json
{
  "cluster_roles": [
    {
      "name": "manager-role",
      "rules": [
        {
          "api_groups": [""],
          "resources": ["secrets"],
          "verbs": ["get", "list", "watch"]
        }
      ],
      "source": "config/rbac/role.yaml"
    }
  ],
  "role_bindings": [],
  "kubebuilder_markers": []
}
```

#### Controller Watches

```json
{
  "controllers": [
    {
      "name": "DSCController",
      "file": "controllers/dsc_controller.go",
      "for": {
        "group": "datasciencecluster.opendatahub.io",
        "version": "v1",
        "kind": "DataScienceCluster"
      },
      "owns": [
        { "group": "apps", "version": "v1", "kind": "Deployment" }
      ],
      "watches": [
        { "group": "", "version": "v1", "kind": "ConfigMap" }
      ]
    }
  ]
}
```

#### Dependencies

```json
{
  "go_version": "1.25",
  "toolchain": "go1.25.0",
  "go_modules": [
    { "module": "sigs.k8s.io/controller-runtime", "version": "v0.23.3" }
  ],
  "replace_directives": [
    {
      "original": "github.com/org/old-module",
      "replacement": "github.com/org/new-module",
      "version": "v1.2.0"
    }
  ],
  "internal_odh": [
    {
      "component": "opendatahub-operator",
      "interaction": "Go module dependency: github.com/opendatahub-io/opendatahub-operator/v2"
    }
  ]
}
```

#### External Connections

```json
[
  {
    "type": "database",
    "service": "postgres",
    "target": "postgres://***@db.example.com:5432/mydb",
    "source": "pkg/storage/db.go:42",
    "function": "NewStore"
  },
  {
    "type": "messaging",
    "service": "kafka",
    "target": "",
    "source": "pkg/events/producer.go:18",
    "function": "InitProducer"
  }
]
```

#### Feature Gates

```json
[
  {
    "name": "PipelineReuse",
    "default": true,
    "pre_release": "Beta",
    "source": "pkg/features/gates.go:15"
  },
  {
    "name": "ExperimentalAPI",
    "default": false,
    "pre_release": "Alpha",
    "source": "pkg/features/gates.go:16"
  },
  {
    "name": "DebugMode",
    "default": true,
    "source": "cmd/main.go:42",
    "runtime_set": true
  }
]
```

#### Cache Config

```json
{
  "scope": "cluster",
  "filtered_types": ["ConfigMap", "Secret"],
  "disabled_types": [],
  "implicit_informers": [
    {
      "type": "Namespace",
      "source": "controllers/dsc_controller.go:145",
      "reason": "client.Get call for unwatched type"
    }
  ],
  "gomemlimit": "512MiB",
  "container_memory_limit": "1Gi",
  "default_transform": false,
  "findings": [
    {
      "severity": "warning",
      "message": "Missing DefaultTransform - managedFields consuming extra memory",
      "recommendation": "Add cache.DefaultTransform to strip managedFields"
    }
  ]
}
```

## Security findings (JSON)

```json
{
  "domain": "security",
  "findings": [
    {
      "id": "CGA-004-001",
      "query": "CGA-004",
      "title": "Hardcoded API key",
      "severity": "high",
      "file": "pkg/config/defaults.go",
      "line": 23,
      "evidence": "String literal matches API key pattern: 'sk-...'",
      "recommendation": "Use environment variable or secret mount"
    }
  ]
}
```

## Security findings (SARIF)

Standard SARIF 2.1.0 format compatible with GitHub Code Scanning:

```json
{
  "$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
  "version": "2.1.0",
  "runs": [
    {
      "tool": {
        "driver": {
          "name": "rhoai-architecture-analyzer",
          "version": "0.2.0",
          "rules": [...]
        }
      },
      "results": [...]
    }
  ]
}
```

## Platform aggregation output

```json
{
  "components": ["repo-a", "repo-b", "repo-c"],
  "aggregated_at": "2026-04-14T10:30:00Z",
  "crd_ownership": {},
  "cross_dependencies": [],
  "rbac_overlap": [],
  "network_mesh": []
}
```
