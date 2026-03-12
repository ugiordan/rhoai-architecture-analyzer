"""Unit tests for all extractors."""

from __future__ import annotations

import os
import sys
import tempfile
import shutil
import unittest
from pathlib import Path

# Ensure project root on path
PROJECT_ROOT = Path(__file__).resolve().parent.parent
sys.path.insert(0, str(PROJECT_ROOT))

FIXTURES_DIR = Path(__file__).resolve().parent / "fixtures"


def _make_repo(structure: dict[str, str]) -> str:
    """Create a temporary repo directory with given file structure.

    Args:
        structure: mapping of relative paths to file content

    Returns:
        Path to the temporary directory (caller must clean up).
    """
    tmpdir = tempfile.mkdtemp(prefix="analyzer-test-")
    for relpath, content in structure.items():
        fpath = Path(tmpdir) / relpath
        fpath.parent.mkdir(parents=True, exist_ok=True)
        fpath.write_text(content, encoding="utf-8")
    return tmpdir


class TestCRDExtractor(unittest.TestCase):
    def test_extract_crds(self) -> None:
        from extractors.extract_crds import CRDExtractor

        crd_content = (FIXTURES_DIR / "sample_crd.yaml").read_text()
        tmpdir = _make_repo({"config/crd/bases/dsc.yaml": crd_content})
        try:
            ext = CRDExtractor(tmpdir)
            result = ext.extract()
            crds = result["crds"]
            self.assertEqual(len(crds), 1)
            crd = crds[0]
            self.assertEqual(crd["group"], "datasciencecluster.opendatahub.io")
            self.assertEqual(crd["version"], "v1")
            self.assertEqual(crd["kind"], "DataScienceCluster")
            self.assertEqual(crd["scope"], "Cluster")
            self.assertGreater(crd["fields_count"], 0)
            self.assertGreater(len(crd["validation_rules"]), 0)
            self.assertIn("has(self.components)", crd["validation_rules"])
        finally:
            shutil.rmtree(tmpdir)

    def test_extract_empty_dir(self) -> None:
        from extractors.extract_crds import CRDExtractor

        tmpdir = _make_repo({})
        try:
            ext = CRDExtractor(tmpdir)
            result = ext.extract()
            self.assertEqual(result["crds"], [])
        finally:
            shutil.rmtree(tmpdir)

    def test_malformed_yaml(self) -> None:
        from extractors.extract_crds import CRDExtractor

        tmpdir = _make_repo({"config/crd/bases/bad.yaml": "{{invalid yaml: [["})
        try:
            ext = CRDExtractor(tmpdir)
            result = ext.extract()
            self.assertEqual(result["crds"], [])
        finally:
            shutil.rmtree(tmpdir)


class TestRBACExtractor(unittest.TestCase):
    def test_extract_rbac(self) -> None:
        from extractors.extract_rbac import RBACExtractor

        rbac_content = (FIXTURES_DIR / "sample_rbac.yaml").read_text()
        controller_content = (FIXTURES_DIR / "sample_controller.go").read_text()
        tmpdir = _make_repo({
            "config/rbac/role.yaml": rbac_content,
            "controllers/dsc_controller.go": controller_content,
        })
        try:
            ext = RBACExtractor(tmpdir)
            result = ext.extract()
            rbac = result["rbac"]

            # ClusterRole
            self.assertEqual(len(rbac["cluster_roles"]), 1)
            cr = rbac["cluster_roles"][0]
            self.assertEqual(cr["name"], "opendatahub-operator-manager-role")
            self.assertGreater(len(cr["rules"]), 0)

            # ClusterRoleBinding
            self.assertEqual(len(rbac["cluster_role_bindings"]), 1)
            crb = rbac["cluster_role_bindings"][0]
            self.assertEqual(crb["role_ref"], "opendatahub-operator-manager-role")
            self.assertEqual(len(crb["subjects"]), 1)
            self.assertEqual(crb["subjects"][0]["name"], "opendatahub-operator")

            # Kubebuilder markers
            self.assertGreater(len(rbac["kubebuilder_markers"]), 0)
        finally:
            shutil.rmtree(tmpdir)


class TestServiceExtractor(unittest.TestCase):
    def test_extract_services(self) -> None:
        from extractors.extract_services import ServiceExtractor

        svc_content = (FIXTURES_DIR / "sample_service.yaml").read_text()
        tmpdir = _make_repo({"config/default/service.yaml": svc_content})
        try:
            ext = ServiceExtractor(tmpdir)
            result = ext.extract()
            services = result["services"]
            self.assertEqual(len(services), 1)
            svc = services[0]
            self.assertEqual(svc["name"], "opendatahub-operator-webhook-service")
            self.assertEqual(svc["type"], "ClusterIP")
            self.assertEqual(len(svc["ports"]), 2)
            self.assertEqual(svc["ports"][0]["port"], 8443)
        finally:
            shutil.rmtree(tmpdir)


