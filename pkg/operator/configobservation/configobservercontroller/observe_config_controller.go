package configobservercontroller

import (
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"

	configinformers "github.com/openshift/client-go/config/informers/externalversions"
	"github.com/openshift/library-go/pkg/controller/factory"
	"github.com/openshift/library-go/pkg/operator/configobserver"
	libgoapiserver "github.com/openshift/library-go/pkg/operator/configobserver/apiserver"
	"github.com/openshift/library-go/pkg/operator/configobserver/proxy"
	"github.com/openshift/library-go/pkg/operator/encryption/observer"
	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/resourcesynccontroller"
	"github.com/openshift/library-go/pkg/operator/v1helpers"

	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation/etcdobserver"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation/images"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation/ingresses"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation/project"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"
)

// NewConfigObserver initializes a new configuration observer.
func NewConfigObserver(
	kubeInformers kubeinformers.SharedInformerFactory,
	kubeInformersForEtcdNamespace kubeinformers.SharedInformerFactory,
	operatorClient v1helpers.OperatorClient,
	resourceSyncer resourcesynccontroller.ResourceSyncer,
	configInformers configinformers.SharedInformerFactory,
	eventRecorder events.Recorder,
) factory.Controller {
	c := configobserver.NewConfigObserver(
		operatorClient,
		eventRecorder,
		configobservation.Listers{
			ResourceSync:        resourceSyncer,
			APIServerLister_:    configInformers.Config().V1().APIServers().Lister(),
			ImageConfigLister:   configInformers.Config().V1().Images().Lister(),
			ProjectConfigLister: configInformers.Config().V1().Projects().Lister(),
			ProxyLister_:        configInformers.Config().V1().Proxies().Lister(),
			IngressConfigLister: configInformers.Config().V1().Ingresses().Lister(),
			EndpointsLister:     kubeInformersForEtcdNamespace.Core().V1().Endpoints().Lister(),
			SecretLister_:       kubeInformers.Core().V1().Secrets().Lister(),
			PreRunCachesSynced: []cache.InformerSynced{
				configInformers.Config().V1().APIServers().Informer().HasSynced,
				configInformers.Config().V1().Images().Informer().HasSynced,
				configInformers.Config().V1().Projects().Informer().HasSynced,
				configInformers.Config().V1().Proxies().Informer().HasSynced,
				configInformers.Config().V1().Ingresses().Informer().HasSynced,
				kubeInformersForEtcdNamespace.Core().V1().Endpoints().Informer().HasSynced,
				kubeInformers.Core().V1().Secrets().Informer().HasSynced,
			},
		},
		images.ObserveInternalRegistryHostname,
		images.ObserveExternalRegistryHostnames,
		images.ObserveAllowedRegistriesForImport,
		ingresses.ObserveIngressDomain,
		etcdobserver.ObserveStorageURLs,
		libgoapiserver.ObserveTLSSecurityProfile,
		project.ObserveProjectRequestMessage,
		project.ObserveProjectRequestTemplateName,
		proxy.NewProxyObserveFunc([]string{"workloadcontroller", "proxy"}),
		observer.NewEncryptionConfigObserver(operatorclient.TargetNamespace, "/var/run/secrets/encryption-config/encryption-config"),
	)

	return c
}
