---
layout: default
title: Deploy Azure resources with Meshery
abstract: "Learn how to deploy and manage Azure resources through Kubernetes with Meshery, utilizing ASO operator to enhance cloud resource management"
permalink: guides/tutorials/deploy-azure-resources-with-meshery
model: azure
kind: resources
type: guides
category: tutorials
language: en
list: include
---

### Introduction

Meshery now supports managing Azure resources declaratively through Kubernetes by integrating with Azure Service Operator (ASO). With this capability, you can visually design, deploy, and manage a variety of Azure resources—such as Storage Accounts, Key Vaults, SQL Servers, and more—directly from Meshery’s UI. In this tutorial, you’ll install the ASO operator (without CRD pattern configurations, as Meshery will handle them), create a Service Principal and a Kubernetes secret with your Azure credentials, and use Meshery to provision Azure resources seamlessly into your subscription.


### Prerequisites

Before you begin, ensure you have the following:

1. **Meshery Installed**
   A self-hosted Meshery instance running on your Kubernetes cluster (in-cluster or out-of-cluster).
2. **Kubernetes Cluster**
   A running Kubernetes cluster (v1.16+) with `kubectl` configured.
3. **Azure Subscription**
   An active Azure subscription where Storage Accounts will be provisioned.
4. **Azure CLI**
   Installed and authenticated (`az login`) in your local shell.
5. **cert-manager**
   Installed in your Kubernetes cluster (required by Azure Service Operator).


### Table of Contents

1. [Create Azure Service Principal](#1-create-azure-service-principal)
2. [Connect Meshery to Your Cluster](#2-connect-meshery-to-your-cluster)
3. [Install Azure Service Operator (Operator Only)](#3-install-azure-service-operator-operator-only)
   
   * 3.1 [Deploy ASO Operator](#31-deploy-aso-operator)
   * 3.2 [Create Azure Credentials Secret](#32-create-azure-credentials-secret)
4. [Deploy ASO Operator using Kanvas](#4-deploy-aso-using-kanvas)
5. [Start deploying azure resources in Kanvas](#5-start-deployment-of-azure-resources-using-kanvas)
6. [Conclusion](#6-conclusion)

### 1. Create Azure Service Principal

If you do not already have a Service Principal (SP) for Meshery, create one using the Azure CLI:

<pre class="codeblock-pre"><div class="codeblock">
<code class="clipboardjs">az ad sp create-for-rbac -n azure-service-operator --role contributor --scopes /subscriptions/&lt;AZURE_SUBSCRIPTION_ID&gt;
</code>
</div></pre>

This command outputs the following credentials:

* `appId`: Application ID (Client ID)
* `displayName`: Service Principal Name
* `name`: Azure Service Principal URL
* `password`: Client Secret
* `tenant`: Tenant ID

To export them, manually enter:

<pre class="codeblock-pre"><div class="codeblock">
<code class="clipboardjs">export AZURE_CLIENT_ID=&lt;appId&gt;
export AZURE_CLIENT_SECRET=&lt;password&gt;
export AZURE_TENANT_ID=&lt;tenant&gt;
export AZURE_SUBSCRIPTION_ID=&lt;subscriptionId&gt;
</code>
</div></pre>


### 2. Connect Meshery to Your Cluster

If you haven’t already connected your cluster to Meshery, run:

<pre class="codeblock-pre"><div class="codeblock">
<code class="clipboardjs">mesheryctl system start
</code>
</div></pre>


Then open the Meshery UI (default: [`http://localhost:9081`](http://localhost:9081)) and ensure your cluster appears under **Lifecycle → Connections**.



### 3. Install Azure Service Operator (Operator Only)

#### Prerequisite

Create a cert-manager that is necessary for deployment of Azure Service operator

<pre class="codeblock-pre"><div class="codeblock">
<code class="clipboardjs">kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.14.1/cert-manager.yaml
</code>
</div></pre>


#### 3.1 Deploy ASO Operator

Apply the official ASO operator manifest (Meshery will manage CRDs):

<pre class="codeblock-pre"><div class="codeblock">
<code class="clipboardjs">kubectl apply -f https://github.com/Azure/azure-service-operator/releases/download/v2.13.0/azureserviceoperator_v2.13.0.yaml
</code>
</div></pre>

#### 3.2 Create Azure Credentials Secret

Azure Service Operator requires a Kubernetes secret with your Azure identity:

<pre class="codeblock-pre"><div class="codeblock">
<code class="clipboardjs">kubectl create secret generic azure-credentials --namespace azureserviceoperator-system --from-literal=AZURE_CLIENT_ID=$AZURE_CLIENT_ID --from-literal=AZURE_CLIENT_SECRET=$AZURE_CLIENT_SECRET --from-literal=AZURE_TENANT_ID=$AZURE_TENANT_ID --from-literal=AZURE_SUBSCRIPTION_ID=$AZURE_SUBSCRIPTION_ID
</code>
</div></pre>


### 4. Deploy Azure Service Operator using Kanvas

1. In the Meshery UI, navigate to **Kanvas**.
2. Click **Catalog**, filter by **Azure**, and select the **Azure Operator** design.
3. Click **Clone** to add it to your canvas.
4. Update the secret **aso-controller-settings** in the design template. The details are also mentioned in the catalog
5. Click **Actions → Deploy**.



### 5. Start deployment of Azure resources using Kanvas

* Go to Kanvas and start by picking up Azure components and putting it to design area.
* Click **Actions → Deploy**.
* **Azure Portal**: Confirm the new Storage Account appears in your specified resource group.



### 6. Conclusion

You have successfully:

* Created an Azure Service Principal for Meshery
* Connected your Kubernetes cluster to Meshery
* Installed the Azure Service Operator (Meshery managed CRDs)
* Created a Kubernetes secret for Azure credentials
* Designed and deployed Azure resources using Meshery’s Kanvas

---

If you want to learn more about Azure Service Operator, visit the [official ASO documentation](https://azure.github.io/azure-service-operator/).
