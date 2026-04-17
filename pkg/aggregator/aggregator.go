// Package aggregator merges per-component architecture JSONs into a
// platform-wide view.
package aggregator

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Aggregate reads all component-architecture.json files under resultsDir,
// merges CRDs, services, secrets, RBAC, dependencies, and cross-component
// watches into a single platform-level map.
func Aggregate(resultsDir string) (map[string]interface{}, error) {
	absDir, err := filepath.Abs(resultsDir)
	if err != nil {
		return nil, fmt.Errorf("resolving results dir: %w", err)
	}
	info, err := os.Stat(absDir)
	if err != nil || !info.IsDir() {
		return nil, fmt.Errorf("results directory does not exist: %s", absDir)
	}

	// Find all component-architecture.json files
	var jsonPaths []string
	const maxComponentFileSize = 50 * 1024 * 1024 // 50 MB
	err = filepath.Walk(absDir, func(path string, fi os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return nil // skip unreadable entries
		}
		if !fi.IsDir() && fi.Name() == "component-architecture.json" {
			if fi.Size() > maxComponentFileSize {
				log.Printf("WARN: skipping oversized file %s (%d bytes)", path, fi.Size())
				return nil
			}
			jsonPaths = append(jsonPaths, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walking results dir: %w", err)
	}
	sort.Strings(jsonPaths)

	// Load all components
	var components []map[string]interface{}
	for _, jp := range jsonPaths {
		raw, readErr := os.ReadFile(jp)
		if readErr != nil {
			log.Printf("WARN: failed to read %s: %v", jp, readErr)
			continue
		}
		var data map[string]interface{}
		if err := json.Unmarshal(raw, &data); err != nil {
			log.Printf("WARN: failed to parse %s: %v", jp, err)
			continue
		}
		components = append(components, data)
	}

	// Build aggregated structures with deduplication
	var (
		allCRDs             []map[string]interface{}
		allServices         []map[string]interface{}
		allSecrets          []map[string]interface{}
		allRBACClusterRoles []map[string]interface{}
		componentNames      []string
		dependencyGraph     []map[string]string
	)
	crdOwners := make(map[string]string)

	// Dedup tracking maps
	seenComponents := make(map[string]bool)
	seenCRDs := make(map[string]bool)       // key: "owner|group|version|kind"
	seenServices := make(map[string]bool)   // key: "owner|name|type"
	seenSecrets := make(map[string]bool)    // key: "owner|name"
	seenRBACRoles := make(map[string]bool)  // key: "owner|name"
	seenDepEdges := make(map[string]bool)   // key: "from|to|type"

	for _, compData := range components {
		compName := getString(compData, "component", "unknown")
		if !seenComponents[compName] {
			seenComponents[compName] = true
			componentNames = append(componentNames, compName)
		}

		// CRDs: dedup on (owner, group, version, kind)
		for _, crd := range getSliceOfMaps(compData, "crds") {
			kind := getString(crd, "kind", "")
			key := compName + "|" + getString(crd, "group", "") + "|" + getString(crd, "version", "") + "|" + kind
			if seenCRDs[key] {
				continue
			}
			seenCRDs[key] = true
			crdWithOwner := copyMap(crd)
			crdWithOwner["owner"] = compName
			allCRDs = append(allCRDs, crdWithOwner)
			if kind != "" {
				if _, exists := crdOwners[kind]; !exists {
					crdOwners[kind] = compName
				}
			}
		}

		// Services: dedup on (owner, name, type)
		for _, svc := range getSliceOfMaps(compData, "services") {
			key := compName + "|" + getString(svc, "name", "") + "|" + getString(svc, "type", "")
			if seenServices[key] {
				continue
			}
			seenServices[key] = true
			s := copyMap(svc)
			s["owner"] = compName
			allServices = append(allServices, s)
		}

		// Secrets: dedup on (owner, name)
		for _, secret := range getSliceOfMaps(compData, "secrets_referenced") {
			key := compName + "|" + getString(secret, "name", "")
			if seenSecrets[key] {
				continue
			}
			seenSecrets[key] = true
			s := copyMap(secret)
			s["owner"] = compName
			allSecrets = append(allSecrets, s)
		}

		// RBAC cluster roles: dedup on (owner, name)
		rbac, _ := compData["rbac"].(map[string]interface{})
		if rbac != nil {
			for _, cr := range getSliceOfMaps(rbac, "cluster_roles") {
				key := compName + "|" + getString(cr, "name", "")
				if seenRBACRoles[key] {
					continue
				}
				seenRBACRoles[key] = true
				c := copyMap(cr)
				c["owner"] = compName
				allRBACClusterRoles = append(allRBACClusterRoles, c)
			}
		}

		// Dependencies (internal ODH): dedup on (from, to, type)
		deps, _ := compData["dependencies"].(map[string]interface{})
		if deps != nil {
			for _, odh := range getSliceOfMaps(deps, "internal_odh") {
				to := getString(odh, "component", "")
				key := compName + "|" + to + "|go-module"
				if seenDepEdges[key] {
					continue
				}
				seenDepEdges[key] = true
				dependencyGraph = append(dependencyGraph, map[string]string{
					"from": compName,
					"to":   to,
					"type": "go-module",
				})
			}
		}

		// Cross-component watches: dedup on (from, to, type)
		for _, watch := range getSliceOfMaps(compData, "controller_watches") {
			if getString(watch, "type", "") != "For" {
				continue
			}
			gvk := getString(watch, "gvk", "")
			kind := gvk
			if idx := strings.LastIndex(gvk, "/"); idx >= 0 {
				kind = gvk[idx+1:]
			}
			if kind == "" {
				log.Printf("WARN: malformed GVK '%s' in component %s, skipping", gvk, compName)
				continue
			}
			if owner, ok := crdOwners[kind]; ok && owner != compName {
				edgeType := "watches-crd:" + kind
				key := compName + "|" + owner + "|" + edgeType
				if seenDepEdges[key] {
					continue
				}
				seenDepEdges[key] = true
				dependencyGraph = append(dependencyGraph, map[string]string{
					"from": compName,
					"to":   owner,
					"type": edgeType,
				})
			}
		}
	}

	// Convert crdOwners to interface map for JSON compatibility
	crdOwnershipIface := make(map[string]interface{}, len(crdOwners))
	for k, v := range crdOwners {
		crdOwnershipIface[k] = v
	}

	// Convert dependencyGraph to interface slice
	depGraphIface := make([]interface{}, len(dependencyGraph))
	for i, d := range dependencyGraph {
		m := make(map[string]interface{}, len(d))
		for k, v := range d {
			m[k] = v
		}
		depGraphIface[i] = m
	}

	// Convert slices to interface slices
	toIfaceSlice := func(s []map[string]interface{}) []interface{} {
		out := make([]interface{}, len(s))
		for i, v := range s {
			out[i] = v
		}
		return out
	}

	compNamesIface := make([]interface{}, len(componentNames))
	for i, n := range componentNames {
		compNamesIface[i] = n
	}

	// Dedup component_data: keep first occurrence per component name
	seenCompData := make(map[string]bool)
	var dedupedComponents []map[string]interface{}
	for _, c := range components {
		name := getString(c, "component", "unknown")
		if seenCompData[name] {
			continue
		}
		seenCompData[name] = true
		dedupedComponents = append(dedupedComponents, c)
	}
	compDataIface := make([]interface{}, len(dedupedComponents))
	for i, c := range dedupedComponents {
		compDataIface[i] = c
	}

	// Sidecar image references: detect when component A embeds component B's
	// image as a sidecar container (e.g., kube-rbac-proxy used by kserve).
	compImagePatterns := make(map[string][]string) // component -> list of image base names
	for _, compData := range dedupedComponents {
		compName := getString(compData, "component", "unknown")
		for _, dep := range getSliceOfMaps(compData, "deployments") {
			for _, c := range getSliceOfMaps(dep, "containers") {
				img := getString(c, "image", "")
				if img == "" {
					continue
				}
				base := img
				if idx := strings.LastIndex(base, ":"); idx > 0 {
					base = base[:idx]
				}
				if idx := strings.LastIndex(base, "@"); idx > 0 {
					base = base[:idx]
				}
				if idx := strings.LastIndex(base, "/"); idx >= 0 {
					base = base[idx+1:]
				}
				// Match if the image base is specific enough and matches the component name.
				// Skip generic bases that cause false positives.
				if len(base) < 5 || isGenericImageBase(base) {
					continue
				}
				normalizedComp := strings.ReplaceAll(compName, "-", "")
				normalizedBase := strings.ReplaceAll(base, "-", "")
				if normalizedBase == normalizedComp {
					compImagePatterns[compName] = append(compImagePatterns[compName], base)
				}
			}
		}
	}
	for comp, patterns := range compImagePatterns {
		seen := make(map[string]bool)
		var unique []string
		for _, p := range patterns {
			if !seen[p] {
				seen[p] = true
				unique = append(unique, p)
			}
		}
		compImagePatterns[comp] = unique
	}
	for _, compData := range dedupedComponents {
		compName := getString(compData, "component", "unknown")
		for _, dep := range getSliceOfMaps(compData, "deployments") {
			for _, c := range getSliceOfMaps(dep, "containers") {
				img := getString(c, "image", "")
				cName := getString(c, "name", "")
				if img == "" && cName == "" {
					continue
				}
				for otherComp, patterns := range compImagePatterns {
					if otherComp == compName {
						continue
					}
					for _, pattern := range patterns {
						// Match on container name or image path component
						matched := cName == pattern
						if !matched && img != "" {
							// Check if pattern appears as a path segment in the image reference
							// e.g., "quay.io/brancz/kube-rbac-proxy:v0.18.0" contains "/kube-rbac-proxy:"
							matched = strings.Contains(img, "/"+pattern+":") ||
								strings.Contains(img, "/"+pattern+"@") ||
								strings.HasSuffix(img, "/"+pattern)
						}
						if matched {
							key := compName + "|" + otherComp + "|uses-image"
							if !seenDepEdges[key] {
								seenDepEdges[key] = true
								dependencyGraph = append(dependencyGraph, map[string]string{
									"from": compName,
									"to":   otherComp,
									"type": "uses-image",
								})
							}
							break
						}
					}
				}
			}
		}
	}

	// Rebuild depGraphIface with new sidecar edges
	depGraphIface = make([]interface{}, len(dependencyGraph))
	for i, d := range dependencyGraph {
		m := make(map[string]interface{}, len(d))
		for k, v := range d {
			m[k] = v
		}
		depGraphIface[i] = m
	}

	return map[string]interface{}{
		"platform":           "OpenShift AI",
		"aggregated_at":      time.Now().UTC().Format(time.RFC3339),
		"components":         compNamesIface,
		"component_count":    len(dedupedComponents),
		"crds":               toIfaceSlice(allCRDs),
		"crd_ownership":      crdOwnershipIface,
		"services":           toIfaceSlice(allServices),
		"secrets_referenced": toIfaceSlice(allSecrets),
		"rbac_cluster_roles": toIfaceSlice(allRBACClusterRoles),
		"dependency_graph":   depGraphIface,
		"component_data":     compDataIface,
	}, nil
}

func getString(m map[string]interface{}, key, fallback string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return fallback
}

func getSliceOfMaps(m map[string]interface{}, key string) []map[string]interface{} {
	v, ok := m[key]
	if !ok {
		return nil
	}
	switch typed := v.(type) {
	case []map[string]interface{}:
		return typed
	case []interface{}:
		out := make([]map[string]interface{}, 0, len(typed))
		for _, item := range typed {
			if mm, ok := item.(map[string]interface{}); ok {
				out = append(out, mm)
			}
		}
		return out
	}
	return nil
}

// isGenericImageBase returns true for image base names that are too generic
// to use for cross-component matching (e.g., "controller", "manager", "proxy").
func isGenericImageBase(base string) bool {
	generic := []string{"controller", "manager", "proxy", "server", "operator", "agent", "sidecar", "init", "busybox", "alpine", "nginx", "redis"}
	lower := strings.ToLower(base)
	for _, g := range generic {
		if lower == g {
			return true
		}
	}
	return false
}

func copyMap(m map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
