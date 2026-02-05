---
title: "Upgrading Meshery"
description: "How to upgrade Meshery and mesheryctl"
weight: 50
aliases:
  - /installation/platforms/upgrades
---

# Upgrading Meshery

This guide covers upgrading both Meshery (the management plane) and mesheryctl (the CLI).

<p style="text-align:center">
<a href="/images/architecture/upgrading-meshery.svg">
    <img src="/images/architecture/upgrading-meshery.svg" style="margin: 1rem;" />
</a><br /><i><small>Figure: Meshery components</small></i>
</p>

## Upgrading Meshery CLI

### Using Homebrew (macOS/Linux)

{{< code >}}
brew upgrade mesheryctl
{{< /code >}}

### Using Scoop (Windows)

{{< code >}}
scoop update mesheryctl
{{< /code >}}

### Using Bash

Re-run the installation script to get the latest version:

{{< code >}}
curl -L https://meshery.io/install | bash -
{{< /code >}}

### Direct Download

Download the latest release from [GitHub Releases](https://github.com/meshery/meshery/releases/latest).

## Upgrading Meshery Server

### Kubernetes (Helm)

{{< code >}}
helm repo update
helm upgrade meshery meshery/meshery --namespace meshery
{{< /code >}}

### Docker

Stop the current instance and start with the latest version:

{{< code >}}
mesheryctl system stop
mesheryctl system start -p docker
{{< /code >}}

### Using mesheryctl

{{< code >}}
mesheryctl system update
{{< /code >}}

## Verify Upgrade

Check the version after upgrading:

{{< code >}}
mesheryctl version
{{< /code >}}

Verify system health:

{{< code >}}
mesheryctl system check
{{< /code >}}

{{< related-discussions tag="meshery" >}}
