package e2e_encryption_kms

import (
	"context"
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

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
func TestKMSEncryptionOnOff(t *testing.T) {
	// Deploy the mock KMS plugin for testing.
	// NOTE: This manual deployment is only required for KMS v1. In the future,
	// the platform will manage the KMS plugins, and this code will no longer be needed.
	librarykms.DeployUpstreamMockKMSPlugin(context.Background(), t, library.GetClients(t).Kube, librarykms.WellKnownUpstreamMockKMSPluginNamespace, librarykms.WellKnownUpstreamMockKMSPluginImage)

	ctx := context.TODO()
	cs := operatorencryption.GetClients(t)

	ns := fmt.Sprintf("test-kms-encryption-on-off-%d", rand.IntN(4))
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

// TestKMSEncryptionProvidersMigration tests migration between KMS and AES encryption providers.
// This test:
// 1. Deploys the mock KMS plugin
// 2. Creates a test OAuth access token (TokenOfLife)
// 3. Randomly picks one AES encryption provider (AESGCM or AESCBC)
// 4. Shuffles the selected AES provider with KMS to create a randomized migration order
// 5. Migrates between the providers in the shuffled order
// 6. Verifies token is correctly encrypted after each migration
func TestKMSEncryptionProvidersMigration(t *testing.T) {
	librarykms.DeployUpstreamMockKMSPlugin(context.Background(), t, library.GetClients(t).Kube, librarykms.WellKnownUpstreamMockKMSPluginNamespace, librarykms.WellKnownUpstreamMockKMSPluginImage)

	ctx := context.TODO()
	cs := operatorencryption.GetClients(t)

	ns := fmt.Sprintf("test-kms-encryption-shuffle-%d", rand.IntN(4))
	_, err := cs.KubeClient.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: ns}}, metav1.CreateOptions{})
	require.NoError(t, err)
	defer cs.KubeClient.CoreV1().Namespaces().Delete(ctx, ns, metav1.DeleteOptions{})

	library.TestEncryptionProvidersMigration(t, library.ProvidersMigrationScenario{
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
		EncryptionProviders:            library.ShuffleEncryptionProviders([]configv1.EncryptionType{configv1.EncryptionTypeKMS, library.SupportedStaticEncryptionProviders[rand.IntN(len(library.SupportedStaticEncryptionProviders))]}),
	})
}
