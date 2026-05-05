package e2e

import (
	"context"
	"fmt"

	g "github.com/onsi/ginkgo/v2"
	o "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/kubernetes"

	configclient "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"

	test "github.com/openshift/cluster-openshift-apiserver-operator/test/library"
)

var _ = g.Describe("[sig-api-machinery] openshift-apiserver operator", func() {
	g.It("[Operator][NetworkPolicy] should enforce NetworkPolicy allow/deny basics in a test namespace [Suite:openshift/cluster-openshift-apiserver-operator/operator/parallel]", func() {
		testGenericNetworkPolicyEnforcement()
	})
	g.It("[Operator][NetworkPolicy] should enforce openshift-apiserver NetworkPolicies [Suite:openshift/cluster-openshift-apiserver-operator/operator/parallel]", func() {
		testAPIServerNetworkPolicyEnforcement()
	})
	g.It("[Operator][NetworkPolicy] should enforce openshift-apiserver-operator NetworkPolicies [Suite:openshift/cluster-openshift-apiserver-operator/operator/parallel]", func() {
		testAPIServerOperatorNetworkPolicyEnforcement()
	})
	g.It("[Operator][NetworkPolicy] should enforce cross-namespace ingress traffic [Suite:openshift/cluster-openshift-apiserver-operator/operator/parallel]", func() {
		testCrossNamespaceIngressEnforcement()
	})
	g.It("[Operator][NetworkPolicy] should block unauthorized namespace traffic [Suite:openshift/cluster-openshift-apiserver-operator/operator/parallel]", func() {
		testUnauthorizedNamespaceBlocking()
	})
})

func testGenericNetworkPolicyEnforcement() {
	ctx := context.Background()
	kubeConfig, err := test.NewClientConfigForTest()
	o.Expect(err).NotTo(o.HaveOccurred())
	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	o.Expect(err).NotTo(o.HaveOccurred())

	g.By("Creating a temporary namespace for policy enforcement checks")
	nsName := "np-enforcement-" + rand.String(5)
	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: nsName}}
	_, err = kubeClient.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	o.Expect(err).NotTo(o.HaveOccurred())
	g.DeferCleanup(func() {
		g.GinkgoWriter.Printf("deleting test namespace %s\n", nsName)
		_ = kubeClient.CoreV1().Namespaces().Delete(ctx, nsName, metav1.DeleteOptions{})
	})

	serverName := "np-server"
	clientLabels := map[string]string{"app": "np-client"}
	serverLabels := map[string]string{"app": "np-server"}

	cleanupSA := ensureServiceAccount(ctx, kubeClient, nsName)
	g.DeferCleanup(cleanupSA)

	g.GinkgoWriter.Printf("creating netexec server pod %s/%s\n", nsName, serverName)
	serverPod := netexecPod(serverName, nsName, serverLabels, 8080)
	_, err = kubeClient.CoreV1().Pods(nsName).Create(ctx, serverPod, metav1.CreateOptions{})
	o.Expect(err).NotTo(o.HaveOccurred())
	o.Expect(waitForPodReady(ctx, kubeClient, nsName, serverName)).NotTo(o.HaveOccurred())

	server, err := kubeClient.CoreV1().Pods(nsName).Get(ctx, serverName, metav1.GetOptions{})
	o.Expect(err).NotTo(o.HaveOccurred())
	o.Expect(server.Status.PodIPs).NotTo(o.BeEmpty())
	serverIPs := podIPs(server)
	g.GinkgoWriter.Printf("server pod %s/%s ips=%v\n", nsName, serverName, serverIPs)

	g.By("Verifying allow-all when no policies select the pod")
	expectConnectivity(ctx, kubeClient, nsName, clientLabels, serverIPs, 8080, true)

	g.By("Applying default deny and verifying traffic is blocked")
	g.GinkgoWriter.Printf("creating default-deny policy in %s\n", nsName)
	_, err = kubeClient.NetworkingV1().NetworkPolicies(nsName).Create(ctx, defaultDenyPolicy("default-deny", nsName), metav1.CreateOptions{})
	o.Expect(err).NotTo(o.HaveOccurred())

	g.By("Adding ingress allow only and verifying traffic is still blocked")
	g.GinkgoWriter.Printf("creating allow-ingress policy in %s\n", nsName)
	_, err = kubeClient.NetworkingV1().NetworkPolicies(nsName).Create(ctx, allowIngressPolicy("allow-ingress", nsName, serverLabels, clientLabels, 8080), metav1.CreateOptions{})
	o.Expect(err).NotTo(o.HaveOccurred())
	expectConnectivity(ctx, kubeClient, nsName, clientLabels, serverIPs, 8080, false)

	g.By("Adding egress allow and verifying traffic is permitted")
	g.GinkgoWriter.Printf("creating allow-egress policy in %s\n", nsName)
	_, err = kubeClient.NetworkingV1().NetworkPolicies(nsName).Create(ctx, allowEgressPolicy("allow-egress", nsName, clientLabels, serverLabels, 8080), metav1.CreateOptions{})
	o.Expect(err).NotTo(o.HaveOccurred())
	expectConnectivity(ctx, kubeClient, nsName, clientLabels, serverIPs, 8080, true)
}

