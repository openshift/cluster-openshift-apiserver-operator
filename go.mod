module github.com/openshift/cluster-openshift-apiserver-operator

go 1.13

require (
	github.com/certifi/gocertifi v0.0.0-20180905225744-ee1a9a0726d2 // indirect
	github.com/coreos/etcd v3.3.15+incompatible
	github.com/getsentry/raven-go v0.2.1-0.20190513200303-c977f96e1095 // indirect
	github.com/ghodss/yaml v0.0.0-20150909031657-73d445a93680
	github.com/gonum/blas v0.0.0-20181208220705-f22b278b28ac // indirect
	github.com/gonum/floats v0.0.0-20181209220543-c233463c7e82 // indirect
	github.com/gonum/graph v0.0.0-20170401004347-50b27dea7ebb
	github.com/gonum/internal v0.0.0-20181124074243-f884aa714029 // indirect
	github.com/gonum/lapack v0.0.0-20181123203213-e4cdc5a0bff9 // indirect
	github.com/gonum/matrix v0.0.0-20181209220409-c518dec07be9 // indirect
	github.com/jteeuwen/go-bindata v0.0.0-00010101000000-000000000000
	github.com/kubernetes-sigs/kube-storage-version-migrator v0.0.0-00010101000000-000000000000 // indirect
	github.com/openshift/api v3.9.1-0.20191105214740-21e87c8db569+incompatible
	github.com/openshift/client-go v0.0.0-20191022152013-2823239d2298
	github.com/openshift/library-go v0.0.0-20191106124920-97392cdc9bd3
	github.com/pkg/errors v0.8.0
	github.com/pkg/profile v1.3.0 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.3
	github.com/stretchr/testify v1.3.0
	go.uber.org/atomic v1.3.3-0.20181018215023-8dc6146f7569 // indirect
	go.uber.org/multierr v1.1.1-0.20180122172545-ddea229ff1df // indirect
	go.uber.org/zap v1.9.2-0.20180814183419-67bc79d13d15 // indirect
	k8s.io/api v0.0.0-kubernetes-1.16.2
	k8s.io/apimachinery v0.0.0-kubernetes-1.16.2
	k8s.io/client-go v0.0.0-kubernetes-1.16.2
	k8s.io/component-base v0.0.0-20190918160511-547f6c5d7090
	k8s.io/klog v0.4.0
	k8s.io/kube-aggregator v0.0.0-kubernetes-1.16.2
	k8s.io/utils v0.0.0-20190801114015-581e00157fb1
)

replace (
	github.com/jteeuwen/go-bindata => github.com/jteeuwen/go-bindata v3.0.8-0.20151023091102-a0ff2567cfb7+incompatible
	github.com/kubernetes-sigs/kube-storage-version-migrator => github.com/openshift/kubernetes-kube-storage-version-migrator v0.0.3-0.20191031181212-bdca3bf7d454
	github.com/openshift/api => github.com/openshift/api v0.0.0-20191201231411-9f834e337466
	github.com/openshift/client-go => github.com/openshift/client-go v0.0.0-20191022152013-2823239d2298
	github.com/openshift/library-go => github.com/openshift/library-go v0.0.0-20191203164047-7cbbcc8aab8b
	github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
	k8s.io/api => k8s.io/api v0.0.0-20190918155943-95b840bb6a1f
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190918161926-8f644eb6e783
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190913080033-27d36303b655
	k8s.io/apiserver => k8s.io/apiserver v0.0.0-20190918160949-bfa5e2e684ad
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190918160344-1fbdaa4c8d90
	k8s.io/component-base => k8s.io/component-base v0.0.0-20190918160511-547f6c5d7090
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.0.0-20190918161219-8c8f079fddc3
)
