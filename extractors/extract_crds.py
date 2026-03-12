"""Extract CRD schemas from Kubernetes manifests."""

from __future__ import annotations

import logging
from typing import Any

from extractors.base import BaseExtractor

logger = logging.getLogger(__name__)

SEARCH_PATTERNS = [
    "config/crd/bases/*.yaml",
    "config/crd/bases/*.yml",
    "deploy/crds/*.yaml",
    "deploy/crds/*.yml",
    "charts/**/crds/*.yaml",
    "charts/**/crds/*.yml",
    "manifests/**/crd*.yaml",
    "manifests/**/crd*.yml",
]


def _count_fields(schema: dict[str, Any], _depth: int = 0) -> int:
    """Recursively count fields in an OpenAPI v3 schema."""
    if not isinstance(schema, dict) or _depth > 50:
        return 0
    count = 0
    props = schema.get("properties")
    if isinstance(props, dict) and props:
        count += len(props)
        for prop_schema in props.values():
            if isinstance(prop_schema, dict) and prop_schema:
                count += _count_fields(prop_schema, _depth + 1)
    # Also count items in arrays
    items = schema.get("items")
    if isinstance(items, dict) and items:
        count += _count_fields(items, _depth + 1)
    # additionalProperties
    additional = schema.get("additionalProperties")
    if isinstance(additional, dict) and additional:
        count += _count_fields(additional, _depth + 1)
    return count


def _extract_cel_rules(schema: dict[str, Any], _depth: int = 0) -> list[str]:
    """Recursively extract CEL validation rules from an OpenAPI v3 schema."""
    if not isinstance(schema, dict) or _depth > 50:
        return []
    rules: list[str] = []
    validations = schema.get("x-kubernetes-validations")
    if isinstance(validations, list):
        for v in validations:
            if isinstance(v, dict) and "rule" in v:
                rules.append(str(v["rule"]))
    # Recurse into properties
    props = schema.get("properties")
    if isinstance(props, dict) and props:
        for prop_schema in props.values():
            if isinstance(prop_schema, dict) and prop_schema:
                rules.extend(_extract_cel_rules(prop_schema, _depth + 1))
    # items
    items = schema.get("items")
    if isinstance(items, dict) and items:
        rules.extend(_extract_cel_rules(items, _depth + 1))
    return rules


class CRDExtractor(BaseExtractor):
    """Extract CustomResourceDefinition data."""

    def extract(self) -> dict[str, Any]:
        files = self.find_yaml_files(SEARCH_PATTERNS)
        crds: list[dict[str, Any]] = []
        for fpath in files:
            for doc in self.parse_yaml_safe(fpath):
                kind = doc.get("kind", "")
                if kind != "CustomResourceDefinition":
                    continue
                spec = doc.get("spec", {})
                if not isinstance(spec, dict):
                    continue
                group = spec.get("group", "")
                names = spec.get("names", {})
                crd_kind = names.get("kind", "") if isinstance(names, dict) else ""
                scope = spec.get("scope", "")
                versions = spec.get("versions", [])
                if not isinstance(versions, list):
                    versions = []
                for ver in versions:
                    if not isinstance(ver, dict):
                        continue
                    ver_name = ver.get("name", "")
                    schema = ver.get("schema", {})
                    if isinstance(schema, dict):
                        openapi = schema.get("openAPIV3Schema", {})
                    else:
                        openapi = {}
                    fields_count = _count_fields(openapi) if isinstance(openapi, dict) else 0
                    cel_rules = _extract_cel_rules(openapi) if isinstance(openapi, dict) else []
                    crds.append(
                        {
                            "group": group,
                            "version": ver_name,
                            "kind": crd_kind,
                            "scope": scope,
                            "fields_count": fields_count,
                            "validation_rules": cel_rules,
                            "source": self._relative(fpath),
                        }
                    )
        return {"crds": crds}
