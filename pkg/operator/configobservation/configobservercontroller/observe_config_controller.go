package configobservercontroller

import (
	"github.com/openshift/library-go/pkg/operator/events"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"

	configinformers "github.com/openshift/client-go/config/informers/externalversions"
	openshiftapiserveroperatorinformers "github.com/openshift/cluster-openshift-apiserver-operator/pkg/generated/informers/externalversions"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation/images"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation/ingresses"
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
	configInformers configinformers.SharedInformerFactory,
	eventRecorder events.Recorder,
) *ConfigObserver {
	c := &ConfigObserver{
		ConfigObserver: configobserver.NewConfigObserver(
			operatorClient,
			eventRecorder,
			configobservation.Listers{
				ImageConfigLister:   configInformers.Config().V1().Images().Lister(),
				IngressConfigLister: configInformers.Config().V1().Ingresses().Lister(),
				EndpointsLister:     kubeInformersForEtcdNamespace.Core().V1().Endpoints().Lister(),
				ImageConfigSynced:   configInformers.Config().V1().Images().Informer().HasSynced,
				IngressConfigSynced: configInformers.Config().V1().Ingresses().Informer().HasSynced,
				PreRunCachesSynced: []cache.InformerSynced{
					operatorConfigInformers.Openshiftapiserver().V1alpha1().OpenShiftAPIServerOperatorConfigs().Informer().HasSynced,
					kubeInformersForEtcdNamespace.Core().V1().Endpoints().Informer().HasSynced,
				},
			},
			images.ObserveInternalRegistryHostname,
			// TODO re-enable once flapping has been sorted out.
			//images.ObserveExternalRegistryHostnames,
			//images.ObserveAllowedRegistriesForImport,
			ingresses.ObserveIngressDomain,
		),
	}
	operatorConfigInformers.Openshiftapiserver().V1alpha1().OpenShiftAPIServerOperatorConfigs().Informer().AddEventHandler(c.EventHandler())
	kubeInformersForEtcdNamespace.Core().V1().Endpoints().Informer().AddEventHandler(c.EventHandler())
	configInformers.Config().V1().Images().Informer().AddEventHandler(c.EventHandler())
	configInformers.Config().V1().Ingresses().Informer().AddEventHandler(c.EventHandler())
	return c
}
