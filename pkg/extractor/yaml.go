package extractor

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// goTemplateRE matches Go template directives like {{ .Foo }} or {{- if ... -}}.
// Used by parseTemplateYAML and shared across extractors.
var goTemplateRE = regexp.MustCompile(`\{\{-?\s*.*?\s*-?\}\}`)

// DefaultExcludedDirs contains directories excluded from all extraction.
// These contain documentation, test fixtures, and build artifacts, not
// production manifests. Validated against KServe, DSPO, and MRO repos:
// none of these directories contain YAML referenced by kustomize overlays
// or deployed to production clusters.
//
// "hack" and "e2e" contain test infrastructure scripts and end-to-end test
// fixtures, not deployable manifests.
var DefaultExcludedDirs = map[string]bool{
	"docs":         true,
	"test":         true,
	"tests":        true,
	"testdata":     true,
	"examples":     true,
	"samples":      true,
	"hack":         true,
	"e2e":          true,
	"vendor":       true,
	".git":         true,
	"node_modules": true,
}

// isExcludedDir checks a directory name against the exclusion set.
// Pass nil for overrides to use DefaultExcludedDirs as-is.
// Overrides can remove defaults (value=false) or add new entries (value=true).
func isExcludedDir(dirName string, overrides map[string]bool) bool {
	if overrides != nil {
		if v, ok := overrides[dirName]; ok {
			return v
		}
	}
	return DefaultExcludedDirs[dirName]
}

// findYAMLFiles locates YAML files matching any of the given glob patterns
// relative to repoPath. Because Go's filepath.Glob does not support "**",
// patterns containing "**" are handled via filepath.WalkDir.
func findYAMLFiles(repoPath string, patterns []string) []string {
	return findFiles(repoPath, patterns)
}

// findFiles locates files matching any of the given glob patterns relative to
// repoPath. Patterns containing "**" are expanded via recursive walk.
func findFiles(repoPath string, patterns []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, pattern := range patterns {
		fullPattern := filepath.Join(repoPath, pattern)
		var matches []string

		if strings.Contains(pattern, "**") {
			matches = globRecursive(repoPath, pattern)
		} else {
			m, err := filepath.Glob(fullPattern)
			if err != nil {
				log.Printf("glob error for pattern %s: %v", pattern, err)
				continue
			}
			matches = m
		}

		for _, p := range matches {
			abs, err := filepath.Abs(p)
			if err != nil {
				continue
			}
			info, err := os.Lstat(abs)
			if err != nil || info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
				continue
			}
			if !seen[abs] {
				seen[abs] = true
				result = append(result, p)
			}
		}
	}

	return result
}

// parseGlobPattern splits a "**" glob pattern into prefix directory and
// suffix match pattern.
func parseGlobPattern(pattern string) (prefix, suffix string) {
	parts := strings.SplitN(pattern, "**", 2)
	prefix = parts[0]
	if len(parts) > 1 {
		suffix = parts[1]
		suffix = strings.TrimPrefix(suffix, "/")
		suffix = strings.TrimPrefix(suffix, string(filepath.Separator))
	}
	return prefix, suffix
}

// matchGlobSuffix checks if a file's relative path matches a suffix pattern.
// Handles simple filename patterns, multi-segment patterns, and nested "**".
func matchGlobSuffix(suffix, relPath, baseName string) bool {
	// Simple filename match
	matched, err := filepath.Match(suffix, baseName)
	if err == nil && matched {
		return true
	}

	// Multi-segment suffix (contains path separators)
	if !strings.Contains(suffix, "/") && !strings.Contains(suffix, string(filepath.Separator)) {
		return false
	}

	// Direct match on full relative path
	matched, err = filepath.Match(suffix, relPath)
	if err == nil && matched {
		return true
	}

	// Nested "**" in suffix
	if strings.Contains(suffix, "**") {
		return matchRecursiveGlob(suffix, relPath)
	}

	return false
}

// globRecursive expands a pattern containing "**" by walking the directory
// tree. The "**" segment matches zero or more intermediate directories.
func globRecursive(root, pattern string) []string {
	prefix, suffix := parseGlobPattern(pattern)
	prefixDir := filepath.Join(root, prefix)

	var results []string
	_ = filepath.WalkDir(prefixDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if isExcludedDir(d.Name(), nil) {
				return fs.SkipDir
			}
			return nil
		}
		if suffix == "" {
			results = append(results, path)
			return nil
		}
		rel, relErr := filepath.Rel(prefixDir, path)
		if relErr != nil {
			return nil
		}
		if matchGlobSuffix(suffix, rel, filepath.Base(path)) {
			results = append(results, path)
		}
		return nil
	})

	return results
}

// matchRecursiveGlob matches a path against a glob pattern containing "**".
// Unlike filepath.Match, ** matches zero or more directory segments.
func matchRecursiveGlob(pattern, path string) bool {
	parts := strings.SplitN(pattern, "**", 2)
	prefix := strings.TrimSuffix(parts[0], "/")
	suffix := ""
	if len(parts) > 1 {
		suffix = strings.TrimPrefix(parts[1], "/")
	}

	// Check prefix matches the start of path
	if prefix != "" {
		if !strings.HasPrefix(path, prefix+"/") && path != prefix {
			return false
		}
		path = strings.TrimPrefix(path, prefix+"/")
	}

	// If no suffix, any remaining path matches
	if suffix == "" {
		return true
	}

	// If suffix contains another **, recurse
	if strings.Contains(suffix, "**") {
		// Try matching suffix against every possible tail of path
		segments := strings.Split(path, "/")
		for i := 0; i <= len(segments); i++ {
			tail := strings.Join(segments[i:], "/")
			if matchRecursiveGlob(suffix, tail) {
				return true
			}
		}
		return false
	}

	// Suffix is a simple glob: try matching against every possible tail
	segments := strings.Split(path, "/")
	for i := 0; i <= len(segments); i++ {
		tail := strings.Join(segments[i:], "/")
		matched, err := filepath.Match(suffix, tail)
		if err == nil && matched {
			return true
		}
	}
	return false
}

