---
layout: default
title: Deploy Azure Storage Account with Meshery
abstract: "Learn how to deploy and manage Azure Storage account through Kubernetes with Meshery, utilizing ASO operator to enhance cloud resource management"
permalink: guides/tutorials/deploy-azure-storage-account-with-meshery
model: azure
kind: StorageAccount
type: guides
category: tutorials
language: en
list: include
---

### Introduction

Meshery enables you to manage Azure Storage Accounts declaratively through Kubernetes by leveraging the Azure Service Operator (ASO). In this tutorial, you'll install the ASO operator (without CRD pattern configurations, which Meshery will manage), create a Service Principal and a Kubernetes secret with your Azure credentials, and use Meshery's UI to visually design and deploy a Storage Account resource to your Azure subscription.


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
4. [Design and Deploy an Azure Storage Account](#4-design-and-deploy-an-azure-storage-account)
5. [Verify Deployment](#5-verify-deployment)
6. [Conclusion](#6-conclusion)

### 1. Create Azure Service Principal

If you do not already have a Service Principal (SP) for Meshery, create one using the Azure CLI:


  <pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">az ad sp create-for-rbac -n azure-service-operator --role contributor --scopes /subscriptions/<AZURE_SUBSCRIPTION_ID>
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
<code class="clipboardjs">
export AZURE_CLIENT_ID=<appId>
export AZURE_CLIENT_SECRET=<password>
export AZURE_TENANT_ID=<tenant>
export AZURE_SUBSCRIPTION_ID=<subscriptionId>
</code>
</div></pre>



### 2. Connect Meshery to Your Cluster

If you haven’t already connected your cluster to Meshery, run:

<pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">
mesheryctl system start
  </code>
  </div></pre>

Then open the Meshery UI (default: [`http://localhost:9081`](http://localhost:9081)) and ensure your cluster appears under **Lifecycle → Connections**.



### 3. Install Azure Service Operator (Operator Only)

#### 3.1 Deploy ASO Operator

Apply the official ASO operator manifest (Meshery will manage CRDs):

<pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">
kubectl apply -f \
  https://github.com/Azure/azure-service-operator/releases/download/v2.13.0/azureserviceoperator_v2.13.0.yaml
  </code>
  </div></pre>
#### 3.2 Create Azure Credentials Secret

Azure Service Operator requires a Kubernetes secret with your Azure identity:

<pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">kubectl create namespace azureserviceoperator-system || true
  </code>
  </div></pre>

<pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">kubectl create secret generic azure-credentials --namespace azureserviceoperator-system --from-literal=AZURE_CLIENT_ID=$AZURE_CLIENT_ID --from-literal=AZURE_CLIENT_SECRET=$AZURE_CLIENT_SECRET --from-literal=AZURE_TENANT_ID=$AZURE_TENANT_ID --from-literal=AZURE_SUBSCRIPTION_ID=$AZURE_SUBSCRIPTION_ID
  </code>
  </div></pre>


### 4. Design and Deploy an Azure Storage Account

1. In the Meshery UI, navigate to **Kanvas**.
2. Click **Catalog**, filter by **Azure**, and select the **StorageAccount** design.
3. Click **Clone** to add it to your canvas.
4. Configure the following properties:

   * `resourceGroupName`
   * `location` (e.g., `eastus`)
   * `accountName`
   * `accessTier` (`Hot` or `Cool`)
5. Click **Actions → Deploy**.



### 5. Verify Deployment

* **Azure Portal**: Confirm the new Storage Account appears in your specified resource group.



### 6. Conclusion

You have successfully:

* Created an Azure Service Principal for Meshery
* Connected your Kubernetes cluster to Meshery
* Installed the Azure Service Operator (Meshery managed CRDs)
* Created a Kubernetes secret for Azure credentials
* Designed and deployed an Azure Storage Account using Meshery’s Kanvas

---

If you want to learn more about Azure Service Operator, visit the [official ASO documentation](https://azure.github.io/azure-service-operator/).
