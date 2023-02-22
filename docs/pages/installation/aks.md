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
Use Meshery to manage your AKS clusters. Deploy Meshery on AKS or outside of AKS.

{% include installation_prerequisites.html %}

#### Prequisites:
1. Access to an active AKS cluster in one of your resource groups. 
2. Azure CLI `az` installed on your machine.

## Manage AKS clusters with Meshery

- [Connect Meshery to your AKS cluster](#automatically-connect-meshery-to-azure-kubernetes-cluster)

### Automatically connect Meshery to your Azure Kubernetes Cluster

Use Meshery's CLI to automatically prepare your connection to your AKS cluster. Configure Meshery to connect to your AKS cluster by executing:

{% capture code_content %}mesheryctl system config aks{% endcapture %}
{% include code.html code=code_content %}

### [Optionally] Manually connect Meshery to your Azure Kubernetes Cluster

Alternatively, you may execute the following steps to manually configure Meshery to connect to your AKS cluster.

1. Install [Azure CLI(az)](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli), and login
  to your azure account using _az login_.

2. After successfull login, identify the subscription with which your AKS is associated.

{% capture code_content %}az account set --subscription [SUBSCRIPTION_ID]{% endcapture %}
{% include code.html code=code_content %}
<br />

3. Get the `kubeconfig` from your AKS cluster

{% capture code_content %}az aks get-credentials --resource-group [RESOURCE_GROUP] --name [AKS_SERVICE_NAME]{% endcapture %}
{% include code.html code=code_content %}
<br />

4. Set your cluster context and check your cluster-info
{% capture code_content %}kubectl config set-context AKS_SERVICE_NAME
kubectl cluster-info{% endcapture %}
{% include code.html code=code_content %}

## Install Meshery on your AKS cluster

{% capture code_content %}helm repo add meshery https://meshery.io/charts/
helm install meshery meshery/meshery --namespace meshery --create-namespace{% endcapture %}
{% include code.html code=code_content %}
<br />

Optionally, Meshery Server supports customizing callback URL for your remote provider, like so:

{% capture code_content %}helm install meshery meshery/meshery --namespace meshery --set env.MESHERY_SERVER_CALLBACK_URL=https://custom-host --create-namespace{% endcapture %}
{% include code.html code=code_content %}

### Access Meshery UI

Access Meshery UI by exposing it as a Kubernetes service or by port forwarding to Meshery UI.

#### [Optional] Port Forward to Meshery UI

{% capture code_content %}export POD_NAME=$(kubectl get pods --namespace meshery -l "app.kubernetes.io/name=meshery,app.kubernetes.io/instance=meshery" -o jsonpath="{.items[0].metadata.name}")
kubectl --namespace meshery port-forward $POD_NAME 9081:8080{% endcapture %}
{% include code.html code=code_content %}
<br />

Meshery should now be running in your AKS cluster and the Meshery UI should be accessible at the specified endpoint you've exposed to. Navigate to the `meshery` service endpoint to log into Meshery.