func testAPIServerNetworkPolicyEnforcement() {
	ctx := context.Background()
	kubeConfig, err := test.NewClientConfigForTest()
	o.Expect(err).NotTo(o.HaveOccurred())
	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	o.Expect(err).NotTo(o.HaveOccurred())
	configClient, err := configclient.NewForConfig(kubeConfig)
	o.Expect(err).NotTo(o.HaveOccurred())

	namespace := "openshift-apiserver"
	serverLabels := map[string]string{"apiserver": "true"}

	g.By("Waiting for openshift-apiserver ClusterOperator to be ready")
	o.Expect(test.WaitForClusterOperatorAvailableNotProgressingNotDegraded(ctx, configClient, "openshift-apiserver")).NotTo(o.HaveOccurred())

	g.By("Creating openshift-apiserver test pods for allow/deny checks")
	g.GinkgoWriter.Printf("creating apiserver server pods in %s\n", namespace)
	allowedServerIPs, cleanupAllowed := createServerPod(ctx, kubeClient, namespace, fmt.Sprintf("np-apiserver-allowed-%s", rand.String(5)), serverLabels, 8443)
	g.DeferCleanup(cleanupAllowed)
	deniedServerIPs, cleanupDenied := createServerPod(ctx, kubeClient, namespace, fmt.Sprintf("np-apiserver-denied-%s", rand.String(5)), serverLabels, 12345)
	g.DeferCleanup(cleanupDenied)

	g.By("Verifying allowed port 8443")
	expectConnectivity(ctx, kubeClient, "default", map[string]string{"test": "any-client"}, allowedServerIPs, 8443, true)

	g.By("Verifying denied port 12345")
	expectConnectivity(ctx, kubeClient, "default", map[string]string{"test": "any-client"}, deniedServerIPs, 12345, false)

	g.By("Verifying denied ports even from allowed namespace")
	for _, port := range []int32{80, 443, 6443, 9090} {
		expectConnectivity(ctx, kubeClient, "default", map[string]string{"test": "any-client"}, allowedServerIPs, port, false)
	}
}

