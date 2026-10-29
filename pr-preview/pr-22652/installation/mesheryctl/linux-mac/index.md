# Install Meshery CLI on Linux or macOS

> Install Meshery CLI on Linux or macOS

---

LLMS index: [llms.txt](/pr-preview/pr-22652/llms.txt)

---

# Overview

To set up and run Meshery on Linux or macOS, install `mesheryctl`, Meshery's command line interface. You use it to install, manage, and operate one or more Meshery deployments. Install `mesheryctl` with [Homebrew](/pr-preview/pr-22652/installation/mesheryctl/linux-mac/brew/) or [Bash](/pr-preview/pr-22652/installation/mesheryctl/linux-mac/bash/), or [download it directly](https://github.com/meshery/meshery/releases/latest) from the latest release.

# Brew



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



# Bash


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

Install the `mesheryctl` binary and deploy Meshery together with one or more [adapters](/pr-preview/pr-22652/concepts/architecture/adapters/). Separate multiple adapters with spaces, for example `ADAPTERS="istio linkerd"`. `mesheryctl` can deploy the `consul`, `istio`, `linkerd`, `nsm`, `kuma`, `traefik`, and `cilium` adapters. The script also accepts `osm`, `nighthawk`, `nginx`, `tanzu`, and `app`, but `mesheryctl` has no deployment definition for those adapters, so `mesheryctl system start` rejects them. The script ignores any other value.

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


<div class="related-discussions">
  <h3>Recent Discussions with "mesheryctl" Tag</h3><ul><li>
          <a href="https://discuss.meshery.io/t/pairing-up-with-a-meshmate-to-gacefully-start-contributing-to-meshery-and-its-project-ecosystem/6957" target="_blank" rel="noopener noreferrer">
            Pairing up with a Meshmate to gacefully start contributing to Meshery and its project ecosystem
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/development-meeting-contributing-to-meshery-cli-april-30-2025/6927" target="_blank" rel="noopener noreferrer">
            [Development Meeting] Contributing to Meshery CLI - April 30, 2025
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/meshery-development-meeting-april-30th-2025/6926" target="_blank" rel="noopener noreferrer">
            Meshery Development Meeting | April 30th, 2025
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/newcomers-meeting-end-to-end-testing-in-meshery-cli-using-bats-april-17-2025/6897" target="_blank" rel="noopener noreferrer">
            [Newcomers’ Meeting] End-to-End Testing in Meshery CLI using BATs – April 17, 2025
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/unable-to-lint-mesheryctl/6854" target="_blank" rel="noopener noreferrer">
            unable to lint mesheryctl 
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/mesheryctl-system-login-problem/6687" target="_blank" rel="noopener noreferrer">
            Mesheryctl system login problem
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/error-while-launching-the-meshery-dashboard/6600" target="_blank" rel="noopener noreferrer">
            Error while launching the Meshery Dashboard
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/cant-find-the-file-path-for-meshery-designs/6319" target="_blank" rel="noopener noreferrer">
            Can&#39;t find the file path for meshery designs
          </a>
        </li></ul><p>
    <a href="https://discuss.meshery.io/tag/mesheryctl" target="_blank" rel="noopener noreferrer">
      View all discussions tagged with <code>mesheryctl</code>
    </a>
  </p>
</div>


### Installation Options
<!-- 
1. You can either use **Bash** or **Brew** to install <a href="/pr-preview/pr-22652/guides/mesheryctl/">mesheryctl</a> ( Meshery command line interface ).
2. To run **Meshery**, execute the following command.

   <pre class="codeblock-pre"><div class="codeblock">
   <div class="clipboardjs">mesheryctl system start</div></div>
   </pre>

Meshery server supports customizing authentication flow callback URL, which can be configured in the following way

<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">
 $ MESHERY_SERVER_CALLBACK_URL=https://custom-host mesheryctl system start

</div></div>
</pre>
-->

---

Section pages:

- [Bash](/pr-preview/pr-22652/installation/mesheryctl/linux-mac/bash/): Install Meshery CLI on Linux or MacOS with Bash
- [Brew](/pr-preview/pr-22652/installation/mesheryctl/linux-mac/brew/): Install Meshery CLI on Linux or macOS with Homebrew
