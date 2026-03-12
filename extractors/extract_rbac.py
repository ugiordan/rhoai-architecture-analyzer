"""Extract RBAC configuration from Kubernetes manifests and kubebuilder markers."""

from __future__ import annotations

import logging
import re
from typing import Any

from extractors.base import BaseExtractor

logger = logging.getLogger(__name__)

YAML_PATTERNS = [
    "config/rbac/*.yaml",
    "config/rbac/*.yml",
    "charts/**/templates/*role*.yaml",
    "charts/**/templates/*role*.yml",
    "deploy/rbac/*.yaml",
    "deploy/rbac/*.yml",
    "manifests/**/*role*.yaml",
    "manifests/**/*role*.yml",
]

KUBEBUILDER_RBAC_RE = re.compile(r"//\+kubebuilder:rbac:(.+)")


def _parse_kubebuilder_marker(marker_body: str) -> dict[str, Any]:
    """Parse key=value pairs from a kubebuilder RBAC marker."""
    result: dict[str, Any] = {}
    for part in marker_body.split(","):
        part = part.strip()
        if "=" in part:
            key, value = part.split("=", 1)
            key = key.strip()
            value = value.strip()
            # Some values are semicolon-separated lists
            if ";" in value:
                result[key] = value.split(";")
            else:
                result[key] = value
        else:
            # Handle bare values that are continuations of a list
            # e.g., groups=apps,resources=deployments
            pass
    return result


def _extract_rules(doc: dict[str, Any]) -> list[dict[str, Any]]:
    """Extract RBAC rules from a Role/ClusterRole document."""
    raw_rules = doc.get("rules", [])
    if not isinstance(raw_rules, list):
        return []
    rules: list[dict[str, Any]] = []
    for rule in raw_rules:
        if not isinstance(rule, dict):
            continue
        rules.append(
            {
                "apiGroups": rule.get("apiGroups", []),
                "resources": rule.get("resources", []),
                "verbs": rule.get("verbs", []),
                "resourceNames": rule.get("resourceNames", []),
            }
        )
    return rules


def _extract_subjects(doc: dict[str, Any]) -> list[dict[str, Any]]:
    """Extract subjects from a RoleBinding/ClusterRoleBinding."""
    raw = doc.get("subjects", [])
    if not isinstance(raw, list):
        return []
    subjects: list[dict[str, Any]] = []
    for subj in raw:
        if not isinstance(subj, dict):
            continue
        subjects.append(
            {
                "kind": subj.get("kind", ""),
                "name": subj.get("name", ""),
                "namespace": subj.get("namespace", ""),
            }
        )
    return subjects


class RBACExtractor(BaseExtractor):
    """Extract RBAC roles, bindings, and kubebuilder markers."""

    def extract(self) -> dict[str, Any]:
        files = self.find_yaml_files(YAML_PATTERNS)
        cluster_roles: list[dict[str, Any]] = []
        cluster_role_bindings: list[dict[str, Any]] = []
        roles: list[dict[str, Any]] = []
        role_bindings: list[dict[str, Any]] = []

        for fpath in files:
            for doc in self.parse_yaml_safe(fpath):
                kind = doc.get("kind", "")
                metadata = doc.get("metadata", {})
                if not isinstance(metadata, dict):
                    metadata = {}
                name = metadata.get("name", "")
                source = self._relative(fpath)

                if kind == "ClusterRole":
                    cluster_roles.append(
                        {
                            "name": name,
                            "source": source,
                            "rules": _extract_rules(doc),
                        }
                    )
                elif kind == "ClusterRoleBinding":
                    role_ref = doc.get("roleRef", {})
                    cluster_role_bindings.append(
                        {
                            "name": name,
                            "role_ref": role_ref.get("name", "") if isinstance(role_ref, dict) else "",
                            "subjects": _extract_subjects(doc),
                            "source": source,
                        }
                    )
                elif kind == "Role":
                    roles.append(
                        {
                            "name": name,
                            "source": source,
                            "rules": _extract_rules(doc),
                        }
                    )
                elif kind == "RoleBinding":
                    role_ref = doc.get("roleRef", {})
                    role_bindings.append(
                        {
                            "name": name,
                            "role_ref": role_ref.get("name", "") if isinstance(role_ref, dict) else "",
                            "subjects": _extract_subjects(doc),
                            "source": source,
                        }
                    )

        # Scan Go files for kubebuilder RBAC markers
        go_files = self.find_go_files()
        markers: list[dict[str, Any]] = []
        for fpath in go_files:
            try:
                lines = fpath.read_text(encoding="utf-8", errors="replace").splitlines()
            except OSError:
                continue
            for line_no, line in enumerate(lines, start=1):
                match = KUBEBUILDER_RBAC_RE.search(line)
                if match:
                    markers.append(
                        {
                            "file": self._relative(fpath),
                            "line": line_no,
                            "marker": match.group(0).lstrip("/"),
                            "parsed": _parse_kubebuilder_marker(match.group(1)),
                        }
                    )

        return {
            "rbac": {
                "cluster_roles": cluster_roles,
                "cluster_role_bindings": cluster_role_bindings,
                "roles": roles,
                "role_bindings": role_bindings,
                "kubebuilder_markers": markers,
            }
        }
