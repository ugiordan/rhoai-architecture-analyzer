"""Aggregate per-component architecture JSONs into a platform-wide view."""

from __future__ import annotations

import json
import logging
from datetime import datetime, timezone
from pathlib import Path
from typing import Any

logger = logging.getLogger(__name__)


class PlatformAggregator:
    """Merge per-repo component-architecture.json files into a platform JSON."""

    def __init__(self, results_dir: str) -> None:
        self.results_dir = Path(results_dir).resolve()
        if not self.results_dir.is_dir():
            raise ValueError(f"Results directory does not exist: {self.results_dir}")

    def aggregate(self) -> dict[str, Any]:
        """Read all component JSONs and merge into platform architecture."""
        components: list[dict[str, Any]] = []

        # Find all component-architecture.json files
        for json_path in sorted(self.results_dir.rglob("component-architecture.json")):
            try:
                data = json.loads(json_path.read_text(encoding="utf-8"))
                components.append(data)
                logger.info("Loaded: %s", json_path)
            except (json.JSONDecodeError, OSError) as exc:
                logger.warning("Failed to load %s: %s", json_path, exc)

        if not components:
            logger.warning("No component architecture files found in %s", self.results_dir)

        # Build platform-wide aggregation
        all_crds: list[dict[str, Any]] = []
        all_services: list[dict[str, Any]] = []
        all_secrets: list[dict[str, Any]] = []
        all_rbac_cluster_roles: list[dict[str, Any]] = []
        component_names: list[str] = []
        dependency_graph: list[dict[str, str]] = []

        # CRD ownership: which component defines which CRDs
        crd_owners: dict[str, str] = {}

        for comp_data in components:
            comp_name = comp_data.get("component", "unknown")
            component_names.append(comp_name)

            # CRDs
            for crd in comp_data.get("crds", []):
                crd_with_owner = {**crd, "owner": comp_name}
                all_crds.append(crd_with_owner)
                kind = crd.get("kind", "")
                if kind:
                    crd_owners[kind] = comp_name

            # Services
            for svc in comp_data.get("services", []):
                all_services.append({**svc, "owner": comp_name})

            # Secrets
            for secret in comp_data.get("secrets_referenced", []):
                all_secrets.append({**secret, "owner": comp_name})

            # RBAC
            rbac = comp_data.get("rbac", {})
            for cr in rbac.get("cluster_roles", []):
                all_rbac_cluster_roles.append({**cr, "owner": comp_name})

            # Dependencies
            deps = comp_data.get("dependencies", {})
            for odh in deps.get("internal_odh", []):
                dependency_graph.append(
                    {
                        "from": comp_name,
                        "to": odh.get("component", ""),
                        "type": "go-module",
                    }
                )

            # Cross-component watches (controller watching CRDs from other components)
            for watch in comp_data.get("controller_watches", []):
                if watch.get("type") == "For":
                    gvk = watch.get("gvk", "")
                    kind = gvk.rsplit("/", 1)[-1] if "/" in gvk else gvk
                    if kind in crd_owners and crd_owners[kind] != comp_name:
                        dependency_graph.append(
                            {
                                "from": comp_name,
                                "to": crd_owners[kind],
                                "type": f"watches-crd:{kind}",
                            }
                        )

        return {
            "platform": "OpenShift AI",
            "aggregated_at": datetime.now(timezone.utc).isoformat(),
            "components": component_names,
            "component_count": len(components),
            "crds": all_crds,
            "crd_ownership": crd_owners,
            "services": all_services,
            "secrets_referenced": all_secrets,
            "rbac_cluster_roles": all_rbac_cluster_roles,
            "dependency_graph": dependency_graph,
            "component_data": components,
        }
