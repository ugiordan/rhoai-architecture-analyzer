// Package sbom generates CycloneDX Software Bill of Materials from
// architecture-analyzer component-architecture.json output.
package sbom

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// CycloneDX 1.5 types

type BOM struct {
	BOMFormat    string       `json:"bomFormat"`
	SpecVersion  string       `json:"specVersion"`
	Version      int          `json:"version"`
	SerialNumber string       `json:"serialNumber,omitempty"`
	Metadata     Metadata     `json:"metadata"`
	Components   []Component  `json:"components"`
}

type Metadata struct {
	Timestamp string     `json:"timestamp"`
	Tools     []Tool     `json:"tools,omitempty"`
	Component *Component `json:"component,omitempty"`
}

type Tool struct {
	Vendor  string `json:"vendor"`
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

type Component struct {
	Type        string       `json:"type"`
	Name        string       `json:"name"`
	Version     string       `json:"version,omitempty"`
	PURL        string       `json:"purl,omitempty"`
	Description string       `json:"description,omitempty"`
	Group       string       `json:"group,omitempty"`
	Scope       string       `json:"scope,omitempty"`
	Hashes      []Hash       `json:"hashes,omitempty"`
	Licenses    []License    `json:"licenses,omitempty"`
	Properties  []Property   `json:"properties,omitempty"`
	Components  []Component  `json:"components,omitempty"`
}

type Hash struct {
	Algorithm string `json:"alg"`
	Content   string `json:"content"`
}

type License struct {
	License struct {
		ID string `json:"id,omitempty"`
	} `json:"license"`
}

type Property struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Generate produces a CycloneDX 1.5 BOM from component-architecture.json data.
func Generate(data map[string]interface{}) *BOM {
	component := getStr(data, "component", "unknown")
	repo := getStr(data, "repo", "")
	version := getStr(data, "analyzer_version", "")
	commitSHA := getStr(data, "commit_sha", "")

	bom := &BOM{
		BOMFormat:   "CycloneDX",
		SpecVersion: "1.5",
		Version:     1,
		Metadata: Metadata{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Tools: []Tool{
				{Vendor: "architecture-analyzer", Name: "arch-analyzer", Version: version},
			},
			Component: &Component{
				Type:    "application",
				Name:    component,
				Version: commitSHA,
				PURL:    fmt.Sprintf("pkg:github/%s@%s", repo, commitSHA),
				Properties: []Property{
					{Name: "arch-analyzer:repo", Value: repo},
					{Name: "arch-analyzer:commit", Value: commitSHA},
				},
			},
		},
		Components: []Component{},
	}

	// Go module dependencies
	deps := getMap(data, "dependencies")
	if deps != nil {
		goMods := getSlice(deps, "go_modules")
		for _, mod := range goMods {
			module := getStr(mod, "module", "")
			ver := getStr(mod, "version", "")
			if module == "" {
				continue
			}
			c := Component{
				Type:    "library",
				Name:    module,
				Version: ver,
				PURL:    fmt.Sprintf("pkg:golang/%s@%s", module, ver),
				Scope:   "required",
				Properties: []Property{
					{Name: "arch-analyzer:source", Value: "go.mod"},
					{Name: "arch-analyzer:ecosystem", Value: "go"},
				},
			}
			bom.Components = append(bom.Components, c)
		}

		// Replace directives
		replaces := getSlice(deps, "replace_directives")
		for _, rep := range replaces {
			old := getStr(rep, "old", "")
			newMod := getStr(rep, "new", "")
			newVer := getStr(rep, "new_version", "")
			if old == "" {
				continue
			}
			c := Component{
				Type:    "library",
				Name:    old,
				Version: newVer,
				PURL:    fmt.Sprintf("pkg:golang/%s@%s", newMod, newVer),
				Scope:   "required",
				Properties: []Property{
					{Name: "arch-analyzer:source", Value: "go.mod (replace)"},
					{Name: "arch-analyzer:ecosystem", Value: "go"},
					{Name: "arch-analyzer:replaces", Value: old},
					{Name: "arch-analyzer:replaced-by", Value: newMod},
				},
			}
			bom.Components = append(bom.Components, c)
		}
	}

	// Python dependencies
	for _, dep := range getSlice(data, "python_deps") {
		name := getStr(dep, "name", "")
		ver := getStr(dep, "version", "")
		source := getStr(dep, "source", "")
		if name == "" {
			continue
		}
		c := Component{
			Type:    "library",
			Name:    name,
			Version: ver,
			PURL:    fmt.Sprintf("pkg:pypi/%s@%s", strings.ToLower(name), ver),
			Scope:   "required",
			Properties: []Property{
				{Name: "arch-analyzer:source", Value: source},
				{Name: "arch-analyzer:ecosystem", Value: "python"},
			},
		}
		bom.Components = append(bom.Components, c)
	}

	// Dockerfile base images
	for _, df := range getSlice(data, "dockerfiles") {
		baseImg := getStr(df, "base_image", "")
		path := getStr(df, "path", "")
		user := getStr(df, "user", "")
		stages := getInt(df, "stages")
		fips := getBool(df, "fips_enabled")
		if baseImg == "" || strings.HasPrefix(baseImg, "$") {
			continue
		}

		// Parse image reference
		imgName, imgTag, imgDigest := parseImageRef(baseImg)
		purl := fmt.Sprintf("pkg:docker/%s", imgName)
		if imgDigest != "" {
			purl += "@" + imgDigest
		} else if imgTag != "" {
			purl += "@" + imgTag
		}

		props := []Property{
			{Name: "arch-analyzer:source", Value: path},
			{Name: "arch-analyzer:type", Value: "dockerfile-base"},
			{Name: "arch-analyzer:stages", Value: fmt.Sprintf("%d", stages)},
		}
		if user != "" {
			props = append(props, Property{Name: "arch-analyzer:user", Value: user})
		}
		if fips {
			props = append(props, Property{Name: "arch-analyzer:fips", Value: "true"})
		}
		for _, arch := range getStringSlice(df, "architectures") {
			props = append(props, Property{Name: "arch-analyzer:architecture", Value: arch})
		}
		for _, issue := range getStringSlice(df, "issues") {
			props = append(props, Property{Name: "arch-analyzer:issue", Value: issue})
		}

		c := Component{
			Type:       "container",
			Name:       imgName,
			Version:    imgTag,
			PURL:       purl,
			Properties: props,
		}
		if imgDigest != "" {
			c.Hashes = []Hash{{Algorithm: "SHA-256", Content: strings.TrimPrefix(imgDigest, "sha256:")}}
		}
		bom.Components = append(bom.Components, c)
	}

	// Container images from deployments
	for _, dep := range getSlice(data, "deployments") {
		depName := getStr(dep, "name", "")
		for _, container := range getSlice(dep, "containers") {
			img := getStr(container, "image", "")
			containerName := getStr(container, "name", "")
			if img == "" || strings.HasPrefix(img, "$") || strings.HasPrefix(img, "ko://") {
				continue
			}

			imgName, imgTag, imgDigest := parseImageRef(img)
			purl := fmt.Sprintf("pkg:docker/%s", imgName)
			if imgDigest != "" {
				purl += "@" + imgDigest
			} else if imgTag != "" {
				purl += "@" + imgTag
			}

			props := []Property{
				{Name: "arch-analyzer:type", Value: "deployment-container"},
				{Name: "arch-analyzer:deployment", Value: depName},
				{Name: "arch-analyzer:container", Value: containerName},
			}

			// Security context
			sc := getMap(container, "security_context")
			if sc != nil {
				if v, ok := sc["runAsNonRoot"]; ok {
					props = append(props, Property{Name: "arch-analyzer:runAsNonRoot", Value: fmt.Sprintf("%v", v)})
				}
				if v, ok := sc["readOnlyRootFilesystem"]; ok {
					props = append(props, Property{Name: "arch-analyzer:readOnlyRootFilesystem", Value: fmt.Sprintf("%v", v)})
				}
				if v, ok := sc["privileged"]; ok {
					props = append(props, Property{Name: "arch-analyzer:privileged", Value: fmt.Sprintf("%v", v)})
				}
			}

			// Resources
			res := getMap(container, "resources")
			if res != nil {
				req := getMap(res, "requests")
				lim := getMap(res, "limits")
				if req != nil {
					if v, ok := req["cpu"]; ok {
						props = append(props, Property{Name: "arch-analyzer:cpu-request", Value: fmt.Sprintf("%v", v)})
					}
					if v, ok := req["memory"]; ok {
						props = append(props, Property{Name: "arch-analyzer:memory-request", Value: fmt.Sprintf("%v", v)})
					}
				}
				if lim != nil {
					if v, ok := lim["cpu"]; ok {
						props = append(props, Property{Name: "arch-analyzer:cpu-limit", Value: fmt.Sprintf("%v", v)})
					}
					if v, ok := lim["memory"]; ok {
						props = append(props, Property{Name: "arch-analyzer:memory-limit", Value: fmt.Sprintf("%v", v)})
					}
				}
			}

			// Probes
			if lp := getMap(container, "liveness_probe"); lp != nil {
				props = append(props, Property{Name: "arch-analyzer:liveness-probe", Value: fmt.Sprintf("%s %s", getStr(lp, "type", ""), getStr(lp, "path", ""))})
			}
			if rp := getMap(container, "readiness_probe"); rp != nil {
				props = append(props, Property{Name: "arch-analyzer:readiness-probe", Value: fmt.Sprintf("%s %s", getStr(rp, "type", ""), getStr(rp, "path", ""))})
			}

			c := Component{
				Type:       "container",
				Name:       imgName,
				Version:    imgTag,
				PURL:       purl,
				Properties: props,
			}
			if imgDigest != "" {
				c.Hashes = []Hash{{Algorithm: "SHA-256", Content: strings.TrimPrefix(imgDigest, "sha256:")}}
			}
			bom.Components = append(bom.Components, c)
		}
	}

	// Operator image constants
	for _, cfg := range getSlice(data, "operator_config") {
		cat := getStr(cfg, "category", "")
		name := getStr(cfg, "name", "")
		value := getStr(cfg, "value", "")
		if (cat != "image" && !strings.Contains(strings.ToLower(name), "image")) || value == "" {
			continue
		}
		// Skip Go-level references (not actual image URIs)
		if strings.Contains(value, ".") && (strings.Contains(value, "/") || strings.Contains(value, ":")) {
			imgName, imgTag, _ := parseImageRef(value)
			c := Component{
				Type:    "container",
				Name:    imgName,
				Version: imgTag,
				PURL:    fmt.Sprintf("pkg:docker/%s@%s", imgName, imgTag),
				Properties: []Property{
					{Name: "arch-analyzer:type", Value: "operator-constant"},
					{Name: "arch-analyzer:constant-name", Value: name},
					{Name: "arch-analyzer:source", Value: getStr(cfg, "source", "")},
				},
			}
			bom.Components = append(bom.Components, c)
		}
	}

	return bom
}

// GenerateJSON returns the SBOM as indented JSON bytes.
func GenerateJSON(data map[string]interface{}) ([]byte, error) {
	bom := Generate(data)
	return json.MarshalIndent(bom, "", "  ")
}

// parseImageRef splits "registry.io/org/name:tag@sha256:abc" into name, tag, digest.
func parseImageRef(ref string) (name, tag, digest string) {
	if idx := strings.Index(ref, "@"); idx >= 0 {
		digest = ref[idx+1:]
		ref = ref[:idx]
	}
	if idx := strings.LastIndex(ref, ":"); idx >= 0 {
		// Distinguish port from tag: if the part after : contains / it's a port
		after := ref[idx+1:]
		if !strings.Contains(after, "/") {
			tag = after
			ref = ref[:idx]
		}
	}
	name = ref
	return
}

// Helpers matching the renderer package pattern.
func getStr(m map[string]interface{}, key, fallback string) string {
	if m == nil {
		return fallback
	}
	v, ok := m[key]
	if !ok {
		return fallback
	}
	s, ok := v.(string)
	if !ok {
		return fallback
	}
	return s
}

func getMap(m map[string]interface{}, key string) map[string]interface{} {
	if m == nil {
		return nil
	}
	v, ok := m[key]
	if !ok {
		return nil
	}
	mm, ok := v.(map[string]interface{})
	if ok {
		return mm
	}
	return nil
}

func getSlice(m map[string]interface{}, key string) []map[string]interface{} {
	if m == nil {
		return nil
	}
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

func getStringSlice(m map[string]interface{}, key string) []string {
	if m == nil {
		return nil
	}
	v, ok := m[key]
	if !ok {
		return nil
	}
	switch typed := v.(type) {
	case []string:
		return typed
	case []interface{}:
		out := make([]string, 0, len(typed))
		for _, item := range typed {
			if s, ok := item.(string); ok {
				out = append(out, s)
			}
		}
		return out
	}
	return nil
}

func getInt(m map[string]interface{}, key string) int {
	if m == nil {
		return 0
	}
	v, ok := m[key]
	if !ok {
		return 0
	}
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	case int64:
		return int(n)
	}
	return 0
}

func getBool(m map[string]interface{}, key string) bool {
	if m == nil {
		return false
	}
	v, ok := m[key]
	if !ok {
		return false
	}
	b, ok := v.(bool)
	return ok && b
}
