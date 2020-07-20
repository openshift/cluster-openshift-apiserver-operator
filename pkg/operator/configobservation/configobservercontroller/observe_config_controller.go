package configobservercontroller

import (
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"

	configinformers "github.com/openshift/client-go/config/informers/externalversions"
	operatorv1informers "github.com/openshift/client-go/operator/informers/externalversions"
	"github.com/openshift/library-go/pkg/controller/factory"
	"github.com/openshift/library-go/pkg/operator/configobserver"
	libgoapiserver "github.com/openshift/library-go/pkg/operator/configobserver/apiserver"
	libgoetcd "github.com/openshift/library-go/pkg/operator/configobserver/etcd"
	"github.com/openshift/library-go/pkg/operator/configobserver/proxy"
	"github.com/openshift/library-go/pkg/operator/encryption/observer"
	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/resourcesynccontroller"
	"github.com/openshift/library-go/pkg/operator/v1helpers"

	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation"
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
	operatorConfigInformers operatorv1informers.SharedInformerFactory,
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
			EndpointsLister_:    kubeInformersForEtcdNamespace.Core().V1().Endpoints().Lister(),
			ConfigmapLister_:    kubeInformersForEtcdNamespace.Core().V1().ConfigMaps().Lister(),
			SecretLister_:       kubeInformers.Core().V1().Secrets().Lister(),
			PreRunCachesSynced: []cache.InformerSynced{
				operatorConfigInformers.Operator().V1().OpenShiftAPIServers().Informer().HasSynced,
				configInformers.Config().V1().APIServers().Informer().HasSynced,
				configInformers.Config().V1().Images().Informer().HasSynced,
				configInformers.Config().V1().Projects().Informer().HasSynced,
				configInformers.Config().V1().Proxies().Informer().HasSynced,
				configInformers.Config().V1().Ingresses().Informer().HasSynced,
				kubeInformersForEtcdNamespace.Core().V1().Endpoints().Informer().HasSynced,
				kubeInformers.Core().V1().Secrets().Informer().HasSynced,
				kubeInformersForEtcdNamespace.Core().V1().ConfigMaps().Informer().HasSynced,
			},
		},
		[]factory.Informer{operatorConfigInformers.Operator().V1().OpenShiftAPIServers().Informer()},
		images.ObserveInternalRegistryHostname,
		images.ObserveExternalRegistryHostnames,
		images.ObserveAllowedRegistriesForImport,
		ingresses.ObserveIngressDomain,
		libgoetcd.ObserveStorageURLs,
		libgoapiserver.ObserveTLSSecurityProfile,
		project.ObserveProjectRequestMessage,
		project.ObserveProjectRequestTemplateName,
		proxy.NewProxyObserveFunc([]string{"workloadcontroller", "proxy"}),
		observer.NewEncryptionConfigObserver(operatorclient.TargetNamespace, "/var/run/secrets/encryption-config/encryption-config"),
	)

	return c
}
