"""Extract Helm chart values and metadata."""

from __future__ import annotations

import logging
from typing import Any

from extractors.base import BaseExtractor

logger = logging.getLogger(__name__)

CHART_PATTERNS = [
    "**/Chart.yaml",
    "**/Chart.yml",
]

VALUES_PATTERNS = [
    "**/values.yaml",
    "**/values.yml",
]

SECURITY_KEYS = {
    "securityContext",
    "podSecurityContext",
    "tls",
    "networkPolicy",
    "serviceAccount",
    "rbac",
    "auth",
}


def _flatten_dict(d: Any, prefix: str = "") -> dict[str, Any]:
    """Flatten a nested dict to dot-notation keys."""
    if not isinstance(d, dict):
        return {prefix: d} if prefix else {}
    items: dict[str, Any] = {}
    for key, value in d.items():
        new_key = f"{prefix}.{key}" if prefix else str(key)
        if isinstance(value, dict):
            items.update(_flatten_dict(value, new_key))
        else:
            items[new_key] = value
    return items


class HelmExtractor(BaseExtractor):
    """Extract Helm chart metadata and security-relevant values defaults."""

    def extract(self) -> dict[str, Any]:
        chart_files = self.find_yaml_files(CHART_PATTERNS)
        values_files = self.find_yaml_files(VALUES_PATTERNS)

        chart_name = ""
        chart_version = ""

        # Parse Chart.yaml for metadata
        for fpath in chart_files:
            for doc in self.parse_yaml_safe(fpath):
                chart_name = doc.get("name", chart_name)
                chart_version = doc.get("version", chart_version)
                break  # Use first Chart.yaml found
            if chart_name:
                break

        # Parse values.yaml for security-relevant defaults
        values_defaults: dict[str, Any] = {}
        for fpath in values_files:
            for doc in self.parse_yaml_safe(fpath):
                if not isinstance(doc, dict):
                    continue
                for key in doc:
                    if key in SECURITY_KEYS:
                        flattened = _flatten_dict(doc[key], key)
                        values_defaults.update(flattened)

        if not chart_name and not values_defaults:
            return {"helm": {}}

        return {
            "helm": {
                "chart_name": chart_name,
                "chart_version": chart_version,
                "values_defaults": values_defaults,
            }
        }
