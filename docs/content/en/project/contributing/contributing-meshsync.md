---
title: Contributing to MeshSync
description: How to build, run, test, and contribute to MeshSync from source
categories: [contributing]
---

[MeshSync]({{< ref "concepts/architecture/meshsync.md" >}}) is Meshery's event-driven,
continuous discovery and synchronization engine. It keeps Meshery Server's view of the
configuration and operational state of Kubernetes (and any supported Meshery platform)
current. Inside a cluster, MeshSync runs as a custom controller managed by
[Meshery Operator]({{< ref "concepts/architecture/operator/index.md" >}}) and publishes
resource changes over the [Meshery Broker]({{< ref "concepts/architecture/broker/index.md" >}})
(NATS). It is written in Go.

Development follows the usual fork-and-pull-request workflow. Every commit must be signed
off; see the [Contributing Overview]({{< ref "project/contributing/_index.md" >}}) and the
[Git workflow guide]({{< ref "project/contributing/contributing-gitflow.md" >}}).

## Prerequisites

- [Go](https://golang.org/dl/) - check `go.mod` for the required version.
- [Docker](https://docs.docker.com/get-docker/) - to build images and run a local NATS server.
- [golangci-lint](https://golangci-lint.run/welcome/install/) - the linter used by CI.
- For integration testing, a Kubernetes cluster; the harness uses [kind](https://kind.sigs.k8s.io/) and [kubectl](https://kubernetes.io/docs/tasks/tools/).

MeshSync uses a Make-driven workflow. Run `make` with no target to print every available
target and its description.

## Building and running

Build the MeshSync binary to `bin/meshsync`:

{{< code code=`make build` >}}

MeshSync runs in one of two modes; run `meshsync --help` for all input parameters.

**NATS mode** (default) connects to a NATS server and publishes Kubernetes resource updates
to a subject that the Meshery Broker forwards to Meshery Server - how MeshSync runs when
deployed alongside Meshery Operator and Broker. For local development, start a NATS server
and run MeshSync against it:

{{< code code=`make nats-run
make run` >}}

**File mode** runs MeshSync with no NATS or CRD dependency, writing Kubernetes resource
updates to disk as manifest-style YAML: an `-extended` snapshot of every event, and a
deduplicated snapshot with one entry per resource, keyed by `metadata.uid`.

## Testing

Run the unit tests with the race detector (linting runs first):

{{< code code=`make test` >}}

Generate an HTML coverage report (`cover.html`):

{{< code code=`make coverage-report` >}}

Integration tests validate end-to-end synchronization against a real cluster; the harness
provisions NATS and a kind cluster. Run the full cycle, or drive the phases individually:

{{< code code=`make integration-tests` >}}

{{< code code=`make integration-tests-setup
make integration-tests-run
make integration-tests-cleanup` >}}

Integration tests carry the keyword `Integration` in their name and are skipped under the
`--short` flag, which is how unit and integration runs are distinguished.

## Linting and dependencies

Lint the codebase with golangci-lint, and tidy Go module dependencies:

{{< code code=`make lint-run
make mod-tidy` >}}

## Building the container image

{{< code code=`make docker-build` >}}

## Error handling

New errors returned from MeshSync should be declared as MeshKit structured errors so they
carry a stable code, severity, probable cause, and suggested remediation. Follow the
conventions in
[How to write MeshKit compatible errors]({{< ref "project/contributing/contributing-error.md" >}}) -
declare the code constant and factory function in an `error.go` file and use
`errors.New(...)` from MeshKit rather than `fmt.Errorf` or the standard-library `errors`
package.
