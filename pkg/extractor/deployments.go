package extractor

import "strings"

// securityRelevantEnvVars lists env var names worth extracting from deployments.
var securityRelevantEnvVars = []string{
	"GOMEMLIMIT", "GOGC", "GOMAXPROCS", "GODEBUG", "GOEXPERIMENT",
}

var deploymentSearchPatterns = []string{
	"**/deployment.yaml",
	"**/deployment.yml",
	"**/deployment*.yaml",
	"**/manager*.yaml",
	"**/manager*.yml",
	"**/statefulset.yaml",
	"**/statefulset.yml",
	"charts/**/templates/deployment*.yaml",
	"charts/**/templates/deployment*.yml",
}

// extractDeployments scans YAML files for Deployment and StatefulSet definitions.
func extractDeployments(repoPath string) []Deployment {
	files := findYAMLFiles(repoPath, deploymentSearchPatterns)
	var deployments []Deployment

	for _, fpath := range files {
		for _, doc := range parseYAMLSafe(fpath) {
			kind, _ := doc["kind"].(string)
			if kind != "Deployment" && kind != "StatefulSet" {
				continue
			}
			metadata, _ := doc["metadata"].(map[string]interface{})
			if metadata == nil {
				metadata = map[string]interface{}{}
			}
			name, _ := metadata["name"].(string)
			spec, ok := doc["spec"].(map[string]interface{})
			if !ok {
				continue
			}
			template, ok := spec["template"].(map[string]interface{})
			if !ok {
				continue
			}
			podSpec, ok := template["spec"].(map[string]interface{})
			if !ok {
				continue
			}

			rawContainers := toSliceOfMaps(podSpec["containers"])
			rawInitContainers := toSliceOfMaps(podSpec["initContainers"])
			volumes := toSliceOfMaps(podSpec["volumes"])

			var containers []Container
			for _, c := range rawContainers {
				containers = append(containers, parseYAMLContainer(c, volumes, repoPath))
			}
			if containers == nil {
				containers = []Container{}
			}

			var initContainers []Container
			for _, c := range rawInitContainers {
				initContainers = append(initContainers, parseYAMLContainer(c, volumes, repoPath))
			}

			replicas := spec["replicas"]
			if replicas == nil {
				replicas = 1
			}
			serviceAccount, _ := podSpec["serviceAccountName"].(string)
			automount := podSpec["automountServiceAccountToken"]
			if automount == nil {
				automount = true
			}

			deployments = append(deployments, Deployment{
				Name:                         name,
				Kind:                         kind,
				Source:                        relativePath(repoPath, fpath),
				Replicas:                     replicas,
				ServiceAccount:               serviceAccount,
				AutomountServiceAccountToken: automount,
				Containers:                   containers,
				InitContainers:               initContainers,
			})
		}
	}

	if deployments == nil {
		deployments = []Deployment{}
	}
	return deployments
}

// parseYAMLContainer extracts a Container from a raw YAML container map.
func parseYAMLContainer(c map[string]interface{}, volumes []map[string]interface{}, repoPath string) Container {
	cName, _ := c["name"].(string)
	cImage, _ := c["image"].(string)

	cSecrets, cConfigmaps := extractEnvRefs([]map[string]interface{}{c})

	rawPorts := toSliceOfMaps(c["ports"])
	var ports []ContainerPort
	for _, p := range rawPorts {
		pName, _ := p["name"].(string)
		pPort := toInt(p["containerPort"])
		pProtocol, _ := p["protocol"].(string)
		if pProtocol == "" {
			pProtocol = "TCP"
		}
		ports = append(ports, ContainerPort{
			Name:          pName,
			ContainerPort: pPort,
			Protocol:      pProtocol,
		})
	}
	if ports == nil {
		ports = []ContainerPort{}
	}

	resources, _ := c["resources"].(map[string]interface{})
	if resources == nil {
		resources = map[string]interface{}{}
	}
	requests, _ := resources["requests"].(map[string]interface{})
	if requests == nil {
		requests = map[string]interface{}{}
	}
	limits, _ := resources["limits"].(map[string]interface{})
	if limits == nil {
		limits = map[string]interface{}{}
	}

	container := Container{
		Name:              cName,
		Image:             cImage,
		Ports:             ports,
		SecurityContext:   extractSecurityContext(c["securityContext"]),
		EnvFromSecrets:    cSecrets,
		EnvFromConfigmaps: cConfigmaps,
		VolumeMounts:      extractVolumeMounts(c, volumes),
		Resources: map[string]interface{}{
			"requests": requests,
			"limits":   limits,
		},
		EnvVars:        extractSecurityEnvVars(c),
		LivenessProbe:  extractProbe(c["livenessProbe"]),
		ReadinessProbe: extractProbe(c["readinessProbe"]),
		StartupProbe:   extractProbe(c["startupProbe"]),
	}

	container.Issues = assessContainerIssues(container)
	return container
}

