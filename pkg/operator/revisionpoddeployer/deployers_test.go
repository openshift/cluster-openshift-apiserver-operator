package revisionpoddeployer_test

import (
	"context"
	"fmt"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/util/diff"
	apiserverv1 "k8s.io/apiserver/pkg/apis/apiserver/v1"
	"k8s.io/client-go/tools/cache"

	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/revisionpoddeployer"
	"github.com/openshift/library-go/pkg/operator/encryption/encryptionconfig"
	encryptiontesting "github.com/openshift/library-go/pkg/operator/encryption/testing"
)

func TestUnionRevisionLabelPodDeployer(t *testing.T) {
	scenarios := []struct {
		name      string
		deployers []revisionpoddeployer.MaybeDisabledDeployer

		expectedSecret            *corev1.Secret
		expectedConverged         bool
		expectedErr               bool
		expectedConstructionError bool
	}{
		{
			name: "happy path",
			deployers: []revisionpoddeployer.MaybeDisabledDeployer{
				newFakeDeployer(createDefaultSecretWithEncryptionConfig(t), true, false, nil),
				newFakeDeployer(createDefaultSecretWithEncryptionConfig(t), true, false, nil),
			},
			expectedSecret:    createDefaultSecretWithEncryptionConfig(t),
			expectedConverged: true,
		},
		{
			name: "encryption config mismatch",
			deployers: []revisionpoddeployer.MaybeDisabledDeployer{
				newFakeDeployer(createDefaultSecretWithEncryptionConfig(t), true, false, nil),
				newFakeDeployer(func() *corev1.Secret {
					ec := createDefaultEncryptionConfig()
					ec.Resources = append(ec.Resources, apiserverv1.ResourceConfiguration{Resources: []string{"pods"}})
					return encryptionCfgToSecret(t, ec)
				}(), true, false, nil),
			},
			expectedSecret:    nil,
			expectedConverged: false,
		},
		{
			name: "deployer2 hasn't converged",
			deployers: []revisionpoddeployer.MaybeDisabledDeployer{
				newFakeDeployer(createDefaultSecretWithEncryptionConfig(t), true, false, nil),
				newFakeDeployer(createDefaultSecretWithEncryptionConfig(t), false, false, nil),
			},
			expectedConverged: false,
		},
		{
			name: "deployer1 reported an error",
			deployers: []revisionpoddeployer.MaybeDisabledDeployer{
				newFakeDeployer(createDefaultSecretWithEncryptionConfig(t), false, false, fmt.Errorf("nasty error")),
				newFakeDeployer(createDefaultSecretWithEncryptionConfig(t), true, false, nil),
			},
			expectedConverged: false,
			expectedErr:       true,
		},
		{
			name: "happy path with a single deployer",
			deployers: []revisionpoddeployer.MaybeDisabledDeployer{
				newFakeDeployer(createDefaultSecretWithEncryptionConfig(t), true, false, nil),
			},
			expectedConverged: true,
			expectedSecret:    createDefaultSecretWithEncryptionConfig(t),
		},
		{
			name:                      "no-op when no deployers",
			deployers:                 []revisionpoddeployer.MaybeDisabledDeployer{},
			expectedConstructionError: true,
		},
		{
			name: "encryption off",
			deployers: []revisionpoddeployer.MaybeDisabledDeployer{
				newFakeDeployer(nil, true, false, nil),
				newFakeDeployer(nil, true, false, nil),
			},
			expectedConverged: true,
		},
		{
			name: "deployer2 hasn't converged - nil secret",
			deployers: []revisionpoddeployer.MaybeDisabledDeployer{
				newFakeDeployer(createDefaultSecretWithEncryptionConfig(t), true, false, nil),
				newFakeDeployer(nil, true, false, nil),
			},
			expectedConverged: false,
		},
		{
			name: "deployer2 is disabled",
			deployers: []revisionpoddeployer.MaybeDisabledDeployer{
				newFakeDeployer(createDefaultSecretWithEncryptionConfig(t), true, false, nil),
				newFakeDeployer(nil, false, true, nil),
			},
			expectedConverged: true,
			expectedSecret:    createDefaultSecretWithEncryptionConfig(t),
		},
		{
			name: "all disabled",
			deployers: []revisionpoddeployer.MaybeDisabledDeployer{
				newFakeDeployer(nil, false, true, nil),
				newFakeDeployer(nil, false, true, nil),
			},
			expectedConverged: true, // encryption is off
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			// test data
			checkError := func(err error, expectedErr bool) {
				if err != nil && !expectedErr {
					t.Errorf("got unexpected error %v", err)
				}
				if err == nil && expectedErr {
					t.Error("expected an error but didn't get one")
				}
			}

			target, err := revisionpoddeployer.NewUnionDeployer(scenario.deployers...)
			checkError(err, scenario.expectedConstructionError)
			if err != nil {
				return
			}

			// act
			actualSecret, actualConverged, actualErr := target.DeployedEncryptionConfigSecret(context.TODO())

			// validate
			checkError(actualErr, scenario.expectedErr)
			if scenario.expectedConverged != actualConverged {
				t.Errorf("expected converged to be %v, got %v", scenario.expectedConverged, actualConverged)
			}
			if !equality.Semantic.DeepEqual(actualSecret, scenario.expectedSecret) {
				t.Error(fmt.Errorf("retruned secret mismatch, diff = %s", diff.ObjectDiff(actualSecret, scenario.expectedSecret)))
			}
		})
	}
}

