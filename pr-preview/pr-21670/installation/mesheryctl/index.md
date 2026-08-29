# Install using mesheryctl

> Use Meshery CLI to install Meshery on supported platforms.

Source: /pr-preview/pr-21670/installation/mesheryctl/

Meshery's command line client is `mesheryctl` and is the recommended tool for configuring and deploying one or more Meshery deployments. To install `mesheryctl` on your system, you may choose from any of the following supported methods.

`mesheryctl` can be installed via [bash](/pr-preview/pr-21670/installation/mesheryctl/linux-mac/bash/), [Homebrew](/pr-preview/pr-21670/installation/mesheryctl/linux-mac/brew/), [Scoop](/pr-preview/pr-21670/installation/mesheryctl/windows/scoop/) or [directly downloaded](https://github.com/meshery/meshery/releases/latest).

<div class="alert alert-info" role="alert"><div class="h4 alert-heading" role="heading">NOTE</div>

 
Mesheryctl is configured for Kubernetes by default. To specify a different supported platform, use the `-p` flag. 
</div>


# Install Meshery CLI with Brew



### Prerequisites

You need to have `Brew` installed on your **Linux** or **macOS** system to perform these actions.

### Install `mesheryctl` using Brew

To install `mesheryctl` using homebrew, execute the following commands.

<pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 $ brew install mesheryctl
 </div></div>
</pre>

You're ready to run Meshery. To do so, execute the following command.

<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">
 $ mesheryctl system start

</div></div>
</pre>

If you are running Meshery on Docker, execute the following command.

<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">
 $ mesheryctl system start -p docker

</div></div>
</pre>

Meshery server supports customizing authentication flow callback URL, which can be configured in the following way

<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">
 $ MESHERY_SERVER_CALLBACK_URL=https://custom-host mesheryctl system start

</div></div>
</pre>

`mesheryctl` uses your current Kubernetes context, your KUBECONFIG environment variable (`~/.kube/config` by default). Confirm if this Kubernetes cluster you want Meshery to interact with by running the following command: `kubectl config get-contexts`.

If there are multiple contexts in your kubeconfig file, specify the one you want to use with the `use-context` subcommand: `kubectl config use-context <context-to-use>`.

### Upgrade `mesheryctl` using Brew

To upgrade `mesheryctl`, execute the following command.

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 $ brew upgrade mesheryctl
 </div></div>
 </pre>

<details>
<summary>
Example output of a successful upgrade.
</summary>

<pre><code>
➜  ~ brew upgrade mesheryctl
==> Upgrading 1 outdated package:
meshery/tap/mesheryctl 0.3.2 -> 0.3.4
==> Upgrading meshery/tap/mesheryctl
==> Downloading https://github.com/meshery/meshery/releases/download/v0.3.4/mesheryctl_0.3.4_Darwin_x86_64.zip
==> Downloading from https://github-production-release-asset-2e65be.s3.amazonaws.com/157554479/17522b00-2af0-11ea-8aef-cbfe8
######################################################################## 100.0%
🍺  /usr/local/Cellar/mesheryctl/0.3.4: 5 files, 10.2MB, built in 4 seconds
Removing: /usr/local/Cellar/mesheryctl/0.3.2... (5 files, 10.2MB)
Removing: /Users/lee/Library/Caches/Homebrew/mesheryctl--0.3.2.zip... (3.9MB)
==> Checking for dependents of upgraded formulae...
==> No dependents found!
</code></pre>
<br />
</details>


# Install Meshery CLI with Bash


To install or upgrade `mesheryctl` using `bash`, execute anyone of the following commands.

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

Install `mesheryctl` binary and include one or more [adapters](/pr-preview/pr-21670/concepts/architecture/adapters/) to be deployed

 <pre class="codeblock-pre">
 <div class="codeblock">
 <div class="clipboardjs">
  $ curl -L https://meshery.io/install | ADAPTERS=consul PLATFORM=kubernetes bash -
 </div></div>
 </pre>
<br />
<br />

### Start Meshery
You are ready to deploy Meshery `mesheryctl`. To do so, execute the following command.

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">mesheryctl system start</div></div>
 </pre>

If you are running Meshery on Docker, execute the following command.

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">mesheryctl system start -p docker</div></div>
 </pre>


# Install Meshery CLI with Scoop



## Prerequisites

You need to have `scoop` installed on your Windows system to perform these actions.

### Install `mesheryctl` with Scoop

To install `mesheryctl` using Scoop, execute the following commands.

<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">scoop bucket add mesheryctl https://github.com/meshery/scoop-bucket.git
scoop install mesheryctl</div></div>
</pre>

You're ready to run Meshery. To do so, execute the following command.

<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">mesheryctl system start</div></div>
</pre>

If you are running Meshery on Docker, execute the following command.

<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">mesheryctl system start -p docker</div></div>
</pre>

### Upgrade `mesheryctl` with Scoop

To upgrade `mesheryctl`, execute the following command.

<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">scoop update mesheryctl</div></div>
</pre>


Continue deploying Meshery onto one of the [Supported Platforms](/pr-preview/pr-21670/installation/).

## Meshery CLI Guides

Guides to using Meshery's various features and components.

<ul>
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
      <li><a href="/pr-preview/pr-21670/guides/mesheryctl/authenticate-with-meshery-via-cli/">Authenticating Meshery via CLI</a></li>
    
  
    
  
    
  
    
  
    
      <li><a href="/pr-preview/pr-21670/guides/mesheryctl/configuring-autocompletion-for-mesheryctl/">Configuring Autocompletion for `mesheryctl`</a></li>
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
      <li><a href="/pr-preview/pr-21670/guides/mesheryctl/system-commands/">Mesheryctl system commands</a></li>
    
  
    
  
    
  
    
  
    
      <li><a href="/pr-preview/pr-21670/guides/mesheryctl/running-system-checks-using-mesheryctl/">Running system checks using Meshery CLI</a></li>
    
  
    
  
    
  
    
  
    
  
    
      <li><a href="/pr-preview/pr-21670/guides/mesheryctl/working-with-mesheryctl/">Using Meshery CLI</a></li>
    
  
    
  
    
  
  <li><a href="/pr-preview/pr-21670/installation/upgrades/#upgrading-meshery-cli">Upgrading Meshery CLI</a></li>
</ul>


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
