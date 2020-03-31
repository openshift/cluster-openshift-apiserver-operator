package encryptionprovider

import (
	"fmt"
	"github.com/openshift/library-go/pkg/operator/encryption/encryptionconfig"
	"github.com/openshift/library-go/pkg/operator/events"
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
		expectedEvents       []string
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
			expectedEvents: []string{"EncryptedGRsChanged"},
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
			eventRecorder := events.NewInMemoryRecorder("")
			fakeSecretsIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
			for _, secret := range scenario.initialSecrets {
				fakeSecretsIndexer.Add(secret)
			}
			fakeSecretsLister := corev1listers.NewSecretLister(fakeSecretsIndexer)

			// act
			target := encryptionProvider{
				oauthAPIServerTargetNamespace:       "openshift-apiserver",
				oauthEncryptionCfgAnnotationKey:     encryptionCfgAnnotationKey,
				allEncryptedGRs:                     scenario.defaultEncryptedGRs,
				encryptedGRsManagedByExternalServer: grsManagedByExternalServer,
				secretLister:                        fakeSecretsLister.Secrets(operatorclient.GlobalMachineSpecifiedConfigNamespace),
				eventRecorder:                       eventRecorder,
			}

			actualEncryptedGRs := target.EncryptedGRs()

			// validate
			if !equality.Semantic.DeepEqual(actualEncryptedGRs, scenario.expectedEncryptedGRs) {
				t.Errorf("incorect GRs returned: %s", diff.ObjectDiff(actualEncryptedGRs, scenario.expectedEncryptedGRs))
			}
			if err := validateEventsReason(eventRecorder.Events(), scenario.expectedEvents); err != nil {
				t.Error(err)
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

func validateEventsReason(actualEvents []*corev1.Event, expectedReasons []string) error {
	if len(actualEvents) != len(expectedReasons) {
		return fmt.Errorf("expected to get %d events but got %d\nexpected=%v \n got=%v", len(expectedReasons), len(actualEvents), expectedReasons, eventReasons(actualEvents))
	}
	for i, e := range actualEvents {
		if got, expected := e.Reason, expectedReasons[i]; got != expected {
			return fmt.Errorf("at %d got %s, expected %s", i, got, expected)
		}
	}
	return nil
}

func eventReasons(events []*corev1.Event) []string {
	ret := make([]string, 0, len(events))
	for _, ev := range events {
		ret = append(ret, ev.Reason)
	}
	return ret
}
