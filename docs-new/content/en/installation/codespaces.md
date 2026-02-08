---
title: "GitHub Codespaces"
description: "Build and contribute to Meshery using GitHub Codespaces"
weight: 30
aliases:
  - /installation/platforms/codespaces
image: /images/platforms/codespaces.png
display_title: "false"
---

<h1>Quick Start with GitHub Codespaces <img src="/images/platforms/codespaces.png" style="width:35px;height:35px;" /></h1>

Use Minikube in GitHub Codespaces to setup your development environment for Meshery.

<div class="prereqs">
  <h4>Prerequisites</h4>
  <ol>
    <li>Install the Meshery command line client, <a href="/installation/mesheryctl" class="meshery-light">mesheryctl</a>.</li>
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


## In-cluster Installation

Follow the steps below to install Meshery in your Minikube cluster.

## Preflight Checks

Read through the following considerations prior to deploying Meshery on Minikube.

### Preflight: Cluster Connectivity

You can develop and run Meshery in GitHub Codespaces using your choice of tool:

- A command shell, via an SSH connection initiated using GitHub CLI.
- One of the JetBrains IDEs, via the JetBrains Gateway.
- The Visual Studio Code desktop application.
- A browser-based version of Visual Studio Code.

> [!NOTE]
> **Choice of Codespaces Tool**  
> For the best experience, run Codespaces in your locally [installed IDE](https://docs.github.com/en/codespaces/developing-in-codespaces/developing-in-a-codespace). Alternatively, you can  
> [![Open in GitHub Codespaces](https://github.com/codespaces/badge.svg)](https://github.com/codespaces/new?hide_repo_select=true&ref=master&repo=157554479&machine=premiumLinux)

Start Minikube, if it is not already running, using the following command:

{{< code >}}
minikube start --cpus 4 --memory 4096
{{< /code >}}

Please allocate CPUs based on the machine you selected in GitHub Codespaces. To check the status of your Minikube cluster:

{{< code >}}
minikube status
{{< /code >}}

Verify your kubeconfig's current context.

{{< code >}}
kubectl cluster-info
{{< /code >}}

## Installation: Using `mesheryctl`

Use Meshery's CLI to streamline your connection to your Minikube cluster. Configure Meshery to connect to your Minikube cluster by executing:

{{< code >}}
mesheryctl system config minikube
{{< /code >}}

Once configured, execute the following command to start Meshery.

{{< code >}}
mesheryctl system start
{{< /code >}}

If you encounter any authentication issues, you can use `mesheryctl system login`. For more information, click [here](/guides/mesheryctl/authenticate-with-meshery-via-cli) to learn more.

## Installation: Using Helm

For detailed instructions on installing Meshery using Helm V3, please refer to the [Helm Installation](/installation/kubernetes/helm) guide.

## Installation: Manual Steps

You may also manually generate and load the kubeconfig file for Meshery to use:

**The following configuration yaml will be used by Meshery. Copy and paste the following in your config file**:

{{< code >}}
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: < cert shortcutted >
    server: https://192.168.99.100:8443
  name: minikube
contexts:
- context:
    cluster: minikube
    user: minikube
  name: minikube
current-context: minikube
kind: Config
preferences: {}
users:
- name: minikube
  user:
    client-certificate-data: < cert shortcutted >
    client-key-data: < key shortcutted >
{{< /code >}}

_Note_: Make sure _current-context_ is set to _minikube_.

**To allow Meshery to auto detect your config file, Run**:

{{< code >}}
kubectl config view --minify --flatten > config_minikube.yaml
{{< /code >}}

Meshery should now be connected with your managed Kubernetes instance. Take a look at the [Meshery guides](/guides) for advanced usage tips.

## Post-Installation Steps

Optionally, you can verify the health of your Meshery deployment using <a href='/reference/mesheryctl/system/check'>mesheryctl system check</a>.

You're ready to use Meshery! Open your browser and navigate to the Meshery UI.

{{< accessing-meshery-ui >}}

For further information to access Meshery UI or port forwarding in GitHub Codespaces, read the [docs](https://docs.github.com/en/codespaces/developing-in-a-codespace/forwarding-ports-in-your-codespace?tool=vscode)

{{< related-discussions tag="meshery" >}}
