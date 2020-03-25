module github.com/openshift/cluster-openshift-apiserver-operator

go 1.13

require (
	github.com/coreos/etcd v3.3.15+incompatible
	github.com/getsentry/raven-go v0.2.1-0.20190513200303-c977f96e1095 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/gonum/graph v0.0.0-20170401004347-50b27dea7ebb
	github.com/jteeuwen/go-bindata v3.0.8-0.20151023091102-a0ff2567cfb7+incompatible
	github.com/kubernetes-sigs/kube-storage-version-migrator v0.0.0-20191127225502-51849bc15f17
	github.com/openshift/api v0.0.0-20200320142426-0de0d539b0c3
	github.com/openshift/build-machinery-go v0.0.0-20200211121458-5e3d6e570160
	github.com/openshift/client-go v0.0.0-20200320150128-a906f3d8e723
	github.com/openshift/library-go v0.0.0-20200325084904-c4f5592456b1
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.4.0
	go.uber.org/atomic v1.3.3-0.20181018215023-8dc6146f7569 // indirect
	go.uber.org/multierr v1.1.1-0.20180122172545-ddea229ff1df // indirect
	gonum.org/v1/netlib v0.0.0-20190331212654-76723241ea4e // indirect
	k8s.io/api v0.18.0-beta.2
	k8s.io/apimachinery v0.18.0-beta.2
	k8s.io/client-go v0.18.0-beta.2
	k8s.io/component-base v0.18.0-beta.2
	k8s.io/klog v1.0.0
	k8s.io/kube-aggregator v0.18.0-beta.2
	k8s.io/utils v0.0.0-20200229041039-0a110f9eb7ab
	sigs.k8s.io/structured-merge-diff v1.0.1-0.20191108220359-b1b620dd3f06 // indirect
)

replace github.com/kubernetes-sigs/kube-storage-version-migrator => github.com/openshift/kubernetes-kube-storage-version-migrator v0.0.3-0.20200312103335-32e07ea4f8ca
