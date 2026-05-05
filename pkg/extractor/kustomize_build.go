package extractor

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/kyaml/filesys"
	"sigs.k8s.io/yaml"
)

// KustomizeBuildResult holds the rendered resources from a kustomize overlay.
type KustomizeBuildResult struct {
	Overlay   string                   `json:"overlay"`
	Resources []map[string]interface{} `json:"resources,omitempty"`
}

// DefaultPreferredOverlays lists overlay directory names in priority order.
// The first match wins. Includes both generic Kubernetes operator overlays
// and common platform-specific names. Configurable via ExtractOptions.OverlayPreference.
var DefaultPreferredOverlays = []string{
	"default",
	"production",
	"prod",
	"staging",
	"odh",
	"rhoai",
	"make-deploy",
}

// kustomizeBuildOverlays discovers and renders kustomize overlays in the repo.
// Returns rendered Kubernetes resources as parsed YAML documents.
// Falls back gracefully if no overlays are found or kustomize build fails.
func kustomizeBuildOverlays(repoPath string, overlayPrefs []string) []KustomizeBuildResult {
	overlayDirs := discoverOverlays(repoPath)
	if len(overlayDirs) == 0 {
		return nil
	}

	// Pick the best overlay to build
	overlayDir := pickOverlay(overlayDirs, overlayPrefs)
	if overlayDir == "" {
		return nil
	}

	result, err := buildOverlay(repoPath, overlayDir)
	if err != nil {
		log.Printf("warning: kustomize build failed for %s: %v", overlayDir, err)
		return nil
	}

	return []KustomizeBuildResult{*result}
}

// discoverOverlays finds directories containing kustomization.yaml under config/overlays/.
func discoverOverlays(repoPath string) []string {
	overlaysRoot := filepath.Join(repoPath, "config", "overlays")
	if _, err := os.Stat(overlaysRoot); err != nil {
		return nil
	}

	var dirs []string
	entries, err := os.ReadDir(overlaysRoot)
	if err != nil {
		return nil
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		dir := filepath.Join(overlaysRoot, entry.Name())
		// Check for kustomization.yaml
		for _, name := range []string{"kustomization.yaml", "kustomization.yml", "Kustomization"} {
			if _, err := os.Stat(filepath.Join(dir, name)); err == nil {
				dirs = append(dirs, dir)
				break
			}
		}
	}

	sort.Strings(dirs)
	return dirs
}

// pickOverlay selects the best overlay directory based on the preference list.
// Falls back to DefaultPreferredOverlays if prefs is nil/empty.
func pickOverlay(overlayDirs []string, prefs []string) string {
	if len(prefs) == 0 {
		prefs = DefaultPreferredOverlays
	}
	for _, pref := range prefs {
		for _, dir := range overlayDirs {
			if filepath.Base(dir) == pref {
				return dir
			}
		}
	}
	// Fallback: use the first one
	if len(overlayDirs) > 0 {
		return overlayDirs[0]
	}
	return ""
}

// buildOverlay runs kustomize build on the given overlay directory and returns
// parsed resources.
func buildOverlay(repoPath, overlayDir string) (*KustomizeBuildResult, error) {
	fSys := filesys.MakeFsOnDisk()
	opts := krusty.MakeDefaultOptions()
	// LoadRestrictionsNone allows loading files from parent directories.
	// This is required because kustomize overlays commonly reference ../base.
	// Security note: when analyzing untrusted repositories, a malicious
	// kustomization.yaml could reference files outside the repo via path
	// traversal. This is acceptable for this tool's threat model (trusted
	// local analysis), but callers running on untrusted repos should be aware.
	opts.LoadRestrictions = 0 // LoadRestrictionsNone

	k := krusty.MakeKustomizer(opts)
	resMap, err := k.Run(fSys, overlayDir)
	if err != nil {
		return nil, fmt.Errorf("kustomize build: %w", err)
	}

	relOverlay, _ := filepath.Rel(repoPath, overlayDir)
	if relOverlay == "" {
		relOverlay = overlayDir
	}

	result := &KustomizeBuildResult{
		Overlay: relOverlay,
	}

	for _, res := range resMap.Resources() {
		yamlBytes, err := res.AsYAML()
		if err != nil {
			continue
		}
		var doc map[string]interface{}
		if err := yaml.Unmarshal(yamlBytes, &doc); err != nil {
			continue
		}
		result.Resources = append(result.Resources, doc)
	}

	return result, nil
}

