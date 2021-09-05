---
layout: default
title: EKS
permalink: installation/platforms/eks
type: installation
display-title: "false"
language: en
list: include
image: /assets/img/platforms/eks.png
---


{% include installation_prerequisites.html %}

## To set up and run Meshery on EKS:

- Connect Meshery to your EKS cluster
  - [Meshery CLI (mesheryctl)](#connect-meshery-to-azure-kubernetes-cluster)
  - [eks CLI (eksctl)](https://eksctl.io/introduction/#installation)
  - [AWS CLI (aws)](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html)
- [Install Meshery on your EKS cluster](#install-meshery-into-your-eks-cluster)
- [Access Meshery's UI](#port-forward-to-the-meshery-ui)

### Connect Meshery to Elastic Kubernetes Cluster

The following set of instructions expects you to have created a EKS cluster in your resource group Configure Meshery to connect to your EKS cluster by executing:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 mesheryctl system config eks
 </div></div>
 </pre>

#### Manual Steps

Alternatively, you may execute the following steps to manually configure Meshery to connect to your AKS cluster.

- Install [AWS CLI(aws)](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html), and login
  to your AWS account. Install [EKS CLI (eksctl)](https://eksctl.io/introduction/#installation) too to access eks cluster easily.

- After successfull login, you have to select the subscription with which your AKS is associated with by configuring your AWS CLI to your AWS account. Refer [this link](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html) for more

- Get/Update the kubeconfig from your EKS cluster by connecting with AWS CLI
<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">
aws eks --region region update-kubeconfig --name cluster_name
</div></div>
</pre>

- Set your cluster context and check your cluster-info
<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">
kubectl set-context EKS_SERVICE_NAME
kubectl cluster-info
</div></div>
</pre>

### Install Meshery into your AKS cluster

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 $ kubectl create ns meshery
 $ helm repo add meshery https://meshery.io/charts/
 $ helm install meshery --namespace meshery
 </div></div>
 </pre>
 - Meshery server supports customizing authentication flow callback URL, which can be configured in the following way
 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 $ helm install meshery --namespace meshery --set env.MESHERY_SERVER_CALLBACK_URL=https://custom-host meshery/meshery
 </div></div>
 </pre>

### Port forward to the Meshery UI

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 export POD_NAME=$(kubectl get pods --namespace meshery -l "app.kubernetes.io/name=meshery,app.kubernetes.io/instance=meshery" -o jsonpath="{.items[0].metadata.name}")
 kubectl --namespace meshery port-forward $POD_NAME 9081:8080
 </div></div>
 </pre>

Meshery should now be running in your EKS cluster and the Meshery UI should be locally accessible. Navigate to [http://localhost:9081](http://localhost:9081) to log into Meshery.

Meshery should now be running in your EKS cluster and the Meshery UI should be accessible. Navigate to the `meshery` service endpoint to log into Meshery.
