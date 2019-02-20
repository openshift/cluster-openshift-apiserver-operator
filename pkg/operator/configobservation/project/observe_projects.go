package project

import (
	"github.com/golang/glog"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/openshift/library-go/pkg/operator/configobserver"
	"github.com/openshift/library-go/pkg/operator/events"

	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation"
)

var (
	// This represents a JSON path for openshiftcontrolplane/v1.OpenShiftAPIServerConfig
	projectRequestMessagePath = []string{"projectConfig", "projectRequestMessage"}

	// This represents a JSON path for openshiftcontrolplane/v1.OpenShiftAPIServerConfig
	// Note: This is not typo, there is no "name" here as this path is not the same as config.openshift.io/v1
	projectRequestTemplateNamePath = []string{"projectConfig", "projectRequestTemplate"}
)

// ObserveProjectRequestTemplateName observers changes to config.openshift.io/Project resource field 'spec.projectRequestTemplate.Name' and update the existing apiserver
// configuration when a change it found.
func ObserveProjectRequestTemplateName(genericListers configobserver.Listers, recorder events.Recorder, existingConfig map[string]interface{}) (map[string]interface{}, []error) {
	listers := genericListers.(configobservation.Listers)
	errs := []error{}

	prevObservedConfig := map[string]interface{}{}

	currentProjectRequestTemplateName, exists, err := unstructured.NestedString(existingConfig, projectRequestTemplateNamePath...)
	if err != nil {
		return prevObservedConfig, append(errs, err)
	}

	if exists && len(currentProjectRequestTemplateName) > 0 {
		if err := unstructured.SetNestedField(prevObservedConfig, currentProjectRequestTemplateName, projectRequestTemplateNamePath...); err != nil {
			return prevObservedConfig, append(errs, err)
		}
	}

	observedConfig := map[string]interface{}{}

	currentClusterInstance, err := listers.ProjectConfigLister.Get("cluster")
	if errors.IsNotFound(err) {
		glog.V(4).Infof("project.config.openshift.io/v1: cluster: not found")
		return observedConfig, errs
	}
	if err != nil {
		return prevObservedConfig, append(errs, err)
	}

	// the openshift-apiserver takes a namespace/name format, but we require a particular namespace.  Prepend it here.
	observedProjectRequestTemplateName := ""
	if len(currentClusterInstance.Spec.ProjectRequestTemplate.Name) > 0 {
		observedProjectRequestTemplateName = operatorclient.GlobalUserSpecifiedConfigNamespace + "/" + currentClusterInstance.Spec.ProjectRequestTemplate.Name
		if err := unstructured.SetNestedField(observedConfig, observedProjectRequestTemplateName, projectRequestTemplateNamePath...); err != nil {
			return prevObservedConfig, append(errs, err)
		}
	}

	// no change, return early to skip the event
	if observedProjectRequestTemplateName == currentProjectRequestTemplateName {
		return observedConfig, errs
	}

	recorder.Eventf("ProjectRequestTemplateChanged", "ProjectRequestTemplate changed from %q to %q", currentProjectRequestTemplateName, observedProjectRequestTemplateName)

	return observedConfig, errs
}

// ObserveProjectRequestMessage observers changes to config.openshift.io/Project resource field 'spec.projectRequestMessage' and update the existing apiserver
// configuration when a change it found.
func ObserveProjectRequestMessage(genericListers configobserver.Listers, recorder events.Recorder, existingConfig map[string]interface{}) (map[string]interface{}, []error) {
	listers := genericListers.(configobservation.Listers)
	errs := []error{}

	prevObservedConfig := map[string]interface{}{}

	currentProjectRequestMessage, exists, err := unstructured.NestedString(existingConfig, projectRequestMessagePath...)
	if err != nil {
		return prevObservedConfig, append(errs, err)
	}

	if exists && len(currentProjectRequestMessage) > 0 {
		if err := unstructured.SetNestedField(prevObservedConfig, currentProjectRequestMessage, projectRequestMessagePath...); err != nil {
			return prevObservedConfig, append(errs, err)
		}
	}

	observedConfig := map[string]interface{}{}

	currentClusterInstance, err := listers.ProjectConfigLister.Get("cluster")
	if errors.IsNotFound(err) {
		glog.V(4).Infof("project.config.openshift.io/v1: cluster: not found")
		return observedConfig, errs
	}
	if err != nil {
		return prevObservedConfig, append(errs, err)
	}

	observedProjectRequestMessage := currentClusterInstance.Spec.ProjectRequestMessage

	if err := unstructured.SetNestedField(observedConfig, observedProjectRequestMessage, projectRequestMessagePath...); err != nil {
		return prevObservedConfig, append(errs, err)
	}

	// no change, return early to skip the event
	if observedProjectRequestMessage == currentProjectRequestMessage {
		return observedConfig, errs
	}

	recorder.Eventf("ProjectRequestMessageChanged", "ProjectRequestMessage changed from %q to %q", currentProjectRequestMessage, observedProjectRequestMessage)

	return observedConfig, errs
}