// mergeKustomizeResources takes rendered kustomize resources and merges them
// into the architecture extraction. For each resource type we support
// (Service, Deployment, NetworkPolicy, Webhook configs), rendered versions
// replace or supplement raw-scanned versions.
//
// After merging, raw-scanned resources that are clearly patches (multiple entries
// with the same name from different files) are deduplicated: if a kustomize-rendered
// version exists, all raw versions with matching names are removed.
func mergeKustomizeResources(arch *ComponentArchitecture, results []KustomizeBuildResult) {
	if len(results) == 0 {
		return
	}

	var renderedServiceNames []string
	var renderedDeploymentNames []string
	var renderedNPNames []string
	var renderedWebhookNames []string

	for _, result := range results {
		source := fmt.Sprintf("kustomize:%s", result.Overlay)

		for _, doc := range result.Resources {
			kind, _ := doc["kind"].(string)
			switch kind {
			case "Service":
				if svc := parseRenderedService(doc, source); svc != nil {
					arch.Services = mergeService(arch.Services, *svc)
					renderedServiceNames = append(renderedServiceNames, svc.Name)
				}
			case "Deployment", "StatefulSet":
				if dep := parseRenderedDeployment(doc, source); dep != nil {
					arch.Deployments = mergeDeployment(arch.Deployments, *dep)
					renderedDeploymentNames = append(renderedDeploymentNames, dep.Name)
				}
			case "NetworkPolicy":
				if np := parseRenderedNetworkPolicy(doc, source); np != nil {
					arch.NetworkPolicies = mergeNetworkPolicy(arch.NetworkPolicies, *np)
					renderedNPNames = append(renderedNPNames, np.Name)
				}
			case "ValidatingWebhookConfiguration":
				mergeRenderedWebhooks(arch, doc, "validating", source)
				renderedWebhookNames = append(renderedWebhookNames, collectWebhookNames(doc)...)
			case "MutatingWebhookConfiguration":
				mergeRenderedWebhooks(arch, doc, "mutating", source)
				renderedWebhookNames = append(renderedWebhookNames, collectWebhookNames(doc)...)
			case "PodDisruptionBudget":
				if pdb := parseRenderedPDB(doc, source); pdb != nil {
					arch.PodDisruptionBudgets = mergePDB(arch.PodDisruptionBudgets, *pdb)
				}
			case "HorizontalPodAutoscaler":
				if hpa := parseRenderedHPA(doc, source); hpa != nil {
					arch.HorizontalPodAutoscalers = mergeHPA(arch.HorizontalPodAutoscalers, *hpa)
				}
			}
		}
	}

	// Remove raw-scanned duplicates that are superseded by kustomize-rendered versions
	arch.Services = deduplicateByRendered(arch.Services, renderedServiceNames,
		func(s Service) string { return s.Name },
		func(s Service) bool { return isKustomizeSource(s.Source) })
	arch.Deployments = deduplicateByRendered(arch.Deployments, renderedDeploymentNames,
		func(d Deployment) string { return d.Name },
		func(d Deployment) bool { return isKustomizeSource(d.Source) })
	arch.NetworkPolicies = deduplicateByRendered(arch.NetworkPolicies, renderedNPNames,
		func(np NetworkPolicy) string { return np.Name },
		func(np NetworkPolicy) bool { return isKustomizeSource(np.Source) })

	// Deduplicate webhooks: kustomize-rendered versions (with full rules) supersede
	// Go-source-scanned versions (which lack rules/operations).
	if len(renderedWebhookNames) > 0 {
		arch.Webhooks = deduplicateByRendered(arch.Webhooks, renderedWebhookNames,
			func(wh WebhookConfig) string { return wh.Name },
			func(wh WebhookConfig) bool { return isKustomizeSource(wh.Source) })
	}
}

// isKustomizeSource checks if a source string indicates a kustomize-rendered resource.
func isKustomizeSource(source string) bool {
	return len(source) >= 10 && source[:10] == "kustomize:"
}

// deduplicateByRendered removes raw-scanned items whose names match (via namesMatch)
// any kustomize-rendered name. Keeps all kustomize-rendered items and any raw items
// that have no rendered equivalent.
func deduplicateByRendered[T any](items []T, renderedNames []string, getName func(T) string, isRendered func(T) bool) []T {
	if len(renderedNames) == 0 {
		return items
	}
	var result []T
	for _, item := range items {
		if isRendered(item) {
			result = append(result, item)
			continue
		}
		// Check if this raw item is superseded by a rendered one
		name := getName(item)
		superseded := false
		for _, rn := range renderedNames {
			if namesMatch(name, rn) {
				superseded = true
				break
			}
		}
		if !superseded {
			result = append(result, item)
		}
	}
	return result
}

