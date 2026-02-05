---
title: "Install Meshery CLI on Windows"
description: "Install mesheryctl on Windows"
weight: 40
aliases:
  - /installation/platforms/windows
image: /images/platforms/wsl2.png
display_title: "true"
---

On Windows systems, `mesheryctl` can be installed via Scoop or [downloaded directly](https://github.com/meshery/meshery/releases/latest).

## Install using Scoop

{{< code >}}
scoop bucket add mesheryctl https://github.com/meshery/scoop-bucket.git
scoop install mesheryctl
{{< /code >}}

## Install `mesheryctl` as a direct download

Follow the installation steps to install the `mesheryctl` CLI. Then, execute:

{{< code >}}
./mesheryctl system start
{{< /code >}}

If you are installing Meshery on Docker, execute the following command:

{{< code >}}
./mesheryctl system start -p docker
{{< /code >}}

Optionally, move the `mesheryctl` binary to a directory in your `PATH`.

# Related Reading

## Meshery CLI Guides

Guides to using Meshery's various features and components.

- [Working with mesheryctl](/guides/mesheryctl/working-with-mesheryctl)
- [Upgrading Meshery CLI](/installation/upgrades#upgrading-meshery-cli)

{{< related-discussions tag="mesheryctl" >}}
