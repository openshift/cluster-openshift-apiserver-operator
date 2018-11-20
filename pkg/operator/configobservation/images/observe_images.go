package images

import (
	"github.com/golang/glog"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation"
	"github.com/openshift/library-go/pkg/operator/configobserver"
)

func ObserveInternalRegistryHostname(genericListers configobserver.Listers, existingConfig map[string]interface{}) (map[string]interface{}, []error) {
	listers := genericListers.(configobservation.Listers)
	var errs []error
	prevObservedConfig := map[string]interface{}{}

	internalRegistryHostnamePath := []string{"imagePolicyConfig", "internalRegistryHostname"}
	if currentInternalRegistryHostname, _, _ := unstructured.NestedString(existingConfig, internalRegistryHostnamePath...); len(currentInternalRegistryHostname) > 0 {
		unstructured.SetNestedField(prevObservedConfig, currentInternalRegistryHostname, internalRegistryHostnamePath...)
	}

	if !listers.ImageConfigSynced() {
		glog.Warning("images.config.openshift.io not synced")
		return prevObservedConfig, errs
	}

	observedConfig := map[string]interface{}{}
	configImage, err := listers.ImageConfigLister.Get("cluster")
	if errors.IsNotFound(err) {
		glog.Warningf("image.config.openshift.io/cluster: not found")
		return observedConfig, errs
	}
	if err != nil {
		return prevObservedConfig, errs
	}
	internalRegistryHostName := configImage.Status.InternalRegistryHostname
	if len(internalRegistryHostName) > 0 {
		unstructured.SetNestedField(observedConfig, internalRegistryHostName, internalRegistryHostnamePath...)
	}
	return observedConfig, errs
}