// parseRenderedService converts a rendered Kubernetes Service document into our Service type.
func parseRenderedService(doc map[string]interface{}, source string) *Service {
	metadata, _ := doc["metadata"].(map[string]interface{})
	if metadata == nil {
		return nil
	}
	name, _ := metadata["name"].(string)
	if name == "" {
		return nil
	}

	spec, _ := doc["spec"].(map[string]interface{})
	if spec == nil {
		spec = map[string]interface{}{}
	}

	rawPorts, _ := spec["ports"].([]interface{})
	var ports []ServicePort
	for _, p := range rawPorts {
		pm, ok := p.(map[string]interface{})
		if !ok {
			continue
		}
		portName, _ := pm["name"].(string)
		protocol, _ := pm["protocol"].(string)
		if protocol == "" {
			protocol = "TCP"
		}
		portVal := pm["port"]
		if portVal == nil {
			portVal = 0
		}
		targetPortVal := pm["targetPort"]
		if targetPortVal == nil {
			targetPortVal = 0
		}
		ports = append(ports, ServicePort{
			Name:       portName,
			Port:       portVal,
			TargetPort: targetPortVal,
			Protocol:   protocol,
		})
	}
	if ports == nil {
		ports = []ServicePort{}
	}

	selector, _ := spec["selector"].(map[string]interface{})
	if selector == nil {
		selector = map[string]interface{}{}
	}

	svcType, _ := spec["type"].(string)
	if svcType == "" {
		svcType = "ClusterIP"
	}

	return &Service{
		Name:     name,
		Source:   source,
		Type:     svcType,
		Ports:    ports,
		Selector: selector,
	}
}

// parseRenderedDeployment converts a rendered Deployment/StatefulSet into our Deployment type.
func parseRenderedDeployment(doc map[string]interface{}, source string) *Deployment {
	metadata, _ := doc["metadata"].(map[string]interface{})
	if metadata == nil {
		return nil
	}
	name, _ := metadata["name"].(string)
	kind, _ := doc["kind"].(string)
	if name == "" {
		return nil
	}

	spec, _ := doc["spec"].(map[string]interface{})
	if spec == nil {
		return nil
	}

	dep := &Deployment{
		Name:   name,
		Kind:   kind,
		Source: source,
	}

	if replicas, ok := spec["replicas"]; ok {
		dep.Replicas = replicas
	}

	// Extract pod template
	template, _ := spec["template"].(map[string]interface{})
	if template == nil {
		return dep
	}
	podSpec, _ := template["spec"].(map[string]interface{})
	if podSpec == nil {
		return dep
	}

	if sa, ok := podSpec["serviceAccountName"].(string); ok {
		dep.ServiceAccount = sa
	}
	if autoMount, ok := podSpec["automountServiceAccountToken"]; ok {
		dep.AutomountServiceAccountToken = autoMount
	}

	volumes := toSliceOfMaps(podSpec["volumes"])

	rawContainers, _ := podSpec["containers"].([]interface{})
	for _, c := range rawContainers {
		cm, ok := c.(map[string]interface{})
		if !ok {
			continue
		}
		container := parseContainerSpec(cm, volumes)
		dep.Containers = append(dep.Containers, container)
	}

	rawInitContainers, _ := podSpec["initContainers"].([]interface{})
	for _, c := range rawInitContainers {
		cm, ok := c.(map[string]interface{})
		if !ok {
			continue
		}
		container := parseContainerSpec(cm, volumes)
		dep.InitContainers = append(dep.InitContainers, container)
	}

	return dep
}

