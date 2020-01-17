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

	"github.com/openshift/library-go/pkg/controller/controllercmd"
	"github.com/openshift/library-go/pkg/operator/encryption"
	"github.com/openshift/library-go/pkg/operator/encryption/controllers/migrators"
	encryptiondeployer "github.com/openshift/library-go/pkg/operator/encryption/deployer"
	"github.com/openshift/library-go/pkg/operator/genericoperatorclient"
	"github.com/openshift/library-go/pkg/operator/loglevel"
	"github.com/openshift/library-go/pkg/operator/revisioncontroller"
	"github.com/openshift/library-go/pkg/operator/status"
	"github.com/openshift/library-go/pkg/operator/unsupportedconfigoverridescontroller"
	"github.com/openshift/library-go/pkg/operator/v1helpers"

	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation/configobservercontroller"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/daemonsetstatuscontroller"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"
	prune "github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/prunecontroller"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/resourcesynccontroller"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/workloadcontroller"
	"github.com/openshift/library-go/pkg/operator/staticpod/controller/revision"
)

func RunOperator(ctx *controllercmd.ControllerContext) error {
	kubeClient, err := kubernetes.NewForConfig(ctx.ProtoKubeConfig)
	if err != nil {
		return err
	}
	migrationClientConfig := dynamic.ConfigFor(ctx.KubeConfig)
	migrationClientConfig.Burst = 40
	migrationClientConfig.QPS = 30
	dynamicClientForMigration, err := dynamic.NewForConfig(migrationClientConfig)
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
	configClient, err := configv1client.NewForConfig(ctx.KubeConfig)
	if err != nil {
		return err
	}

	operatorConfigInformers := operatorv1informers.NewSharedInformerFactory(operatorConfigClient, 10*time.Minute)
	kubeInformersForNamespaces := v1helpers.NewKubeInformersForNamespaces(kubeClient,
		"",
		operatorclient.GlobalUserSpecifiedConfigNamespace,
		operatorclient.GlobalMachineSpecifiedConfigNamespace,
		operatorclient.OperatorNamespace,
		operatorclient.TargetNamespace,
	)
	apiregistrationInformers := apiregistrationinformers.NewSharedInformerFactory(apiregistrationv1Client, 10*time.Minute)
	configInformers := configinformers.NewSharedInformerFactory(configClient, 10*time.Minute)

	operatorClient, dynamicInformers, err := genericoperatorclient.NewClusterScopedOperatorClient(ctx.KubeConfig, operatorv1.GroupVersion.WithResource("openshiftapiservers"))
	if err != nil {
		return err
	}

	resourceSyncController, debugHandler, err := resourcesynccontroller.NewResourceSyncController(
		operatorClient,
		kubeInformersForNamespaces,
		v1helpers.CachedConfigMapGetter(kubeClient.CoreV1(), kubeInformersForNamespaces),
		v1helpers.CachedSecretGetter(kubeClient.CoreV1(), kubeInformersForNamespaces),
		ctx.EventRecorder,
	)
	if err != nil {
		return err
	}

	revisionController := revisioncontroller.NewRevisionController(
		operatorclient.TargetNamespace,
		nil,
		[]revision.RevisionResource{{
			Name:     "encryption-config",
			Optional: true,
		}},
		kubeInformersForNamespaces.InformersFor(operatorclient.TargetNamespace),
		OpenshiftDeploymentLatestRevisionClient{OperatorClient: operatorClient, TypedClient: operatorConfigClient.OperatorV1()},
		v1helpers.CachedConfigMapGetter(kubeClient.CoreV1(), kubeInformersForNamespaces),
		v1helpers.CachedSecretGetter(kubeClient.CoreV1(), kubeInformersForNamespaces),
		ctx.EventRecorder,
	)

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
		os.Getenv("IMAGE"), os.Getenv("OPERATOR_IMAGE"),
		versionRecorder,
		operatorConfigInformers.Operator().V1().OpenShiftAPIServers(),
		kubeInformersForNamespaces.InformersFor(operatorclient.TargetNamespace),
		kubeInformersForNamespaces.InformersFor(operatorclient.GlobalUserSpecifiedConfigNamespace),
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
		kubeInformersForNamespaces.InformersFor(operatorclient.TargetNamespace),
		operatorClient,
		resourceSyncController,
		operatorConfigInformers,
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
		configInformers.Config().V1().ClusterOperators(),
		operatorClient,
		versionRecorder,
		ctx.EventRecorder,
	)

	nodeProvider := DaemonSetNodeProvider{
		TargetNamespaceDaemonSetInformer: kubeInformersForNamespaces.InformersFor(operatorclient.TargetNamespace).Apps().V1().DaemonSets(),
		NodeInformer:                     kubeInformersForNamespaces.InformersFor("").Core().V1().Nodes(),
	}
	deployer, err := encryptiondeployer.NewRevisionLabelPodDeployer("revision", operatorclient.TargetNamespace, kubeInformersForNamespaces, resourceSyncController, kubeClient.CoreV1(), kubeClient.CoreV1(), nodeProvider)
	if err != nil {
		return err
	}
	migrator := migrators.NewInProcessMigrator(dynamicClientForMigration, kubeClient.Discovery())

	encryptionControllers, err := encryption.NewControllers(
		operatorclient.TargetNamespace,
		deployer,
		migrator,
		operatorClient,
		configClient.ConfigV1().APIServers(),
		configInformers.Config().V1().APIServers(),
		kubeInformersForNamespaces,
		kubeClient.CoreV1(),
		ctx.EventRecorder,
		schema.GroupResource{Group: "route.openshift.io", Resource: "routes"}, // routes can contain embedded TLS private keys
		schema.GroupResource{Group: "oauth.openshift.io", Resource: "oauthaccesstokens"},
		schema.GroupResource{Group: "oauth.openshift.io", Resource: "oauthauthorizetokens"},
	)
	if err != nil {
		return err
	}

	pruneController := prune.NewPruneController(
		operatorclient.TargetNamespace,
		[]string{"encryption-config-"},
		kubeClient.CoreV1(),
		kubeClient.CoreV1(),
		kubeInformersForNamespaces,
		ctx.EventRecorder,
	)

	daemonsetStatusController := daemonsetstatuscontroller.NewDaemonSetStatusController(
		operatorclient.TargetNamespace, "apiserver",
		kubeInformersForNamespaces.InformersFor(operatorclient.TargetNamespace),
		operatorClient,
		ctx.EventRecorder,
	)

	configUpgradeableController := unsupportedconfigoverridescontroller.NewUnsupportedConfigOverridesController(operatorClient, ctx.EventRecorder)
	logLevelController := loglevel.NewClusterOperatorLoggingController(operatorClient, ctx.EventRecorder)

	if ctx.Server != nil {
		ctx.Server.Handler.NonGoRestfulMux.Handle("/debug/controllers/resourcesync", debugHandler)
	}

	operatorConfigInformers.Start(ctx.Done())
	kubeInformersForNamespaces.Start(ctx.Done())
	apiregistrationInformers.Start(ctx.Done())
	configInformers.Start(ctx.Done())
	dynamicInformers.Start(ctx.Done())

	go workloadController.Run(1, ctx.Done())
	go configObserver.Run(1, ctx.Done())
	go clusterOperatorStatus.Run(1, ctx.Done())
	go finalizerController.Run(1, ctx.Done())
	go resourceSyncController.Run(1, ctx.Done())
	go revisionController.Run(1, ctx.Done())
	go encryptionControllers.Run(ctx.Done())
	go pruneController.Run(1, ctx.Done())
	go daemonsetStatusController.Run(1, ctx.Done())
	go configUpgradeableController.Run(1, ctx.Done())
	go logLevelController.Run(1, ctx.Done())

	<-ctx.Done()
	return fmt.Errorf("stopped")
}
