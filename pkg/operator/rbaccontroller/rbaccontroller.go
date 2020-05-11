package rbaccontroller

import (
	"context"

	"github.com/openshift/library-go/pkg/controller/factory"
	"github.com/openshift/library-go/pkg/operator/events"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	listerrbacv1 "k8s.io/client-go/listers/rbac/v1"
	"k8s.io/client-go/tools/cache"
)

const (
	// This annotation is set to a value of "true" on RBAC resources created by the openshift-apiserver.
	openShiftManagedRBACAnnotation = "openshiftapiserver.rbac.authorization.openshift.io/managed"
)

var (
	// The named cluster role bindings are created by openshift-apiserver
	// and in previous releases may have included the unauthenticated subject.
	clusterRoleBindingNames = sets.NewString("discovery", "basic-users", "system:openshift:discovery",
		"cluster-status-binding")
)

// RBACController reconciles RBAC resources created by openshift-apiserver.
//
// If openshiftManagedRBACAnnotation is set to "true" for a cluster role binding,
// RBACController ensures the removal of the unauthenticated subject from its list of subjects.
//
// A user can opt-out of this reconciliation by setting the annotation to "false".
type RBACController struct {
	eventRecorder            events.Recorder
	kubeClient               kubernetes.Interface
	clusterRoleBindingLister listerrbacv1.ClusterRoleBindingLister
}

func NewRBACController(kubeClient kubernetes.Interface, informerFactory informers.SharedInformerFactory, recorder events.Recorder) factory.Controller {
	c := &RBACController{
		kubeClient:               kubeClient,
		clusterRoleBindingLister: informerFactory.Rbac().V1().ClusterRoleBindings().Lister(),
		eventRecorder:            recorder,
	}

	syncCtx := factory.NewSyncContext("RBACController", c.eventRecorder)

	informer := informerFactory.Rbac().V1().ClusterRoleBindings().Informer()
	informer.AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: hasMatchingAttributes,
		Handler: cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				key, err := cache.MetaNamespaceKeyFunc(obj)
				if err == nil {
					syncCtx.Queue().Add(key)
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				key, err := cache.MetaNamespaceKeyFunc(newObj)
				if err == nil {
					syncCtx.Queue().Add(key)
				}
			},
		},
	})

	return factory.New().
		WithBareInformers(informer).
		WithSyncContext(syncCtx).
		WithSync(c.sync).
		ToController("", nil)
}

// sync ensures that cluster role bindings created by openshift-apiserver do not have the unauthenticated subject.
func (c *RBACController) sync(ctx context.Context, syncContext factory.SyncContext) error {
	_, clusterRoleBindingName, err := cache.SplitMetaNamespaceKey(syncContext.QueueKey())
	if err != nil {
		return err
	}

	clusterRoleBinding, err := c.clusterRoleBindingLister.Get(clusterRoleBindingName)
	if err != nil {
		return err
	}

	updatedCopy := withUnauthenticatedSubjectRemoved(clusterRoleBinding)
	if updatedCopy == nil {
		return nil
	}

	_, err = c.kubeClient.RbacV1().ClusterRoleBindings().Update(ctx, updatedCopy, metav1.UpdateOptions{})
	return err
}

// withUnauthenticatedSubjectRemoved returns a new copy of given cluster role binding without unauthenticated user in subjects.
func withUnauthenticatedSubjectRemoved(clusterRoleBinding *rbacv1.ClusterRoleBinding) *rbacv1.ClusterRoleBinding {
	for i, subject := range clusterRoleBinding.Subjects {
		if subject.Name == user.AllUnauthenticated {
			clusterRoleBindingCopy := clusterRoleBinding.DeepCopy()
			clusterRoleBindingCopy.Subjects = append(clusterRoleBindingCopy.Subjects[:i],
				clusterRoleBindingCopy.Subjects[i+1:]...)
			return clusterRoleBindingCopy
		}
	}

	return nil
}

// hasMatchingAttributes returns true if the obj name equals one of
// clusterRoleBindingNames and it's RBAC annotation is set to "true".
func hasMatchingAttributes(obj interface{}) bool {
	clusterRoleBinding, ok := obj.(*rbacv1.ClusterRoleBinding)
	if !ok {
		return false
	}

	if !clusterRoleBindingNames.Has(clusterRoleBinding.Name) {
		return false
	}

	return (clusterRoleBinding.Annotations[openShiftManagedRBACAnnotation] == "true")
}