// assessContainerIssues checks a container for missing resource constraints and probes.
func assessContainerIssues(c Container) []string {
	var issues []string
	res := c.Resources
	if res != nil {
		limits, _ := res["limits"].(map[string]interface{})
		if len(limits) == 0 {
			issues = append(issues, "no resource limits configured")
		} else {
			if _, ok := limits["memory"]; !ok {
				issues = append(issues, "no memory limit configured")
			}
			if _, ok := limits["cpu"]; !ok {
				issues = append(issues, "no CPU limit configured")
			}
		}
		requests, _ := res["requests"].(map[string]interface{})
		if len(requests) == 0 {
			issues = append(issues, "no resource requests configured")
		}
	}
	if c.LivenessProbe == nil {
		issues = append(issues, "no liveness probe configured")
	}
	if c.ReadinessProbe == nil {
		issues = append(issues, "no readiness probe configured")
	}
	sc := c.SecurityContext
	if len(sc) == 0 {
		issues = append(issues, "no security context configured")
	}
	return issues
}

// extractSecurityContext extracts security context fields from a container spec.
func extractSecurityContext(v interface{}) map[string]interface{} {
	sc, ok := v.(map[string]interface{})
	if !ok {
		return map[string]interface{}{}
	}
	result := make(map[string]interface{})
	for _, key := range []string{
		"allowPrivilegeEscalation",
		"readOnlyRootFilesystem",
		"runAsNonRoot",
		"runAsUser",
		"runAsGroup",
		"privileged",
	} {
		if val, exists := sc[key]; exists {
			result[key] = val
		}
	}
	if caps, ok := sc["capabilities"].(map[string]interface{}); ok {
		dropList := toStringSlice(caps["drop"])
		addList := toStringSlice(caps["add"])
		result["capabilities"] = map[string]interface{}{
			"drop": dropList,
			"add":  addList,
		}
	}
	if seccomp, ok := sc["seccompProfile"].(map[string]interface{}); ok {
		seccompType, _ := seccomp["type"].(string)
		result["seccompProfile"] = map[string]interface{}{
			"type": seccompType,
		}
	}
	return result
}

// extractEnvRefs extracts secret and configmap names from env and envFrom.
func extractEnvRefs(containers []map[string]interface{}) ([]string, []string) {
	var secrets, configmaps []string
	secretSet := make(map[string]bool)
	cmSet := make(map[string]bool)

	for _, container := range containers {
		// env[].valueFrom
		envVars := toSliceOfMaps(container["env"])
		for _, envVar := range envVars {
			valueFrom, ok := envVar["valueFrom"].(map[string]interface{})
			if !ok {
				continue
			}
			if secretRef, ok := valueFrom["secretKeyRef"].(map[string]interface{}); ok {
				if name, ok := secretRef["name"].(string); ok && !secretSet[name] {
					secrets = append(secrets, name)
					secretSet[name] = true
				}
			}
			if cmRef, ok := valueFrom["configMapKeyRef"].(map[string]interface{}); ok {
				if name, ok := cmRef["name"].(string); ok && !cmSet[name] {
					configmaps = append(configmaps, name)
					cmSet[name] = true
				}
			}
		}
		// envFrom[]
		envFroms := toSliceOfMaps(container["envFrom"])
		for _, envFrom := range envFroms {
			if secretRef, ok := envFrom["secretRef"].(map[string]interface{}); ok {
				if name, ok := secretRef["name"].(string); ok && !secretSet[name] {
					secrets = append(secrets, name)
					secretSet[name] = true
				}
			}
			if cmRef, ok := envFrom["configMapRef"].(map[string]interface{}); ok {
				if name, ok := cmRef["name"].(string); ok && !cmSet[name] {
					configmaps = append(configmaps, name)
					cmSet[name] = true
				}
			}
		}
	}

	if secrets == nil {
		secrets = []string{}
	}
	if configmaps == nil {
		configmaps = []string{}
	}
	return secrets, configmaps
}

