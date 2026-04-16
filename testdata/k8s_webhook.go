package testdata

import (
	"context"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type DSCInitializationWebhook struct{}

func (w *DSCInitializationWebhook) Handle(ctx context.Context, req admission.Request) admission.Response {
	switch req.Operation {
	case admissionv1.Create:
		return admission.Allowed("ok")
	case admissionv1.Delete:
		return admission.Allowed("ok")
	default:
		return admission.Allowed("ok")
	}
}

func regularHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
