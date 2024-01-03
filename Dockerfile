FROM registry.ci.openshift.org/ocp/builder:rhel-9-golang-1.20-openshift-4.15 AS builder
WORKDIR /go/src/github.com/openshift/cluster-openshift-apiserver-operator
COPY . .
RUN GODEBUG=tls13=1 go build ./cmd/cluster-openshift-apiserver-operator

FROM registry.ci.openshift.org/ocp/4.15:base-rhel9
COPY --from=builder /go/src/github.com/openshift/cluster-openshift-apiserver-operator/cluster-openshift-apiserver-operator /usr/bin/
COPY manifests /manifests
COPY vendor/github.com/openshift/api/operator/v1/*_openshift-apiserver-operator_*.crd.yaml manifests/
LABEL io.openshift.release.operator true
