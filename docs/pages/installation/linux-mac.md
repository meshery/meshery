---
layout: default
title: Linux or Mac
permalink: installation/linux-mac
type: installation
display-title: "false"
language: en
list: include
image: /assets/img/platforms/linux_mac.png 
---

# Overview

To set up and run Meshery on Linux or macOS, you will need to install `mesheryctl`. `mesheryctl` is the command line interface (CLI) for Meshery. It is used to install, manage, and operate one or more Meshery deployments. `mesheryctl` can be installed via `bash` is also available [directly](https://github.com/meshery/meshery/releases/latest) or through [Homebrew]({{site.baseurl}}/installation/linux-mac/brew) or [Scoop]({{site.baseurl}}/installation/windows/scoop).

# Brew

{% include mesheryctl/installation-brew.md %}

# Bash

{% include mesheryctl/installation-bash.md %}

<!-- 
1. You can either use **Bash** or **Brew** to install <a href="/guides/mesheryctl">mesheryctl</a> ( Meshery command line interface ).
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