func createDefaultSecretWithEncryptionConfig(t *testing.T) *corev1.Secret {
	ec := createDefaultEncryptionConfig()
	return encryptionCfgToSecret(t, ec)
}

func encryptionCfgToSecret(t *testing.T, ec *apiserverv1.EncryptionConfiguration) *corev1.Secret {
	s, err := encryptionconfig.ToSecret("targetNs", fmt.Sprintf("%s-%s", "encryption-config", "1"), ec)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func createDefaultEncryptionConfig() *apiserverv1.EncryptionConfiguration {
	keysResForSecrets := encryptiontesting.EncryptionKeysResourceTuple{
		Resource: "secrets",
		Keys: []apiserverv1.Key{
			{
				Name:   "1",
				Secret: "NzFlYTdjOTE0MTlhNjhmZDEyMjRmODhkNTAzMTZiNGU=",
			},
		},
	}
	keysResForConfigMaps := encryptiontesting.EncryptionKeysResourceTuple{
		Resource: "configmaps",
		Keys: []apiserverv1.Key{
			{
				Name:   "1",
				Secret: "NzFlYTdjOTE0MTlhNjhmZDEyMjRmODhkNTAzMTZiNGU=",
			},
		},
	}

	return encryptiontesting.CreateEncryptionCfgWithWriteKey([]encryptiontesting.EncryptionKeysResourceTuple{keysResForConfigMaps, keysResForSecrets})
}

type fakeDeployer struct {
	secret    *corev1.Secret
	converged bool
	disabled  bool
	err       error
}

func newFakeDeployer(secret *corev1.Secret, converged bool, disabled bool, err error) *fakeDeployer {
	return &fakeDeployer{secret: secret, converged: converged, disabled: disabled, err: err}
}

func (d *fakeDeployer) DeployedEncryptionConfigSecret(context.Context) (secret *corev1.Secret, converged bool, err error) {
	return d.secret, d.converged, d.err
}

func (d *fakeDeployer) AddEventHandler(handler cache.ResourceEventHandler) (cache.ResourceEventHandlerRegistration, error) {
	panic("implement me")
}

func (d *fakeDeployer) HasSynced() bool {
	return true
}

func (d *fakeDeployer) Disabled() bool {
	return d.disabled
}
