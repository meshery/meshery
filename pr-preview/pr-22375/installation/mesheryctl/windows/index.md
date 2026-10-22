# Install Meshery CLI on Windows

> Install Meshery CLI on Windows

---

LLMS index: [llms.txt](/pr-preview/pr-22375/llms.txt)

---

On Windows, install `mesheryctl` with Scoop or [download it directly](https://github.com/meshery/meshery/releases/latest) from the latest release.



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



## Install `mesheryctl` as a direct download

Download the `mesheryctl_<version>_Windows_x86_64.zip` archive from a [Meshery release](https://github.com/meshery/meshery/releases) and extract `mesheryctl.exe`. Patch releases do not always include `mesheryctl` archives, so if the latest release has none, use the most recent release that does. Then, from the folder that contains it, start Meshery on Kubernetes:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">.\mesheryctl.exe system start</code>
	</div>
</pre>


To run Meshery on Docker instead, run:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">.\mesheryctl.exe system start -p docker</code>
	</div>
</pre>


Optionally, move `mesheryctl.exe` to a folder in your `PATH` so you can run it from anywhere.


<!-- Meshery server supports customizing authentication flow callback URL, which can be configured in the following way
  <pre class="codeblock-pre">
  <div class="codeblock"><div class="clipboardjs">MESHERY_SERVER_CALLBACK_URL=https://custom-host ./mesheryctl system start</div></div>
  </pre>

Type `yes` when prompted to choose to configure a file. To get started, choose Docker as your platform to deploy Meshery. -->

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

---

Section pages:

- [Scoop](/pr-preview/pr-22375/installation/mesheryctl/windows/scoop/): Install Meshery CLI on Windows with Scoop
