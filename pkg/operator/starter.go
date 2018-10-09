package operator

import (
	"fmt"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	apiregistrationclient "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
	apiregistrationinformers "k8s.io/kube-aggregator/pkg/client/informers/externalversions"

	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/apis/openshiftapiserver/v1alpha1"
	operatorconfigclient "github.com/openshift/cluster-openshift-apiserver-operator/pkg/generated/clientset/versioned"
	operatorsv1alpha1client "github.com/openshift/cluster-openshift-apiserver-operator/pkg/generated/clientset/versioned/typed/openshiftapiserver/v1alpha1"
	operatorclientinformers "github.com/openshift/cluster-openshift-apiserver-operator/pkg/generated/informers/externalversions"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/v311_00_assets"
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

	ensureOperatorConfigExists(operator.operatorConfigClient, "v3.11.0/openshift-apiserver/operator-config.yaml")

	configObserver := NewConfigObserver(
		operatorConfigInformers.Openshiftapiserver().V1alpha1().OpenShiftAPIServerOperatorConfigs(),
		kubeInformersKubeAPIServerNamespaced,
		kubeInformersEtcdNamespaced,
		operatorConfigClient.OpenshiftapiserverV1alpha1(),
		kubeClient,
	)

	operatorConfigInformers.Start(stopCh)
	kubeInformersLocallyNamespaced.Start(stopCh)
	kubeInformersKubeAPIServerNamespaced.Start(stopCh)
	kubeInformersEtcdNamespaced.Start(stopCh)
	apiregistrationInformers.Start(stopCh)

	go operator.Run(1, stopCh)
	go configObserver.Run(1, stopCh)

	<-stopCh
	return fmt.Errorf("stopped")
}

func ensureOperatorConfigExists(client operatorsv1alpha1client.OpenShiftAPIServerOperatorConfigsGetter, path string) {
	v1alpha1Scheme := runtime.NewScheme()
	v1alpha1.Install(v1alpha1Scheme)
	v1alpha1Codecs := serializer.NewCodecFactory(v1alpha1Scheme)
	operatorConfigBytes := v311_00_assets.MustAsset(path)
	operatorConfigObj, err := runtime.Decode(v1alpha1Codecs.UniversalDecoder(v1alpha1.GroupVersion), operatorConfigBytes)
	if err != nil {
		panic(err)
	}
	requiredOperatorConfig, ok := operatorConfigObj.(*v1alpha1.OpenShiftAPIServerOperatorConfig)
	if !ok {
		panic(fmt.Sprintf("unexpected object in %s: %t", path, operatorConfigObj))
	}

	hasImageEnvVar := false
	if imagePullSpecFromEnv := os.Getenv("IMAGE"); len(imagePullSpecFromEnv) > 0 {
		hasImageEnvVar = true
		requiredOperatorConfig.Spec.ImagePullSpec = imagePullSpecFromEnv
	}

	existing, err := client.OpenShiftAPIServerOperatorConfigs().Get(requiredOperatorConfig.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		if _, err := client.OpenShiftAPIServerOperatorConfigs().Create(requiredOperatorConfig); err != nil {
			panic(err)
		}
		return
	}
	if err != nil {
		panic(err)
	}

	if !hasImageEnvVar {
		return
	}

	// If ImagePullSpec changed, update the existing config instance
	if existing.Spec.ImagePullSpec != requiredOperatorConfig.Spec.ImagePullSpec {
		existing.Spec.ImagePullSpec = requiredOperatorConfig.Spec.ImagePullSpec
		if _, err := client.OpenShiftAPIServerOperatorConfigs().Update(existing); err != nil {
			panic(err)
		}
	}
}
