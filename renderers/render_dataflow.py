"""Render dataflow sequence diagram in Mermaid."""

from __future__ import annotations

from renderers.base import BaseRenderer


class DataflowRenderer(BaseRenderer):
    """Render a Mermaid sequence diagram showing static connections."""

    def filename(self) -> str:
        return "dataflow.mmd"

    def render(self) -> str:
        lines = [
            "sequenceDiagram",
            f"    %% Static dataflow for {self.component}",
            "",
        ]

        deployments = self.data.get("deployments", [])
        services = self.data.get("services", [])
        watches = self.data.get("controller_watches", [])
        crds = self.data.get("crds", [])

        # Collect participants
        participants: list[str] = []

        # Add K8s API
        participants.append("KubernetesAPI")
        lines.append("    participant KubernetesAPI as Kubernetes API")

        # Add deployments as participants
        for dep in deployments:
            name = dep.get("name", "")
            pid = self._sanitize_id(name)
            if pid not in participants:
                participants.append(pid)
                lines.append(f"    participant {pid} as {self._escape_label(name)}")

        # Add a generic controller if no deployments
        if not deployments:
            ctrl_id = self._sanitize_id(self.component)
            participants.append(ctrl_id)
            lines.append(f"    participant {ctrl_id} as {self._escape_label(self.component)}")

        lines.append("")

        controller_id = (
            self._sanitize_id(deployments[0]["name"])
            if deployments
            else self._sanitize_id(self.component)
        )

        # Controller watches (For)
        for w in watches:
            wtype = w.get("type", "")
            gvk = w.get("gvk", "")
            kind = gvk.rsplit("/", 1)[-1] if "/" in gvk else gvk
            if wtype == "For":
                lines.append(
                    f"    KubernetesAPI->>+{controller_id}: Watch {self._escape_label(kind)} (reconcile)"
                )
            elif wtype == "Owns":
                lines.append(
                    f"    {controller_id}->>KubernetesAPI: Create/Update {self._escape_label(kind)}"
                )
            elif wtype == "Watches":
                lines.append(
                    f"    KubernetesAPI-->>+{controller_id}: Watch {self._escape_label(kind)} (informer)"
                )

        # Service connections
        if services:
            lines.append("")
            lines.append(f"    Note over {controller_id}: Exposed Services")
            for svc in services:
                name = svc.get("name", "")
                for port in svc.get("ports", []):
                    port_num = port.get("port", 0)
                    port_name = port.get("name", "")
                    protocol = port.get("protocol", "TCP")
                    lines.append(
                        f"    Note right of {controller_id}: "
                        f"{self._escape_label(name)}:{port_num}/{protocol} [{self._escape_label(port_name)}]"
                    )

        # CRD definitions
        if crds:
            lines.append("")
            lines.append(f"    Note over KubernetesAPI: Defined CRDs")
            for crd in crds:
                kind = crd.get("kind", "")
                group = crd.get("group", "")
                version = crd.get("version", "")
                lines.append(
                    f"    Note right of KubernetesAPI: "
                    f"{self._escape_label(kind)} ({self._escape_label(group)}/{self._escape_label(version)})"
                )

        return "\n".join(lines)
