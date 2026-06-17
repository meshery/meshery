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

# Related Reading

## Meshery CLI Guides

Guides to using Meshery's various features and components.

{{< mesheryctl-guides-list >}}

{{< related-discussions tag="mesheryctl" >}}
