--- 
layout: default
title: Scoop
permalink: installation/platforms/scoop
type: installation
display-title: "false"
language: en
list: include
---

{% include installation_prerequisites.html %}

Meshery can now be installed via Scoop (a package manager for Windows, just like apt for Ubuntu) You need to have `scoop` installed on your Windows to perform these actions.

### Install Meshery using Scoop
To install `mesheryctl` using Scoop, execute the following commands.
<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">
scoop bucket add mesheryctl https://github.com/layer5io/scoop-bucket.git
scoop install mesheryctl

</div></div>
</pre>

### Upgrading
To upgrade `mesheryctl`, just execute the following command.
<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">
scoop update mesheryctl

</div></div>
</pre>

