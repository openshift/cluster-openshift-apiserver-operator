package e2e

import (
	"context"

	g "github.com/onsi/ginkgo/v2"
	o "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	configclient "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"
	test "github.com/openshift/cluster-openshift-apiserver-operator/test/library"
)

const (
	apiserverNamespace         = operatorclient.TargetNamespace   // "openshift-apiserver"
	apiserverOperatorNamespace = operatorclient.OperatorNamespace // "openshift-apiserver-operator"
	defaultDenyPolicyName      = "default-deny"
	operatorAllowPolicyName    = "allow-operator"
	operandAllowPolicyName     = "allow-apiserver"
)

var _ = g.Describe("[sig-api-machinery] openshift-apiserver operator", func() {
	g.It("[Operator][NetworkPolicy] should ensure apiserver NetworkPolicies are defined [Suite:openshift/cluster-openshift-apiserver-operator/operator/parallel]", func() {
		testAPIServerNetworkPolicies()
	})
	g.It("[Serial][Operator][NetworkPolicy] should restore apiserver NetworkPolicies after delete or mutation[Timeout:30m] [Suite:openshift/cluster-openshift-apiserver-operator/operator/serial]", func() {
		testAPIServerNetworkPolicyReconcile()
	})
})

func testAPIServerNetworkPolicies() {
	ctx := context.Background()
	g.By("Creating Kubernetes clients")
	kubeConfig, err := test.NewClientConfigForTest()
	o.Expect(err).NotTo(o.HaveOccurred())
	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	o.Expect(err).NotTo(o.HaveOccurred())
	configClient, err := configclient.NewForConfig(kubeConfig)
	o.Expect(err).NotTo(o.HaveOccurred())

	g.By("Waiting for openshift-apiserver ClusterOperator to be stable")
	err = test.WaitForClusterOperatorAvailableNotProgressingNotDegraded(ctx, configClient, "openshift-apiserver")
	o.Expect(err).NotTo(o.HaveOccurred())

	g.By("Validating NetworkPolicies in openshift-apiserver-operator")
	operatorDefaultDeny := getNetworkPolicy(ctx, kubeClient, apiserverOperatorNamespace, defaultDenyPolicyName)
	logNetworkPolicySummary("apiserver-operator/default-deny", operatorDefaultDeny)
	logNetworkPolicyDetails("apiserver-operator/default-deny", operatorDefaultDeny)
	requireDefaultDenyAll(operatorDefaultDeny)

	operatorAllowPolicy := getNetworkPolicy(ctx, kubeClient, apiserverOperatorNamespace, operatorAllowPolicyName)
	logNetworkPolicySummary("apiserver-operator/"+operatorAllowPolicyName, operatorAllowPolicy)
	logNetworkPolicyDetails("apiserver-operator/"+operatorAllowPolicyName, operatorAllowPolicy)
	requirePodSelectorLabel(operatorAllowPolicy, "app", "openshift-apiserver-operator")
	requireIngressPort(operatorAllowPolicy, corev1.ProtocolTCP, 8443)
	requireIngressAllowAll(operatorAllowPolicy, 8443)
	logEgressAllowAllTCP(operatorAllowPolicy)

	g.By("Validating NetworkPolicies in openshift-apiserver")
	operandDefaultDeny := getNetworkPolicy(ctx, kubeClient, apiserverNamespace, defaultDenyPolicyName)
	logNetworkPolicySummary("apiserver/default-deny", operandDefaultDeny)
	logNetworkPolicyDetails("apiserver/default-deny", operandDefaultDeny)
	requireDefaultDenyAll(operandDefaultDeny)

	operandAllowPolicy := getNetworkPolicy(ctx, kubeClient, apiserverNamespace, operandAllowPolicyName)
	logNetworkPolicySummary("apiserver/"+operandAllowPolicyName, operandAllowPolicy)
	logNetworkPolicyDetails("apiserver/"+operandAllowPolicyName, operandAllowPolicy)
	requirePodSelectorLabel(operandAllowPolicy, "apiserver", "true")
	requireIngressPort(operandAllowPolicy, corev1.ProtocolTCP, 8443)
	requireIngressAllowAll(operandAllowPolicy, 8443)
	logEgressAllowAllTCP(operandAllowPolicy)

	g.By("Verifying pods are ready in apiserver namespaces")
	waitForPodsReadyByLabel(ctx, kubeClient, apiserverNamespace, "apiserver=true")
	waitForPodsReadyByLabel(ctx, kubeClient, apiserverOperatorNamespace, "app=openshift-apiserver-operator")
}

func testAPIServerNetworkPolicyReconcile() {
	ctx := context.Background()
	g.By("Creating Kubernetes clients")
	kubeConfig, err := test.NewClientConfigForTest()
	o.Expect(err).NotTo(o.HaveOccurred())
	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	o.Expect(err).NotTo(o.HaveOccurred())
	configClient, err := configclient.NewForConfig(kubeConfig)
	o.Expect(err).NotTo(o.HaveOccurred())

	g.By("Waiting for openshift-apiserver ClusterOperator to be stable")
	err = test.WaitForClusterOperatorAvailableNotProgressingNotDegraded(ctx, configClient, "openshift-apiserver")
	o.Expect(err).NotTo(o.HaveOccurred())

	g.By("Capturing expected NetworkPolicy specs")
	expectedOperatorPolicy := getNetworkPolicy(ctx, kubeClient, apiserverOperatorNamespace, operatorAllowPolicyName)
	expectedOperandPolicy := getNetworkPolicy(ctx, kubeClient, apiserverNamespace, operandAllowPolicyName)
	expectedOperatorDefaultDeny := getNetworkPolicy(ctx, kubeClient, apiserverOperatorNamespace, defaultDenyPolicyName)
	expectedOperandDefaultDeny := getNetworkPolicy(ctx, kubeClient, apiserverNamespace, defaultDenyPolicyName)

	g.By("Deleting main policies and waiting for restoration")
	restoreNetworkPolicy(ctx, kubeClient, expectedOperatorPolicy)
	restoreNetworkPolicy(ctx, kubeClient, expectedOperandPolicy)

	g.By("Deleting default-deny policies and waiting for restoration")
	restoreNetworkPolicy(ctx, kubeClient, expectedOperatorDefaultDeny)
	restoreNetworkPolicy(ctx, kubeClient, expectedOperandDefaultDeny)

	g.By("Mutating main policies and waiting for reconciliation")
	mutateAndRestoreNetworkPolicy(ctx, kubeClient, apiserverOperatorNamespace, operatorAllowPolicyName)
	mutateAndRestoreNetworkPolicy(ctx, kubeClient, apiserverNamespace, operandAllowPolicyName)

	g.By("Mutating default-deny policies and waiting for reconciliation")
	mutateAndRestoreNetworkPolicy(ctx, kubeClient, apiserverOperatorNamespace, defaultDenyPolicyName)
	mutateAndRestoreNetworkPolicy(ctx, kubeClient, apiserverNamespace, defaultDenyPolicyName)

	g.By("Checking NetworkPolicy-related events (best-effort)")
	logNetworkPolicyEvents(ctx, kubeClient, []string{apiserverOperatorNamespace, apiserverNamespace}, operatorAllowPolicyName)
}
