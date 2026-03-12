#!/usr/bin/env python3
"""CLI entry point for the RHOAI Architecture Analyzer.

Usage:
    python3 analyze.py extract <repo-path> [--output <output.json>]
    python3 analyze.py render <architecture.json> [--output-dir <dir>] [--formats mermaid,security,c4,all]
    python3 analyze.py analyze <repo-path> [--output-dir <dir>]
    python3 analyze.py aggregate <results-dir> [--output-dir <dir>]
"""

from __future__ import annotations

import argparse
import json
import logging
import sys
from datetime import datetime, timezone
from pathlib import Path
from typing import Any

# Ensure the project root is on the path
PROJECT_ROOT = Path(__file__).resolve().parent
if str(PROJECT_ROOT) not in sys.path:
    sys.path.insert(0, str(PROJECT_ROOT))

from extractors import ALL_EXTRACTORS
from renderers import (
    C4Renderer,
    ComponentRenderer,
    DataflowRenderer,
    DependencyRenderer,
    RBACRenderer,
    SecurityNetworkRenderer,
)
from aggregator.aggregate import PlatformAggregator
from aggregator.render_platform import PlatformRenderer

VERSION = "0.1.0"

logger = logging.getLogger("analyze")

# Map format names to renderer classes
FORMAT_RENDERERS: dict[str, list[type]] = {
    "mermaid": [RBACRenderer, ComponentRenderer, DependencyRenderer, DataflowRenderer],
    "security": [SecurityNetworkRenderer],
    "c4": [C4Renderer],
}


def run_extract(repo_path: str) -> dict[str, Any]:
    """Run all extractors on a repository and return the combined data."""
    result: dict[str, Any] = {}

    # Derive component name from repo path
    repo_dir = Path(repo_path).resolve()
    component_name = repo_dir.name

    result["component"] = component_name
    result["repo"] = f"opendatahub-io/{component_name}"
    result["extracted_at"] = datetime.now(timezone.utc).isoformat()
    result["analyzer_version"] = VERSION

    for extractor_cls in ALL_EXTRACTORS:
        try:
            extractor = extractor_cls(repo_path)
            data = extractor.extract()
            result.update(data)
        except Exception as exc:
            logger.error(
                "Extractor %s failed: %s", extractor_cls.__name__, exc
            )
            # Continue with other extractors

    return result


def run_render(
    data: dict[str, Any],
    output_dir: str,
    formats: list[str] | None = None,
) -> list[str]:
    """Run renderers and write output files. Returns list of output paths."""
    out = Path(output_dir)
    out.mkdir(parents=True, exist_ok=True)
    output_files: list[str] = []

    # Determine which renderers to use
    if formats is None or "all" in formats:
        renderer_classes = [
            RBACRenderer,
            ComponentRenderer,
            SecurityNetworkRenderer,
            DependencyRenderer,
            C4Renderer,
            DataflowRenderer,
        ]
    else:
        renderer_classes = []
        for fmt in formats:
            renderer_classes.extend(FORMAT_RENDERERS.get(fmt, []))

    for renderer_cls in renderer_classes:
        try:
            renderer = renderer_cls(data)
            content = renderer.render()
            filepath = out / renderer.filename()
            filepath.write_text(content, encoding="utf-8")
            output_files.append(str(filepath))
            logger.info("Rendered: %s", filepath)
        except Exception as exc:
            logger.error(
                "Renderer %s failed: %s", renderer_cls.__name__, exc
            )

    return output_files


def cmd_extract(args: argparse.Namespace) -> int:
    """Handle the 'extract' command."""
    data = run_extract(args.repo_path)
    output = args.output or "component-architecture.json"
    output_path = Path(output)
    output_path.parent.mkdir(parents=True, exist_ok=True)
    output_path.write_text(
        json.dumps(data, indent=2, default=str), encoding="utf-8"
    )
    logger.info("Extracted architecture to: %s", output_path)
    return 0


def cmd_render(args: argparse.Namespace) -> int:
    """Handle the 'render' command."""
    arch_path = Path(args.architecture_json)
    if not arch_path.is_file():
        logger.error("Architecture file not found: %s", arch_path)
        return 1

    data = json.loads(arch_path.read_text(encoding="utf-8"))
    formats = args.formats.split(",") if args.formats else None
    output_dir = args.output_dir or str(arch_path.parent / "diagrams")
    files = run_render(data, output_dir, formats)
    logger.info("Rendered %d diagram(s) to %s", len(files), output_dir)
    return 0


