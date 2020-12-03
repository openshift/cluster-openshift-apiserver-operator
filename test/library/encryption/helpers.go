package encryption

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.etcd.io/etcd/clientv3"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	routev1 "github.com/openshift/api/route/v1"
	operatorv1client "github.com/openshift/client-go/operator/clientset/versioned/typed/operator/v1"
	routeclient "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	operatorlibrary "github.com/openshift/cluster-openshift-apiserver-operator/test/library"
	library "github.com/openshift/library-go/test/library/encryption"
)

type ClientSet struct {
	KubeClient     kubernetes.Interface
	OperatorClient operatorv1client.OpenShiftAPIServerInterface
	RouteClient    routeclient.RoutesGetter
}

func GetClients(t testing.TB) ClientSet {
	t.Helper()

	kubeConfig, err := operatorlibrary.NewClientConfigForTest()
	require.NoError(t, err)

	return GetClientsFor(t, kubeConfig)
}

func GetClientsFor(t testing.TB, kubeConfig *rest.Config) ClientSet {
	t.Helper()

	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	require.NoError(t, err)

	operatorClient, err := operatorv1client.NewForConfig(kubeConfig)
	require.NoError(t, err)

	rc, err := routeclient.NewForConfig(kubeConfig)
	require.NoError(t, err)

	return ClientSet{OperatorClient: operatorClient.OpenShiftAPIServers(), RouteClient: rc, KubeClient: kubeClient}
}

func CreateAndStoreRouteOfLife(ctx context.Context, t testing.TB, cs ClientSet, ns string) runtime.Object {
	t.Helper()
	t.Logf("Creating %q in %q namespace", "route-of-life", ns)
	rawRouteOfLife := RouteOfLife(t, ns)
	routeOfLife, err := cs.RouteClient.Routes(ns).Create(ctx, rawRouteOfLife.(*routev1.Route), metav1.CreateOptions{})
	require.NoError(t, err)
	return routeOfLife
}

func RouteOfLife(t testing.TB, ns string) runtime.Object {
	t.Helper()
	return &routev1.Route{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "route-of-life",
			Namespace: ns,
		},
		Spec: routev1.RouteSpec{
			Host: "devcluster.openshift.io",
			Port: &routev1.RoutePort{
				TargetPort: intstr.FromInt(2014),
			},
			To: routev1.RouteTargetReference{
				Name: "dummyroute",
			},
		},
	}
}

func GetRawRouteOfLife(t testing.TB, clientSet library.ClientSet, ns string) string {
	t.Helper()
	timeout, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	routeOfLifeEtcdPrefix := fmt.Sprintf("/openshift.io/routes/%s/%s", ns, "route-of-life")
	resp, err := clientSet.Etcd.Get(timeout, routeOfLifeEtcdPrefix, clientv3.WithPrefix())
	require.NoError(t, err)

	if len(resp.Kvs) != 1 {
		t.Errorf("Expected to get a single key from etcd, got %d", len(resp.Kvs))
	}

	return string(resp.Kvs[0].Value)
}
