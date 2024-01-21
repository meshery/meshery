---
layout: default
title: Gitpod
permalink: installation/gitpod
type: installation
category: kubernetes
redirect_from:
- installation/platforms/gitpod
display-title: "false"
language: en
list: include
image: /assets/img/platforms/gitpod.png
abstract: Build and contribute to Meshery using Gitpod
---

<h1>Quick Start with {{ page.title }} <img src="{{ page.image }}" style="width:40px;height:40px;" /></h1>

Use Minikube in Gitpod to setup your development environment for Meshery.

<div class="prereqs"><p><strong style="font-size: 20px;">Prerequisites</strong> </p> 
  <ol>
    <li>Make an account on <a href="https://www.gitpod.io/" class="meshery-light">Gitpod</a> and add the <a href="https://chromewebstore.google.com/detail/gitpod/dodmmooeoklaejobgleioelladacbeki" class="meshery-light">Gitpod extension</a> to your browser.</li>
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
  - [Access Meshery UI](#access-meshery-ui)

# In-cluster Installation

Follow the steps below to install Meshery in your Minikube cluster.

## Preflight Checks

Read through the following considerations prior to deploying Meshery on Minikube.

### Preflight: Cluster Connectivity


You can develop and run Meshery in a Gitpod using your choice of tool:

- A command shell, via an SSH connection initiated using GitHub CLI.
- One of the JetBrains IDEs, via the JetBrains Gateway.
- The Visual Studio Code desktop application.
- A browser-based version of Visual Studio Code.
#### Fork and Open in Gitpod:
- Open your forked Meshery repository on GitHub.
- Click on the "Gitpod" button in the top right corner of the repository page (only visible with the <a href="https://chromewebstore.google.com/detail/gitpod/dodmmooeoklaejobgleioelladacbeki" class="meshery-light">Gitpod browser extension</a> installed).

{% include alert.html type="dark" title="About Gitpod" content="Gitpod will automatically clone and open the repository for you in VSCode by default. It will also automatically build the project for you on opening and comes with Docker and other tools pre-installed making it one of the fastest ways to spin up an environment for <a href='https://github.com/meshery/meshery'>Meshery.</a>" %}

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

For detailed instructions on installing Meshery using Helm V3, please refer to the [Helm Installation](/installation/kubernetes/helm) guide.

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

## Access Meshery UI

To access Meshery's UI, please refer to the [instruction](/tasks/accessing-meshery-ui) for detailed guidance.

For accessing the UI in your localhost, read the [docs](https://www.gitpod.io/docs/configure/workspaces/ports#local-port-forwarding-via-ssh).

{% include suggested-reading.html language="en" %}

{% include related-discussions.html tag="meshery" %}