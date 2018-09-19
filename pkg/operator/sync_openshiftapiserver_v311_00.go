package operator

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	appsclientv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	coreclientv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	apiregistrationv1client "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset/typed/apiregistration/v1"

	operatorsv1alpha1 "github.com/openshift/api/operator/v1alpha1"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/apis/openshiftapiserver/v1alpha1"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/v311_00_assets"
	"github.com/openshift/library-go/pkg/operator/resource/resourceapply"
	"github.com/openshift/library-go/pkg/operator/resource/resourcecread"
	"github.com/openshift/library-go/pkg/operator/resource/resourcemerge"
)

// syncOpenShiftAPIServer_v311_00_to_latest takes care of synchronizing (not upgrading) the thing we're managing.
// most of the time the sync method will be good for a large span of minor versions
func syncOpenShiftAPIServer_v311_00_to_latest(c OpenShiftAPIServerOperator, operatorConfig *v1alpha1.OpenShiftAPIServerOperatorConfig, previousAvailability *operatorsv1alpha1.VersionAvailablity) (operatorsv1alpha1.VersionAvailablity, []error) {
	versionAvailability := operatorsv1alpha1.VersionAvailablity{
		Version: operatorConfig.Spec.Version,
	}

	errors := []error{}
	var err error

	requiredNamespace := resourceread.ReadNamespaceV1OrDie(v311_00_assets.MustAsset("v3.11.0/openshift-apiserver/ns.yaml"))
	_, _, err = resourceapply.ApplyNamespace(c.corev1Client, requiredNamespace)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q: %v", "ns", err))
	}

	requiredPublicRole := resourceread.ReadRoleV1OrDie(v311_00_assets.MustAsset("v3.11.0/openshift-apiserver/public-info-role.yaml"))
	_, _, err = resourceapply.ApplyRole(c.rbacv1Client, requiredPublicRole)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q: %v", "publicrole", err))
	}

	requiredPublicRoleBinding := resourceread.ReadRoleBindingV1OrDie(v311_00_assets.MustAsset("v3.11.0/openshift-apiserver/public-info-rolebinding.yaml"))
	_, _, err = resourceapply.ApplyRoleBinding(c.rbacv1Client, requiredPublicRoleBinding)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q: %v", "publicrolebinding", err))
	}

	requiredAPIServerClusterRoleBinding := resourceread.ReadClusterRoleBindingV1OrDie(v311_00_assets.MustAsset("v3.11.0/openshift-apiserver/apiserver-clusterrolebinding.yaml"))
	_, _, err = resourceapply.ApplyClusterRoleBinding(c.rbacv1Client, requiredAPIServerClusterRoleBinding)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q: %v", "apiserverclusterrole", err))
	}

	requiredService := resourceread.ReadServiceV1OrDie(v311_00_assets.MustAsset("v3.11.0/openshift-apiserver/svc.yaml"))
	_, _, err = resourceapply.ApplyService(c.corev1Client, requiredService)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q: %v", "svc", err))
	}

	requiredSA := resourceread.ReadServiceAccountV1OrDie(v311_00_assets.MustAsset("v3.11.0/openshift-apiserver/sa.yaml"))
	_, saModified, err := resourceapply.ApplyServiceAccount(c.corev1Client, requiredSA)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q: %v", "sa", err))
	}

	apiserverConfig, configMapModified, err := manageOpenShiftAPIServerConfigMap_v311_00_to_latest(c.corev1Client, operatorConfig)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q: %v", "configmap", err))
	}

	// the kube-apiserver is the source of truth for etcd serving CA bundles and etcd write keys.  We copy both so they can properly mounted
	etcdModified, err := manageOpenShiftAPIServerEtcdCerts_v311_00_to_latest(c.corev1Client)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q: %v", "etcd-certs", err))
	}
	// the kube-apiserver is the source of truth for client CA bundles
	clientCAModified, err := manageOpenShiftAPIServerClientCA_v311_00_to_latest(c.corev1Client)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q: %v", "client-ca", err))
	}

	forceDeployment := operatorConfig.ObjectMeta.Generation != operatorConfig.Status.ObservedGeneration
	forceDeployment = forceDeployment || saModified || configMapModified || etcdModified || clientCAModified

	// our configmaps and secrets are in order, now it is time to create the DS
	// TODO check basic preconditions here
	actualDeployment, _, err := manageOpenShiftAPIServerDeployment_v311_00_to_latest(c.appsv1Client, operatorConfig, previousAvailability, forceDeployment)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q: %v", "deployment", err))
	}

	err = manageAPIServices_v311_00_to_latest(c.apiregistrationv1Client)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q: %v", "apiservices", err))
	}

	configData := ""
	if apiserverConfig != nil {
		configData = apiserverConfig.Data["config.yaml"]
	}
	_, _, err = manageOpenShiftAPIServerPublicConfigMap_v311_00_to_latest(c.corev1Client, configData, operatorConfig)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q: %v", "configmap/public-info", err))
	}

	return resourcemerge.ApplyGenerationAvailability(versionAvailability, actualDeployment, errors...), errors
}