func testAPIServerOperatorNetworkPolicyEnforcement() {
	ctx := context.Background()
	kubeConfig, err := test.NewClientConfigForTest()
	o.Expect(err).NotTo(o.HaveOccurred())
	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	o.Expect(err).NotTo(o.HaveOccurred())
	configClient, err := configclient.NewForConfig(kubeConfig)
	o.Expect(err).NotTo(o.HaveOccurred())

	namespace := "openshift-apiserver-operator"
	serverLabels := map[string]string{"app": "openshift-apiserver-operator"}

	g.By("Waiting for openshift-apiserver ClusterOperator to be ready")
	o.Expect(test.WaitForClusterOperatorAvailableNotProgressingNotDegraded(ctx, configClient, "openshift-apiserver")).NotTo(o.HaveOccurred())

	g.By("Creating openshift-apiserver-operator test pods for policy checks")
	g.GinkgoWriter.Printf("creating apiserver-operator server pod in %s\n", namespace)
	serverIPs, cleanupServer := createServerPod(ctx, kubeClient, namespace, fmt.Sprintf("np-apiserver-op-server-%s", rand.String(5)), serverLabels, 8443)
	g.DeferCleanup(cleanupServer)

	g.By("Verifying within-namespace traffic is allowed (metrics port allows all)")
	expectConnectivity(ctx, kubeClient, namespace, map[string]string{"app": "openshift-apiserver-operator"}, serverIPs, 8443, true)

	g.By("Verifying cross-namespace traffic from monitoring is allowed")
	expectConnectivity(ctx, kubeClient, "openshift-monitoring", map[string]string{"app.kubernetes.io/name": "prometheus"}, serverIPs, 8443, true)

	g.By("Verifying unauthorized ports are denied")
	expectConnectivity(ctx, kubeClient, "openshift-monitoring", map[string]string{"app.kubernetes.io/name": "prometheus"}, serverIPs, 12345, false)

	g.By("Verifying metrics port 8443 from openshift-etcd with custom app label: should be denied")
	expectConnectivity(ctx, kubeClient, "openshift-etcd", map[string]string{"test": "client"}, serverIPs, 8443, false)

	g.By("Verifying metrics port 8443 from openshift-console: should be allowed")
	expectConnectivity(ctx, kubeClient, "openshift-console", map[string]string{"custom-app": "test-client"}, serverIPs, 8443, true)
}

func testCrossNamespaceIngressEnforcement() {
	ctx := context.Background()
	kubeConfig, err := test.NewClientConfigForTest()
	o.Expect(err).NotTo(o.HaveOccurred())
	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	o.Expect(err).NotTo(o.HaveOccurred())
	configClient, err := configclient.NewForConfig(kubeConfig)
	o.Expect(err).NotTo(o.HaveOccurred())

	g.By("Waiting for openshift-apiserver ClusterOperator to be ready")
	o.Expect(test.WaitForClusterOperatorAvailableNotProgressingNotDegraded(ctx, configClient, "openshift-apiserver")).NotTo(o.HaveOccurred())

	g.By("Creating test server pods in apiserver namespaces")
	apiserverServerIPs, cleanupAPIServer := createServerPod(ctx, kubeClient, "openshift-apiserver", fmt.Sprintf("np-apiserver-xns-%s", rand.String(5)), map[string]string{"apiserver": "true"}, 8443)
	g.DeferCleanup(cleanupAPIServer)
	operatorIPs, cleanupOperator := createServerPod(ctx, kubeClient, "openshift-apiserver-operator", fmt.Sprintf("np-apiserver-op-xns-%s", rand.String(5)), map[string]string{"app": "openshift-apiserver-operator"}, 8443)
	g.DeferCleanup(cleanupOperator)

	g.By("Testing cross-namespace ingress: monitoring -> openshift-apiserver:8443")
	expectConnectivity(ctx, kubeClient, "openshift-monitoring", map[string]string{"app.kubernetes.io/name": "prometheus"}, apiserverServerIPs, 8443, true)

	g.By("Testing cross-namespace ingress: monitoring -> apiserver-operator:8443")
	expectConnectivity(ctx, kubeClient, "openshift-monitoring", map[string]string{"app.kubernetes.io/name": "prometheus"}, operatorIPs, 8443, true)

	g.By("Testing allow-all ingress: arbitrary namespace -> openshift-apiserver:8443")
	expectConnectivity(ctx, kubeClient, "default", map[string]string{"test": "arbitrary-client"}, apiserverServerIPs, 8443, true)

	g.By("Testing denied cross-namespace: unauthorized namespace -> apiserver on unauthorized port")
	expectConnectivity(ctx, kubeClient, "default", map[string]string{"test": "arbitrary-client"}, apiserverServerIPs, 8080, false)

	g.By("Testing ingress from allowed namespace (labels may vary)")
	expectConnectivity(ctx, kubeClient, "openshift-monitoring", map[string]string{"app": "wrong-app"}, apiserverServerIPs, 8443, true)

	g.By("Testing cross-namespace: openshift-apiserver -> apiserver-operator allowed (metrics allow-all)")
	expectConnectivity(ctx, kubeClient, "openshift-apiserver", map[string]string{"apiserver": "true"}, operatorIPs, 8443, true)
}

