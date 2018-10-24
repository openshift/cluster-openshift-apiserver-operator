package operator

import (
	"testing"

	configv1 "github.com/openshift/api/config/v1"
	configlistersv1 "github.com/openshift/client-go/config/listers/config/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/tools/cache"
)

func TestObserveRegistryConfig(t *testing.T) {
	const (
		expectedInternalRegistryHostname = "docker-registry.openshift-image-registry.svc.cluster.local:5000"
	)
	indexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	imageConfig := &configv1.Image{
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster",
		},
		Status: configv1.ImageStatus{
			InternalRegistryHostname: expectedInternalRegistryHostname,
		},
	}
	indexer.Add(imageConfig)
	listers := Listers{
		imageConfigLister: configlistersv1.NewImageLister(indexer),
	}
	result, err := observeInternalRegistryHostname(listers, map[string]interface{}{})
	if err != nil {
		t.Error("expected err == nil")
	}
	internalRegistryHostname, _, err := unstructured.NestedString(result, "imagePolicyConfig", "internalRegistryHostname")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if internalRegistryHostname != expectedInternalRegistryHostname {
		t.Errorf("expected internal registry hostname: %s, got %s", expectedInternalRegistryHostname, internalRegistryHostname)
	}
}
