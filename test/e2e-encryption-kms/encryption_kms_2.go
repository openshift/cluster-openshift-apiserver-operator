package e2e_encryption_kms

import (
	"context"
	"fmt"
	"math/rand/v2"
	"testing"

	g "github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"
	library "github.com/openshift/library-go/test/library/encryption"
	librarykms "github.com/openshift/library-go/test/library/encryption/kms"
)

var _ = g.Describe("[sig-openshift-apiserver] cluster-openshift-apiserver-operator", func() {
	g.It("TestKMSEncryptionKMSToKMSMigration [OCPFeatureGate:KMSEncryption][Serial][Timeout:120m][Suite:encryption-kms-2]", func(ctx context.Context) {
		testKMSEncryptionKMSToKMSMigration(ctx, g.GinkgoTB())
	})
})

// testKMSEncryptionKMSToKMSMigration tests migration between two distinct KMS providers
// (default Vault instance and secondary Vault instance).
// This test:
// 1. Shuffles the two KMS providers and one AES provider to create a randomized migration order
// 2. Migrates between the providers in the shuffled order
// 3. Verifies route is correctly encrypted after each migration
// 4. Switches to identity (off) to verify the resource is re-written unencrypted
func testKMSEncryptionKMSToKMSMigration(ctx context.Context, t testing.TB) {
	cs := library.GetClients(t)

	ns := fmt.Sprintf("test-kms-encryption-kms-to-kms-%d", rand.IntN(4))
	_, err := cs.Kube.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: ns}}, metav1.CreateOptions{})
	require.NoError(t, err)
	defer cs.Kube.CoreV1().Namespaces().Delete(ctx, ns, metav1.DeleteOptions{})

	library.TestEncryptionProvidersMigration(ctx, t, library.ProvidersMigrationScenario{
		BasicScenario: library.BasicScenario{
			Namespace:                       operatorclient.GlobalMachineSpecifiedConfigNamespace,
			LabelSelector:                   "encryption.apiserver.operator.openshift.io/component" + "=" + operatorclient.TargetNamespace,
			EncryptionConfigSecretName:      fmt.Sprintf("encryption-config-%s", operatorclient.TargetNamespace),
			EncryptionConfigSecretNamespace: operatorclient.GlobalMachineSpecifiedConfigNamespace,
			OperatorNamespace:               operatorclient.OperatorNamespace,
			TargetGRs:                       library.OASTargetGRs,
			AssertFunc:                      library.AssertRoutes,
		},
		CreateResourceFunc: func(t testing.TB, _ library.ClientSet, namespace string) runtime.Object {
			return library.CreateAndStoreRouteOfLife(context.TODO(), t, library.GetClients(t), ns)
		},
		AssertResourceEncryptedFunc:    library.AssertRouteOfLifeEncrypted,
		AssertResourceNotEncryptedFunc: library.AssertRouteOfLifeNotEncrypted,
		ResourceFunc:                   func(t testing.TB, _ string) runtime.Object { return library.RouteOfLife(t, ns) },
		ResourceName:                   "RouteOfLife",
		EncryptionProviders: library.ShuffleEncryptionProviders([]library.EncryptionProvider{
			librarykms.DefaultVaultEncryptionProvider(ctx, t),
			librarykms.SecondaryVaultEncryptionProvider(ctx, t),
		}),
	})
}
