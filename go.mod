module github.com/openshift/cluster-openshift-apiserver-operator

go 1.13

require (
	github.com/getsentry/raven-go v0.2.1-0.20190513200303-c977f96e1095 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/go-bindata/go-bindata v3.1.2+incompatible
	github.com/gonum/graph v0.0.0-20170401004347-50b27dea7ebb
	github.com/kubernetes-sigs/kube-storage-version-migrator v0.0.0-20191127225502-51849bc15f17
	github.com/openshift/api v0.0.0-20200728200559-811027b63048
	github.com/openshift/build-machinery-go v0.0.0-20200713135615-1f43d26dccc7
	github.com/openshift/client-go v0.0.0-20200722173614-5a1b0aaeff15
	github.com/openshift/library-go v0.0.0-20200730143437-a1811581365b
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.4.0
	go.etcd.io/etcd v0.5.0-alpha.5.0.20200716221620-18dfb9cca345
	go.uber.org/multierr v1.1.1-0.20180122172545-ddea229ff1df // indirect
	k8s.io/api v0.19.0-rc.2
	k8s.io/apimachinery v0.19.0-rc.2
	k8s.io/apiserver v0.19.0-rc.2
	k8s.io/client-go v0.19.0-rc.2
	k8s.io/component-base v0.19.0-rc.2
	k8s.io/klog v1.0.0
	k8s.io/kube-aggregator v0.19.0-rc.2
	k8s.io/utils v0.0.0-20200720150651-0bdb4ca86cbc
)

replace (
	github.com/coreos/etcd => github.com/coreos/etcd v3.3.17+incompatible // Pinned for openshift
	github.com/kubernetes-sigs/kube-storage-version-migrator => github.com/openshift/kubernetes-kube-storage-version-migrator v0.0.3-0.20200312103335-32e07ea4f8ca
	github.com/openshift/library-go => github.com/tkashem/library-go v0.0.0-20200730225339-19c5b9919658
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
