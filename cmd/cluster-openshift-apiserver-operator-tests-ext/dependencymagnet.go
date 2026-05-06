// This file imports test packages to ensure they are included in the build.
// These imports are necessary to register Ginkgo tests with the OpenShift Tests Extension framework.
package main

import (
	// Import test packages to register Ginkgo tests
	// The import below is necessary to ensure that the OAS operator tests are registered with the extension.
	_ "github.com/openshift/cluster-openshift-apiserver-operator/test/e2e"
	_ "github.com/openshift/cluster-openshift-apiserver-operator/test/extended"
)
