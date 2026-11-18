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

`mesheryctl`, Meshery's command line client, is the recommended tool for configuring and running one or more Meshery deployments. Install it with [Bash]({{< ref "installation/mesheryctl/linux-mac/bash.md" >}}), [Homebrew]({{< ref "installation/mesheryctl/linux-mac/brew.md" >}}), or [Scoop]({{< ref "installation/mesheryctl/windows/scoop.md" >}}), or [download it directly](https://github.com/meshery/meshery/releases/latest) from the latest release.

{{% alert color="info" title="NOTE" %}} 
`mesheryctl` deploys Meshery to Kubernetes by default. To deploy to Docker instead, pass `-p docker` (or `--platform docker`) to `mesheryctl system start`.
{{% /alert %}}

# Install Meshery CLI with Brew

{{% mesheryctl/installation-brew %}}

# Install Meshery CLI with Bash

{{% mesheryctl/installation-bash %}}

# Install Meshery CLI with Scoop

{{% mesheryctl/installation-scoop %}}

Continue deploying Meshery onto one of the [Supported Platforms]({{< ref "installation/_index.md" >}}).

{{< related-discussions tag="mesheryctl" >}}
