---
title: Install using mesheryctl
categories: [mesheryctl]
aliases:
- /installation/mesheryctl/
- /installation/platforms/mesheryctl
suggested-reading: false
description: Use Meshery CLI to install Meshery on supported platforms.
weight: 5
---

Meshery's command line client is `mesheryctl` and is the recommended tool for configuring and deploying one or more Meshery deployments. To install `mesheryctl` on your system, you may choose from any of the following supported methods.

`mesheryctl` can be installed via [bash]({{< ref "installation/mesheryctl/linux-mac/bash.md" >}}), [Homebrew]({{< ref "installation/mesheryctl/linux-mac/brew.md" >}}), [Scoop]({{< ref "installation/mesheryctl/windows/scoop.md" >}}) or [directly downloaded](https://github.com/meshery/meshery/releases/latest).

{{% alert color="info" title="NOTE" %}} 
Mesheryctl is configured for Kubernetes by default. To specify a different supported platform, use the `-p` flag. 
{{% /alert %}}

# Install Meshery CLI with Brew

{{% mesheryctl/installation-brew %}}

# Install Meshery CLI with Bash

{{% mesheryctl/installation-bash %}}

# Install Meshery CLI with Scoop

{{% mesheryctl/installation-scoop %}}

Continue deploying Meshery onto one of the [Supported Platforms]({{< ref "installation/_index.md" >}}).

# Uninstalling `mesheryctl`

To uninstall `mesheryctl`, use the uninstall method corresponding to your package manager (e.g., `brew uninstall mesheryctl` or `scoop uninstall mesheryctl`).

If you installed `mesheryctl` via direct download, remove `mesheryctl` by deleting the binary file from wherever it is located on your `PATH`.

{{% alert color="info" title="NOTE" %}}
Running `mesheryctl system stop` removes a Meshery deployment (containers and platform resources) from your system. This is distinct from removing the `mesheryctl` CLI binary itself.
{{% /alert %}}


{{< related-discussions tag="mesheryctl" >}}

