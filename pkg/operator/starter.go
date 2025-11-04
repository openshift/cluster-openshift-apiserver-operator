package operator

import (
	"context"
	"fmt"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/runtime"

	configv1 "github.com/openshift/api/config/v1"
	operatorv1 "github.com/openshift/api/operator/v1"
	configv1client "github.com/openshift/client-go/config/clientset/versioned"
	configinformers "github.com/openshift/client-go/config/informers/externalversions"
	configlisterv1 "github.com/openshift/client-go/config/listers/config/v1"
	applyoperatorv1 "github.com/openshift/client-go/operator/applyconfigurations/operator/v1"
	operatorv1client "github.com/openshift/client-go/operator/clientset/versioned"
	operatorv1informers "github.com/openshift/client-go/operator/informers/externalversions"
	operatorcontrolplaneclient "github.com/openshift/client-go/operatorcontrolplane/clientset/versioned"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation/configobservercontroller"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/connectivitycheckcontroller"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/resourcesynccontroller"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/v311_00_assets"
	operatorworkload "github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/workload"
	"github.com/openshift/library-go/pkg/controller/controllercmd"
	workloadcontroller "github.com/openshift/library-go/pkg/operator/apiserver/controller/workload"
	apiservercontrollerset "github.com/openshift/library-go/pkg/operator/apiserver/controllerset"
	libgoetcd "github.com/openshift/library-go/pkg/operator/configobserver/etcd"
	"github.com/openshift/library-go/pkg/operator/configobserver/featuregates"
	"github.com/openshift/library-go/pkg/operator/encryption"
	"github.com/openshift/library-go/pkg/operator/encryption/controllers/migrators"
	encryptiondeployer "github.com/openshift/library-go/pkg/operator/encryption/deployer"
	"github.com/openshift/library-go/pkg/operator/genericoperatorclient"
	"github.com/openshift/library-go/pkg/operator/staleconditions"
	staticpodcommon "github.com/openshift/library-go/pkg/operator/staticpod/controller/common"
	"github.com/openshift/library-go/pkg/operator/staticpod/controller/revision"
	"github.com/openshift/library-go/pkg/operator/status"
	"github.com/openshift/library-go/pkg/operator/v1helpers"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apiextensionsinformers "k8s.io/apiextensions-apiserver/pkg/client/informers/externalversions"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	apiregistrationclient "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
	apiregistrationinformers "k8s.io/kube-aggregator/pkg/client/informers/externalversions"
	utilpointer "k8s.io/utils/pointer"
	kubemigratorclient "sigs.k8s.io/kube-storage-version-migrator/pkg/clients/clientset"
	migrationv1alpha1informer "sigs.k8s.io/kube-storage-version-migrator/pkg/clients/informer"
)

const (
	oauthAPIServerTargetNamespace = "openshift-oauth-apiserver"
)

var apiServiceGroupVersions = []schema.GroupVersion{
	// these are all the apigroups we manage
	{Group: "apps.openshift.io", Version: "v1"},
	{Group: "authorization.openshift.io", Version: "v1"},
	{Group: "build.openshift.io", Version: "v1"},
	{Group: "image.openshift.io", Version: "v1"},
	{Group: "project.openshift.io", Version: "v1"},
	{Group: "quota.openshift.io", Version: "v1"},
	{Group: "route.openshift.io", Version: "v1"},
	{Group: "security.openshift.io", Version: "v1"},
	{Group: "template.openshift.io", Version: "v1"},
}

