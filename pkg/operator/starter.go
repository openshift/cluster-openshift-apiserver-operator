package operator

import (
	"fmt"
	"os"
	"time"

	"github.com/openshift/library-go/pkg/controller/controllercmd"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	apiregistrationclient "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
	apiregistrationinformers "k8s.io/kube-aggregator/pkg/client/informers/externalversions"

	configv1client "github.com/openshift/client-go/config/clientset/versioned"
	configinformers "github.com/openshift/client-go/config/informers/externalversions"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/apis/openshiftapiserver/v1alpha1"
	operatorconfigclient "github.com/openshift/cluster-openshift-apiserver-operator/pkg/generated/clientset/versioned"
	operatorclientinformers "github.com/openshift/cluster-openshift-apiserver-operator/pkg/generated/informers/externalversions"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation/configobservercontroller"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/v311_00_assets"
	"github.com/openshift/library-go/pkg/operator/status"
	"github.com/openshift/library-go/pkg/operator/v1alpha1helpers"
)

func RunOperator(controllerContext *controllercmd.ControllerContext) error {
	kubeClient, err := kubernetes.NewForConfig(controllerContext.KubeConfig)
	if err != nil {
		return err
	}
	apiregistrationv1Client, err := apiregistrationclient.NewForConfig(controllerContext.KubeConfig)
	if err != nil {
		return err
	}
	operatorConfigClient, err := operatorconfigclient.NewForConfig(controllerContext.KubeConfig)
	if err != nil {
		return err
	}
	dynamicClient, err := dynamic.NewForConfig(controllerContext.KubeConfig)
	if err != nil {
		return err
	}
	configClient, err := configv1client.NewForConfig(controllerContext.KubeConfig)
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
	configInformers := configinformers.NewSharedInformerFactory(configClient, 10*time.Minute)

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
		controllerContext.EventRecorder,
	)

	configObserver := configobservercontroller.NewConfigObserver(
		operatorClient,
		operatorConfigInformers,
		kubeInformersForEtcdNamespace,
		configInformers,
		controllerContext.EventRecorder,
	)

	clusterOperatorStatus := status.NewClusterOperatorStatusController(
		"openshift-cluster-openshift-apiserver-operator",
		configClient.ConfigV1(),
		operatorClient,
		controllerContext.EventRecorder,
	)

	operatorConfigInformers.Start(controllerContext.StopCh)
	kubeInformersForOpenShiftAPIServerNamespace.Start(controllerContext.StopCh)
	kubeInformersForKubeAPIServerNamespace.Start(controllerContext.StopCh)
	kubeInformersForEtcdNamespace.Start(controllerContext.StopCh)
	apiregistrationInformers.Start(controllerContext.StopCh)
	configInformers.Start(controllerContext.StopCh)

	go workloadController.Run(1, controllerContext.StopCh)
	go configObserver.Run(1, controllerContext.StopCh)
	go clusterOperatorStatus.Run(1, controllerContext.StopCh)

	<-controllerContext.StopCh
	return fmt.Errorf("stopped")
}
