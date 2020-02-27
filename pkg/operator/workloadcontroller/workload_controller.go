package workloadcontroller

import (
	"context"
	"fmt"
	"time"

	"github.com/openshift/library-go/pkg/operator/v1helpers"

	"github.com/openshift/library-go/pkg/operator/status"

	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"

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
	"k8s.io/klog"
	apiregistrationinformers "k8s.io/kube-aggregator/pkg/client/informers/externalversions"

	operatorsv1 "github.com/openshift/api/operator/v1"
	openshiftconfigclientv1 "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"
	configinformers "github.com/openshift/client-go/config/informers/externalversions"
	operatorv1client "github.com/openshift/client-go/operator/clientset/versioned/typed/operator/v1"
	operatorv1informers "github.com/openshift/client-go/operator/informers/externalversions/operator/v1"
	clusteroperatorv1helpers "github.com/openshift/library-go/pkg/config/clusteroperator/v1helpers"
	"github.com/openshift/library-go/pkg/operator/events"
)

const (
	imageImportCAName = "image-import-ca"
	workQueueKey      = "key"
)

type OpenShiftAPIServerOperator struct {
	targetImagePullSpec, targetOperandVersion, operatorImagePullSpec string

	operatorClient        v1helpers.OperatorClient
	versionRecorder       status.VersionGetter
	operatorConfigClient  operatorv1client.OpenShiftAPIServersGetter
	openshiftConfigClient openshiftconfigclientv1.ConfigV1Interface
	kubeClient            kubernetes.Interface
	eventRecorder         events.Recorder

	// queue only ever has one item, but it has nice error handling backoff/retry semantics
	queue workqueue.RateLimitingInterface

	// haveObservedExtensionConfigMap preserves the state so that we don't ask the server on every sync
	haveObservedExtensionConfigMap bool
}

func NewWorkloadController(
	targetImagePullSpec, targetOperandVersion, operatorImagePullSpec string,
	operatorClient v1helpers.OperatorClient,
	versionRecorder status.VersionGetter,
	operatorConfigInformer operatorv1informers.OpenShiftAPIServerInformer,
	kubeInformersForOpenShiftAPIServerNamespace kubeinformers.SharedInformerFactory,
	kubeInformersForEtcdNamespace kubeinformers.SharedInformerFactory,
	kubeInformersForOpenShiftConfigNamespace kubeinformers.SharedInformerFactory,
	apiregistrationInformers apiregistrationinformers.SharedInformerFactory,
	configInformers configinformers.SharedInformerFactory,
	operatorConfigClient operatorv1client.OpenShiftAPIServersGetter,
	openshiftConfigClient openshiftconfigclientv1.ConfigV1Interface,
	kubeClient kubernetes.Interface,
	eventRecorder events.Recorder,
) *OpenShiftAPIServerOperator {
	c := &OpenShiftAPIServerOperator{
		targetImagePullSpec:   targetImagePullSpec,
		targetOperandVersion:  targetOperandVersion,
		operatorImagePullSpec: operatorImagePullSpec,

		operatorClient:        operatorClient,
		versionRecorder:       versionRecorder,
		operatorConfigClient:  operatorConfigClient,
		openshiftConfigClient: openshiftConfigClient,
		kubeClient:            kubeClient,
		eventRecorder:         eventRecorder.WithComponentSuffix("workload-controller"),

		queue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "OpenShiftAPIServerOperator"),
	}

	operatorConfigInformer.Informer().AddEventHandler(c.eventHandler())
	kubeInformersForEtcdNamespace.Core().V1().ConfigMaps().Informer().AddEventHandler(c.eventHandler())
	kubeInformersForEtcdNamespace.Core().V1().Secrets().Informer().AddEventHandler(c.eventHandler())
	kubeInformersForOpenShiftAPIServerNamespace.Core().V1().ConfigMaps().Informer().AddEventHandler(c.eventHandler())
	kubeInformersForOpenShiftAPIServerNamespace.Core().V1().ServiceAccounts().Informer().AddEventHandler(c.eventHandler())
	kubeInformersForOpenShiftAPIServerNamespace.Core().V1().Services().Informer().AddEventHandler(c.eventHandler())
	kubeInformersForOpenShiftAPIServerNamespace.Apps().V1().Deployments().Informer().AddEventHandler(c.eventHandler())
	kubeInformersForOpenShiftConfigNamespace.Core().V1().ConfigMaps().Informer().AddEventHandler(c.eventHandler())
	configInformers.Config().V1().Images().Informer().AddEventHandler(c.eventHandler())
	apiregistrationInformers.Apiregistration().V1().APIServices().Informer().AddEventHandler(c.eventHandler())

	// we only watch some namespaces
	kubeInformersForOpenShiftAPIServerNamespace.Core().V1().Namespaces().Informer().AddEventHandler(c.namespaceEventHandler())

	return c
}

