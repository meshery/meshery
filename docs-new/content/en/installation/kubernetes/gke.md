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
  - [Preflight Checks](#preflight-checks)
    - [Preflight: Cluster Connectivity](#preflight-cluster-connectivity)
    - [Preflight: Plan your access to Meshery UI](#preflight-plan-your-access-to-meshery-ui)
  - [Installation: Using `mesheryctl`](#installation-using-mesheryctl)
  - [Installation: Using Helm](#installation-using-helm)
  - [Post-Installation Steps](#post-installation-steps)

# In-cluster Installation

Follow the steps below to install Meshery in your GKE cluster.

## Preflight Checks

Read through the following considerations prior to deploying Meshery on GKE.

### Preflight: Cluster Connectivity

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

### Preflight: Plan your access to Meshery UI

1. If you are using port-forwarding, please refer to the [port-forwarding](/reference/mesheryctl/system/dashboard) guide for detailed instructions.
2. If you are using a LoadBalancer, please refer to the [LoadBalancer](/installation/kubernetes#exposing-meshery-serviceloadbalancer) guide for detailed instructions.
3. Customize your Meshery Provider Callback URL. Meshery Server supports customizing authentication flow callback URL, which can be configured in the following way:

{{< code >}}
MESHERY_SERVER_CALLBACK_URL=https://custom-host mesheryctl system start
{{< /code >}}

Meshery should now be running in your GKE cluster, and the Meshery UI should be accessible at the `EXTERNAL IP` of the `meshery` service.

## Installation: Using `mesheryctl`

Use Meshery's CLI to streamline your connection to your GKE cluster. Configure Meshery to connect to your GKE cluster by executing:

{{< code >}}
mesheryctl system config gke
{{< /code >}}

Once configured, execute the following command to start Meshery:

{{< code >}}
mesheryctl system start
{{< /code >}}

If you encounter any authentication issues, you can use `mesheryctl system login`. For more information, click [here](/guides/mesheryctl/authenticate-with-meshery-via-cli) to learn more.

## Installation: Using Helm

For detailed instructions on installing Meshery using Helm V3, please refer to the [Helm Installation](/installation/kubernetes/helm) guide.

## Post-Installation Steps

Optionally, you can verify the health of your Meshery deployment using [mesheryctl system check](/reference/mesheryctl/system/check).

You're ready to use Meshery! Open your browser and navigate to the Meshery UI.

{{< accessing-meshery-ui >}}

{{< related-discussions tag="meshery" >}}