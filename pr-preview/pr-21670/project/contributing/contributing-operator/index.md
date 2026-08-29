# Contributing to Meshery Operator

> How to build, test, and deploy Meshery Operator from source

Source: /pr-preview/pr-21670/project/contributing/contributing-operator/

[Meshery Operator](/pr-preview/pr-21670/concepts/architecture/operator/) is a
Kubernetes operator that manages the lifecycle of
[MeshSync](/pr-preview/pr-21670/concepts/architecture/meshsync/) and the
[Meshery Broker](/pr-preview/pr-21670/concepts/architecture/broker/). It is built with
[Kubebuilder](https://book.kubebuilder.io/) and the Operator SDK, and follows the standard
Kubebuilder `go/v4` project layout - the manager entrypoint lives at `cmd/main.go`.

Development follows the usual fork-and-pull-request workflow. Every commit must be signed
off; see the [Contributing Overview](/pr-preview/pr-21670/project/contributing/) and the
[Git workflow guide](/pr-preview/pr-21670/project/contributing/contributing-gitflow/).

## Prerequisites

You only need Go and Docker installed locally. The `Makefile` downloads its pinned build
tools into `bin/` on demand (kustomize, controller-gen, setup-envtest, kind, operator-sdk,
opm), so their versions are reproducible and you do not install them globally. A running
Kubernetes cluster - or `make kind` to provision one with kind - is required for the deploy
and integration-test targets.

## Building and testing

Generate manifests and deepcopy code, then build the manager binary:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">make build</code>
	</div>
</pre>


Run the controller against the cluster in your current `~/.kube/config`:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">make run</code>
	</div>
</pre>


Run the unit and envtest suites. `make test` provisions an envtest control plane
automatically, so no cluster is required:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">make test</code>
	</div>
</pre>


Lint the codebase (and auto-fix where possible with `make lint-fix`):



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">make lint</code>
	</div>
</pre>


## Building the container image



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">make docker-build docker-push IMG=&lt;registry&gt;/meshery-operator:&lt;tag&gt;</code>
	</div>
</pre>


## Deploying to a cluster

Install the CRDs and deploy the operator into your current kube context:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">make install
make deploy IMG=&lt;registry&gt;/meshery-operator:&lt;tag&gt;</code>
	</div>
</pre>


To remove it:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">make undeploy
make uninstall</code>
	</div>
</pre>


## Building the OLM bundle

To generate and build the Operator Lifecycle Manager bundle image:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">make bundle bundle-build bundle-push BUNDLE_IMG=&lt;registry&gt;/meshery-operator-bundle:&lt;tag&gt;</code>
	</div>
</pre>


## Error handling

New errors returned from a controller or package must be declared as MeshKit structured
errors so they carry a stable code, severity, probable cause, and suggested remediation.
Follow the conventions in
[How to write MeshKit compatible errors](/pr-preview/pr-21670/project/contributing/contributing-error/) -
declare the code constant and factory function in an `error.go` file, and use
`errors.New(...)` from MeshKit rather than `fmt.Errorf` or the standard-library `errors`
package. Error names and codes must be unique across the whole component, so give each
constructor a distinct name (for example `ErrGettingBrokerResource` and
`ErrGettingMeshsyncResource`).

After adding or changing an error, run `make error` to validate that codes and names are
unique and to regenerate the error reference; `make error-util` assigns codes to new
placeholder constants and bumps `next_error_code` in `helpers/component_info.json`.
