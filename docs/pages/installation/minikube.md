---
layout: default
title: Minikube
permalink: installation/kubernetes/minikube
type: installation
category: kubernetes
redirect_from:
- installation/platforms/minikube
display-title: "false"
language: en
list: include
image: /assets/img/platforms/minikube.png
---

<h1>Quick Start with {{ page.title }} <img src="{{ page.image }}" style="width:35px;height:35px;" /></h1>

Manage your Minikube clusters with Meshery. Deploy Meshery in Minikube [in-cluster](#in-cluster-installation) or outside of Minikube [out-of-cluster](#out-of-cluster-installation). **_Note: It is advisable to [Install Meshery in your Minikube clusters](#install-meshery-into-your-minikube-cluster)_**

<div class="prereqs"><p><strong style="font-size: 20px;">Prerequisites</strong> </p> 
  <ol>
    <li>Install the Meshery command line client, <a href="{{ site.baseurl }}/installation/mesheryctl" class="meshery-light">mesheryctl</a>.</li>
    <li>Install <a href="https://kubernetes.io/docs/tasks/tools/">kubectl</a> installed on your local machine.</li>
  </ol>
</div>

Also see: [Install Meshery on Kubernetes]({{ site.baseurl }}/installation/kubernetes)

## Available Deployment Methods

- [In-cluster Installation](#in-cluster-installation)
  - [Preflight Checks](#preflight-checks)
    - [Preflight: Cluster Connectivity](#preflight-cluster-connectivity)
    - [Preflight: Plan your access to Meshery UI](#preflight-plan-your-access-to-meshery-ui)
  - [Installation: Using `mesheryctl`](#installation-using-mesheryctl)
  - [Installation: Using Helm](#installation-using-helm)
  - [Installation: Manual Steps](#installation-manual-steps)
  - [Installation: Docker Driver Users](#installation-docker-driver-users)
- [Out-of-cluster Installation](#out-of-cluster-installation)
  - [Installation: Install Meshery on Docker](#installation-install-meshery-on-docker)
  - [Installation: Upload Config File in Meshery Web UI](#installation-upload-config-file-in-meshery-web-ui)
- [Post-Installation Steps](#post-installation-steps)
  - [Access Meshery UI](#access-meshery-ui)

# In-cluster Installation

Follow the steps below to install Meshery in your Minikube cluster.

## Preflight Checks

Read through the following considerations prior to deploying Meshery on Minikube.

### Preflight: Cluster Connectivity

Start the minikube, if not started using the following command:
{% capture code_content %}minikube start --cpus 4 --memory 8192 --kubernetes-version=v1.14.1{% endcapture %}
{% include code.html code=code_content %}
Check up on your minikube cluster :
{% capture code_content %}minikube status{% endcapture %}
{% include code.html code=code_content %}
Verify your kubeconfig's current context.
{% capture code_content %}kubectl config current-context{% endcapture %}
{% include code.html code=code_content %}

### Preflight: Plan your access to Meshery UI

1. If you are using port-forwarding, please refer to the [port-forwarding]({{ site.baseurl }}/reference/mesheryctl/system/dashboard) guide for detailed instructions.
2. Customize your Meshery Provider Callback URL. Meshery Server supports customizing authentication flow callback URL, which can be configured in the following way:

{% capture code_content %}$ MESHERY_SERVER_CALLBACK_URL=https://custom-host mesheryctl system start{% endcapture %}
{% include code.html code=code_content %}

Meshery should now be running in your Minikube cluster and Meshery UI should be accessible at the `INTERNAL IP` of `meshery` service.

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

## Installation: Docker Driver Users

Follow the [installation steps](/installation/quick-start) to setup the mesheryctl CLI and install Meshery.

**Users using docker driver**:
After completing the Meshery installation, execute the following commands to establish connectivity between Meshery Server and Kubernetes cluster:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">docker network connect bridge meshery_meshery_1</div></div>
 </pre>

<br/>

<pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">docker network connect minikube meshery_meshery_1</div></div>
 </pre>

To establish connectivity between a particular Meshery Adapter and Kubernetes server, use _"docker ps"_ to identify the name of the desired container, and execute the following commands:

<pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">docker network connect bridge &#60; container name of the desired adapter &#62;</div></div>
 </pre>

<br/>

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">docker network connect minikube &#60; container name of the desired adapter &#62;</div></div>
 </pre>

# Out-of-cluster Installation

Install Meshery on Docker (out-of-cluster) and connect it to your Minikube cluster.

## Installation: Install Meshery on Docker

{% capture code_content %}$ mesheryctl system start -p docker{% endcapture %}
{% include code.html code=code_content %}

Configure Meshery to connect to your cluster by executing:

{% capture code_content %}$ mesheryctl system config minikube{% endcapture %}
{% include code.html code=code_content %}

Once you have verified that all the services are up and running, Meshery UI will be accessible on your local machine on port 9081. Open your browser and access Meshery at [`http://localhost:9081`](http://localhost:9081).

## Installation: Upload Config File in Meshery Web UI

- Run the below command to generate the _"config_minikube.yaml"_ file for your cluster:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">kubectl config view --minify --flatten > config_minikube.yaml</div></div>
 </pre>

- Upload the generated config file by navigating to _Settings > Environment > Out of Cluster Deployment_ in the Web UI and using the _"Upload kubeconfig"_ option.

# Post-Installation Steps

## Access Meshery UI

To access Meshery's UI, please refer to the [instruction](/tasks/accessing-meshery-ui) for detailed guidance.

{% include suggested-reading.html language="en" %}

{% include related-discussions.html tag="meshery" %}

