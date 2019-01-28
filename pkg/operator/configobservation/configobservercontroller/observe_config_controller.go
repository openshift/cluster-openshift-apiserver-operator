package configobservercontroller

import (
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"

	"github.com/openshift/library-go/pkg/operator/configobserver"
	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/resourcesynccontroller"
	"github.com/openshift/library-go/pkg/operator/v1helpers"

	configinformers "github.com/openshift/client-go/config/informers/externalversions"
	operatorv1informers "github.com/openshift/client-go/operator/informers/externalversions"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation/images"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation/ingresses"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation/project"
)

type ConfigObserver struct {
	*configobserver.ConfigObserver
}

// NewConfigObserver initializes a new configuration observer.
func NewConfigObserver(
	operatorClient v1helpers.OperatorClient,
	resourceSyncer resourcesynccontroller.ResourceSyncer,
	operatorConfigInformers operatorv1informers.SharedInformerFactory,
	kubeInformersForEtcdNamespace kubeinformers.SharedInformerFactory,
	configInformers configinformers.SharedInformerFactory,
	eventRecorder events.Recorder,
) *ConfigObserver {
	c := &ConfigObserver{
		ConfigObserver: configobserver.NewConfigObserver(
			operatorClient,
			eventRecorder,
			configobservation.Listers{
				ResourceSync:        resourceSyncer,
				ImageConfigLister:   configInformers.Config().V1().Images().Lister(),
				ProjectConfigLister: configInformers.Config().V1().Projects().Lister(),
				IngressConfigLister: configInformers.Config().V1().Ingresses().Lister(),
				EndpointsLister:     kubeInformersForEtcdNamespace.Core().V1().Endpoints().Lister(),
				ImageConfigSynced:   configInformers.Config().V1().Images().Informer().HasSynced,
				ProjectConfigSynced: configInformers.Config().V1().Projects().Informer().HasSynced,
				IngressConfigSynced: configInformers.Config().V1().Ingresses().Informer().HasSynced,
				PreRunCachesSynced: []cache.InformerSynced{
					operatorConfigInformers.Operator().V1().OpenShiftAPIServers().Informer().HasSynced,
					kubeInformersForEtcdNamespace.Core().V1().Endpoints().Informer().HasSynced,
				},
			},
			images.ObserveInternalRegistryHostname,
			images.ObserveExternalRegistryHostnames,
			images.ObserveAllowedRegistriesForImport,
			ingresses.ObserveIngressDomain,
			project.ObserveProjectRequestMessage,
			project.ObserveProjectRequestTemplateName,
		),
	}
	operatorConfigInformers.Operator().V1().OpenShiftAPIServers().Informer().AddEventHandler(c.EventHandler())
	kubeInformersForEtcdNamespace.Core().V1().Endpoints().Informer().AddEventHandler(c.EventHandler())
	configInformers.Config().V1().Images().Informer().AddEventHandler(c.EventHandler())
	configInformers.Config().V1().Ingresses().Informer().AddEventHandler(c.EventHandler())
	configInformers.Config().V1().Projects().Informer().AddEventHandler(c.EventHandler())
	return c
}
