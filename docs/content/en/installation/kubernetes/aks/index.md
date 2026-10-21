---
title: AKS
categories: [kubernetes]
aliases:
- /installation/platforms/aks
display_title: false
image: installation/kubernetes/aks/images/aks.svg
description: Manage your AKS clusters with Meshery. Deploy Meshery in AKS in-cluster or out-of-cluster.
---

<h1>Quick Start with AKS <img src="images/aks.svg" style="width:35px;height:35px;" /></h1>

Manage your AKS clusters with Meshery. Deploy Meshery inside your AKS cluster ([in-cluster](#in-cluster-installation)), or run it outside of AKS ([out-of-cluster]({{< ref "installation/kubernetes/_index.md#out-of-cluster-installation" >}})) and connect it to your cluster. **_Note: We recommend installing Meshery in your AKS clusters._**

<div class="prereqs"><h4>Prerequisites</h4>
<ol>
<li>Install the Meshery command line client, <a href="{{< ref "installation/mesheryctl/_index.md" >}}" class="meshery-light">mesheryctl</a>.</li>
<li>Install <a href="https://kubernetes.io/docs/tasks/tools/">kubectl</a> on your local machine.</li>
<li>Install <a href="https://learn.microsoft.com/en-us/cli/azure/install-azure-cli">Azure CLI</a>, configured for your environment.</li>
<li>Access to an active AKS cluster in one of your resource groups.</li>
</ol>
</div>

Also see: [Install Meshery on Kubernetes]({{< ref "installation/kubernetes/_index.md" >}})

## Available Deployment Methods

- [In-cluster Installation](#in-cluster-installation)
    - [Preflight Checks](#preflight-checks)
    - [Preflight: Cluster Connectivity](#preflight-cluster-connectivity)
    - [Installation: Using `mesheryctl`](#installation-using-mesheryctl)
    - [Installation: Using Helm](#installation-using-helm)
  - [Post-Installation Steps](#post-installation-steps)

# In-cluster Installation

Follow the steps below to install Meshery in your AKS cluster.

### Preflight Checks

Before you deploy Meshery on AKS, check the following.

### Preflight: Cluster Connectivity

Use the Azure CLI to connect `kubectl` to your Azure Kubernetes Service (AKS) cluster.

1. Sign in to your Azure account with [az login](https://learn.microsoft.com/en-us/cli/azure/authenticate-azure-cli).
1. Select the subscription that contains your AKS cluster:

{{< code code="az account set --subscription [SUBSCRIPTION_ID]" >}}

1. Download the cluster's credentials and make it your current context:

{{< code code="az aks get-credentials --resource-group [RESOURCE_GROUP] --name [AKS_SERVICE_NAME]" >}}

### Installation: Using `mesheryctl`

With your AKS cluster as the current context, start Meshery on Kubernetes:

{{< code code="mesheryctl system start -p kubernetes" >}}

Meshery Server runs in the cluster and connects to it automatically, so you do not need to upload a kubeconfig. To connect a Meshery Server that runs elsewhere to this cluster, sign in with `mesheryctl system login` and run `mesheryctl connection create --type aks`, which prompts for the resource group and cluster name. (`mesheryctl system config aks` is deprecated.)

If `mesheryctl` reports an authentication error, sign in with `mesheryctl system login`. See [Authenticate with Meshery via CLI]({{< ref "guides/mesheryctl/authenticate-with-meshery-via-cli/index.md" >}}).

### Installation: Using Helm

To install Meshery with Helm 3, see the [Helm Installation]({{< ref "installation/kubernetes/helm.md" >}}) guide.

## Post-Installation Steps

Optionally, verify the health of your Meshery deployment with <a href='{{< ref "reference/references/mesheryctl/system/check.md" >}}'>mesheryctl system check</a>.

You're ready to use Meshery. Open Meshery UI in your browser.

{{< installation/accessing-meshery-ui >}}

{{< related-discussions tag="meshery" >}}