// parseContainerSpec extracts container details from a container map.
func parseContainerSpec(cm map[string]interface{}, volumes []map[string]interface{}) Container {
	container := Container{}
	container.Name, _ = cm["name"].(string)
	container.Image, _ = cm["image"].(string)

	// Ports
	if rawPorts, ok := cm["ports"].([]interface{}); ok {
		for _, p := range rawPorts {
			if pm, ok := p.(map[string]interface{}); ok {
				cp := ContainerPort{}
				cp.Name, _ = pm["name"].(string)
				if v, ok := pm["containerPort"].(float64); ok {
					cp.ContainerPort = int(v)
				} else if v, ok := pm["containerPort"].(int64); ok {
					cp.ContainerPort = int(v)
				}
				cp.Protocol, _ = pm["protocol"].(string)
				container.Ports = append(container.Ports, cp)
			}
		}
	}

	// Security context
	if sc, ok := cm["securityContext"].(map[string]interface{}); ok {
		container.SecurityContext = sc
	}

	// Resources (normalize to always have requests/limits keys)
	if res, ok := cm["resources"].(map[string]interface{}); ok {
		requests, _ := res["requests"].(map[string]interface{})
		if requests == nil {
			requests = map[string]interface{}{}
		}
		limits, _ := res["limits"].(map[string]interface{})
		if limits == nil {
			limits = map[string]interface{}{}
		}
		container.Resources = map[string]interface{}{
			"requests": requests,
			"limits":   limits,
		}
	} else {
		container.Resources = map[string]interface{}{
			"requests": map[string]interface{}{},
			"limits":   map[string]interface{}{},
		}
	}

	// Environment variables
	if envList, ok := cm["env"].([]interface{}); ok {
		envVars := make(map[string]string)
		for _, e := range envList {
			if em, ok := e.(map[string]interface{}); ok {
				name, _ := em["name"].(string)
				value, _ := em["value"].(string)
				if name != "" {
					if value != "" {
						envVars[name] = value
					}
					// Check secretKeyRef
					if vf, ok := em["valueFrom"].(map[string]interface{}); ok {
						if skr, ok := vf["secretKeyRef"].(map[string]interface{}); ok {
							secretName, _ := skr["name"].(string)
							container.EnvFromSecrets = appendUnique(container.EnvFromSecrets, secretName)
						}
						if cmr, ok := vf["configMapKeyRef"].(map[string]interface{}); ok {
							cmName, _ := cmr["name"].(string)
							container.EnvFromConfigmaps = appendUnique(container.EnvFromConfigmaps, cmName)
						}
					}
				}
			}
		}
		if len(envVars) > 0 {
			container.EnvVars = envVars
		}
	}

	// Volume mounts
	container.VolumeMounts = extractVolumeMounts(cm, volumes)

	// Probes
	container.LivenessProbe = parseProbe(cm, "livenessProbe")
	container.ReadinessProbe = parseProbe(cm, "readinessProbe")
	container.StartupProbe = parseProbe(cm, "startupProbe")

	container.Issues = assessContainerIssues(container)
	return container
}

// parseProbe extracts probe info from a container map.
func parseProbe(cm map[string]interface{}, key string) *ProbeInfo {
	probe, ok := cm[key].(map[string]interface{})
	if !ok {
		return nil
	}
	info := &ProbeInfo{}
	if hp, ok := probe["httpGet"].(map[string]interface{}); ok {
		info.Type = "httpGet"
		info.Path, _ = hp["path"].(string)
		info.Port = hp["port"]
	} else if tp, ok := probe["tcpSocket"].(map[string]interface{}); ok {
		info.Type = "tcpSocket"
		info.Port = tp["port"]
	} else if _, ok := probe["exec"].(map[string]interface{}); ok {
		info.Type = "exec"
	} else if gp, ok := probe["grpc"].(map[string]interface{}); ok {
		info.Type = "grpc"
		info.Port = gp["port"]
	}
	return info
}

// parseRenderedNetworkPolicy converts a rendered NetworkPolicy into our type.
func parseRenderedNetworkPolicy(doc map[string]interface{}, source string) *NetworkPolicy {
	np := parseNetworkPolicyDoc(doc, "", "")
	if np == nil {
		return nil
	}
	np.Source = source
	return np
}

