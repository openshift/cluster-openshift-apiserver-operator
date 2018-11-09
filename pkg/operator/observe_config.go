package operator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/imdario/mergo"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/diff"
	"k8s.io/apimachinery/pkg/util/rand"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	kubeinformers "k8s.io/client-go/informers"
	corelistersv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/flowcontrol"
	"k8s.io/client-go/util/workqueue"

	"github.com/openshift/api/operator/v1alpha1"
	imageconfiginformers "github.com/openshift/client-go/config/informers/externalversions"
	configlistersv1 "github.com/openshift/client-go/config/listers/config/v1"
	operatorconfigclientv1alpha1 "github.com/openshift/cluster-openshift-apiserver-operator/pkg/generated/clientset/versioned/typed/openshiftapiserver/v1alpha1"
	openshiftapiserveroperatorinformers "github.com/openshift/cluster-openshift-apiserver-operator/pkg/generated/informers/externalversions"
	"github.com/openshift/library-go/pkg/operator/v1alpha1helpers"
)

const configObservationErrorConditionReason = "ConfigObservationError"

type Listers struct {
	imageConfigLister configlistersv1.ImageLister
	endpointsLister   corelistersv1.EndpointsLister
	imageConfigSynced cache.InformerSynced
}

// observeConfigFunc observes configuration and returns the observedConfig. This function should not return an
// observedConfig that would cause the service being managed by the operator to crash. For example, if a required
// configuration key cannot be observed, consider reusing the configuration key's previous value. Errors that occur
// while attempting to generate the observedConfig should be returned in the errs slice.
type observeConfigFunc func(listers Listers, existingConfig map[string]interface{}) (observedConfig map[string]interface{}, errs []error)

type ConfigObserver struct {
	operatorConfigClient operatorconfigclientv1alpha1.OpenshiftapiserverV1alpha1Interface

	// queue only ever has one item, but it has nice error handling backoff/retry semantics
	queue workqueue.RateLimitingInterface

	listers Listers

	rateLimiter flowcontrol.RateLimiter

	// observers are called in an undefined order and their results are merged to
	// determine the observed configuration.
	observers []observeConfigFunc

	cachesSynced []cache.InformerSynced
}

func NewConfigObserver(
	operatorConfigInformer openshiftapiserveroperatorinformers.SharedInformerFactory,
	kubeInformersForEtcdNamespace kubeinformers.SharedInformerFactory,
	imageConfigInformer imageconfiginformers.SharedInformerFactory,
	operatorConfigClient operatorconfigclientv1alpha1.OpenshiftapiserverV1alpha1Interface,
) *ConfigObserver {
	c := &ConfigObserver{
		operatorConfigClient: operatorConfigClient,

		queue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ConfigObserver"),

		rateLimiter: flowcontrol.NewTokenBucketRateLimiter(0.05 /*3 per minute*/, 4),
		observers: []observeConfigFunc{
			observeStorageURLs,
			observeInternalRegistryHostname,
		},
		listers: Listers{
			imageConfigLister: imageConfigInformer.Config().V1().Images().Lister(),
			endpointsLister:   kubeInformersForEtcdNamespace.Core().V1().Endpoints().Lister(),
			imageConfigSynced: imageConfigInformer.Config().V1().Images().Informer().HasSynced,
		},
		cachesSynced: []cache.InformerSynced{
			operatorConfigInformer.Openshiftapiserver().V1alpha1().OpenShiftAPIServerOperatorConfigs().Informer().HasSynced,
			kubeInformersForEtcdNamespace.Core().V1().Endpoints().Informer().HasSynced,
		},
	}

	operatorConfigInformer.Openshiftapiserver().V1alpha1().OpenShiftAPIServerOperatorConfigs().Informer().AddEventHandler(c.eventHandler())
	kubeInformersForEtcdNamespace.Core().V1().Endpoints().Informer().AddEventHandler(c.eventHandler())
	imageConfigInformer.Config().V1().Images().Informer().AddEventHandler(c.eventHandler())

	return c
}

