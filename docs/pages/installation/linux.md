---
layout: default
title: Linux or Mac
permalink: installation/linux
type: installation
display-title: "false"
language: en
list: include
image: /assets/img/platforms/linux_mac.png 
---

{% include installation_prerequisites.html %}

# Overview
To set up and run Meshery on Linux or Mac

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

Stuck at another error? [Tell us about it](https://slack.meshery.io/)