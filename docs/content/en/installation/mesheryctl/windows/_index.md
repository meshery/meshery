---
title: Install Meshery CLI on Windows
categories: [mesheryctl]
aliases:
- /installation/platforms/windows
image: installation/mesheryctl/windows/images/wsl2.png
description: Install Meshery CLI on Windows
---

On Windows, install `mesheryctl` with Scoop or [download it directly](https://github.com/meshery/meshery/releases/latest) from the latest release.

{{% mesheryctl/installation-scoop %}}


## Install `mesheryctl` as a direct download

Download the `mesheryctl_<version>_Windows_x86_64.zip` archive from a [Meshery release](https://github.com/meshery/meshery/releases) and extract `mesheryctl.exe`. Patch releases do not always include `mesheryctl` archives, so if the latest release has none, use the most recent release that does. Then, from the folder that contains it, start Meshery on Kubernetes:

{{< code code="./mesheryctl system start" >}}

To run Meshery on Docker instead, run:

{{< code code="./mesheryctl system start -p docker" >}}

Optionally, move `mesheryctl.exe` to a folder in your `PATH` so you can run it from anywhere.


<!-- Meshery server supports customizing authentication flow callback URL, which can be configured in the following way
  <pre class="codeblock-pre">
  <div class="codeblock"><div class="clipboardjs">MESHERY_SERVER_CALLBACK_URL=https://custom-host ./mesheryctl system start</div></div>
  </pre>

Type `yes` when prompted to choose to configure a file. To get started, choose Docker as your platform to deploy Meshery. -->

{{< related-discussions tag="mesheryctl" >}}

### Installation Options
