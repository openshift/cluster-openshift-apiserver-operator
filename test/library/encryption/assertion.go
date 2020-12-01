package encryption

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	configv1 "github.com/openshift/api/config/v1"
	routev1 "github.com/openshift/api/route/v1"
	library "github.com/openshift/library-go/test/library/encryption"
)

var DefaultTargetGRs = []schema.GroupResource{
	{Group: "route.openshift.io", Resource: "routes"},
}

func AssertRouteOfLifeEncrypted(t testing.TB, clientSet library.ClientSet, rawRouteOfLife runtime.Object) {
	t.Helper()
	routeOfLife := rawRouteOfLife.(*routev1.Route)
	rawRouteValue := GetRawRouteOfLife(t, clientSet, routeOfLife.Namespace)
	if strings.Contains(rawRouteValue, routeOfLife.Spec.To.Name) {
		t.Errorf("route not encrypted, route received from etcd have %q (plain text), raw content in etcd is %s", routeOfLife.Spec.To.Name, rawRouteValue)
	}
}

func AssertRouteOfLifeNotEncrypted(t testing.TB, clientSet library.ClientSet, rawRouteOfLife runtime.Object) {
	t.Helper()
	routeOfLife := rawRouteOfLife.(*routev1.Route)
	rawRouteValue := GetRawRouteOfLife(t, clientSet, routeOfLife.Namespace)
	if !strings.Contains(rawRouteValue, routeOfLife.Spec.To.Name) {
		t.Errorf("route received from etcd doesnt have %q (plain text), raw content in etcd is %s", routeOfLife.Spec.TLS, rawRouteValue)
	}
}

func AssertRoutes(t testing.TB, clientSet library.ClientSet, expectedMode configv1.EncryptionType, namespace, labelSelector string) {
	t.Helper()
	assertRoutes(t, clientSet.Etcd, string(expectedMode))
	library.AssertLastMigratedKey(t, clientSet.Kube, DefaultTargetGRs, namespace, labelSelector)
}

func assertRoutes(t testing.TB, etcdClient library.EtcdClient, expectedMode string) {
	t.Logf("Checking if all Routes where encrypted/decrypted for %q mode", expectedMode)
	totalRoutes, err := library.VerifyResources(t, etcdClient, "/openshift.io/routes/", expectedMode, false)
	t.Logf("Verified %d Routes", totalRoutes)
	require.NoError(t, err)
}
