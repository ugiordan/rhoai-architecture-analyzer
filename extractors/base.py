"""Base extractor with common YAML/file parsing utilities."""

from __future__ import annotations

import logging
import re
from abc import ABC, abstractmethod
from pathlib import Path
from typing import Any

import yaml

logger = logging.getLogger(__name__)


class BaseExtractor(ABC):
    """Abstract base class for all extractors."""

    def __init__(self, repo_path: str) -> None:
        self.repo_path = Path(repo_path).resolve()
        if not self.repo_path.is_dir():
            raise ValueError(f"Repository path does not exist: {self.repo_path}")

    @abstractmethod
    def extract(self) -> dict[str, Any]:
        """Run the extraction and return the data dictionary."""

    def _relative(self, path: Path) -> str:
        """Return a path relative to the repo root."""
        try:
            return str(path.relative_to(self.repo_path))
        except ValueError:
            return str(path)

    # ----- File discovery helpers -----

    def find_yaml_files(self, patterns: list[str]) -> list[Path]:
        """Find YAML files matching any of the given glob patterns."""
        found: list[Path] = []
        for pattern in patterns:
            found.extend(sorted(self.repo_path.glob(pattern)))
        # Deduplicate while preserving order
        seen: set[Path] = set()
        result: list[Path] = []
        for p in found:
            resolved = p.resolve()
            if resolved not in seen and resolved.is_file():
                seen.add(resolved)
                result.append(p)
        return result

    def find_go_files(self, patterns: list[str] | None = None) -> list[Path]:
        """Find Go source files, optionally filtered by glob patterns."""
        if patterns:
            found: list[Path] = []
            for pat in patterns:
                found.extend(sorted(self.repo_path.glob(pat)))
            seen: set[Path] = set()
            result: list[Path] = []
            for p in found:
                resolved = p.resolve()
                if resolved not in seen and p.is_file():
                    seen.add(resolved)
                    result.append(p)
            return result
        return sorted(self.repo_path.rglob("*.go"))

    def find_files(self, patterns: list[str]) -> list[Path]:
        """Find files matching any of the given glob patterns."""
        found: list[Path] = []
        for pattern in patterns:
            found.extend(sorted(self.repo_path.glob(pattern)))
        seen: set[Path] = set()
        result: list[Path] = []
        for p in found:
            resolved = p.resolve()
            if resolved not in seen and resolved.is_file():
                seen.add(resolved)
                result.append(p)
        return result

    # ----- YAML parsing -----

    def parse_yaml_safe(self, path: Path) -> list[dict[str, Any]]:
        """Parse a YAML file, returning a list of documents.

        Skips files that contain unresolvable Helm template syntax.
        Returns an empty list on parse errors.
        """
        try:
            content = path.read_text(encoding="utf-8", errors="replace")
        except OSError as exc:
            logger.warning("Cannot read %s: %s", path, exc)
            return []

        # If file contains Helm template directives that would break parsing,
        # attempt to parse anyway but gracefully handle failures.
        has_helm_templates = "{{" in content and "}}" in content

        docs: list[dict[str, Any]] = []
        try:
            for doc in yaml.safe_load_all(content):
                if isinstance(doc, dict):
                    docs.append(doc)
        except yaml.YAMLError:
            if has_helm_templates:
                logger.debug("Skipping Helm template file: %s", path)
            else:
                logger.warning("Failed to parse YAML: %s", path)
            return []

        return docs

    # ----- Regex scanning -----

    def regex_scan(
        self, files: list[Path], pattern: str
    ) -> list[dict[str, Any]]:
        """Scan files for regex matches, returning file/line/match info."""
        compiled = re.compile(pattern)
        results: list[dict[str, Any]] = []
        for fpath in files:
            try:
                lines = fpath.read_text(encoding="utf-8", errors="replace").splitlines()
            except OSError:
                continue
            for line_no, line in enumerate(lines, start=1):
                match = compiled.search(line)
                if match:
                    results.append(
                        {
                            "file": self._relative(fpath),
                            "line": line_no,
                            "match": match.group(0),
                            "groups": match.groups(),
                        }
                    )
        return results
