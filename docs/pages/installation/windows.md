---
layout: default
title: Install Meshery CLI on Windows
permalink: installation/windows
type: installation
display-title: "true"
language: en
list: include
image: /assets/img/platforms/wsl2.png
---


On Windows systems, `mesheryctl` can be installed via Scoop or can be [downloaded directly](https://github.com/meshery/meshery/releases/latest).

{% include mesheryctl/installation-scoop.md %}

## Install `mesheryctl` as a direct download

Follow the [installation steps]({{ site.baseurl }}/installation#windows) to install the mesheryctl CLI. Then, execute:
<pre class="codeblock-pre">
<div class="codeblock"><div class="clipboardjs">./mesheryctl system start</div></div>
</pre>

Optionally, move the mesheryctl binary to a directory in your PATH.


<!-- Meshery server supports customizing authentication flow callback URL, which can be configured in the following way
  <pre class="codeblock-pre">
  <div class="codeblock"><div class="clipboardjs">MESHERY_SERVER_CALLBACK_URL=https://custom-host ./mesheryctl system start</div></div>
  </pre>

Type `yes` when prompted to choose to configure a file. To get started, choose Docker as your platform to deploy Meshery. -->