class TestDeploymentExtractor(unittest.TestCase):
    def test_extract_deployments(self) -> None:
        from extractors.extract_deployments import DeploymentExtractor

        dep_content = (FIXTURES_DIR / "sample_deployment.yaml").read_text()
        tmpdir = _make_repo({"config/manager/deployment.yaml": dep_content})
        try:
            ext = DeploymentExtractor(tmpdir)
            result = ext.extract()
            deps = result["deployments"]
            self.assertEqual(len(deps), 1)
            dep = deps[0]
            self.assertEqual(dep["name"], "opendatahub-operator-controller-manager")
            self.assertEqual(dep["service_account"], "opendatahub-operator")
            self.assertEqual(dep["replicas"], 1)

            containers = dep["containers"]
            self.assertEqual(len(containers), 1)
            c = containers[0]
            self.assertEqual(c["name"], "manager")
            self.assertIn("quay.io", c["image"])
            self.assertEqual(len(c["ports"]), 2)

            # Security context
            sc = c["security_context"]
            self.assertFalse(sc["allowPrivilegeEscalation"])
            self.assertTrue(sc["readOnlyRootFilesystem"])
            self.assertTrue(sc["runAsNonRoot"])
            self.assertIn("ALL", sc["capabilities"]["drop"])

            # Env refs
            self.assertIn("db-credentials", c["env_from_secrets"])
            self.assertIn("extra-secrets", c["env_from_secrets"])
            self.assertIn("operator-config", c["env_from_configmaps"])

            # Volume mounts
            self.assertGreater(len(c["volume_mounts"]), 0)
            tls_mount = next(
                (m for m in c["volume_mounts"] if m["name"] == "tls-certs"), None
            )
            self.assertIsNotNone(tls_mount)
            self.assertEqual(tls_mount["secret"], "webhook-server-cert")
        finally:
            shutil.rmtree(tmpdir)


class TestNetworkPolicyExtractor(unittest.TestCase):
    def test_extract_network_policies(self) -> None:
        from extractors.extract_network_policies import NetworkPolicyExtractor

        np_content = (FIXTURES_DIR / "sample_networkpolicy.yaml").read_text()
        tmpdir = _make_repo({"config/default/networkpolicy.yaml": np_content})
        try:
            ext = NetworkPolicyExtractor(tmpdir)
            result = ext.extract()
            policies = result["network_policies"]
            self.assertEqual(len(policies), 1)
            pol = policies[0]
            self.assertEqual(pol["name"], "opendatahub-operator-netpol")
            self.assertIn("Ingress", pol["policy_types"])
            self.assertIn("Egress", pol["policy_types"])
            self.assertGreater(len(pol["ingress_rules"]), 0)
            self.assertGreater(len(pol["egress_rules"]), 0)
        finally:
            shutil.rmtree(tmpdir)


class TestControllerWatchExtractor(unittest.TestCase):
    def test_extract_watches(self) -> None:
        from extractors.extract_controller_watches import ControllerWatchExtractor

        ctrl_content = (FIXTURES_DIR / "sample_controller.go").read_text()
        tmpdir = _make_repo({
            "controllers/dsc_controller.go": ctrl_content,
        })
        try:
            ext = ControllerWatchExtractor(tmpdir)
            result = ext.extract()
            watches = result["controller_watches"]
            self.assertGreater(len(watches), 0)

            types_found = {w["type"] for w in watches}
            self.assertIn("For", types_found)
            self.assertIn("Owns", types_found)
            self.assertIn("Watches", types_found)

            # Check specific watches
            for_watches = [w for w in watches if w["type"] == "For"]
            self.assertTrue(
                any("DataScienceCluster" in w["gvk"] for w in for_watches)
            )

            owns_watches = [w for w in watches if w["type"] == "Owns"]
            self.assertTrue(any("Deployment" in w["gvk"] for w in owns_watches))
            self.assertTrue(any("Service" in w["gvk"] for w in owns_watches))
        finally:
            shutil.rmtree(tmpdir)


