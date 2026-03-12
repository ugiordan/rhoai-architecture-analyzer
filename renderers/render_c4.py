"""Render C4 context diagram in Structurizr DSL."""

from __future__ import annotations

from renderers.base import BaseRenderer


class C4Renderer(BaseRenderer):
    """Render a C4 context diagram in Structurizr DSL."""

    def filename(self) -> str:
        return "c4-context.dsl"

    def render(self) -> str:
        lines: list[str] = []
        component = self.component
        sanitized_component = self._sanitize_id(component)

        lines.append("workspace {")
        lines.append("")
        lines.append("    model {")

        # Persons
        lines.append('        admin = person "Platform Admin" "Manages the OpenShift AI platform"')
        lines.append('        user = person "Data Scientist" "Uses ML tools and services"')
        lines.append("")

        # System boundary
        lines.append(
            f'        {sanitized_component}_system = softwareSystem "{self._dsl_escape(component)}" {{'
        )

        # Containers from deployments
        deployments = self.data.get("deployments", [])
        container_ids: list[str] = []
        for dep in deployments:
            name = dep.get("name", "")
            dep_id = self._sanitize_id(f"container_{name}")
            container_ids.append(dep_id)
            image = ""
            ports_str = ""
            for c in dep.get("containers", []):
                image = c.get("image", image)
                for p in c.get("ports", []):
                    port = p.get("containerPort", 0)
                    proto = p.get("protocol", "TCP")
                    ports_str += f" port {port}/{proto}"

            desc = f"Container image: {image}" if image else "Controller container"
            tech = ports_str.strip() if ports_str else "Kubernetes"
            lines.append(
                f'            {dep_id} = container "{self._dsl_escape(name)}" '
                f'"{self._dsl_escape(desc)}" "{self._dsl_escape(tech)}"'
            )

        if not deployments:
            ctrl_id = self._sanitize_id(f"container_{component}")
            container_ids.append(ctrl_id)
            lines.append(
                f'            {ctrl_id} = container "{self._dsl_escape(component)}" '
                f'"Main controller" "Kubernetes"'
            )

        lines.append("        }")
        lines.append("")

        # External systems from dependencies
        deps = self.data.get("dependencies", {})
        ext_ids: list[str] = []
        for odh in deps.get("internal_odh", []):
            comp_name = odh.get("component", "")
            ext_id = self._sanitize_id(f"ext_{comp_name}")
            ext_ids.append(ext_id)
            lines.append(
                f'        {ext_id} = softwareSystem "{self._dsl_escape(comp_name)}" '
                f'"ODH component" "Existing System"'
            )

        # Kubernetes API as external system
        k8s_id = "kubernetes_api"
        lines.append(
            f'        {k8s_id} = softwareSystem "Kubernetes API" '
            f'"Cluster API server" "Existing System"'
        )
        lines.append("")

        # Relationships
        lines.append(f"        admin -> {sanitized_component}_system \"Manages\"")
        lines.append(f"        user -> {sanitized_component}_system \"Uses\"")

        for cid in container_ids:
            lines.append(f"        {cid} -> {k8s_id} \"API calls\" \"HTTPS\"")

        for ext_id in ext_ids:
            lines.append(
                f"        {sanitized_component}_system -> {ext_id} \"Depends on\""
            )

        # CRD relationships
        crds = self.data.get("crds", [])
        for crd in crds:
            kind = crd.get("kind", "")
            group = crd.get("group", "")
            lines.append(
                f"        admin -> {sanitized_component}_system "
                f'"Creates {self._dsl_escape(kind)} ({self._dsl_escape(group)})" "YAML/kubectl"'
            )
            break  # One relationship is enough for context diagram

        lines.append("    }")
        lines.append("")

        # Views
        lines.append("    views {")
        lines.append(
            f"        systemContext {sanitized_component}_system \"SystemContext\" {{"
        )
        lines.append("            include *")
        lines.append("            autoLayout")
        lines.append("        }")
        lines.append(
            f"        container {sanitized_component}_system \"Containers\" {{"
        )
        lines.append("            include *")
        lines.append("            autoLayout")
        lines.append("        }")
        lines.append("    }")
        lines.append("")
        lines.append("}")

        return "\n".join(lines)

    @staticmethod
    def _dsl_escape(text: str) -> str:
        """Escape special characters for Structurizr DSL strings."""
        return text.replace('"', '\\"').replace("\n", " ")
