package workloadcontroller

import (
	"sort"
	"strings"
	"testing"

	"github.com/openshift/library-go/pkg/operator/status"

	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"

	"github.com/pkg/errors"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/diff"
	"k8s.io/client-go/kubernetes/fake"
	kubetesting "k8s.io/client-go/testing"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	kubeaggregatorfake "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset/fake"

	operatorv1 "github.com/openshift/api/operator/v1"
	configfake "github.com/openshift/client-go/config/clientset/versioned/fake"
	operatorfake "github.com/openshift/client-go/operator/clientset/versioned/fake"
	"github.com/openshift/library-go/pkg/operator/events"
	operatorv1helpers "github.com/openshift/library-go/pkg/operator/v1helpers"
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
			kubeAggregatorClient := kubeaggregatorfake.NewSimpleClientset()

			operator := OpenShiftAPIServerOperator{
				kubeClient:              kubeClient,
				eventRecorder:           events.NewInMemoryRecorder(""),
				operatorConfigClient:    apiServiceOperatorClient.OperatorV1(),
				openshiftConfigClient:   openshiftConfigClient.ConfigV1(),
				apiregistrationv1Client: kubeAggregatorClient.ApiregistrationV1(),
				versionRecorder:         status.NewVersionGetter(),
			}

			_, _ = syncOpenShiftAPIServer_v410_00_to_latest(operator, operatorConfig)

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
		expectedReason          string
		expectedMessages        []string
		expectedFailingMessages []string
		apiServiceReactor       kubetesting.ReactionFunc
		daemonReactor           kubetesting.ReactionFunc
	}{
		{
			name:           "Default",
			expectedStatus: operatorv1.ConditionTrue,
		},
		{
			name:                    "APIServiceCreateFailure",
			expectedStatus:          operatorv1.ConditionFalse,
			expectedReason:          "NoRegisteredAPIServices",
			expectedMessages:        []string{"registered apiservices could not be retrieved"},
			expectedFailingMessages: []string{"\"apiservices\": TEST ERROR: fail to create apiservice"},

			apiServiceReactor: func(action kubetesting.Action) (handled bool, ret runtime.Object, err error) {
				if action.GetVerb() == "get" && action.(kubetesting.GetAction).GetName() == "v1.build.openshift.io" {
					return true, nil, apierrors.NewNotFound(apiregistrationv1.Resource("apiservices"), "v1.build.openshift.io")
				}
				if action.GetVerb() != "create" {
					return false, nil, nil
				}
				if action.(kubetesting.CreateAction).GetObject().(*apiregistrationv1.APIService).Name == "v1.build.openshift.io" {
					return true, nil, errors.New("TEST ERROR: fail to create apiservice")
				}
				return false, nil, nil
			},
		},
		{
			name:                    "APIServiceGetFailure",
			expectedStatus:          operatorv1.ConditionFalse,
			expectedReason:          "NoRegisteredAPIServices",
			expectedMessages:        []string{"registered apiservices could not be retrieved"},
			expectedFailingMessages: []string{"\"apiservices\": TEST ERROR: fail to get apiservice"},

			apiServiceReactor: func(action kubetesting.Action) (handled bool, ret runtime.Object, err error) {
				if action.GetVerb() == "get" && action.(kubetesting.GetAction).GetName() == "v1.build.openshift.io" {
					return true, nil, errors.New("TEST ERROR: fail to get apiservice")
				}
				return false, nil, nil
			},
		},
		{
			name:                    "DaemonSetGetFailure",
			expectedStatus:          operatorv1.ConditionFalse,
			expectedReason:          "NoDaemon",
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
			expectedReason:   "NoAPIServerPod",
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
		{
			name:             "APIServiceNotAvailable",
			expectedStatus:   operatorv1.ConditionFalse,
			expectedReason:   "APIServiceNotAvailable",
			expectedMessages: []string{"apiservice/v1.build.openshift.io: not available: TEST MESSAGE"},

			apiServiceReactor: func(action kubetesting.Action) (handled bool, ret runtime.Object, err error) {
				if action.GetVerb() == "get" && action.(kubetesting.GetAction).GetName() == "v1.build.openshift.io" {
					return true, &apiregistrationv1.APIService{
						ObjectMeta: metav1.ObjectMeta{Name: "v1.build.openshift.io", Annotations: map[string]string{"service.alpha.openshift.io/inject-cabundle": "true"}},
						Spec: apiregistrationv1.APIServiceSpec{
							Group:                "build.openshift.io",
							Version:              "v1",
							Service:              &apiregistrationv1.ServiceReference{Namespace: operatorclient.TargetNamespace, Name: "api"},
							GroupPriorityMinimum: 9900,
							VersionPriority:      15,
						},
						Status: apiregistrationv1.APIServiceStatus{
							Conditions: []apiregistrationv1.APIServiceCondition{
								{Type: apiregistrationv1.Available, Status: apiregistrationv1.ConditionFalse, Message: "TEST MESSAGE"},
							},
						},
					}, nil
				}
				return false, nil, nil
			},
		},
		{
			name:           "MultipleAPIServiceNotAvailable",
			expectedStatus: operatorv1.ConditionFalse,
			expectedReason: "Multiple",
			expectedMessages: []string{
				"apiservice/v1.build.openshift.io: not available: TEST MESSAGE",
				"apiservice/v1.project.openshift.io: not available: TEST MESSAGE",
			},

			apiServiceReactor: func(action kubetesting.Action) (handled bool, ret runtime.Object, err error) {
				if action.GetVerb() != "get" {
					return false, nil, nil
				}

				switch action.(kubetesting.GetAction).GetName() {
				case "v1.build.openshift.io":
					fallthrough
				case "v1.project.openshift.io":
					return true, &apiregistrationv1.APIService{
						ObjectMeta: metav1.ObjectMeta{Name: action.(kubetesting.GetAction).GetName(), Annotations: map[string]string{"service.alpha.openshift.io/inject-cabundle": "true"}},
						Spec: apiregistrationv1.APIServiceSpec{
							Group:                action.GetResource().Group,
							Version:              action.GetResource().Version,
							Service:              &apiregistrationv1.ServiceReference{Namespace: operatorclient.TargetNamespace, Name: "api"},
							GroupPriorityMinimum: 9900,
							VersionPriority:      15,
						},
						Status: apiregistrationv1.APIServiceStatus{
							Conditions: []apiregistrationv1.APIServiceCondition{
								{Type: apiregistrationv1.Available, Status: apiregistrationv1.ConditionFalse, Message: "TEST MESSAGE"},
							},
						},
					}, nil
				default:
					return false, nil, nil
				}
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
			kubeAggregatorClient := kubeaggregatorfake.NewSimpleClientset()
			kubeAggregatorClient.PrependReactor("get", "apiservices", func(action kubetesting.Action) (handled bool, ret runtime.Object, err error) {
				return true,
					&apiregistrationv1.APIService{
						ObjectMeta: metav1.ObjectMeta{Name: action.(kubetesting.GetAction).GetName(), Annotations: map[string]string{"service.alpha.openshift.io/inject-cabundle": "true"}},
						Spec:       apiregistrationv1.APIServiceSpec{Group: action.GetResource().Group, Version: action.GetResource().Version, Service: &apiregistrationv1.ServiceReference{Namespace: operatorclient.TargetNamespace, Name: "api"}, GroupPriorityMinimum: 9900, VersionPriority: 15},
						Status:     apiregistrationv1.APIServiceStatus{Conditions: []apiregistrationv1.APIServiceCondition{{Type: apiregistrationv1.Available, Status: apiregistrationv1.ConditionTrue}}},
					}, nil
			})

			if tc.daemonReactor != nil {
				kubeClient.PrependReactor("*", "daemonsets", tc.daemonReactor)
			}
			if tc.apiServiceReactor != nil {
				kubeAggregatorClient.PrependReactor("*", "apiservices", tc.apiServiceReactor)
			}

			operator := OpenShiftAPIServerOperator{
				kubeClient:              kubeClient,
				eventRecorder:           events.NewInMemoryRecorder(""),
				operatorConfigClient:    apiServiceOperatorClient.OperatorV1(),
				openshiftConfigClient:   openshiftConfigClient.ConfigV1(),
				apiregistrationv1Client: kubeAggregatorClient.ApiregistrationV1(),
				versionRecorder:         status.NewVersionGetter(),
			}

			_, _ = syncOpenShiftAPIServer_v410_00_to_latest(operator, operatorConfig)

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
			if tc.expectedReason != "" && condition.Reason != tc.expectedReason {
				t.Error(diff.ObjectGoPrintSideBySide(condition.Reason, tc.expectedReason))
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
