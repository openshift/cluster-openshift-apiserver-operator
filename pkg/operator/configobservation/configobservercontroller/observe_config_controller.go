package configobservercontroller

import (
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"

	imageconfiginformers "github.com/openshift/client-go/config/informers/externalversions"
	openshiftapiserveroperatorinformers "github.com/openshift/cluster-openshift-apiserver-operator/pkg/generated/informers/externalversions"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation/images"
	"github.com/openshift/library-go/pkg/operator/configobserver"
)

type ConfigObserver struct {
	*configobserver.ConfigObserver
}

// NewConfigObserver initializes a new configuration observer.
func NewConfigObserver(
	operatorClient configobserver.OperatorClient,
	operatorConfigInformers openshiftapiserveroperatorinformers.SharedInformerFactory,
	kubeInformersForEtcdNamespace kubeinformers.SharedInformerFactory,
	imageConfigInformers imageconfiginformers.SharedInformerFactory,
) *ConfigObserver {
	c := &ConfigObserver{
		ConfigObserver: configobserver.NewConfigObserver(
			operatorClient,
			configobservation.Listers{
				ImageConfigLister: imageConfigInformers.Config().V1().Images().Lister(),
				EndpointsLister:   kubeInformersForEtcdNamespace.Core().V1().Endpoints().Lister(),
				ImageConfigSynced: imageConfigInformers.Config().V1().Images().Informer().HasSynced,
				PreRunCachesSynced: []cache.InformerSynced{
					operatorConfigInformers.Openshiftapiserver().V1alpha1().OpenShiftAPIServerOperatorConfigs().Informer().HasSynced,
					kubeInformersForEtcdNamespace.Core().V1().Endpoints().Informer().HasSynced,
				},
			},
			images.ObserveInternalRegistryHostname,
			images.ObserveExternalRegistryHostnames,
			images.ObserveAllowedRegistriesForImport,
		),
	}
	operatorConfigInformers.Openshiftapiserver().V1alpha1().OpenShiftAPIServerOperatorConfigs().Informer().AddEventHandler(c.EventHandler())
	kubeInformersForEtcdNamespace.Core().V1().Endpoints().Informer().AddEventHandler(c.EventHandler())
	imageConfigInformers.Config().V1().Images().Informer().AddEventHandler(c.EventHandler())
	return c
}
