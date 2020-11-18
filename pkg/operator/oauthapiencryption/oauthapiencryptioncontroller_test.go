package oauthapiencryption

import (
	"context"
	"fmt"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/meta"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/diff"
	"k8s.io/client-go/kubernetes/fake"
	corev1listers "k8s.io/client-go/listers/core/v1"
	clientgotesting "k8s.io/client-go/testing"
	"k8s.io/client-go/tools/cache"

	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"
	"github.com/openshift/library-go/pkg/controller/factory"
	"github.com/openshift/library-go/pkg/operator/encryption/encryptionconfig"
	encryptionsecret "github.com/openshift/library-go/pkg/operator/encryption/secrets"
	encryptionstate "github.com/openshift/library-go/pkg/operator/encryption/state"
	"github.com/openshift/library-go/pkg/operator/events"
)

func TestOAuthAPIServerController(t *testing.T) {
	oAuthAPIServerTargetNamespace := "openshift-oauth-apiserver"

	scenarios := []struct {
		name                      string
		fakeDeployer              *fakeDeployer
		initialSecretsGlobal      []*corev1.Secret
		initialSecretsOAuthTarget []*corev1.Secret
		validateFunc              func(ts *testing.T, actions []clientgotesting.Action)

		expectedActions []string
		expectedEvents  []string
	}{
		{
			name: "test cases 1,4 - the encryption config in the global ns doesn't exist, encryption is off",
		},
		{
			name:         "test cases 2, 2.1 - the encryption configs for oauth-apiserver in the global and target ns exists and are annotated",
			fakeDeployer: newFakeDeployer(true, nil),
			initialSecretsGlobal: []*corev1.Secret{
				defaultSecret(fmt.Sprintf("%s-%s", encryptionconfig.EncryptionConfSecretName, oAuthAPIServerTargetNamespace), operatorclient.GlobalMachineSpecifiedConfigNamespace),
			},
			initialSecretsOAuthTarget: []*corev1.Secret{
				defaultSecret("encryption-config", oAuthAPIServerTargetNamespace),
			},
			expectedActions: []string{"update:secrets:openshift-config-managed:encryption-config-openshift-oauth-apiserver", "update:secrets:openshift-oauth-apiserver:encryption-config"},
			expectedEvents:  []string{"EncryptionConfigUpdated", "EncryptionConfigUpdated"},
			validateFunc: func(ts *testing.T, actions []clientgotesting.Action) {
				validatedSecrets := []bool{}
				for _, action := range actions {
					if action.Matches("update", "secrets") {
						var expectedSecret *corev1.Secret
						updateAction := action.(clientgotesting.UpdateAction)
						actualSecret := updateAction.GetObject().(*corev1.Secret)

						if actualSecret.Namespace == operatorclient.GlobalMachineSpecifiedConfigNamespace {
							expectedSecret = defaultSecret(fmt.Sprintf("%s-%s", encryptionconfig.EncryptionConfSecretName, oAuthAPIServerTargetNamespace), operatorclient.GlobalMachineSpecifiedConfigNamespace)
						} else {
							expectedSecret = defaultSecret(encryptionconfig.EncryptionConfSecretName, oAuthAPIServerTargetNamespace)
						}
						delete(expectedSecret.Annotations, EncryptionConfigManagedBy)

						if !equality.Semantic.DeepEqual(actualSecret, expectedSecret) {
							ts.Errorf(diff.ObjectDiff(actualSecret, expectedSecret))
						}
						validatedSecrets = append(validatedSecrets, true)
					}
				}
				if len(validatedSecrets) != 2 {
					ts.Errorf("unexpected secrets were validate, expected 2, got %d", len(validatedSecrets))
				}
			},
		},
		{
			name:         "test case 3 - no-op the secret was created by CAO in 4.8 OR this is downgrade",
			fakeDeployer: newFakeDeployer(true, nil),
			initialSecretsGlobal: []*corev1.Secret{
				func() *corev1.Secret {
					s := defaultSecret(fmt.Sprintf("%s-%s", encryptionconfig.EncryptionConfSecretName, oAuthAPIServerTargetNamespace), operatorclient.GlobalMachineSpecifiedConfigNamespace)
					delete(s.Annotations, EncryptionConfigManagedBy)
					return s
				}(),
			},
			initialSecretsOAuthTarget: []*corev1.Secret{
				func() *corev1.Secret {
					s := defaultSecret("encryption-config", oAuthAPIServerTargetNamespace)
					delete(s.Annotations, EncryptionConfigManagedBy)
					return s
				}(),
			},
		},

		{
			name:         "no-op if the encryption is in progress",
			fakeDeployer: newFakeDeployer(false, nil),
			initialSecretsGlobal: []*corev1.Secret{
				defaultSecret(fmt.Sprintf("%s-%s", encryptionconfig.EncryptionConfSecretName, oAuthAPIServerTargetNamespace), operatorclient.GlobalMachineSpecifiedConfigNamespace),
			},
			initialSecretsOAuthTarget: []*corev1.Secret{
				defaultSecret("encryption-config", oAuthAPIServerTargetNamespace),
			},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			// test data
			eventRecorder := events.NewInMemoryRecorder("")
			syncContext := factory.NewSyncContext("", eventRecorder)
			fakeSecretsIndexerGlobal := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
			for _, secret := range scenario.initialSecretsGlobal {
				fakeSecretsIndexerGlobal.Add(secret)
			}
			fakeSecretsListerGlobal := corev1listers.NewSecretLister(fakeSecretsIndexerGlobal)

			fakeSecretsIndexerOAuthTarget := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
			for _, secret := range scenario.initialSecretsOAuthTarget {
				fakeSecretsIndexerOAuthTarget.Add(secret)
			}
			fakeSecretsListerOAuthTarget := corev1listers.NewSecretLister(fakeSecretsIndexerOAuthTarget)

			rawSecrets := []runtime.Object{}
			for _, secret := range scenario.initialSecretsGlobal {
				rawSecrets = append(rawSecrets, secret)
			}
			for _, secret := range scenario.initialSecretsOAuthTarget {
				rawSecrets = append(rawSecrets, secret)
			}
			fakeKubeClient := fake.NewSimpleClientset(rawSecrets...)

			target := oauthEncryptionConfigSyncController{
				oauthAPIServerTargetNamespace:       "openshift-oauth-apiserver",
				secretListerConfigManaged:           fakeSecretsListerGlobal.Secrets(operatorclient.GlobalMachineSpecifiedConfigNamespace),
				secretListerForOAuthTargetNamespace: fakeSecretsListerOAuthTarget.Secrets("openshift-oauth-apiserver"),
				secretClient:                        fakeKubeClient.CoreV1(),
				deployer:                            scenario.fakeDeployer,
			}

			// act
			err := target.sync(context.TODO(), syncContext)
			if err != nil {
				t.Fatal(err)
			}

			// validate
			if err := validateActionsVerbs(fakeKubeClient.Actions(), scenario.expectedActions); err != nil {
				t.Fatal(err)
			}

			if err := validateEventsReason(eventRecorder.Events(), scenario.expectedEvents); err != nil {
				t.Error(err)
			}
			if scenario.validateFunc != nil {
				scenario.validateFunc(t, fakeKubeClient.Actions())
			}
		})
	}
}

