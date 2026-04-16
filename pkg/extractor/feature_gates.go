package extractor

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Feature gate registration patterns in Go source.
// Matches controller-runtime, k8s.io/component-base, and common custom patterns.
var featureGateRegistrationPatterns = []*regexp.Regexp{
	// utilfeature.DefaultMutableFeatureGate.Add(map[featuregate.Feature]featuregate.FeatureSpec{
	regexp.MustCompile(`(?:DefaultMutableFeatureGate|MutableFeatureGate)\.Add\(`),
	// featuregate.NewFeatureGate() / featuregate.New()
	regexp.MustCompile(`featuregate\.(?:NewFeatureGate|New)\(`),
}

// featureSpecRE matches a feature gate definition inside a map literal:
//   FeatureName: {Default: true, PreRelease: featuregate.Beta},
// Captures: name, default value, pre-release stage
var featureSpecRE = regexp.MustCompile(
	`"?(\w+)"?\s*:\s*\{[^}]*Default\s*:\s*(true|false)` +
		`(?:[^}]*PreRelease\s*:\s*(?:featuregate\.)?(\w+))?`,
)

// featureConstRE matches const declarations that define feature gate names:
//   const MyFeature featuregate.Feature = "MyFeature"
//   MyFeature featuregate.Feature = "MyFeature"
var featureConstRE = regexp.MustCompile(
	`(\w+)\s+featuregate\.Feature\s*=\s*"([^"]+)"`,
)

// setFeatureGateRE matches runtime feature gate enable/disable calls:
//   utilfeature.DefaultMutableFeatureGate.Set("FeatureName=true")
//   featureGate.SetFromMap(...)
var setFeatureGateRE = regexp.MustCompile(
	`(?:DefaultMutableFeatureGate|MutableFeatureGate|featureGate)\.(?:Set|SetFromMap|Overrides)\(`,
)

// setArgRE extracts the gate name and value from Set("GateName=true") calls.
var setArgRE = regexp.MustCompile(`\.Set\(\s*"(\w+)=(true|false)"`)


// extractFeatureGates scans Go source for feature gate definitions and registrations.
func extractFeatureGates(repoPath string) []FeatureGate {
	var gates []FeatureGate
	seen := make(map[string]bool)

	goFiles := findFiles(repoPath, []string{"**/*.go"})

	// First pass: collect const definitions that map names to string values
	constMap := make(map[string]string) // Go identifier -> string value
	for _, fpath := range goFiles {
		if strings.Contains(fpath, "_test.go") || strings.Contains(fpath, "/vendor/") {
			continue
		}
		data, err := readFileNoSymlink(fpath)
		if err != nil {
			continue
		}
		for _, m := range featureConstRE.FindAllStringSubmatch(string(data), -1) {
			constMap[m[1]] = m[2]
		}
	}

	// Second pass: find gate registrations and extract specs
	for _, fpath := range goFiles {
		if strings.Contains(fpath, "_test.go") || strings.Contains(fpath, "/vendor/") {
			continue
		}
		data, err := readFileNoSymlink(fpath)
		if err != nil {
			log.Printf("warning: skipping %s: %v", fpath, err)
			continue
		}

		content := string(data)
		lines := strings.Split(content, "\n")
		source := relativePath(repoPath, fpath)

		// Check for registration calls
		hasRegistration := false
		for _, pat := range featureGateRegistrationPatterns {
			if pat.MatchString(content) {
				hasRegistration = true
				break
			}
		}

		if hasRegistration {
			// Extract feature specs from the map literal
			for _, m := range featureSpecRE.FindAllStringSubmatch(content, -1) {
				gateName := m[1]
				defaultVal := m[2] == "true"
				preRelease := ""
				if len(m) > 3 && m[3] != "" {
					preRelease = m[3]
				}

				// Resolve const reference to string value if available
				resolvedName := gateName
				if v, ok := constMap[gateName]; ok {
					resolvedName = v
				}

				if seen[resolvedName] {
					continue
				}
				seen[resolvedName] = true

				// Find the line number
				lineNum := findPatternLine(lines, gateName+":")

				gates = append(gates, FeatureGate{
					Name:       resolvedName,
					Default:    defaultVal,
					PreRelease: preRelease,
					Source:     source + ":" + strconv.Itoa(lineNum),
					LockToDefault: false,
				})
			}
		}

		// Also detect SetFromMap/Set calls as runtime overrides (informational)
		if setFeatureGateRE.MatchString(content) {
			for lineNum, line := range lines {
				stripped := strings.TrimSpace(line)
				if strings.HasPrefix(stripped, "//") {
					continue
				}
				if setFeatureGateRE.MatchString(stripped) {
					if sm := setArgRE.FindStringSubmatch(stripped); sm != nil {
						gateName := sm[1]
						if seen[gateName] {
							continue
						}
						seen[gateName] = true
						gates = append(gates, FeatureGate{
							Name:       gateName,
							Default:    sm[2] == "true",
							Source:     source + ":" + strconv.Itoa(lineNum+1),
							RuntimeSet: true,
						})
					}
				}
			}
		}
	}

	if gates == nil {
		gates = []FeatureGate{}
	}
	return gates
}

// readFileNoSymlink reads a file after verifying it is not a symlink and not
// oversized. Uses os.Lstat (does not follow symlinks) for the check, then
// opens and reads via the file descriptor to minimize the TOCTOU window.
func readFileNoSymlink(fpath string) ([]byte, error) {
	info, err := os.Lstat(fpath)
	if err != nil {
		return nil, err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("skipping symlink: %s", fpath)
	}
	if info.Size() > maxFileSize {
		return nil, fmt.Errorf("skipping oversized file: %s", fpath)
	}
	f, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

// findPatternLine returns the 1-based line number of the first occurrence of pattern.
func findPatternLine(lines []string, pattern string) int {
	for i, line := range lines {
		if strings.Contains(line, pattern) {
			return i + 1
		}
	}
	return 0
}
