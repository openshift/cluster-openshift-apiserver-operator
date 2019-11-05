package e2e

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"
	test "github.com/openshift/cluster-openshift-apiserver-operator/test/library"
	operatorencryption "github.com/openshift/cluster-openshift-apiserver-operator/test/library/encryption"
	library "github.com/openshift/library-go/test/library/encryption"
)

func TestOperatorNamespace(t *testing.T) {
	kubeConfig, err := test.NewClientConfigForTest()
	if err != nil {
		t.Fatal(err)
	}
	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		t.Fatal(err)
	}
	_, err = kubeClient.CoreV1().Namespaces().Get(operatorclient.OperatorNamespace, metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestEncryptionTypeAESCBC(t *testing.T) {
	library.TestEncryptionTypeAESCBC(t, library.BasicScenario{
		Namespace:     operatorclient.GlobalMachineSpecifiedConfigNamespace,
		LabelSelector: "encryption.apiserver.operator.openshift.io/component" + "=" + operatorclient.TargetNamespace,
		TargetGRs:     operatorencryption.DefaultTargetGRs,
		AssertFunc:    operatorencryption.AssertRoutesAndTokens,
	})
}
