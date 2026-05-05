package extractor

import (
	"fmt"
	"strings"
)

var pdbSearchPatterns = []string{
	"**/poddisruptionbudget*.yaml",
	"**/poddisruptionbudget*.yml",
	"**/pdb*.yaml",
	"**/pdb*.yml",
	"config/**/poddisruptionbudget*.yaml",
	"config/**/pdb*.yaml",
	"charts/**/templates/pdb*.yaml",
	"charts/**/templates/poddisruptionbudget*.yaml",
}

var hpaSearchPatterns = []string{
	"**/horizontalpodautoscaler*.yaml",
	"**/horizontalpodautoscaler*.yml",
	"**/hpa*.yaml",
	"**/hpa*.yml",
	"config/**/hpa*.yaml",
	"config/**/horizontalpodautoscaler*.yaml",
	"charts/**/templates/hpa*.yaml",
	"charts/**/templates/horizontalpodautoscaler*.yaml",
}

// extractPDBs scans YAML files for PodDisruptionBudget definitions.
func extractPDBs(repoPath string) []PodDisruptionBudget {
	files := findYAMLFiles(repoPath, pdbSearchPatterns)
	var pdbs []PodDisruptionBudget

	for _, fpath := range files {
		for _, doc := range parseYAMLSafe(fpath) {
			kind, _ := doc["kind"].(string)
			if kind != "PodDisruptionBudget" {
				continue
			}
			metadata, _ := doc["metadata"].(map[string]interface{})
			if metadata == nil {
				continue
			}
			name, _ := metadata["name"].(string)
			spec, _ := doc["spec"].(map[string]interface{})
			if spec == nil {
				continue
			}

			pdb := PodDisruptionBudget{
				Name:   name,
				Source: relativePath(repoPath, fpath),
			}

			if v, ok := spec["minAvailable"]; ok {
				pdb.MinAvailable = v
			}
			if v, ok := spec["maxUnavailable"]; ok {
				pdb.MaxUnavailable = v
			}
			if sel, ok := spec["selector"].(map[string]interface{}); ok {
				if ml, ok := sel["matchLabels"].(map[string]interface{}); ok {
					pdb.Selector = ml
				}
			}

			pdbs = append(pdbs, pdb)
		}
	}
	return pdbs
}

// extractHPAs scans YAML files for HorizontalPodAutoscaler definitions.
func extractHPAs(repoPath string) []HorizontalPodAutoscaler {
	files := findYAMLFiles(repoPath, hpaSearchPatterns)
	var hpas []HorizontalPodAutoscaler

	for _, fpath := range files {
		for _, doc := range parseYAMLSafe(fpath) {
			kind, _ := doc["kind"].(string)
			if kind != "HorizontalPodAutoscaler" {
				continue
			}
			metadata, _ := doc["metadata"].(map[string]interface{})
			if metadata == nil {
				continue
			}
			name, _ := metadata["name"].(string)
			spec, _ := doc["spec"].(map[string]interface{})
			if spec == nil {
				continue
			}

			hpa := HorizontalPodAutoscaler{
				Name:   name,
				Source: relativePath(repoPath, fpath),
			}

			if ref, ok := spec["scaleTargetRef"].(map[string]interface{}); ok {
				refKind, _ := ref["kind"].(string)
				refName, _ := ref["name"].(string)
				hpa.TargetRef = fmt.Sprintf("%s/%s", refKind, refName)
			}
			hpa.MinReplicas = spec["minReplicas"]
			hpa.MaxReplicas = spec["maxReplicas"]

			if metrics, ok := spec["metrics"].([]interface{}); ok {
				for _, m := range metrics {
					mm, ok := m.(map[string]interface{})
					if !ok {
						continue
					}
					mType, _ := mm["type"].(string)
					if res, ok := mm["resource"].(map[string]interface{}); ok {
						resName, _ := res["name"].(string)
						hpa.Metrics = append(hpa.Metrics, fmt.Sprintf("%s/%s", mType, resName))
					} else {
						hpa.Metrics = append(hpa.Metrics, mType)
					}
				}
			}

			hpas = append(hpas, hpa)
		}
	}
	return hpas
}

