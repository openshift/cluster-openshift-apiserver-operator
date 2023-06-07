package e2e_encryption_rotation

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
	operatorencryption "github.com/openshift/cluster-openshift-apiserver-operator/test/library/encryption"
	library "github.com/openshift/library-go/test/library/encryption"
)

var provider = flag.String("provider", "aescbc", "encryption provider used by the tests")

// TestEncryptionRotation first encrypts data with aescbc key
// then it forces a key rotation by setting the "encyrption.Reason" in the operator's configuration file
func TestEncryptionRotation(t *testing.T) {
	ctx := context.TODO()
	cs := operatorencryption.GetClients(t)

	ns := fmt.Sprintf("test-encryption-on-off-%s", rand.String(4))
	_, err := cs.KubeClient.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: ns}}, metav1.CreateOptions{})
	require.NoError(t, err)
	defer cs.KubeClient.CoreV1().Namespaces().Delete(ctx, ns, metav1.DeleteOptions{})

	library.TestEncryptionRotation(t, library.RotationScenario{
		BasicScenario: library.BasicScenario{
			Namespace:                       operatorclient.GlobalMachineSpecifiedConfigNamespace,
			LabelSelector:                   "encryption.apiserver.operator.openshift.io/component" + "=" + operatorclient.TargetNamespace,
			EncryptionConfigSecretName:      fmt.Sprintf("encryption-config-%s", operatorclient.TargetNamespace),
			EncryptionConfigSecretNamespace: operatorclient.GlobalMachineSpecifiedConfigNamespace,
			OperatorNamespace:               operatorclient.OperatorNamespace,
			TargetGRs:                       operatorencryption.DefaultTargetGRs,
			AssertFunc:                      operatorencryption.AssertRoutes,
		},
		CreateResourceFunc: func(t testing.TB, _ library.ClientSet, _ string) runtime.Object {
			return operatorencryption.CreateAndStoreRouteOfLife(ctx, t, operatorencryption.GetClients(t), ns)
		},
		GetRawResourceFunc: func(t testing.TB, clientSet library.ClientSet, _ string) string {
			return operatorencryption.GetRawRouteOfLife(t, clientSet, ns)
		},
		UnsupportedConfigFunc: func(raw []byte) error {
			cs := operatorencryption.GetClients(t)
			apiServerOperator, err := cs.OperatorClient.Get(ctx, "cluster", metav1.GetOptions{})
			if err != nil {
				return err
			}
			apiServerOperator.Spec.UnsupportedConfigOverrides.Raw = raw
			_, err = cs.OperatorClient.Update(ctx, apiServerOperator, metav1.UpdateOptions{})
			return err
		},
		EncryptionProvider: configv1.EncryptionType(*provider),
	})
}
