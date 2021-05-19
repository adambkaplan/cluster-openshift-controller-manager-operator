package e2e

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	"k8s.io/apimachinery/pkg/api/equality"

	"github.com/openshift/cluster-openshift-controller-manager-operator/pkg/operator/v311_00_assets"
	"github.com/openshift/cluster-openshift-controller-manager-operator/test/framework"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestOpenshiftControllersRBAC(t *testing.T) {
	ctx := context.Background()
	client := framework.MustNewClientset(t, nil)
	framework.MustEnsureClusterOperatorStatusIsSet(t, client)

	expectedRBACs := []struct {
		name                     string
		expectedRolesYAML        []string
		expectedRoleBindingsYAML []string
	}{
		{
			name:                     "serviceaccount-controller",
			expectedRolesYAML:        []string{"v3.11.0/openshift-controller-manager/serviceaccount-controller-clusterrole.yaml"},
			expectedRoleBindingsYAML: []string{"v3.11.0/openshift-controller-manager/serviceaccount-controller-clusterrolebinding.yaml"},
		},
		{
			name:                     "build-controller",
			expectedRolesYAML:        []string{"v3.11.0/openshift-controller-manager/build-controller-clusterrole.yaml"},
			expectedRoleBindingsYAML: []string{"v3.11.0/openshift-controller-manager/build-controller-clusterrolebinding.yaml"},
		},
	}

	for _, tc := range expectedRBACs {
		t.Run(tc.name, func(t *testing.T) {
			for _, expectedRoleYAML := range tc.expectedRolesYAML {
				expectedRole := framework.MustDecodeClusterRole(t, v311_00_assets.Asset, expectedRoleYAML)
				actualRole, err := client.RbacV1Interface.ClusterRoles().Get(ctx, expectedRole.Name, metav1.GetOptions{})
				if err != nil {
					t.Errorf("failed to get clusterrole %s: %v", expectedRole.Name, err)
					continue
				}
				if !equality.Semantic.DeepEqual(expectedRole.Rules, actualRole.Rules) {
					t.Errorf("rules for cluster role %s do match expected value: %s",
						expectedRole.Name,
						cmp.Diff(expectedRole.Rules, actualRole.Rules))
				}
			}
			for _, expectedRoleBindingYAML := range tc.expectedRoleBindingsYAML {
				expectedRoleBinding := framework.MustDecodeClusterRoleBinding(t, v311_00_assets.Asset, expectedRoleBindingYAML)
				actualRoleBinding, err := client.RbacV1Interface.ClusterRoleBindings().Get(ctx, expectedRoleBinding.Name, metav1.GetOptions{})
				if err != nil {
					t.Errorf("failed to get clusterrolebinding %s: %v", expectedRoleBinding.Name, err)
					continue
				}
				if !equality.Semantic.DeepEqual(expectedRoleBinding.RoleRef, actualRoleBinding.RoleRef) ||
					!equality.Semantic.DeepEqual(expectedRoleBinding.Subjects, actualRoleBinding.Subjects) {
					t.Errorf("clusterrolebinding %s does match expected value: %s",
						expectedRoleBinding.Name,
						cmp.Diff(expectedRoleBinding, actualRoleBinding))
				}
			}
		})
	}

}
