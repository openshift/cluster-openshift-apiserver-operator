package workloadcontroller

import (
	"fmt"
	"net/http"

	"k8s.io/client-go/rest"

	configv1 "github.com/openshift/api/config/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var apiServiceGroupVersions = []schema.GroupVersion{
	// these are all the apigroups we manage
	{Group: "apps.openshift.io", Version: "v1"},
	{Group: "authorization.openshift.io", Version: "v1"},
	{Group: "build.openshift.io", Version: "v1"},
	{Group: "image.openshift.io", Version: "v1"},
	{Group: "oauth.openshift.io", Version: "v1"},
	{Group: "project.openshift.io", Version: "v1"},
	{Group: "quota.openshift.io", Version: "v1"},
	{Group: "route.openshift.io", Version: "v1"},
	{Group: "security.openshift.io", Version: "v1"},
	{Group: "template.openshift.io", Version: "v1"},
	{Group: "user.openshift.io", Version: "v1"},
}

func checkForAPIs(restclient rest.Interface, groupVersions ...schema.GroupVersion) []string {
	missingMessages := []string{}
	for _, groupVersion := range groupVersions {
		url := "/apis/" + groupVersion.Group + "/" + groupVersion.Version

		statusCode := 0
		restclient.Get().AbsPath(url).Do().StatusCode(&statusCode)
		if statusCode != http.StatusOK {
			missingMessages = append(missingMessages, fmt.Sprintf("%s.%s is not ready: %v", groupVersion.Version, groupVersion.Group, statusCode))
		}
	}

	return missingMessages
}

func APIServiceReferences() []configv1.ObjectReference {
	ret := []configv1.ObjectReference{}
	for _, gv := range apiServiceGroupVersions {
		ret = append(ret, configv1.ObjectReference{Group: "apiregistration.k8s.io", Resource: "apiservices", Name: gv.Version + "." + gv.Group})
	}
	return ret
}
