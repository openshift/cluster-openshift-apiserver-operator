package workload

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"
	"time"

	configv1 "github.com/openshift/api/config/v1"
	openshiftcontrolplanev1 "github.com/openshift/api/openshiftcontrolplane/v1"
	operatorv1 "github.com/openshift/api/operator/v1"
	configfake "github.com/openshift/client-go/config/clientset/versioned/fake"
	configinformers "github.com/openshift/client-go/config/informers/externalversions"
	operatorfake "github.com/openshift/client-go/operator/clientset/versioned/fake"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"
	"github.com/openshift/library-go/pkg/controller/factory"
	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/status"
	operatorv1helpers "github.com/openshift/library-go/pkg/operator/v1helpers"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/diff"
	"k8s.io/client-go/kubernetes/fake"
)

func fakeCountNodes(_ map[string]string) (*int32, error) {
	masterNodeCount := int32(3)
	return &masterNodeCount, nil
}

func TestOperatorConfigProgressingCondition(t *testing.T) {
	testCases := []struct {
		name                             string
		operatorConfigGeneration         int64
		operatorConfigObservedGeneration int64
		deploymentGeneration             int64
		expectedStatus                   operatorv1.ConditionStatus
		expectedMessage                  string
	}{
		{
			name:                             "HappyPath",
			operatorConfigGeneration:         100,
			operatorConfigObservedGeneration: 100,
			deploymentGeneration:             200,
			expectedStatus:                   operatorv1.ConditionFalse,
		},
		{
			name:                             "ConfigObservedAhead",
			operatorConfigGeneration:         100,
			operatorConfigObservedGeneration: 101,
			deploymentGeneration:             200,
			expectedStatus:                   operatorv1.ConditionTrue,
			expectedMessage:                  "openshiftapiserveroperatorconfigs/instance: observed generation is 101, desired generation is 100.",
		},
		{
			name:                             "DeploymentObservedBehind",
			operatorConfigGeneration:         101,
			operatorConfigObservedGeneration: 100,
			deploymentGeneration:             201,
			expectedStatus:                   operatorv1.ConditionTrue,
			expectedMessage:                  "openshiftapiserveroperatorconfigs/instance: observed generation is 100, desired generation is 101.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			kubeClient := fake.NewSimpleClientset(
				&corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "serving-cert", Namespace: "openshift-apiserver"}},
				&corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "etcd-client", Namespace: operatorclient.GlobalUserSpecifiedConfigNamespace}},
				&appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Name:       "apiserver",
						Namespace:  "openshift-apiserver",
						Generation: tc.deploymentGeneration,
					},
					Status: appsv1.DeploymentStatus{
						AvailableReplicas: 100,
					},
				})

			operatorConfig := &operatorv1.OpenShiftAPIServer{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "cluster",
					Generation: tc.operatorConfigGeneration,
				},
				Spec: operatorv1.OpenShiftAPIServerSpec{
					OperatorSpec: operatorv1.OperatorSpec{},
				},
				Status: operatorv1.OpenShiftAPIServerStatus{
					OperatorStatus: operatorv1.OperatorStatus{
						ObservedGeneration: tc.operatorConfigObservedGeneration,
					},
				},
			}
			apiServiceOperatorClient := operatorfake.NewSimpleClientset(operatorConfig)
			openshiftConfigClient := configfake.NewSimpleClientset(
				&configv1.Image{ObjectMeta: metav1.ObjectMeta{Name: "cluster"}},
				&configv1.ClusterVersion{ObjectMeta: metav1.ObjectMeta{Name: "version"}},
			)

			fakeOperatorClient := operatorv1helpers.NewFakeOperatorClient(&operatorv1.OperatorSpec{ManagementState: operatorv1.Managed}, &operatorv1.OperatorStatus{}, nil)
			configInformers := configinformers.NewSharedInformerFactory(openshiftConfigClient, 10*time.Minute)

			target := OpenShiftAPIServerWorkload{
				kubeClient:                kubeClient,
				operatorClient:            fakeOperatorClient,
				operatorConfigClient:      apiServiceOperatorClient.OperatorV1(),
				openshiftConfigClient:     openshiftConfigClient.ConfigV1(),
				clusterVersionInformer:    configInformers.Config().V1().ClusterVersions(),
				versionRecorder:           status.NewVersionGetter(),
				countNodes:                fakeCountNodes,
				ensureAtMostOnePodPerNode: func(spec *appsv1.DeploymentSpec, componentName string) error { return nil },
			}

			// register the informer
			configInformers.Config().V1().ClusterVersions().Informer()

			ctx := context.Background()
			configInformers.Start(ctx.Done())
			configInformers.WaitForCacheSync(ctx.Done())

			if _, _, err := target.Sync(ctx, factory.NewSyncContext("TestSyncCOntext", events.NewInMemoryRecorder(""))); len(err) > 0 {
				t.Fatal(err)
			}

			_, resultStatus, _, err := fakeOperatorClient.GetOperatorState()
			if err != nil {
				t.Fatal(err)
			}
			if resultStatus.ObservedGeneration != tc.operatorConfigGeneration {
				t.Fatalf("expected operator.ObservedGeneration of %d, bot got %d", tc.operatorConfigGeneration, resultStatus.ObservedGeneration)
			}

			condition := operatorv1helpers.FindOperatorCondition(resultStatus.Conditions, "OperatorConfigProgressing")
			if condition == nil {
				t.Fatalf("No %v condition found.", operatorv1.OperatorStatusTypeProgressing)
			}
			if condition.Status != tc.expectedStatus {
				t.Errorf("expected status == %v, actual status == %v", tc.expectedStatus, condition.Status)
			}
			if condition.Message != tc.expectedMessage {
				t.Errorf("expected message:\n%v\nactual message:\n%v", tc.expectedMessage, condition.Message)
			}
			actualGenerationStatus := resultStatus.Generations[0]
			expectedGenerationStatus := operatorv1.GenerationStatus{
				Group:          "apps",
				Resource:       "deployments",
				Namespace:      "openshift-apiserver",
				Name:           "apiserver",
				LastGeneration: tc.deploymentGeneration,
			}
			if !equality.Semantic.DeepEqual(actualGenerationStatus, expectedGenerationStatus) {
				t.Errorf("generation status mismatch, diff = %s", diff.ObjectDiff(actualGenerationStatus, expectedGenerationStatus))
			}
		})
	}
}

