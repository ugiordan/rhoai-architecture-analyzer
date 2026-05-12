package extractor

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// requirementLineRE matches a pip requirements line: package[extras]>=version,<version
var requirementLineRE = regexp.MustCompile(`^([a-zA-Z0-9][a-zA-Z0-9._-]*)\s*(?:\[.*?\])?\s*([><=!~].+)?`)

// pyprojectDepRE matches dependencies in pyproject.toml arrays: "package>=version"
var pyprojectDepRE = regexp.MustCompile(`^\s*"([a-zA-Z0-9][a-zA-Z0-9._-]*)\s*(?:\[.*?\])?\s*([><=!~][^"]*)?"\s*,?\s*$`)

// poetryDepRE matches Poetry-style dependencies: package = "^version" or package = ">=version"
var poetryDepRE = regexp.MustCompile(`^([a-zA-Z0-9][a-zA-Z0-9._-]*)\s*=\s*"([^"]*)"`)


// pythonPackageCategoryRules maps package name prefixes/names to categories.
var pythonPackageCategoryRules = []struct {
	name     string
	exact    bool
	category string
}{
	// ML/AI frameworks
	{"torch", true, "ml-framework"},
	{"pytorch", false, "ml-framework"},
	{"transformers", true, "ml-framework"},
	{"tensorflow", false, "ml-framework"},
	{"jax", true, "ml-framework"},

	// Web frameworks
	{"fastapi", true, "web-framework"},
	{"flask", true, "web-framework"},
	{"uvicorn", true, "web-server"},
	{"gunicorn", true, "web-server"},
	{"aiohttp", true, "web-framework"},
	{"starlette", true, "web-framework"},

	// gRPC/serialization
	{"grpcio", false, "grpc"},
	{"protobuf", true, "serialization"},
	{"pydantic", true, "serialization"},

	// Observability
	{"prometheus", false, "observability"},
	{"opentelemetry", false, "observability"},

	// Database
	{"sqlalchemy", true, "database"},
	{"psycopg", false, "database"},
	{"pymysql", true, "database"},
	{"redis", true, "database"},

	// Cloud/storage
	{"boto3", true, "cloud-provider"},
	{"minio", true, "storage"},
	{"s3fs", true, "storage"},

	// Kubernetes
	{"kubernetes", true, "k8s-client"},
	{"openshift", false, "openshift"},

	// Image/multimodal
	{"pillow", true, "image-processing"},
	{"opencv", false, "image-processing"},

	// Tokenizer/NLP
	{"sentencepiece", true, "tokenizer"},
	{"tokenizers", true, "tokenizer"},
	{"tiktoken", true, "tokenizer"},

	// Testing
	{"pytest", false, "testing"},
	{"unittest", false, "testing"},
	{"tox", true, "testing"},

	// OpenAI/Anthropic SDKs
	{"openai", true, "llm-sdk"},
	{"anthropic", true, "llm-sdk"},

	// Numerical
	{"numpy", true, "numerical"},
	{"scipy", true, "numerical"},
	{"pandas", true, "data"},
}

