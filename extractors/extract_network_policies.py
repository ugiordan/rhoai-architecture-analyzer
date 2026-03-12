"""Extract NetworkPolicy definitions from Kubernetes manifests."""

from __future__ import annotations

import logging
from typing import Any

from extractors.base import BaseExtractor

logger = logging.getLogger(__name__)

SEARCH_PATTERNS = [
    "**/networkpolicy.yaml",
    "**/networkpolicy.yml",
    "**/network-policy*.yaml",
    "**/network-policy*.yml",
    "**/netpol*.yaml",
    "**/netpol*.yml",
    "config/**/networkpolicy*.yaml",
    "charts/**/templates/networkpolicy*.yaml",
]


def _extract_ingress_rules(rules: Any) -> list[dict[str, Any]]:
    """Extract ingress rules."""
    if not isinstance(rules, list):
        return []
    result: list[dict[str, Any]] = []
    for rule in rules:
        if not isinstance(rule, dict):
            continue
        ports = []
        for p in rule.get("ports", []):
            if isinstance(p, dict):
                ports.append(
                    {
                        "port": p.get("port", 0),
                        "protocol": p.get("protocol", "TCP"),
                    }
                )
        from_selectors = rule.get("from", [])
        if not isinstance(from_selectors, list):
            from_selectors = []
        result.append({"ports": ports, "from": from_selectors})
    return result


def _extract_egress_rules(rules: Any) -> list[dict[str, Any]]:
    """Extract egress rules."""
    if not isinstance(rules, list):
        return []
    result: list[dict[str, Any]] = []
    for rule in rules:
        if not isinstance(rule, dict):
            continue
        ports = []
        for p in rule.get("ports", []):
            if isinstance(p, dict):
                ports.append(
                    {
                        "port": p.get("port", 0),
                        "protocol": p.get("protocol", "TCP"),
                    }
                )
        to_selectors = rule.get("to", [])
        if not isinstance(to_selectors, list):
            to_selectors = []
        result.append({"ports": ports, "to": to_selectors})
    return result


class NetworkPolicyExtractor(BaseExtractor):
    """Extract NetworkPolicy definitions."""

    def extract(self) -> dict[str, Any]:
        files = self.find_yaml_files(SEARCH_PATTERNS)
        policies: list[dict[str, Any]] = []

        for fpath in files:
            for doc in self.parse_yaml_safe(fpath):
                kind = doc.get("kind", "")
                if kind != "NetworkPolicy":
                    continue
                metadata = doc.get("metadata", {})
                if not isinstance(metadata, dict):
                    metadata = {}
                spec = doc.get("spec", {})
                if not isinstance(spec, dict):
                    continue

                pod_selector = spec.get("podSelector", {})
                if not isinstance(pod_selector, dict):
                    pod_selector = {}
                match_labels = pod_selector.get("matchLabels", {})
                if not isinstance(match_labels, dict):
                    match_labels = {}

                policy_types = spec.get("policyTypes", [])
                if not isinstance(policy_types, list):
                    policy_types = []

                policies.append(
                    {
                        "name": metadata.get("name", ""),
                        "source": self._relative(fpath),
                        "pod_selector": match_labels,
                        "policy_types": policy_types,
                        "ingress_rules": _extract_ingress_rules(spec.get("ingress")),
                        "egress_rules": _extract_egress_rules(spec.get("egress")),
                    }
                )

        return {"network_policies": policies}
