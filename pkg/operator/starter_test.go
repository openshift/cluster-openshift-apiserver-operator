package operator

import (
	"testing"
	"time"

	operatorv1 "github.com/openshift/api/operator/v1"
)

func TestNewAvailableInertia(t *testing.T) {
	inertia := newAvailableInertia()

	tests := []struct {
		conditionType string
		expected      time.Duration
	}{
		{
			conditionType: "APIServicesAvailable",
			expected:      10 * time.Second,
		},
		{
			conditionType: "APIServerDeploymentAvailable",
			expected:      10 * time.Second,
		},
		{
			// Any other Available sub-condition must not be affected.
			conditionType: "SomeOtherAvailable",
			expected:      0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.conditionType, func(t *testing.T) {
			got := inertia(operatorv1.OperatorCondition{Type: tt.conditionType})
			if got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}
