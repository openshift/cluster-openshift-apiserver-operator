package operator

import (
	"fmt"
	"time"

	"github.com/blang/semver"
	"github.com/golang/glog"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/flowcontrol"
	"k8s.io/client-go/util/workqueue"
	apiregistrationv1client "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset/typed/apiregistration/v1"
	apiregistrationinformers "k8s.io/kube-aggregator/pkg/client/informers/externalversions"

	operatorsv1 "github.com/openshift/api/operator/v1"
	operatorconfigclientv1alpha1 "github.com/openshift/cluster-openshift-apiserver-operator/pkg/generated/clientset/versioned/typed/openshiftapiserver/v1alpha1"
	operatorconfiginformerv1alpha1 "github.com/openshift/cluster-openshift-apiserver-operator/pkg/generated/informers/externalversions/openshiftapiserver/v1alpha1"
	"github.com/openshift/library-go/pkg/operator/v1alpha1helpers"
	"github.com/openshift/library-go/pkg/operator/versioning"
)

const (
	etcdNamespaceName          = "kube-system"
	kubeAPIServerNamespaceName = "openshift-kube-apiserver"
	targetNamespaceName        = "openshift-apiserver"
	workQueueKey               = "key"
)

type OpenShiftAPIServerOperator struct {
	operatorConfigClient operatorconfigclientv1alpha1.OpenshiftapiserverV1alpha1Interface

	kubeClient              kubernetes.Interface
	apiregistrationv1Client apiregistrationv1client.ApiregistrationV1Interface

	// queue only ever has one item, but it has nice error handling backoff/retry semantics
	queue workqueue.RateLimitingInterface

	rateLimiter flowcontrol.RateLimiter
}

func NewKubeApiserverOperator(
	operatorConfigInformer operatorconfiginformerv1alpha1.OpenShiftAPIServerOperatorConfigInformer,
	kubeInformersForOpenShiftAPIServerNamespace kubeinformers.SharedInformerFactory,
	kubeInformersForEtcdNamespace kubeinformers.SharedInformerFactory,
	kubeInformersForKubeAPIServerNamespace kubeinformers.SharedInformerFactory,
	apiregistrationInformers apiregistrationinformers.SharedInformerFactory,
	operatorConfigClient operatorconfigclientv1alpha1.OpenshiftapiserverV1alpha1Interface,
	kubeClient kubernetes.Interface,
	apiregistrationv1Client apiregistrationv1client.ApiregistrationV1Interface,
) *OpenShiftAPIServerOperator {
	c := &OpenShiftAPIServerOperator{
		operatorConfigClient:    operatorConfigClient,
		kubeClient:              kubeClient,
		apiregistrationv1Client: apiregistrationv1Client,

		queue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "OpenShiftAPIServerOperator"),

		rateLimiter: flowcontrol.NewTokenBucketRateLimiter(0.05 /*3 per minute*/, 4),
	}

	operatorConfigInformer.Informer().AddEventHandler(c.eventHandler())
	kubeInformersForEtcdNamespace.Core().V1().ConfigMaps().Informer().AddEventHandler(c.eventHandler())
	kubeInformersForEtcdNamespace.Core().V1().Secrets().Informer().AddEventHandler(c.eventHandler())
	kubeInformersForKubeAPIServerNamespace.Core().V1().ConfigMaps().Informer().AddEventHandler(c.eventHandler())
	kubeInformersForOpenShiftAPIServerNamespace.Core().V1().ConfigMaps().Informer().AddEventHandler(c.eventHandler())
	kubeInformersForOpenShiftAPIServerNamespace.Core().V1().ServiceAccounts().Informer().AddEventHandler(c.eventHandler())
	kubeInformersForOpenShiftAPIServerNamespace.Core().V1().Services().Informer().AddEventHandler(c.eventHandler())
	kubeInformersForOpenShiftAPIServerNamespace.Apps().V1().Deployments().Informer().AddEventHandler(c.eventHandler())
	apiregistrationInformers.Apiregistration().V1().APIServices().Informer().AddEventHandler(c.eventHandler())

	// we only watch some namespaces
	kubeInformersForOpenShiftAPIServerNamespace.Core().V1().Namespaces().Informer().AddEventHandler(c.namespaceEventHandler())

	return c
}

