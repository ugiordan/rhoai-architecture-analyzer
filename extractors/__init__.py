"""Extractors for Kubernetes/OpenShift component architecture data."""

from extractors.extract_crds import CRDExtractor
from extractors.extract_rbac import RBACExtractor
from extractors.extract_services import ServiceExtractor
from extractors.extract_deployments import DeploymentExtractor
from extractors.extract_network_policies import NetworkPolicyExtractor
from extractors.extract_controller_watches import ControllerWatchExtractor
from extractors.extract_dependencies import DependencyExtractor
from extractors.extract_secrets import SecretExtractor
from extractors.extract_helm import HelmExtractor
from extractors.extract_dockerfiles import DockerfileExtractor

ALL_EXTRACTORS = [
    CRDExtractor,
    RBACExtractor,
    ServiceExtractor,
    DeploymentExtractor,
    NetworkPolicyExtractor,
    ControllerWatchExtractor,
    DependencyExtractor,
    SecretExtractor,
    HelmExtractor,
    DockerfileExtractor,
]

__all__ = [
    "CRDExtractor",
    "RBACExtractor",
    "ServiceExtractor",
    "DeploymentExtractor",
    "NetworkPolicyExtractor",
    "ControllerWatchExtractor",
    "DependencyExtractor",
    "SecretExtractor",
    "HelmExtractor",
    "DockerfileExtractor",
    "ALL_EXTRACTORS",
]