func manageOpenShiftAPIServerEtcdCerts_v311_00_to_latest(client coreclientv1.CoreV1Interface) (bool, error) {
	const etcdServingCAName = "etcd-serving-ca"
	const etcdClientCertKeyPairName = "etcd-client"

	_, caChanged, err := resourceapply.SyncConfigMap(client, kubeAPIServerNamespaceName, etcdServingCAName, targetNamespaceName, etcdServingCAName)
	if err != nil{
		return false, err
	}
	_, certKeyPairChanged, err := resourceapply.SyncSecret(client, kubeAPIServerNamespaceName, etcdClientCertKeyPairName, targetNamespaceName, etcdClientCertKeyPairName)
	if err != nil {
		return false, err
	}
	return caChanged || certKeyPairChanged, nil
}

func manageOpenShiftAPIServerClientCA_v311_00_to_latest(client coreclientv1.CoreV1Interface) (bool, error) {
	const apiserverClientCA = "client-ca"
	_, caChanged, err := resourceapply.SyncConfigMap(client, kubeAPIServerNamespaceName, apiserverClientCA, targetNamespaceName, apiserverClientCA)
	if err != nil{
		return false, err
	}
	return caChanged, nil
}

func manageOpenShiftAPIServerConfigMap_v311_00_to_latest(client coreclientv1.ConfigMapsGetter, operatorConfig *v1alpha1.OpenShiftAPIServerOperatorConfig) (*corev1.ConfigMap, bool, error) {
	configMap := resourceread.ReadConfigMapV1OrDie(v311_00_assets.MustAsset("v3.11.0/openshift-apiserver/cm.yaml"))
	defaultConfig := v311_00_assets.MustAsset("v3.11.0/openshift-apiserver/defaultconfig.yaml")
	requiredConfigMap, _, err := resourcemerge.MergeConfigMap(configMap, "config.yaml", nil, defaultConfig, operatorConfig.Spec.OpenShiftAPIServerConfig.Raw)
	if err != nil {
		return nil, false, err
	}
	return resourceapply.ApplyConfigMap(client, requiredConfigMap)
}

func manageOpenShiftAPIServerDeployment_v311_00_to_latest(client appsclientv1.DeploymentsGetter, options *v1alpha1.OpenShiftAPIServerOperatorConfig, previousAvailability *operatorsv1alpha1.VersionAvailablity, forceDeployment bool) (*appsv1.Deployment, bool, error) {
	required := resourceread.ReadDeploymentV1OrDie(v311_00_assets.MustAsset("v3.11.0/openshift-apiserver/deployment.yaml"))
	required.Spec.Template.Spec.Containers[0].Image = options.Spec.ImagePullSpec
	required.Spec.Template.Spec.Containers[0].Args = append(required.Spec.Template.Spec.Containers[0].Args, fmt.Sprintf("-v=%d", options.Spec.Logging.Level))

	return resourceapply.ApplyDeployment(client, required, resourcemerge.ExpectedDeploymentGeneration(required, previousAvailability), forceDeployment)
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
		{Group: "apps.openshift.io", Version: "v1"},
		{Group: "authorization.openshift.io", Version: "v1"},
		{Group: "build.openshift.io", Version: "v1"},
		{Group: "image.openshift.io", Version: "v1"},
		{Group: "network.openshift.io", Version: "v1"},
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
