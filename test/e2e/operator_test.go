package e2e

import (
	"context"
	"fmt"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
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
	_, err = kubeClient.CoreV1().Namespaces().Get(context.TODO(), operatorclient.OperatorNamespace, metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestRedeployOnConfigChange(t *testing.T) {
	ctx := context.TODO()
	kubeConfig, err := test.NewClientConfigForTest()
	if err != nil {
		t.Fatal(err)
	}
	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		t.Fatal(err)
	}

	deployment, err := kubeClient.AppsV1().Deployments(operatorclient.TargetNamespace).Get(ctx, "apiserver", metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}

	prevGeneration := deployment.Generation

	configCMClient := kubeClient.CoreV1().ConfigMaps(operatorclient.GlobalUserSpecifiedConfigNamespace)

	etcdClientCM, err := configCMClient.Get(ctx, "etcd-serving-ca", metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	etcdClientCM.Data["some-key"] = "non-random data"

	_, err = configCMClient.Update(ctx, etcdClientCM, metav1.UpdateOptions{})
	if err != nil {
		t.Fatal(err)
	}

	err = wait.PollImmediate(1*time.Second, 2*time.Minute, func() (done bool, err error) {
		deployment, err := kubeClient.AppsV1().Deployments(operatorclient.TargetNamespace).Get(ctx, "apiserver", metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		if deployment.Generation == prevGeneration {
			return false, nil
		}

		return true, nil
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestEncryptionTypeAESCBC(t *testing.T) {
	library.TestEncryptionTypeAESCBC(t, library.BasicScenario{
		Namespace:                       operatorclient.GlobalMachineSpecifiedConfigNamespace,
		LabelSelector:                   "encryption.apiserver.operator.openshift.io/component" + "=" + operatorclient.TargetNamespace,
		EncryptionConfigSecretName:      fmt.Sprintf("encryption-config-%s", operatorclient.TargetNamespace),
		EncryptionConfigSecretNamespace: operatorclient.GlobalMachineSpecifiedConfigNamespace,
		OperatorNamespace:               operatorclient.OperatorNamespace,
		TargetGRs:                       operatorencryption.DefaultTargetGRs,
		AssertFunc:                      operatorencryption.AssertRoutesAndTokens,
	})
}