var withETCDServerListJSON = `
{
  "storageConfig": {
    "urls": [
      "https://10.0.131.191:2379",
      "https://10.0.159.206:2379"
    ]
  }
}
`

var withSomeDataJSON = `
{
  "cloudProviderFile": "path"
}
`

func TestPreconditionFulfilled(t *testing.T) {
	scenarios := []struct {
		name            string
		operator        *operatorv1.OpenShiftAPIServer
		expectError     bool
		preconditionMet bool
	}{
		// scenario 1
		{
			name: "mandatory etcd-servers are specified",
			operator: &operatorv1.OpenShiftAPIServer{
				Spec: operatorv1.OpenShiftAPIServerSpec{OperatorSpec: operatorv1.OperatorSpec{
					ObservedConfig: runtime.RawExtension{Raw: []byte(withETCDServerListJSON)},
				}},
			},
			preconditionMet: true,
		},

		// scenario 2
		{
			name: "no etcd-servers were specified",
			operator: &operatorv1.OpenShiftAPIServer{
				Spec: operatorv1.OpenShiftAPIServerSpec{OperatorSpec: operatorv1.OperatorSpec{
					ObservedConfig: runtime.RawExtension{Raw: []byte(withSomeDataJSON)},
				}},
			},
			expectError: true,
		},

		// scenario 3
		{
			name:        "no cfg",
			operator:    &operatorv1.OpenShiftAPIServer{},
			expectError: true,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			// test data
			target := &OpenShiftAPIServerWorkload{}

			// act
			actualPreconditions, err := target.preconditionFulfilledInternal(scenario.operator)

			// validate
			if err != nil && !scenario.expectError {
				t.Fatalf("unexpected error returned %v", err)
			}
			if err == nil && scenario.expectError {
				t.Fatal("expected an error")
			}
			if scenario.preconditionMet != actualPreconditions {
				t.Fatalf("unexpected precondtions = %v, expected = %v", actualPreconditions, scenario.preconditionMet)
			}
		})
	}
}

