package connectivitycheckcontroller

import (
	"context"
	"net"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"
	v1 "github.com/openshift/api/config/v1"
	operatorv1 "github.com/openshift/api/operator/v1"
	"github.com/openshift/api/operatorcontrolplane/v1alpha1"
	operatorcontrolplaneclient "github.com/openshift/client-go/operatorcontrolplane/clientset/versioned"
	"github.com/openshift/library-go/pkg/controller/factory"
	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/resource/resourcehelper"
	"github.com/openshift/library-go/pkg/operator/v1helpers"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	corev1listers "k8s.io/client-go/listers/core/v1"

	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"
)

type ConnectivityCheckController interface {
	factory.Controller
}

func NewConnectivityCheckController(
	kubeClient kubernetes.Interface,
	operatorClient v1helpers.OperatorClient,
	kubeInformersForNamespaces v1helpers.KubeInformersForNamespaces,
	operatorcontrolplaneClient *operatorcontrolplaneclient.Clientset,
	recorder events.Recorder,
) ConnectivityCheckController {
	c := &connectivityCheckController{
		kubeClient:                 kubeClient,
		operatorClient:             operatorClient,
		operatorcontrolplaneClient: operatorcontrolplaneClient,
		podLister:                  kubeInformersForNamespaces.InformersFor("openshift-apiserver").Core().V1().Pods().Lister(),
		endpointsLister:            kubeInformersForNamespaces.InformersFor("openshift-kube-apiserver").Core().V1().Endpoints().Lister(),
		serviceLister:              kubeInformersForNamespaces.InformersFor("openshift-kube-apiserver").Core().V1().Services().Lister(),
	}
	c.Controller = factory.New().
		WithSync(c.Sync).
		WithInformers(
			operatorClient.Informer(),
			kubeInformersForNamespaces.InformersFor("openshift-apiserver").Core().V1().Pods().Informer(),
			kubeInformersForNamespaces.InformersFor("openshift-kube-apiserver").Core().V1().Endpoints().Informer(),
			kubeInformersForNamespaces.InformersFor("openshift-kube-apiserver").Core().V1().Services().Informer(),
		).
		ToController("ConnectivityCheckController", recorder.WithComponentSuffix("connectivity-check-controller"))
	return c
}

type connectivityCheckController struct {
	factory.Controller
	kubeClient                 kubernetes.Interface
	operatorClient             v1helpers.OperatorClient
	operatorcontrolplaneClient *operatorcontrolplaneclient.Clientset
	endpointsLister            corev1listers.EndpointsLister
	serviceLister              corev1listers.ServiceLister
	podLister                  corev1listers.PodLister
}

func (c *connectivityCheckController) Sync(ctx context.Context, syncContext factory.SyncContext) error {
	operatorSpec, _, _, err := c.operatorClient.GetOperatorState()
	if err != nil {
		return err
	}
	switch operatorSpec.ManagementState {
	case operatorv1.Managed:
	case operatorv1.Unmanaged:
		return nil
	case operatorv1.Removed:
		return nil
	default:
		syncContext.Recorder().Warningf("ManagementStateUnknown", "Unrecognized operator management state %q", operatorSpec.ManagementState)
		return nil
	}
	c.managePodNetworkConnectivityChecks(ctx, operatorSpec, syncContext.Recorder())
	return nil
}

func (c *connectivityCheckController) managePodNetworkConnectivityChecks(ctx context.Context, operatorSpec *operatorv1.OperatorSpec, recorder events.Recorder) {

	var templates []*v1alpha1.PodNetworkConnectivityCheck
	// each storage endpoint
	templates = append(templates, getTemplatesForStorageEndpoints(operatorSpec, recorder)...)
	// kas service IP
	templates = append(templates, getTemplatesForKubernetesServiceMonitorService(c.serviceLister, recorder)...)
	// kas default service IP
	templates = append(templates, getTemplatesForKubernetesDefaultService(recorder)...)
	// each kas endpoint IP
	templates = append(templates, getTemplatesForKubernetesEndpoints(c.endpointsLister, recorder)...)

	pods, err := c.podLister.List(labels.Set{"apiserver": "true"}.AsSelector())
	if err != nil {
		recorder.Warningf("EndpointDetectionFailure", "failed to list openshift-apiserver pods: %v", err)
		return
	}

	// create each check per static pod
	var checks []*v1alpha1.PodNetworkConnectivityCheck
	for _, pod := range pods {
		for _, template := range templates {
			check := template.DeepCopy()
			check.Name = strings.Replace(check.Name, "$(SOURCE)", pod.Name, -1)
			check.Spec.SourcePod = pod.Name
			checks = append(checks, check)
		}
	}

	pnccClient := c.operatorcontrolplaneClient.ControlplaneV1alpha1().PodNetworkConnectivityChecks(operatorclient.TargetNamespace)
	for _, check := range checks {
		_, err := pnccClient.Get(ctx, check.Name, metav1.GetOptions{})
		if err == nil {
			// already exists, skip
			continue
		}
		if apierrors.IsNotFound(err) {
			_, err = pnccClient.Create(ctx, check, metav1.CreateOptions{})
		}
		if err != nil {
			recorder.Warningf("EndpointDetectionFailure", "%s: %v", resourcehelper.FormatResourceForCLIWithNamespace(check), err)
			continue
		}
		recorder.Eventf("EndpointCheckCreated", "Created %s because it was missing.", resourcehelper.FormatResourceForCLIWithNamespace(check))
	}
}

