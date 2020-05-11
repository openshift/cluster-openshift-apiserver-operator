package rbaccontroller

import (
	"context"
	"testing"
	"time"

	"github.com/openshift/library-go/pkg/operator/events"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes/fake"
)

func TestReconcileClusterRoleBinding(t *testing.T) {

	tests := []struct {
		clusterRoleBinding rbacv1.ClusterRoleBinding
		expNumOfSubjects   int
	}{
		{
			clusterRoleBinding: rbacv1.ClusterRoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "basic-users",
					Annotations: map[string]string{openShiftManagedRBACAnnotation: "true"},
				},
				Subjects: []rbacv1.Subject{
					{
						Name: user.AllUnauthenticated,
					},
					{
						Name: user.AllAuthenticated,
					},
				},
			},
			expNumOfSubjects: 1,
		},
		{
			clusterRoleBinding: rbacv1.ClusterRoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "discovery",
					Annotations: map[string]string{openShiftManagedRBACAnnotation: "true"},
				},
				Subjects: []rbacv1.Subject{
					{
						Name: user.AllAuthenticated,
					},
				},
			},
			expNumOfSubjects: 1,
		},
		{
			clusterRoleBinding: rbacv1.ClusterRoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "system:openshift:discovery",
					Annotations: map[string]string{openShiftManagedRBACAnnotation: "false"},
				},
				Subjects: []rbacv1.Subject{
					{
						Name: user.AllUnauthenticated,
					},
				},
			},
			expNumOfSubjects: 1,
		},
		{
			clusterRoleBinding: rbacv1.ClusterRoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name: "cluster-status-binding",
				},
				Subjects: []rbacv1.Subject{
					{
						Name: user.AllUnauthenticated,
					},
				},
			},
			expNumOfSubjects: 1,
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fakeClient := fake.NewSimpleClientset()
	for _, obj := range tests {
		_, err := fakeClient.RbacV1().ClusterRoleBindings().Create(ctx, &obj.clusterRoleBinding, metav1.CreateOptions{})
		if err != nil {
			t.Fatal(err)
		}
	}

	informersFactory := informers.NewSharedInformerFactory(fakeClient, 10*time.Minute)

	recorder := events.NewInMemoryRecorder("rbac_controller_test")
	c := NewRBACController(fakeClient, informersFactory, recorder)

	informersFactory.Start(ctx.Done())

	go c.Run(ctx, 1)

	time.Sleep(1 * time.Second) // allowing for reconcile to happen

	for _, obj := range tests {
		clusterRoleBinding, err := fakeClient.RbacV1().ClusterRoleBindings().Get(ctx, obj.clusterRoleBinding.Name, metav1.GetOptions{})
		if err != nil {
			t.Fatal(err)
		}

		if len(clusterRoleBinding.Subjects) != obj.expNumOfSubjects {
			t.Fatalf("%s: mismatch in number of subjects. expected %v observed %v", clusterRoleBinding.Name,
				obj.expNumOfSubjects, len(clusterRoleBinding.Subjects))
		}
	}
}
