package prune

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	corev1informer "k8s.io/client-go/informers/core/v1"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"

	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/v1helpers"
)

// PruneController is a controller that watches the operand pods and deletes old
// revisioned secrets that are not used anymore.
type PruneController struct {
	targetNamespace string
	secretPrefixes  []string

	secretGetter   corev1client.SecretsGetter
	podGetter      corev1client.PodsGetter
	podInformer    corev1informer.PodInformer
	secretInformer corev1informer.SecretInformer

	cachesToSync  []cache.InformerSynced
	queue         workqueue.RateLimitingInterface
	eventRecorder events.Recorder
}

const (
	pruneControllerWorkQueueKey = "key"
	numOldRevisionsToPreserve   = 5
)

// NewPruneController creates a new pruning controller
func NewPruneController(
	targetNamespace string,
	secretPrefixes []string,
	secretGetter corev1client.SecretsGetter,
	podGetter corev1client.PodsGetter,
	informers v1helpers.KubeInformersForNamespaces,
	eventRecorder events.Recorder,
) *PruneController {
	c := &PruneController{
		targetNamespace: targetNamespace,
		secretPrefixes:  secretPrefixes,

		secretGetter:   secretGetter,
		podGetter:      podGetter,
		podInformer:    informers.InformersFor(targetNamespace).Core().V1().Pods(),
		secretInformer: informers.InformersFor(targetNamespace).Core().V1().Secrets(),
		eventRecorder:  eventRecorder.WithComponentSuffix("prune-controller"),

		queue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "PruneController"),
	}

	c.podInformer.Informer().AddEventHandler(c.eventHandler())
	c.secretInformer.Informer().AddEventHandler(c.eventHandler())

	c.cachesToSync = append(
		c.cachesToSync,
		c.podInformer.Informer().HasSynced,
		c.secretInformer.Informer().HasSynced,
	)

	return c
}

func (c *PruneController) sync() error {
	klog.V(5).Info("Syncing revision pruner")

	pods, err := c.podInformer.Lister().Pods(c.targetNamespace).List(labels.SelectorFromSet(map[string]string{"apiserver": "true"}))
	if err != nil {
		return err
	}

	minRevision := minPodRevision(pods)
	if minRevision == 0 {
		return nil
	}

	secrets, err := c.secretInformer.Lister().Secrets(c.targetNamespace).List(labels.Everything())
	if err != nil {
		return err
	}

	for _, s := range secretsToBePruned(minRevision, c.secretPrefixes, secrets) {
		klog.V(4).Infof("Pruning old secret %q", s.Name)
		if err := c.secretGetter.Secrets(s.Namespace).Delete(s.Name, nil); err != nil {
			return err
		}
	}

	return nil
}

func secretsToBePruned(minRevision int, secretPrefixes []string, secrets []*corev1.Secret) []*corev1.Secret {
	// filter secrets by prefix and by revision < minRevision
	filtered := map[int][]*corev1.Secret{}
	for _, s := range secrets {
		for _, p := range secretPrefixes {
			if strings.HasPrefix(s.Name, p) {
				comps := strings.SplitAfter(s.Name, "-")
				if len(comps) == 1 {
					// skip, we cannot derive a revision
					klog.Warningf("Unexpected %q prefixed secret without a dash: %q", p, s.Name)
					break
				}
				revString := comps[len(comps)-1]
				rev, err := strconv.ParseInt(revString, 10, 32)
				if err != nil {
					// skip, we cannot derive a revision
					klog.Warningf("Unexpected %q prefixed secret %q with invalid trailing revision: %v", p, s.Name, err)
					break
				}

				if int(rev) >= minRevision {
					break
				}

				filtered[int(rev)] = append(filtered[int(rev)], s)

				break
			}
		}
	}

	sortedRevs := sortedRevisionsRecentLast(filtered)
	if len(sortedRevs) < numOldRevisionsToPreserve {
		// not enough old revisions found, nothing to prune
		return nil
	}

	revsToBePruned := sortedRevs[:len(sortedRevs)-numOldRevisionsToPreserve]

	ret := []*corev1.Secret{}
	for _, r := range revsToBePruned {
		secrets := filtered[r]
		for _, s := range secrets {
			ret = append(ret, s)
		}
	}

	return ret
}

func sortedRevisionsRecentLast(revs map[int][]*corev1.Secret) []int {
	ret := make([]int, 0, len(revs))
	for r := range revs {
		ret = append(ret, r)
	}
	sort.Ints(ret)
	return ret
}

func minPodRevision(pods []*corev1.Pod) int {
	minRevision := int64(0)
	for _, p := range pods {
		l := p.Labels["revision"]
		if len(l) == 0 {
			continue
		}
		rev, err := strconv.ParseInt(l, 10, 32)
		if err != nil || rev < 0 {
			klog.Warningf("Invalid revision label on pod %s: %q", p.Name, l)
			continue
		}
		if minRevision == 0 || rev < minRevision {
			minRevision = rev
		}
	}
	return int(minRevision)
}

func (c *PruneController) Run(workers int, stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	klog.Infof("Starting PruneController")
	defer klog.Infof("Shutting down PruneController")
	if !cache.WaitForCacheSync(stopCh, c.cachesToSync...) {
		return
	}

	// doesn't matter what workers say, only start one.
	go wait.Until(c.runWorker, time.Second, stopCh)

	<-stopCh
}

func (c *PruneController) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *PruneController) processNextWorkItem() bool {
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
func (c *PruneController) eventHandler() cache.ResourceEventHandler {
	return cache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj interface{}) { c.queue.Add(pruneControllerWorkQueueKey) },
		UpdateFunc: func(old, new interface{}) { c.queue.Add(pruneControllerWorkQueueKey) },
		DeleteFunc: func(obj interface{}) { c.queue.Add(pruneControllerWorkQueueKey) },
	}
}
