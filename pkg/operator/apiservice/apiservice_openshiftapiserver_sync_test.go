package apiservice

import (
	"testing"

	operatorv1 "github.com/openshift/api/operator/v1"
	operatorlistersv1 "github.com/openshift/client-go/operator/listers/operator/v1"
	"github.com/openshift/library-go/pkg/operator/events"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/diff"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/tools/cache"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	apiregistrationv1lister "k8s.io/kube-aggregator/pkg/client/listers/apiregistration/v1"
)

func TestDiffAPIServices(t *testing.T) {
	testCases := []struct {
		name              string
		oldAPIServices    []*apiregistrationv1.APIService
		newAPIServices    []*apiregistrationv1.APIService
		resultList        []string
		resultListChanged bool
	}{
		// scenario 1
		{
			name: "oauth removed",
			oldAPIServices: []*apiregistrationv1.APIService{
				newAPIService("authorization.openshift.io", "v1"),
				newAPIService("build.openshift.io", "v1"),
				newAPIService("image.openshift.io", "v1"),
				newAPIService("oauth.openshift.io", "v1"),
				newAPIService("route.openshift.io", "v1"),
				newAPIService("template.openshift.io", "v1"),
				newAPIService("user.openshift.io", "v1"),
			},
			newAPIServices: []*apiregistrationv1.APIService{
				newAPIService("authorization.openshift.io", "v1"),
				newAPIService("build.openshift.io", "v1"),
				newAPIService("image.openshift.io", "v1"),
				newAPIService("route.openshift.io", "v1"),
				newAPIService("template.openshift.io", "v1"),
				newAPIService("user.openshift.io", "v1"),
			},
			resultList: []string{
				"v1.authorization.openshift.io",
				"v1.build.openshift.io",
				"v1.image.openshift.io",
				"v1.route.openshift.io",
				"v1.template.openshift.io",
				"v1.user.openshift.io",
			},
			resultListChanged: true,
		},
		// scenario 2
		{
			name: "oauth added",
			oldAPIServices: []*apiregistrationv1.APIService{
				newAPIService("authorization.openshift.io", "v1"),
				newAPIService("build.openshift.io", "v1"),
				newAPIService("image.openshift.io", "v1"),
				newAPIService("route.openshift.io", "v1"),
				newAPIService("template.openshift.io", "v1"),
				newAPIService("user.openshift.io", "v1"),
			},
			newAPIServices: []*apiregistrationv1.APIService{
				newAPIService("authorization.openshift.io", "v1"),
				newAPIService("build.openshift.io", "v1"),
				newAPIService("image.openshift.io", "v1"),
				newAPIService("oauth.openshift.io", "v1"),
				newAPIService("route.openshift.io", "v1"),
				newAPIService("template.openshift.io", "v1"),
				newAPIService("user.openshift.io", "v1"),
			},
			resultList: []string{
				"v1.authorization.openshift.io",
				"v1.build.openshift.io",
				"v1.image.openshift.io",
				"v1.oauth.openshift.io",
				"v1.route.openshift.io",
				"v1.template.openshift.io",
				"v1.user.openshift.io",
			},
			resultListChanged: true,
		},
		// scenario 3
		{
			name: "oauth added, user removed",
			oldAPIServices: []*apiregistrationv1.APIService{
				newAPIService("authorization.openshift.io", "v1"),
				newAPIService("build.openshift.io", "v1"),
				newAPIService("image.openshift.io", "v1"),
				newAPIService("route.openshift.io", "v1"),
				newAPIService("template.openshift.io", "v1"),
				newAPIService("user.openshift.io", "v1"),
			},
			newAPIServices: []*apiregistrationv1.APIService{
				newAPIService("authorization.openshift.io", "v1"),
				newAPIService("build.openshift.io", "v1"),
				newAPIService("image.openshift.io", "v1"),
				newAPIService("oauth.openshift.io", "v1"),
				newAPIService("route.openshift.io", "v1"),
				newAPIService("template.openshift.io", "v1"),
			},
			resultList: []string{
				"v1.authorization.openshift.io",
				"v1.build.openshift.io",
				"v1.image.openshift.io",
				"v1.oauth.openshift.io",
				"v1.route.openshift.io",
				"v1.template.openshift.io",
			},
			resultListChanged: true,
		},
		// scenario 4
		{
			name: "no diff",
			oldAPIServices: []*apiregistrationv1.APIService{
				newAPIService("authorization.openshift.io", "v1"),
				newAPIService("build.openshift.io", "v1"),
				newAPIService("image.openshift.io", "v1"),
				newAPIService("oauth.openshift.io", "v1"),
				newAPIService("route.openshift.io", "v1"),
				newAPIService("template.openshift.io", "v1"),
				newAPIService("user.openshift.io", "v1"),
			},
			newAPIServices: []*apiregistrationv1.APIService{
				newAPIService("authorization.openshift.io", "v1"),
				newAPIService("build.openshift.io", "v1"),
				newAPIService("image.openshift.io", "v1"),
				newAPIService("oauth.openshift.io", "v1"),
				newAPIService("route.openshift.io", "v1"),
				newAPIService("template.openshift.io", "v1"),
				newAPIService("user.openshift.io", "v1"),
			},
			resultList: []string{
				"v1.authorization.openshift.io",
				"v1.build.openshift.io",
				"v1.image.openshift.io",
				"v1.oauth.openshift.io",
				"v1.route.openshift.io",
				"v1.template.openshift.io",
				"v1.user.openshift.io",
			},
			resultListChanged: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			changed, newAPIServicesSet := apiServicesChanged(tc.oldAPIServices, tc.newAPIServices)

			if changed != tc.resultListChanged {
				t.Errorf("result list chaned = %v, expected it to change = %v", changed, tc.resultListChanged)
			}

			if !equality.Semantic.DeepEqual(tc.resultList, newAPIServicesSet.List()) {
				t.Errorf("incorect api services list returned: %s", diff.ObjectDiff(tc.resultList, newAPIServicesSet.List()))
			}
		})
	}
}

