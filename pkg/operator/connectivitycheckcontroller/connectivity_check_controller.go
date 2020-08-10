package connectivitycheckcontroller

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"

	"github.com/ghodss/yaml"
	operatorv1 "github.com/openshift/api/operator/v1"
	"github.com/openshift/api/operatorcontrolplane/v1alpha1"
	configinformers "github.com/openshift/client-go/config/informers/externalversions"
	configv1listers "github.com/openshift/client-go/config/listers/config/v1"
	operatorcontrolplaneclient "github.com/openshift/client-go/operatorcontrolplane/clientset/versioned"
	"github.com/openshift/library-go/pkg/controller/factory"
	"github.com/openshift/library-go/pkg/operator/connectivitycheckcontroller"
	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/resource/resourcehelper"
	"github.com/openshift/library-go/pkg/operator/v1helpers"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	corev1listers "k8s.io/client-go/listers/core/v1"

	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"
)

type OpenshiftAPIServerConnectivityCheckController interface {
	connectivitycheckcontroller.ConnectivityCheckController
}

func NewOpenshiftAPIServerConnectivityCheckController(
	kubeClient kubernetes.Interface,
	operatorClient v1helpers.OperatorClient,
	kubeInformersForNamespaces v1helpers.KubeInformersForNamespaces,
	configInformers configinformers.SharedInformerFactory,
	operatorcontrolplaneClient *operatorcontrolplaneclient.Clientset,
	recorder events.Recorder,
) OpenshiftAPIServerConnectivityCheckController {
	c := openshiftAPIServerConnectivityCheckController{
		ConnectivityCheckController: connectivitycheckcontroller.NewConnectivityCheckController(
			operatorclient.TargetNamespace,
			operatorClient,
			operatorcontrolplaneClient,
			[]factory.Informer{
				operatorClient.Informer(),
				kubeInformersForNamespaces.InformersFor("openshift-apiserver").Core().V1().Pods().Informer(),
				kubeInformersForNamespaces.InformersFor("openshift-kube-apiserver").Core().V1().Endpoints().Informer(),
				kubeInformersForNamespaces.InformersFor("openshift-kube-apiserver").Core().V1().Services().Informer(),
				kubeInformersForNamespaces.InformersFor("").Core().V1().Nodes().Informer(),
				configInformers.Config().V1().Infrastructures().Informer(),
			},
			recorder,
		),
	}
	generator := &connectivityCheckTemplateProvider{
		operatorClient:             operatorClient,
		operatorcontrolplaneClient: operatorcontrolplaneClient,
		endpointsLister:            kubeInformersForNamespaces.InformersFor("openshift-kube-apiserver").Core().V1().Endpoints().Lister(),
		serviceLister:              kubeInformersForNamespaces.InformersFor("openshift-kube-apiserver").Core().V1().Services().Lister(),
		podLister:                  kubeInformersForNamespaces.InformersFor("openshift-apiserver").Core().V1().Pods().Lister(),
		nodeLister:                 kubeInformersForNamespaces.InformersFor("").Core().V1().Nodes().Lister(),
		infrastructureLister:       configInformers.Config().V1().Infrastructures().Lister(),
	}
	return c.WithPodNetworkConnectivityCheckFn(generator.generate)
}

type openshiftAPIServerConnectivityCheckController struct {
	connectivitycheckcontroller.ConnectivityCheckController
}

type connectivityCheckTemplateProvider struct {
	operatorClient             v1helpers.OperatorClient
	operatorcontrolplaneClient *operatorcontrolplaneclient.Clientset
	endpointsLister            corev1listers.EndpointsLister
	serviceLister              corev1listers.ServiceLister
	podLister                  corev1listers.PodLister
	nodeLister                 corev1listers.NodeLister
	infrastructureLister       configv1listers.InfrastructureLister
}

func (c *connectivityCheckTemplateProvider) generate(ctx context.Context, syncContext factory.SyncContext) ([]*v1alpha1.PodNetworkConnectivityCheck, error) {
	return nil, nil
}

