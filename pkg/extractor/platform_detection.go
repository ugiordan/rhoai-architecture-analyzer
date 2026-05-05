package extractor

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

// platformCapabilityFieldPatterns matches boolean struct fields that indicate platform detection.
var platformCapabilityFieldPatterns = []string{
	"IsOpenShift", "HasUserAPI", "IsBYOIDC", "IsAvailable",
	"HasRoute", "HasIstio", "HasKnative", "HasGatewayAPI",
}

// apiDiscoveryRE matches API discovery calls that check runtime capability.
var apiDiscoveryRE = regexp.MustCompile(
	`(?:discovery\.ServerResourcesForGroupVersion|` +
		`RESTMapper\(\)\.ResourcesFor|` +
		`apiutil\.IsGVKNamespaced|` +
		`discoveryClient\.ServerGroups)` +
		`\(`,
)

// platformSearchPaths lists directories to scan for platform detection patterns.
var platformSearchPaths = []string{
	"controllers/",
	"internal/controller/",
	"pkg/controller/",
	"pkg/reconciler/",
	"pkg/config/",
}

// resourceKindFromMethodRE extracts a resource kind from a method name.
var resourceKindFromMethodRE = regexp.MustCompile(
	`^(?:create|deploy|ensure|setup|delete|remove|add|register)([A-Z]\w*)$`,
)

// extractPlatformDetection scans Go source for platform capability checks
// and conditional resource creation patterns.
func extractPlatformDetection(repoPath string) *PlatformDetection {
	var goFiles []string
	for _, dir := range platformSearchPaths {
		fullDir := filepath.Join(repoPath, dir)
		if info, err := os.Stat(fullDir); err != nil || !info.IsDir() {
			continue
		}
		filepath.Walk(fullDir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") &&
				!strings.Contains(path, "/vendor/") {
				goFiles = append(goFiles, path)
			}
			return nil
		})
	}

	pd := &PlatformDetection{}
	capabilityNames := make(map[string]bool)

	// Pass 1: Find capability structs
	for _, fpath := range goFiles {
		caps := findCapabilityStructs(fpath, repoPath)
		for _, c := range caps {
			if !capabilityNames[c.Name] {
				capabilityNames[c.Name] = true
				pd.Capabilities = append(pd.Capabilities, c)
			}
		}
	}

	// Pass 2: API discovery calls
	for _, fpath := range goFiles {
		data, err := os.ReadFile(fpath)
		if err != nil {
			continue
		}
		content := string(data)
		if apiDiscoveryRE.MatchString(content) {
			relPath := relativePath(repoPath, fpath)
			locs := apiDiscoveryRE.FindAllStringIndex(content, -1)
			for range locs {
				pd.Capabilities = append(pd.Capabilities, PlatformCapability{
					Name:   "APIDiscoveryCheck",
					Check:  "runtime API availability check",
					Source: relPath,
				})
				break // One entry per file is enough
			}
		}
	}

	// Pass 3: Conditional resource creation
	for _, fpath := range goFiles {
		conditionals := findPlatformConditionals(fpath, repoPath, capabilityNames)
		pd.Conditionals = append(pd.Conditionals, conditionals...)
	}

	// Return nil if nothing found
	if len(pd.Capabilities) == 0 && len(pd.Conditionals) == 0 {
		return nil
	}

	return pd
}

// findCapabilityStructs scans a file for struct types with platform-detection boolean fields.
func findCapabilityStructs(fpath, repoPath string) []PlatformCapability {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, fpath, nil, parser.ParseComments)
	if err != nil {
		return nil
	}

	relPath := relativePath(repoPath, fpath)
	var caps []PlatformCapability

	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok || structType.Fields == nil {
				continue
			}

			for _, field := range structType.Fields.List {
				if len(field.Names) == 0 {
					continue
				}
				fieldName := field.Names[0].Name
				// Check if field is bool type and matches a platform capability pattern
				if ident, ok := field.Type.(*ast.Ident); ok && ident.Name == "bool" {
					if isPlatformCapabilityField(fieldName) {
						doc := ""
						if field.Doc != nil {
							doc = strings.TrimSpace(field.Doc.Text())
						}
						caps = append(caps, PlatformCapability{
							Name:   fieldName,
							Check:  doc,
							Source: relPath,
						})
					}
				}
			}
		}
	}

	return caps
}

