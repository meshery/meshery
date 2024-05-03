---
layout: default
title: Codespaces
permalink: installation/codespaces
type: installation
category: kubernetes
redirect_from:
- installation/platforms/codespaces
display-title: "false"
language: en
list: include
image: /assets/img/platforms/codespaces.png
abstract: Build and contribute to Meshery using GitHub Codespaces
---

<h1>Quick Start with {{ page.title }} <img src="{{ page.image }}" style="width:35px;height:35px;" /></h1>

Use Minikube in GitHub Codespace to setup your development environment for Meshery.

<div class="prereqs"><p><strong style="font-size: 20px;">Prerequisites</strong> </p> 
  <ol>
    <li>Install the Meshery command line client, <a href="{{ site.baseurl }}/installation/mesheryctl" class="meshery-light">mesheryctl</a>.</li>
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

Follow the steps below to install Meshery in your Minikube cluster.

## Preflight Checks

Read through the following considerations prior to deploying Meshery on Minikube.

### Preflight: Cluster Connectivity


You can develop and run Meshery in a GitHub Codespace using your choice of tool:

- A command shell, via an SSH connection initiated using GitHub CLI.
- One of the JetBrains IDEs, via the JetBrains Gateway.
- The Visual Studio Code desktop application.
- A browser-based version of Visual Studio Code.

{% include alert.html type="dark" title="Choice of Codespace Tool" content="For the best experience, run Codespace in your locally <a href='https://docs.github.com/en/codespaces/developing-in-codespaces/developing-in-a-codespace'>installed IDE</a>. Alternatively, you can <br /><a href='https://github.com/codespaces/new?hide_repo_select=true&ref=master&repo=157554479&machine=premiumLinux'><img alt='Open in GitHub Codespaces' src='https://github.com/codespaces/badge.svg' /></a>" %}

Start the minikube, if not started using the following command:
{% capture code_content %}minikube start --cpus 4 --memory 4096{% endcapture %}
{% include code.html code=code_content %}
Please allocate cpus based on the machine you selected in the Github codespaces and to check up on your minikube cluster :
{% capture code_content %}minikube status{% endcapture %}
{% include code.html code=code_content %}
Verify your kubeconfig's current context.
{% capture code_content %}kubectl cluster-info{% endcapture %}
{% include code.html code=code_content %}

## Installation: Using `mesheryctl`

Use Meshery's CLI to streamline your connection to your Minikube cluster. Configure Meshery to connect to your Minikube cluster by executing:

{% capture code_content %}$ mesheryctl system config minikube{% endcapture %}
{% include code.html code=code_content %}

Once configured, execute the following command to start Meshery.

{% capture code_content %}$ mesheryctl system start{% endcapture %}
{% include code.html code=code_content %}

If you encounter any authentication issues, you can use `mesheryctl system login`. For more information, click [here](/guides/mesheryctl/authenticate-with-meshery-via-cli) to learn more.

## Installation: Using Helm

For detailed instructions on installing Meshery using Helm V3, please refer to the [Helm Installation](/installation/helm) guide.

## Installation: Manual Steps

You may also manually generate and load the kubeconfig file for Meshery to use:

**The following configuration yaml will be used by Meshery. Copy and paste the following in your config file** :

{% capture code_content %}apiVersion: v1
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
  client-key-data: < key shortcutted >{% endcapture %}
  {% include code.html code=code_content %}

_Note_: Make sure _current-context_ is set to _minikube_.

<br />
**To allow Meshery to auto detect your config file, Run** :
{% capture code_content %}kubectl config view --minify --flatten > config_minikube.yaml{% endcapture %}
{% include code.html code=code_content %}

<br />
Meshery should now be connected with your managed Kubernetes instance. Take a look at the [Meshery guides]({{ site.baseurl }}/guides) for advanced usage tips.

# Post-Installation Steps

Otionally, you can verify the health of your Meshery deployment, using <a href='/reference/mesheryctl/system/check'>mesheryctl system check</a>.

You're ready to use Meshery! Open your browser and navigate to the Meshery UI.

{% include_cached installation/accessing-meshery-ui.md %}

For further information to access meshery-ui/port-forwarding in Github Codespace, read the [docs](https://docs.github.com/en/codespaces/developing-in-a-codespace/forwarding-ports-in-your-codespace?tool=vscode)

{% include related-discussions.html tag="meshery" %}