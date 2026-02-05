---
title: "Install Meshery CLI on Linux or macOS"
description: "Install mesheryctl on Linux or macOS"
weight: 30
aliases:
  - /installation/platforms/linux-mac
image: /images/platforms/linux_mac.png
display_title: "false"
---

# Overview

To set up and run Meshery on Linux or macOS, you will need to install `mesheryctl`. `mesheryctl` is the command-line interface (CLI) for Meshery. It is used to install, manage, and operate one or more Meshery deployments. `mesheryctl` can be installed via `bash`, is also available [directly](https://github.com/meshery/meshery/releases/latest), or through [Homebrew](/installation/linux-mac/brew) or [Scoop](/installation/windows/scoop).

## Brew

Install `mesheryctl` using Homebrew:

{{< code >}}
brew install mesheryctl
{{< /code >}}

Then start Meshery:

{{< code >}}
mesheryctl system start
{{< /code >}}

## Bash

Install `mesheryctl` using Bash:

{{< code >}}
curl -L https://meshery.io/install | bash -
{{< /code >}}

Then start Meshery:

{{< code >}}
mesheryctl system start
{{< /code >}}

# Related Reading

## Meshery CLI Guides

Guides to using Meshery's various features and components.

- [Working with mesheryctl](/guides/mesheryctl/working-with-mesheryctl)
- [Upgrading Meshery CLI](/installation/upgrades#upgrading-meshery-cli)

{{< related-discussions tag="mesheryctl" >}}
