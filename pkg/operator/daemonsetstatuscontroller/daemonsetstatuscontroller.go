package daemonsetstatuscontroller

import (
	"fmt"
	"time"

	operatorv1 "github.com/openshift/api/operator/v1"
	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/v1helpers"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	appsv1listers "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
)

var (
	controllerWorkQueueKey = "key"
)

type DaemonSetStatusController struct {
	targetNamespace string
	targetName      string

	operatorClient  v1helpers.OperatorClient
	daemonSetLister appsv1listers.DaemonSetLister

	cachesToSync  []cache.InformerSynced
	queue         workqueue.RateLimitingInterface
	eventRecorder events.Recorder
}

// goes available=false when missing/no pods
// goes progressing=true when rolling out
// goes degraded=true when pod not available

func NewDaemonSetStatusController(
	targetNamespace string,
	targetName string,
	kubeInformersForTargetNamespace informers.SharedInformerFactory,
	operatorClient v1helpers.OperatorClient,
	eventRecorder events.Recorder,
) *DaemonSetStatusController {
	c := &DaemonSetStatusController{
		targetNamespace: targetNamespace,
		targetName:      targetName,

		operatorClient:  operatorClient,
		daemonSetLister: kubeInformersForTargetNamespace.Apps().V1().DaemonSets().Lister(),
		eventRecorder:   eventRecorder.WithComponentSuffix("daemonset-status-controller"),

		queue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "DaemonSetStatusController"),
	}

	operatorClient.Informer().AddEventHandler(c.eventHandler())
	kubeInformersForTargetNamespace.Apps().V1().DaemonSets().Informer().AddEventHandler(c.eventHandler())

	c.cachesToSync = append(c.cachesToSync, operatorClient.Informer().HasSynced)
	c.cachesToSync = append(c.cachesToSync, kubeInformersForTargetNamespace.Apps().V1().DaemonSets().Informer().HasSynced)

	return c
}

