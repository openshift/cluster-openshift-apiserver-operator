package configobservation

import (
	corelistersv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"

	configlistersv1 "github.com/openshift/client-go/config/listers/config/v1"
)

type Listers struct {
	ImageConfigLister   configlistersv1.ImageLister
	IngressConfigLister configlistersv1.IngressLister
	EndpointsLister     corelistersv1.EndpointsLister
	ImageConfigSynced   cache.InformerSynced
	IngressConfigSynced cache.InformerSynced
	PreRunCachesSynced  []cache.InformerSynced
}

func (l Listers) PreRunHasSynced() []cache.InformerSynced {
	return l.PreRunCachesSynced
}
