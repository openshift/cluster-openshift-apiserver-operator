package workloadcontroller

import (
	"fmt"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
	appsclientv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	coreclientv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	apiregistrationv1client "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset/typed/apiregistration/v1"

	operatorv1 "github.com/openshift/api/operator/v1"
	openshiftconfigclientv1 "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"
	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/resource/resourceapply"
	"github.com/openshift/library-go/pkg/operator/resource/resourcehash"
	"github.com/openshift/library-go/pkg/operator/resource/resourcemerge"
	"github.com/openshift/library-go/pkg/operator/resource/resourceread"
	"github.com/openshift/library-go/pkg/operator/status"
	"github.com/openshift/library-go/pkg/operator/v1helpers"

	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/v311_00_assets"
)

// syncOpenShiftAPIServer_v311_00_to_latest takes care of synchronizing (not upgrading) the thing we're managing.
// most of the time the sync method will be good for a large span of minor versions
func syncOpenShiftAPIServer_v311_00_to_latest(c OpenShiftAPIServerOperator, originalOperatorConfig *operatorv1.OpenShiftAPIServer) (bool, error) {
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

	imageImportCAModified, err := manageOpenShiftAPIServerImageImportCA_v311_00_to_latest(c.openshiftConfigClient, c.kubeClient.CoreV1(), c.eventRecorder)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q: %v", "client-ca", err))
	}

	forceRollingUpdate = forceRollingUpdate || operatorConfig.ObjectMeta.Generation != operatorConfig.Status.ObservedGeneration
	forceRollingUpdate = forceRollingUpdate || configMapModified || imageImportCAModified

	// our configmaps and secrets are in order, now it is time to create the DS
	// TODO check basic preconditions here
	actualDaemonSet, _, err := manageOpenShiftAPIServerDaemonSet_v311_00_to_latest(c.kubeClient.AppsV1(), c.eventRecorder, c.targetImagePullSpec, operatorConfig, operatorConfig.Status.Generations, forceRollingUpdate)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q: %v", "daemonsets", err))
	}

	// only manage the apiservices if we have ready pods for the daemonset.  This makes sure that if we're taking over for
	// something else, we don't stomp their apiservices until ours have a reasonable chance at working.
	var actualAPIServices []*apiregistrationv1.APIService
	manageAPIServices := actualDaemonSet != nil && actualDaemonSet.Status.NumberAvailable > 0
	if manageAPIServices {
		actualAPIServices, err = manageAPIServices_v311_00_to_latest(c.apiregistrationv1Client)
		if err != nil {
			errors = append(errors, fmt.Errorf("%q: %v", "apiservices", err))
		}
	}

	// manage status
	var availableConditions []operatorv1.OperatorCondition

	switch {
	case actualDaemonSet == nil:
		availableConditions = append(availableConditions, operatorv1.OperatorCondition{
			Reason:  "NoDaemon",
			Message: "daemonset/apiserver.openshift-apiserver: could not be retrieved",
		})
	case actualDaemonSet.Status.NumberAvailable == 0:
		availableConditions = append(availableConditions, operatorv1.OperatorCondition{
			Reason:  "NoAPIServerPod",
			Message: "no openshift-apiserver daemon pods available on any node.",
		})
	}

	//even if we can't manage the apiservices we still need to examine their availability
	if !manageAPIServices {
		for _, groupVersion := range apiServiceGroupVersions {
			apiService, err := c.apiregistrationv1Client.APIServices().Get(groupVersion.Version+"."+groupVersion.Group, metav1.GetOptions{})
			if err != nil {
				errors = append(errors, fmt.Errorf("%q: %v.%v: %v", "apiservices", groupVersion.Version, groupVersion.Group, err))
				availableConditions = append(availableConditions, operatorv1.OperatorCondition{
					Reason:  "APIServiceError",
					Message: fmt.Sprintf("apiservice/%v.%v: error retrieving: %v", groupVersion.Version, groupVersion.Group, err),
				})
			} else {
				actualAPIServices = append(actualAPIServices, apiService)
			}
		}
	}

	if len(actualAPIServices) == 0 {
		availableConditions = append(availableConditions, operatorv1.OperatorCondition{
			Reason:  "NoRegisteredAPIServices",
			Message: "registered apiservices could not be retrieved",
		})
	}

	for _, apiService := range actualAPIServices {
		for _, condition := range apiService.Status.Conditions {
			if condition.Type == apiregistrationv1.Available {
				if condition.Status == apiregistrationv1.ConditionFalse {
					availableConditions = append(availableConditions, operatorv1.OperatorCondition{
						Reason:  "APIServiceNotAvailable",
						Message: fmt.Sprintf("apiservice/%v: not available: %v", apiService.Name, condition.Message),
					})
				}
				break
			}
		}
	}

	// if the apiservices themselves check out ok, try to actually hit the discovery endpoints.  We have a history in clusterup
	// of something delaying them.  This isn't perfect because of round-robining, but let's see if we get an improvement
	if len(availableConditions) == 0 && c.kubeClient.Discovery().RESTClient() != nil {
		if missingAPIMessages := checkForAPIs(c.kubeClient.Discovery().RESTClient(), apiServiceGroupVersions...); len(missingAPIMessages) > 0 {
			availableConditions = append(availableConditions, operatorv1.OperatorCondition{
				Reason:  "APIServiceNotAvailable",
				Message: strings.Join(missingAPIMessages, "\n"),
			})
		}
	}

	availableCondition := operatorv1.OperatorCondition{Type: operatorv1.OperatorStatusTypeAvailable, Status: operatorv1.ConditionUnknown}
	if len(availableConditions) > 0 {
		availableCondition.Status = operatorv1.ConditionFalse
		var messages []string
		for _, condition := range availableConditions {
			if len(condition.Message) == 0 {
				continue
			}
			messages = append(messages, condition.Message)
		}
		if len(messages) > 0 {
			availableCondition.Message = strings.Join(messages, "\n")
		}
		if len(availableConditions) == 1 {
			availableCondition.Reason = availableConditions[0].Reason
		} else {
			availableCondition.Reason = "Multiple"
		}
	} else {
		availableCondition = operatorv1.OperatorCondition{
			Type:   operatorv1.OperatorStatusTypeAvailable,
			Status: operatorv1.ConditionTrue,
		}
	}
	v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, availableCondition)

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

	// if we are available, we need to try to set our versions correctly.
	if v1helpers.IsOperatorConditionTrue(operatorConfig.Status.Conditions, operatorv1.OperatorStatusTypeAvailable) {
		// we have the actual daemonset and we need the pull spec
		operandVersion := status.VersionForOperand(
			operatorclient.OperatorNamespace,
			actualDaemonSet.Spec.Template.Spec.Containers[0].Image,
			c.kubeClient.CoreV1(),
			c.eventRecorder)
		c.versionRecorder.SetVersion("openshift-apiserver", operandVersion)

	}
	if !equality.Semantic.DeepEqual(operatorConfig.Status, originalOperatorConfig.Status) {
		if _, err := c.operatorConfigClient.OpenShiftAPIServers().UpdateStatus(operatorConfig); err != nil {
			return false, err
		}
	}

	if len(errors) > 0 {
		return true, nil
	}
	if !v1helpers.IsOperatorConditionFalse(operatorConfig.Status.Conditions, operatorv1.OperatorStatusTypeFailing) {
		return true, nil
	}
	if !v1helpers.IsOperatorConditionFalse(operatorConfig.Status.Conditions, operatorv1.OperatorStatusTypeProgressing) {
		return true, nil
	}
	if !v1helpers.IsOperatorConditionTrue(operatorConfig.Status.Conditions, operatorv1.OperatorStatusTypeAvailable) {
		return true, nil
	}

	return false, nil
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
		err := client.ConfigMaps(operatorclient.TargetNamespace).Delete(imageImportCAName, nil)
		if apierrors.IsNotFound(err) {
			return false, nil
		}
		if err != nil {
			return false, err
		}
		return true, nil
	}
	_, caChanged, err := resourceapply.SyncConfigMap(client, recorder, operatorclient.GlobalUserSpecifiedConfigNamespace, imageConfig.Spec.AdditionalTrustedCA.Name, operatorclient.TargetNamespace, imageImportCAName, nil)
	if err != nil {
		return false, err
	}
	return caChanged, nil
}

