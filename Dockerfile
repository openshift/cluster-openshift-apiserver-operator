#
# This is the integrated OpenShift Service Serving Cert Signer.  It signs serving certificates for use inside the platform.
#
# The standard name for this image is openshift/origin-cluster-openshift-apiserver-operator
#
FROM openshift/origin-release:golang-1.10
COPY . /go/src/github.com/openshift/cluster-openshift-apiserver-operator
RUN cd /go/src/github.com/openshift/cluster-openshift-apiserver-operator && go build ./cmd/cluster-openshift-apiserver-operator

FROM centos:7
COPY --from=0 /go/src/github.com/openshift/cluster-openshift-apiserver-operator/cluster-openshift-apiserver-operator /usr/bin/cluster-openshift-apiserver-operator
