---
title: AKS
description: "Manage your AKS clusters with Meshery. Deploy Meshery in AKS in-cluster or out-of-cluster."
weight: 20
aliases:
  - /installation/platforms/aks
image: /images/platforms/aks.svg
display_title: false
---

<h1>Quick Start with AKS <img src="/images/platforms/aks.svg" style="width:35px;height:35px;" /></h1>

Manage your AKS clusters with Meshery. Deploy Meshery in AKS [in-cluster](#in-cluster-installation) or outside of AKS [out-of-cluster](#out-of-cluster-installation). **_Note: It is advisable to install Meshery in your AKS clusters_**

<div class="prereqs"><h4>Prerequisites</h4>
<ol>
<li>Install the Meshery command line client, <a href="/installation/mesheryctl" class="meshery-light">mesheryctl</a>.</li>
<li>Install <a href="https://kubernetes.io/docs/tasks/tools/">kubectl</a> on your local machine.</li>
<li>Install <a href="https://learn.microsoft.com/en-us/cli/azure/install-azure-cli">Azure CLI</a>, configured for your environment.</li>
<li>Access to an active AKS cluster in one of your resource groups.</li>
</ol>
</div>

Also see: [Install Meshery on Kubernetes](/installation/kubernetes)

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

Read through the following considerations prior to deploying Meshery on AKS.

### Preflight: Cluster Connectivity

1. Verify your connection to an Azure Kubernetes Services Cluster using Azure CLI.
2. Login to Azure account using [az login](https://learn.microsoft.com/en-us/cli/azure/authenticate-azure-cli).
3. After a successful login, identify the subscription associated with your AKS cluster:

{{< code >}}
az account set --subscription [SUBSCRIPTION_ID]
{{< /code >}}

4. After setting the subscription, set the cluster context:

{{< code >}}
az aks get-credentials --resource-group [RESOURCE_GROUP] --name [AKS_SERVICE_NAME]
{{< /code >}}

### Installation: Using `mesheryctl`

Use Meshery's CLI to streamline your connection to your AKS cluster. Configure Meshery to connect to your AKS cluster by executing:

{{< code >}}
mesheryctl system config aks
{{< /code >}}

Once configured, execute the following command to start Meshery:

{{< code >}}
mesheryctl system start
{{< /code >}}

If you encounter any authentication issues, you can use `mesheryctl system login`. For more information, click [here](/guides/mesheryctl/authenticate-with-meshery-via-cli) to learn more.

### Installation: Using Helm

For detailed instructions on installing Meshery using Helm V3, please refer to the [Helm Installation](/installation/kubernetes/helm) guide.

## Post-Installation Steps

Optionally, you can verify the health of your Meshery deployment, using <a href='/reference/mesheryctl/system/check'>mesheryctl system check</a>.

You're ready to use Meshery! Open your browser and navigate to the Meshery UI.

{{< accessing-meshery-ui >}}

{{< related-discussions tag="meshery" >}}