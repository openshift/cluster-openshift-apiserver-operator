package configobservation

import (
	corelistersv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"

	configlistersv1 "github.com/openshift/client-go/config/listers/config/v1"
	libgoetcd "github.com/openshift/library-go/pkg/operator/configobserver/etcd"
	"github.com/openshift/library-go/pkg/operator/resourcesynccontroller"
)

var _ libgoetcd.ConfigMapLister = Listers{}
var _ libgoetcd.EndpointsLister = Listers{}

type Listers struct {
	ResourceSync resourcesynccontroller.ResourceSyncer

	APIServerLister_    configlistersv1.APIServerLister
	ImageConfigLister   configlistersv1.ImageLister
	ProjectConfigLister configlistersv1.ProjectLister
	ProxyLister_        configlistersv1.ProxyLister
	IngressConfigLister configlistersv1.IngressLister
	EndpointsLister_    corelistersv1.EndpointsLister
	PreRunCachesSynced  []cache.InformerSynced
	SecretLister_       corelistersv1.SecretLister
	ConfigmapLister_    corelistersv1.ConfigMapLister
}

func (l Listers) ResourceSyncer() resourcesynccontroller.ResourceSyncer {
	return l.ResourceSync
}

func (l Listers) SecretLister() corelistersv1.SecretLister {
	return l.SecretLister_
}

func (l Listers) PreRunHasSynced() []cache.InformerSynced {
	return l.PreRunCachesSynced
}

func (l Listers) APIServerLister() configlistersv1.APIServerLister {
	return l.APIServerLister_
}

func (l Listers) ProxyLister() configlistersv1.ProxyLister {
	return l.ProxyLister_
}

func (l Listers) ConfigMapLister() corelistersv1.ConfigMapLister {
	return l.ConfigmapLister_
}

func (l Listers) EndpointsLister() corelistersv1.EndpointsLister {
	return l.EndpointsLister_
}