func getTemplatesForKubernetesDefaultService(recorder events.Recorder) []*v1alpha1.PodNetworkConnectivityCheck {
	var templates []*v1alpha1.PodNetworkConnectivityCheck
	host := os.Getenv("KUBERNETES_SERVICE_HOST")
	port := os.Getenv("KUBERNETES_SERVICE_PORT")
	if len(host) == 0 || len(port) == 0 {
		recorder.Warningf("EndpointDetectionFailure", "unable to determine kubernetes service endpoint: in-cluster configuration not found")
		return templates
	}
	return append(templates, newPodNetworkProductivityCheck("kubernetes-default-service", net.JoinHostPort(host, port)))
}

func getTemplatesForKubernetesServiceMonitorService(serviceLister corev1listers.ServiceLister, recorder events.Recorder) []*v1alpha1.PodNetworkConnectivityCheck {
	var templates []*v1alpha1.PodNetworkConnectivityCheck
	for _, address := range listAddressesForKubernetesServiceMonitorService(serviceLister, recorder) {
		templates = append(templates, newPodNetworkProductivityCheck("kubernetes-apiserver-service", address))
	}
	return templates
}

func listAddressesForKubernetesServiceMonitorService(serviceLister corev1listers.ServiceLister, recorder events.Recorder) []string {
	service, err := serviceLister.Services("openshift-kube-apiserver").Get("apiserver")
	if err != nil {
		recorder.Warningf("EndpointDetectionFailure", "unable to determine openshift-kube-apiserver apiserver service endpoint: %v", err)
		return nil
	}
	for _, port := range service.Spec.Ports {
		if port.TargetPort.IntValue() == 6443 {
			return []string{net.JoinHostPort(service.Spec.ClusterIP, strconv.Itoa(int(port.Port)))}
		}
	}
	return []string{net.JoinHostPort(service.Spec.ClusterIP, "443")}
}

func getTemplatesForKubernetesEndpoints(endpointsLister corev1listers.EndpointsLister, recorder events.Recorder) []*v1alpha1.PodNetworkConnectivityCheck {
	var templates []*v1alpha1.PodNetworkConnectivityCheck
	for _, address := range listAddressesForKubeAPIServerServiceEndpoints(endpointsLister, recorder) {
		templates = append(templates, newPodNetworkProductivityCheck("kubernetes-apiserver-endpoint", address))
	}
	return templates
}

// listAddressesForKubeAPIServerServiceEndpoints returns kas api service endpoints ip
func listAddressesForKubeAPIServerServiceEndpoints(endpointsLister corev1listers.EndpointsLister, recorder events.Recorder) []string {
	var results []string
	endpoints, err := endpointsLister.Endpoints("openshift-kube-apiserver").Get("apiserver")
	if err != nil {
		recorder.Warningf("EndpointDetectionFailure", "unable to determine openshift-kube-apiserver apiserver endpoints: %v", err)
		return nil
	}
	for _, subset := range endpoints.Subsets {
		for _, address := range subset.Addresses {
			for _, port := range subset.Ports {
				results = append(results, net.JoinHostPort(address.IP, strconv.Itoa(int(port.Port))))
			}
		}
	}
	return results
}

func getTemplatesForStorageEndpoints(operatorSpec *operatorv1.OperatorSpec, recorder events.Recorder) []*v1alpha1.PodNetworkConnectivityCheck {
	var templates []*v1alpha1.PodNetworkConnectivityCheck
	for _, address := range listAddressesForStorageEndpoints(operatorSpec, recorder) {
		templates = append(templates, newPodNetworkProductivityCheck("storage-endpoint", address, withTLSClientCert("etcd-client")))
	}
	return templates
}

func listAddressesForStorageEndpoints(operatorSpec *operatorv1.OperatorSpec, recorder events.Recorder) []string {
	var results []string
	var observedConfig map[string]interface{}
	if err := yaml.Unmarshal(operatorSpec.ObservedConfig.Raw, &observedConfig); err != nil {
		recorder.Warningf("EndpointDetectionFailure", "failed to unmarshal the observedConfig: %v", err)
		return nil
	}
	urls, _, err := unstructured.NestedStringSlice(observedConfig, "storageConfig", "urls")
	if err != nil {
		recorder.Warningf("EndpointDetectionFailure", "couldn't get the storage config urls from observedConfig: %v", err)
		return nil
	}
	for _, rawStorageConfigURL := range urls {
		storageConfigURL, err := url.Parse(rawStorageConfigURL)
		if err != nil {
			recorder.Warningf("EndpointDetectionFailure", "couldn't parse a storage config url from observedConfig: %v", err)
			continue
		}
		results = append(results, storageConfigURL.Host)
	}
	return results
}

var checkNameRegex = regexp.MustCompile(`[.:\[\]]+`)

func newPodNetworkProductivityCheck(label, address string, options ...func(*v1alpha1.PodNetworkConnectivityCheck)) *v1alpha1.PodNetworkConnectivityCheck {
	check := &v1alpha1.PodNetworkConnectivityCheck{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "$(SOURCE)-to-" + label + "-" + checkNameRegex.ReplaceAllLiteralString(address, "-"),
			Namespace: operatorclient.TargetNamespace,
		},
		Spec: v1alpha1.PodNetworkConnectivityCheckSpec{
			TargetEndpoint: address,
		},
	}
	for _, option := range options {
		option(check)
	}
	return check
}

func withTLSClientCert(secretName string) func(*v1alpha1.PodNetworkConnectivityCheck) {
	return func(check *v1alpha1.PodNetworkConnectivityCheck) {
		if len(secretName) > 0 {
			check.Spec.TLSClientCert = v1.SecretNameReference{Name: secretName}
		}
	}
}
