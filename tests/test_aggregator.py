"""Unit tests for the platform aggregator."""

from __future__ import annotations

import json
import shutil
import sys
import tempfile
import unittest
from pathlib import Path

PROJECT_ROOT = Path(__file__).resolve().parent.parent
sys.path.insert(0, str(PROJECT_ROOT))


def _sample_component(name: str, crds: list[dict] | None = None) -> dict:
    """Create a minimal component architecture dict."""
    return {
        "component": name,
        "repo": f"opendatahub-io/{name}",
        "extracted_at": "2025-01-01T00:00:00Z",
        "analyzer_version": "0.1.0",
        "crds": crds or [],
        "rbac": {
            "cluster_roles": [
                {
                    "name": f"{name}-role",
                    "source": "config/rbac/role.yaml",
                    "rules": [
                        {"apiGroups": [""], "resources": ["configmaps"], "verbs": ["get"]}
                    ],
                }
            ],
            "cluster_role_bindings": [],
            "roles": [],
            "role_bindings": [],
            "kubebuilder_markers": [],
        },
        "deployments": [],
        "services": [
            {
                "name": f"{name}-svc",
                "source": "config/svc.yaml",
                "type": "ClusterIP",
                "ports": [{"name": "https", "port": 8443, "targetPort": 8443, "protocol": "TCP"}],
                "selector": {"app": name},
            }
        ],
        "network_policies": [],
        "secrets_referenced": [
            {
                "name": f"{name}-tls",
                "type": "kubernetes.io/tls",
                "referenced_by": [f"deployment/{name}"],
                "provisioned_by": "cert-manager",
            }
        ],
        "controller_watches": [],
        "dependencies": {
            "go_modules": [],
            "internal_odh": [],
        },
        "dockerfiles": [],
        "helm": {},
    }


class TestPlatformAggregator(unittest.TestCase):
    def test_aggregate_multiple_components(self) -> None:
        from aggregator.aggregate import PlatformAggregator

        tmpdir = tempfile.mkdtemp(prefix="aggregator-test-")
        try:
            # Create component directories with JSON files
            comp1 = _sample_component(
                "operator-a",
                crds=[
                    {
                        "group": "a.opendatahub.io",
                        "version": "v1",
                        "kind": "ResourceA",
                        "scope": "Namespaced",
                        "fields_count": 10,
                        "validation_rules": [],
                        "source": "config/crd/bases/a.yaml",
                    }
                ],
            )
            comp2 = _sample_component("operator-b")
            # operator-b watches ResourceA from operator-a
            comp2["controller_watches"] = [
                {"type": "For", "gvk": "a.opendatahub.io/v1/ResourceA", "source": "ctrl.go:10"}
            ]
            comp2["dependencies"]["internal_odh"] = [
                {"component": "operator-a", "interaction": "Go module dependency"}
            ]

            dir1 = Path(tmpdir) / "operator-a"
            dir2 = Path(tmpdir) / "operator-b"
            dir1.mkdir()
            dir2.mkdir()
            (dir1 / "component-architecture.json").write_text(
                json.dumps(comp1), encoding="utf-8"
            )
            (dir2 / "component-architecture.json").write_text(
                json.dumps(comp2), encoding="utf-8"
            )

            agg = PlatformAggregator(tmpdir)
            result = agg.aggregate()

            self.assertEqual(result["component_count"], 2)
            self.assertIn("operator-a", result["components"])
            self.assertIn("operator-b", result["components"])

            # CRDs
            self.assertEqual(len(result["crds"]), 1)
            self.assertEqual(result["crd_ownership"]["ResourceA"], "operator-a")

            # Services
            self.assertEqual(len(result["services"]), 2)

            # Secrets
            self.assertEqual(len(result["secrets_referenced"]), 2)

            # Dependency graph should have go-module dependency
            dep_types = {d["type"] for d in result["dependency_graph"]}
            self.assertIn("go-module", dep_types)

        finally:
            shutil.rmtree(tmpdir)

    def test_aggregate_empty_dir(self) -> None:
        from aggregator.aggregate import PlatformAggregator

        tmpdir = tempfile.mkdtemp(prefix="aggregator-test-")
        try:
            agg = PlatformAggregator(tmpdir)
            result = agg.aggregate()
            self.assertEqual(result["component_count"], 0)
            self.assertEqual(result["components"], [])
        finally:
            shutil.rmtree(tmpdir)

    def test_invalid_dir(self) -> None:
        from aggregator.aggregate import PlatformAggregator

        with self.assertRaises(ValueError):
            PlatformAggregator("/nonexistent/path/xyz")


class TestPlatformRenderer(unittest.TestCase):
    def test_render_all(self) -> None:
        from aggregator.render_platform import PlatformRenderer

        platform_data = {
            "platform": "OpenShift AI",
            "components": ["operator-a", "operator-b"],
            "crds": [
                {
                    "group": "a.opendatahub.io",
                    "version": "v1",
                    "kind": "ResourceA",
                    "owner": "operator-a",
                }
            ],
            "crd_ownership": {"ResourceA": "operator-a"},
            "services": [
                {
                    "name": "svc-a",
                    "owner": "operator-a",
                    "type": "ClusterIP",
                    "ports": [{"port": 8443, "protocol": "TCP"}],
                }
            ],
            "secrets_referenced": [],
            "rbac_cluster_roles": [
                {
                    "name": "role-a",
                    "owner": "operator-a",
                    "rules": [
                        {"apiGroups": [""], "resources": ["configmaps"], "verbs": ["get"]}
                    ],
                }
            ],
            "dependency_graph": [
                {"from": "operator-b", "to": "operator-a", "type": "go-module"},
                {"from": "operator-b", "to": "operator-a", "type": "watches-crd:ResourceA"},
            ],
            "component_data": [],
        }

        renderer = PlatformRenderer(platform_data)
        diagrams = renderer.render_all()

        self.assertEqual(len(diagrams), 4)
        self.assertIn("platform-dependencies.mmd", diagrams)
        self.assertIn("platform-crd-ownership.mmd", diagrams)
        self.assertIn("platform-rbac-overview.mmd", diagrams)
        self.assertIn("platform-network-topology.mmd", diagrams)

        # Check content
        dep_graph = diagrams["platform-dependencies.mmd"]
        self.assertIn("graph LR", dep_graph)
        self.assertIn("operator_a", dep_graph)
        self.assertIn("operator_b", dep_graph)

        crd_map = diagrams["platform-crd-ownership.mmd"]
        self.assertIn("ResourceA", crd_map)

        rbac = diagrams["platform-rbac-overview.mmd"]
        self.assertIn("role-a", rbac)

        network = diagrams["platform-network-topology.mmd"]
        self.assertIn("svc-a", network)


if __name__ == "__main__":
    unittest.main()
