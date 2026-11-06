# Bash

> Install Meshery CLI on Linux or MacOS with Bash

---

LLMS index: [llms.txt](/pr-preview/pr-22689/llms.txt)

---

# Install Meshery CLI with Bash


To install or upgrade `mesheryctl` with Bash, run one of the following commands.

#### Option 1: Only install `mesheryctl` binary

 <pre class="codeblock-pre">
 <div class="codeblock">
 <div class="clipboardjs">
  $ curl -L https://meshery.io/install | DEPLOY_MESHERY=false bash -
 </div></div>
 </pre>
<br />
<br />

#### Option 2: Install `mesheryctl` binary and deploy Meshery on Docker

 <pre class="codeblock-pre">
 <div class="codeblock">
 <div class="clipboardjs">
  $ curl -L https://meshery.io/install | PLATFORM=docker bash -
 </div></div>
 </pre>
<br />
<br />

#### Option 3: Install `mesheryctl` binary and deploy Meshery on Kubernetes

 <pre class="codeblock-pre">
 <div class="codeblock">
 <div class="clipboardjs">
  $ curl -L https://meshery.io/install | PLATFORM=kubernetes bash -
 </div></div>
 </pre>
<br />
<br />

#### Option 4: Install `mesheryctl` binary and Meshery adapter(s)

Install the `mesheryctl` binary and deploy Meshery together with one or more [adapters](/pr-preview/pr-22689/concepts/architecture/adapters/). Separate multiple adapters with spaces, for example `ADAPTERS="istio linkerd"`. `mesheryctl` can deploy the `consul`, `istio`, `linkerd`, `nsm`, `kuma`, `traefik`, and `cilium` adapters. The script also accepts `osm`, `nighthawk`, `nginx`, `tanzu`, and `app`, but `mesheryctl` has no deployment definition for those adapters, so `mesheryctl system start` rejects them. The script ignores any other value.

 <pre class="codeblock-pre">
 <div class="codeblock">
 <div class="clipboardjs">
  $ curl -L https://meshery.io/install | ADAPTERS=consul PLATFORM=kubernetes bash -
 </div></div>
 </pre>
<br />
<br />

### Start Meshery
If you installed only the binary (Option 1), start Meshery on Kubernetes with `mesheryctl`:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">mesheryctl system start</div></div>
 </pre>

To run Meshery on Docker instead, run:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">mesheryctl system start -p docker</div></div>
 </pre>
