# AGENTS.md

OpenShift operator that manages the OpenShift API Server for OpenShift clusters.

## What This Operator Does

This operator manages the OpenShift API Server operand and its supporting resources:

- **OpenShift API Server** (`openshift-apiserver` namespace) — serves the OpenShift API groups and runs as a deployment
- **Configuration observation** — watches cluster config and writes `OpenShiftAPIServer.spec.observedConfig`
- **APIService registration** — registers the OpenShift API groups served by the operand
- **Encryption** — manages encryption configuration, key rotation, migration, and KMS preflight
- **Connectivity checks** — creates `PodNetworkConnectivityCheck` resources for etcd and kube-apiserver dependencies
- **Status** — reports operator and operand health through the `openshift-apiserver` `ClusterOperator`

The operator runs in `openshift-apiserver-operator` and is installed by the Cluster Version Operator.

## Commands

```bash
# Build
make build                          # Builds the operator and tests-extension binaries

# Test
make test                           # Same as make test-unit
make test-unit                      # go test -race ./pkg/... ./cmd/...
make test-e2e                       # General e2e tests; requires a cluster

# Slow, serial, cluster-mutating encryption tests
make test-e2e-encryption
make test-e2e-encryption-rotation
make test-e2e-encryption-kms

# Verification
make verify                         # gofmt, go vet, Go version checks, and bindata verification
make verify-bindata

# Generated bindata
make update-bindata                 # Required after changing files under bindata/

# OpenShift tests extension
make tests-ext-build
make tests-ext-update               # Required when adding, renaming, or deleting extended tests

# Inspect resource dependencies
make build
./cluster-openshift-apiserver-operator resource-graph
```

Cluster tests require a valid `KUBECONFIG`.

## Tech Stack

- **Go** — See `go.mod` for version
- **OpenShift library-go** — Controller factory, `resourceapply`, operator helpers
- **OpenShift api** — OpenShift API Server CRDs and config types
- **klog/v2** — Structured logging

## Feature Gates

- Feature gate state is loaded in `pkg/operator/starter.go` from `ClusterVersion` and `FeatureGate`.
- Use the existing `featureGateAccessor`; do not read `FeatureGate` resources directly from controllers.
- The config observer renders known feature gates and their enabled or disabled values into `apiServerArguments.feature-gates`.
- A post-startup change to the feature gate set causes the operator process to exit and restart.
- `KMSEncryption` controls KMS plugin sidecar injection in the operand `Deployment`.
- `WatchList` is force-disabled in the config observer because the `Project` API does not support it.
- When changing feature-gated behavior, test both the enabled and disabled states where feasible.

## Always Do

- **Use informers and listers for reads** — controllers should avoid direct API reads on every sync
- **Use `resourceapply` helpers for writes** — for example, `ApplyDeployment()` and `ApplyConfigMap()`
- **Propagate `context.Context`** — pass an existing context through to downstream calls that accept one
- **Return errors from `sync()`** — let the controller factory handle retries
- **Run `make update-bindata`** after changing any file under `bindata/`
- **Run `make test-unit` and `make verify`** before considering Go or manifest changes complete
- **Run `make tests-ext-update`** after adding, renaming, or deleting extended tests
- **Prefer table-driven unit tests** — existing tests in `pkg/operator` follow that style
- **Log with `klog`** — use `klog.V()` levels rather than ad hoc printing
- **Keep RBAC in sync** — new controllers or resources may require changes in `manifests/`

## Ask First

- **Adding a controller** — requires registration in `pkg/operator/starter.go` and likely RBAC changes
- **Changing RBAC or service accounts** — inspect `manifests/04_roles.yaml` and `manifests/05_serviceaccount.yaml`
- **Changing encryption or KMS behavior** — security-sensitive and requires the serial encryption suites
- **Changing operand rollout or availability semantics** — affects cluster API availability
- **Changing namespaces** — this operator depends on fixed operator, operand, and config namespaces
- **Changing generated bindata structure** — confirm the Makefile generation path still matches the source tree
- **Changing metrics, TLS, or network-policy behavior** — verify both unit tests and cluster behavior
- **Changing tests-extension suite mapping** — update `cmd/cluster-openshift-apiserver-operator-tests-ext/main.go` and regenerate metadata

## Never Do

