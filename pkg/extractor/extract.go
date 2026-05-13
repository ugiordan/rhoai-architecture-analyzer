package extractor

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// AnalyzerVersion is the version of the analyzer, set by the CLI entry point.
// Defaults to "dev" if not set.
var AnalyzerVersion = "dev"

// ExtractAll runs all extractors on the given repository path and returns the
// combined ComponentArchitecture. Pass nil for opts to use defaults.
func ExtractAll(repoPath string, opts *ExtractOptions) (*ComponentArchitecture, error) {
	absPath, err := filepath.Abs(repoPath)
	if err != nil {
		return nil, fmt.Errorf("invalid repo path: %w", err)
	}
	info, err := os.Stat(absPath)
	if err != nil || !info.IsDir() {
		return nil, fmt.Errorf("repository path does not exist or is not a directory: %s", absPath)
	}

	if opts == nil {
		opts = &ExtractOptions{}
	}

	componentName := filepath.Base(absPath)

	org := opts.Org
	if org == "" {
		org = detectOrg(absPath)
	}

	modulePrefixes := opts.ModulePrefixes
	if len(modulePrefixes) == 0 {
		modulePrefixes = DefaultModulePrefixes()
	}

	// Load Go packages for AST-based extraction (requires go.mod).
	// Returns nil for non-Go repos, "full" when type info is available,
	// "syntax-only" when >50% of packages have type-checking errors.
	goPackages := loadGoPackages(repoPath)

	arch := &ComponentArchitecture{
		Component:       componentName,
		Aliases:         opts.Aliases,
		Repo:            fmt.Sprintf("%s/%s", org, componentName),
		CommitSHA:       detectHEAD(absPath),
		ExtractedAt:     time.Now().UTC().Format(time.RFC3339),
		AnalyzerVersion: AnalyzerVersion,
		SchemaVersion:   "2",
		CRDs:            extractCRDs(absPath),
		RBAC:            extractRBAC(absPath),
		Services:        extractServices(absPath),
		Deployments:     extractDeployments(absPath),
		NetworkPolicies: extractNetworkPolicies(absPath),
		ControllerWatch: extractControllerWatches(absPath),
		Dependencies:    extractDependencies(absPath, modulePrefixes),
		Secrets:         extractSecrets(absPath),
		Dockerfiles:     extractDockerfiles(absPath),
		Helm:            extractHelm(absPath),
		Webhooks:        extractWebhooks(absPath),
		ConfigMaps:      extractConfigMaps(absPath),
		HTTPEndpoints:   extractHTTPEndpoints(absPath),
		IngressRouting:      extractIngress(absPath),
		ExternalConnections: extractExternalConnections(absPath),
		FeatureGates:        extractFeatureGates(absPath),
		RuntimeDependencies: extractRuntimeDependencies(absPath),
	}

	// Go AST enrichment: set mode/warning and merge AST-extracted data
	// with YAML-extracted data. Only "full" mode (type info available) is
	// used for CRD, webhook, and resource-op extraction.
	if goPackages != nil {
		arch.GoASTMode = goPackages.Mode
		arch.GoASTWarning = goPackages.Warning
	}
	if goPackages != nil && goPackages.Mode == "full" {
		goCRDs := extractCRDsFromGo(goPackages)
		arch.CRDs = mergeCRDs(arch.CRDs, goCRDs)
	}

	// Python source port detection: scan .py files for listening ports
	if pythonPorts := extractPythonPorts(absPath); len(pythonPorts) > 0 {
		arch.Services = append(arch.Services, pythonPorts...)
	}

	// Cache analysis runs after watches and deployments are extracted
	arch.CacheConfig = extractCacheConfig(absPath, arch.ControllerWatch, arch.Deployments)

	// Kustomize build: render overlays and merge rendered resources into extraction.
	// Rendered resources replace or supplement raw-scanned ones, giving us fully
	// resolved manifests with patches and substitutions applied.
	kustomizeResults := kustomizeBuildOverlays(absPath, opts.OverlayPreference)
	mergeKustomizeResources(arch, kustomizeResults)

	// Extract webhook server port from Go source (controller-runtime webhook.Options).
	// Runs AFTER kustomize merge so rendered webhooks also get the port.
	// If not explicitly configured, controller-runtime defaults to 9443.
	webhookPort := extractWebhookServerPort(absPath)
	if webhookPort == 0 && len(arch.Webhooks) > 0 {
		webhookPort = 9443 // controller-runtime default
	}
	if webhookPort > 0 {
		for i := range arch.Webhooks {
			if arch.Webhooks[i].Port == 0 {
				arch.Webhooks[i].Port = webhookPort
			}
		}
	}

	// Go source enrichment: handler mapping, data_read, enable_condition.
	// Runs after kustomize merge and port assignment.
	enrichWebhooks(arch.Webhooks, absPath)

	// Go AST webhook behavior: extract field-level mutations and validations
	// from Default/Validate* methods on kubebuilder webhook types.
	if goPackages != nil && goPackages.Mode == "full" {
		behaviors := extractWebhookBehavior(goPackages)
		mergeWebhookBehavior(arch, behaviors)
	}

	// Go AST resource operations: extract programmatic Create/Update/Patch/Delete
	// calls from Reconcile methods with resolved type information.
	if goPackages != nil && goPackages.Mode == "full" {
		resourceOps := extractResourceOps(goPackages)
		if len(resourceOps) > 0 {
			arch.ControllerWatch = append(arch.ControllerWatch, ControllerWatch{
				Type:        "resource_ops",
				GVK:         "programmatic",
				Source:      "go_ast",
				ResourceOps: resourceOps,
			})
		}
	}

	// Kustomize component discovery (for operator repos with *_support.go files)
	arch.KustomizeComponents = extractKustomizeComponents(absPath)

	// Serving runtime discovery (KServe/ModelMesh)
	arch.ServingRuntimes = extractServingRuntimes(absPath)

	// Serving runtime image-to-CRD mapping: scan all YAML and Go source for
	// ServingRuntime/ClusterServingRuntime/InferenceService container images
	arch.ServingRuntimeRefs = extractServingRuntimeRefs(absPath)

	// Resource defaults from configmaps (inference config, deployment defaults)
	arch.ResourceDefaults = extractResourceDefaults(absPath)

	// Availability: PDB and HPA extraction
	arch.PodDisruptionBudgets = extractPDBs(absPath)
	arch.HorizontalPodAutoscalers = extractHPAs(absPath)

	// API types: parse *_types.go files for CR struct definitions
	arch.APITypes = extractAPITypes(absPath)

	// Status conditions: extract condition type/reason constants.
	// Also returns Go constant names for dedup with operator config.
	var statusConditionConstNames map[string]bool
	arch.StatusConditions, statusConditionConstNames = extractStatusConditions(absPath)

	// Operator config: extract const/var blocks (dedup with status conditions)
	arch.OperatorConfig = extractOperatorConfig(absPath, statusConditionConstNames)

	// Reconcile sequences: extract ordered sub-resource reconciliation steps
	arch.ReconcileSequences = extractReconcileSequences(absPath)

	// Prometheus metrics: extract metric registrations
	arch.PrometheusMetrics = extractPrometheusMetrics(absPath)

	// Platform detection: extract capability checks and conditional resource creation
	arch.PlatformDetection = extractPlatformDetection(absPath)

	// Template file enumeration: list .tmpl files for operators that use
	// Go templates to define runtime-rendered Kubernetes resources.
	arch.TemplateFiles = findTemplateFiles(absPath)

	// Label/annotation contracts: detect well-known labels that imply
	// cross-component integration (Kueue, KServe, Istio, etc.)
	arch.LabelContracts = extractLabelContracts(absPath)

	// Python k8s API calls: detect kubernetes client usage in Python source
	// (e.g., codeflare-sdk creating LocalQueue objects via CustomObjectsApi)
	arch.PythonK8sCalls = extractPythonK8sCalls(absPath)

	// Kustomize overlay cross-references: parse all kustomization.yaml files
	// to extract resources, patches, generators, and image transforms
	arch.KustomizeOverlayRefs = extractKustomizeOverlayRefs(absPath)

	// Component cross-references: detect provider/adapter directories referencing
	// other known components (e.g., llama-stack's providers/remote/inference/vllm/)
	arch.ComponentRefs = extractComponentRefs(absPath, componentName, opts.KnownComponents)

	// Cross-reference pass: link services to deployments, detect runtime deps
	buildCrossReferences(arch)

	// ConfigMap volume mount correlation: explicit ConfigMap → container links
	arch.ConfigMapVolumes = extractConfigMapVolumes(arch)

	// Availability assessment: flag deployments missing PDB/HPA
	assessAvailability(arch)

	// Data coverage: assess richness of each section for LLM context
	arch.DataCoverage = computeDataCoverage(arch)

	// Generate natural-language summary
	arch.Summary = generateSummary(arch)

	// Normalize output ordering for deterministic JSON
	SortOutput(arch)

	return arch, nil
}