func manageOpenShiftAPIServerConfigMap_v311_00_to_latest(kubeClient kubernetes.Interface, client coreclientv1.ConfigMapsGetter, recorder events.Recorder, operatorConfig *operatorv1.OpenShiftAPIServer) (*corev1.ConfigMap, bool, error) {
	configMap := resourceread.ReadConfigMapV1OrDie(v311_00_assets.MustAsset("v3.11.0/openshift-apiserver/cm.yaml"))
	defaultConfig := v311_00_assets.MustAsset("v3.11.0/openshift-apiserver/defaultconfig.yaml")
	requiredConfigMap, _, err := resourcemerge.MergeConfigMap(configMap, "config.yaml", nil, defaultConfig, operatorConfig.Spec.ObservedConfig.Raw, operatorConfig.Spec.UnsupportedConfigOverrides.Raw)
	if err != nil {
		return nil, false, err
	}

	// we can embed input hashes on our main configmap to drive rollouts when they change.
	inputHashes, err := resourcehash.MultipleObjectHashStringMapForObjectReferences(
		kubeClient,
		resourcehash.NewObjectRef().ForConfigMap().InNamespace(operatorclient.TargetNamespace).Named("aggregator-client-ca"),
		resourcehash.NewObjectRef().ForConfigMap().InNamespace(operatorclient.TargetNamespace).Named("client-ca"),
		resourcehash.NewObjectRef().ForSecret().InNamespace(operatorclient.TargetNamespace).Named("etcd-client"),
		resourcehash.NewObjectRef().ForConfigMap().InNamespace(operatorclient.TargetNamespace).Named("etcd-serving-ca"),
		resourcehash.NewObjectRef().ForSecret().InNamespace(operatorclient.TargetNamespace).Named("serving-cert"),
	)
	if err != nil {
		return nil, false, err
	}
	for k, v := range inputHashes {
		requiredConfigMap.Data[k] = v
	}

	return resourceapply.ApplyConfigMap(client, recorder, requiredConfigMap)
}

