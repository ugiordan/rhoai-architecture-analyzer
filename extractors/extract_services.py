"""Extract Service definitions from Kubernetes manifests."""

from __future__ import annotations

import logging
from typing import Any

from extractors.base import BaseExtractor

logger = logging.getLogger(__name__)

SEARCH_PATTERNS = [
    "**/service.yaml",
    "**/service.yml",
    "**/service*.yaml",
    "**/service*.yml",
    "config/**/service*.yaml",
    "config/**/service*.yml",
]


def _has_unresolved_templates(value: Any) -> bool:
    """Check if a value contains unresolvable Helm template syntax."""
    if isinstance(value, str):
        return "{{" in value and "}}" in value
    return False


class ServiceExtractor(BaseExtractor):
    """Extract Kubernetes Service definitions."""

    def extract(self) -> dict[str, Any]:
        files = self.find_yaml_files(SEARCH_PATTERNS)
        services: list[dict[str, Any]] = []

        for fpath in files:
            for doc in self.parse_yaml_safe(fpath):
                kind = doc.get("kind", "")
                if kind != "Service":
                    continue
                metadata = doc.get("metadata", {})
                if not isinstance(metadata, dict):
                    metadata = {}
                name = metadata.get("name", "")
                # Skip entries with unresolvable template names
                if _has_unresolved_templates(name):
                    logger.debug("Skipping templated service name in %s", fpath)
                    continue

                spec = doc.get("spec", {})
                if not isinstance(spec, dict):
                    spec = {}

                raw_ports = spec.get("ports", [])
                if not isinstance(raw_ports, list):
                    raw_ports = []
                ports: list[dict[str, Any]] = []
                for p in raw_ports:
                    if not isinstance(p, dict):
                        continue
                    # Skip ports with template values
                    if _has_unresolved_templates(p.get("port")) or _has_unresolved_templates(
                        p.get("targetPort")
                    ):
                        continue
                    ports.append(
                        {
                            "name": p.get("name", ""),
                            "port": p.get("port", 0),
                            "targetPort": p.get("targetPort", 0),
                            "protocol": p.get("protocol", "TCP"),
                        }
                    )

                selector = spec.get("selector", {})
                if not isinstance(selector, dict):
                    selector = {}

                services.append(
                    {
                        "name": name,
                        "source": self._relative(fpath),
                        "type": spec.get("type", "ClusterIP"),
                        "ports": ports,
                        "selector": selector,
                    }
                )

        return {"services": services}