// detectOrg tries to determine the GitHub organization from the repo's go.mod
// module path, then from .git/config remote origin, then falls back to
// "opendatahub-io".
//
// Note: The detected org name is embedded in output artifacts (ComponentArchitecture.Repo).
// When analyzing internal/private forks, use ExtractOptions.Org to override
// auto-detection and avoid disclosing internal organization names.
func detectOrg(repoPath string) string {
	// Try go.mod first
	goModPath := filepath.Join(repoPath, "go.mod")
	if f, err := os.Open(goModPath); err == nil {
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(line, "module ") {
				module := strings.TrimPrefix(line, "module ")
				module = strings.TrimSpace(module)
				// Parse github.com/org/repo format
				parts := strings.Split(module, "/")
				if len(parts) >= 2 && parts[0] == "github.com" {
					return parts[1]
				}
			}
		}
	}

	// Try .git/config remote origin
	gitConfigPath := filepath.Join(repoPath, ".git", "config")
	if f, err := os.Open(gitConfigPath); err == nil {
		defer f.Close()
		scanner := bufio.NewScanner(f)
		inOrigin := false
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			// Track which [remote] section we're in
			if strings.HasPrefix(line, "[") {
				inOrigin = line == `[remote "origin"]`
				continue
			}
			if !inOrigin {
				continue
			}
			if !strings.HasPrefix(line, "url = ") {
				continue
			}
			url := strings.TrimPrefix(line, "url = ")
			url = strings.TrimSpace(url)
			if org := orgFromGitURL(url); org != "" {
				return org
			}
		}
	}

	return "opendatahub-io"
}

