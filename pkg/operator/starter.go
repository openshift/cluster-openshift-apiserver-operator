package operator

import (
	"context"
	"os"
	"time"

	kubemigratorclient "github.com/kubernetes-sigs/kube-storage-version-migrator/pkg/clients/clientset"
	migrationv1alpha1informer "github.com/kubernetes-sigs/kube-storage-version-migrator/pkg/clients/informer"
	configv1 "github.com/openshift/api/config/v1"
	operatorv1 "github.com/openshift/api/operator/v1"
	configv1client "github.com/openshift/client-go/config/clientset/versioned"
	configinformers "github.com/openshift/client-go/config/informers/externalversions"
	operatorv1client "github.com/openshift/client-go/operator/clientset/versioned"
	operatorv1informers "github.com/openshift/client-go/operator/informers/externalversions"
	operatorcontrolplaneclient "github.com/openshift/client-go/operatorcontrolplane/clientset/versioned"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/connectivitycheckcontroller"
	"github.com/openshift/library-go/pkg/controller/controllercmd"
	"github.com/openshift/library-go/pkg/operator/encryption/controllers/migrators"
	encryptiondeployer "github.com/openshift/library-go/pkg/operator/encryption/deployer"
	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/genericoperatorclient"
	"github.com/openshift/library-go/pkg/operator/staleconditions"
	"github.com/openshift/library-go/pkg/operator/staticpod/controller/revision"
	"github.com/openshift/library-go/pkg/operator/status"
	"github.com/openshift/library-go/pkg/operator/v1helpers"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	apiregistrationclient "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
	apiregistrationinformers "k8s.io/kube-aggregator/pkg/client/informers/externalversions"
	utilpointer "k8s.io/utils/pointer"

	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/apiservice"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation/configobservercontroller"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/encryptionprovider"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/oauthapiencryption"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/resourcesynccontroller"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/revisionpoddeployer"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/v311_00_assets"
	operatorworkload "github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/workload"
	libgoassets "github.com/openshift/library-go/pkg/operator/apiserver/audit"
	workloadcontroller "github.com/openshift/library-go/pkg/operator/apiserver/controller/workload"
	apiservercontrollerset "github.com/openshift/library-go/pkg/operator/apiserver/controllerset"
	libgoetcd "github.com/openshift/library-go/pkg/operator/configobserver/etcd"
)

