package encryption

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/stretchr/testify/require"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	oauthapiv1 "github.com/openshift/api/oauth/v1"
	oauthclient "github.com/openshift/client-go/oauth/clientset/versioned/typed/oauth/v1"
	operatorv1client "github.com/openshift/client-go/operator/clientset/versioned/typed/operator/v1"
	operatorlibrary "github.com/openshift/cluster-openshift-apiserver-operator/test/library"
	library "github.com/openshift/library-go/test/library/encryption"
)

type ClientSet struct {
	OperatorClient operatorv1client.OpenShiftAPIServerInterface
	TokenClient    oauthclient.OAuthAccessTokensGetter
}

func GetClients(t testing.TB) ClientSet {
	t.Helper()

	kubeConfig, err := operatorlibrary.NewClientConfigForTest()
	require.NoError(t, err)

	operatorClient, err := operatorv1client.NewForConfig(kubeConfig)
	require.NoError(t, err)

	oc, err := oauthclient.NewForConfig(kubeConfig)
	require.NoError(t, err)

	return ClientSet{OperatorClient: operatorClient.OpenShiftAPIServers(), TokenClient: oc}
}

func CreateAndStoreTokenOfLife(t testing.TB, cs ClientSet) runtime.Object {
	t.Helper()
	{
		oldTokenOfLife, err := cs.TokenClient.OAuthAccessTokens().Get("token-aaaaaaaa-of-aaaaaaaa-life-aaaaaaaa", metav1.GetOptions{})
		if err != nil && !errors.IsNotFound(err) {
			t.Errorf("Failed to check if the route already exists, due to %v", err)
		}
		if len(oldTokenOfLife.Name) > 0 {
			t.Log("The access token already exist, removing it first")
			err := cs.TokenClient.OAuthAccessTokens().Delete(oldTokenOfLife.Name, &metav1.DeleteOptions{})
			if err != nil {
				t.Errorf("Failed to delete %s, err %v", oldTokenOfLife.Name, err)
			}
		}
	}
	t.Logf("Creating %q at cluster scope level", "token-aaaaaaaa-of-aaaaaaaa-life-aaaaaaaa")
	rawTokenOfLife := TokenOfLife(t)
	tokenOfLife, err := cs.TokenClient.OAuthAccessTokens().Create(rawTokenOfLife.(*oauthapiv1.OAuthAccessToken))
	require.NoError(t, err)
	return tokenOfLife
}

func TokenOfLife(t testing.TB) runtime.Object {
	t.Helper()
	return &oauthapiv1.OAuthAccessToken{
		ObjectMeta: metav1.ObjectMeta{
			Name: "token-aaaaaaaa-of-aaaaaaaa-life-aaaaaaaa",
		},
		RefreshToken: "I have no special talents. I am only passionately curious",
		UserName:     "kube:admin",
		Scopes:       []string{"user:full"},
		RedirectURI:  "redirect.me.to.token.of.life",
		ClientName:   "console",
		UserUID:      "non-existing-user-id",
	}
}

func GetRawTokenOfLife(t testing.TB, clientSet library.ClientSet) string {
	t.Helper()
	timeout, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	tokenOfLifeEtcdPrefix := fmt.Sprintf("/openshift.io/oauth/accesstokens/%s", "token-aaaaaaaa-of-aaaaaaaa-life-aaaaaaaa")
	resp, err := clientSet.Etcd.Get(timeout, tokenOfLifeEtcdPrefix, clientv3.WithPrefix())
	require.NoError(t, err)

	if len(resp.Kvs) != 1 {
		t.Errorf("Expected to get a single key from etcd, got %d", len(resp.Kvs))
	}

	return string(resp.Kvs[0].Value)
}
