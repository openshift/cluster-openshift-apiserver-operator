package operator

import (
	operatorv1client "github.com/openshift/client-go/operator/clientset/versioned/typed/operator/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"

	operatorv1 "github.com/openshift/api/operator/v1"

	"github.com/openshift/library-go/pkg/operator/v1helpers"

	"github.com/openshift/library-go/pkg/operator/revisioncontroller"
)

// OpenshiftDeploymentLatestRevisionClient is an LatestRevisionClient implementation for openshift-apiserver..
type OpenshiftDeploymentLatestRevisionClient struct {
	v1helpers.OperatorClient
	TypedClient operatorv1client.OpenShiftAPIServersGetter
}

var _ revisioncontroller.LatestRevisionClient = OpenshiftDeploymentLatestRevisionClient{}

func (c OpenshiftDeploymentLatestRevisionClient) GetLatestRevisionState() (*operatorv1.OperatorSpec, *operatorv1.OperatorStatus, int32, string, error) {
	o, err := c.TypedClient.OpenShiftAPIServers().Get("cluster", metav1.GetOptions{})
	if err != nil {
		return nil, nil, 0, "", err
	}
	return &o.Spec.OperatorSpec, &o.Status.OperatorStatus, o.Status.LatestAvailableRevision, o.ResourceVersion, nil
}

func (c OpenshiftDeploymentLatestRevisionClient) UpdateLatestRevisionOperatorStatus(latestAvailableRevision int32, updateFuncs ...v1helpers.UpdateStatusFunc) (*operatorv1.OperatorStatus, bool, error) {
	updated := false
	var updatedOperatorStatus *operatorv1.OperatorStatus
	err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		old, err := c.TypedClient.OpenShiftAPIServers().Get("cluster", metav1.GetOptions{})
		if err != nil {
			return err
		}

		modified := old.DeepCopy()
		for _, update := range updateFuncs {
			if err := update(&modified.Status.OperatorStatus); err != nil {
				return err
			}
		}
		modified.Status.LatestAvailableRevision = latestAvailableRevision

		if equality.Semantic.DeepEqual(old, modified) {
			// We return the newStatus which is a deep copy of oldStatus but with all update funcs applied.
			updatedOperatorStatus = &modified.Status.OperatorStatus
			return nil
		}

		modified, err = c.TypedClient.OpenShiftAPIServers().UpdateStatus(modified)
		if err != nil {
			return err
		}
		updated = true
		updatedOperatorStatus = &modified.Status.OperatorStatus
		return nil
	})

	return updatedOperatorStatus, updated, err
}