func manageOpenShiftAPIServerDaemonSet_v311_00_to_latest(client appsclientv1.DaemonSetsGetter, recorder events.Recorder, imagePullSpec string, operatorConfig *operatorv1.OpenShiftAPIServer, generationStatus []operatorv1.GenerationStatus, forceRollingUpdate bool) (*appsv1.DaemonSet, bool, error) {
	required := resourceread.ReadDaemonSetV1OrDie(v311_00_assets.MustAsset("v3.11.0/openshift-apiserver/ds.yaml"))
	if len(imagePullSpec) > 0 {
		required.Spec.Template.Spec.Containers[0].Image = imagePullSpec
	}
	// we set this so that when the requested image pull spec changes, we always have a diff.  Remember that we don't directly
	// diff any fields on the daemonset because they can be rewritten by admission and we don't want to constantly be fighting
	// against admission or defaults.  That was a problem with original versions of apply.
	if required.Annotations == nil {
		required.Annotations = map[string]string{}
	}
	required.Annotations["openshiftapiservers.operator.openshift.io/pull-spec"] = imagePullSpec

	switch operatorConfig.Spec.LogLevel {
	case operatorv1.Normal:
		required.Spec.Template.Spec.Containers[0].Args = append(required.Spec.Template.Spec.Containers[0].Args, fmt.Sprintf("-v=%d", 2))
	case operatorv1.Debug:
		required.Spec.Template.Spec.Containers[0].Args = append(required.Spec.Template.Spec.Containers[0].Args, fmt.Sprintf("-v=%d", 4))
	case operatorv1.Trace:
		required.Spec.Template.Spec.Containers[0].Args = append(required.Spec.Template.Spec.Containers[0].Args, fmt.Sprintf("-v=%d", 6))
	case operatorv1.TraceAll:
		required.Spec.Template.Spec.Containers[0].Args = append(required.Spec.Template.Spec.Containers[0].Args, fmt.Sprintf("-v=%d", 8))
	default:
		required.Spec.Template.Spec.Containers[0].Args = append(required.Spec.Template.Spec.Containers[0].Args, fmt.Sprintf("-v=%d", 2))
	}

	return resourceapply.ApplyDaemonSet(client, recorder, required, resourcemerge.ExpectedDaemonSetGeneration(required, generationStatus), forceRollingUpdate)
}

func manageAPIServices_v311_00_to_latest(client apiregistrationv1client.APIServicesGetter) ([]*apiregistrationv1.APIService, error) {
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
					Namespace: operatorclient.TargetNamespace,
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
