package revisionpoddeployer

import (
	"context"
	"fmt"
	"reflect"

	"github.com/openshift/library-go/pkg/operator/encryption/encryptionconfig"
	"github.com/openshift/library-go/pkg/operator/encryption/statemachine"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

type MaybeDisabledDeployer interface {
	statemachine.Deployer
	Disabled() bool
}

// UnionDeployer provides unified state from multiple distinct deployers.
// It's converged if all delegates are converged. Each deployer can be disabled.
type UnionDeployer struct {
	delegates []MaybeDisabledDeployer
	hasSynced []cache.InformerSynced
}

var _ statemachine.Deployer = &UnionDeployer{}

// NewUnionDeployer creates a deployer that returns a unified state from multiple distinct deployers.
// That means:
//  - none has reported an error
//  - all have converged
//  - all have observed exactly the same encryption configuration otherwise it returns converged=false
func NewUnionDeployer(delegates ...MaybeDisabledDeployer) (*UnionDeployer, error) {
	if len(delegates) == 0 {
		return nil, fmt.Errorf("no deployers were configured")
	}
	return &UnionDeployer{delegates: delegates}, nil
}

// DeployedEncryptionConfigSecret returns the actual encryption configuration across multiple deployers if they all agree.
func (d *UnionDeployer) DeployedEncryptionConfigSecret(ctx context.Context) (secret *corev1.Secret, converged bool, err error) {
	seenSecrets := []*corev1.Secret{}

	for _, delegate := range d.delegates {
		if delegate.Disabled() {
			continue
		}
		secret, converged, err := delegate.DeployedEncryptionConfigSecret(ctx)
		if !converged || err != nil {
			return nil, converged, err
		}

		seenSecrets = append(seenSecrets, secret)
	}

	// at this point an empty secret (nil) can actually mean two things
	// 1. encryption is off
	// 2. a replica hasn't converged, it's been stuck, it's been slow etc
	//
	// in either of this case we should report !converged and let the encryption state machine handle it
	potentiallyConvergedSecrets := []*corev1.Secret{}
	for _, secret := range seenSecrets {
		if secret != nil {
			potentiallyConvergedSecrets = append(potentiallyConvergedSecrets, secret)
		}
	}

	if len(potentiallyConvergedSecrets) == 0 {
		return nil, true, nil // encryption is off
	}
	if len(potentiallyConvergedSecrets) != len(seenSecrets) {
		return nil, false, nil // not all replicas have converged
	}

	// we need to check that the encryption configuration is exactly the same among deployers
	// so we promote the fist secret and compare it with the rest
	goldenSecret := potentiallyConvergedSecrets[0]
	potentiallyConvergedSecrets = potentiallyConvergedSecrets[1:]

	goldenEncryptionCfg, err := encryptionconfig.FromSecret(goldenSecret)
	if err != nil {
		return nil, false, err
	}

	for _, secret := range potentiallyConvergedSecrets {
		currentEncryptionCfg, err := encryptionconfig.FromSecret(secret)
		if err != nil {
			return nil, false, err
		}

		if !reflect.DeepEqual(goldenEncryptionCfg.Resources, currentEncryptionCfg.Resources) {
			return nil, false, nil
		}
	}

	return goldenSecret, true, nil
}

func (d *UnionDeployer) HasSynced() bool {
	for _, hasSynced := range d.hasSynced {
		if !hasSynced() {
			return false
		}
	}
	return true
}

// AddEventHandler registers a event handler that might influence the result of DeployedEncryptionConfigSecret for all configured deployers.
func (d *UnionDeployer) AddEventHandler(handler cache.ResourceEventHandler) {
	d.hasSynced = []cache.InformerSynced{}
	for _, delegate := range d.delegates {
		delegate.AddEventHandler(handler)
		d.hasSynced = append(d.hasSynced, delegate.HasSynced)
	}
}

type disabledByPredicateDeployer struct {
	statemachine.Deployer
	enabled func() bool
}

var _ MaybeDisabledDeployer = &disabledByPredicateDeployer{}

// NewDisabledByPredicateDeployer returns a deployer used by the encryption controllers.
// Whether this deployer is on/off is determined by enabled
// TODO: remove this deployer in 4.8
func NewDisabledByPredicateDeployer(
	enabled func() bool,
	delegate statemachine.Deployer) *disabledByPredicateDeployer {
	return &disabledByPredicateDeployer{
		Deployer: delegate,
		enabled:  enabled,
	}
}

// Disabled indicates whether this deployer is disabled
// Note: OAuthAPIServer deployer changes its availability - see enabled implementation
func (d *disabledByPredicateDeployer) Disabled() bool {
	return !d.enabled()
}

type AlwaysEnabledDeployer struct {
	statemachine.Deployer
}

var _ MaybeDisabledDeployer = &AlwaysEnabledDeployer{}

// Disabled indicates whether this deployer is disabled
// Note: OpenShift deployer is always enabled
func (d *AlwaysEnabledDeployer) Disabled() bool {
	return false
}
