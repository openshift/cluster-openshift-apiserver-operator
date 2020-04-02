package encryptionprovider

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
	corev1listers "k8s.io/client-go/listers/core/v1"

	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"
	"github.com/openshift/library-go/pkg/operator/encryption/controllers"
	"github.com/openshift/library-go/pkg/operator/encryption/encryptionconfig"
	operatorv1helpers "github.com/openshift/library-go/pkg/operator/v1helpers"
)

type encryptionProvider struct {
	allEncryptedGRs                     []schema.GroupResource
	encryptedGRsManagedByExternalServer sets.String

	isOAuthEncryptionConfigManagedByThisOperator func() bool
}

var _ controllers.Provider = &encryptionProvider{}

func New(
	oauthAPIServerTargetNamespace string,
	oauthEncryptionCfgAnnotationKey string,
	allEncryptedGRs []schema.GroupResource,
	encryptedGRsManagedByExternalServer sets.String,
	kubeInformersForNamespaces operatorv1helpers.KubeInformersForNamespaces) *encryptionProvider {
	return &encryptionProvider{
		allEncryptedGRs:                     allEncryptedGRs,
		encryptedGRsManagedByExternalServer: encryptedGRsManagedByExternalServer,
		isOAuthEncryptionConfigManagedByThisOperator: IsOAuthEncryptionConfigManagedByThisOperator(
			kubeInformersForNamespaces.InformersFor(operatorclient.GlobalMachineSpecifiedConfigNamespace).Core().V1().Secrets().Lister().Secrets(operatorclient.GlobalMachineSpecifiedConfigNamespace),
			oauthAPIServerTargetNamespace,
			oauthEncryptionCfgAnnotationKey,
		),
	}
}

// EncryptedGRs returns resources that need to be encrypted
// Note: the list can change depending on the existence and attached annotations of encryption-config-openshift-oauth-apiserver in openshift-config-managed namespace as described in https://github.com/openshift/enhancements/blob/master/enhancements/etcd/etcd-encryption-for-separate-oauth-apis.md
//
// case 1 encryption off or the secret was annotated - return authoritative list of EncryptedGRs
// case 2 otherwise reduce the authoritative list and let CAO manage its own encryption configuration
//
// TODO:
// - change the code in 4.7 so that it only returns a static list (https://bugzilla.redhat.com/show_bug.cgi?id=1819723)
func (p *encryptionProvider) EncryptedGRs() []schema.GroupResource {
	// case 1 - we are in charge
	if p.isOAuthEncryptionConfigManagedByThisOperator() {
		return p.allEncryptedGRs
	}

	// case 2 - CAO is in charge, reduce the list
	newEncryptedGRsToManage := []schema.GroupResource{}
	for _, gr := range p.allEncryptedGRs {
		if p.encryptedGRsManagedByExternalServer.Has(gr.String()) {
			continue
		}
		newEncryptedGRsToManage = append(newEncryptedGRsToManage, gr)
	}
	return newEncryptedGRsToManage
}

// ShouldRunEncryptionControllers indicates whether external preconditions are satisfied so that encryption controllers can start synchronizing
func (p *encryptionProvider) ShouldRunEncryptionControllers() (bool, error) {
	return true, nil // always ready
}

// IsOAuthEncryptionConfigManagedByThisOperator determines whether this operator is in charge of encryption-config-openshift-oauth-apiserver
//
// case 1 encryption off or the secret was annotated - we are in charge
// case 2 otherwise let CAO manage its own encryption configuration
// TODO:
// - change case 1 in in 4.7 so that this operator doesn't manage CAO's encryption config when encryption is off
func IsOAuthEncryptionConfigManagedByThisOperator(secretLister corev1listers.SecretNamespaceLister, oauthAPIServerTargetNamespace string, oauthEncryptionCfgAnnotationKey string) func() bool {
	return func() bool {
		oauthAPIServerEncryptionCfg, err := secretLister.Get(fmt.Sprintf("%s-%s", encryptionconfig.EncryptionConfSecretName, oauthAPIServerTargetNamespace))
		if err != nil {
			// note that it's okay to return true on an error because:
			// - the only type of error we can get here (cache) is NotFound which means that the encryption is off
			return true // case 1 - we are in charge
		}

		if _, exist := oauthAPIServerEncryptionCfg.Annotations[oauthEncryptionCfgAnnotationKey]; exist {
			return true // case 1 - we are in charge
		}
		return false // case 2 - CAO is in charge
	}
}