const (
	oauthAPIServerTargetNamespace = "openshift-oauth-apiserver"
)

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

	operatorClient, dynamicInformers, err := genericoperatorclient.NewClusterScopedOperatorClient(controllerConfig.KubeConfig, operatorv1.GroupVersion.WithResource("openshiftapiservers"))
	if err != nil {
		return err
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
	openshiftDeployer, err := encryptiondeployer.NewRevisionLabelPodDeployer("revision", operatorclient.TargetNamespace, kubeInformersForNamespaces, resourceSyncController, kubeClient.CoreV1(), kubeClient.CoreV1(), openshiftNodeProvider)
	if err != nil {
		return err
	}
	oauthNodeProvider := encryptiondeployer.NewDeploymentNodeProvider(oauthAPIServerTargetNamespace, kubeInformersForNamespaces)
	// note: do not pass resourceSyncController - CAO is responsible for synchronization
	oauthDeployer, err := encryptiondeployer.NewRevisionLabelPodDeployer("revision", oauthAPIServerTargetNamespace, kubeInformersForNamespaces, nil, kubeClient.CoreV1(), kubeClient.CoreV1(), oauthNodeProvider)
	if err != nil {
		return err
	}
	weManageOauthConfig := encryptionprovider.IsOAuthEncryptionConfigManagedByThisOperator(kubeInformersForNamespaces.InformersFor(operatorclient.GlobalMachineSpecifiedConfigNamespace).Core().V1().Secrets().Lister().Secrets(operatorclient.GlobalMachineSpecifiedConfigNamespace), oauthAPIServerTargetNamespace, oauthapiencryption.EncryptionConfigManagedBy)
	deployer, err := revisionpoddeployer.NewUnionDeployer(&revisionpoddeployer.AlwaysEnabledDeployer{Deployer: openshiftDeployer}, revisionpoddeployer.NewDisabledByPredicateDeployer(weManageOauthConfig, oauthDeployer))

	migrationClient := kubemigratorclient.NewForConfigOrDie(controllerConfig.KubeConfig)
	migrationInformer := migrationv1alpha1informer.NewSharedInformerFactory(migrationClient, time.Minute*30)
	migrator := migrators.NewKubeStorageVersionMigrator(migrationClient, migrationInformer.Migration().V1alpha1(), kubeClient.Discovery())
	encryptionProvider := encryptionprovider.New(
		oauthAPIServerTargetNamespace,
		oauthapiencryption.EncryptionConfigManagedBy,
		[]schema.GroupResource{
			{Group: "route.openshift.io", Resource: "routes"}, // routes can contain embedded TLS private keys
			{Group: "oauth.openshift.io", Resource: "oauthaccesstokens"},
			{Group: "oauth.openshift.io", Resource: "oauthauthorizetokens"},
		},
		sets.NewString("oauthaccesstokens.oauth.openshift.io", "oauthauthorizetokens.oauth.openshift.io"),
		kubeInformersForNamespaces,
	)

	openShiftAPIServerWorkload := operatorworkload.NewOpenShiftAPIServerWorkload(
		operatorClient,
		operatorConfigClient.OperatorV1(),
		configClient.ConfigV1(),
		workloadcontroller.CountNodesFuncWrapper(kubeInformersForNamespaces.InformersFor("").Core().V1().Nodes().Lister()),
		workloadcontroller.EnsureAtMostOnePodPerNode,
		"openshift-apiserver",
		os.Getenv("IMAGE"),
		os.Getenv("OPERATOR_IMAGE"),
		kubeClient,
		controllerConfig.EventRecorder,
		versionRecorder)

	apiServerControllers := apiservercontrollerset.NewAPIServerControllerSet(
		operatorClient,
		controllerConfig.EventRecorder,
	).WithAPIServiceController(
		"openshift-apiserver",
		apiservice.NewAPIServicesToManage(
			apiregistrationInformers.Apiregistration().V1().APIServices().Lister(),
			operatorConfigInformers.Operator().V1().Authentications().Lister(),
			apiServices(),
			controllerConfig.EventRecorder,
			sets.NewString("v1.oauth.openshift.io", "v1.user.openshift.io"),
			"authentication.operator.openshift.io/managed",
		).GetAPIServicesToManage,
		apiregistrationInformers,
		apiregistrationv1Client.ApiregistrationV1(),
		kubeInformersForNamespaces.InformersFor(operatorclient.TargetNamespace),
		kubeClient,
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
	).WithWorkloadController(
		"OpenShiftAPIServer",
		operatorclient.OperatorNamespace,
		operatorclient.TargetNamespace,
		os.Getenv("OPERATOR_IMAGE_VERSION"),
		"openshift",
		"APIServer",
		kubeClient,
		openShiftAPIServerWorkload,
		configClient.ConfigV1().ClusterOperators(),
		versionRecorder,
		kubeInformersForNamespaces,
		operatorConfigInformers.Operator().V1().OpenShiftAPIServers().Informer(),
		configInformers.Config().V1().Images().Informer(),
	).WithStaticResourcesController(
		"APIServerStaticResources",
		libgoassets.WithAuditPolicies("audit", operatorclient.TargetNamespace, v311_00_assets.Asset),
		[]string{
			"v3.11.0/openshift-apiserver/ns.yaml",
			"v3.11.0/openshift-apiserver/apiserver-clusterrolebinding.yaml",
			"v3.11.0/openshift-apiserver/svc.yaml",
			"v3.11.0/openshift-apiserver/sa.yaml",
			"v3.11.0/openshift-apiserver/trusted_ca_cm.yaml",
			libgoassets.AuditPoliciesConfigMapFileName,
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
		OpenshiftDeploymentLatestRevisionClient{OperatorClient: operatorClient, TypedClient: operatorConfigClient.OperatorV1()},
		v1helpers.CachedConfigMapGetter(kubeClient.CoreV1(), kubeInformersForNamespaces),
		v1helpers.CachedSecretGetter(kubeClient.CoreV1(), kubeInformersForNamespaces),
	).WithEncryptionControllers(
		operatorclient.TargetNamespace,
		encryptionProvider,
		deployer,
		migrator,
		kubeClient.CoreV1(),
		configClient.ConfigV1().APIServers(),
		configInformers.Config().V1().APIServers(),
		kubeInformersForNamespaces,
	).WithSecretRevisionPruneController(
		operatorclient.TargetNamespace,
		[]string{"encryption-config-"},
		kubeClient.CoreV1(),
		kubeClient.CoreV1(),
		kubeInformersForNamespaces,
	).
		WithConfigUpgradableController().
		WithLogLevelController()

	runnableAPIServerControllers, err := apiServerControllers.PrepareRun()
	if err != nil {
		return err
	}

	oauthEncryptionController := oauthapiencryption.NewEncryptionConfigSyncController(
		"OAuthAPIServerController",
		oauthAPIServerTargetNamespace,
		kubeClient.CoreV1(),
		kubeInformersForNamespaces,
		controllerConfig.EventRecorder,
	)
	auditPolicyPahGetter, err := libgoassets.NewAuditPolicyPathGetter("/var/run/configmaps/audit")
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
		auditPolicyPahGetter,
		controllerConfig.EventRecorder,
	)

	staleConditions := staleconditions.NewRemoveStaleConditionsController(
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

	connectivityCheckController := connectivitycheckcontroller.NewOpenshiftAPIServerConnectivityCheckController(
		kubeClient,
		operatorClient,
		kubeInformersForNamespaces,
		configInformers,
		operatorcontrolplaneClient,
		controllerConfig.EventRecorder,
	)

	ensureDaemonSetCleanup(ctx, kubeClient, controllerConfig.EventRecorder)

	operatorConfigInformers.Start(ctx.Done())
	kubeInformersForNamespaces.Start(ctx.Done())
	apiregistrationInformers.Start(ctx.Done())
	configInformers.Start(ctx.Done())
	dynamicInformers.Start(ctx.Done())
	migrationInformer.Start(ctx.Done())

	go oauthEncryptionController.Run(ctx, 1)
	go configObserver.Run(ctx, 1)
	go resourceSyncController.Run(ctx, 1)
	go runnableAPIServerControllers.Run(ctx)
	go staleConditions.Run(ctx, 1)
	go connectivityCheckController.Run(ctx, 1)

	<-ctx.Done()
	return nil
}

func apiServices() []*apiregistrationv1.APIService {
	var apiServiceGroupVersions = []schema.GroupVersion{
		// these are all the apigroups we manage
		{Group: "apps.openshift.io", Version: "v1"},
		{Group: "authorization.openshift.io", Version: "v1"},
		{Group: "build.openshift.io", Version: "v1"},
		{Group: "image.openshift.io", Version: "v1"},
		{Group: "oauth.openshift.io", Version: "v1"},
		{Group: "project.openshift.io", Version: "v1"},
		{Group: "quota.openshift.io", Version: "v1"},
		{Group: "route.openshift.io", Version: "v1"},
		{Group: "security.openshift.io", Version: "v1"},
		{Group: "template.openshift.io", Version: "v1"},
		{Group: "user.openshift.io", Version: "v1"},
	}

	ret := []*apiregistrationv1.APIService{}
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
		ret = append(ret, obj)
	}

	return ret
}

