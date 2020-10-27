module github.com/openshift/cluster-openshift-apiserver-operator

go 1.13

require (
	github.com/getsentry/raven-go v0.2.1-0.20190513200303-c977f96e1095 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/go-bindata/go-bindata v3.1.2+incompatible
	github.com/gonum/graph v0.0.0-20170401004347-50b27dea7ebb
	github.com/openshift/api v0.0.0-20201022163842-96b75c633fdd
	github.com/openshift/build-machinery-go v0.0.0-20200917070002-f171684f77ab
	github.com/openshift/client-go v0.0.0-20200827190008-3062137373b5
	github.com/openshift/library-go v0.0.0-20201013192036-5bd7c282e3e7
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.4.0
	go.etcd.io/etcd v0.5.0-alpha.5.0.20200819165624-17cef6e3e9d5
	go.uber.org/multierr v1.1.1-0.20180122172545-ddea229ff1df // indirect
	k8s.io/api v0.19.2
	k8s.io/apiextensions-apiserver v0.19.0
	k8s.io/apimachinery v0.19.2
	k8s.io/apiserver v0.19.0
	k8s.io/client-go v0.19.0
	k8s.io/component-base v0.19.0
	k8s.io/klog/v2 v2.3.0
	k8s.io/kube-aggregator v0.19.0
	k8s.io/utils v0.0.0-20200729134348-d5654de09c73
	sigs.k8s.io/kube-storage-version-migrator v0.0.3
)

replace (
	github.com/coreos/etcd => github.com/coreos/etcd v3.3.17+incompatible // Pinned for openshift
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
