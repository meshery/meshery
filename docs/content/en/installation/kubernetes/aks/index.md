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

Manage your AKS clusters with Meshery. Deploy Meshery in AKS [in-cluster](#in-cluster-installation) or outside of AKS [out-of-cluster](#out-of-cluster-installation). **_Note: It is advisable to install Meshery in your AKS clusters_**

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
    - [Preflight: Plan your access to Meshery UI](#preflight-plan-your-access-to-meshery-ui)
  - [Installation: Using `mesheryctl`](#installation-using-mesheryctl)
  - [Installation: Using Helm](#installation-using-helm)
  - [Post-Installation Steps](#post-installation-steps)
  - [Accessing Meshery UI](#accessing-meshery-ui)
- [Out-of-cluster Installation](#out-of-cluster-installation)
  - [Install Meshery on Docker](#install-meshery-on-docker)

# In-cluster Installation

Follow the steps below to install Meshery in your AKS cluster.

## Preflight Checks

Read through the following considerations prior to deploying Meshery on AKS.

### Preflight: Cluster Connectivity

Verify your connection to an Azure Kubernetes Service cluster using Azure CLI.

1. Login to Azure account using [az login](https://learn.microsoft.com/en-us/cli/azure/authenticate-azure-cli).
2. After a successful login, identify the subscription associated with your AKS cluster:

{{< code code="az account set --subscription [SUBSCRIPTION_ID]" >}}

3. After setting the subscription, set the cluster context:

{{< code code="az aks get-credentials --resource-group [RESOURCE_GROUP] --name [AKS_CLUSTER_NAME]" >}}

4. Verify your kubeconfig's current context:

{{< code code="kubectl config current-context" >}}

### Preflight: Plan your access to Meshery UI

1. If you are using port-forwarding, please refer to the [Port Forwarding Guide]({{< ref "installation/kubernetes/port-forwarding.md" >}}) for detailed instructions. This guide is common for all cloud providers including AWS, Azure, GCP, and Minikube cluster services.
2. If you are using a LoadBalancer, please refer to the [LoadBalancer]({{< ref "installation/kubernetes/_index.md#accessing-meshery-ui" >}}) guide for detailed instructions.
3. Customize your Meshery Provider Callback URL. Meshery Server supports customizing authentication flow callback URL, which can be configured in the following way:

{{< code code="MESHERY_SERVER_CALLBACK_URL=https://custom-host mesheryctl system start" >}}

## Installation: Using `mesheryctl`

Use Meshery's CLI to streamline your connection to your AKS cluster. Configure Meshery to connect to your AKS cluster by executing:

{{< code code="mesheryctl system config aks" >}}

Once configured, execute the following command to start Meshery.

{{< code code="mesheryctl system start" >}}

If you encounter any authentication issues, you can use `mesheryctl system login`. For more information, click [here]({{< ref "guides/mesheryctl/authenticate-with-meshery-via-cli/index.md" >}}) to learn more.

## Installation: Using Helm

For detailed instructions on installing Meshery using Helm V3, please refer to the [Helm Installation]({{< ref "installation/kubernetes/helm.md" >}}) guide.

## Post-Installation Steps

Optionally, you can verify the health of your Meshery deployment using <a href='{{< ref "reference/mesheryctl/system/check.md" >}}'>mesheryctl system check</a>.

You're ready to use Meshery! Open your browser and navigate to the Meshery UI.

{{< installation/accessing-meshery-ui >}}

## Accessing Meshery UI

After installing Meshery in-cluster, you can access the Meshery UI using one of several methods depending on your environment and requirements.

### Option 1: LoadBalancer Service (Recommended for Production)

When Meshery is installed in-cluster, Meshery UI is served by a Kubernetes `Service` named `meshery` in the `meshery` namespace. This `Service` is created as type `LoadBalancer` by default, forwarding port `9081` (Meshery UI) to the Meshery Server container on port `8080`.

On AKS, a `LoadBalancer` `Service` instructs Azure to provision an external load balancer and assign it a routable `EXTERNAL-IP`. Retrieve the address with:

{{< code code="kubectl get service meshery --namespace meshery" >}}

Once the `EXTERNAL-IP` column shows an address instead of `<pending>`, open Meshery UI in your browser at `http://[EXTERNAL-IP]:9081`.

If the `Service` was previously set to another type (for example, `ClusterIP`), switch it back to `LoadBalancer`:

{{< code code=`kubectl patch service meshery --namespace meshery --type merge -p '{"spec":{"type":"LoadBalancer"}}'` >}}

{{% alert title="EXTERNAL-IP not assigned?" color="warning" %}}
A `LoadBalancer` `Service` is only assigned an external address when the cluster runs a load balancer controller. Azure provides one out of the box for AKS; however, if you encounter issues, ensure your AKS cluster has the cloud provider integration enabled. Alternatively, you can use port-forwarding instead.
{{% /alert %}}

### Option 2: Port Forwarding (Universal)

Port forwarding is a universal method to access Meshery UI that works across all Kubernetes environments, including local clusters and managed clouds. This is particularly useful when:

- A LoadBalancer is not available or practical
- You're working in a local development environment
- You need temporary access without exposing Meshery externally
- You're behind a firewall or in a restricted network

For detailed instructions on using port forwarding, see the [Port Forwarding Guide]({{< ref "installation/kubernetes/port-forwarding.md" >}}).

**Quick example using mesheryctl:**

{{< code code="mesheryctl system dashboard --port-forward" >}}

**Quick example using kubectl:**

{{< code code="kubectl port-forward svc/meshery 9081:9081 --namespace meshery" >}}

Then open `http://localhost:9081` in your browser.

{{% alert title="Port Forwarding vs LoadBalancer" color="info" %}}

- **LoadBalancer**: Best for production deployments requiring external access on AKS
- **Port Forwarding**: Best for development, testing, and local access

Choose the method that best fits your use case. You can use both methods simultaneously if needed.
{{% /alert %}}

# Out-of-cluster Installation

Install Meshery on Docker (out-of-cluster) and connect it to your AKS cluster.

## Install Meshery on Docker

{{< code code="mesheryctl system start -p docker" >}}

Configure Meshery to connect to your cluster by executing:

{{< code code="mesheryctl system config aks" >}}

Once you have verified that all the services are up and running, Meshery UI will be accessible on your local machine on port 9081. Open your browser and access Meshery at [`http://localhost:9081`](http://localhost:9081).

{{< related-discussions tag="meshery" >}}
