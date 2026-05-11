package extractor

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Directories to skip when scanning Python files.
var pythonSkipDirs = map[string]bool{
	"venv":        true,
	".venv":       true,
	"env":         true,
	".env":        true,
	"__pycache__": true,
	".tox":        true,
	".eggs":       true,
	"node_modules": true,
	"vendor":      true,
	".git":        true,
}

// pythonPortPattern holds a compiled regex and metadata for a Python port pattern.
type pythonPortPattern struct {
	re          *regexp.Regexp
	serviceName string // static name, or empty to derive from match
	nameFromArg int    // capture group index for dynamic service name (0 = unused)
	portGroup   int    // capture group index containing the port number
}

// uvicorn.run(..., port=NNNN, ...)
// uvicorn.run("app:app", port=NNNN)
var uvicornRunRE = regexp.MustCompile(`uvicorn\.run\([^)]*\bport\s*=\s*(\d+)`)

// app.run(..., port=NNNN, ...) covers Flask and similar frameworks
var appRunPortRE = regexp.MustCompile(`\.run\([^)]*\bport\s*=\s*(\d+)`)

// gRPC: server.add_insecure_port('[::]:NNNN') or server.add_insecure_port('0.0.0.0:NNNN')
var grpcInsecurePortRE = regexp.MustCompile(`\.add_insecure_port\(\s*['"]([^'"]+)['"]\s*\)`)

// gRPC: server.add_secure_port('[::]:NNNN', creds)
var grpcSecurePortRE = regexp.MustCompile(`\.add_secure_port\(\s*['"]([^'"]+)['"]`)

// Gunicorn config: bind = "0.0.0.0:NNNN" or bind = ":NNNN"
var gunicornBindRE = regexp.MustCompile(`^\s*bind\s*=\s*['"]([^'"]+)['"]`)

// argparse: add_argument("--port", ... default=NNNN)
// Match both single-arg and multi-arg forms
var argparsePortRE = regexp.MustCompile(`add_argument\([^)]*--port[^)]*default\s*=\s*(\d+)`)

// os.environ.get("PORT", NNNN) or os.environ.get("PORT", "NNNN")
var environGetPortRE = regexp.MustCompile(`os\.environ\.get\(\s*['"]PORT['"]\s*,\s*['"]?(\d+)['"]?\s*\)`)

// os.getenv("PORT", NNNN) or os.getenv("PORT", "NNNN")
var getenvPortRE = regexp.MustCompile(`os\.getenv\(\s*['"]PORT['"]\s*,\s*['"]?(\d+)['"]?\s*\)`)

