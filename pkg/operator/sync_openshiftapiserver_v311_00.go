package operator

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
	appsclientv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	coreclientv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	apiregistrationv1client "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset/typed/apiregistration/v1"

	operatorsv1alpha1 "github.com/openshift/api/operator/v1alpha1"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/apis/openshiftapiserver/v1alpha1"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/v311_00_assets"
	"github.com/openshift/library-go/pkg/operator/resource/resourceapply"
	"github.com/openshift/library-go/pkg/operator/resource/resourcemerge"
	"github.com/openshift/library-go/pkg/operator/resource/resourceread"
)

// syncOpenShiftAPIServer_v311_00_to_latest takes care of synchronizing (not upgrading) the thing we're managing.
// most of the time the sync method will be good for a large span of minor versions
func syncOpenShiftAPIServer_v311_00_to_latest(c OpenShiftAPIServerOperator, operatorConfig *v1alpha1.OpenShiftAPIServerOperatorConfig, previousAvailability *operatorsv1alpha1.VersionAvailability) (operatorsv1alpha1.VersionAvailability, []error) {
	versionAvailability := operatorsv1alpha1.VersionAvailability{
		Version: operatorConfig.Spec.Version,
	}

	errors := []error{}
	var err error

	directResourceResults := resourceapply.ApplyDirectly(c.kubeClient, v311_00_assets.Asset,
		"v3.11.0/openshift-apiserver/ns.yaml",
		"v3.11.0/openshift-apiserver/public-info-role.yaml",
		"v3.11.0/openshift-apiserver/public-info-rolebinding.yaml",
		"v3.11.0/openshift-apiserver/apiserver-clusterrolebinding.yaml",
		"v3.11.0/openshift-apiserver/svc.yaml",
		"v3.11.0/openshift-apiserver/sa.yaml",
	)
	resourcesThatForceRedeployment := sets.NewString("v3.11.0/openshift-apiserver/sa.yaml")
	forceRollingUpdate := false

	for _, currResult := range directResourceResults {
		if currResult.Error != nil {
			errors = append(errors, fmt.Errorf("%q (%T): %v", currResult.File, currResult.Type, currResult.Error))
			continue
		}

		if currResult.Changed && resourcesThatForceRedeployment.Has(currResult.File) {
			forceRollingUpdate = true
		}
	}

	apiserverConfig, configMapModified, err := manageOpenShiftAPIServerConfigMap_v311_00_to_latest(c.kubeClient.CoreV1(), operatorConfig)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q: %v", "configmap", err))
	}

	// the kube-apiserver is the source of truth for etcd serving CA bundles and etcd write keys.  We copy both so they can properly mounted
	etcdModified, err := manageOpenShiftAPIServerEtcdCerts_v311_00_to_latest(c.kubeClient.CoreV1())
	if err != nil {
		errors = append(errors, fmt.Errorf("%q: %v", "etcd-certs", err))
	}
	// the kube-apiserver is the source of truth for client CA bundles
	clientCAModified, err := manageOpenShiftAPIServerClientCA_v311_00_to_latest(c.kubeClient.CoreV1())
	if err != nil {
		errors = append(errors, fmt.Errorf("%q: %v", "client-ca", err))
	}

	forceRollingUpdate = forceRollingUpdate || operatorConfig.ObjectMeta.Generation != operatorConfig.Status.ObservedGeneration
	forceRollingUpdate = forceRollingUpdate || configMapModified || etcdModified || clientCAModified

	// our configmaps and secrets are in order, now it is time to create the DS
	// TODO check basic preconditions here
	actualDaemonSet, _, err := manageOpenShiftAPIServerDaemonSet_v311_00_to_latest(c.kubeClient.AppsV1(), operatorConfig, previousAvailability, forceRollingUpdate)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q: %v", "daemonsets", err))
	}

	// only manage the apiservices if we have ready pods for the daemonset.  This makes sure that if we're taking over for
	// something else, we don't stomp their apiservices until ours have a reasonable chance at working.
	if actualDaemonSet != nil && actualDaemonSet.Status.NumberReady > 0 {
		err = manageAPIServices_v311_00_to_latest(c.apiregistrationv1Client)
		if err != nil {
			errors = append(errors, fmt.Errorf("%q: %v", "apiservices", err))
		}
	}

	configData := ""
	if apiserverConfig != nil {
		configData = apiserverConfig.Data["config.yaml"]
	}
	_, _, err = manageOpenShiftAPIServerPublicConfigMap_v311_00_to_latest(c.kubeClient.CoreV1(), configData, operatorConfig)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q: %v", "configmap/public-info", err))
	}

	return resourcemerge.ApplyDaemonSetGenerationAvailability(versionAvailability, actualDaemonSet, errors...), errors
}

