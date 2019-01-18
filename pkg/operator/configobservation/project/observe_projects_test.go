package project

import (
	"testing"

	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/util/diff"
	"k8s.io/client-go/tools/cache"

	projectv1 "github.com/openshift/api/config/v1"
	configlistersv1 "github.com/openshift/client-go/config/listers/config/v1"
	"github.com/openshift/library-go/pkg/operator/events"

	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation"
)

func fakeProjectConfig(name string, spec projectv1.ProjectSpec) projectv1.Project {
	p := projectv1.Project{}
	p.Name = "cluster"
	p.Spec = spec
	return p
}

func TestObserveProjectRequestMessage(t *testing.T) {
	tests := []struct {
		name                 string
		existingConfig       map[string]interface{}
		expectedConfig       map[string]interface{}
		currentProjectConfig projectv1.Project
		expectErrorsCount    int
		expectEventCount     int
		listersNotSynced     bool
	}{
		{
			name:                 "simple update",
			existingConfig:       map[string]interface{}{"projectConfig": map[string]interface{}{"projectRequestMessage": "foo"}},
			expectedConfig:       map[string]interface{}{"projectConfig": map[string]interface{}{"projectRequestMessage": "bar"}},
			currentProjectConfig: fakeProjectConfig("cluster", projectv1.ProjectSpec{ProjectRequestMessage: "bar"}),
			expectEventCount:     1,
		},
		{
			name:                 "empty field",
			existingConfig:       map[string]interface{}{"projectConfig": map[string]interface{}{"projectRequestMessage": "foo"}},
			expectedConfig:       map[string]interface{}{"projectConfig": map[string]interface{}{"projectRequestMessage": ""}},
			currentProjectConfig: fakeProjectConfig("cluster", projectv1.ProjectSpec{ProjectRequestMessage: ""}),
			expectEventCount:     1,
		},
		{
			name:             "lister not synced",
			existingConfig:   map[string]interface{}{"projectConfig": map[string]interface{}{"projectRequestMessage": "foo"}},
			expectedConfig:   map[string]interface{}{"projectConfig": map[string]interface{}{"projectRequestMessage": "foo"}},
			listersNotSynced: true,
			expectEventCount: 0,
		},
		{
			name:                 "no existing",
			expectedConfig:       map[string]interface{}{"projectConfig": map[string]interface{}{"projectRequestMessage": "foo"}},
			currentProjectConfig: fakeProjectConfig("cluster", projectv1.ProjectSpec{ProjectRequestMessage: "foo"}),
			expectEventCount:     1,
		},
		{
			name:                 "no change",
			existingConfig:       map[string]interface{}{"projectConfig": map[string]interface{}{"projectRequestMessage": "foo"}},
			expectedConfig:       map[string]interface{}{"projectConfig": map[string]interface{}{"projectRequestMessage": "foo"}},
			currentProjectConfig: fakeProjectConfig("cluster", projectv1.ProjectSpec{ProjectRequestMessage: "foo"}),
			expectEventCount:     0, // Do not fire events on no-op change
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			indexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
			if err := indexer.Add(&test.currentProjectConfig); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			listers := configobservation.Listers{
				ProjectConfigLister: configlistersv1.NewProjectLister(indexer),
				ProjectConfigSynced: func() bool { return !test.listersNotSynced },
			}

			eventRecorder := events.NewInMemoryRecorder("")

			result, errs := ObserveProjectRequestMessage(listers, eventRecorder, test.existingConfig)
			if len(errs) != test.expectErrorsCount {
				t.Errorf("unexpected error count: %d != %d (errors: %#v)", len(errs), test.expectErrorsCount, errs)
				return
			}
			if len(eventRecorder.Events()) != test.expectEventCount {
				t.Errorf("unexpected event count: %d != %d (events: %#v)", len(eventRecorder.Events()), test.expectEventCount, eventRecorder.Events())
			}
			if !equality.Semantic.DeepEqual(test.expectedConfig, result) {
				t.Errorf("result does not match expected config: %s", diff.ObjectDiff(test.expectedConfig, result))
			}

		})
	}
}

