---
title: "GKE"
description: "Install Meshery on Google Kubernetes Engine. Deploy Meshery in-cluster or out-of-cluster."
weight: 40
aliases:
  - /installation/platforms/gke
image: /images/platforms/gke.png
display_title: "false"
---

<h1>Quick Start with GKE <img src="/images/platforms/gke.png" style="width:35px;height:35px;" /></h1>

Manage your GKE clusters with Meshery. Deploy Meshery in GKE [in-cluster](#in-cluster-installation) or outside of GKE [out-of-cluster](#out-of-cluster-installation). **_Note: It is advisable to install Meshery in your GKE clusters_**

<div class="prereqs"><h4>Prerequisites</h4>
  <ol>
    <li>Install the Meshery command line client, <a href="/installation/mesheryctl" class="meshery-light">mesheryctl</a>.</li>
    <li>Install <a href="https://kubernetes.io/docs/tasks/tools/">kubectl</a> on your local machine.</li>
    <li>Install <a href="https://cloud.google.com/sdk/docs/install">gcloud CLI</a>, configured for your environment.</li>
    <li>Access to an active GKE cluster in your Google Cloud project.</li>
  </ol>
</div>

Also see: [Install Meshery on Kubernetes](/installation/kubernetes)

## Available Deployment Methods

- [In-cluster Installation](#in-cluster-installation)
- [Installation: Using `mesheryctl`](#installation-using-mesheryctl)
- [Installation: Using Helm](#installation-using-helm)

# In-cluster Installation

Follow the steps below to install Meshery in your GKE cluster.

## Preflight: Cluster Connectivity

1. Verify your connection to a Google Kubernetes Engine Cluster using the gcloud CLI.
2. Log in to your GCP account using [gcloud auth login](https://cloud.google.com/sdk/gcloud/reference/auth/login).
3. After a successful login, set the Project Id:

{{< code >}}
gcloud config set project [PROJECT_ID]
{{< /code >}}

4. After setting the Project Id, set the cluster context:

{{< code >}}
gcloud container clusters get-credentials [CLUSTER_NAME] --zone [CLUSTER_ZONE]
{{< /code >}}

5. Verify your kubeconfig's current context:

{{< code >}}
kubectl config current-context
{{< /code >}}

## Installation: Using `mesheryctl`

Use Meshery's CLI to streamline your connection to your GKE cluster:

{{< code >}}
$ mesheryctl system config gke
{{< /code >}}

Once configured, execute the following command to start Meshery:

{{< code >}}
$ mesheryctl system start
{{< /code >}}

## Installation: Using Helm

For detailed instructions on installing Meshery using Helm V3, please refer to the [Helm Installation](/installation/kubernetes/helm) guide.

## Post-Installation Steps

Optionally, you can verify the health of your Meshery deployment using [mesheryctl system check](/reference/mesheryctl/system/check).

You're ready to use Meshery! Open your browser and navigate to the Meshery UI.

{{< accessing-meshery-ui >}}

{{< related-discussions tag="meshery" >}}
