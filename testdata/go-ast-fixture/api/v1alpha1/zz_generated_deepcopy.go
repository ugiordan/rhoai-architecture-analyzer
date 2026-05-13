package v1alpha1

import "k8s.io/apimachinery/pkg/runtime"

// DeepCopyObject implements runtime.Object for Widget.
func (in *Widget) DeepCopyObject() runtime.Object {
	return &Widget{}
}

// DeepCopyObject implements runtime.Object for WidgetList.
func (in *WidgetList) DeepCopyObject() runtime.Object {
	return &WidgetList{}
}
