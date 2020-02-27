package workloadcontroller

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	utilerrors "k8s.io/apimachinery/pkg/util/errors"

	"github.com/ghodss/yaml"
	"github.com/google/uuid"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	appsclientv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	coreclientv1 "k8s.io/client-go/kubernetes/typed/core/v1"

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
	"github.com/openshift/library-go/pkg/operator/v1helpers"
)

// syncOpenShiftAPIServer_v311_00_to_latest takes care of synchronizing (not upgrading) the thing we're managing.
// most of the time the sync method will be good for a large span of minor versions
func syncOpenShiftAPIServer_v311_00_to_latest(c OpenShiftAPIServerOperator, originalOperatorConfig *operatorv1.OpenShiftAPIServer) error {
	errors := []error{}
	operatorConfig := originalOperatorConfig.DeepCopy()

	_, _, err := manageOpenShiftAPIServerConfigMap_v311_00_to_latest(c.kubeClient, c.kubeClient.CoreV1(), c.eventRecorder, operatorConfig)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q: %v", "configmap", err))
	}

	_, _, err = manageOpenShiftAPIServerImageImportCA_v311_00_to_latest(c.openshiftConfigClient, c.kubeClient.CoreV1(), c.eventRecorder)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q: %v", "image-import-ca", err))
	}

	// our configmaps and secrets are in order, now it is time to create the deployment
	// TODO check basic preconditions here
	actualDeployment, _, err := manageOpenShiftAPIServerDeployment_v311_00_to_latest(c.kubeClient, c.kubeClient.AppsV1(), c.eventRecorder, c.targetImagePullSpec, c.operatorImagePullSpec, operatorConfig, operatorConfig.Status.Generations)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q: %v", "deployments", err))
	}

	if len(errors) > 0 {
		message := ""
		for _, err := range errors {
			message = message + err.Error() + "\n"
		}
		errors = append(errors,
			appendErrors(v1helpers.UpdateStatus(c.operatorClient, v1helpers.UpdateConditionFn(operatorv1.OperatorCondition{
				Type:    "WorkloadDegraded",
				Status:  operatorv1.ConditionTrue,
				Reason:  "SyncError",
				Message: message,
			})))...,
		)
	} else {
		errors = append(errors,
			appendErrors(v1helpers.UpdateStatus(c.operatorClient, v1helpers.UpdateConditionFn(operatorv1.OperatorCondition{
				Type:   "WorkloadDegraded",
				Status: operatorv1.ConditionFalse,
			})))...,
		)
	}
	if operatorConfig.ObjectMeta.Generation != operatorConfig.Status.ObservedGeneration {
		errors = append(errors,
			appendErrors(v1helpers.UpdateStatus(c.operatorClient, v1helpers.UpdateConditionFn(operatorv1.OperatorCondition{
				Type:    "OperatorConfigProgressing",
				Status:  operatorv1.ConditionTrue,
				Reason:  "NewGeneration",
				Message: fmt.Sprintf("openshiftapiserveroperatorconfigs/instance: observed generation is %d, desired generation is %d.", operatorConfig.Status.ObservedGeneration, operatorConfig.ObjectMeta.Generation),
			})))...,
		)
	} else {
		errors = append(errors,
			appendErrors(v1helpers.UpdateStatus(c.operatorClient, v1helpers.UpdateConditionFn(operatorv1.OperatorCondition{
				Type:   "OperatorConfigProgressing",
				Status: operatorv1.ConditionFalse,
				Reason: "AsExpected",
			})))...,
		)
	}

	if actualDeployment == nil {
		message := "deployment/apiserver.openshift-apiserver: could not be retrieved"
		errors = append(errors,
			appendErrors(v1helpers.UpdateStatus(c.operatorClient, v1helpers.UpdateConditionFn(operatorv1.OperatorCondition{
				Type:    "APIServerDeploymentAvailable",
				Status:  operatorv1.ConditionFalse,
				Reason:  "NoDeployment",
				Message: message,
			})))...,
		)
		errors = append(errors,
			appendErrors(v1helpers.UpdateStatus(c.operatorClient, v1helpers.UpdateConditionFn(operatorv1.OperatorCondition{
				Type:    "APIServerDeploymentProgressing",
				Status:  operatorv1.ConditionTrue,
				Reason:  "NoDeployment",
				Message: message,
			})))...,
		)
		errors = append(errors,
			appendErrors(v1helpers.UpdateStatus(c.operatorClient, v1helpers.UpdateConditionFn(operatorv1.OperatorCondition{
				Type:    "APIServerDeploymentDegraded",
				Status:  operatorv1.ConditionTrue,
				Reason:  "NoDeployment",
				Message: message,
			})))...,
		)

		return utilerrors.NewAggregate(errors)
	}

	// manage status
	if actualDeployment.Status.AvailableReplicas == 0 {
		errors = append(errors,
			appendErrors(v1helpers.UpdateStatus(c.operatorClient, v1helpers.UpdateConditionFn(operatorv1.OperatorCondition{
				Type:    "APIServerDeploymentAvailable",
				Status:  operatorv1.ConditionFalse,
				Reason:  "NoAPIServerPod",
				Message: "no openshift-apiserver pods available.",
			})))...,
		)
	} else {
		errors = append(errors,
			appendErrors(v1helpers.UpdateStatus(c.operatorClient, v1helpers.UpdateConditionFn(operatorv1.OperatorCondition{
				Type:   "APIServerDeploymentAvailable",
				Status: operatorv1.ConditionTrue,
				Reason: "AsExpected",
			})))...,
		)
	}

	// If the deployment is up to date and the operatorConfig are up to date, then we are no longer progressing
	deploymentAtHighestGeneration := actualDeployment.ObjectMeta.Generation == actualDeployment.Status.ObservedGeneration
	if !deploymentAtHighestGeneration {
		errors = append(errors,
			appendErrors(v1helpers.UpdateStatus(c.operatorClient, v1helpers.UpdateConditionFn(operatorv1.OperatorCondition{
				Type:    "APIServerDeploymentProgressing",
				Status:  operatorv1.ConditionTrue,
				Reason:  "NewGeneration",
				Message: fmt.Sprintf("deployment/apiserver.openshift-operator: observed generation is %d, desired generation is %d.", actualDeployment.Status.ObservedGeneration, actualDeployment.ObjectMeta.Generation),
			})))...,
		)
	} else {
		errors = append(errors,
			appendErrors(v1helpers.UpdateStatus(c.operatorClient, v1helpers.UpdateConditionFn(operatorv1.OperatorCondition{
				Type:   "APIServerDeploymentProgressing",
				Status: operatorv1.ConditionFalse,
				Reason: "AsExpected",
			})))...,
		)
	}

	// TODO this is changing too early and it was before too.
	errors = append(errors,
		appendErrors(v1helpers.UpdateStatus(c.operatorClient, func(status *operatorv1.OperatorStatus) error {
			status.ObservedGeneration = operatorConfig.ObjectMeta.Generation
			return nil
		}))...,
	)
	errors = append(errors,
		appendErrors(v1helpers.UpdateStatus(c.operatorClient, func(status *operatorv1.OperatorStatus) error {
			resourcemerge.SetDeploymentGeneration(&status.Generations, actualDeployment)
			return nil
		}))...,
	)

	desiredReplicas := int32(1)
	if actualDeployment.Spec.Replicas != nil {
		desiredReplicas = *(actualDeployment.Spec.Replicas)
	}

	// During a rollout with default maxSurge (25%) will allow the available
	// replicas to temporarily exceed the desired replica count.
	deploymentHasAllPodsAvailable := actualDeployment.Status.AvailableReplicas >= desiredReplicas
	if !deploymentHasAllPodsAvailable {
		numNonAvailablePods := desiredReplicas - actualDeployment.Status.AvailableReplicas
		errors = append(errors,
			appendErrors(v1helpers.UpdateStatus(c.operatorClient, v1helpers.UpdateConditionFn(operatorv1.OperatorCondition{
				Type:    "APIServerDeploymentDegraded",
				Status:  operatorv1.ConditionTrue,
				Reason:  "UnavailablePod",
				Message: fmt.Sprintf("%v of %v requested instances are unavailable", numNonAvailablePods, desiredReplicas),
			})))...,
		)
	} else {
		errors = append(errors,
			appendErrors(v1helpers.UpdateStatus(c.operatorClient, v1helpers.UpdateConditionFn(operatorv1.OperatorCondition{
				Type:   "APIServerDeploymentDegraded",
				Status: operatorv1.ConditionFalse,
				Reason: "AsExpected",
			})))...,
		)
	}

	// if the deployment is all available and at the expected generation, then update the version to the latest
	// when we update, the image pull spec should immediately be different, which should immediately cause a deployment rollout
	// which should immediately result in a deployment generation diff, which should cause this block to be skipped until it is ready.
	deploymentHasAllPodsUpdated := actualDeployment.Status.UpdatedReplicas == desiredReplicas
	operatorConfigAtHighestGeneration := operatorConfig.Status.ObservedGeneration == operatorConfig.ObjectMeta.Generation
	if operatorConfigAtHighestGeneration && deploymentAtHighestGeneration && deploymentHasAllPodsAvailable && deploymentHasAllPodsUpdated {
		c.versionRecorder.SetVersion("openshift-apiserver", c.targetOperandVersion)
	}

	return utilerrors.NewAggregate(errors)
}

