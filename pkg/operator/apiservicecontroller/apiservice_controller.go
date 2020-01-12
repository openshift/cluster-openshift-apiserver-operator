package apiservicecontroller

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	operatorsv1 "github.com/openshift/api/operator/v1"
	operatorv1 "github.com/openshift/api/operator/v1"
	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/resource/resourceapply"
	"github.com/openshift/library-go/pkg/operator/status"
	"github.com/openshift/library-go/pkg/operator/v1helpers"
	"k8s.io/apimachinery/pkg/util/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	apiregistrationv1client "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset/typed/apiregistration/v1"
	apiregistrationinformers "k8s.io/kube-aggregator/pkg/client/informers/externalversions"
)

const (
	workQueueKey = "key"
)

type APIServiceController struct {
	name        string
	apiServices []*apiregistrationv1.APIService
	// precondition must return true before the apiservices will be created
	precondition wait.ConditionFunc

	versionRecorder         status.VersionGetter
	operatorClient          v1helpers.OperatorClient
	kubeClient              kubernetes.Interface
	apiregistrationv1Client apiregistrationv1client.ApiregistrationV1Interface
	eventRecorder           events.Recorder

	// queue only ever has one item, but it has nice error handling backoff/retry semantics
	queue workqueue.RateLimitingInterface
}

func NewAPIServiceController(
	name string,
	apiServices []*apiregistrationv1.APIService,
	operatorClient v1helpers.OperatorClient,
	apiregistrationInformers apiregistrationinformers.SharedInformerFactory,
	apiregistrationv1Client apiregistrationv1client.ApiregistrationV1Interface,
	kubeInformersForOperandNamespace kubeinformers.SharedInformerFactory,
	kubeClient kubernetes.Interface,
	eventRecorder events.Recorder,
) *APIServiceController {
	fullname := "APIServiceController_" + name
	c := &APIServiceController{
		name:         fullname,
		apiServices:  apiServices,
		precondition: NewEndpointPrecondition(kubeInformersForOperandNamespace, apiServices),

		operatorClient:          operatorClient,
		apiregistrationv1Client: apiregistrationv1Client,
		kubeClient:              kubeClient,
		eventRecorder:           eventRecorder.WithComponentSuffix("apiservice-" + name + "-controller"),

		queue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), fullname),
	}

	kubeInformersForOperandNamespace.Core().V1().Services().Informer().AddEventHandler(c.eventHandler())
	kubeInformersForOperandNamespace.Core().V1().Endpoints().Informer().AddEventHandler(c.eventHandler())
	apiregistrationInformers.Apiregistration().V1().APIServices().Informer().AddEventHandler(c.eventHandler())

	return c
}

func (c *APIServiceController) sync() error {
	operatorConfigSpec, _, _, err := c.operatorClient.GetOperatorState()
	if err != nil {
		return err
	}

	switch operatorConfigSpec.ManagementState {
	case operatorsv1.Managed:
	case operatorsv1.Unmanaged:
		return nil
	case operatorsv1.Removed:
		errs := []error{}
		for _, apiService := range c.apiServices {
			if err := c.apiregistrationv1Client.APIServices().Delete(apiService.Name, nil); err != nil {
				errs = append(errs, err)
			}
		}
		return errors.NewAggregate(errs)
	default:
		c.eventRecorder.Warningf("ManagementStateUnknown", "Unrecognized operator management state %q", operatorConfigSpec.ManagementState)
		return nil
	}

	ready, err := c.precondition()
	if err != nil {
		v1helpers.UpdateStatus(c.operatorClient, v1helpers.UpdateConditionFn(operatorv1.OperatorCondition{
			Type:    "APIServicesAvailable",
			Status:  operatorv1.ConditionFalse,
			Reason:  "ErrorCheckingPrecondition",
			Message: err.Error(),
		}))
		return err
	}
	if !ready {
		v1helpers.UpdateStatus(c.operatorClient, v1helpers.UpdateConditionFn(operatorv1.OperatorCondition{
			Type:    "APIServicesAvailable",
			Status:  operatorv1.ConditionFalse,
			Reason:  "PreconditionNotReady",
			Message: "PreconditionNotReady",
		}))
		return err
	}

	err = c.syncAPIServices()

	// update failing condition
	cond := operatorv1.OperatorCondition{
		Type:   "APIServicesAvailable",
		Status: operatorv1.ConditionTrue,
	}
	if err != nil {
		cond.Status = operatorv1.ConditionFalse
		cond.Reason = "Error"
		cond.Message = err.Error()
	}
	if _, _, updateError := v1helpers.UpdateStatus(c.operatorClient, v1helpers.UpdateConditionFn(cond)); updateError != nil {
		if err == nil {
			return updateError
		}
	}

	return err
}

func (c *APIServiceController) syncAPIServices() error {
	errs := []error{}
	var availableConditionMessages []string

	for _, apiService := range c.apiServices {
		apiregistrationv1.SetDefaults_ServiceReference(apiService.Spec.Service)
		apiService, _, err := resourceapply.ApplyAPIService(c.apiregistrationv1Client, c.eventRecorder, apiService)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		for _, condition := range apiService.Status.Conditions {
			if condition.Type == apiregistrationv1.Available {
				if condition.Status != apiregistrationv1.ConditionTrue {
					availableConditionMessages = append(availableConditionMessages, fmt.Sprintf("apiservices.apiregistration.k8s.io/%v: not available: %v", apiService.Name, condition.Message))
				}
				break
			}
		}
	}
	if len(errs) > 0 {
		return errors.NewAggregate(errs)
	}
	if len(availableConditionMessages) > 0 {
		sort.Sort(sort.StringSlice(availableConditionMessages))
		return fmt.Errorf(strings.Join(availableConditionMessages, "\n"))
	}

	// if the apiservices themselves check out ok, try to actually hit the discovery endpoints.  We have a history in clusterup
	// of something delaying them.  This isn't perfect because of round-robining, but let's see if we get an improvement
	if c.kubeClient.Discovery().RESTClient() != nil {
		missingAPIMessages := checkDiscoveryForByAPIServices(c.eventRecorder, c.kubeClient.Discovery().RESTClient(), c.apiServices)
		availableConditionMessages = append(availableConditionMessages, missingAPIMessages...)
	}

	if len(availableConditionMessages) > 0 {
		sort.Sort(sort.StringSlice(availableConditionMessages))
		return fmt.Errorf(strings.Join(availableConditionMessages, "\n"))
	}

	return nil
}

// Run starts the openshift-apiserver and blocks until stopCh is closed.
func (c *APIServiceController) Run(ctx context.Context) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	klog.Infof("Starting %v", c.name)
	defer klog.Infof("Shutting down %v", c.name)

	// doesn't matter what workers say, only start one.
	go wait.Until(c.runWorker, time.Second, ctx.Done())

	<-ctx.Done()
}

func (c *APIServiceController) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *APIServiceController) processNextWorkItem() bool {
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
func (c *APIServiceController) eventHandler() cache.ResourceEventHandler {
	return cache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj interface{}) { c.queue.Add(workQueueKey) },
		UpdateFunc: func(old, new interface{}) { c.queue.Add(workQueueKey) },
		DeleteFunc: func(obj interface{}) { c.queue.Add(workQueueKey) },
	}
}
