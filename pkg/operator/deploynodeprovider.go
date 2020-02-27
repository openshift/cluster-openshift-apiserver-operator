package operator

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	appsv1informers "k8s.io/client-go/informers/apps/v1"
	corev1informers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/tools/cache"

	encryptiondeployer "github.com/openshift/library-go/pkg/operator/encryption/deployer"

	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"
)

// DeploymentNodeProvider returns the node list from nodes matching the node selector of a Deployment
type DeploymentNodeProvider struct {
	TargetNamespaceDeploymentInformer appsv1informers.DeploymentInformer
	NodeInformer                      corev1informers.NodeInformer
}

var (
	_ encryptiondeployer.MasterNodeProvider = &DeploymentNodeProvider{}
)

func (p DeploymentNodeProvider) MasterNodeNames() ([]string, error) {
	deploy, err := p.TargetNamespaceDeploymentInformer.Lister().Deployments(operatorclient.TargetNamespace).Get("apiserver")
	if err != nil && errors.IsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	nodes, err := p.NodeInformer.Lister().List(labels.SelectorFromSet(deploy.Spec.Template.Spec.NodeSelector))
	if err != nil {
		return nil, err
	}

	ret := make([]string, 0, len(nodes))
	for _, n := range nodes {
		ret = append(ret, n.Name)
	}

	return ret, nil
}

func (p DeploymentNodeProvider) AddEventHandler(handler cache.ResourceEventHandler) []cache.InformerSynced {
	p.TargetNamespaceDeploymentInformer.Informer().AddEventHandler(handler)
	p.NodeInformer.Informer().AddEventHandler(handler)

	return []cache.InformerSynced{
		p.TargetNamespaceDeploymentInformer.Informer().HasSynced,
		p.NodeInformer.Informer().HasSynced,
	}
}
