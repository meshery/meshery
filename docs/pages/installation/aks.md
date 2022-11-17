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

{% include installation_prerequisites.html %}

## To set up and run Meshery on AKS:

- Connect Meshery to your AKS cluster
  - [Meshery CLI (mesheryctl)](#connect-meshery-to-azure-kubernetes-cluster)
  - [Azure CLI (az)](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli)
- [Install Meshery on your AKS cluster](#install-meshery-into-your-aks-cluster)
- [Access Meshery's UI](#port-forward-to-the-meshery-ui)

### Connect Meshery to Azure Kubernetes Cluster

The following set of instructions expects you to have created a AKS cluster in your resource group Configure Meshery to connect to your AKS cluster by executing:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 mesheryctl system config aks
 </div></div>
 </pre>

#### Manual Steps

Alternatively, you may execute the following steps to manually configure Meshery to connect to your AKS cluster.

- Install [Azure CLI(az)](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli), and login
  to your azure account using _az login_.

- After successfull login, you have to select the subscription with which your AKS is associated with
<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">
az account set --subscription [SUBSCRIPTION_ID]
</div></div>
</pre>
- Get the kubeconfig from your AKS cluster
<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">
az aks get-credentials --resource-group [RESOURCE_GROUP] --name [AKS_SERVICE_NAME]
</div></div>
</pre>
- Set your cluster context and check your cluster-info
<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">
kubectl set-context AKS_SERVICE_NAME
kubectl cluster-info
</div></div>
</pre>

### Install Meshery into your AKS cluster

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 $ helm repo add meshery https://meshery.io/charts/
 $ helm install meshery meshery/meshery --namespace meshery --create-namespace
 </div></div>
 </pre>
 - Meshery server supports customizing authentication flow callback URL, which can be configured in the following way
 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 $ helm install meshery meshery/meshery --namespace meshery --set env.MESHERY_SERVER_CALLBACK_URL=https://custom-host --create-namespace
 </div></div>
 </pre>

### Port forward to the Meshery UI

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 export POD_NAME=$(kubectl get pods --namespace meshery -l "app.kubernetes.io/name=meshery,app.kubernetes.io/instance=meshery" -o jsonpath="{.items[0].metadata.name}")
 kubectl --namespace meshery port-forward $POD_NAME 9081:8080
 </div></div>
 </pre>

Meshery should now be running in your AKS cluster and the Meshery UI should be accessible at the specified endpoint you've exposed to. Navigate to the `meshery` service endpoint to log into Meshery.
