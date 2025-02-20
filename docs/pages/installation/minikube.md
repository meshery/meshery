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
abstract: Install Meshery on Minikube. Deploy Meshery in Minikube in-cluster or outside of Minikube out-of-cluster.
---

<h1>Quick Start with {{ page.title }} <img src="{{ page.image }}" style="width:35px;height:35px;" /></h1>

Meshery can manage your minikube clusters and is particularly useful for multi-cluster management and deployments.

For Meshery to manage your Minikube cluster, it has to be discovered and added as a Kubernetes connection in the Meshery server.

After your cluster has been added as a connection, you can use Meshery to make infrastructure deployments of your [Meshery Designs](https://cloud.layer5.io/academy/learning-paths/mastering-meshery/introduction-to-meshery?chapter=creating-designs) to your cluster. To learn more about this, See [Deploying Meshery Designs](https://cloud.layer5.io/academy/learning-paths/mastering-meshery/introduction-to-meshery?chapter=deploying-meshery-designs). You can also Visualize the resources in your cluster in Operator Mode, See [Views in Visualizer](https://docs.layer5.io/kanvas/visualizer/visualizer-views/).

**There are two ways to create this connection:**

1. Deploying Meshery in minikube [(in-cluster)](#in-cluster-installation).
1. Deploying Meshery using Docker and connect it to minikube [(out-of-cluster)](#out-of-cluster-installation).

**_Note: It is advisable to [Install Meshery in your Minikube clusters](#install-meshery-into-your-minikube-cluster)_**

<div class="prereqs"><p><strong style="font-size: 20px;">Prerequisites</strong> </p> 
  <ol>
    <li>Install the Meshery command-line client, <a href="{{ site.baseurl }}/installation/mesheryctl" class="meshery-light">mesheryctl</a>.</li>
    <li>Install <a href="https://kubernetes.io/docs/tasks/tools/">kubectl</a> on your local machine.</li>
    <li>Install <a href="https://minikube.sigs.k8s.io/docs/start/?arch=%2Fmacos%2Fx86-64%2Fstable%2Fbinary+download"> Minikube</a> on your local machine.</li>
    <li>Install <a href="https://helm.sh/docs/intro/install/">Helm</a>.</li>
  </ol>
</div>

Also see: [Install Meshery on Kubernetes]({{ site.baseurl }}/installation/kubernetes)

## Available Deployment Methods

- [Preflight Checks](#preflight-checks)
  - [Preflight: Cluster Connectivity](#1-preflight-cluster-connectivity)
  - [Preflight: Meshery Authentication](#2-preflight-meshery-authentication)

- [In-cluster Installation](#in-cluster-installation)
  - [Installation: Install Meshery on Kubernetes Using `mesheryctl`](#installation-install-meshery-on-kubernetes-using-mesheryctl)
  - [Installation: Using Helm](#installation-using-helm)

- [Out-of-cluster Installation](#out-of-cluster-installation)
  - [Installation: Install Meshery on Docker Using `mesheryctl`](#installation-install-meshery-on-docker)
  - [Minikube Docker Driver Users](#minikube-docker-driver-users)

- [Uploading Configuration File in the Meshery Web UI](#uploading-configuration-file-in-the-meshery-web-ui)
- [Post-Installation Steps](#post-installation-steps)
  
## Preflight Checks

Before deploying Meshery on minikube, complete the following initial setup tasks to prepare your environment.

### 1. Preflight: Cluster Connectivity

Start minikube using the following command if it is not already running:
{% capture code_content %}$ minikube start{% endcapture %}
{% include code.html code=code_content %}
Check the status of your minikube cluster by running:
{% capture code_content %}$ minikube status{% endcapture %}
{% include code.html code=code_content %}
Verify that the current context is set to minikube by running:
{% capture code_content %}$ kubectl config current-context{% endcapture %}
{% include code.html code=code_content %}

### 2. Preflight: Meshery Authentication

Ensure you are logged in and [authenticated with Meshery](https://docs.meshery.io/guides/mesheryctl/authenticate-with-meshery-via-cli) by running the following command:

{% capture code_content %}$ mesheryctl system login{% endcapture %}
{% include code.html code=code_content %}

# In-cluster Installation

## Installation: Install Meshery on Kubernetes Using `mesheryctl`

To install Meshery inside your minikube cluster, run the command:
{% capture code_content %}$ mesheryctl system start -p kubernetes{% endcapture %}
{% include code.html code=code_content %}
This command deploys the Meshery Helm chart in the Meshery namespace.

To verify your deployment, run:
{% capture code_content %}$ helm list -A -n meshery{% endcapture %}
{% include code.html code=code_content %}
After deployment, access the Meshery UI using port forwarding, with the command:
{% capture code_content %}$  mesheryctl system dashboard --port-forward {% endcapture %}
{% include code.html code=code_content %}
For detailed instructions on port forwarding, refer to the [port-forwarding]({{ site.baseurl }}/reference/mesheryctl/system/dashboard) guide.

By default, Meshery auto-detects your Minikube cluster and establishes a connection. However, if this doesn’t happen, you can connect by running the following command:
{% capture code_content %}$  mesheryctl system config minikube {% endcapture %}
{% include code.html code=code_content %}

The `mesheryctl system config minikube` command properly configures and uploads your kubeconfig file to the Meshery UI.

<a href="{{ site.baseurl }}/assets/img/applications/minikube-upload.png"><img alt="Minikube KubeConfig Upload" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/applications/minikube-upload.png" /></a>

## Installation: Using Helm

You can deploy Meshery directly using the Helm CLI.
For detailed instructions on installing Meshery using Helm V3, please refer to the [Helm Installation](/installation/kubernetes/helm) guide.

# Out-of-cluster Installation

To install Meshery on Docker(out-of-cluster) and connect it to your Minikube cluster, follow these steps:

## Installation: Install Meshery on Docker

Run the following command to start Meshery in a Docker environment:
{% capture code_content %}$ mesheryctl system start -p docker{% endcapture %}
{% include code.html code=code_content %}
This will spin up meshery docker containers. To verify that Meshery  is running, use
{% capture code_content %}$ docker ps{% endcapture %}
{% include code.html code=code_content %}
Meshery UI will be accessible on your local machine on port 9081. Open your browser and access Meshery at http://localhost:9081.

Configure Meshery to connect with your minikube cluster by running the command:
{% capture code_content %}$ mesheryctl system config minikube {% endcapture %}
{% include code.html code=code_content %}

### Minikube Docker Driver Users

For users running minikube with the Docker driver, specific steps are needed to ensure that Meshery can connect properly to your minikube cluster.

If you set up your minikube cluster using the [Docker driver](https://minikube.sigs.k8s.io/docs/drivers/docker/), both minikube and Meshery will be running in Docker containers. So, you need to ensure that the Meshery and minikube containers can communicate with each other by placing them in the same Docker network.

To configure this, run the following commands:

<pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">$ docker network connect bridge meshery-meshery-1</div></div>
</pre>

<br/>

<pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">$ docker network connect minikube meshery-meshery-1</div></div>
</pre>

Next, update the Kubernetes API server address in your kubeconfig file before running the `mesheryctl system config minikube` command. The steps to do this are outlined below.

#### Next Step: Update the Kubernetes API Server Address for Meshery Access

To allow the Meshery container to access your Minikube cluster (since both are running in containers), you need to update the Kubernetes API server address in your `kubeconfig file` to the `external minikube IP address`. This is necessary because Docker typically forwards ports to a localhost address, which isn’t accessible between containers.

To retrieve the Minikube IP, run the command `minikube ip`. To check which port minikube is using, run `docker ps` to view the container's port, which is typically `8443`.

Open the kubeconfig file and update the server address.

{% capture code_content %}$ nano ~/.kube/config.yaml

## Change the server address

server: https://{minikubeIP}:{port}

{% endcapture %}
{% include code.html code=code_content %}

`Ctrl + X` then enter `Y` to save and close the file.

Next run this command to configure Meshery to access your cluster.
{% capture code_content %}$  mesheryctl system config minikube{% endcapture %}
{% include code.html code=code_content %}

**Note**: An alternative to running the mesheryctl system config minikube command for meshery to discover your cluster is manually uploading your config file to the UI.

# Uploading Configuration File in the Meshery Web UI

**Note**: Meshery can only connect to your cluster if it is running locally (Kubernetes or Docker). Direct connections are not possible when using the hosted [Meshery Playground](https://playground.meshery.io/).

**To manually upload your kubeconfig file after running Meshery locally**:

1. In the Meshery UI, navigate to **Lifecycle** from the menu on the left.
2. Click on **Connections**.
3. Click on **Add Cluster** and search for your kubeconfig file.
4. Click **Import**.

**Note**:  If you encounter a connections refused error while uploading your kubeconfig, try changing your cluster server URL to the external API address of minikube. To do this follow the steps listed in the [Minikube Docker Driver Users Section](#docker-driver-update-the-kubernetes-api-server-address-for-meshery-access).

#### Troubleshooting Meshery Installation

If you experience any issues during installation, refer to the [Troubleshooting Meshery Installations](https://docs.meshery.io/guides/troubleshooting/installation#setting-up-meshery-using-kind-or-minikube) guide for help.

# Post-Installation Steps

Verify the health of your Meshery deployment, using <a href='/reference/mesheryctl/system/check'>mesheryctl system check</a>.

{% include_cached installation/accessing-meshery-ui.md %}

{% include related-discussions.html tag="meshery" %}
