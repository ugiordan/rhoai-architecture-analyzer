# Renderers

The analyzer includes 7 renderers that produce visualizations and reports from extracted architecture data.

## Renderer reference

| Renderer | Output File | Format | Description |
|----------|-------------|--------|-------------|
| RBAC | `rbac.mmd` | Mermaid | ServiceAccounts -> Bindings -> Roles -> Resources |
| Component | `component.mmd` | Mermaid | CRDs watched/owned, dependency relationships |
| Security/Network | `security-network.txt` | ASCII | Layered view: network, RBAC, secrets, security contexts |
| Dependencies | `dependencies.mmd` | Mermaid | Go module dependencies (internal ODH highlighted) |
| C4 | `c4-context.dsl` | Structurizr DSL | C4 context diagram |
| Dataflow | `dataflow.mmd` | Mermaid | Controller watches and service connections |
| Report | `report.md` | Markdown | Structured tables for all extracted data |

## Selecting renderers

```bash
# All renderers (default)
rhoai-analyzer render component-architecture.json --output-dir diagrams/

# Specific renderers
rhoai-analyzer render component-architecture.json --formats rbac,component,report
```

## RBAC renderer

Produces a Mermaid graph showing the RBAC chain:

```
ServiceAccount --> ClusterRoleBinding --> ClusterRole --> resources
```

Includes:

- All ServiceAccounts referenced in bindings
- All ClusterRoles/Roles with their rules
- Resource types with verbs (get, list, watch, create, update, delete)
- Kubebuilder-generated RBAC markers

## Component renderer

Produces a Mermaid diagram showing:

- The component and its owned CRDs
- Controller watch relationships (For/Owns/Watches)
- Dependencies on other ODH components
- Webhook registrations

## Security/Network renderer

Produces an ASCII layered view:

```
=== Network Layer ===
NetworkPolicy: allow-webhook (ingress from kube-apiserver on 9443)
NetworkPolicy: allow-metrics (ingress from monitoring on 8443)

=== RBAC Layer ===
ClusterRole: manager-role
  - secrets: get, list, watch
  - deployments: get, list, watch, create, update, patch, delete

=== Secrets Layer ===
Secret: webhook-certs (referenced by deployment/manager)

=== Security Context ===
Container: manager
  runAsNonRoot: true
  allowPrivilegeEscalation: false
```

## Dependencies renderer

Produces a Mermaid dependency graph:

- Direct Go module dependencies
- Internal ODH dependencies highlighted with a distinct color
- Replace directives shown as dashed lines

## C4 renderer

Produces a Structurizr C4 context diagram in DSL format:

```
workspace {
    model {
        component = softwareSystem "my-operator" {
            description "Manages ML workloads"
        }
        kubernetes = softwareSystem "Kubernetes API"
        component -> kubernetes "watches CRDs, manages resources"
    }
    views {
        systemContext component "Context" {
            include *
        }
    }
}
```

Load into [Structurizr](https://structurizr.com/) for rendering.

## Dataflow renderer

Produces a Mermaid sequence diagram showing:

- Controller watch setup (what each controller watches)
- Service connections between components
- Webhook call flows

## Report renderer

Produces a comprehensive markdown report with tables for:

- All extracted CRDs, RBAC rules, services, deployments
- Network policies and ingress routes
- Controller watches with GVK resolution
- Cache analysis findings with severity
- Webhook configurations
- Dependency list

## Platform renderers

Additional renderers for aggregated platform output:

- `platform.go`: Cross-component topology diagram
- `platform_report.go`: Platform-level report with cross-repo findings
