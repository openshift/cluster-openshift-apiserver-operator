# Contributing

Thanks for contributing to the OpenShift API Server operator.

This repository is part of the OpenShift control-plane component family. For repository-specific development context, see:

- [AGENTS.md](AGENTS.md) — operational commands, conventions, and guardrails
- [ARCHITECTURE.md](ARCHITECTURE.md) — controller and operand structure
- [Makefile](Makefile) — canonical build, test, and verification targets
- [OWNERS](OWNERS) / [OWNERS_ALIASES](OWNERS_ALIASES) — reviewers and approvers

This guide contains the following sections:

- [Coding conventions](#coding-conventions)
- [Before opening a pull request](#before-opening-a-pull-request)
- [Repository-specific workflows](#repository-specific-workflows)
- [Verifying your changes / creating an OpenShift cluster from a PR](#verifying-your-changes--creating-an-openshift-cluster-from-a-pr)
- [Pull request expectations](#pull-request-expectations)
- [Review](#review)

## Coding Conventions

Follow the Kubernetes Go conventions and match the style of nearby code.

Reference guides:

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)
- [Go Landmines](https://gist.github.com/lavalamp/4bd23295a9f32706a48f)
- [Go comment conventions](https://go.dev/doc/comment)

Additional repository-specific guidance:

- Keep changes narrowly scoped and avoid unrelated refactors.
- Prefer table-driven tests when covering multiple scenarios.
- Test behavior, not implementation details.
- In controllers, use informers and listers for reads and `resourceapply` for declarative writes.
- Propagate an existing `context.Context`; do not replace it with `context.Background()` or `context.TODO()`.
- Wrap errors with enough context to identify the failing operation.
- Use `klog` for logging, with `klog.V()` levels for debug output.
- Add comments for non-obvious behavior; exported symbols should have doc comments.
- Avoid new dependencies unless the standard library or an already-vendored package cannot do the job.
- Keep new controller code under the existing `pkg/operator/` structure and register it in `pkg/operator/starter.go`.
- For feature-gated behavior, validate both the enabled and disabled states where feasible.

## Before Opening a Pull Request

Run the checks that match your change:

```sh
make test-unit
make verify
```

`make verify` checks formatting, vet, Go-version consistency, and bindata; it does not run unit tests.

If you changed files under [bindata/v3.11.0](bindata/v3.11.0), also run:

```sh
make update-bindata
make verify-bindata
```

If you added, renamed, or deleted tests under [test/extended](test/extended), also run:

```sh
make tests-ext-update
make tests-ext-build
```

If you changed `go.mod` or dependencies, run:

```sh
go mod tidy
go mod vendor
```

## Repository-Specific Workflows

### Bindata

This repository uses generated bindata, not `//go:embed`.

- Source files live in [bindata/v3.11.0](bindata/v3.11.0)
- Generated output lives in [pkg/operator/v311_00_assets/bindata.go](pkg/operator/v311_00_assets/bindata.go)
- Never edit the generated file directly
- Run `make update-bindata` after changing any source YAML under `bindata/`

### Operator and Operand Manifests

Operator payload manifests live in [manifests](manifests).
Operand manifests and default configuration live in [bindata/v3.11.0](bindata/v3.11.0).

YAML parsing is covered by [pkg/test/assets_test.go](pkg/test/assets_test.go).

### Tests

| Area | Location | Notes |
|---|---|---|
| Unit tests | `pkg/` and `cmd/` | Run with `make test-unit` |
| General e2e | [test/e2e](test/e2e) | Cluster required |
| Encryption | [test/e2e-encryption](test/e2e-encryption) | Serial, 4h, mutates global state |
| Encryption rotation | [test/e2e-encryption-rotation](test/e2e-encryption-rotation) | Serial, 4h, mutates global state |
| KMS encryption | [test/e2e-encryption-kms](test/e2e-encryption-kms) | Serial, 4h |
| Extended tests | [test/extended](test/extended) | See [test/extended/README.md](test/extended/README.md) |

Cluster tests require a valid `KUBECONFIG`.

### Extended Tests

Extended tests are built and run through the OpenShift tests extension framework.

For building, listing, running, renaming, and deleting tests, see:

- [test/extended/README.md](test/extended/README.md)
- [cmd/cluster-openshift-apiserver-operator-tests-ext](cmd/cluster-openshift-apiserver-operator-tests-ext)

After changing test names or metadata, run `make tests-ext-update`.

## Verifying Your Changes / Creating an OpenShift Cluster from a PR

Some changes need verification on a real OpenShift cluster, especially encryption, KMS, TLS, ServiceMonitor, network-policy, or operand rollout changes.

You can use the Cluster Bot Slack app to create a cluster from a pull request.

Send Cluster Bot a message in this format:

```text
launch ${OCP_VERSION},${PR_LINK} ${PLATFORM},${VARIANT}
```

For example:

```text
launch 4.23,https://github.com/openshift/cluster-openshift-apiserver-operator/pull/1234 aws,techpreview
```

For more options, send Cluster Bot:

```text
help
```

Once you have a cluster, set `KUBECONFIG` to the downloaded cluster credentials and run the relevant test target from the [Makefile](Makefile).

If existing CI is sufficient, mark the pull request with:

```text
/verified by ci
```

If you verified manually, mark it with:

```text
/verified by @<github-handle>
```

## Pull Request Expectations

Use the control-plane title convention:

```text
CNTRLPLANE-XXXX: short description
NO-JIRA: short description
```

Include:

- the problem being fixed
- the change being made
- why the approach is appropriate
- validation performed
- links to relevant Jira, enhancements, or upstream issues

Commits must be signed.

## Review

Reviewers and approvers are listed in [OWNERS](OWNERS) and [OWNERS_ALIASES](OWNERS_ALIASES).

Before requesting human review:

- address AI reviewer feedback
- respond to any comments you are not taking
- make sure the pull request description is current
