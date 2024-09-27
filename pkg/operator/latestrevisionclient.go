package operator

import (
	"context"

	operatorv1 "github.com/openshift/api/operator/v1"
	"github.com/openshift/library-go/pkg/operator/revisioncontroller"
	"github.com/openshift/library-go/pkg/operator/v1helpers"
)

// OpenshiftDeploymentLatestRevisionClient implements LatestRevisionClient for openshift-apiserver.
type OpenshiftDeploymentLatestRevisionClient struct {
	v1helpers.OperatorClient
}

var _ revisioncontroller.LatestRevisionClient = OpenshiftDeploymentLatestRevisionClient{}

func (c OpenshiftDeploymentLatestRevisionClient) GetLatestRevisionState() (*operatorv1.OperatorSpec, *operatorv1.OperatorStatus, int32, string, error) {
	spec, status, resourceVersion, err := c.OperatorClient.GetOperatorState()
	if err != nil {
		return nil, nil, 0, "", err
	}
	return spec, status, status.LatestAvailableRevision, resourceVersion, nil
}

func (c OpenshiftDeploymentLatestRevisionClient) UpdateLatestRevisionOperatorStatus(ctx context.Context, latestAvailableRevision int32, updateFuncs ...v1helpers.UpdateStatusFunc) (*operatorv1.OperatorStatus, bool, error) {
	updateFuncs = append(updateFuncs, func(status *operatorv1.OperatorStatus) error {
		status.LatestAvailableRevision = latestAvailableRevision
		return nil
	})

	return v1helpers.UpdateStatus(ctx, c.OperatorClient, updateFuncs...)
}
