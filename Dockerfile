FROM registry.ci.openshift.org/ocp/builder:rhel-9-golang-1.24-openshift-4.20 AS builder
WORKDIR /go/src/github.com/openshift/cluster-openshift-apiserver-operator
COPY . .
RUN GODEBUG=tls13=1 make build \
    && make tests-ext-build \
    && gzip cluster-openshift-apiserver-operator-tests-ext

FROM registry.ci.openshift.org/ocp/4.20:base-rhel9
COPY --from=builder /go/src/github.com/openshift/cluster-openshift-apiserver-operator/cluster-openshift-apiserver-operator /usr/bin/
COPY --from=builder /go/src/github.com/openshift/cluster-openshift-apiserver-operator/cluster-openshift-apiserver-operator-tests-ext.gz /usr/bin/
COPY manifests /manifests
LABEL io.openshift.release.operator true