// extractVolumeMounts extracts volume mounts with source info from volumes.
func extractVolumeMounts(container map[string]interface{}, volumes []map[string]interface{}) []map[string]interface{} {
	volMap := make(map[string]map[string]interface{})
	for _, vol := range volumes {
		volName, _ := vol["name"].(string)
		info := make(map[string]interface{})
		if secret, ok := vol["secret"].(map[string]interface{}); ok {
			secretName, _ := secret["secretName"].(string)
			info["secret"] = secretName
		}
		if cm, ok := vol["configMap"].(map[string]interface{}); ok {
			cmName, _ := cm["name"].(string)
			info["configMap"] = cmName
		}
		if _, ok := vol["projected"]; ok {
			info["projected"] = true
		}
		if _, ok := vol["emptyDir"]; ok {
			info["emptyDir"] = true
		}
		volMap[volName] = info
	}

	mounts := toSliceOfMaps(container["volumeMounts"])
	var result []map[string]interface{}
	for _, vm := range mounts {
		vmName, _ := vm["name"].(string)
		mountPath, _ := vm["mountPath"].(string)
		entry := map[string]interface{}{
			"name":      vmName,
			"mountPath": mountPath,
		}
		if extra, exists := volMap[vmName]; exists {
			for k, v := range extra {
				entry[k] = v
			}
		}
		result = append(result, entry)
	}
	if result == nil {
		result = []map[string]interface{}{}
	}
	return result
}

// toSliceOfMaps converts an interface{} to []map[string]interface{}.
func toSliceOfMaps(v interface{}) []map[string]interface{} {
	items, ok := v.([]interface{})
	if !ok {
		return nil
	}
	var result []map[string]interface{}
	for _, item := range items {
		if m, ok := item.(map[string]interface{}); ok {
			result = append(result, m)
		}
	}
	return result
}

// toInt converts an interface{} to int, handling YAML's tendency to decode
// numbers as int or float64.
func toInt(v interface{}) int {
	switch n := v.(type) {
	case int:
		return n
	case int64:
		return int(n)
	case float64:
		return int(n)
	default:
		return 0
	}
}

// extractSecurityEnvVars extracts security-relevant env vars with direct values.
func extractSecurityEnvVars(container map[string]interface{}) map[string]string {
	envVars := toSliceOfMaps(container["env"])
	result := make(map[string]string)
	for _, envVar := range envVars {
		name, _ := envVar["name"].(string)
		value, _ := envVar["value"].(string)
		if name == "" || value == "" {
			continue
		}
		for _, relevant := range securityRelevantEnvVars {
			if strings.EqualFold(name, relevant) {
				result[name] = value
				break
			}
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

// extractProbe extracts probe metadata from a container spec.
func extractProbe(v interface{}) *ProbeInfo {
	probe, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	info := &ProbeInfo{}
	if httpGet, ok := probe["httpGet"].(map[string]interface{}); ok {
		info.Type = "httpGet"
		info.Path, _ = httpGet["path"].(string)
		info.Port = httpGet["port"]
	} else if tcpSocket, ok := probe["tcpSocket"].(map[string]interface{}); ok {
		info.Type = "tcpSocket"
		info.Port = tcpSocket["port"]
	} else if _, ok := probe["exec"].(map[string]interface{}); ok {
		info.Type = "exec"
	} else if grpc, ok := probe["grpc"].(map[string]interface{}); ok {
		info.Type = "grpc"
		info.Port = grpc["port"]
	} else {
		return nil
	}
	return info
}
