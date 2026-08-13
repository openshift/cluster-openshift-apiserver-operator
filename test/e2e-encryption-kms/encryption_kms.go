package e2e_encryption_kms

import (
	"context"
	"testing"

	g "github.com/onsi/ginkgo/v2"

	library "github.com/openshift/library-go/test/library/encryption"
	librarykms "github.com/openshift/library-go/test/library/encryption/kms"
)

var _ = g.Describe("[sig-openshift-apiserver] cluster-openshift-apiserver-operator", func() {
	g.It("TestKMSEncryptionOnOff [OCPFeatureGate:KMSEncryption][Serial][Timeout:120m]", func(ctx context.Context) {
		testKMSEncryptionOnOff(ctx, g.GinkgoTB())
	})

	g.It("TestKMSEncryptionProvidersMigration [OCPFeatureGate:KMSEncryption][Serial][Timeout:120m]", func(ctx context.Context) {
		testKMSEncryptionProvidersMigration(ctx, g.GinkgoTB())
	})

})

// testKMSEncryptionOnOff tests KMS encryption on/off cycle.
// This test:
// 1. Creates a test route (RouteOfLife)
// 2. Enables KMS encryption
// 3. Verifies route is encrypted
// 4. Disables encryption (Identity)
// 5. Verifies route is NOT encrypted
// 6. Re-enables KMS encryption
// 7. Verifies route is encrypted again
// 8. Disables encryption (Identity) again
// 9. Verifies route is NOT encrypted again
func testKMSEncryptionOnOff(ctx context.Context, t testing.TB) {
	library.TestEncryptionTurnOnAndOff(ctx, t, librarykms.EncryptionTurnOnAndOffScenarios(ctx, t)...)
}

// testKMSEncryptionProvidersMigration tests migration between KMS and AES encryption providers.
// This test:
// 1. Creates a test route (RouteOfLife)
// 2. Randomly picks one AES encryption provider (AESGCM or AESCBC)
// 3. Shuffles the selected AES provider with KMS to create a randomized migration order
// 4. Migrates between the providers in the shuffled order
// 5. Verifies route is correctly encrypted after each migration
func testKMSEncryptionProvidersMigration(ctx context.Context, t testing.TB) {
	library.TestEncryptionProvidersMigration(ctx, t, librarykms.EncryptionProvidersMigrationScenarios(ctx, t)...)
}
