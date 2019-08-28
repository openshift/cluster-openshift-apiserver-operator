package resourcesynccontroller

import (
	"net/http"

	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/resourcesynccontroller"
	"github.com/openshift/library-go/pkg/operator/v1helpers"

	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"
)

func NewResourceSyncController(
	operatorConfigClient v1helpers.OperatorClient,
	kubeInformersForNamespaces v1helpers.KubeInformersForNamespaces,
	configMapsGetter corev1client.ConfigMapsGetter,
	secretsGetter corev1client.SecretsGetter,
	eventRecorder events.Recorder) (*resourcesynccontroller.ResourceSyncController, http.Handler, error) {

	resourceSyncController := resourcesynccontroller.NewResourceSyncController(
		operatorConfigClient,
		kubeInformersForNamespaces,
		secretsGetter,
		configMapsGetter,
		eventRecorder,
	)
	if err := resourceSyncController.SyncConfigMap(
		resourcesynccontroller.ResourceLocation{Namespace: operatorclient.TargetNamespace, Name: "etcd-serving-ca"},
		resourcesynccontroller.ResourceLocation{Namespace: operatorclient.GlobalUserSpecifiedConfigNamespace, Name: "etcd-serving-ca"},
	); err != nil {
		return nil, nil, err
	}
	if err := resourceSyncController.SyncSecret(
		resourcesynccontroller.ResourceLocation{Namespace: operatorclient.TargetNamespace, Name: "etcd-client"},
		resourcesynccontroller.ResourceLocation{Namespace: operatorclient.GlobalUserSpecifiedConfigNamespace, Name: "etcd-client"},
	); err != nil {
		return nil, nil, err
	}
	if err := resourceSyncController.SyncConfigMap(
		resourcesynccontroller.ResourceLocation{Namespace: operatorclient.TargetNamespace, Name: "client-ca"},
		resourcesynccontroller.ResourceLocation{Namespace: operatorclient.GlobalMachineSpecifiedConfigNamespace, Name: "kube-apiserver-client-ca"},
	); err != nil {
		return nil, nil, err
	}
	if err := resourceSyncController.SyncConfigMap(
		resourcesynccontroller.ResourceLocation{Namespace: operatorclient.TargetNamespace, Name: "aggregator-client-ca"},
		resourcesynccontroller.ResourceLocation{Namespace: operatorclient.GlobalMachineSpecifiedConfigNamespace, Name: "kube-apiserver-aggregator-client-ca"},
	); err != nil {
		return nil, nil, err
	}
	if err := resourceSyncController.SyncConfigMap(
		resourcesynccontroller.ResourceLocation{Namespace: operatorclient.TargetNamespace, Name: "trusted-ca-bundle"},
		resourcesynccontroller.ResourceLocation{Namespace: operatorclient.OperatorNamespace, Name: "trusted-ca-bundle"},
	); err != nil {
		return nil, nil, err
	}

	return resourceSyncController, resourcesynccontroller.NewDebugHandler(resourceSyncController), nil
}