func TestCapabilities(t *testing.T) {
	testCases := []struct {
		name                    string
		knownCapabilities       []configv1.ClusterVersionCapability
		enabledCapabilities     []configv1.ClusterVersionCapability
		expectedPerGroupOptions []openshiftcontrolplanev1.PerGroupOptions
	}{
		{
			name:                    "Build and DC capabilities enabled implicitly",
			knownCapabilities:       []configv1.ClusterVersionCapability{},
			enabledCapabilities:     []configv1.ClusterVersionCapability{},
			expectedPerGroupOptions: []openshiftcontrolplanev1.PerGroupOptions{},
		},
		{
			name: "Build and DC capabilities enabled explicitly",
			knownCapabilities: []configv1.ClusterVersionCapability{
				configv1.ClusterVersionCapabilityBuild,
				configv1.ClusterVersionCapabilityDeploymentConfig,
			},
			enabledCapabilities: []configv1.ClusterVersionCapability{
				configv1.ClusterVersionCapabilityBuild,
				configv1.ClusterVersionCapabilityDeploymentConfig,
			},
			expectedPerGroupOptions: []openshiftcontrolplanev1.PerGroupOptions{},
		},
		{
			name: "Build capability disabled",
			knownCapabilities: []configv1.ClusterVersionCapability{
				configv1.ClusterVersionCapabilityBuild,
				configv1.ClusterVersionCapabilityDeploymentConfig,
			},
			enabledCapabilities: []configv1.ClusterVersionCapability{
				configv1.ClusterVersionCapabilityDeploymentConfig,
			},
			expectedPerGroupOptions: []openshiftcontrolplanev1.PerGroupOptions{
				{Name: openshiftcontrolplanev1.OpenShiftBuildAPIserver, DisabledVersions: []string{"v1"}},
			},
		},
		{
			name: "DC capability disabled",
			knownCapabilities: []configv1.ClusterVersionCapability{
				configv1.ClusterVersionCapabilityBuild,
				configv1.ClusterVersionCapabilityDeploymentConfig,
			},
			enabledCapabilities: []configv1.ClusterVersionCapability{
				configv1.ClusterVersionCapabilityBuild,
			},
			expectedPerGroupOptions: []openshiftcontrolplanev1.PerGroupOptions{
				{Name: openshiftcontrolplanev1.OpenShiftAppsAPIserver, DisabledVersions: []string{"v1"}},
			},
		},
		{
			name: "Build and DC capabilities disabled",
			knownCapabilities: []configv1.ClusterVersionCapability{
				configv1.ClusterVersionCapabilityBuild,
				configv1.ClusterVersionCapabilityDeploymentConfig,
			},
			enabledCapabilities: []configv1.ClusterVersionCapability{},
			expectedPerGroupOptions: []openshiftcontrolplanev1.PerGroupOptions{
				{Name: openshiftcontrolplanev1.OpenShiftBuildAPIserver, DisabledVersions: []string{"v1"}},
				{Name: openshiftcontrolplanev1.OpenShiftAppsAPIserver, DisabledVersions: []string{"v1"}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			kubeClient := fake.NewSimpleClientset(
				&corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "serving-cert", Namespace: "openshift-apiserver"}},
				&corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "etcd-client", Namespace: operatorclient.GlobalUserSpecifiedConfigNamespace}},
				&appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Name:       "apiserver",
						Namespace:  "openshift-apiserver",
						Generation: 100,
					},
					Status: appsv1.DeploymentStatus{
						AvailableReplicas: 100,
					},
				})

			operatorConfig := &operatorv1.OpenShiftAPIServer{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "cluster",
					Generation: 100,
				},
				Spec: operatorv1.OpenShiftAPIServerSpec{
					OperatorSpec: operatorv1.OperatorSpec{},
				},
				Status: operatorv1.OpenShiftAPIServerStatus{
					OperatorStatus: operatorv1.OperatorStatus{
						ObservedGeneration: 100,
					},
				},
			}
			apiServiceOperatorClient := operatorfake.NewSimpleClientset(operatorConfig)
			openshiftConfigClient := configfake.NewSimpleClientset(
				&configv1.Image{ObjectMeta: metav1.ObjectMeta{Name: "cluster"}},
				&configv1.ClusterVersion{
					ObjectMeta: metav1.ObjectMeta{Name: "version"},
					Status: configv1.ClusterVersionStatus{
						Capabilities: configv1.ClusterVersionCapabilitiesStatus{
							EnabledCapabilities: tc.enabledCapabilities,
							KnownCapabilities:   tc.knownCapabilities,
						},
					},
				},
			)

			fakeOperatorClient := operatorv1helpers.NewFakeOperatorClient(&operatorv1.OperatorSpec{ManagementState: operatorv1.Managed}, &operatorv1.OperatorStatus{}, nil)
			configInformers := configinformers.NewSharedInformerFactory(openshiftConfigClient, 10*time.Minute)

			target := OpenShiftAPIServerWorkload{
				kubeClient:                kubeClient,
				operatorClient:            fakeOperatorClient,
				operatorConfigClient:      apiServiceOperatorClient.OperatorV1(),
				openshiftConfigClient:     openshiftConfigClient.ConfigV1(),
				clusterVersionInformer:    configInformers.Config().V1().ClusterVersions(),
				versionRecorder:           status.NewVersionGetter(),
				countNodes:                fakeCountNodes,
				ensureAtMostOnePodPerNode: func(spec *appsv1.DeploymentSpec, componentName string) error { return nil },
			}

			// register the informer
			configInformers.Config().V1().ClusterVersions().Informer()

			ctx := context.Background()
			configInformers.Start(ctx.Done())
			configInformers.WaitForCacheSync(ctx.Done())

			if _, _, err := target.Sync(ctx, factory.NewSyncContext("TestSyncCOntext", events.NewInMemoryRecorder(""))); len(err) > 0 {
				t.Fatal(err)
			}

			cm, err := kubeClient.CoreV1().ConfigMaps("openshift-apiserver").Get(ctx, "config", metav1.GetOptions{})
			if err != nil {
				t.Fatalf("Unable to get 'config' configmap: %v", err)
			}

			config := openshiftcontrolplanev1.OpenShiftAPIServerConfig{}
			if err := json.NewDecoder(bytes.NewBuffer([]byte(cm.Data["config.yaml"]))).Decode(&config); err != nil {
				t.Fatalf("Unable to decode OpenShiftAPIServerConfig: %v", err)
			}

			if !equality.Semantic.DeepEqual(config.APIServers.PerGroupOptions, tc.expectedPerGroupOptions) {
				t.Errorf("generation status mismatch, diff = %s", diff.ObjectDiff(config.APIServers.PerGroupOptions, tc.expectedPerGroupOptions))
			}
		})
	}
}
