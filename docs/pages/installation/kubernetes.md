---
layout: default
title: Kubernetes
permalink: installation/kubernetes
type: installation
category: kubernetes
redirect_from:
- installation/platforms/kubernetes
display-title: "false"
language: en
list: include
image: /assets/img/platforms/kubernetes.svg
abstract: Install Meshery on Kubernetes. Deploy Meshery in Kubernetes in-cluster or outside of Kubernetes out-of-cluster.
---

<h1>Quick Start with {{ page.title }} <img src="{{ page.image }}" style="width:35px;height:35px;" /></h1>

Manage your kubernetes clusters with Meshery. Deploy Meshery in kubernetes [in-cluster](#in-cluster-installation) or outside of kubernetes [out-of-cluster](#out-of-cluster-installation). **_Note: It is advisable to [Install Meshery in your kubernetes clusters](#install-meshery-into-your-kubernetes-cluster)_**

<div class="prereqs"><h4>Prerequisites</h4>
  <ol>
    <li>Install the Meshery command line client, <a href="{{ site.baseurl }}/installation/mesheryctl" class="meshery-light">mesheryctl</a>.</li>
    <li>Install <a href="https://kubernetes.io/docs/tasks/tools/">kubectl</a> on your local machine.</li>
    <li>Access to an active kubernetes cluster.</li>
  </ol>
</div>

## Available Deployment Methods

- [In-cluster Installation](#in-cluster-installation)
  - [Preflight Checks](#preflight-checks)
    - [Preflight: Cluster Connectivity](#preflight-cluster-connectivity)
    - [Preflight: Plan your access to Meshery UI](#preflight-plan-your-access-to-meshery-ui)
  - [Installation: Using `mesheryctl`](#installation-using-mesheryctl)
  - [Installation: Using Helm](#installation-using-helm)
- [Post-Installation Steps](#post-installation-steps)
  - [Access Meshery UI](#access-meshery-ui)
- [Out-of-cluster Installation](#out-of-cluster-installation)
  - [Installation: Upload Config File in Meshery Web UI](#installation-upload-config-file-in-meshery-web-ui)

# In-cluster Installation

Follow the steps below to install Meshery in your kubernetes cluster.

## Preflight Checks

Read through the following considerations prior to deploying Meshery on kubernetes.

### Preflight: Cluster Connectivity

Verify your kubeconfig's current context is set the kubernetes cluster you want to deploy Meshery.
{% capture code_content %}kubectl config current-context{% endcapture %}
{% include code.html code=code_content %}

### Preflight: Plan your access to Meshery UI

1. If you are using port-forwarding, please refer to the [port-forwarding](/tasks/accessing-meshery-ui) guide for detailed instructions.
2. Customize your Meshery Provider Callback URL. Meshery Server supports customizing authentication flow callback URL, which can be configured in the following way:

{% capture code_content %}$ MESHERY_SERVER_CALLBACK_URL=https://custom-host mesheryctl system start{% endcapture %}
{% include code.html code=code_content %}

Meshery should now be running in your Kubernetes cluster and Meshery UI should be accessible at the `EXTERNAL IP` of `meshery` service.

## Installation: Using `mesheryctl`

Once configured, execute the following command to start Meshery.

Before executing the below command, go to ~/.meshery/config.yaml and ensure that current platform is set to kubernetes.
{% capture code_content %}$ mesheryctl system start{% endcapture %}
{% include code.html code=code_content %}

If you encounter any authentication issues, you can use `mesheryctl system login`. For more information, click [here](/guides/mesheryctl/authenticate-with-meshery-via-cli) to learn more.

## Installation: Using Helm

For detailed instructions on installing Meshery using Helm V3, please refer to the [Helm Installation](/installation/kubernetes/helm) guide.

# Post-Installation Steps

## Access Meshery UI

To access Meshery's UI, please refer to the [instruction](/tasks/accessing-meshery-ui) for detailed guidance.

# Out-of-cluster Installation

Install Meshery on Docker (out-of-cluster) and connect it to your Kubernetes cluster.

## Installation: Upload Config File in Meshery Web UI

- Run the below command to generate the _"config_minikube.yaml"_ file for your cluster:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">kubectl config view --minify --flatten > config_minikube.yaml</div></div>
 </pre>

- Upload the generated config file by navigating to _Settings > Environment > Out of Cluster Deployment_ in the Web UI and using the _"Upload kubeconfig"_ option.

{% if page.suggested-reading != false and page.title and page.type and page.category and page.url %}
{% include_cached suggested-reading.html  title=page.title type=page.type category=page.category url=page.url language="en" %}
{% endif %}

{% include related-discussions.html tag="meshery" %}
