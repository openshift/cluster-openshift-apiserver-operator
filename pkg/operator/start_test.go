package operator

import (
	"testing"

	"github.com/davecgh/go-spew/spew"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	configv1 "github.com/openshift/api/config/v1"
	fakeconfigclientset "github.com/openshift/client-go/config/clientset/versioned/fake"
)

func TestNothing(t *testing.T) {
}

func TestFakeClientsetBroken(t *testing.T) {
	// just instantiate the object
	originalObj := &configv1.ClusterOperator{ObjectMeta: metav1.ObjectMeta{Name: "cool-operator"}}
	originalObj.Status.Conditions = []configv1.ClusterOperatorStatusCondition{{Type: "my-condition"}}
	fakeclient := fakeconfigclientset.NewSimpleClientset(originalObj)
	mutatedObj := originalObj.DeepCopy()
	mutatedObj.Status.Conditions = []configv1.ClusterOperatorStatusCondition{{Type: "different-condition"}}
	if _, err := fakeclient.ConfigV1().ClusterOperators().UpdateStatus(mutatedObj); err != nil {
		t.Fatal(err)
	}

	finalObj, err := fakeclient.ConfigV1().ClusterOperators().Get("cool-operator", metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if finalObj.Status.Conditions[0].Type != "different-condition" {
		t.Fatal(spew.Sdump(finalObj))
	}
}
