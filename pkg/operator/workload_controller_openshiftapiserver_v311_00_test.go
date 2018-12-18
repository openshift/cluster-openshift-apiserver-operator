package operator

import (
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	operatorv1 "github.com/openshift/api/operator/v1"
	configfake "github.com/openshift/client-go/config/clientset/versioned/fake"
	openshiftapiserveroperator "github.com/openshift/cluster-openshift-apiserver-operator/pkg/apis/openshiftapiserver/v1alpha1"
	operatorfake "github.com/openshift/cluster-openshift-apiserver-operator/pkg/generated/clientset/versioned/fake"
	"github.com/openshift/library-go/pkg/operator/events"
	operatorv1helpers "github.com/openshift/library-go/pkg/operator/v1helpers"
)

func TestSyncOpenShiftAPIServer_v311_00_to_latestProgressingCondition(t *testing.T) {

	testCases := []struct {
		name                        string
		daemonSetGeneration         int64
		daemonSetObservedGeneration int64
		configGeneration            int64
		configObservedGeneration    int64
		expectedStatus              operatorv1.ConditionStatus
		expectedMessage             string
	}{
		{
			name:                        "HappyPath",
			daemonSetGeneration:         100,
			daemonSetObservedGeneration: 100,
			configGeneration:            100,
			configObservedGeneration:    100,
			expectedStatus:              operatorv1.ConditionFalse,
		},
		{
			name:                        "DaemonSetObservedAhead",
			daemonSetGeneration:         100,
			daemonSetObservedGeneration: 101,
			configGeneration:            100,
			configObservedGeneration:    100,
			expectedStatus:              operatorv1.ConditionTrue,
			expectedMessage:             "daemonset/apiserver.openshift-operator: observed generation is 101, desired generation is 100.",
		},
		{
			name:                        "DaemonSetObservedBehind",
			daemonSetGeneration:         101,
			daemonSetObservedGeneration: 100,
			configGeneration:            100,
			configObservedGeneration:    100,
			expectedStatus:              operatorv1.ConditionTrue,
			expectedMessage:             "daemonset/apiserver.openshift-operator: observed generation is 100, desired generation is 101.",
		},
		{
			name:                        "ConfigObservedAhead",
			daemonSetGeneration:         100,
			daemonSetObservedGeneration: 100,
			configGeneration:            100,
			configObservedGeneration:    101,
			expectedStatus:              operatorv1.ConditionTrue,
			expectedMessage:             "openshiftapiserveroperatorconfigs/instance: observed generation is 101, desired generation is 100.",
		},
		{
			name:                        "ConfigObservedBehind",
			daemonSetGeneration:         100,
			daemonSetObservedGeneration: 100,
			configGeneration:            101,
			configObservedGeneration:    100,
			expectedStatus:              operatorv1.ConditionTrue,
			expectedMessage:             "openshiftapiserveroperatorconfigs/instance: observed generation is 100, desired generation is 101.",
		},
		{
			name:                        "MultipleObservedAhead",
			daemonSetGeneration:         100,
			daemonSetObservedGeneration: 101,
			configGeneration:            100,
			configObservedGeneration:    101,
			expectedStatus:              operatorv1.ConditionTrue,
			expectedMessage:             "daemonset/apiserver.openshift-operator: observed generation is 101, desired generation is 100.\nopenshiftapiserveroperatorconfigs/instance: observed generation is 101, desired generation is 100.",
		},
		{
			name:                        "ConfigAndDaemonSetGenerationMismatch",
			daemonSetGeneration:         100,
			daemonSetObservedGeneration: 100,
			configGeneration:            101,
			configObservedGeneration:    101,
			expectedStatus:              operatorv1.ConditionFalse,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			kubeClient := fake.NewSimpleClientset(
				&corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "serving-cert", Namespace: "openshift-apiserver"}},
				&corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "etcd-client", Namespace: "kube-system"}},
				&appsv1.DaemonSet{
					ObjectMeta: metav1.ObjectMeta{
						Name:       "apiserver",
						Namespace:  "openshift-apiserver",
						Generation: tc.daemonSetGeneration,
					},
					Status: appsv1.DaemonSetStatus{
						NumberAvailable:    100,
						ObservedGeneration: tc.daemonSetObservedGeneration,
					},
				})

			operatorConfig := &openshiftapiserveroperator.OpenShiftAPIServerOperatorConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "instance",
					Generation: tc.configGeneration,
				},
				Spec: openshiftapiserveroperator.OpenShiftAPIServerOperatorConfigSpec{
					OperatorSpec: operatorv1.OperatorSpec{},
				},
				Status: openshiftapiserveroperator.OpenShiftAPIServerOperatorConfigStatus{
					OperatorStatus: operatorv1.OperatorStatus{
						ObservedGeneration: tc.configObservedGeneration,
					},
				},
			}
			apiServiceOperatorClient := operatorfake.NewSimpleClientset(operatorConfig)
			openshiftConfigClient := configfake.NewSimpleClientset()

			operator := OpenShiftAPIServerOperator{
				kubeClient:            kubeClient,
				eventRecorder:         events.NewInMemoryRecorder(""),
				operatorConfigClient:  apiServiceOperatorClient.Openshiftapiserver(),
				openshiftConfigClient: openshiftConfigClient.ConfigV1(),
			}

			syncOpenShiftAPIServer_v311_00_to_latest(operator, operatorConfig)

			result, err := apiServiceOperatorClient.Openshiftapiserver().OpenShiftAPIServerOperatorConfigs().Get("instance", metav1.GetOptions{})
			if err != nil {
				t.Fatal(err)
			}

			condition := operatorv1helpers.FindOperatorCondition(result.Status.Conditions, operatorv1.OperatorStatusTypeProgressing)
			if condition == nil {
				t.Fatalf("No %v condition found.", operatorv1.OperatorStatusTypeProgressing)
			}
			if condition.Status != tc.expectedStatus {
				t.Errorf("expected status == %v, actual status == %v", tc.expectedStatus, condition.Status)
			}
			if condition.Message != tc.expectedMessage {
				t.Errorf("expected message:\n%v\nactual message:\n%v", tc.expectedMessage, condition.Message)
			}

		})
	}

}
