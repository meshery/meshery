---
title: Install on Kubernetes
categories: [kubernetes]
aliases:
- /installation/platforms/kubernetes
display_title: false
image: installation/kubernetes/images/kubernetes.svg
description: Install Meshery on Kubernetes. Deploy Meshery in Kubernetes in-cluster or outside of Kubernetes out-of-cluster.
weight: 15
---

<h1>Quick Start with Kubernetes <img src="images/kubernetes.svg" style="width:35px;height:35px;" /></h1>

Manage your Kubernetes clusters with Meshery. Deploy Meshery inside a Kubernetes cluster ([in-cluster](#in-cluster-installation)) or outside of it ([out-of-cluster](#out-of-cluster-installation)). **_Note: We recommend installing Meshery in your Kubernetes clusters._**

<div class="prereqs"><h4>Prerequisites</h4>
  <ol>
    <li>Install the Meshery command line client, <a href="{{< ref "installation/mesheryctl/_index.md" >}}" class="meshery-light">mesheryctl</a>.</li>
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
  - [Exposing Meshery Service (LoadBalancer)](#exposing-meshery-serviceloadbalancer)
- [Out-of-cluster Installation](#out-of-cluster-installation)
  - [Set up Ingress on Minikube with the NGINX Ingress Controller](#set-up-ingress-on-minikube-with-the-nginx-ingress-controller)
  - [Installing cert-manager with kubectl](#installing-cert-manager-with-kubectl)

# In-cluster Installation

Follow the steps below to install Meshery in your Kubernetes cluster.

## Preflight Checks

Before you deploy Meshery on Kubernetes, check the following.

### Preflight: Cluster Connectivity

Verify that your kubeconfig's current context points at the cluster you want to deploy Meshery to.

{{< code code="kubectl config current-context" >}}

## Installation: Using `mesheryctl`

Start Meshery. `mesheryctl` deploys to the platform of your current `mesheryctl` context, which is `kubernetes` by default. To check it, run `mesheryctl system context view`. To deploy to Kubernetes regardless of the context, pass `-p kubernetes`:

{{< code code="mesheryctl system start -p kubernetes" >}}

`mesheryctl` installs the Meshery Helm chart into the `meshery` namespace of your current cluster.

## Installation: Using Helm

To install Meshery with Helm 3, see the [Helm Installation]({{< ref "installation/kubernetes/helm.md" >}}) guide.

## Post-Installation Steps

Optionally, verify the health of your Meshery deployment with <a href='{{< ref "reference/references/mesheryctl/system/check.md" >}}'>mesheryctl system check</a>.

You're ready to use Meshery. Open Meshery UI in your browser.

{{< installation/accessing-meshery-ui >}}

## Exposing Meshery Service (LoadBalancer) {#exposing-meshery-serviceloadbalancer}

When you install Meshery in-cluster, a Kubernetes `Service` named `meshery` in the `meshery` namespace serves Meshery UI. The Helm chart creates this `Service` as type `LoadBalancer` by default (`service.type`), forwarding port `9081` to the Meshery Server container on port `8080`.

On a managed Kubernetes offering, such as GKE, EKS, AKS, or DigitalOcean Kubernetes, a `LoadBalancer` `Service` tells the cloud provider to provision an external load balancer and assign it a routable `EXTERNAL-IP`. Retrieve the address with:

{{< code code="kubectl get service meshery --namespace meshery" >}}

Once the `EXTERNAL-IP` column shows an address instead of `<pending>`, open Meshery UI in your browser at `http://[EXTERNAL-IP]:9081`.

If the `Service` was previously set to another type (for example, `ClusterIP`), switch it back to `LoadBalancer`:

{{< code code=`kubectl patch service meshery --namespace meshery --type merge -p '{"spec":{"type":"LoadBalancer"}}'` >}}

{{% alert title="EXTERNAL-IP not assigned?" color="warning" %}}
A `LoadBalancer` `Service` gets an external address only when the cluster runs a load balancer controller. Managed clouds provide one out of the box, but bare-metal and local clusters (for example, Minikube or kind) do not. On those clusters, install a load balancer implementation such as [MetalLB](https://metallb.universe.tf/), expose Meshery through a `NodePort` `Service` instead, or reach the UI with port-forwarding by following the [mesheryctl system dashboard]({{< ref "reference/references/mesheryctl/system/dashboard.md" >}}) guide.
{{% /alert %}}

# Out-of-cluster Installation

Install Meshery on Docker (out-of-cluster) and connect it to your Kubernetes cluster. Follow [Install on Docker]({{< ref "installation/docker/_index.md" >}}), then sign in with `mesheryctl system login` and add your cluster with `mesheryctl connection create --type <aks|eks|gke|minikube>`, or upload its kubeconfig in Meshery UI under **Lifecycle** > **Connections**.

If you expose Meshery through an ingress on Minikube, the following steps set up an ingress controller and cert-manager.

<!-- ## Installation: Upload Config File in Meshery Web UI

- Run the below command to generate the _"config_minikube.yaml"_ file for your cluster:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">kubectl config view --minify --flatten > config_minikube.yaml</div></div>
 </pre>

- Upload the generated config file by navigating to _Settings > Environment > Out of Cluster Deployment_ in the Web UI and using the _"Upload kubeconfig"_ option. -->

## Set up Ingress on Minikube with the NGINX Ingress Controller

- Enable the NGINX Ingress controller for your cluster:

{{< code code="minikube addons enable ingress" >}}

- Check that the NGINX Ingress controller is running:

{{< code code="kubectl get pods -n ingress-nginx" >}}

## Installing cert-manager with kubectl

- Install cert-manager in your cluster:

{{< code code="kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.15.3/cert-manager.yaml" >}}

{{< related-discussions tag="meshery" >}}