func manageOpenShiftAPIServerEtcdCerts_v311_00_to_latest(client coreclientv1.CoreV1Interface) (bool, error) {
	const etcdServingCAName = "etcd-serving-ca"
	const etcdClientCertKeyPairName = "etcd-client"

	_, caChanged, err := resourceapply.SyncConfigMap(client, etcdNamespaceName, etcdServingCAName, targetNamespaceName, etcdServingCAName)
	if err != nil {
		return false, err
	}
	_, certKeyPairChanged, err := resourceapply.SyncSecret(client, etcdNamespaceName, etcdClientCertKeyPairName, targetNamespaceName, etcdClientCertKeyPairName)
	if err != nil {
		return false, err
	}
	return caChanged || certKeyPairChanged, nil
}

func manageOpenShiftAPIServerClientCA_v311_00_to_latest(client coreclientv1.CoreV1Interface) (bool, error) {
	const apiserverClientCA = "client-ca"
	_, caChanged, err := resourceapply.SyncConfigMap(client, kubeAPIServerNamespaceName, apiserverClientCA, targetNamespaceName, apiserverClientCA)
	if err != nil {
		return false, err
	}
	return caChanged, nil
}

func manageOpenShiftAPIServerConfigMap_v311_00_to_latest(client coreclientv1.ConfigMapsGetter, operatorConfig *v1alpha1.OpenShiftAPIServerOperatorConfig) (*corev1.ConfigMap, bool, error) {
	configMap := resourceread.ReadConfigMapV1OrDie(v311_00_assets.MustAsset("v3.11.0/openshift-apiserver/cm.yaml"))
	defaultConfig := v311_00_assets.MustAsset("v3.11.0/openshift-apiserver/defaultconfig.yaml")
	requiredConfigMap, _, err := resourcemerge.MergeConfigMap(configMap, "config.yaml", nil, defaultConfig, operatorConfig.Spec.UserConfig.Raw, operatorConfig.Spec.ObservedConfig.Raw)
	if err != nil {
		return nil, false, err
	}
	return resourceapply.ApplyConfigMap(client, requiredConfigMap)
}

func manageOpenShiftAPIServerDaemonSet_v311_00_to_latest(client appsclientv1.DaemonSetsGetter, options *v1alpha1.OpenShiftAPIServerOperatorConfig, previousAvailability *operatorsv1alpha1.VersionAvailability, forceRollingUpdate bool) (*appsv1.DaemonSet, bool, error) {
	required := resourceread.ReadDaemonSetV1OrDie(v311_00_assets.MustAsset("v3.11.0/openshift-apiserver/ds.yaml"))
	required.Spec.Template.Spec.Containers[0].Image = options.Spec.ImagePullSpec
	required.Spec.Template.Spec.Containers[0].Args = append(required.Spec.Template.Spec.Containers[0].Args, fmt.Sprintf("-v=%d", options.Spec.Logging.Level))

	return resourceapply.ApplyDaemonSet(client, required, resourcemerge.ExpectedDaemonSetGeneration(required, previousAvailability), forceRollingUpdate)
}

func manageOpenShiftAPIServerPublicConfigMap_v311_00_to_latest(client coreclientv1.ConfigMapsGetter, apiserverConfigString string, operatorConfig *v1alpha1.OpenShiftAPIServerOperatorConfig) (*corev1.ConfigMap, bool, error) {
	configMap := resourceread.ReadConfigMapV1OrDie(v311_00_assets.MustAsset("v3.11.0/openshift-apiserver/public-info.yaml"))
	if operatorConfig.Status.CurrentAvailability != nil {
		configMap.Data["version"] = operatorConfig.Status.CurrentAvailability.Version
	} else {
		configMap.Data["version"] = ""
	}

	return resourceapply.ApplyConfigMap(client, configMap)
}

func manageAPIServices_v311_00_to_latest(client apiregistrationv1client.APIServicesGetter) error {
	apiServices := []schema.GroupVersion{
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

	for _, apiService := range apiServices {
		obj := &apiregistrationv1.APIService{
			ObjectMeta: metav1.ObjectMeta{
				Name: apiService.Version + "." + apiService.Group,
				Annotations: map[string]string{
					"service.alpha.openshift.io/inject-cabundle": "true",
				},
			},
			Spec: apiregistrationv1.APIServiceSpec{
				Group:   apiService.Group,
				Version: apiService.Version,
				Service: &apiregistrationv1.ServiceReference{
					Namespace: targetNamespaceName,
					Name:      "api",
				},
				GroupPriorityMinimum: 9900,
				VersionPriority:      15,
			},
		}

		_, _, err := resourceapply.ApplyAPIService(client, obj)
		if err != nil {
			return err
		}
	}

	return nil
}
