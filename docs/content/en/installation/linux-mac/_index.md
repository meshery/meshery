---
title: Install Meshery CLI on Linux or macOS
categories: [mesheryctl]
aliases:
- /installation/platforms/linux-mac
display_title: false
image: /installation/linux-mac/images/linux_mac.png
description: Install Meshery CLI on Linux or macOS
---

# Overview

To set up and run Meshery on Linux or macOS, you will need to install `mesheryctl`. `mesheryctl` is the command-line interface (CLI) for Meshery. It is used to install, manage, and operate one or more Meshery deployments. `mesheryctl` can be installed via `bash`, is also available [directly](https://github.com/meshery/meshery/releases/latest), or through [Homebrew](/installation/linux-mac/brew) or [Scoop](/installation/windows/scoop).

# Brew

{{% mesheryctl/installation-brew %}}

# Bash

{{% mesheryctl/installation-bash %}}

# Related Reading

## Meshery CLI Guides

Guides to using Meshery's various features and components.

{{< mesheryctl-guides-list >}}

{{< related-discussions tag="mesheryctl" >}}

### Installation Options
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