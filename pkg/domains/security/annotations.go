package security

// SecurityAnnotationPrefix is the namespace prefix for all security annotations.
// Queries like CGA-010 use this prefix to match annotations dynamically.
const SecurityAnnotationPrefix = "sec:"

const (
	AnnotCreatesRBAC           = "sec:creates_rbac"
	AnnotHandlesAdmission      = "sec:handles_admission"
	AnnotGeneratesCert         = "sec:generates_cert"
	AnnotAccessesSecret        = "sec:accesses_secret"
	AnnotCrossesNamespace      = "sec:crosses_namespace"
	AnnotConfiguresCache       = "sec:configures_cache"
	AnnotBindsSubject          = "sec:binds_subject"
	AnnotWritesPlaintextSecret = "sec:writes_plaintext_secret"

	// Cross-language annotations
	AnnotHandlesRequest    = "sec:handles_request"
	AnnotExecutesSQL       = "sec:executes_sql"
	AnnotDeserializesInput = "sec:deserializes_input"

	// Python-specific annotations
	AnnotSubprocessCall = "sec:subprocess_call"
	AnnotFileAccess     = "sec:file_access"
	AnnotTemplateRender = "sec:template_render"

	// TypeScript-specific annotations
	AnnotRendersHTML = "sec:renders_html"
	AnnotEvalUsage   = "sec:eval_usage"
	AnnotRedirect    = "sec:redirect"

	// Rust-specific annotations
	AnnotUnsafeBlock      = "sec:unsafe_block"
	AnnotFFICall          = "sec:ffi_call"
	AnnotCommandExecution = "sec:command_execution"
)
