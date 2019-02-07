package configobservation

import (
	"github.com/openshift/library-go/pkg/operator/resourcesynccontroller"
	corelistersv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"

	configlistersv1 "github.com/openshift/client-go/config/listers/config/v1"
)

type Listers struct {
	ResourceSync resourcesynccontroller.ResourceSyncer

	ImageConfigLister   configlistersv1.ImageLister
	ProjectConfigLister configlistersv1.ProjectLister
	IngressConfigLister configlistersv1.IngressLister
	EndpointsLister     corelistersv1.EndpointsLister
	PreRunCachesSynced  []cache.InformerSynced
}

func (l Listers) ResourceSyncer() resourcesynccontroller.ResourceSyncer {
	return l.ResourceSync
}

func (l Listers) PreRunHasSynced() []cache.InformerSynced {
	return l.PreRunCachesSynced
}
