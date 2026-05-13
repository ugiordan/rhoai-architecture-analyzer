package controller

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	_ "example.com/widget-operator/api/v1alpha1"
)

type WidgetReconciler struct {
	client.Client
}

func (r *WidgetReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	svc := &corev1.Service{}
	r.Client.Create(ctx, svc)
	deploy := &appsv1.Deployment{}
	r.Client.Create(ctx, deploy)
	return reconcile.Result{}, nil
}
