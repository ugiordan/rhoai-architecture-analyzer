"""Render RBAC hierarchy as a Mermaid graph."""

from __future__ import annotations

from renderers.base import BaseRenderer


class RBACRenderer(BaseRenderer):
    """Render RBAC relationships as a Mermaid graph."""

    def filename(self) -> str:
        return "rbac.mmd"

    def render(self) -> str:
        rbac = self.data.get("rbac", {})
        if not rbac:
            return f"graph TD\n    note[No RBAC data for {self._escape_label(self.component)}]"

        lines = [
            "graph TD",
            f"    %% RBAC hierarchy for {self.component}",
        ]

        # Style classes
        lines.append("    classDef sa fill:#4a90d9,stroke:#2c5f8a,color:#fff")
        lines.append("    classDef role fill:#e8a838,stroke:#b07828,color:#fff")
        lines.append("    classDef resource fill:#5cb85c,stroke:#3d8b3d,color:#fff")
        lines.append("")

        node_counter = 0

        def next_id(prefix: str) -> str:
            nonlocal node_counter
            node_counter += 1
            return f"{prefix}_{node_counter}"

        # ClusterRoles and ClusterRoleBindings
        for binding in rbac.get("cluster_role_bindings", []):
            binding_name = binding.get("name", "")
            role_ref = binding.get("role_ref", "")
            binding_id = next_id("crb")
            role_id = self._sanitize_id(f"cr_{role_ref}")

            for subject in binding.get("subjects", []):
                subj_kind = subject.get("kind", "")
                subj_name = subject.get("name", "")
                subj_ns = subject.get("namespace", "")
                sa_label = f"{subj_name}"
                if subj_ns:
                    sa_label = f"{subj_name} ({subj_ns})"
                sa_id = next_id("sa")
                lines.append(
                    f'    {sa_id}["{self._escape_label(subj_kind)}: {self._escape_label(sa_label)}"]'
                    f" -->|bound via {self._escape_label(binding_name)}| "
                    f'{binding_id}["{self._escape_label(binding_name)}"]'
                )
                lines.append(f"    class {sa_id} sa")

            lines.append(
                f'    {binding_id} -->|grants| {role_id}["CR: {self._escape_label(role_ref)}"]'
            )
            lines.append(f"    class {role_id} role")

        # RoleBindings
        for binding in rbac.get("role_bindings", []):
            binding_name = binding.get("name", "")
            role_ref = binding.get("role_ref", "")
            binding_id = next_id("rb")
            role_id = self._sanitize_id(f"r_{role_ref}")

            for subject in binding.get("subjects", []):
                subj_name = subject.get("name", "")
                subj_kind = subject.get("kind", "")
                sa_id = next_id("sa")
                lines.append(
                    f'    {sa_id}["{self._escape_label(subj_kind)}: {self._escape_label(subj_name)}"]'
                    f" -->|bound via {self._escape_label(binding_name)}| "
                    f'{binding_id}["{self._escape_label(binding_name)}"]'
                )
                lines.append(f"    class {sa_id} sa")

            lines.append(
                f'    {binding_id} -->|grants| {role_id}["Role: {self._escape_label(role_ref)}"]'
            )
            lines.append(f"    class {role_id} role")

        # Render rules for each ClusterRole
        for role in rbac.get("cluster_roles", []):
            role_name = role.get("name", "")
            role_id = self._sanitize_id(f"cr_{role_name}")
            for rule in role.get("rules", []):
                api_groups = rule.get("apiGroups", [])
                resources = rule.get("resources", [])
                verbs = rule.get("verbs", [])
                for res in resources:
                    group = api_groups[0] if api_groups else "core"
                    if group == "":
                        group = "core"
                    res_id = next_id("res")
                    verb_str = ", ".join(verbs)
                    lines.append(
                        f'    {role_id} -->|{self._escape_label(verb_str)}| '
                        f'{res_id}["{self._escape_label(group)}: {self._escape_label(res)}"]'
                    )
                    lines.append(f"    class {res_id} resource")

        # Render rules for each Role
        for role in rbac.get("roles", []):
            role_name = role.get("name", "")
            role_id = self._sanitize_id(f"r_{role_name}")
            for rule in role.get("rules", []):
                api_groups = rule.get("apiGroups", [])
                resources = rule.get("resources", [])
                verbs = rule.get("verbs", [])
                for res in resources:
                    group = api_groups[0] if api_groups else "core"
                    if group == "":
                        group = "core"
                    res_id = next_id("res")
                    verb_str = ", ".join(verbs)
                    lines.append(
                        f'    {role_id} -->|{self._escape_label(verb_str)}| '
                        f'{res_id}["{self._escape_label(group)}: {self._escape_label(res)}"]'
                    )
                    lines.append(f"    class {res_id} resource")

        return "\n".join(lines)
