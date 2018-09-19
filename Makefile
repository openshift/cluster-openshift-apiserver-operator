all build:
	go build ./cmd/cluster-openshift-apiserver-operator
.PHONY: all build

verify-govet:
	go vet ./...
.PHONY: verify-govet

verify: verify-govet
	hack/verify-gofmt.sh
	hack/verify-codegen.sh
	hack/verify-generated-bindata.sh
.PHONY: verify

test test-unit:
	go test -race ./...
.PHONY: test-unit

images:
	imagebuilder -f Dockerfile -t openshift/cluster-openshift-apiserver-operator .
.PHONY: images

clean:
	$(RM) ./cluster-openshift-apiserver-operator
.PHONY: clean