func RunOperator(ctx context.Context, controllerConfig *controllercmd.ControllerContext) error {
	kubeClient, err := kubernetes.NewForConfig(controllerConfig.ProtoKubeConfig)
	if err != nil {
		return err
	}
	apiregistrationv1Client, err := apiregistrationclient.NewForConfig(controllerConfig.ProtoKubeConfig)
	if err != nil {
		return err
	}
	operatorConfigClient, err := operatorv1client.NewForConfig(controllerConfig.KubeConfig)
	if err != nil {
		return err
	}
	configClient, err := configv1client.NewForConfig(controllerConfig.KubeConfig)
	if err != nil {
		return err
	}
	operatorcontrolplaneClient, err := operatorcontrolplaneclient.NewForConfig(controllerConfig.KubeConfig)
	if err != nil {
		return err
	}
	apiextensionsClient, err := apiextensionsclient.NewForConfig(controllerConfig.KubeConfig)
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
		libgoetcd.EtcdEndpointNamespace,
		metav1.NamespaceSystem,
		"openshift-kube-apiserver",
		oauthAPIServerTargetNamespace,
	)
	apiregistrationInformers := apiregistrationinformers.NewSharedInformerFactory(apiregistrationv1Client, 10*time.Minute)
	configInformers := configinformers.NewSharedInformerFactory(configClient, 10*time.Minute)

	operatorClient, dynamicInformers, err := genericoperatorclient.NewClusterScopedOperatorClient(
		controllerConfig.Clock,
		controllerConfig.KubeConfig,
		operatorv1.GroupVersion.WithResource("openshiftapiservers"),
		operatorv1.GroupVersion.WithKind("OpenShiftAPIServer"),
		extractOperatorSpec,
		extractOperatorStatus,
	)
	if err != nil {
		return err
	}

	desiredVersion := status.VersionForOperatorFromEnv()
	missingVersion := "0.0.1-snapshot"

	// By default, this will exit(0) the process if the featuregates ever change to a different set of values.
	featureGateAccessor := featuregates.NewFeatureGateAccess(
		desiredVersion, missingVersion,
		configInformers.Config().V1().ClusterVersions(), configInformers.Config().V1().FeatureGates(),
		controllerConfig.EventRecorder,
	)
	go featureGateAccessor.Run(ctx)
	// this second start is intentional.  This must start in order to sync the clusterversion and if it isn't started later
	// then any additional resources being watched won't have their informers started.
	configInformers.Start(ctx.Done())

	select {
	case <-featureGateAccessor.InitialFeatureGatesObserved():
		featureGates, _ := featureGateAccessor.CurrentFeatureGates()
		klog.Infof("FeatureGates initialized: knownFeatureGates=%v", featureGates.KnownFeatures())
	case <-time.After(1 * time.Minute):
		klog.Errorf("timed out waiting for FeatureGate detection")
		return fmt.Errorf("timed out waiting for FeatureGate detection")
	}

	resourceSyncController, debugHandler, err := resourcesynccontroller.NewResourceSyncController(
		operatorClient,
		kubeInformersForNamespaces,
		v1helpers.CachedConfigMapGetter(kubeClient.CoreV1(), kubeInformersForNamespaces),
		v1helpers.CachedSecretGetter(kubeClient.CoreV1(), kubeInformersForNamespaces),
		controllerConfig.EventRecorder,
	)
	if err != nil {
		return err
	}

	// don't change any versions until we sync
	versionRecorder := status.NewVersionGetter()
	clusterOperator, err := configClient.ConfigV1().ClusterOperators().Get(ctx, "openshift-apiserver", metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	for _, version := range clusterOperator.Status.Versions {
		versionRecorder.SetVersion(version.Name, version.Version)
	}
	versionRecorder.SetVersion("operator", os.Getenv("OPERATOR_IMAGE_VERSION"))

	openshiftNodeProvider := encryptiondeployer.NewDeploymentNodeProvider(operatorclient.TargetNamespace, kubeInformersForNamespaces)
	openshiftDeployer, err := encryptiondeployer.NewRevisionLabelPodDeployer("revision", operatorclient.TargetNamespace, kubeInformersForNamespaces, kubeClient.CoreV1(), kubeClient.CoreV1(), openshiftNodeProvider)
	if err != nil {
		return err
	}
	migrationClient := kubemigratorclient.NewForConfigOrDie(controllerConfig.KubeConfig)
	migrationInformer := migrationv1alpha1informer.NewSharedInformerFactory(migrationClient, time.Minute*30)
	migrator := migrators.NewKubeStorageVersionMigrator(migrationClient, migrationInformer.Migration().V1alpha1(), kubeClient.Discovery())

	openShiftAPIServerWorkload := operatorworkload.NewOpenShiftAPIServerWorkload(
		operatorClient,
		operatorConfigClient.OperatorV1(),
		configClient.ConfigV1(),
		configInformers.Config().V1().ClusterVersions().Lister(),
		configInformers.Config().V1().APIServers().Lister(),
		kubeInformersForNamespaces.InformersFor(operatorclient.TargetNamespace).Core().V1().Secrets().Lister(),
		workloadcontroller.CountNodesFuncWrapper(kubeInformersForNamespaces.InformersFor("").Core().V1().Nodes().Lister()),
		workloadcontroller.EnsureAtMostOnePodPerNode,
		"openshift-apiserver",
		os.Getenv("IMAGE"),
		os.Getenv("OPERATOR_IMAGE"),
		os.Getenv("KMS_PLUGIN_IMAGE"),
		kubeClient,
		versionRecorder)

	infra, err := configClient.ConfigV1().Infrastructures().Get(ctx, "cluster", metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		klog.Warningf("unexpectedly no infrastructure resource found, assuming non SingleReplicaTopologyMode controlPlaneTopology: %v", err)
	} else if err != nil {
		return err
	}
	var statusControllerOptions []func(*status.StatusSyncer) *status.StatusSyncer
	if infra == nil || infra.Status.ControlPlaneTopology != configv1.SingleReplicaTopologyMode {
		statusControllerOptions = append(statusControllerOptions, apiservercontrollerset.WithStatusControllerPdbCompatibleHighInertia("APIServer"))
	}

	apiServerControllers := apiservercontrollerset.NewAPIServerControllerSet(
		"openshift-apiserver",
		operatorClient,
		controllerConfig.EventRecorder,
		controllerConfig.Clock,
	).WithAPIServiceController(
		"openshift-apiserver",
		operatorclient.TargetNamespace,
		func() ([]*apiregistrationv1.APIService, []*apiregistrationv1.APIService, error) {
			return apiServices(configInformers.Config().V1().ClusterVersions().Lister())
		},
		apiregistrationInformers,
		apiregistrationv1Client.ApiregistrationV1(),
		kubeInformersForNamespaces,
		kubeClient,
		configInformers.Config().V1().ClusterVersions().Informer(),
	).WithFinalizerController(
		operatorclient.TargetNamespace,
		kubeInformersForNamespaces.InformersFor(operatorclient.TargetNamespace),
		kubeClient.CoreV1(),
	).WithClusterOperatorStatusController(
		"openshift-apiserver",
		append(
			[]configv1.ObjectReference{
				{Group: "operator.openshift.io", Resource: "openshiftapiservers", Name: "cluster"},
				{Resource: "namespaces", Name: operatorclient.GlobalUserSpecifiedConfigNamespace},
				{Resource: "namespaces", Name: operatorclient.GlobalMachineSpecifiedConfigNamespace},
				{Resource: "namespaces", Name: operatorclient.OperatorNamespace},
				{Resource: "namespaces", Name: operatorclient.TargetNamespace},
				{Resource: "namespaces", Name: "openshift-etcd-operator"}, // Capture events from etcd operator
				{Resource: "endpoints", Name: libgoetcd.EtcdEndpointName, Namespace: libgoetcd.EtcdEndpointNamespace},
				{Group: "controlplane.operator.openshift.io", Resource: "podnetworkconnectivitychecks", Namespace: "openshift-apiserver"},
			},
			apiServicesReferences()...,
		),
		configClient.ConfigV1(),
		configInformers.Config().V1().ClusterOperators(),
		versionRecorder,
		statusControllerOptions...,
	).WithWorkloadController(
		"OpenShiftAPIServer",
		operatorclient.OperatorNamespace,
		operatorclient.TargetNamespace,
		os.Getenv("OPERATOR_IMAGE_VERSION"),
		"openshift",
		"APIServer",
		kubeClient,
		openShiftAPIServerWorkload,
		versionRecorder,
		kubeInformersForNamespaces,
		operatorConfigInformers.Operator().V1().OpenShiftAPIServers().Informer(),
		configInformers.Config().V1().Images().Informer(),
		configInformers.Config().V1().ClusterVersions().Informer(),
	).WithStaticResourcesController(
		"APIServerStaticResources",
		v311_00_assets.Asset,
		[]apiservercontrollerset.ConditionalFiles{
			{
				Files: []string{
					"v3.11.0/openshift-apiserver/ns.yaml",
					"v3.11.0/openshift-apiserver/apiserver-clusterrolebinding.yaml",
					"v3.11.0/openshift-apiserver/svc.yaml",
					"v3.11.0/openshift-apiserver/sa.yaml",
					"v3.11.0/openshift-apiserver/trusted_ca_cm.yaml",
				},
			},
			{
				Files: []string{
					"v3.11.0/openshift-apiserver/pdb.yaml",
				},
				ShouldCreateFn: func() bool {
					isSNO, precheckSucceeded, err := staticpodcommon.NewIsSingleNodePlatformFn(configInformers.Config().V1().Infrastructures())()
					if err != nil {
						klog.Errorf("IsSNOCheckFnc failed: %v", err)
						return false
					}
					if !precheckSucceeded {
						klog.V(4).Infof("IsSNOCheckFnc precheck did not succeed, skipping")
						return false
					}
					return !isSNO
				},
				ShouldDeleteFn: func() bool {
					isSNO, precheckSucceeded, err := staticpodcommon.NewIsSingleNodePlatformFn(configInformers.Config().V1().Infrastructures())()
					if err != nil {
						klog.Errorf("IsSNOCheckFnc failed: %v", err)
						return false
					}
					if !precheckSucceeded {
						klog.V(4).Infof("IsSNOCheckFnc precheck did not succeed, skipping")
						return false
					}
					return isSNO
				},
			},
		},
		kubeInformersForNamespaces,
		kubeClient,
	).WithRevisionController(
		operatorclient.TargetNamespace,
		[]revision.RevisionResource{
			{
				Name: "audit",
			},
		},
		[]revision.RevisionResource{
			{
				Name:     "encryption-config",
				Optional: true,
			},
		},
		kubeInformersForNamespaces.InformersFor(operatorclient.TargetNamespace),
		operatorClient,
		v1helpers.CachedConfigMapGetter(kubeClient.CoreV1(), kubeInformersForNamespaces),
		v1helpers.CachedSecretGetter(kubeClient.CoreV1(), kubeInformersForNamespaces),
	).WithEncryptionControllers(
		operatorclient.TargetNamespace,
		encryption.StaticEncryptionProvider{
			schema.GroupResource{Group: "route.openshift.io", Resource: "routes"}, // routes can contain embedded TLS private keys
		},
		openshiftDeployer,
		migrator,
		kubeClient.CoreV1(),
		configClient.ConfigV1().APIServers(),
		configInformers.Config().V1().APIServers(),
		kubeInformersForNamespaces,
		resourceSyncController,
	).WithSecretRevisionPruneController(
		operatorclient.TargetNamespace,
		[]string{"encryption-config-"},
		kubeClient.CoreV1(),
		kubeClient.CoreV1(),
		kubeInformersForNamespaces,
	).WithAuditPolicyController(
		operatorclient.TargetNamespace,
		"audit",
		configInformers,
		kubeInformersForNamespaces.InformersFor(operatorclient.TargetNamespace),
		kubeClient,
	).
		WithConfigUpgradableController().
		WithLogLevelController()

	runnableAPIServerControllers, err := apiServerControllers.PrepareRun()
	if err != nil {
		return err
	}

	configObserver := configobservercontroller.NewConfigObserver(
		kubeInformersForNamespaces.InformersFor(operatorclient.TargetNamespace),
		kubeInformersForNamespaces.InformersFor(libgoetcd.EtcdEndpointNamespace),
		operatorClient,
		resourceSyncController,
		operatorConfigInformers,
		configInformers,
		featureGateAccessor,
		controllerConfig.EventRecorder,
	)

	staleConditions := staleconditions.NewRemoveStaleConditionsController(
		"openshift-apiserver",
		[]string{
			// in 4.1.0-4.3.0 this was used for indicating the apiserver daemonset was progressing
			"Progressing",
			// in 4.1.0-4.3.0 this was used for indicating the apiserver daemonset was available
			"Available",
			// in 4.1.0-4.3.z this was used for indicating the conditions of the apiserver daemonset.
			"APIServerDaemonSetAvailable",
			"APIServerDaemonSetProgressing",
			"APIServerDaemonSetDegraded",
			// in 4.5.z we changed (renamed) the following conditions
			"OpenshiftAPIServerStaticResourcesDegraded", // became APIServerStaticResourcesDegraded
			"WorkloadDegraded",                          // became APIServerWorkloadDegraded
		},
		operatorClient,
		controllerConfig.EventRecorder,
	)

	if controllerConfig.Server != nil {
		controllerConfig.Server.Handler.NonGoRestfulMux.Handle("/debug/controllers/resourcesync", debugHandler)
	}

	apiextensionsInformers := apiextensionsinformers.NewSharedInformerFactory(apiextensionsClient, 10*time.Minute)
	connectivityCheckController := connectivitycheckcontroller.NewOpenshiftAPIServerConnectivityCheckController(
		kubeClient,
		operatorClient,
		operatorcontrolplaneClient,
		apiextensionsClient,
		kubeInformersForNamespaces,
		configInformers,
		apiextensionsInformers,
		controllerConfig.EventRecorder,
	)

	operatorConfigInformers.Start(ctx.Done())
	kubeInformersForNamespaces.Start(ctx.Done())
	apiregistrationInformers.Start(ctx.Done())
	configInformers.Start(ctx.Done())
	dynamicInformers.Start(ctx.Done())
	migrationInformer.Start(ctx.Done())
	apiextensionsInformers.Start(ctx.Done())

	go configObserver.Run(ctx, 1)
	go resourceSyncController.Run(ctx, 1)
	go runnableAPIServerControllers.Run(ctx)
	go staleConditions.Run(ctx, 1)
	go connectivityCheckController.Run(ctx, 1)

	<-ctx.Done()
	return nil
}

