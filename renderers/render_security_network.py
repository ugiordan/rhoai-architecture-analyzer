"""Render ASCII security and network diagram."""

from __future__ import annotations

from renderers.base import BaseRenderer


class SecurityNetworkRenderer(BaseRenderer):
    """Render a layered ASCII security and network diagram."""

    def filename(self) -> str:
        return "security-network.txt"

    def render(self) -> str:
        lines: list[str] = []
        width = 80
        component = self.component

        def header(title: str) -> None:
            lines.append("")
            lines.append("=" * width)
            lines.append(f"  {title}")
            lines.append("=" * width)

        def section(title: str) -> None:
            lines.append("")
            lines.append(f"--- {title} {'-' * (width - len(title) - 5)}")

        lines.append("+" + "=" * (width - 2) + "+")
        lines.append(f"|{'SECURITY & NETWORK ARCHITECTURE':^{width - 2}}|")
        lines.append(f"|{component:^{width - 2}}|")
        lines.append("+" + "=" * (width - 2) + "+")

        # ---- Services & Ports ----
        header("NETWORK TOPOLOGY")
        services = self.data.get("services", [])
        if services:
            section("Services")
            for svc in services:
                name = svc.get("name", "")
                svc_type = svc.get("type", "ClusterIP")
                lines.append(f"  [{svc_type}] {name}")
                for port in svc.get("ports", []):
                    port_name = port.get("name", "")
                    port_num = port.get("port", 0)
                    target = port.get("targetPort", 0)
                    protocol = port.get("protocol", "TCP")
                    tls_hint = " (TLS)" if "https" in port_name.lower() or port_num == 443 else ""
                    lines.append(
                        f"    Port: {port_num} -> {target}/{protocol} [{port_name}]{tls_hint}"
                    )
        else:
            section("Services")
            lines.append("  (none found)")

        # ---- Network Policies ----
        policies = self.data.get("network_policies", [])
        section("Network Policies")
        if policies:
            for pol in policies:
                name = pol.get("name", "")
                types = ", ".join(pol.get("policy_types", []))
                selector = pol.get("pod_selector", {})
                sel_str = ", ".join(f"{k}={v}" for k, v in selector.items()) or "(all pods)"
                lines.append(f"  Policy: {name}")
                lines.append(f"    Selector: {sel_str}")
                lines.append(f"    Types: {types}")
                for ing in pol.get("ingress_rules", []):
                    for p in ing.get("ports", []):
                        lines.append(
                            f"    INGRESS: port {p.get('port', '?')}/{p.get('protocol', 'TCP')}"
                        )
                    if not ing.get("from"):
                        lines.append("    INGRESS: from ALL sources")
                for eg in pol.get("egress_rules", []):
                    for p in eg.get("ports", []):
                        lines.append(
                            f"    EGRESS: port {p.get('port', '?')}/{p.get('protocol', 'TCP')}"
                        )
        else:
            lines.append("  (none found - all traffic allowed by default)")

        # ---- RBAC Summary ----
        header("RBAC SUMMARY")
        rbac = self.data.get("rbac", {})
        cr_count = len(rbac.get("cluster_roles", []))
        crb_count = len(rbac.get("cluster_role_bindings", []))
        r_count = len(rbac.get("roles", []))
        rb_count = len(rbac.get("role_bindings", []))
        marker_count = len(rbac.get("kubebuilder_markers", []))

        lines.append(f"  ClusterRoles:        {cr_count}")
        lines.append(f"  ClusterRoleBindings: {crb_count}")
        lines.append(f"  Roles:               {r_count}")
        lines.append(f"  RoleBindings:        {rb_count}")
        lines.append(f"  Kubebuilder markers: {marker_count}")

        # List all ClusterRoles with resource counts
        for cr in rbac.get("cluster_roles", []):
            name = cr.get("name", "")
            rules = cr.get("rules", [])
            total_resources = sum(len(r.get("resources", [])) for r in rules)
            lines.append(f"    CR: {name} ({total_resources} resource types)")

        # ---- Secrets Inventory ----
        header("SECRETS INVENTORY")
        secrets = self.data.get("secrets_referenced", [])
        if secrets:
            for secret in secrets:
                name = secret.get("name", "")
                stype = secret.get("type", "Opaque")
                refs = secret.get("referenced_by", [])
                prov = secret.get("provisioned_by", "")
                lines.append(f"  Secret: {name}")
                lines.append(f"    Type: {stype}")
                lines.append(f"    Provisioned by: {prov or 'unknown'}")
                for ref in refs:
                    lines.append(f"    Referenced by: {ref}")
        else:
            lines.append("  (no secret references found)")

        # ---- Deployment Security Controls ----
        header("DEPLOYMENT SECURITY CONTROLS")
        deployments = self.data.get("deployments", [])
        if deployments:
            for dep in deployments:
                name = dep.get("name", "")
                sa = dep.get("service_account", "")
                automount = dep.get("automount_service_account_token", True)
                lines.append(f"  Deployment: {name}")
                lines.append(f"    Service Account: {sa or '(default)'}")
                lines.append(f"    Automount SA Token: {automount}")
                for container in dep.get("containers", []):
                    c_name = container.get("name", "")
                    sc = container.get("security_context", {})
                    lines.append(f"    Container: {c_name}")
                    if sc:
                        for key, val in sc.items():
                            if isinstance(val, dict):
                                val_str = ", ".join(f"{k}={v}" for k, v in val.items())
                            elif isinstance(val, list):
                                val_str = ", ".join(str(v) for v in val)
                            else:
                                val_str = str(val)
                            lines.append(f"      {key}: {val_str}")
                    else:
                        lines.append("      (no security context)")
        else:
            lines.append("  (no deployments found)")

        # ---- Dockerfile Security ----
        header("DOCKERFILE SECURITY")
        dockerfiles = self.data.get("dockerfiles", [])
        if dockerfiles:
            for df in dockerfiles:
                path = df.get("path", "")
                base = df.get("base_image", "")
                user = df.get("user", "")
                stages = df.get("stages", 0)
                issues = df.get("issues", [])
                lines.append(f"  {path}")
                lines.append(f"    Base image: {base}")
                lines.append(f"    Stages: {stages}")
                lines.append(f"    User: {user or '(not set)'}")
                if issues:
                    for issue in issues:
                        lines.append(f"    [!] {issue}")
        else:
            lines.append("  (no Dockerfiles found)")

        lines.append("")
        lines.append("+" + "=" * (width - 2) + "+")

        return "\n".join(lines)
