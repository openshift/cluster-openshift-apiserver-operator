package e2e_encryption_kms

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/rand"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"
	operatorencryption "github.com/openshift/cluster-openshift-apiserver-operator/test/library/encryption"
	library "github.com/openshift/library-go/test/library/encryption"
	librarykms "github.com/openshift/library-go/test/library/encryption/kms"
)

// TestKMSEncryptionOnOff tests KMS encryption on/off cycle.
// This test:
// 1. Deploys the mock KMS plugin
// 2. Creates a test OAuth access token (TokenOfLife)
// 3. Enables KMS encryption
// 4. Verifies token is encrypted
// 5. Disables encryption (Identity)
// 6. Verifies token is NOT encrypted
// 7. Re-enables KMS encryption
// 8. Verifies token is encrypted again
// 9. Disables encryption (Identity) again
// 10. Verifies token is NOT encrypted again
// 11. Cleans up the KMS plugin
func TestKMSEncryptionOnOff(t *testing.T) {
	// Deploy the mock KMS plugin for testing.
	// NOTE: This manual deployment is only required for KMS v1. In the future,
	// the platform will manage the KMS plugins, and this code will no longer be needed.
	librarykms.DeployUpstreamMockKMSPlugin(context.Background(), t, library.GetClients(t).Kube, librarykms.WellKnownUpstreamMockKMSPluginNamespace, librarykms.WellKnownUpstreamMockKMSPluginImage)

	ctx := context.TODO()
	cs := operatorencryption.GetClients(t)

	ns := fmt.Sprintf("test-kms-encryption-on-off-%s", rand.String(4))
	_, err := cs.KubeClient.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: ns}}, metav1.CreateOptions{})
	require.NoError(t, err)
	defer cs.KubeClient.CoreV1().Namespaces().Delete(ctx, ns, metav1.DeleteOptions{})

	library.TestEncryptionTurnOnAndOff(t, library.OnOffScenario{
		BasicScenario: library.BasicScenario{
			Namespace:                       operatorclient.GlobalMachineSpecifiedConfigNamespace,
			LabelSelector:                   "encryption.apiserver.operator.openshift.io/component" + "=" + operatorclient.TargetNamespace,
			EncryptionConfigSecretName:      fmt.Sprintf("encryption-config-%s", operatorclient.TargetNamespace),
			EncryptionConfigSecretNamespace: operatorclient.GlobalMachineSpecifiedConfigNamespace,
			OperatorNamespace:               operatorclient.OperatorNamespace,
			TargetGRs:                       operatorencryption.DefaultTargetGRs,
			AssertFunc:                      operatorencryption.AssertRoutes,
		},
		CreateResourceFunc: func(t testing.TB, _ library.ClientSet, namespace string) runtime.Object {
			return operatorencryption.CreateAndStoreRouteOfLife(context.TODO(), t, operatorencryption.GetClients(t), ns)
		},
		AssertResourceEncryptedFunc:    operatorencryption.AssertRouteOfLifeEncrypted,
		AssertResourceNotEncryptedFunc: operatorencryption.AssertRouteOfLifeNotEncrypted,
		ResourceFunc:                   func(t testing.TB, _ string) runtime.Object { return operatorencryption.RouteOfLife(t, ns) },
		ResourceName:                   "TokenOfLife",
		EncryptionProvider:             configv1.EncryptionTypeKMS,
	})
}