func (c OpenShiftAPIServerOperator) sync() error {
	operatorConfig, err := c.operatorConfigClient.OpenShiftAPIServerOperatorConfigs().Get("instance", metav1.GetOptions{})
	if err != nil {
		return err
	}
	switch operatorConfig.Spec.ManagementState {
	case operatorsv1.Unmanaged:
		return nil

	case operatorsv1.Removed:
		// TODO probably need to watch until the NS is really gone
		if err := c.kubeClient.CoreV1().Namespaces().Delete(targetNamespaceName, nil); err != nil && !apierrors.IsNotFound(err) {
			return err
		}
		operatorConfig.Status.TaskSummary = "Remove"
		operatorConfig.Status.TargetAvailability = nil
		operatorConfig.Status.CurrentAvailability = nil
		operatorConfig.Status.Conditions = []operatorsv1.OperatorCondition{
			{
				Type:   operatorsv1.OperatorStatusTypeAvailable,
				Status: operatorsv1.ConditionFalse,
			},
		}
		if _, err := c.operatorConfigClient.OpenShiftAPIServerOperatorConfigs().Update(operatorConfig); err != nil {
			return err
		}
		return nil
	}

	var currentActualVerion *semver.Version

	if operatorConfig.Status.CurrentAvailability != nil {
		ver, err := semver.Parse(operatorConfig.Status.CurrentAvailability.Version)
		if err != nil {
			utilruntime.HandleError(err)
		} else {
			currentActualVerion = &ver
		}
	}
	desiredVersion, err := semver.Parse(operatorConfig.Spec.Version)
	if err != nil {
		// TODO report failing status, we may actually attempt to do this in the "normal" error handling
		return err
	}

	v311_00_to_unknown := versioning.NewRangeOrDie("3.11.0", "3.12.0")

	var versionAvailability operatorsv1.VersionAvailability
	errors := []error{}
	switch {
	case v311_00_to_unknown.BetweenOrEmpty(currentActualVerion) && v311_00_to_unknown.Between(&desiredVersion):
		operatorConfig.Status.TaskSummary = "sync-[3.11.0,3.12.0)"
		operatorConfig.Status.TargetAvailability = nil
		versionAvailability, errors = syncOpenShiftAPIServer_v311_00_to_latest(c, operatorConfig, operatorConfig.Status.CurrentAvailability)

	default:
		operatorConfig.Status.TaskSummary = "unrecognized"
		if _, err := c.operatorConfigClient.OpenShiftAPIServerOperatorConfigs().UpdateStatus(operatorConfig); err != nil {
			utilruntime.HandleError(err)
		}

		return fmt.Errorf("unrecognized state")
	}

	v1alpha1helpers.SetStatusFromAvailability(&operatorConfig.Status.OperatorStatus, operatorConfig.ObjectMeta.Generation, &versionAvailability)
	if _, err := c.operatorConfigClient.OpenShiftAPIServerOperatorConfigs().UpdateStatus(operatorConfig); err != nil {
		errors = append(errors, err)
	}

	return utilerrors.NewAggregate(errors)
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

	// before we call sync, we want to wait for token.  We do this to avoid hot looping.
	c.rateLimiter.Accept()

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
var interestingNamespaces = sets.NewString(targetNamespaceName)

func (c *OpenShiftAPIServerOperator) namespaceEventHandler() cache.ResourceEventHandler {
	return cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			ns, ok := obj.(*corev1.Namespace)
			if !ok {
				c.queue.Add(workQueueKey)
			}
			if ns.Name == targetNamespaceName {
				c.queue.Add(workQueueKey)
			}
		},
		UpdateFunc: func(old, new interface{}) {
			ns, ok := old.(*corev1.Namespace)
			if !ok {
				c.queue.Add(workQueueKey)
			}
			if ns.Name == targetNamespaceName {
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
			if ns.Name == targetNamespaceName {
				c.queue.Add(workQueueKey)
			}
		},
	}
}
