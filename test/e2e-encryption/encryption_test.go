package e2e_encryption

import (
	"context"
	"flag"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/rand"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"
	library "github.com/openshift/library-go/test/library/encryption"
)

var provider = flag.String("provider", "aescbc", "encryption provider used by the tests")

func TestEncryptionTypeIdentity(t *testing.T) {
	library.TestEncryptionTypeIdentity(t.Context(), t, library.BasicScenario{
		Namespace:                       operatorclient.GlobalMachineSpecifiedConfigNamespace,
		LabelSelector:                   "encryption.apiserver.operator.openshift.io/component" + "=" + operatorclient.TargetNamespace,
		EncryptionConfigSecretName:      fmt.Sprintf("encryption-config-%s", operatorclient.TargetNamespace),
		EncryptionConfigSecretNamespace: operatorclient.GlobalMachineSpecifiedConfigNamespace,
		OperatorNamespace:               operatorclient.OperatorNamespace,
		TargetGRs:                       library.WellKnownOASTargetGRs,
		AssertFunc:                      library.AssertWellKnownRoutes,
	})
}

func TestEncryptionTypeUnset(t *testing.T) {
	library.TestEncryptionTypeUnset(t.Context(), t, library.BasicScenario{
		Namespace:                       operatorclient.GlobalMachineSpecifiedConfigNamespace,
		LabelSelector:                   "encryption.apiserver.operator.openshift.io/component" + "=" + operatorclient.TargetNamespace,
		EncryptionConfigSecretName:      fmt.Sprintf("encryption-config-%s", operatorclient.TargetNamespace),
		EncryptionConfigSecretNamespace: operatorclient.GlobalMachineSpecifiedConfigNamespace,
		OperatorNamespace:               operatorclient.OperatorNamespace,
		TargetGRs:                       library.WellKnownOASTargetGRs,
		AssertFunc:                      library.AssertWellKnownRoutes,
	})
}

func TestEncryptionTurnOnAndOff(t *testing.T) {
	ctx := context.TODO()
	cs := library.GetClients(t)

	ns := fmt.Sprintf("test-encryption-on-off-%s", rand.String(4))
	_, err := cs.Kube.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: ns}}, metav1.CreateOptions{})
	require.NoError(t, err)
	defer cs.Kube.CoreV1().Namespaces().Delete(ctx, ns, metav1.DeleteOptions{})

	library.TestEncryptionTurnOnAndOff(t.Context(), t, library.OnOffScenario{
		BasicScenario: library.BasicScenario{
			Namespace:                       operatorclient.GlobalMachineSpecifiedConfigNamespace,
			LabelSelector:                   "encryption.apiserver.operator.openshift.io/component" + "=" + operatorclient.TargetNamespace,
			EncryptionConfigSecretName:      fmt.Sprintf("encryption-config-%s", operatorclient.TargetNamespace),
			EncryptionConfigSecretNamespace: operatorclient.GlobalMachineSpecifiedConfigNamespace,
			OperatorNamespace:               operatorclient.OperatorNamespace,
			TargetGRs:                       library.WellKnownOASTargetGRs,
			AssertFunc:                      library.AssertWellKnownRoutes,
		},
		CreateResourceFunc: func(t testing.TB, _ library.ClientSet, namespace string) runtime.Object {
			return library.CreateAndStoreWellKnownRouteOfLife(context.TODO(), t, library.GetClients(t), ns)
		},
		AssertResourceEncryptedFunc:    library.AssertWellKnownRouteOfLifeEncrypted,
		AssertResourceNotEncryptedFunc: library.AssertWellKnownRouteOfLifeNotEncrypted,
		ResourceFunc:                   func(t testing.TB, _ string) runtime.Object { return library.WellKnownRouteOfLife(t, ns) },
		ResourceName:                   "RouteOfLife",
		EncryptionProvider:             library.EncryptionProvider{APIServerEncryption: configv1.APIServerEncryption{Type: configv1.EncryptionType(*provider)}},
	})
}
