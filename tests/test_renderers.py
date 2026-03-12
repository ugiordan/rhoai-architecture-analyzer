"""Unit tests for all renderers."""

from __future__ import annotations

import sys
import unittest
from pathlib import Path

PROJECT_ROOT = Path(__file__).resolve().parent.parent
sys.path.insert(0, str(PROJECT_ROOT))


def _sample_data() -> dict:
    """Return a sample component-architecture.json structure."""
    return {
        "component": "test-operator",
        "repo": "opendatahub-io/test-operator",
        "extracted_at": "2025-01-01T00:00:00Z",
        "analyzer_version": "0.1.0",
        "crds": [
            {
                "group": "test.opendatahub.io",
                "version": "v1alpha1",
                "kind": "TestResource",
                "scope": "Namespaced",
                "fields_count": 15,
                "validation_rules": ["has(self.spec)"],
                "source": "config/crd/bases/test.yaml",
            }
        ],
        "rbac": {
            "cluster_roles": [
                {
                    "name": "test-manager-role",
                    "source": "config/rbac/role.yaml",
                    "rules": [
                        {
                            "apiGroups": [""],
                            "resources": ["configmaps", "secrets"],
                            "verbs": ["get", "list", "watch"],
                            "resourceNames": [],
                        },
                        {
                            "apiGroups": ["apps"],
                            "resources": ["deployments"],
                            "verbs": ["get", "list", "create", "update"],
                            "resourceNames": [],
                        },
                    ],
                }
            ],
            "cluster_role_bindings": [
                {
                    "name": "test-manager-rolebinding",
                    "role_ref": "test-manager-role",
                    "subjects": [
                        {
                            "kind": "ServiceAccount",
                            "name": "test-operator",
                            "namespace": "test-system",
                        }
                    ],
                    "source": "config/rbac/binding.yaml",
                }
            ],
            "roles": [],
            "role_bindings": [],
            "kubebuilder_markers": [
                {
                    "file": "controllers/test_controller.go",
                    "line": 10,
                    "marker": "+kubebuilder:rbac:groups=test.opendatahub.io,resources=testresources,verbs=get;list",
                    "parsed": {"groups": "test.opendatahub.io", "resources": "testresources", "verbs": ["get", "list"]},
                }
            ],
        },
        "deployments": [
            {
                "name": "test-controller-manager",
                "kind": "Deployment",
                "source": "config/manager/deployment.yaml",
                "replicas": 1,
                "service_account": "test-operator",
                "automount_service_account_token": True,
                "containers": [
                    {
                        "name": "manager",
                        "image": "quay.io/opendatahub/test-operator:v1.0.0",
                        "ports": [
                            {"name": "https", "containerPort": 8443, "protocol": "TCP"}
                        ],
                        "security_context": {
                            "allowPrivilegeEscalation": False,
                            "readOnlyRootFilesystem": True,
                            "runAsNonRoot": True,
                            "capabilities": {"drop": ["ALL"], "add": []},
                            "seccompProfile": {"type": "RuntimeDefault"},
                        },
                        "env_from_secrets": ["db-creds"],
                        "env_from_configmaps": ["config"],
                        "volume_mounts": [
                            {"name": "tls", "mountPath": "/certs", "secret": "tls-cert"}
                        ],
                        "resources": {
                            "requests": {"cpu": "100m", "memory": "128Mi"},
                            "limits": {"cpu": "500m", "memory": "256Mi"},
                        },
                    }
                ],
            }
        ],
        "services": [
            {
                "name": "test-webhook-service",
                "source": "config/default/service.yaml",
                "type": "ClusterIP",
                "ports": [
                    {"name": "https", "port": 8443, "targetPort": 8443, "protocol": "TCP"}
                ],
                "selector": {"app": "test-operator"},
            }
        ],
        "network_policies": [
            {
                "name": "test-netpol",
                "source": "config/default/netpol.yaml",
                "pod_selector": {"app": "test-operator"},
                "policy_types": ["Ingress"],
                "ingress_rules": [
                    {"ports": [{"port": 8443, "protocol": "TCP"}], "from": []}
                ],
                "egress_rules": [],
            }
        ],
        "secrets_referenced": [
            {
                "name": "tls-cert",
                "type": "kubernetes.io/tls",
                "referenced_by": ["deployment/test-controller-manager"],
                "provisioned_by": "cert-manager",
            }
        ],
        "controller_watches": [
            {"type": "For", "gvk": "test.opendatahub.io/v1alpha1/TestResource", "source": "controllers/test_controller.go:30"},
            {"type": "Owns", "gvk": "apps/v1/Deployment", "source": "controllers/test_controller.go:31"},
            {"type": "Watches", "gvk": "/v1/ConfigMap", "source": "controllers/test_controller.go:32"},
        ],
        "dependencies": {
            "go_modules": [
                {"module": "sigs.k8s.io/controller-runtime", "version": "v0.19.0"},
                {"module": "k8s.io/api", "version": "v0.31.0"},
            ],
            "internal_odh": [
                {"component": "model-registry", "interaction": "Go module dependency"},
            ],
        },
        "dockerfiles": [
            {
                "path": "Dockerfile",
                "base_image": "golang:1.22",
                "stages": 2,
                "user": "65532",
                "exposed_ports": [8443],
                "issues": [],
            }
        ],
        "helm": {
            "chart_name": "test-chart",
            "chart_version": "0.1.0",
            "values_defaults": {"tls.enabled": True},
        },
    }


