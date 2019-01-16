package operator

import (
	"fmt"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
	appsclientv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	coreclientv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	apiregistrationv1client "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset/typed/apiregistration/v1"

	operatorv1 "github.com/openshift/api/operator/v1"
	openshiftconfigclientv1 "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/apis/openshiftapiserver/v1alpha1"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/v311_00_assets"
	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/resource/resourceapply"
	"github.com/openshift/library-go/pkg/operator/resource/resourcehash"
	"github.com/openshift/library-go/pkg/operator/resource/resourcemerge"
	"github.com/openshift/library-go/pkg/operator/resource/resourceread"
	"github.com/openshift/library-go/pkg/operator/v1helpers"
)

// syncOpenShiftAPIServer_v311_00_to_latest takes care of synchronizing (not upgrading) the thing we're managing.
// most of the time the sync method will be good for a large span of minor versions
func syncOpenShiftAPIServer_v311_00_to_latest(c OpenShiftAPIServerOperator, originalOperatorConfig *v1alpha1.OpenShiftAPIServerOperatorConfig) (bool, error) {
	errors := []error{}
	var err error
	operatorConfig := originalOperatorConfig.DeepCopy()

	directResourceResults := resourceapply.ApplyDirectly(c.kubeClient, c.eventRecorder, v311_00_assets.Asset,
		"v3.11.0/openshift-apiserver/ns.yaml",
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

	_, configMapModified, err := manageOpenShiftAPIServerConfigMap_v311_00_to_latest(c.kubeClient, c.kubeClient.CoreV1(), c.eventRecorder, operatorConfig)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q: %v", "configmap", err))
	}

	// the kube-apiserver is the source of truth for etcd serving CA bundles and etcd write keys.  We copy both so they can properly mounted
	etcdModified, err := manageOpenShiftAPIServerEtcdCerts_v311_00_to_latest(c.kubeClient.CoreV1(), c.eventRecorder)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q: %v", "etcd-certs", err))
	}
	// the kube-apiserver is the source of truth for client CA bundles
	clientCAModified, err := manageOpenShiftAPIServerClientCA_v311_00_to_latest(c.kubeClient.CoreV1(), c.eventRecorder)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q: %v", "client-ca", err))
	}

	imageImportCAModified, err := manageOpenShiftAPIServerImageImportCA_v311_00_to_latest(c.openshiftConfigClient, c.kubeClient.CoreV1(), c.eventRecorder)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q: %v", "client-ca", err))
	}

	forceRollingUpdate = forceRollingUpdate || operatorConfig.ObjectMeta.Generation != operatorConfig.Status.ObservedGeneration
	forceRollingUpdate = forceRollingUpdate || configMapModified || etcdModified || clientCAModified || imageImportCAModified

	// our configmaps and secrets are in order, now it is time to create the DS
	// TODO check basic preconditions here
	actualDaemonSet, _, err := manageOpenShiftAPIServerDaemonSet_v311_00_to_latest(c.kubeClient.AppsV1(), c.eventRecorder, c.targetImagePullSpec, operatorConfig.Status.Generations, forceRollingUpdate)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q: %v", "daemonsets", err))
	}

	// only manage the apiservices if we have ready pods for the daemonset.  This makes sure that if we're taking over for
	// something else, we don't stomp their apiservices until ours have a reasonable chance at working.
	var actualAPIServices []*apiregistrationv1.APIService
	if actualDaemonSet != nil && actualDaemonSet.Status.NumberAvailable > 0 {
		actualAPIServices, err = manageAPIServices_v311_00_to_latest(c.apiregistrationv1Client)
		if err != nil {
			errors = append(errors, fmt.Errorf("%q: %v", "apiservices", err))
		}
	}

	// manage status
	var availableConditionReason string
	var availableConditionMessages []string

	switch {
	case actualDaemonSet == nil:
		availableConditionReason = "NoDaemon"
		availableConditionMessages = append(availableConditionMessages, "daemonset/apiserver.openshift-apiserver: could not be retrieved")
	case actualDaemonSet.Status.NumberAvailable == 0:
		availableConditionReason = "NoAPIServerPod"
		availableConditionMessages = append(availableConditionMessages, "no openshift-apiserver daemon pods available on any node.")
	case actualDaemonSet.Status.NumberAvailable > 0 && len(actualAPIServices) == 0:
		availableConditionReason = "NoRegisteredAPIServices"
		availableConditionMessages = append(availableConditionMessages, "registered apiservices could not be retrieved")
	}
	for _, apiService := range actualAPIServices {
		for _, condition := range apiService.Status.Conditions {
			if condition.Type == apiregistrationv1.Available {
				if condition.Status == apiregistrationv1.ConditionFalse {
					availableConditionReason = "APIServiceNotAvailable"
					availableConditionMessages = append(availableConditionMessages, fmt.Sprintf("apiservice/%v: not available: %v", apiService.Name, condition.Message))
				}
				break
			}
		}
	}
	switch {
	case len(availableConditionMessages) == 1:
		v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorv1.OperatorCondition{
			Type:    operatorv1.OperatorStatusTypeAvailable,
			Status:  operatorv1.ConditionFalse,
			Reason:  availableConditionReason,
			Message: availableConditionMessages[0],
		})
	case len(availableConditionMessages) > 1:
		v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorv1.OperatorCondition{
			Type:    operatorv1.OperatorStatusTypeAvailable,
			Status:  operatorv1.ConditionFalse,
			Reason:  "Multiple",
			Message: strings.Join(availableConditionMessages, "\n"),
		})
	default:
		v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorv1.OperatorCondition{
			Type:   operatorv1.OperatorStatusTypeAvailable,
			Status: operatorv1.ConditionTrue,
		})
	}

	// If the daemonset is up to date and the operatorConfig are up to date, then we are no longer progressing
	var progressingMessages []string
	if actualDaemonSet != nil && actualDaemonSet.ObjectMeta.Generation != actualDaemonSet.Status.ObservedGeneration {
		progressingMessages = append(progressingMessages, fmt.Sprintf("daemonset/apiserver.openshift-operator: observed generation is %d, desired generation is %d.", actualDaemonSet.Status.ObservedGeneration, actualDaemonSet.ObjectMeta.Generation))
	}
	if operatorConfig.ObjectMeta.Generation != operatorConfig.Status.ObservedGeneration {
		progressingMessages = append(progressingMessages, fmt.Sprintf("openshiftapiserveroperatorconfigs/instance: observed generation is %d, desired generation is %d.", operatorConfig.Status.ObservedGeneration, operatorConfig.ObjectMeta.Generation))
	}

	if len(progressingMessages) == 0 {
		v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorv1.OperatorCondition{
			Type:   operatorv1.OperatorStatusTypeProgressing,
			Status: operatorv1.ConditionFalse,
		})
	} else {
		v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorv1.OperatorCondition{
			Type:    operatorv1.OperatorStatusTypeProgressing,
			Status:  operatorv1.ConditionTrue,
			Reason:  "DesiredStateNotYetAchieved",
			Message: strings.Join(progressingMessages, "\n"),
		})
	}

	// TODO this is changing too early and it was before too.
	operatorConfig.Status.ObservedGeneration = operatorConfig.ObjectMeta.Generation
	resourcemerge.SetDaemonSetGeneration(&operatorConfig.Status.Generations, actualDaemonSet)
	if len(errors) > 0 {
		message := ""
		for _, err := range errors {
			message = message + err.Error() + "\n"
		}
		v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorv1.OperatorCondition{
			Type:    workloadFailingCondition,
			Status:  operatorv1.ConditionTrue,
			Message: message,
			Reason:  "SyncError",
		})
	} else {
		v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorv1.OperatorCondition{
			Type:   workloadFailingCondition,
			Status: operatorv1.ConditionFalse,
		})
	}
	if !equality.Semantic.DeepEqual(operatorConfig.Status, originalOperatorConfig.Status) {
		if _, err := c.operatorConfigClient.OpenShiftAPIServerOperatorConfigs().UpdateStatus(operatorConfig); err != nil {
			return false, err
		}
	}

	if len(errors) > 0 {
		return true, nil
	}
	return false, nil
}

