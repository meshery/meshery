---
title: Install Meshery CLI on Windows
description: "Install mesheryctl on Windows"
weight: 40
aliases:
  - /installation/platforms/windows
image: /images/platforms/wsl2.png
display_title: true
---

On Windows systems, `mesheryctl` can be installed via Scoop or [downloaded directly](https://github.com/meshery/meshery/releases/latest).

## Prerequisites

You need to have `scoop` installed on your Windows system to perform these actions.

### Install `mesheryctl` with Scoop

To install `mesheryctl` using Scoop, execute the following commands.

{{< code >}}
scoop bucket add mesheryctl https://github.com/meshery/scoop-bucket.git
scoop install mesheryctl
{{< /code >}}

You're ready to run Meshery. To do so, execute the following command.

{{< code >}}
mesheryctl system start
{{< /code >}}

If you are running Meshery on Docker, execute the following command.

{{< code >}}
mesheryctl system start -p docker
{{< /code >}}

### Upgrade `mesheryctl` with Scoop

To upgrade `mesheryctl`, execute the following command.

{{< code >}}
scoop update mesheryctl
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

- [Authenticating Meshery via CLI](/guides/mesheryctl/authenticate-with-meshery-via-cli)
- [Configuring Autocompletion for `mesheryctl`](/guides/mesheryctl/configuring-autocompletion-for-mesheryctl)
- [Running system checks using Meshery CLI](/guides/mesheryctl/running-system-checks-using-mesheryctl)
- [Mesheryctl system commands](/guides/mesheryctl/system-commands)
- [Using Meshery CLI](/guides/mesheryctl/working-with-mesheryctl)
- [Upgrading Meshery CLI](/installation/upgrades#upgrading-meshery-cli)

{{< related-discussions tag="mesheryctl" >}}
