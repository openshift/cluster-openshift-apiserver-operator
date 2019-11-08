# OpenShift API Server Operator


The OpenShift API Server operator manages and updates the [OpenShift API Server](https://github.com/openshift/origin). The operator is based on OpenShift [library-go](https://github.com/openshift/library-go) framework
 and it
 is installed via [Cluster Version Operator](https://github.com/openshift/cluster-version-operator) (CVO).

It contains the following sub-systems:

* Operator
* Configuration observer

By default, the operator exposes [Prometheus](https://prometheus.io) metrics via `metrics` service.
The metrics are collected from following components:

* OpenShift API Server Operator

## Configuration

The configuration observer component is responsible for reacting on external configuration changes.
For example, this allows external components ([registry](https://github.com/openshift/cluster-image-registry-operator), etc..)
to interact with the OpenShift API server configuration ([OpenShiftAPIServerConfig](https://github.com/openshift/api/blob/master/openshiftcontrolplane/v1/types.go#L13) custom resource).

Currently changes in following external components are being observed:

* `cluster` *images.config.openshift.io* custom resource
  - The observed CR resource is used to configure the `imagePolicyConfig.internalRegistryHostname` in Kubernetes API server configuration
* `cluster` *projects.config.openshift.io* custom resource
  - The observed CR resource is used to configure the Project request defaults
* `cluster` *ingress.config.openshift.io* custom resource
  - The observed CR resource is used to set `routingConfig.subdomain` in the OpenShift API server configuration.
   

The configuration for the OpenShift API server is the result of merging:

* a [default config](https://github.com/openshift/cluster-openshift-apiserver-operator/blob/master/bindata/v3.11.0/openshift-apiserver/defaultconfig.yaml)
* observed config (compare observed values above) `spec.spec.unsupportedConfigOverrides` from the `openshiftapiserveroperatorconfig`.

All of these are sparse configurations, i.e. unvalidated json snippets which are merged in order to form a valid configuration at the end.

## Debugging

To gather all information necessary for debugging operator please use the [must-gather](https://github.com/openshift/must-gather) tool.