func apiServicesReferences() []configv1.ObjectReference {
	ret := []configv1.ObjectReference{}
	for _, apiService := range apiServices() {
		ret = append(ret, configv1.ObjectReference{Group: "apiregistration.k8s.io", Resource: "apiservices", Name: apiService.Spec.Version + "." + apiService.Spec.Group})
	}
	return ret
}

// ensureDaemonSetCleanup continually ensures the removal of the daemonset
// used to manage apiserver pods in releases prior to 4.5. The daemonset is
// removed once the deployment now managing apiserver pods reports at least
// one pod available.
func ensureDaemonSetCleanup(ctx context.Context, kubeClient *kubernetes.Clientset, eventRecorder events.Recorder) {
	// daemonset and deployment both use the same name
	resourceName := "apiserver"

	dsClient := kubeClient.AppsV1().DaemonSets(operatorclient.TargetNamespace)
	deployClient := kubeClient.AppsV1().Deployments(operatorclient.TargetNamespace)

	go wait.UntilWithContext(ctx, func(_ context.Context) {
		// This function isn't expected to take long enough to suggest
		// checking that the context is done. The wait method will do that
		// checking.

		// Check whether the legacy daemonset exists and is not marked for deletion
		ds, err := dsClient.Get(ctx, resourceName, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			// Done - daemonset does not exist
			return
		}
		if err != nil {
			klog.Warningf("Error retrieving legacy daemonset: %v", err)
			return
		}
		if ds.ObjectMeta.DeletionTimestamp != nil {
			// Done - daemonset has been marked for deletion
			return
		}

		// Check that the deployment managing the apiserver pods has at last one available replica
		deploy, err := deployClient.Get(ctx, resourceName, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			// No available replicas if the deployment doesn't exist
			return
		}
		if err != nil {
			klog.Warningf("Error retrieving the deployment that manages apiserver pods: %v", err)
			return
		}
		if deploy.Status.AvailableReplicas == 0 {
			eventRecorder.Warning("LegacyDaemonSetCleanup", "the deployment replacing the daemonset does not have available replicas yet")
			return
		}

		// Safe to remove legacy daemonset since the deployment has at least one available replica
		err = dsClient.Delete(ctx, resourceName, metav1.DeleteOptions{})
		if err != nil && !errors.IsNotFound(err) {
			klog.Warningf("Failed to delete legacy daemonset: %v", err)
			return
		}
		eventRecorder.Event("LegacyDaemonSetCleanup", "legacy daemonset has been removed")
	}, time.Minute)
}
