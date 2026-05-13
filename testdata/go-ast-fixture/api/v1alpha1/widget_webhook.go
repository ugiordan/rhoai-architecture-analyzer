package v1alpha1

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/mutate-v1alpha1-widget,mutating=true,failurePolicy=fail,sideEffects=None,groups=apps.example.com,resources=widgets,verbs=create;update,versions=v1alpha1,name=mwidget.example.com,admissionReviewVersions=v1

func (r *Widget) Default(ctx context.Context, obj runtime.Object) error {
	w := obj.(*Widget)
	if w.Spec.Image == "" {
		w.Spec.Image = "default-image:latest"
	}
	if w.Spec.GPU == "" {
		r.setGPUDefaults(w)
	}
	return nil
}

func (r *Widget) setGPUDefaults(w *Widget) {
	w.Spec.GPU = "nvidia-t4"
}

// +kubebuilder:webhook:path=/validate-v1alpha1-widget,mutating=false,failurePolicy=fail,sideEffects=None,groups=apps.example.com,resources=widgets,verbs=create;update,versions=v1alpha1,name=vwidget.example.com,admissionReviewVersions=v1

func (r *Widget) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	w := obj.(*Widget)
	var allErrs field.ErrorList
	if w.Spec.Replicas < 1 {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "replicas"), w.Spec.Replicas, "must be >= 1"))
	}
	if len(allErrs) > 0 {
		return nil, fmt.Errorf("validation failed: %v", allErrs)
	}
	return nil, nil
}

func (r *Widget) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	return r.ValidateCreate(ctx, newObj)
}

func (r *Widget) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	return nil, nil
}
