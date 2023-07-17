---
layout: default
title: Scoop
permalink: installation/platforms/scoop
type: installation
display-title: "false"
language: en
list: exclude
image: /assets/img/platforms/scoop.png
---

{% include installation_prerequisites.html %}

`mesheryctl` can be installed via Scoop (a package manager for Windows, just like apt for Ubuntu). `mesheryctl` is also available through Homebrew.

### Prerequisites

You need to have `scoop` installed on your Windows system to perform these actions.

### Install

To install `mesheryctl` using Scoop, execute the following commands.

<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">scoop bucket add mesheryctl https://github.com/layer5io/scoop-bucket.git
scoop install mesheryctl</div></div>
</pre>

You're ready to run Meshery. To do so, execute the following command.

<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">mesheryctl system start</div></div>
</pre>

### Upgrade

To upgrade `mesheryctl`, just execute the following command.

<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">scoop update mesheryctl</div></div>
</pre>