func appendErrors(_ *operatorv1.OperatorStatus, _ bool, err error) []error {
	if err != nil {
		return []error{err}
	}
	return []error{}
}

// mergeImageRegistryCertificates merges two distinct ConfigMap, both containing
// trusted CAs for Image Registries. The first one is the default CA bundle for
// OpenShift internal registry access, the latter is a custom config map that may
// be configured by the user on image.config.openshift.io/cluster.
func mergeImageRegistryCertificates(cfgCli openshiftconfigclientv1.ConfigV1Interface, cli coreclientv1.CoreV1Interface) (map[string]string, error) {
	cas := make(map[string]string)

	internalRegistryCAs, err := cli.ConfigMaps("openshift-image-registry").Get(
		"image-registry-certificates", metav1.GetOptions{},
	)
	if err != nil && !apierrors.IsNotFound(err) {
		return nil, err
	} else if err == nil {
		for key, value := range internalRegistryCAs.Data {
			cas[key] = value
		}
	}

	imageConfig, err := cfgCli.Images().Get(
		"cluster", metav1.GetOptions{},
	)
	if err != nil {
		return nil, err
	}

	// No custom config map, return.
	if len(imageConfig.Spec.AdditionalTrustedCA.Name) == 0 {
		return cas, nil
	}

	additionalImageRegistryCAs, err := cli.ConfigMaps(
		operatorclient.GlobalUserSpecifiedConfigNamespace,
	).Get(
		imageConfig.Spec.AdditionalTrustedCA.Name,
		metav1.GetOptions{},
	)
	if err != nil && !apierrors.IsNotFound(err) {
		return nil, err
	} else if err == nil {
		for key, value := range additionalImageRegistryCAs.Data {
			cas[key] = value
		}
	}
	return cas, nil
}

