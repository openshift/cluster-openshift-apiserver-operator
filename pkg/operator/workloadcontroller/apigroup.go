package workloadcontroller

import (
	"fmt"
	"net/http"

	"k8s.io/client-go/rest"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

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
