package encryptionstatusprovider

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"

	operatorv1 "github.com/openshift/api/operator/v1"
	applyoperatorv1 "github.com/openshift/client-go/operator/applyconfigurations/operator/v1"
	operatorclient "github.com/openshift/client-go/operator/clientset/versioned"
	operatorv1typed "github.com/openshift/client-go/operator/clientset/versioned/typed/operator/v1"

	"github.com/openshift/library-go/pkg/operator/encryption/kms"
)

// NewOpenShiftAPIServerEncryptionStatusProvider builds a kms.EncryptionStatusProvider for
// OpenShiftAPIServer/cluster from an operator client.
func NewOpenShiftAPIServerEncryptionStatusProvider(client operatorclient.Interface) (kms.EncryptionStatusProvider, error) {
	return &openShiftAPIServerEncryptionStatusProvider{client: client.OperatorV1().OpenShiftAPIServers()}, nil
}

var _ kms.EncryptionStatusProvider = &openShiftAPIServerEncryptionStatusProvider{}

type openShiftAPIServerEncryptionStatusProvider struct {
	client operatorv1typed.OpenShiftAPIServerInterface
}

func (p *openShiftAPIServerEncryptionStatusProvider) GetKMSEncryptionStatus(ctx context.Context) (*operatorv1.KMSEncryptionStatus, error) {
	obj, err := p.client.Get(ctx, "cluster", metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return &obj.Status.EncryptionStatus, nil
}

func (p *openShiftAPIServerEncryptionStatusProvider) ApplyKMSEncryptionStatus(ctx context.Context, fieldManager string, status *applyoperatorv1.KMSEncryptionStatusApplyConfiguration) error {
	_, err := p.client.ApplyStatus(
		ctx,
		applyoperatorv1.OpenShiftAPIServer("cluster").WithStatus(applyoperatorv1.OpenShiftAPIServerStatus().WithEncryptionStatus(status)),
		metav1.ApplyOptions{FieldManager: fieldManager, Force: true},
	)
	return err
}

func (p *openShiftAPIServerEncryptionStatusProvider) UpdateKMSEncryptionStatus(ctx context.Context, mutateFn func(*operatorv1.KMSEncryptionStatus)) error {
	return retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		obj, err := p.client.Get(ctx, "cluster", metav1.GetOptions{})
		if err != nil {
			return err
		}
		mutateFn(&obj.Status.EncryptionStatus)
		_, err = p.client.UpdateStatus(ctx, obj, metav1.UpdateOptions{})
		return err
	})
}