class TestRBACRenderer(unittest.TestCase):
    def test_render_rbac(self) -> None:
        from renderers.render_rbac import RBACRenderer

        r = RBACRenderer(_sample_data())
        output = r.render()
        self.assertIn("graph TD", output)
        self.assertIn("test-manager-role", output)
        self.assertIn("ServiceAccount", output)
        self.assertIn("configmaps", output)
        self.assertIn("deployments", output)
        self.assertEqual(r.filename(), "rbac.mmd")

    def test_render_empty_rbac(self) -> None:
        from renderers.render_rbac import RBACRenderer

        data = {"component": "empty", "rbac": {}}
        r = RBACRenderer(data)
        output = r.render()
        self.assertIn("graph TD", output)
        self.assertIn("No RBAC data", output)


class TestComponentRenderer(unittest.TestCase):
    def test_render_component(self) -> None:
        from renderers.render_component import ComponentRenderer

        r = ComponentRenderer(_sample_data())
        output = r.render()
        self.assertIn("graph LR", output)
        self.assertIn("TestResource", output)
        self.assertIn("Deployment", output)
        self.assertIn("model-registry", output)
        self.assertEqual(r.filename(), "component.mmd")


class TestSecurityNetworkRenderer(unittest.TestCase):
    def test_render_security(self) -> None:
        from renderers.render_security_network import SecurityNetworkRenderer

        r = SecurityNetworkRenderer(_sample_data())
        output = r.render()
        self.assertIn("SECURITY & NETWORK ARCHITECTURE", output)
        self.assertIn("test-webhook-service", output)
        self.assertIn("8443", output)
        self.assertIn("RBAC SUMMARY", output)
        self.assertIn("SECRETS INVENTORY", output)
        self.assertIn("tls-cert", output)
        self.assertIn("DEPLOYMENT SECURITY CONTROLS", output)
        self.assertEqual(r.filename(), "security-network.txt")


class TestDependencyRenderer(unittest.TestCase):
    def test_render_dependencies(self) -> None:
        from renderers.render_dependencies import DependencyRenderer

        r = DependencyRenderer(_sample_data())
        output = r.render()
        self.assertIn("graph LR", output)
        self.assertIn("test-operator", output)
        self.assertIn("model-registry", output)
        self.assertEqual(r.filename(), "dependencies.mmd")


class TestC4Renderer(unittest.TestCase):
    def test_render_c4(self) -> None:
        from renderers.render_c4 import C4Renderer

        r = C4Renderer(_sample_data())
        output = r.render()
        self.assertIn("workspace", output)
        self.assertIn("model", output)
        self.assertIn("softwareSystem", output)
        self.assertIn("container", output)
        self.assertIn("Kubernetes API", output)
        self.assertIn("test-controller-manager", output)
        self.assertEqual(r.filename(), "c4-context.dsl")


class TestDataflowRenderer(unittest.TestCase):
    def test_render_dataflow(self) -> None:
        from renderers.render_dataflow import DataflowRenderer

        r = DataflowRenderer(_sample_data())
        output = r.render()
        self.assertIn("sequenceDiagram", output)
        self.assertIn("KubernetesAPI", output)
        self.assertIn("TestResource", output)
        self.assertIn("Deployment", output)
        self.assertEqual(r.filename(), "dataflow.mmd")


class TestRendererEdgeCases(unittest.TestCase):
    def test_empty_data(self) -> None:
        from renderers.render_rbac import RBACRenderer
        from renderers.render_component import ComponentRenderer
        from renderers.render_security_network import SecurityNetworkRenderer
        from renderers.render_dependencies import DependencyRenderer
        from renderers.render_c4 import C4Renderer
        from renderers.render_dataflow import DataflowRenderer

        data = {"component": "empty"}
        for cls in [RBACRenderer, ComponentRenderer, SecurityNetworkRenderer,
                    DependencyRenderer, C4Renderer, DataflowRenderer]:
            r = cls(data)
            output = r.render()
            self.assertIsInstance(output, str)
            self.assertGreater(len(output), 0)


if __name__ == "__main__":
    unittest.main()
