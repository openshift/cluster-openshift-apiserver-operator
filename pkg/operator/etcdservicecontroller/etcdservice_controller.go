package etcdservicecontroller

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"

	operatorv1 "github.com/openshift/api/operator/v1"
	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/v1helpers"
)

var (
	controllerWorkQueueKey = "key"

	degradedConditionNames = []string{
		"EtcdServiceDegraded",
		"EtcdEndpointDegraded",
	}
)

type EtcdServiceController struct {
	targetNamespace string

	operatorClient  v1helpers.OperatorClient
	endpointsGetter corev1client.EndpointsGetter
	servicesGetter  corev1client.ServicesGetter

	cachesToSync  []cache.InformerSynced
	queue         workqueue.RateLimitingInterface
	eventRecorder events.Recorder
}

func NewEtcdServiceController(
	targetNamespace string,
	kubeInformersForTargetNamespace informers.SharedInformerFactory,
	operatorClient v1helpers.OperatorClient,
	servicesGetter corev1client.ServicesGetter,
	endpointsGetter corev1client.EndpointsGetter,
	eventRecorder events.Recorder,
) *EtcdServiceController {
	c := &EtcdServiceController{
		targetNamespace: targetNamespace,
		operatorClient:  operatorClient,
		servicesGetter:  servicesGetter,
		endpointsGetter: endpointsGetter,
		eventRecorder:   eventRecorder.WithComponentSuffix("etcd-service-controller"),

		queue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "EtcdServiceController"),
	}

	operatorClient.Informer().AddEventHandler(c.eventHandler())

	kubeInformersForTargetNamespace.Core().V1().Services().Informer().AddEventHandler(c.eventHandler())
	kubeInformersForTargetNamespace.Core().V1().Endpoints().Informer().AddEventHandler(c.eventHandler())

	c.cachesToSync = append(c.cachesToSync, operatorClient.Informer().HasSynced)
	c.cachesToSync = append(c.cachesToSync, kubeInformersForTargetNamespace.Core().V1().Services().Informer().HasSynced)
	c.cachesToSync = append(c.cachesToSync, kubeInformersForTargetNamespace.Core().V1().Endpoints().Informer().HasSynced)

	return c
}

func (c *EtcdServiceController) sync() error {
	foundConditions := []operatorv1.OperatorCondition{}

	if _, err := c.servicesGetter.Services(c.targetNamespace).Get("etcd", metav1.GetOptions{}); err != nil {
		foundConditions = append(foundConditions, operatorv1.OperatorCondition{
			Type:    "EtcdServiceDegraded",
			Status:  operatorv1.ConditionTrue,
			Reason:  "EtcdServiceError",
			Message: fmt.Sprintf("Error getting etcd service: %v", err),
		})
	}

	etcdEndpoints, err := c.endpointsGetter.Endpoints(c.targetNamespace).Get("etcd", metav1.GetOptions{})
	if err != nil {
		foundConditions = append(foundConditions, operatorv1.OperatorCondition{
			Type:    "EtcdEndpointDegraded",
			Status:  operatorv1.ConditionTrue,
			Reason:  "EtcdEndpointError",
			Message: fmt.Sprintf("Error getting etcd service: %v", err),
		})
	} else {
		if len(etcdEndpoints.Subsets) == 0 {
			foundConditions = append(foundConditions, operatorv1.OperatorCondition{
				Type:    "EtcdEndpointDegraded",
				Status:  operatorv1.ConditionTrue,
				Reason:  "EtcdEndpointNoSubnets",
				Message: "Etcd endpoint has empty subnets",
			})
		} else {
			if len(etcdEndpoints.Subsets[0].Addresses) == 0 {
				foundConditions = append(foundConditions, operatorv1.OperatorCondition{
					Type:    "EtcdEndpointDegraded",
					Status:  operatorv1.ConditionTrue,
					Reason:  "EtcdEndpointNoAddresses",
					Message: "Etcd endpoint has empty addresses",
				})
			}
		}
	}

	updateConditionFuncs := []v1helpers.UpdateStatusFunc{}

	// check the supported degraded foundConditions and check if any pending pod matching them.
	for _, degradedConditionName := range degradedConditionNames {
		// clean up existing foundConditions
		updatedCondition := operatorv1.OperatorCondition{
			Type:   degradedConditionName,
			Status: operatorv1.ConditionFalse,
		}
		if condition := v1helpers.FindOperatorCondition(foundConditions, degradedConditionName); condition != nil {
			updatedCondition = *condition
		}
		updateConditionFuncs = append(updateConditionFuncs, v1helpers.UpdateConditionFn(updatedCondition))
	}

	if _, _, err := v1helpers.UpdateStatus(c.operatorClient, updateConditionFuncs...); err != nil {
		return err
	}

	return err
}

// Run starts the kube-apiserver and blocks until stopCh is closed.
func (c *EtcdServiceController) Run(ctx context.Context, workers int) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	klog.Infof("Starting EtcdServiceController")
	defer klog.Infof("Shutting down EtcdServiceController")
	if !cache.WaitForCacheSync(ctx.Done(), c.cachesToSync...) {
		return
	}

	// doesn't matter what workers say, only start one.
	go wait.UntilWithContext(ctx, c.runWorker, time.Second)

	// add time based trigger
	go wait.UntilWithContext(ctx, func(context.Context) { c.queue.Add(controllerWorkQueueKey) }, time.Minute)

	<-ctx.Done()
}

func (c *EtcdServiceController) runWorker(ctx context.Context) {
	for c.processNextWorkItem() {
	}
}

func (c *EtcdServiceController) processNextWorkItem() bool {
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
func (c *EtcdServiceController) eventHandler() cache.ResourceEventHandler {
	return cache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj interface{}) { c.queue.Add(controllerWorkQueueKey) },
		UpdateFunc: func(old, new interface{}) { c.queue.Add(controllerWorkQueueKey) },
		DeleteFunc: func(obj interface{}) { c.queue.Add(controllerWorkQueueKey) },
	}
}
