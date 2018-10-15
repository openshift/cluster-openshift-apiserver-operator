package operator

import (
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	apiregistrationclient "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
	apiregistrationinformers "k8s.io/kube-aggregator/pkg/client/informers/externalversions"

	operatorv1alpha1 "github.com/openshift/api/operator/v1alpha1"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/apis/openshiftapiserver/v1alpha1"
	operatorconfigclient "github.com/openshift/cluster-openshift-apiserver-operator/pkg/generated/clientset/versioned"
	operatorclientinformers "github.com/openshift/cluster-openshift-apiserver-operator/pkg/generated/informers/externalversions"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/v311_00_assets"
	"github.com/openshift/library-go/pkg/operator/status"
	"github.com/openshift/library-go/pkg/operator/v1alpha1helpers"
)

func RunOperator(clientConfig *rest.Config, stopCh <-chan struct{}) error {
	kubeClient, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return err
	}
	apiregistrationv1Client, err := apiregistrationclient.NewForConfig(clientConfig)
	if err != nil {
		return err
	}
	operatorConfigClient, err := operatorconfigclient.NewForConfig(clientConfig)
	if err != nil {
		return err
	}
	dynamicClient, err := dynamic.NewForConfig(clientConfig)
	if err != nil {
		return err
	}

	v1alpha1helpers.EnsureOperatorConfigExists(
		dynamicClient,
		v311_00_assets.MustAsset("v3.11.0/openshift-apiserver/operator-config.yaml"),
		schema.GroupVersionResource{Group: v1alpha1.GroupName, Version: "v1alpha1", Resource: "openshiftapiserveroperatorconfigs"},
		v1alpha1helpers.GetImageEnv(),
	)

	operatorConfigInformers := operatorclientinformers.NewSharedInformerFactory(operatorConfigClient, 10*time.Minute)
	kubeInformersLocallyNamespaced := kubeinformers.NewFilteredSharedInformerFactory(kubeClient, 10*time.Minute, targetNamespaceName, nil)
	kubeInformersKubeAPIServerNamespaced := kubeinformers.NewFilteredSharedInformerFactory(kubeClient, 10*time.Minute, kubeAPIServerNamespaceName, nil)
	kubeInformersEtcdNamespaced := kubeinformers.NewFilteredSharedInformerFactory(kubeClient, 10*time.Minute, etcdNamespaceName, nil)
	apiregistrationInformers := apiregistrationinformers.NewSharedInformerFactory(apiregistrationv1Client, 10*time.Minute)

	operator := NewKubeApiserverOperator(
		operatorConfigInformers.Openshiftapiserver().V1alpha1().OpenShiftAPIServerOperatorConfigs(),
		kubeInformersLocallyNamespaced,
		kubeInformersKubeAPIServerNamespaced,
		apiregistrationInformers,
		operatorConfigClient.OpenshiftapiserverV1alpha1(),
		kubeClient,
		apiregistrationv1Client.ApiregistrationV1(),
	)

	configObserver := NewConfigObserver(
		operatorConfigInformers.Openshiftapiserver().V1alpha1().OpenShiftAPIServerOperatorConfigs(),
		kubeInformersKubeAPIServerNamespaced,
		kubeInformersEtcdNamespaced,
		operatorConfigClient.OpenshiftapiserverV1alpha1(),
		kubeClient,
	)

	clusterOperatorStatus := status.NewClusterOperatorStatusController(
		"openshift-apiserver",
		"openshift-apiserver",
		dynamicClient,
		&operatorStatusProvider{informers: operatorConfigInformers},
	)

	operatorConfigInformers.Start(stopCh)
	kubeInformersLocallyNamespaced.Start(stopCh)
	kubeInformersKubeAPIServerNamespaced.Start(stopCh)
	kubeInformersEtcdNamespaced.Start(stopCh)
	apiregistrationInformers.Start(stopCh)

	go operator.Run(1, stopCh)
	go configObserver.Run(1, stopCh)
	go clusterOperatorStatus.Run(1, stopCh)

	<-stopCh
	return fmt.Errorf("stopped")
}

type operatorStatusProvider struct {
	informers operatorclientinformers.SharedInformerFactory
}

func (p *operatorStatusProvider) Informer() cache.SharedIndexInformer {
	return p.informers.Openshiftapiserver().V1alpha1().OpenShiftAPIServerOperatorConfigs().Informer()
}

func (p *operatorStatusProvider) CurrentStatus() (operatorv1alpha1.OperatorStatus, error) {
	instance, err := p.informers.Openshiftapiserver().V1alpha1().OpenShiftAPIServerOperatorConfigs().Lister().Get("instance")
	if err != nil {
		return operatorv1alpha1.OperatorStatus{}, err
	}

	return instance.Status.OperatorStatus, nil
}
