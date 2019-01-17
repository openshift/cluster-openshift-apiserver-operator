package resourcesynccontroller

import (
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"
	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/resourcesynccontroller"
	"github.com/openshift/library-go/pkg/operator/v1helpers"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
)

func NewResourceSyncController(
	operatorConfigClient v1helpers.OperatorClient,
	kubeInformersForNamespaces map[string]informers.SharedInformerFactory,
	kubeClient kubernetes.Interface,
	eventRecorder events.Recorder) (*resourcesynccontroller.ResourceSyncController, error) {

	resourceSyncController := resourcesynccontroller.NewResourceSyncController(
		operatorConfigClient,
		kubeInformersForNamespaces,
		kubeClient,
		eventRecorder,
	)
	if err := resourceSyncController.SyncConfigMap(
		resourcesynccontroller.ResourceLocation{Namespace: operatorclient.TargetNamespaceName, Name: "etcd-serving-ca"},
		resourcesynccontroller.ResourceLocation{Namespace: operatorclient.EtcdNamespaceName, Name: "etcd-serving-ca"},
	); err != nil {
		return nil, err
	}
	if err := resourceSyncController.SyncSecret(
		resourcesynccontroller.ResourceLocation{Namespace: operatorclient.TargetNamespaceName, Name: "etcd-client"},
		resourcesynccontroller.ResourceLocation{Namespace: operatorclient.EtcdNamespaceName, Name: "etcd-client"},
	); err != nil {
		return nil, err
	}
	if err := resourceSyncController.SyncConfigMap(
		resourcesynccontroller.ResourceLocation{Namespace: operatorclient.TargetNamespaceName, Name: "client-ca"},
		resourcesynccontroller.ResourceLocation{Namespace: operatorclient.KubeAPIServerNamespaceName, Name: "client-ca"},
	); err != nil {
		return nil, err
	}

	return resourceSyncController, nil
}
