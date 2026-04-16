# RHOAI Architecture Model

> This document describes RHOAI component interactions, data flows, and design principles.
> It serves as both a human-readable reference and AI context for architecture review.

## Component Overview

| Component | Role | Primary CRDs |
|-----------|------|-------------|
| opendatahub-operator | Platform lifecycle management | DataScienceCluster, DSCInitialization |
| kserve | Model serving inference runtime | InferenceService, ServingRuntime |
| odh-model-controller | Model serving orchestration | (watches InferenceService, ServingRuntime) |
| odh-dashboard | Web UI for RHOAI | (reads DSC, DSCI, InferenceService) |
| model-registry-operator | ML model registry | ModelRegistry |
| data-science-pipelines-operator | ML pipeline orchestration | DataSciencePipelinesApplication |
| training-operator | Distributed training jobs | (TODO) |
| trustyai-service-operator | AI explainability | (TODO) |

## Component Interactions

<!-- Each component team should fill in their section -->

### opendatahub-operator
- **Provides:** DataScienceCluster (v1), DSCInitialization (v1)
- **Consumes:** Component CRDs for lifecycle management
- **Behavioral contracts:** Reconciles DSC to manage component installations. Owns namespace-scoped resources.

### kserve
- **Provides:** InferenceService (v1beta1), ServingRuntime (v1alpha1)
- **Consumes:** None (standalone)
- **Behavioral contracts:** TODO

<!-- Add more components as needed -->

## Design Principles

1. **CRD Versioning:** New fields must be optional. Breaking changes require a version bump.
2. **Namespace Isolation:** Components should not cross namespace boundaries without explicit contracts.
3. **RBAC Boundaries:** Each component should use its own ServiceAccount with minimal required permissions.
4. **Label Conventions:** TODO

## Known Integration Points

- **DSC Lifecycle:** opendatahub-operator reconciles DSC to install/remove components
- **Dashboard API Reads:** odh-dashboard reads DSC, DSCI, InferenceService via API
- **Model Controller Write Paths:** odh-model-controller creates and manages InferenceService CRs