// assessAvailability cross-references deployments with PDBs and HPAs,
// adding issues to deployments that lack availability protections.
func assessAvailability(arch *ComponentArchitecture) {
	for i := range arch.Deployments {
		dep := &arch.Deployments[i]
		hasPDB := false
		for _, pdb := range arch.PodDisruptionBudgets {
			if pdbMatchesDeployment(pdb, dep.Name) {
				hasPDB = true
				break
			}
		}
		if !hasPDB {
			dep.Issues = append(dep.Issues, "no PodDisruptionBudget configured")
		}

		hasHPA := false
		for _, hpa := range arch.HorizontalPodAutoscalers {
			if hpaTargetsDeployment(hpa, dep.Name, dep.Kind) {
				hasHPA = true
				break
			}
		}
		if !hasHPA {
			dep.Issues = append(dep.Issues, "no HorizontalPodAutoscaler configured")
		}
	}
}

// pdbMatchesDeployment checks if a PDB's selector labels match a deployment name.
// Platform-agnostic: checks all selector label values against the deployment name.
func pdbMatchesDeployment(pdb PodDisruptionBudget, deploymentName string) bool {
	if pdb.Selector == nil {
		return false
	}
	for _, v := range pdb.Selector {
		val, ok := v.(string)
		if !ok || val == "" || val == "true" || val == "false" {
			continue
		}
		if val == deploymentName || namesMatch(val, deploymentName) || namesMatch(deploymentName, val) {
			return true
		}
		if len(val) >= 3 && strings.HasPrefix(deploymentName, val+"-") && len(val)*100/len(deploymentName) >= 40 {
			return true
		}
	}
	return false
}

// hpaTargetsDeployment checks if an HPA targets a specific deployment.
func hpaTargetsDeployment(hpa HorizontalPodAutoscaler, deploymentName, deploymentKind string) bool {
	expected := fmt.Sprintf("%s/%s", deploymentKind, deploymentName)
	return hpa.TargetRef == expected || namesMatch(deploymentName, hpa.TargetRef)
}

// parseRenderedPDB converts a kustomize-rendered PDB document.
func parseRenderedPDB(doc map[string]interface{}, source string) *PodDisruptionBudget {
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
		return nil
	}
	pdb := &PodDisruptionBudget{
		Name:   name,
		Source: source,
	}
	if v, ok := spec["minAvailable"]; ok {
		pdb.MinAvailable = v
	}
	if v, ok := spec["maxUnavailable"]; ok {
		pdb.MaxUnavailable = v
	}
	if sel, ok := spec["selector"].(map[string]interface{}); ok {
		if ml, ok := sel["matchLabels"].(map[string]interface{}); ok {
			pdb.Selector = ml
		}
	}
	return pdb
}

// parseRenderedHPA converts a kustomize-rendered HPA document.
func parseRenderedHPA(doc map[string]interface{}, source string) *HorizontalPodAutoscaler {
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
		return nil
	}
	hpa := &HorizontalPodAutoscaler{
		Name:   name,
		Source: source,
	}
	if ref, ok := spec["scaleTargetRef"].(map[string]interface{}); ok {
		refKind, _ := ref["kind"].(string)
		refName, _ := ref["name"].(string)
		hpa.TargetRef = fmt.Sprintf("%s/%s", refKind, refName)
	}
	hpa.MinReplicas = spec["minReplicas"]
	hpa.MaxReplicas = spec["maxReplicas"]
	if metrics, ok := spec["metrics"].([]interface{}); ok {
		for _, m := range metrics {
			mm, ok := m.(map[string]interface{})
			if !ok {
				continue
			}
			mType, _ := mm["type"].(string)
			if res, ok := mm["resource"].(map[string]interface{}); ok {
				resName, _ := res["name"].(string)
				hpa.Metrics = append(hpa.Metrics, fmt.Sprintf("%s/%s", mType, resName))
			} else {
				hpa.Metrics = append(hpa.Metrics, mType)
			}
		}
	}
	return hpa
}

// mergePDB adds or replaces a PDB by name.
func mergePDB(existing []PodDisruptionBudget, rendered PodDisruptionBudget) []PodDisruptionBudget {
	for i, p := range existing {
		if namesMatch(p.Name, rendered.Name) {
			existing[i] = rendered
			return existing
		}
	}
	return append(existing, rendered)
}

// mergeHPA adds or replaces an HPA by name.
func mergeHPA(existing []HorizontalPodAutoscaler, rendered HorizontalPodAutoscaler) []HorizontalPodAutoscaler {
	for i, h := range existing {
		if namesMatch(h.Name, rendered.Name) {
			existing[i] = rendered
			return existing
		}
	}
	return append(existing, rendered)
}
