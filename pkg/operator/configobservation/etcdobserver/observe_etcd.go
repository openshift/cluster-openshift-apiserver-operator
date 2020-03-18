package etcdobserver

import (
	"fmt"
	"net"
	"reflect"
	"sort"
	"strings"

	"github.com/openshift/library-go/pkg/operator/configobserver"
	"github.com/openshift/library-go/pkg/operator/events"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation"
)

const (
	EtcdEndpointNamespace = "openshift-etcd"
	EtcdEndpointName      = "host-etcd-2"
)

// ObserveStorageURLs observes the storage URL config.
func ObserveStorageURLs(genericListers configobserver.Listers, recorder events.Recorder, currentConfig map[string]interface{}) (map[string]interface{}, []error) {
	listers := genericListers.(configobservation.Listers)
	storageConfigURLsPath := []string{"storageConfig", "urls"}
	var errs []error

	previouslyObservedConfig := map[string]interface{}{}
	currentEtcdURLs, _, err := unstructured.NestedStringSlice(currentConfig, storageConfigURLsPath...)
	if err != nil {
		errs = append(errs, err)
	}
	if len(currentEtcdURLs) > 0 {
		if err := unstructured.SetNestedStringSlice(previouslyObservedConfig, currentEtcdURLs, storageConfigURLsPath...); err != nil {
			errs = append(errs, err)
		}
	}

	observedConfig := map[string]interface{}{}

	var etcdURLs []string
	etcdEndpoints, err := listers.EndpointsLister.Endpoints(EtcdEndpointNamespace).Get(EtcdEndpointName)
	if errors.IsNotFound(err) {
		recorder.Warningf("ObserveStorageFailed", "Required endpoints/%s in the %s namespace not found.", EtcdEndpointName, EtcdEndpointNamespace)
		errs = append(errs, fmt.Errorf("endpoints/%s in the %s namespace: not found", EtcdEndpointName, EtcdEndpointNamespace))
		return previouslyObservedConfig, errs
	}
	if err != nil {
		recorder.Warningf("ObserveStorageFailed", "Error getting endpoints/%s in the %s namespace: %v", EtcdEndpointName, EtcdEndpointNamespace, err)
		errs = append(errs, err)
		return previouslyObservedConfig, errs
	}

	for subsetIndex, subset := range etcdEndpoints.Subsets {
		for addressIndex, address := range subset.Addresses {
			ip := net.ParseIP(address.IP)
			if ip == nil {
				ipErr := fmt.Errorf("endpoints/%s in the %s namespace: subsets[%v]addresses[%v].IP is not a valid IP address", EtcdEndpointName, EtcdEndpointNamespace, subsetIndex, addressIndex)
				errs = append(errs, ipErr)
				continue
			}
			// skip placeholder ip addresses used in previous versions where the hostname was used instead
			if strings.HasPrefix(ip.String(), "192.0.2.") || strings.HasPrefix(ip.String(), "2001:db8:") {
				// not considered an error
				continue
			}
			// use the canonical representation of the ip address (not original input) when constructing the url
			if ip.To4() != nil {
				etcdURLs = append(etcdURLs, fmt.Sprintf("https://%s:2379", ip))
			} else {
				etcdURLs = append(etcdURLs, fmt.Sprintf("https://[%s]:2379", ip))
			}
		}
	}

	// do not add empty storage urls slice to observed config, we don't want override defaults with an empty slice
	if len(etcdURLs) > 0 {
		sort.Strings(etcdURLs)
		if err := unstructured.SetNestedStringSlice(observedConfig, etcdURLs, storageConfigURLsPath...); err != nil {
			errs = append(errs, err)
			return previouslyObservedConfig, errs
		}
	} else {
		err := fmt.Errorf("endpoints/%s in the %s namespace: no etcd endpoint addresses found, falling back to default etcd service", EtcdEndpointName, EtcdEndpointNamespace)
		recorder.Warningf("ObserveStorageFallback", err.Error())
		errs = append(errs, err)
	}

	sort.Strings(currentEtcdURLs)
	if !reflect.DeepEqual(currentEtcdURLs, etcdURLs) {
		recorder.Eventf("ObserveStorageUpdated", "Updated storage urls to %s", strings.Join(etcdURLs, ","))
	}

	return observedConfig, errs
}