func TestHandlingControlOverTheAPI(t *testing.T) {
	const externalServerAnnotation = "authentication.operator.openshift.io/managed"

	testCases := []struct {
		name                    string
		existingAPIServices     []runtime.Object
		expectedAPIServices     []*apiregistrationv1.APIService
		expectedEventMsg        string
		expectsEvent            bool
		authOperatorUnavailable bool
		authOperatorManages     bool
	}{
		// scenario 1
		{
			name:                "checks if user.openshift.io and oauth.openshift.io groups are managed by an external server if all preconditions (authentication-operator status field set, annotations added) are valid",
			authOperatorManages: true,
			existingAPIServices: []runtime.Object{
				runtime.Object(newAPIService("build.openshift.io", "v1")),
				runtime.Object(newAPIService("apps.openshift.io", "v1")),
				runtime.Object(func() *apiregistrationv1.APIService {
					apiService := newAPIService("user.openshift.io", "v1")
					apiService.Annotations = map[string]string{}
					apiService.Annotations[externalServerAnnotation] = "true"
					return apiService
				}()),
				runtime.Object(func() *apiregistrationv1.APIService {
					apiService := newAPIService("oauth.openshift.io", "v1")
					apiService.Annotations = map[string]string{}
					apiService.Annotations[externalServerAnnotation] = "true"
					return apiService
				}()),
			},
			expectedAPIServices: []*apiregistrationv1.APIService{
				newAPIService("build.openshift.io", "v1"),
				newAPIService("apps.openshift.io", "v1"),
			},
			expectedEventMsg: "The new API Services list this operator will manage is [v1.apps.openshift.io v1.build.openshift.io]",
		},

		// scenario 2
		{
			name:                "checks that oauth.openshift.io group is not managed by an external server if it's missing the annotation",
			authOperatorManages: true,
			existingAPIServices: []runtime.Object{
				runtime.Object(newAPIService("build.openshift.io", "v1")),
				runtime.Object(newAPIService("apps.openshift.io", "v1")),
				runtime.Object(func() *apiregistrationv1.APIService {
					apiService := newAPIService("user.openshift.io", "v1")
					apiService.Annotations = map[string]string{}
					apiService.Annotations[externalServerAnnotation] = "true"
					return apiService
				}()),
				runtime.Object(newAPIService("oauth.openshift.io", "v1")),
			},
			expectedAPIServices: []*apiregistrationv1.APIService{
				newAPIService("build.openshift.io", "v1"),
				newAPIService("apps.openshift.io", "v1"),
				newAPIService("oauth.openshift.io", "v1"),
			},
			expectedEventMsg: "The new API Services list this operator will manage is [v1.apps.openshift.io v1.build.openshift.io v1.oauth.openshift.io]",
		},

		// scenario 3
		{
			name:                "authoritative/initial list is taken if authentication-operator wasn't found BUT the annotations were added",
			authOperatorManages: true,
			existingAPIServices: []runtime.Object{
				runtime.Object(newAPIService("build.openshift.io", "v1")),
				runtime.Object(newAPIService("apps.openshift.io", "v1")),
				runtime.Object(func() *apiregistrationv1.APIService {
					apiService := newAPIService("user.openshift.io", "v1")
					apiService.Annotations = map[string]string{}
					apiService.Annotations[externalServerAnnotation] = "true"
					return apiService
				}()),
				runtime.Object(func() *apiregistrationv1.APIService {
					apiService := newAPIService("oauth.openshift.io", "v1")
					apiService.Annotations = map[string]string{}
					apiService.Annotations[externalServerAnnotation] = "true"
					return apiService
				}()),
			},
			expectedAPIServices: []*apiregistrationv1.APIService{
				newAPIService("build.openshift.io", "v1"),
				newAPIService("apps.openshift.io", "v1"),
				func() *apiregistrationv1.APIService {
					apiService := newAPIService("user.openshift.io", "v1")
					apiService.Annotations = map[string]string{}
					apiService.Annotations[externalServerAnnotation] = "true"
					return apiService
				}(),
				func() *apiregistrationv1.APIService {
					apiService := newAPIService("oauth.openshift.io", "v1")
					apiService.Annotations = map[string]string{}
					apiService.Annotations[externalServerAnnotation] = "true"
					return apiService
				}(),
			},
			authOperatorUnavailable: true,
		},

		// scenario 4
		{
			name:                "authoritative/initial list is taken when ManagingOAuthAPIServer field set to false",
			authOperatorManages: false,
			existingAPIServices: []runtime.Object{
				runtime.Object(newAPIService("build.openshift.io", "v1")),
				runtime.Object(newAPIService("apps.openshift.io", "v1")),
				runtime.Object(func() *apiregistrationv1.APIService {
					apiService := newAPIService("user.openshift.io", "v1")
					apiService.Annotations = map[string]string{}
					apiService.Annotations[externalServerAnnotation] = "true"
					return apiService
				}()),
				runtime.Object(func() *apiregistrationv1.APIService {
					apiService := newAPIService("oauth.openshift.io", "v1")
					apiService.Annotations = map[string]string{}
					apiService.Annotations[externalServerAnnotation] = "true"
					return apiService
				}()),
			},
			expectedAPIServices: []*apiregistrationv1.APIService{
				newAPIService("build.openshift.io", "v1"),
				newAPIService("apps.openshift.io", "v1"),
				func() *apiregistrationv1.APIService {
					apiService := newAPIService("user.openshift.io", "v1")
					apiService.Annotations = map[string]string{}
					apiService.Annotations[externalServerAnnotation] = "true"
					return apiService
				}(),
				func() *apiregistrationv1.APIService {
					apiService := newAPIService("oauth.openshift.io", "v1")
					apiService.Annotations = map[string]string{}
					apiService.Annotations[externalServerAnnotation] = "true"
					return apiService
				}(),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			eventRecorder := events.NewInMemoryRecorder("")

			fakeAuthOperatorIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
			{
				authOperator := &operatorv1.Authentication{
					TypeMeta:   metav1.TypeMeta{},
					ObjectMeta: metav1.ObjectMeta{Name: "cluster"},
					Spec:       operatorv1.AuthenticationSpec{OperatorSpec: operatorv1.OperatorSpec{ManagementState: operatorv1.Managed}},
					Status:     operatorv1.AuthenticationStatus{ManagingOAuthAPIServer: tc.authOperatorManages, OperatorStatus: operatorv1.OperatorStatus{}},
				}

				if !tc.authOperatorUnavailable {
					err := fakeAuthOperatorIndexer.Add(authOperator)
					if err != nil {
						t.Fatal(err)
					}
				}
			}

			apiServices := []*apiregistrationv1.APIService{}
			fakeAPIRegistrationIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
			for _, rawService := range tc.existingAPIServices {
				service, ok := rawService.(*apiregistrationv1.APIService)
				if !ok {
					t.Fatal("unable to convert an api service to *apiregistrationv1.APIService")
				}
				apiServices = append(apiServices, service)
				if err := fakeAPIRegistrationIndexer.Add(service); err != nil {
					t.Fatal(err)
				}
			}

			target := NewAPIServicesToManage(
				apiregistrationv1lister.NewAPIServiceLister(fakeAPIRegistrationIndexer),
				operatorlistersv1.NewAuthenticationLister(fakeAuthOperatorIndexer),
				apiServices,
				eventRecorder,
				sets.NewString("v1.oauth.openshift.io", "v1.user.openshift.io"),
				externalServerAnnotation,
			)

			actualAPIServicesToManage, err := target.GetAPIServicesToManage()
			if err != nil {
				t.Fatal(err)
			}

			eventValidated := false
			for _, ev := range eventRecorder.Events() {
				if ev.Reason == "APIServicesToManageChanged" {
					if ev.Message != tc.expectedEventMsg {
						t.Errorf("unexpected APIServicesToManageChanged event message = %v, expected = %v", tc.expectedEventMsg, ev.Message)
					}
					eventValidated = true
				}
			}
			if !eventValidated && tc.expectsEvent {
				t.Error("APIServicesToManageChanged hasn't been found")
			}

			// validate
			if !equality.Semantic.DeepEqual(actualAPIServicesToManage, tc.expectedAPIServices) {
				t.Errorf("incorect api services list returned: %s", diff.ObjectDiff(actualAPIServicesToManage, tc.expectedAPIServices))
			}
		})
	}
}

func newAPIService(group, version string) *apiregistrationv1.APIService {
	return &apiregistrationv1.APIService{
		ObjectMeta: metav1.ObjectMeta{Name: version + "." + group, Annotations: map[string]string{"service.alpha.openshift.io/inject-cabundle": "true"}},
		Spec:       apiregistrationv1.APIServiceSpec{Group: group, Version: version, Service: &apiregistrationv1.ServiceReference{Namespace: "target-namespace", Name: "api"}, GroupPriorityMinimum: 9900, VersionPriority: 15},
		Status:     apiregistrationv1.APIServiceStatus{Conditions: []apiregistrationv1.APIServiceCondition{{Type: apiregistrationv1.Available, Status: apiregistrationv1.ConditionTrue}}},
	}
}
