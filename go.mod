module github.com/openshift/cluster-openshift-apiserver-operator

go 1.13

require (
	github.com/coreos/etcd v3.3.15+incompatible
	github.com/getsentry/raven-go v0.2.1-0.20190513200303-c977f96e1095 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/gonum/graph v0.0.0-20170401004347-50b27dea7ebb
	github.com/jteeuwen/go-bindata v3.0.8-0.20151023091102-a0ff2567cfb7+incompatible
	github.com/openshift/api v3.9.1-0.20191212095247-c1898f32de35+incompatible
	github.com/openshift/client-go v0.0.0-20191205152420-9faca5198b4f
	github.com/openshift/library-go v0.0.0-20191212204559-ecee9828e7b3
	github.com/pkg/errors v0.8.1
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.4.0
	go.uber.org/atomic v1.3.3-0.20181018215023-8dc6146f7569 // indirect
	go.uber.org/multierr v1.1.1-0.20180122172545-ddea229ff1df // indirect
	k8s.io/api v0.0.0-kubernetes-1.17.0-rc.2
	k8s.io/apiextensions-apiserver v0.0.0-kubernetes-1.17.0-rc.2 // indirect
	k8s.io/apimachinery v0.17.1-beta.0
	k8s.io/apiserver v0.0.0-kubernetes-1.17.0-rc.2 // indirect
	k8s.io/client-go v0.0.0-kubernetes-1.17.0-rc.2
	k8s.io/component-base v0.0.0-kubernetes-1.17.0-rc.2
	k8s.io/klog v1.0.0
	k8s.io/kube-aggregator v0.0.0-kubernetes-1.17.0-rc.2
	k8s.io/utils v0.0.0-20191114184206-e782cd3c129f
)

replace (
	github.com/jteeuwen/go-bindata => github.com/jteeuwen/go-bindata v3.0.8-0.20151023091102-a0ff2567cfb7+incompatible
	github.com/openshift/api => github.com/openshift/api v3.9.1-0.20191212095247-c1898f32de35+incompatible
	github.com/openshift/client-go => github.com/openshift/client-go v0.0.0-20191205152420-9faca5198b4f
	github.com/openshift/library-go => github.com/openshift/library-go v0.0.0-20191212204559-ecee9828e7b3
	github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
	k8s.io/api => k8s.io/api v0.0.0-20191204082340-384b28a90b2b
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20191204090904-aab77140f100 // indirect
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.1-beta.0
	k8s.io/apiserver => k8s.io/apiserver v0.0.0-20191204085103-2ce178ac32b7 // indirect
	k8s.io/client-go => k8s.io/client-go v0.0.0-20191204083517-ea72ff2b5b2f
	k8s.io/component-base => k8s.io/component-base v0.0.0-20191204084121-18d14e17701e
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.0.0-20191204085536-307dc9fddc57
)
