package operator

import (
	"fmt"
	"time"

	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	apiregistrationclient "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
	apiregistrationinformers "k8s.io/kube-aggregator/pkg/client/informers/externalversions"

	operatorconfigclient "github.com/openshift/cluster-openshift-apiserver-operator/pkg/generated/clientset/versioned"
	operatorclientinformers "github.com/openshift/cluster-openshift-apiserver-operator/pkg/generated/informers/externalversions"
)

func RunOperator(clientConfig *rest.Config, stopCh <-chan struct{}) error {
	kubeClient, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		panic(err)
	}
	apiregistrationv1Client, err := apiregistrationclient.NewForConfig(clientConfig)
	if err != nil {
		panic(err)
	}
	operatorConfigClient, err := operatorconfigclient.NewForConfig(clientConfig)
	if err != nil {
		panic(err)
	}

	operatorConfigInformers := operatorclientinformers.NewSharedInformerFactory(operatorConfigClient, 10*time.Minute)
	kubeInformersLocallyNamespaced := kubeinformers.NewFilteredSharedInformerFactory(kubeClient, 10*time.Minute, targetNamespaceName, nil)
	kubeInformersKubeAPIServerNamespaced := kubeinformers.NewFilteredSharedInformerFactory(kubeClient, 10*time.Minute, kubeAPIServerNamespaceName, nil)
	apiregistrationInformers := apiregistrationinformers.NewSharedInformerFactory(apiregistrationv1Client, 10*time.Minute)

	operator := NewKubeApiserverOperator(
		operatorConfigInformers.Openshiftapiserver().V1alpha1().OpenShiftAPIServerOperatorConfigs(),
		kubeInformersLocallyNamespaced,
		kubeInformersKubeAPIServerNamespaced,
		apiregistrationInformers,
		operatorConfigClient.OpenshiftapiserverV1alpha1(),
		kubeClient.AppsV1(),
		kubeClient.CoreV1(),
		kubeClient.RbacV1(),
		apiregistrationv1Client.ApiregistrationV1(),
	)

	operatorConfigInformers.Start(stopCh)
	kubeInformersLocallyNamespaced.Start(stopCh)
	apiregistrationInformers.Start(stopCh)

	operator.Run(1, stopCh)
	return fmt.Errorf("stopped")
}
