"""Render dependency graph as a Mermaid diagram."""

from __future__ import annotations

from renderers.base import BaseRenderer


class DependencyRenderer(BaseRenderer):
    """Render a dependency graph in Mermaid."""

    def filename(self) -> str:
        return "dependencies.mmd"

    def render(self) -> str:
        deps = self.data.get("dependencies", {})
        go_modules = deps.get("go_modules", [])
        internal_odh = deps.get("internal_odh", [])

        lines = [
            "graph LR",
            f"    %% Dependency graph for {self.component}",
            "",
            "    classDef component fill:#3498db,stroke:#2980b9,color:#fff,stroke-width:3px",
            "    classDef internal fill:#2ecc71,stroke:#27ae60,color:#fff",
            "    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#333",
            "",
        ]

        comp_id = self._sanitize_id(f"comp_{self.component}")
        lines.append(f'    {comp_id}["{self._escape_label(self.component)}"]')
        lines.append(f"    class {comp_id} component")
        lines.append("")

        node_counter = 0

        def next_id(prefix: str) -> str:
            nonlocal node_counter
            node_counter += 1
            return f"{prefix}_{node_counter}"

        # Internal ODH dependencies
        if internal_odh:
            lines.append("    %% Internal ODH dependencies")
            for odh in internal_odh:
                comp_name = odh.get("component", "")
                interaction = odh.get("interaction", "")
                dep_id = next_id("odh")
                lines.append(f'    {dep_id}["{self._escape_label(comp_name)}"]')
                lines.append(f"    class {dep_id} internal")
                lines.append(
                    f'    {comp_id} -->|"{self._escape_label(interaction)}"| {dep_id}'
                )
            lines.append("")

        # Key external dependencies (filter to notable ones)
        notable_prefixes = [
            "sigs.k8s.io/controller-runtime",
            "k8s.io/api",
            "k8s.io/apimachinery",
            "k8s.io/client-go",
            "github.com/operator-framework",
            "github.com/prometheus",
            "google.golang.org/grpc",
            "github.com/go-logr",
        ]
        if go_modules:
            lines.append("    %% Key external dependencies")
            seen: set[str] = set()
            for mod in go_modules:
                module = mod.get("module", "")
                version = mod.get("version", "")
                # Skip internal ODH (already rendered)
                if module.startswith("github.com/opendatahub-io/"):
                    continue
                # Only show notable dependencies
                is_notable = any(module.startswith(p) for p in notable_prefixes)
                if not is_notable:
                    continue
                # Deduplicate by base module
                base = module.split("/")[0] + "/" + module.split("/")[1] if "/" in module else module
                if base in seen:
                    continue
                seen.add(base)
                dep_id = next_id("ext")
                short_name = module.rsplit("/", 1)[-1]
                lines.append(
                    f'    {dep_id}["{self._escape_label(short_name)}\\n{self._escape_label(version)}"]'
                )
                lines.append(f"    class {dep_id} external")
                lines.append(f"    {comp_id} --> {dep_id}")

        return "\n".join(lines)
