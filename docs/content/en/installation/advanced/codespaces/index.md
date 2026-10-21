---
title: Codespaces
categories: [kubernetes]
aliases:
- /installation/platforms/codespaces
# display_title: false
image: ./images/codespaces.png
description: Build and contribute to Meshery using GitHub Codespaces
---

Use Minikube in GitHub Codespaces to set up a development environment for Meshery.

<div class="prereqs"><p><strong style="font-size: 20px;">Prerequisites</strong> </p> 
  <ol>
    <li>Install the Meshery command line client, <a href="{{< ref "installation/mesheryctl/_index.md" >}}" class="meshery-light">mesheryctl</a>.</li>
  </ol>
</div>

## Available Deployment Methods

- [In-cluster Installation](#in-cluster-installation)
  - [Preflight Checks](#preflight-checks)
    - [Preflight: Cluster Connectivity](#preflight-cluster-connectivity)
  - [Installation: Using `mesheryctl`](#installation-using-mesheryctl)
  - [Installation: Using Helm](#installation-using-helm)
  - [Installation: Manual Steps](#installation-manual-steps)
- [Post-Installation Steps](#post-installation-steps)

# In-cluster Installation

Follow these steps to install Meshery in your Minikube cluster.

## Preflight Checks

Before you deploy Meshery on Minikube, check the following.

### Preflight: Cluster Connectivity


You can develop and run Meshery in GitHub Codespaces using your choice of tool:

- A command shell, via an SSH connection initiated using GitHub CLI.
- One of the JetBrains IDEs, via the JetBrains Gateway.
- The Visual Studio Code desktop application.
- A browser-based version of Visual Studio Code.

{{% alert color="dark" title="Choice of Codespaces Tool" %}}For the best experience, run Codespaces in your locally <a href='https://docs.github.com/en/codespaces/developing-in-codespaces/developing-in-a-codespace'>installed IDE</a>. Alternatively, you can <br /><a href='https://github.com/codespaces/new?hide_repo_select=true&ref=master&repo=157554479&machine=premiumLinux'><img alt='Open in GitHub Codespaces' src='https://github.com/codespaces/badge.svg' /></a>{{% /alert %}}

Start Minikube if it is not already running:

{{< code code="minikube start --cpus 4 --memory 4096" >}}

Set `--cpus` and `--memory` to fit the machine type you chose in GitHub Codespaces. To check the status of your Minikube cluster:

{{< code code="minikube status" >}}

Verify that `kubectl` can reach the cluster:

{{< code code="kubectl cluster-info" >}}

## Installation: Using `mesheryctl`

Start Meshery in your Minikube cluster:

{{< code code="mesheryctl system start -p kubernetes" >}}

Meshery Server running in the cluster connects to it automatically. If the connection does not appear, sign in with `mesheryctl system login` and create it with `mesheryctl connection create --type minikube`, which uploads your kubeconfig to Meshery Server. (`mesheryctl system config minikube` is deprecated.) For more on signing in, see [Authenticate with Meshery via CLI]({{< ref "guides/mesheryctl/authenticate-with-meshery-via-cli/index.md" >}}).

## Installation: Using Helm

To install Meshery with Helm 3, see the [Helm Installation]({{< ref "installation/kubernetes/helm.md" >}}) guide.

## Installation: Manual Steps

You can also generate a self-contained kubeconfig for your Minikube cluster and upload it to Meshery yourself. Export the `minikube` context, with its certificates embedded:

{{< code code="kubectl config view --minify --flatten --context minikube > config_minikube.yaml" >}}

The file looks like this, with your own server address and certificate data:

{{< code code=`apiVersion: v1
kind: Config
clusters:
- cluster:
    certificate-authority-data: <certificate data>
    server: https://192.168.49.2:8443
  name: minikube
contexts:
- context:
    cluster: minikube
    user: minikube
  name: minikube
current-context: minikube
preferences: {}
users:
- name: minikube
  user:
    client-certificate-data: <certificate data>
    client-key-data: <key data>` >}}

_Note_: Make sure that `current-context` is `minikube`.

Then upload `config_minikube.yaml` in Meshery UI: go to **Lifecycle** > **Connections**, click **Create Connection**, choose **Kubernetes**, and upload the file in the **Import Kubeconfig** step. See [Registering a Connection]({{< ref "guides/infrastructure-management/registering-a-connection.md" >}}). Take a look at the [Meshery guides]({{< ref "guides/_index.md" >}}) for advanced usage tips.

# Post-Installation Steps

Optionally, verify the health of your Meshery deployment with <a href='{{< ref "reference/references/mesheryctl/system/check.md" >}}'>mesheryctl system check</a>.

You're ready to use Meshery. Open Meshery UI in your browser.

{{< installation/accessing-meshery-ui >}}

To reach Meshery UI from your browser, forward its port from the codespace. See [Forwarding ports in your codespace](https://docs.github.com/en/codespaces/developing-in-a-codespace/forwarding-ports-in-your-codespace?tool=vscode).

{{< related-discussions tag="meshery" >}}