func manageOpenShiftAPIServerEtcdCerts_v311_00_to_latest(client coreclientv1.CoreV1Interface, recorder events.Recorder) (bool, error) {
	const etcdServingCAName = "etcd-serving-ca"
	const etcdClientCertKeyPairName = "etcd-client"

	_, caChanged, err := resourceapply.SyncConfigMap(client, recorder, etcdNamespaceName, etcdServingCAName, targetNamespaceName, etcdServingCAName, nil)
	if err != nil {
		return false, err
	}
	_, certKeyPairChanged, err := resourceapply.SyncSecret(client, recorder, etcdNamespaceName, etcdClientCertKeyPairName, targetNamespaceName, etcdClientCertKeyPairName, nil)
	if err != nil {
		return false, err
	}
	return caChanged || certKeyPairChanged, nil
}

func manageOpenShiftAPIServerClientCA_v311_00_to_latest(client coreclientv1.CoreV1Interface, recorder events.Recorder) (bool, error) {
	const apiserverClientCA = "client-ca"
	_, caChanged, err := resourceapply.SyncConfigMap(client, recorder, kubeAPIServerNamespaceName, apiserverClientCA, targetNamespaceName, apiserverClientCA, nil)
	if err != nil {
		return false, err
	}
	return caChanged, nil
}

