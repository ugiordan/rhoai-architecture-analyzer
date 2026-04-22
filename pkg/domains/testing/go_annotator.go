package testing

import (
	"strings"

	"github.com/ugiordan/architecture-analyzer/pkg/domains"
	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

type GoAnnotator struct{}

func (a *GoAnnotator) Annotate(g *graph.CPG, archData *domains.ArchitectureData) error {
	for _, fn := range g.NodesByKind(graph.NodeFunction) {
		a.annotateTestFunction(g, fn)
	}
	return nil
}

func (a *GoAnnotator) annotateTestFunction(g *graph.CPG, fn *graph.Node) {
	name := fn.Name
	file := fn.File

	// test:is_test_func: functions named Test* or Benchmark* in *_test.go
	if strings.HasSuffix(file, "_test.go") && (strings.HasPrefix(name, "Test") || strings.HasPrefix(name, "Benchmark")) {
		g.SetAnnotation(fn.ID, AnnotIsTestFunc, true)
	}

	// Check contained call sites for test patterns
	for _, edge := range g.OutEdges(fn.ID) {
		target := g.GetNode(edge.To)
		if target == nil || target.Kind != graph.NodeCallSite {
			continue
		}
		callName := target.Name

		// test:is_test_helper: function body contains t.Helper()
		if callName == "t.Helper" || callName == "b.Helper" {
			g.SetAnnotation(fn.ID, AnnotIsTestHelper, true)
		}

		// test:uses_fake_client
		if strings.Contains(callName, "fakeclient.New") ||
			strings.Contains(callName, "fake.NewSimpleClientset") ||
			strings.Contains(callName, "fake.NewClientBuilder") {
			g.SetAnnotation(fn.ID, AnnotUsesFakeClient, true)
		}

		// test:uses_envtest
		if strings.Contains(callName, "envtest.Environment") ||
			strings.Contains(callName, "testEnv.Start") {
			g.SetAnnotation(fn.ID, AnnotUsesEnvtest, true)
		}

		// test:subtests: t.Run calls
		if callName == "t.Run" || callName == "b.Run" {
			g.SetAnnotation(fn.ID, AnnotSubtests, true)
		}

		// test:error_path: error assertions
		if strings.Contains(callName, "assert.Error") ||
			strings.Contains(callName, "require.Error") ||
			strings.Contains(callName, "assert.NotNil") {
			g.SetAnnotation(fn.ID, AnnotErrorPath, true)
		}
	}

	// test:table_driven: test function with multiple struct literals and subtests
	if fn.Annotations[AnnotIsTestFunc] && fn.Annotations[AnnotSubtests] {
		structCount := 0
		for _, edge := range g.OutEdges(fn.ID) {
			target := g.GetNode(edge.To)
			if target != nil && target.Kind == graph.NodeStructLiteral {
				structCount++
			}
		}
		if structCount >= 2 {
			g.SetAnnotation(fn.ID, AnnotTableDriven, true)
		}
	}
}
