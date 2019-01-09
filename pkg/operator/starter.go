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

	configv1 "github.com/openshift/api/config/v1"
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

const (
	etcdNamespaceName                     = "kube-system"
	userSpecifiedGlobalConfigNamespace    = "openshift-config"
	machineSpecifiedGlobalConfigNamespace = "openshift-config-managed"
	kubeAPIServerNamespaceName            = "openshift-kube-apiserver"
	operatorNamespace                     = "openshift-apiserver-operator"
	targetNamespaceName                   = "openshift-apiserver"
	workQueueKey                          = "key"
)

func RunOperator(ctx *controllercmd.ControllerContext) error {
	kubeClient, err := kubernetes.NewForConfig(ctx.KubeConfig)
	if err != nil {
		return err
	}
	apiregistrationv1Client, err := apiregistrationclient.NewForConfig(ctx.KubeConfig)
	if err != nil {
		return err
	}
	operatorConfigClient, err := operatorconfigclient.NewForConfig(ctx.KubeConfig)
	if err != nil {
		return err
	}
	dynamicClient, err := dynamic.NewForConfig(ctx.KubeConfig)
	if err != nil {
		return err
	}
	configClient, err := configv1client.NewForConfig(ctx.KubeConfig)
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
	kubeInformersForOpenShiftConfigNamespace := kubeinformers.NewSharedInformerFactoryWithOptions(kubeClient, 10*time.Minute, kubeinformers.WithNamespace(userSpecifiedGlobalConfigNamespace))
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
		kubeInformersForOpenShiftConfigNamespace,
		apiregistrationInformers,
		configInformers,
		operatorConfigClient.OpenshiftapiserverV1alpha1(),
		configClient.ConfigV1(),
		kubeClient,
		apiregistrationv1Client.ApiregistrationV1(),
		ctx.EventRecorder,
	)
	finalizerController := NewFinalizerController(
		kubeInformersForOpenShiftAPIServerNamespace,
		kubeClient,
		ctx.EventRecorder,
	)

	configObserver := configobservercontroller.NewConfigObserver(
		operatorClient,
		operatorConfigInformers,
		kubeInformersForEtcdNamespace,
		configInformers,
		ctx.EventRecorder,
	)

	clusterOperatorStatus := status.NewClusterOperatorStatusController(
		"openshift-apiserver",
		[]configv1.ObjectReference{
			{Group: "openshiftapiserver.operator.openshift.io", Resource: "openshiftapiserveroperatorconfigs", Name: "cluster"},
			{Resource: "namespaces", Name: userSpecifiedGlobalConfigNamespace},
			{Resource: "namespaces", Name: machineSpecifiedGlobalConfigNamespace},
			{Resource: "namespaces", Name: operatorNamespace},
			{Resource: "namespaces", Name: targetNamespaceName},
		},
		configClient.ConfigV1(),
		operatorClient,
		ctx.EventRecorder,
	)

	operatorConfigInformers.Start(ctx.StopCh)
	kubeInformersForOpenShiftAPIServerNamespace.Start(ctx.StopCh)
	kubeInformersForKubeAPIServerNamespace.Start(ctx.StopCh)
	kubeInformersForEtcdNamespace.Start(ctx.StopCh)
	apiregistrationInformers.Start(ctx.StopCh)
	configInformers.Start(ctx.StopCh)

	go workloadController.Run(1, ctx.StopCh)
	go configObserver.Run(1, ctx.StopCh)
	go clusterOperatorStatus.Run(1, ctx.StopCh)
	go finalizerController.Run(1, ctx.StopCh)

	<-ctx.StopCh
	return fmt.Errorf("stopped")
}