// findPlatformConditionals scans reconciler functions for if-blocks guarded by capability checks.
func findPlatformConditionals(fpath, repoPath string, capabilityNames map[string]bool) []PlatformConditional {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, fpath, nil, parser.ParseComments)
	if err != nil {
		return nil
	}

	relPath := relativePath(repoPath, fpath)
	var conditionals []PlatformConditional

	for _, decl := range node.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			continue
		}

		// Only scan reconciler and setup methods
		name := funcDecl.Name.Name
		isReconciler := strings.HasPrefix(name, "Reconcile") || name == "SetupWithManager" ||
			strings.HasPrefix(name, "reconcile") || strings.HasPrefix(name, "setup")
		if !isReconciler {
			continue
		}

		// Walk function body for if-statements with capability conditions
		ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
			ifStmt, ok := n.(*ast.IfStmt)
			if !ok {
				return true
			}

			condStr := renderExpr(ifStmt.Cond, fset)
			if !isPlatformCondition(condStr, capabilityNames) {
				return true
			}

			// Extract resource kinds and actions from if-body
			resourceKind, action := extractResourceFromBody(ifStmt.Body)

			conditionals = append(conditionals, PlatformConditional{
				Condition:    condStr,
				ResourceKind: resourceKind,
				Action:       action,
				Source:       relPath,
			})

			return true
		})
	}

	return conditionals
}

// isPlatformCondition checks if a condition string references a known capability.
func isPlatformCondition(cond string, capabilityNames map[string]bool) bool {
	for name := range capabilityNames {
		if strings.Contains(cond, name) {
			return true
		}
	}
	// Also match common patterns even without a capability struct
	platformKeywords := []string{
		"IsOpenShift", "isOpenShift", "HasRoute", "hasRoute",
		"HasIstio", "hasIstio", "HasKnative",
	}
	for _, kw := range platformKeywords {
		if strings.Contains(cond, kw) {
			return true
		}
	}
	return false
}

// extractResourceFromBody extracts the resource kind and action from an if-body.
func extractResourceFromBody(body *ast.BlockStmt) (string, string) {
	resourceKind := ""
	action := "create"

	ast.Inspect(body, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.CallExpr:
			// Method name extraction
			var methodName string
			switch fn := node.Fun.(type) {
			case *ast.SelectorExpr:
				methodName = fn.Sel.Name
			case *ast.Ident:
				methodName = fn.Name
			}
			if methodName != "" {
				if match := resourceKindFromMethodRE.FindStringSubmatch(methodName); len(match) >= 2 {
					resourceKind = match[1]
					// Derive action from method prefix
					lower := strings.ToLower(methodName)
					for _, verb := range []string{"create", "deploy", "ensure", "setup", "delete", "remove"} {
						if strings.HasPrefix(lower, verb) {
							action = verb
							break
						}
					}
				}
				// Check for Watch registration
				if methodName == "Watches" || methodName == "Watch" || methodName == "WatchesRawSource" {
					action = "watch"
				}
			}

		case *ast.UnaryExpr:
			// Type reference: &routev1.Route{}
			if node.Op.String() == "&" {
				if comp, ok := node.X.(*ast.CompositeLit); ok {
					switch t := comp.Type.(type) {
					case *ast.SelectorExpr:
						if resourceKind == "" {
							resourceKind = t.Sel.Name
							action = "allocate"
						}
					case *ast.Ident:
						if resourceKind == "" {
							resourceKind = t.Name
							action = "allocate"
						}
					}
				}
			}
		}
		return true
	})

	if resourceKind != "" {
		// Clean up: strip trailing s for plurals, capitalize first letter
		resourceKind = strings.TrimSuffix(resourceKind, "s")
		if len(resourceKind) > 0 {
			runes := []rune(resourceKind)
			runes[0] = unicode.ToUpper(runes[0])
			resourceKind = string(runes)
		}
	}

	return resourceKind, action
}

// platformCapabilityBlocklist lists boolean field names that are NOT platform capabilities.
var platformCapabilityBlocklist = map[string]bool{
	"IsDeleted": true, "IsReady": true, "IsReconciling": true,
	"HasError": true, "IsEnabled": true, "IsDefault": true,
	"HasFinalizer": true, "IsTerminating": true,
}

// isPlatformCapabilityField checks if a boolean struct field name indicates a platform capability.
func isPlatformCapabilityField(name string) bool {
	if platformCapabilityBlocklist[name] {
		return false
	}
	for _, pattern := range platformCapabilityFieldPatterns {
		if name == pattern {
			return true
		}
	}
	// Platform-related prefixes that indicate runtime environment detection
	platformPrefixes := []string{
		"IsOpenShift", "HasRoute", "HasIstio", "HasKnative",
		"HasGateway", "HasUser", "HasConfig", "Supports",
		"IsBYO",
	}
	for _, prefix := range platformPrefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}
