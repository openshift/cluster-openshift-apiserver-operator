package operator

import (
	"context"
	"fmt"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/apiservicecontroller"
	"k8s.io/apimachinery/pkg/util/sets"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	apiregistrationclient "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
	apiregistrationinformers "k8s.io/kube-aggregator/pkg/client/informers/externalversions"
	utilpointer "k8s.io/utils/pointer"

	configv1 "github.com/openshift/api/config/v1"
	operatorv1 "github.com/openshift/api/operator/v1"
	configv1client "github.com/openshift/client-go/config/clientset/versioned"
	configinformers "github.com/openshift/client-go/config/informers/externalversions"
	operatorv1client "github.com/openshift/client-go/operator/clientset/versioned"
	operatorv1informers "github.com/openshift/client-go/operator/informers/externalversions"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation/configobservercontroller"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"
	prune "github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/prunecontroller"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/resourcesynccontroller"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/workloadcontroller"
	"github.com/openshift/library-go/pkg/controller/controllercmd"
	apiservercontrollerset "github.com/openshift/library-go/pkg/operator/apiserver/controllerset"
	"github.com/openshift/library-go/pkg/operator/genericoperatorclient"
	"github.com/openshift/library-go/pkg/operator/revisioncontroller"
	"github.com/openshift/library-go/pkg/operator/staticpod/controller/revision"
	"github.com/openshift/library-go/pkg/operator/status"
	"github.com/openshift/library-go/pkg/operator/v1helpers"
)

func RunOperator(ctx context.Context, controllerConfig *controllercmd.ControllerContext) error {
	kubeClient, err := kubernetes.NewForConfig(controllerConfig.ProtoKubeConfig)
	if err != nil {
		return err
	}
	migrationClientConfig := dynamic.ConfigFor(controllerConfig.KubeConfig)
	migrationClientConfig.Burst = 40
	migrationClientConfig.QPS = 30
	dynamicClientForMigration, err := dynamic.NewForConfig(migrationClientConfig)
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
		controllerConfig.EventRecorder,
	)

	apiServerControllers, err := apiservercontrollerset.NewAPIServerControllerSet(
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
	).WithEncryptionControllers(
		operatorclient.TargetNamespace,
		[]schema.GroupResource{
			{Group: "route.openshift.io", Resource: "routes"}, // routes can contain embedded TLS private keys
			{Group: "oauth.openshift.io", Resource: "oauthaccesstokens"},
			{Group: "oauth.openshift.io", Resource: "oauthauthorizetokens"},
		},
		dynamicClientForMigration,
		configClient.ConfigV1().APIServers(),
		configInformers.Config().V1().APIServers(),
		kubeClient,
		kubeInformersForNamespaces,
		resourceSyncController,
	).WithClusterOperatorStatusController(
		"openshift-apiserver",
		append(
			[]configv1.ObjectReference{
				{Group: "operator.openshift.io", Resource: "openshiftapiservers", Name: "cluster"},
				{Resource: "namespaces", Name: operatorclient.GlobalUserSpecifiedConfigNamespace},
				{Resource: "namespaces", Name: operatorclient.GlobalMachineSpecifiedConfigNamespace},
				{Resource: "namespaces", Name: operatorclient.OperatorNamespace},
				{Resource: "namespaces", Name: operatorclient.TargetNamespace},
			},
			apiServicesReferences()...,
		),
		configClient.ConfigV1(),
		configInformers.Config().V1().ClusterOperators(),
		versionRecorder,
	).WithConfigUpgradableController().
		WithLogLevelController().
		PrepareRun()
	if err != nil {
		return err
	}

	configObserver := configobservercontroller.NewConfigObserver(
		kubeInformersForNamespaces.InformersFor(operatorclient.TargetNamespace),
		operatorClient,
		resourceSyncController,
		operatorConfigInformers,
		configInformers,
		controllerConfig.EventRecorder,
	)

	pruneController := prune.NewPruneController(
		operatorclient.TargetNamespace,
		[]string{"encryption-config-"},
		kubeClient.CoreV1(),
		kubeClient.CoreV1(),
		kubeInformersForNamespaces,
		controllerConfig.EventRecorder,
	)

	if controllerConfig.Server != nil {
		controllerConfig.Server.Handler.NonGoRestfulMux.Handle("/debug/controllers/resourcesync", debugHandler)
	}

	operatorConfigInformers.Start(ctx.Done())
	kubeInformersForNamespaces.Start(ctx.Done())
	apiregistrationInformers.Start(ctx.Done())
	configInformers.Start(ctx.Done())
	dynamicInformers.Start(ctx.Done())

	go workloadController.Run(ctx, 1)
	go configObserver.Run(ctx, 1)
	go resourceSyncController.Run(ctx, 1)
	go revisionController.Run(ctx, 1)
	go pruneController.Run(ctx)
	go apiServerControllers.Run(ctx)

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
