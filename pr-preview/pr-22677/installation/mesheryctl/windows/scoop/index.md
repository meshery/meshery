# Scoop

> Install Meshery CLI on Windows with Scoop

---

LLMS index: [llms.txt](/pr-preview/pr-22677/llms.txt)

---

# Install Meshery CLI with Scoop



## Prerequisites

You need [Scoop](https://scoop.sh) installed on your Windows system.

### Install `mesheryctl` with Scoop

`mesheryctl` is not in Scoop's default buckets. Add the Meshery bucket, then install `mesheryctl`:

<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">scoop bucket add mesheryctl https://github.com/meshery/scoop-bucket.git
scoop install mesheryctl</div></div>
</pre>

You're ready to run Meshery. To start it on Kubernetes, run:

<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">mesheryctl system start</div></div>
</pre>

To run Meshery on Docker instead, run:

<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">mesheryctl system start -p docker</div></div>
</pre>

### Upgrade `mesheryctl` with Scoop

To upgrade `mesheryctl`, run:

<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">scoop update mesheryctl</div></div>
</pre>
