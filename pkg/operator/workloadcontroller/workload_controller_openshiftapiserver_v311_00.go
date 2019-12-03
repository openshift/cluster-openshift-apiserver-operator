package workloadcontroller

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	appsclientv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	coreclientv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	apiregistrationv1client "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset/typed/apiregistration/v1"

	openshiftapi "github.com/openshift/api"
	openshiftcontrolplanev1 "github.com/openshift/api/openshiftcontrolplane/v1"
	operatorv1 "github.com/openshift/api/operator/v1"
	openshiftconfigclientv1 "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/v311_00_assets"
	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/resource/resourceapply"
	"github.com/openshift/library-go/pkg/operator/resource/resourcehash"
	"github.com/openshift/library-go/pkg/operator/resource/resourcemerge"
	"github.com/openshift/library-go/pkg/operator/resource/resourceread"
	"github.com/openshift/library-go/pkg/operator/status"
	"github.com/openshift/library-go/pkg/operator/v1helpers"
)

// syncOpenShiftAPIServer_v311_00_to_latest takes care of synchronizing (not upgrading) the thing we're managing.
// most of the time the sync method will be good for a large span of minor versions
func syncOpenShiftAPIServer_v311_00_to_latest(c OpenShiftAPIServerOperator, originalOperatorConfig *operatorv1.OpenShiftAPIServer) (bool, error) {
	errors := []error{}
	operatorConfig := originalOperatorConfig.DeepCopy()

	directResourceResults := resourceapply.ApplyDirectly(c.kubeClient, c.eventRecorder, v311_00_assets.Asset,
		"v3.11.0/openshift-apiserver/ns.yaml",
		"v3.11.0/openshift-apiserver/apiserver-clusterrolebinding.yaml",
		"v3.11.0/openshift-apiserver/svc.yaml",
		"v3.11.0/openshift-apiserver/sa.yaml",
	)
	resourcesThatForceRedeployment := sets.NewString("v3.11.0/openshift-apiserver/sa.yaml")
	var reasonsForForcedRollingUpdate []string

	for _, currResult := range directResourceResults {
		if currResult.Error != nil {
			errors = append(errors, fmt.Errorf("%q (%T): %v", currResult.File, currResult.Type, currResult.Error))
			continue
		}

		if currResult.Changed && resourcesThatForceRedeployment.Has(currResult.File) {
			reasonsForForcedRollingUpdate = append(reasonsForForcedRollingUpdate, "modified: "+resourceSelectorForCLI(currResult.Result))
		}
	}

	configMapModifiedObject, configMapModified, err := manageOpenShiftAPIServerConfigMap_v311_00_to_latest(c.kubeClient, c.kubeClient.CoreV1(), c.eventRecorder, operatorConfig)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q: %v", "configmap", err))
	}
	if configMapModified {
		reasonsForForcedRollingUpdate = append(reasonsForForcedRollingUpdate, "modified: "+resourceSelectorForCLI(configMapModifiedObject))
	}

	imageImportCAModifiedObject, imageImportCAModified, err := manageOpenShiftAPIServerImageImportCA_v311_00_to_latest(c.openshiftConfigClient, c.kubeClient.CoreV1(), c.eventRecorder)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q: %v", "client-ca", err))
	}
	if imageImportCAModified {
		reasonsForForcedRollingUpdate = append(reasonsForForcedRollingUpdate, "modified: "+resourceSelectorForCLI(imageImportCAModifiedObject))
	}

	err = syncOpenShiftAPIServerTrustedCA_v311_00_to_latest(c.kubeClient.CoreV1(), c.eventRecorder)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q: %v", "trusted-ca-bundle", err))
	}

	if operatorConfig.ObjectMeta.Generation != operatorConfig.Status.ObservedGeneration {
		reasonsForForcedRollingUpdate = append(reasonsForForcedRollingUpdate, fmt.Sprintf("operator config spec generation %d does not match status generation %d", operatorConfig.ObjectMeta.Generation,
			operatorConfig.Status.ObservedGeneration))
	}

	forceRollingUpdate := len(reasonsForForcedRollingUpdate) > 0

	if forceRollingUpdate {
		sort.Sort(sort.StringSlice(reasonsForForcedRollingUpdate))
		c.eventRecorder.Eventf("RollingUpdateForced", "Rolling update forced because:\n* %s", strings.Join(reasonsForForcedRollingUpdate, "\n"))
	}

	// our configmaps and secrets are in order, now it is time to create the DS
	// TODO check basic preconditions here
	actualDaemonSet, _, err := manageOpenShiftAPIServerDaemonSet_v311_00_to_latest(c.kubeClient.AppsV1(), c.eventRecorder, c.targetImagePullSpec, c.operatorImagePullSpec, operatorConfig, operatorConfig.Status.Generations, forceRollingUpdate)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q: %v", "daemonsets", err))
	}

	// only manage the apiservices if we have ready pods for the daemonset.  This makes sure that if we're taking over for
	// something else, we don't stomp their apiservices until ours have a reasonable chance at working.
	var actualAPIServices []*apiregistrationv1.APIService
	if actualDaemonSet != nil && actualDaemonSet.Status.NumberAvailable > 0 {
		actualAPIServices, err = manageAPIServices_v311_00_to_latest(c.apiregistrationv1Client, c.eventRecorder)
		if err != nil {
			errors = append(errors, fmt.Errorf("%q: %v", "apiservices", err))
		}
	}

	if actualDaemonSet == nil {
		v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorv1.OperatorCondition{
			Type:    daemonSetAvailableCondition,
			Status:  operatorv1.ConditionFalse,
			Message: "daemonset/apiserver.openshift-apiserver: could not be retrieved",
			Reason:  "NoDaemon",
		})
	} else {
		v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorv1.OperatorCondition{
			Type:   daemonSetAvailableCondition,
			Status: operatorv1.ConditionTrue,
		})
	}
	switch {
	case actualDaemonSet == nil:
		v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorv1.OperatorCondition{
			Type:   podAvailableCondition,
			Status: operatorv1.ConditionUnknown,
		})
	case actualDaemonSet.Status.NumberAvailable == 0:
		v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorv1.OperatorCondition{
			Type:    podAvailableCondition,
			Status:  operatorv1.ConditionFalse,
			Message: "no openshift-apiserver daemon pods available on any node.",
			Reason:  "NoAPIServerPod",
		})
	default:
		v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorv1.OperatorCondition{
			Type:   podAvailableCondition,
			Status: operatorv1.ConditionTrue,
		})
	}

	var apiServiceRegistrationMessages []string
	v1helpers.RemoveOperatorCondition(&operatorConfig.Status.Conditions, apiServicesAvailableCondition)
	for _, apiService := range actualAPIServices {
		for _, condition := range apiService.Status.Conditions {
			if condition.Type == apiregistrationv1.Available {
				if condition.Status == apiregistrationv1.ConditionFalse {
					apiServiceRegistrationMessages = append(apiServiceRegistrationMessages, fmt.Sprintf("apiservice registration/%v: not available: %v", apiService.Name, condition.Message))
				}
			}
		}
	}
	// if the apiservices themselves check out ok, try to actually hit the discovery endpoints.  We have a history in clusterup
	// of something delaying them.  This isn't perfect because of round-robining, but let's see if we get an improvement
	var missingAPIMessages []string
	if len(actualAPIServices) > 0 && len(apiServiceRegistrationMessages) == 0 && c.kubeClient.Discovery().RESTClient() != nil {
		missingAPIMessages = checkForAPIs(c.eventRecorder, c.kubeClient.Discovery().RESTClient(), apiServiceGroupVersions...)
	}

	switch {
	case actualDaemonSet == nil || actualDaemonSet.Status.NumberAvailable == 0:
		v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorv1.OperatorCondition{
			Type:   apiServicesAvailableCondition,
			Status: operatorv1.ConditionUnknown,
		})
	case len(actualAPIServices) == 0:
		v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorv1.OperatorCondition{
			Type:    apiServicesAvailableCondition,
			Status:  operatorv1.ConditionFalse,
			Message: "registered apiservices could not be retrieved",
			Reason:  "NoRegisteredAPIServices",
		})
	case len(missingAPIMessages) > 0:
		v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorv1.OperatorCondition{
			Type:    apiServicesAvailableCondition,
			Status:  operatorv1.ConditionFalse,
			Message: strings.Join(missingAPIMessages, "\n"),
			Reason:  "APIServiceNotReachable",
		})
	case len(apiServiceRegistrationMessages) > 0:
		sort.Strings(apiServiceRegistrationMessages)
		v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorv1.OperatorCondition{
			Type:    apiServicesAvailableCondition,
			Status:  operatorv1.ConditionFalse,
			Message: strings.Join(apiServiceRegistrationMessages, "\n"),
			Reason:  "APIServiceRegistrationNotAvailable",
		})
	default:
		v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorv1.OperatorCondition{
			Type:   apiServicesAvailableCondition,
			Status: operatorv1.ConditionTrue,
		})
	}

	// If the daemonset is up to date and the operatorConfig are up to date, then we are no longer progressing
	if actualDaemonSet != nil && actualDaemonSet.ObjectMeta.Generation != actualDaemonSet.Status.ObservedGeneration {
		v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorv1.OperatorCondition{
			Type:    daemonSetGenerationProgressingCondition,
			Status:  operatorv1.ConditionTrue,
			Reason:  "DesiredStateNotYetAchieved",
			Message: fmt.Sprintf("daemonset/apiserver.openshift-operator: observed generation is %d, desired generation is %d.", actualDaemonSet.Status.ObservedGeneration, actualDaemonSet.ObjectMeta.Generation),
		})
	}
	if operatorConfig.ObjectMeta.Generation != operatorConfig.Status.ObservedGeneration {
		v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorv1.OperatorCondition{
			Type:    operatorConfigGenerationProgressingCondition,
			Status:  operatorv1.ConditionTrue,
			Reason:  "DesiredStateNotYetAchieved",
			Message: fmt.Sprintf("openshiftapiserveroperatorconfigs/instance: observed generation is %d, desired generation is %d.", operatorConfig.Status.ObservedGeneration, operatorConfig.ObjectMeta.Generation),
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
			Type:    workloadDegradedCondition,
			Status:  operatorv1.ConditionTrue,
			Message: message,
			Reason:  "SyncError",
		})
	} else {
		v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorv1.OperatorCondition{
			Type:   workloadDegradedCondition,
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
	if !v1helpers.IsOperatorConditionFalse(operatorConfig.Status.Conditions, operatorv1.OperatorStatusTypeDegraded) {
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

// manageOpenShiftAPIServerImageImportCA_v311_00_to_latest synchronizes image import ca-bundle. Returns the modified
// ca-bundle ConfigMap.
func manageOpenShiftAPIServerImageImportCA_v311_00_to_latest(openshiftConfigClient openshiftconfigclientv1.ConfigV1Interface, client coreclientv1.CoreV1Interface, recorder events.Recorder) (*corev1.ConfigMap, bool, error) {
	imageConfig, err := openshiftConfigClient.Images().Get("cluster", metav1.GetOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		return nil, false, err
	}
	if apierrors.IsNotFound(err) {
		return nil, false, nil
	}
	if len(imageConfig.Spec.AdditionalTrustedCA.Name) == 0 {
		existingCM, err := client.ConfigMaps(operatorclient.TargetNamespace).Get(imageImportCAName, metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			return nil, false, nil
		}
		if err != nil {
			return nil, false, err
		}
		err = client.ConfigMaps(operatorclient.TargetNamespace).Delete(imageImportCAName, nil)
		if err != nil {
			return nil, false, err
		}
		return existingCM, true, nil
	}
	configMap, caChanged, err := resourceapply.SyncConfigMap(client, recorder, operatorclient.GlobalUserSpecifiedConfigNamespace, imageConfig.Spec.AdditionalTrustedCA.Name, operatorclient.TargetNamespace, imageImportCAName, nil)
	switch {
	case err != nil:
		return nil, false, err
	case caChanged:
		return configMap, true, nil
	default:
		return nil, false, nil
	}
}

func syncOpenShiftAPIServerTrustedCA_v311_00_to_latest(client coreclientv1.CoreV1Interface, recorder events.Recorder) error {
	// get the CM created by the CVO in the operator namespace
	sourceCM, err := client.ConfigMaps(operatorclient.OperatorNamespace).Get("trusted-ca-bundle", metav1.GetOptions{})
	if err != nil {
		return err
	}

	// it's optional, don't sync it if it has no data -> prevent trust store being removed upon startup/upgrade
	if len(sourceCM.Data) == 0 {
		return nil
	}

	// get the copy of the sourceCM, remove the annotations and labels we set for CVO and network operator
	requiredCM := sourceCM.DeepCopy()
	requiredCM.UID = ""
	requiredCM.ObjectMeta.ResourceVersion = ""
	requiredCM.ObjectMeta.Namespace = operatorclient.TargetNamespace
	delete(requiredCM.Annotations, "release.openshift.io/create-only")
	delete(requiredCM.Labels, "config.openshift.io/inject-trusted-cabundle")

	targetClient := client.ConfigMaps(operatorclient.TargetNamespace)
	imageCM, err := targetClient.Get(requiredCM.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		_, err = targetClient.Create(requiredCM)
		return err
	}

	if err != nil {
		return err
	}

	// if the target CM Data differs, update it
	if !reflect.DeepEqual(sourceCM.Data, imageCM.Data) {
		requiredCM.UID = imageCM.UID
		requiredCM.ResourceVersion = imageCM.ResourceVersion
		_, err = targetClient.Update(requiredCM)
	}

	return err
}

func manageOpenShiftAPIServerConfigMap_v311_00_to_latest(kubeClient kubernetes.Interface, client coreclientv1.ConfigMapsGetter, recorder events.Recorder, operatorConfig *operatorv1.OpenShiftAPIServer) (*corev1.ConfigMap, bool, error) {
	configMap := resourceread.ReadConfigMapV1OrDie(v311_00_assets.MustAsset("v3.11.0/openshift-apiserver/cm.yaml"))
	defaultConfig := v311_00_assets.MustAsset("v3.11.0/openshift-apiserver/defaultconfig.yaml")
	requiredConfigMap, _, err := resourcemerge.MergePrunedConfigMap(
		&openshiftcontrolplanev1.OpenShiftAPIServerConfig{},
		configMap,
		"config.yaml",
		nil,
		defaultConfig,
		operatorConfig.Spec.ObservedConfig.Raw,
		operatorConfig.Spec.UnsupportedConfigOverrides.Raw,
	)
	if err != nil {
		return nil, false, err
	}

	// we can embed input hashes on our main configmap to drive rollouts when they change.
	inputHashes, err := resourcehash.MultipleObjectHashStringMapForObjectReferences(
		kubeClient,
		resourcehash.NewObjectRef().ForSecret().InNamespace(operatorclient.TargetNamespace).Named("etcd-client"),
		resourcehash.NewObjectRef().ForConfigMap().InNamespace(operatorclient.TargetNamespace).Named("etcd-serving-ca"),
		resourcehash.NewObjectRef().ForConfigMap().InNamespace(operatorclient.TargetNamespace).Named("trusted-ca-bundle"),
	)
	if err != nil {
		return nil, false, err
	}
	for k, v := range inputHashes {
		requiredConfigMap.Data[k] = v
	}

	return resourceapply.ApplyConfigMap(client, recorder, requiredConfigMap)
}

func manageOpenShiftAPIServerDaemonSet_v311_00_to_latest(
	client appsclientv1.DaemonSetsGetter,
	recorder events.Recorder,
	imagePullSpec string,
	operatorImagePullSpec string,
	operatorConfig *operatorv1.OpenShiftAPIServer,
	generationStatus []operatorv1.GenerationStatus,
	forceRollingUpdate bool) (*appsv1.DaemonSet, bool, error) {
	tmpl := v311_00_assets.MustAsset("v3.11.0/openshift-apiserver/ds.yaml")

	r := strings.NewReplacer(
		"${IMAGE}", imagePullSpec,
		"${OPERATOR_IMAGE}", operatorImagePullSpec,
		"${REVISION}", strconv.Itoa(int(operatorConfig.Status.LatestAvailableRevision)),
	)
	tmpl = []byte(r.Replace(string(tmpl)))

	re := regexp.MustCompile("\\$\\{[^}]*}")
	if match := re.Find(tmpl); len(match) > 0 {
		return nil, false, fmt.Errorf("invalid template reference %q", string(match))
	}

	required := resourceread.ReadDaemonSetV1OrDie(tmpl)

	// we set this so that when the requested image pull spec changes, we always have a diff.  Remember that we don't directly
	// diff any fields on the daemonset because they can be rewritten by admission and we don't want to constantly be fighting
	// against admission or defaults.  That was a problem with original versions of apply.
	if required.Annotations == nil {
		required.Annotations = map[string]string{}
	}
	required.Annotations["openshiftapiservers.operator.openshift.io/pull-spec"] = imagePullSpec
	required.Annotations["openshiftapiservers.operator.openshift.io/operator-pull-spec"] = operatorImagePullSpec

	required.Labels["revision"] = strconv.Itoa(int(operatorConfig.Status.LatestAvailableRevision))
	required.Spec.Template.Labels["revision"] = strconv.Itoa(int(operatorConfig.Status.LatestAvailableRevision))

	containerArgsWithLoglevel := required.Spec.Template.Spec.Containers[0].Args
	if argsCount := len(containerArgsWithLoglevel); argsCount > 1 {
		return nil, false, fmt.Errorf("expected only one container argument, got %d", argsCount)
	}
	if !strings.Contains(containerArgsWithLoglevel[0], "exec openshift-apiserver") {
		return nil, false, fmt.Errorf("exec openshift-apiserver not found in first argument %q", containerArgsWithLoglevel[0])
	}
	containerArgsWithLoglevel[0] = strings.TrimSpace(containerArgsWithLoglevel[0])
	switch operatorConfig.Spec.LogLevel {
	case operatorv1.Normal:
		containerArgsWithLoglevel[0] += fmt.Sprintf(" -v=%d", 2)
	case operatorv1.Debug:
		containerArgsWithLoglevel[0] += fmt.Sprintf(" -v=%d", 4)
	case operatorv1.Trace:
		containerArgsWithLoglevel[0] += fmt.Sprintf(" -v=%d", 6)
	case operatorv1.TraceAll:
		containerArgsWithLoglevel[0] += fmt.Sprintf(" -v=%d", 8)
	default:
		containerArgsWithLoglevel[0] += fmt.Sprintf(" -v=%d", 2)
	}
	required.Spec.Template.Spec.Containers[0].Args = containerArgsWithLoglevel

	var observedConfig map[string]interface{}
	if err := yaml.Unmarshal(operatorConfig.Spec.ObservedConfig.Raw, &observedConfig); err != nil {
		return nil, false, fmt.Errorf("failed to unmarshal the observedConfig: %v", err)
	}
	proxyConfig, _, err := unstructured.NestedStringMap(observedConfig, "workloadcontroller", "proxy")
	if err != nil {
		return nil, false, fmt.Errorf("couldn't get the proxy config from observedConfig: %v", err)
	}

	proxyEnvVars := proxyMapToEnvVars(proxyConfig)
	for i, container := range required.Spec.Template.Spec.Containers {
		required.Spec.Template.Spec.Containers[i].Env = append(container.Env, proxyEnvVars...)
	}

	return resourceapply.ApplyDaemonSet(client, recorder, required, resourcemerge.ExpectedDaemonSetGeneration(required, generationStatus), forceRollingUpdate)
}

func manageAPIServices_v311_00_to_latest(client apiregistrationv1client.APIServicesGetter, recorder events.Recorder) ([]*apiregistrationv1.APIService, error) {
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
		apiregistrationv1.SetDefaults_ServiceReference(obj.Spec.Service)

		apiService, _, err := resourceapply.ApplyAPIService(client, recorder, obj)
		if err != nil {
			return nil, err
		}
		apiServices = append(apiServices, apiService)
	}

	return apiServices, nil
}

var openshiftScheme = runtime.NewScheme()

func init() {
	if err := openshiftapi.Install(openshiftScheme); err != nil {
		panic(err)
	}
}

func resourceSelectorForCLI(obj runtime.Object) string {
	groupVersionKind := obj.GetObjectKind().GroupVersionKind()
	if len(groupVersionKind.Kind) == 0 {
		if kinds, _, _ := scheme.Scheme.ObjectKinds(obj); len(kinds) > 0 {
			groupVersionKind = kinds[0]
		}
	}
	if len(groupVersionKind.Kind) == 0 {
		if kinds, _, _ := openshiftScheme.ObjectKinds(obj); len(kinds) > 0 {
			groupVersionKind = kinds[0]
		}
	}
	if len(groupVersionKind.Kind) == 0 {
		groupVersionKind = schema.GroupVersionKind{Kind: "Unknown"}
	}
	kind := groupVersionKind.Kind
	group := groupVersionKind.Group
	var name string
	accessor, err := meta.Accessor(obj)
	if err != nil {
		name = "unknown"
	}
	name = accessor.GetName()
	if len(group) > 0 {
		group = "." + group
	}
	return kind + group + "/" + name
}

func proxyMapToEnvVars(proxyConfig map[string]string) []corev1.EnvVar {
	if proxyConfig == nil {
		return nil
	}

	envVars := []corev1.EnvVar{}
	for k, v := range proxyConfig {
		envVars = append(envVars, corev1.EnvVar{Name: k, Value: v})
	}

	// sort the env vars to prevent update hotloops
	sort.Slice(envVars, func(i, j int) bool { return envVars[i].Name < envVars[j].Name })
	return envVars
}
