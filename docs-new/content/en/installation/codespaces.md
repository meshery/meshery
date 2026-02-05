---
title: "GitHub Codespaces"
description: "Install Meshery in GitHub Codespaces"
weight: 30
aliases:
  - /installation/platforms/codespaces
image: /images/platforms/codespaces.png
display_title: "false"
---

<h1>Quick Start with GitHub Codespaces <img src="/images/platforms/codespaces.png" style="width:35px;height:35px;" /></h1>

Deploy Meshery directly in GitHub Codespaces for a cloud-based development environment.

## Prerequisites

- A GitHub account with access to Codespaces
- A repository with a devcontainer configuration or use the Meshery repository

## Getting Started

1. Navigate to the [Meshery repository](https://github.com/meshery/meshery)
2. Click the **Code** button and select **Codespaces**
3. Click **Create codespace on master**

## Install Meshery

Once your Codespace is running, install mesheryctl:

{{< code >}}
curl -L https://meshery.io/install | bash -
{{< /code >}}

Start Meshery:

{{< code >}}
mesheryctl system start
{{< /code >}}

## Access Meshery UI

In Codespaces, Meshery UI will be available through the forwarded port. Check the **Ports** tab in VS Code to find the forwarded URL for port 9081.

{{< related-discussions tag="meshery" >}}
