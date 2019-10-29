package encryption

import (
	"testing"

	"github.com/stretchr/testify/require"

	"k8s.io/apimachinery/pkg/runtime/schema"

	configv1 "github.com/openshift/api/config/v1"
	library "github.com/openshift/library-go/test/library/encryption"
)

var DefaultTargetGRs = []schema.GroupResource{
	{Group: "route.openshift.io", Resource: "routes"},
	{Group: "oauth.openshift.io", Resource: "oauthaccesstokens"},
	{Group: "oauth.openshift.io", Resource: "oauthauthorizetokens"},
}

func AssertRoutesAndTokens(t testing.TB, clientSet library.ClientSet, expectedMode configv1.EncryptionType, namespace, labelSelector string) {
	t.Helper()
	assertRoutes(t, clientSet.Etcd, string(expectedMode))
	assertAccessTokens(t, clientSet.Etcd, string(expectedMode))
	assertAuthTokens(t, clientSet.Etcd, string(expectedMode))
	library.AssertLastMigratedKey(t, clientSet.Kube, DefaultTargetGRs, namespace, labelSelector)
}

func assertRoutes(t testing.TB, etcdClient library.EtcdClient, expectedMode string) {
	t.Logf("Checking if all Routes where encrypted/decrypted for %q mode", expectedMode)
	totalSecrets, err := library.VerifyResources(t, etcdClient, "/kubernetes.io/routes/", expectedMode)
	t.Logf("Verified %d Secrets, err %v", totalSecrets, err)
	require.NoError(t, err)
}

func assertAccessTokens(t testing.TB, etcdClient library.EtcdClient, expectedMode string) {
	t.Logf("Checking if all OauthAccessTokens where encrypted/decrypted for %q mode", expectedMode)
	totalConfigMaps, err := library.VerifyResources(t, etcdClient, "/kubernetes.io/oauthaccesstokens/", expectedMode)
	t.Logf("Verified %d ConfigMaps, err %v", totalConfigMaps, err)
	require.NoError(t, err)
}

func assertAuthTokens(t testing.TB, etcdClient library.EtcdClient, expectedMode string) {
	t.Logf("Checking if all OauthAuthorizeTokens where encrypted/decrypted for %q mode", expectedMode)
	totalConfigMaps, err := library.VerifyResources(t, etcdClient, "/kubernetes.io/oauthauthorizetokens/", expectedMode)
	t.Logf("Verified %d ConfigMaps, err %v", totalConfigMaps, err)
	require.NoError(t, err)
}
