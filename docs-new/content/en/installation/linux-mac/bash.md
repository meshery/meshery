---
title: "Bash"
description: "Install mesheryctl using Bash"
weight: 10
aliases:
  - /installation/platforms/bash
---

# Install mesheryctl using Bash

The fastest way to install `mesheryctl` on Linux or macOS is using the Bash installation script.

## Installation

Run the following command to install `mesheryctl`:

{{< code >}}
curl -L https://meshery.io/install | bash -
{{< /code >}}

## Start Meshery

After installing `mesheryctl`, you can start Meshery with:

{{< code >}}
mesheryctl system start
{{< /code >}}

## Verify Installation

Verify your installation by checking the version:

{{< code >}}
mesheryctl version
{{< /code >}}

{{< related-discussions tag="mesheryctl" >}}
