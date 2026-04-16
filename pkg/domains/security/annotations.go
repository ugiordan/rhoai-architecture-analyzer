package security

const (
	AnnotCreatesRBAC           = "sec:creates_rbac"
	AnnotHandlesAdmission      = "sec:handles_admission"
	AnnotGeneratesCert         = "sec:generates_cert"
	AnnotAccessesSecret        = "sec:accesses_secret"
	AnnotCrossesNamespace      = "sec:crosses_namespace"
	AnnotConfiguresCache       = "sec:configures_cache"
	AnnotBindsSubject          = "sec:binds_subject"
	AnnotWritesPlaintextSecret = "sec:writes_plaintext_secret"
)