func apiServices(clusterVersionLister configlisterv1.ClusterVersionLister) ([]*apiregistrationv1.APIService, []*apiregistrationv1.APIService, error) {
	clusterVersion, err := clusterVersionLister.Get("version")
	if err != nil {
		return nil, nil, err
	}

	knownCaps := sets.New[configv1.ClusterVersionCapability](clusterVersion.Status.Capabilities.KnownCapabilities...)
	capsEnabled := sets.New[configv1.ClusterVersionCapability](clusterVersion.Status.Capabilities.EnabledCapabilities...)

	groupDisabled := make(map[string]bool)
	if knownCaps.Has(configv1.ClusterVersionCapabilityBuild) && !capsEnabled.Has(configv1.ClusterVersionCapabilityBuild) {
		groupDisabled["build.openshift.io"] = true
	}
	if knownCaps.Has(configv1.ClusterVersionCapabilityDeploymentConfig) && !capsEnabled.Has(configv1.ClusterVersionCapabilityDeploymentConfig) {
		groupDisabled["apps.openshift.io"] = true
	}

	disabled := []*apiregistrationv1.APIService{}
	enabled := []*apiregistrationv1.APIService{}

	for _, apiServiceGroupVersion := range apiServiceGroupVersions {
		obj := &apiregistrationv1.APIService{
			ObjectMeta: metav1.ObjectMeta{
				Name: apiServiceGroupVersion.Version + "." + apiServiceGroupVersion.Group,
				Annotations: map[string]string{
					"service.alpha.openshift.io/inject-cabundle": "true",
				},
			},
			Spec: apiregistrationv1.APIServiceSpec{
				Group:   apiServiceGroupVersion.Group,
				Version: apiServiceGroupVersion.Version,
				Service: &apiregistrationv1.ServiceReference{
					Namespace: operatorclient.TargetNamespace,
					Name:      "api",
					Port:      utilpointer.Int32Ptr(443),
				},
				GroupPriorityMinimum: 9900,
				VersionPriority:      15,
			},
		}
		if groupDisabled[apiServiceGroupVersion.Group] {
			disabled = append(disabled, obj)
		} else {
			enabled = append(enabled, obj)
		}
	}

	return enabled, disabled, nil
}