func TestObserveProjectRequestTemplateName(t *testing.T) {
	tests := []struct {
		name                 string
		existingConfig       map[string]interface{}
		expectedConfig       map[string]interface{}
		currentProjectConfig projectv1.Project
		expectErrorsCount    int
		expectEventCount     int
		listersNotSynced     bool
	}{
		{
			name:                 "simple update",
			existingConfig:       map[string]interface{}{"projectConfig": map[string]interface{}{"projectRequestTemplate": map[string]interface{}{"name": "foo-template"}}},
			expectedConfig:       map[string]interface{}{"projectConfig": map[string]interface{}{"projectRequestTemplate": map[string]interface{}{"name": "bar-template"}}},
			currentProjectConfig: fakeProjectConfig("cluster", projectv1.ProjectSpec{ProjectRequestTemplate: projectv1.TemplateReference{Name: "bar-template"}}),
			expectEventCount:     1,
		},
		{
			name:                 "empty field",
			existingConfig:       map[string]interface{}{"projectConfig": map[string]interface{}{"projectRequestTemplate": map[string]interface{}{"name": "foo-template"}}},
			expectedConfig:       map[string]interface{}{"projectConfig": map[string]interface{}{"projectRequestTemplate": map[string]interface{}{"name": ""}}},
			currentProjectConfig: fakeProjectConfig("cluster", projectv1.ProjectSpec{ProjectRequestTemplate: projectv1.TemplateReference{Name: ""}}),
			expectEventCount:     1,
		},
		{
			name:             "lister not synced",
			existingConfig:   map[string]interface{}{"projectConfig": map[string]interface{}{"projectRequestTemplate": map[string]interface{}{"name": "foo-template"}}},
			expectedConfig:   map[string]interface{}{"projectConfig": map[string]interface{}{"projectRequestTemplate": map[string]interface{}{"name": "foo-template"}}},
			listersNotSynced: true,
			expectEventCount: 0,
		},
		{
			name:                 "no existing",
			expectedConfig:       map[string]interface{}{"projectConfig": map[string]interface{}{"projectRequestTemplate": map[string]interface{}{"name": "bar-template"}}},
			currentProjectConfig: fakeProjectConfig("cluster", projectv1.ProjectSpec{ProjectRequestTemplate: projectv1.TemplateReference{Name: "bar-template"}}),
			expectEventCount:     1,
		},
		{
			name:                 "no change",
			existingConfig:       map[string]interface{}{"projectConfig": map[string]interface{}{"projectRequestTemplate": map[string]interface{}{"name": "bar-template"}}},
			expectedConfig:       map[string]interface{}{"projectConfig": map[string]interface{}{"projectRequestTemplate": map[string]interface{}{"name": "bar-template"}}},
			currentProjectConfig: fakeProjectConfig("cluster", projectv1.ProjectSpec{ProjectRequestTemplate: projectv1.TemplateReference{Name: "bar-template"}}),
			expectEventCount:     0, // Do not fire events on no-op change
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			indexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
			if err := indexer.Add(&test.currentProjectConfig); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			listers := configobservation.Listers{
				ProjectConfigLister: configlistersv1.NewProjectLister(indexer),
				ProjectConfigSynced: func() bool { return !test.listersNotSynced },
			}

			eventRecorder := events.NewInMemoryRecorder("")

			result, errs := ObserveProjectRequestTemplateName(listers, eventRecorder, test.existingConfig)
			if len(errs) != test.expectErrorsCount {
				t.Errorf("unexpected error count: %d != %d (errors: %#v)", len(errs), test.expectErrorsCount, errs)
				return
			}
			if len(eventRecorder.Events()) != test.expectEventCount {
				t.Errorf("unexpected event count: %d != %d (events: %#v)", len(eventRecorder.Events()), test.expectEventCount, eventRecorder.Events())
			}

			if !equality.Semantic.DeepEqual(test.expectedConfig, result) {
				t.Errorf("result does not match expected config: %s", diff.ObjectDiff(test.expectedConfig, result))
				return
			}

		})
	}
}
