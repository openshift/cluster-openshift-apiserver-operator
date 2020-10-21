package e2e

import (
	"context"
	"regexp"
	"testing"

	operatorcontrolplaneclient "github.com/openshift/client-go/operatorcontrolplane/clientset/versioned"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"
	test "github.com/openshift/cluster-openshift-apiserver-operator/test/library"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestConnectivityChecksCreated(t *testing.T) {
	t.Skip("Connectivity Checks disabled.")
	kubeConfig, err := test.NewClientConfigForTest()
	require.NoError(t, err)
	operatorControlPlaneClient, err := operatorcontrolplaneclient.NewForConfig(kubeConfig)
	require.NoError(t, err)
	checks, err := operatorControlPlaneClient.ControlplaneV1alpha1().PodNetworkConnectivityChecks(operatorclient.TargetNamespace).List(context.TODO(), metav1.ListOptions{})
	require.NoError(t, err)
	testCases := []struct {
		name    string
		pattern string
		count   int
	}{
		{
			name:    "etcd-server",
			pattern: `^apiserver-.*-to-etcd-server-.*$`,
			count:   9,
		},
		{
			name:    "kubernetes-apiserver-service",
			pattern: `^apiserver-.*-to-kubernetes-apiserver-service-.*$`,
			count:   3,
		},
		{
			name:    "kubernetes-apiserver-endpoint",
			pattern: `^apiserver-.*-to-kubernetes-apiserver-endpoint-.*$`,
			count:   9,
		},
		{
			name:    "kubernetes-default-service",
			pattern: `^apiserver-.*-to-kubernetes-default-service-.*$`,
			count:   3,
		},
		{
			name:    "load-balancer-api-external",
			pattern: `^apiserver-.*-to-load-balancer-api-external$`,
			count:   3,
		},
		{
			name:    "load-balancer-api-internal",
			pattern: `^apiserver-.*-to-load-balancer-api-internal$`,
			count:   3,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var count int
			regex := regexp.MustCompile(tc.pattern)
			for _, check := range checks.Items {
				if regex.MatchString(check.Name) {
					count++
				}
			}
			assert.Equal(t, tc.count, count)
		})
	}
}