// manageOpenShiftAPIServerImageImportCA_v311_00_to_latest synchronizes image import ca-bundle. Returns the modified
// ca-bundle ConfigMap.
func manageOpenShiftAPIServerImageImportCA_v311_00_to_latest(openshiftConfigClient openshiftconfigclientv1.ConfigV1Interface, client coreclientv1.CoreV1Interface, recorder events.Recorder) (*corev1.ConfigMap, bool, error) {
	mergedCAs, err := mergeImageRegistryCertificates(openshiftConfigClient, client)
	if err != nil {
		return nil, false, err
	}
	requiredConfigMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: operatorclient.TargetNamespace,
			Name:      imageImportCAName,
		},
		Data: mergedCAs,
	}

	// this can leave configmaps mounted without any content, but that should not have an impact on functionality since empty and missing
	// should logically be treated the same in the case of trust.
	return resourceapply.ApplyConfigMap(client, recorder, requiredConfigMap)
}

func manageOpenShiftAPIServerConfigMap_v311_00_to_latest(kubeClient kubernetes.Interface, client coreclientv1.ConfigMapsGetter, recorder events.Recorder, operatorConfig *operatorv1.OpenShiftAPIServer) (*corev1.ConfigMap, bool, error) {
	configMap := resourceread.ReadConfigMapV1OrDie(v311_00_assets.MustAsset("v3.11.0/openshift-apiserver/cm.yaml"))
	defaultConfig := v311_00_assets.MustAsset("v3.11.0/config/defaultconfig.yaml")
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

	return resourceapply.ApplyConfigMap(client, recorder, requiredConfigMap)
}

func loglevelToKlog(logLevel operatorv1.LogLevel) string {
	switch logLevel {
	case operatorv1.Normal:
		return "2"
	case operatorv1.Debug:
		return "4"
	case operatorv1.Trace:
		return "6"
	case operatorv1.TraceAll:
		return "8"
	default:
		return "2"
	}
}