func (c OpenShiftAPIServerOperator) sync() error {
	operatorConfig, err := c.operatorConfigClient.OpenShiftAPIServers().Get("cluster", metav1.GetOptions{})
	if err != nil {
		return err
	}

	switch operatorConfig.Spec.ManagementState {
	case operatorsv1.Managed:
	case operatorsv1.Unmanaged:
		return nil
	case operatorsv1.Removed:
		// TODO probably need to watch until the NS is really gone
		if err := c.kubeClient.CoreV1().Namespaces().Delete(operatorclient.TargetNamespace, nil); err != nil && !apierrors.IsNotFound(err) {
			return err
		}
		return nil
	default:
		c.eventRecorder.Warningf("ManagementStateUnknown", "Unrecognized operator management state %q", operatorConfig.Spec.ManagementState)
		return nil
	}

	kubeAPIServerOperator, err := c.openshiftConfigClient.ClusterOperators().Get("kube-apiserver", metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		kubeAPIServerOperator, err = c.openshiftConfigClient.ClusterOperators().Get("openshift-kube-apiserver-operator", metav1.GetOptions{})
	}
	if apierrors.IsNotFound(err) {
		message := "clusteroperator/kube-apiserver not found"
		c.eventRecorder.Warning("PrereqNotReady", message)
		return fmt.Errorf(message)
	}
	if err != nil {
		return err
	}
	if !clusteroperatorv1helpers.IsStatusConditionTrue(kubeAPIServerOperator.Status.Conditions, "Available") {
		message := fmt.Sprintf("clusteroperator/%s is not Available", kubeAPIServerOperator.Name)
		c.eventRecorder.Warning("PrereqNotReady", message)
		return fmt.Errorf(message)
	}

	// block until config is obvserved
	if len(operatorConfig.Spec.ObservedConfig.Raw) == 0 {
		klog.Info("Waiting for observed configuration to be available")
		return nil
	}

	// block until extension-apiserver-authentication configmap is available
	// see https://bugzilla.redhat.com/show_bug.cgi?id=1795163#c19 to check why we have to wait for it
	// in the future we need to change upstream code to be more dynamic
	if !c.haveObservedExtensionConfigMap {
		_, err = c.kubeClient.CoreV1().ConfigMaps(metav1.NamespaceSystem).Get("extension-apiserver-authentication", metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			klog.Infof("Waiting for %q configmap in %q namespace to be available", "extension-apiserver-authentication", metav1.NamespaceSystem)
			return nil
		}
		if err != nil {
			return err
		}
		c.haveObservedExtensionConfigMap = true
	}

	return syncOpenShiftAPIServer_v311_00_to_latest(c, operatorConfig)
}

// Run starts the openshift-apiserver and blocks until stopCh is closed.
func (c *OpenShiftAPIServerOperator) Run(ctx context.Context, workers int) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	klog.Infof("Starting OpenShiftAPIServerOperator")
	defer klog.Infof("Shutting down OpenShiftAPIServerOperator")

	// doesn't matter what workers say, only start one.
	go wait.Until(c.runWorker, time.Second, ctx.Done())

	go wait.Until(func() { c.queue.Add(workQueueKey) }, time.Minute, ctx.Done())

	<-ctx.Done()
}

func (c *OpenShiftAPIServerOperator) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *OpenShiftAPIServerOperator) processNextWorkItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(key)

	err := c.sync()
	if err == nil {
		c.queue.Forget(key)
		return true
	}

	utilruntime.HandleError(fmt.Errorf("%v failed with : %v", key, err))
	c.queue.AddRateLimited(key)

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
