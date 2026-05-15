#!/usr/bin/env python3
"""Generate mkdocs nav entries for platform docs directories.

Scans site/docs/*-platform/ directories and produces YAML nav entries
for each platform's generated docs. Called by CI after docs generation
and before mkdocs build.

Usage:
    python3 scripts/generate-nav.py site/docs site/mkdocs.yml
"""

import os
import sys
import re


def find_platforms(docs_dir):
    """Find all *-platform/ directories under docs_dir/use-cases/."""
    use_cases_dir = os.path.join(docs_dir, "use-cases")
    platforms = []
    if not os.path.isdir(use_cases_dir):
        return platforms
    for entry in sorted(os.listdir(use_cases_dir)):
        if entry.endswith("-platform") and os.path.isdir(os.path.join(use_cases_dir, entry)):
            platforms.append(entry)
    return platforms


def build_platform_nav(docs_dir, platform_dir):
    """Build nav entries for a platform directory."""
    platform_path = os.path.join(docs_dir, "use-cases", platform_dir)
    # Derive display name and org: "odh-platform" -> ("ODH", "opendatahub-io"), "rhoai-platform" -> ("RHOAI", "red-hat-data-services")
    name_part = platform_dir.replace("-platform", "")
    if name_part == "odh":
        display_name = "ODH (opendatahub-io)"
    elif name_part == "rhoai":
        display_name = "RHOAI (red-hat-data-services)"
    else:
        display_name = name_part.upper() + " Platform"

    lines = []
    lines.append(f"    - {display_name}:")

    # Top-level files
    top_files = {
        "index.md": "Overview",
        "platform-architecture.md": "Platform Architecture",
        "network-topology.md": "Network Topology",
        "rbac-surface.md": "RBAC Surface",
        "secrets-inventory.md": "Secrets Inventory",
    }
    for fname, label in top_files.items():
        if os.path.exists(os.path.join(platform_path, fname)):
            lines.append(f"      - {label}: use-cases/{platform_dir}/{fname}")

    # Components
    components_dir = os.path.join(platform_path, "components")
    if os.path.isdir(components_dir):
        lines.append("      - Components:")
        for comp in sorted(os.listdir(components_dir)):
            comp_path = os.path.join(components_dir, comp)
            if not os.path.isdir(comp_path):
                continue
            # Display name: underscores to hyphens
            comp_display = comp.replace("_", "-")
            lines.append(f"        - {comp_display}:")

            # Standard doc files in order
            doc_files = {
                "index.md": "Overview",
                "network.md": "Network",
                "rbac.md": "RBAC",
                "security.md": "Security",
                "cache.md": "Cache",
                "dataflow.md": "Dataflow",
            }
            for fname, label in doc_files.items():
                if os.path.exists(os.path.join(comp_path, fname)):
                    lines.append(f"          - {label}: use-cases/{platform_dir}/components/{comp}/{fname}")

    return lines


def update_mkdocs(mkdocs_path, platform_nav_lines):
    """Replace platform nav sections in mkdocs.yml."""
    with open(mkdocs_path, "r") as f:
        content = f.read()

    # Replace everything between "Use Cases:" and "Contributing:" with platform nav.
    # Works regardless of where Use Cases sits in the nav order.
    pattern = r"(  - Use Cases:\n    - use-cases/index\.md\n).*?(\n  - Contributing:)"
    replacement = r"\1" + "\n".join(platform_nav_lines) + r"\2"
    match = re.search(pattern, content, flags=re.DOTALL)
    if not match:
        print("WARNING: Could not find nav section to replace (check mkdocs.yml nav order)", file=sys.stderr)
        return False

    new_content = re.sub(pattern, replacement, content, flags=re.DOTALL)

    if new_content == content:
        print(f"  {mkdocs_path} already up to date")
        return True

    with open(mkdocs_path, "w") as f:
        f.write(new_content)
    return True


def main():
    if len(sys.argv) != 3:
        print(f"Usage: {sys.argv[0]} <docs_dir> <mkdocs_yml>", file=sys.stderr)
        sys.exit(1)

    docs_dir = sys.argv[1]
    mkdocs_path = sys.argv[2]

    platforms = find_platforms(docs_dir)
    if not platforms:
        print("No *-platform/ directories found, skipping nav generation")
        return

    all_nav_lines = []
    for platform in platforms:
        nav_lines = build_platform_nav(docs_dir, platform)
        all_nav_lines.extend(nav_lines)
        comp_count = len([d for d in os.listdir(os.path.join(docs_dir, "use-cases", platform, "components"))
                         if os.path.isdir(os.path.join(docs_dir, "use-cases", platform, "components", d))]) \
            if os.path.isdir(os.path.join(docs_dir, "use-cases", platform, "components")) else 0
        print(f"  {platform}: {comp_count} components")

    if update_mkdocs(mkdocs_path, all_nav_lines):
        print(f"Updated {mkdocs_path} with {len(platforms)} platform(s)")
    else:
        sys.exit(1)


if __name__ == "__main__":
    main()
