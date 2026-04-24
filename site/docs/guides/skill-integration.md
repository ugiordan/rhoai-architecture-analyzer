# Integrating Analyzer Output into AI Agent Skills

This guide explains how to consume architecture-analyzer output from Claude Code skills, agent workflows, or any AI-powered tooling that needs architectural context about a component.

## What the analyzer produces

The analyzer extracts three categories of output from a repository:

| File | What it contains | When to use it |
|------|-----------------|----------------|
| `component-architecture.json` | RBAC, CRDs, webhooks, controllers, external connections, feature gates, cache config, Dockerfiles, network policies, secrets, HTTP endpoints | Skill needs to understand what a component does, what it talks to, what permissions it has |
| `security-findings.json` | Architectural security findings (webhook gaps, RBAC issues, cache OOM risks, plaintext secrets) | Skill needs to know existing risks or audit results |
| `code-graph.json` | Raw code property graph: functions, call sites, HTTP handlers, DB operations, edges | Skill needs fine-grained code structure (rarely needed directly) |

For most skills, **`component-architecture.json` is the only file you need**.

## How to get the output

### Option 1: Pre-computed snapshots (recommended)

The CI pipeline runs weekly and publishes versioned snapshots. Fetch the latest for any component:

```bash
# From the snapshots branch
COMPONENT="opendatahub-operator"
VERSION="2026-04-24"  # or use latest tag

gh api repos/ugiordan/architecture-analyzer/contents/output/${VERSION}/odh/${COMPONENT}/component-architecture.json \
  --jq '.content' | base64 -d > arch-context.json
```

Or clone the snapshots branch:

```bash
git clone --branch snapshots --depth 1 \
  https://github.com/ugiordan/architecture-analyzer.git /tmp/arch-snapshots

# All component outputs are under output/<version>/<platform>/<component>/
ls /tmp/arch-snapshots/output/2026-04-24/odh/
```

### Option 2: On-demand analysis

If the skill already has the repo cloned (e.g., auto-bug-fix skills), run the analyzer directly:

```bash
# Install (one-time)
go install github.com/ugiordan/architecture-analyzer/cmd/arch-analyzer@latest

# Extract architecture (takes 2-10 seconds per repo)
arch-analyzer extract <repo-path> -output-dir /tmp/arch-context

# Or run full analysis (architecture + code graph + security scan)
arch-analyzer full-analysis <repo-path> -output-dir /tmp/arch-context
```

### Option 3: Embed in repo CI

Have each component's CI run the analyzer and upload the output as an artifact or commit to a known path:

```yaml
- name: Extract architecture context
  run: |
    go install github.com/ugiordan/architecture-analyzer/cmd/arch-analyzer@latest
    arch-analyzer extract . -output-dir .arch-context
- uses: actions/upload-artifact@v4
  with:
    name: arch-context
    path: .arch-context/
```

## Feeding output into a skill

### Prompt pattern

The JSON is self-describing with human-readable keys. No schema documentation is needed in the prompt. A minimal integration:

```
The following is the extracted architecture context for the {component_name}
component. Use it to understand the component's APIs, permissions, dependencies,
security posture, and infrastructure topology.

<architecture-context>
{contents of component-architecture.json}
</architecture-context>
```

### What each skill type gets from the context

**RFE creation and review:**

- CRDs: what custom resources the component owns, field counts, validation rules (CEL)
- Controller watches: what resources trigger reconciliation
- Dependencies: Go modules, internal cross-component references
- Feature gates: existing gating mechanisms and their default states

**Strategy creation:**

- RBAC: full permission surface (ClusterRoles, RoleBindings, kubebuilder markers)
- Network policies: ingress/egress rules, pod selectors
- Webhooks: admission controllers, failure policies, side effects
- External connections: databases, gRPC services, message queues, object storage
- Cache config: scope, filtered types, memory limits

**Security review and audit:**

- Security findings: pre-computed architectural issues with severity and evidence
- Secrets inventory: referenced secrets and where they're consumed (never values)
- Dockerfiles: base images, USER directives, FIPS indicators, exposed ports
- RBAC surface: who can do what to which resources

**Auto-bug-fix:**

- Controller watches: understand reconciliation triggers before changing controller logic
- CRD schemas: field structure and validation constraints
- External connections: service dependencies the fix must preserve
- Feature gates: whether the affected code path is gated

### Selective loading

For token-sensitive skills, load only the sections you need:

```bash
# Extract just RBAC and CRDs (using jq)
jq '{rbac, crds}' component-architecture.json > slim-context.json

# Extract just external connections and security findings
jq '{external_connections}' component-architecture.json > connections.json
```

### Multi-component context

For skills that operate across components (platform-level RFEs, cross-component strategies), use the aggregated platform output:

```bash
# Platform-level architecture (all components merged)
cat output/2026-04-24/odh/platform-architecture.json
```

This includes cross-component discovery: shared CRD groups, dependency graphs, network topology across components.

## Output format reference

See the [output format reference](../reference/output-format.md) for the complete JSON schema with examples of every section.

Key sections of `component-architecture.json`:

| Section | Type | Description |
|---------|------|-------------|
| `crds` | array | Custom Resource Definitions with group, version, kind, scope, field count, CEL rules |
| `rbac` | object | ClusterRoles, Roles, Bindings, kubebuilder RBAC markers |
| `services` | array | Kubernetes Services with ports, selectors |
| `deployments` | array | Deployments/StatefulSets with containers, security contexts, probes |
| `network_policies` | array | NetworkPolicy ingress/egress rules |
| `controller_watches` | object | Controllers with For/Owns/Watches GVKs |
| `dependencies` | object | Go modules, toolchain version, internal cross-references |
| `secrets_referenced` | array | Secret names and where they're consumed |
| `dockerfiles` | array | Base images, build stages, security indicators |
| `webhooks` | array | Admission webhooks with rules, failure policies |
| `configmaps` | array | ConfigMap names, data keys, references |
| `http_endpoints` | array | HTTP routes with method, path, handler function |
| `ingress_routing` | array | Gateway API, Istio, Ingress resources |
| `external_connections` | array | Database, gRPC, messaging, object storage connections |
| `feature_gates` | array | Feature gate names, defaults, pre-release stage |
| `cache_config` | object | Controller-runtime cache scope, filters, memory limits |

## Freshness and versioning

- Weekly snapshots are tagged with `arch/YYYY-MM-DD` in the snapshots branch
- Named version snapshots (e.g., `arch/v2.15.0`) are preserved indefinitely
- Date-stamped snapshots are pruned to the latest 10
- Each snapshot includes `snapshot-metadata.json` with the analyzer version and timestamp
- On-demand analysis always produces fresh results from the current repo state