def cmd_analyze(args: argparse.Namespace) -> int:
    """Handle the 'analyze' command (extract + render)."""
    data = run_extract(args.repo_path)
    output_dir = args.output_dir or "output"
    out = Path(output_dir)
    out.mkdir(parents=True, exist_ok=True)

    # Write JSON
    json_path = out / "component-architecture.json"
    json_path.write_text(
        json.dumps(data, indent=2, default=str), encoding="utf-8"
    )
    logger.info("Extracted architecture to: %s", json_path)

    # Render all diagrams
    diagrams_dir = out / "diagrams"
    files = run_render(data, str(diagrams_dir))
    logger.info("Rendered %d diagram(s) to %s", len(files), diagrams_dir)
    return 0


def cmd_aggregate(args: argparse.Namespace) -> int:
    """Handle the 'aggregate' command."""
    aggregator = PlatformAggregator(args.results_dir)
    platform_data = aggregator.aggregate()

    output_dir = args.output_dir or "platform-output"
    out = Path(output_dir)
    out.mkdir(parents=True, exist_ok=True)

    # Write platform JSON
    json_path = out / "platform-architecture.json"
    json_path.write_text(
        json.dumps(platform_data, indent=2, default=str), encoding="utf-8"
    )
    logger.info("Aggregated platform architecture to: %s", json_path)

    # Render platform diagrams
    renderer = PlatformRenderer(platform_data)
    diagrams = renderer.render_all()
    diagrams_dir = out / "diagrams"
    diagrams_dir.mkdir(parents=True, exist_ok=True)
    for filename, content in diagrams.items():
        filepath = diagrams_dir / filename
        filepath.write_text(content, encoding="utf-8")
        logger.info("Rendered: %s", filepath)

    return 0


def main() -> int:
    """Main CLI entry point."""
    parser = argparse.ArgumentParser(
        description="RHOAI Architecture Analyzer - Static analysis tool for Kubernetes/OpenShift components",
        formatter_class=argparse.RawDescriptionHelpFormatter,
    )
    parser.add_argument(
        "-v", "--verbose", action="store_true", help="Enable verbose logging"
    )
    subparsers = parser.add_subparsers(dest="command", help="Command to run")

    # extract
    p_extract = subparsers.add_parser(
        "extract", help="Extract architecture data from a repository"
    )
    p_extract.add_argument("repo_path", help="Path to the repository")
    p_extract.add_argument(
        "--output", "-o", help="Output JSON file (default: component-architecture.json)"
    )

    # render
    p_render = subparsers.add_parser(
        "render", help="Render diagrams from architecture JSON"
    )
    p_render.add_argument("architecture_json", help="Path to component-architecture.json")
    p_render.add_argument(
        "--output-dir", help="Output directory for diagrams"
    )
    p_render.add_argument(
        "--formats",
        help="Comma-separated formats: mermaid,security,c4,all (default: all)",
    )

    # analyze
    p_analyze = subparsers.add_parser(
        "analyze", help="Extract + render in one step"
    )
    p_analyze.add_argument("repo_path", help="Path to the repository")
    p_analyze.add_argument(
        "--output-dir", help="Output directory (default: output)"
    )

    # aggregate
    p_aggregate = subparsers.add_parser(
        "aggregate", help="Aggregate multiple component JSONs into platform view"
    )
    p_aggregate.add_argument(
        "results_dir", help="Directory containing per-component results"
    )
    p_aggregate.add_argument(
        "--output-dir", help="Output directory (default: platform-output)"
    )

    args = parser.parse_args()

    # Configure logging
    level = logging.DEBUG if args.verbose else logging.INFO
    logging.basicConfig(
        level=level,
        format="%(asctime)s [%(levelname)s] %(name)s: %(message)s",
        datefmt="%H:%M:%S",
    )

    if not args.command:
        parser.print_help()
        return 1

    handlers = {
        "extract": cmd_extract,
        "render": cmd_render,
        "analyze": cmd_analyze,
        "aggregate": cmd_aggregate,
    }

    handler = handlers.get(args.command)
    if handler is None:
        parser.print_help()
        return 1

    return handler(args)


if __name__ == "__main__":
    sys.exit(main())