// observeStorageURLs observes the storage config URLs. If there is a problem observing the current storage config URLs,
// then the previously observed storage config URLs will be re-used.
func observeStorageURLs(listers Listers, existingConfig map[string]interface{}) (map[string]interface{}, []error) {
	storageConfigURLsPath := []string{"storageConfig", "urls"}
	previouslyObservedConfig := map[string]interface{}{}
	if currentStorageURLs, _, _ := unstructured.NestedStringSlice(existingConfig, storageConfigURLsPath...); len(currentStorageURLs) > 0 {
		unstructured.SetNestedStringSlice(previouslyObservedConfig, currentStorageURLs, storageConfigURLsPath...)
	}

	var errs []error

	var storageURLs []string
	etcdEndpoints, err := listers.endpointsLister.Endpoints(etcdNamespaceName).Get("etcd")
	if errors.IsNotFound(err) {
		errs = append(errs, fmt.Errorf("endpoints/etcd.kube-system: not found"))
		return previouslyObservedConfig, errs
	}
	if err != nil {
		errs = append(errs, err)
		return previouslyObservedConfig, errs
	}
	dnsSuffix := etcdEndpoints.Annotations["alpha.installer.openshift.io/dns-suffix"]
	if len(dnsSuffix) == 0 {
		errs = append(errs, fmt.Errorf("endpoints/etcd.kube-system: alpha.installer.openshift.io/dns-suffix annotation not found"))
		return previouslyObservedConfig, errs
	}
	for subsetIndex, subset := range etcdEndpoints.Subsets {
		for addressIndex, address := range subset.Addresses {
			if address.Hostname == "" {
				errs = append(errs, fmt.Errorf("endpoints/etcd.kube-system: subsets[%v]addresses[%v].hostname not found", subsetIndex, addressIndex))
				continue
			}
			storageURLs = append(storageURLs, "https://"+address.Hostname+"."+dnsSuffix+":2379")
		}
	}

	if len(storageURLs) == 0 {
		errs = append(errs, fmt.Errorf("endpoints/etcd.kube-system: no etcd endpoint addresses found"))
	}
	if len(errs) > 0 {
		return previouslyObservedConfig, errs
	}
	observedConfig := map[string]interface{}{}
	unstructured.SetNestedStringSlice(observedConfig, storageURLs, storageConfigURLsPath...)
	return observedConfig, errs
}

