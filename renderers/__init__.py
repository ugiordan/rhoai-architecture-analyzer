"""Renderers that produce diagrams from extracted architecture data."""

from renderers.render_rbac import RBACRenderer
from renderers.render_component import ComponentRenderer
from renderers.render_security_network import SecurityNetworkRenderer
from renderers.render_dependencies import DependencyRenderer
from renderers.render_c4 import C4Renderer
from renderers.render_dataflow import DataflowRenderer

ALL_RENDERERS = [
    RBACRenderer,
    ComponentRenderer,
    SecurityNetworkRenderer,
    DependencyRenderer,
    C4Renderer,
    DataflowRenderer,
]

__all__ = [
    "RBACRenderer",
    "ComponentRenderer",
    "SecurityNetworkRenderer",
    "DependencyRenderer",
    "C4Renderer",
    "DataflowRenderer",
    "ALL_RENDERERS",
]
