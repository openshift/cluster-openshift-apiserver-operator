package workload

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	appsclientv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	coreclientv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/klog/v2"

	openshiftapi "github.com/openshift/api"
	openshiftcontrolplanev1 "github.com/openshift/api/openshiftcontrolplane/v1"
	operatorv1 "github.com/openshift/api/operator/v1"
	openshiftconfigclientv1 "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"
	operatorv1client "github.com/openshift/client-go/operator/clientset/versioned/typed/operator/v1"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/v311_00_assets"
	"github.com/openshift/library-go/pkg/controller/factory"
	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/resource/resourceapply"
	"github.com/openshift/library-go/pkg/operator/resource/resourcehash"
	"github.com/openshift/library-go/pkg/operator/resource/resourcemerge"
	"github.com/openshift/library-go/pkg/operator/resource/resourceread"
	"github.com/openshift/library-go/pkg/operator/status"
	"github.com/openshift/library-go/pkg/operator/v1helpers"
)

const (
	imageImportCAName = "image-import-ca"
)

// nodeCountFunction a function to return count of nodes
type nodeCountFunc func(nodeSelector map[string]string) (*int32, error)

// ensureAtMostOnePodPerNode a function that updates the deployment spec to prevent more than
// one pod of a given replicaset from landing on a node.
type ensureAtMostOnePodPerNodeFunc func(spec *appsv1.DeploymentSpec, component string) error

// OpenShiftAPIServerWorkload is a struct that holds necessary data to install OpenShiftAPIServer
type OpenShiftAPIServerWorkload struct {
	operatorClient        v1helpers.OperatorClient
	operatorConfigClient  operatorv1client.OpenShiftAPIServersGetter
	openshiftConfigClient openshiftconfigclientv1.ConfigV1Interface
	kubeClient            kubernetes.Interface

	// countNodes a function to return count of nodes on which the workload will be installed
	countNodes nodeCountFunc

	// ensureAtMostOnePodPerNode a function that updates the deployment spec to prevent more than
	// one pod of a given replicaset from landing on a node.
	ensureAtMostOnePodPerNode ensureAtMostOnePodPerNodeFunc

	targetNamespace       string
	targetImagePullSpec   string
	operatorImagePullSpec string

	versionRecorder status.VersionGetter
}

// NewOpenShiftAPIServerWorkload creates new OpenShiftAPIServerWorkload struct
func NewOpenShiftAPIServerWorkload(
	operatorClient v1helpers.OperatorClient,
	operatorConfigClient operatorv1client.OpenShiftAPIServersGetter,
	openshiftConfigClient openshiftconfigclientv1.ConfigV1Interface,
	countNodes nodeCountFunc,
	ensureAtMostOnePodPerNode ensureAtMostOnePodPerNodeFunc,
	targetNamespace string,
	targetImagePullSpec string,
	operatorImagePullSpec string,
	kubeClient kubernetes.Interface,
	versionRecorder status.VersionGetter,
) *OpenShiftAPIServerWorkload {
	return &OpenShiftAPIServerWorkload{
		operatorClient:            operatorClient,
		operatorConfigClient:      operatorConfigClient,
		openshiftConfigClient:     openshiftConfigClient,
		countNodes:                countNodes,
		ensureAtMostOnePodPerNode: ensureAtMostOnePodPerNode,
		targetNamespace:           targetNamespace,
		targetImagePullSpec:       targetImagePullSpec,
		operatorImagePullSpec:     operatorImagePullSpec,
		kubeClient:                kubeClient,
		versionRecorder:           versionRecorder,
	}
}

