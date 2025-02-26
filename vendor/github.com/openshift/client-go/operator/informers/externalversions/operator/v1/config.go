// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	context "context"
	time "time"

	apioperatorv1 "github.com/openshift/api/operator/v1"
	versioned "github.com/openshift/client-go/operator/clientset/versioned"
	internalinterfaces "github.com/openshift/client-go/operator/informers/externalversions/internalinterfaces"
	operatorv1 "github.com/openshift/client-go/operator/listers/operator/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// ConfigInformer provides access to a shared informer and lister for
// Configs.
type ConfigInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() operatorv1.ConfigLister
}

type configInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewConfigInformer constructs a new informer for Config type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewConfigInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredConfigInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredConfigInformer constructs a new informer for Config type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredConfigInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.OperatorV1().Configs().List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.OperatorV1().Configs().Watch(context.TODO(), options)
			},
		},
		&apioperatorv1.Config{},
		resyncPeriod,
		indexers,
	)
}

func (f *configInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredConfigInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *configInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&apioperatorv1.Config{}, f.defaultInformer)
}

func (f *configInformer) Lister() operatorv1.ConfigLister {
	return operatorv1.NewConfigLister(f.Informer().GetIndexer())
}
