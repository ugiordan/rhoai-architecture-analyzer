"""Extract Deployment/StatefulSet definitions from Kubernetes manifests."""

from __future__ import annotations

import logging
from typing import Any

from extractors.base import BaseExtractor

logger = logging.getLogger(__name__)

SEARCH_PATTERNS = [
    "**/deployment.yaml",
    "**/deployment.yml",
    "**/deployment*.yaml",
    "**/manager*.yaml",
    "**/manager*.yml",
    "**/statefulset.yaml",
    "**/statefulset.yml",
    "charts/**/templates/deployment*.yaml",
    "charts/**/templates/deployment*.yml",
]


def _extract_security_context(sc: Any) -> dict[str, Any]:
    """Extract security context fields."""
    if not isinstance(sc, dict):
        return {}
    result: dict[str, Any] = {}
    for key in (
        "allowPrivilegeEscalation",
        "readOnlyRootFilesystem",
        "runAsNonRoot",
        "runAsUser",
        "runAsGroup",
        "privileged",
    ):
        if key in sc:
            result[key] = sc[key]
    caps = sc.get("capabilities", {})
    if isinstance(caps, dict):
        result["capabilities"] = {
            "drop": caps.get("drop", []),
            "add": caps.get("add", []),
        }
    seccomp = sc.get("seccompProfile", {})
    if isinstance(seccomp, dict):
        result["seccompProfile"] = {"type": seccomp.get("type", "")}
    return result


def _extract_env_refs(
    containers: list[dict[str, Any]],
) -> tuple[list[str], list[str]]:
    """Extract secret and configmap names from env and envFrom."""
    secrets: list[str] = []
    configmaps: list[str] = []
    for container in containers:
        if not isinstance(container, dict):
            continue
        # env[].valueFrom
        for env_var in container.get("env", []):
            if not isinstance(env_var, dict):
                continue
            value_from = env_var.get("valueFrom", {})
            if isinstance(value_from, dict):
                secret_ref = value_from.get("secretKeyRef", {})
                if isinstance(secret_ref, dict) and "name" in secret_ref:
                    name = secret_ref["name"]
                    if isinstance(name, str) and name not in secrets:
                        secrets.append(name)
                cm_ref = value_from.get("configMapKeyRef", {})
                if isinstance(cm_ref, dict) and "name" in cm_ref:
                    name = cm_ref["name"]
                    if isinstance(name, str) and name not in configmaps:
                        configmaps.append(name)
        # envFrom[]
        for env_from in container.get("envFrom", []):
            if not isinstance(env_from, dict):
                continue
            secret_ref = env_from.get("secretRef", {})
            if isinstance(secret_ref, dict) and "name" in secret_ref:
                name = secret_ref["name"]
                if isinstance(name, str) and name not in secrets:
                    secrets.append(name)
            cm_ref = env_from.get("configMapRef", {})
            if isinstance(cm_ref, dict) and "name" in cm_ref:
                name = cm_ref["name"]
                if isinstance(name, str) and name not in configmaps:
                    configmaps.append(name)
    return secrets, configmaps


def _extract_volume_mounts(
    container: dict[str, Any], volumes: list[dict[str, Any]]
) -> list[dict[str, Any]]:
    """Extract volume mounts with source info from volumes list."""
    vol_map: dict[str, dict[str, Any]] = {}
    for vol in volumes:
        if not isinstance(vol, dict):
            continue
        vol_name = vol.get("name", "")
        info: dict[str, Any] = {}
        if "secret" in vol and isinstance(vol["secret"], dict):
            info["secret"] = vol["secret"].get("secretName", "")
        if "configMap" in vol and isinstance(vol["configMap"], dict):
            info["configMap"] = vol["configMap"].get("name", "")
        if "projected" in vol:
            info["projected"] = True
        if "emptyDir" in vol:
            info["emptyDir"] = True
        vol_map[vol_name] = info

    mounts: list[dict[str, Any]] = []
    for vm in container.get("volumeMounts", []):
        if not isinstance(vm, dict):
            continue
        name = vm.get("name", "")
        entry: dict[str, Any] = {
            "name": name,
            "mountPath": vm.get("mountPath", ""),
        }
        if name in vol_map:
            entry.update(vol_map[name])
        mounts.append(entry)
    return mounts


class DeploymentExtractor(BaseExtractor):
    """Extract Deployment and StatefulSet definitions."""

    def extract(self) -> dict[str, Any]:
        files = self.find_yaml_files(SEARCH_PATTERNS)
        deployments: list[dict[str, Any]] = []

        for fpath in files:
            for doc in self.parse_yaml_safe(fpath):
                kind = doc.get("kind", "")
                if kind not in ("Deployment", "StatefulSet"):
                    continue
                metadata = doc.get("metadata", {})
                if not isinstance(metadata, dict):
                    metadata = {}
                name = metadata.get("name", "")
                spec = doc.get("spec", {})
                if not isinstance(spec, dict):
                    continue

                template = spec.get("template", {})
                if not isinstance(template, dict):
                    continue
                pod_spec = template.get("spec", {})
                if not isinstance(pod_spec, dict):
                    continue

                raw_containers = pod_spec.get("containers", [])
                if not isinstance(raw_containers, list):
                    raw_containers = []
                init_containers = pod_spec.get("initContainers", [])
                if not isinstance(init_containers, list):
                    init_containers = []
                all_containers = raw_containers + init_containers

                volumes = pod_spec.get("volumes", [])
                if not isinstance(volumes, list):
                    volumes = []

                env_secrets, env_configmaps = _extract_env_refs(all_containers)

                containers: list[dict[str, Any]] = []
                for c in raw_containers:
                    if not isinstance(c, dict):
                        continue
                    raw_ports = c.get("ports", [])
                    if not isinstance(raw_ports, list):
                        raw_ports = []
                    ports = []
                    for p in raw_ports:
                        if isinstance(p, dict):
                            ports.append(
                                {
                                    "name": p.get("name", ""),
                                    "containerPort": p.get("containerPort", 0),
                                    "protocol": p.get("protocol", "TCP"),
                                }
                            )

                    resources = c.get("resources", {})
                    if not isinstance(resources, dict):
                        resources = {}

                    containers.append(
                        {
                            "name": c.get("name", ""),
                            "image": c.get("image", ""),
                            "ports": ports,
                            "security_context": _extract_security_context(
                                c.get("securityContext", {})
                            ),
                            "env_from_secrets": env_secrets,
                            "env_from_configmaps": env_configmaps,
                            "volume_mounts": _extract_volume_mounts(c, volumes),
                            "resources": {
                                "requests": resources.get("requests", {}),
                                "limits": resources.get("limits", {}),
                            },
                        }
                    )

                deployments.append(
                    {
                        "name": name,
                        "kind": kind,
                        "source": self._relative(fpath),
                        "replicas": spec.get("replicas", 1),
                        "service_account": pod_spec.get("serviceAccountName", ""),
                        "automount_service_account_token": pod_spec.get(
                            "automountServiceAccountToken", True
                        ),
                        "containers": containers,
                    }
                )

        return {"deployments": deployments}