// parseYAMLSafe parses a YAML file and returns a slice of documents (each as
// map[string]interface{}). Files with Helm template syntax ("{{") that fail to
// parse are silently skipped. Returns nil on errors.
const maxFileSize = 50 * 1024 * 1024 // 50MB

func parseYAMLSafe(path string) []map[string]interface{} {
	info, err := os.Lstat(path)
	if err != nil {
		log.Printf("cannot stat %s: %v", path, err)
		return nil
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return nil
	}
	if info.Size() > maxFileSize {
		log.Printf("skipping oversized file %s: %d bytes", path, info.Size())
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("cannot read %s: %v", path, err)
		return nil
	}

	content := string(data)
	hasHelmTemplates := strings.Contains(content, "{{") && strings.Contains(content, "}}")

	var docs []map[string]interface{}
	decoder := yaml.NewDecoder(strings.NewReader(content))
	for {
		var doc interface{}
		err := decoder.Decode(&doc)
		if err != nil {
			break
		}
		if m, ok := doc.(map[string]interface{}); ok {
			docs = append(docs, m)
		}
	}

	if len(docs) == 0 && !hasHelmTemplates {
		// Only warn for non-helm files that produced no docs
		// (helm template files are expected to fail)
	}

	return docs
}

// parseYAMLFromBytes parses raw YAML bytes (e.g. from an embedded configmap
// data value) and returns a slice of documents. Returns nil on errors.
func parseYAMLFromBytes(data []byte) []map[string]interface{} {
	var docs []map[string]interface{}
	decoder := yaml.NewDecoder(strings.NewReader(string(data)))
	for {
		var doc interface{}
		if err := decoder.Decode(&doc); err != nil {
			break
		}
		if m, ok := doc.(map[string]interface{}); ok {
			docs = append(docs, m)
		}
	}
	return docs
}

// templateCondRE matches Go template {{if ...}} and {{else if ...}} directives.
var templateCondRE = regexp.MustCompile(`\{\{-?\s*(?:if|else if)\s+(.+?)\s*-?\}\}`)

// findTemplateFiles returns enriched template file info for all .yaml.tmpl
// and .yml.tmpl files. For each file, extracts the Kubernetes resource kinds
// defined in the template and any Go template conditional guards.
func findTemplateFiles(repoPath string) []TemplateFile {
	patterns := []string{
		"**/*.yaml.tmpl",
		"**/*.yml.tmpl",
	}
	files := findFiles(repoPath, patterns)
	var result []TemplateFile
	for _, f := range files {
		relPath := relativePath(repoPath, f)
		tf := TemplateFile{Path: relPath}

		// Parse template to extract resource kinds
		docs := parseTemplateYAML(f)
		kindSeen := make(map[string]bool)
		for _, doc := range docs {
			if kind, ok := doc["kind"].(string); ok && kind != "" && !kindSeen[kind] {
				kindSeen[kind] = true
				tf.ResourceKinds = append(tf.ResourceKinds, kind)
			}
		}

		// Extract conditional guards from template directives
		data, err := os.ReadFile(f)
		if err == nil {
			condSeen := make(map[string]bool)
			for _, match := range templateCondRE.FindAllStringSubmatch(string(data), -1) {
				cond := strings.TrimSpace(match[1])
				if !condSeen[cond] {
					condSeen[cond] = true
					tf.Conditionals = append(tf.Conditionals, cond)
				}
			}
		}

		result = append(result, tf)
	}
	return result
}

// parseTemplateYAML reads a Go template file (.tmpl), strips template
// directives (if/else/end become empty, value expressions become
// "template-value"), and parses the result as YAML. Used for .yaml.tmpl
// files in operator repos that define Kubernetes resources with Go
// template placeholders.
func parseTemplateYAML(path string) []map[string]interface{} {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	content := string(data)
	cleaned := stripGoTemplateDirectives(content)
	return parseYAMLFromBytes([]byte(cleaned))
}

// stripGoTemplateDirectives replaces Go template directives with either
// empty strings (control flow) or "template-value" (expressions).
func stripGoTemplateDirectives(content string) string {
	return goTemplateRE.ReplaceAllStringFunc(content, func(match string) string {
		trimmed := strings.TrimSpace(strings.TrimPrefix(strings.TrimSuffix(match, "}}"), "{{"))
		trimmed = strings.TrimPrefix(trimmed, "-")
		trimmed = strings.TrimSuffix(trimmed, "-")
		trimmed = strings.TrimSpace(trimmed)
		if strings.HasPrefix(trimmed, "if ") || strings.HasPrefix(trimmed, "else") ||
			strings.HasPrefix(trimmed, "end") || strings.HasPrefix(trimmed, "range") ||
			strings.HasPrefix(trimmed, "define") || strings.HasPrefix(trimmed, "template") ||
			strings.HasPrefix(trimmed, "block") || strings.HasPrefix(trimmed, "with") ||
			trimmed == "-" {
			return ""
		}
		return "template-value"
	})
}

// relativePath returns path relative to repoPath, falling back to the
// absolute path on error.
func relativePath(repoPath, path string) string {
	rel, err := filepath.Rel(repoPath, path)
	if err != nil {
		return path
	}
	return rel
}
