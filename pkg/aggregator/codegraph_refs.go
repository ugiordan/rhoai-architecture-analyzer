package aggregator

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// codeGraphRef represents a cross-component reference found in code-graph file paths.
type codeGraphRef struct {
	From     string // source component
	To       string // referenced component
	Evidence string // file path containing the reference
}

// detectCodeGraphRefs scans code-graph.json files for file paths that mention
// other component names as path segments. This catches implicit relationships
// like llama-stack having a providers/remote/inference/vllm/ directory.
//
// resultsDir is the root directory containing per-component subdirectories.
// jsonPaths are the paths to component-architecture.json files (used to locate
// sibling code-graph.json files). componentNames is the set of all known
// component names (normalized to lowercase with dashes).
func detectCodeGraphRefs(jsonPaths []string, componentNames []string) []codeGraphRef {
	// Build lookup: normalized name -> canonical name
	// Also build a set of name segments to match against path segments
	nameToCanonical := make(map[string]string)
	for _, name := range componentNames {
		norm := strings.ToLower(strings.ReplaceAll(name, "_", "-"))
		nameToCanonical[norm] = name
		// Also add without dashes for fuzzy matching (e.g., "vllm" in path "vllm/")
		noDash := strings.ReplaceAll(norm, "-", "")
		if noDash != norm {
			nameToCanonical[noDash] = name
		}
	}

	var refs []codeGraphRef
	seen := make(map[string]bool) // "from|to" dedup

	for _, jp := range jsonPaths {
		dir := filepath.Dir(jp)
		cgPath := filepath.Join(dir, "code-graph.json")

		fi, err := os.Stat(cgPath)
		if err != nil {
			continue // no code-graph.json for this component
		}
		// Skip very large files (>500MB) to avoid memory pressure
		if fi.Size() > 500*1024*1024 {
			log.Printf("WARN: skipping oversized code-graph.json %s (%d bytes)", cgPath, fi.Size())
			continue
		}

		// Determine component name from the architecture JSON
		compName := componentNameFromArchJSON(jp)
		if compName == "" {
			continue
		}

		filePaths, err := extractCodeGraphFilePaths(cgPath)
		if err != nil {
			log.Printf("WARN: failed to extract file paths from %s: %v", cgPath, err)
			continue
		}

		// Check each file path for references to other components
		for _, fp := range filePaths {
			segments := pathSegments(fp)
			for _, seg := range segments {
				normSeg := strings.ToLower(strings.ReplaceAll(seg, "_", "-"))
				if canonical, ok := nameToCanonical[normSeg]; ok && canonical != compName {
					key := compName + "|" + canonical
					if !seen[key] {
						seen[key] = true
						refs = append(refs, codeGraphRef{
							From:     compName,
							To:       canonical,
							Evidence: fp,
						})
					}
				}
				// Also try without dashes
				noDash := strings.ReplaceAll(normSeg, "-", "")
				if len(noDash) >= 4 { // skip very short segments
					if canonical, ok := nameToCanonical[noDash]; ok && canonical != compName {
						key := compName + "|" + canonical
						if !seen[key] {
							seen[key] = true
							refs = append(refs, codeGraphRef{
								From:     compName,
								To:       canonical,
								Evidence: fp,
							})
						}
					}
				}
			}
		}
	}

	return refs
}

// componentNameFromArchJSON reads the "component" field from a component-architecture.json.
func componentNameFromArchJSON(path string) string {
	raw, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	var data map[string]interface{}
	if err := json.Unmarshal(raw, &data); err != nil {
		return ""
	}
	if name, ok := data["component"].(string); ok {
		return name
	}
	return ""
}

// extractCodeGraphFilePaths extracts unique file paths from code-graph.json nodes.
// Uses streaming JSON decoding for memory efficiency.
func extractCodeGraphFilePaths(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dec := json.NewDecoder(f)

	// Navigate to the "nodes" array
	// Structure: {"nodes": [...], "edges": [...], ...}
	t, err := dec.Token() // opening {
	if err != nil {
		return nil, fmt.Errorf("reading opening brace: %w", err)
	}
	if delim, ok := t.(json.Delim); !ok || delim != '{' {
		return nil, fmt.Errorf("expected opening brace, got %v", t)
	}

	seen := make(map[string]bool)
	var files []string

	for dec.More() {
		t, err := dec.Token() // field name
		if err != nil {
			return nil, fmt.Errorf("reading field name: %w", err)
		}
		key, ok := t.(string)
		if !ok {
			continue
		}

		if key == "nodes" {
			// Parse the nodes array, extracting only file fields
			t, err := dec.Token() // opening [
			if err != nil {
				return nil, fmt.Errorf("reading nodes array: %w", err)
			}
			if delim, ok := t.(json.Delim); !ok || delim != '[' {
				return nil, fmt.Errorf("expected nodes array, got %v", t)
			}

			for dec.More() {
				var node struct {
					File string `json:"file"`
				}
				if err := dec.Decode(&node); err != nil {
					return nil, fmt.Errorf("decoding node: %w", err)
				}
				if node.File != "" && !seen[node.File] {
					seen[node.File] = true
					files = append(files, node.File)
				}
			}

			// Read closing ]
			if _, err := dec.Token(); err != nil {
				return nil, fmt.Errorf("reading nodes closing bracket: %w", err)
			}
		} else {
			// Skip other fields by consuming the value
			if err := skipJSONValue(dec); err != nil {
				return nil, fmt.Errorf("skipping field %s: %w", key, err)
			}
		}
	}

	return files, nil
}

// skipJSONValue consumes one JSON value from the decoder.
func skipJSONValue(dec *json.Decoder) error {
	t, err := dec.Token()
	if err != nil {
		return err
	}
	delim, ok := t.(json.Delim)
	if !ok {
		return nil // primitive value, already consumed
	}
	// It's an array or object; consume until the matching close
	for dec.More() {
		if delim == '{' {
			// Skip key
			if _, err := dec.Token(); err != nil {
				return err
			}
		}
		if err := skipJSONValue(dec); err != nil {
			return err
		}
	}
	// Consume closing delimiter
	_, err = dec.Token()
	return err
}

// pathSegments splits a file path into directory and filename segments.
func pathSegments(path string) []string {
	// Normalize separators
	path = strings.ReplaceAll(path, "\\", "/")
	parts := strings.Split(path, "/")
	// Also split on dots for package names like "evalhub.adapter"
	var segments []string
	for _, p := range parts {
		if p == "" {
			continue
		}
		segments = append(segments, p)
		// If the segment contains dots (not a file extension), split on dots too
		if strings.Contains(p, ".") {
			dotParts := strings.Split(p, ".")
			for _, dp := range dotParts {
				if dp != "" && dp != p {
					segments = append(segments, dp)
				}
			}
		}
	}
	return segments
}