// detectHEAD reads the current commit SHA from .git/HEAD.
// Returns empty string if the repo has no git directory or HEAD cannot be resolved.
func detectHEAD(repoPath string) string {
	headPath := filepath.Join(repoPath, ".git", "HEAD")
	data, err := os.ReadFile(headPath)
	if err != nil {
		return ""
	}
	content := strings.TrimSpace(string(data))

	// Detached HEAD: raw SHA
	if len(content) == 40 && !strings.Contains(content, " ") {
		return content
	}

	// Symbolic ref: "ref: refs/heads/main"
	if strings.HasPrefix(content, "ref: ") {
		refPath := strings.TrimPrefix(content, "ref: ")
		shaPath := filepath.Join(repoPath, ".git", refPath)
		shaData, err := os.ReadFile(shaPath)
		if err != nil {
			// Try packed-refs
			return readPackedRef(repoPath, refPath)
		}
		return strings.TrimSpace(string(shaData))
	}

	return ""
}

// readPackedRef looks up a ref in .git/packed-refs (used when refs are packed).
func readPackedRef(repoPath, refPath string) string {
	packedPath := filepath.Join(repoPath, ".git", "packed-refs")
	f, err := os.Open(packedPath)
	if err != nil {
		return ""
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) == 2 && parts[1] == refPath {
			return parts[0]
		}
	}
	return ""
}

// orgFromGitURL extracts the GitHub organization from a git remote URL.
// Supports HTTPS (https://github.com/org/repo.git) and SSH (git@github.com:org/repo.git).
func orgFromGitURL(url string) string {
	if !strings.Contains(url, "github.com") {
		return ""
	}
	url = strings.TrimSuffix(url, ".git")

	// SSH format: git@github.com:org/repo
	if strings.Contains(url, ":") && !strings.Contains(url, "://") {
		colonParts := strings.SplitN(url, ":", 2)
		if len(colonParts) == 2 {
			orgRepo := strings.SplitN(colonParts[1], "/", 2)
			if len(orgRepo) >= 1 && orgRepo[0] != "" {
				return orgRepo[0]
			}
		}
		return ""
	}

	// HTTPS format: https://github.com/org/repo
	parts := strings.Split(url, "/")
	// Find "github.com" in parts and return the next segment
	for i, part := range parts {
		if part == "github.com" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}
