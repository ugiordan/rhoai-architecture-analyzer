# SBOM Generation & Image Reporting

Architecture Analyzer can generate CycloneDX SBOMs and comprehensive image/container analysis reports from the same extraction data used for architecture diagrams. No additional scanning or extraction is needed.

## Software Bill of Materials (SBOM)

The `sbom` command produces a [CycloneDX 1.5](https://cyclonedx.org/) JSON file from the extracted `component-architecture.json`.

### What's Included

The SBOM catalogs every dependency detected during extraction:

- **Go modules** from `go.mod` with exact versions and Package URLs
- **Python packages** from `requirements.txt` / `pyproject.toml`
- **Dockerfile base images** with tags, SHA-256 digests, multi-arch info, build stages, and security issues
- **Deployment container images** with operational metadata:
    - Security context (runAsNonRoot, readOnlyRootFilesystem, privileged, capabilities)
    - Resource limits and requests (CPU, memory)
    - Health probes (liveness, readiness endpoints)
- **Operator image constants** (Go const values referencing container images)

### Usage

```bash
# Generate SBOM from existing extraction
arch-analyzer sbom component-architecture.json --output sbom.json

# View component count
cat sbom.json | jq '.components | length'

# Filter by type
cat sbom.json | jq '[.components[] | select(.type == "container")]'
```

### Example Output

```json
{
  "bomFormat": "CycloneDX",
  "specVersion": "1.5",
  "components": [
    {
      "type": "library",
      "name": "github.com/Azure/azure-sdk-for-go/sdk/azcore",
      "version": "v1.19.1",
      "purl": "pkg:golang/github.com/Azure/azure-sdk-for-go/sdk/azcore@v1.19.1",
      "scope": "required"
    },
    {
      "type": "container",
      "name": "registry.access.redhat.com/ubi9/ubi-minimal",
      "version": "latest",
      "purl": "pkg:docker/registry.access.redhat.com/ubi9/ubi-minimal@latest",
      "properties": [
        { "name": "arch-analyzer:type", "value": "dockerfile-base" },
        { "name": "arch-analyzer:stages", "value": "2" },
        { "name": "arch-analyzer:user", "value": "1001" }
      ]
    }
  ]
}
```

### Pre-generated SBOMs

SBOMs for all 42 analyzed RHOAI components are available in the repository:

- [All SBOMs](https://github.com/ugiordan/architecture-analyzer/tree/main/docs/sboms)
- [KServe (64 deps)](https://github.com/ugiordan/architecture-analyzer/blob/main/docs/sboms/kserve-sbom.json)
- [ODH Dashboard (131 deps)](https://github.com/ugiordan/architecture-analyzer/blob/main/docs/sboms/odh-dashboard-sbom.json)
- [Data Science Pipelines (104 deps)](https://github.com/ugiordan/architecture-analyzer/blob/main/docs/sboms/data-science-pipelines-sbom.json)

## Image & Container Analysis Report

The `report` command generates a comprehensive markdown report covering security posture, GPU dependencies, resource configuration, and deployment issues.

### Report Sections

| Section | What It Covers |
|---------|---------------|
| **GPU / CUDA Dependencies** | Components requiring NVIDIA CUDA, Intel Gaudi (HPU), AMD ROCm. Detects CUDA versions from Dockerfile paths and base images. |
| **Base Image Registry Distribution** | Which container registries are used across all Dockerfiles (Red Hat UBI, Distroless, Docker Hub, etc.) |
| **Multi-Architecture Support** | Dockerfiles declaring multi-arch builds |
| **Dockerfile Security Issues** | Unpinned base images, missing USER directive, running as root, etc. |
| **Container Security Contexts** | RunAsNonRoot, ReadOnlyRootFilesystem, Privileged mode, DROP ALL capabilities |
| **Resource Limits** | CPU and memory requests/limits per container |
| **Health Probes** | Liveness and readiness probe coverage and endpoints |
| **Sidecar Containers** | Deployments with kube-rbac-proxy, oauth-proxy, and other sidecars |
| **Deployment Issues** | Missing PodDisruptionBudget, HorizontalPodAutoscaler |
| **Operator Image Constants** | Go constants defining default container images |

### Usage

```bash
# Single component report
arch-analyzer report component-architecture.json --output report.md

# Cross-component analysis (pass multiple JSON files)
arch-analyzer report results/*/component-architecture.json --output platform-report.md
```

### Pre-generated Report

The full platform report covering all 42 RHOAI components is available:

- [Image Dependencies Report](https://github.com/ugiordan/architecture-analyzer/blob/main/docs/reports/image-dependencies.md)

## Integration with CI

Both commands can be integrated into CI pipelines:

```yaml
# GitHub Actions example
- name: Extract architecture
  run: arch-analyzer extract . --output component-architecture.json

- name: Generate SBOM
  run: arch-analyzer sbom component-architecture.json --output sbom.json

- name: Generate report
  run: arch-analyzer report component-architecture.json --output report.md

- name: Upload artifacts
  uses: actions/upload-artifact@v4
  with:
    name: architecture-analysis
    path: |
      component-architecture.json
      sbom.json
      report.md
```
