# Use Cases

Architecture Analyzer works with any Kubernetes/OpenShift operator ecosystem. Below are live platform analyses generated from real-world projects, demonstrating the tool's capabilities at scale.

## Available Platforms

### [Open Data Hub (upstream)](../odh-platform/index.md)

Analysis of the **opendatahub-io** GitHub organization: 29 components including operators, controllers, dashboards, and inference infrastructure. This is the upstream community project.

- 29 components, 35 CRDs, 95 cluster roles
- Full dependency graph, RBAC overlap detection, network topology

### [RHOAI (downstream)](../rhoai-platform/index.md)

Analysis of the **red-hat-data-services** GitHub organization: 31 components representing the downstream productized distribution. Includes additional components and downstream-specific configurations.

- 31 components, 48 CRDs, 91 cluster roles
- Separate CRD ownership chain, downstream-specific RBAC surface

## Analyzing your own platform

To generate a similar analysis for your platform:

```bash
# 1. Define repos in scan-config.yaml
# 2. Analyze each repo
for repo in $(yq '.repos[].name' scan-config.yaml); do
  ./scripts/analyze-repo.sh "$repo" results/
done

# 3. Aggregate into platform view
arch-analyzer aggregate results/ --output-dir platform-output/

# 4. Generate browsable docs
arch-analyzer docs --output-dir site/docs/my-platform \
  platform-output/platform-architecture.json
```

See the [Platform Aggregation](../guides/platform-aggregation.md) guide for details.
