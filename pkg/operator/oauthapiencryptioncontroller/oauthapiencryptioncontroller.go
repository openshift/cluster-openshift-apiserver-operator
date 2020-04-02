package oauthapiencryptioncontroller

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"
	"github.com/openshift/library-go/pkg/controller/factory"
	"github.com/openshift/library-go/pkg/operator/encryption/encryptionconfig"
	encryptionsecret "github.com/openshift/library-go/pkg/operator/encryption/secrets"
	encryptionstate "github.com/openshift/library-go/pkg/operator/encryption/state"
	"github.com/openshift/library-go/pkg/operator/events"
	operatorv1helpers "github.com/openshift/library-go/pkg/operator/v1helpers"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	corev1listers "k8s.io/client-go/listers/core/v1"
)

const (
	EncryptionConfigManagedBy      = "encryption.apiserver.operator.openshift.io/managed-by"
	encryptionConfigManagedByValue = `WARNING: DO NOT REMOVE.
This annotation indicates that OAS-O manages this secret.`
)

type oauthAPIServerController struct {
	oauthAPIServerTargetNamespace string

	secretLister corev1listers.SecretNamespaceLister
	secretClient corev1client.SecretInterface
}

// New creates OAuthAPIServerController that will manage encryption-config-openshift-oauth-apiserver in openshift-config-managed namespace as described in https://github.com/openshift/enhancements/blob/master/enhancements/etcd/etcd-encryption-for-separate-oauth-apis.md
// Note that this code will be removed in the future release (4.6)
func New(
	name string,
	oauthAPIServerTargetNamespace string,
	secretClient corev1client.SecretsGetter,
	kubeInformersForNamespaces operatorv1helpers.KubeInformersForNamespaces,
	eventRecorder events.Recorder) factory.Controller {

	controllerFactory := factory.New()
	target := &oauthAPIServerController{
		oauthAPIServerTargetNamespace: oauthAPIServerTargetNamespace,
		secretLister:                  kubeInformersForNamespaces.InformersFor(operatorclient.GlobalMachineSpecifiedConfigNamespace).Core().V1().Secrets().Lister().Secrets(operatorclient.GlobalMachineSpecifiedConfigNamespace),
		secretClient:                  secretClient.Secrets(operatorclient.GlobalMachineSpecifiedConfigNamespace),
	}

	controllerFactory.WithSync(target.sync)
	controllerFactory.WithInformers(kubeInformersForNamespaces.InformersFor(operatorclient.GlobalMachineSpecifiedConfigNamespace).Core().V1().Secrets().Informer())
	return controllerFactory.ToController(name, eventRecorder.WithComponentSuffix("oauth-apiserver-encryption-cfg-sync-controller"))
}

// sync starts managing oauth-apiserver encryption config (encryption-config-openshift-oauth-apiserver in openshift-config-managed namespace) as described in https://github.com/openshift/enhancements/blob/master/enhancements/etcd/etcd-encryption-for-separate-oauth-apis.md
//
// case 1: if the secret doesn't exist and encryption is on then:
//         - it will add the secret with the annotation (must be atomic operation) because CAO will start managing its own config iff the secret exist without the annotation
//
// case 2: if the secret exists and it is annotated
//         - it will simply start synchronisation
//
// case 3: no-op: when the secret exits but it doesn't have the annotation - that means it was created by CAO in 4.6 and this is downgrade
// case 4: no-op: when the secret doesn't exist and encryption is off
//
// drawbacks:
// - it will not recover when the annotation was manually removed by a user,
//   to recover we would have to put a value in the annotation instead and coordinate OAS-A and CAO but then we would have to remember to add it in CAO (4.6) as well
func (c *oauthAPIServerController) sync(ctx context.Context, controllerContext factory.SyncContext) error {
	openshiftAPIServerEncryptionCfg, err := c.secretLister.Get(fmt.Sprintf("%s-%s", encryptionconfig.EncryptionConfSecretName, operatorclient.TargetNamespace))
	if apierrors.IsNotFound(err) {
		return nil // case 4: encryption off
	}
	if err != nil {
		return err
	}

	oauthAPIServerEncryptionCfgName := fmt.Sprintf("%s-%s", encryptionconfig.EncryptionConfSecretName, c.oauthAPIServerTargetNamespace)
	oauthAPIServerEncryptionCfg, err := c.secretLister.Get(oauthAPIServerEncryptionCfgName)
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}
	if apierrors.IsNotFound(err) {
		// case 1: create with annotation
		if _, exists := openshiftAPIServerEncryptionCfg.Data[encryptionconfig.EncryptionConfSecretKey]; !exists {
			return fmt.Errorf("%s/%s doesn't contain the required key %q", openshiftAPIServerEncryptionCfg.Namespace, openshiftAPIServerEncryptionCfg.Name, encryptionconfig.EncryptionConfSecretKey)
		}
		encryptionCfg := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      oauthAPIServerEncryptionCfgName,
				Namespace: operatorclient.GlobalMachineSpecifiedConfigNamespace,
				Annotations: map[string]string{
					EncryptionConfigManagedBy:                encryptionConfigManagedByValue,
					encryptionstate.KubernetesDescriptionKey: encryptionstate.KubernetesDescriptionScaryValue,
				},
				Finalizers: []string{encryptionsecret.EncryptionSecretFinalizer},
			},
			Data: map[string][]byte{},
		}
		encryptionCfg.Data[encryptionconfig.EncryptionConfSecretKey] = openshiftAPIServerEncryptionCfg.Data[encryptionconfig.EncryptionConfSecretKey]

		_, err := c.secretClient.Create(ctx, encryptionCfg, metav1.CreateOptions{})
		if err != nil {
			return err
		}
		controllerContext.Recorder().Eventf("SecretCreated", "Created %s in %s namespace because it was missing", oauthAPIServerEncryptionCfgName, operatorclient.GlobalMachineSpecifiedConfigNamespace)
		return nil
	}
	if _, exist := oauthAPIServerEncryptionCfg.Annotations[EncryptionConfigManagedBy]; exist {
		// case 2: exists and it is annotated
		oauthEncryptionCfgData := oauthAPIServerEncryptionCfg.Data[encryptionconfig.EncryptionConfSecretKey]
		oasEncryptionCfgData := openshiftAPIServerEncryptionCfg.Data[encryptionconfig.EncryptionConfSecretKey]
		if !equality.Semantic.DeepEqual(oauthEncryptionCfgData, oasEncryptionCfgData) {
			encryptionCfg := oauthAPIServerEncryptionCfg.DeepCopy()
			encryptionCfg.Data[encryptionconfig.EncryptionConfSecretName] = oasEncryptionCfgData
			_, err := c.secretClient.Update(ctx, encryptionCfg, metav1.UpdateOptions{})
			if err != nil {
				return err
			}
			controllerContext.Recorder().Eventf("SecretUpdated", "Updates %s in %s namespace because it was out of date with %s ", oauthAPIServerEncryptionCfgName, operatorclient.GlobalMachineSpecifiedConfigNamespace, openshiftAPIServerEncryptionCfg.Name)
			return nil
		}
		return nil
	}

	// case 3: no-op the secret is managed by CAO
	return nil
}
