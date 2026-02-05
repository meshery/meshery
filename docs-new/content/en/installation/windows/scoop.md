---
title: "Scoop"
description: "Install mesheryctl using Scoop on Windows"
weight: 10
aliases:
  - /installation/platforms/scoop
---

# Install mesheryctl using Scoop

Install `mesheryctl` using Scoop on Windows.

## Installation

First, add the mesheryctl bucket:

{{< code >}}
scoop bucket add mesheryctl https://github.com/meshery/scoop-bucket.git
{{< /code >}}

Then install mesheryctl:

{{< code >}}
scoop install mesheryctl
{{< /code >}}

## Start Meshery

After installing `mesheryctl`, you can start Meshery with:

{{< code >}}
mesheryctl system start
{{< /code >}}

## Upgrade mesheryctl

To upgrade to the latest version:

{{< code >}}
scoop update mesheryctl
{{< /code >}}

## Verify Installation

Verify your installation by checking the version:

{{< code >}}
mesheryctl version
{{< /code >}}

{{< related-discussions tag="mesheryctl" >}}