- **Never hand-edit `pkg/operator/v311_00_assets/bindata.go`** — it is generated from `bindata/v3.11.0/`
- **Never assume `//go:embed`** — this repository uses generated bindata
- **Never edit `vendor/` directly** — use `go get`, `go mod tidy`, and `go mod vendor`
- **Never treat `make verify` as a substitute for unit tests** — it does not run `go test`
- **Never run encryption tests in parallel** — they are serial and mutate global cluster state
- **Never commit secrets or credentials**
- **Never bypass informers/listers for controller reads** when an informer is already available
- **Never skip `make tests-ext-update`** after changing extended-test names or metadata

## Controller Pattern

Repo-local controllers use the library-go factory pattern:

```go
func NewMyController(
    operatorClient v1helpers.OperatorClient,
    kubeInformers kubeinformers.SharedInformerFactory,
    recorder events.Recorder,
) factory.Controller {
    c := &myController{
        operatorClient: operatorClient,
        lister:        kubeInformers.Apps().V1().Deployments().Lister(),
    }

    return factory.New().
        WithInformers(kubeInformers.Apps().V1().Deployments().Informer()).
        WithSync(c.sync).
        ToController("MyController", recorder)
}

func (c *myController) sync(ctx context.Context, syncCtx factory.SyncContext) error {
    deployment, err := c.lister.Deployments("openshift-apiserver").Get("openshift-apiserver")
    if err != nil {
        return err
    }

    // Reconcile using resourceapply helpers.
    return nil
}
```

Existing examples:

- `pkg/operator/nsfinalizercontroller/finalizer_controller.go`
- `pkg/operator/connectivitycheckcontroller/connectivity_check_controller.go`
- `pkg/operator/workload/workload_openshiftapiserver_v311_00_sync.go`

## Key File Locations

```text
cmd/cluster-openshift-apiserver-operator/             Operator entrypoint
cmd/cluster-openshift-apiserver-operator-tests-ext/   Tests-extension entrypoint

pkg/operator/
  starter.go                                          Controller registration and startup
  workload/                                           Operand workload controller
  configobservation/                                  Config observers
  resourcesynccontroller/                             Shared resource sync
  connectivitycheckcontroller/                        Network connectivity checks
  nsfinalizercontroller/                              Namespace finalizer
  v311_00_assets/bindata.go                           Generated bindata; do not edit directly

pkg/test/assets_test.go                               YAML correctness tests

bindata/v3.11.0/                                      Operand manifests and default config
manifests/                                            Operator deployment, RBAC, and monitoring
test/
  e2e/                                                General e2e tests
  e2e-encryption/                                     Encryption e2e tests
  e2e-encryption-rotation/                            Encryption rotation e2e tests
  e2e-encryption-kms/                                 KMS encryption e2e tests
  extended/                                           OpenShift tests extension suite
  library/                                            Shared test utilities
```

## Namespaces

- `openshift-apiserver-operator` — the operator runs here
- `openshift-apiserver` — the operand runs here
- `openshift-config` — user-specified configuration
- `openshift-config-managed` — machine-managed configuration

## Common Workflows

### Modify Operand Manifests

1. Edit the YAML under `bindata/v3.11.0/`
2. Run `make update-bindata`
3. Run `make test-unit`
4. Run `make verify`

`pkg/test/assets_test.go` validates YAML parsing for both `manifests/` and `bindata/`.

### Add or Change a Controller

1. Implement the controller under `pkg/operator/`
2. Register it in `pkg/operator/starter.go`
3. Add or update RBAC in `manifests/`
4. Add unit tests beside the controller
5. Run `make test-unit` and `make verify`

### Add or Update Extended Tests

1. Edit tests under `test/extended/`
2. Run `make tests-ext-update`
3. Run `make tests-ext-build`
4. Run the relevant suite with `./cluster-openshift-apiserver-operator-tests-ext`

### Inspect Resource Dependencies

```bash
make build
./cluster-openshift-apiserver-operator resource-graph
```

This prints a DOT graph of the operator’s resource dependencies.

## Security Notes

- This operator manages a cluster API server and encryption-at-rest configuration
- Never log secrets, private keys, or credentials
- Encryption and KMS changes require the serial encryption test suites
- Treat TLS, ServiceMonitor, and network-policy changes as cluster-sensitive
