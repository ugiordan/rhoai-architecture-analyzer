"""Render platform-wide diagrams from aggregated data."""

from __future__ import annotations

from typing import Any


class PlatformRenderer:
    """Generate platform-wide diagrams from aggregated architecture data."""

    def __init__(self, platform_data: dict[str, Any]) -> None:
        self.data = platform_data

    @staticmethod
    def _sanitize_id(text: str) -> str:
        result = ""
        for ch in text:
            if ch.isalnum() or ch == "_":
                result += ch
            else:
                result += "_"
        if result and not result[0].isalpha():
            result = "n_" + result
        return result or "node"

    @staticmethod
    def _escape_label(text: str) -> str:
        return text.replace('"', "'").replace("<", "&lt;").replace(">", "&gt;")

    def render_dependency_graph(self) -> str:
        """Render full platform dependency graph in Mermaid."""
        lines = [
            "graph LR",
            "    %% Platform-wide dependency graph",
            "",
            "    classDef component fill:#3498db,stroke:#2980b9,color:#fff,stroke-width:2px",
            "    classDef gomod fill:#2ecc71,stroke:#27ae60,color:#fff",
            "    classDef crdwatch fill:#e74c3c,stroke:#c0392b,color:#fff",
            "",
        ]

        components = self.data.get("components", [])
        for comp in components:
            cid = self._sanitize_id(comp)
            lines.append(f'    {cid}["{self._escape_label(comp)}"]')
            lines.append(f"    class {cid} component")

        lines.append("")

        dep_graph = self.data.get("dependency_graph", [])
        for dep in dep_graph:
            from_id = self._sanitize_id(dep.get("from", ""))
            to_id = self._sanitize_id(dep.get("to", ""))
            dep_type = dep.get("type", "")
            if dep_type.startswith("watches-crd:"):
                crd_kind = dep_type.split(":", 1)[1]
                lines.append(
                    f'    {from_id} -->|"watches {self._escape_label(crd_kind)}"| {to_id}'
                )
            else:
                lines.append(f'    {from_id} -.->|"{self._escape_label(dep_type)}"| {to_id}')

        return "\n".join(lines)

    def render_crd_ownership(self) -> str:
        """Render CRD ownership map in Mermaid."""
        lines = [
            "graph TD",
            "    %% CRD ownership map",
            "",
            "    classDef component fill:#3498db,stroke:#2980b9,color:#fff",
            "    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff",
            "",
        ]

        crd_owners = self.data.get("crd_ownership", {})
        node_counter = 0

        for kind, owner in sorted(crd_owners.items()):
            node_counter += 1
            owner_id = self._sanitize_id(owner)
            crd_id = f"crd_{node_counter}"
            lines.append(f'    {owner_id}["{self._escape_label(owner)}"]')
            lines.append(f"    class {owner_id} component")
            lines.append(
                f'    {owner_id} -->|"defines"| {crd_id}{{{{"{self._escape_label(kind)}"}}}}'
            )
            lines.append(f"    class {crd_id} crd")

        return "\n".join(lines)

    def render_rbac_overview(self) -> str:
        """Render platform RBAC overview."""
        lines = [
            "graph TD",
            "    %% Platform RBAC overview",
            "",
            "    classDef component fill:#3498db,stroke:#2980b9,color:#fff",
            "    classDef role fill:#e8a838,stroke:#b07828,color:#fff",
            "",
        ]

        roles = self.data.get("rbac_cluster_roles", [])
        node_counter = 0

        for role in roles:
            node_counter += 1
            owner = role.get("owner", "")
            name = role.get("name", "")
            rules = role.get("rules", [])
            total_resources = sum(len(r.get("resources", [])) for r in rules)
            owner_id = self._sanitize_id(owner)
            role_id = f"role_{node_counter}"
            lines.append(f'    {owner_id}["{self._escape_label(owner)}"]')
            lines.append(f"    class {owner_id} component")
            lines.append(
                f'    {owner_id} --> {role_id}["{self._escape_label(name)}\\n'
                f'({total_resources} resources)"]'
            )
            lines.append(f"    class {role_id} role")

        return "\n".join(lines)

    def render_network_topology(self) -> str:
        """Render platform network topology in Mermaid."""
        lines = [
            "graph LR",
            "    %% Platform network topology",
            "",
            "    classDef component fill:#3498db,stroke:#2980b9,color:#fff",
            "    classDef service fill:#2ecc71,stroke:#27ae60,color:#fff",
            "",
        ]

        services = self.data.get("services", [])
        node_counter = 0

        for svc in services:
            node_counter += 1
            owner = svc.get("owner", "")
            name = svc.get("name", "")
            svc_type = svc.get("type", "ClusterIP")
            ports = svc.get("ports", [])
            port_str = ", ".join(
                f"{p.get('port', 0)}/{p.get('protocol', 'TCP')}" for p in ports
            )

            owner_id = self._sanitize_id(owner)
            svc_id = f"svc_{node_counter}"

            lines.append(f'    {owner_id}["{self._escape_label(owner)}"]')
            lines.append(f"    class {owner_id} component")
            lines.append(
                f'    {owner_id} --> {svc_id}["{self._escape_label(name)}\\n'
                f'{self._escape_label(svc_type)}: {self._escape_label(port_str)}"]'
            )
            lines.append(f"    class {svc_id} service")

        return "\n".join(lines)

    def render_all(self) -> dict[str, str]:
        """Render all platform diagrams, returning filename -> content mapping."""
        return {
            "platform-dependencies.mmd": self.render_dependency_graph(),
            "platform-crd-ownership.mmd": self.render_crd_ownership(),
            "platform-rbac-overview.mmd": self.render_rbac_overview(),
            "platform-network-topology.mmd": self.render_network_topology(),
        }
