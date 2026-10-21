---
title: GKE
categories: [kubernetes]
aliases:
- /installation/platforms/gke
display_title: false
image: installation/kubernetes/gke/images/gke.png
description: Install Meshery on Google Kubernetes Engine. Deploy Meshery in GKE in-cluster or outside of GKE out-of-cluster.
---

<h1>Quick Start with GKE <img src="images/gke.png" style="width:35px;height:35px;" /></h1>

Manage your GKE clusters with Meshery. Deploy Meshery inside your GKE cluster ([in-cluster](#in-cluster-installation)), or run it outside of GKE ([out-of-cluster]({{< ref "installation/kubernetes/_index.md#out-of-cluster-installation" >}})) and connect it to your cluster. **_Note: We recommend installing Meshery in your GKE clusters._**

<div class="prereqs"><p><strong style="font-size: 20px;">Prerequisites</strong> </p> 
  <ol>
    <li>Install the Meshery command line client, <a href="{{< ref "installation/mesheryctl/_index.md" >}}" class="meshery-light">mesheryctl</a>.</li>
    <li>Install <a href="https://kubernetes.io/docs/tasks/tools/">kubectl</a> on your local machine.</li>
    <li>Install <a href="https://cloud.google.com/sdk/docs/install">gcloud CLI</a>, configured for your environment.</li>
    <li>Access to an active GKE cluster in your Google Cloud project.</li>
  </ol>
</div>

Also see: [Install Meshery on Kubernetes]({{< ref "installation/kubernetes/_index.md" >}})

## Available Deployment Methods

- [In-cluster Installation](#in-cluster-installation)
  - [Preflight Checks](#preflight-checks)
    - [Preflight: Cluster Connectivity](#preflight-cluster-connectivity)
    - [Preflight: Plan your access to Meshery UI](#preflight-plan-your-access-to-meshery-ui)
  - [Installation: Using `mesheryctl`](#installation-using-mesheryctl)
  - [Installation: Using Helm](#installation-using-helm)
  - [Post-Installation Steps](#post-installation-steps)

# In-cluster Installation

Follow the steps below to install Meshery in your GKE cluster.

## Preflight Checks

Before you deploy Meshery on GKE, check the following.

### Preflight: Cluster Connectivity

Use the gcloud CLI to connect `kubectl` to your Google Kubernetes Engine (GKE) cluster.

1. Sign in to your Google Cloud account with [gcloud auth login](https://cloud.google.com/sdk/gcloud/reference/auth/login).
2. Set the project that contains your cluster:

{{< code code="gcloud config set project [PROJECT_ID]" >}}
3. Download the cluster's credentials and make it your current context:

{{< code code="gcloud container clusters get-credentials [CLUSTER_NAME] --zone [CLUSTER_ZONE]" >}}
4. Verify your kubeconfig's current context:

{{< code code="kubectl config current-context" >}}

### Preflight: Plan your access to Meshery UI

1. To reach Meshery UI through port-forwarding, see the [port-forwarding]({{< ref "reference/references/mesheryctl/system/dashboard.md" >}}) guide.
2. To reach it through a load balancer, see [Exposing Meshery Service (LoadBalancer)]({{< ref "installation/kubernetes/_index.md#exposing-meshery-serviceloadbalancer" >}}).
3. If users reach Meshery through a custom host name, set the callback URL that your Provider redirects to after sign-in when you start Meshery:

{{< code code="MESHERY_SERVER_CALLBACK_URL=https://custom-host mesheryctl system start -p kubernetes" >}}

With a load balancer, Meshery UI is then available at the `EXTERNAL-IP` of the `meshery` service.

## Installation: Using `mesheryctl`

With your GKE cluster as the current context, start Meshery on Kubernetes:

{{< code code="mesheryctl system start -p kubernetes" >}}

Meshery Server runs in the cluster and connects to it automatically, so you do not need to upload a kubeconfig. To connect a Meshery Server that runs elsewhere to this cluster, sign in with `mesheryctl system login` and run `mesheryctl connection create --type gke`. (`mesheryctl system config gke` is deprecated.)

If `mesheryctl` reports an authentication error, sign in with `mesheryctl system login`. See [Authenticate with Meshery via CLI]({{< ref "guides/mesheryctl/authenticate-with-meshery-via-cli/index.md" >}}).

## Installation: Using Helm

To install Meshery with Helm 3, see the [Helm Installation]({{< ref "installation/kubernetes/helm.md" >}}) guide.

## Post-Installation Steps

Optionally, verify the health of your Meshery deployment with <a href='{{< ref "reference/references/mesheryctl/system/check.md" >}}'>mesheryctl system check</a>.

You're ready to use Meshery. Open Meshery UI in your browser.

{{< installation/accessing-meshery-ui >}}

{{< related-discussions tag="meshery" >}}