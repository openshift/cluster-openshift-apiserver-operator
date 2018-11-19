package v1alpha1

import (
	operatorsv1 "github.com/openshift/api/operator/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OpenShiftAPIServerConfig provides information to configure openshift-apiserver
type OpenShiftAPIServerConfig struct {
	metav1.TypeMeta `json:",inline"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OpenShiftAPIServerOperatorConfig provides information to configure an operator to manage openshift-apiserver.
type OpenShiftAPIServerOperatorConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   OpenShiftAPIServerOperatorConfigSpec   `json:"spec"`
	Status OpenShiftAPIServerOperatorConfigStatus `json:"status"`
}

type OpenShiftAPIServerOperatorConfigSpec struct {
	operatorsv1.OperatorSpec `json:",inline"`
}

type OpenShiftAPIServerOperatorConfigStatus struct {
	operatorsv1.OperatorStatus `json:",inline"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OpenShiftAPIServerOperatorConfigList is a collection of items
type OpenShiftAPIServerOperatorConfigList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items contains the items
	Items []OpenShiftAPIServerOperatorConfig `json:"items"`
}
