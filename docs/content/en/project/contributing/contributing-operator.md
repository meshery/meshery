---
title: Contributing to Meshery Operator
description: How to build, test, and deploy Meshery Operator from source
categories: [contributing]
---

[Meshery Operator]({{< ref "concepts/architecture/operator/index.md" >}}) is a
Kubernetes operator that manages the lifecycle of
[MeshSync]({{< ref "concepts/architecture/meshsync.md" >}}) and the
[Meshery Broker]({{< ref "concepts/architecture/broker/index.md" >}}). It is built with
[Kubebuilder](https://book.kubebuilder.io/) and the Operator SDK, and follows the standard
Kubebuilder `go/v4` project layout - the manager entrypoint lives at `cmd/main.go`.

Development follows the usual fork-and-pull-request workflow. Every commit must be signed
off; see the [Contributing Overview]({{< ref "project/contributing/_index.md" >}}) and the
[Git workflow guide]({{< ref "project/contributing/contributing-gitflow.md" >}}).

## Prerequisites

You only need Go and Docker installed locally. The `Makefile` downloads its pinned build
tools into `bin/` on demand (kustomize, controller-gen, setup-envtest, kind, operator-sdk,
opm), so their versions are reproducible and you do not install them globally. A running
Kubernetes cluster - or `make kind` to provision one with kind - is required for the deploy
and integration-test targets.

## Building and testing

Generate manifests and deepcopy code, then build the manager binary:

{{< code code=`make build` >}}

Run the controller against the cluster in your current `~/.kube/config`:

{{< code code=`make run` >}}

Run the unit and envtest suites. `make test` provisions an envtest control plane
automatically, so no cluster is required:

{{< code code=`make test` >}}

Lint the codebase (and auto-fix where possible with `make lint-fix`):

{{< code code=`make lint` >}}

## Building the container image

{{< code code=`make docker-build docker-push IMG=<registry>/meshery-operator:<tag>` >}}

## Deploying to a cluster

Install the CRDs and deploy the operator into your current kube context:

{{< code code=`make install
make deploy IMG=<registry>/meshery-operator:<tag>` >}}

To remove it:

{{< code code=`make undeploy
make uninstall` >}}

## Building the OLM bundle

To generate and build the Operator Lifecycle Manager bundle image:

{{< code code=`make bundle bundle-build bundle-push BUNDLE_IMG=<registry>/meshery-operator-bundle:<tag>` >}}

## Error handling

New errors returned from a controller or package must be declared as MeshKit structured
errors so they carry a stable code, severity, probable cause, and suggested remediation.
Follow the conventions in
[How to write MeshKit compatible errors]({{< ref "project/contributing/contributing-error.md" >}}) -
declare the code constant and factory function in an `error.go` file, and use
`errors.New(...)` from MeshKit rather than `fmt.Errorf` or the standard-library `errors`
package. Error names and codes must be unique across the whole component, so give each
constructor a distinct name (for example `ErrGettingBrokerResource` and
`ErrGettingMeshsyncResource`).

After adding or changing an error, run `make error` to validate that codes and names are
unique and to regenerate the error reference; `make error-util` assigns codes to new
placeholder constants and bumps `next_error_code` in `helpers/component_info.json`.
