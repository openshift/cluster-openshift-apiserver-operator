package operator

import (
	"fmt"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	"github.com/openshift/library-go/pkg/operator/unsupportedconfigoverridescontroller"
	"github.com/openshift/library-go/pkg/operator/v1helpers"
)

func RunOperator(ctx *controllercmd.ControllerContext) error {
	kubeClient, err := kubernetes.NewForConfig(ctx.ProtoKubeConfig)
	if err != nil {
		return err
	}
	apiregistrationv1Client, err := apiregistrationclient.NewForConfig(ctx.ProtoKubeConfig)
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
		operatorclient.GlobalUserSpecifiedConfigNamespace,
		operatorclient.GlobalMachineSpecifiedConfigNamespace,
		operatorclient.OperatorNamespace,
		operatorclient.TargetNamespace,
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
		v1helpers.CachedConfigMapGetter(kubeClient.CoreV1(), kubeInformersForNamespaces),
		v1helpers.CachedSecretGetter(kubeClient.CoreV1(), kubeInformersForNamespaces),
		ctx.EventRecorder,
	)
	if err != nil {
		return err
	}

	// don't change any versions until we sync
	versionRecorder := status.NewVersionGetter()
	clusterOperator, err := configClient.ConfigV1().ClusterOperators().Get("openshift-apiserver", metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	for _, version := range clusterOperator.Status.Versions {
		versionRecorder.SetVersion(version.Name, version.Version)
	}
	versionRecorder.SetVersion("operator", os.Getenv("OPERATOR_IMAGE_VERSION"))

	workloadController := workloadcontroller.NewWorkloadController(
		os.Getenv("IMAGE"),
		versionRecorder,
		operatorConfigInformers.Operator().V1().OpenShiftAPIServers(),
		kubeInformersForNamespaces.InformersFor(operatorclient.TargetNamespace),
		kubeInformersForNamespaces.InformersFor(operatorclient.EtcdNamespaceName),
		kubeInformersForNamespaces.InformersFor(operatorclient.GlobalUserSpecifiedConfigNamespace),
		apiregistrationInformers,
		configInformers,
		operatorConfigClient.OperatorV1(),
		configClient.ConfigV1(),
		kubeClient,
		apiregistrationv1Client.ApiregistrationV1(),
		ctx.EventRecorder,
	)
	finalizerController := NewFinalizerController(
		kubeInformersForNamespaces.InformersFor(operatorclient.TargetNamespace),
		kubeClient.CoreV1(),
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
				{Resource: "namespaces", Name: operatorclient.GlobalUserSpecifiedConfigNamespace},
				{Resource: "namespaces", Name: operatorclient.GlobalMachineSpecifiedConfigNamespace},
				{Resource: "namespaces", Name: operatorclient.OperatorNamespace},
				{Resource: "namespaces", Name: operatorclient.TargetNamespace},
			},
			workloadcontroller.APIServiceReferences()...,
		),
		configClient.ConfigV1(),
		operatorClient,
		versionRecorder,
		ctx.EventRecorder,
	)

	configUpgradeableController := unsupportedconfigoverridescontroller.NewUnsupportedConfigOverridesController(operatorClient, ctx.EventRecorder)

	operatorConfigInformers.Start(ctx.Done())
	kubeInformersForNamespaces.Start(ctx.Done())
	apiregistrationInformers.Start(ctx.Done())
	configInformers.Start(ctx.Done())

	go workloadController.Run(1, ctx.Done())
	go configObserver.Run(1, ctx.Done())
	go clusterOperatorStatus.Run(1, ctx.Done())
	go finalizerController.Run(1, ctx.Done())
	go resourceSyncController.Run(1, ctx.Done())
	go configUpgradeableController.Run(1, ctx.Done())

	<-ctx.Done()
	return fmt.Errorf("stopped")
}
