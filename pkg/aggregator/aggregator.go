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

	// Build aggregated structures
	var (
		allCRDs             []map[string]interface{}
		allServices         []map[string]interface{}
		allSecrets          []map[string]interface{}
		allRBACClusterRoles []map[string]interface{}
		componentNames      []string
		dependencyGraph     []map[string]string
	)
	crdOwners := make(map[string]string)

	for _, compData := range components {
		compName := getString(compData, "component", "unknown")
		componentNames = append(componentNames, compName)

		// CRDs
		for _, crd := range getSliceOfMaps(compData, "crds") {
			crdWithOwner := copyMap(crd)
			crdWithOwner["owner"] = compName
			allCRDs = append(allCRDs, crdWithOwner)
			if kind := getString(crd, "kind", ""); kind != "" {
				crdOwners[kind] = compName
			}
		}

		// Services
		for _, svc := range getSliceOfMaps(compData, "services") {
			s := copyMap(svc)
			s["owner"] = compName
			allServices = append(allServices, s)
		}

		// Secrets
		for _, secret := range getSliceOfMaps(compData, "secrets_referenced") {
			s := copyMap(secret)
			s["owner"] = compName
			allSecrets = append(allSecrets, s)
		}

		// RBAC cluster roles
		rbac, _ := compData["rbac"].(map[string]interface{})
		if rbac != nil {
			for _, cr := range getSliceOfMaps(rbac, "cluster_roles") {
				c := copyMap(cr)
				c["owner"] = compName
				allRBACClusterRoles = append(allRBACClusterRoles, c)
			}
		}

		// Dependencies (internal ODH)
		deps, _ := compData["dependencies"].(map[string]interface{})
		if deps != nil {
			for _, odh := range getSliceOfMaps(deps, "internal_odh") {
				dependencyGraph = append(dependencyGraph, map[string]string{
					"from": compName,
					"to":   getString(odh, "component", ""),
					"type": "go-module",
				})
			}
		}

		// Cross-component watches
		for _, watch := range getSliceOfMaps(compData, "controller_watches") {
			if getString(watch, "type", "") != "For" {
				continue
			}
			gvk := getString(watch, "gvk", "")
			kind := gvk
			if idx := strings.LastIndex(gvk, "/"); idx >= 0 {
				kind = gvk[idx+1:]
			}
			if owner, ok := crdOwners[kind]; ok && owner != compName {
				dependencyGraph = append(dependencyGraph, map[string]string{
					"from": compName,
					"to":   owner,
					"type": "watches-crd:" + kind,
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

	compDataIface := make([]interface{}, len(components))
	for i, c := range components {
		compDataIface[i] = c
	}

	return map[string]interface{}{
		"platform":           "OpenShift AI",
		"aggregated_at":      time.Now().UTC().Format(time.RFC3339),
		"components":         compNamesIface,
		"component_count":    len(components),
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

func copyMap(m map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