func manageOpenShiftAPIServerImageImportCA_v311_00_to_latest(openshiftConfigClient openshiftconfigclientv1.ConfigV1Interface, client coreclientv1.CoreV1Interface, recorder events.Recorder) (bool, error) {
	imageConfig, err := openshiftConfigClient.Images().Get("cluster", metav1.GetOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		return false, err
	}
	if apierrors.IsNotFound(err) {
		return false, nil
	}
	if len(imageConfig.Spec.AdditionalTrustedCA.Name) == 0 {
		err := client.ConfigMaps(targetNamespaceName).Delete(imageImportCAName, nil)
		if apierrors.IsNotFound(err) {
			return false, nil
		}
		if err != nil {
			return false, err
		}
		return true, nil
	}
	_, caChanged, err := resourceapply.SyncConfigMap(client, recorder, userSpecifiedGlobalConfigNamespace, imageConfig.Spec.AdditionalTrustedCA.Name, targetNamespaceName, imageImportCAName, nil)
	if err != nil {
		return false, err
	}
	return caChanged, nil
}

func manageOpenShiftAPIServerConfigMap_v311_00_to_latest(kubeClient kubernetes.Interface, client coreclientv1.ConfigMapsGetter, recorder events.Recorder, operatorConfig *v1alpha1.OpenShiftAPIServerOperatorConfig) (*corev1.ConfigMap, bool, error) {
	configMap := resourceread.ReadConfigMapV1OrDie(v311_00_assets.MustAsset("v3.11.0/openshift-apiserver/cm.yaml"))
	defaultConfig := v311_00_assets.MustAsset("v3.11.0/openshift-apiserver/defaultconfig.yaml")
	requiredConfigMap, _, err := resourcemerge.MergeConfigMap(configMap, "config.yaml", nil, defaultConfig, operatorConfig.Spec.ObservedConfig.Raw, operatorConfig.Spec.UnsupportedConfigOverrides.Raw)
	if err != nil {
		return nil, false, err
	}

	// we can embed input hashes on our main configmap to drive rollouts when they change.
	inputHashes, err := resourcehash.MultipleObjectHashStringMapForObjectReferences(
		kubeClient,
		resourcehash.ObjectReference{
			Resource:  schema.GroupResource{Resource: "secret"},
			Namespace: targetNamespaceName,
			Name:      "serving-cert",
		},
		resourcehash.ObjectReference{
			Resource:  schema.GroupResource{Resource: "configmap"},
			Namespace: "kube-system",
			Name:      "extension-apiserver-authentication",
		},
	)
	if err != nil {
		return nil, false, err
	}
	for k, v := range inputHashes {
		requiredConfigMap.Data[k] = v
	}

	return resourceapply.ApplyConfigMap(client, recorder, requiredConfigMap)
}

func manageOpenShiftAPIServerDaemonSet_v311_00_to_latest(client appsclientv1.DaemonSetsGetter, recorder events.Recorder, imagePullSpec string, generationStatus []operatorv1.GenerationStatus, forceRollingUpdate bool) (*appsv1.DaemonSet, bool, error) {
	required := resourceread.ReadDaemonSetV1OrDie(v311_00_assets.MustAsset("v3.11.0/openshift-apiserver/ds.yaml"))
	required.Spec.Template.Spec.Containers[0].ImagePullPolicy = corev1.PullAlways
	if len(imagePullSpec) > 0 {
		required.Spec.Template.Spec.Containers[0].Image = imagePullSpec
	}
	required.Spec.Template.Spec.Containers[0].Args = append(required.Spec.Template.Spec.Containers[0].Args, fmt.Sprintf("-v=%d", 2))

	return resourceapply.ApplyDaemonSet(client, recorder, required, resourcemerge.ExpectedDaemonSetGeneration(required, generationStatus), forceRollingUpdate)
}

func manageAPIServices_v311_00_to_latest(client apiregistrationv1client.APIServicesGetter) ([]*apiregistrationv1.APIService, error) {
	apiServiceGroupVersions := []schema.GroupVersion{
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

	var apiServices []*apiregistrationv1.APIService
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
					Namespace: targetNamespaceName,
					Name:      "api",
				},
				GroupPriorityMinimum: 9900,
				VersionPriority:      15,
			},
		}

		apiService, _, err := resourceapply.ApplyAPIService(client, obj)
		if err != nil {
			return nil, err
		}
		apiServices = append(apiServices, apiService)
	}

	return apiServices, nil
}
