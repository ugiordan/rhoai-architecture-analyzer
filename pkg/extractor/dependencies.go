package extractor

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

var singleRequireRE = regexp.MustCompile(`^require\s+(\S+)\s+(\S+)`)
var singleReplaceRE = regexp.MustCompile(`^replace\s+(\S+)\s+\S*\s*=>\s*(\S+)\s+(\S+)`)

// extractDependencies parses go.mod files and extracts module dependencies,
// identifying internal dependencies based on the provided module prefixes.
func extractDependencies(repoPath string, modulePrefixes []string) *DependencyData {
	goModFiles := findFiles(repoPath, []string{"go.mod", "**/go.mod"})

	var goModules []GoModule
	var internalODH []InternalODH
	var replaceDirectives []ReplaceDirective
	var goVersion, toolchain string

	for _, fpath := range goModFiles {
		data, err := os.ReadFile(fpath)
		if err != nil {
			log.Printf("warning: skipping %s: %v", fpath, err)
			continue
		}

		content := string(data)

		// First pass: collect replace directives
		replaces := parseReplaceDirectives(content)

		// Surface replace directives in output
		for orig, rep := range replaces {
			replaceDirectives = append(replaceDirectives, ReplaceDirective{
				Original:    orig,
				Replacement: rep.module,
				Version:     rep.version,
			})
		}

		// Second pass: collect require directives and metadata
		inRequire := false
		for _, line := range strings.Split(content, "\n") {
			stripped := strings.TrimSpace(line)

			// Skip comments (go.mod uses // for comments)
			if strings.HasPrefix(stripped, "//") {
				continue
			}

			// Extract go version directive (e.g., "go 1.21")
			if goVersion == "" && strings.HasPrefix(stripped, "go ") && !strings.HasPrefix(stripped, "go.") {
				fields := strings.Fields(stripped)
				if len(fields) == 2 {
					goVersion = fields[1]
				}
				continue
			}

			// Extract toolchain directive (e.g., "toolchain go1.21.0")
			if toolchain == "" && strings.HasPrefix(stripped, "toolchain ") {
				fields := strings.Fields(stripped)
				if len(fields) == 2 {
					toolchain = fields[1]
				}
				continue
			}

			if strings.HasPrefix(stripped, "require (") || strings.HasPrefix(stripped, "require(") {
				inRequire = true
				continue
			}
			if inRequire && stripped == ")" {
				inRequire = false
				continue
			}

			// Single-line require
			if match := singleRequireRE.FindStringSubmatch(stripped); match != nil {
				if strings.HasSuffix(stripped, "// indirect") {
					continue
				}
				module, version := match[1], match[2]
				// Apply replace directive if present
				if rep, ok := replaces[module]; ok {
					module = rep.module
					version = rep.version
				}
				goModules = append(goModules, GoModule{Module: module, Version: version})
				if comp, ok := matchInternalModule(module, modulePrefixes); ok {
					internalODH = append(internalODH, InternalODH{
						Component:   comp,
						Interaction: fmt.Sprintf("Go module dependency: %s", module),
					})
				}
				continue
			}

			// Lines inside require block
			if inRequire {
				if strings.HasSuffix(stripped, "// indirect") {
					continue
				}
				parts := strings.Fields(stripped)
				if len(parts) >= 2 {
					module, version := parts[0], parts[1]
					if rep, ok := replaces[module]; ok {
						module = rep.module
						version = rep.version
					}
					goModules = append(goModules, GoModule{Module: module, Version: version})
					if comp, ok := matchInternalModule(module, modulePrefixes); ok {
						internalODH = append(internalODH, InternalODH{
							Component:   comp,
							Interaction: fmt.Sprintf("Go module dependency: %s", module),
						})
					}
				}
			}
		}
	}

	if goModules == nil {
		goModules = []GoModule{}
	}
	if internalODH == nil {
		internalODH = []InternalODH{}
	}

	return &DependencyData{
		GoVersion:         goVersion,
		Toolchain:         toolchain,
		GoModules:         goModules,
		ReplaceDirectives: replaceDirectives,
		InternalODH:       internalODH,
	}
}

type replaceTarget struct {
	module  string
	version string
}

// parseReplaceDirectives extracts replace directives from go.mod content.
// Handles both single-line (replace A => B v1) and block (replace (...)) forms.
func parseReplaceDirectives(content string) map[string]replaceTarget {
	replaces := make(map[string]replaceTarget)
	inReplace := false

	for _, line := range strings.Split(content, "\n") {
		stripped := strings.TrimSpace(line)

		if strings.HasPrefix(stripped, "replace (") || strings.HasPrefix(stripped, "replace(") {
			inReplace = true
			continue
		}
		if inReplace && stripped == ")" {
			inReplace = false
			continue
		}

		// Single-line replace
		if match := singleReplaceRE.FindStringSubmatch(stripped); match != nil {
			original, replacement, version := match[1], match[2], match[3]
			// Skip local path replacements (no version, starts with . or /)
			if !strings.HasPrefix(replacement, ".") && !strings.HasPrefix(replacement, "/") {
				replaces[original] = replaceTarget{module: replacement, version: version}
			}
			continue
		}

		// Lines inside replace block: A v1 => B v2
		if inReplace {
			parts := strings.Fields(stripped)
			// Format: module [version] => replacement version
			arrowIdx := -1
			for i, p := range parts {
				if p == "=>" {
					arrowIdx = i
					break
				}
			}
			if arrowIdx >= 1 && arrowIdx+2 <= len(parts) {
				original := parts[0]
				replacement := parts[arrowIdx+1]
				version := ""
				if arrowIdx+2 < len(parts) {
					version = parts[arrowIdx+2]
				}
				if !strings.HasPrefix(replacement, ".") && !strings.HasPrefix(replacement, "/") {
					replaces[original] = replaceTarget{module: replacement, version: version}
				}
			}
		}
	}

	return replaces
}

// matchInternalModule checks if a module matches any of the internal prefixes.
// Returns the component name and true if matched.
func matchInternalModule(module string, prefixes []string) (string, bool) {
	for _, prefix := range prefixes {
		if strings.HasPrefix(module, prefix) {
			component := strings.SplitN(module[len(prefix):], "/", 2)[0]
			return component, true
		}
	}
	return "", false
}
