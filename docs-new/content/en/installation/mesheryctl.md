---
title: "Install mesheryctl"
description: "Install Meshery CLI"
weight: 25
aliases:
  - /installation/mesheryctl/
  - /installation/platforms/mesheryctl
display_title: "true"
---

# Install mesheryctl

Meshery's command line client is `mesheryctl` and is the recommended tool for configuring and deploying one or more Meshery deployments. To install `mesheryctl` on your system, you may choose from any of the following supported methods.

`mesheryctl` can be installed via [Bash](/installation/linux-mac/bash), [Homebrew](/installation/linux-mac/brew), [Scoop](/installation/windows/scoop) or [directly downloaded](https://github.com/meshery/meshery/releases/latest).

> **Note:** Mesheryctl is configured for Kubernetes by default. To specify a different supported platform, use the `-p` flag.

## Install with Homebrew (macOS/Linux)

{{< code >}}
brew install mesheryctl
mesheryctl system start
{{< /code >}}

## Install with Bash (Linux/macOS)

{{< code >}}
curl -L https://meshery.io/install | bash -
mesheryctl system start
{{< /code >}}

## Install with Scoop (Windows)

{{< code >}}
scoop bucket add mesheryctl https://github.com/meshery/scoop-bucket.git
scoop install mesheryctl
mesheryctl system start
{{< /code >}}

## Direct Download

Download the latest release from [GitHub Releases](https://github.com/meshery/meshery/releases/latest).

Continue deploying Meshery onto one of the [Supported Platforms](/installation).

# Related Reading

## Meshery CLI Guides

Guides to using Meshery's various features and components.

- [Working with mesheryctl](/guides/mesheryctl/working-with-mesheryctl)
- [Upgrading Meshery CLI](/installation/upgrades#upgrading-meshery-cli)
- [Authenticating with Meshery via CLI](/guides/mesheryctl/authenticate-with-meshery-via-cli)

{{< related-discussions tag="mesheryctl" >}}