// sync reacts to a change in prereqs by finding information that is required to match another value in the cluster. This
// must be information that is logically "owned" by another component.
func (c ConfigObserver) sync() error {
	operatorConfig, err := c.operatorConfigClient.OpenShiftAPIServerOperatorConfigs().Get("instance", metav1.GetOptions{})
	if err != nil {
		return err
	}
	// don't worry about errors
	currentConfig := map[string]interface{}{}
	json.NewDecoder(bytes.NewBuffer(operatorConfig.Spec.ObservedConfig.Raw)).Decode(&currentConfig)

	var errs []error
	var observedConfigs []map[string]interface{}
	for _, i := range rand.Perm(len(c.observers)) {
		var currErrs []error
		observedConfig, currErrs := c.observers[i](c.listers, currentConfig)
		observedConfigs = append(observedConfigs, observedConfig)
		errs = append(errs, currErrs...)
	}

	mergedObservedConfig := map[string]interface{}{}
	for _, observedConfig := range observedConfigs {
		mergo.Merge(&mergedObservedConfig, observedConfig)
	}

	if !equality.Semantic.DeepEqual(currentConfig, mergedObservedConfig) {
		glog.Infof("writing updated observedConfig: %v", diff.ObjectDiff(operatorConfig.Spec.ObservedConfig.Object, mergedObservedConfig))
		operatorConfig.Spec.ObservedConfig = runtime.RawExtension{Object: &unstructured.Unstructured{Object: mergedObservedConfig}}
		updatedOperatorConfig, err := c.operatorConfigClient.OpenShiftAPIServerOperatorConfigs().Update(operatorConfig)
		if err != nil {
			errs = append(errs, fmt.Errorf("openshiftapiserveroperatorconfigs/instance: error writing updated observed config: %v", err))
		} else {
			operatorConfig = updatedOperatorConfig
		}
	}

	status := operatorConfig.Status.DeepCopy()
	if len(errs) > 0 {
		var messages []string
		for _, currentError := range errs {
			messages = append(messages, currentError.Error())
		}
		v1alpha1helpers.SetOperatorCondition(&status.Conditions, v1alpha1.OperatorCondition{
			Type:    v1alpha1.OperatorStatusTypeFailing,
			Status:  v1alpha1.ConditionTrue,
			Reason:  configObservationErrorConditionReason,
			Message: strings.Join(messages, "\n"),
		})
	} else {
		condition := v1alpha1helpers.FindOperatorCondition(status.Conditions, v1alpha1.OperatorStatusTypeFailing)
		if condition != nil && condition.Status != v1alpha1.ConditionFalse && condition.Reason == configObservationErrorConditionReason {
			condition.Status = v1alpha1.ConditionFalse
			condition.Reason = ""
			condition.Message = ""
		}
	}

	if !equality.Semantic.DeepEqual(operatorConfig.Status, status) {
		operatorConfig.Status = *status
		_, err = c.operatorConfigClient.OpenShiftAPIServerOperatorConfigs().UpdateStatus(operatorConfig)
		if err != nil {
			return err
		}
	}

	return nil
}

func observeInternalRegistryHostname(listers Listers, existingConfig map[string]interface{}) (map[string]interface{}, []error) {
	errs := []error{}
	prevObservedConfig := map[string]interface{}{}

	internalRegistryHostnamePath := []string{"imagePolicyConfig", "internalRegistryHostname"}
	if currentInternalRegistryHostname, _, _ := unstructured.NestedString(existingConfig, internalRegistryHostnamePath...); len(currentInternalRegistryHostname) > 0 {
		unstructured.SetNestedField(prevObservedConfig, currentInternalRegistryHostname, internalRegistryHostnamePath...)
	}

	if !listers.imageConfigSynced() {
		glog.Warning("images.config.openshift.io not synced")
		return prevObservedConfig, errs
	}

	observedConfig := map[string]interface{}{}
	configImage, err := listers.imageConfigLister.Get("cluster")
	if errors.IsNotFound(err) {
		glog.Warningf("image.config.openshift.io/cluster: not found")
		return observedConfig, errs
	}
	if err != nil {
		return prevObservedConfig, errs
	}
	internalRegistryHostName := configImage.Status.InternalRegistryHostname
	if len(internalRegistryHostName) > 0 {
		unstructured.SetNestedField(observedConfig, internalRegistryHostName, internalRegistryHostnamePath...)
	}
	return observedConfig, errs
}

func (c *ConfigObserver) Run(workers int, stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	glog.Infof("Starting ConfigObserver")
	defer glog.Infof("Shutting down ConfigObserver")

	if !cache.WaitForCacheSync(stopCh, c.cachesSynced...) {
		utilruntime.HandleError(fmt.Errorf("caches did not sync"))
		return
	}

	// doesn't matter what workers say, only start one.
	go wait.Until(c.runWorker, time.Second, stopCh)

	<-stopCh
}

func (c *ConfigObserver) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *ConfigObserver) processNextWorkItem() bool {
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
func (c *ConfigObserver) eventHandler() cache.ResourceEventHandler {
	return cache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj interface{}) { c.queue.Add(workQueueKey) },
		UpdateFunc: func(old, new interface{}) { c.queue.Add(workQueueKey) },
		DeleteFunc: func(obj interface{}) { c.queue.Add(workQueueKey) },
	}
}

func (c *ConfigObserver) namespaceEventHandler() cache.ResourceEventHandler {
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
