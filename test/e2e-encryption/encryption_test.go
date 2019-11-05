package e2e_encryption

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"
	operatorencryption "github.com/openshift/cluster-openshift-apiserver-operator/test/library/encryption"
	library "github.com/openshift/library-go/test/library/encryption"
)

func TestEncryptionTypeIdentity(t *testing.T) {
	library.TestEncryptionTypeIdentity(t, library.BasicScenario{
		Namespace:     operatorclient.GlobalMachineSpecifiedConfigNamespace,
		LabelSelector: "encryption.apiserver.operator.openshift.io/component" + "=" + operatorclient.TargetNamespace,
		TargetGRs:     operatorencryption.DefaultTargetGRs,
		AssertFunc:    operatorencryption.AssertRoutesAndTokens,
	})
}

func TestEncryptionTypeUnset(t *testing.T) {
	library.TestEncryptionTypeUnset(t, library.BasicScenario{
		Namespace:     operatorclient.GlobalMachineSpecifiedConfigNamespace,
		LabelSelector: "encryption.apiserver.operator.openshift.io/component" + "=" + operatorclient.TargetNamespace,
		TargetGRs:     operatorencryption.DefaultTargetGRs,
		AssertFunc:    operatorencryption.AssertRoutesAndTokens,
	})
}

func TestEncryptionTurnOnAndOff(t *testing.T) {
	library.TestEncryptionTurnOnAndOff(t, library.OnOffScenario{
		BasicScenario: library.BasicScenario{
			Namespace:     operatorclient.GlobalMachineSpecifiedConfigNamespace,
			LabelSelector: "encryption.apiserver.operator.openshift.io/component" + "=" + operatorclient.TargetNamespace,
			TargetGRs:     operatorencryption.DefaultTargetGRs,
			AssertFunc:    operatorencryption.AssertRoutesAndTokens,
		},
		CreateResourceFunc: func(t testing.TB, _ library.ClientSet, namespace string) runtime.Object {
			return operatorencryption.CreateAndStoreTokenOfLife(t, operatorencryption.GetClients(t))
		},
		AssertResourceEncryptedFunc:    operatorencryption.AssertTokenOfLifeEncrypted,
		AssertResourceNotEncryptedFunc: operatorencryption.AssertTokenOfLifeNotEncrypted,
		ResourceFunc:                   func(t testing.TB, _ string) runtime.Object { return operatorencryption.TokenOfLife(t) },
		ResourceName:                   "TokenOfLife",
	})
}

// TestEncryptionRotation first encrypts data with aescbc key
// then it forces a key rotation by setting the "encyrption.Reason" in the operator's configuration file
func TestEncryptionRotation(t *testing.T) {
	library.TestEncryptionRotation(t, library.RotationScenario{
		BasicScenario: library.BasicScenario{
			Namespace:     operatorclient.GlobalMachineSpecifiedConfigNamespace,
			LabelSelector: "encryption.apiserver.operator.openshift.io/component" + "=" + operatorclient.TargetNamespace,
			TargetGRs:     operatorencryption.DefaultTargetGRs,
			AssertFunc:    operatorencryption.AssertRoutesAndTokens,
		},
		CreateResourceFunc: func(t testing.TB, _ library.ClientSet, _ string) runtime.Object {
			return operatorencryption.CreateAndStoreTokenOfLife(t, operatorencryption.GetClients(t))
		},
		GetRawResourceFunc: func(t testing.TB, clientSet library.ClientSet, _ string) string {
			return operatorencryption.GetRawTokenOfLife(t, clientSet)
		},
		UnsupportedConfigFunc: func(raw []byte) error {
			cs := operatorencryption.GetClients(t)
			apiServerOperator, err := cs.OperatorClient.Get("cluster", metav1.GetOptions{})
			if err != nil {
				return err
			}
			apiServerOperator.Spec.UnsupportedConfigOverrides.Raw = raw
			_, err = cs.OperatorClient.Update(apiServerOperator)
			return err
		},
	})
}
