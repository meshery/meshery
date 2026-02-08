---
title: "Kubernetes"
description: "Install Meshery on Kubernetes. Deploy Meshery in-cluster or out-of-cluster."
weight: 20
aliases:
  - /installation/platforms/kubernetes
image: /images/platforms/kubernetes.svg
display_title: "false"
---

<h1>Quick Start with Kubernetes <img src="/images/platforms/kubernetes.svg" style="width:35px;height:35px;" /></h1>

Manage your Kubernetes clusters with Meshery. Deploy Meshery in Kubernetes [in-cluster](#in-cluster-installation) or outside of Kubernetes [out-of-cluster](#out-of-cluster-installation). **_Note: It is advisable to install Meshery in your Kubernetes clusters_**

<div class="prereqs"><h4>Prerequisites</h4>
  <ol>
    <li>Install the Meshery command line client, <a href="/installation/mesheryctl" class="meshery-light">mesheryctl</a>.</li>
    <li>Install <a href="https://kubernetes.io/docs/tasks/tools/">kubectl</a> on your local machine.</li>
    <li>Access to an active Kubernetes cluster.</li>
  </ol>
</div>

## Available Deployment Methods

- [In-cluster Installation](#in-cluster-installation)
  - [Preflight Checks](#preflight-checks)
    - [Preflight: Cluster Connectivity](#preflight-cluster-connectivity)
  - [Installation: Using `mesheryctl`](#installation-using-mesheryctl)
  - [Installation: Using Helm](#installation-using-helm)
  - [Post-Installation Steps](#post-installation-steps)
- [Out-of-cluster Installation](#out-of-cluster-installation)
  - [Set up Ingress on Minikube with the NGINX Ingress Controller](#set-up-ingress-on-minikube-with-the-nginx-ingress-controller)
  - [Installing cert-manager with kubectl](#installing-cert-manager-with-kubectl)

# In-cluster Installation

Follow the steps below to install Meshery in your Kubernetes cluster.

## Preflight Checks

Read through the following considerations prior to deploying Meshery on Kubernetes.

### Preflight: Cluster Connectivity

Verify your kubeconfig's current context is set to the Kubernetes cluster you want to deploy Meshery to.

{{< code >}}
kubectl config current-context
{{< /code >}}

## Installation: Using `mesheryctl`

Once configured, execute the following command to start Meshery.

Before executing the below command, go to ~/.meshery/config.yaml and ensure that the current platform is set to Kubernetes.

{{< code >}}
mesheryctl system start
{{< /code >}}

## Installation: Using Helm

For detailed instructions on installing Meshery using Helm V3, please refer to the [Helm Installation](/installation/kubernetes/helm) guide.

## Post-Installation Steps

Optionally, you can verify the health of your Meshery deployment using <a href='/reference/mesheryctl/system/check'>mesheryctl system check</a>.

You're ready to use Meshery! Open your browser and navigate to the Meshery UI.

{{< accessing-meshery-ui >}}

# Out-of-cluster Installation

Install Meshery on Docker (out-of-cluster) and connect it to your Kubernetes cluster.

## Set up Ingress on Minikube with the NGINX Ingress Controller

Run the below command to enable the NGINX Ingress controller for your cluster:

{{< code >}}
minikube addons enable ingress
{{< /code >}}

To check if NGINX Ingress controller is running:

{{< code >}}
kubectl get pods -n ingress-nginx
{{< /code >}}

## Installing cert-manager with kubectl

Run the below command to install cert-manager for your cluster:

{{< code >}}
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/latest/download/cert-manager.yaml
{{< /code >}}

{{< related-discussions tag="meshery" >}}