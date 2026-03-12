"""Extract controller watch configurations from Go source code."""

from __future__ import annotations

import logging
import re
from typing import Any

from extractors.base import BaseExtractor

logger = logging.getLogger(__name__)

GO_FILE_PATTERNS = [
    "**/*_controller.go",
    "**/*_reconciler.go",
    "**/setup.go",
    "**/controller.go",
    "**/reconciler.go",
]

# Regex for import aliases: `alias "path/to/package"`
IMPORT_ALIAS_RE = re.compile(r'(\w+)\s+"([^"]+)"')

# Watch patterns in SetupWithManager
# The dot before For/Owns/Watches may be on the previous line in chained calls,
# so we match with an optional leading dot.
FOR_RE = re.compile(r"\.?For\(\s*&(\w+)\.(\w+)\{")
OWNS_RE = re.compile(r"\.?Owns\(\s*&(\w+)\.(\w+)\{")
WATCHES_RE = re.compile(r"\.?Watches\(\s*&?(?:source\.Kind\{Type:\s*&)?(\w+)\.(\w+)\{")

# Known Go package to API group mappings
KNOWN_GROUPS: dict[str, str] = {
    "k8s.io/api/core/v1": "/v1",
    "k8s.io/api/apps/v1": "apps/v1",
    "k8s.io/api/batch/v1": "batch/v1",
    "k8s.io/api/networking/v1": "networking.k8s.io/v1",
    "k8s.io/api/rbac/v1": "rbac.authorization.k8s.io/v1",
    "k8s.io/api/policy/v1": "policy/v1",
    "k8s.io/api/autoscaling/v1": "autoscaling/v1",
    "k8s.io/api/autoscaling/v2": "autoscaling/v2",
    "k8s.io/api/admissionregistration/v1": "admissionregistration.k8s.io/v1",
    "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1": "apiextensions.k8s.io/v1",
}


def _resolve_import_alias(alias: str, imports: dict[str, str]) -> str:
    """Resolve an import alias to an API group/version string."""
    import_path = imports.get(alias, "")
    if not import_path:
        return alias

    # Check known mappings
    if import_path in KNOWN_GROUPS:
        return KNOWN_GROUPS[import_path]

    # Try to infer group/version from import path
    # e.g., github.com/opendatahub-io/opendatahub-operator/v2/api/dscinitialization/v1
    # -> dscinitialization.opendatahub.io/v1
    parts = import_path.rstrip("/").split("/")
    if len(parts) >= 2:
        version = parts[-1]
        if re.match(r"v\d+", version):
            group_part = parts[-2] if len(parts) >= 3 else ""
            # Build a rough group
            return f"{group_part}/{version}" if group_part else f"/{version}"

    return import_path


def _parse_imports(content: str) -> dict[str, str]:
    """Parse Go import block to build alias -> import path mapping."""
    imports: dict[str, str] = {}
    # Find import blocks
    import_blocks = re.findall(r"import\s*\((.*?)\)", content, re.DOTALL)
    for block in import_blocks:
        for line in block.splitlines():
            line = line.strip()
            if not line or line.startswith("//"):
                continue
            match = IMPORT_ALIAS_RE.match(line)
            if match:
                alias, path = match.groups()
                imports[alias] = path
            else:
                # Non-aliased import: `"path/to/package"`
                path_match = re.match(r'"([^"]+)"', line)
                if path_match:
                    path = path_match.group(1)
                    # Use last path component as implicit alias
                    last = path.rstrip("/").rsplit("/", 1)[-1]
                    imports[last] = path
    return imports


class ControllerWatchExtractor(BaseExtractor):
    """Extract For/Owns/Watches from Go controller code."""

    def extract(self) -> dict[str, Any]:
        files = self.find_go_files(GO_FILE_PATTERNS)
        watches: list[dict[str, Any]] = []

        for fpath in files:
            try:
                content = fpath.read_text(encoding="utf-8", errors="replace")
            except OSError:
                continue

            imports = _parse_imports(content)
            lines = content.splitlines()

            for line_no, line in enumerate(lines, start=1):
                source = f"{self._relative(fpath)}:{line_no}"

                for match in FOR_RE.finditer(line):
                    alias, kind = match.groups()
                    gv = _resolve_import_alias(alias, imports)
                    watches.append(
                        {"type": "For", "gvk": f"{gv}/{kind}", "source": source}
                    )

                for match in OWNS_RE.finditer(line):
                    alias, kind = match.groups()
                    gv = _resolve_import_alias(alias, imports)
                    watches.append(
                        {"type": "Owns", "gvk": f"{gv}/{kind}", "source": source}
                    )

                for match in WATCHES_RE.finditer(line):
                    alias, kind = match.groups()
                    gv = _resolve_import_alias(alias, imports)
                    watches.append(
                        {"type": "Watches", "gvk": f"{gv}/{kind}", "source": source}
                    )

        return {"controller_watches": watches}
