"""Extract dependencies from go.mod and identify internal ODH dependencies."""

from __future__ import annotations

import logging
import re
from typing import Any

from extractors.base import BaseExtractor

logger = logging.getLogger(__name__)

ODH_MODULE_PREFIX = "github.com/opendatahub-io/"


class DependencyExtractor(BaseExtractor):
    """Extract Go module dependencies."""

    def extract(self) -> dict[str, Any]:
        go_mod_files = self.find_files(["go.mod", "**/go.mod"])
        go_modules: list[dict[str, str]] = []
        internal_odh: list[dict[str, str]] = []

        for fpath in go_mod_files:
            try:
                content = fpath.read_text(encoding="utf-8", errors="replace")
            except OSError:
                continue

            in_require = False
            for line in content.splitlines():
                stripped = line.strip()

                # Track require block
                if stripped.startswith("require ("):
                    in_require = True
                    continue
                if in_require and stripped == ")":
                    in_require = False
                    continue

                # Single-line require
                single_match = re.match(r"require\s+(\S+)\s+(\S+)", stripped)
                if single_match:
                    module, version = single_match.groups()
                    go_modules.append({"module": module, "version": version})
                    if module.startswith(ODH_MODULE_PREFIX):
                        component = module[len(ODH_MODULE_PREFIX):].split("/")[0]
                        internal_odh.append(
                            {
                                "component": component,
                                "interaction": f"Go module dependency: {module}",
                            }
                        )
                    continue

                # Lines inside require block
                if in_require:
                    # Skip indirect dependencies
                    if "// indirect" in stripped:
                        continue
                    parts = stripped.split()
                    if len(parts) >= 2:
                        module, version = parts[0], parts[1]
                        go_modules.append({"module": module, "version": version})
                        if module.startswith(ODH_MODULE_PREFIX):
                            component = module[len(ODH_MODULE_PREFIX):].split("/")[0]
                            internal_odh.append(
                                {
                                    "component": component,
                                    "interaction": f"Go module dependency: {module}",
                                }
                            )

        return {
            "dependencies": {
                "go_modules": go_modules,
                "internal_odh": internal_odh,
            }
        }
