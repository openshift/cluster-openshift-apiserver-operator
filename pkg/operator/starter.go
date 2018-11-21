package operator

import (
	"fmt"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	apiregistrationclient "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
	apiregistrationinformers "k8s.io/kube-aggregator/pkg/client/informers/externalversions"

	configv1client "github.com/openshift/client-go/config/clientset/versioned"
	imageconfiginformers "github.com/openshift/client-go/config/informers/externalversions"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/apis/openshiftapiserver/v1alpha1"
	operatorconfigclient "github.com/openshift/cluster-openshift-apiserver-operator/pkg/generated/clientset/versioned"
	operatorclientinformers "github.com/openshift/cluster-openshift-apiserver-operator/pkg/generated/informers/externalversions"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation/configobservercontroller"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/v311_00_assets"
	"github.com/openshift/library-go/pkg/operator/status"
	"github.com/openshift/library-go/pkg/operator/v1alpha1helpers"
)

func RunOperator(clientConfig *rest.Config, stopCh <-chan struct{}) error {
	kubeClient, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return err
	}
	apiregistrationv1Client, err := apiregistrationclient.NewForConfig(clientConfig)
	if err != nil {
		return err
	}
	operatorConfigClient, err := operatorconfigclient.NewForConfig(clientConfig)
	if err != nil {
		return err
	}
	dynamicClient, err := dynamic.NewForConfig(clientConfig)
	if err != nil {
		return err
	}
	configClient, err := configv1client.NewForConfig(clientConfig)
	if err != nil {
		return err
	}

	v1alpha1helpers.EnsureOperatorConfigExists(
		dynamicClient,
		v311_00_assets.MustAsset("v3.11.0/openshift-apiserver/operator-config.yaml"),
		schema.GroupVersionResource{Group: v1alpha1.GroupName, Version: "v1alpha1", Resource: "openshiftapiserveroperatorconfigs"},
		v1alpha1helpers.GetImageEnv,
	)

	operatorConfigInformers := operatorclientinformers.NewSharedInformerFactory(operatorConfigClient, 10*time.Minute)
	kubeInformersForOpenShiftAPIServerNamespace := kubeinformers.NewSharedInformerFactoryWithOptions(kubeClient, 10*time.Minute, kubeinformers.WithNamespace(targetNamespaceName))
	kubeInformersForKubeAPIServerNamespace := kubeinformers.NewSharedInformerFactoryWithOptions(kubeClient, 10*time.Minute, kubeinformers.WithNamespace(kubeAPIServerNamespaceName))
	kubeInformersForEtcdNamespace := kubeinformers.NewSharedInformerFactoryWithOptions(kubeClient, 10*time.Minute, kubeinformers.WithNamespace(etcdNamespaceName))
	apiregistrationInformers := apiregistrationinformers.NewSharedInformerFactory(apiregistrationv1Client, 10*time.Minute)
	imageConfigInformers := imageconfiginformers.NewSharedInformerFactory(configClient, 10*time.Minute)

	operatorClient := &operatorClient{
		informers: operatorConfigInformers,
		client:    operatorConfigClient.OpenshiftapiserverV1alpha1(),
	}

	workloadController := NewWorkloadController(
		os.Getenv("IMAGE"),
		operatorConfigInformers.Openshiftapiserver().V1alpha1().OpenShiftAPIServerOperatorConfigs(),
		kubeInformersForOpenShiftAPIServerNamespace,
		kubeInformersForEtcdNamespace,
		kubeInformersForKubeAPIServerNamespace,
		apiregistrationInformers,
		operatorConfigClient.OpenshiftapiserverV1alpha1(),
		kubeClient,
		apiregistrationv1Client.ApiregistrationV1(),
	)

	configObserver := configobservercontroller.NewConfigObserver(
		operatorClient,
		operatorConfigInformers,
		kubeInformersForEtcdNamespace,
		imageConfigInformers,
	)

	clusterOperatorStatus := status.NewClusterOperatorStatusController(
		"openshift-cluster-openshift-apiserver-operator",
		"openshift-cluster-openshift-apiserver-operator",
		dynamicClient,
		operatorClient,
	)

	operatorConfigInformers.Start(stopCh)
	kubeInformersForOpenShiftAPIServerNamespace.Start(stopCh)
	kubeInformersForKubeAPIServerNamespace.Start(stopCh)
	kubeInformersForEtcdNamespace.Start(stopCh)
	apiregistrationInformers.Start(stopCh)
	imageConfigInformers.Start(stopCh)

	go workloadController.Run(1, stopCh)
	go configObserver.Run(1, stopCh)
	go clusterOperatorStatus.Run(1, stopCh)

	<-stopCh
	return fmt.Errorf("stopped")
}