// PreconditionFulfilled is a function that indicates whether all prerequisites are met and we can Sync.
func (c *OpenShiftAPIServerWorkload) PreconditionFulfilled(ctx context.Context) (bool, error) {
	operatorConfig, err := c.operatorConfigClient.OpenShiftAPIServers().Get(ctx, "cluster", metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	return c.preconditionFulfilledInternal(operatorConfig)
}

func (c *OpenShiftAPIServerWorkload) preconditionFulfilledInternal(operator *operatorv1.OpenShiftAPIServer) (bool, error) {
	// block until config is obvserved
	if len(operator.Spec.ObservedConfig.Raw) == 0 {
		klog.Info("Waiting for observed configuration to be available")
		return false, errors.New("waiting for observed configuration to be available")
	}

	// specifying etcd servers list is mandatory, without it the pods will be crashlooping, so wait for it.
	// TODO: once OAS moves to using args directly (not via the config)
	//       we could have a helper function for checking mandatory args for oauth-apiserver and openshift-apiserver
	//       alternatively we could change the generic workload ctrl to accept a list of common preconditions

	var typedObservedConfig openshiftcontrolplanev1.OpenShiftAPIServerConfig
	if err := json.Unmarshal(operator.Spec.ObservedConfig.Raw, &typedObservedConfig); err != nil {
		return false, err
	}

	if len(typedObservedConfig.StorageConfig.URLs) == 0 {
		klog.Info("Waiting for observed configuration to have mandatory StorageConfig.URLs")
		return false, errors.New("waiting for observed configuration to have mandatory StorageConfig.URLs")
	}

	return true, nil
}

// Sync takes care of synchronizing (not upgrading) the thing we're managing.
// most of the time the sync method will be good for a large span of minor versions
func (c *OpenShiftAPIServerWorkload) Sync(ctx context.Context, syncContext factory.SyncContext) (*appsv1.Deployment, bool, []error) {
	errors := []error{}

	originalOperatorConfig, err := c.operatorConfigClient.OpenShiftAPIServers().Get(ctx, "cluster", metav1.GetOptions{})
	if err != nil {
		errors = append(errors, err)
		return nil, false, errors
	}
	operatorConfig := originalOperatorConfig.DeepCopy()

	_, _, err = manageOpenShiftAPIServerConfigMap_v311_00_to_latest(ctx, c.kubeClient.CoreV1(), syncContext.Recorder(), operatorConfig)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q: %v", "configmap", err))
	}

	_, _, err = manageOpenShiftAPIServerImageImportCA_v311_00_to_latest(ctx, c.openshiftConfigClient, c.kubeClient.CoreV1(), syncContext.Recorder())
	if err != nil {
		errors = append(errors, fmt.Errorf("%q: %v", "image-import-ca", err))
	}

	// our configmaps and secrets are in order, now it is time to create the deployment
	// TODO check basic preconditions here
	actualDeployment, _, err := manageOpenShiftAPIServerDeployment_v311_00_to_latest(
		ctx,
		c.kubeClient,
		c.kubeClient.AppsV1(),
		c.countNodes,
		syncContext.Recorder(),
		c.targetImagePullSpec,
		c.operatorImagePullSpec,
		operatorConfig,
		operatorConfig.Status.Generations,
		c.ensureAtMostOnePodPerNode)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q: %v", "deployments", err))
	}

	if operatorConfig.ObjectMeta.Generation != operatorConfig.Status.ObservedGeneration {
		handleErrorForOperatorStatus(v1helpers.UpdateStatus(c.operatorClient, v1helpers.UpdateConditionFn(operatorv1.OperatorCondition{
			Type:    "OperatorConfigProgressing",
			Status:  operatorv1.ConditionTrue,
			Reason:  "NewGeneration",
			Message: fmt.Sprintf("openshiftapiserveroperatorconfigs/instance: observed generation is %d, desired generation is %d.", operatorConfig.Status.ObservedGeneration, operatorConfig.ObjectMeta.Generation),
		})),
		)
	} else {
		handleErrorForOperatorStatus(v1helpers.UpdateStatus(c.operatorClient, v1helpers.UpdateConditionFn(operatorv1.OperatorCondition{
			Type:   "OperatorConfigProgressing",
			Status: operatorv1.ConditionFalse,
			Reason: "AsExpected",
		})),
		)
	}

	// TODO this is changing too early and it was before too.
	handleErrorForOperatorStatus(v1helpers.UpdateStatus(c.operatorClient, func(status *operatorv1.OperatorStatus) error {
		status.ObservedGeneration = operatorConfig.ObjectMeta.Generation
		return nil
	}),
	)
	handleErrorForOperatorStatus(v1helpers.UpdateStatus(c.operatorClient, func(status *operatorv1.OperatorStatus) error {
		resourcemerge.SetDeploymentGeneration(&status.Generations, actualDeployment)
		return nil
	}),
	)

	return actualDeployment, operatorConfig.Status.ObservedGeneration == operatorConfig.ObjectMeta.Generation, errors
}

