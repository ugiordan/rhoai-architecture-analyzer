"""Render component architecture as a Mermaid graph."""

from __future__ import annotations

from renderers.base import BaseRenderer


class ComponentRenderer(BaseRenderer):
    """Render component architecture showing CRDs, controllers, and resources."""

    def filename(self) -> str:
        return "component.mmd"

    def render(self) -> str:
        lines = [
            "graph LR",
            f"    %% Component architecture for {self.component}",
            "",
            "    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff",
            "    classDef controller fill:#3498db,stroke:#2980b9,color:#fff",
            "    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff",
            "    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff",
            "    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff",
            "",
        ]

        node_counter = 0

        def next_id(prefix: str) -> str:
            nonlocal node_counter
            node_counter += 1
            return f"{prefix}_{node_counter}"

        # CRDs (primary watched resources)
        crds = self.data.get("crds", [])
        watches = self.data.get("controller_watches", [])
        deployments = self.data.get("deployments", [])

        # Build sets for For/Owns/Watches
        for_kinds: set[str] = set()
        owns_kinds: set[str] = set()
        watches_kinds: set[str] = set()
        for w in watches:
            gvk = w.get("gvk", "")
            kind = gvk.rsplit("/", 1)[-1] if "/" in gvk else gvk
            wtype = w.get("type", "")
            if wtype == "For":
                for_kinds.add(kind)
            elif wtype == "Owns":
                owns_kinds.add(kind)
            elif wtype == "Watches":
                watches_kinds.add(kind)

        # Controller subgraph
        lines.append(f'    subgraph controller["{self._escape_label(self.component)} Controller"]')
        for dep in deployments:
            dep_name = dep.get("name", "")
            dep_id = next_id("dep")
            lines.append(f'        {dep_id}["{self._escape_label(dep_name)}"]')
            lines.append(f"        class {dep_id} controller")
        if not deployments:
            ctrl_id = next_id("ctrl")
            lines.append(f'        {ctrl_id}["Controller"]')
            lines.append(f"        class {ctrl_id} controller")
        lines.append("    end")
        lines.append("")

        # CRDs - watched via For
        for crd in crds:
            kind = crd.get("kind", "")
            group = crd.get("group", "")
            version = crd.get("version", "")
            crd_id = self._sanitize_id(f"crd_{kind}")
            label = f"{kind}\\n{group}/{version}"
            lines.append(f'    {crd_id}{{{{"{self._escape_label(label)}"}}}}')
            lines.append(f"    class {crd_id} crd")
            if kind in for_kinds:
                lines.append(f'    {crd_id} -->|"For (reconciles)"| controller')

        # Owned resources
        for kind in sorted(owns_kinds):
            res_id = next_id("owned")
            lines.append(f'    controller -->|"Owns"| {res_id}["{self._escape_label(kind)}"]')
            lines.append(f"    class {res_id} owned")

        # Watched resources (not For, not Owns)
        for kind in sorted(watches_kinds):
            res_id = next_id("watch")
            lines.append(f'    {res_id}["{self._escape_label(kind)}"] -->|"Watches"| controller')
            lines.append(f"    class {res_id} external")

        # Internal ODH dependencies
        deps = self.data.get("dependencies", {})
        for odh_dep in deps.get("internal_odh", []):
            component = odh_dep.get("component", "")
            dep_id = next_id("odh")
            lines.append(
                f'    controller -.->|"depends on"| {dep_id}["{self._escape_label(component)}"]'
            )
            lines.append(f"    class {dep_id} dep")

        return "\n".join(lines)