func manageOpenShiftAPIServerDeployment_v311_00_to_latest(
	kubeClient kubernetes.Interface,
	client appsclientv1.DeploymentsGetter,
	recorder events.Recorder,
	imagePullSpec string,
	operatorImagePullSpec string,
	operatorConfig *operatorv1.OpenShiftAPIServer,
	generationStatus []operatorv1.GenerationStatus,
) (*appsv1.Deployment, bool, error) {
	tmpl := v311_00_assets.MustAsset("v3.11.0/openshift-apiserver/deploy.yaml")

	r := strings.NewReplacer(
		"${IMAGE}", imagePullSpec,
		"${OPERATOR_IMAGE}", operatorImagePullSpec,
		"${REVISION}", strconv.Itoa(int(operatorConfig.Status.LatestAvailableRevision)),
		"${VERBOSITY}", loglevelToKlog(operatorConfig.Spec.LogLevel),
	)
	tmpl = []byte(r.Replace(string(tmpl)))

	re := regexp.MustCompile("\\$\\{[^}]*}")
	if match := re.Find(tmpl); len(match) > 0 {
		return nil, false, fmt.Errorf("invalid template reference %q", string(match))
	}

	required := resourceread.ReadDeploymentV1OrDie(tmpl)

	// we set this so that when the requested image pull spec changes, we always have a diff.  Remember that we don't directly
	// diff any fields on the deployment because they can be rewritten by admission and we don't want to constantly be fighting
	// against admission or defaults.  That was a problem with original versions of apply.
	if required.Annotations == nil {
		required.Annotations = map[string]string{}
	}
	required.Annotations["openshiftapiservers.operator.openshift.io/pull-spec"] = imagePullSpec
	required.Annotations["openshiftapiservers.operator.openshift.io/operator-pull-spec"] = operatorImagePullSpec

	required.Labels["revision"] = strconv.Itoa(int(operatorConfig.Status.LatestAvailableRevision))
	required.Spec.Template.Labels["revision"] = strconv.Itoa(int(operatorConfig.Status.LatestAvailableRevision))

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

	// we watch some resources so that our deployment will redeploy without explicitly and carefully ordered resource creation
	inputHashes, err := resourcehash.MultipleObjectHashStringMapForObjectReferences(
		kubeClient,
		resourcehash.NewObjectRef().ForConfigMap().InNamespace(operatorclient.TargetNamespace).Named("config"),
		resourcehash.NewObjectRef().ForSecret().InNamespace(operatorclient.TargetNamespace).Named("etcd-client"),
		resourcehash.NewObjectRef().ForConfigMap().InNamespace(operatorclient.TargetNamespace).Named("etcd-serving-ca"),
		resourcehash.NewObjectRef().ForConfigMap().InNamespace(operatorclient.TargetNamespace).Named("image-import-ca"),
		resourcehash.NewObjectRef().ForConfigMap().InNamespace(operatorclient.TargetNamespace).Named("trusted-ca-bundle"),
	)
	if err != nil {
		return nil, false, fmt.Errorf("invalid dependency reference: %q", err)
	}
	inputHashes["desired.generation"] = fmt.Sprintf("%d", operatorConfig.ObjectMeta.Generation)
	for k, v := range inputHashes {
		annotationKey := fmt.Sprintf("operator.openshift.io/dep-%s", k)
		required.Annotations[annotationKey] = v
		if required.Spec.Template.Annotations == nil {
			required.Spec.Template.Annotations = map[string]string{}
		}
		required.Spec.Template.Annotations[annotationKey] = v
	}

	err = ensureAtMostOnePodPerNode(&required.Spec)
	if err != nil {
		return nil, false, fmt.Errorf("unable to ensure at most one pod per node: %v", err)
	}

	return resourceapply.ApplyDeployment(client, recorder, required, resourcemerge.ExpectedDeploymentGeneration(required, generationStatus), false)
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

// ensureAtMostOnePodPerNode updates the deployment spec to prevent more than
// one pod of a given replicaset from landing on a node. It accomplishes this
// by adding a uuid as a label on the template and updates the pod
// anti-affinity term to include that label. Since the deployment is only
// written (via ApplyDeployment) when the metadata differs or the generations
// don't match, the uuid should only be updated in the API when a new
// replicaset is created.
func ensureAtMostOnePodPerNode(spec *appsv1.DeploymentSpec) error {
	uuidKey := "anti-affinity-uuid"
	uuidValue := uuid.New().String()

	// Label the pod template with the template hash
	spec.Template.Labels[uuidKey] = uuidValue

	// Ensure that match labels are defined
	if spec.Selector == nil {
		return fmt.Errorf("deployment is missing spec.selector")
	}
	if len(spec.Selector.MatchLabels) == 0 {
		return fmt.Errorf("deployment is missing spec.selector.matchLabels")
	}

	// Ensure anti-affinity selects on the uuid
	antiAffinityMatchLabels := map[string]string{
		uuidKey: uuidValue,
	}
	// Ensure anti-affinity selects on the same labels as the deployment
	for key, value := range spec.Selector.MatchLabels {
		antiAffinityMatchLabels[key] = value
	}

	// Add an anti-affinity rule to the pod template that precludes more than
	// one pod for a uuid from being scheduled to a node.
	spec.Template.Spec.Affinity = &corev1.Affinity{
		PodAntiAffinity: &corev1.PodAntiAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
				{
					TopologyKey: "kubernetes.io/hostname",
					LabelSelector: &metav1.LabelSelector{
						MatchLabels: antiAffinityMatchLabels,
					},
				},
			},
		},
	}

	return nil
}
