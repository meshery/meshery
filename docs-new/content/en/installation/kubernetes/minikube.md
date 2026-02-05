---
title: "Minikube"
description: "Install Meshery on Minikube. Deploy Meshery in-cluster or out-of-cluster."
weight: 60
aliases:
  - /installation/platforms/minikube
image: /images/platforms/minikube.png
display_title: "false"
---

<h1>Quick Start with Minikube <img src="/images/platforms/minikube.png" style="width:35px;height:35px;" /></h1>

Meshery can manage your minikube clusters and is particularly useful for multi-cluster management and deployments.

**There are two ways to create this connection:**

1. Deploying Meshery in minikube [(in-cluster)](#in-cluster-installation).
2. Deploying Meshery using Docker and connect it to minikube [(out-of-cluster)](#out-of-cluster-installation).

**_Note: It is advisable to install Meshery in your Minikube clusters_**

<div class="prereqs"><h4>Prerequisites</h4>
  <ol>
    <li>Install the Meshery command-line client, <a href="/installation/mesheryctl" class="meshery-light">mesheryctl</a>.</li>
    <li>Install <a href="https://kubernetes.io/docs/tasks/tools/">kubectl</a> on your local machine.</li>
    <li>Install <a href="https://minikube.sigs.k8s.io/docs/start/">Minikube</a> on your local machine.</li>
    <li>Install <a href="https://helm.sh/docs/intro/install/">Helm</a>.</li>
  </ol>
</div>

Also see: [Install Meshery on Kubernetes](/installation/kubernetes)

## Preflight: Cluster Connectivity

Start minikube if it is not already running:

{{< code >}}
$ minikube start
{{< /code >}}

Check the status of your minikube cluster:

{{< code >}}
$ minikube status
{{< /code >}}

Verify that the current context is set to minikube:

{{< code >}}
$ kubectl config current-context
{{< /code >}}

# In-cluster Installation

## Installation: Using `mesheryctl`

To install Meshery inside your minikube cluster:

{{< code >}}
$ mesheryctl system start -p kubernetes
{{< /code >}}

To verify your deployment:

{{< code >}}
$ helm list -A -n meshery
{{< /code >}}

After deployment, access the Meshery UI using port forwarding:

{{< code >}}
$ mesheryctl system dashboard --port-forward
{{< /code >}}

## Installation: Using Helm

For detailed instructions, please refer to the [Helm Installation](/installation/kubernetes/helm) guide.

# Out-of-cluster Installation

To install Meshery on Docker and connect it to your Minikube cluster:

{{< code >}}
$ mesheryctl system start -p docker
{{< /code >}}

Configure Meshery to connect with your minikube cluster:

{{< code >}}
$ mesheryctl system config minikube
{{< /code >}}

The `mesheryctl system config minikube` command properly configures and uploads your kubeconfig file to the Meshery UI.

<a href="/images/applications/minikube-upload.png"><img alt="Minikube KubeConfig Upload" style="width:500px;height:auto;" src="/images/applications/minikube-upload.png" /></a>

# Post-Installation Steps

Verify the health of your Meshery deployment using [mesheryctl system check](/reference/mesheryctl/system/check).

{{< accessing-meshery-ui >}}

{{< related-discussions tag="meshery" >}}