func (c *DaemonSetStatusController) sync() error {
	actualDaemonSet, err := c.daemonSetLister.DaemonSets(c.targetNamespace).Get(c.targetName)
	if err != nil {
		return err
	}

	errors := []error{}
	if actualDaemonSet == nil {
		message := fmt.Sprintf("daemonset.v1.apps/%v: could not be retrieved", c.targetName)
		errors = append(errors,
			appendErrors(v1helpers.UpdateStatus(c.operatorClient, v1helpers.UpdateConditionFn(operatorv1.OperatorCondition{
				Type:    "APIServerDaemonSetAvailable",
				Status:  operatorv1.ConditionFalse,
				Reason:  "NoDaemon",
				Message: message,
			})))...,
		)
		errors = append(errors,
			appendErrors(v1helpers.UpdateStatus(c.operatorClient, v1helpers.UpdateConditionFn(operatorv1.OperatorCondition{
				Type:    "APIServerDaemonSetProgressing",
				Status:  operatorv1.ConditionTrue,
				Reason:  "NoDaemon",
				Message: message,
			})))...,
		)
		errors = append(errors,
			appendErrors(v1helpers.UpdateStatus(c.operatorClient, v1helpers.UpdateConditionFn(operatorv1.OperatorCondition{
				Type:    "APIServerDaemonSetDegraded",
				Status:  operatorv1.ConditionTrue,
				Reason:  "NoDaemon",
				Message: message,
			})))...,
		)

		return utilerrors.NewAggregate(errors)
	}

	if actualDaemonSet.Status.NumberAvailable == 0 {
		errors = append(errors,
			appendErrors(v1helpers.UpdateStatus(c.operatorClient, v1helpers.UpdateConditionFn(operatorv1.OperatorCondition{
				Type:    "APIServerDaemonSetAvailable",
				Status:  operatorv1.ConditionFalse,
				Reason:  "NoAPIServerPod",
				Message: "no openshift-apiserver daemon pods available on any node.",
			})))...,
		)
	} else {
		errors = append(errors,
			appendErrors(v1helpers.UpdateStatus(c.operatorClient, v1helpers.UpdateConditionFn(operatorv1.OperatorCondition{
				Type:   "APIServerDaemonSetAvailable",
				Status: operatorv1.ConditionTrue,
				Reason: "AsExpected",
			})))...,
		)
	}

	// If the daemonset is up to date and the operatorConfig are up to date, then we are no longer progressing
	daemonSetAtHighestGeneration := actualDaemonSet.ObjectMeta.Generation == actualDaemonSet.Status.ObservedGeneration
	if !daemonSetAtHighestGeneration {
		errors = append(errors,
			appendErrors(v1helpers.UpdateStatus(c.operatorClient, v1helpers.UpdateConditionFn(operatorv1.OperatorCondition{
				Type:    "APIServerDaemonSetProgressing",
				Status:  operatorv1.ConditionTrue,
				Reason:  "NewGeneration",
				Message: fmt.Sprintf("daemonset/apiserver.openshift-operator: observed generation is %d, desired generation is %d.", actualDaemonSet.Status.ObservedGeneration, actualDaemonSet.ObjectMeta.Generation),
			})))...,
		)
	} else {
		errors = append(errors,
			appendErrors(v1helpers.UpdateStatus(c.operatorClient, v1helpers.UpdateConditionFn(operatorv1.OperatorCondition{
				Type:   "APIServerDaemonSetProgressing",
				Status: operatorv1.ConditionFalse,
				Reason: "AsExpected",
			})))...,
		)
	}

	daemonSetHasAllPodsAvailable := actualDaemonSet.Status.NumberAvailable == actualDaemonSet.Status.DesiredNumberScheduled
	if !daemonSetHasAllPodsAvailable {
		numNonAvailablePods := actualDaemonSet.Status.DesiredNumberScheduled - actualDaemonSet.Status.NumberAvailable
		errors = append(errors,
			appendErrors(v1helpers.UpdateStatus(c.operatorClient, v1helpers.UpdateConditionFn(operatorv1.OperatorCondition{
				Type:    "APIServerDaemonSetDegraded",
				Status:  operatorv1.ConditionTrue,
				Reason:  "UnavailablePod",
				Message: fmt.Sprintf("%v of %v requested instances are unavailable", numNonAvailablePods, actualDaemonSet.Status.DesiredNumberScheduled),
			})))...,
		)
	} else {
		errors = append(errors,
			appendErrors(v1helpers.UpdateStatus(c.operatorClient, v1helpers.UpdateConditionFn(operatorv1.OperatorCondition{
				Type:   "APIServerDaemonSetDegraded",
				Status: operatorv1.ConditionFalse,
				Reason: "AsExpected",
			})))...,
		)
	}

	return utilerrors.NewAggregate(errors)
}

func appendErrors(_ *operatorv1.OperatorStatus, _ bool, err error) []error {
	if err != nil {
		return []error{err}
	}
	return []error{}
}

// Run starts the kube-apiserver and blocks until stopCh is closed.
func (c *DaemonSetStatusController) Run(workers int, stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	klog.Infof("Starting DaemonSetStatusController")
	defer klog.Infof("Shutting down DaemonSetStatusController")
	if !cache.WaitForCacheSync(stopCh, c.cachesToSync...) {
		return
	}

	// doesn't matter what workers say, only start one.
	go wait.Until(c.runWorker, time.Second, stopCh)

	// add time based trigger
	go wait.Until(func() { c.queue.Add(controllerWorkQueueKey) }, time.Minute, stopCh)

	<-stopCh
}

func (c *DaemonSetStatusController) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *DaemonSetStatusController) processNextWorkItem() bool {
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
func (c *DaemonSetStatusController) eventHandler() cache.ResourceEventHandler {
	return cache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj interface{}) { c.queue.Add(controllerWorkQueueKey) },
		UpdateFunc: func(old, new interface{}) { c.queue.Add(controllerWorkQueueKey) },
		DeleteFunc: func(obj interface{}) { c.queue.Add(controllerWorkQueueKey) },
	}
}