// mergeRenderedWebhooks extracts webhook configurations from rendered webhook docs.
func mergeRenderedWebhooks(arch *ComponentArchitecture, doc map[string]interface{}, whType, source string) {
	metadata, _ := doc["metadata"].(map[string]interface{})
	if metadata == nil {
		return
	}
	configName, _ := metadata["name"].(string)

	webhooks, _ := doc["webhooks"].([]interface{})
	for _, wh := range webhooks {
		whMap, ok := wh.(map[string]interface{})
		if !ok {
			continue
		}

		name, _ := whMap["name"].(string)
		failurePolicy, _ := whMap["failurePolicy"].(string)

		var serviceRef, path string
		if cs, ok := whMap["clientConfig"].(map[string]interface{}); ok {
			if svc, ok := cs["service"].(map[string]interface{}); ok {
				svcName, _ := svc["name"].(string)
				svcNs, _ := svc["namespace"].(string)
				if svcNs != "" {
					serviceRef = svcNs + "/" + svcName
				} else {
					serviceRef = svcName
				}
				path, _ = svc["path"].(string)
			}
		}

		var rules []WebhookRule
		if rawRules, ok := whMap["rules"].([]interface{}); ok {
			for _, r := range rawRules {
				rm, ok := r.(map[string]interface{})
				if !ok {
					continue
				}
				rules = append(rules, WebhookRule{
					APIGroups:   toStringSlice(rm["apiGroups"]),
					APIVersions: toStringSlice(rm["apiVersions"]),
					Resources:   toStringSlice(rm["resources"]),
					Operations:  toStringSlice(rm["operations"]),
				})
			}
		}

		wc := WebhookConfig{
			Name:          name,
			Type:          whType,
			ServiceRef:    serviceRef,
			Path:          path,
			FailurePolicy: failurePolicy,
			Rules:         rules,
			Source:        fmt.Sprintf("%s (%s)", source, configName),
		}

		arch.Webhooks = mergeWebhook(arch.Webhooks, wc)
	}
}

// mergeService adds or replaces a service by name. Kustomize-rendered versions
// take precedence over raw-scanned ones. Also handles kustomize namePrefix by
// checking if the rendered name contains the raw name as a suffix.
func mergeService(existing []Service, rendered Service) []Service {
	for i, s := range existing {
		if namesMatch(s.Name, rendered.Name) {
			existing[i] = rendered
			return existing
		}
	}
	return append(existing, rendered)
}

// mergeDeployment adds or replaces a deployment by name.
func mergeDeployment(existing []Deployment, rendered Deployment) []Deployment {
	for i, d := range existing {
		if namesMatch(d.Name, rendered.Name) {
			existing[i] = rendered
			return existing
		}
	}
	return append(existing, rendered)
}

// mergeNetworkPolicy adds or replaces a network policy by name.
func mergeNetworkPolicy(existing []NetworkPolicy, rendered NetworkPolicy) []NetworkPolicy {
	for i, np := range existing {
		if namesMatch(np.Name, rendered.Name) {
			existing[i] = rendered
			return existing
		}
	}
	return append(existing, rendered)
}

// mergeWebhook adds or replaces a webhook config by name.
func mergeWebhook(existing []WebhookConfig, rendered WebhookConfig) []WebhookConfig {
	for i, wh := range existing {
		if wh.Name == rendered.Name {
			existing[i] = rendered
			return existing
		}
	}
	return append(existing, rendered)
}

// namesMatch checks if two resource names refer to the same resource,
// accounting for kustomize namePrefix and nameSuffix transforms.
// e.g. "webhook-service" matches "model-registry-operator-webhook-service" (prefix)
// e.g. "controller-manager" matches "controller-manager-metrics" (suffix)
func namesMatch(rawName, kustomizeName string) bool {
	if rawName == kustomizeName {
		return true
	}
	// Check if kustomize name is a prefixed version of the raw name
	if hasSeparatedAffix(kustomizeName, rawName) {
		return true
	}
	// Check if kustomize name is a suffixed version of the raw name
	if hasSeparatedAffix(rawName, kustomizeName) {
		return true
	}
	return false
}

// hasSeparatedAffix checks if longer contains shorter as a prefix or suffix
// separated by '-'.
func hasSeparatedAffix(longer, shorter string) bool {
	if len(longer) <= len(shorter) {
		return false
	}
	// Check suffix: longer ends with shorter, preceded by '-'
	if longer[len(longer)-len(shorter):] == shorter {
		prefix := longer[:len(longer)-len(shorter)]
		if len(prefix) > 0 && prefix[len(prefix)-1] == '-' {
			return true
		}
	}
	return false
}

// collectWebhookNames extracts individual webhook names from a webhook configuration document.
func collectWebhookNames(doc map[string]interface{}) []string {
	webhooks, _ := doc["webhooks"].([]interface{})
	var names []string
	for _, wh := range webhooks {
		if whMap, ok := wh.(map[string]interface{}); ok {
			if name, ok := whMap["name"].(string); ok {
				names = append(names, name)
			}
		}
	}
	return names
}

// appendUnique appends a string to a slice only if not already present.
func appendUnique(slice []string, s string) []string {
	if s == "" {
		return slice
	}
	for _, existing := range slice {
		if existing == s {
			return slice
		}
	}
	return append(slice, s)
}
