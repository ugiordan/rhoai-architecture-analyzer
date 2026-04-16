package extractor

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

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

// globRecursive expands a pattern containing "**" by walking the directory
// tree. The "**" segment matches zero or more intermediate directories.
func globRecursive(root, pattern string) []string {
	// Split pattern on "**" to get prefix and suffix parts.
	parts := strings.SplitN(pattern, "**", 2)
	prefix := parts[0]
	suffix := ""
	if len(parts) > 1 {
		suffix = parts[1]
		// Remove leading separator from suffix
		suffix = strings.TrimPrefix(suffix, "/")
		suffix = strings.TrimPrefix(suffix, string(filepath.Separator))
	}

	prefixDir := filepath.Join(root, prefix)

	var results []string
	_ = filepath.WalkDir(prefixDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		if suffix == "" {
			results = append(results, path)
			return nil
		}
		// Check if the file name (or relative tail) matches the suffix pattern
		rel, relErr := filepath.Rel(prefixDir, path)
		if relErr != nil {
			return nil
		}
		matched, matchErr := filepath.Match(suffix, filepath.Base(path))
		if matchErr != nil {
			return nil
		}
		if matched {
			results = append(results, path)
			return nil
		}
		// For multi-segment suffixes (e.g., "network-policies/**/*.yaml"),
		// try matching each tail segment of the relative path.
		// filepath.Match doesn't support ** as zero-or-more dirs, so we
		// check all possible sub-paths to handle arbitrary nesting.
		if strings.Contains(suffix, "/") || strings.Contains(suffix, string(filepath.Separator)) {
			// Direct match on full relative path
			matched, matchErr = filepath.Match(suffix, rel)
			if matchErr == nil && matched {
				results = append(results, path)
				return nil
			}
			// If suffix itself contains **, expand recursively
			if strings.Contains(suffix, "**") {
				if matchRecursiveGlob(suffix, rel) {
					results = append(results, path)
				}
			}
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

// relativePath returns path relative to repoPath, falling back to the
// absolute path on error.
func relativePath(repoPath, path string) string {
	rel, err := filepath.Rel(repoPath, path)
	if err != nil {
		return path
	}
	return rel
}
