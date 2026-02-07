---
title: "Bash"
description: "Install mesheryctl using Bash"
weight: 10
aliases:
  - /installation/platforms/bash
---

# Install mesheryctl using Bash

To install or upgrade `mesheryctl` using `bash`, execute anyone of the following commands.

#### Option 1: Only install `mesheryctl` binary

{{< code >}}
$ curl -L https://meshery.io/install | DEPLOY_MESHERY=false bash -
{{< /code >}}

<br />
<br />

#### Option 2: Install `mesheryctl` binary and deploy Meshery on Docker

{{< code >}}
$ curl -L https://meshery.io/install | PLATFORM=docker bash -
{{< /code >}}

<br />
<br />

#### Option 3: Install `mesheryctl` binary and deploy Meshery on Kubernetes

{{< code >}}
$ curl -L https://meshery.io/install | PLATFORM=kubernetes bash -
{{< /code >}}

<br />
<br />

#### Option 4: Install `mesheryctl` binary and Meshery adapter(s)

Install `mesheryctl` binary and include one or more [adapters](/concepts/architecture/adapters) to be deployed

{{< code >}}
$ curl -L https://meshery.io/install | ADAPTERS=consul PLATFORM=kubernetes bash -
{{< /code >}}

<br />
<br />

### Start Meshery

You are ready to deploy Meshery `mesheryctl`. To do so, execute the following command.

{{< code >}}
mesheryctl system start
{{< /code >}}

If you are running Meshery on Docker, execute the following command.

{{< code >}}
mesheryctl system start -p docker
{{< /code >}}
