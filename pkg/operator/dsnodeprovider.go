package operator

import (
	"k8s.io/apimachinery/pkg/labels"
	appsv1informers "k8s.io/client-go/informers/apps/v1"
	corev1informers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/tools/cache"

	encryptiondeployer "github.com/openshift/library-go/pkg/operator/encryption/deployer"

	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"
)

// DaemonSetNodeProvider returns the node list from nodes matching the node selector of a DaemonSet
type DaemonSetNodeProvider struct {
	TargetNamespaceDaemonSetInformer appsv1informers.DaemonSetInformer
	NodeInformer                     corev1informers.NodeInformer
}

var (
	_ encryptiondeployer.MasterNodeProvider = &DaemonSetNodeProvider{}
)

func (p DaemonSetNodeProvider) MasterNodeNames() ([]string, error) {
	ds, err := p.TargetNamespaceDaemonSetInformer.Lister().DaemonSets(operatorclient.TargetNamespace).Get("apiserver")
	if err != nil {
		return nil, err
	}

	nodes, err := p.NodeInformer.Lister().List(labels.SelectorFromSet(ds.Spec.Template.Spec.NodeSelector))
	if err != nil {
		return nil, err
	}

	ret := make([]string, 0, len(nodes))
	for _, n := range nodes {
		ret = append(ret, n.Name)
	}

	return ret, nil
}

func (p DaemonSetNodeProvider) AddEventHandler(handler cache.ResourceEventHandler) []cache.InformerSynced {
	p.TargetNamespaceDaemonSetInformer.Informer().AddEventHandler(handler)
	p.NodeInformer.Informer().AddEventHandler(handler)

	return []cache.InformerSynced{
		p.TargetNamespaceDaemonSetInformer.Informer().HasSynced,
		p.NodeInformer.Informer().HasSynced,
	}
}