// extractPythonDeps scans a repository for Python dependency files and extracts
// package dependencies with version constraints.
func extractPythonDeps(repoPath string) []PythonPackage {
	seen := make(map[string]PythonPackage) // normalized name -> package

	// 1. Parse requirements/*.txt and requirements*.txt at root
	reqFiles := findPythonRequirementFiles(repoPath)
	for _, fpath := range reqFiles {
		deps := parseRequirementsTxt(fpath, repoPath)
		for _, dep := range deps {
			mergePackage(seen, dep)
		}
	}

	// 2. Parse pyproject.toml
	pyprojectPath := filepath.Join(repoPath, "pyproject.toml")
	if data, err := os.ReadFile(pyprojectPath); err == nil {
		deps := parsePyprojectToml(string(data), "pyproject.toml")
		for _, dep := range deps {
			mergePackage(seen, dep)
		}
	}

	// 3. Parse setup.py / setup.cfg install_requires
	setupPyPath := filepath.Join(repoPath, "setup.py")
	if data, err := os.ReadFile(setupPyPath); err == nil {
		deps := parseSetupPy(string(data), "setup.py")
		for _, dep := range deps {
			mergePackage(seen, dep)
		}
	}

	setupCfgPath := filepath.Join(repoPath, "setup.cfg")
	if data, err := os.ReadFile(setupCfgPath); err == nil {
		deps := parseSetupCfg(string(data), "setup.cfg")
		for _, dep := range deps {
			mergePackage(seen, dep)
		}
	}

	// Build result list, categorize, sort
	result := make([]PythonPackage, 0, len(seen))
	for _, pkg := range seen {
		categorizePythonPackage(&pkg)
		result = append(result, pkg)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result
}

// findPythonRequirementFiles finds all requirements*.txt files in a repo.
func findPythonRequirementFiles(repoPath string) []string {
	var files []string

	// Root-level requirements*.txt
	matches, _ := filepath.Glob(filepath.Join(repoPath, "requirements*.txt"))
	files = append(files, matches...)

	// requirements/ directory
	reqDir := filepath.Join(repoPath, "requirements")
	if info, err := os.Stat(reqDir); err == nil && info.IsDir() {
		entries, _ := os.ReadDir(reqDir)
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".txt") {
				files = append(files, filepath.Join(reqDir, e.Name()))
			}
		}
	}

	return files
}

// parseRequirementsTxt parses a pip requirements.txt file.
func parseRequirementsTxt(fpath, repoPath string) []PythonPackage {
	data, err := os.ReadFile(fpath)
	if err != nil {
		return nil
	}

	relPath := relativePath(repoPath, fpath)
	var packages []PythonPackage

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		// Skip comments, empty lines, flags, and constraint references
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "-") {
			continue
		}

		// Strip inline comments
		if idx := strings.Index(line, " #"); idx >= 0 {
			line = strings.TrimSpace(line[:idx])
		}

		// Strip environment markers (e.g., ; platform_machine == "x86_64")
		if idx := strings.Index(line, ";"); idx >= 0 {
			line = strings.TrimSpace(line[:idx])
		}

		if match := requirementLineRE.FindStringSubmatch(line); match != nil {
			name := normalizePackageName(match[1])
			version := strings.TrimSpace(match[2])
			packages = append(packages, PythonPackage{
				Name:     name,
				Version:  version,
				Source:   relPath,
				Required: !isTestRequirement(relPath),
			})
		}
	}

	return packages
}

// parsePyprojectToml extracts dependencies from pyproject.toml.
func parsePyprojectToml(content, source string) []PythonPackage {
	var packages []PythonPackage
	inDeps := false
	isPoetrySection := false

	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)

		// Detect [project] dependencies or [tool.poetry.dependencies]
		if trimmed == "dependencies = [" ||
			strings.HasPrefix(trimmed, "dependencies = [") {
			inDeps = true
			isPoetrySection = false
			continue
		}
		if trimmed == "[tool.poetry.dependencies]" {
			inDeps = true
			isPoetrySection = true
			continue
		}

		// End of array (PEP 621 style)
		if inDeps && !isPoetrySection && trimmed == "]" {
			inDeps = false
			continue
		}

		// New section header ends deps
		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			inDeps = false
			isPoetrySection = false
			continue
		}

		if inDeps {
			// PEP 621 array format: "package>=version"
			if match := pyprojectDepRE.FindStringSubmatch(line); match != nil {
				name := normalizePackageName(match[1])
				version := strings.TrimSpace(match[2])
				packages = append(packages, PythonPackage{
					Name:     name,
					Version:  version,
					Source:   source,
					Required: true,
				})
				continue
			}
			// Poetry key-value format: package = "^version"
			if isPoetrySection {
				if match := poetryDepRE.FindStringSubmatch(trimmed); match != nil {
					name := normalizePackageName(match[1])
					if name == "python" {
						continue // skip python version constraint
					}
					version := match[2]
					packages = append(packages, PythonPackage{
						Name:     name,
						Version:  version,
						Source:   source,
						Required: true,
					})
				}
			}
		}
	}

	return packages
}

// setupPyInstallRequiresRE matches install_requires=[...] blocks.
var setupPyInstallRequiresRE = regexp.MustCompile(`install_requires\s*=\s*\[`)

