"""Extract secret references from Kubernetes manifests (names only, never values)."""

from __future__ import annotations

import logging
from typing import Any

from extractors.base import BaseExtractor

logger = logging.getLogger(__name__)

DEPLOYMENT_PATTERNS = [
    "**/deployment.yaml",
    "**/deployment*.yaml",
    "**/statefulset.yaml",
    "**/manager*.yaml",
    "charts/**/templates/deployment*.yaml",
]

SERVICE_PATTERNS = [
    "**/service.yaml",
    "**/service*.yaml",
]


class SecretExtractor(BaseExtractor):
    """Extract references to Kubernetes Secrets (names and types, never values)."""

    def extract(self) -> dict[str, Any]:
        secrets_map: dict[str, dict[str, Any]] = {}

        # Scan deployments for secret references
        dep_files = self.find_yaml_files(DEPLOYMENT_PATTERNS)
        for fpath in dep_files:
            for doc in self.parse_yaml_safe(fpath):
                kind = doc.get("kind", "")
                if kind not in ("Deployment", "StatefulSet"):
                    continue
                metadata = doc.get("metadata", {})
                dep_name = metadata.get("name", "") if isinstance(metadata, dict) else ""
                ref_label = f"{kind.lower()}/{dep_name}"

                spec = doc.get("spec", {})
                if not isinstance(spec, dict):
                    continue
                template = spec.get("template", {})
                if not isinstance(template, dict):
                    continue
                pod_spec = template.get("spec", {})
                if not isinstance(pod_spec, dict):
                    continue

                # Volumes with secrets
                for vol in pod_spec.get("volumes", []):
                    if not isinstance(vol, dict):
                        continue
                    secret = vol.get("secret", {})
                    if isinstance(secret, dict) and "secretName" in secret:
                        secret_name = secret["secretName"]
                        if isinstance(secret_name, str):
                            self._add_secret(
                                secrets_map, secret_name, ref_label, "volume-mounted"
                            )

                # Containers env and envFrom
                all_containers = list(pod_spec.get("containers", [])) + list(
                    pod_spec.get("initContainers", [])
                )
                for container in all_containers:
                    if not isinstance(container, dict):
                        continue
                    # env[].valueFrom.secretKeyRef
                    for env_var in container.get("env", []):
                        if not isinstance(env_var, dict):
                            continue
                        value_from = env_var.get("valueFrom", {})
                        if isinstance(value_from, dict):
                            secret_ref = value_from.get("secretKeyRef", {})
                            if isinstance(secret_ref, dict) and "name" in secret_ref:
                                name = secret_ref["name"]
                                if isinstance(name, str):
                                    self._add_secret(
                                        secrets_map, name, ref_label, "env-var"
                                    )
                    # envFrom[].secretRef
                    for env_from in container.get("envFrom", []):
                        if not isinstance(env_from, dict):
                            continue
                        secret_ref = env_from.get("secretRef", {})
                        if isinstance(secret_ref, dict) and "name" in secret_ref:
                            name = secret_ref["name"]
                            if isinstance(name, str):
                                self._add_secret(
                                    secrets_map, name, ref_label, "envFrom"
                                )

        # Scan services for cert annotations
        svc_files = self.find_yaml_files(SERVICE_PATTERNS)
        for fpath in svc_files:
            for doc in self.parse_yaml_safe(fpath):
                if doc.get("kind") != "Service":
                    continue
                metadata = doc.get("metadata", {})
                if not isinstance(metadata, dict):
                    continue
                annotations = metadata.get("annotations", {})
                if not isinstance(annotations, dict):
                    continue
                cert_secret = annotations.get(
                    "service.beta.openshift.io/serving-cert-secret-name"
                )
                if isinstance(cert_secret, str):
                    svc_name = metadata.get("name", "")
                    self._add_secret(
                        secrets_map,
                        cert_secret,
                        f"service/{svc_name}",
                        "OpenShift serving cert",
                        secret_type="kubernetes.io/tls",
                    )

        return {
            "secrets_referenced": [
                {
                    "name": name,
                    "type": info.get("type", "Opaque"),
                    "referenced_by": info["referenced_by"],
                    "provisioned_by": info.get("provisioned_by", ""),
                }
                for name, info in secrets_map.items()
            ]
        }

    @staticmethod
    def _add_secret(
        secrets_map: dict[str, dict[str, Any]],
        name: str,
        referenced_by: str,
        provisioned_by: str,
        secret_type: str = "Opaque",
    ) -> None:
        """Add or update a secret reference in the map."""
        if name not in secrets_map:
            secrets_map[name] = {
                "type": secret_type,
                "referenced_by": [],
                "provisioned_by": provisioned_by,
            }
        if referenced_by not in secrets_map[name]["referenced_by"]:
            secrets_map[name]["referenced_by"].append(referenced_by)
        # Upgrade type if we have more specific info
        if secret_type != "Opaque":
            secrets_map[name]["type"] = secret_type