func testUnauthorizedNamespaceBlocking() {
	ctx := context.Background()
	kubeConfig, err := test.NewClientConfigForTest()
	o.Expect(err).NotTo(o.HaveOccurred())
	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	o.Expect(err).NotTo(o.HaveOccurred())
	configClient, err := configclient.NewForConfig(kubeConfig)
	o.Expect(err).NotTo(o.HaveOccurred())

	g.By("Waiting for openshift-apiserver ClusterOperator to be ready")
	o.Expect(test.WaitForClusterOperatorAvailableNotProgressingNotDegraded(ctx, configClient, "openshift-apiserver")).NotTo(o.HaveOccurred())

	g.By("Creating test server pods in apiserver namespaces")
	apiserverServerIPs, cleanupAPIServer := createServerPod(ctx, kubeClient, "openshift-apiserver", fmt.Sprintf("np-apiserver-unauth-%s", rand.String(5)), map[string]string{"apiserver": "true"}, 8443)
	g.DeferCleanup(cleanupAPIServer)
	operatorIPs, cleanupOperator := createServerPod(ctx, kubeClient, "openshift-apiserver-operator", fmt.Sprintf("np-apiserver-op-unauth-%s", rand.String(5)), map[string]string{"app": "openshift-apiserver-operator"}, 8443)
	g.DeferCleanup(cleanupOperator)

	g.By("Testing allow-all rules: openshift-apiserver:8443 (allow from any namespace)")
	expectConnectivity(ctx, kubeClient, "default", map[string]string{"test": "any-pod"}, apiserverServerIPs, 8443, true)

	g.By("Testing allow-all metrics: any namespace -> apiserver-operator:8443")
	expectConnectivity(ctx, kubeClient, "default", map[string]string{"test": "unauthorized"}, operatorIPs, 8443, true)

	g.By("Testing metrics policy: etcd namespace with custom app label should be denied")
	expectConnectivity(ctx, kubeClient, "openshift-etcd", map[string]string{"test": "unauthorized"}, operatorIPs, 8443, false)

	g.By("Testing metrics policy: console namespace with custom app label can access metrics")
	expectConnectivity(ctx, kubeClient, "openshift-console", map[string]string{"custom-app": "test-client"}, operatorIPs, 8443, true)

	g.By("Testing port-based blocking: unauthorized port even from any namespace")
	expectConnectivity(ctx, kubeClient, "default", map[string]string{"test": "any-pod"}, apiserverServerIPs, 9999, false)

	g.By("Testing allow-all metrics: any labels from any namespace")
	expectConnectivity(ctx, kubeClient, "openshift-monitoring", map[string]string{"app": "wrong-label"}, operatorIPs, 8443, true)

	g.By("Testing multiple unauthorized ports on apiserver")
	for _, port := range []int32{80, 443, 6443, 22, 3306, 8080} {
		if port == 8443 {
			continue
		}
		expectConnectivity(ctx, kubeClient, "default", map[string]string{"test": "any-pod"}, apiserverServerIPs, port, false)
	}

	g.By("Testing cross-namespace blocking: apiserver cannot reach operator on wrong port")
	expectConnectivity(ctx, kubeClient, "openshift-apiserver", map[string]string{"apiserver": "true"}, operatorIPs, 9999, false)
}
