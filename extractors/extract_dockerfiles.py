"""Extract information from Dockerfiles/Containerfiles."""

from __future__ import annotations

import logging
import re
from typing import Any

from extractors.base import BaseExtractor

logger = logging.getLogger(__name__)

DOCKERFILE_PATTERNS = [
    "Dockerfile",
    "Dockerfile.*",
    "Containerfile",
    "Containerfile.*",
    "**/Dockerfile",
    "**/Dockerfile.*",
    "**/Containerfile",
    "**/Containerfile.*",
]


class DockerfileExtractor(BaseExtractor):
    """Extract metadata from Dockerfiles."""

    def extract(self) -> dict[str, Any]:
        files = self.find_files(DOCKERFILE_PATTERNS)
        dockerfiles: list[dict[str, Any]] = []

        for fpath in files:
            try:
                content = fpath.read_text(encoding="utf-8", errors="replace")
            except OSError:
                continue

            lines = content.splitlines()
            from_images: list[str] = []
            user: str = ""
            exposed_ports: list[int] = []
            issues: list[str] = []

            for line in lines:
                stripped = line.strip()
                if not stripped or stripped.startswith("#"):
                    continue

                # FROM instruction
                from_match = re.match(
                    r"FROM\s+(?:--platform=\S+\s+)?(\S+)(?:\s+[Aa][Ss]\s+\S+)?",
                    stripped,
                )
                if from_match:
                    image = from_match.group(1)
                    from_images.append(image)
                    # Check for unpinned tags
                    if image.endswith(":latest") or (":" not in image and "@" not in image):
                        issues.append(f"Unpinned base image: {image}")

                # USER instruction
                user_match = re.match(r"USER\s+(\S+)", stripped)
                if user_match:
                    user = user_match.group(1)

                # EXPOSE instruction
                expose_match = re.match(r"EXPOSE\s+(.*)", stripped)
                if expose_match:
                    for part in expose_match.group(1).split():
                        # Handle port/protocol format (e.g., 8080/tcp)
                        port_str = part.split("/")[0]
                        try:
                            exposed_ports.append(int(port_str))
                        except ValueError:
                            pass

            # Check for root user
            if user == "root" or user == "0":
                issues.append("Container runs as root user")
            elif not user:
                issues.append("No USER directive found (defaults to root)")

            stages = len(from_images)
            base_image = from_images[0] if from_images else ""

            dockerfiles.append(
                {
                    "path": self._relative(fpath),
                    "base_image": base_image,
                    "stages": stages,
                    "user": user,
                    "exposed_ports": exposed_ports,
                    "issues": issues,
                }
            )

        return {"dockerfiles": dockerfiles}
