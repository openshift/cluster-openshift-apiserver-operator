package etcd

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation"
	"github.com/openshift/library-go/pkg/operator/configobserver"
)

// ObserveEtcdEndpoints observes the storage config URLs. If there is a problem observing the current storage config URLs,
// then the previously observed storage config URLs will be re-used.
func ObserveEtcdEndpoints(genericListers configobserver.Listers, existingConfig map[string]interface{}) (map[string]interface{}, []error) {
	listers := genericListers.(configobservation.Listers)
	storageConfigURLsPath := []string{"storageConfig", "urls"}
	previouslyObservedConfig := map[string]interface{}{}
	if currentStorageURLs, _, _ := unstructured.NestedStringSlice(existingConfig, storageConfigURLsPath...); len(currentStorageURLs) > 0 {
		unstructured.SetNestedStringSlice(previouslyObservedConfig, currentStorageURLs, storageConfigURLsPath...)
	}

	var errs []error

	var storageURLs []string
	etcdEndpoints, err := listers.EndpointsLister.Endpoints("kube-system").Get("etcd")
	if errors.IsNotFound(err) {
		errs = append(errs, fmt.Errorf("endpoints/etcd.kube-system: not found"))
		return previouslyObservedConfig, errs
	}
	if err != nil {
		errs = append(errs, err)
		return previouslyObservedConfig, errs
	}
	dnsSuffix := etcdEndpoints.Annotations["alpha.installer.openshift.io/dns-suffix"]
	if len(dnsSuffix) == 0 {
		errs = append(errs, fmt.Errorf("endpoints/etcd.kube-system: alpha.installer.openshift.io/dns-suffix annotation not found"))
		return previouslyObservedConfig, errs
	}
	for subsetIndex, subset := range etcdEndpoints.Subsets {
		for addressIndex, address := range subset.Addresses {
			if address.Hostname == "" {
				errs = append(errs, fmt.Errorf("endpoints/etcd.kube-system: subsets[%v]addresses[%v].hostname not found", subsetIndex, addressIndex))
				continue
			}
			storageURLs = append(storageURLs, "https://"+address.Hostname+"."+dnsSuffix+":2379")
		}
	}

	if len(storageURLs) == 0 {
		errs = append(errs, fmt.Errorf("endpoints/etcd.kube-system: no etcd endpoint addresses found"))
	}
	if len(errs) > 0 {
		return previouslyObservedConfig, errs
	}
	observedConfig := map[string]interface{}{}
	unstructured.SetNestedStringSlice(observedConfig, storageURLs, storageConfigURLsPath...)
	return observedConfig, errs
}