// parseSetupPy extracts install_requires from setup.py.
func parseSetupPy(content, source string) []PythonPackage {
	var packages []PythonPackage

	// Find install_requires = [ ... ] block
	if !setupPyInstallRequiresRE.MatchString(content) {
		return nil
	}

	inBlock := false
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)

		if !inBlock {
			if setupPyInstallRequiresRE.MatchString(trimmed) {
				inBlock = true
				continue
			}
			continue
		}

		if trimmed == "]" || trimmed == "]," {
			break
		}

		// Parse quoted dependency strings
		dep := extractQuotedDep(trimmed)
		if dep != "" {
			if match := requirementLineRE.FindStringSubmatch(dep); match != nil {
				name := normalizePackageName(match[1])
				version := strings.TrimSpace(match[2])
				packages = append(packages, PythonPackage{
					Name:     name,
					Version:  version,
					Source:   source,
					Required: true,
				})
			}
		}
	}

	return packages
}

// parseSetupCfg extracts install_requires from setup.cfg.
func parseSetupCfg(content, source string) []PythonPackage {
	var packages []PythonPackage
	inInstallRequires := false

	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)

		if trimmed == "install_requires =" || strings.HasPrefix(trimmed, "install_requires=") {
			inInstallRequires = true
			// Handle inline value
			if strings.Contains(trimmed, "=") {
				parts := strings.SplitN(trimmed, "=", 2)
				val := strings.TrimSpace(parts[1])
				if val != "" {
					if match := requirementLineRE.FindStringSubmatch(val); match != nil {
						packages = append(packages, PythonPackage{
							Name:     normalizePackageName(match[1]),
							Version:  strings.TrimSpace(match[2]),
							Source:   source,
							Required: true,
						})
					}
				}
			}
			continue
		}

		if inInstallRequires {
			// Lines starting without whitespace are new sections
			if !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") && trimmed != "" {
				break
			}
			if trimmed == "" {
				continue
			}
			if match := requirementLineRE.FindStringSubmatch(trimmed); match != nil {
				packages = append(packages, PythonPackage{
					Name:     normalizePackageName(match[1]),
					Version:  strings.TrimSpace(match[2]),
					Source:   source,
					Required: true,
				})
			}
		}
	}

	return packages
}

// normalizePackageName normalizes a Python package name (PEP 503).
func normalizePackageName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "_", "-")
	name = strings.ReplaceAll(name, ".", "-")
	return name
}

// extractQuotedDep extracts a dependency string from a quoted line like '    "package>=1.0",'
func extractQuotedDep(line string) string {
	line = strings.TrimSpace(line)
	for _, q := range []string{`"`, `'`} {
		if strings.HasPrefix(line, q) {
			end := strings.LastIndex(line[1:], q)
			if end >= 0 {
				return line[1 : end+1]
			}
		}
	}
	return ""
}

// isTestRequirement checks if a requirements file is for testing/dev.
func isTestRequirement(path string) bool {
	lower := strings.ToLower(path)
	return strings.Contains(lower, "test") || strings.Contains(lower, "dev") ||
		strings.Contains(lower, "lint") || strings.Contains(lower, "doc")
}

// categorizePythonPackage assigns a category based on package name.
func categorizePythonPackage(pkg *PythonPackage) {
	nameLower := strings.ToLower(pkg.Name)
	for _, rule := range pythonPackageCategoryRules {
		if rule.exact {
			if nameLower == rule.name {
				pkg.Category = rule.category
				return
			}
		} else {
			if strings.Contains(nameLower, rule.name) {
				pkg.Category = rule.category
				return
			}
		}
	}
}

// mergePackage adds a package to the seen map, preferring entries with version info.
func mergePackage(seen map[string]PythonPackage, pkg PythonPackage) {
	existing, ok := seen[pkg.Name]
	if !ok {
		seen[pkg.Name] = pkg
		return
	}
	// Prefer the entry with version info, or required over not required
	if existing.Version == "" && pkg.Version != "" {
		seen[pkg.Name] = pkg
	} else if !existing.Required && pkg.Required {
		seen[pkg.Name] = pkg
	}
}
