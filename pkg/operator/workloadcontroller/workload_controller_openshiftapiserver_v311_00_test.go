package workloadcontroller

import (
	"sort"
	"strings"
	"testing"

	"github.com/openshift/library-go/pkg/operator/status"

	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"

	"github.com/pkg/errors"

	operatorv1 "github.com/openshift/api/operator/v1"
	configfake "github.com/openshift/client-go/config/clientset/versioned/fake"
	operatorfake "github.com/openshift/client-go/operator/clientset/versioned/fake"
	"github.com/openshift/library-go/pkg/operator/events"
	operatorv1helpers "github.com/openshift/library-go/pkg/operator/v1helpers"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/diff"
	"k8s.io/client-go/kubernetes/fake"
	kubetesting "k8s.io/client-go/testing"
)

func TestProgressingCondition(t *testing.T) {

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
				&corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "etcd-client", Namespace: operatorclient.GlobalUserSpecifiedConfigNamespace}},
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

			operatorConfig := &operatorv1.OpenShiftAPIServer{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "cluster",
					Generation: tc.configGeneration,
				},
				Spec: operatorv1.OpenShiftAPIServerSpec{
					OperatorSpec: operatorv1.OperatorSpec{},
				},
				Status: operatorv1.OpenShiftAPIServerStatus{
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
				operatorConfigClient:  apiServiceOperatorClient.OperatorV1(),
				openshiftConfigClient: openshiftConfigClient.ConfigV1(),
				versionRecorder:       status.NewVersionGetter(),
			}

			_, _ = syncOpenShiftAPIServer_v311_00_to_latest(operator, operatorConfig)

			result, err := apiServiceOperatorClient.OperatorV1().OpenShiftAPIServers().Get("cluster", metav1.GetOptions{})
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

func TestAvailableStatus(t *testing.T) {
	testCases := []struct {
		name                    string
		expectedStatus          operatorv1.ConditionStatus
		expectedReasons         []string
		expectedMessages        []string
		expectedFailingMessages []string
		daemonReactor           kubetesting.ReactionFunc
	}{
		{
			name:           "Default",
			expectedStatus: operatorv1.ConditionTrue,
		},
		{
			name:                    "DaemonSetGetFailure",
			expectedStatus:          operatorv1.ConditionFalse,
			expectedReasons:         []string{"NoDaemon"},
			expectedMessages:        []string{"daemonset/apiserver.openshift-apiserver: could not be retrieved"},
			expectedFailingMessages: []string{"\"daemonsets\": TEST ERROR: fail to get daemonset/apiserver.openshift-apiserver"},

			daemonReactor: func(action kubetesting.Action) (handled bool, ret runtime.Object, err error) {
				if action.GetVerb() == "get" && action.GetNamespace() == "openshift-apiserver" && action.(kubetesting.GetAction).GetName() == "apiserver" {
					return true, nil, errors.New("TEST ERROR: fail to get daemonset/apiserver.openshift-apiserver")
				}
				return false, nil, nil
			},
		},
		{
			name:             "NoDaemonSetPods",
			expectedStatus:   operatorv1.ConditionFalse,
			expectedReasons:  []string{"NoAPIServerPod"},
			expectedMessages: []string{"no openshift-apiserver daemon pods available on any node."},

			daemonReactor: func(action kubetesting.Action) (handled bool, ret runtime.Object, err error) {
				if action.GetVerb() == "get" && action.GetNamespace() == "openshift-apiserver" && action.(kubetesting.GetAction).GetName() == "apiserver" {
					return true, &appsv1.DaemonSet{
						ObjectMeta: metav1.ObjectMeta{Name: "apiserver", Namespace: "openshift-apiserver"},
						Status:     appsv1.DaemonSetStatus{NumberAvailable: 0},
					}, nil
				}
				return false, nil, nil
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			kubeClient := fake.NewSimpleClientset(
				&corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "serving-cert", Namespace: "openshift-apiserver"}},
				&corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "etcd-client", Namespace: operatorclient.GlobalUserSpecifiedConfigNamespace}},
				&appsv1.DaemonSet{
					ObjectMeta: metav1.ObjectMeta{
						Name:       "apiserver",
						Namespace:  "openshift-apiserver",
						Generation: 99,
					},
					Status: appsv1.DaemonSetStatus{
						NumberAvailable:    100,
						ObservedGeneration: 99,
					},
				})

			operatorConfig := &operatorv1.OpenShiftAPIServer{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "cluster",
					Generation: 99,
				},
				Spec: operatorv1.OpenShiftAPIServerSpec{
					OperatorSpec: operatorv1.OperatorSpec{},
				},
				Status: operatorv1.OpenShiftAPIServerStatus{
					OperatorStatus: operatorv1.OperatorStatus{
						ObservedGeneration: 99,
					},
				},
			}
			apiServiceOperatorClient := operatorfake.NewSimpleClientset(operatorConfig)
			openshiftConfigClient := configfake.NewSimpleClientset()

			if tc.daemonReactor != nil {
				kubeClient.PrependReactor("*", "daemonsets", tc.daemonReactor)
			}

			operator := OpenShiftAPIServerOperator{
				kubeClient:            kubeClient,
				eventRecorder:         events.NewInMemoryRecorder(""),
				operatorConfigClient:  apiServiceOperatorClient.OperatorV1(),
				openshiftConfigClient: openshiftConfigClient.ConfigV1(),
				versionRecorder:       status.NewVersionGetter(),
			}

			_, _ = syncOpenShiftAPIServer_v311_00_to_latest(operator, operatorConfig)

			result, err := apiServiceOperatorClient.OperatorV1().OpenShiftAPIServers().Get("cluster", metav1.GetOptions{})
			if err != nil {
				t.Fatal(err)
			}

			condition := operatorv1helpers.FindOperatorCondition(result.Status.Conditions, operatorv1.OperatorStatusTypeAvailable)
			if condition == nil {
				t.Fatal("Available condition not found")
			}
			if condition.Status != tc.expectedStatus {
				t.Error(diff.ObjectGoPrintSideBySide(condition.Status, tc.expectedStatus))
			}
			expectedReasons := strings.Join(tc.expectedReasons, "\n")
			if len(expectedReasons) > 0 && condition.Reason != expectedReasons {
				t.Error(diff.ObjectGoPrintSideBySide(condition.Reason, expectedReasons))
			}
			if len(tc.expectedMessages) > 0 {
				actualMessages := strings.Split(condition.Message, "\n")
				a := make([]string, len(tc.expectedMessages))
				b := make([]string, len(actualMessages))
				copy(a, tc.expectedMessages)
				copy(b, actualMessages)
				sort.Strings(a)
				sort.Strings(b)
				if !equality.Semantic.DeepEqual(a, b) {
					t.Error("\n" + diff.ObjectDiff(a, b))
				}
			}
			if len(tc.expectedFailingMessages) > 0 {
				failingCondition := operatorv1helpers.FindOperatorCondition(result.Status.Conditions, workloadDegradedCondition)
				for _, expected := range tc.expectedFailingMessages {
					if failingCondition == nil {
						t.Errorf("expected failing message not found: %q", expected)
						continue
					}
					found := false
					for _, actual := range strings.Split(failingCondition.Message, "\n") {
						if expected == actual {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("expected failing message not found: %q", expected)
					}
				}
			}
		})
	}

}
