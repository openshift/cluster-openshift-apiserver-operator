package encryption

import (
	"testing"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"
	library "github.com/openshift/library-go/test/library/encryption"
)

func TestEncryptionTypeAESCBC(t *testing.T) {
	e := library.NewE(t)
	ns := operatorclient.GlobalMachineSpecifiedConfigNamespace
	labelSelector := "encryption.apiserver.operator.openshift.io/component" + "=" + operatorclient.TargetNamespace
	clientSet := library.SetAndWaitForEncryptionType(e, configv1.EncryptionTypeAESCBC, DefaultTargetGRs, ns, labelSelector)
	AssertRoutesAndTokens(e, clientSet, configv1.EncryptionTypeAESCBC, ns, labelSelector)
}
