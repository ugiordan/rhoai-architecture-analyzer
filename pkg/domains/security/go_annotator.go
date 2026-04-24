package security

import (
	"strings"

	"github.com/ugiordan/architecture-analyzer/pkg/domains"
	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

type GoAnnotator struct{}

func (a *GoAnnotator) Annotate(g *graph.CPG, archData *domains.ArchitectureData) error {
	// First pass: annotate individual nodes
	for _, cs := range g.NodesByKind(graph.NodeCallSite) {
		a.annotateCallSite(g, cs)
	}
	for _, sl := range g.NodesByKind(graph.NodeStructLiteral) {
		a.annotateStructLiteral(g, sl)
	}
	// Second pass: annotate functions based on contained nodes
	for _, fn := range g.NodesByKind(graph.NodeFunction) {
		a.annotateFunction(g, fn)
	}
	return nil
}

func (a *GoAnnotator) annotateFunction(g *graph.CPG, fn *graph.Node) {
	// sec:handles_admission: function with admission.Request parameter type
	paramTypes := strings.Join(fn.ParamTypes, ",")
	if strings.Contains(paramTypes, "admission.Request") {
		g.SetAnnotation(fn.ID, AnnotHandlesAdmission, true)
	}

	// Propagate annotations from contained call sites and struct literals
	for _, edge := range g.OutEdges(fn.ID) {
		if edge.Kind != graph.EdgeDataFlow {
			continue
		}
		target := g.GetNode(edge.To)
		if target == nil {
			continue
		}

		if target.Kind == graph.NodeCallSite {
			if target.Annotations[AnnotCreatesRBAC] {
				g.SetAnnotation(fn.ID, AnnotCreatesRBAC, true)
			}
			if target.Annotations[AnnotAccessesSecret] {
				g.SetAnnotation(fn.ID, AnnotAccessesSecret, true)
			}
			if target.Annotations[AnnotConfiguresCache] {
				g.SetAnnotation(fn.ID, AnnotConfiguresCache, true)
			}
			if target.Annotations[AnnotWritesPlaintextSecret] {
				g.SetAnnotation(fn.ID, AnnotWritesPlaintextSecret, true)
			}
		}
		if target.Kind == graph.NodeStructLiteral {
			if target.Annotations[AnnotGeneratesCert] {
				g.SetAnnotation(fn.ID, AnnotGeneratesCert, true)
			}
			if target.Annotations[AnnotConfiguresCache] {
				g.SetAnnotation(fn.ID, AnnotConfiguresCache, true)
			}
		}
	}

	// sec:binds_subject: check for system:authenticated strings in RBAC functions
	if fn.Annotations[AnnotCreatesRBAC] {
		a.checkBindsSubject(g, fn)
	}
}

func (a *GoAnnotator) annotateCallSite(g *graph.CPG, cs *graph.Node) {
	name := cs.Name
	argTypes := cs.Properties["arg_types"]
	stringArgs := cs.Properties["string_args"]

	// sec:creates_rbac
	if isClientMutation(name) && containsRBACType(argTypes) {
		g.SetAnnotation(cs.ID, AnnotCreatesRBAC, true)
	}

	// sec:accesses_secret
	if isClientAccess(name) && strings.Contains(argTypes, "Secret") {
		g.SetAnnotation(cs.ID, AnnotAccessesSecret, true)
	}

	// sec:configures_cache
	if isCacheConfig(name) {
		g.SetAnnotation(cs.ID, AnnotConfiguresCache, true)
	}

	// sec:writes_plaintext_secret
	if isFileWrite(name) && hasSecretArg(stringArgs, argTypes) {
		g.SetAnnotation(cs.ID, AnnotWritesPlaintextSecret, true)
	}
}

func (a *GoAnnotator) annotateStructLiteral(g *graph.CPG, sl *graph.Node) {
	typeName := sl.StructType
	fields := strings.Join(sl.FieldNames, ",")

	// sec:generates_cert
	if strings.Contains(typeName, "Certificate") {
		if strings.Contains(fields, "IsCA") || strings.Contains(fields, "KeyUsage") || strings.Contains(fields, "SerialNumber") {
			g.SetAnnotation(sl.ID, AnnotGeneratesCert, true)
		}
	}

	// sec:configures_cache
	if strings.Contains(typeName, "ByObject") {
		g.SetAnnotation(sl.ID, AnnotConfiguresCache, true)
	}
}

func (a *GoAnnotator) checkBindsSubject(g *graph.CPG, fn *graph.Node) {
	for _, edge := range g.OutEdges(fn.ID) {
		target := g.GetNode(edge.To)
		if target == nil {
			continue
		}
		if target.Kind == graph.NodeStructLiteral {
			sv := target.Properties["string_values"]
			if containsSubjectString(sv) {
				g.SetAnnotation(fn.ID, AnnotBindsSubject, true)
				return
			}
		}
		if target.Kind == graph.NodeCallSite {
			sa := target.Properties["string_args"]
			if containsSubjectString(sa) {
				g.SetAnnotation(fn.ID, AnnotBindsSubject, true)
				return
			}
		}
	}
}

func isClientMutation(name string) bool {
	for _, s := range []string{".Create", ".Update", ".Patch"} {
		if strings.HasSuffix(name, s) {
			return true
		}
	}
	return false
}

func isClientAccess(name string) bool {
	for _, s := range []string{".Get", ".List", ".Create", ".Update"} {
		if strings.HasSuffix(name, s) {
			return true
		}
	}
	return false
}

func containsRBACType(argTypes string) bool {
	for _, rt := range []string{"Role", "ClusterRole", "RoleBinding", "ClusterRoleBinding"} {
		if strings.Contains(argTypes, rt) {
			return true
		}
	}
	return false
}

func isCacheConfig(name string) bool {
	for _, p := range []string{"cache.New", "ctrl.NewManager", "NewCache"} {
		if strings.Contains(name, p) {
			return true
		}
	}
	return false
}

func isFileWrite(name string) bool {
	for _, p := range []string{"WriteFile", "ReplaceStringsInFile"} {
		if strings.HasSuffix(name, p) {
			return true
		}
	}
	return false
}

func hasSecretArg(stringArgs, argTypes string) bool {
	combined := strings.ToLower(stringArgs + " " + argTypes)
	for _, w := range []string{"secret", "password", "key", "token", "credential"} {
		if strings.Contains(combined, w) {
			return true
		}
	}
	return false
}

func containsSubjectString(s string) bool {
	for _, sub := range []string{"system:authenticated", "system:unauthenticated", "system:serviceaccount:"} {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}
