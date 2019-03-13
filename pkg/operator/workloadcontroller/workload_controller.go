package workloadcontroller

import (
	"fmt"
	"time"

	"github.com/golang/glog"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	apiregistrationv1client "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset/typed/apiregistration/v1"
	apiregistrationinformers "k8s.io/kube-aggregator/pkg/client/informers/externalversions"

	operatorv1 "github.com/openshift/api/operator/v1"
	openshiftconfigclientv1 "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"
	configinformers "github.com/openshift/client-go/config/informers/externalversions"
	operatorv1client "github.com/openshift/client-go/operator/clientset/versioned/typed/operator/v1"
	operatorv1informers "github.com/openshift/client-go/operator/informers/externalversions/operator/v1"
	clusteroperatorv1helpers "github.com/openshift/library-go/pkg/config/clusteroperator/v1helpers"
	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/management"
	"github.com/openshift/library-go/pkg/operator/status"
	"github.com/openshift/library-go/pkg/operator/v1helpers"
	operatorv1helpers "github.com/openshift/library-go/pkg/operator/v1helpers"

	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"
)

const (
	workloadFailingCondition = "WorkloadFailing"
	imageImportCAName        = "image-import-ca"
	workQueueKey             = "key"
)

type OpenShiftAPIServerOperator struct {
	targetImagePullSpec string
	versionRecorder     status.VersionGetter

	operatorConfigClient    operatorv1client.OpenShiftAPIServersGetter
	openshiftConfigClient   openshiftconfigclientv1.ConfigV1Interface
	operatorStatusProvider  operatorv1helpers.OperatorClient
	kubeClient              kubernetes.Interface
	apiregistrationv1Client apiregistrationv1client.ApiregistrationV1Interface
	eventRecorder           events.Recorder

	// queue only ever has one item, but it has nice error handling backoff/retry semantics
	queue workqueue.RateLimitingInterface
}

func NewWorkloadController(
	targetImagePullSpec string,
	versionRecorder status.VersionGetter,
	openshiftAPIOperatorConfigInformer operatorv1informers.OpenShiftAPIServerInformer,
	kubeAPIOperatorConfigInformer operatorv1informers.KubeAPIServerInformer,
	kubeInformersForOpenShiftAPIServerNamespace kubeinformers.SharedInformerFactory,
	kubeInformersForEtcdNamespace kubeinformers.SharedInformerFactory,
	kubeInformersForOpenShiftConfigNamespace kubeinformers.SharedInformerFactory,
	apiregistrationInformers apiregistrationinformers.SharedInformerFactory,
	configInformers configinformers.SharedInformerFactory,
	operatorConfigClient operatorv1client.OpenShiftAPIServersGetter,
	openshiftConfigClient openshiftconfigclientv1.ConfigV1Interface,
	kubeClient kubernetes.Interface,
	apiregistrationv1Client apiregistrationv1client.ApiregistrationV1Interface,
	operatorStatusProvider operatorv1helpers.OperatorClient,
	eventRecorder events.Recorder,
) *OpenShiftAPIServerOperator {
	c := &OpenShiftAPIServerOperator{
		targetImagePullSpec: targetImagePullSpec,
		versionRecorder:     versionRecorder,

		operatorConfigClient:    operatorConfigClient,
		openshiftConfigClient:   openshiftConfigClient,
		operatorStatusProvider:  operatorStatusProvider,
		kubeClient:              kubeClient,
		apiregistrationv1Client: apiregistrationv1Client,
		eventRecorder:           eventRecorder.WithComponentSuffix("workload-controller"),

		queue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "OpenShiftAPIServerOperator"),
	}

	openshiftAPIOperatorConfigInformer.Informer().AddEventHandler(c.eventHandler())
	kubeInformersForEtcdNamespace.Core().V1().ConfigMaps().Informer().AddEventHandler(c.eventHandler())
	kubeInformersForEtcdNamespace.Core().V1().Secrets().Informer().AddEventHandler(c.eventHandler())
	kubeInformersForOpenShiftAPIServerNamespace.Core().V1().ConfigMaps().Informer().AddEventHandler(c.eventHandler())
	kubeInformersForOpenShiftAPIServerNamespace.Core().V1().ServiceAccounts().Informer().AddEventHandler(c.eventHandler())
	kubeInformersForOpenShiftAPIServerNamespace.Core().V1().Services().Informer().AddEventHandler(c.eventHandler())
	kubeInformersForOpenShiftAPIServerNamespace.Apps().V1().DaemonSets().Informer().AddEventHandler(c.eventHandler())
	kubeInformersForOpenShiftConfigNamespace.Core().V1().ConfigMaps().Informer().AddEventHandler(c.eventHandler())
	configInformers.Config().V1().Images().Informer().AddEventHandler(c.eventHandler())
	apiregistrationInformers.Apiregistration().V1().APIServices().Informer().AddEventHandler(c.eventHandler())

	// requeue when the cluster operators kube-apiserver change
	kubeAPIOperatorConfigInformer.Informer().AddEventHandler(c.eventHandler())

	// we only watch some namespaces
	kubeInformersForOpenShiftAPIServerNamespace.Core().V1().Namespaces().Informer().AddEventHandler(c.namespaceEventHandler())

	return c
}

