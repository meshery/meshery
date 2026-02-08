---
title: "Brew"
description: "Install mesheryctl using Homebrew"
weight: 20
aliases:
  - /installation/platforms/brew
---

# Install Meshery CLI with Homebrew

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