func apiServicesReferences() []configv1.ObjectReference {
	ret := []configv1.ObjectReference{}
	for _, apiService := range apiServiceGroupVersions {
		ret = append(ret, configv1.ObjectReference{Group: "apiregistration.k8s.io", Resource: "apiservices", Name: apiService.Version + "." + apiService.Group})
	}
	return ret
}

func extractOperatorSpec(obj *unstructured.Unstructured, fieldManager string) (*applyoperatorv1.OperatorSpecApplyConfiguration, error) {
	castObj := &operatorv1.OpenShiftAPIServer{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.Object, castObj); err != nil {
		return nil, fmt.Errorf("unable to convert to OpenShiftAPIServer: %w", err)
	}
	ret, err := applyoperatorv1.ExtractOpenShiftAPIServer(castObj, fieldManager)
	if err != nil {
		return nil, fmt.Errorf("unable to extract fields for %q: %w", fieldManager, err)
	}
	if ret.Spec == nil {
		return nil, nil
	}
	return &ret.Spec.OperatorSpecApplyConfiguration, nil
}

func extractOperatorStatus(obj *unstructured.Unstructured, fieldManager string) (*applyoperatorv1.OperatorStatusApplyConfiguration, error) {
	castObj := &operatorv1.OpenShiftAPIServer{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.Object, castObj); err != nil {
		return nil, fmt.Errorf("unable to convert to OpenShiftAPIServer: %w", err)
	}
	ret, err := applyoperatorv1.ExtractOpenShiftAPIServerStatus(castObj, fieldManager)
	if err != nil {
		return nil, fmt.Errorf("unable to extract fields for %q: %w", fieldManager, err)
	}

	if ret.Status == nil {
		return nil, nil
	}
	return &ret.Status.OperatorStatusApplyConfiguration, nil
}
