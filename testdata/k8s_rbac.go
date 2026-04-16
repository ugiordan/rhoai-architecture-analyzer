package testdata

import (
	"context"

	rbacv1 "k8s.io/api/rbac/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func bindRole(ctx context.Context, c client.Client) error {
	binding := &rbacv1.ClusterRoleBinding{
		Subjects: []rbacv1.Subject{
			{Kind: "Group", Name: "system:authenticated"},
		},
	}
	return c.Create(ctx, binding)
}
