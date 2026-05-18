// Package renderer converts architecture JSON data into diagram formats
// (Mermaid, ASCII, Structurizr C4 DSL).
package renderer

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

// Renderer produces a diagram string from architecture data.
type Renderer interface {
	// Render produces the diagram content from the given data.
	Render(data map[string]interface{}) string
	// Filename returns the suggested output filename for this renderer.
	Filename() string
}

// rendererRegistry is the single source of truth for all renderers.
// Add new renderers here only.
var rendererRegistry = map[string]Renderer{
	"rbac":             &RBACRenderer{},
	"component":        &ComponentRenderer{},
	"security_network": &SecurityNetworkRenderer{},
	"dependencies":     &DependencyRenderer{},
	"c4":               &C4Renderer{},
	"dataflow":         &DataflowRenderer{},
	"report":           &ReportRenderer{},
	"flow":             &FlowRenderer{},
}

// allRenderers returns every available per-component renderer in stable order.
func allRenderers() []Renderer {
	names := make([]string, 0, len(rendererRegistry))
	for name := range rendererRegistry {
		names = append(names, name)
	}
	sort.Strings(names)
	out := make([]Renderer, 0, len(names))
	for _, name := range names {
		out = append(out, rendererRegistry[name])
	}
	return out
}

// RenderAll runs the selected renderers (by short name) against data and
// returns a filename->content map. If formats is empty, all renderers run.
func RenderAll(data map[string]interface{}, formats []string) map[string]string {
	results := make(map[string]string)

	var renderers []Renderer
	if len(formats) == 0 {
		renderers = allRenderers()
	} else {
		for _, name := range formats {
			if r, ok := rendererRegistry[name]; ok {
				renderers = append(renderers, r)
			}
		}
	}

	for _, r := range renderers {
		results[r.Filename()] = r.Render(data)
	}
	return results
}

// sanitizeID replaces non-alphanumeric characters with underscores and
// prepends "n_" if the result starts with a non-letter.
func sanitizeID(text string) string {
	if text == "" {
		return "node"
	}
	var b strings.Builder
	for _, ch := range text {
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' {
			b.WriteRune(ch)
		} else {
			b.WriteByte('_')
		}
	}
	result := b.String()
	if result == "" {
		return "node"
	}
	if !unicode.IsLetter(rune(result[0])) {
		result = "n_" + result
	}
	return result
}

// escapeLabel escapes special characters for Mermaid labels.
func escapeLabel(text string) string {
	text = strings.ReplaceAll(text, `"`, "'")
	text = strings.ReplaceAll(text, "<", "&lt;")
	text = strings.ReplaceAll(text, ">", "&gt;")
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\r", " ")
	text = strings.ReplaceAll(text, "|", "/")
	return text
}

// escapeMdCell escapes characters that break markdown table cells.
func escapeMdCell(text string) string {
	text = strings.ReplaceAll(text, "|", "\\|")
	text = strings.ReplaceAll(text, "\n", " ")
	return text
}


// RoleSummary holds computed permission breadth for a single RBAC role.
type RoleSummary struct {
	Name          string
	Kind          string // "ClusterRole" or "Role"
	ResourceCount int
	HasWildcard   bool
}

// computeRoleSummary extracts permission breadth from a role data map.
func computeRoleSummary(role map[string]interface{}, kind string) RoleSummary {
	name := getStr(role, "name", "")
	rules := getSlice(role, "rules")
	resourceCount := 0
	hasWildcard := false
	for _, rule := range rules {
		resourceCount += len(getStringSlice(rule, "resources"))
		for _, v := range getStringSlice(rule, "verbs") {
			if v == "*" {
				hasWildcard = true
			}
		}
	}
	return RoleSummary{Name: name, Kind: kind, ResourceCount: resourceCount, HasWildcard: hasWildcard}
}

// sourceLink generates a GitHub permalink for a source file reference.
// If commit_sha and repo are available, returns a clickable markdown link.
// Otherwise returns the source path in backticks.
func sourceLink(data map[string]interface{}, source string) string {
	sha := getStr(data, "commit_sha", "")
	repo := getStr(data, "repo", "")
	if sha == "" || repo == "" || source == "" {
		return fmt.Sprintf("`%s`", source)
	}
	// Split source into file:line if present
	file := source
	anchor := ""
	if idx := strings.LastIndex(source, ":"); idx > 0 {
		possibleLine := source[idx+1:]
		if _, err := fmt.Sscanf(possibleLine, "%d", new(int)); err == nil {
			file = source[:idx]
			anchor = fmt.Sprintf("#L%s", possibleLine)
		}
	}
	url := fmt.Sprintf("https://github.com/%s/blob/%s/%s%s", repo, sha, file, anchor)
	return fmt.Sprintf("[`%s`](%s)", source, url)
}

// GetStr is the exported version of getStr for use by CLI commands.
func GetStr(m map[string]interface{}, key, fallback string) string {
	return getStr(m, key, fallback)
}

// GetSlice is the exported version of getSlice for use by CLI commands.
func GetSlice(m map[string]interface{}, key string) []map[string]interface{} {
	return getSlice(m, key)
}

// helper to extract a string from a map with a default.
func getStr(m map[string]interface{}, key, fallback string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return fallback
}

// helper to extract a slice of maps from a map.
func getSlice(m map[string]interface{}, key string) []map[string]interface{} {
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

// helper to extract a map from a map.
func getMap(m map[string]interface{}, key string) map[string]interface{} {
	if v, ok := m[key]; ok {
		if mm, ok := v.(map[string]interface{}); ok {
			return mm
		}
	}
	return nil
}

// helper to extract a string slice from a map value.
func getStringSlice(m map[string]interface{}, key string) []string {
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

// helper to get an int from a map.
func getInt(m map[string]interface{}, key string) int {
	v, ok := m[key]
	if !ok {
		return 0
	}
	switch n := v.(type) {
	case int:
		return n
	case float64:
		return int(n)
	case int64:
		return int(n)
	}
	return 0
}

// helper to get a bool from a map, with a default.
func getBool(m map[string]interface{}, key string, fallback bool) bool {
	v, ok := m[key]
	if !ok {
		return fallback
	}
	if b, ok := v.(bool); ok {
		return b
	}
	return fallback
}
