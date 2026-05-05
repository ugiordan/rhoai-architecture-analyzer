package extractor

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var dockerfilePatterns = []string{
	"Dockerfile",
	"Dockerfile.*",
	"Containerfile",
	"Containerfile.*",
	"**/Dockerfile",
	"**/Dockerfile.*",
	"**/Containerfile",
	"**/Containerfile.*",
}

var (
	fromRE       = regexp.MustCompile(`^FROM\s+(?:--platform=(\S+)\s+)?(\S+)(?:\s+[Aa][Ss]\s+\S+)?`)
	userRE       = regexp.MustCompile(`^USER\s+(\S+)`)
	exposeRE     = regexp.MustCompile(`^EXPOSE\s+(.*)`)
	argRE        = regexp.MustCompile(`^ARG\s+(\w+)(?:=(.*))?`)
	envRE        = regexp.MustCompile(`^ENV\s+(\w+)(?:=|\s+)(.*)`)
	targetArchRE = regexp.MustCompile(`TARGETARCH|TARGETOS|TARGETPLATFORM`)
)

// FIPS-related environment variables and build args.
var fipsIndicators = []string{
	"GOEXPERIMENT=boringcrypto",
	"FIPS_MODE",
	"OPENSSL_FIPS",
	"FIPS_ENABLED",
}

// extractDockerfiles scans Dockerfiles/Containerfiles for metadata including
// base images, user directives, exposed ports, architectures, and FIPS flags.
func extractDockerfiles(repoPath string) []DockerfileInfo {
	files := findFiles(repoPath, dockerfilePatterns)
	var dockerfiles []DockerfileInfo

	for _, fpath := range files {
		data, err := os.ReadFile(fpath)
		if err != nil {
			log.Printf("warning: skipping %s: %v", fpath, err)
			continue
		}

		lines := joinContinuationLines(strings.Split(string(data), "\n"))
		var fromImages []string
		user := ""
		var exposedPorts []int
		var issues []string
		archSet := make(map[string]bool)
		fipsEnabled := false
		buildArgs := make(map[string]string)
		content := string(data)

		for _, line := range lines {
			stripped := strings.TrimSpace(line)
			if stripped == "" || strings.HasPrefix(stripped, "#") {
				continue
			}

			// FROM instruction with optional --platform
			if match := fromRE.FindStringSubmatch(stripped); match != nil {
				platform := match[1]
				image := match[2]
				fromImages = append(fromImages, image)

				if platform != "" {
					// --platform=linux/amd64 or $TARGETPLATFORM
					if strings.Contains(platform, "$") {
						// Variable reference, will check for TARGETARCH usage
					} else {
						parts := strings.Split(platform, "/")
						if len(parts) >= 2 {
							archSet[parts[1]] = true
						}
					}
				}

				if strings.HasSuffix(image, ":latest") ||
					(!strings.Contains(image, ":") && !strings.Contains(image, "@")) {
					issues = append(issues, fmt.Sprintf("Unpinned base image: %s", image))
				}
			}

			// USER instruction
			if match := userRE.FindStringSubmatch(stripped); match != nil {
				user = match[1]
			}

			// EXPOSE instruction
			if match := exposeRE.FindStringSubmatch(stripped); match != nil {
				for _, part := range strings.Fields(match[1]) {
					portStr := strings.SplitN(part, "/", 2)[0]
					if port, err := strconv.Atoi(portStr); err == nil {
						exposedPorts = append(exposedPorts, port)
					}
				}
			}

			// ARG instruction
			if match := argRE.FindStringSubmatch(stripped); match != nil {
				argName := match[1]
				argVal := ""
				if len(match) > 2 {
					argVal = strings.TrimSpace(match[2])
				}
				// Track security-relevant build args
				if isSecurityRelevantArg(argName) {
					buildArgs[argName] = argVal
				}
			}

			// ENV instruction
			if match := envRE.FindStringSubmatch(stripped); match != nil {
				envName := match[1]
				envVal := strings.TrimSpace(match[2])
				for _, indicator := range fipsIndicators {
					if strings.Contains(envName+"="+envVal, indicator) {
						fipsEnabled = true
					}
				}
			}
		}

		// Check for TARGETARCH usage (multi-arch build)
		if targetArchRE.MatchString(content) {
			archSet["multi-arch"] = true
		}

		// Check for root user
		if user == "root" || user == "0" {
			issues = append(issues, "Container runs as root user")
		} else if user == "" {
			issues = append(issues, "No USER directive found (defaults to root)")
		}

		// Check for FIPS in build args
		for argName, argVal := range buildArgs {
			for _, indicator := range fipsIndicators {
				if strings.Contains(argName+"="+argVal, indicator) {
					fipsEnabled = true
				}
			}
		}

		stages := len(fromImages)
		baseImage := ""
		var buildStageImages []string
		if len(fromImages) > 0 {
			baseImage = fromImages[len(fromImages)-1] // runtime image is the last FROM stage
			// Capture all non-final stage images (build stages)
			if len(fromImages) > 1 {
				buildStageImages = fromImages[:len(fromImages)-1]
			}
		}

		if exposedPorts == nil {
			exposedPorts = []int{}
		}
		if issues == nil {
			issues = []string{}
		}

		var architectures []string
		for arch := range archSet {
			architectures = append(architectures, arch)
		}
		sort.Strings(architectures)
		if architectures == nil {
			architectures = []string{}
		}

		if len(buildArgs) == 0 {
			buildArgs = nil
		}

		dockerfiles = append(dockerfiles, DockerfileInfo{
			Path:             relativePath(repoPath, fpath),
			BaseImage:        baseImage,
			BuildStageImages: buildStageImages,
			Stages:           stages,
			User:             user,
			ExposedPorts:     exposedPorts,
			Issues:           issues,
			Architectures:    architectures,
			FIPSEnabled:      fipsEnabled,
			BuildArgs:        buildArgs,
		})
	}

	if dockerfiles == nil {
		dockerfiles = []DockerfileInfo{}
	}
	return dockerfiles
}

// joinContinuationLines merges Dockerfile lines joined with trailing backslashes
// into single logical lines.
func joinContinuationLines(lines []string) []string {
	var result []string
	var current strings.Builder
	for _, line := range lines {
		trimmed := strings.TrimRight(line, " \t")
		if strings.HasSuffix(trimmed, "\\") {
			current.WriteString(strings.TrimSuffix(trimmed, "\\"))
			current.WriteByte(' ')
		} else {
			current.WriteString(line)
			result = append(result, current.String())
			current.Reset()
		}
	}
	if current.Len() > 0 {
		result = append(result, current.String())
	}
	return result
}

func isSecurityRelevantArg(name string) bool {
	upper := strings.ToUpper(name)
	relevant := []string{
		"GOEXPERIMENT", "CGO_ENABLED", "FIPS", "TARGETARCH", "TARGETOS",
		"TARGETPLATFORM", "GO_VERSION", "GOLANG_VERSION", "USER_ID",
		"GROUP_ID", "OPENSSL",
	}
	for _, r := range relevant {
		if strings.Contains(upper, r) {
			return true
		}
	}
	return false
}
