"""Base renderer class."""

from __future__ import annotations

from abc import ABC, abstractmethod
from typing import Any


class BaseRenderer(ABC):
    """Abstract base class for diagram renderers."""

    def __init__(self, data: dict[str, Any]) -> None:
        self.data = data
        self.component = data.get("component", "unknown")

    @abstractmethod
    def render(self) -> str:
        """Produce the diagram as a string."""

    @abstractmethod
    def filename(self) -> str:
        """Return the suggested output filename."""

    @staticmethod
    def _sanitize_id(text: str) -> str:
        """Sanitize a string for use as a Mermaid node ID."""
        # Replace non-alphanumeric chars with underscores
        result = ""
        for ch in text:
            if ch.isalnum() or ch == "_":
                result += ch
            else:
                result += "_"
        # Ensure it starts with a letter
        if result and not result[0].isalpha():
            result = "n_" + result
        return result or "node"

    @staticmethod
    def _escape_label(text: str) -> str:
        """Escape special characters for Mermaid labels."""
        return text.replace('"', "'").replace("<", "&lt;").replace(">", "&gt;")
