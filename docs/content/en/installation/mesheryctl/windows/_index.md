---
title: Install Meshery CLI on Windows
categories: [mesheryctl]
aliases:
- /installation/platforms/windows
image: installation/mesheryctl/windows/images/wsl2.png
description: Install Meshery CLI on Windows
---

On Windows systems, `mesheryctl` can be installed via Scoop or [downloaded directly](https://github.com/meshery/meshery/releases/latest).

{{% mesheryctl/installation-scoop %}}


## Install `mesheryctl` as a direct download

Follow the installation steps to install the `mesheryctl` CLI. Then, execute:

{{< code code="./mesheryctl system start" >}}

If you are installing Meshery on Docker, execute the following command:

{{< code code="./mesheryctl system start -p docker" >}}

Optionally, move the `mesheryctl` binary to a directory in your `PATH`.


## Uninstalling `mesheryctl`

To uninstall `mesheryctl`, use `scoop uninstall mesheryctl` if installed via Scoop. If you installed `mesheryctl` via direct download, remove `mesheryctl` by deleting the binary file from wherever it is located on your `PATH`.

{{% alert color="info" title="NOTE" %}}
Running `mesheryctl system stop` removes a Meshery deployment from your system. This is distinct from removing the `mesheryctl` CLI binary itself.
{{% /alert %}}


<!-- Meshery server supports customizing authentication flow callback URL, which can be configured in the following way

  <pre class="codeblock-pre">
  <div class="codeblock"><div class="clipboardjs">MESHERY_SERVER_CALLBACK_URL=https://custom-host ./mesheryctl system start</div></div>
  </pre>

Type `yes` when prompted to choose to configure a file. To get started, choose Docker as your platform to deploy Meshery. -->

{{< related-discussions tag="mesheryctl" >}}

### Installation Options