func (c OpenShiftAPIServerOperator) sync() error {
	operatorConfig, err := c.operatorConfigClient.OpenShiftAPIServers().Get("cluster", metav1.GetOptions{})
	if err != nil {
		return err
	}

	if !management.IsOperatorManaged(operatorConfig.Spec.ManagementState) {
		return nil
	}

	cond := operatorv1.OperatorCondition{
		Type:   "PrerequisiteNotReady",
		Status: operatorv1.ConditionFalse,
	}

	kubeAPIServerOperator, err := c.openshiftConfigClient.ClusterOperators().Get("kube-apiserver", metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		kubeAPIServerOperator, err = c.openshiftConfigClient.ClusterOperators().Get("openshift-kube-apiserver-operator", metav1.GetOptions{})
	}
	if apierrors.IsNotFound(err) {
		return c.reportNotReadyPreRequirements("NotFound", fmt.Sprintf("clusteroperator/kube-apiserver: %v", err))
	}
	if err != nil {
		return err
	}
	if !clusteroperatorv1helpers.IsStatusConditionTrue(kubeAPIServerOperator.Status.Conditions, "Available") {
		return c.reportNotReadyPreRequirements("NotAvailable", fmt.Sprintf("clusteroperator/%s is not Available", kubeAPIServerOperator.Name))
	}

	// block until config is obvserved
	if len(operatorConfig.Spec.ObservedConfig.Raw) == 0 {
		return c.reportNotReadyPreRequirements("NoConfig", "waiting for config")
	}

	// clean up the prereq not ready condition
	_, _, updateError := v1helpers.UpdateStatus(c.operatorStatusProvider, v1helpers.UpdateConditionFn(cond))
	if updateError != nil {
		return err
	}

	forceRequeue, err := syncOpenShiftAPIServer_v311_00_to_latest(c, operatorConfig)
	if forceRequeue && err != nil {
		c.queue.AddRateLimited(workQueueKey)
	}

	return err
}

func (c OpenShiftAPIServerOperator) reportNotReadyPreRequirements(reason, message string) error {
	cond := operatorv1.OperatorCondition{
		Type:    "PrerequisiteNotReady",
		Status:  operatorv1.ConditionTrue,
		Reason:  reason,
		Message: message,
	}
	_, _, updateError := v1helpers.UpdateStatus(c.operatorStatusProvider, v1helpers.UpdateConditionFn(cond))
	if updateError == nil {
		c.eventRecorder.Warning("PrerequisiteNotReady", cond.Message)
	} else {
		return updateError
	}
	return fmt.Errorf(cond.Message)
}

// Run starts the openshift-apiserver and blocks until stopCh is closed.
func (c *OpenShiftAPIServerOperator) Run(workers int, stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	glog.Infof("Starting OpenShiftAPIServerOperator")
	defer glog.Infof("Shutting down OpenShiftAPIServerOperator")

	// doesn't matter what workers say, only start one.
	go wait.Until(c.runWorker, time.Second, stopCh)

	<-stopCh
}

func (c *OpenShiftAPIServerOperator) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *OpenShiftAPIServerOperator) processNextWorkItem() bool {
	dsKey, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(dsKey)

	err := c.sync()
	if err == nil {
		c.queue.Forget(dsKey)
		return true
	}

	utilruntime.HandleError(fmt.Errorf("%v failed with : %v", dsKey, err))
	c.queue.AddRateLimited(dsKey)

	return true
}

// eventHandler queues the operator to check spec and status
func (c *OpenShiftAPIServerOperator) eventHandler() cache.ResourceEventHandler {
	return cache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj interface{}) { c.queue.Add(workQueueKey) },
		UpdateFunc: func(old, new interface{}) { c.queue.Add(workQueueKey) },
		DeleteFunc: func(obj interface{}) { c.queue.Add(workQueueKey) },
	}
}

// this set of namespaces will include things like logging and metrics which are used to drive
var interestingNamespaces = sets.NewString(operatorclient.TargetNamespace)

func (c *OpenShiftAPIServerOperator) namespaceEventHandler() cache.ResourceEventHandler {
	return cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			ns, ok := obj.(*corev1.Namespace)
			if !ok {
				c.queue.Add(workQueueKey)
			}
			if ns.Name == operatorclient.TargetNamespace {
				c.queue.Add(workQueueKey)
			}
		},
		UpdateFunc: func(old, new interface{}) {
			ns, ok := old.(*corev1.Namespace)
			if !ok {
				c.queue.Add(workQueueKey)
			}
			if ns.Name == operatorclient.TargetNamespace {
				c.queue.Add(workQueueKey)
			}
		},
		DeleteFunc: func(obj interface{}) {
			ns, ok := obj.(*corev1.Namespace)
			if !ok {
				tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
				if !ok {
					utilruntime.HandleError(fmt.Errorf("couldn't get object from tombstone %#v", obj))
					return
				}
				ns, ok = tombstone.Obj.(*corev1.Namespace)
				if !ok {
					utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not a Namespace %#v", obj))
					return
				}
			}
			if ns.Name == operatorclient.TargetNamespace {
				c.queue.Add(workQueueKey)
			}
		},
	}
}
