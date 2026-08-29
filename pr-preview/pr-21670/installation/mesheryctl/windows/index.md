# Install Meshery CLI on Windows

> Install Meshery CLI on Windows

Source: /pr-preview/pr-21670/installation/mesheryctl/windows/

On Windows systems, `mesheryctl` can be installed via Scoop or [downloaded directly](https://github.com/meshery/meshery/releases/latest).



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



## Install `mesheryctl` as a direct download

Follow the installation steps to install the `mesheryctl` CLI. Then, execute:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">./mesheryctl system start</code>
	</div>
</pre>


If you are installing Meshery on Docker, execute the following command:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">./mesheryctl system start -p docker</code>
	</div>
</pre>


Optionally, move the `mesheryctl` binary to a directory in your `PATH`.


<!-- Meshery server supports customizing authentication flow callback URL, which can be configured in the following way
  <pre class="codeblock-pre">
  <div class="codeblock"><div class="clipboardjs">MESHERY_SERVER_CALLBACK_URL=https://custom-host ./mesheryctl system start</div></div>
  </pre>

Type `yes` when prompted to choose to configure a file. To get started, choose Docker as your platform to deploy Meshery. -->

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


### Installation Options
