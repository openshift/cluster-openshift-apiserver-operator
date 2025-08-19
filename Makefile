all: build
.PHONY: all

# Include the library makefile
include $(addprefix ./vendor/github.com/openshift/build-machinery-go/make/, \
	golang.mk \
	targets/openshift/bindata.mk \
	targets/openshift/deps.mk \
	targets/openshift/images.mk \
	targets/openshift/operator/telepresence.mk \
)

IMAGE_REGISTRY :=registry.svc.ci.openshift.org

ENCRYPTION_PROVIDERS=aescbc aesgcm
ENCRYPTION_PROVIDER?=aescbc

# -------------------------------------------------------------------
# OpenShift Tests Extension (Cluster OpenShift API Server Operator)
# -------------------------------------------------------------------
TESTS_EXT_BINARY := cluster-openshift-apiserver-operator-tests-ext
TESTS_EXT_PACKAGE := ./cmd/cluster-openshift-apiserver-operator-tests-ext

TESTS_EXT_GIT_COMMIT := $(shell git rev-parse --short HEAD)
TESTS_EXT_BUILD_DATE := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
TESTS_EXT_GIT_TREE_STATE := $(shell if git diff --quiet; then echo clean; else echo dirty; fi)

TESTS_EXT_LDFLAGS := -X 'github.com/openshift-eng/openshift-tests-extension/pkg/version.CommitFromGit=$(TESTS_EXT_GIT_COMMIT)' \
                     -X 'github.com/openshift-eng/openshift-tests-extension/pkg/version.BuildDate=$(TESTS_EXT_BUILD_DATE)' \
                     -X 'github.com/openshift-eng/openshift-tests-extension/pkg/version.GitTreeState=$(TESTS_EXT_GIT_TREE_STATE)'


# This will call a macro called "build-image" which will generate image specific targets based on the parameters:
# $0 - macro name
# $1 - target name
# $2 - image ref
# $3 - Dockerfile path
# $4 - context directory for image build
$(call build-image,ocp-cluster-openshift-apiserver-operator,$(IMAGE_REGISTRY)/ocp/4.2:cluster-openshift-apiserver-operator,./Dockerfile.rhel7,.)

# This will call a macro called "add-bindata" which will generate bindata specific targets based on the parameters:
# $0 - macro name
# $1 - target suffix
# $2 - input dirs
# $3 - prefix
# $4 - pkg
# $5 - output
# It will generate targets {update,verify}-bindata-$(1) logically grouping them in unsuffixed versions of these targets
# and also hooked into {update,verify}-generated for broader integration.
$(call add-bindata,v3.11.0,./bindata/v3.11.0/...,bindata,v311_00_assets,pkg/operator/v311_00_assets/bindata.go)

$(call verify-golang-versions,Dockerfile)

clean:
	$(RM) ./cluster-openshift-apiserver-operator
.PHONY: clean

GO_TEST_PACKAGES :=./pkg/... ./cmd/...

TEST_E2E_ENCRYPTION_TARGETS=$(addprefix test-e2e-encryption-,$(ENCRYPTION_PROVIDERS))

# these are extremely slow serial e2e encryption tests that modify the cluster's global state
test-e2e-encryption: GO_TEST_PACKAGES :=./test/e2e-encryption/...
test-e2e-encryption: GO_TEST_FLAGS += -v
test-e2e-encryption: GO_TEST_FLAGS += -timeout 4h
test-e2e-encryption: GO_TEST_FLAGS += -p 1
test-e2e-encryption: GO_TEST_ARGS += -args -provider=$(ENCRYPTION_PROVIDER)
test-e2e-encryption: test-unit
.PHONY: test-e2e-encryption

.PHONY: $(TEST_E2E_ENCRYPTION_TARGETS)
$(TEST_E2E_ENCRYPTION_TARGETS): test-e2e-encryption-%:
	ENCRYPTION_PROVIDER=$* $(MAKE) test-e2e-encryption

TEST_E2E_ENCRYPTION_ROTATION_TARGETS=$(addprefix test-e2e-encryption-rotation-,$(ENCRYPTION_PROVIDERS))

# these are extremely slow serial e2e encryption rotation tests that modify the cluster's global state
test-e2e-encryption-rotation: GO_TEST_PACKAGES :=./test/e2e-encryption-rotation/...
test-e2e-encryption-rotation: GO_TEST_FLAGS += -v
test-e2e-encryption-rotation: GO_TEST_FLAGS += -timeout 4h
test-e2e-encryption-rotation: GO_TEST_FLAGS += -p 1
test-e2e-encryption-rotation: GO_TEST_ARGS += -args -provider=$(ENCRYPTION_PROVIDER)
test-e2e-encryption-rotation: test-unit
.PHONY: test-e2e-encryption-rotation

.PHONY: $(TEST_E2E_ENCRYPTION_ROTATION_TARGETS)
$(TEST_E2E_ENCRYPTION_ROTATION_TARGETS): test-e2e-encryption-rotation-%:
	ENCRYPTION_PROVIDER=$* $(MAKE) test-e2e-encryption-rotation

.PHONY: test-e2e
test-e2e: GO_TEST_PACKAGES :=./test/e2e/...
test-e2e: GO_TEST_FLAGS += -timeout 1h
test-e2e: test-unit

# -------------------------------------------------------------------
# Build binary with metadata (CI-compliant)
# -------------------------------------------------------------------
.PHONY: tests-ext-build
tests-ext-build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) GO_COMPLIANCE_POLICY=exempt_all CGO_ENABLED=0 \
	go build -o $(TESTS_EXT_BINARY) -ldflags "$(TESTS_EXT_LDFLAGS)" $(TESTS_EXT_PACKAGE)

# -------------------------------------------------------------------
# Run "update" and strip env-specific metadata
# -------------------------------------------------------------------
.PHONY: tests-ext-update
tests-ext-update: tests-ext-build
	./$(TESTS_EXT_BINARY) update
	for f in .openshift-tests-extension/*.json; do \
		jq 'map(del(.codeLocations))' "$$f" > tmpp && mv tmpp "$$f"; \
	done

# Configure the 'telepresence' target
# See vendor/github.com/openshift/build-machinery-go/scripts/run-telepresence.sh for usage and configuration details
export TP_DEPLOYMENT_YAML ?=./manifests/0000_30_openshift-apiserver-operator_07_deployment.yaml
export TP_CMD_PATH ?=./cmd/cluster-openshift-apiserver-operator
