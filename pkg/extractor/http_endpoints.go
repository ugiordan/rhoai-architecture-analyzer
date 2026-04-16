package extractor

import (
	"encoding/json"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Route registration patterns for common Go HTTP frameworks.
var httpRoutePatterns = []*regexp.Regexp{
	// net/http: http.HandleFunc("/path", handler)
	regexp.MustCompile(`(?:http|mux)\.HandleFunc\(\s*"([^"]+)"`),
	// net/http: http.Handle("/path", handler)
	regexp.MustCompile(`(?:http|mux)\.Handle\(\s*"([^"]+)"`),
	// gorilla/mux: r.HandleFunc("/path", handler).Methods("GET", "POST")
	regexp.MustCompile(`\.HandleFunc\(\s*"([^"]+)"[^)]*\)\.Methods\(([^)]+)\)`),
	// gin/echo/chi: r.GET("/path", handler)
	regexp.MustCompile(`\.\s*(GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS)\(\s*"([^"]+)"`),
	// chi: r.Route("/path", func(...) { ... })
	regexp.MustCompile(`\.Route\(\s*"([^"]+)"`),
	// chi: r.Mount("/path", handler)
	regexp.MustCompile(`\.Mount\(\s*"([^"]+)"`),
	// Group: r.Group("/path")
	regexp.MustCompile(`\.Group\(\s*"([^"]+)"`),
}

var openAPIPatterns = []string{
	"swagger.yaml",
	"swagger.yml",
	"swagger.json",
	"openapi.yaml",
	"openapi.yml",
	"openapi.json",
	"api/swagger.yaml",
	"api/swagger.json",
	"api/openapi.yaml",
	"api/openapi.json",
	"docs/swagger.yaml",
	"docs/openapi.yaml",
}

// extractHTTPEndpoints scans Go source for route registrations and OpenAPI specs.
func extractHTTPEndpoints(repoPath string) []HTTPEndpoint {
	var endpoints []HTTPEndpoint

	// Scan Go files for route registrations
	goFiles := findFiles(repoPath, []string{
		"**/*.go",
	})

	for _, fpath := range goFiles {
		// Skip test files and vendor
		if strings.Contains(fpath, "_test.go") || strings.Contains(fpath, "/vendor/") {
			continue
		}

		data, err := os.ReadFile(fpath)
		if err != nil {
			log.Printf("warning: skipping %s: %v", fpath, err)
			continue
		}

		content := string(data)
		lines := strings.Split(content, "\n")

		for lineNum, line := range lines {
			stripped := strings.TrimSpace(line)
			if strings.HasPrefix(stripped, "//") {
				continue
			}

			for _, pattern := range httpRoutePatterns {
				matches := pattern.FindStringSubmatch(stripped)
				if matches == nil {
					continue
				}

				ep := HTTPEndpoint{
					Source: relativePath(repoPath, fpath) + ":" + strconv.Itoa(lineNum+1),
				}

				switch {
				case len(matches) == 3 && isHTTPMethod(matches[1]):
					// Patterns like .GET("/path", ...)
					ep.Method = matches[1]
					ep.Path = matches[2]
					endpoints = append(endpoints, ep)
				case len(matches) == 3 && !isHTTPMethod(matches[1]):
					// gorilla/mux: HandleFunc("/path", ...).Methods("GET", "POST")
					// matches[2] contains the full argument list, extract quoted methods
					ep.Path = matches[1]
					methods := extractQuotedStrings(matches[2])
					if len(methods) == 0 {
						ep.Method = "*"
						endpoints = append(endpoints, ep)
					} else {
						for _, m := range methods {
							mep := ep
							mep.Method = m
							endpoints = append(endpoints, mep)
						}
					}
				default:
					ep.Path = matches[1]
					ep.Method = "*"
					endpoints = append(endpoints, ep)
				}
			}
		}
	}

	// Parse OpenAPI/Swagger specs
	specFiles := findFiles(repoPath, openAPIPatterns)
	for _, fpath := range specFiles {
		data, err := os.ReadFile(fpath)
		if err != nil {
			continue
		}

		endpoints = append(endpoints, parseOpenAPIEndpoints(fpath, data, repoPath)...)
	}

	if endpoints == nil {
		endpoints = []HTTPEndpoint{}
	}
	return endpoints
}

func parseOpenAPIEndpoints(fpath string, data []byte, repoPath string) []HTTPEndpoint {
	var endpoints []HTTPEndpoint
	source := relativePath(repoPath, fpath)

	// Try JSON first
	var spec map[string]interface{}
	if err := json.Unmarshal(data, &spec); err != nil {
		// Try YAML
		docs := parseYAMLSafe(fpath)
		if len(docs) > 0 {
			spec = docs[0]
		}
	}
	if spec == nil {
		return nil
	}

	paths, ok := spec["paths"].(map[string]interface{})
	if !ok {
		return nil
	}

	for path, methods := range paths {
		methodMap, ok := methods.(map[string]interface{})
		if !ok {
			continue
		}
		for method := range methodMap {
			method = strings.ToUpper(method)
			if isHTTPMethod(method) {
				endpoints = append(endpoints, HTTPEndpoint{
					Method: method,
					Path:   path,
					Source: source,
				})
			}
		}
	}

	return endpoints
}

func isHTTPMethod(s string) bool {
	switch strings.ToUpper(s) {
	case "GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS":
		return true
	}
	return false
}

var quotedStringRE = regexp.MustCompile(`"([^"]+)"`)

// extractQuotedStrings pulls all double-quoted strings from a text fragment.
func extractQuotedStrings(s string) []string {
	matches := quotedStringRE.FindAllStringSubmatch(s, -1)
	out := make([]string, 0, len(matches))
	for _, m := range matches {
		out = append(out, m[1])
	}
	return out
}