class TestDependencyExtractor(unittest.TestCase):
    def test_extract_dependencies(self) -> None:
        from extractors.extract_dependencies import DependencyExtractor

        gomod_content = (FIXTURES_DIR / "sample_go.mod").read_text()
        tmpdir = _make_repo({"go.mod": gomod_content})
        try:
            ext = DependencyExtractor(tmpdir)
            result = ext.extract()
            deps = result["dependencies"]
            go_mods = deps["go_modules"]
            internal = deps["internal_odh"]

            # Should have direct dependencies (not indirect)
            self.assertGreater(len(go_mods), 0)
            modules = {m["module"] for m in go_mods}
            self.assertIn("sigs.k8s.io/controller-runtime", modules)
            self.assertIn("github.com/opendatahub-io/model-registry", modules)

            # prometheus is indirect, should be excluded
            self.assertNotIn("github.com/prometheus/client_golang", modules)

            # Internal ODH deps
            self.assertGreater(len(internal), 0)
            components = {d["component"] for d in internal}
            self.assertIn("model-registry", components)
            self.assertIn("odh-model-controller", components)
        finally:
            shutil.rmtree(tmpdir)


class TestSecretExtractor(unittest.TestCase):
    def test_extract_secrets(self) -> None:
        from extractors.extract_secrets import SecretExtractor

        dep_content = (FIXTURES_DIR / "sample_deployment.yaml").read_text()
        svc_content = (FIXTURES_DIR / "sample_service.yaml").read_text()
        tmpdir = _make_repo({
            "config/manager/deployment.yaml": dep_content,
            "config/default/service.yaml": svc_content,
        })
        try:
            ext = SecretExtractor(tmpdir)
            result = ext.extract()
            secrets = result["secrets_referenced"]
            self.assertGreater(len(secrets), 0)
            names = {s["name"] for s in secrets}
            self.assertIn("db-credentials", names)
            self.assertIn("webhook-server-cert", names)
            self.assertIn("extra-secrets", names)

            # webhook-server-cert should be TLS type from annotation
            tls_secret = next(
                s for s in secrets if s["name"] == "webhook-server-cert"
            )
            self.assertEqual(tls_secret["type"], "kubernetes.io/tls")
        finally:
            shutil.rmtree(tmpdir)


class TestHelmExtractor(unittest.TestCase):
    def test_extract_helm(self) -> None:
        from extractors.extract_helm import HelmExtractor

        values_content = (FIXTURES_DIR / "sample_values.yaml").read_text()
        chart_content = "name: test-chart\nversion: 0.1.0\n"
        tmpdir = _make_repo({
            "charts/operator/values.yaml": values_content,
            "charts/operator/Chart.yaml": chart_content,
        })
        try:
            ext = HelmExtractor(tmpdir)
            result = ext.extract()
            helm = result["helm"]
            self.assertEqual(helm["chart_name"], "test-chart")
            self.assertEqual(helm["chart_version"], "0.1.0")
            defaults = helm["values_defaults"]
            self.assertIn("tls.enabled", defaults)
            self.assertTrue(defaults["tls.enabled"])
            self.assertIn("securityContext.runAsNonRoot", defaults)
        finally:
            shutil.rmtree(tmpdir)


class TestDockerfileExtractor(unittest.TestCase):
    def test_extract_dockerfiles(self) -> None:
        from extractors.extract_dockerfiles import DockerfileExtractor

        df_content = (FIXTURES_DIR / "sample_dockerfile").read_text()
        tmpdir = _make_repo({"Dockerfile": df_content})
        try:
            ext = DockerfileExtractor(tmpdir)
            result = ext.extract()
            dfs = result["dockerfiles"]
            self.assertEqual(len(dfs), 1)
            df = dfs[0]
            self.assertEqual(df["stages"], 2)
            self.assertIn("golang", df["base_image"])
            self.assertEqual(df["user"], "65532:65532")
            self.assertIn(8443, df["exposed_ports"])
            self.assertIn(8080, df["exposed_ports"])
            # No issues for this well-formed Dockerfile
            # (golang:1.22 is pinned to a version, so no unpinned tag issue for the final stage)
        finally:
            shutil.rmtree(tmpdir)


class TestBaseExtractor(unittest.TestCase):
    def test_invalid_repo_path(self) -> None:
        from extractors.base import BaseExtractor

        class DummyExtractor(BaseExtractor):
            def extract(self):
                return {}

        with self.assertRaises(ValueError):
            DummyExtractor("/nonexistent/path/xyz")

    def test_helm_template_yaml(self) -> None:
        from extractors.extract_services import ServiceExtractor

        content = """apiVersion: v1
kind: Service
metadata:
  name: {{ .Release.Name }}-svc
spec:
  ports:
    - port: {{ .Values.service.port }}
"""
        tmpdir = _make_repo({"config/default/service.yaml": content})
        try:
            ext = ServiceExtractor(tmpdir)
            result = ext.extract()
            # Should not crash, may return empty or skip
            self.assertIsInstance(result["services"], list)
        finally:
            shutil.rmtree(tmpdir)


if __name__ == "__main__":
    unittest.main()