func defaultSecret(name, ns string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: ns,
			Annotations: map[string]string{
				EncryptionConfigManagedBy:                encryptionConfigManagedByValue,
				encryptionstate.KubernetesDescriptionKey: encryptionstate.KubernetesDescriptionScaryValue,
			},
			Finalizers: []string{encryptionsecret.EncryptionSecretFinalizer},
		},
		Data: map[string][]byte{"encryption-config": {0xFF}},
	}
}

func validateActionsVerbs(actualActions []clientgotesting.Action, expectedActions []string) error {
	if len(actualActions) != len(expectedActions) {
		return fmt.Errorf("expected to get %d actions but got %d\nexpected=%v \n got=%v", len(expectedActions), len(actualActions), expectedActions, actionStrings(actualActions))
	}
	for i, a := range actualActions {
		if got, expected := actionString(a), expectedActions[i]; got != expected {
			return fmt.Errorf("at %d got %s, expected %s", i, got, expected)
		}
	}
	return nil
}

func actionString(a clientgotesting.Action) string {
	involvedObjectName := ""
	if updateAction, isUpdateAction := a.(clientgotesting.UpdateAction); isUpdateAction {
		rawObj := updateAction.GetObject()
		if objMeta, err := meta.Accessor(rawObj); err == nil {
			involvedObjectName = objMeta.GetName()
		}
	}
	if getAction, isGetAction := a.(clientgotesting.GetAction); isGetAction {
		involvedObjectName = getAction.GetName()
	}
	action := a.GetVerb() + ":" + a.GetResource().Resource
	if len(a.GetNamespace()) > 0 {
		action = action + ":" + a.GetNamespace()
	}
	if len(involvedObjectName) > 0 {
		action = action + ":" + involvedObjectName
	}
	return action
}

func actionStrings(actions []clientgotesting.Action) []string {
	res := make([]string, 0, len(actions))
	for _, a := range actions {
		res = append(res, actionString(a))
	}
	return res
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

type fakeDeployer struct {
	converged bool
	err       error
}

func newFakeDeployer(converged bool, err error) *fakeDeployer {
	return &fakeDeployer{converged: converged, err: err}
}

func (d *fakeDeployer) DeployedEncryptionConfigSecret() (secret *corev1.Secret, converged bool, err error) {
	return nil, d.converged, d.err
}

func (d *fakeDeployer) AddEventHandler(handler cache.ResourceEventHandler) {}

func (d *fakeDeployer) HasSynced() bool {
	return true
}
