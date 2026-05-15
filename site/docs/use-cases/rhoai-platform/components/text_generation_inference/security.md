# text-generation-inference: Security

## Secrets

Kubernetes secrets referenced by this component. Only names and types are shown, not values.

## Deployment Security Controls

SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.

### Container Security Contexts

| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |
|------------|-----------|--------------|------------|------------|--------|
| inference-server | server | true | ? | false | [`deployment/base/deployment.yaml`](https://github.com/red-hat-data-services/text-generation-inference/blob/fded01861025fff09ba5f9a49cda710fcfd3ca93/deployment/base/deployment.yaml) |

## Build Security

Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.

| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |
|------|------------|--------|------|-------|---------------|------|--------|
| `Dockerfile` | base | 14 | tgis |  |  |  | Unpinned base image: base; Unpinned base image: cuda-base; Unpinned base image: rust-builder; Unpinned base image: rust-builder; Unpinned base image: base; Unpinned base image: test-base; Unpinned base image: cuda-devel; Unpinned base image: python-builder; Unpinned base image: python-builder; Unpinned base image: base; Unpinned base image: python-builder; Unpinned base image: base |

