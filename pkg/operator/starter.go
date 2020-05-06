package operator

import (
	"context"
	"fmt"
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
	"github.com/openshift/library-go/pkg/controller/controllercmd"
	"github.com/openshift/library-go/pkg/operator/encryption"
	"github.com/openshift/library-go/pkg/operator/encryption/controllers/migrators"
	encryptiondeployer "github.com/openshift/library-go/pkg/operator/encryption/deployer"
	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/genericoperatorclient"
	"github.com/openshift/library-go/pkg/operator/resource/resourceapply"
	"github.com/openshift/library-go/pkg/operator/revisioncontroller"
	"github.com/openshift/library-go/pkg/operator/staleconditions"
	"github.com/openshift/library-go/pkg/operator/staticpod/controller/revision"
	"github.com/openshift/library-go/pkg/operator/staticresourcecontroller"
	"github.com/openshift/library-go/pkg/operator/status"
	"github.com/openshift/library-go/pkg/operator/v1helpers"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	apiregistrationclient "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
	apiregistrationinformers "k8s.io/kube-aggregator/pkg/client/informers/externalversions"
	utilpointer "k8s.io/utils/pointer"

	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/apiservercontrollerset"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/apiservicecontroller"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation/configobservercontroller"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation/etcdobserver"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"
	prune "github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/prunecontroller"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/resourcesynccontroller"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/v311_00_assets"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/workloadcontroller"
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

	operatorConfigInformers := operatorv1informers.NewSharedInformerFactory(operatorConfigClient, 10*time.Minute)
	kubeInformersForNamespaces := v1helpers.NewKubeInformersForNamespaces(kubeClient,
		"",
		operatorclient.GlobalUserSpecifiedConfigNamespace,
		operatorclient.GlobalMachineSpecifiedConfigNamespace,
		operatorclient.OperatorNamespace,
		operatorclient.TargetNamespace,
		etcdobserver.EtcdEndpointNamespace,
	)
	apiregistrationInformers := apiregistrationinformers.NewSharedInformerFactory(apiregistrationv1Client, 10*time.Minute)
	configInformers := configinformers.NewSharedInformerFactory(configClient, 10*time.Minute)

	operatorClient, dynamicInformers, err := genericoperatorclient.NewClusterScopedOperatorClient(controllerConfig.KubeConfig, operatorv1.GroupVersion.WithResource("openshiftapiservers"))
	if err != nil {
		return err
	}

	staticResourceController := staticresourcecontroller.NewStaticResourceController(
		"OpenshiftAPIServerStaticResources",
		v311_00_assets.Asset,
		[]string{
			"v3.11.0/openshift-apiserver/ns.yaml",
			"v3.11.0/openshift-apiserver/apiserver-clusterrolebinding.yaml",
			"v3.11.0/openshift-apiserver/svc.yaml",
			"v3.11.0/openshift-apiserver/sa.yaml",
			"v3.11.0/openshift-apiserver/trusted_ca_cm.yaml",
		},
		resourceapply.NewKubeClientHolder(kubeClient),
		operatorClient,
		controllerConfig.EventRecorder,
	).AddKubeInformers(kubeInformersForNamespaces)

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
		controllerConfig.EventRecorder,
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

	nodeInformer := kubeInformersForNamespaces.InformersFor("").Core().V1().Nodes()

	workloadController := workloadcontroller.NewWorkloadController(
		os.Getenv("IMAGE"), os.Getenv("OPERATOR_IMAGE_VERSION"), os.Getenv("OPERATOR_IMAGE"),
		operatorClient,
		versionRecorder,
		operatorConfigInformers.Operator().V1().OpenShiftAPIServers(),
		kubeInformersForNamespaces.InformersFor(operatorclient.TargetNamespace),
		kubeInformersForNamespaces.InformersFor(operatorclient.GlobalUserSpecifiedConfigNamespace),
		kubeInformersForNamespaces.InformersFor(operatorclient.GlobalUserSpecifiedConfigNamespace),
		apiregistrationInformers,
		configInformers,
		nodeInformer,
		operatorConfigClient.OperatorV1(),
		configClient.ConfigV1(),
		kubeClient,
		controllerConfig.EventRecorder,
	)

	apiServerControllers := apiservercontrollerset.NewAPIServerControllerSet(
		operatorClient,
		controllerConfig.EventRecorder,
	).WithAPIServiceController(
		"openshift-apiserver",
		apiservicecontroller.NewAPIServicesToManage(
			apiregistrationv1Client.ApiregistrationV1(),
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
				{Resource: "endpoints", Name: etcdobserver.EtcdEndpointName, Namespace: etcdobserver.EtcdEndpointNamespace},
			},
			apiServicesReferences()...,
		),
		configClient.ConfigV1(),
		configInformers.Config().V1().ClusterOperators(),
		versionRecorder,
	).WithConfigUpgradableController().
		WithLogLevelController()

	runnableAPIServerControllers, err := apiServerControllers.PrepareRun()
	if err != nil {
		return err
	}

	configObserver := configobservercontroller.NewConfigObserver(
		kubeInformersForNamespaces.InformersFor(operatorclient.TargetNamespace),
		kubeInformersForNamespaces.InformersFor(etcdobserver.EtcdEndpointNamespace),
		operatorClient,
		resourceSyncController,
		operatorConfigInformers,
		configInformers,
		controllerConfig.EventRecorder,
	)

	nodeProvider := DeploymentNodeProvider{
		TargetNamespaceDeploymentInformer: kubeInformersForNamespaces.InformersFor(operatorclient.TargetNamespace).Apps().V1().Deployments(),
		NodeInformer:                      nodeInformer,
	}

	deployer, err := encryptiondeployer.NewRevisionLabelPodDeployer("revision", operatorclient.TargetNamespace, kubeInformersForNamespaces, resourceSyncController, kubeClient.CoreV1(), kubeClient.CoreV1(), nodeProvider)
	if err != nil {
		return err
	}

	migrationClient := kubemigratorclient.NewForConfigOrDie(controllerConfig.KubeConfig)
	migrationInformer := migrationv1alpha1informer.NewSharedInformerFactory(migrationClient, time.Minute*30)
	migrator := migrators.NewKubeStorageVersionMigrator(migrationClient, migrationInformer.Migration().V1alpha1(), kubeClient.Discovery())

	encryptionControllers, err := encryption.NewControllers(
		operatorclient.TargetNamespace,
		deployer,
		migrator,
		operatorClient,
		configClient.ConfigV1().APIServers(),
		configInformers.Config().V1().APIServers(),
		kubeInformersForNamespaces,
		kubeClient.CoreV1(),
		controllerConfig.EventRecorder,
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
		controllerConfig.EventRecorder,
	)

	staleConditions := staleconditions.NewRemoveStaleConditions(
		[]string{
			// in 4.1.0-4.3.0 this was used for indicating the apiserver daemonset was progressing
			"Progressing",
			// in 4.1.0-4.3.0 this was used for indicating the apiserver daemonset was available
			"Available",
			// in 4.1.0-4.3.z this was used for indicating the conditions of the apiserver daemonset.
			"APIServerDaemonSetAvailable",
			"APIServerDaemonSetProgressing",
			"APIServerDaemonSetDegraded",
		},
		operatorClient,
		controllerConfig.EventRecorder,
	)

	if controllerConfig.Server != nil {
		controllerConfig.Server.Handler.NonGoRestfulMux.Handle("/debug/controllers/resourcesync", debugHandler)
	}

	ensureDaemonSetCleanup(kubeClient, ctx, controllerConfig.EventRecorder)

	operatorConfigInformers.Start(ctx.Done())
	kubeInformersForNamespaces.Start(ctx.Done())
	apiregistrationInformers.Start(ctx.Done())
	configInformers.Start(ctx.Done())
	dynamicInformers.Start(ctx.Done())
	migrationInformer.Start(ctx.Done())

	go staticResourceController.Run(ctx, 1)
	go workloadController.Run(ctx, 1)
	go configObserver.Run(ctx, 1)
	go resourceSyncController.Run(ctx, 1)
	go revisionController.Run(ctx, 1)
	go encryptionControllers.Run(ctx.Done())
	go pruneController.Run(ctx)
	go runnableAPIServerControllers.Run(ctx)
	go staleConditions.Run(ctx, 1)

	<-ctx.Done()
	return fmt.Errorf("stopped")
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
func ensureDaemonSetCleanup(kubeClient *kubernetes.Clientset, ctx context.Context, eventRecorder events.Recorder) {
	// daemonset and deployment both use the same name
	resourceName := "apiserver"

	dsClient := kubeClient.AppsV1().DaemonSets(operatorclient.TargetNamespace)
	deployClient := kubeClient.AppsV1().Deployments(operatorclient.TargetNamespace)

	go wait.UntilWithContext(ctx, func(_ context.Context) {
		// This function isn't expected to take long enough to suggest
		// checking that the context is done. The wait method will do that
		// checking.

		// Check whether the legacy daemonset exists and is not marked for deletion
		ds, err := dsClient.Get(resourceName, metav1.GetOptions{})
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
		deploy, err := deployClient.Get(resourceName, metav1.GetOptions{})
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
		err = dsClient.Delete(resourceName, &metav1.DeleteOptions{})
		if err != nil && !errors.IsNotFound(err) {
			klog.Warningf("Failed to delete legacy daemonset: %v", err)
			return
		}
		eventRecorder.Event("LegacyDaemonSetCleanup", "legacy daemonset has been removed")
	}, time.Minute)
}
