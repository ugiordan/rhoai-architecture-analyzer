// Package testing provides test quality analysis for code property graphs.
package testing

const (
	AnnotIsTestFunc     = "test:is_test_func"
	AnnotIsTestHelper   = "test:is_test_helper"
	AnnotUsesFakeClient = "test:uses_fake_client"
	AnnotUsesEnvtest    = "test:uses_envtest"
	AnnotTableDriven    = "test:table_driven"
	AnnotErrorPath      = "test:error_path"
	AnnotSubtests       = "test:subtests"
)
