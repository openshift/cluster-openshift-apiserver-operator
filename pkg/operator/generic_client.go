package operator

import (
	"k8s.io/client-go/tools/cache"

	"github.com/openshift/api/operator/v1alpha1"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/generated/informers/externalversions"
)

/*
type OperatorStatusProvider interface {
	Informer() cache.SharedIndexInformer
	CurrentStatus() (operatorv1.OperatorStatus, error)
}
type : interface {
	GetOperatorState() (spec *operatorv1.OperatorSpec, status *operatorv1.OperatorStatus, resourceVersion string, err error)
	UpdateOperatorSpec(string, *operatorv1.OperatorSpec) (spec *operatorv1.OperatorSpec, resourceVersion string, err error)
	UpdateOperatorStatus(string, *operatorv1.OperatorStatus) (status *operatorv1.OperatorStatus, resourceVersion string, err error)
}
*/

type statusProvider struct {
	informers externalversions.SharedInformerFactory
}

func (p *statusProvider) Informer() cache.SharedIndexInformer {
	return p.informers.Openshiftapiserver().V1alpha1().OpenShiftAPIServerOperatorConfigs().Informer()
}

func (p *statusProvider) CurrentStatus() (v1alpha1.OperatorStatus, error) {
	instance, err := p.informers.Openshiftapiserver().V1alpha1().OpenShiftAPIServerOperatorConfigs().Lister().Get("instance")
	if err != nil {
		return v1alpha1.OperatorStatus{}, err
	}
	return instance.Status.OperatorStatus, nil
}
