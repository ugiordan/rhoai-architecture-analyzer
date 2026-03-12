module github.com/opendatahub-io/opendatahub-operator/v2

go 1.22

require (
	github.com/opendatahub-io/model-registry v0.2.3
	github.com/opendatahub-io/odh-model-controller v0.12.0
	sigs.k8s.io/controller-runtime v0.19.0
	k8s.io/api v0.31.0
	k8s.io/apimachinery v0.31.0
	k8s.io/client-go v0.31.0
	github.com/go-logr/logr v1.4.2
	github.com/prometheus/client_golang v1.20.0 // indirect
)
