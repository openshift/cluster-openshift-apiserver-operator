package operator

import (
	"k8s.io/client-go/tools/cache"

	operatorv1 "github.com/openshift/api/operator/v1"
	operatorconfigclientv1alpha1 "github.com/openshift/cluster-openshift-apiserver-operator/pkg/generated/clientset/versioned/typed/openshiftapiserver/v1alpha1"
	operatorclientinformers "github.com/openshift/cluster-openshift-apiserver-operator/pkg/generated/informers/externalversions"
)

type operatorClient struct {
	informers operatorclientinformers.SharedInformerFactory
	client    operatorconfigclientv1alpha1.OpenshiftapiserverV1alpha1Interface
}

func (p *operatorClient) Informer() cache.SharedIndexInformer {
	return p.informers.Openshiftapiserver().V1alpha1().OpenShiftAPIServerOperatorConfigs().Informer()
}

func (c *operatorClient) GetOperatorState() (*operatorv1.OperatorSpec, *operatorv1.OperatorStatus, string, error) {
	instance, err := c.informers.Openshiftapiserver().V1alpha1().OpenShiftAPIServerOperatorConfigs().Lister().Get("cluster")
	if err != nil {
		return nil, nil, "", err
	}

	return &instance.Spec.OperatorSpec, &instance.Status.OperatorStatus, instance.ResourceVersion, nil
}

func (c *operatorClient) UpdateOperatorSpec(resourceVersion string, spec *operatorv1.OperatorSpec) (*operatorv1.OperatorSpec, string, error) {
	original, err := c.informers.Openshiftapiserver().V1alpha1().OpenShiftAPIServerOperatorConfigs().Lister().Get("cluster")
	if err != nil {
		return nil, "", err
	}
	copy := original.DeepCopy()
	copy.ResourceVersion = resourceVersion
	copy.Spec.OperatorSpec = *spec

	ret, err := c.client.OpenShiftAPIServerOperatorConfigs().Update(copy)
	if err != nil {
		return nil, "", err
	}

	return &ret.Spec.OperatorSpec, ret.ResourceVersion, nil
}
func (c *operatorClient) UpdateOperatorStatus(resourceVersion string, status *operatorv1.OperatorStatus) (*operatorv1.OperatorStatus, error) {
	original, err := c.informers.Openshiftapiserver().V1alpha1().OpenShiftAPIServerOperatorConfigs().Lister().Get("cluster")
	if err != nil {
		return nil, err
	}
	copy := original.DeepCopy()
	copy.ResourceVersion = resourceVersion
	copy.Status.OperatorStatus = *status

	ret, err := c.client.OpenShiftAPIServerOperatorConfigs().UpdateStatus(copy)
	if err != nil {
		return nil, err
	}

	return &ret.Status.OperatorStatus, nil
}