// extractPythonPorts scans Python files for listening port definitions and
// returns them as Services with type "python-source".
func extractPythonPorts(repoPath string) []Service {
	pyFiles := findPythonFiles(repoPath)
	if len(pyFiles) == 0 {
		return nil
	}

	// Deduplicate: same port + same service name = single entry.
	type portKey struct {
		name string
		port int
	}
	seen := make(map[portKey]bool)
	var services []Service

	for _, fpath := range pyFiles {
		info, err := os.Lstat(fpath)
		if err != nil || info.Mode()&os.ModeSymlink != 0 {
			continue
		}
		if info.Size() > maxFileSize {
			log.Printf("skipping oversized Python file %s: %d bytes", fpath, info.Size())
			continue
		}

		data, err := os.ReadFile(fpath)
		if err != nil {
			log.Printf("warning: skipping Python file %s: %v", fpath, err)
			continue
		}

		content := string(data)
		lines := strings.Split(content, "\n")
		source := relativePath(repoPath, fpath)

		for lineNum, line := range lines {
			stripped := strings.TrimSpace(line)
			// Skip comments
			if strings.HasPrefix(stripped, "#") {
				continue
			}

			loc := source + ":" + strconv.Itoa(lineNum+1)

			// uvicorn.run(..., port=NNNN)
			if m := uvicornRunRE.FindStringSubmatch(stripped); m != nil {
				if port, err := strconv.Atoi(m[1]); err == nil {
					key := portKey{"uvicorn-server", port}
					if !seen[key] {
						seen[key] = true
						services = append(services, pythonService("uvicorn-server", port, loc))
					}
				}
			}

			// app.run(port=NNNN) but NOT uvicorn.run (already handled)
			if !strings.Contains(stripped, "uvicorn.run") {
				if m := appRunPortRE.FindStringSubmatch(stripped); m != nil {
					if port, err := strconv.Atoi(m[1]); err == nil {
						name := inferAppRunName(source)
						key := portKey{name, port}
						if !seen[key] {
							seen[key] = true
							services = append(services, pythonService(name, port, loc))
						}
					}
				}
			}

			// gRPC insecure port
			if m := grpcInsecurePortRE.FindStringSubmatch(stripped); m != nil {
				if port := extractPortFromAddr(m[1]); port > 0 {
					key := portKey{"grpc-server", port}
					if !seen[key] {
						seen[key] = true
						services = append(services, pythonService("grpc-server", port, loc))
					}
				}
			}

			// gRPC secure port
			if m := grpcSecurePortRE.FindStringSubmatch(stripped); m != nil {
				if port := extractPortFromAddr(m[1]); port > 0 {
					key := portKey{"grpc-server-tls", port}
					if !seen[key] {
						seen[key] = true
						services = append(services, pythonService("grpc-server-tls", port, loc))
					}
				}
			}

			// Gunicorn bind
			if m := gunicornBindRE.FindStringSubmatch(stripped); m != nil {
				if port := extractPortFromAddr(m[1]); port > 0 {
					key := portKey{"gunicorn-server", port}
					if !seen[key] {
						seen[key] = true
						services = append(services, pythonService("gunicorn-server", port, loc))
					}
				}
			}

			// argparse --port default
			if m := argparsePortRE.FindStringSubmatch(stripped); m != nil {
				if port, err := strconv.Atoi(m[1]); err == nil {
					key := portKey{"cli-port-default", port}
					if !seen[key] {
						seen[key] = true
						services = append(services, pythonService("cli-port-default", port, loc))
					}
				}
			}

			// os.environ.get("PORT", default)
			if m := environGetPortRE.FindStringSubmatch(stripped); m != nil {
				if port, err := strconv.Atoi(m[1]); err == nil {
					key := portKey{"env-port-default", port}
					if !seen[key] {
						seen[key] = true
						services = append(services, pythonService("env-port-default", port, loc))
					}
				}
			}

			// os.getenv("PORT", default)
			if m := getenvPortRE.FindStringSubmatch(stripped); m != nil {
				if port, err := strconv.Atoi(m[1]); err == nil {
					key := portKey{"env-port-default", port}
					if !seen[key] {
						seen[key] = true
						services = append(services, pythonService("env-port-default", port, loc))
					}
				}
			}
		}
	}

	return services
}

// pythonService builds a Service with type "python-source".
func pythonService(name string, port int, source string) Service {
	return Service{
		Name:   name,
		Type:   "python-source",
		Source: source,
		Ports: []ServicePort{
			{Port: port, Protocol: "TCP"},
		},
	}
}

// extractPortFromAddr extracts the port number from an address like "[::]:50051",
// "0.0.0.0:8080", or ":8000".
func extractPortFromAddr(addr string) int {
	// Handle unix sockets
	if strings.HasPrefix(addr, "unix:") {
		return 0
	}
	idx := strings.LastIndex(addr, ":")
	if idx < 0 {
		// Try parsing the whole string as a port number
		port, err := strconv.Atoi(strings.TrimSpace(addr))
		if err != nil {
			return 0
		}
		return port
	}
	portStr := addr[idx+1:]
	port, err := strconv.Atoi(strings.TrimSpace(portStr))
	if err != nil {
		return 0
	}
	return port
}

// inferAppRunName generates a service name from the Python file path.
// e.g. "src/server.py" -> "server", "app.py" -> "app-server"
func inferAppRunName(source string) string {
	base := filepath.Base(source)
	name := strings.TrimSuffix(base, ".py")
	if name == "" || name == "__main__" {
		// Use parent directory name
		dir := filepath.Dir(source)
		name = filepath.Base(dir)
	}
	if name == "app" || name == "main" || name == "server" || name == "run" {
		return name + "-server"
	}
	return name + "-server"
}

// findPythonFiles returns all .py files under repoPath, skipping venv,
// __pycache__, test directories, and other non-source directories.
func findPythonFiles(repoPath string) []string {
	var result []string
	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // skip errors
		}
		if info.IsDir() {
			base := filepath.Base(path)
			if pythonSkipDirs[base] {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".py") {
			return nil
		}
		// Skip test files
		name := info.Name()
		if strings.HasPrefix(name, "test_") || strings.HasSuffix(name, "_test.py") {
			return nil
		}
		// Skip files inside test directories
		rel, _ := filepath.Rel(repoPath, path)
		if strings.Contains(rel, "/tests/") || strings.Contains(rel, "/test/") {
			return nil
		}
		// Skip symlinks
		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}
		result = append(result, path)
		return nil
	})
	if err != nil {
		log.Printf("warning: walking Python files in %s: %v", repoPath, err)
	}
	return result
}
