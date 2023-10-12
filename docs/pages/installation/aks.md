---
layout: page
title: AKS
permalink: installation/platforms/aks
type: installation
display-title: "false"
language: en
list: include
image: /assets/img/platforms/aks.svg
---
Manage your AKS clusters with Meshery. Deploy Meshery on AKS or outside of AKS.
{% include installation_prerequisites.html %}

### General Prerequisites:
1. Access to an active AKS cluster in one of your resource groups.
2. Ensure you have the Azure CLI `az` installed on your machine.

### Install Meshery on an AKS Cluster using mesheryctl

## Connect to Azure Kubernetes Services Cluster using mesheryctl

Use Meshery's CLI to streamline your connection to your AKS cluster. Configure Meshery to connect to your AKS cluster by executing:

{% capture code_content %}$ mesheryctl system config aks{% endcapture %}
{% include code.html code=code_content %}
 <br>

Once configured, install Meshery with this command:

{% capture code_content %}$ mesheryctl system start{% endcapture %}
{% include code.html code=code_content %}
 <br>

<button class="toggle-button" onclick="HideToggleFunction()">Optional</button> Customize your Meshery Provider Callback URL

<div id="hiddendiv">
Meshery Server supports customizing authentication flow callback URL, which can be configured in the following way
<br>
{% capture code_content %}$ MESHERY_SERVER_CALLBACK_URL=https://custom-host mesheryctl system start{% endcapture %}
{% include code.html code=code_content %}
 <br>
Meshery should now be running in your AKS cluster and Meshery UI should be accessible at the `EXTERNAL IP` of `meshery` service.
<br/>
</div>
<br/>

If you encounter any authentication issues, you can use `mesheryctl system login`. For more information, click [here](/guides/mesheryctl/authenticate-with-meshery-via-cli) to learn more.

### [Optional] Connecting to an Azure Kubernetes Services Cluster using Azure CLI

Alternatively, you may choose to configure Meshery to connect to your AKS cluster manually. Follow these steps:

1. Install [Azure CLI(az)](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli) if not installed, and login to your azure account using _az login_.
2. After a successful login, identify the subscription associated with your AKS cluster:

{% capture code_content %}az account set --subscription [SUBSCRIPTION_ID]{% endcapture %}
{% include code.html code=code_content %}
<br />

3. Obtain the kubeconfig for your AKS cluster:

{% capture code_content %}az aks get-credentials --resource-group [RESOURCE_GROUP] --name [AKS_SERVICE_NAME]{% endcapture %}
{% include code.html code=code_content %}
<br />

4. Verify the current context of the cluster:
{% capture code_content %}kubectl cluster-info{% endcapture %}
{% include code.html code=code_content %}

## Install Meshery on your Azure Kubernetes Services Cluster using Helm V3

For detailed instructions on installing Meshery using Helm V3, please refer to the [Helm Installation](/installation/helm) guide.
<br />

### Access Meshery UI

To access Meshery's UI via port-forwarding, please refer to the [port-forwarding](/services/port-forward) guide for detailed instructions.