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
ifndef JUNITFILE
	go test -race ./...
else
ifeq (, $(shell which gotest2junit 2>/dev/null))
$(error gotest2junit not found! Get it by `go get -u github.com/openshift/release/tools/gotest2junit`.)
endif
	go test -race -json ./... | gotest2junit > $(JUNITFILE)
endif
.PHONY: test-unit

images:
	imagebuilder -f Dockerfile -t openshift/cluster-openshift-apiserver-operator .
.PHONY: images

clean:
	$(RM) ./cluster-openshift-apiserver-operator
.PHONY: clean
