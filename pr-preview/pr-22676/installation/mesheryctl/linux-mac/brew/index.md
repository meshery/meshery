# Brew

> Install Meshery CLI on Linux or macOS with Homebrew

---

LLMS index: [llms.txt](/pr-preview/pr-22676/llms.txt)

---

# Install Meshery CLI with Homebrew



### Prerequisites

You need [Homebrew](https://brew.sh) installed on your **Linux** or **macOS** system.

### Install `mesheryctl` using Brew

To install `mesheryctl` with Homebrew, run:

<pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 $ brew install mesheryctl
 </div></div>
</pre>

You're ready to run Meshery. To start it on Kubernetes, run:

<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">
 $ mesheryctl system start

</div></div>
</pre>

To run Meshery on Docker instead, run:

<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">
 $ mesheryctl system start -p docker

</div></div>
</pre>

To customize the callback URL that your Provider redirects to after sign-in, set `MESHERY_SERVER_CALLBACK_URL` when you start Meshery:

<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">
 $ MESHERY_SERVER_CALLBACK_URL=https://custom-host mesheryctl system start

</div></div>
</pre>

On Kubernetes, `mesheryctl` uses the current context of your kubeconfig (`~/.kube/config` by default). To confirm that it points at the cluster you want Meshery to use, run `kubectl config get-contexts`.

If your kubeconfig has several contexts, switch to the one you want with `kubectl config use-context <context-to-use>`.

### Upgrade `mesheryctl` using Brew

To upgrade `mesheryctl`, run:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 $ brew upgrade mesheryctl
 </div></div>
 </pre>