func (c *connectivityCheckTemplateProvider) getPodNetworkConnectivityChecks(ctx context.Context, operatorSpec *operatorv1.OperatorSpec, recorder events.Recorder) {

	var templates []*v1alpha1.PodNetworkConnectivityCheck
	// each storage endpoint
	templates = append(templates, c.getTemplatesForStorageChecks(recorder)...)
	// kas service IP
	templates = append(templates, c.getTemplatesForKubernetesServiceMonitorService(recorder)...)
	// kas default service IP
	templates = append(templates, c.getTemplatesForKubernetesDefaultServiceCheck(recorder)...)
	// each kas endpoint IP
	templates = append(templates, c.getTemplatesForKubernetesServiceEndpointsChecks(recorder)...)
	// each api load balancer hostname
	templates = append(templates, c.getTemplatesForApiLoadBalancerChecks(recorder)...)

	pods, err := c.podLister.List(labels.Set{"apiserver": "true"}.AsSelector())
	if err != nil {
		recorder.Warningf("EndpointDetectionFailure", "failed to list openshift-apiserver pods: %v", err)
		return
	}

	// create each check per static pod
	var checks []*v1alpha1.PodNetworkConnectivityCheck
	for _, pod := range pods {
		if len(pod.Spec.NodeName) == 0 {
			// apiserver pod hasn't been assigned a node yet, skip
			continue
		}
		for _, template := range templates {
			check := template.DeepCopy()
			WithSource("apiserver-" + pod.Spec.NodeName)(check)
			check.Spec.SourcePod = pod.Name
			checks = append(checks, check)
		}
	}

	pnccClient := c.operatorcontrolplaneClient.ControlplaneV1alpha1().PodNetworkConnectivityChecks(operatorclient.TargetNamespace)
	for _, check := range checks {
		existing, err := pnccClient.Get(ctx, check.Name, metav1.GetOptions{})
		if err == nil {
			if equality.Semantic.DeepEqual(existing.Spec, check.Spec) {
				// already exists, no changes, skip
				continue
			}
			updated := existing.DeepCopy()
			updated.Spec = *check.Spec.DeepCopy()
			_, err := pnccClient.Update(ctx, updated, metav1.UpdateOptions{})
			if err != nil {
				recorder.Warningf("EndpointDetectionFailure", "%s: %v", resourcehelper.FormatResourceForCLIWithNamespace(check), err)
				continue
			}
			recorder.Eventf("EndpointCheckUpdated", "Updated %s because it changed.", resourcehelper.FormatResourceForCLIWithNamespace(check))
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

func (c *connectivityCheckTemplateProvider) getTemplatesForKubernetesDefaultServiceCheck(recorder events.Recorder) []*v1alpha1.PodNetworkConnectivityCheck {
	var templates []*v1alpha1.PodNetworkConnectivityCheck
	host := os.Getenv("KUBERNETES_SERVICE_HOST")
	port := os.Getenv("KUBERNETES_SERVICE_PORT")
	if len(host) == 0 || len(port) == 0 {
		recorder.Warningf("EndpointDetectionFailure", "unable to determine kubernetes service endpoint: in-cluster configuration not found")
		return templates
	}
	return append(templates, NewPodNetworkConnectivityCheckTemplate(net.JoinHostPort(host, port), operatorclient.TargetNamespace, withTarget("kubernetes-default-service", "cluster")))
}

func (c *connectivityCheckTemplateProvider) getTemplatesForKubernetesServiceMonitorService(recorder events.Recorder) []*v1alpha1.PodNetworkConnectivityCheck {
	var templates []*v1alpha1.PodNetworkConnectivityCheck
	for _, address := range c.listAddressesForKubernetesServiceMonitorService(recorder) {
		templates = append(templates, NewPodNetworkConnectivityCheckTemplate(address, operatorclient.TargetNamespace, withTarget("kubernetes-apiserver-service", "cluster")))
	}
	return templates
}

func (c *connectivityCheckTemplateProvider) listAddressesForKubernetesServiceMonitorService(recorder events.Recorder) []string {
	service, err := c.serviceLister.Services("openshift-kube-apiserver").Get("apiserver")
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

func (c *connectivityCheckTemplateProvider) getTemplatesForKubernetesServiceEndpointsChecks(recorder events.Recorder) []*v1alpha1.PodNetworkConnectivityCheck {
	var templates []*v1alpha1.PodNetworkConnectivityCheck
	addresses, err := c.listAddressesForKubeAPIServerServiceEndpoints(recorder)
	if err != nil {
		recorder.Warningf("EndpointDetectionFailure", "unable to determine openshift-kube-apiserver apiserver endpoints: %v", err)
		return nil
	}

	for _, address := range addresses {
		templates = append(templates, NewPodNetworkConnectivityCheckTemplate(net.JoinHostPort(address.hostName, address.port), operatorclient.TargetNamespace, withTarget("kubernetes-apiserver-endpoint", address.nodeName)))
	}
	return templates
}

// listAddressesForKubeAPIServerServiceEndpoints returns kas api service endpoints ip
func (c *connectivityCheckTemplateProvider) listAddressesForKubeAPIServerServiceEndpoints(recorder events.Recorder) ([]endpointInfo, error) {
	var results []endpointInfo
	endpoints, err := c.endpointsLister.Endpoints("openshift-kube-apiserver").Get("apiserver")
	if err != nil {
		return nil, err
	}
	for _, subset := range endpoints.Subsets {
		for _, address := range subset.Addresses {
			for _, port := range subset.Ports {
				results = append(results, endpointInfo{
					hostName: address.IP,
					port:     strconv.Itoa(int(port.Port)),
					nodeName: *address.NodeName,
				})
			}
		}
	}
	return results, nil
}

func (c *connectivityCheckTemplateProvider) getTemplatesForStorageChecks(recorder events.Recorder) []*v1alpha1.PodNetworkConnectivityCheck {
	operatorSpec, _, _, err := c.operatorClient.GetOperatorState()
	if err != nil {
		recorder.Warningf("EndpointDetectionFailure", "unable to determine storage endpoints: %v", err)
		return nil
	}
	var templates []*v1alpha1.PodNetworkConnectivityCheck
	for _, endpointInfo := range c.listAddressesForStorageEndpoints(operatorSpec, recorder) {
		templates = append(templates, NewPodNetworkConnectivityCheckTemplate(
			net.JoinHostPort(endpointInfo.hostName, endpointInfo.port),
			operatorclient.TargetNamespace,
			withTarget("storage-endpoint", endpointInfo.nodeName),
			WithTlsClientCert("etcd-client")))
	}
	return templates
}

func (c *connectivityCheckTemplateProvider) listAddressesForStorageEndpoints(operatorSpec *operatorv1.OperatorSpec, recorder events.Recorder) []endpointInfo {
	var results []endpointInfo
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
		switch storageConfigURL.Hostname() {
		case "localhost", "127.0.0.1", "::1":
			results = append(results, endpointInfo{
				hostName: storageConfigURL.Hostname(),
				port:     storageConfigURL.Port(),
				nodeName: "localhost",
			})
			continue
		}
		node, err := c.findNodeForInternalIP(storageConfigURL.Hostname())
		if err != nil {
			recorder.Warningf("EndpointDetectionFailure", "unable to determine node for storage server: %v", err)
			continue
		}
		results = append(results, endpointInfo{
			hostName: storageConfigURL.Hostname(),
			port:     storageConfigURL.Port(),
			nodeName: node.Name,
		})
	}
	return results
}

func (c *connectivityCheckTemplateProvider) findNodeForInternalIP(internalIP string) (*corev1.Node, error) {
	nodes, err := c.nodeLister.List(labels.Everything())
	if err != nil {
		return nil, err
	}
	for _, node := range nodes {
		for _, nodeAddress := range node.Status.Addresses {
			if nodeAddress.Type != corev1.NodeInternalIP {
				continue
			}
			if internalIP == nodeAddress.Address {
				return node, nil
			}
		}
	}
	return nil, fmt.Errorf("no node found with internal IP %s", internalIP)
}

func (c *connectivityCheckTemplateProvider) getTemplatesForApiLoadBalancerChecks(recorder events.Recorder) []*v1alpha1.PodNetworkConnectivityCheck {
	var templates []*v1alpha1.PodNetworkConnectivityCheck
	infrastructure, err := c.infrastructureLister.Get("cluster")
	if err != nil {
		recorder.Warningf("EndpointDetectionFailure", "error detecting api load balancer endpoints: %v", err)
		return nil
	}

	apiUrl, err := url.Parse(infrastructure.Status.APIServerURL)
	if err != nil {
		recorder.Warningf("EndpointDetectionFailure", "error detecting external api load balancer endpoint: %v", err)

	} else {
		templates = append(templates, NewPodNetworkConnectivityCheckTemplate(apiUrl.Host, operatorclient.TargetNamespace, withTarget("load-balancer", "api-external")))
	}

	apiInternalUrl, err := url.Parse(infrastructure.Status.APIServerInternalURL)
	if err != nil {
		recorder.Warningf("EndpointDetectionFailure", "error detecting internal api load balancer endpoint: %v", err)
	} else {
		templates = append(templates, NewPodNetworkConnectivityCheckTemplate(apiInternalUrl.Host, operatorclient.TargetNamespace, withTarget("load-balancer", "api-internal")))
	}
	return templates
}

type endpointInfo struct {
	hostName string
	port     string
	nodeName string
}

func withTarget(label, nodeName string) func(check *v1alpha1.PodNetworkConnectivityCheck) {
	return WithTarget(label + "-" + nodeName)
}
