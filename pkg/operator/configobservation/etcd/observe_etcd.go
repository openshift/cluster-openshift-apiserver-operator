package etcd

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation"
	"github.com/openshift/library-go/pkg/operator/configobserver"
	"github.com/openshift/library-go/pkg/operator/events"
)

const (
	etcdNamespace   = "openshift-etcd"
	etcdServiceName = "etcd"
)

func ObserveEtcd(genericListers configobserver.Listers, recorder events.Recorder, currentConfig map[string]interface{}) (observedConfig map[string]interface{}, errs []error) {
	listers := genericListers.(configobservation.Listers)
	observedConfig = map[string]interface{}{}
	storageConfigURLsPath := []string{"storageConfig", "urls"}

	currentEtcdURLs, found, err := unstructured.NestedStringSlice(currentConfig, storageConfigURLsPath...)
	if err != nil {
		errs = append(errs, err)
	}
	if found {
		if err := unstructured.SetNestedStringSlice(observedConfig, currentEtcdURLs, storageConfigURLsPath...); err != nil {
			errs = append(errs, err)
		}
	}

	if _, err := listers.EndpointsLister.Endpoints(etcdNamespace).Get(etcdServiceName); err != nil {
		return
	}

	if err := unstructured.SetNestedStringSlice(observedConfig, []string{"https://etcd.openshift-etcd.svc:2379"}, storageConfigURLsPath...); err != nil {
		errs = append(errs, err)
		return
	}

	return
}
