package encryptionprovider

import (
	"fmt"
	"github.com/openshift/library-go/pkg/operator/encryption/encryptionconfig"
	"k8s.io/apimachinery/pkg/api/equality"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/diff"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
	corev1listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"

	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"
)

func TestEncryptionProvider(t *testing.T) {
	encryptionCfgAnnotationKey := "ec-key"
	defaultGRs := []schema.GroupResource{
		{Group: "route.openshift.io", Resource: "routes"},
		{Group: "oauth.openshift.io", Resource: "oauthaccesstokens"},
		{Group: "oauth.openshift.io", Resource: "oauthauthorizetokens"},
	}
	grsManagedByExternalServer := sets.NewString("oauthaccesstokens.oauth.openshift.io", "oauthauthorizetokens.oauth.openshift.io")

	scenarios := []struct {
		name                 string
		initialSecrets       []*corev1.Secret
		defaultEncryptedGRs  []schema.GroupResource
		expectedEncryptedGRs []schema.GroupResource
	}{
		{
			name:                 "encryption off, default GRs returned",
			defaultEncryptedGRs:  defaultGRs,
			expectedEncryptedGRs: defaultGRs,
		},
		{
			name: "encryption on, secret without the annotation, reduced GRs returned",
			initialSecrets: []*corev1.Secret{
				func() *corev1.Secret {
					s := defaultSecret("openshift-apiserver", encryptionCfgAnnotationKey)
					delete(s.Annotations, encryptionCfgAnnotationKey)
					return s
				}(),
			},
			defaultEncryptedGRs: defaultGRs,
			expectedEncryptedGRs: []schema.GroupResource{
				{Group: "route.openshift.io", Resource: "routes"},
			},
		},
		{
			name:                 "encryption on, secret with the annotation, default GRs returned",
			initialSecrets:       []*corev1.Secret{defaultSecret("openshift-apiserver", encryptionCfgAnnotationKey)},
			defaultEncryptedGRs:  defaultGRs,
			expectedEncryptedGRs: defaultGRs,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			// test data
			fakeSecretsIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
			for _, secret := range scenario.initialSecrets {
				fakeSecretsIndexer.Add(secret)
			}
			fakeSecretsLister := corev1listers.NewSecretLister(fakeSecretsIndexer)

			// act
			target := encryptionProvider{
				allEncryptedGRs:                     scenario.defaultEncryptedGRs,
				encryptedGRsManagedByExternalServer: grsManagedByExternalServer,
				isOAuthEncryptionConfigManagedByThisOperator: IsOAuthEncryptionConfigManagedByThisOperator(
					fakeSecretsLister.Secrets(operatorclient.GlobalMachineSpecifiedConfigNamespace),
					"openshift-apiserver",
					encryptionCfgAnnotationKey,
				),
			}

			actualEncryptedGRs := target.EncryptedGRs()

			// validate
			if !equality.Semantic.DeepEqual(actualEncryptedGRs, scenario.expectedEncryptedGRs) {
				t.Errorf("incorect GRs returned: %s", diff.ObjectDiff(actualEncryptedGRs, scenario.expectedEncryptedGRs))
			}
		})
	}
}

func TestEncryptionProviderAlwaysReady(t *testing.T) {
	target := encryptionProvider{}
	ready, err := target.ShouldRunEncryptionControllers()
	if err != nil {
		t.Errorf("got an unexpected error from the provider, err = %v", err)
	}
	if !ready {
		t.Error("the provider is not ready!")
	}
}

func defaultSecret(name, annotation string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", encryptionconfig.EncryptionConfSecretName, name),
			Namespace: operatorclient.GlobalMachineSpecifiedConfigNamespace,
			Annotations: map[string]string{
				annotation: "value",
			},
		},
		Data: map[string][]byte{"encryption-config": {0xFF}},
	}
}
