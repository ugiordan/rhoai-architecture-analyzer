package extractor

import (
	"log"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Database connection patterns: scheme-based URIs and dial functions.
var dbConnectionPatterns = []struct {
	re      *regexp.Regexp
	service string
	connType string
}{
	// PostgreSQL
	{regexp.MustCompile(`"(postgres(?:ql)?://[^"]+)"`), "postgres", "database"},
	{regexp.MustCompile(`(?:sql\.Open|sqlx\.(?:Open|Connect|MustConnect))\(\s*"postgres(?:ql)?"\s*,`), "postgres", "database"},
	{regexp.MustCompile(`pgx\.Connect\(`), "postgres", "database"},
	{regexp.MustCompile(`pq\.Open\(`), "postgres", "database"},

	// MySQL
	{regexp.MustCompile(`"(mysql://[^"]+)"`), "mysql", "database"},
	{regexp.MustCompile(`(?:sql\.Open|sqlx\.(?:Open|Connect|MustConnect))\(\s*"mysql"\s*,`), "mysql", "database"},

	// MongoDB
	{regexp.MustCompile(`"(mongodb(?:\+srv)?://[^"]+)"`), "mongodb", "database"},
	{regexp.MustCompile(`mongo\.Connect\(`), "mongodb", "database"},

	// Redis
	{regexp.MustCompile(`"(redis(?:s)?://[^"]+)"`), "redis", "database"},
	{regexp.MustCompile(`redis\.NewClient\(`), "redis", "database"},
	{regexp.MustCompile(`redis\.NewClusterClient\(`), "redis", "database"},
	{regexp.MustCompile(`redis\.NewFailoverClient\(`), "redis", "database"},

	// SQLite
	{regexp.MustCompile(`(?:sql\.Open|sqlx\.(?:Open|Connect|MustConnect))\(\s*"sqlite[3]?"\s*,`), "sqlite", "database"},

	// etcd
	{regexp.MustCompile(`clientv3\.New\(`), "etcd", "database"},
}

// Object storage patterns.
var objectStoragePatterns = []struct {
	re      *regexp.Regexp
	service string
}{
	// AWS S3
	{regexp.MustCompile(`s3\.New(?:Client|FromConfig)\(`), "s3"},
	{regexp.MustCompile(`"(s3://[^"]+)"`), "s3"},
	{regexp.MustCompile(`s3\.amazonaws\.com`), "s3"},

	// MinIO
	{regexp.MustCompile(`minio\.New\(`), "minio"},

	// GCS
	{regexp.MustCompile(`storage\.NewClient\(`), "gcs"},
	{regexp.MustCompile(`"(gs://[^"]+)"`), "gcs"},

	// Azure Blob
	{regexp.MustCompile(`azblob\.NewClient\(`), "azure-blob"},
	{regexp.MustCompile(`blob\.core\.windows\.net`), "azure-blob"},
}

// gRPC patterns.
var grpcPatterns = []*regexp.Regexp{
	regexp.MustCompile(`grpc\.(?:Dial|NewClient)\(\s*"?([^",)]+)"?`),
	regexp.MustCompile(`grpc\.DialContext\(\s*[^,]+,\s*"([^"]+)"`),
}

// Messaging patterns.
var messagingPatterns = []struct {
	re      *regexp.Regexp
	service string
}{
	// Kafka
	{regexp.MustCompile(`sarama\.NewSyncProducer\(`), "kafka"},
	{regexp.MustCompile(`sarama\.NewAsyncProducer\(`), "kafka"},
	{regexp.MustCompile(`sarama\.NewConsumer\(`), "kafka"},
	{regexp.MustCompile(`kafka\.NewReader\(`), "kafka"},
	{regexp.MustCompile(`kafka\.NewWriter\(`), "kafka"},
	{regexp.MustCompile(`confluent.*Producer`), "kafka"},

	// NATS
	{regexp.MustCompile(`nats\.Connect\(\s*"([^"]+)"`), "nats"},

	// RabbitMQ / AMQP
	{regexp.MustCompile(`"(amqp(?:s)?://[^"]+)"`), "rabbitmq"},
	{regexp.MustCompile(`amqp\.Dial\(`), "rabbitmq"},
}

// extractExternalConnections scans Go source for references to external services.
func extractExternalConnections(repoPath string) []ExternalConnection {
	var connections []ExternalConnection

	goFiles := findFiles(repoPath, []string{"**/*.go"})

	for _, fpath := range goFiles {
		if strings.Contains(fpath, "_test.go") || strings.Contains(fpath, "/vendor/") {
			continue
		}

		info, err := os.Lstat(fpath)
		if err != nil || info.Mode()&os.ModeSymlink != 0 {
			continue
		}
		if info.Size() > maxFileSize {
			log.Printf("skipping oversized file %s: %d bytes", fpath, info.Size())
			continue
		}
		data, err := os.ReadFile(fpath)
		if err != nil {
			log.Printf("warning: skipping %s: %v", fpath, err)
			continue
		}

		content := string(data)
		lines := strings.Split(content, "\n")
		source := relativePath(repoPath, fpath)

		// Track which function we're in (simple brace-counting heuristic)
		currentFunc := ""
		braceDepth := 0
		funcBraceDepth := 0

		for lineNum, line := range lines {
			stripped := strings.TrimSpace(line)
			if strings.HasPrefix(stripped, "//") {
				continue
			}

			// Track function boundaries
			if fn := extractFuncName(stripped); fn != "" {
				currentFunc = fn
				funcBraceDepth = braceDepth
			}
			// Strip comments and string literals before counting braces
			braceLine := stripStringsAndComments(line)
			braceDepth += strings.Count(braceLine, "{") - strings.Count(braceLine, "}")
			if braceDepth <= funcBraceDepth && currentFunc != "" {
				currentFunc = ""
			}

			loc := source + ":" + strconv.Itoa(lineNum+1)

			// Database connections
			for _, p := range dbConnectionPatterns {
				matches := p.re.FindStringSubmatch(stripped)
				if matches == nil {
					continue
				}
				target := ""
				if len(matches) > 1 && strings.Contains(matches[1], "://") {
					target = redactConnectionString(matches[1])
				}
				connections = append(connections, ExternalConnection{
					Type:     p.connType,
					Service:  p.service,
					Target:   target,
					Source:   loc,
					Function: currentFunc,
				})
			}

			// Object storage
			for _, p := range objectStoragePatterns {
				matches := p.re.FindStringSubmatch(stripped)
				if matches == nil {
					continue
				}
				target := ""
				if len(matches) > 1 {
					target = redactConnectionString(matches[1])
				}
				connections = append(connections, ExternalConnection{
					Type:     "object-storage",
					Service:  p.service,
					Target:   target,
					Source:   loc,
					Function: currentFunc,
				})
			}

			// gRPC
			for _, p := range grpcPatterns {
				matches := p.FindStringSubmatch(stripped)
				if matches == nil {
					continue
				}
				target := ""
				if len(matches) > 1 {
					target = matches[1]
				}
				connections = append(connections, ExternalConnection{
					Type:     "grpc",
					Service:  "grpc",
					Target:   redactTarget(target),
					Source:   loc,
					Function: currentFunc,
				})
			}

			// Messaging
			for _, p := range messagingPatterns {
				matches := p.re.FindStringSubmatch(stripped)
				if matches == nil {
					continue
				}
				target := ""
				if len(matches) > 1 && strings.Contains(matches[1], "://") {
					target = redactConnectionString(matches[1])
				}
				connections = append(connections, ExternalConnection{
					Type:     "messaging",
					Service:  p.service,
					Target:   target,
					Source:   loc,
					Function: currentFunc,
				})
			}
		}
	}

	if connections == nil {
		connections = []ExternalConnection{}
	}
	return connections
}

var stringLiteralRE = regexp.MustCompile(`"([^"\\]|\\.)*"`)
var backtickLiteralRE = regexp.MustCompile("`[^`]*`")

// stripStringsAndComments removes string literals and line comments from a line
// so that braces inside them don't affect depth counting.
func stripStringsAndComments(line string) string {
	// Strip line comments first
	if idx := strings.Index(line, "//"); idx >= 0 {
		line = line[:idx]
	}
	// Strip backtick strings
	line = backtickLiteralRE.ReplaceAllString(line, "")
	// Strip double-quoted strings
	line = stringLiteralRE.ReplaceAllString(line, "")
	return line
}

var funcDeclRE = regexp.MustCompile(`^func\s+(?:\([^)]+\)\s+)?(\w+)\s*\(`)

// extractFuncName returns the function name from a Go func declaration line.
func extractFuncName(line string) string {
	if matches := funcDeclRE.FindStringSubmatch(line); matches != nil {
		return matches[1]
	}
	return ""
}

var credentialInTargetRE = regexp.MustCompile(`[^/@]+:[^/@]+@`)

// redactConnectionString strips credentials from a connection URI.
// postgres://user:password@host:5432/db -> postgres://***@host:5432/db
func redactConnectionString(connStr string) string {
	if connStr == "" {
		return ""
	}
	parsed, err := url.Parse(connStr)
	if err != nil || parsed.Scheme == "" {
		return redactTarget(connStr)
	}
	if parsed.User != nil {
		parsed.User = url.User("***")
	}
	// Strip query parameters that commonly contain credentials
	q := parsed.Query()
	for key := range q {
		lower := strings.ToLower(key)
		if strings.Contains(lower, "password") ||
			strings.Contains(lower, "secret") ||
			strings.Contains(lower, "token") ||
			strings.Contains(lower, "key") ||
			strings.Contains(lower, "credential") {
			q.Set(key, "***")
		}
	}
	parsed.RawQuery = q.Encode()
	result := parsed.String()
	// Replace URL-encoded *** with literal *** for readability
	result = strings.ReplaceAll(result, "%2A%2A%2A", "***")
	return result
}

// redactTarget strips credential patterns from non-URI target strings.
func redactTarget(target string) string {
	if target == "" {
		return ""
	}
	// Env var interpolation is safe to keep as-is
	if strings.Contains(target, "${") || strings.Contains(target, "os.Getenv") {
		return target
	}
	// Redact user:pass@ patterns
	return credentialInTargetRE.ReplaceAllString(target, "***@")
}
