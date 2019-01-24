package operator

import (
	"fmt"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	apiregistrationclient "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
	apiregistrationinformers "k8s.io/kube-aggregator/pkg/client/informers/externalversions"

	configv1 "github.com/openshift/api/config/v1"
	operatorv1 "github.com/openshift/api/operator/v1"
	configv1client "github.com/openshift/client-go/config/clientset/versioned"
	configinformers "github.com/openshift/client-go/config/informers/externalversions"
	operatorv1client "github.com/openshift/client-go/operator/clientset/versioned"
	operatorv1informers "github.com/openshift/client-go/operator/informers/externalversions"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation/configobservercontroller"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/resourcesynccontroller"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/v311_00_assets"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/workloadcontroller"
	"github.com/openshift/library-go/pkg/controller/controllercmd"
	"github.com/openshift/library-go/pkg/operator/status"
	"github.com/openshift/library-go/pkg/operator/v1helpers"
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
	operatorConfigClient, err := operatorv1client.NewForConfig(ctx.KubeConfig)
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

	v1helpers.EnsureOperatorConfigExists(
		dynamicClient,
		v311_00_assets.MustAsset("v3.11.0/openshift-apiserver/operator-config.yaml"),
		schema.GroupVersionResource{Group: operatorv1.GroupName, Version: operatorv1.GroupVersion.Version, Resource: "openshiftapiservers"},
	)

	operatorConfigInformers := operatorv1informers.NewSharedInformerFactory(operatorConfigClient, 10*time.Minute)
	kubeInformersForNamespaces := v1helpers.NewKubeInformersForNamespaces(kubeClient,
		"",
		operatorclient.UserSpecifiedGlobalConfigNamespace,
		operatorclient.MachineSpecifiedGlobalConfigNamespace,
		operatorclient.KubeAPIServerNamespaceName,
		operatorclient.OperatorNamespace,
		operatorclient.TargetNamespaceName,
		"kube-system",
	)
	apiregistrationInformers := apiregistrationinformers.NewSharedInformerFactory(apiregistrationv1Client, 10*time.Minute)
	configInformers := configinformers.NewSharedInformerFactory(configClient, 10*time.Minute)

	operatorClient := &operatorclient.OperatorClient{
		Informers: operatorConfigInformers,
		Client:    operatorConfigClient.OperatorV1(),
	}

	resourceSyncController, err := resourcesynccontroller.NewResourceSyncController(
		operatorClient,
		kubeInformersForNamespaces,
		kubeClient,
		ctx.EventRecorder,
	)
	if err != nil {
		return err
	}

	workloadController := workloadcontroller.NewWorkloadController(
		os.Getenv("IMAGE"),
		operatorConfigInformers.Operator().V1().OpenShiftAPIServers(),
		kubeInformersForNamespaces.InformersFor(operatorclient.TargetNamespaceName),
		kubeInformersForNamespaces.InformersFor(operatorclient.EtcdNamespaceName),
		kubeInformersForNamespaces.InformersFor(operatorclient.KubeAPIServerNamespaceName),
		kubeInformersForNamespaces.InformersFor(operatorclient.UserSpecifiedGlobalConfigNamespace),
		apiregistrationInformers,
		configInformers,
		operatorConfigClient.OperatorV1(),
		configClient.ConfigV1(),
		kubeClient,
		apiregistrationv1Client.ApiregistrationV1(),
		ctx.EventRecorder,
	)
	finalizerController := NewFinalizerController(
		kubeInformersForNamespaces.InformersFor(operatorclient.TargetNamespaceName),
		kubeClient,
		ctx.EventRecorder,
	)

	configObserver := configobservercontroller.NewConfigObserver(
		operatorClient,
		resourceSyncController,
		operatorConfigInformers,
		kubeInformersForNamespaces.InformersFor("kube-system"),
		configInformers,
		ctx.EventRecorder,
	)

	clusterOperatorStatus := status.NewClusterOperatorStatusController(
		"openshift-apiserver",
		append(
			[]configv1.ObjectReference{
				{Group: "operator.openshift.io", Resource: "openshiftapiservers", Name: "cluster"},
				{Resource: "namespaces", Name: operatorclient.UserSpecifiedGlobalConfigNamespace},
				{Resource: "namespaces", Name: operatorclient.MachineSpecifiedGlobalConfigNamespace},
				{Resource: "namespaces", Name: operatorclient.OperatorNamespace},
				{Resource: "namespaces", Name: operatorclient.TargetNamespaceName},
			},
			workloadcontroller.APIServiceReferences()...,
		),
		configClient.ConfigV1(),
		operatorClient,
		status.NewVersionGetter(),
		ctx.EventRecorder,
	)

	operatorConfigInformers.Start(ctx.StopCh)
	kubeInformersForNamespaces.Start(ctx.StopCh)
	apiregistrationInformers.Start(ctx.StopCh)
	configInformers.Start(ctx.StopCh)

	go workloadController.Run(1, ctx.StopCh)
	go configObserver.Run(1, ctx.StopCh)
	go clusterOperatorStatus.Run(1, ctx.StopCh)
	go finalizerController.Run(1, ctx.StopCh)
	go resourceSyncController.Run(1, ctx.StopCh)

	<-ctx.StopCh
	return fmt.Errorf("stopped")
}