// mergeImageRegistryCertificates merges two distinct ConfigMap, both containing
// trusted CAs for Image Registries. The first one is the default CA bundle for
// OpenShift internal registry access, the latter is a custom config map that may
// be configured by the user on image.config.openshift.io/cluster.
func mergeImageRegistryCertificates(ctx context.Context, cfgCli openshiftconfigclientv1.ConfigV1Interface, cli coreclientv1.CoreV1Interface) (map[string]string, error) {
	cas := make(map[string]string)

	internalRegistryCAs, err := cli.ConfigMaps("openshift-image-registry").Get(
		ctx, "image-registry-certificates", metav1.GetOptions{},
	)
	if err != nil && !apierrors.IsNotFound(err) {
		return nil, err
	} else if err == nil {
		for key, value := range internalRegistryCAs.Data {
			cas[key] = value
		}
	}

	imageConfig, err := cfgCli.Images().Get(
		ctx, "cluster", metav1.GetOptions{},
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
		ctx,
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
func manageOpenShiftAPIServerImageImportCA_v311_00_to_latest(ctx context.Context, openshiftConfigClient openshiftconfigclientv1.ConfigV1Interface, client coreclientv1.CoreV1Interface, recorder events.Recorder) (*corev1.ConfigMap, bool, error) {
	mergedCAs, err := mergeImageRegistryCertificates(ctx, openshiftConfigClient, client)
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
	return resourceapply.ApplyConfigMap(ctx, client, recorder, requiredConfigMap)
}

func manageOpenShiftAPIServerConfigMap_v311_00_to_latest(ctx context.Context, client coreclientv1.ConfigMapsGetter, recorder events.Recorder, operatorConfig *operatorv1.OpenShiftAPIServer) (*corev1.ConfigMap, bool, error) {
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

	return resourceapply.ApplyConfigMap(ctx, client, recorder, requiredConfigMap)
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
	ctx context.Context,
	kubeClient kubernetes.Interface,
	client appsclientv1.DeploymentsGetter,
	countNodes nodeCountFunc,
	recorder events.Recorder,
	imagePullSpec string,
	operatorImagePullSpec string,
	operatorConfig *operatorv1.OpenShiftAPIServer,
	generationStatus []operatorv1.GenerationStatus,
	ensureAtMostOnePodPerNodeFn ensureAtMostOnePodPerNodeFunc,
) (*appsv1.Deployment, bool, error) {
	tmpl := v311_00_assets.MustAsset("v3.11.0/openshift-apiserver/deploy.yaml")

	r := strings.NewReplacer(
		"${IMAGE}", imagePullSpec,
		"${OPERATOR_IMAGE}", operatorImagePullSpec,
		"${REVISION}", strconv.Itoa(int(operatorConfig.Status.LatestAvailableRevision)),
		"${VERBOSITY}", loglevelToKlog(operatorConfig.Spec.LogLevel),
		"${KUBE_APISERVER_OPERATOR_IMAGE}", os.Getenv("KUBE_APISERVER_OPERATOR_IMAGE"),
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

	err = ensureAtMostOnePodPerNodeFn(&required.Spec, operatorclient.TargetNamespace)
	if err != nil {
		return nil, false, fmt.Errorf("unable to ensure at most one pod per node: %v", err)
	}

	// Set the replica count to the number of master nodes.
	masterNodeCount, err := countNodes(required.Spec.Template.Spec.NodeSelector)
	if err != nil {
		return nil, false, fmt.Errorf("failed to determine number of master nodes: %v", err)
	}
	required.Spec.Replicas = masterNodeCount

	return resourceapply.ApplyDeployment(ctx, client, recorder, required, resourcemerge.ExpectedDeploymentGeneration(required, generationStatus))
}

var openshiftScheme = runtime.NewScheme()

func init() {
	if err := openshiftapi.Install(openshiftScheme); err != nil {
		panic(err)
	}
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

func handleErrorForOperatorStatus(_ *operatorv1.OperatorStatus, _ bool, err error) {
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("failed to update the operator status, err %v", err))
	}
}
