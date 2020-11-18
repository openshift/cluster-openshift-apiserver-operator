package oauthapiencryption

import (
	"context"
	"fmt"

	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"
	"github.com/openshift/library-go/pkg/controller/factory"
	"github.com/openshift/library-go/pkg/operator/encryption/encryptionconfig"
	"github.com/openshift/library-go/pkg/operator/encryption/statemachine"
	"github.com/openshift/library-go/pkg/operator/events"

	operatorv1helpers "github.com/openshift/library-go/pkg/operator/v1helpers"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	corev1listers "k8s.io/client-go/listers/core/v1"
)

const (
	EncryptionConfigManagedBy      = "encryption.apiserver.operator.openshift.io/managed-by"
	encryptionConfigManagedByValue = `WARNING: DO NOT REMOVE.
This annotation indicates that OAS-O manages this secret.`
)

type oauthEncryptionConfigSyncController struct {
	oauthAPIServerTargetNamespace string

	secretListerConfigManaged           corev1listers.SecretNamespaceLister
	secretListerForOAuthTargetNamespace corev1listers.SecretNamespaceLister
	secretClient                        corev1client.SecretsGetter

	deployer statemachine.Deployer
}

// NewEncryptionConfigSyncController creates OAuthAPIServerController that will clean up the encryption configuration for oauth-apiserver as described in https://github.com/openshift/enhancements/blob/master/enhancements/etcd/etcd-encryption-for-separate-oauth-apis.md
// TODO: remove it in 4.8
func NewEncryptionConfigSyncController(
	name string,
	oauthAPIServerTargetNamespace string,
	secretClient corev1client.SecretsGetter,
	kubeInformersForNamespaces operatorv1helpers.KubeInformersForNamespaces,
	deployer statemachine.Deployer,
	eventRecorder events.Recorder) factory.Controller {

	controllerFactory := factory.New()
	target := &oauthEncryptionConfigSyncController{
		oauthAPIServerTargetNamespace: oauthAPIServerTargetNamespace,

		secretListerConfigManaged:           kubeInformersForNamespaces.InformersFor(operatorclient.GlobalMachineSpecifiedConfigNamespace).Core().V1().Secrets().Lister().Secrets(operatorclient.GlobalMachineSpecifiedConfigNamespace),
		secretListerForOAuthTargetNamespace: kubeInformersForNamespaces.InformersFor(oauthAPIServerTargetNamespace).Core().V1().Secrets().Lister().Secrets(oauthAPIServerTargetNamespace),

		secretClient: secretClient,
		deployer:     deployer,
	}

	controllerFactory.WithSync(target.sync)
	controllerFactory.WithInformers(kubeInformersForNamespaces.InformersFor(operatorclient.GlobalMachineSpecifiedConfigNamespace).Core().V1().Secrets().Informer())
	return controllerFactory.ToController(name, eventRecorder.WithComponentSuffix("oauth-apiserver-encryption-cfg-sync-controller"))
}

// sync starts cleaning up oauth-apiserver encryption config (encryption-config-openshift-oauth-apiserver in openshift-config-managed namespace AND encryption-config in openshift-oauth-apiserver)
//
// case 1: no-op: when the encryption config (oauth-apiserver) doesn't exist in the global namespace and the encryption is off, let CAO manage its own encryption config
//
// case 2: remove the annotation from the encryption configs
//      a: if the encryption config (oauth-apiserver) in the global namespace exists and it is annotated (an upgrade from 4.6)
//      b: if the encryption config (oauth-apiserver) in the target namespace also exists and it is annotated (an upgrade from 4.6)
//
// case 3: no-op: when both encryption configs exit but don't have the annotation - that means they were created by CAO in 4.7 OR this is downgrade
func (c *oauthEncryptionConfigSyncController) sync(ctx context.Context, controllerContext factory.SyncContext) error {
	oauthAPIServerEncryptionCfgName := fmt.Sprintf("%s-%s", encryptionconfig.EncryptionConfSecretName, c.oauthAPIServerTargetNamespace)
	oauthAPIGlobalConfig, err := c.secretListerConfigManaged.Get(oauthAPIServerEncryptionCfgName)
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}
	if apierrors.IsNotFound(err) {
		// case 1: no-op let CAO manage its own encryption config
		return nil
	}

	// before messing with the encryption configs wait for stability
	_, converged, err := c.deployer.DeployedEncryptionConfigSecret()
	if err != nil || !converged {
		return err
	}

	if _, exist := oauthAPIGlobalConfig.Annotations[EncryptionConfigManagedBy]; exist {
		// case 2a: the encryptionCfg exists and it is annotated in the global namespace, we need to remove the annotation
		if err := c.removeAnnotationAndRecordEvent(ctx, operatorclient.GlobalMachineSpecifiedConfigNamespace, oauthAPIGlobalConfig.DeepCopy(), controllerContext.Recorder()); err != nil {
			return err
		}
	}

	oauthAPIConfig, err := c.secretListerForOAuthTargetNamespace.Get(encryptionconfig.EncryptionConfSecretName)
	if err != nil {
		// this probably means that the encryption cfg hasn't been synced to the target namespace
		return err
	}
	if _, exist := oauthAPIConfig.Annotations[EncryptionConfigManagedBy]; exist {
		// case 2b: the encryptionCfg exists and it is annotated in the target (oauthAPIServerTargetNamespace) namespace, we need to remove the annotation
		//
		// note: this is nice to have since both operators watch the encryption config in the global namespace
		return c.removeAnnotationAndRecordEvent(ctx, c.oauthAPIServerTargetNamespace, oauthAPIConfig.DeepCopy(), controllerContext.Recorder())
	}

	// case 3: no-op the encryption config is already managed by CAO
	return nil
}

func (c *oauthEncryptionConfigSyncController) removeAnnotationAndRecordEvent(ctx context.Context, ns string, encryptionCfg *corev1.Secret, recorder events.Recorder) error {
	delete(encryptionCfg.Annotations, EncryptionConfigManagedBy)
	_, err := c.secretClient.Secrets(ns).Update(ctx, encryptionCfg, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	recorder.Eventf("EncryptionConfigUpdated", "Removed the %q annotation from %s in %s namespace", EncryptionConfigManagedBy, encryptionCfg.Name, ns)
	return nil
}
