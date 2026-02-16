package library

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"

	configclient "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"
)

// WaitForClusterOperatorAvailableNotProgressingNotDegraded waits for a ClusterOperator to report
// Available=True, Progressing=False, and Degraded=False status conditions.
func WaitForClusterOperatorAvailableNotProgressingNotDegraded(ctx context.Context, client configclient.ConfigV1Interface, name string) error {
	return wait.PollUntilContextTimeout(ctx, 5*time.Second, 5*time.Minute, true, func(ctx context.Context) (bool, error) {
		co, err := client.ClusterOperators().Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		available := false
		progressing := true
		degraded := true
		for _, condition := range co.Status.Conditions {
			if condition.Type == "Available" && condition.Status == "True" {
				available = true
			}
			if condition.Type == "Progressing" && condition.Status == "False" {
				progressing = false
			}
			if condition.Type == "Degraded" && condition.Status == "False" {
				degraded = false
			}
		}
		return available && !progressing && !degraded, nil
	})
}
