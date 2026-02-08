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

# Brew

## Prerequisites

You need to have `Brew` installed on your **Linux** or **macOS** system to perform these actions.

### Install `mesheryctl` using Brew

To install `mesheryctl` using homebrew, execute the following commands.

{{< code >}}
brew install mesheryctl
{{< /code >}}

You're ready to run Meshery. To do so, execute the following command.

{{< code >}}
mesheryctl system start
{{< /code >}}

If you are running Meshery on Docker, execute the following command.

{{< code >}}
mesheryctl system start -p docker
{{< /code >}}

Meshery server supports customizing authentication flow callback URL, which can be configured in the following way

{{< code >}}
MESHERY_SERVER_CALLBACK_URL=https://custom-host mesheryctl system start
{{< /code >}}

`mesheryctl` uses your current Kubernetes context, your KUBECONFIG environment variable (`~/.kube/config` by default). Confirm if this Kubernetes cluster you want Meshery to interact with by running the following command: `kubectl config get-contexts`.

If there are multiple contexts in your kubeconfig file, specify the one you want to use with the `use-context` subcommand: `kubectl config use-context <context-to-use>`.

### Upgrade `mesheryctl` using Brew

To upgrade `mesheryctl`, execute the following command.

{{< code >}}
brew upgrade mesheryctl
{{< /code >}}

<details>
<summary>
Example output of a successful upgrade.
</summary>

<pre><code>
➜  ~ brew upgrade mesheryctl
==> Upgrading 1 outdated package:
meshery/tap/mesheryctl 0.3.2 -> 0.3.4
==> Upgrading meshery/tap/mesheryctl
==> Downloading https://github.com/meshery/meshery/releases/download/v0.3.4/mesheryctl_0.3.4_Darwin_x86_64.zip
==> Downloading from https://github-production-release-asset-2e65be.s3.amazonaws.com/157554479/17522b00-2af0-11ea-8aef-cbfe8
######################################################################## 100.0%
🍺  /usr/local/Cellar/mesheryctl/0.3.4: 5 files, 10.2MB, built in 4 seconds
Removing: /usr/local/Cellar/mesheryctl/0.3.2... (5 files, 10.2MB)
Removing: /Users/lee/Library/Caches/Homebrew/mesheryctl--0.3.2.zip... (3.9MB)
==> Checking for dependents of upgraded formulae...
==> No dependents found!
</code></pre>
<br />
</details>

# Bash

To install or upgrade `mesheryctl` using `bash`, execute anyone of the following commands.

#### Option 1: Only install `mesheryctl` binary

{{< code >}}
curl -L https://meshery.io/install | DEPLOY_MESHERY=false bash -
{{< /code >}}

<br />
<br />

#### Option 2: Install `mesheryctl` binary and deploy Meshery on Docker

{{< code >}}
curl -L https://meshery.io/install | PLATFORM=docker bash -
{{< /code >}}

<br />
<br />

#### Option 3: Install `mesheryctl` binary and deploy Meshery on Kubernetes

{{< code >}}
curl -L https://meshery.io/install | PLATFORM=kubernetes bash -
{{< /code >}}

<br />
<br />

#### Option 4: Install `mesheryctl` binary and Meshery adapter(s)

Install `mesheryctl` binary and include one or more [adapters](/concepts/architecture/adapters) to be deployed

{{< code >}}
curl -L https://meshery.io/install | ADAPTERS=consul PLATFORM=kubernetes bash -
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

# Related Reading

## Meshery CLI Guides

Guides to using Meshery's various features and components.

- [Authenticating Meshery via CLI](/guides/mesheryctl/authenticate-with-meshery-via-cli)
- [Configuring Autocompletion for `mesheryctl`](/guides/mesheryctl/configuring-autocompletion-for-mesheryctl)
- [Running system checks using Meshery CLI](/guides/mesheryctl/running-system-checks-using-mesheryctl)
- [Mesheryctl system commands](/guides/mesheryctl/system-commands)
- [Using Meshery CLI](/guides/mesheryctl)
- [Upgrading Meshery CLI](/installation/upgrades#upgrading-meshery-cli)

{{< related-discussions tag="mesheryctl" >}}
