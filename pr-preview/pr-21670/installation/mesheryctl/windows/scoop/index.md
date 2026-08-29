# Scoop

> Install Meshery CLI on Windows with Scoop

Source: /pr-preview/pr-21670/installation/mesheryctl/windows/scoop/

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


## Mesheryctl Guides

Guides to using Meshery's various features and components.

<ul>
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
      <li><a href="/pr-preview/pr-21670/guides/mesheryctl/authenticate-with-meshery-via-cli/">Authenticating Meshery via CLI</a></li>
    
  
    
  
    
  
    
  
    
      <li><a href="/pr-preview/pr-21670/guides/mesheryctl/configuring-autocompletion-for-mesheryctl/">Configuring Autocompletion for `mesheryctl`</a></li>
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
  
    
      <li><a href="/pr-preview/pr-21670/guides/mesheryctl/system-commands/">Mesheryctl system commands</a></li>
    
  
    
  
    
  
    
  
    
      <li><a href="/pr-preview/pr-21670/guides/mesheryctl/running-system-checks-using-mesheryctl/">Running system checks using Meshery CLI</a></li>
    
  
    
  
    
  
    
  
    
  
    
      <li><a href="/pr-preview/pr-21670/guides/mesheryctl/working-with-mesheryctl/">Using Meshery CLI</a></li>
    
  
    
  
    
  
  <li><a href="/pr-preview/pr-21670/installation/upgrades/#upgrading-meshery-cli">Upgrading Meshery CLI</a></li>
</ul>
